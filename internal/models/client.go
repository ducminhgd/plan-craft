package models

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
