package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func AttendanceRoutes(db *gorm.DB, f fiber.Router) {

	attendanceRepo := repository.NewAttendanceRepository(db)
	attendanceAuditRepo := repository.NewAttendanceAuditRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo)
	attendanceAuditService := service.NewAttendanceAuditService(attendanceAuditRepo)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService, auditService, attendanceAuditService)

	attendance := f.Group("/Attendance")

	attendance.Post("/in", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), attendanceHandler.CreateAttendance)
	attendance.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), attendanceHandler.GetAttendance)
	attendance.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), attendanceHandler.GetAttendanceById)
	attendance.Patch("/out/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), attendanceHandler.PatchAttendance)
	attendance.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), attendanceHandler.DeleteAttendance)
}
