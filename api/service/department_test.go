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

type DepartmentMockRepository struct {
	mock.Mock
}

func (m *DepartmentMockRepository) CreateDepartment(row models.Department) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err
}

func (m *DepartmentMockRepository) GetDepartment(row []models.Department, page dto.Pagination, filter dto.DepartmentFilter) ([]models.Department, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.Department), args.Int(1), args.Get(2).(int64), err
}

func (m *DepartmentMockRepository) GetDepartmentById(row models.Department, id uuid.UUID) (models.Department, int, dto.Error) {
	args := m.Called(row, id)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Department), args.Int(1), err
}

func (m *DepartmentMockRepository) PatchDepartment(row models.Department, userID uuid.UUID) (models.Department, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Department), args.Int(1), err
}

func (m *DepartmentMockRepository) DeleteDepartment(row models.Department, id uuid.UUID) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err

}

type DepartmentServiceTestSuite struct {
	suite.Suite
	mockRepo *DepartmentMockRepository
	service  DepartmentService
}

func (s *DepartmentServiceTestSuite) SetupTest() {
	s.mockRepo = new(DepartmentMockRepository)
	s.service = NewDepartmentService(s.mockRepo)
}
func DepartmentTestServiceSuite(t *testing.T) {
	suite.Run(t, new(DepartmentServiceTestSuite))
}

func (s *DepartmentServiceTestSuite) TestCreateDepartment() {

	inputID := helper.NewUUIDv7()

	input := models.Department{
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
			testName:       "Success - Department Created",
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

			s.mockRepo.On("CreateDepartment", mock.MatchedBy(func(a models.Department) bool {
				return a.ID == input.ID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			status, err := s.service.CreateDepartment(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *DepartmentServiceTestSuite) TestGetDepartment() {

	sampleAudit := models.Department{ID: helper.NewUUIDv7()}
	sampleData := []models.Department{sampleAudit}

	var getDepartmentTests = []struct {
		testName       string
		inputData      []models.Department
		inputPage      dto.Pagination
		inputFilter    dto.DepartmentFilter
		mockReturn     []models.Department
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.Department
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Department",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.DepartmentFilter{},
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
			inputData:      []models.Department{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.DepartmentFilter{},
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

	for _, tt := range getDepartmentTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetDepartment", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetDepartment(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *DepartmentServiceTestSuite) TestGetDepartmentById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Department{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Department
		mockReturn     models.Department
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
			mockReturn:     models.Department{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetDepartmentById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetDepartmentById(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *DepartmentServiceTestSuite) TestPatchDepartment() {
	rowID := helper.NewUUIDv7()

	baseRow := models.Department{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName          string
		inputRow          models.Department
		inputDepartmentID uuid.UUID
		mockReturn        models.Department
		mockStatus        int
		mockErr           dto.Error
		expectedStatus    int
		expectedMsg       string
		mockRepoCall      bool
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
			mockReturn:     models.Department{},
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
				s.mockRepo.On("PatchDepartment", tt.inputRow, tt.inputDepartmentID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchDepartment(tt.inputRow, tt.inputDepartmentID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputDepartmentID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *DepartmentServiceTestSuite) TestDeleteDepartment() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Department{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Department
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
			s.mockRepo.On("DeleteDepartment", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			status, err := s.service.DeleteDepartment(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
