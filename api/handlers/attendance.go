package handlers

import (
	"fmt"
	"hrms/api/service"
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofrs/uuid"
)

type AttendanceHandler struct {
	service         service.AttendanceService
	attendanceAduit service.AttendanceAuditService
	audit           service.AuditService
}

func NewAttendanceHandler(service service.AttendanceService, audit service.AuditService, attendanceAduit service.AttendanceAuditService) *AttendanceHandler {

	return &AttendanceHandler{
		service:         service,
		audit:           audit,
		attendanceAduit: attendanceAduit,
	}

}

func (h *AttendanceHandler) CreateAttendance(f fiber.Ctx) error {

	currentUserID, ok := f.Locals("user_id").(string)
	if !ok {
		errorResponse := &dto.Error{
			Message:    "Need Authentication",
			StatusCode: fiber.StatusForbidden,
			Error:      "Forbidden",
		}
		log := models.AuditTable{
			UserID: currentUserID,
			Action: "Create",
			Level:  "Handler",
			Data:   fmt.Sprint(currentUserID),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
	}

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "Create",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	userID, errorResponse := helper.StringToUUID(currentUserID)
	if errorResponse.Error != "" {
		audit(userID, "Invalid ID format")
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	payload := models.Attendance{
		UserId: userID,
		TimeIn: time.Now(),
	}

	id, code, errorResponse := h.service.CreateAttendance(payload)
	if errorResponse.Error != "" {
		audit(code, errorResponse)
		return f.Status(code).JSON(errorResponse)

	}

	attendanceAudit := models.AttendanceAudit{
		AttendanceID: id,
		UserId:       payload.UserId,
		TimeIn:       payload.TimeIn,
		TimeOut:      payload.TimeOut,
	}

	code, errorResponse = h.attendanceAduit.CreateAttendanceAudit(attendanceAudit)
	if errorResponse.Error != "" {
		audit(code, errorResponse)
		return f.Status(code).JSON(errorResponse)

	}

	successResponse := &dto.Success{
		Message:    "Successfully Created",
		StatusCode: fiber.StatusCreated,
	}
	audit("Success", "nil")

	return f.Status(fiber.StatusCreated).JSON(successResponse)

}

func (h *AttendanceHandler) GetAttendance(f fiber.Ctx) error {

	var payload, result []models.Attendance
	var code int
	var count int64
	var errorResponse dto.Error

	currentUserID, _ := f.Locals("user_id").(string)
	role, _ := f.Locals("role").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "View",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	pageParams := &dto.Pagination{Page: 1, Limit: 5, Offset: 0, UserID: uuid.Nil}

	if err := f.Bind().Query(pageParams); err != nil {
		errorResponse := &dto.Error{
			Message:    "Invalid Query Parameters",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}
	pageParams = helper.FindOffset(pageParams.Page, pageParams.Limit)

	filter := &dto.AttendanceFilter{}

	if err := f.Bind().Query(filter); err != nil {
		errorResponse := &dto.Error{
			Message:    "Invalid Query Parameters",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if role == "ADMIN" {
		result, code, count, errorResponse = h.service.GetAttendanceAdmin(payload, *pageParams, *filter)
	} else {
		if currentUserID == "" {
			return f.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Authentication Required"})
		}

		id, uuidErr := helper.StringToUUID(currentUserID)
		if uuidErr.Error != "" {
			audit(currentUserID, uuidErr)
			return f.Status(fiber.StatusBadRequest).JSON(uuidErr)
		}

		filter.UserId = id
		result, code, count, errorResponse = h.service.GetAttendance(payload, *pageParams, *filter)
	}

	if errorResponse.Error != "" {
		audit(code, errorResponse)
		return f.Status(code).JSON(errorResponse)
	}

	successResponse := &dto.Success{
		Message:    "Successfully Received",
		StatusCode: fiber.StatusOK,
		Data:       result,
		Pagination: dto.PaginationResult{
			Page:        pageParams.Page,
			Limit:       pageParams.Limit,
			TotalCount:  int(count),
			HasNextPage: helper.FindNextPage(pageParams.Page, pageParams.Limit, count),
		},
	}

	audit(successResponse, "nil")
	return f.Status(fiber.StatusOK).JSON(successResponse)
}

func (h *AttendanceHandler) GetAttendanceById(f fiber.Ctx) error {

	var payload models.Attendance

	currentUserID, _ := f.Locals("user_id").(string)
	role, _ := f.Locals("role").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "ViewByID",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if currentUserID == "" || role == "" {
		errorResponse := dto.Error{Message: "Forbidden", StatusCode: fiber.StatusForbidden, Error: "Missing Auth"}
		audit("AuthCheck", errorResponse)
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

	id, errorResponse := helper.GetId(f.Params("id"))
	if errorResponse.Error != "" {
		audit(f.Params("id"), errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	userUUID, errorResponse := helper.StringToUUID(currentUserID)
	if errorResponse.Error != "" {
		audit(currentUserID, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	filters := dto.GetByID{
		ID:     id,
		UserID: userUUID,
		Role:   role,
	}

	result, code, errorResponse := h.service.GetAttendanceById(payload, filters)
	if errorResponse.Error != "" {
		audit(filters, errorResponse)
		return f.Status(code).JSON(errorResponse)
	}

	successResponse := &dto.Success{
		Message:    "Successfully Received",
		StatusCode: fiber.StatusOK,
		Data:       result,
	}

	audit(successResponse, "nil")
	return f.Status(fiber.StatusOK).JSON(successResponse)
}

func (h *AttendanceHandler) DeleteAttendance(f fiber.Ctx) error {

	var payload models.Attendance

	currentUserID, _ := f.Locals("user_id").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "Delete",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if currentUserID == "" {
		errorResponse := dto.Error{
			Message:    "Forbidden",
			StatusCode: fiber.StatusForbidden,
			Error:      "Missing Auth"}
		audit("AuthCheck", errorResponse)
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

	id, errorResponse := helper.GetId(f.Params("id"))
	if errorResponse.Error != "" {
		audit(f.Params("id"), errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	code, errorResponse := h.service.DeleteAttendance(payload, id)
	if errorResponse.Error != "" {
		audit(code, errorResponse)
		return f.Status(code).JSON(errorResponse)
	}

	successResponse := &dto.Success{
		Message:    "Successfully Deleted",
		StatusCode: fiber.StatusOK,
	}
	audit(successResponse, "nil")
	return f.Status(fiber.StatusOK).JSON(successResponse)

}

func (h *AttendanceHandler) PatchAttendance(f fiber.Ctx) error {

	currentUserID, ok := f.Locals("user_id").(string)
	if !ok {
		errorResponse := &dto.Error{
			Message:    "Need Authentication",
			StatusCode: fiber.StatusForbidden,
			Error:      "Forbidden",
		}
		log := models.AuditTable{
			UserID: currentUserID,
			Action: "Create",
			Level:  "Handler",
			Data:   fmt.Sprint(currentUserID),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
	}

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "Update",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if currentUserID == "" {
		errorResponse := dto.Error{
			Message:    "Forbidden",
			StatusCode: fiber.StatusForbidden,
			Error:      "Missing Auth"}
		audit("AuthCheck", errorResponse)
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

	id, errorResponse := helper.GetId(f.Params("id"))
	if errorResponse.Error != "" {
		audit(f.Params("id"), errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	userID, errorResponse := helper.StringToUUID(currentUserID)
	if errorResponse.Error != "" {
		audit(userID, "Invalid ID format")
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	payload := models.Attendance{
		UserId:  userID,
		TimeOut: time.Now(),
	}

	result, code, errorResponse := h.service.PatchAttendance(payload, id)
	if errorResponse.Error != "" {
		audit(result, errorResponse)
		return f.Status(code).JSON(errorResponse)
	}

	attendanceAudit := models.AttendanceAudit{
		AttendanceID: result.ID,
		TimeOut:      result.TimeOut,
	}

	_, code, errorResponse = h.attendanceAduit.PatchAttendanceAudit(attendanceAudit, userID)
	if errorResponse.Error != "" {
		audit(code, errorResponse)
		return f.Status(code).JSON(errorResponse)

	}

	successResponse := &dto.Success{
		Message:    "Successfully Updated",
		StatusCode: fiber.StatusOK,
		Data:       result,
	}
	audit(successResponse, "nil")

	return f.Status(fiber.StatusOK).JSON(successResponse)
}
