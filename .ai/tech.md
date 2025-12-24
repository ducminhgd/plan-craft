# Technology

## Tech Stack

1. Backend:
   1. Programing language: Golang
   2. Framework: go-chi
   3. Database: SQLite first, and be able to be extended to other databases (PostgreSQL, MySQL, etc.)
   4. Cache: Redis
   5. Logging: Uber Zap
   6. Testing: Go test
   7. CI/CD: GitHub Actions
   8. DB Migration: golang-migrate
2. Frontend:
   1. Programing language: TypeScript
   2. Framework: React
   3. UI Library: Material UI
   4. Testing: Jest
   5. CI/CD: GitHub Actions
3. Deployment:
   1. Docker: Docker
4. This project can be run locally, so it should be compatible with Windows, Linux, and macOS, and should be compiled to a single binary file.

## Principals

1. Backend is RESTFul API.
2. Frontend is a single page application.
3. Backend and frontend are separated.
4. Backend and frontend are deployed separately.
5. Frontend sends requests to Backend via RESTFul API, and they should be authenticated via JWT.

## Source code

There are two main folders in the source code:

1. `server`: contains the backend source code.
   1. `cmd`: Contains commands to run the application (server entry point, migration runner). **Required**.
   2. `config`: Configuration loading and parsing. **Optional**.
   3. `internal`: Internal packages used exclusively within this project. **Required**.
      1. `models`: Data models and ORM definitions (project, task, resource, cost, user, dependency). **Required**.
      2. `repositories`: Repository interfaces and implementations for database operations. **Required**.
      3. `services`: Business logic services (project, task, timeline, resource, cost, auth). **Required**.
      4. `handlers`: HTTP handlers for RESTful API endpoints. **Required**.
      5. `middleware`: HTTP middleware (auth, CORS, logging, recovery, rate limiting). **Required**.
      6. `utils`: Utility functions (crypto, JWT, time, response helpers). **Optional**.
      7. `db`: Database initialization and connection management. **Required**.
      8. `cache`: Cache initialization and operations (Redis). **Optional**.
   4. `pkg`: Packages exposed for use by other services or clients (errors, logger). **Optional**.
   5. `migrations`: Database migration files (up/down SQL files). **Required**.
   6. `tests`: Test files organized by package (integration, unit, testdata). **Optional**.
   7. `docs`: API documentation (Swagger/OpenAPI). **Optional**.
2. `site`: contains the frontend source code.
   1. `public`: Static assets (index.html, favicon, images). **Required**.
   2. `src`: Main source code for the React application. **Required**.
      1. `assets`: Images, fonts, and other static resources. **Optional**.
      2. `components`: Reusable React components organized by feature (common, layout, projects, tasks, timeline, resources, costs). **Required**.
      3. `pages`: Page-level components representing routes (ProjectsPage, ProjectDetailPage, TimelinePage, ResourcesPage, CostsPage, LoginPage, DashboardPage). **Required**.
      4. `services`: API service modules for backend communication (auth, project, task, resource, cost). **Required**.
      5. `hooks`: Custom React hooks (useAuth, useProject, useDebounce). **Optional**.
      6. `contexts`: React Context providers for global state management (AuthContext, ThemeContext). **Optional**.
      7. `types`: TypeScript type definitions and interfaces (project, task, resource, cost, api types). **Required**.
      8. `utils`: Utility functions and helpers (date, validation, calculations, constants). **Optional**.
      9. `routes`: Routing configuration (main routing setup, PrivateRoute). **Required**.
      10. `styles`: Global styles and theme configuration (Material UI theme, global CSS). **Optional**.
   3. `tests`: Test files organized by feature (components, services, utils). **Optional**.