package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func AttendanceAuditRoutes(db *gorm.DB, f fiber.Router) {

	attendanceAuditRepo := repository.NewAttendanceAuditRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	attendanceAuditService := service.NewAttendanceAuditService(attendanceAuditRepo)
	attendanceAuditHandler := handlers.NewAttendanceAuditHandler(attendanceAuditService, auditService)

	attendanceAudit := f.Group("/AttendanceAudit")

	attendanceAudit.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), attendanceAuditHandler.GetAttendanceAudit)
	attendanceAudit.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), attendanceAuditHandler.GetAttendanceAuditById)
}
