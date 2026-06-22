package service

import (
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) Login(email string) (models.User, int, dto.Error) {
	args := m.Called(email)

	var err dto.Error
	var user models.User
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}
	if args.Get(0) != nil {
		user = args.Get(0).(models.User)
	}

	return user, args.Int(1), err
}

func (m *MockAuthRepo) Register(user models.User) (int, dto.Error) {
	args := m.Called(user)

	var err dto.Error
	if args.Get(1) != nil {
		err = args.Get(1).(dto.Error)
	}

	return args.Int(0), err
}

type MockAuditRepo struct {
	mock.Mock
}

func (m *MockAuditRepo) CreateAudit(a models.AuditTable) {
	m.Called(a)
}

type AuthServiceTestSuite struct {
	suite.Suite
	authRepo  *MockAuthRepo
	auditRepo *MockAuditRepo
	service   authservice
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.authRepo = new(MockAuthRepo)
	s.auditRepo = new(MockAuditRepo)
	s.service = authservice{repo: s.authRepo, audit: s.auditRepo}
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (s *AuthServiceTestSuite) TestLogin() {
	userID := "ADMIN"
	userUUID := helper.NewUUIDv7()

	hashedPassword := "$2a$10$OSHVWhEuHjfQczqbbkkJ2e3dEKp/4KbXp7IS/XCtO0N2zKR/XYJ2e"
	correctPassword := "P@ssw0rd"

	tests := []struct {
		name        string
		input       dto.LoginRequest
		repoReturn  models.User
		mockStatus  int
		mockErr     dto.Error
		expectCode  int
		expectedMsg string
		expectErr   bool
	}{
		{
			name:  "Success",
			input: dto.LoginRequest{Email: "test@mail.com", Password: correctPassword},
			repoReturn: models.User{
				ID:       userUUID,
				Password: hashedPassword,
				Role:     models.Roles{Name: "ADMIN"},
			},
			mockStatus:  200,
			mockErr:     dto.Error{},
			expectCode:  200,
			expectErr:   false,
			expectedMsg: "",
		},
		{
			name:  "Invalid Password",
			input: dto.LoginRequest{Email: "test@mail.com", Password: "wrong_password"},
			repoReturn: models.User{
				ID:       userUUID,
				Password: hashedPassword,
				Role:     models.Roles{Name: "ADMIN"},
			},
			mockStatus:  401,
			mockErr:     dto.Error{Message: "Enter Correct Password"},
			expectCode:  401,
			expectErr:   true,
			expectedMsg: "Enter Correct Password",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.authRepo.On("Login", tt.input.Email).
				Return(tt.repoReturn, tt.mockStatus, tt.mockErr).
				Once()

			s.auditRepo.On("CreateAudit", mock.Anything).Return().Maybe()

			_, code, _, err := s.service.Login(tt.input, userID)

			s.Equal(tt.expectCode, code)
			if tt.expectErr {
				s.Equal(tt.expectedMsg, err.Message)

			} else {
				s.Equal(tt.expectedMsg, err.Message)

			}
		})
	}
}

func (s *AuthServiceTestSuite) TestRegister() {

	roleUUID := helper.NewUUIDv7()

	tests := []struct {
		name        string
		input       models.User
		mockStatus  int
		mockErr     dto.Error
		expectCode  int
		expectedMsg string
		expectErr   bool
	}{
		{
			name: "Success",
			input: models.User{
				Email:    "test@mail.com",
				Password: "P@ssw0rd",
				RoleID:   roleUUID,
			},
			mockStatus:  200,
			mockErr:     dto.Error{},
			expectCode:  200,
			expectErr:   false,
			expectedMsg: "",
		},
		{
			name: "Invalid Password",
			input: models.User{
				Email:    "test@mail.com",
				Password: "wrong_password",
			},
			mockStatus:  400,
			mockErr:     dto.Error{Message: "Enter Correct Password"},
			expectCode:  400,
			expectErr:   true,
			expectedMsg: "Validation failed",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.authRepo.On("Register", mock.Anything).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			s.auditRepo.On("CreateAudit", mock.Anything).Return().Maybe()

			code, err := s.service.Register(tt.input)

			s.Equal(tt.expectCode, code)
			if tt.expectErr {
				s.Equal(tt.expectedMsg, err.Message)

			} else {
				s.Equal(tt.expectedMsg, err.Message)

			}
		})
	}
}
