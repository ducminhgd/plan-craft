package entities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Client represents a client or customer in the system
type Client struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Basic Information
	Name  string `gorm:"type:varchar(255);not null;index" json:"name"`
	Email string `gorm:"type:varchar(255);index" json:"email,omitempty"`

	// Contact Information
	Phone   string `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Address string `gorm:"type:text" json:"address,omitempty"`

	// Status
	IsActive bool   `gorm:"default:true;index" json:"is_active"` // Whether client is active
	Notes    string `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Projects []Project `gorm:"foreignKey:ClientID;constraint:OnDelete:SET NULL" json:"projects,omitempty"`
}

// TableName specifies the table name for Client model
func (Client) TableName() string {
	return "clients"
}

// BeforeSave is a GORM hook that runs before saving
func (c *Client) BeforeSave(tx *gorm.DB) error {
	// Validate name is not empty
	if c.Name == "" {
		return errors.New("client name cannot be empty")
	}

	return nil
}

// IsActiveClient returns true if the client is active
func (c *Client) IsActiveClient() bool {
	return c.IsActive
}

// ClientQuery is the query for searching clients
type ClientQuery struct {
	// ID_In is a list of IDs to search for
	ID_In []uint `json:"id__in"`

	// Name is the name of the client
	Name string `json:"name"`
	// Name_Like is the name of the client to search for (case-insensitive)
	Name_Like string `json:"name__like"`

	// Email is the email of the client
	Email string `json:"email"`
	// Email_Like is the email of the client to search for (case-insensitive)
	Email_Like string `json:"email__like"`

	// Phone is the phone number
	Phone string `json:"phone"`
	// Phone_Like is the phone number to search for (case-insensitive)
	Phone_Like string `json:"phone__like"`

	// IsActive filters for active/inactive clients
	IsActive *bool `json:"is_active"`

	// CreatedAt_Gte is the start time of the client creation time to search for
	CreatedAt_Gte *time.Time `json:"created_at__gte"`
	// CreatedAt_Lte is the end time of the client creation time to search for
	CreatedAt_Lte *time.Time `json:"created_at__lte"`

	// UpdatedAt_Gte is the start time of the client update time to search for
	UpdatedAt_Gte *time.Time `json:"updated_at__gte"`
	// UpdatedAt_Lte is the end time of the client update time to search for
	UpdatedAt_Lte *time.Time `json:"updated_at__lte"`

	// QueryParams holds pagination, sorting, and filtering options
	QueryParams `json:",inline"`
}

// AllowedSortFields returns the allowed fields for sorting
func (q *ClientQuery) AllowedSortFields() map[string]string {
	return map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"is_active":  "is_active",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
}

// AllowedFilterFields returns the allowed fields for filtering
func (q *ClientQuery) AllowedFilterFields() map[string]string {
	return map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"phone":      "phone",
		"is_active":  "is_active",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
}
