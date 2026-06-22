package repository

import (
	"errors"
	"hrms/common/dto"
	"hrms/pkg/models"
	"net/http"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(row models.User) (int, dto.Error)
	Login(email string) (models.User, int, dto.Error)
}

type authdatabase struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authdatabase{
		DB: db}
}

func (d *authdatabase) Login(email string) (models.User, int, dto.Error) {

	var row models.User

	err := d.DB.Preload("Role").Where("email = ?", email).First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := dto.Error{
				Message:    "Enter valid user/Register before login",
				StatusCode: http.StatusUnauthorized,
				Error:      "User not found :" + email,
			}
			return models.User{}, http.StatusUnauthorized, errorResponse
		}
		errorResponse := dto.Error{
			Message:    "InternalServerError",
			StatusCode: http.StatusInternalServerError,
			Error:      "Database error : " + err.Error(),
		}
		return models.User{}, http.StatusInternalServerError, errorResponse
	}
	return row, http.StatusOK, dto.Error{}
}

func (d *authdatabase) Register(row models.User) (int, dto.Error) {

	if err := d.DB.Create(&row).Error; err != nil {
		errorResponse := dto.Error{
			Message:    "Failed to Register",
			StatusCode: http.StatusInternalServerError,
			Error:      "Failed inserting the row : " + err.Error(),
		}
		return http.StatusInternalServerError, errorResponse
	}

	return http.StatusCreated, dto.Error{}
}
