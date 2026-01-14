package entities

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	HumanResourceStatusUnknown  = 0
	HumanResourceStatusInactive = 1
	HumanResourceStatusActive   = 2
)

var (
	ErrHumanResourceNameRequired  = errors.New("human resource name is required")
	ErrHumanResourceTitleRequired = errors.New("human resource title is required")
	ErrHumanResourceLevelRequired = errors.New("human resource level is required")
	ErrHumanResourceInvalidStatus = errors.New("human resource status must be 1 (inactive) or 2 (active)")

	HumanResourceAllowedSortField = map[string]string{
		"id":         "id",
		"name":       "name",
		"title":      "title",
		"level":      "level",
		"status":     "status",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
)

// HumanResource represents a human resource entity
type HumanResource struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Title     string    `gorm:"not null" json:"title"`
	Level     string    `gorm:"not null" json:"level"`
	Status    uint      `gorm:"not null;default:2" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

// TableName returns the table name for the human resource entity
func (HumanResource) TableName() string {
	return "human_resources"
}

// IsActive returns true if the human resource is active
func (hr *HumanResource) IsActive() bool {
	return hr.Status == HumanResourceStatusActive
}

// Validate validates the human resource fields
func (hr *HumanResource) Validate() error {
	// Trim whitespace from string fields
	hr.Name = strings.TrimSpace(hr.Name)
	hr.Title = strings.TrimSpace(hr.Title)
	hr.Level = strings.TrimSpace(hr.Level)

	// Validate required fields
	if hr.Name == "" {
		return ErrHumanResourceNameRequired
	}

	if hr.Title == "" {
		return ErrHumanResourceTitleRequired
	}

	if hr.Level == "" {
		return ErrHumanResourceLevelRequired
	}

	// Validate status
	if err := hr.validateStatus(); err != nil {
		return err
	}

	return nil
}

func (hr *HumanResource) validateStatus() error {
	switch hr.Status {
	case HumanResourceStatusActive, HumanResourceStatusInactive:
		return nil
	}
	return ErrHumanResourceInvalidStatus
}

// BeforeCreate is a GORM hook that runs before creating a human resource
func (hr *HumanResource) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not valid
	if err := hr.validateStatus(); err != nil {
		hr.Status = HumanResourceStatusActive
	}

	return hr.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a human resource
func (hr *HumanResource) BeforeUpdate(tx *gorm.DB) error {
	return hr.Validate()
}

type HumanResourceQueryParams struct {
	ID_In         []uint     `json:"id_in"`
	Name          string     `json:"name"`
	Name_Like     string     `json:"name_like"`
	Title         string     `json:"title"`
	Title_Like    string     `json:"title_like"`
	Level         string     `json:"level"`
	Level_Like    string     `json:"level_like"`
	Status        uint       `json:"status"`
	Status_In     []uint     `json:"status_in"`
	CreatedAt_Gte *time.Time `json:"created_at_gte"`
	CreatedAt_Lte *time.Time `json:"created_at_lte"`
	UpdatedAt_Gte *time.Time `json:"updated_at_gte"`
	UpdatedAt_Lte *time.Time `json:"updated_at_lte"`
	*QueryParams
}

// HumanResourceListResponse represents the response for GetHumanResources
type HumanResourceListResponse struct {
	Data  []*HumanResource `json:"data"`
	Total int64            `json:"total"`
}
