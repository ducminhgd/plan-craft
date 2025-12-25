package entities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Task represents a work item in a project (supports hierarchical WBS)
type Task struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Basic Information
	ProjectID   uint   `gorm:"not null;index:idx_task_project" json:"project_id"`
	MilestoneID *uint  `gorm:"index:idx_task_milestone" json:"milestone_id,omitempty"` // Optional milestone association
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Code        string `gorm:"type:varchar(50);index" json:"code,omitempty"`        // WBS code like 1.2.3
	Description string `gorm:"type:text" json:"description,omitempty"`

	// Hierarchical Structure (WBS)
	ParentID *uint  `gorm:"index:idx_task_parent" json:"parent_id,omitempty"` // For epics → tasks → subtasks
	Level    int    `gorm:"default:1" json:"level"`                            // 1 = epic, 2 = task, 3 = subtask
	Order    int    `gorm:"default:0" json:"order"`                            // Order within parent

	// Status and Priority
	Status   TaskStatus `gorm:"type:varchar(50);not null;default:'not_started';index" json:"status"`
	Priority Priority   `gorm:"type:varchar(50);default:'medium'" json:"priority"`

	// Effort Estimation (in hours)
	EstimatedHours float64 `gorm:"type:decimal(10,2);default:0" json:"estimated_hours"` // Estimated effort in hours
	ActualHours    float64 `gorm:"type:decimal(10,2);default:0" json:"actual_hours"`    // Actual time spent

	// Timeline
	PlannedStartDate  *time.Time `gorm:"type:datetime" json:"planned_start_date,omitempty"`
	PlannedEndDate    *time.Time `gorm:"type:datetime" json:"planned_end_date,omitempty"`
	ActualStartDate   *time.Time `gorm:"type:datetime" json:"actual_start_date,omitempty"`
	ActualEndDate     *time.Time `gorm:"type:datetime" json:"actual_end_date,omitempty"`

	// Critical Path & Slack
	IsCriticalPath bool    `gorm:"default:false" json:"is_critical_path"` // Auto-calculated
	SlackDays      float64 `gorm:"type:decimal(10,2);default:0" json:"slack_days"` // Float slack in days

	// Progress
	Progress float64 `gorm:"type:decimal(5,2);default:0" json:"progress"` // Percentage (0-100)

	// Metadata
	Tags  StringArray `gorm:"type:json" json:"tags,omitempty"`
	Notes string      `gorm:"type:text" json:"notes,omitempty"`

	// Assignee
	AssigneeID *uint `gorm:"index" json:"assignee_id,omitempty"`

	// Relationships
	Project      *Project          `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Milestone    *Milestone        `gorm:"foreignKey:MilestoneID" json:"milestone,omitempty"`
	Parent       *Task             `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Subtasks     []Task            `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"subtasks,omitempty"`
	Dependencies []TaskDependency  `gorm:"foreignKey:DependentTaskID;constraint:OnDelete:CASCADE" json:"dependencies,omitempty"`
	Dependents   []TaskDependency  `gorm:"foreignKey:PredecessorTaskID;constraint:OnDelete:CASCADE" json:"dependents,omitempty"`
	Assignments  []TaskAssignment  `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"assignments,omitempty"`
	Costs        []Cost            `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"costs,omitempty"`
}

// TableName specifies the table name for Task model
func (Task) TableName() string {
	return "tasks"
}

// BeforeSave is a GORM hook that runs before saving
func (t *Task) BeforeSave(tx *gorm.DB) error {
	// Validate status
	if !IsValidTaskStatus(t.Status) {
		return errors.New("invalid task status")
	}

	// Validate priority
	if !IsValidPriority(t.Priority) {
		return errors.New("invalid priority")
	}

	// Validate progress
	if t.Progress < 0 || t.Progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}

	// Validate level
	if t.Level < 1 || t.Level > 10 {
		return errors.New("level must be between 1 and 10")
	}

	// Validate dates
	if t.PlannedStartDate != nil && t.PlannedEndDate != nil {
		if t.PlannedEndDate.Before(*t.PlannedStartDate) {
			return errors.New("planned end date cannot be before planned start date")
		}
	}

	if t.ActualStartDate != nil && t.ActualEndDate != nil {
		if t.ActualEndDate.Before(*t.ActualStartDate) {
			return errors.New("actual end date cannot be before actual start date")
		}
	}

	// Prevent circular parent reference
	if t.ParentID != nil && *t.ParentID == t.ID {
		return errors.New("task cannot be its own parent")
	}

	return nil
}

