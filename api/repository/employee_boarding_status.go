package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type EmployeeBoardingStatusRepository interface {
	CreateEmployeeBoardingStatus(row models.EmployeeBoardingStatus) (int, dto.Error)
	GetEmployeeBoardingStatusAdmin(row []models.EmployeeBoardingStatus, page dto.Pagination, fliter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error)
	GetEmployeeBoardingStatus(row []models.EmployeeBoardingStatus, page dto.Pagination, fliter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error)
	GetEmployeeBoardingStatusById(row models.EmployeeBoardingStatus, filters dto.GetByID) (models.EmployeeBoardingStatus, int, dto.Error)
	DeleteEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (int, dto.Error)
	PatchEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (models.EmployeeBoardingStatus, int, dto.Error)
}

type employeeBoardingStatusdatabase struct {
	DB *gorm.DB
}

func NewEmployeeBoardingStatusRepository(db *gorm.DB) EmployeeBoardingStatusRepository {
	return &employeeBoardingStatusdatabase{
		DB: db}
}

func (d *employeeBoardingStatusdatabase) CreateEmployeeBoardingStatus(row models.EmployeeBoardingStatus) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Boarding Status",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *employeeBoardingStatusdatabase) GetEmployeeBoardingStatusAdmin(row []models.EmployeeBoardingStatus, page dto.Pagination, fliter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.EmployeeBoardingStatus{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("User.Role")

	if !fliter.Offboarding.IsZero() {
		query = query.Where("off_boarding = ?", fliter.Offboarding)
	}
	if !fliter.Onboarding.IsZero() {
		query = query.Where("on_boarding = ?", fliter.Onboarding)
	}
	if fliter.UserId != uuid.Nil {
		query = query.Where("user_id = ?", fliter.UserId)
	}
	if fliter.Status != nil {
		query = query.Where("status = ?", fliter.Status)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.EmployeeBoardingStatus{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.EmployeeBoardingStatus{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.EmployeeBoardingStatus{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *employeeBoardingStatusdatabase) GetEmployeeBoardingStatus(row []models.EmployeeBoardingStatus, page dto.Pagination, fliter dto.EmployeeBoardingStatusFilter) ([]models.EmployeeBoardingStatus, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.EmployeeBoardingStatus{}).Where("user_id = ?", fliter.UserId).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("User.Role").Where("user_id = ?", fliter.UserId)

	if !fliter.Offboarding.IsZero() {
		query = query.Where("off_boarding = ?", fliter.Offboarding)
	}
	if !fliter.Onboarding.IsZero() {
		query = query.Where("on_boarding = ?", fliter.Onboarding)
	}
	if fliter.Status != nil {
		query = query.Where("status = ?", fliter.Status)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.EmployeeBoardingStatus{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.EmployeeBoardingStatus{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.EmployeeBoardingStatus{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *employeeBoardingStatusdatabase) GetEmployeeBoardingStatusById(row models.EmployeeBoardingStatus, filters dto.GetByID) (models.EmployeeBoardingStatus, int, dto.Error) {

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {

		if err := d.DB.Preload("User.Role").First(&row, filters.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errorResponse := dto.Error{
					Message:    "No Record found",
					StatusCode: http.StatusBadRequest,
					Error:      "User not found :" + filters.ID.String(),
				}
				return models.EmployeeBoardingStatus{}, http.StatusBadRequest, errorResponse
			}
			errorResponse := dto.Error{
				Message:    "InternalServerError",
				StatusCode: http.StatusInternalServerError,
				Error:      "Database error : " + err.Error(),
			}
			return models.EmployeeBoardingStatus{}, http.StatusInternalServerError, errorResponse
		}

	}
	if err := d.DB.Preload("User.Role").Where("user_id = ?", filters.UserID).First(&row, filters.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found :" + filters.ID.String(),
			}
			return models.EmployeeBoardingStatus{}, http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.EmployeeBoardingStatus{}, http.StatusInternalServerError, errorResponse
	}

	return row, http.StatusOK, dto.Error{}
}

func (d *employeeBoardingStatusdatabase) DeleteEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (int, dto.Error) {

	var rows models.EmployeeBoardingStatus

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

func (d *employeeBoardingStatusdatabase) PatchEmployeeBoardingStatus(row models.EmployeeBoardingStatus, id uuid.UUID) (models.EmployeeBoardingStatus, int, dto.Error) {

	var rows models.EmployeeBoardingStatus

	if err := d.DB.Preload("User.Role").First(&rows, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found :" + id.String(),
			}
			return models.EmployeeBoardingStatus{}, http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.EmployeeBoardingStatus{}, http.StatusInternalServerError, errorResponse
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
