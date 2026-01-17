package entities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	ProjectResourceStatusUnknown  = 0
	ProjectResourceStatusInactive = 1
	ProjectResourceStatusActive   = 2
)

var (
	ErrProjectResourceInvalidProjectID       = errors.New("project resource must have a valid project ID")
	ErrProjectResourceInvalidHumanResourceID = errors.New("project resource must have a valid human resource ID")
	ErrProjectResourceInvalidStatus          = errors.New("project resource status must be 1 (inactive) or 2 (active)")
	ErrProjectResourceInvalidAllocation      = errors.New("allocation percentage must be between 0 and 100")
	ErrProjectResourceInvalidDates           = errors.New("project resource end date must be after start date")

	ProjectResourceAllowedSortField = map[string]string{
		"id":                "id",
		"project_id":        "project_id",
		"human_resource_id": "human_resource_id",
		"role":              "role",
		"allocation":        "allocation",
		"start_date":        "start_date",
		"end_date":          "end_date",
		"status":            "status",
		"created_at":        "created_at",
		"updated_at":        "updated_at",
	}
)

// ProjectResource represents the allocation of a human resource to a project
type ProjectResource struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	ProjectID       uint       `gorm:"not null;index;uniqueIndex:idx_project_human_resource" json:"project_id"`
	HumanResourceID uint       `gorm:"not null;index;uniqueIndex:idx_project_human_resource" json:"human_resource_id"`
	Role            string     `gorm:"" json:"role"`                  // Role in the project (e.g., "Developer", "Tech Lead", "QA")
	Allocation      float64    `gorm:"default:100" json:"allocation"` // Allocation percentage (0-100)
	StartDate       *time.Time `gorm:"" json:"start_date"`            // When the resource starts on the project
	EndDate         *time.Time `gorm:"" json:"end_date"`              // When the resource ends on the project
	Notes           string     `gorm:"type:text" json:"notes"`        // Additional notes
	Status          uint       `gorm:"not null;default:2" json:"status"`
	CreatedAt       time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime:milli" json:"updated_at"`

	// Relationships
	Project       *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	HumanResource *HumanResource `gorm:"foreignKey:HumanResourceID" json:"human_resource,omitempty"`
}

// TableName returns the table name for the project resource entity
func (ProjectResource) TableName() string {
	return "project_resources"
}

// IsActive returns true if the project resource allocation is active
func (pr *ProjectResource) IsActive() bool {
	return pr.Status == ProjectResourceStatusActive
}

// Validate validates the project resource fields
func (pr *ProjectResource) Validate() error {
	// Validate required fields
	if pr.ProjectID == 0 {
		return ErrProjectResourceInvalidProjectID
	}

	if pr.HumanResourceID == 0 {
		return ErrProjectResourceInvalidHumanResourceID
	}

	// Validate allocation percentage
	if pr.Allocation < 0 || pr.Allocation > 100 {
		return ErrProjectResourceInvalidAllocation
	}

	// Validate dates
	if pr.StartDate != nil && pr.EndDate != nil {
		if pr.EndDate.Before(*pr.StartDate) {
			return ErrProjectResourceInvalidDates
		}
	}

	// Validate status
	if err := pr.validateStatus(); err != nil {
		return err
	}

	return nil
}

func (pr *ProjectResource) validateStatus() error {
	switch pr.Status {
	case ProjectResourceStatusActive, ProjectResourceStatusInactive:
		return nil
	}
	return ErrProjectResourceInvalidStatus
}

// BeforeCreate is a GORM hook that runs before creating a project resource
func (pr *ProjectResource) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not valid
	if err := pr.validateStatus(); err != nil {
		pr.Status = ProjectResourceStatusActive
	}

	return pr.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a project resource
func (pr *ProjectResource) BeforeUpdate(tx *gorm.DB) error {
	return pr.Validate()
}

type ProjectResourceQueryParams struct {
	ID_In              []uint     `json:"id_in"`
	ProjectID          uint       `json:"project_id"`
	ProjectID_In       []uint     `json:"project_id_in"`
	HumanResourceID    uint       `json:"human_resource_id"`
	HumanResourceID_In []uint     `json:"human_resource_id_in"`
	Role               string     `json:"role"`
	Role_Like          string     `json:"role_like"`
	Allocation_Gte     *float64   `json:"allocation_gte"`
	Allocation_Lte     *float64   `json:"allocation_lte"`
	Status             uint       `json:"status"`
	Status_In          []uint     `json:"status_in"`
	StartDate_Gte      *time.Time `json:"start_date_gte"`
	StartDate_Lte      *time.Time `json:"start_date_lte"`
	EndDate_Gte        *time.Time `json:"end_date_gte"`
	EndDate_Lte        *time.Time `json:"end_date_lte"`
	CreatedAt_Gte      *time.Time `json:"created_at_gte"`
	CreatedAt_Lte      *time.Time `json:"created_at_lte"`
	UpdatedAt_Gte      *time.Time `json:"updated_at_gte"`
	UpdatedAt_Lte      *time.Time `json:"updated_at_lte"`
	*QueryParams
}

// ProjectResourceListResponse represents the response for GetProjectResources
type ProjectResourceListResponse struct {
	Data  []*ProjectResource `json:"data"`
	Total int64              `json:"total"`
}
