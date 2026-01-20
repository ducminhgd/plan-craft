# Plan Craft

A desktop project management and estimation tool built with Go and Wails, designed to help teams plan software projects with work breakdown structures, timeline estimation, resource planning, and cost estimation.

## Features

- **Project Management** - Create and manage projects with metadata, methodology tracking (Waterfall/Agile/Hybrid), and milestone planning
- **Work Breakdown Structure (WBS)** - Hierarchical task organization (epics → tasks → subtasks) with effort estimation
- **Task Dependencies** - Support for Finish-to-Start, Start-to-Start, Finish-to-Finish, and Start-to-Finish dependencies
- **Resource Planning** - Role definitions, resource assignments, and capacity tracking
- **Cost Estimation** - Track costs by category (labor, materials, equipment, overhead)

## Tech Stack

| Layer | Technology |
|-------|------------|
| Desktop Framework | [Wails v2](https://wails.io/) |
| Backend | Go 1.23+ |
| Frontend | React 18 + TypeScript + Vite |
| UI Library | Ant Design |
| Database | SQLite (with GORM) |
| Migrations | golang-migrate |

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.23+** - [Download](https://go.dev/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **Wails CLI** - Install after Go is set up:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```
- **golang-migrate** (for database migrations):
  ```bash
  go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

### Platform-specific Requirements

Wails requires additional dependencies depending on your OS. Run this command to check:

```bash
wails doctor
```

See the [Wails installation guide](https://wails.io/docs/gettingstarted/installation) for platform-specific requirements.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/ducminhgd/plan-craft.git
cd plan-craft
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` if needed. Default configuration uses SQLite at `data/plancraft.db`.

### 3. Install Dependencies

```bash
# Install Go dependencies
make deps

# Install frontend dependencies
make frontend-install
```

### 4. Run Database Migrations

```bash
make migrate
```

### 5. Run in Development Mode

```bash
make dev
```

This starts Wails in development mode with hot-reload for both Go backend and React frontend.

## Development Workflow

### Running the Application

| Command | Description |
|---------|-------------|
| `make dev` | Run in development mode with hot-reload |
| `make run` | Alias for `make dev` |
| `make build` | Build production binary |

### Database Management

| Command | Description |
|---------|-------------|
| `make migrate` | Run pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-status` | Show current migration version |
| `make db-clean` | Delete database files |

### Code Quality

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make fmt` | Format Go code |
| `make lint` | Run golangci-lint |

### Frontend Only

| Command | Description |
|---------|-------------|
| `make frontend-install` | Install npm dependencies |
| `make frontend-dev` | Run Vite dev server (standalone) |
| `make frontend-build` | Build frontend for production |

### Cleanup

| Command | Description |
|---------|-------------|
| `make clean` | Remove all build artifacts |

## Project Structure

```
plan-craft/
├── cmd/app/                 # Wails application entry point
├── config/                  # Environment configuration
├── frontend/                # React frontend
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   ├── layouts/         # Page layouts
│   │   ├── pages/           # Route pages
│   │   ├── router/          # React Router configuration
│   │   └── utils/           # Utility functions
│   └── wailsjs/             # Auto-generated Wails bindings
├── internal/
│   ├── entities/            # GORM domain models
│   ├── handlers/            # Wails handler methods (exposed to frontend)
│   ├── infrastructures/db/  # Database initialization
│   ├── repositories/        # Data access layer
│   └── services/            # Business logic layer
├── migrations/              # SQL migration files
├── data/                    # SQLite database (created at runtime)
├── wails.json               # Wails configuration
└── Makefile                 # Build and development tasks
```

## Database

The application uses SQLite with these optimizations:
- WAL (Write-Ahead Logging) mode for better concurrency
- 64MB cache for improved performance
- Foreign key constraints enabled
- Incremental auto-vacuum

Database file location: `data/plancraft.db`

## Roadmap

- **v1.0** - Project and work items management (current)
- **v1.1** - Timeline estimation and critical path
- **v1.2** - Resource planning and allocation
- **v1.3** - Cost estimation and tracking

## License

This project is licensed under the BSL 1.1 License - see the [LICENSE.md](LICENSE.md) file for details.

## Support

For issues and feature requests, please use the [GitHub issue tracker](https://github.com/ducminhgd/plan-craft/issues).