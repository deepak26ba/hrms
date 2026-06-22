package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type LeaveBalanceRepository interface {
	GetLeaveBalanceAdmin(row []models.LeaveBalance, page dto.Pagination, fliter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error)
	GetLeaveBalance(row []models.LeaveBalance, page dto.Pagination, fliter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error)
	GetLeaveBalanceById(row models.LeaveBalance, filters dto.GetByID) (models.LeaveBalance, int, dto.Error)
	PatchLeaveBalance(row models.LeaveBalance, id uuid.UUID) (models.LeaveBalance, int, dto.Error)
}

type leaveBalancedatabase struct {
	DB *gorm.DB
}

func NewLeaveBalanceRepository(db *gorm.DB) LeaveBalanceRepository {
	return &leaveBalancedatabase{
		DB: db}
}

func (d *leaveBalancedatabase) GetLeaveBalanceAdmin(row []models.LeaveBalance, page dto.Pagination, fliter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.LeaveBalance{}).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("User.Role")

	if fliter.UserId != uuid.Nil {
		query = query.Where("user_id = ?", fliter.UserId)
	}
	if fliter.MaxLeave != 0 {
		query = query.Where("max_leave = ?", fliter.MaxLeave)
	}
	if fliter.LeaveTaken != 0 {
		query = query.Where("leave_taken = ?", fliter.LeaveTaken)
	}
	if fliter.LeaveRemaining != 0 {
		query = query.Where("leave_remaining = ?", fliter.LeaveRemaining)
	}
	if fliter.SickLeave != 0 {
		query = query.Where("sick_leave = ?", fliter.SickLeave)
	}
	if fliter.CausalLeave != 0 {
		query = query.Where("causal_leave = ?", fliter.CausalLeave)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.LeaveBalance{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.LeaveBalance{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.LeaveBalance{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *leaveBalancedatabase) GetLeaveBalance(row []models.LeaveBalance, page dto.Pagination, fliter dto.LeaveBalanceFilter) ([]models.LeaveBalance, int, int64, dto.Error) {

	var totalRows int64

	d.DB.Model(&models.LeaveBalance{}).Where("user_id = ?", fliter.UserId).Count(&totalRows)

	query := d.DB.Offset(page.Offset).Limit(page.Limit).Preload("User.Role").Where("user_id = ?", fliter.UserId)

	if fliter.MaxLeave != 0 {
		query = query.Where("max_leave = ?", fliter.MaxLeave)
	}
	if fliter.LeaveTaken != 0 {
		query = query.Where("leave_taken = ?", fliter.LeaveTaken)
	}
	if fliter.LeaveRemaining != 0 {
		query = query.Where("leave_remaining = ?", fliter.LeaveRemaining)
	}
	if fliter.SickLeave != 0 {
		query = query.Where("sick_leave = ?", fliter.SickLeave)
	}
	if fliter.CausalLeave != 0 {
		query = query.Where("causal_leave = ?", fliter.CausalLeave)
	}

	if err := query.Find(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.LeaveBalance{}, http.StatusNotFound, totalRows, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "Record not found in database",
			}
		}
		return []models.LeaveBalance{}, http.StatusInternalServerError, totalRows, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	if len(row) == 0 {
		return []models.LeaveBalance{}, http.StatusOK, totalRows, dto.Error{}
	}

	return row, http.StatusOK, totalRows, dto.Error{}
}

func (d *leaveBalancedatabase) GetLeaveBalanceById(row models.LeaveBalance, filters dto.GetByID) (models.LeaveBalance, int, dto.Error) {

	var err error

	if filters.Role == "ADMIN" || filters.Role == "VP" || filters.Role == "HR" {
		err = d.DB.Preload("User.Role").
			First(&row, "id = ?", filters.ID).Error
	} else {
		err = d.DB.Preload("User.Role").
			Where("user_id = ?", filters.UserID).First(&row, "id = ?", filters.ID).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LeaveBalance{}, http.StatusNotFound, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusNotFound,
				Error:      "User not found: " + filters.ID.String(),
			}
		}
		return models.LeaveBalance{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error: " + err.Error(),
		}
	}

	return row, http.StatusOK, dto.Error{}

}

func (d *leaveBalancedatabase) PatchLeaveBalance(row models.LeaveBalance, id uuid.UUID) (models.LeaveBalance, int, dto.Error) {

	var existing models.LeaveBalance

	if err := d.DB.Where("user_id = ?", id).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LeaveBalance{}, http.StatusBadRequest, dto.Error{
				Message:    "No Record found",
				StatusCode: http.StatusBadRequest,
				Error:      "User not found: " + id.String(),
			}
		}
		return models.LeaveBalance{}, http.StatusInternalServerError, dto.Error{
			Message:    "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
	}

	result := d.DB.Model(&existing).Select("*").Updates(&row)

	if result.Error != nil {
		return models.LeaveBalance{}, http.StatusInternalServerError, dto.Error{
			Message:    "Failed updating row",
			StatusCode: http.StatusInternalServerError,
			Error:      result.Error.Error(),
		}
	}

	if result.RowsAffected == 0 {
		return models.LeaveBalance{}, http.StatusBadRequest, dto.Error{
			Message:    "No changes made",
			StatusCode: http.StatusBadRequest,
			Error:      "Update failed",
		}
	}

	d.DB.Where("user_id = ?", id).First(&existing)

	return existing, http.StatusOK, dto.Error{}
}
