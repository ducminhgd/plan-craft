package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Milestone represents a project milestone or phase
type Milestone struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Basic Information
	ProjectID   uint   `gorm:"not null;index:idx_milestone_project" json:"project_id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Order       int    `gorm:"default:0" json:"order"` // Order within project

	// Status
	Status TaskStatus `gorm:"type:varchar(50);not null;default:'not_started';index" json:"status"`

	// Timeline
	PlannedStartDate *time.Time `gorm:"type:datetime" json:"planned_start_date,omitempty"`
	PlannedEndDate   *time.Time `gorm:"type:datetime" json:"planned_end_date,omitempty"`
	ActualStartDate  *time.Time `gorm:"type:datetime" json:"actual_start_date,omitempty"`
	ActualEndDate    *time.Time `gorm:"type:datetime" json:"actual_end_date,omitempty"`

	// Estimation (in hours)
	EstimatedEffort float64 `gorm:"type:decimal(10,2);default:0" json:"estimated_effort"` // Estimated hours (can be summed from tasks)
	ActualEffort    float64 `gorm:"type:decimal(10,2);default:0" json:"actual_effort"`    // Actual hours spent

	// Cost
	EstimatedCost float64 `gorm:"type:decimal(15,2);default:0" json:"estimated_cost"` // Estimated cost (can be summed from tasks)
	ActualCost    float64 `gorm:"type:decimal(15,2);default:0" json:"actual_cost"`    // Actual cost

	// Progress
	Progress float64 `gorm:"type:decimal(5,2);default:0" json:"progress"` // Percentage (0-100)

	// Metadata
	Notes string      `gorm:"type:text" json:"notes,omitempty"`
	Tags  StringArray `gorm:"type:json" json:"tags,omitempty"`

	// Relationships
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Tasks   []Task   `gorm:"foreignKey:MilestoneID;constraint:OnDelete:SET NULL" json:"tasks,omitempty"`
	Costs   []Cost   `gorm:"foreignKey:MilestoneID;constraint:OnDelete:SET NULL" json:"costs,omitempty"`
}

// TableName specifies the table name for Milestone model
func (Milestone) TableName() string {
	return "milestones"
}

// BeforeSave is a GORM hook that runs before saving
func (m *Milestone) BeforeSave(tx *gorm.DB) error {
	// Validate status
	if !IsValidTaskStatus(m.Status) {
		return errors.New("invalid milestone status")
	}

	// Validate progress
	if m.Progress < 0 || m.Progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}

	// Validate dates
	if m.PlannedStartDate != nil && m.PlannedEndDate != nil {
		if m.PlannedEndDate.Before(*m.PlannedStartDate) {
			return errors.New("planned end date cannot be before planned start date")
		}
	}

	if m.ActualStartDate != nil && m.ActualEndDate != nil {
		if m.ActualEndDate.Before(*m.ActualStartDate) {
			return errors.New("actual end date cannot be before actual start date")
		}
	}

	return nil
}

// IsCompleted returns true if the milestone is completed
func (m *Milestone) IsCompleted() bool {
	return m.Status == TaskStatusCompleted
}

// PlannedDuration returns the planned duration
func (m *Milestone) PlannedDuration() *time.Duration {
	if m.PlannedStartDate == nil || m.PlannedEndDate == nil {
		return nil
	}
	duration := m.PlannedEndDate.Sub(*m.PlannedStartDate)
	return &duration
}

// ActualDuration returns the actual duration
func (m *Milestone) ActualDuration() *time.Duration {
	if m.ActualStartDate == nil || m.ActualEndDate == nil {
		return nil
	}
	duration := m.ActualEndDate.Sub(*m.ActualStartDate)
	return &duration
}

// EstimatedDays converts estimated hours to days
func (m *Milestone) EstimatedDays() float64 {
	return HoursToDays(m.EstimatedEffort)
}

// EstimatedMonths converts estimated hours to man-months
func (m *Milestone) EstimatedMonths() float64 {
	return HoursToMonths(m.EstimatedEffort)
}

// ActualDays converts actual hours to days
func (m *Milestone) ActualDays() float64 {
	return HoursToDays(m.ActualEffort)
}

// ActualMonths converts actual hours to man-months
func (m *Milestone) ActualMonths() float64 {
	return HoursToMonths(m.ActualEffort)
}
