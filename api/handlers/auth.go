package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"hrms/api/middleware"
	"hrms/api/service"
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/internals/config"
	"hrms/pkg/models"
	"io"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	service service.AuthService
	audit   service.AuditService
}

var (
	googleOAuthRegister = oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/v1/googleoauth/register/callback",
		ClientID:     config.GetClientID(),
		ClientSecret: config.GetClientSecret(),
		Scopes:       []string{"email", "profile", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	googleOAuthLogin = oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/v1/googleoauth/login/callback",
		ClientID:     config.GetClientID(),
		ClientSecret: config.GetClientSecret(),
		Scopes:       []string{"email", "profile", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
)

func NewAuthHandler(service service.AuthService, audit service.AuditService) *AuthHandler {

	return &AuthHandler{
		service: service,
		audit:   audit,
	}

}

func (h *AuthHandler) Register(f fiber.Ctx) error {

	var payload models.User

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: "",
			Action: "Register",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if err := f.Bind().Body(&payload); err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed converting json to struct",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(payload, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)

	}

	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Validation failed",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(payload, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if helper.ValidatedPassword(payload.Password) {
		errorResponse := &dto.Error{
			Message:    "Validation failed",
			StatusCode: fiber.StatusBadRequest,
			Error:      "Enter required aleast one uppercase/lowercase letter, one number and one special character",
		}
		audit(payload, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	code, errorResponse := h.service.Register(payload)
	if errorResponse.Error != "" {
		audit(code, errorResponse)
		return f.Status(code).JSON(errorResponse)

	}

	successResponse := &dto.Success{
		Message:    "Successfully Created",
		StatusCode: fiber.StatusCreated,
	}
	audit(successResponse, "nil")

	return f.Status(fiber.StatusCreated).JSON(successResponse)

}

func (h *AuthHandler) Login(f fiber.Ctx) error {

	var loginCredentials dto.LoginRequest
	currentUserID, _ := f.Locals("user_id").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "Login",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if err := f.Bind().Body(&loginCredentials); err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed converting json to struct",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(loginCredentials, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)

	}

	validate := validator.New()
	err := validate.Struct(loginCredentials)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Validation failed",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(loginCredentials, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	id, code, role, errorResponse := h.service.Login(loginCredentials, currentUserID)
	if errorResponse.Error != "" {
		audit(id, errorResponse)
		return f.Status(code).JSON(errorResponse)
	}

	token, err := middleware.GenerateJWT(role, id)
	if err != nil {
		audit(token, errorResponse)
		return f.Status(code).JSON(err.Error())
	}

	successResponse := &dto.Login{
		Message:    "Successfully Logged in",
		StatusCode: fiber.StatusOK,
		Data:       token,
	}
	audit(successResponse, "nil")

	return f.Status(fiber.StatusOK).JSON(successResponse)
}

func (h *AuthHandler) GoogleLogin(f fiber.Ctx) error {
	url := googleOAuthLogin.AuthCodeURL("randomstate")

	return f.Redirect().To(url)
}

func (h *AuthHandler) GoogleRegister(f fiber.Ctx) error {
	url := googleOAuthRegister.AuthCodeURL("randomstate")

	return f.Redirect().To(url)
}

func (h *AuthHandler) GoogleRegisterCallback(f fiber.Ctx) error {

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: "",
			Action: "Google Register CallBack",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	code := f.Query("code")

	token, err := googleOAuthRegister.Exchange(context.Background(), code)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Token exchange failed",
			StatusCode: fiber.StatusInternalServerError,
			Error:      err.Error(),
		}
		audit(token, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	client := googleOAuthRegister.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed to get user info",
			StatusCode: fiber.StatusInternalServerError,
			Error:      err.Error(),
		}
		audit(resp, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	if resp.StatusCode != 200 {
		errorResponse := &dto.Error{
			Message:    "Failed to get user info",
			StatusCode: fiber.StatusInternalServerError,
			Error:      fmt.Sprintf("Google API returned status code %d", resp.StatusCode),
		}
		audit(resp, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed to read user info",
			StatusCode: fiber.StatusInternalServerError,
			Error:      err.Error(),
		}
		audit(body, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	var userInfo map[string]interface{}

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed to unmarshal user info",
			StatusCode: fiber.StatusInternalServerError,
			Error:      err.Error(),
		}
		audit(body, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}
	UserRoleID, errorResponse := helper.StringToUUID("019d0076-f4bf-79e9-aac4-e6dafd4d6670")
	if errorResponse.Error != "" {
		audit(UserRoleID, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	RegisterRequest := dto.OAuthRequest{
		Email:    userInfo["email"].(string),
		Name:     userInfo["name"].(string),
		Provider: "google",
		UserID:   userInfo["id"].(string),
		RoleID:   UserRoleID,
	}

	status, errorResponse := h.service.GoogleRegister(RegisterRequest)
	if errorResponse.Error != "" {
		audit(RegisterRequest, errorResponse)
		return f.Status(status).JSON(errorResponse)
	}

	successResponse := &dto.Success{
		Message:    "Successfully Registered",
		StatusCode: fiber.StatusCreated,
	}
	audit(successResponse, "nil")
	return f.Status(fiber.StatusCreated).JSON(successResponse)
}

func (h *AuthHandler) GoogleLoginCallback(f fiber.Ctx) error {

	currentUserID, _ := f.Locals("user_id").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "Google Login CallBack",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	code := f.Query("code")

	token, err := googleOAuthLogin.Exchange(context.Background(), code)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Token exchange failed",
			StatusCode: fiber.StatusInternalServerError,
			Error:      err.Error(),
		}
		audit(token, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	client := googleOAuthLogin.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed to get user info",
			StatusCode: fiber.StatusInternalServerError,
			Error:      err.Error(),
		}
		audit(resp, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	if resp.StatusCode != 200 {
		errorResponse := &dto.Error{
			Message:    "Failed to get user info",
			StatusCode: fiber.StatusInternalServerError,
			Error:      fmt.Sprintf("Google API returned status code %d", resp.StatusCode),
		}
		audit(resp, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed to read user info",
			StatusCode: fiber.StatusInternalServerError,
			Error:      err.Error(),
		}
		audit(body, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	var userInfo map[string]interface{}

	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed to unmarshal user info",
			StatusCode: fiber.StatusInternalServerError,
			Error:      err.Error(),
		}
		audit(body, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	LoginRequest := dto.OAuthRequest{
		Email:    userInfo["email"].(string),
		Name:     userInfo["name"].(string),
		Provider: "google",
		UserID:   userInfo["id"].(string),
	}

	id, statusCode, role, errorResponse := h.service.GoogleLogin(LoginRequest, currentUserID)
	if errorResponse.Error != "" {
		audit(id, errorResponse)
		return f.Status(statusCode).JSON(errorResponse)
	}

	JWTtoken, err := middleware.GenerateJWT(role, id)
	if err != nil {
		audit(JWTtoken, errorResponse)
		return f.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	successResponse := &dto.Login{
		Message:    "Successfully Logged in",
		StatusCode: fiber.StatusOK,
		Data:       JWTtoken,
	}
	audit(successResponse, "nil")

	return f.Status(fiber.StatusOK).JSON(successResponse)
}
