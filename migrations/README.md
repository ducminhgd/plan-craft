# Database Migrations

This directory contains database migration files for Plan Craft using [golang-migrate](https://github.com/golang-migrate/migrate).

## Migration Files

Migrations follow the naming convention: `{version}_{description}.{up|down}.sql`

### Current Migrations

- **000001_initial_schema** - Creates the complete database schema
  - **Tables created (in order):**
    1. `clients` - Client/customer management
    2. `projects` - Project management with metadata, timeline, and costs
    3. `milestones` - Project milestones and phases
    4. `tasks` - Work breakdown structure (WBS) with hierarchical tasks
    5. `task_dependencies` - Task dependencies (finish-to-start, etc.)
    6. `resources` - Resource/team member management
    7. `project_roles` - Resource role assignments in projects
    8. `resource_allocations` - Time-based resource allocation with percentages (requirement 4.3)
    9. `task_assignments` - Resource assignments to specific tasks
    10. `costs` - Cost tracking (labor, infrastructure, services, etc.)
  - **Includes:**
    - All foreign keys with CASCADE constraints
    - CHECK constraints for data validation
    - Comprehensive indexes for query performance
    - Triggers for automatic `updated_at` timestamp updates

## Prerequisites

Install golang-migrate:

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Windows
scoop install migrate

# Or using Go
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Usage

### Apply Migrations

```bash
# Apply all pending migrations
migrate -path migrations -database "sqlite3://data/plancraft.db" up

# Apply next N migrations
migrate -path migrations -database "sqlite3://data/plancraft.db" up 1
```

### Rollback Migrations

```bash
# Rollback last migration
migrate -path migrations -database "sqlite3://data/plancraft.db" down 1

# Rollback all migrations
migrate -path migrations -database "sqlite3://data/plancraft.db" down
```

### Check Migration Status

```bash
# Check current version
migrate -path migrations -database "sqlite3://data/plancraft.db" version

# Force to specific version (use with caution!)
migrate -path migrations -database "sqlite3://data/plancraft.db" force {version}
```

## Creating New Migrations

```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq {description}
```

This will create two files:
- `{version}_{description}.up.sql` - Contains changes to apply
- `{version}_{description}.down.sql` - Contains changes to rollback

## Best Practices

1. **Always create both up and down migrations** - This allows rollbacks
2. **Test migrations on a copy of production data** before deploying
3. **Never modify existing migrations** - Create new ones instead
4. **Keep migrations atomic** - One migration should do one thing
5. **Use transactions** where possible for data safety
6. **Add constraints at the database level** - Don't rely only on application validation

## Migration Checklist

Before running migrations in production:

- Migration has been tested on development environment
- Migration has been tested on staging environment with production-like data
- Both up and down migrations have been tested
- Migration is idempotent (can be run multiple times safely)
- Backup of production database exists
- Rollback plan is documented
- Team has been notified of planned migration

## Troubleshooting

### Dirty State

If a migration fails midway, the database may be in a "dirty" state:

```bash
# Check version (will show "dirty" flag)
migrate -path migrations -database "sqlite3://data/plancraft.db" version

# Force to a specific version (clears dirty flag)
migrate -path migrations -database "sqlite3://data/plancraft.db" force {version}

# Then manually fix the database and retry
```

### Clean Slate (Development Only)

```bash
# Drop all tables and re-migrate
rm data/plancraft.db
migrate -path migrations -database "sqlite3://data/plancraft.db" up
```

## SQLite-Specific Notes

1. **Foreign Keys**: Enabled via `_foreign_keys=on` in connection string
2. **ALTER TABLE Limitations**: SQLite has limited ALTER TABLE support
   - Cannot drop columns (requires table recreation)
   - Cannot modify column constraints
3. **Triggers**: Used for `updated_at` timestamps
4. **AUTOINCREMENT**: Used for primary keys

## Database Schema

The complete schema includes 10 tables:

| Table | Purpose | Key Features |
|-------|---------|--------------|
| `clients` | Client/customer management | Active status, contact info |
| `projects` | Project tracking | Timeline, budget, work config, metadata |
| `milestones` | Project phases | Timeline, cost, progress tracking |
| `tasks` | Work breakdown structure | Hierarchical (epics, tasks, subtasks), dependencies |
| `task_dependencies` | Task relationships | Finish-to-start, lag/lead times |
| `resources` | Team members | Skills, capacity, default rates |
| `project_roles` | Role in projects | Project-specific rates, capacity, timeline |
| `resource_allocations` | Time-based allocation | Percentage allocation per role per time range |
| `task_assignments` | Task-level assignments | Man-days estimation, allocation % |
| `costs` | Cost tracking | Labor, infrastructure, services, materials |

### Relationships

- Projects → Milestones (1:N)
- Projects → Tasks (1:N)
- Milestones → Tasks (1:N)
- Tasks → Subtasks (1:N, hierarchical)
- Tasks → Task Dependencies (M:N)
- Resources → Project Roles (1:N)
- Project Roles → Resource Allocations (1:N)
- Project Roles → Task Assignments (1:N)
- Projects/Milestones/Tasks → Costs (1:N)

See [internal/entities/](../internal/entities/) for Go model definitions.
