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