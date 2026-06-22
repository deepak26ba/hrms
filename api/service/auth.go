package service

import (
	"fmt"
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofrs/uuid"
)

type AuthService interface {
	Register(row models.User) (int, dto.Error)
	GoogleRegister(payload dto.OAuthRequest) (int, dto.Error)
	Login(row dto.LoginRequest, userID string) (uuid.UUID, int, string, dto.Error)
	GoogleLogin(payload dto.OAuthRequest, userID string) (uuid.UUID, int, string, dto.Error)
}

type authservice struct {
	repo  repository.AuthRepository
	audit repository.AuditRepository
}

func NewAuthService(repo repository.AuthRepository, audit repository.AuditRepository) AuthService {

	return &authservice{
		repo:  repo,
		audit: audit,
	}

}

func (s *authservice) Login(credentials dto.LoginRequest, userID string) (uuid.UUID, int, string, dto.Error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			UserID: userID,
			Action: "Login",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	result, code, err := s.repo.Login(credentials.Email)
	if err.Error != "" {
		audit(result, err)
		return uuid.Nil, code, "", err
	}

	if helper.IsValidPassword(result.Password, credentials.Password) {
		errorResponse := dto.Error{
			Message:    "Enter Correct Password",
			StatusCode: http.StatusUnauthorized,
			Error:      "Invalid Password",
		}
		audit(credentials, errorResponse)
		return uuid.Nil, code, "", errorResponse
	}
	audit(credentials, "nil")

	return result.ID, http.StatusOK, result.Role.Name, dto.Error{}
}

func (s *authservice) Register(credentials models.User) (int, dto.Error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			UserID: "",
			Action: "Register",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	validate := validator.New()
	err := validate.Struct(credentials)
	if err != nil {
		errorResponse := dto.Error{
			Message:    "Validation failed",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(credentials, errorResponse)
		return http.StatusBadRequest, errorResponse
	}

	if helper.ValidatedPassword(credentials.Password) {
		errorResponse := dto.Error{
			Message:    "Validation failed",
			StatusCode: fiber.StatusBadRequest,
			Error:      "Enter required aleast one uppercase/lowercase letter, one number and one special character",
		}
		audit(credentials, errorResponse)
		return http.StatusBadRequest, errorResponse
	}

	code, password, errorResponse := helper.HashPassword(credentials.Password)
	if errorResponse.Error != "" {
		audit(code, errorResponse)
		return code, errorResponse
	}

	result := models.User{
		Email:    credentials.Email,
		Password: password,
		RoleID:   credentials.RoleID,
	}
	audit(result, "nil")

	return s.repo.Register(result)

}

func (s *authservice) GoogleRegister(payload dto.OAuthRequest) (int, dto.Error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			UserID: "",
			Action: "Register",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		errorResponse := dto.Error{
			Message:    "Validation failed",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(payload, errorResponse)
		return http.StatusBadRequest, errorResponse
	}

	result := models.User{
		Email:   payload.Email,
		RoleID:  payload.RoleID,
		Provider: payload.Provider,
		OAuthID: payload.UserID,
	}
	audit(result, "nil")

	return s.repo.Register(result)

}

func (s *authservice) GoogleLogin(payload dto.OAuthRequest, userID string) (uuid.UUID, int, string, dto.Error) {

	audit := func(data any, err any) {
		s.audit.CreateAudit(models.AuditTable{
			UserID: userID,
			Action: "Login",
			Level:  "Service",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	result, code, err := s.repo.Login(payload.Email)
	if err.Error != "" {
		audit(result, err)
		return uuid.Nil, code, "", err
	}

	if payload.UserID != result.OAuthID {
		errorResponse := dto.Error{
			Message:    "Unauthorized Access",
			StatusCode: http.StatusUnauthorized,
			Error:      "Invalid User",
		}
		audit(payload, errorResponse)
		return uuid.Nil, code, "", errorResponse
	}

	audit(payload, "nil")

	return result.ID, http.StatusOK, result.Role.Name, dto.Error{}
}
