# Plan Craft

A comprehensive project management and estimation tool designed to help teams plan, estimate, and track software projects effectively.

## Overview

Plan Craft is a project management tool that enables teams to:
- Define and manage work items with hierarchical breakdown structures
- Estimate effort in man-days or man-months for tasks, milestones, and projects
- Calculate timelines and roadmaps based on task dependencies
- Determine resource requirements (number of people and duration)
- Estimate project costs based on resources and rates

## Core Features

### 1. Project Management
- Project metadata (name, type, methodology: Waterfall/Agile/Hybrid)
- Start date and target end date tracking
- Assumptions and constraints documentation

### 2. Work Breakdown Structure (WBS)
- Hierarchical task organization (epics → tasks → subtasks)
- Estimated effort tracking (hours/days)
- Task dependency management:
  - Finish-to-Start dependencies
  - Start-to-Start dependencies
- Critical path calculation

### 3. Timeline & Dependencies
- Auto-calculated project duration
- Gantt-style timeline visualization
- Slack/buffer visibility
- Roadmap planning

### 4. Resource Planning
- Role definition (PM, Backend, QA, Designer, etc.)
- Resource assignment to tasks
- Capacity limits per resource
- Resource utilization tracking

### 5. Cost Estimation
- Hourly/daily cost per role
- Total cost per task, phase, and project
- Cost breakdown by category

## Tech Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: go-chi
- **Database**: SQLite (extensible to PostgreSQL, MySQL)
- **Cache**: Redis
- **Logging**: Uber Zap
- **Testing**: Go test
- **DB Migration**: golang-migrate

### Frontend
- **Language**: TypeScript
- **Framework**: React
- **UI Library**: Material UI
- **Testing**: Jest

### Deployment
- **Containerization**: Docker
- **CI/CD**: GitHub Actions
- **Distribution**: Single binary for Windows, Linux, and macOS

## Architecture

- RESTful API backend
- Single Page Application (SPA) frontend
- JWT-based authentication
- Backend and frontend deployed separately

## Roadmap

- **Version 1.0**: Project management, work items management
- **Version 1.1**: Timeline estimation
- **Version 1.2**: Resource planning
- **Version 1.3**: Cost estimation

## License

This project is licensed under the BSL 1.1 License - see the [LICENSE.md](LICENSE.md) file for details.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- SQLite (included with most systems)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/ducminhgd/plan-craft.git
cd plan-craft
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Install dependencies:
```bash
make deps
```

4. Build the application:
```bash
make build
```

5. Run the application:
```bash
make run
```

The application will:
- Create the SQLite database at `data/plancraft.db`
- Run auto-migrations to create all tables
- Create sample data (project, tasks, resources)

### Development

Run in development mode with auto-reload (requires [air](https://github.com/air-verse/air)):
```bash
make dev
```

### Database

The application uses SQLite with the following optimizations:
- WAL (Write-Ahead Logging) mode for better concurrency
- 64MB cache for improved performance
- Foreign key constraints enabled
- Incremental auto-vacuum

To clean the database:
```bash
make db-clean
```

### Project Structure

```
plan-craft/
├── cmd/
│   └── server/          # Application entry point
├── config/              # Configuration management
├── internal/
│   ├── db/             # Database initialization
│   └── models/         # GORM models
├── data/               # SQLite database files
├── .env.example        # Example environment variables
├── Makefile           # Build and development tasks
└── README.md
```

### Available Make Commands

- `make help` - Show available commands
- `make deps` - Download dependencies
- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make db-clean` - Remove database file
- `make dev` - Run in development mode
- `make fmt` - Format code
- `make lint` - Lint code

## Contributing

(Coming soon)

## Support

For issues and feature requests, please use the GitHub issue tracker.