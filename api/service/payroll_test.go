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

type PayrollMockRepository struct {
	mock.Mock
}

func (m *PayrollMockRepository) CreatePayroll(row models.Payroll) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err
}

func (m *PayrollMockRepository) GetPayrollAdmin(row []models.Payroll, page dto.Pagination, filter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.Payroll), args.Int(1), args.Get(2).(int64), err
}

func (m *PayrollMockRepository) GetPayroll(row []models.Payroll, page dto.Pagination, filter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.Payroll), args.Int(1), args.Get(2).(int64), err
}

func (m *PayrollMockRepository) GetPayrollById(row models.Payroll, filter dto.GetByID) (models.Payroll, int, dto.Error) {
	args := m.Called(row, filter)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Payroll), args.Int(1), err
}

func (m *PayrollMockRepository) PatchPayroll(row models.Payroll, userID uuid.UUID) (models.Payroll, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Payroll), args.Int(1), err
}

func (m *PayrollMockRepository) DeletePayroll(row models.Payroll, id uuid.UUID) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err

}

type PayrollServiceTestSuite struct {
	suite.Suite
	mockRepo *PayrollMockRepository
	service  PayrollService
}

func (s *PayrollServiceTestSuite) SetupTest() {
	s.mockRepo = new(PayrollMockRepository)
	s.service = NewPayrollService(s.mockRepo)
}
func PayrollTestServiceSuite(t *testing.T) {
	suite.Run(t, new(PayrollServiceTestSuite))
}

func (s *PayrollServiceTestSuite) TestCreatePayroll() {

	inputID := helper.NewUUIDv7()

	input := models.Payroll{
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
			testName:       "Success - Payroll Created",
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

			s.mockRepo.On("CreatePayroll", mock.MatchedBy(func(a models.Payroll) bool {
				return a.ID == input.ID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			status, err := s.service.CreatePayroll(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PayrollServiceTestSuite) TestGetPayroll() {

	sampleAudit := models.Payroll{ID: helper.NewUUIDv7()}
	sampleData := []models.Payroll{sampleAudit}

	var getPayrollTests = []struct {
		testName       string
		inputData      []models.Payroll
		inputPage      dto.Pagination
		inputFilter    dto.PayrollFilter
		mockReturn     []models.Payroll
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.Payroll
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Payroll",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.PayrollFilter{},
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
			inputData:      []models.Payroll{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.PayrollFilter{},
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

	for _, tt := range getPayrollTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetPayroll", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetPayroll(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PayrollServiceTestSuite) TestGetPayrollAdmin() {

	sampleAudit := models.Payroll{ID: helper.NewUUIDv7()}
	sampleData := []models.Payroll{sampleAudit}

	var getPayrollTests = []struct {
		testName       string
		inputData      []models.Payroll
		inputPage      dto.Pagination
		inputFilter    dto.PayrollFilter
		mockReturn     []models.Payroll
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.Payroll
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Attendance Audit",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.PayrollFilter{},
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
			inputData:      []models.Payroll{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.PayrollFilter{},
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

	for _, tt := range getPayrollTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetPayrollAdmin", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetPayrollAdmin(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PayrollServiceTestSuite) TestGetPayrollById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Payroll{ID: sampleID}
	sampleFilter := dto.GetByID{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Payroll
		inputFilter    dto.GetByID
		mockReturn     models.Payroll
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
			mockReturn:     models.Payroll{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetPayrollById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetPayrollById(tt.inputRow, tt.inputFilter)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PayrollServiceTestSuite) TestPatchPayroll() {
	rowID := helper.NewUUIDv7()

	baseRow := models.Payroll{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName       string
		inputRow       models.Payroll
		inputPayrollID uuid.UUID
		mockReturn     models.Payroll
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
			mockReturn:     models.Payroll{},
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
				s.mockRepo.On("PatchPayroll", tt.inputRow, tt.inputPayrollID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchPayroll(tt.inputRow, tt.inputPayrollID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputPayrollID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PayrollServiceTestSuite) TestDeletePayroll() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Payroll{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Payroll
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
			s.mockRepo.On("DeletePayroll", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			status, err := s.service.DeletePayroll(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
