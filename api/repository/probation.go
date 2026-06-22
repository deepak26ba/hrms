package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ProbationRepository interface {
	CreateProbation(row models.Probation) (int, dto.Error)
	GetProbationAdmin(row []models.Probation, page dto.Pagination, fliter dto.ProbationFilter) ([]models.Probation, int, int64, dto.Error)
	GetProbation(row []models.Probation, page dto.Pagination, fliter dto.ProbationFilter) ([]models.Probation, int, int64, dto.Error)
	GetProbationById(row models.Probation, filters dto.GetByID) (models.Probation, int, dto.Error)
	DeleteProbation(row models.Probation, id uuid.UUID) (int, dto.Error)
	PatchProbation(row models.Probation, id uuid.UUID) (models.Probation, int, dto.Error)
}

type probationdatabase struct {
	DB *gorm.DB
}

func NewProbationRepository(db *gorm.DB) ProbationRepository {
	return &probationdatabase{
		DB: db}
}

func (d *probationdatabase) CreateProbation(row models.Probation) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Probation",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *probationdatabase) GetProbationAdmin(row []models.Probation, page dto.Pagination, fliter dto.ProbationFilter) ([]models.Probation, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Probation{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).
		Preload("Employee.Position").Preload("Employee.Department").Preload("Employee.User.Role").
		Preload("Employer.Position").Preload("Employer.Department").Preload("Employer.User.Role").
		Preload("User.Role").Preload("User.Role")

	if fliter.EmployeeId != uuid.Nil {
		query = query.Where("employee_id = ?", fliter.EmployeeId)
	}
	if fliter.UserId != uuid.Nil {
		query = query.Where("user_id = ?", fliter.UserId)
	}
	if fliter.ReviewedBy != uuid.Nil {
		query = query.Where("reviewed_by = ?", fliter.ReviewedBy)
	}
	if fliter.ProbationStatus != "" {
		query = query.Where("probation_status = ?", fliter.ProbationStatus)
	}
	if !fliter.ProbationStartDate.IsZero() {
		query = query.Where("probation_start_date = ?", fliter.ProbationStartDate)
	}
	if !fliter.ProbationEndDate.IsZero() {
		query = query.Where("probation_end_date = ?", fliter.ProbationEndDate)
	}
	if fliter.TaskCompleted != nil {
		query = query.Where("task_completed = ?", fliter.TaskCompleted)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Probation{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Probation{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Probation{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *probationdatabase) GetProbation(row []models.Probation, page dto.Pagination, fliter dto.ProbationFilter) ([]models.Probation, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Probation{}).Where("user_id = ?", fliter.UserId).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).
		Preload("Employee.Position").Preload("Employee.Department").Preload("Employee.User.Role").
		Preload("Employer.Position").Preload("Employer.Department").Preload("Employer.User.Role").
		Preload("User.Role").Preload("User.Role").Where("user_id = ?", fliter.UserId)

	if fliter.EmployeeId != uuid.Nil {
		query = query.Where("employee_id = ?", fliter.EmployeeId)
	}
	if fliter.ReviewedBy != uuid.Nil {
		query = query.Where("reviewed_by = ?", fliter.ReviewedBy)
	}
	if fliter.ProbationStatus != "" {
		query = query.Where("probation_status = ?", fliter.ProbationStatus)
	}
	if !fliter.ProbationStartDate.IsZero() {
		query = query.Where("probation_start_date = ?", fliter.ProbationStartDate)
	}
	if !fliter.ProbationEndDate.IsZero() {
		query = query.Where("probation_end_date = ?", fliter.ProbationEndDate)
	}
	if fliter.TaskCompleted != nil {
		query = query.Where("task_completed = ?", fliter.TaskCompleted)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Probation{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Probation{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Probation{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *probationdatabase) GetProbationById(row models.Probation, filters dto.GetByID) (models.Probation, int, dto.Error) {

	var err error

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {
		err = d.DB.
			Preload("Employee.Position").Preload("Employee.Department").Preload("Employee.User.Role").
			Preload("Employer.Position").Preload("Employer.Department").Preload("Employer.User.Role").
			Preload("User.Role").Preload("User.Role").First(&row, "id = ?", filters.ID).Error
	} else {
		err = d.DB.
			Preload("Employee.Position").Preload("Employee.Department").Preload("Employee.User.Role").
			Preload("Employer.Position").Preload("Employer.Department").Preload("Employer.User.Role").
			Preload("User.Role").Preload("User.Role").
			Where("user_id = ?", filters.UserID).First(&row, "id = ?", filters.ID).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Probation{}, http.StatusNotFound, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "User not found: " + filters.ID.String(),
			}
		}
		return models.Probation{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	return row, http.StatusOK, dto.Error{}
}

func (d *probationdatabase) DeleteProbation(row models.Probation, id uuid.UUID) (int, dto.Error) {

	var rows models.Probation

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

func (d *probationdatabase) PatchProbation(row models.Probation, id uuid.UUID) (models.Probation, int, dto.Error) {

	var rows models.Probation

	if err := d.DB.
		Preload("Employee.Position").Preload("Employee.Department").Preload("Employee.User.Role").
		Preload("Employer.Position").Preload("Employer.Department").Preload("Employer.User.Role").
		Preload("User.Role").Preload("User.Role").
		First(&rows, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found :" + id.String(),
			}
			return models.Probation{}, http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.Probation{}, http.StatusInternalServerError, errorResponse
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
