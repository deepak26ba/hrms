package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type RolesRepository interface {
	CreateRoles(row models.Roles) (int, dto.Error)
	GetRoles(row []models.Roles, page dto.Pagination, filter dto.RoleFilter) ([]models.Roles, int, int64, dto.Error)
	GetRolesById(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error)
	DeleteRoles(row models.Roles, id uuid.UUID) (int, dto.Error)
	PatchRoles(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error)
}

type roledatabase struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RolesRepository {
	return &roledatabase{
		DB: db}
}

func (d *roledatabase) CreateRoles(row models.Roles) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Roles",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *roledatabase) GetRoles(row []models.Roles, page dto.Pagination, filter dto.RoleFilter) ([]models.Roles, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Roles{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit)

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Roles{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Roles{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Roles{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *roledatabase) GetRolesById(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error) {

	result := d.DB.First(&row, id)
	if result.Error != nil {
		errorResponse := dto.Error{
			Message:    "ID not found",
			StatusCode: http.StatusBadRequest,
			Error:      result.Error.Error(),
		}
		return row, http.StatusBadRequest, errorResponse
	}

	return row, http.StatusOK, dto.Error{}
}

func (d *roledatabase) DeleteRoles(row models.Roles, id uuid.UUID) (int, dto.Error) {

	var rows models.Roles

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

func (d *roledatabase) PatchRoles(row models.Roles, id uuid.UUID) (models.Roles, int, dto.Error) {

	var rows models.Roles

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
