package routes

import (
	"fmt"
	"hrms/common/dto"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, f *fiber.App) {

	api := f.Group("/api/v1")

	AuthRoutes(db, api)
	UserRoutes(db, api)
	RoleRoutes(db, api)
	DepartmentRoutes(db, api)
	PositionRoutes(db, api)
	EmployeeRoutes(db, api)
	PayrollRoutes(db, api)
	EmployeeBoardingStatusRoutes(db, api)
	LeaveRequestRoutes(db, api)
	ProbationRoutes(db, api)
	AttendanceRoutes(db, api)
	AttendanceAuditRoutes(db, api)
	LeaveBalanceRoutes(db, api)

	f.Use(func(c fiber.Ctx) error {
		errorResponse := &dto.Error{
			Message:    "Use valid path",
			StatusCode: http.StatusNotFound,
			Error:      fmt.Errorf("Invaild Path").Error(),
		}
		return c.Status(http.StatusNotFound).JSON(errorResponse)
	})

}
