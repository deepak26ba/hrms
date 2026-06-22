package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PositionRepository interface {
	CreatePosition(row models.Position) (int, dto.Error)
	GetPosition(row []models.Position, page dto.Pagination, filter dto.PositionFilter) ([]models.Position, int, int64, dto.Error)
	GetPositionById(row models.Position, id uuid.UUID) (models.Position, int, dto.Error)
	DeletePosition(row models.Position, id uuid.UUID) (int, dto.Error)
	PatchPosition(row models.Position, id uuid.UUID) (models.Position, int, dto.Error)
}

type positiondatabase struct {
	DB *gorm.DB
}

func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &positiondatabase{
		DB: db}
}

func (d *positiondatabase) CreatePosition(row models.Position) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Position",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *positiondatabase) GetPosition(row []models.Position, page dto.Pagination, filter dto.PositionFilter) ([]models.Position, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Position{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit)

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Position{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Position{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Position{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d positiondatabase) GetPositionById(row models.Position, id uuid.UUID) (models.Position, int, dto.Error) {

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

func (d *positiondatabase) DeletePosition(row models.Position, id uuid.UUID) (int, dto.Error) {

	var rows models.Position

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

func (d positiondatabase) PatchPosition(row models.Position, id uuid.UUID) (models.Position, int, dto.Error) {

	var rows models.Position

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
