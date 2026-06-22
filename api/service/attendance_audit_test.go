package service

import (
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AttendanceAuditMockRepository struct {
	mock.Mock
}

func (m *AttendanceAuditMockRepository) CreateAttendanceAudit(row models.AttendanceAudit) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err
}

func (m *AttendanceAuditMockRepository) GetAttendanceAuditAdmin(row []models.AttendanceAudit, page dto.Pagination, filter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.AttendanceAudit), args.Int(1), args.Get(2).(int64), err
}

func (m *AttendanceAuditMockRepository) GetAttendanceAudit(row []models.AttendanceAudit, page dto.Pagination, filter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.AttendanceAudit), args.Int(1), args.Get(2).(int64), err
}

func (m *AttendanceAuditMockRepository) GetAttendanceAuditById(row models.AttendanceAudit, filter dto.GetByID) (models.AttendanceAudit, int, dto.Error) {
	args := m.Called(row, filter)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.AttendanceAudit), args.Int(1), err
}

func (m *AttendanceAuditMockRepository) PatchAttendanceAudit(row models.AttendanceAudit, userID uuid.UUID) (models.AttendanceAudit, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.AttendanceAudit), args.Int(1), err
}

type AttendanceAuditServiceTestSuite struct {
	suite.Suite
	mockRepo *AttendanceAuditMockRepository
	service  AttendanceAuditService
}

func (s *AttendanceAuditServiceTestSuite) SetupTest() {
	s.mockRepo = new(AttendanceAuditMockRepository)
	s.service = NewAttendanceAuditService(s.mockRepo)
}

func AttendanceAuditTestServiceSuite(t *testing.T) {
	suite.Run(t, new(AttendanceAuditServiceTestSuite))
}

func (s *AttendanceAuditServiceTestSuite) TestCreateAttendanceAudit() {

	inputID := helper.NewUUIDv7()
	attendanceID := helper.NewUUIDv7()
	userID := helper.NewUUIDv7()

	input := models.AttendanceAudit{
		ID:           inputID,
		AttendanceID: attendanceID,
		UserId:       userID,
		TimeIn:       time.Now(),
		TimeOut:      time.Now(),
	}

	var createTests = []struct {
		testName       string
		mockStatus     int
		mockErr        dto.Error
		expectedStatus int
		expectedMsg    string
	}{
		{
			testName:       "Success - Audit Created",
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

			s.mockRepo.On("CreateAttendanceAudit", mock.MatchedBy(func(a models.AttendanceAudit) bool {
				return a.ID == input.ID && a.AttendanceID == input.AttendanceID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			status, err := s.service.CreateAttendanceAudit(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceAuditServiceTestSuite) TestGetAttendanceAuditAdmin() {

	sampleAudit := models.AttendanceAudit{ID: helper.NewUUIDv7()}
	sampleData := []models.AttendanceAudit{sampleAudit}

	var getAttendanceAuditTests = []struct {
		testName       string
		inputData      []models.AttendanceAudit
		inputPage      dto.Pagination
		inputFilter    dto.AttendanceAuditFilter
		mockReturn     []models.AttendanceAudit
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.AttendanceAudit
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Attendance Audit",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.AttendanceAuditFilter{},
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
			inputData:      []models.AttendanceAudit{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.AttendanceAuditFilter{},
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

	for _, tt := range getAttendanceAuditTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetAttendanceAuditAdmin", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetAttendanceAuditAdmin(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceAuditServiceTestSuite) TestGetAttendanceAudit() {

	sampleAudit := models.AttendanceAudit{ID: helper.NewUUIDv7()}
	sampleData := []models.AttendanceAudit{sampleAudit}

	var getAttendanceAuditTests = []struct {
		testName       string
		inputData      []models.AttendanceAudit
		inputPage      dto.Pagination
		inputFilter    dto.AttendanceAuditFilter
		mockReturn     []models.AttendanceAudit
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.AttendanceAudit
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Attendance Audit",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.AttendanceAuditFilter{},
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
			inputData:      []models.AttendanceAudit{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.AttendanceAuditFilter{},
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

	for _, tt := range getAttendanceAuditTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetAttendanceAudit", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetAttendanceAudit(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceAuditServiceTestSuite) TestGetAttendanceAuditById() {
	sampleID := helper.NewUUIDv7()
	sampleRow := models.AttendanceAudit{ID: sampleID}
	sampleFilter := dto.GetByID{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.AttendanceAudit
		inputFilter    dto.GetByID
		mockReturn     models.AttendanceAudit
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
			mockReturn:     models.AttendanceAudit{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetAttendanceAuditById", tt.inputRow, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetAttendanceAuditById(tt.inputRow, tt.inputFilter)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 200 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceAuditServiceTestSuite) TestPatchAttendanceAudit() {
	rowID := helper.NewUUIDv7()
	userID := helper.NewUUIDv7()

	baseRow := models.AttendanceAudit{ID: rowID}
	updatedRow := baseRow
	updatedRow.UserId = userID

	var patchTests = []struct {
		testName       string
		inputRow       models.AttendanceAudit
		inputUserID    uuid.UUID
		mockReturn     models.AttendanceAudit
		mockStatus     int
		mockErr        dto.Error
		expectedStatus int
		expectedMsg    string
		mockRepoCall   bool
	}{
		{
			testName:       "Success - Partial Update",
			inputRow:       baseRow,
			inputUserID:    userID,
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
			inputUserID:    userID,
			mockReturn:     models.AttendanceAudit{},
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
				s.mockRepo.On("PatchAttendanceAudit", tt.inputRow, tt.inputUserID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchAttendanceAudit(tt.inputRow, tt.inputUserID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputUserID, result.UserId)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
