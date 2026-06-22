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

type LeaveRequestMockRepository struct {
	mock.Mock
}

func (m *LeaveRequestMockRepository) CreateLeaveRequest(row models.LeaveRequest) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err
}

func (m *LeaveRequestMockRepository) GetLeaveRequestAdmin(row []models.LeaveRequest, page dto.Pagination, filter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.LeaveRequest), args.Int(1), args.Get(2).(int64), err
}

func (m *LeaveRequestMockRepository) GetLeaveRequest(row []models.LeaveRequest, page dto.Pagination, filter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.LeaveRequest), args.Int(1), args.Get(2).(int64), err
}

func (m *LeaveRequestMockRepository) GetLeaveRequestById(row models.LeaveRequest, filter dto.GetByID) (models.LeaveRequest, int, dto.Error) {
	args := m.Called(row, filter)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.LeaveRequest), args.Int(1), err
}

func (m *LeaveRequestMockRepository) PatchLeaveRequest(row models.LeaveRequest, userID uuid.UUID) (models.LeaveRequest, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.LeaveRequest), args.Int(1), err
}

type LeaveRequestServiceTestSuite struct {
	suite.Suite
	mockLeaveRequestRepo *LeaveRequestMockRepository
	mockLeaveBalanceRepo *LeaveBalanceMockRepository
	service              LeaveRequestService
}

func (s *LeaveRequestServiceTestSuite) SetupTest() {
	s.mockLeaveRequestRepo = new(LeaveRequestMockRepository)
	s.mockLeaveBalanceRepo = new(LeaveBalanceMockRepository)
	s.service = NewLeaveRequestService(s.mockLeaveRequestRepo, s.mockLeaveBalanceRepo)
}
func LeaveRequestTestServiceSuite(t *testing.T) {
	suite.Run(t, new(LeaveRequestServiceTestSuite))
}

func (s *LeaveRequestServiceTestSuite) TestCreateLeaveRequest() {

	inputID := helper.NewUUIDv7()

	input := models.LeaveRequest{
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
			testName:       "Success - LeaveRequest Created",
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

			s.mockLeaveRequestRepo.On("CreateLeaveRequest", mock.MatchedBy(func(a models.LeaveRequest) bool {
				return a.ID == input.ID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			status, err := s.service.CreateLeaveRequest(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockLeaveRequestRepo.AssertExpectations(s.T())
		})
	}
}

func (s *LeaveRequestServiceTestSuite) TestGetLeaveRequest() {

	sampleAudit := models.LeaveRequest{ID: helper.NewUUIDv7()}
	sampleData := []models.LeaveRequest{sampleAudit}

	var getLeaveRequestTests = []struct {
		testName       string
		inputData      []models.LeaveRequest
		inputPage      dto.Pagination
		inputFilter    dto.LeaveRequestFilter
		mockReturn     []models.LeaveRequest
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.LeaveRequest
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get LeaveRequest",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.LeaveRequestFilter{},
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
			inputData:      []models.LeaveRequest{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.LeaveRequestFilter{},
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

	for _, tt := range getLeaveRequestTests {
		s.Run(tt.testName, func() {

			s.mockLeaveRequestRepo.On("GetLeaveRequest", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetLeaveRequest(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockLeaveRequestRepo.AssertExpectations(s.T())
		})
	}
}

func (s *LeaveRequestServiceTestSuite) TestGetLeaveRequestAdmin() {

	sampleAudit := models.LeaveRequest{ID: helper.NewUUIDv7()}
	sampleData := []models.LeaveRequest{sampleAudit}

	var getLeaveRequestTests = []struct {
		testName       string
		inputData      []models.LeaveRequest
		inputPage      dto.Pagination
		inputFilter    dto.LeaveRequestFilter
		mockReturn     []models.LeaveRequest
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.LeaveRequest
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Attendance Audit",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.LeaveRequestFilter{},
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
			inputData:      []models.LeaveRequest{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.LeaveRequestFilter{},
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

	for _, tt := range getLeaveRequestTests {
		s.Run(tt.testName, func() {

			s.mockLeaveRequestRepo.On("GetLeaveRequestAdmin", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetLeaveRequestAdmin(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockLeaveRequestRepo.AssertExpectations(s.T())
		})
	}
}

func (s *LeaveRequestServiceTestSuite) TestGetLeaveRequestById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.LeaveRequest{ID: sampleID}
	sampleFilter := dto.GetByID{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.LeaveRequest
		inputFilter    dto.GetByID
		mockReturn     models.LeaveRequest
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
			mockReturn:     models.LeaveRequest{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockLeaveRequestRepo.On("GetLeaveRequestById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetLeaveRequestById(tt.inputRow, tt.inputFilter)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockLeaveRequestRepo.AssertExpectations(s.T())
		})
	}
}

func (s *LeaveRequestServiceTestSuite) TestPatchLeaveRequest() {
	rowID := helper.NewUUIDv7()

	baseRow := models.LeaveRequest{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName            string
		inputRow            models.LeaveRequest
		inputLeaveRequestID uuid.UUID
		mockReturn          models.LeaveRequest
		mockStatus          int
		mockErr             dto.Error
		expectedStatus      int
		expectedMsg         string
		mockRepoCall        bool
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
			mockReturn:     models.LeaveRequest{},
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
				s.mockLeaveRequestRepo.On("PatchLeaveRequest", tt.inputRow, tt.inputLeaveRequestID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchLeaveRequest(tt.inputRow, tt.inputLeaveRequestID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputLeaveRequestID, result.ID)
			}

			s.mockLeaveRequestRepo.AssertExpectations(s.T())
		})
	}
}
