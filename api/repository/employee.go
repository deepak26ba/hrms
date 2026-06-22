package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateEmployee(row models.Employee) (int, dto.Error)
	GetEmployee(row []models.Employee, page dto.Pagination, fliter dto.EmployeeFilter) ([]models.Employee, int, int64, dto.Error)
	GetEmployeeById(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error)
	DeleteEmployee(row models.Employee, id uuid.UUID) (int, dto.Error)
	PatchEmployee(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error)
}

type employeedatabase struct {
	DB *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeedatabase{
		DB: db}
}

func (d *employeedatabase) CreateEmployee(row models.Employee) (int, dto.Error) {

	var leave models.LeaveBalance

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Employee",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}
	
	leave.UserId = row.UserId
	
	if err := d.DB.Create(&leave).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Employee",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *employeedatabase) GetEmployee(row []models.Employee, page dto.Pagination, filter dto.EmployeeFilter) ([]models.Employee, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Employee{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).
		Preload("Position").Preload("Department").Preload("User.Role")

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}
	if filter.Age != 0 {
		query = query.Where("age = ?", filter.Age)
	}
	if filter.PhoneNumber != "" {
		query = query.Where("phone_number = ?", filter.PhoneNumber)
	}
	if filter.DepartmentId != uuid.Nil {
		query = query.Where("department_id = ?", filter.DepartmentId)
	}
	if filter.PositionId != uuid.Nil {
		query = query.Where("position_id = ?", filter.PositionId)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", filter.IsActive)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Employee{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Employee{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}
	if len(row) == 0 {
		return []models.Employee{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *employeedatabase) GetEmployeeById(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error) {

	var err error

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {
		err = d.DB.Preload("Position").Preload("Department").Preload("User.Role").First(&row, "id = ?", filters.ID).Error
	} else {
		err = d.DB.Preload("Position").Preload("Department").Preload("User.Role").
			Where("user_id = ?", filters.UserID).First(&row, "id = ?", filters.ID).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Employee{}, http.StatusNotFound, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "User not found: " + filters.ID.String(),
			}
		}
		return models.Employee{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	return row, http.StatusOK, dto.Error{}
}

func (d *employeedatabase) DeleteEmployee(row models.Employee, id uuid.UUID) (int, dto.Error) {

	var rows models.Employee

	if err := d.DB.First(&rows, id).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "ID not found  ",
			StatusCode: http.StatusUnauthorized,
			Error:      err.Error(),
		}
		return http.StatusUnauthorized, errorResponse

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

func (d *employeedatabase) PatchEmployee(row models.Employee, filters dto.GetByID) (models.Employee, int, dto.Error) {

	var err error
	var rows models.Employee

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {
		err = d.DB.Preload("Position").Preload("Department").Preload("User.Role").First(&rows, "id = ?", filters.ID).Error
	} else {
		err = d.DB.Preload("Position").Preload("Department").Preload("User.Role").
			Where("user_id = ?", filters.UserID).First(&rows, "id = ?", filters.ID).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Employee{}, http.StatusNotFound, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "User not found: " + filters.ID.String(),
			}
		}
		return models.Employee{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	result := d.DB.Model(&rows).Updates(&row)

	if err := result.Error; err != nil {
		return row, http.StatusInternalServerError, dto.Error{
			Message:    "Failed updating the row ",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}

	if result.RowsAffected == 0 {
		return row, http.StatusInternalServerError, dto.Error{
			Message:    "Failed updating the row ",
			StatusCode: http.StatusInternalServerError,
			Error:      "ID not found/no changes made",
		}
	}

	return rows, http.StatusOK, dto.Error{}
}
