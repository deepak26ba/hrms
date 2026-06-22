package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func ProbationRoutes(db *gorm.DB, f fiber.Router) {

	probationRepo := repository.NewProbationRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	probationService := service.NewProbationService(probationRepo)
	probationHandler := handlers.NewProbationHandler(probationService, auditService)

	probation := f.Group("/Probation")

	probation.Post("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP"), probationHandler.CreateProbation)
	probation.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), probationHandler.GetProbation)
	probation.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), probationHandler.GetProbationById)
	probation.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), probationHandler.PatchProbation)
	probation.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP"), probationHandler.DeleteProbation)
}
