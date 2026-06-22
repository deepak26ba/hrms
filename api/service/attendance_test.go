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

type AttendanceMockRepository struct {
	mock.Mock
}

func (m *AttendanceMockRepository) CreateAttendance(row models.Attendance) (uuid.UUID, int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(2).(dto.Error); ok {
		err = val
	}
	return args.Get(0).(uuid.UUID), args.Int(1), err
}

func (m *AttendanceMockRepository) GetAttendanceAdmin(row []models.Attendance, page dto.Pagination, filter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.Attendance), args.Int(1), args.Get(2).(int64), err
}

func (m *AttendanceMockRepository) GetAttendance(row []models.Attendance, page dto.Pagination, filter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.Attendance), args.Int(1), args.Get(2).(int64), err
}

func (m *AttendanceMockRepository) GetAttendanceById(row models.Attendance, filter dto.GetByID) (models.Attendance, int, dto.Error) {
	args := m.Called(row, filter)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Attendance), args.Int(1), err
}

func (m *AttendanceMockRepository) PatchAttendance(row models.Attendance, userID uuid.UUID) (models.Attendance, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Attendance), args.Int(1), err
}

func (m *AttendanceMockRepository) DeleteAttendance(row models.Attendance, id uuid.UUID) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err

}

type AttendanceServiceTestSuite struct {
	suite.Suite
	mockRepo *AttendanceMockRepository
	service  AttendanceService
}

func (s *AttendanceServiceTestSuite) SetupTest() {
	s.mockRepo = new(AttendanceMockRepository)
	s.service = NewAttendanceService(s.mockRepo)
}
func AttendanceTestServiceSuite(t *testing.T) {
	suite.Run(t, new(AttendanceServiceTestSuite))
}

func (s *AttendanceServiceTestSuite) TestCreateAttendance() {

	inputID := helper.NewUUIDv7()

	input := models.Attendance{
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
			testName:       "Success - Attendance Created",
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

			s.mockRepo.On("CreateAttendance", mock.MatchedBy(func(a models.Attendance) bool {
				return a.ID == input.ID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			_, status, err := s.service.CreateAttendance(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceServiceTestSuite) TestGetAttendance() {

	sampleAudit := models.Attendance{ID: helper.NewUUIDv7()}
	sampleData := []models.Attendance{sampleAudit}

	var getAttendanceTests = []struct {
		testName       string
		inputData      []models.Attendance
		inputPage      dto.Pagination
		inputFilter    dto.AttendanceFilter
		mockReturn     []models.Attendance
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.Attendance
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Attendance",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.AttendanceFilter{},
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
			inputData:      []models.Attendance{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.AttendanceFilter{},
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

	for _, tt := range getAttendanceTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetAttendance", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetAttendance(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceServiceTestSuite) TestGetAttendanceAdmin() {

	sampleAudit := models.Attendance{ID: helper.NewUUIDv7()}
	sampleData := []models.Attendance{sampleAudit}

	var getAttendanceTests = []struct {
		testName       string
		inputData      []models.Attendance
		inputPage      dto.Pagination
		inputFilter    dto.AttendanceFilter
		mockReturn     []models.Attendance
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.Attendance
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Attendance Audit",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.AttendanceFilter{},
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
			inputData:      []models.Attendance{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.AttendanceFilter{},
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

	for _, tt := range getAttendanceTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetAttendanceAdmin", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetAttendanceAdmin(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceServiceTestSuite) TestGetAttendanceById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Attendance{ID: sampleID}
	sampleFilter := dto.GetByID{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Attendance
		inputFilter    dto.GetByID
		mockReturn     models.Attendance
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
			mockReturn:     models.Attendance{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetAttendanceById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetAttendanceById(tt.inputRow, tt.inputFilter)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceServiceTestSuite) TestPatchAttendance() {
	rowID := helper.NewUUIDv7()

	baseRow := models.Attendance{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName          string
		inputRow          models.Attendance
		inputAttendanceID uuid.UUID
		mockReturn        models.Attendance
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
			mockReturn:     models.Attendance{},
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
				s.mockRepo.On("PatchAttendance", tt.inputRow, tt.inputAttendanceID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchAttendance(tt.inputRow, tt.inputAttendanceID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputAttendanceID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AttendanceServiceTestSuite) TestDeleteAttendance() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Attendance{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Attendance
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
			s.mockRepo.On("DeleteAttendance", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			status, err := s.service.DeleteAttendance(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
