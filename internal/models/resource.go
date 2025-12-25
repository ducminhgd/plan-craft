package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Resource represents a team member or resource in a project
type Resource struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Basic Information
	Name  string `gorm:"type:varchar(255);not null" json:"name"`
	Email string `gorm:"type:varchar(255);index" json:"email,omitempty"`

	// Role Information (can be used globally across projects)
	Role     string `gorm:"type:varchar(100);not null;index" json:"role"` // e.g., "Backend Developer", "QA", "Designer"
	IsActive bool   `gorm:"default:true;index" json:"is_active"`          // Whether resource is currently available

	// Capacity (default capacity, can be overridden per project)
	DefaultHoursPerDay   float64 `gorm:"type:decimal(5,2);default:8" json:"default_hours_per_day"`     // Default working hours per day
	DefaultDaysPerWeek   float64 `gorm:"type:decimal(5,2);default:5" json:"default_days_per_week"`     // Default working days per week
	DefaultDaysPerMonth  float64 `gorm:"type:decimal(5,2);default:20" json:"default_days_per_month"`   // Default working days per month

	// Cost Information (default rates, can be overridden per project)
	DefaultHourlyRate  float64  `gorm:"type:decimal(10,2);default:0" json:"default_hourly_rate"`
	DefaultDailyRate   float64  `gorm:"type:decimal(10,2);default:0" json:"default_daily_rate"`
	DefaultMonthlyRate float64  `gorm:"type:decimal(10,2);default:0" json:"default_monthly_rate"`
	Currency           string   `gorm:"type:varchar(3);default:'USD'" json:"currency"` // ISO 4217 currency code

	// Metadata
	Skills StringArray `gorm:"type:json" json:"skills,omitempty"` // List of skills
	Notes  string      `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	ProjectRoles []ProjectRole    `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE" json:"project_roles,omitempty"`
	Assignments  []TaskAssignment `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE" json:"assignments,omitempty"`
}

// TableName specifies the table name for Resource model
func (Resource) TableName() string {
	return "resources"
}

// BeforeSave is a GORM hook that runs before saving
func (r *Resource) BeforeSave(tx *gorm.DB) error {
	// Validate capacity
	if r.DefaultHoursPerDay < 0 || r.DefaultHoursPerDay > 24 {
		return errors.New("default hours per day must be between 0 and 24")
	}

	if r.DefaultDaysPerWeek < 0 || r.DefaultDaysPerWeek > 7 {
		return errors.New("default days per week must be between 0 and 7")
	}

	if r.DefaultDaysPerMonth < 0 || r.DefaultDaysPerMonth > 31 {
		return errors.New("default days per month must be between 0 and 31")
	}

	// Validate rates are non-negative
	if r.DefaultHourlyRate < 0 {
		return errors.New("default hourly rate cannot be negative")
	}

	if r.DefaultDailyRate < 0 {
		return errors.New("default daily rate cannot be negative")
	}

	if r.DefaultMonthlyRate < 0 {
		return errors.New("default monthly rate cannot be negative")
	}

	return nil
}

// MonthlyCapacityHours returns the monthly capacity in hours
func (r *Resource) MonthlyCapacityHours() float64 {
	return r.DefaultDaysPerMonth * r.DefaultHoursPerDay
}

