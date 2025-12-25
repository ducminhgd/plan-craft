package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Project represents a project in the system
type Project struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Basic Information
	Name        string `gorm:"type:varchar(255);not null;index" json:"name"`
	Code        string `gorm:"type:varchar(50);uniqueIndex" json:"code,omitempty"` // Optional project code/identifier
	Description string `gorm:"type:text" json:"description,omitempty"`

	// Project Classification
	Type   ProjectType `gorm:"type:varchar(50);not null;default:'product'" json:"type"`
	Status TaskStatus  `gorm:"type:varchar(50);not null;default:'not_started';index" json:"status"`

	// Timeline
	StartDate     *time.Time `gorm:"type:datetime" json:"start_date,omitempty"`
	TargetEndDate *time.Time `gorm:"type:datetime" json:"target_end_date,omitempty"`
	ActualEndDate *time.Time `gorm:"type:datetime" json:"actual_end_date,omitempty"`

	// Estimation (in hours)
	EstimatedEffort float64 `gorm:"type:decimal(10,2);default:0" json:"estimated_effort"` // Total estimated hours
	ActualEffort    float64 `gorm:"type:decimal(10,2);default:0" json:"actual_effort"`    // Total actual hours spent

	// Budget and Cost
	EstimatedCost float64 `gorm:"type:decimal(15,2);default:0" json:"estimated_cost"` // Total estimated cost
	ActualCost    float64 `gorm:"type:decimal(15,2);default:0" json:"actual_cost"`    // Total actual cost
	Currency      string  `gorm:"type:varchar(3);default:'USD'" json:"currency"`      // ISO 4217 currency code

	// Progress Tracking
	Progress float64 `gorm:"type:decimal(5,2);default:0" json:"progress"` // Percentage (0-100)

	// Work Time Configuration (nullable to use defaults if not set)
	HoursPerDay   *float64 `gorm:"type:decimal(5,2)" json:"hours_per_day,omitempty"`   // Working hours per day (default: 8.0)
	DaysPerWeek   *float64 `gorm:"type:decimal(5,2)" json:"days_per_week,omitempty"`   // Working days per week (default: 5.0)
	DaysPerMonth  *float64 `gorm:"type:decimal(5,2)" json:"days_per_month,omitempty"`  // Working days per month (default: 20.0)

	// Metadata
	Assumptions StringArray `gorm:"type:json" json:"assumptions,omitempty"` // List of project assumptions
	Constraints StringArray `gorm:"type:json" json:"constraints,omitempty"` // List of project constraints
	Tags        StringArray `gorm:"type:json" json:"tags,omitempty"`        // Tags for categorization

	// Ownership
	OwnerID  *uint `gorm:"index" json:"owner_id,omitempty"`  // User ID of project owner (project manager)
	ClientID *uint `gorm:"index" json:"client_id,omitempty"` // Client/customer ID

	// Relationships
	Client              *Client              `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	Tasks               []Task               `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"`
	Milestones          []Milestone          `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"milestones,omitempty"`
	ProjectRoles        []ProjectRole        `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"project_roles,omitempty"`
	Costs               []Cost               `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"costs,omitempty"`
	ResourceAllocations []ResourceAllocation `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"resource_allocations,omitempty"`
}

// TableName specifies the table name for Project model
func (Project) TableName() string {
	return "projects"
}

// BeforeSave is a GORM hook that runs before saving
func (p *Project) BeforeSave(tx *gorm.DB) error {
	// Validate project type
	if !IsValidProjectType(p.Type) {
		return errors.New("invalid project type")
	}

	// Validate status
	if !IsValidTaskStatus(p.Status) {
		return errors.New("invalid status")
	}

	// Validate progress
	if p.Progress < 0 || p.Progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}

	// Validate dates
	if p.StartDate != nil && p.TargetEndDate != nil {
		if p.TargetEndDate.Before(*p.StartDate) {
			return errors.New("target end date cannot be before start date")
		}
	}

	if p.StartDate != nil && p.ActualEndDate != nil {
		if p.ActualEndDate.Before(*p.StartDate) {
			return errors.New("actual end date cannot be before start date")
		}
	}

	// Validate work time configuration
	if p.HoursPerDay != nil && *p.HoursPerDay <= 0 {
		return errors.New("hours per day must be positive")
	}
	if p.DaysPerWeek != nil && *p.DaysPerWeek <= 0 {
		return errors.New("days per week must be positive")
	}
	if p.DaysPerMonth != nil && *p.DaysPerMonth <= 0 {
		return errors.New("days per month must be positive")
	}

	return nil
}

