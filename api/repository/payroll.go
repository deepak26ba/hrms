package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PayrollRepository interface {
	CreatePayroll(row models.Payroll) (int, dto.Error)
	GetPayrollAdmin(row []models.Payroll, page dto.Pagination, fliter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error)
	GetPayroll(row []models.Payroll, page dto.Pagination, fliter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error)
	GetPayrollById(row models.Payroll, filters dto.GetByID) (models.Payroll, int, dto.Error)
	DeletePayroll(row models.Payroll, id uuid.UUID) (int, dto.Error)
	PatchPayroll(row models.Payroll, id uuid.UUID) (models.Payroll, int, dto.Error)
}

type payrolldatabase struct {
	DB *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &payrolldatabase{
		DB: db}
}

func (d *payrolldatabase) CreatePayroll(row models.Payroll) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Payroll",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *payrolldatabase) GetPayrollAdmin(row []models.Payroll, page dto.Pagination, fliter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Payroll{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("User.Role").Preload("Position")

	if fliter.Allowance != 0 {
		query = query.Where("allowance = ?", fliter.Allowance)
	}
	if fliter.BasicPay != 0 {
		query = query.Where("basic_pay = ?", fliter.BasicPay)
	}
	if fliter.Salary != 0 {
		query = query.Where("salary = ?", fliter.Salary)
	}
	if fliter.PositionId != uuid.Nil {
		query = query.Where("position_id = ?", fliter.PositionId)
	}
	if fliter.UserId != uuid.Nil {
		query = query.Where("user_id = ?", fliter.UserId)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Payroll{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Payroll{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Payroll{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *payrolldatabase) GetPayroll(row []models.Payroll, page dto.Pagination, fliter dto.PayrollFilter) ([]models.Payroll, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Payroll{}).Where("user_id = ?", fliter.UserId).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("User.Role").Preload("Position").
		Where("user_id = ?", fliter.UserId)

	if fliter.Allowance != 0 {
		query = query.Where("allowance = ?", fliter.Allowance)
	}
	if fliter.BasicPay != 0 {
		query = query.Where("basic_pay = ?", fliter.BasicPay)
	}
	if fliter.Salary != 0 {
		query = query.Where("salary = ?", fliter.Salary)
	}
	if fliter.PositionId != uuid.Nil {
		query = query.Where("position_id = ?", fliter.PositionId)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Payroll{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Payroll{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Payroll{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *payrolldatabase) GetPayrollById(row models.Payroll, filters dto.GetByID) (models.Payroll, int, dto.Error) {

	var err error

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {
		err = d.DB.
			Preload("User.Role").Preload("Position").
			First(&row, "id = ?", filters.ID).Error
	} else {
		err = d.DB.
			Preload("User.Role").Preload("Position").
			Where("user_id = ?", filters.UserID).First(&row, "id = ?", filters.ID).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Payroll{}, http.StatusNotFound, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "User not found: " + filters.ID.String(),
			}
		}
		return models.Payroll{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	return row, http.StatusOK, dto.Error{}

}

func (d *payrolldatabase) DeletePayroll(row models.Payroll, id uuid.UUID) (int, dto.Error) {

	var rows models.Payroll

	if err := d.DB.First(&rows, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found :" + id.String(),
			}
			return http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse

	}

	if err := d.DB.Delete(&row, id).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed deleting the row",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusOK, dto.Error{}
}

func (d *payrolldatabase) PatchPayroll(row models.Payroll, id uuid.UUID) (models.Payroll, int, dto.Error) {

	var rows models.Payroll

	if err := d.DB.Preload("User.Role").Preload("Position").First(&rows, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found :" + id.String(),
			}
			return models.Payroll{}, http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.Payroll{}, http.StatusInternalServerError, errorResponse
	}

	result := d.DB.Model(&rows).Updates(&row)

	if err := result.Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed updating the row ",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		return row, http.StatusInternalServerError, errorResponse
	}

	if result.RowsAffected == 0 {
		errorResponse := dto.Error{
			Message:    "Failed updating the row ",
			StatusCode: http.StatusInternalServerError,
			Error:      "ID not found/no changes made",
		}
		return row, http.StatusInternalServerError, errorResponse
	}

	return rows, http.StatusOK, dto.Error{}
}
