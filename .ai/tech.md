# Technology

## Tech Stack

1. Programing language: Golang
2. UI: https://github.com/wailsapp/wails
3. Database: SQLite first, and be able to be extended to other databases (PostgreSQL, MySQL, etc.)
4. Logging: Uber Zap
5. Testing: Go test
6. CI/CD: GitHub Actions
7. DB Migration: [golang-migrate](https://github.com/golang-migrate/migrate)
8. This project can be run locally, so it should be compatible with Windows, Linux, and macOS, and should be compiled to a single binary file.


## Source code

1. In a package, the index file or the base file should be named as the same with the package name. For example entities should have `entities.go` as the index file.
2. `cmd`: Contains commands to run the application (server entry point, migration runner). **Required**.
3. `config`: Configuration loading and parsing. **Optional**.
4. `internal`: Internal packages used exclusively within this project. **Required**.
   1. `entities`: Data models and ORM definitions (project, task, resource, cost, user, dependency). **Required**.
      1. I don't need `deleted_at` field.
      2. I don't need to create base model.
   2. `repositories`: Repository interfaces and implementations for database operations. **Required**.
   3. `services`: Business logic services (project, task, timeline, resource, cost, auth). **Required**.
   4. `usecases`: Complex orchestrations and business rules. **Optional**.
   5. `infrastructures/db`: Database initialization and connection management. **Required**.
   6. `infrastructures/cache`: Cache initialization and operations (Redis). **Optional**.
   7. `presentations/wails`: Wails UI code. **Required**.
5. `pkg`: Packages exposed for use by other services or clients (errors, logger). **Optional**.
6. `migrations`: Database migration files (up/down SQL files). **Required**.
7. `tests`: Test files organized by package (integration, unit, testdata). **Optional**.
8. `docs`: API documentation (Swagger/OpenAPI). **Optional**.

```
plan-craft/
├── docs/
├── cmd/
│   └── app/
│       └── main.go                 # Application entry point
├── internal/
│   ├── entities/
│   ├── repositories/
│   ├── services/
│   ├── usecases/                   # Optional: complex orchestration
│   ├── infrastructures/             # Frameworks & drivers
│   │   ├── db/
│   │   ├── cache/
│   │   │   └── redis.go            # Or in-memory cache
│   │
│   └── presentations/               # Interface adapters
│       ├── wails/
│       │   ├── app.go              # Wails App struct
│       │   ├── client_handler.go     # Wails-bound methods
│       │   └── project_handler.go
│
├── frontend/                       # Wails frontend (React/Vue/Svelte)
│   ├── src/
│   ├── wailsjs/                    # Auto-generated Wails bindings
│   └── package.json
│
├── migrations/
│
├── pkg/                            # Public shared utilities
│   ├── x/                          # Utilities package
│   ├── logger/
│   │   └── logger.go
│   └── validator/
│       └── validator.go
│
├── config/
│   └── config.yaml
├── scripts/
│   └── build.sh
├── wails.json                      # Wails configuration
├── go.mod
└── go.sum
```