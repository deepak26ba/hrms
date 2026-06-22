package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
)

type LeaveRequestService interface {
	CreateLeaveRequest(row models.LeaveRequest) (int, dto.Error)
	GetLeaveRequestAdmin(row []models.LeaveRequest, page dto.Pagination, fliter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error)
	GetLeaveRequest(row []models.LeaveRequest, page dto.Pagination, fliter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error)
	GetLeaveRequestById(row models.LeaveRequest, filter dto.GetByID) (models.LeaveRequest, int, dto.Error)
	PatchLeaveRequest(row models.LeaveRequest, id uuid.UUID) (models.LeaveRequest, int, dto.Error)
}

type leaveRequestservice struct {
	repo  repository.LeaveRequestRepository
	leave repository.LeaveBalanceRepository
}

func NewLeaveRequestService(repo repository.LeaveRequestRepository, leave repository.LeaveBalanceRepository) LeaveRequestService {

	return &leaveRequestservice{
		repo:  repo,
		leave: leave,
	}

}

func (s *leaveRequestservice) CreateLeaveRequest(row models.LeaveRequest) (int, dto.Error) {

	var leaveBalance []models.LeaveBalance
	var leaveNeedtoBeAdded models.LeaveBalance

	if row.From.After(row.To) {
		return http.StatusBadRequest, dto.Error{
			Message:    "Invalid date range: 'From' date must be before 'To' date.",
			StatusCode: http.StatusBadRequest,
			Error:      "Bad Request",
		}
	}

	totalDays, err := helper.FindDays(row.From.Format("2006-01-02"), row.To.Format("2006-01-02"))
	if err.Error != "" {
		return http.StatusInternalServerError, dto.Error{
			Message:    "Failed to calculate leave days",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error}
	}

	leaveBalanceRecord, code, _, err := s.leave.GetLeaveBalanceAdmin(
		leaveBalance,
		dto.Pagination{Page: 1,
			Limit:  1,
			Offset: 0,
			UserID: row.UserId},
		dto.LeaveBalanceFilter{
			UserId: row.UserId})
	if err.Error != "" {
		return code, dto.Error{
			Message:    "Failed to retrieve leave balance",
			StatusCode: code,
			Error:      err.Error}
	}

	leaveNeedtoBeAdded.ID = leaveBalanceRecord[0].ID
	leaveNeedtoBeAdded.MaxLeave = leaveBalanceRecord[0].MaxLeave
	leaveNeedtoBeAdded.UserId = row.UserId
	leaveNeedtoBeAdded.LeaveHolding = leaveBalanceRecord[0].LeaveHolding + uint(totalDays)
	leaveNeedtoBeAdded.LeaveRemaining = leaveBalanceRecord[0].LeaveRemaining - uint(totalDays)
	leaveNeedtoBeAdded.LeaveHolding = uint(totalDays)
	leaveNeedtoBeAdded.LeaveTaken = uint(totalDays)

	if row.Type == "Sick Leave" {
		if leaveBalanceRecord[0].SickLeave < uint(totalDays) {
			return http.StatusBadRequest, dto.Error{
				Message:    "Insufficient sick leave balance",
				StatusCode: http.StatusBadRequest,
				Error:      "Not enough sick leave available",
			}
		}
		leaveNeedtoBeAdded.SickLeave = leaveBalanceRecord[0].SickLeave - uint(totalDays)
		leaveNeedtoBeAdded.CausalLeave = leaveBalanceRecord[0].CausalLeave
		leaveNeedtoBeAdded.LossOfPay = leaveBalanceRecord[0].LossOfPay
	}

	if row.Type == "Causal Leave" {
		if leaveBalanceRecord[0].CausalLeave < uint(totalDays) {
			return http.StatusBadRequest, dto.Error{
				Message:    "Insufficient Causal leave balance",
				StatusCode: http.StatusBadRequest,
				Error:      "Not enough Causal leave available",
			}
		}
		leaveNeedtoBeAdded.CausalLeave = leaveBalanceRecord[0].CausalLeave - uint(totalDays)
		leaveNeedtoBeAdded.SickLeave = leaveBalanceRecord[0].SickLeave
		leaveNeedtoBeAdded.LossOfPay = leaveBalanceRecord[0].LossOfPay
	}
	if row.Type == "Loss Of Pay" {
		leaveNeedtoBeAdded.LossOfPay = leaveBalanceRecord[0].LossOfPay + uint(totalDays)
		leaveNeedtoBeAdded.SickLeave = leaveBalanceRecord[0].SickLeave
		leaveNeedtoBeAdded.CausalLeave = leaveBalanceRecord[0].CausalLeave
	}

	code, err = s.repo.CreateLeaveRequest(row)
	if err.Error != "" {
		return code, err
	}

	_, code, err = s.leave.PatchLeaveBalance(leaveNeedtoBeAdded, row.UserId)
	if err.Error != "" {
		return code, err
	}

	return code, dto.Error{}

}
func (s *leaveRequestservice) GetLeaveRequestAdmin(row []models.LeaveRequest, page dto.Pagination, fliter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error) {

	return s.repo.GetLeaveRequestAdmin(row, page, fliter)

}

func (s *leaveRequestservice) GetLeaveRequest(row []models.LeaveRequest, page dto.Pagination, fliter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error) {

	return s.repo.GetLeaveRequest(row, page, fliter)

}

func (s *leaveRequestservice) GetLeaveRequestById(row models.LeaveRequest, filter dto.GetByID) (models.LeaveRequest, int, dto.Error) {

	return s.repo.GetLeaveRequestById(row, filter)

}

func (s *leaveRequestservice) PatchLeaveRequest(row models.LeaveRequest, id uuid.UUID) (models.LeaveRequest, int, dto.Error) {

	result, code, err := s.repo.PatchLeaveRequest(row, id)
	if err.Error != "" {
		return result, code, err
	}

	leaveBalanceRecord, code, _, err := s.leave.GetLeaveBalanceAdmin(
		nil,
		dto.Pagination{Page: 1, Limit: 1, UserID: result.UserId},
		dto.LeaveBalanceFilter{UserId: result.UserId},
	)
	if err.Error != "" || len(leaveBalanceRecord) == 0 {
		return result, code, err
	}

	current := leaveBalanceRecord[0]
	update := models.LeaveBalance{
		ID:           current.ID,
		UserId:       current.UserId,
		MaxLeave:     current.MaxLeave,
		LeaveHolding: 0,
	}

	if result.IsApproved != nil && !*result.IsApproved {
		holding := current.LeaveHolding
		update.LeaveTaken = current.LeaveTaken - holding
		update.LeaveRemaining = current.LeaveRemaining + holding

		switch result.Type {
		case "Sick Leave":
			update.SickLeave = current.SickLeave + holding
		case "Causal Leave":
			update.CausalLeave = current.CausalLeave + holding
		case "Loss Of Pay":
			update.LossOfPay = current.LossOfPay - holding
		}
	} else {
		update.LeaveTaken = current.LeaveTaken
		update.LeaveRemaining = current.LeaveRemaining
		update.SickLeave = current.SickLeave
		update.CausalLeave = current.CausalLeave
		update.LossOfPay = current.LossOfPay
	}

	_, code, err = s.leave.PatchLeaveBalance(update, current.UserId)

	return result, code, err
}
