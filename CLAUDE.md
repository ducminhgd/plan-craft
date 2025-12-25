# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Plan Craft is a **desktop project management and estimation tool** built with Go and Wails, designed to help teams plan software projects with work breakdown structures, timeline estimation, resource planning, and cost estimation. It's a **local-first desktop application** targeting Windows, Linux, and macOS with a single binary distribution.

## Common Commands

### Development Workflow
```bash
make deps          # Download and tidy dependencies
make build         # Build to bin/plancraft (requires cmd/server/main.go)
make run           # Build and run the application
make dev           # Auto-reload dev mode (requires air: go install github.com/air-verse/air@latest)
make test          # Run all tests with coverage
```

### Database Management
```bash
make db-clean      # Remove SQLite database files (data/plancraft.db*)
migrate -path migrations -database "sqlite3://data/plancraft.db" up     # Apply migrations
migrate -path migrations -database "sqlite3://data/plancraft.db" down 1 # Rollback last migration
```

### Code Quality
```bash
make fmt           # Format code with go fmt
make lint          # Run golangci-lint (if installed)
make mocks         # Generate mocks with mockery (if installed)
```

### Running Single Tests
```bash
go test -v ./internal/entities -run TestProjectValidation
go test -v ./internal/entities/... -cover
```

## Architecture Overview

### Layered Architecture
Plan Craft follows a **clean architecture with clear separation of concerns**:

1. **Entities Layer** (`internal/entities/`) - GORM domain entities with business validation, pagination, filtering, and sorting capabilities
2. **Repository Layer** (`internal/repositories/`) - Database operations (planned, not yet implemented)
3. **Service Layer** (`internal/services/`) - Business logic (planned, not yet implemented)
4. **UI Layer** (`internal/ui/`) - Wails desktop UI (planned, not yet implemented)

### Current State (as of feature/init-with-ai branch)
- ✅ Database schema and migrations complete
- ✅ All GORM entities implemented with 90% test coverage
- ✅ Database initialization and configuration complete
- 🚧 Application entry point (`cmd/`) exists but empty
- ⏳ Repository, service, and UI layers not yet implemented

### Key Architectural Decisions

**No Soft Deletes**: There is no `deleted_at` field in any model. All deletes are permanent. This is an explicit design choice documented in `.ai/tech.md`.

**No Base Model**: Each model explicitly defines `ID`, `CreatedAt`, and `UpdatedAt` fields rather than embedding a base model.

**Global Database Instance**: Database connection is managed as a singleton via `requires.DB` (global GORM instance). Initialize with `requires.InitializeDatabase()` at application startup.

**Environment Configuration**: All configuration via environment variables using `sethvargo/go-envconfig`. See [.env.example](.env.example) for available settings. Load with `config.Cfg` singleton.

**SQLite First, Extensible Design**: Uses SQLite for simplicity but designed to support PostgreSQL/MySQL. Connection string and pragma settings in [config/config.go](config/config.go).

## Domain Model Architecture

### Core Entities and Relationships

The domain model consists of 10 interrelated entities representing project management concepts:

**Hierarchical Structures:**
- `Client` → `Project` (1:N) - Clients own multiple projects
- `Project` → `Milestone` (1:N) - Projects have phases
- `Project` → `Task` (1:N) - Projects have work items
- `Milestone` → `Task` (1:N) - Tasks can belong to milestones
- `Task` → `Task` (hierarchical) - Tasks can have subtasks via `ParentID` and `Level` fields

**Task Dependencies:**
- `Task` ↔ `Task` (M:N via `task_dependencies` table)
- Supports 4 dependency types: Finish-to-Start, Start-to-Start, Finish-to-Finish, Start-to-Finish
- Includes `LagDays` (positive = delay, negative = lead time)

**Resource Management (3-tier system):**
1. `Resource` - Global resource with default rates, capacity, skills
2. `ProjectRole` - Resource assigned to project with project-specific rates, role, level
3. `ResourceAllocation` - Time-based allocation (e.g., "50% allocated from Jan-Mar 2024")
4. `TaskAssignment` - Specific task assignment with estimated man-days and allocation percentage

**Cost Tracking:**
- `Cost` - Polymorphic cost tracking linked to Project, Milestone, or Task
- Supports labor, infrastructure, services, materials, equipment, overhead costs
- Fields: `CostableType`, `CostableID` for polymorphic associations

### Custom Work Time Configuration

Projects can override default work time units via nullable fields in [internal/entities/project.go](internal/entities/project.go):

```go
// Project model
HoursPerDay   *float64  // Override default 8 hours/day
DaysPerWeek   *float64  // Override default 5 days/week
DaysPerMonth  *float64  // Override default 20 days/month

// Helper methods automatically use project-specific or default values
func (p *Project) GetHoursPerDay() float64   // Returns project override or DefaultHoursPerDay
func (p *Project) EstimatedMonths() float64  // Converts EstimatedHours using project settings
```

