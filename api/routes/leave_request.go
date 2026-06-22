package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func LeaveRequestRoutes(db *gorm.DB, f fiber.Router) {

	leaveRequestRepo := repository.NewLeaveRequestRepository(db)
	leaveBalanceRepo := repository.NewLeaveBalanceRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	leaveBalanceService := service.NewLeaveBalanceService(leaveBalanceRepo)
	leaveRequestService := service.NewLeaveRequestService(leaveRequestRepo, leaveBalanceRepo)
	leaveRequestHandler := handlers.NewLeaveRequestHandler(leaveRequestService, auditService, leaveBalanceService)

	leaveRequest := f.Group("/LeaveRequest")

	leaveRequest.Post("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), leaveRequestHandler.CreateLeaveRequest)
	leaveRequest.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), leaveRequestHandler.GetLeaveRequest)
	leaveRequest.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), leaveRequestHandler.GetLeaveRequestById)
	leaveRequest.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), leaveRequestHandler.PatchLeaveRequest)
}
