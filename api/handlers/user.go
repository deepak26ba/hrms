package handlers

import (
	"fmt"
	"hrms/api/service"
	"hrms/common/dto"
	"hrms/common/helper"
	"hrms/pkg/models"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	service service.UserService
	audit   service.AuditService
}

func NewUserHandler(service service.UserService, audit service.AuditService) *UserHandler {

	return &UserHandler{
		service: service,
		audit:   audit,
	}

}

func (h *UserHandler) GetUser(f fiber.Ctx) error {

	var payload []models.User

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
	filter := &dto.UserFilter{}

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
	pageParams = helper.FindOffset(pageParams.Page, pageParams.Limit)

	result, code, count, errorResp := h.service.GetUser(payload, *pageParams, *filter)

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

func (h *UserHandler) GetUserById(f fiber.Ctx) error {

	var usercredentials models.User

	idStr := f.Params("id")
	id, errorResponse := helper.GetId(idStr)
	if errorResponse.Error != "" {
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "View",
			Level:  "Handler",
			Data:   fmt.Sprint(id),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	userIDStr, ok := f.Locals("user_id").(string)
	if !ok {
		errorResponse := &dto.Error{
			Message:    "Need Authentication",
			StatusCode: fiber.StatusForbidden,
			Error:      "Forbidden",
		}
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "View",
			Level:  "Handler",
			Data:   fmt.Sprint(userIDStr),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

	userID, errorResponse := helper.StringToUUID(userIDStr)
	if errorResponse.Error != "" {
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "View",
			Level:  "Handler",
			Data:   fmt.Sprint(userID),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	roles, ok := f.Locals("roles").([]string)
	if !ok {
		errorResponse := &dto.Error{
			Message:    "Need Authentication",
			StatusCode: fiber.StatusForbidden,
			Error:      "Forbidden",
		}
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "View",
			Level:  "Handler",
			Data:   fmt.Sprint(roles),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

	role, ok := f.Locals("role").(string)
	if !ok {
		errorResponse := &dto.Error{
			Message:    "Need Authentication",
			StatusCode: fiber.StatusForbidden,
			Error:      "Forbidden",
		}
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "View",
			Level:  "Handler",
			Data:   fmt.Sprint(role),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

	credentials := dto.Authorize{
		Role:    role,
		Roles:   roles,
		QueryID: id,
		UserID:  userID}

	if helper.Authorize(credentials) {

		result, code, errorResponse := h.service.GetUserById(usercredentials, id)
		if errorResponse.Error != "" {
			userID, _ := f.Locals("user_id").(string)
			log := models.AuditTable{
				UserID: userID,
				Action: "View",
				Level:  "Handler",
				Data:   fmt.Sprint(result),
				Error:  fmt.Sprint(errorResponse),
			}
			h.audit.CreateAudit(log)
			return f.Status(code).JSON(errorResponse)
		}

		successResponse := &dto.Success{
			Message:    "Successfully Received",
			StatusCode: fiber.StatusOK,
			Data:       result,
		}

		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "View",
			Level:  "Handler",
			Data:   fmt.Sprint(successResponse),
			Error:  "nil",
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusOK).JSON(successResponse)
	}

	err := &dto.Error{
		Message:    "Need Permission",
		StatusCode: fiber.StatusForbidden,
		Error:      "Access Denied",
	}

	userIDLog, _ := f.Locals("user_id").(string)
	log := models.AuditTable{
		UserID: userIDLog,
		Action: "View",
		Level:  "Handler",
		Data:   fmt.Sprint(usercredentials),
		Error:  fmt.Sprint(err),
	}
	h.audit.CreateAudit(log)
	return f.Status(fiber.StatusForbidden).JSON(err)

}

func (h *UserHandler) DeleteUser(f fiber.Ctx) error {

	var UserRow models.User

	idStr := f.Params("id")
	id, errorResponse := helper.GetId(idStr)
	if errorResponse.Error != "" {
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "Delete",
			Level:  "Handler",
			Data:   fmt.Sprint(id),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	code, errorResponse := h.service.DeleteUser(UserRow, id)
	if errorResponse.Error != "" {
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "Delete",
			Level:  "Handler",
			Data:   fmt.Sprint(code),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(code).JSON(errorResponse)
	}

	successResponse := &dto.Success{
		Message:    "Successfully Deleted",
		StatusCode: fiber.StatusOK,
	}
	userID, _ := f.Locals("user_id").(string)
	log := models.AuditTable{
		UserID: userID,
		Action: "Delete",
		Level:  "Handler",
		Data:   fmt.Sprint(successResponse),
		Error:  "nil",
	}
	h.audit.CreateAudit(log)
	return f.Status(fiber.StatusOK).JSON(successResponse)

}

func (h *UserHandler) PatchUser(f fiber.Ctx) error {

	var usercredentials models.User

	idStr := f.Params("id")
	id, errorResponse := helper.GetId(idStr)
	if errorResponse.Error != "" {
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "Update",
			Level:  "Handler",
			Data:   fmt.Sprint(id),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if err := f.Bind().Body(&usercredentials); err != nil {
		errorResponse := &dto.Error{
			Message:    "Failed converting json to struct",
			StatusCode: fiber.StatusBadRequest,
			Error:      err.Error(),
		}
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "Update",
			Level:  "Handler",
			Data:   fmt.Sprint(userID),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	result, code, errorResponse := h.service.PatchUser(usercredentials, id)
	if errorResponse.Error != "" {
		userID, _ := f.Locals("user_id").(string)
		log := models.AuditTable{
			UserID: userID,
			Action: "Update",
			Level:  "Handler",
			Data:   fmt.Sprint(result),
			Error:  fmt.Sprint(errorResponse),
		}
		h.audit.CreateAudit(log)
		return f.Status(code).JSON(errorResponse)
	}

	successResponse := &dto.Success{
		Message:    "Successfully Updated",
		StatusCode: fiber.StatusOK,
		Data:       result,
	}

	userID, _ := f.Locals("user_id").(string)
	log := models.AuditTable{
		UserID: userID,
		Action: "Update",
		Level:  "Handler",
		Data:   fmt.Sprint(successResponse),
		Error:  "nil",
	}
	h.audit.CreateAudit(log)
	return f.Status(fiber.StatusOK).JSON(successResponse)
}
