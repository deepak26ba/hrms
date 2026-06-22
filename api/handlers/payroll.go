package handlers

import (
	"fmt"
	"hrms/api/service"
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"

	"github.com/gofiber/fiber/v3"
	"github.com/gofrs/uuid"
)

type PayrollHandler struct {
	service service.PayrollService
	audit   service.AuditService
}

func NewPayrollHandler(service service.PayrollService, audit service.AuditService) *PayrollHandler {

	return &PayrollHandler{
		service: service,
		audit:   audit,
	}

}

func (h *PayrollHandler) CreatePayroll(f fiber.Ctx) error {

	var payload models.Payroll

	currentUserID, _ := f.Locals("user_id").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "Create",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if err := f.Bind().Body(&payload); err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed converting json to struct",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(payload, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)

	}

	code, errorResponse := h.service.CreatePayroll(payload)
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

func (h *PayrollHandler) GetPayroll(f fiber.Ctx) error {

	var payload, result []models.Payroll
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

	filter := &dto.PayrollFilter{}

	if err := f.Bind().Query(filter); err != nil {
		errorResponse := &dto.Error{
			Message:    "Invalid Query Parameters",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}
	fmt.Println(filter)

	if role == "ADMIN" {
		result, code, count, errorResponse = h.service.GetPayrollAdmin(payload, *pageParams, *filter)
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
		result, code, count, errorResponse = h.service.GetPayroll(payload, *pageParams, *filter)
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

func (h *PayrollHandler) GetPayrollById(f fiber.Ctx) error {

	var payload models.Payroll

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

	result, code, errorResponse := h.service.GetPayrollById(payload, filters)
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

func (h *PayrollHandler) DeletePayroll(f fiber.Ctx) error {

	var payload models.Payroll

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

	code, errorResponse := h.service.DeletePayroll(payload, id)
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

func (h *PayrollHandler) PatchPayroll(f fiber.Ctx) error {

	var payload models.Payroll

	currentUserID, _ := f.Locals("user_id").(string)

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

	if err := f.Bind().Body(&payload); err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed converting json to struct",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(payload, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	result, code, errorResponse := h.service.PatchPayroll(payload, id)
	if errorResponse.Error != "" {
		audit(result, errorResponse)
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
