# Models Package

This package contains all GORM data models for the Plan Craft application.

## Overview

The models package defines the database schema and business logic for:
- **Projects** - Main project entities with metadata, timeline, and cost tracking
- **Tasks** - Work items (placeholder for future implementation)
- **Resources** - Team members and roles (placeholder for future implementation)
- **Costs** - Cost tracking (placeholder for future implementation)

## Design Principles

- **No Base Model** - Each model defines its own `ID`, `CreatedAt`, and `UpdatedAt` fields
- **No Soft Deletes** - Records are permanently deleted (no `DeletedAt` field)
- **Explicit Fields** - All fields are explicitly defined for clarity

## Project Model

The `Project` model is the core entity representing a project in the system.

### Fields

#### Standard Fields
- `ID` (uint, primary key) - Auto-incrementing ID
- `CreatedAt` (time.Time) - Record creation timestamp
- `UpdatedAt` (time.Time) - Record update timestamp

#### Basic Information
- `Name` (string, required) - Project name
- `Code` (string, unique, optional) - Project code/identifier
- `Description` (text, optional) - Detailed project description

#### Classification
- `Type` (ProjectType, required, default: "product") - Type of project
  - `product`, `service`, `internal`, `consulting`, `research`, `maintenance`
- `Methodology` (Methodology, required, default: "agile") - Project methodology
  - `waterfall`, `agile`, `hybrid`, `kanban`, `scrum`
- `Status` (TaskStatus, required, default: "not_started") - Current status
  - `not_started`, `in_progress`, `on_hold`, `completed`, `cancelled`

#### Timeline
- `StartDate` (*time.Time, optional) - Project start date
- `TargetEndDate` (*time.Time, optional) - Planned end date
- `ActualEndDate` (*time.Time, optional) - Actual completion date

#### Estimation & Tracking
- `EstimatedEffort` (float64, hours) - Total estimated effort in hours
- `ActualEffort` (float64, hours) - Total actual effort spent
- `EstimatedCost` (decimal) - Total estimated cost
- `ActualCost` (decimal) - Total actual cost incurred
- `Currency` (string, default: "USD") - ISO 4217 currency code
- `Progress` (float64, 0-100) - Completion percentage

#### Metadata
- `Assumptions` (StringArray, JSON) - List of project assumptions
- `Constraints` (StringArray, JSON) - List of project constraints
- `Risks` (StringArray, JSON) - List of identified risks
- `Tags` (StringArray, JSON) - Tags for categorization
- `CustomFields` (JSONB) - Flexible custom fields

#### Ownership
- `OwnerID` (*uint, optional) - User ID of project owner
- `OwnerName` (string, optional) - Owner name (denormalized)
- `ClientName` (string, optional) - Client/customer name
- `Department` (string, optional) - Department/team

#### Relationships
- `Tasks` ([]Task) - Associated tasks (one-to-many)
- `Resources` ([]Resource) - Associated resources (one-to-many)
- `Costs` ([]Cost) - Associated costs (one-to-many)

### Methods

#### Validation
- `BeforeSave()` - GORM hook that validates the project before saving
  - Validates project type, methodology, and status
  - Ensures progress is between 0-100
  - Validates date logic (end date after start date)

#### Status Checks
- `IsActive() bool` - Returns true if status is "in_progress"
- `IsCompleted() bool` - Returns true if status is "completed"

#### Duration Calculations
- `Duration() *time.Duration` - Returns planned duration (TargetEndDate - StartDate)
- `ActualDuration() *time.Duration` - Returns actual duration (ActualEndDate - StartDate)

### Usage Examples

#### Create a Project

```go
import (
    "time"
    "github.com/ducminhgd/plan-craft/internal/models"
    "github.com/ducminhgd/plan-craft/internal/requires"
)

startDate := time.Now()
endDate := startDate.AddDate(0, 3, 0) // 3 months later

project := models.Project{
    Name:          "E-Commerce Platform",
    Code:          "ECOM-2024",
    Description:   "Build a modern e-commerce platform",
    Type:          models.ProjectTypeProduct,
    Methodology:   models.MethodologyAgile,
    Status:        models.TaskStatusInProgress,
    StartDate:     &startDate,
    TargetEndDate: &endDate,
    EstimatedEffort: 1600, // hours
    EstimatedCost:   150000,
    Currency:        "USD",
    Progress:        25.5,
    Assumptions: models.StringArray{
        "Team of 5 developers available",
        "Requirements are stable",
    },
    Tags: models.StringArray{"web", "ecommerce", "high-priority"},
    CustomFields: models.JSONB{
        "client_priority": "high",
        "contract_type":   "fixed-price",
    },
}

result := requires.DB.Create(&project)
if result.Error != nil {
    log.Fatal(result.Error)
}
```

