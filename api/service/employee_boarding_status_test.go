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

type EmployeeBoardingStatusMockRepository struct {
	mock.Mock
}

func (m *EmployeeBoardingStatusMockRepository) CreateEmployeeBoardingStatus(row models.EmployeeBoardingStatus) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err
}

func (m *EmployeeBoardingStatusMockRepository) GetEmployeeBoardingStatusAdmin(row []models.EmployeeBoardingStatus, page dto.Pagination, filter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.EmployeeBoardingStatus), args.Int(1), args.Get(2).(int64), err
}

func (m *EmployeeBoardingStatusMockRepository) GetEmployeeBoardingStatus(row []models.EmployeeBoardingStatus, page dto.Pagination, filter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.EmployeeBoardingStatus), args.Int(1), args.Get(2).(int64), err
}

func (m *EmployeeBoardingStatusMockRepository) GetEmployeeBoardingStatusById(row models.EmployeeBoardingStatus, filter dto.GetByID) (models.EmployeeBoardingStatus, int, dto.Error) {
	args := m.Called(row, filter)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.EmployeeBoardingStatus), args.Int(1), err
}

func (m *EmployeeBoardingStatusMockRepository) PatchEmployeeBoardingStatus(row models.EmployeeBoardingStatus, userID uuid.UUID) (models.EmployeeBoardingStatus, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.EmployeeBoardingStatus), args.Int(1), err
}

func (m *EmployeeBoardingStatusMockRepository) DeleteEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err

}

type EmployeeBoardingStatusServiceTestSuite struct {
	suite.Suite
	mockRepo *EmployeeBoardingStatusMockRepository
	service  EmployeeBoardingStatusService
}

func (s *EmployeeBoardingStatusServiceTestSuite) SetupTest() {
	s.mockRepo = new(EmployeeBoardingStatusMockRepository)
	s.service = NewEmployeeBoardingStatusService(s.mockRepo)
}
func EmployeeBoardingStatusTestServiceSuite(t *testing.T) {
	suite.Run(t, new(EmployeeBoardingStatusServiceTestSuite))
}

func (s *EmployeeBoardingStatusServiceTestSuite) TestCreateEmployeeBoardingStatus() {

	inputID := helper.NewUUIDv7()

	input := models.EmployeeBoardingStatus{
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
			testName:       "Success - EmployeeBoardingStatus Created",
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

			s.mockRepo.On("CreateEmployeeBoardingStatus", mock.MatchedBy(func(a models.EmployeeBoardingStatus) bool {
				return a.ID == input.ID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			status, err := s.service.CreateEmployeeBoardingStatus(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeBoardingStatusServiceTestSuite) TestGetEmployeeBoardingStatus() {

	sampleAudit := models.EmployeeBoardingStatus{ID: helper.NewUUIDv7()}
	sampleData := []models.EmployeeBoardingStatus{sampleAudit}

	var getEmployeeBoardingStatusTests = []struct {
		testName       string
		inputData      []models.EmployeeBoardingStatus
		inputPage      dto.Pagination
		inputFilter    dto.EmployeeBoardingStatusFilter
		mockReturn     []models.EmployeeBoardingStatus
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.EmployeeBoardingStatus
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get EmployeeBoardingStatus",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.EmployeeBoardingStatusFilter{},
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
			inputData:      []models.EmployeeBoardingStatus{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.EmployeeBoardingStatusFilter{},
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

	for _, tt := range getEmployeeBoardingStatusTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetEmployeeBoardingStatus", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetEmployeeBoardingStatus(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeBoardingStatusServiceTestSuite) TestGetEmployeeBoardingStatusAdmin() {

	sampleAudit := models.EmployeeBoardingStatus{ID: helper.NewUUIDv7()}
	sampleData := []models.EmployeeBoardingStatus{sampleAudit}

	var getEmployeeBoardingStatusTests = []struct {
		testName       string
		inputData      []models.EmployeeBoardingStatus
		inputPage      dto.Pagination
		inputFilter    dto.EmployeeBoardingStatusFilter
		mockReturn     []models.EmployeeBoardingStatus
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.EmployeeBoardingStatus
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Attendance Audit",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.EmployeeBoardingStatusFilter{},
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
			inputData:      []models.EmployeeBoardingStatus{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.EmployeeBoardingStatusFilter{},
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

	for _, tt := range getEmployeeBoardingStatusTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetEmployeeBoardingStatusAdmin", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetEmployeeBoardingStatusAdmin(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeBoardingStatusServiceTestSuite) TestGetEmployeeBoardingStatusById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.EmployeeBoardingStatus{ID: sampleID}
	sampleFilter := dto.GetByID{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.EmployeeBoardingStatus
		inputFilter    dto.GetByID
		mockReturn     models.EmployeeBoardingStatus
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
			mockReturn:     models.EmployeeBoardingStatus{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetEmployeeBoardingStatusById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetEmployeeBoardingStatusById(tt.inputRow, tt.inputFilter)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeBoardingStatusServiceTestSuite) TestPatchEmployeeBoardingStatus() {
	rowID := helper.NewUUIDv7()

	baseRow := models.EmployeeBoardingStatus{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName       string
		inputRow       models.EmployeeBoardingStatus
		inputEmployeeBoardingStatusID uuid.UUID
		mockReturn     models.EmployeeBoardingStatus
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
			mockReturn:     models.EmployeeBoardingStatus{},
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
				s.mockRepo.On("PatchEmployeeBoardingStatus", tt.inputRow, tt.inputEmployeeBoardingStatusID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchEmployeeBoardingStatus(tt.inputRow, tt.inputEmployeeBoardingStatusID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputEmployeeBoardingStatusID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *EmployeeBoardingStatusServiceTestSuite) TestDeleteEmployeeBoardingStatus() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.EmployeeBoardingStatus{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.EmployeeBoardingStatus
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
			s.mockRepo.On("DeleteEmployeeBoardingStatus", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			status, err := s.service.DeleteEmployeeBoardingStatus(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
