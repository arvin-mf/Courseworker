package service

import (
	"courseworker/internal/db/sqlc"
	"courseworker/internal/dto"
	"courseworker/internal/repository"
	"courseworker/pkg/bcrypt"
	_error "courseworker/pkg/error"
	_jwt "courseworker/pkg/jwt"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

type UserService interface {
	GetUsers() ([]dto.UserResponse, error)
	GetUserByID(userID string) (*dto.UserResponse, error)
	EmailExists(email string) (bool, error)
	GenerateToken(email string) (*dto.TokenResp, error)
	CreateUser(arg dto.CreateUserParams) (*dto.ResponseID, error)
	HashPassword(pw string) (string, error)
	SendConfirmationEmail(c *gin.Context, arg dto.CreateUserParams) (*dto.RegisterUserResp, error)
	ValidateTokenAndClaims(c *gin.Context, reqToken string) (*dto.RegistrationClaims, error)
	GetTempUser(c *gin.Context, tempUserID string) (*dto.CreateUserParams, error)
	LoginUser(arg dto.LoginUserReq) (*dto.TokenResp, error)
}

type userService struct {
	repo repository.UserRepository
	rd   *redis.Client
}

func NewUserService(r repository.UserRepository, rdc *redis.Client) UserService {
	return &userService{
		repo: r,
		rd:   rdc,
	}
}

func (s *userService) GetUsers() ([]dto.UserResponse, error) {
	const op _error.Op = "serv/GetUsers"
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get users"), err)
	}
	return dto.ToUserResponses(&users), nil
}

func (s *userService) GetUserByID(userID string) (*dto.UserResponse, error) {
	const op _error.Op = "serv/GetUserByID"
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get user"), err)
	}
	return dto.ToUserResponse(user), nil
}

func (s *userService) EmailExists(email string) (bool, error) {
	const op _error.Op = "serv/EmailExists"
	count, err := s.repo.EmailExists(email)
	if err != nil {
		return true, _error.E(op, _error.Title("Failed to check email existence"))
	}
	if count < 1 {
		return false, nil
	}
	return true, nil
}

func (s *userService) GenerateToken(email string) (*dto.TokenResp, error) {
	const op _error.Op = "serv/GenerateToken"
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get user"), err)
	}

	token, err := _jwt.GenerateToken(*user)
	if err != nil {
		return nil, _error.E(op, _error.Internal, _error.Title("Failed to generate token"), err)
	}
	return dto.ToTokenResp(token), nil
}

func (s *userService) CreateUser(arg dto.CreateUserParams) (*dto.ResponseID, error) {
	const op _error.Op = "serv/CreateUser"

	input := sqlc.CreateUserParams{
		ID:       uuid.New().String(),
		Name:     arg.Name,
		Email:    arg.Email,
		Password: arg.HashedPw,
	}

	_, err := s.repo.CreateUser(input)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to create user"), err)
	}
	return &dto.ResponseID{ID: input.ID}, nil
}

func (s *userService) HashPassword(pw string) (string, error) {
	const op _error.Op = "serv/HashPassword"
	hashedPassword, err := bcrypt.HashValue(pw)
	if err != nil {
		return "", _error.E(op, _error.Internal, _error.Title("Failed to create user"), err)
	}

	return hashedPassword, err
}

func (s *userService) SendConfirmationEmail(c *gin.Context, arg dto.CreateUserParams) (*dto.RegisterUserResp, error) {
	const op _error.Op = "serv/SendConfirmationEmail"
	domain := strings.Split(arg.Email, "@")[1]
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Failed to send email"), err)
	}

	tempUserID := uuid.New().String()
	key := "temp-user:" + tempUserID
	if err := s.rd.HSet(c, key, map[string]interface{}{
		"name":      arg.Name,
		"email":     arg.Email,
		"hashed_pw": arg.HashedPw,
	}).Err(); err != nil {
		return nil, _error.E(op, _error.Cache, _error.Title("Failed to store data"), err)
	}
	if err := s.rd.Expire(c, key, 15*time.Minute).Err(); err != nil {
		log.Printf("Failed to set Redis expiration for key: %s", key)
	}

	token, err := _jwt.GenerateConfirmationToken(tempUserID)
	if err != nil {
		return nil, _error.E(op, _error.Internal, _error.Title("Failed to send email"), err)
	}

	link := fmt.Sprintf("%s/account-confirm?token=%s", os.Getenv("BASE_URL"), token)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, fromEmail))
	m.SetHeader("To", arg.Email)
	m.SetHeader("Subject", "Email Confirmation")
	m.SetBody("text/html", fmt.Sprintf("<p>Please confirm your email by clicking <a href='%s'>here</a>.</p>", link))

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return nil, _error.E(op, _error.Internal, _error.Title("Failed to send email"), err)
	}

	go func() {
		d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
		err = d.DialAndSend(m)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", arg.Email, err)
		}
	}()

	return &dto.RegisterUserResp{
		Email: arg.Email,
	}, nil
}

func (s *userService) ValidateTokenAndClaims(c *gin.Context, reqToken string) (*dto.RegistrationClaims, error) {
	const op _error.Op = "serv/ValidateTokenAndClaims"
	claims := &dto.RegistrationClaims{}
	token, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, _error.E(op, _error.InvalidRequest, _error.Title("Invalid token"), err)
	}
	return claims, nil
}

func (s *userService) GetTempUser(c *gin.Context, tempUserID string) (*dto.CreateUserParams, error) {
	const op _error.Op = "serv/GetTempUser"
	key := "temp-user:" + tempUserID
	result, err := s.rd.HGetAll(c, key).Result()
	if err != nil {
		return nil, _error.E(op, _error.Cache, _error.Title("Failed to get temp user"), err)
	}
	return &dto.CreateUserParams{
		Name:     result["name"],
		Email:    result["email"],
		HashedPw: result["hashed_pw"],
	}, nil
}

func (s *userService) LoginUser(arg dto.LoginUserReq) (*dto.TokenResp, error) {
	const op _error.Op = "serv/GetUserByEmail"
	user, err := s.repo.GetUserByEmail(arg.Email)
	if err != nil {
		return nil, _error.E(op, _error.Title("Failed to get user"), err)
	}

	if err := bcrypt.ValidateHash(arg.Password, user.Password); err != nil {
		return nil, _error.E(op, _error.Validation, _error.Title("Failed to validate password"), err)
	}

	token, err := _jwt.GenerateToken(*user)
	if err != nil {
		return nil, _error.E(op, _error.Internal, _error.Title("Failed to generate token"), err)
	}

	return dto.ToTokenResp(token), nil
}