// IsActive returns true if the project is currently active
func (p *Project) IsActive() bool {
	return p.Status == TaskStatusInProgress
}

// IsCompleted returns true if the project is completed
func (p *Project) IsCompleted() bool {
	return p.Status == TaskStatusCompleted
}

// Duration returns the planned duration of the project
func (p *Project) Duration() *time.Duration {
	if p.StartDate == nil || p.TargetEndDate == nil {
		return nil
	}
	duration := p.TargetEndDate.Sub(*p.StartDate)
	return &duration
}

// ActualDuration returns the actual duration of the project
func (p *Project) ActualDuration() *time.Duration {
	if p.StartDate == nil || p.ActualEndDate == nil {
		return nil
	}
	duration := p.ActualEndDate.Sub(*p.StartDate)
	return &duration
}

// GetHoursPerDay returns the hours per day for this project (or default if not set)
func (p *Project) GetHoursPerDay() float64 {
	if p.HoursPerDay != nil {
		return *p.HoursPerDay
	}
	return DefaultHoursPerDay
}

// GetDaysPerWeek returns the days per week for this project (or default if not set)
func (p *Project) GetDaysPerWeek() float64 {
	if p.DaysPerWeek != nil {
		return *p.DaysPerWeek
	}
	return DefaultDaysPerWeek
}

// GetDaysPerMonth returns the days per month for this project (or default if not set)
func (p *Project) GetDaysPerMonth() float64 {
	if p.DaysPerMonth != nil {
		return *p.DaysPerMonth
	}
	return DefaultDaysPerMonth
}

// GetHoursPerMonth returns the calculated hours per month for this project
func (p *Project) GetHoursPerMonth() float64 {
	return p.GetHoursPerDay() * p.GetDaysPerMonth()
}

// EstimatedDays converts estimated effort from hours to days using project settings
func (p *Project) EstimatedDays() float64 {
	return p.EstimatedEffort / p.GetHoursPerDay()
}

// EstimatedMonths converts estimated effort from hours to man-months using project settings
func (p *Project) EstimatedMonths() float64 {
	return p.EstimatedEffort / p.GetHoursPerMonth()
}

// ActualDays converts actual effort from hours to days using project settings
func (p *Project) ActualDays() float64 {
	return p.ActualEffort / p.GetHoursPerDay()
}

// ActualMonths converts actual effort from hours to man-months using project settings
func (p *Project) ActualMonths() float64 {
	return p.ActualEffort / p.GetHoursPerMonth()
}

// StringArray is a custom type for storing string arrays in JSON format
type StringArray []string

// Scan implements the sql.Scanner interface
func (sa *StringArray) Scan(value any) error {
	if value == nil {
		*sa = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}

	return json.Unmarshal(bytes, sa)
}

// Value implements the driver.Valuer interface
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return nil, nil
	}
	return json.Marshal(sa)
}

// JSONB is a custom type for storing arbitrary JSON data
type JSONB map[string]any

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = make(map[string]any)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONB")
	}

	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// ProjectQuery is the query for searching projects
