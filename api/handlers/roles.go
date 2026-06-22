package handlers

import (
	"fmt"
	"hrms/api/service"
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"

	"github.com/gofiber/fiber/v3"
)

type RoleHandler struct {
	service service.RolesService
	audit   service.AuditService
}

func NewRoleHandler(service service.RolesService, audit service.AuditService) *RoleHandler {

	return &RoleHandler{
		service: service,
		audit:   audit,
	}

}

func (h *RoleHandler) CreateRoles(f fiber.Ctx) error {

	var payload models.Roles

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

	code, errorResponse := h.service.CreateRoles(payload)
	if errorResponse.Error != "" {
		audit(code, errorResponse)
		return f.Status(code).JSON(errorResponse)

	}

	successResponse := &dto.Success{
		Message:    "Successfully Created",
		StatusCode: fiber.StatusCreated,
	}
	audit(successResponse, "nil")

	return f.Status(fiber.StatusCreated).JSON(successResponse)

}

func (h *RoleHandler) GetRoles(f fiber.Ctx) error {

	var payload []models.Roles

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

	pageParams := &dto.Pagination{Page: 1, Limit: 5}
	filter := &dto.RoleFilter{}

	if err := f.Bind().Query(pageParams); err != nil {
		errorResponse := &dto.Error{
			Message:    "Invalid Query Parameters",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if err := f.Bind().Query(filter); err != nil {
		errorResponse := &dto.Error{
			Message:    "Invalid Filter Parameters",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}
	pageParams= helper.FindOffset(pageParams.Page, pageParams.Limit)

	result, code, count, errorResp := h.service.GetRoles(payload, *pageParams, *filter)

	if errorResp.Error != "" {
		audit(nil, errorResp.Error)
		return f.Status(code).JSON(errorResp)
	}

	successResponse := dto.Success{
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

	audit(successResponse, nil)

	return f.Status(fiber.StatusOK).JSON(successResponse)
}

func (h *RoleHandler) GetRolesById(f fiber.Ctx) error {

	var payload models.Roles

	currentUserID, _ := f.Locals("user_id").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "ViewByID",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	idStr := f.Params("id")
	id, errorResponse := helper.GetId(idStr)
	if errorResponse.Error != "" {
		audit(id, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	result, code, errorResponse := h.service.GetRolesById(payload, id)
	if errorResponse.Error != "" {
		audit(result, errorResponse)
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

func (h *RoleHandler) DeleteRoles(f fiber.Ctx) error {

	var usercredentials models.Roles

	currentUserID, _ := f.Locals("user_id").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "ViewByID",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	idStr := f.Params("id")
	id, errorResponse := helper.GetId(idStr)
	if errorResponse.Error != "" {
		audit(id, errorResponse)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	code, errorResponse := h.service.DeleteRoles(usercredentials, id)
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

func (h *RoleHandler) PatchRoles(f fiber.Ctx) error {

	var payload models.Roles

	currentUserID, _ := f.Locals("user_id").(string)

	audit := func(data any, err any) {
		h.audit.CreateAudit(models.AuditTable{
			UserID: currentUserID,
			Action: "ViewByID",
			Level:  "Handler",
			Data:   fmt.Sprint(data),
			Error:  fmt.Sprint(err),
		})
	}

	idStr := f.Params("id")
	id, errorResponse := helper.GetId(idStr)
	if errorResponse.Error != "" {
		audit(id, errorResponse)
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

	result, code, errorResponse := h.service.PatchRoles(payload, id)
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