// ProjectRole represents a resource's role assignment in a specific project
// This allows the same resource to have different roles, levels, and rates in different projects
type ProjectRole struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Project and Resource
	ProjectID  uint `gorm:"not null;index:idx_project_role_project" json:"project_id"`
	ResourceID uint `gorm:"not null;index:idx_project_role_resource" json:"resource_id"`

	// Role Information in this Project
	Role  string `gorm:"type:varchar(100);not null" json:"role"`         // Role in this project (can differ from Resource.Role)
	Level string `gorm:"type:varchar(50);not null" json:"level"`         // e.g., "Junior", "Mid", "Senior", "Lead"

	// Capacity in this Project (overrides Resource defaults)
	HoursPerDay  *float64 `gorm:"type:decimal(5,2)" json:"hours_per_day,omitempty"`   // If null, uses Resource.DefaultHoursPerDay
	DaysPerWeek  *float64 `gorm:"type:decimal(5,2)" json:"days_per_week,omitempty"`   // If null, uses Resource.DefaultDaysPerWeek
	DaysPerMonth *float64 `gorm:"type:decimal(5,2)" json:"days_per_month,omitempty"`  // If null, uses Resource.DefaultDaysPerMonth

	// Estimated Allocation for the entire project
	EstimatedManMonths float64 `gorm:"type:decimal(10,2);default:0" json:"estimated_man_months"` // Total man-months for this role in project

	// Cost Rates in this Project (overrides Resource defaults)
	HourlyRate  *float64 `gorm:"type:decimal(10,2)" json:"hourly_rate,omitempty"`
	DailyRate   *float64 `gorm:"type:decimal(10,2)" json:"daily_rate,omitempty"`
	MonthlyRate *float64 `gorm:"type:decimal(10,2)" json:"monthly_rate,omitempty"`

	// Timeline in Project
	StartDate *time.Time `gorm:"type:datetime" json:"start_date,omitempty"`
	EndDate   *time.Time `gorm:"type:datetime" json:"end_date,omitempty"`

	// Status
	IsActive bool   `gorm:"default:true" json:"is_active"` // Whether this role assignment is currently active
	Notes    string `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Project     *Project         `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Resource    *Resource        `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
	Assignments []TaskAssignment `gorm:"foreignKey:ProjectRoleID;constraint:OnDelete:CASCADE" json:"assignments,omitempty"`
}

// TableName specifies the table name for ProjectRole model
func (ProjectRole) TableName() string {
	return "project_roles"
}

// BeforeSave is a GORM hook that runs before saving
func (pr *ProjectRole) BeforeSave(tx *gorm.DB) error {
	// Validate capacity if set
	if pr.HoursPerDay != nil && (*pr.HoursPerDay < 0 || *pr.HoursPerDay > 24) {
		return errors.New("hours per day must be between 0 and 24")
	}

	if pr.DaysPerWeek != nil && (*pr.DaysPerWeek < 0 || *pr.DaysPerWeek > 7) {
		return errors.New("days per week must be between 0 and 7")
	}

	if pr.DaysPerMonth != nil && (*pr.DaysPerMonth < 0 || *pr.DaysPerMonth > 31) {
		return errors.New("days per month must be between 0 and 31")
	}

	// Validate rates if set
	if pr.HourlyRate != nil && *pr.HourlyRate < 0 {
		return errors.New("hourly rate cannot be negative")
	}

	if pr.DailyRate != nil && *pr.DailyRate < 0 {
		return errors.New("daily rate cannot be negative")
	}

	if pr.MonthlyRate != nil && *pr.MonthlyRate < 0 {
		return errors.New("monthly rate cannot be negative")
	}

	// Validate estimated man-months
	if pr.EstimatedManMonths < 0 {
		return errors.New("estimated man-months cannot be negative")
	}

	// Validate dates
	if pr.StartDate != nil && pr.EndDate != nil {
		if pr.EndDate.Before(*pr.StartDate) {
			return errors.New("end date cannot be before start date")
		}
	}

	return nil
}

// GetEffectiveHoursPerDay returns the effective hours per day (project-specific or default)
func (pr *ProjectRole) GetEffectiveHoursPerDay(resource *Resource) float64 {
	if pr.HoursPerDay != nil {
		return *pr.HoursPerDay
	}
	if resource != nil {
		return resource.DefaultHoursPerDay
	}
	return 8.0 // fallback default
}

// GetEffectiveDaysPerMonth returns the effective days per month (project-specific or default)
func (pr *ProjectRole) GetEffectiveDaysPerMonth(resource *Resource) float64 {
	if pr.DaysPerMonth != nil {
		return *pr.DaysPerMonth
	}
	if resource != nil {
		return resource.DefaultDaysPerMonth
	}
	return 20.0 // fallback default
}

// GetEffectiveHourlyRate returns the effective hourly rate (project-specific or default)
func (pr *ProjectRole) GetEffectiveHourlyRate(resource *Resource) float64 {
	if pr.HourlyRate != nil {
		return *pr.HourlyRate
	}
	if resource != nil {
		return resource.DefaultHourlyRate
	}
	return 0.0
}

// TaskAssignment represents the assignment of a resource to a task
type TaskAssignment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Task and Resource Assignment
	TaskID        uint  `gorm:"not null;index:idx_task_assignment_task" json:"task_id"`
	ProjectRoleID uint  `gorm:"not null;index:idx_task_assignment_role" json:"project_role_id"`
	ResourceID    uint  `gorm:"not null;index:idx_task_assignment_resource" json:"resource_id"` // Denormalized for easier querying

	// Effort Estimation (in man-days)
	EstimatedManDays float64 `gorm:"type:decimal(10,2);not null" json:"estimated_man_days"` // Estimated effort in man-days
	ActualManDays    float64 `gorm:"type:decimal(10,2);default:0" json:"actual_man_days"`   // Actual effort spent

	// Allocation Percentage (0-100)
	AllocationPercent float64 `gorm:"type:decimal(5,2);default:100" json:"allocation_percent"` // How much of the resource's time is allocated

	// Status
	IsActive bool   `gorm:"default:true" json:"is_active"`
	Notes    string `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Task        *Task        `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	ProjectRole *ProjectRole `gorm:"foreignKey:ProjectRoleID" json:"project_role,omitempty"`
	Resource    *Resource    `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
}

// TableName specifies the table name for TaskAssignment model
func (TaskAssignment) TableName() string {
	return "task_assignments"
}

// BeforeSave is a GORM hook that runs before saving
func (ta *TaskAssignment) BeforeSave(tx *gorm.DB) error {
	// Validate estimated man-days
	if ta.EstimatedManDays < 0 {
		return errors.New("estimated man-days cannot be negative")
	}

	// Validate actual man-days
	if ta.ActualManDays < 0 {
		return errors.New("actual man-days cannot be negative")
	}

	// Validate allocation percentage
	if ta.AllocationPercent < 0 || ta.AllocationPercent > 100 {
		return errors.New("allocation percent must be between 0 and 100")
	}

	return nil
}

// EstimatedHours converts estimated man-days to hours (assuming 8-hour workday)
func (ta *TaskAssignment) EstimatedHours() float64 {
	return ta.EstimatedManDays * 8.0
}

// ActualHours converts actual man-days to hours (assuming 8-hour workday)
func (ta *TaskAssignment) ActualHours() float64 {
	return ta.ActualManDays * 8.0
}
