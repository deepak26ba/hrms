package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type DepartmentService interface {
	CreateDepartment(row models.Department) (int, dto.Error)
	GetDepartment(row []models.Department, page dto.Pagination, filter dto.DepartmentFilter) ([]models.Department, int, int64, dto.Error)
	GetDepartmentById(row models.Department, id uuid.UUID) (models.Department, int, dto.Error)
	DeleteDepartment(row models.Department, id uuid.UUID) (int, dto.Error)
	PatchDepartment(row models.Department, id uuid.UUID) (models.Department, int, dto.Error)
}

type departmentservice struct {
	repo repository.DepartmentRepository
}

func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {

	return &departmentservice{
		repo: repo,
	}

}

func (s *departmentservice) CreateDepartment(row models.Department) (int, dto.Error) {

	return s.repo.CreateDepartment(row)

}

func (s *departmentservice) GetDepartment(row []models.Department, page dto.Pagination, filter dto.DepartmentFilter) ([]models.Department, int, int64, dto.Error) {

	return s.repo.GetDepartment(row, page, filter)

}

func (s *departmentservice) GetDepartmentById(row models.Department, id uuid.UUID) (models.Department, int, dto.Error) {

	return s.repo.GetDepartmentById(row, id)

}

func (s *departmentservice) PatchDepartment(row models.Department, id uuid.UUID) (models.Department, int, dto.Error) {

	return s.repo.PatchDepartment(row, id)

}

func (s *departmentservice) DeleteDepartment(row models.Department, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeleteDepartment(row, id)

}
