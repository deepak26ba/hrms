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

type EmployeeMockRepository struct {
	mock.Mock
}

func (m *EmployeeMockRepository) CreateEmployee(row models.Employee) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err
}

func (m *EmployeeMockRepository) GetEmployee(row []models.Employee, page dto.Pagination, filter dto.EmployeeFilter) ([]models.Employee, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.Employee), args.Int(1), args.Get(2).(int64), err
}

func (m *EmployeeMockRepository) GetEmployeeById(row models.Employee, filter dto.GetByID) (models.Employee, int, dto.Error) {
	args := m.Called(row, filter)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Employee), args.Int(1), err
}

func (m *EmployeeMockRepository) PatchEmployee(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error) {
	args := m.Called(row, filters)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Employee), args.Int(1), err
}

func (m *EmployeeMockRepository) DeleteEmployee(row models.Employee, id uuid.UUID) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err

}

type EmployeeServiceTestSuite struct {
	suite.Suite
	mockRepo *EmployeeMockRepository
	service  EmployeeService
}

func (s *EmployeeServiceTestSuite) SetupTest() {
	s.mockRepo = new(EmployeeMockRepository)
	s.service = NewEmployeeService(s.mockRepo)
}
func EmployeeTestServiceSuite(t *testing.T) {
	suite.Run(t, new(EmployeeServiceTestSuite))
}

func (s *EmployeeServiceTestSuite) TestCreateEmployee() {

	inputID := helper.NewUUIDv7()

	input := models.Employee{
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
			testName:       "Success - Employee Created",
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

			s.mockRepo.On("CreateEmployee", mock.MatchedBy(func(a models.Employee) bool {
				return a.ID == input.ID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			status, err := s.service.CreateEmployee(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeServiceTestSuite) TestGetEmployee() {

	sampleAudit := models.Employee{ID: helper.NewUUIDv7()}
	sampleData := []models.Employee{sampleAudit}

	var getEmployeeTests = []struct {
		testName       string
		inputData      []models.Employee
		inputPage      dto.Pagination
		inputFilter    dto.EmployeeFilter
		mockReturn     []models.Employee
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.Employee
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Employee",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.EmployeeFilter{},
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
			inputData:      []models.Employee{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.EmployeeFilter{},
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

	for _, tt := range getEmployeeTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetEmployee", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetEmployee(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeServiceTestSuite) Employee() {
	sampleID := helper.NewUUIDv7()
	sampleRow := models.Employee{ID: sampleID}
	sampleFilter := dto.GetByID{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Employee
		inputFilter    dto.GetByID
		mockReturn     models.Employee
		mockStatus     int
		mockErr        dto.Error
		expectedStatus int
		expectedMsg    string
	}{
		{
			testName:       "Success - Found Audit",
			inputRow:       sampleRow,
			inputFilter:    sampleFilter,
			mockReturn:     sampleRow,
			mockStatus:     200,
			mockErr:        dto.Error{},
			expectedStatus: 200,
			expectedMsg:    "",
		},
		{
			testName:       "Failure - Not Found",
			inputRow:       sampleRow,
			inputFilter:    sampleFilter,
			mockReturn:     models.Employee{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetEmployeeById", tt.inputRow, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetEmployeeById(tt.inputRow, tt.inputFilter)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 200 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeServiceTestSuite) TestPatchEmployee() {
	rowID := helper.NewUUIDv7()
	sampleFilter := dto.GetByID{ID: rowID}

	baseRow := models.Employee{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName        string
		inputRow        models.Employee
		inputEmployeeID uuid.UUID
		inputFilter     dto.GetByID
		mockReturn      models.Employee
		mockStatus      int
		mockErr         dto.Error
		expectedStatus  int
		expectedMsg     string
		mockRepoCall    bool
	}{
		{
			testName:       "Success - Updated",
			inputRow:       baseRow,
			inputFilter:    sampleFilter,
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
			inputFilter:    sampleFilter,
			mockReturn:     models.Employee{},
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
				s.mockRepo.On("PatchEmployee", tt.inputRow, tt.inputEmployeeID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchEmployee(tt.inputRow, tt.inputFilter)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputEmployeeID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeServiceTestSuite) TestDeleteEmployee() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Employee{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Employee
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
			s.mockRepo.On("DeleteEmployee", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			status, err := s.service.DeleteEmployee(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
