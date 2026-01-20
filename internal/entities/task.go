package entities

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Task status constants (numeric values for database storage)
const (
	TaskWorkStatusUnknown    = 0
	TaskWorkStatusToDo       = 1
	TaskWorkStatusInProgress = 2
	TaskWorkStatusDone       = 3
	TaskWorkStatusCancelled  = 4
)

// Task priority constants (numeric values for database storage)
const (
	TaskPriorityUnknown  = 0
	TaskPriorityLow      = 1
	TaskPriorityMedium   = 2
	TaskPriorityHigh     = 3
	TaskPriorityCritical = 4
)

var (
	ErrTaskNameRequired       = errors.New("task name is required")
	ErrTaskInvalidStatus      = errors.New("task status must be 1 (to do), 2 (in progress), 3 (done), or 4 (cancelled)")
	ErrTaskInvalidPriority    = errors.New("task priority must be 1 (low), 2 (medium), 3 (high), or 4 (critical)")
	ErrTaskInvalidProjectID   = errors.New("task must belong to a project")
	ErrTaskInvalidLevel       = errors.New("task level must be at least 1")
	ErrTaskInvalidEffort      = errors.New("task estimated effort must be non-negative")
	ErrTaskCircularDependency = errors.New("task cannot be its own parent")

	TaskAllowedSortField = map[string]string{
		"id":               "id",
		"name":             "name",
		"level":            "level",
		"project_id":       "project_id",
		"milestone_id":     "milestone_id",
		"parent_id":        "parent_id",
		"priority":         "priority",
		"status":           "status",
		"estimated_effort": "estimated_effort",
		"created_at":       "created_at",
		"updated_at":       "updated_at",
	}
)

// Task represents a task entity within a project
type Task struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	Name            string     `gorm:"not null" json:"name"`
	Description     string     `gorm:"type:text" json:"description"`
	Level           int        `gorm:"not null;default:1" json:"level"`
	ProjectID       uint       `gorm:"not null;index" json:"project_id"`
	MilestoneID     *uint      `gorm:"index" json:"milestone_id"`
	ParentID        *uint      `gorm:"index" json:"parent_id"`
	Priority        uint       `gorm:"not null;default:2" json:"priority"`
	EstimatedEffort float64    `gorm:"not null;default:0" json:"estimated_effort"`
	Status          uint       `gorm:"not null;default:1" json:"status"`
	CreatedAt       time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime:milli" json:"updated_at"`

	// Relationships
	Project   *Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Milestone *Milestone `gorm:"foreignKey:MilestoneID" json:"milestone,omitempty"`
	Parent    *Task      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []*Task    `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName returns the table name for the task entity
func (Task) TableName() string {
	return "tasks"
}

// IsToDo returns true if the task is in to do status
func (t *Task) IsToDo() bool {
	return t.Status == TaskWorkStatusToDo
}

// IsInProgress returns true if the task is in progress
func (t *Task) IsInProgress() bool {
	return t.Status == TaskWorkStatusInProgress
}

// IsDone returns true if the task is done
func (t *Task) IsDone() bool {
	return t.Status == TaskWorkStatusDone
}

// IsCancelled returns true if the task is cancelled
func (t *Task) IsCancelled() bool {
	return t.Status == TaskWorkStatusCancelled
}

// Validate validates the task fields
func (t *Task) Validate() error {
	// Trim whitespace from string fields
	t.Name = strings.TrimSpace(t.Name)
	t.Description = strings.TrimSpace(t.Description)

	// Validate required fields
	if t.Name == "" {
		return ErrTaskNameRequired
	}

	// Validate project ID
	if t.ProjectID == 0 {
		return ErrTaskInvalidProjectID
	}

	// Validate level
	if t.Level < 1 {
		return ErrTaskInvalidLevel
	}

	// Validate estimated effort
	if t.EstimatedEffort < 0 {
		return ErrTaskInvalidEffort
	}

	// Validate parent is not self
	if t.ParentID != nil && *t.ParentID == t.ID && t.ID != 0 {
		return ErrTaskCircularDependency
	}

	// Validate status
	if err := t.validateStatus(); err != nil {
		return err
	}

	// Validate priority
	if err := t.validatePriority(); err != nil {
		return err
	}

	return nil
}

func (t *Task) validateStatus() error {
	switch t.Status {
	case TaskWorkStatusToDo, TaskWorkStatusInProgress, TaskWorkStatusDone, TaskWorkStatusCancelled:
		return nil
	}
	return ErrTaskInvalidStatus
}

func (t *Task) validatePriority() error {
	switch t.Priority {
	case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh, TaskPriorityCritical:
		return nil
	}
	return ErrTaskInvalidPriority
}

// BeforeCreate is a GORM hook that runs before creating a task
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not valid
	if err := t.validateStatus(); err != nil {
		t.Status = TaskWorkStatusToDo
	}

	// Set default priority if not valid
	if err := t.validatePriority(); err != nil {
		t.Priority = TaskPriorityMedium
	}

	// Set default level if not valid
	if t.Level < 1 {
		t.Level = 1
	}

	return t.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a task
func (t *Task) BeforeUpdate(tx *gorm.DB) error {
	return t.Validate()
}

// TaskQueryParams defines query parameters for filtering tasks
type TaskQueryParams struct {
	ID_In               []uint     `json:"id_in"`
	Name                string     `json:"name"`
	Name_Like           string     `json:"name_like"`
	Description_Like    string     `json:"description_like"`
	Level               int        `json:"level"`
	Level_Gte           *int       `json:"level_gte"`
	Level_Lte           *int       `json:"level_lte"`
	ProjectID           uint       `json:"project_id"`
	ProjectID_In        []uint     `json:"project_id_in"`
	MilestoneID         *uint      `json:"milestone_id"`
	MilestoneID_In      []uint     `json:"milestone_id_in"`
	MilestoneID_IsNull  *bool      `json:"milestone_id_is_null"`
	ParentID            *uint      `json:"parent_id"`
	ParentID_In         []uint     `json:"parent_id_in"`
	ParentID_IsNull     *bool      `json:"parent_id_is_null"`
	Priority            uint       `json:"priority"`
	Priority_In         []uint     `json:"priority_in"`
	Status              uint       `json:"status"`
	Status_In           []uint     `json:"status_in"`
	EstimatedEffort_Gte *float64   `json:"estimated_effort_gte"`
	EstimatedEffort_Lte *float64   `json:"estimated_effort_lte"`
	CreatedAt_Gte       *time.Time `json:"created_at_gte"`
	CreatedAt_Lte       *time.Time `json:"created_at_lte"`
	UpdatedAt_Gte       *time.Time `json:"updated_at_gte"`
	UpdatedAt_Lte       *time.Time `json:"updated_at_lte"`
	*QueryParams
}

// TaskListResponse represents the response for GetTasks
type TaskListResponse struct {
	Data  []*Task `json:"data"`
	Total int64   `json:"total"`
}