#### Query Projects

```go
// Find all active projects
var activeProjects []models.Project
requires.DB.Where("status = ?", models.TaskStatusInProgress).Find(&activeProjects)

// Find projects by type
var productProjects []models.Project
requires.DB.Where("type = ?", models.ProjectTypeProduct).Find(&productProjects)

// Find with relationships
var project models.Project
requires.DB.Preload("Tasks").Preload("Resources").First(&project, 1)

// Complex query
var projects []models.Project
requires.DB.Where("progress > ? AND status = ?", 50, models.TaskStatusInProgress).
    Order("target_end_date ASC").
    Find(&projects)
```

#### Update a Project

```go
// Update single field
requires.DB.Model(&project).Update("Progress", 75.0)

// Update multiple fields
requires.DB.Model(&project).Updates(models.Project{
    Progress: 80.0,
    ActualEffort: 1200,
})

// Update with map
requires.DB.Model(&project).Updates(map[string]interface{}{
    "status":   models.TaskStatusCompleted,
    "progress": 100,
})
```

#### Delete a Project

```go
// Permanent delete
requires.DB.Delete(&project)

// Delete by ID
requires.DB.Delete(&models.Project{}, 1)
```

#### Use Helper Methods

```go
// Check status
if project.IsActive() {
    fmt.Println("Project is currently active")
}

if project.IsCompleted() {
    fmt.Println("Project is completed")
}

// Calculate duration
if duration := project.Duration(); duration != nil {
    fmt.Printf("Planned duration: %v days\n", duration.Hours()/24)
}

if actualDuration := project.ActualDuration(); actualDuration != nil {
    fmt.Printf("Actual duration: %v days\n", actualDuration.Hours()/24)
}
```

## Custom Types

### StringArray

A custom type for storing string arrays in JSON format.

```go
assumptions := models.StringArray{"Assumption 1", "Assumption 2"}
project.Assumptions = assumptions
```

### JSONB

A custom type for storing arbitrary JSON data.

```go
customFields := models.JSONB{
    "field1": "value1",
    "field2": 123,
    "nested": map[string]interface{}{
        "key": "value",
    },
}
project.CustomFields = customFields
```

## Enums and Constants

All enum types have validation helper functions:

- `IsValidProjectType(pt ProjectType) bool`
- `IsValidMethodology(m Methodology) bool`
- `IsValidTaskStatus(ts TaskStatus) bool`
- `IsValidDependencyType(dt DependencyType) bool`
- `IsValidPriority(p Priority) bool`
- `IsValidCostType(ct CostType) bool`
- `IsValidRateType(rt RateType) bool`

## Testing

Run tests with:

```bash
go test ./internal/models/...
go test -cover ./internal/models/...
```

Current test coverage: **90%**

## Database Schema

The Project model creates the following table:

```sql
CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    description TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'product',
    methodology VARCHAR(50) NOT NULL DEFAULT 'agile',
    status VARCHAR(50) NOT NULL DEFAULT 'not_started',
    start_date DATETIME,
    target_end_date DATETIME,
    actual_end_date DATETIME,
    estimated_effort DECIMAL(10,2) DEFAULT 0,
    actual_effort DECIMAL(10,2) DEFAULT 0,
    estimated_cost DECIMAL(15,2) DEFAULT 0,
    actual_cost DECIMAL(15,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    progress DECIMAL(5,2) DEFAULT 0,
    assumptions JSON,
    constraints JSON,
    risks JSON,
    tags JSON,
    custom_fields JSON,
    owner_id INTEGER,
    owner_name VARCHAR(255),
    client_name VARCHAR(255),
    department VARCHAR(100)
);

CREATE INDEX idx_projects_name ON projects(name);
CREATE INDEX idx_projects_code ON projects(code);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_owner_id ON projects(owner_id);
```

## Best Practices

1. **Always validate** - Use the validation helper functions before setting enum values
2. **Use pointers for optional dates** - Allows NULL values in the database
3. **Set defaults** - Use GORM tags to set sensible defaults
4. **Leverage hooks** - BeforeSave validates data automatically
5. **Use transactions** - For operations affecting multiple models
6. **Preload relationships** - Avoid N+1 queries
7. **Index frequently queried fields** - Status, dates, owner_id

## Future Enhancements

- [ ] Add audit trail (who created/updated)
- [ ] Add project templates
- [ ] Add project archiving
- [ ] Add project cloning
- [ ] Add budget tracking
- [ ] Add milestone tracking
- [ ] Add project health indicators

