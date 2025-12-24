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

(Coming soon)

## Contributing

(Coming soon)

## Support

For issues and feature requests, please use the GitHub issue tracker.