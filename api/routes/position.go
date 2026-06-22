package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func PositionRoutes(db *gorm.DB, f fiber.Router) {

	positionRepo := repository.NewPositionRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	positionService := service.NewPositionService(positionRepo)
	positionHandler := handlers.NewPositionHandler(positionService, auditService)

	position := f.Group("/position")

	position.Post("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), positionHandler.CreatePosition)
	position.Get("/", middleware.ValidateJWT(), positionHandler.GetPosition)
	position.Get("/:id", middleware.ValidateJWT(), positionHandler.GetPositionById)
	position.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), positionHandler.PatchPosition)
	position.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), positionHandler.DeletePosition)
}
