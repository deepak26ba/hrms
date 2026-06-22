package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type EmployeeService interface {
	CreateEmployee(row models.Employee) (int, dto.Error)
	GetEmployee(row []models.Employee, page dto.Pagination, fliter dto.EmployeeFilter) ([]models.Employee, int, int64, dto.Error)
	GetEmployeeById(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error)
	DeleteEmployee(row models.Employee, id uuid.UUID) (int, dto.Error)
	PatchEmployee(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error)
}

type employeeservice struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {

	return &employeeservice{
		repo: repo,
	}

}

func (s *employeeservice) CreateEmployee(row models.Employee) (int, dto.Error) {

	return s.repo.CreateEmployee(row)

}

func (s *employeeservice) GetEmployee(row []models.Employee, page dto.Pagination, fliter dto.EmployeeFilter) ([]models.Employee, int, int64, dto.Error) {

	return s.repo.GetEmployee(row, page, fliter)

}

func (s *employeeservice) GetEmployeeById(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error) {

	return s.repo.GetEmployeeById(row, filters)

}

func (s *employeeservice) PatchEmployee(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error) {

	return s.repo.PatchEmployee(row, filters)

}

func (s *employeeservice) DeleteEmployee(row models.Employee, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeleteEmployee(row, id)

}
