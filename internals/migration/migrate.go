package migration

import (
	"fmt"
	"hrms/pkg/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {

	err := db.AutoMigrate(
		&models.Roles{}, &models.User{}, &models.Department{},
		&models.Position{}, &models.Employee{}, &models.LeaveBalance{},
		&models.EmployeeBoardingStatus{}, &models.LeaveRequest{},
		&models.Probation{}, &models.Payroll{}, &models.Attendance{},
		&models.AttendanceAudit{}, &models.AuditTable{})
	if err != nil {
		return fmt.Errorf("Failed to migrate : %v", err)
	}

	return nil
}