type ProjectQuery struct {
	// ID_In is a list of IDs to search for
	ID_In []uint `json:"id__in"`

	// Name is the name of the project
	Name string `json:"name"`
	// Name_Like is the name of the project to search for (case-insensitive)
	Name_Like string `json:"name__like"`

	// Code is the project code
	Code string `json:"code"`
	// Code_Like is the project code to search for (case-insensitive)
	Code_Like string `json:"code__like"`

	// Type is the type of the project
	Type ProjectType `json:"type"`
	// Type_In is a list of types to search for
	Type_In []ProjectType `json:"type__in"`

	// Status is the status of the project
	Status TaskStatus `json:"status"`
	// Status_In is a list of statuses to search for
	Status_In []TaskStatus `json:"status__in"`

	// ClientID is the client ID
	ClientID *uint `json:"client_id"`
	// ClientID_In is a list of client IDs to search for
	ClientID_In []uint `json:"client_id__in"`

	// OwnerID is the owner ID
	OwnerID *uint `json:"owner_id"`
	// OwnerID_In is a list of owner IDs to search for
	OwnerID_In []uint `json:"owner_id__in"`

	// Currency is the currency code
	Currency string `json:"currency"`
	// Currency_In is a list of currency codes to search for
	Currency_In []string `json:"currency__in"`

	// Progress_Gte is the minimum progress percentage
	Progress_Gte *float64 `json:"progress__gte"`
	// Progress_Lte is the maximum progress percentage
	Progress_Lte *float64 `json:"progress__lte"`

	// EstimatedEffort_Gte is the minimum estimated effort
	EstimatedEffort_Gte *float64 `json:"estimated_effort__gte"`
	// EstimatedEffort_Lte is the maximum estimated effort
	EstimatedEffort_Lte *float64 `json:"estimated_effort__lte"`

	// ActualEffort_Gte is the minimum actual effort
	ActualEffort_Gte *float64 `json:"actual_effort__gte"`
	// ActualEffort_Lte is the maximum actual effort
	ActualEffort_Lte *float64 `json:"actual_effort__lte"`

	// EstimatedCost_Gte is the minimum estimated cost
	EstimatedCost_Gte *float64 `json:"estimated_cost__gte"`
	// EstimatedCost_Lte is the maximum estimated cost
	EstimatedCost_Lte *float64 `json:"estimated_cost__lte"`

	// ActualCost_Gte is the minimum actual cost
	ActualCost_Gte *float64 `json:"actual_cost__gte"`
	// ActualCost_Lte is the maximum actual cost
	ActualCost_Lte *float64 `json:"actual_cost__lte"`

	// StartDate_Gte is the minimum start date
	StartDate_Gte *time.Time `json:"start_date__gte"`
	// StartDate_Lte is the maximum start date
	StartDate_Lte *time.Time `json:"start_date__lte"`

	// TargetEndDate_Gte is the minimum target end date
	TargetEndDate_Gte *time.Time `json:"target_end_date__gte"`
	// TargetEndDate_Lte is the maximum target end date
	TargetEndDate_Lte *time.Time `json:"target_end_date__lte"`

	// ActualEndDate_Gte is the minimum actual end date
	ActualEndDate_Gte *time.Time `json:"actual_end_date__gte"`
	// ActualEndDate_Lte is the maximum actual end date
	ActualEndDate_Lte *time.Time `json:"actual_end_date__lte"`

	// CreatedAt_Gte is the start time of the project creation time to search for
	CreatedAt_Gte *time.Time `json:"created_at__gte"`
	// CreatedAt_Lte is the end time of the project creation time to search for
	CreatedAt_Lte *time.Time `json:"created_at__lte"`

	// UpdatedAt_Gte is the start time of the project update time to search for
	UpdatedAt_Gte *time.Time `json:"updated_at__gte"`
	// UpdatedAt_Lte is the end time of the project update time to search for
	UpdatedAt_Lte *time.Time `json:"updated_at__lte"`

	// Tags_Contains searches for projects containing all specified tags
	Tags_Contains []string `json:"tags__contains"`

	// QueryParams holds pagination, sorting, and filtering options
	QueryParams `json:",inline"`
}

// AllowedSortFields returns the allowed fields for sorting
func (q *ProjectQuery) AllowedSortFields() map[string]string {
	return map[string]string{
		"id":                "id",
		"name":              "name",
		"code":              "code",
		"type":              "type",
		"status":            "status",
		"progress":          "progress",
		"estimated_effort":  "estimated_effort",
		"actual_effort":     "actual_effort",
		"estimated_cost":    "estimated_cost",
		"actual_cost":       "actual_cost",
		"start_date":        "start_date",
		"target_end_date":   "target_end_date",
		"actual_end_date":   "actual_end_date",
		"created_at":        "created_at",
		"updated_at":        "updated_at",
	}
}

// AllowedFilterFields returns the allowed fields for filtering
func (q *ProjectQuery) AllowedFilterFields() map[string]string {
	return map[string]string{
		"id":                "id",
		"name":              "name",
		"code":              "code",
		"type":              "type",
		"status":            "status",
		"client_id":         "client_id",
		"owner_id":          "owner_id",
		"currency":          "currency",
		"progress":          "progress",
		"estimated_effort":  "estimated_effort",
		"actual_effort":     "actual_effort",
		"estimated_cost":    "estimated_cost",
		"actual_cost":       "actual_cost",
		"start_date":        "start_date",
		"target_end_date":   "target_end_date",
		"actual_end_date":   "actual_end_date",
		"created_at":        "created_at",
		"updated_at":        "updated_at",
	}
}

