package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func UserRoutes(db *gorm.DB, f fiber.Router) {

	userRepo := repository.NewUserRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, auditService)

	user := f.Group("/user")

	user.Get("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), userHandler.GetUser)
	user.Get("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR", "EMPLOYEE"), userHandler.GetUserById)
	user.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), userHandler.PatchUser)
	user.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN", "VP", "HR"), userHandler.DeleteUser)
}
