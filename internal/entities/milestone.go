package entities

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	MilestoneStatusUnknown  = 0
	MilestoneStatusInactive = 1
	MilestoneStatusActive   = 2
)

var (
	ErrMilestoneNameRequired      = errors.New("milestone name is required")
	ErrMilestoneInvalidStatus     = errors.New("milestone status must be 1 (inactive) or 2 (active)")
	ErrMilestoneInvalidProjectID  = errors.New("milestone must belong to a project")
	ErrMilestoneInvalidDates      = errors.New("milestone end date must be on or after start date")

	MilestoneAllowedSortField = map[string]string{
		"id":          "id",
		"name":        "name",
		"description": "description",
		"project_id":  "project_id",
		"start_date":  "start_date",
		"end_date":    "end_date",
		"status":      "status",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
	}
)

// Milestone represents a milestone entity within a project
type Milestone struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	ProjectID   uint       `gorm:"not null;index" json:"project_id"`
	StartDate   *time.Time `gorm:"" json:"start_date"`
	EndDate     *time.Time `gorm:"" json:"end_date"`
	Status      uint       `gorm:"not null;default:2" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime:milli" json:"updated_at"`

	// Relationships
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName returns the table name for the milestone entity
func (Milestone) TableName() string {
	return "milestones"
}

// IsActive returns true if the milestone is active
func (m *Milestone) IsActive() bool {
	return m.Status == MilestoneStatusActive
}

// Validate validates the milestone fields
func (m *Milestone) Validate() error {
	// Trim whitespace from string fields
	m.Name = strings.TrimSpace(m.Name)
	m.Description = strings.TrimSpace(m.Description)

	// Validate required fields
	if m.Name == "" {
		return ErrMilestoneNameRequired
	}

	// Validate project ID
	if m.ProjectID == 0 {
		return ErrMilestoneInvalidProjectID
	}

	// Validate dates
	if m.StartDate != nil && m.EndDate != nil {
		if m.EndDate.Before(*m.StartDate) {
			return ErrMilestoneInvalidDates
		}
	}

	// Validate status
	if err := m.validateStatus(); err != nil {
		return err
	}

	return nil
}

func (m *Milestone) validateStatus() error {
	switch m.Status {
	case MilestoneStatusActive, MilestoneStatusInactive:
		return nil
	}
	return ErrMilestoneInvalidStatus
}

// BeforeCreate is a GORM hook that runs before creating a milestone
func (m *Milestone) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not valid
	if err := m.validateStatus(); err != nil {
		m.Status = MilestoneStatusActive
	}

	return m.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a milestone
func (m *Milestone) BeforeUpdate(tx *gorm.DB) error {
	return m.Validate()
}

type MilestoneQueryParams struct {
	ID_In            []uint     `json:"id_in"`
	Name             string     `json:"name"`
	Name_Like        string     `json:"name_like"`
	Description_Like string     `json:"description_like"`
	ProjectID        uint       `json:"project_id"`
	ProjectID_In     []uint     `json:"project_id_in"`
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

// MilestoneListResponse represents the response for GetMilestones
type MilestoneListResponse struct {
	Data  []*Milestone `json:"data"`
	Total int64        `json:"total"`
}
