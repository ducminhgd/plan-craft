# Database Configuration Examples

This document shows various ways to configure the SQLite database for different scenarios.

## Default Configuration (No Environment Variables)

When no environment variables are set, the application uses defaults from `config.DBConfig`:

```bash
# No configuration needed
go run main.go
```

**Resulting DSN:**
```
data/plancraft.db?_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON&_busy_timeout=5000&_cache_size=-64000&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL
```

## Minimal Configuration (Only DSN)

Set only the database path, let SQLite use its defaults for everything else:

```bash
export DB_DSN=/var/lib/plancraft/app.db
# Don't set any other DB_* variables
```

**Resulting DSN:**
```
/var/lib/plancraft/app.db?_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON&_busy_timeout=5000&_cache_size=-64000&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL
```

(Still uses defaults from config struct tags)

## Custom Configuration

Override specific parameters while keeping others at default:

```bash
export DB_DSN=data/plancraft.db
export DB_CACHE_SIZE=-128000      # 128MB cache instead of 64MB
export DB_BUSY_TIMEOUT=10000      # 10 seconds instead of 5
```

**Resulting DSN:**
```
data/plancraft.db?_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON&_busy_timeout=10000&_cache_size=-128000&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL
```

## In-Memory Database (Testing)

Use an in-memory database for testing:

```bash
export DB_DSN=:memory:
export DB_JOURNAL_MODE=MEMORY
```

**Resulting DSN:**
```
:memory:?_journal_mode=MEMORY
```

## Read-Heavy Workload

Optimize for read-heavy workloads:

```bash
export DB_DSN=data/plancraft.db
export DB_JOURNAL_MODE=WAL        # Better for concurrent reads
export DB_CACHE_SIZE=-256000      # 256MB cache
export DB_SYNCHRONOUS=NORMAL      # Balance safety/performance
export DB_TEMP_STORE=MEMORY       # Fast temp operations
```

**Resulting DSN:**
```
data/plancraft.db?_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON&_busy_timeout=5000&_cache_size=-256000&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL
```

## Write-Heavy Workload

Optimize for write-heavy workloads:

```bash
export DB_DSN=data/plancraft.db
export DB_JOURNAL_MODE=WAL        # Better write concurrency
export DB_SYNCHRONOUS=NORMAL      # Faster writes (still safe)
export DB_CACHE_SIZE=-128000      # 128MB cache
export DB_BUSY_TIMEOUT=15000      # Wait longer for locks
```

**Resulting DSN:**
```
data/plancraft.db?_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON&_busy_timeout=15000&_cache_size=-128000&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL
```

## Maximum Safety (Production)

Prioritize data safety over performance:

```bash
export DB_DSN=data/plancraft.db
export DB_JOURNAL_MODE=WAL
export DB_SYNCHRONOUS=FULL        # Maximum safety
export DB_FOREIGN_KEYS=ON         # Enforce constraints
export DB_BUSY_TIMEOUT=30000      # Wait up to 30 seconds
```

**Resulting DSN:**
```
data/plancraft.db?_journal_mode=WAL&_synchronous=FULL&_foreign_keys=ON&_busy_timeout=30000&_cache_size=-64000&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL
```

## Maximum Performance (Development)

Prioritize performance over safety (development only):

```bash
export DB_DSN=data/plancraft.db
export DB_JOURNAL_MODE=MEMORY     # Fastest (data loss on crash)
export DB_SYNCHRONOUS=OFF         # Fastest (risky)
export DB_CACHE_SIZE=-512000      # 512MB cache
export DB_TEMP_STORE=MEMORY
```

**Resulting DSN:**
```
data/plancraft.db?_journal_mode=MEMORY&_synchronous=OFF&_foreign_keys=ON&_busy_timeout=5000&_cache_size=-512000&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL
```

⚠️ **Warning:** This configuration risks data loss. Only use in development!

## Embedded/IoT Device (Low Memory)

Optimize for devices with limited memory:

```bash
export DB_DSN=data/plancraft.db
export DB_JOURNAL_MODE=WAL
export DB_CACHE_SIZE=-8000        # Only 8MB cache
export DB_TEMP_STORE=FILE         # Use disk for temp tables
export DB_SYNCHRONOUS=NORMAL
```

**Resulting DSN:**
```
data/plancraft.db?_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON&_busy_timeout=5000&_cache_size=-8000&_temp_store=FILE&_auto_vacuum=INCREMENTAL
```

## Docker Container

Configuration for containerized deployment:

```bash
# In docker-compose.yml or Dockerfile
environment:
  - DB_DSN=/data/plancraft.db
  - DB_JOURNAL_MODE=WAL
  - DB_SYNCHRONOUS=NORMAL
  - DB_CACHE_SIZE=-128000
  - DB_BUSY_TIMEOUT=10000
  - LOG_LEVEL=INFO
```

## Kubernetes ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: plancraft-config
data:
  DB_DSN: "/data/plancraft.db"
  DB_JOURNAL_MODE: "WAL"
  DB_SYNCHRONOUS: "NORMAL"
  DB_FOREIGN_KEYS: "ON"
  DB_BUSY_TIMEOUT: "10000"
  DB_CACHE_SIZE: "-128000"
  DB_TEMP_STORE: "MEMORY"
  DB_AUTO_VACUUM: "INCREMENTAL"
  LOG_LEVEL: "INFO"
```

## Parameter Reference

| Parameter | Values | Description |
|-----------|--------|-------------|
| `DB_DSN` | File path or `:memory:` | Database file location |
| `DB_JOURNAL_MODE` | `DELETE`, `TRUNCATE`, `PERSIST`, `MEMORY`, `WAL`, `OFF` | Journal mode |
| `DB_SYNCHRONOUS` | `OFF`, `NORMAL`, `FULL`, `EXTRA` | Synchronization level |
| `DB_FOREIGN_KEYS` | `ON`, `OFF` | Foreign key enforcement |
| `DB_BUSY_TIMEOUT` | Milliseconds (e.g., `5000`) | Lock wait timeout |
| `DB_CACHE_SIZE` | Negative KB (e.g., `-64000` = 64MB) | Page cache size |
| `DB_TEMP_STORE` | `DEFAULT`, `FILE`, `MEMORY` | Temp table storage |
| `DB_AUTO_VACUUM` | `NONE`, `FULL`, `INCREMENTAL` | Auto-vacuum mode |

## Best Practices

1. **Always set `DB_DSN`** - Don't rely on the default path in production
2. **Use WAL mode** - Better concurrency for most workloads
3. **Set appropriate cache size** - Based on available memory
4. **Tune busy timeout** - Based on expected lock contention
5. **Enable foreign keys** - Maintain data integrity
6. **Use NORMAL synchronous** - Good balance for most cases
7. **Monitor performance** - Adjust based on actual workload

## Testing Different Configurations

```bash
# Test with different configurations
DB_CACHE_SIZE=-32000 go test ./...
DB_JOURNAL_MODE=MEMORY go test ./...
DB_DSN=:memory: go test ./...
```

