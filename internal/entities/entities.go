package entities

import (
	"fmt"

	"gorm.io/gorm"
)

// Pagination defines pagination parameters for queries
type Pagination struct {
	Page     int `json:"page"`      // Current page number (1-indexed)
	PageSize int `json:"page_size"` // Number of items per page
	Total    int `json:"total"`     // Total number of items
}

// NewPagination creates a new Pagination with defaults
func NewPagination(page, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20 // Default page size
	}
	if pageSize > 100 {
		pageSize = 100 // Maximum page size
	}
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset calculates the offset for database queries
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// TotalPages calculates the total number of pages
func (p *Pagination) TotalPages() int {
	if p.Total == 0 || p.PageSize == 0 {
		return 0
	}
	pages := p.Total / p.PageSize
	if p.Total%p.PageSize > 0 {
		pages++
	}
	return pages
}

// HasNext returns true if there is a next page
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages()
}

// HasPrev returns true if there is a previous page
func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

// Apply applies pagination to a GORM query
func (p *Pagination) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(p.Offset()).Limit(p.PageSize)
}

// SortOrder represents sort direction
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// IsValid checks if the sort order is valid
func (s SortOrder) IsValid() bool {
	return s == SortOrderAsc || s == SortOrderDesc
}

// String returns the string representation
func (s SortOrder) String() string {
	return string(s)
}

// Sort defines sorting parameters
type Sort struct {
	Field string    `json:"field"` // Field name to sort by
	Order SortOrder `json:"order"` // Sort order (asc/desc)
}

// NewSort creates a new Sort with defaults
func NewSort(field string, order SortOrder) *Sort {
	if order == "" {
		order = SortOrderAsc
	}
	return &Sort{
		Field: field,
		Order: order,
	}
}

// Apply applies sorting to a GORM query
// allowedFields is a whitelist of fields that can be sorted
func (s *Sort) Apply(db *gorm.DB, allowedFields map[string]string) *gorm.DB {
	if s.Field == "" {
		return db
	}

	// Validate field against whitelist
	dbField, ok := allowedFields[s.Field]
	if !ok {
		return db // Ignore invalid fields
	}

	// Validate order
	if !s.Order.IsValid() {
		s.Order = SortOrderAsc
	}

	orderClause := fmt.Sprintf("%s %s", dbField, s.Order)
	return db.Order(orderClause)
}

// QueryParams combines pagination, sorting, and filtering
type QueryParams struct {
	Pagination *Pagination `json:"pagination,omitempty"`
	Sorts      []*Sort     `json:"sorts,omitempty"`
}

// NewQueryParams creates a new QueryParams with defaults
func NewQueryParams() *QueryParams {
	return &QueryParams{
		Pagination: NewPagination(1, 20),
		Sorts:      []*Sort{},
	}
}

// Enums and constants

// ProjectType represents the type of project
type ProjectType string

const (
	ProjectTypeProduct     ProjectType = "product"
	ProjectTypeService     ProjectType = "service"
	ProjectTypeInternal    ProjectType = "internal"
	ProjectTypeConsulting  ProjectType = "consulting"
	ProjectTypeResearch    ProjectType = "research"
	ProjectTypeMaintenance ProjectType = "maintenance"
)

// TaskStatus represents the status of a task or project
type TaskStatus string

const (
	TaskStatusNotStarted TaskStatus = "not_started"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusOnHold     TaskStatus = "on_hold"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// DependencyType represents the type of task dependency
type DependencyType string

const (
	DependencyFinishToStart  DependencyType = "finish_to_start"  // Task B starts after Task A finishes
	DependencyStartToStart   DependencyType = "start_to_start"   // Task B starts when Task A starts
	DependencyFinishToFinish DependencyType = "finish_to_finish" // Task B finishes when Task A finishes
	DependencyStartToFinish  DependencyType = "start_to_finish"  // Task B finishes when Task A starts
)

// Priority represents task priority
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// CostType represents the type of cost
type CostType string

const (
	CostTypeLabor          CostType = "labor"
	CostTypeMaterial       CostType = "material"
	CostTypeEquipment      CostType = "equipment"
	CostTypeOverhead       CostType = "overhead"
	CostTypeInfrastructure CostType = "infrastructure" // Cloud, hosting, servers, etc.
	CostTypeService        CostType = "service"        // Third-party services, SaaS, APIs, etc.
	CostTypeOther          CostType = "other"
)

// RateType represents how a resource is billed
type RateType string

const (
	RateTypeHourly  RateType = "hourly"
	RateTypeDaily   RateType = "daily"
	RateTypeMonthly RateType = "monthly"
	RateTypeFixed   RateType = "fixed"
)

// Helper functions for validation

// IsValidProjectType checks if the project type is valid
func IsValidProjectType(pt ProjectType) bool {
	switch pt {
	case ProjectTypeProduct, ProjectTypeService, ProjectTypeInternal,
		ProjectTypeConsulting, ProjectTypeResearch, ProjectTypeMaintenance:
		return true
	}
	return false
}

// IsValidTaskStatus checks if the task status is valid
func IsValidTaskStatus(ts TaskStatus) bool {
	switch ts {
	case TaskStatusNotStarted, TaskStatusInProgress, TaskStatusOnHold,
		TaskStatusCompleted, TaskStatusCancelled:
		return true
	}
	return false
}

// IsValidDependencyType checks if the dependency type is valid
func IsValidDependencyType(dt DependencyType) bool {
	switch dt {
	case DependencyFinishToStart, DependencyStartToStart,
		DependencyFinishToFinish, DependencyStartToFinish:
		return true
	}
	return false
}

// IsValidPriority checks if the priority is valid
func IsValidPriority(p Priority) bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	}
	return false
}

