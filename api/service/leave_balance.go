package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"
)

type LeaveBalanceService interface {
	GetLeaveBalanceAdmin(row []models.LeaveBalance, page dto.Pagination, fliter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error)
	GetLeaveBalance(row []models.LeaveBalance, page dto.Pagination, fliter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error)
	GetLeaveBalanceById(row models.LeaveBalance, filter dto.GetByID) (models.LeaveBalance, int, dto.Error)
}

type leaveBalanceservice struct {
	repo repository.LeaveBalanceRepository
}

func NewLeaveBalanceService(repo repository.LeaveBalanceRepository) LeaveBalanceService {

	return &leaveBalanceservice{
		repo: repo,
	}

}

func (s *leaveBalanceservice) GetLeaveBalanceAdmin(row []models.LeaveBalance, page dto.Pagination, fliter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error) {

	return s.repo.GetLeaveBalanceAdmin(row, page, fliter)

}

func (s *leaveBalanceservice) GetLeaveBalance(row []models.LeaveBalance, page dto.Pagination, fliter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error) {

	return s.repo.GetLeaveBalance(row, page, fliter)

}

func (s *leaveBalanceservice) GetLeaveBalanceById(row models.LeaveBalance, filter dto.GetByID) (models.LeaveBalance, int, dto.Error) {

	return s.repo.GetLeaveBalanceById(row, filter)

}