**Default Constants** (in [internal/entities/models.go](internal/entities/models.go)):
- `DefaultHoursPerDay = 8.0`
- `DefaultDaysPerWeek = 5.0`
- `DefaultDaysPerMonth = 20.0`

### Hierarchical Task Structure (WBS)

Tasks support unlimited hierarchical depth via `Level` and `ParentID`:

```go
Level    int    // 1 = epic, 2 = task, 3 = subtask, 4 = sub-subtask, etc.
ParentID *uint  // Self-referencing foreign key
```

When working with tasks, always consider the hierarchical implications for cost/effort rollup.

### Custom GORM Types

**StringArray** - Stores string slices as JSON in SQLite:
```go
type StringArray []string  // Stored as JSON array in database
// Used for: Project.Tags, Resource.Skills
```

**JSONB** - Stores arbitrary JSON data:
```go
type JSONB map[string]interface{}  // Stored as JSON object
// Used for: Project.Metadata, Task.Metadata
```

See [internal/entities/entities.go](internal/entities/entities.go) for implementation.

### Query Patterns: Pagination, Filtering, and Sorting

The entities package provides a flexible query system for building database queries with pagination, filtering, and sorting. All patterns are defined in [internal/entities/entities.go](internal/entities/entities.go).

**Pagination**:
```go
// Create pagination with defaults (page 1, size 20, max 100)
pagination := entities.NewPagination(1, 20)

// Apply to GORM query
db = pagination.Apply(db)

// Get pagination metadata
totalPages := pagination.TotalPages()
hasNext := pagination.HasNext()
```

**Sorting**:
```go
// Create sort parameters
sort := entities.NewSort("created_at", entities.SortOrderDesc)

// Define allowed fields (whitelist for security)
allowedFields := map[string]string{
    "name": "name",
    "created_at": "created_at",
}

// Apply to GORM query
db = sort.Apply(db, allowedFields)
```

**Filtering**:
```go
// Create filter conditions
filters := entities.NewFilters([]entities.Filter{
    {Field: "status", Operator: entities.FilterOpEqual, Value: "active"},
    {Field: "created_at", Operator: entities.FilterOpGreaterThan, Value: time.Now().AddDate(0, -1, 0)},
}, "AND")

// Define allowed fields (whitelist for security)
allowedFields := map[string]string{
    "status": "status",
    "created_at": "created_at",
}

// Apply to GORM query
db = filters.Apply(db, allowedFields)
```

**Combined Query Parameters**:
```go
// Use QueryParams to combine all three
params := entities.NewQueryParams()
params.Pagination = entities.NewPagination(2, 50)
params.Sort = entities.NewSort("name", entities.SortOrderAsc)
params.Filters = entities.NewFilters([]entities.Filter{
    {Field: "type", Operator: entities.FilterOpEqual, Value: "product"},
}, "AND")

// Apply all at once
db = params.Apply(db, allowedSortFields, allowedFilterFields)
```

**Available Filter Operators**:
- `FilterOpEqual`, `FilterOpNotEqual` - Equality checks
- `FilterOpGreaterThan`, `FilterOpGreaterOrEqual`, `FilterOpLessThan`, `FilterOpLessOrEqual` - Comparisons
- `FilterOpLike`, `FilterOpNotLike` - Pattern matching
- `FilterOpIn`, `FilterOpNotIn` - List membership
- `FilterOpIsNull`, `FilterOpIsNotNull` - Null checks
- `FilterOpBetween` - Range queries (requires 2-element array value)
- `FilterOpContains` - Case-insensitive substring search

**Security Note**: Always use field whitelists (`allowedFields` maps) to prevent SQL injection and unauthorized field access.

## Database and Migrations

### Database Configuration

SQLite is configured with performance optimizations in [internal/requires/database.go](internal/requires/database.go):
- WAL (Write-Ahead Logging) mode for concurrency
- 64MB cache size (`cache_size=-64000`)
- Foreign keys enabled
- Incremental auto-vacuum
- Connection pooling (max 10 connections, max 5 idle)

### Migration Files

Migrations use golang-migrate with format: `{version}_{description}.{up|down}.sql`

**Current migration**: [migrations/000001_initial_schema.up.sql](migrations/000001_initial_schema.up.sql)
- Creates all 10 tables with foreign keys, indexes, CHECK constraints
- Includes triggers for automatic `updated_at` timestamps
- Comprehensive schema documentation in [migrations/README.md](migrations/README.md)

**Important**: Never modify existing migrations. Create new migration files for schema changes.

### Creating New Migrations

```bash
migrate create -ext sql -dir migrations -seq description_of_change
# Edit both .up.sql and .down.sql files
# Test locally before committing
```