// IsValidCostType checks if the cost type is valid
func IsValidCostType(ct CostType) bool {
	switch ct {
	case CostTypeLabor, CostTypeMaterial, CostTypeEquipment,
		CostTypeOverhead, CostTypeInfrastructure, CostTypeService, CostTypeOther:
		return true
	}
	return false
}

// IsValidRateType checks if the rate type is valid
func IsValidRateType(rt RateType) bool {
	switch rt {
	case RateTypeHourly, RateTypeDaily, RateTypeMonthly, RateTypeFixed:
		return true
	}
	return false
}

// Default unit conversion constants
// These are used when projects don't specify custom work time settings
const (
	DefaultHoursPerDay   = 8.0  // Standard working hours per day
	DefaultDaysPerWeek   = 5.0  // Standard working days per week
	DefaultDaysPerMonth  = 20.0 // Standard working days per month (approximately)
	DefaultHoursPerWeek  = DefaultHoursPerDay * DefaultDaysPerWeek
	DefaultHoursPerMonth = DefaultHoursPerDay * DefaultDaysPerMonth
)

// Legacy constants for backward compatibility
// Deprecated: Use DefaultHoursPerDay instead
const HoursPerDay = DefaultHoursPerDay

// Deprecated: Use DefaultDaysPerWeek instead
const DaysPerWeek = DefaultDaysPerWeek

// Deprecated: Use DefaultDaysPerMonth instead
const DaysPerMonth = DefaultDaysPerMonth

// Deprecated: Use DefaultHoursPerWeek instead
const HoursPerWeek = DefaultHoursPerWeek

// Deprecated: Use DefaultHoursPerMonth instead
const HoursPerMonth = DefaultHoursPerMonth

// Unit conversion helper functions using default constants

// HoursToDays converts hours to days (using default 8-hour workday)
func HoursToDays(hours float64) float64 {
	return hours / DefaultHoursPerDay
}

// DaysToHours converts days to hours (using default 8-hour workday)
func DaysToHours(days float64) float64 {
	return days * DefaultHoursPerDay
}

// HoursToMonths converts hours to man-months (using default 160 hours per month)
func HoursToMonths(hours float64) float64 {
	return hours / DefaultHoursPerMonth
}

// MonthsToHours converts man-months to hours (using default 160 hours per month)
func MonthsToHours(months float64) float64 {
	return months * DefaultHoursPerMonth
}

// DaysToMonths converts days to man-months (using default 20 days per month)
func DaysToMonths(days float64) float64 {
	return days / DefaultDaysPerMonth
}

// MonthsToDays converts man-months to days (using default 20 days per month)
func MonthsToDays(months float64) float64 {
	return months * DefaultDaysPerMonth
}

// Custom unit conversion helper functions
// These functions accept custom work time parameters for project-specific calculations

// HoursToDaysCustom converts hours to days using custom hours per day
func HoursToDaysCustom(hours, hoursPerDay float64) float64 {
	if hoursPerDay <= 0 {
		return hours / DefaultHoursPerDay
	}
	return hours / hoursPerDay
}

// DaysToHoursCustom converts days to hours using custom hours per day
func DaysToHoursCustom(days, hoursPerDay float64) float64 {
	if hoursPerDay <= 0 {
		return days * DefaultHoursPerDay
	}
	return days * hoursPerDay
}

// HoursToMonthsCustom converts hours to man-months using custom hours per day and days per month
func HoursToMonthsCustom(hours, hoursPerDay, daysPerMonth float64) float64 {
	if hoursPerDay <= 0 {
		hoursPerDay = DefaultHoursPerDay
	}
	if daysPerMonth <= 0 {
		daysPerMonth = DefaultDaysPerMonth
	}
	hoursPerMonth := hoursPerDay * daysPerMonth
	return hours / hoursPerMonth
}

// MonthsToHoursCustom converts man-months to hours using custom hours per day and days per month
func MonthsToHoursCustom(months, hoursPerDay, daysPerMonth float64) float64 {
	if hoursPerDay <= 0 {
		hoursPerDay = DefaultHoursPerDay
	}
	if daysPerMonth <= 0 {
		daysPerMonth = DefaultDaysPerMonth
	}
	hoursPerMonth := hoursPerDay * daysPerMonth
	return months * hoursPerMonth
}

// DaysToMonthsCustom converts days to man-months using custom days per month
func DaysToMonthsCustom(days, daysPerMonth float64) float64 {
	if daysPerMonth <= 0 {
		return days / DefaultDaysPerMonth
	}
	return days / daysPerMonth
}

// MonthsToDaysCustom converts man-months to days using custom days per month
func MonthsToDaysCustom(months, daysPerMonth float64) float64 {
	if daysPerMonth <= 0 {
		return months * DefaultDaysPerMonth
	}
	return months * daysPerMonth
}
