# Requires Package

This package contains initialization and setup functions for application dependencies.

## Database

The database module provides GORM connection management for SQLite with optimized settings.

### Features

- **SQLite with GORM**: Direct GORM integration without additional abstraction layers
- **Optimized Configuration**: WAL mode, connection pooling, and performance tuning
- **Simple API**: Easy initialization and health checking
- **Global Access**: Singleton pattern for database access throughout the application

### Usage

#### Initialize Database

```go
import (
    "log"
    "github.com/ducminhgd/plan-craft/internal/requires"
)

func main() {
    // Initialize database
    if err := requires.InitializeDatabase(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer requires.CloseDatabase()

    // Use the global DB instance
    var count int64
    requires.DB.Model(&models.Project{}).Count(&count)
}
```

#### Health Check

```go
if err := requires.HealthCheck(); err != nil {
    log.Printf("Database health check failed: %v", err)
}
```

#### Direct GORM Access

```go
// Create
project := &models.Project{Name: "My Project"}
requires.DB.Create(project)

// Read
var project models.Project
requires.DB.First(&project, 1)

// Update
requires.DB.Model(&project).Update("Status", "completed")

// Delete (soft delete)
requires.DB.Delete(&project)

// Query with conditions
var projects []models.Project
requires.DB.Where("status = ?", "in_progress").Find(&projects)
```

### Configuration

Database configuration is loaded from environment variables via `config.Config`:

```bash
# Database file path (DSN)
DB_DSN=data/plancraft.db

# SQLite optimization parameters (all optional)
# If not set, the parameter will be omitted from the DSN
# The go-envconfig library provides defaults in config.DBConfig
DB_JOURNAL_MODE=WAL           # Journal mode (default: WAL)
DB_SYNCHRONOUS=NORMAL         # Synchronous mode (default: NORMAL)
DB_FOREIGN_KEYS=ON            # Foreign keys (default: ON)
DB_BUSY_TIMEOUT=5000          # Busy timeout in ms (default: 5000)
DB_CACHE_SIZE=-64000          # Cache size in KB (default: -64000 = 64MB)
DB_TEMP_STORE=MEMORY          # Temp store location (default: MEMORY)
DB_AUTO_VACUUM=INCREMENTAL    # Auto vacuum mode (default: INCREMENTAL)

# Log level (affects database query logging)
LOG_LEVEL=WARN  # Options: ERROR, WARN, INFO, DEBUG
```

**Note on Empty Parameters:**
- If a parameter is not set or empty, it will be **omitted** from the DSN
- This allows SQLite to use its own defaults for those parameters
- The `config.DBConfig` struct provides sensible defaults via struct tags
- You can override any parameter by setting the corresponding environment variable

### SQLite Optimizations

The database connection is configured with the following optimizations:

| Setting | Value | Purpose |
|---------|-------|---------|
| Journal Mode | WAL | Write-Ahead Logging for better concurrency |
| Synchronous | NORMAL | Balance between safety and performance |
| Foreign Keys | ON | Enforce referential integrity |
| Busy Timeout | 5000ms | Wait time when database is locked |
| Cache Size | 64MB | In-memory cache for better performance |
| Temp Store | MEMORY | Store temporary tables in memory |
| Auto Vacuum | INCREMENTAL | Reclaim space without blocking |

### Connection Pool Settings

Default connection pool configuration:

```go
MaxIdleConns:    5              // Maximum idle connections
MaxOpenConns:    10             // Maximum open connections
ConnMaxLifetime: 5 * time.Minute // Maximum connection lifetime
ConnMaxIdleTime: 5 * time.Minute // Maximum idle time
```

### GORM Configuration

The GORM instance is configured with:

- **SkipDefaultTransaction**: Disabled for better performance
- **PrepareStmt**: Enabled to cache prepared statements
- **NowFunc**: Uses UTC time for consistency
- **Logger**: Configured based on LOG_LEVEL environment variable

### Database Files

SQLite creates multiple files in WAL mode:

- `plancraft.db` - Main database file
- `plancraft.db-shm` - Shared memory file
- `plancraft.db-wal` - Write-ahead log file

**Important**: Don't delete the `-shm` and `-wal` files while the application is running.

### Error Handling

```go
if err := requires.InitializeDatabase(); err != nil {
    // Handle initialization error
    log.Fatal(err)
}

if err := requires.HealthCheck(); err != nil {
    // Handle health check failure
    log.Printf("Database unhealthy: %v", err)
}
```

### Best Practices

1. **Initialize Once**: Call `InitializeDatabase()` once at application startup
2. **Defer Close**: Always defer `CloseDatabase()` after initialization
3. **Use Context**: Pass context to GORM operations for cancellation support
4. **Check Health**: Periodically check database health in production
5. **Handle Errors**: Always check and handle database errors

### Example: Complete Setup

```go
package main

import (
    "log"
    "log/slog"
    
    "github.com/ducminhgd/plan-craft/internal/requires"
    "github.com/ducminhgd/plan-craft/internal/entities"
)

func main() {
    // Initialize database
    if err := requires.InitializeDatabase(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer requires.CloseDatabase()

    // Check health
    if err := requires.HealthCheck(); err != nil {
        log.Fatalf("Database health check failed: %v", err)
    }

    slog.Info("Database initialized and healthy")

    // Use database
    var projects []models.Project
    result := requires.DB.Find(&projects)
    if result.Error != nil {
        log.Fatalf("Failed to query projects: %v", result.Error)
    }

    slog.Info("Found projects", slog.Int("count", len(projects)))
}
```

### Switching to PostgreSQL or MySQL

To switch to a different database:

1. Update the driver import in `database.go`:
```go
import (
    "gorm.io/driver/postgres" // or "gorm.io/driver/mysql"
)
```

2. Update the `gorm.Open` call:
```go
// PostgreSQL
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{...})

// MySQL
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{...})
```

3. Update the DSN building logic for the new database type

4. Update environment variables and configuration

### Troubleshooting

#### Database is locked

**Cause**: Another process is accessing the database

**Solution**:
- Close any SQLite browser/viewer
- Check for long-running transactions
- The 5-second busy timeout should handle most cases

#### Cannot create database file

**Cause**: Permission issues or directory doesn't exist

**Solution**:
```bash
mkdir -p data
chmod 755 data
```

#### Connection pool exhausted

**Cause**: Too many concurrent operations

**Solution**:
- Increase `MaxOpenConns` in the configuration
- Ensure connections are properly released
- Use connection pooling wisely

## References

- [GORM Documentation](https://gorm.io/docs/)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [SQLite WAL Mode](https://www.sqlite.org/wal.html)