## Model Validation and Enums

All enums have corresponding validation functions in [internal/entities/models.go](internal/entities/models.go):

```go
// Enums
type ProjectType string       // product, service, internal, consulting, research, maintenance
type TaskStatus string        // not_started, in_progress, on_hold, completed, cancelled
type DependencyType string    // finish_to_start, start_to_start, finish_to_finish, start_to_finish
type Priority string          // low, medium, high, critical
type CostType string          // labor, material, equipment, overhead, infrastructure, service, other
type RateType string          // hourly, daily, monthly, fixed

// Validation functions
func IsValidProjectType(pt ProjectType) bool
func IsValidTaskStatus(ts TaskStatus) bool
func IsValidDependencyType(dt DependencyType) bool
// ... etc
```

**GORM BeforeCreate/BeforeUpdate hooks** validate enums and required fields. See individual model files for validation logic.

## Testing Practices

### Current Test Coverage
- Entities package: 90% coverage
- 8 test files: `*_test.go` in [internal/entities/](internal/entities/)
- Uses `github.com/stretchr/testify/assert` for assertions

### Test Patterns Used

**Table-Driven Tests** (preferred pattern):
```go
func TestIsValidProjectType(t *testing.T) {
    tests := []struct {
        name string
        pt   ProjectType
        want bool
    }{
        {"Valid: product", ProjectTypeProduct, true},
        {"Invalid: unknown", ProjectType("unknown"), false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := IsValidProjectType(tt.pt); got != tt.want {
                t.Errorf("IsValidProjectType(%v) = %v, want %v", tt.pt, got, tt.want)
            }
        })
    }
}
```

**GORM Hook Testing**: Tests validate that BeforeCreate/BeforeUpdate hooks reject invalid data.

### Running Tests

```bash
go test ./...                          # All tests
go test -v ./internal/entities/...     # Verbose output for entities
go test -cover ./...                   # With coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out  # HTML coverage
```

## Technology Stack

### Core Dependencies (from go.mod)
- Go 1.23+
- GORM v1.30.0 (`gorm.io/gorm`)
- SQLite driver v1.6.0 (`gorm.io/driver/sqlite`)
- Environment config v1.3.0 (`github.com/sethvargo/go-envconfig`)
- Testify v1.9.0 (`github.com/stretchr/testify`)

### Planned Technologies (from .ai/tech.md)
- UI Framework: Wails (https://github.com/wailsapp/wails)
- Web Framework: go-chi (for API layer)
- Logging: Uber Zap (replacing current `log/slog`)
- Cache: Redis (optional)
- CI/CD: GitHub Actions

## Development Roadmap

The project follows a phased delivery approach:

**Version 1 (Desktop App):**
- v1.0: Project and work items management ← **Current focus**
- v1.1: Timeline estimation and critical path
- v1.2: Resource planning and allocation
- v1.3: Cost estimation and tracking

**Version 2:** Web application with REST API (future)

See [.ai/product-features.md](.ai/product-features.md) for detailed feature specifications.

## Important Files Reference

### Configuration & Initialization
- [config/config.go](config/config.go) - Environment-based configuration loading
- [internal/requires/database.go](internal/requires/database.go) - Database initialization and health checks
- [.env.example](.env.example) - Environment variable template

### Entities (all in internal/entities/)
- [entities.go](internal/entities/entities.go) - Enums, constants, validation functions, custom types, pagination, filtering, sorting
- [project.go](internal/entities/project.go) - Project entity with custom work time support
- [task.go](internal/entities/task.go) - Task entity with hierarchical WBS
- [resource.go](internal/entities/resource.go) - Resource, ProjectRole, ResourceAllocation, TaskAssignment
- [cost.go](internal/entities/cost.go) - Polymorphic cost tracking
- [milestone.go](internal/entities/milestone.go), [client.go](internal/entities/client.go)

### Database
- [migrations/000001_initial_schema.up.sql](migrations/000001_initial_schema.up.sql) - Complete database schema
- [migrations/README.md](migrations/README.md) - Migration documentation and best practices

### Build & Development
- [Makefile](Makefile) - All build, test, and development commands
- [go.mod](go.mod) - Dependency management

## AI Assistant Guidelines (from .ai/instructions.md)

When working with this codebase:

1. **No Markdown Files**: Don't generate additional markdown documentation unless explicitly requested
2. **Package Index Files**: The base file in a package should be named the same as the package (e.g., `entities/entities.go`)
3. **No Soft Deletes**: Never add `deleted_at` fields to entities
4. **No Base Model**: Don't create or use a base model with embedded fields
5. **Refer to Product Specs**: Check [.ai/product-features.md](.ai/product-features.md) for feature requirements
6. **Refer to Tech Specs**: Check [.ai/tech.md](.ai/tech.md) for technology choices and architecture decisions
