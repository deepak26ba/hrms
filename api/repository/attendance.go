package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	CreateAttendance(row models.Attendance) (uuid.UUID, int, dto.Error)
	GetAttendanceAdmin(row []models.Attendance, page dto.Pagination, fliter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error)
	GetAttendance(row []models.Attendance, page dto.Pagination, fliter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error)
	GetAttendanceById(row models.Attendance, filters dto.GetByID) (models.Attendance, int, dto.Error)
	DeleteAttendance(row models.Attendance, id uuid.UUID) (int, dto.Error)
	PatchAttendance(row models.Attendance, id uuid.UUID) (models.Attendance, int, dto.Error)
}

type attendancedatabase struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendancedatabase{
		DB: db}
}

func (d *attendancedatabase) CreateAttendance(row models.Attendance) (uuid.UUID, int, dto.Error) {

	ID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, http.StatusInternalServerError, dto.Error{
			Message:    "Failed to generate a valid uuid",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}
	row.ID = ID



	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to create Attandance",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return uuid.Nil, http.StatusInternalServerError, errorResponse
	}

	return row.ID, http.StatusCreated, dto.Error{}
}

func (d *attendancedatabase) GetAttendanceAdmin(row []models.Attendance, page dto.Pagination, fliter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Attendance{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("User.Role")

	if fliter.EmployeeId != uuid.Nil {
		query = query.Where("employee_id = ?", fliter.EmployeeId)
	}
	if fliter.UserId != uuid.Nil {
		query = query.Where("user_id = ?", fliter.UserId)
	}
	if !fliter.TimeIn.IsZero() {
		query = query.Where("time_in = ?", fliter.TimeIn)
	}
	if !fliter.TimeOut.IsZero() {
		query = query.Where("time_out = ?", fliter.TimeOut)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Attendance{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Attendance{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Attendance{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *attendancedatabase) GetAttendance(row []models.Attendance, page dto.Pagination, fliter dto.AttendanceFilter) ([]models.Attendance, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.Attendance{}).Where("user_id = ?", fliter.UserId).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("User.Role").
		Where("user_id = ?", fliter.UserId)

	if fliter.EmployeeId != uuid.Nil {
		query = query.Where("employee_id = ?", fliter.EmployeeId)
	}
	if !fliter.TimeIn.IsZero() {
		query = query.Where("time_in = ?", fliter.TimeIn)
	}
	if !fliter.TimeOut.IsZero() {
		query = query.Where("time_out = ?", fliter.TimeOut)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Attendance{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.Attendance{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.Attendance{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *attendancedatabase) GetAttendanceById(row models.Attendance, filters dto.GetByID) (models.Attendance, int, dto.Error) {

	var err error

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {
		err = d.DB.Preload("User.Role").First(&row, "id = ?", filters.ID).Error
	} else {
		err = d.DB.Preload("User.Role").Where("user_id = ?", filters.UserID).First(&row, "id = ?", filters.ID).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Attendance{}, http.StatusNotFound, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "User not found: " + filters.ID.String(),
			}
		}
		return models.Attendance{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	return row, http.StatusOK, dto.Error{}

}

func (d *attendancedatabase) DeleteAttendance(row models.Attendance, id uuid.UUID) (int, dto.Error) {

	var rows models.Attendance

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

func (d *attendancedatabase) PatchAttendance(row models.Attendance, id uuid.UUID) (models.Attendance, int, dto.Error) {

	var rows models.Attendance

	if err := d.DB.Preload("User.Role").First(&rows, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found :" + id.String(),
			}
			return models.Attendance{}, http.StatusBadRequest, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.Attendance{}, http.StatusInternalServerError, errorResponse
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
