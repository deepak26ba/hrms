package service

import (
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserMockRepository struct {
	mock.Mock
}

func (m *UserMockRepository) GetUser(row []models.User, page dto.Pagination, filter dto.UserFilter) ([]models.User, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.User), args.Int(1), args.Get(2).(int64), err
}

func (m *UserMockRepository) GetUserById(row models.User, id uuid.UUID) (models.User, int, dto.Error) {
	args := m.Called(row, id)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.User), args.Int(1), err
}

func (m *UserMockRepository) PatchUser(row models.User, userID uuid.UUID) (models.User, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.User), args.Int(1), err
}

func (m *UserMockRepository) DeleteUser(row models.User, id uuid.UUID) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err

}

type UserServiceTestSuite struct {
	suite.Suite
	mockRepo *UserMockRepository
	service  UserService
}

func (s *UserServiceTestSuite) SetupTest() {
	s.mockRepo = new(UserMockRepository)
	s.service = NewUserService(s.mockRepo)
}
func UserTestServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) TestGetUser() {

	sampleAudit := models.Roles{ID: helper.NewUUIDv7()}
	sampleData := []models.Roles{sampleAudit}

	var getRolesTests = []struct {
		testName       string
		inputData      []models.User
		inputPage      dto.Pagination
		inputFilter    dto.UserFilter
		mockReturn     []models.Roles
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.Roles
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Roles",
			inputData:      []models.User{},
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.UserFilter{},
			mockReturn:     sampleData,
			mockStatus:     200,
			mockCount:      1,
			mockErr:        dto.Error{},
			expectedRows:   sampleData,
			expectedStatus: 200,
			expectedCount:  1,
			expectedErr:    dto.Error{},
		},
		{
			testName:       "Failure - Database Error",
			inputData:      []models.User{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.UserFilter{},
			mockReturn:     nil,
			mockStatus:     500,
			mockCount:      0,
			mockErr:        dto.Error{Error: "Internal Server Error"},
			expectedRows:   nil,
			expectedStatus: 500,
			expectedCount:  0,
			expectedErr:    dto.Error{Error: "Internal Server Error"},
		},
	}

	for _, tt := range getRolesTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetRoles", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetUser(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *UserServiceTestSuite) TestGetUserById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.User{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.User
		mockReturn     models.User
		mockStatus     int
		mockErr        dto.Error
		expectedStatus int
		expectedMsg    string
	}{
		{
			testName:       "Success - Found Audit",
			inputRow:       sampleRow,
			mockReturn:     sampleRow,
			mockStatus:     200,
			mockErr:        dto.Error{},
			expectedStatus: 200,
			expectedMsg:    "",
		},
		{
			testName:       "Failure - Not Found",
			inputRow:       sampleRow,
			mockReturn:     models.User{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetUserById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetUserById(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *UserServiceTestSuite) TestPatchUser() {
	rowID := helper.NewUUIDv7()

	baseRow := models.User{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName       string
		inputRow       models.User
		inputUserID    uuid.UUID
		mockReturn     models.User
		mockStatus     int
		mockErr        dto.Error
		expectedStatus int
		expectedMsg    string
		mockRepoCall   bool
	}{
		{
			testName:       "Success - Updated",
			inputRow:       baseRow,
			mockReturn:     updatedRow,
			mockStatus:     200,
			mockErr:        dto.Error{},
			expectedStatus: 200,
			expectedMsg:    "",
			mockRepoCall:   true,
		},
		{
			testName:       "Failure - Update Error",
			inputRow:       baseRow,
			mockReturn:     models.User{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "update failed"},
			expectedStatus: 500,
			expectedMsg:    "update failed",
			mockRepoCall:   true,
		},
	}

	for _, tt := range patchTests {
		s.Run(tt.testName, func() {
			if tt.mockRepoCall {
				s.mockRepo.On("PatchUser", tt.inputRow, tt.inputUserID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchUser(tt.inputRow, tt.inputUserID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputUserID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *UserServiceTestSuite) TestDeleteUser() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.User{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.User
		mockStatus     int
		mockErr        dto.Error
		expectedStatus int
		expectedMsg    string
	}{
		{
			testName:       "Success - Deleted",
			inputRow:       sampleRow,
			mockStatus:     200,
			mockErr:        dto.Error{},
			expectedStatus: 200,
			expectedMsg:    "",
		},
		{
			testName:       "Failure - Not Found",
			inputRow:       sampleRow,
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("DeleteUser", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			status, err := s.service.DeleteUser(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
