package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Cost represents a cost entry for a project or task
type Cost struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Project, Milestone, and Task Association
	ProjectID   *uint `gorm:"index:idx_cost_project" json:"project_id,omitempty"`
	MilestoneID *uint `gorm:"index:idx_cost_milestone" json:"milestone_id,omitempty"`
	TaskID      *uint `gorm:"index:idx_cost_task" json:"task_id,omitempty"`

	// Cost Classification
	Type     CostType `gorm:"type:varchar(50);not null;index" json:"type"`           // labor, material, equipment, overhead, other
	Category string   `gorm:"type:varchar(100);index" json:"category,omitempty"`     // Custom category for grouping
	Name     string   `gorm:"type:varchar(255);not null" json:"name"`                // Description of the cost

	// Cost Details
	Amount       float64  `gorm:"type:decimal(15,2);not null" json:"amount"`              // Cost amount
	Currency     string   `gorm:"type:varchar(3);default:'USD'" json:"currency"`          // ISO 4217 currency code
	Quantity     float64  `gorm:"type:decimal(10,2);default:1" json:"quantity"`           // Quantity (for unit costs)
	UnitCost     float64  `gorm:"type:decimal(15,2);default:0" json:"unit_cost"`          // Cost per unit

	// For Labor Costs - linked to resources
	ResourceID    *uint    `gorm:"index:idx_cost_resource" json:"resource_id,omitempty"`       // If labor cost, which resource
	ProjectRoleID *uint    `gorm:"index:idx_cost_project_role" json:"project_role_id,omitempty"` // If labor cost, which project role
	RateType      RateType `gorm:"type:varchar(50)" json:"rate_type,omitempty"`                // hourly, daily, monthly, fixed
	Hours         float64  `gorm:"type:decimal(10,2);default:0" json:"hours,omitempty"`        // Hours worked (for hourly/daily rates)

	// Status
	IsEstimated bool       `gorm:"default:true;index" json:"is_estimated"` // true = estimated, false = actual
	Date        *time.Time `gorm:"type:datetime;index" json:"date,omitempty"` // Date when cost was incurred

	// Additional Information
	Notes string `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Project     *Project     `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Milestone   *Milestone   `gorm:"foreignKey:MilestoneID" json:"milestone,omitempty"`
	Task        *Task        `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Resource    *Resource    `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
	ProjectRole *ProjectRole `gorm:"foreignKey:ProjectRoleID" json:"project_role,omitempty"`
}

// TableName specifies the table name for Cost model
func (Cost) TableName() string {
	return "costs"
}

// BeforeSave is a GORM hook that runs before saving
func (c *Cost) BeforeSave(tx *gorm.DB) error {
	// Validate cost type
	if !IsValidCostType(c.Type) {
		return errors.New("invalid cost type")
	}

	// Validate rate type if provided
	if c.RateType != "" && !IsValidRateType(c.RateType) {
		return errors.New("invalid rate type")
	}

	// Validate that at least one of ProjectID, MilestoneID, or TaskID is set
	if c.ProjectID == nil && c.MilestoneID == nil && c.TaskID == nil {
		return errors.New("cost must be associated with a project, milestone, or task")
	}

	// Validate amount and quantity
	if c.Amount < 0 {
		return errors.New("amount cannot be negative")
	}

	if c.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	if c.UnitCost < 0 {
		return errors.New("unit cost cannot be negative")
	}

	if c.Hours < 0 {
		return errors.New("hours cannot be negative")
	}

	// Auto-calculate amount from quantity and unit cost if applicable
	if c.UnitCost > 0 && c.Quantity > 0 {
		calculatedAmount := c.UnitCost * c.Quantity
		// Only update if amount is not set or differs significantly
		if c.Amount == 0 || (c.Amount > 0 && c.Amount != calculatedAmount) {
			c.Amount = calculatedAmount
		}
	}

	// Labor costs should have resource or project role
	if c.Type == CostTypeLabor {
		if c.ResourceID == nil && c.ProjectRoleID == nil {
			return errors.New("labor costs must have a resource or project role")
		}
	}

	return nil
}

// IsLabor returns true if this is a labor cost
func (c *Cost) IsLabor() bool {
	return c.Type == CostTypeLabor
}

// IsActual returns true if this is an actual (not estimated) cost
func (c *Cost) IsActual() bool {
	return !c.IsEstimated
}

// TotalAmount returns the total cost amount (considering quantity)
func (c *Cost) TotalAmount() float64 {
	if c.UnitCost > 0 && c.Quantity > 0 {
		return c.UnitCost * c.Quantity
	}
	return c.Amount
}

// CalculateLaborCost calculates labor cost based on hours and rate
func (c *Cost) CalculateLaborCost(hourlyRate float64) float64 {
	if c.Type != CostTypeLabor || c.Hours <= 0 {
		return 0
	}

	switch c.RateType {
	case RateTypeHourly:
		return c.Hours * hourlyRate
	case RateTypeDaily:
		// Assuming 8 hours per day
		days := c.Hours / 8.0
		return days * hourlyRate
	case RateTypeMonthly:
		// Assuming 160 hours per month (20 days * 8 hours)
		months := c.Hours / 160.0
		return months * hourlyRate
	default:
		return c.Amount
	}
}
