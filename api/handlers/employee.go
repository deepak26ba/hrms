package handlers

import (
	"fmt"
	"hrms/api/service"
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofrs/uuid"
)

type EmployeeHandler struct {
	service service.EmployeeService
	audit   service.AuditService
}

func NewEmployeeHandler(service service.EmployeeService, audit service.AuditService) *EmployeeHandler {

	return &EmployeeHandler{
		service: service,
		audit:   audit,
	}

}

func (h *EmployeeHandler) CreateEmployee(f fiber.Ctx) error {

	var userCreadentails dto.EmployeeRequest

	if err := f.Bind().Body(&userCreadentails); err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed converting json to struct",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "Create",
			Level:  "Handler",
			Data:   fmt.Sprint(userCreadentails),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)

	}

	validate := validator.New()
	err := validate.Struct(userCreadentails)
	if err != nil {
		errorResponse := &dto.Error{
			Message:    "Validation failed",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "Create",
			Level:  "Handler",
			Data:   fmt.Sprint(userCreadentails),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	credentials := models.Employee{
		Name:         userCreadentails.Name,
		Age:          userCreadentails.Age,
		PhoneNumber:  userCreadentails.PhoneNumber,
		UserId:       userCreadentails.UserId,
		DepartmentId: userCreadentails.DepartmentId,
		PositionId:   userCreadentails.PositionId,
	}

	code, errorResponse := h.service.CreateEmployee(credentials)
	if errorResponse.Error != "" {
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "Create",
			Level:  "Handler",
			Data:   fmt.Sprint(code),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(code).JSON(errorResponse)

	}

	successResponse := &dto.Success{
		Message:    "Successfully Created",
		StatusCode: fiber.StatusCreated,
	}

	userID, _ := f.Locals("user_id").(string)
	log := models.AuditTable{
		UserID: userID,
		Action: "Create",
		Level:  "Handler",
		Data:   fmt.Sprint(successResponse),
		Error:  "nil",
	}
	h.audit.CreateAudit(log)

	return f.Status(fiber.StatusCreated).JSON(successResponse)

}

func (h *EmployeeHandler) GetEmployee(f fiber.Ctx) error {

	var payload []models.Employee

	currentUserID, _ := f.Locals("user_id").(string)

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
		audit(pageParams, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}
	pageParams = helper.FindOffset(pageParams.Page, pageParams.Limit)

	filter := &dto.EmployeeFilter{}

	if err := f.Bind().Query(filter); err != nil {
		errorResponse := &dto.Error{
			Message:    "Invalid Query Parameters",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		audit(filter, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}
	fmt.Println(filter)

	result, code, count, errorResponse := h.service.GetEmployee(payload, *pageParams, *filter)

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

func (h *EmployeeHandler) DeleteEmployee(f fiber.Ctx) error {

	var usercredentials models.Employee

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

	idStr := f.Params("id")
	id, errorResponse := helper.GetId(idStr)
	if errorResponse.Error != "" {
		audit(id, "Invalid ID format")
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	code, errorResponse := h.service.DeleteEmployee(usercredentials, id)
	if errorResponse.Error != "" {
		audit(code, "Invalid ID format")
		return f.Status(code).JSON(errorResponse)
	}

	successResponse := &dto.Success{
		Message:    "Successfully Deleted",
		StatusCode: fiber.StatusOK,
	}
	audit(successResponse, "nil")

	return f.Status(fiber.StatusOK).JSON(successResponse)

}

func (h *EmployeeHandler) PatchEmployee(f fiber.Ctx) error {

	var payload models.Employee

	currentUserID, ok := f.Locals("user_id").(string)
	if !ok {
		errorResponse := &dto.Error{
			Message:    "Need Authentication",
			StatusCode: fiber.StatusForbidden,
			Error:      "Forbidden",
		}
		log := models.AuditTable{
			UserID: currentUserID,
			Action: "Update",
			Level:  "Handler",
			Data:   fmt.Sprint(currentUserID),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
	}
	roles, _ := f.Locals("roles").([]string)
	role, _ := f.Locals("role").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "Update",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	if currentUserID == "" || role == "" {
		errorResponse := dto.Error{
			Message:    "Need Authentication",
			StatusCode: fiber.StatusForbidden,
			Error:      "Forbidden"}
		audit("AuthCheck", errorResponse)
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

	if err := f.Bind().Body(&payload); err != nil {
		errorResponse := dto.Error{
			Message:    "Invalid JSON format",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error()}
		audit(payload, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	id, errorResponse := helper.GetId(f.Params("id"))
	userID, errUUID := helper.StringToUUID(currentUserID)
	if errorResponse.Error != "" || errUUID.Error != "" {
		audit(id, "Invalid ID format")
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	credentials := dto.Authorize{
		Role:    role,
		Roles:   roles,
		QueryID: id,
		UserID:  userID}

	if !helper.Authorize(credentials) {
		errorResponse := dto.Error{
			Message:    "Need Permission",
			StatusCode: fiber.StatusForbidden,
			Error:      "Access Denied"}
		audit(id, errorResponse)
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

	filter := dto.GetByID{
		ID:     id,
		UserID: userID,
		Role:   role,
	}

	result, code, errorResponse := h.service.PatchEmployee(payload, filter)
	if errorResponse.Error != "" {
		audit(result, errorResponse)
		return f.Status(code).JSON(errorResponse)
	}

	success := dto.Success{
		Message:    "Successfully Updated",
		StatusCode: fiber.StatusOK,
		Data:       result,
	}
	audit(success, "nil")
	return f.Status(fiber.StatusOK).JSON(success)
}

func (h *EmployeeHandler) GetEmployeeById(f fiber.Ctx) error {

	var payload models.Employee

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
		errorResponse := dto.Error{
			Message:    "Forbidden",
			StatusCode: fiber.StatusForbidden,
			Error:      "Missing Auth",
		}
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

	result, code, errorResponse := h.service.GetEmployeeById(payload, filters)
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
