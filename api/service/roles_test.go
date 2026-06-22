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

type RolesMockRepository struct {
	mock.Mock
}

func (m *RolesMockRepository) CreateRoles(row models.Roles) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err
}

func (m *RolesMockRepository) GetRoles(row []models.Roles, page dto.Pagination, filter dto.RoleFilter) ([]models.Roles, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.Roles), args.Int(1), args.Get(2).(int64), err
}

func (m *RolesMockRepository) GetRolesById(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error) {
	args := m.Called(row, id)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Roles), args.Int(1), err
}

func (m *RolesMockRepository) PatchRoles(row models.Roles, userID uuid.UUID) (models.Roles, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Roles), args.Int(1), err
}

func (m *RolesMockRepository) DeleteRoles(row models.Roles, id uuid.UUID) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err

}

type RolesServiceTestSuite struct {
	suite.Suite
	mockRepo *RolesMockRepository
	service  RolesService
}

func (s *RolesServiceTestSuite) SetupTest() {
	s.mockRepo = new(RolesMockRepository)
	s.service = NewRoleService(s.mockRepo)
}
func RolesTestServiceSuite(t *testing.T) {
	suite.Run(t, new(RolesServiceTestSuite))
}

func (s *RolesServiceTestSuite) TestCreateRoles() {

	inputID := helper.NewUUIDv7()

	input := models.Roles{
		ID: inputID,
	}

	var createTests = []struct {
		testName       string
		mockStatus     int
		mockErr        dto.Error
		expectedStatus int
		expectedMsg    string
	}{
		{
			testName:       "Success - Roles Created",
			mockStatus:     201,
			mockErr:        dto.Error{},
			expectedStatus: 201,
			expectedMsg:    "",
		},
		{
			testName:       "Failure - DB Error",
			mockStatus:     500,
			mockErr:        dto.Error{Message: "db error", StatusCode: 500},
			expectedStatus: 500,
			expectedMsg:    "db error",
		},
	}

	for _, tt := range createTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("CreateRoles", mock.MatchedBy(func(a models.Roles) bool {
				return a.ID == input.ID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			status, err := s.service.CreateRoles(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *RolesServiceTestSuite) TestGetRoles() {

	sampleAudit := models.Roles{ID: helper.NewUUIDv7()}
	sampleData := []models.Roles{sampleAudit}

	var getRolesTests = []struct {
		testName       string
		inputData      []models.Roles
		inputPage      dto.Pagination
		inputFilter    dto.RoleFilter
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
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.RoleFilter{},
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
			inputData:      []models.Roles{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.RoleFilter{},
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

			rows, status, count, err := s.service.GetRoles(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *RolesServiceTestSuite) TestGetRolesById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Roles{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Roles
		mockReturn     models.Roles
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
			mockReturn:     models.Roles{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetRolesById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetRolesById(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *RolesServiceTestSuite) TestPatchRoles() {
	rowID := helper.NewUUIDv7()

	baseRow := models.Roles{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName       string
		inputRow       models.Roles
		inputRolesID   uuid.UUID
		mockReturn     models.Roles
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
			mockReturn:     models.Roles{},
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
				s.mockRepo.On("PatchRoles", tt.inputRow, tt.inputRolesID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchRoles(tt.inputRow, tt.inputRolesID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputRolesID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *RolesServiceTestSuite) TestDeleteRoles() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Roles{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Roles
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
			s.mockRepo.On("DeleteRoles", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			status, err := s.service.DeleteRoles(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
