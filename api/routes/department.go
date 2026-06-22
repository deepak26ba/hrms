package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func DepartmentRoutes(db *gorm.DB, f fiber.Router) {

	departmentRepo := repository.NewDepartmentRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	departmentService := service.NewDepartmentService(departmentRepo)
	departmentHandler := handlers.NewDepartmentHandler(departmentService, auditService)

	department := f.Group("/department")

	department.Post("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP"), departmentHandler.CreateDepartment)
	department.Get("/", middleware.ValidateJWT(), departmentHandler.GetDepartment)
	department.Get("/:id", middleware.ValidateJWT(), departmentHandler.GetDepartmentById)
	department.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), departmentHandler.PatchDepartment)
	department.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP"), departmentHandler.DeleteDepartment)
}
