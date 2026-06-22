package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type EmployeeBoardingStatusService interface {
	CreateEmployeeBoardingStatus(row models.EmployeeBoardingStatus) (int, dto.Error)
	GetEmployeeBoardingStatusAdmin(row []models.EmployeeBoardingStatus, page dto.Pagination, fliter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error)
	GetEmployeeBoardingStatus(row []models.EmployeeBoardingStatus, page dto.Pagination, fliter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error)
	GetEmployeeBoardingStatusById(row models.EmployeeBoardingStatus, filter dto.GetByID) (models.EmployeeBoardingStatus, int, dto.Error)
	DeleteEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (int, dto.Error)
	PatchEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (models.EmployeeBoardingStatus, int, dto.Error)
}

type employeeBoardingStatusservice struct {
	repo repository.EmployeeBoardingStatusRepository
}

func NewEmployeeBoardingStatusService(repo repository.EmployeeBoardingStatusRepository) EmployeeBoardingStatusService {

	return &employeeBoardingStatusservice{
		repo: repo,
	}

}

func (s *employeeBoardingStatusservice) CreateEmployeeBoardingStatus(row models.EmployeeBoardingStatus) (int, dto.Error) {

	return s.repo.CreateEmployeeBoardingStatus(row)

}
func (s *employeeBoardingStatusservice) GetEmployeeBoardingStatusAdmin(row []models.EmployeeBoardingStatus, page dto.Pagination, fliter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error) {

	return s.repo.GetEmployeeBoardingStatusAdmin(row, page, fliter)

}

func (s *employeeBoardingStatusservice) GetEmployeeBoardingStatus(row []models.EmployeeBoardingStatus, page dto.Pagination, fliter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error) {

	return s.repo.GetEmployeeBoardingStatus(row, page, fliter)

}

func (s *employeeBoardingStatusservice) GetEmployeeBoardingStatusById(row models.EmployeeBoardingStatus, filter dto.GetByID) (models.EmployeeBoardingStatus, int, dto.Error) {

	return s.repo.GetEmployeeBoardingStatusById(row, filter)

}

func (s *employeeBoardingStatusservice) PatchEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (models.EmployeeBoardingStatus, int, dto.Error) {

	return s.repo.PatchEmployeeBoardingStatus(row, id)

}

func (s *employeeBoardingStatusservice) DeleteEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeleteEmployeeBoardingStatus(row, id)

}
