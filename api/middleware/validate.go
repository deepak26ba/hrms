package middleware

import (
	"hrms/common/dto"
	"hrms/internals/config"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT() fiber.Handler {

	return func(f fiber.Ctx) error {

		var jwtSecret = config.GetKey()

		authHeader := f.Get("Authorization")
		if authHeader == "" {
			errorResponse := &dto.Error{
				Message:    "Enter the Token",
				StatusCode: fiber.StatusForbidden,
				Error:      "Missing token",
			}
			return f.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		tokenString := strings.Split(authHeader, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			errorResponse := &dto.Error{
				Message:    "Enter valid Token",
				StatusCode: fiber.StatusForbidden,
				Error:      "Invalid token : " + err.Error(),
			}
			return f.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		claims := token.Claims.(jwt.MapClaims)
		role := f.Locals("role", claims["role"])
		userID := f.Locals("user_id", claims["user_id"])

		if role.(string) == "" && userID.(string) == "" {
			errorResponse := &dto.Error{
				Message:    "Access Deined",
				StatusCode: fiber.StatusForbidden,
				Error:      "Forbidden",
			}
			return f.Status(fiber.StatusForbidden).JSON(errorResponse)
		}

		return f.Next()
	}
}
