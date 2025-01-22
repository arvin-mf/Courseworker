package handler

import (
	"context"
	"courseworker/internal/dto"
	"courseworker/internal/service"
	_error "courseworker/pkg/error"
	"courseworker/pkg/response"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserHandler struct {
	serv service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	resp, err := h.serv.GetUsers()
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, usersFetchSuccess, resp)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("userId")
	resp, err := h.serv.GetUserByID(userID)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, userFetchSuccess, resp)
}

var googleOauthConfig = &oauth2.Config{
	RedirectURL: "http://localhost:8000/auth/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func (*UserHandler) LoginWithGoogle(c *gin.Context) {
	googleOauthConfig.ClientID = os.Getenv("CLIENT_ID")
	googleOauthConfig.ClientSecret = os.Getenv("CLIENT_SECRET")

	oauthState := generateStateOauthCookie()
	authURL := googleOauthConfig.AuthCodeURL(oauthState)

	fmt.Println("url = " + authURL)

	// temporary with string
	c.String(http.StatusOK, authURL)
}

func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}

func (h *UserHandler) GetGoogleDetails(c *gin.Context) {
	token, err := googleOauthConfig.Exchange(context.Background(), c.Request.FormValue("code"))
	if err != nil {
		response.HttpError(c, err)
		return
	}

	rsp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		response.HttpError(c, err)
		return
	}

	content, err := io.ReadAll(rsp.Body)
	if err != nil {
		response.HttpError(c, err)
		return
	}

	var container dto.AuthenticatedUser
	err = json.Unmarshal(content, &container)
	if err != nil {
		response.HttpError(c, err)
		return
	}

	isExist, err := h.serv.EmailExists(container.Email)
	if err != nil {
		response.HttpError(c, err)
	}
	if !isExist {
		h.serv.CreateUser(dto.CreateUserParams{
			Name:     container.Name,
			Email:    container.Email,
			HashedPw: "",
		})
	}

	resp, err := h.serv.GenerateToken(container.Email)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, 200, userLoginSuccess, resp)
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	const op _error.Op = "hand/RegisterUser"
	var req dto.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HttpBindingError(c, err, req)
		return
	}

	isExist, err := h.serv.EmailExists(req.Email)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	if isExist {
		response.HttpError(c, _error.E(
			op, _error.Forbidden,
			_error.Title("Failed to register user"),
			"email has been used",
		))
		return
	}

	if req.Password != req.ConfirmPassword {
		response.HttpError(c, _error.E(
			op, _error.InvalidRequest,
			_error.Title("Failed to register user"),
			"password confirmation does not match",
		))
		return
	}

	hashed_pw, err := h.serv.HashPassword(req.Password)
	if err != nil {
		response.HttpError(c, err)
		return
	}

	resp, err := h.serv.SendConfirmationEmail(c, dto.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		HashedPw: hashed_pw,
	})
	if err != nil {
		response.HttpError(c, err)
		return
	}

	response.Success(c, http.StatusOK, userRegisterSuccess, resp)
}

func (h *UserHandler) CreateConfirmedUser(c *gin.Context) {
	tokenTemp := c.Query("token")

	claims, err := h.serv.ValidateTokenAndClaims(c, tokenTemp)
	if err != nil {
		response.HttpError(c, err)
		return
	}

	temp_user, err := h.serv.GetTempUser(c, claims.TempUserID)
	if err != nil {
		response.HttpError(c, err)
		return
	}

	resp, err := h.serv.CreateUser(*temp_user)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, userCreateSuccess, resp)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HttpBindingError(c, err, req)
		return
	}

	resp, err := h.serv.LoginUser(req)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, userLoginSuccess, resp)
}
