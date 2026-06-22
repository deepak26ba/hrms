package routes

import (
	"hrms/api/handlers"
	"hrms/api/repository"
	"hrms/api/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func AuthRoutes(db *gorm.DB, f fiber.Router) {

	// Initialize repository
	authRepo := repository.NewAuthRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	auditService := service.NewAuditService(auditRepo)
	authService := service.NewAuthService(authRepo, auditRepo)

	authHandler := handlers.NewAuthHandler(authService, auditService)

	//auth using Email and Password
	register := f.Group("/auth")

	register.Post("/register", authHandler.Register)
	register.Post("/login", authHandler.Login)

	//auth using OAuth
	goauth := f.Group("/googleoauth")

	goauth.Get("/register", authHandler.GoogleRegister)
	goauth.Get("/login", authHandler.GoogleLogin)
	goauth.Get("/register/callback", authHandler.GoogleRegisterCallback)
	goauth.Get("/login/callback", authHandler.GoogleLoginCallback)

}
