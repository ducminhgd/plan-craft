package entities

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Role level constants
const (
	RoleLevelUnknown  = 0
	RoleLevelJunior   = 1
	RoleLevelMid      = 2
	RoleLevelSenior   = 3
	RoleLevelLead     = 4
	RoleLevelManager  = 5
	RoleLevelDirector = 6
	RoleLevelVP       = 7
	RoleLevelCLevel   = 8
)

var (
	ErrProjectRoleNameRequired     = errors.New("project role name is required")
	ErrProjectRoleInvalidProjectID = errors.New("project role must belong to a project")
	ErrProjectRoleInvalidLevel     = errors.New("project role level must be 1 (junior), 2 (mid), 3 (senior), or 4 (lead)")
	ErrProjectRoleInvalidHeadcount = errors.New("project role headcount must be non-negative")

	ProjectRoleAllowedSortField = map[string]string{
		"id":         "id",
		"name":       "name",
		"project_id": "project_id",
		"level":      "level",
		"headcount":  "headcount",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
)

// RoleLevelName returns the string name for a role level
func RoleLevelName(level uint) string {
	switch level {
	case RoleLevelJunior:
		return "Junior"
	case RoleLevelMid:
		return "Mid"
	case RoleLevelSenior:
		return "Senior"
	case RoleLevelLead:
		return "Lead"
	case RoleLevelManager:
		return "Manager"
	case RoleLevelDirector:
		return "Director"
	case RoleLevelVP:
		return "VP"
	case RoleLevelCLevel:
		return "C-Level"
	default:
		return "Unknown"
	}
}

// ProjectRole represents a role within a project with its level and headcount
type ProjectRole struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ProjectID uint      `gorm:"not null;index;uniqueIndex:idx_project_role_name_level,priority:1" json:"project_id"`
	Name      string    `gorm:"not null;uniqueIndex:idx_project_role_name_level,priority:2" json:"name"`
	Level     uint      `gorm:"not null;uniqueIndex:idx_project_role_name_level,priority:3" json:"level"`
	Headcount int       `gorm:"not null;default:1" json:"headcount"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`

	// Relationships
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName returns the table name for the project role entity
func (ProjectRole) TableName() string {
	return "project_roles"
}

// GetLevelName returns the human-readable name for this role's level
func (pr *ProjectRole) GetLevelName() string {
	return RoleLevelName(pr.Level)
}

// Validate validates the project role fields
func (pr *ProjectRole) Validate() error {
	// Trim whitespace from string fields
	pr.Name = strings.TrimSpace(pr.Name)

	// Validate required fields
	if pr.Name == "" {
		return ErrProjectRoleNameRequired
	}

	// Validate project ID
	if pr.ProjectID == 0 {
		return ErrProjectRoleInvalidProjectID
	}

	// Validate level
	if err := pr.validateLevel(); err != nil {
		return err
	}

	// Validate headcount
	if pr.Headcount < 0 {
		return ErrProjectRoleInvalidHeadcount
	}

	return nil
}

func (pr *ProjectRole) validateLevel() error {
	switch pr.Level {
	case RoleLevelJunior, RoleLevelMid, RoleLevelSenior, RoleLevelLead, RoleLevelManager, RoleLevelDirector, RoleLevelVP, RoleLevelCLevel:
		return nil
	}
	return ErrProjectRoleInvalidLevel
}

// BeforeCreate is a GORM hook that runs before creating a project role
func (pr *ProjectRole) BeforeCreate(tx *gorm.DB) error {
	// Set default level if not valid
	if err := pr.validateLevel(); err != nil {
		pr.Level = RoleLevelMid
	}

	// Set default headcount if not set
	if pr.Headcount == 0 {
		pr.Headcount = 1
	}

	return pr.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a project role
func (pr *ProjectRole) BeforeUpdate(tx *gorm.DB) error {
	return pr.Validate()
}

// ProjectRoleQueryParams defines query parameters for filtering project roles
type ProjectRoleQueryParams struct {
	ID_In         []uint     `json:"id_in"`
	ProjectID     uint       `json:"project_id"`
	ProjectID_In  []uint     `json:"project_id_in"`
	Name          string     `json:"name"`
	Name_Like     string     `json:"name_like"`
	Level         uint       `json:"level"`
	Level_In      []uint     `json:"level_in"`
	Headcount_Gte *int       `json:"headcount_gte"`
	Headcount_Lte *int       `json:"headcount_lte"`
	CreatedAt_Gte *time.Time `json:"created_at_gte"`
	CreatedAt_Lte *time.Time `json:"created_at_lte"`
	UpdatedAt_Gte *time.Time `json:"updated_at_gte"`
	UpdatedAt_Lte *time.Time `json:"updated_at_lte"`
	*QueryParams
}

// ProjectRoleListResponse represents the response for GetProjectRoles
type ProjectRoleListResponse struct {
	Data  []*ProjectRole `json:"data"`
	Total int64          `json:"total"`
}
