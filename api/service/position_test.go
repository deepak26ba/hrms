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

type PositionMockRepository struct {
	mock.Mock
}

func (m *PositionMockRepository) CreatePosition(row models.Position) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err
}

func (m *PositionMockRepository) GetPosition(row []models.Position, page dto.Pagination, filter dto.PositionFilter) ([]models.Position, int, int64, dto.Error) {
	args := m.Called(row, page, filter)

	var err dto.Error
	if args.Get(3) != nil {
		err = args.Get(3).(dto.Error)
	}

	return args.Get(0).([]models.Position), args.Int(1), args.Get(2).(int64), err
}

func (m *PositionMockRepository) GetPositionById(row models.Position, id uuid.UUID) (models.Position, int, dto.Error) {
	args := m.Called(row, id)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Position), args.Int(1), err
}

func (m *PositionMockRepository) PatchPosition(row models.Position, userID uuid.UUID) (models.Position, int, dto.Error) {
	args := m.Called(row, userID)

	var err dto.Error
	if args.Get(2) != nil {
		err = args.Get(2).(dto.Error)
	}

	return args.Get(0).(models.Position), args.Int(1), err
}

func (m *PositionMockRepository) DeletePosition(row models.Position, id uuid.UUID) (int, dto.Error) {
	args := m.Called(row)

	var err dto.Error
	if val, ok := args.Get(1).(dto.Error); ok {
		err = val
	}
	return args.Int(0), err

}

type PositionServiceTestSuite struct {
	suite.Suite
	mockRepo *PositionMockRepository
	service  PositionService
}

func (s *PositionServiceTestSuite) SetupTest() {
	s.mockRepo = new(PositionMockRepository)
	s.service = NewPositionService(s.mockRepo)
}
func PositionTestServiceSuite(t *testing.T) {
	suite.Run(t, new(PositionServiceTestSuite))
}

func (s *PositionServiceTestSuite) TestCreatePosition() {

	inputID := helper.NewUUIDv7()

	input := models.Position{
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
			testName:       "Success - Position Created",
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

			s.mockRepo.On("CreatePosition", mock.MatchedBy(func(a models.Position) bool {
				return a.ID == input.ID
			})).Return(tt.mockStatus, tt.mockErr).Once()

			status, err := s.service.CreatePosition(input)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PositionServiceTestSuite) TestGetPosition() {

	sampleAudit := models.Position{ID: helper.NewUUIDv7()}
	sampleData := []models.Position{sampleAudit}

	var getPositionTests = []struct {
		testName       string
		inputData      []models.Position
		inputPage      dto.Pagination
		inputFilter    dto.PositionFilter
		mockReturn     []models.Position
		mockStatus     int
		mockCount      int64
		mockErr        dto.Error
		expectedRows   []models.Position
		expectedStatus int
		expectedCount  int64
		expectedErr    dto.Error
	}{
		{
			testName:       "Success - Get Position",
			inputData:      sampleData,
			inputPage:      dto.Pagination{Page: 1, Limit: 10},
			inputFilter:    dto.PositionFilter{},
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
			inputData:      []models.Position{},
			inputPage:      dto.Pagination{},
			inputFilter:    dto.PositionFilter{},
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

	for _, tt := range getPositionTests {
		s.Run(tt.testName, func() {

			s.mockRepo.On("GetPosition", tt.inputData, tt.inputPage, tt.inputFilter).
				Return(tt.mockReturn, tt.mockStatus, tt.mockCount, tt.mockErr).
				Once()

			rows, status, count, err := s.service.GetPosition(tt.inputData, tt.inputPage, tt.inputFilter)

			s.Equal(tt.expectedRows, rows)
			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedCount, count)
			s.Equal(tt.expectedErr, err)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PositionServiceTestSuite) TestGetPositionById() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Position{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Position
		mockReturn     models.Position
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
			mockReturn:     models.Position{},
			mockStatus:     500,
			mockErr:        dto.Error{Message: "not found"},
			expectedStatus: 500,
			expectedMsg:    "not found",
		},
	}

	for _, tt := range getByIdTests {
		s.Run(tt.testName, func() {
			s.mockRepo.On("GetPositionById", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
				Once()

			result, status, err := s.service.GetPositionById(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)
			if tt.expectedStatus == 1 {
				s.Equal(tt.mockReturn.ID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PositionServiceTestSuite) TestPatchPosition() {
	rowID := helper.NewUUIDv7()

	baseRow := models.Position{ID: rowID}
	updatedRow := baseRow

	var patchTests = []struct {
		testName          string
		inputRow          models.Position
		inputPositionID uuid.UUID
		mockReturn        models.Position
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
			mockReturn:     models.Position{},
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
				s.mockRepo.On("PatchPosition", tt.inputRow, tt.inputPositionID).
					Return(tt.mockReturn, tt.mockStatus, tt.mockErr).
					Once()
			}

			result, status, err := s.service.PatchPosition(tt.inputRow, tt.inputPositionID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			if tt.expectedStatus == 200 {
				s.Equal(tt.inputPositionID, result.ID)
			}

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PositionServiceTestSuite) TestDeletePosition() {

	sampleID := helper.NewUUIDv7()
	sampleRow := models.Position{ID: sampleID}

	var getByIdTests = []struct {
		testName       string
		inputRow       models.Position
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
			s.mockRepo.On("DeletePosition", tt.inputRow, tt.inputRow.ID).
				Return(tt.mockStatus, tt.mockErr).
				Once()

			status, err := s.service.DeletePosition(tt.inputRow, tt.inputRow.ID)

			s.Equal(tt.expectedStatus, status)
			s.Equal(tt.expectedMsg, err.Message)

			s.mockRepo.AssertExpectations(s.T())
		})
	}
}
