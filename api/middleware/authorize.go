package middleware

import (
	"hrms/common/dto"
	"slices"

	"github.com/gofiber/fiber/v3"
)

func Authorize(rolesAllowed ...string) fiber.Handler {

	return func(f fiber.Ctx) error {

		role, ok := f.Locals("role").(string)
		if !ok {
			errorResponse := &dto.Error{
				Message:    "Need Authentication",
				StatusCode: fiber.StatusForbidden,
				Error:      "Forbidden",
			}
			return f.Status(fiber.StatusForbidden).JSON(errorResponse)
		}
		f.Locals("roles", rolesAllowed)

		if role == "" {
			errorResponse := &dto.Error{
				Message:    "Access Deined",
				StatusCode: fiber.StatusForbidden,
				Error:      "Forbidden",
			}
			return f.Status(fiber.StatusForbidden).JSON(errorResponse)
		}

		if slices.Contains(rolesAllowed, role) {
			return f.Next()
		}

		errorResponse := &dto.Error{
			Message:    "Access Deined",
			StatusCode: fiber.StatusForbidden,
			Error:      "Forbidden",
		}
		return f.Status(fiber.StatusForbidden).JSON(errorResponse)
	}

}
