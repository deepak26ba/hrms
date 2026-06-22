package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func EmployeeRoutes(db *gorm.DB, f fiber.Router) {

	employeeRepo := repository.NewEmployeeRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	employeeService := service.NewEmployeeService(employeeRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService, auditService)

	employee := f.Group("/Employee")

	employee.Post("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), employeeHandler.CreateEmployee)
	employee.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), employeeHandler.GetEmployee)
	employee.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), employeeHandler.GetEmployeeById)
	employee.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), employeeHandler.PatchEmployee)
	employee.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP"), employeeHandler.DeleteEmployee)
}
