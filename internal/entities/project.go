package entities

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	DefaultProjectHoursPerDay = 8
	DefaultProjectDaysPerWeek = 5
)

const (
	ProjectStatusUnknown  = 0
	ProjectStatusInactive = 1
	ProjectStatusActive   = 2
)

var (
	ErrProjectNameRequired           = errors.New("project name is required")
	ErrProjectInvalidStatus          = errors.New("project status must be 1 (inactive) or 2 (active)")
	ErrProjectInvalidClientID        = errors.New("project must belong to a client")
	ErrProjectInvalidDates           = errors.New("project end date must be on or after start date")
	ErrProjectInvalidHoursPerDay     = errors.New("hours per day must be between 1 and 24")
	ErrProjectInvalidDaysPerWeek     = errors.New("days per week must be between 1 and 7")
	ErrProjectInvalidWorkingDays     = errors.New("working days must contain valid weekdays (Sunday=0 to Saturday=6)")
	ErrProjectDuplicateWorkingDays   = errors.New("working days must not contain duplicates")
	ErrProjectWorkingDaysExceedsWeek = errors.New("working days cannot exceed 7 days")

	ProjectAllowedSortField = map[string]string{
		"id":          "id",
		"name":        "name",
		"description": "description",
		"client_id":   "client_id",
		"start_date":  "start_date",
		"end_date":    "end_date",
		"status":      "status",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
	}
)

// Project represents a project entity
type Project struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	ClientID    uint       `gorm:"not null;index" json:"client_id"`
	StartDate   *time.Time `gorm:"" json:"start_date"`
	EndDate     *time.Time `gorm:"" json:"end_date"`
	Status      uint       `gorm:"not null;default:2" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime:milli" json:"updated_at"`

	// Configurations
	HoursPerDay        int          `gorm:"default:8" json:"hours_per_day"`
	DaysPerWeek        int          `gorm:"default:5" json:"days_per_week"`
	WorkingDaysPerWeek WeekdayArray `gorm:"type:text" json:"working_days_per_week"`
	Timezone           string       `gorm:"default:''" json:"timezone"`
	Currency           string       `gorm:"default:''" json:"currency"`

	// Relationships
	Client           *Client            `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	ProjectResources []*ProjectResource `gorm:"foreignKey:ProjectID" json:"project_resources,omitempty"`
}

// TableName returns the table name for the project entity
func (Project) TableName() string {
	return "projects"
}

// IsActive returns true if the project is active
func (p *Project) IsActive() bool {
	return p.Status == ProjectStatusActive
}

// GetHoursPerDay returns the project's hours per day or the default if not set
func (p *Project) GetHoursPerDay() int {
	if p.HoursPerDay == 0 {
		return DefaultProjectHoursPerDay
	}
	return p.HoursPerDay
}

// GetDaysPerWeek returns the project's days per week or the default if not set
func (p *Project) GetDaysPerWeek() int {
	if p.DaysPerWeek == 0 {
		return DefaultProjectDaysPerWeek
	}
	return p.DaysPerWeek
}

// GetWorkingDaysPerWeek returns the project's working days or the default if not set
func (p *Project) GetWorkingDaysPerWeek() WeekdayArray {
	if len(p.WorkingDaysPerWeek) == 0 {
		return DefaultWorkingDays()
	}
	return p.WorkingDaysPerWeek
}

// Validate validates the project fields
func (p *Project) Validate() error {
	// Trim whitespace from string fields
	p.Name = strings.TrimSpace(p.Name)
	p.Description = strings.TrimSpace(p.Description)
	p.Timezone = strings.TrimSpace(p.Timezone)
	p.Currency = strings.TrimSpace(p.Currency)

	// Validate required fields
	if p.Name == "" {
		return ErrProjectNameRequired
	}

	// Validate client ID
	if p.ClientID == 0 {
		return ErrProjectInvalidClientID
	}

	// Validate dates
	if p.StartDate != nil && p.EndDate != nil {
		if p.EndDate.Before(*p.StartDate) {
			return ErrProjectInvalidDates
		}
	}

	// Validate status
	if err := p.validateStatus(); err != nil {
		return err
	}

	// Validate configuration fields
	if p.HoursPerDay != 0 && (p.HoursPerDay < 1 || p.HoursPerDay > 24) {
		return ErrProjectInvalidHoursPerDay
	}

	if p.DaysPerWeek != 0 && (p.DaysPerWeek < 1 || p.DaysPerWeek > 7) {
		return ErrProjectInvalidDaysPerWeek
	}

	// Validate working days
	if err := p.validateWorkingDays(); err != nil {
		return err
	}

	return nil
}

func (p *Project) validateWorkingDays() error {
	if len(p.WorkingDaysPerWeek) == 0 {
		return nil // Empty is allowed (will use defaults)
	}

	if len(p.WorkingDaysPerWeek) > 7 {
		return ErrProjectWorkingDaysExceedsWeek
	}

	seen := make(map[time.Weekday]bool)
	for _, day := range p.WorkingDaysPerWeek {
		// Validate weekday value is in range 0-6
		if day < time.Sunday || day > time.Saturday {
			return ErrProjectInvalidWorkingDays
		}
		// Check for duplicates
		if seen[day] {
			return ErrProjectDuplicateWorkingDays
		}
		seen[day] = true
	}

	return nil
}

func (p *Project) validateStatus() error {
	switch p.Status {
	case ProjectStatusActive, ProjectStatusInactive:
		return nil
	}
	return ErrProjectInvalidStatus
}

// BeforeCreate is a GORM hook that runs before creating a project
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not valid
	if err := p.validateStatus(); err != nil {
		p.Status = ProjectStatusActive
	}

	// Set default working days if not set (HoursPerDay and DaysPerWeek use GORM defaults)
	if len(p.WorkingDaysPerWeek) == 0 {
		p.WorkingDaysPerWeek = DefaultWorkingDays()
	}

	return p.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a project
func (p *Project) BeforeUpdate(tx *gorm.DB) error {
	return p.Validate()
}

type ProjectQueryParams struct {
	ID_In            []uint     `json:"id_in"`
	Name             string     `json:"name"`
	Name_Like        string     `json:"name_like"`
	Description_Like string     `json:"description_like"`
	ClientID         uint       `json:"client_id"`
	ClientID_In      []uint     `json:"client_id_in"`
	Status           uint       `json:"status"`
	Status_In        []uint     `json:"status_in"`
	StartDate_Gte    *time.Time `json:"start_date_gte"`
	StartDate_Lte    *time.Time `json:"start_date_lte"`
	EndDate_Gte      *time.Time `json:"end_date_gte"`
	EndDate_Lte      *time.Time `json:"end_date_lte"`
	CreatedAt_Gte    *time.Time `json:"created_at_gte"`
	CreatedAt_Lte    *time.Time `json:"created_at_lte"`
	UpdatedAt_Gte    *time.Time `json:"updated_at_gte"`
	UpdatedAt_Lte    *time.Time `json:"updated_at_lte"`
	*QueryParams
}

// ProjectListResponse represents the response for GetProjects
type ProjectListResponse struct {
	Data  []*Project `json:"data"`
	Total int64      `json:"total"`
}
