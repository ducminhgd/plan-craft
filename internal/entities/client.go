package entities

import (
	"errors"
	"strings"
	"time"

	"github.com/ducminhgd/plan-craft/pkg/x"
	"gorm.io/gorm"
)

const (
	ClientStatusActive   = 1
	ClientStatusInactive = 0
)

var (
	ErrClientNameRequired  = errors.New("client name is required")
	ErrClientEmailRequired = errors.New("client email is required")
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrClientInvalidStatus = errors.New("client status must be 0 (inactive) or 1 (active)")
)

// Client represents a client entity
type Client struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Email         string    `gorm:"not null;size:255" json:"email"`
	Phone         string    `gorm:"size:50" json:"phone"`
	Address       string    `gorm:"type:text" json:"address"`
	ContactPerson string    `gorm:"" json:"contact_person"`
	Notes         string    `gorm:"type:text" json:"notes"`
	Status        uint      `gorm:"not null;default:1" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

// TableName returns the table name for the client entity
func (Client) TableName() string {
	return "clients"
}

// IsActive returns true if the client is active
func (c *Client) IsActive() bool {
	return c.Status == ClientStatusActive
}

// Validate validates the client fields
func (c *Client) Validate() error {
	// Trim whitespace from string fields
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	c.Phone = strings.TrimSpace(c.Phone)
	c.ContactPerson = strings.TrimSpace(c.ContactPerson)

	// Validate required fields
	if c.Name == "" {
		return ErrClientNameRequired
	}

	if c.Email == "" {
		return ErrClientEmailRequired
	}

	// Validate email format
	if !x.IsValidEmail(c.Email) {
		return ErrInvalidEmail
	}

	// Validate status
	if err := c.validateStatus(); err != nil {
		return err
	}

	return nil
}

func (c *Client) validateStatus() error {
	switch c.Status {
	case ClientStatusActive, ClientStatusInactive:
		return nil
	}
	return ErrClientInvalidStatus
}

// BeforeCreate is a GORM hook that runs before creating a client
func (c *Client) BeforeCreate(tx *gorm.DB) error {
	// Set default status if not valid
	if err := c.validateStatus(); err != nil {
		c.Status = ClientStatusActive
	}

	return c.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a client
func (c *Client) BeforeUpdate(tx *gorm.DB) error {
	return c.Validate()
}
