package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type ConnectionString struct {
	Host     string
	User     string
	DBName   string
	Password string
	SslMode  string
	Port     string
}

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" validate:"required,email" gorm:"size:100;not null;unique"`
	Password  string         `json:"password" validate:"required,min=8" gorm:"size:255"`
	RoleID    uuid.UUID      `json:"role_id" gorm:"not null"`
	Login     bool           `json:"login"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Provider  string         `json:"provider" gorm:"size:50;not null;default:'local'"`
	OAuthID   string         `json:"oauth_id" gorm:"uniqueIndex;default:null"`

	Role Roles `json:"role"`
}

type Roles struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Department struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null;unique"`
	Description string         `json:"description" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Position struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null;unique"`
	Description string         `json:"description" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Employee struct {
	ID           uuid.UUID      `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"size:50;not null;unique"`
	Age          uint           `json:"age"`
	PhoneNumber  string         `json:"phone_number" gorm:"size:15;not null;unique"`
	DepartmentId uuid.UUID      `json:"department_id" gorm:"not null"`
	PositionId   uuid.UUID      `json:"position_id" gorm:"not null"`
	UserId       uuid.UUID      `json:"user_id" gorm:"not null"`
	IsActive     bool           `json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	Position   Position   `json:"position"`
	Department Department `json:"department"`
	User       User       `json:"user"`
}

type EmployeeBoardingStatus struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	UserId      uuid.UUID      `json:"user_id" gorm:"not null"`
	Onboarding  time.Time      `json:"on_boarding"`
	Offboarding time.Time      `json:"off_boarding"`
	Status      bool           `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"user"`
}

type Payroll struct {
	ID         uuid.UUID      `json:"id" gorm:"primaryKey"`
	UserId     uuid.UUID      `json:"user_id" gorm:"not null"`
	PositionId uuid.UUID      `json:"position_id" gorm:"not null"`
	BasicPay   int            `json:"basic_pay"`
	Allowance  int            `json:"allowance"`
	Salary     int            `json:"salary"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	User     User     `json:"user"`
	Position Position `json:"position"`
}

type LeaveRequest struct {
	ID         uuid.UUID      `json:"id" gorm:"primaryKey"`
	UserId     uuid.UUID      `json:"user_id" gorm:"not null"`
	EmployeeId uuid.UUID      `json:"employee_id" gorm:"not null"`
	ReportedTo uuid.UUID      `json:"reported_to" gorm:"not null"`
	Type       string         `json:"type"`
	From       time.Time      `json:"from"`
	To         time.Time      `json:"to"`
	IsApproved *bool          `json:"is_approved"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	User     User     `json:"user"`
	Employee Employee `json:"employee" gorm:"foreignKey:EmployeeId"`
	Employer Employee `json:"employer" gorm:"foreignKey:ReportedTo"`
}

type LeaveBalance struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey"`
	UserId         uuid.UUID `json:"user_id" gorm:"not null;unique"`
	MaxLeave       uint      `json:"max_leave" gorm:"default:10"`
	LeaveTaken     uint      `json:"leave_taken"`
	LeaveRemaining uint      `json:"leave_remaining" gorm:"default:10"`
	LeaveHolding   uint      `json:"leave_holding"`
	SickLeave      uint      `json:"sick_leave" gorm:"default:5"`
	CausalLeave    uint      `json:"causal_leave" gorm:"default:5"`
	LossOfPay      uint      `json:"loss_of_pay" gorm:"default:0"`

	User User `json:"user"`
}

type Probation struct {
	ID                 uuid.UUID      `json:"id" gorm:"primaryKey"`
	EmployeeId         uuid.UUID      `json:"employee_id" gorm:"not null"`
	UserId             uuid.UUID      `json:"user_id" gorm:"not null"`
	ProbationStatus    string         `json:"probation_status"`
	ProbationStartDate time.Time      `json:"probation_start_date"`
	ProbationEndDate   time.Time      `json:"probation_end_date"`
	ReviewedBy         uuid.UUID      `json:"reviewed_by"`
	TaskCompleted      bool           `json:"task_completed"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`

	User     User     `json:"user"`
	Employee Employee `json:"employee"`
	Employer Employee `json:"employer" gorm:"foreignKey:ReviewedBy"`
}

type Attendance struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey"`
	UserId    uuid.UUID      `json:"user_id" gorm:"not null"`
	TimeIn    time.Time      `json:"time_in"`
	TimeOut   time.Time      `json:"time_out"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"user"`
}

type AttendanceAudit struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey"`
	AttendanceID uuid.UUID `json:"attendance_id"`
	UserId       uuid.UUID `json:"user_id" gorm:"not null"`
	TimeIn       time.Time `json:"time_in"`
	TimeOut      time.Time `json:"time_out"`

	Attendance Attendance `json:"attendance"`
}

type AuditTable struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Action    string    `json:"action" gorm:"not null"`
	Level     string    `json:"level" gorm:"not null"`
	Data      string    `json:"data" gorm:"not null"`
	Error     string    `json:"error" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *AuditTable) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (a *AttendanceAudit) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (l *LeaveRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (l *LeaveBalance) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (p *Payroll) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (p *Probation) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (e *EmployeeBoardingStatus) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (e *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (d *Department) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (p *Position) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

func (r *Roles) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}
