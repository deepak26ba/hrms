package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type LeaveRequestRepository interface {
	CreateLeaveRequest(row models.LeaveRequest) (int, dto.Error)
	GetLeaveRequestAdmin(row []models.LeaveRequest, page dto.Pagination, fliter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error)
	GetLeaveRequest(row []models.LeaveRequest, page dto.Pagination, fliter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error)
	GetLeaveRequestById(row models.LeaveRequest, filters dto.GetByID) (models.LeaveRequest, int, dto.Error)
	PatchLeaveRequest(row models.LeaveRequest, id uuid.UUID) (models.LeaveRequest, int, dto.Error)
}

type leaveRequestdatabase struct {
	DB *gorm.DB
}

func NewLeaveRequestRepository(db *gorm.DB) LeaveRequestRepository {
	return &leaveRequestdatabase{
		DB: db}
}

func (d *leaveRequestdatabase) CreateLeaveRequest(row models.LeaveRequest) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Leave Request",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *leaveRequestdatabase) GetLeaveRequestAdmin(row []models.LeaveRequest, page dto.Pagination, fliter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.LeaveRequest{}).Count(&totalRows)

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
	if fliter.ReportedTo != uuid.Nil {
		query = query.Where("reported_to = ?", fliter.ReportedTo)
	}
	if fliter.Type != "" {
		query = query.Where("type = ?", fliter.Type)
	}
	if !fliter.From.IsZero() {
		query = query.Where("from = ?", fliter.From)
	}
	if !fliter.To.IsZero() {
		query = query.Where("to = ?", fliter.To)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.LeaveRequest{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.LeaveRequest{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.LeaveRequest{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *leaveRequestdatabase) GetLeaveRequest(row []models.LeaveRequest, page dto.Pagination, fliter dto.LeaveRequestFilter) ([]models.LeaveRequest, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.LeaveRequest{}).Where("user_id = ?", fliter.UserId).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).
		Preload("Employee.Position").Preload("Employee.Department").Preload("Employee.User.Role").
		Preload("Employer.Position").Preload("Employer.Department").Preload("Employer.User.Role").
		Preload("User.Role").Preload("User.Role").
		Where("user_id = ?", fliter.UserId)

	if fliter.EmployeeId != uuid.Nil {
		query = query.Where("employee_id = ?", fliter.EmployeeId)
	}
	if fliter.ReportedTo != uuid.Nil {
		query = query.Where("reported_to = ?", fliter.ReportedTo)
	}
	if fliter.Type != "" {
		query = query.Where("type = ?", fliter.Type)
	}
	if !fliter.From.IsZero() {
		query = query.Where("from = ?", fliter.From)
	}
	if !fliter.To.IsZero() {
		query = query.Where("to = ?", fliter.To)
	}
	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.LeaveRequest{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.LeaveRequest{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.LeaveRequest{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *leaveRequestdatabase) GetLeaveRequestById(row models.LeaveRequest, filters dto.GetByID) (models.LeaveRequest, int, dto.Error) {

	var err error

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {
		err = d.DB.
			Preload("User.Role").
			Preload("Employee.Position").Preload("Employee.Department").Preload("Employee.User.Role").
			Preload("Employer.Position").Preload("Employer.Department").Preload("Employer.User.Role").
			Preload("User.Role").Preload("User.Role").
			First(&row, "id = ?", filters.ID).Error
	} else {
		err = d.DB.
			Preload("Employee.Position").Preload("Employee.Department").Preload("Employee.User.Role").
			Preload("Employer.Position").Preload("Employer.Department").Preload("Employer.User.Role").
			Preload("User.Role").Preload("User.Role").
			Where("user_id = ?", filters.UserID).First(&row, "id = ?", filters.ID).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LeaveRequest{}, http.StatusNotFound, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "User not found: " + filters.ID.String(),
			}
		}
		return models.LeaveRequest{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	return row, http.StatusOK, dto.Error{}

}

func (d *leaveRequestdatabase) PatchLeaveRequest(row models.LeaveRequest, id uuid.UUID) (models.LeaveRequest, int, dto.Error) {

	var rows models.LeaveRequest

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
			return models.LeaveRequest{}, http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.LeaveRequest{}, http.StatusInternalServerError, errorResponse
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
