package service

import (
	"hrms/api/repository"
	"hrms/common/dto"
	"hrms/pkg/models"

	"github.com/gofrs/uuid"
)

type PayrollService interface {
	CreatePayroll(row models.Payroll) (int, dto.Error)
	GetPayrollAdmin(row []models.Payroll, page dto.Pagination, fliter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error)
	GetPayroll(row []models.Payroll, page dto.Pagination, fliter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error)
	GetPayrollById(row models.Payroll, filter dto.GetByID) (models.Payroll, int, dto.Error)
	DeletePayroll(row models.Payroll, id uuid.UUID) (int, dto.Error)
	PatchPayroll(row models.Payroll, id uuid.UUID) (models.Payroll, int, dto.Error)
}

type payrollservice struct {
	repo repository.PayrollRepository
}

func NewPayrollService(repo repository.PayrollRepository) PayrollService {

	return &payrollservice{
		repo: repo,
	}

}

func (s *payrollservice) CreatePayroll(row models.Payroll) (int, dto.Error) {

	return s.repo.CreatePayroll(row)

}
func (s *payrollservice) GetPayrollAdmin(row []models.Payroll, page dto.Pagination, fliter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error) {

	return s.repo.GetPayrollAdmin(row, page, fliter)

}

func (s *payrollservice) GetPayroll(row []models.Payroll, page dto.Pagination, fliter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error) {

	return s.repo.GetPayroll(row, page, fliter)

}

func (s *payrollservice) GetPayrollById(row models.Payroll, filter dto.GetByID) (models.Payroll, int, dto.Error) {

	return s.repo.GetPayrollById(row, filter)

}

func (s *payrollservice) PatchPayroll(row models.Payroll, id uuid.UUID) (models.Payroll, int, dto.Error) {

	return s.repo.PatchPayroll(row, id)

}

func (s *payrollservice) DeletePayroll(row models.Payroll, id uuid.UUID) (int, dto.Error) {

	return s.repo.DeletePayroll(row, id)

}
