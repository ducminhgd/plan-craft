package entities

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
	Type     CostType `gorm:"type:varchar(50);not null;index" json:"type"`       // labor, material, equipment, overhead, other
	Category string   `gorm:"type:varchar(100);index" json:"category,omitempty"` // Custom category for grouping
	Name     string   `gorm:"type:varchar(255);not null" json:"name"`            // Description of the cost

	// Cost Details
	Amount   float64 `gorm:"type:decimal(15,2);not null" json:"amount"`     // Cost amount
	Currency string  `gorm:"type:varchar(3);default:'USD'" json:"currency"` // ISO 4217 currency code
	Quantity float64 `gorm:"type:decimal(10,2);default:1" json:"quantity"`  // Quantity (for unit costs)
	UnitCost float64 `gorm:"type:decimal(15,2);default:0" json:"unit_cost"` // Cost per unit

	// For Labor Costs - linked to resources
	ResourceID    *uint    `gorm:"index:idx_cost_resource" json:"resource_id,omitempty"`         // If labor cost, which resource
	ProjectRoleID *uint    `gorm:"index:idx_cost_project_role" json:"project_role_id,omitempty"` // If labor cost, which project role
	RateType      RateType `gorm:"type:varchar(50)" json:"rate_type,omitempty"`                  // hourly, daily, monthly, fixed
	Hours         float64  `gorm:"type:decimal(10,2);default:0" json:"hours,omitempty"`          // Hours worked (for hourly/daily rates)

	// Status
	IsEstimated bool       `gorm:"default:true;index" json:"is_estimated"`    // true = estimated, false = actual
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

// CostQuery is the query for searching costs
type CostQuery struct {
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

	// TaskID is the task ID
	TaskID *uint `json:"task_id"`
	// TaskID_In is a list of task IDs to search for
	TaskID_In []uint `json:"task_id__in"`

	// ResourceID is the resource ID
	ResourceID *uint `json:"resource_id"`
	// ResourceID_In is a list of resource IDs to search for
	ResourceID_In []uint `json:"resource_id__in"`

	// ProjectRoleID is the project role ID
	ProjectRoleID *uint `json:"project_role_id"`
	// ProjectRoleID_In is a list of project role IDs to search for
	ProjectRoleID_In []uint `json:"project_role_id__in"`

	// Type is the cost type
	Type CostType `json:"type"`
	// Type_In is a list of cost types to search for
	Type_In []CostType `json:"type__in"`

	// Category is the cost category
	Category string `json:"category"`
	// Category_Like is the cost category to search for (case-insensitive)
	Category_Like string `json:"category__like"`

	// Name is the name of the cost
	Name string `json:"name"`
	// Name_Like is the name of the cost to search for (case-insensitive)
	Name_Like string `json:"name__like"`

	// Currency is the currency code
	Currency string `json:"currency"`
	// Currency_In is a list of currency codes to search for
	Currency_In []string `json:"currency__in"`

	// RateType is the rate type
	RateType RateType `json:"rate_type"`
	// RateType_In is a list of rate types to search for
	RateType_In []RateType `json:"rate_type__in"`

	// IsEstimated filters for estimated/actual costs
	IsEstimated *bool `json:"is_estimated"`

	// Amount_Gte is the minimum amount
	Amount_Gte *float64 `json:"amount__gte"`
	// Amount_Lte is the maximum amount
	Amount_Lte *float64 `json:"amount__lte"`

	// Quantity_Gte is the minimum quantity
	Quantity_Gte *float64 `json:"quantity__gte"`
	// Quantity_Lte is the maximum quantity
	Quantity_Lte *float64 `json:"quantity__lte"`

	// UnitCost_Gte is the minimum unit cost
	UnitCost_Gte *float64 `json:"unit_cost__gte"`
	// UnitCost_Lte is the maximum unit cost
	UnitCost_Lte *float64 `json:"unit_cost__lte"`

	// Hours_Gte is the minimum hours
	Hours_Gte *float64 `json:"hours__gte"`
	// Hours_Lte is the maximum hours
	Hours_Lte *float64 `json:"hours__lte"`

	// Date_Gte is the minimum cost date
	Date_Gte *time.Time `json:"date__gte"`
	// Date_Lte is the maximum cost date
	Date_Lte *time.Time `json:"date__lte"`

	// CreatedAt_Gte is the start time of the cost creation time to search for
	CreatedAt_Gte *time.Time `json:"created_at__gte"`
	// CreatedAt_Lte is the end time of the cost creation time to search for
	CreatedAt_Lte *time.Time `json:"created_at__lte"`

	// UpdatedAt_Gte is the start time of the cost update time to search for
	UpdatedAt_Gte *time.Time `json:"updated_at__gte"`
	// UpdatedAt_Lte is the end time of the cost update time to search for
	UpdatedAt_Lte *time.Time `json:"updated_at__lte"`

	// QueryParams holds pagination, sorting, and filtering options
	QueryParams `json:",inline"`
}

// AllowedSortFields returns the allowed fields for sorting
func (q *CostQuery) AllowedSortFields() map[string]string {
	return map[string]string{
		"id":          "id",
		"type":        "type",
		"category":    "category",
		"name":        "name",
		"amount":      "amount",
		"quantity":    "quantity",
		"unit_cost":   "unit_cost",
		"hours":       "hours",
		"date":        "date",
		"is_estimated": "is_estimated",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
	}
}

// AllowedFilterFields returns the allowed fields for filtering
func (q *CostQuery) AllowedFilterFields() map[string]string {
	return map[string]string{
		"id":              "id",
		"project_id":      "project_id",
		"milestone_id":    "milestone_id",
		"task_id":         "task_id",
		"resource_id":     "resource_id",
		"project_role_id": "project_role_id",
		"type":            "type",
		"category":        "category",
		"name":            "name",
		"currency":        "currency",
		"rate_type":       "rate_type",
		"is_estimated":    "is_estimated",
		"amount":          "amount",
		"quantity":        "quantity",
		"unit_cost":       "unit_cost",
		"hours":           "hours",
		"date":            "date",
		"created_at":      "created_at",
		"updated_at":      "updated_at",
	}
}
