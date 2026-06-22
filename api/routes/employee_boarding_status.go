package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func EmployeeBoardingStatusRoutes(db *gorm.DB, f fiber.Router) {

	employeeBoardingStatusRepo := repository.NewEmployeeBoardingStatusRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	employeeBoardingStatusService := service.NewEmployeeBoardingStatusService(employeeBoardingStatusRepo)
	employeeBoardingStatusHandler := handlers.NewEmployeeBoardingStatusHandler(employeeBoardingStatusService, auditService)

	employeeBoardingStatus := f.Group("/EmployeeBoardingStatus")

	employeeBoardingStatus.Post("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), employeeBoardingStatusHandler.CreateEmployeeBoardingStatus)
	employeeBoardingStatus.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), employeeBoardingStatusHandler.GetEmployeeBoardingStatus)
	employeeBoardingStatus.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), employeeBoardingStatusHandler.GetEmployeeBoardingStatusById)
	employeeBoardingStatus.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), employeeBoardingStatusHandler.PatchEmployeeBoardingStatus)
	employeeBoardingStatus.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), employeeBoardingStatusHandler.DeleteEmployeeBoardingStatus)
}
