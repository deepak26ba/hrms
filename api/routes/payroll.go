package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func PayrollRoutes(db *gorm.DB, f fiber.Router) {

	payrollRepo := repository.NewPayrollRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	payrollService := service.NewPayrollService(payrollRepo)
	payrollHandler := handlers.NewPayrollHandler(payrollService, auditService)

	payroll := f.Group("/Payroll")

	payroll.Post("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP"), payrollHandler.CreatePayroll)
	payroll.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), payrollHandler.GetPayroll)
	payroll.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), payrollHandler.GetPayrollById)
	payroll.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), payrollHandler.PatchPayroll)
	payroll.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP"), payrollHandler.DeletePayroll)
}
