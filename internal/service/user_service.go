package service

import (
	"courseworker/internal/db/sqlc"
	"courseworker/internal/dto"
	"courseworker/internal/repository"
	"courseworker/pkg/bcrypt"
	_error "courseworker/pkg/error"
	"courseworker/pkg/jwt"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type UserService interface {
	GetUsers() ([]dto.UserResponse, error)
	GetUserByID(userID string) (*dto.UserResponse, error)
	EmailExists(email string) (bool, error)
	GenerateToken(email string) (*dto.TokenResp, error)
	CreateUser(arg dto.CreateUserParams) (*dto.ResponseID, error)
	HashPassword(pw string) (string, error)
	SendConfirmationEmail(arg sqlc.CreateUserParams) (*dto.RegisterUserResp, error)
	LoginUser(arg dto.LoginUserReq) (*dto.TokenResp, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{r}
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
		return false, _error.E(op, _error.Title("Failed to check email existance"))
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

	token, err := jwt.GenerateToken(*user)
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
	return dto.NewResponseID(input.ID), nil
}

func (s *userService) HashPassword(pw string) (string, error) {
	const op _error.Op = "serv/HashPassword"
	hashedPassword, err := bcrypt.HashValue(pw)
	if err != nil {
		return "", _error.E(op, _error.Internal, _error.Title("Failed to create user"), err)
	}

	return hashedPassword, err
}

func (s *userService) SendConfirmationEmail(arg sqlc.CreateUserParams) (*dto.RegisterUserResp, error) {
	const op _error.Op = "serv/SendConfirmationEmail"
	domain := strings.Split(arg.Email, "@")[1]
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return nil, _error.E(op, _error.Forbidden, _error.Title("Failed to send email"), err)
	}

	token, err := jwt.GenerateConfirmationToken(arg)
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
	m.SetBody("text/plain", fmt.Sprintf("Please confirm your email by clicking on the following link: %s", link))

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return nil, _error.E(op, _error.Internal, _error.Title("Failes to send email"), err)
	}

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	err = d.DialAndSend(m)
	if err != nil {
		return nil, _error.E(op, _error.Internal, _error.Title("Failed to send email"), err)
	}

	return &dto.RegisterUserResp{
		Email: arg.Email,
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

	token, err := jwt.GenerateToken(*user)
	if err != nil {
		return nil, _error.E(op, _error.Internal, _error.Title("Failed to generate token"), err)
	}

	return dto.ToTokenResp(token), nil
}
