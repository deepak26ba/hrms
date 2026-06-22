package dto

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type ClaimsJWT struct {
	Role   string    `json:"role"`
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type Success struct {
	Message    string           `json:"message"`
	StatusCode int              `json:"status_code"`
	Data       any              `json:"data,omitempty"`
	Pagination PaginationResult `json:"pagination,omitzero"`
}

type PaginationResult struct {
	Page        int  `json:"page,omitempty"`
	Limit       int  `json:"limit,omitempty"`
	HasNextPage bool `json:"has_next_page"`
	TotalCount  int  `json:"total_count,omitempty"`
}

type Login struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Data       string `json:"data"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type EmployeeRequest struct {
	Name         string    `json:"name" gorm:"size:50;not null;unique"`
	Age          uint      `json:"age" validate:"required,numeric,min=18,max=50"`
	PhoneNumber  string    `json:"phone_number" validate:"required,numeric,len=10" gorm:"size:15;not null;unique"`
	DepartmentId uuid.UUID `json:"department_id" gorm:"not null"`
	PositionId   uuid.UUID `json:"position_id" gorm:"not null"`
	UserId       uuid.UUID `json:"user_id" gorm:"not null"`
}

type EmployeeBoardingStatusRequest struct {
	UserId      uuid.UUID `json:"user_id" gorm:"not null"`
	Onboarding  time.Time `json:"on_boarding"`
	Offboarding time.Time `json:"off_boarding"`
	Status      bool      `json:"status"`
}

type PayrollRequest struct {
	UserId     uuid.UUID `json:"user_id" gorm:"not null"`
	PositionId uuid.UUID `json:"position_id" gorm:"not null"`
	BasicPay   int       `json:"basic_pay"`
	Allowance  int       `json:"allowance"`
	Salary     int       `json:"salary"`
}

type LeavePermissionRequest struct {
	UserId     uuid.UUID `json:"user_id" gorm:"not null"`
	EmployeeId uuid.UUID `json:"employee_id" gorm:"not null"`
	ReportedTo uuid.UUID `json:"reported_to" gorm:"not null"`
	Type       string    `json:"type"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
}

type LeaveRequest struct {
	LeaveRequestId uuid.UUID `json:"leave_request_id" gorm:"not null"`
	UserId         uuid.UUID `json:"user_id" gorm:"not null"`
	EmployeeId     uuid.UUID `json:"employee_id" gorm:"not null"`
	MaxLeave       int       `json:"max_leave"`
	LeaveTaken     int       `json:"leave_taken"`
	LeaveRemaining int       `json:"leave_remaining"`
}

type AttendanceAuditRequest struct {
	AttendanceID uuid.UUID `json:"attendance_id"`
}

type Authorize struct {
	Role    string
	Roles   []string
	QueryID uuid.UUID
	UserID  uuid.UUID
}

type GetByID struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Role   string
}

type Pagination struct {
	Page   int       `json:"page"`
	Offset int       `json:"offset"`
	Limit  int       `json:"limit"`
	UserID uuid.UUID `json:"user_id"`
}

type EmployeeBoardingStatusFilter struct {
	UserId      uuid.UUID `json:"user_id"`
	Onboarding  time.Time `json:"on_boarding"`
	Offboarding time.Time `json:"off_boarding"`
	Status      *bool     `json:"status"`
}

type EmployeeFilter struct {
	Name         string    `json:"name" `
	Age          uint      `json:"age"`
	PhoneNumber  string    `json:"phone_number" `
	DepartmentId uuid.UUID `json:"department_id" `
	PositionId   uuid.UUID `json:"position_id" `
	IsActive     *bool     `json:"is_active"`
}

type RoleFilter struct {
	Name string `json:"name" `
}

type DepartmentFilter struct {
	Name string `json:"name" `
}
type PositionFilter struct {
	Name string `json:"name" `
}

type PayrollFilter struct {
	UserId     uuid.UUID `json:"user_id" `
	PositionId uuid.UUID `json:"position_id" `
	BasicPay   int       `json:"basic_pay"`
	Allowance  int       `json:"allowance"`
	Salary     int       `json:"salary"`
}

type ProbationFilter struct {
	EmployeeId         uuid.UUID `json:"employee_id"`
	UserId             uuid.UUID `json:"user_id"`
	ProbationStatus    string    `json:"probation_status"`
	ProbationStartDate time.Time `json:"probation_start_date"`
	ProbationEndDate   time.Time `json:"probation_end_date"`
	ReviewedBy         uuid.UUID `json:"reviewed_by"`
	TaskCompleted      *bool     `json:"task_completed"`
}

type LeaveFilter struct {
	LeaveRequestId uuid.UUID `json:"leave_request_id" `
	UserId         uuid.UUID `json:"user_id"`
	EmployeeId     uuid.UUID `json:"employee_id"`
}

type LeaveRequestFilter struct {
	UserId     uuid.UUID `json:"user_id"`
	EmployeeId uuid.UUID `json:"employee_id"`
	ReportedTo uuid.UUID `json:"reported_to"`
	Type       string    `json:"type"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
}

type AttendanceFilter struct {
	EmployeeId uuid.UUID `json:"employee_id"`
	UserId     uuid.UUID `json:"user_id"`
	TimeIn     time.Time `json:"time_in"`
	TimeOut    time.Time `json:"time_out"`
}

type AttendanceAuditFilter struct {
	AttendanceID uuid.UUID `json:"attendance_id"`
	UserId       uuid.UUID `json:"user_id"`
	TimeIn       time.Time `json:"time_in"`
	TimeOut      time.Time `json:"time_out"`
}

type LeaveBalanceFilter struct {
	UserId         uuid.UUID `json:"user_id"`
	MaxLeave       uint      `json:"max_leave"`
	LeaveTaken     uint      `json:"leave_taken"`
	LeaveRemaining uint      `json:"leave_remaining"`
	SickLeave      uint      `json:"sick_leave"`
	CausalLeave    uint      `json:"causal_leave"`
	LossOfPay      uint      `json:"loss_of_pay"`
}

type UserFilter struct {
	Email    string    `json:"email"`
	RoleID   uuid.UUID `json:"role_id"`
	IsActive *bool     `json:"is_active"`
}

type OAuthRequest struct {
	Email    string    `json:"email" validate:"required,email"`
	Name     string    `json:"name" `
	Provider string    `json:"provider" `
	UserID   string    `json:"user_id" `
	RoleID   uuid.UUID `json:"role_id" `
}