// IsEpic returns true if the task is an epic (level 1)
func (t *Task) IsEpic() bool {
	return t.Level == 1 && t.ParentID == nil
}

// IsSubtask returns true if the task has a parent
func (t *Task) IsSubtask() bool {
	return t.ParentID != nil
}

// IsCompleted returns true if the task is completed
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusCompleted
}

// PlannedDuration returns the planned duration in days
func (t *Task) PlannedDuration() *time.Duration {
	if t.PlannedStartDate == nil || t.PlannedEndDate == nil {
		return nil
	}
	duration := t.PlannedEndDate.Sub(*t.PlannedStartDate)
	return &duration
}

// ActualDuration returns the actual duration in days
func (t *Task) ActualDuration() *time.Duration {
	if t.ActualStartDate == nil || t.ActualEndDate == nil {
		return nil
	}
	duration := t.ActualEndDate.Sub(*t.ActualStartDate)
	return &duration
}

// EstimatedDays converts estimated hours to days (assuming 8-hour workday)
func (t *Task) EstimatedDays() float64 {
	return HoursToDays(t.EstimatedHours)
}

// ActualDays converts actual hours to days (assuming 8-hour workday)
func (t *Task) ActualDays() float64 {
	return HoursToDays(t.ActualHours)
}

// EstimatedMonths converts estimated hours to man-months
func (t *Task) EstimatedMonths() float64 {
	return HoursToMonths(t.EstimatedHours)
}

// ActualMonths converts actual hours to man-months
func (t *Task) ActualMonths() float64 {
	return HoursToMonths(t.ActualHours)
}

