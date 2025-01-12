package handler

import (
	"context"
	"courseworker/internal/db/sqlc"
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
	response.Success(c, http.StatusOK, "Users retrieved successfully", resp)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("userId")
	resp, err := h.serv.GetUserByID(userID)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "User retrieved successfully", resp)
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

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		response.HttpError(c, err)
		return
	}

	content, err := io.ReadAll(resp.Body)
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
		h.serv.CreateUser(sqlc.CreateUserParams{
			Name:     container.Name,
			Email:    container.Email,
			Password: "",
		})
	}

	data, err := h.serv.GenerateToken(container.Email)
	if err != nil {
		response.HttpError(c, err)
		return
	}
	response.Success(c, 200, "success", gin.H{"token": data})
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
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
			_error.Op("hand/RegisterUser"),
			_error.Forbidden,
			_error.Title("Failed to register user"),
			_error.Detail("email has been used"),
		))
		return
	}

	if req.Password != req.ConfirmPassword {
		response.HttpError(c, _error.E(
			_error.Op("hand/RegisterUser"),
			_error.InvalidRequest,
			_error.Title("Failed to register user"),
			_error.Detail("password confirmation does not match"),
		))
		return
	}

	hashed_pw, err := h.serv.HashPassword(req.Password)
	if err != nil {
		response.HttpError(c, err)
		return
	}

	resp, err := h.serv.SendConfirmationEmail(sqlc.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed_pw,
	})
	if err != nil {
		response.HttpError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Confirmation email sent", resp)
}
