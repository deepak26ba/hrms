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

type LeaveBalanceMockRepository struct {
	mock.Mock
}

func (m *LeaveBalanceMockRepository) GetLeaveBalanceAdmin(row []models.LeaveBalance, page dto.Pagination, filter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.LeaveBalance), args.Int(1), args.Get(2).(int64), err
}

func (m *LeaveBalanceMockRepository) GetLeaveBalance(row []models.LeaveBalance, page dto.Pagination, filter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.LeaveBalance), args.Int(1), args.Get(2).(int64), err
}

func (m *LeaveBalanceMockRepository) GetLeaveBalanceById(row models.LeaveBalance, filter dto.GetByID) (models.LeaveBalance, int, dto.Error) {
	args := m.Called(row, filter)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.LeaveBalance), args.Int(1), err
}

func (m *LeaveBalanceMockRepository) PatchLeaveBalance(row models.LeaveBalance, id uuid.UUID) (models.LeaveBalance, int, dto.Error) {
	args := m.Called(row, id)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.LeaveBalance), args.Int(1), err
}

type LeaveBalanceServiceTestSuite struct {
	suite.Suite
	mockLeaveBalanceRepo *LeaveBalanceMockRepository
	service  LeaveBalanceService
}

func (s *LeaveBalanceServiceTestSuite) SetupTest() {
	s.mockLeaveBalanceRepo = new(LeaveBalanceMockRepository)
	s.service = NewLeaveBalanceService(s.mockLeaveBalanceRepo)
}
func LeaveBalanceTestServiceSuite(t *testing.T) {
	suite.Run(t, new(LeaveBalanceServiceTestSuite))
}

func (s *LeaveBalanceServiceTestSuite) TestGetLeaveBalance() {

	sampleAudit := models.LeaveBalance{ID: helper.NewUUIDv7()}
	sampleData := []models.LeaveBalance{sampleAudit}

	var getLeaveBalanceTests = []struct {
		testName       string
		inputData      []models.LeaveBalance
		inputPage      dto.Pagination
		inputFilter    dto.LeaveBalanceFilter
		mockReturn     []models.LeaveBalance
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.LeaveBalance
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get LeaveBalance",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.LeaveBalanceFilter{},
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
			inputData:      []models.LeaveBalance{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.LeaveBalanceFilter{},
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

	for _, tt := range getLeaveBalanceTests {
		s.Run(tt.testName, func() {

			s.mockLeaveBalanceRepo.On("GetLeaveBalance", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetLeaveBalance(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockLeaveBalanceRepo.AssertExpectations(s.T())
		})
	}
}

func (s *LeaveBalanceServiceTestSuite) TestGetLeaveBalanceAdmin() {

	sampleAudit := models.LeaveBalance{ID: helper.NewUUIDv7()}
	sampleData := []models.LeaveBalance{sampleAudit}

	var getLeaveBalanceTests = []struct {
		testName       string
		inputData      []models.LeaveBalance
		inputPage      dto.Pagination
		inputFilter    dto.LeaveBalanceFilter
		mockReturn     []models.LeaveBalance
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.LeaveBalance
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Attendance Audit",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.LeaveBalanceFilter{},
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
			inputData:      []models.LeaveBalance{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.LeaveBalanceFilter{},
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

	for _, tt := range getLeaveBalanceTests {
		s.Run(tt.testName, func() {

			s.mockLeaveBalanceRepo.On("GetLeaveBalanceAdmin", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetLeaveBalanceAdmin(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockLeaveBalanceRepo.AssertExpectations(s.T())
		})
	}
}

func (s *LeaveBalanceServiceTestSuite) TestGetLeaveBalanceById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.LeaveBalance{ID: sampleID}
	sampleFilter := dto.GetByID{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.LeaveBalance
		inputFilter    dto.GetByID
		mockReturn     models.LeaveBalance
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
			mockReturn:     models.LeaveBalance{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockLeaveBalanceRepo.On("GetLeaveBalanceById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetLeaveBalanceById(tt.inputRow, tt.inputFilter)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockLeaveBalanceRepo.AssertExpectations(s.T())
		})
	}
}
