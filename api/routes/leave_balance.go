package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func LeaveBalanceRoutes(db *gorm.DB, f fiber.Router) {

	leaveBalanceRepo := repository.NewLeaveBalanceRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	leaveBalanceService := service.NewLeaveBalanceService(leaveBalanceRepo)
	leaveBalanceHandler := handlers.NewLeaveBalanceHandler(leaveBalanceService, auditService)

	leaveBalance := f.Group("/LeaveBalance")

	leaveBalance.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), leaveBalanceHandler.GetLeaveBalance)
	leaveBalance.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), leaveBalanceHandler.GetLeaveBalanceById)
}
