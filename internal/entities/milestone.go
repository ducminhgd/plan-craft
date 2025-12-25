package entities

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

// MilestoneQuery is the query for searching milestones
type MilestoneQuery struct {
	// ID_In is a list of IDs to search for
	ID_In []uint `json:"id__in"`

	// ProjectID is the project ID
	ProjectID *uint `json:"project_id"`
	// ProjectID_In is a list of project IDs to search for
	ProjectID_In []uint `json:"project_id__in"`

	// Name is the name of the milestone
	Name string `json:"name"`
	// Name_Like is the name of the milestone to search for (case-insensitive)
	Name_Like string `json:"name__like"`

	// Status is the status of the milestone
	Status TaskStatus `json:"status"`
	// Status_In is a list of statuses to search for
	Status_In []TaskStatus `json:"status__in"`

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

	// PlannedStartDate_Gte is the minimum planned start date
	PlannedStartDate_Gte *time.Time `json:"planned_start_date__gte"`
	// PlannedStartDate_Lte is the maximum planned start date
	PlannedStartDate_Lte *time.Time `json:"planned_start_date__lte"`

	// PlannedEndDate_Gte is the minimum planned end date
	PlannedEndDate_Gte *time.Time `json:"planned_end_date__gte"`
	// PlannedEndDate_Lte is the maximum planned end date
	PlannedEndDate_Lte *time.Time `json:"planned_end_date__lte"`

	// ActualStartDate_Gte is the minimum actual start date
	ActualStartDate_Gte *time.Time `json:"actual_start_date__gte"`
	// ActualStartDate_Lte is the maximum actual start date
	ActualStartDate_Lte *time.Time `json:"actual_start_date__lte"`

	// ActualEndDate_Gte is the minimum actual end date
	ActualEndDate_Gte *time.Time `json:"actual_end_date__gte"`
	// ActualEndDate_Lte is the maximum actual end date
	ActualEndDate_Lte *time.Time `json:"actual_end_date__lte"`

	// CreatedAt_Gte is the start time of the milestone creation time to search for
	CreatedAt_Gte *time.Time `json:"created_at__gte"`
	// CreatedAt_Lte is the end time of the milestone creation time to search for
	CreatedAt_Lte *time.Time `json:"created_at__lte"`

	// UpdatedAt_Gte is the start time of the milestone update time to search for
	UpdatedAt_Gte *time.Time `json:"updated_at__gte"`
	// UpdatedAt_Lte is the end time of the milestone update time to search for
	UpdatedAt_Lte *time.Time `json:"updated_at__lte"`

	// Tags_Contains searches for milestones containing all specified tags
	Tags_Contains []string `json:"tags__contains"`

	// QueryParams holds pagination, sorting, and filtering options
	QueryParams `json:",inline"`
}

// AllowedSortFields returns the allowed fields for sorting
func (q *MilestoneQuery) AllowedSortFields() map[string]string {
	return map[string]string{
		"id":                 "id",
		"name":               "name",
		"status":             "status",
		"order":              "order",
		"progress":           "progress",
		"estimated_effort":   "estimated_effort",
		"actual_effort":      "actual_effort",
		"estimated_cost":     "estimated_cost",
		"actual_cost":        "actual_cost",
		"planned_start_date": "planned_start_date",
		"planned_end_date":   "planned_end_date",
		"actual_start_date":  "actual_start_date",
		"actual_end_date":    "actual_end_date",
		"created_at":         "created_at",
		"updated_at":         "updated_at",
	}
}

// AllowedFilterFields returns the allowed fields for filtering
func (q *MilestoneQuery) AllowedFilterFields() map[string]string {
	return map[string]string{
		"id":                 "id",
		"project_id":         "project_id",
		"name":               "name",
		"status":             "status",
		"progress":           "progress",
		"estimated_effort":   "estimated_effort",
		"actual_effort":      "actual_effort",
		"estimated_cost":     "estimated_cost",
		"actual_cost":        "actual_cost",
		"planned_start_date": "planned_start_date",
		"planned_end_date":   "planned_end_date",
		"actual_start_date":  "actual_start_date",
		"actual_end_date":    "actual_end_date",
		"created_at":         "created_at",
		"updated_at":         "updated_at",
	}
}
