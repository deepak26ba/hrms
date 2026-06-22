package helper

import (
	"hrms/common/dto"
	"net/http"
	"regexp"
	"slices"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetId(idStr string) (uuid.UUID, dto.Error) {

	if idStr == "" {
		errorResponse := dto.Error{
			Message:    "Pass value",
			StatusCode: http.StatusBadRequest,
			Error:      "Value should not be empty",
		}
		return uuid.Nil, errorResponse
	}

	id, err := StringToUUID(idStr)
	if err.Error != "" {
		return uuid.Nil, err
	}
	return id, dto.Error{}
}

func StringToUUID(idStr string) (uuid.UUID, dto.Error) {

	id, err := uuid.FromString(idStr)
	if err != nil {
		errorResponse := dto.Error{
			Message:    "Invalid Format",
			StatusCode: http.StatusBadRequest,
			Error:      "Failed to parsed : " + err.Error(),
		}
		return uuid.Nil, errorResponse
	}
	return id, dto.Error{}
}

func ValidatedPassword(password string) bool {

	emailRegex := regexp.MustCompile(`^[a-z0-9A-Z._%+\-]{8,}$`)
	return emailRegex.MatchString(password)
}

func HashPassword(password string) (int, string, dto.Error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errorResponse := dto.Error{
			Message:    "Failed Hashing",
			StatusCode: http.StatusInternalServerError,
			Error:      "Internal Error",
		}
		return http.StatusInternalServerError, "", errorResponse
	}
	return http.StatusOK, string(bytes), dto.Error{}
}

func Validator(loginCredentials any) dto.Error {

	validate := validator.New()
	err := validate.Struct(loginCredentials)
	if err != nil {
		errorResponse := dto.Error{
			Message:    "Validation failed",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		return errorResponse
	}
	return dto.Error{}
}

func IsValidPassword(storedHash, enteredPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(enteredPassword))
	if err != nil {
		return true
	}
	return false
}

func IsValidID(queryID, userID uuid.UUID) bool {

	if queryID == userID {
		return true
	}
	return false
}

func Authorize(credentials dto.Authorize) bool {

	if slices.Contains(credentials.Roles, credentials.Role) {
		if credentials.Role == "EMPLOYEE" {
			if IsValidID(credentials.QueryID, credentials.UserID) {
				return true
			}
		}
		return true
	}
	return false
}

func FindOffset(page, limit int) *dto.Pagination {

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	pageParams := dto.Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
	return &pageParams
}

func FindNextPage(page, limit int, totalRow int64) bool {

	if page*limit < int(totalRow) {
		return true
	}
	return false
}

func FindDays(startDate, endDate string) (int, dto.Error) {

	layout := "2006-01-02"

	start, err := time.Parse(layout, startDate)
	if err != nil {
		return 0, dto.Error{Message: err.Error()}
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return 0, dto.Error{Message: err.Error()}
	}

	hours := end.Sub(start).Hours()
	days := int(hours/24+0.5) + 1
	return days, dto.Error{}
}

func NewUUIDv7() uuid.UUID {
	id, _ := uuid.NewV7()
	return id
}