// TaskDependency represents a dependency between two tasks
type TaskDependency struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Dependency Relationship
	PredecessorTaskID uint           `gorm:"not null;index:idx_dependency_predecessor" json:"predecessor_task_id"` // Task that must be completed first
	DependentTaskID   uint           `gorm:"not null;index:idx_dependency_dependent" json:"dependent_task_id"`     // Task that depends on predecessor
	DependencyType    DependencyType `gorm:"type:varchar(50);not null;default:'finish_to_start'" json:"dependency_type"`

	// Lead/Lag Time (in days)
	LagDays  float64 `gorm:"type:decimal(10,2);default:0" json:"lag_days"`  // Delay after predecessor (positive value)
	LeadDays float64 `gorm:"type:decimal(10,2);default:0" json:"lead_days"` // Overlap before predecessor completes (positive value)

	// Metadata
	Notes string `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	PredecessorTask *Task `gorm:"foreignKey:PredecessorTaskID" json:"predecessor_task,omitempty"`
	DependentTask   *Task `gorm:"foreignKey:DependentTaskID" json:"dependent_task,omitempty"`
}

// TableName specifies the table name for TaskDependency model
func (TaskDependency) TableName() string {
	return "task_dependencies"
}

// BeforeSave is a GORM hook that runs before saving
func (td *TaskDependency) BeforeSave(tx *gorm.DB) error {
	// Validate dependency type
	if !IsValidDependencyType(td.DependencyType) {
		return errors.New("invalid dependency type")
	}

	// Prevent self-dependency
	if td.PredecessorTaskID == td.DependentTaskID {
		return errors.New("task cannot depend on itself")
	}

	// Validate that lag and lead are not both set
	if td.LagDays > 0 && td.LeadDays > 0 {
		return errors.New("cannot have both lag and lead time")
	}

	return nil
}

// TaskQuery is the query for searching tasks
type TaskQuery struct {
	// ID_In is a list of IDs to search for
	ID_In []uint `json:"id__in"`

	// ProjectID is the project ID
	ProjectID *uint `json:"project_id"`
	// ProjectID_In is a list of project IDs to search for
	ProjectID_In []uint `json:"project_id__in"`

	// MilestoneID is the milestone ID
	MilestoneID *uint `json:"milestone_id"`
	// MilestoneID_In is a list of milestone IDs to search for
	MilestoneID_In []uint `json:"milestone_id__in"`

	// ParentID is the parent task ID
	ParentID *uint `json:"parent_id"`
	// ParentID_IsNull filters for tasks without a parent
	ParentID_IsNull *bool `json:"parent_id__isnull"`

	// Name is the name of the task
	Name string `json:"name"`
	// Name_Like is the name of the task to search for (case-insensitive)
	Name_Like string `json:"name__like"`

	// Code is the task code (WBS code)
	Code string `json:"code"`
	// Code_Like is the task code to search for (case-insensitive)
	Code_Like string `json:"code__like"`

	// Status is the status of the task
	Status TaskStatus `json:"status"`
	// Status_In is a list of statuses to search for
	Status_In []TaskStatus `json:"status__in"`

	// Priority is the priority of the task
	Priority Priority `json:"priority"`
	// Priority_In is a list of priorities to search for
	Priority_In []Priority `json:"priority__in"`

	// Level is the hierarchical level
	Level *int `json:"level"`
	// Level_In is a list of levels to search for
	Level_In []int `json:"level__in"`

	// AssigneeID is the assignee ID
	AssigneeID *uint `json:"assignee_id"`
	// AssigneeID_In is a list of assignee IDs to search for
	AssigneeID_In []uint `json:"assignee_id__in"`

	// IsCriticalPath filters for critical path tasks
	IsCriticalPath *bool `json:"is_critical_path"`

	// Progress_Gte is the minimum progress percentage
	Progress_Gte *float64 `json:"progress__gte"`
	// Progress_Lte is the maximum progress percentage
	Progress_Lte *float64 `json:"progress__lte"`

	// EstimatedHours_Gte is the minimum estimated hours
	EstimatedHours_Gte *float64 `json:"estimated_hours__gte"`
	// EstimatedHours_Lte is the maximum estimated hours
	EstimatedHours_Lte *float64 `json:"estimated_hours__lte"`

	// ActualHours_Gte is the minimum actual hours
	ActualHours_Gte *float64 `json:"actual_hours__gte"`
	// ActualHours_Lte is the maximum actual hours
	ActualHours_Lte *float64 `json:"actual_hours__lte"`

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

	// CreatedAt_Gte is the start time of the task creation time to search for
	CreatedAt_Gte *time.Time `json:"created_at__gte"`
	// CreatedAt_Lte is the end time of the task creation time to search for
	CreatedAt_Lte *time.Time `json:"created_at__lte"`

	// UpdatedAt_Gte is the start time of the task update time to search for
	UpdatedAt_Gte *time.Time `json:"updated_at__gte"`
	// UpdatedAt_Lte is the end time of the task update time to search for
	UpdatedAt_Lte *time.Time `json:"updated_at__lte"`

	// Tags_Contains searches for tasks containing all specified tags
	Tags_Contains []string `json:"tags__contains"`

	// QueryParams holds pagination, sorting, and filtering options
	QueryParams `json:",inline"`
}

// AllowedSortFields returns the allowed fields for sorting
func (q *TaskQuery) AllowedSortFields() map[string]string {
	return map[string]string{
		"id":                 "id",
		"name":               "name",
		"code":               "code",
		"status":             "status",
		"priority":           "priority",
		"level":              "level",
		"order":              "order",
		"progress":           "progress",
		"estimated_hours":    "estimated_hours",
		"actual_hours":       "actual_hours",
		"planned_start_date": "planned_start_date",
		"planned_end_date":   "planned_end_date",
		"actual_start_date":  "actual_start_date",
		"actual_end_date":    "actual_end_date",
		"created_at":         "created_at",
		"updated_at":         "updated_at",
	}
}

// AllowedFilterFields returns the allowed fields for filtering
func (q *TaskQuery) AllowedFilterFields() map[string]string {
	return map[string]string{
		"id":                 "id",
		"project_id":         "project_id",
		"milestone_id":       "milestone_id",
		"parent_id":          "parent_id",
		"name":               "name",
		"code":               "code",
		"status":             "status",
		"priority":           "priority",
		"level":              "level",
		"assignee_id":        "assignee_id",
		"is_critical_path":   "is_critical_path",
		"progress":           "progress",
		"estimated_hours":    "estimated_hours",
		"actual_hours":       "actual_hours",
		"planned_start_date": "planned_start_date",
		"planned_end_date":   "planned_end_date",
		"actual_start_date":  "actual_start_date",
		"actual_end_date":    "actual_end_date",
		"created_at":         "created_at",
		"updated_at":         "updated_at",
	}
}

