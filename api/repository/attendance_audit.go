package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AttendanceAuditRepository interface {
	CreateAttendanceAudit(row models.AttendanceAudit) (int, dto.Error)
	GetAttendanceAuditAdmin(row []models.AttendanceAudit, page dto.Pagination, fliter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error)
	GetAttendanceAudit(row []models.AttendanceAudit, page dto.Pagination, fliter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error)
	GetAttendanceAuditById(row models.AttendanceAudit, filters dto.GetByID) (models.AttendanceAudit, int, dto.Error)
	PatchAttendanceAudit(row models.AttendanceAudit, userID uuid.UUID) (models.AttendanceAudit, int, dto.Error)
}

type attendanceAuditdatabase struct {
	DB *gorm.DB
}

func NewAttendanceAuditRepository(db *gorm.DB) AttendanceAuditRepository {
	return &attendanceAuditdatabase{
		DB: db}
}

func (d *attendanceAuditdatabase) CreateAttendanceAudit(row models.AttendanceAudit) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Create Attandance Audit",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}

func (d *attendanceAuditdatabase) GetAttendanceAuditAdmin(row []models.AttendanceAudit, page dto.Pagination, fliter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.AttendanceAudit{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("Attendance.User.Role")

	if !fliter.TimeIn.IsZero() {
		query = query.Where("time_in = ?", fliter.TimeIn)
	}
	if !fliter.TimeOut.IsZero() {
		query = query.Where("time_out = ?", fliter.TimeOut)
	}
	if fliter.AttendanceID != uuid.Nil {
		query = query.Where("attendance_id = ?", fliter.AttendanceID)
	}
	if fliter.UserId != uuid.Nil {
		query = query.Where("user_id = ?", fliter.UserId)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.AttendanceAudit{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.AttendanceAudit{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.AttendanceAudit{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *attendanceAuditdatabase) GetAttendanceAudit(row []models.AttendanceAudit, page dto.Pagination, fliter dto.AttendanceAuditFilter) ([]models.AttendanceAudit, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.AttendanceAudit{}).Where("user_id = ?", fliter.UserId).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("Attendance.User.Role").Where("user_id = ?", fliter.UserId)

	if !fliter.TimeIn.IsZero() {
		query = query.Where("time_in = ?", fliter.TimeIn)
	}
	if !fliter.TimeOut.IsZero() {
		query = query.Where("time_out = ?", fliter.TimeOut)
	}
	if fliter.AttendanceID != uuid.Nil {
		query = query.Where("attendance_id = ?", fliter.AttendanceID)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.AttendanceAudit{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.AttendanceAudit{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.AttendanceAudit{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *attendanceAuditdatabase) GetAttendanceAuditById(row models.AttendanceAudit, filters dto.GetByID) (models.AttendanceAudit, int, dto.Error) {

	var err error

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {
		err = d.DB.Preload("Attendance.User.Role").First(&row, "id = ?", filters.ID).Error
	} else {
		err = d.DB.Preload("Attendance.User.Role").
			Where("user_id = ?", filters.UserID).First(&row, "id = ?", filters.ID).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.AttendanceAudit{}, http.StatusNotFound, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "User not found: " + filters.ID.String(),
			}
		}
		return models.AttendanceAudit{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	return row, http.StatusOK, dto.Error{}

}

func (d *attendanceAuditdatabase) PatchAttendanceAudit(row models.AttendanceAudit, userID uuid.UUID) (models.AttendanceAudit, int, dto.Error) {

	var rows models.AttendanceAudit

	if err := d.DB.Preload("Attendance.User.Role").
		Where("user_id = ?", userID).Where("attendance_id = ?", row.AttendanceID).
		First(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found :" + row.AttendanceID.String(),
			}
			return models.AttendanceAudit{}, http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.AttendanceAudit{}, http.StatusInternalServerError, errorResponse
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
