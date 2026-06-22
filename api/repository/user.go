package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(row []models.User, page dto.Pagination, filter dto.UserFilter) ([]models.User, int, int64, dto.Error)
	GetUserById(row models.User, id uuid.UUID) (models.User, int, dto.Error)
	DeleteUser(row models.User, id uuid.UUID) (int, dto.Error)
	PatchUser(row models.User, id uuid.UUID) (models.User, int, dto.Error)
}

type userdatabase struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userdatabase{
		DB: db}
}

func (d *userdatabase) GetUser(row []models.User, page dto.Pagination, fliter dto.UserFilter) ([]models.User, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.User{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("Role")

	if fliter.Email != "" {
		query = query.Where("email = ?", fliter.Email)
	}
	if fliter.RoleID != uuid.Nil {
		query = query.Where("role_id = ?", fliter.RoleID)
	}
	if fliter.IsActive != nil {
		query = query.Where("is_active = ?", fliter.IsActive)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.User{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.User{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.User{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *userdatabase) GetUserById(row models.User, id uuid.UUID) (models.User, int, dto.Error) {

	result := d.DB.Preload("Role").First(&row, id)
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

func (d *userdatabase) DeleteUser(row models.User, id uuid.UUID) (int, dto.Error) {

	var rows models.User

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

func (d *userdatabase) PatchUser(row models.User, id uuid.UUID) (models.User, int, dto.Error) {

	var rows models.User

	if err := d.DB.Preload("Role").First(&rows, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found :" + id.String(),
			}
			return models.User{}, http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.User{}, http.StatusInternalServerError, errorResponse
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
