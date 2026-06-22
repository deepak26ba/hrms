package routes

import (
	"hrms/api/handlers"
	"hrms/api/middleware"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RoleRoutes(db *gorm.DB, f fiber.Router) {

	roleRepo := repository.NewRoleRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditService := service.NewAuditService(auditRepo)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := handlers.NewRoleHandler(roleService, auditService)

	role := f.Group("/role")

	role.Post("/", middleware.ValidateJWT(), middleware.Authorize("ADMIN"), roleHandler.CreateRoles)
	role.Get("/", middleware.ValidateJWT(), roleHandler.GetRoles)
	role.Get("/:id", middleware.ValidateJWT(), roleHandler.GetRolesById)
	role.Patch("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN"), roleHandler.PatchRoles)
	role.Delete("/:id", middleware.ValidateJWT(), middleware.Authorize("ADMIN"), roleHandler.DeleteRoles)
}
