package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	CreateDepartment(row models.Department) (int, dto.Error)
	GetDepartment(row []models.Department, page dto.Pagination, filter dto.DepartmentFilter) ([]models.Department, int, int64, dto.Error)
	GetDepartmentById(row models.Department, id uuid.UUID) (models.Department, int, dto.Error)
	DeleteDepartment(row models.Department, id uuid.UUID) (int, dto.Error)
	PatchDepartment(row models.Department, id uuid.UUID) (models.Department, int, dto.Error)
}

type departmentdatabase struct {
	DB *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentdatabase{
		DB: db}
}

func (d departmentdatabase) CreateDepartment(row models.Department) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Department",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *departmentdatabase) GetDepartment(row []models.Department, page dto.Pagination, filter dto.DepartmentFilter) ([]models.Department, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Department{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit)

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Department{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Department{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Department{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *departmentdatabase) GetDepartmentById(row models.Department, id uuid.UUID) (models.Department, int, dto.Error) {

	result := d.DB.First(&row, id)
	if result.Error != nil {
		errorResponse := dto.Error{
			Message:    "ID not found",
			StatusCode: http.StatusUnauthorized,
			Error:      result.Error.Error(),
		}
		return row, http.StatusUnauthorized, errorResponse
	}

	return row, http.StatusOK, dto.Error{}
}

func (d *departmentdatabase) DeleteDepartment(row models.Department, id uuid.UUID) (int, dto.Error) {

	var rows models.Department

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

func (d *departmentdatabase) PatchDepartment(row models.Department, id uuid.UUID) (models.Department, int, dto.Error) {

	var rows models.Department

	if err := d.DB.First(&rows, id).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "ID not found",
			StatusCode: http.StatusUnauthorized,
			Error:      err.Error(),
		}
		return row, http.StatusUnauthorized, errorResponse
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
