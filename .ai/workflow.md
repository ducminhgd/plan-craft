# Workflow

1. When a new entity is created:
   1. Create the entity and unit tests in `internal/entities/`.
   2. Create the repository interface and unit tests in `internal/repositories/`.
   3. Create the service interface and unit tests in `internal/services/`.
   4. Create the handler and unit tests in `internal/handlers/`.
   5. Create the router in `internal/routers/`.
   6. Create the migration in `migrations/`.
   7. Update `/internal/infrastructures/database.go` for auto-migration.
2. For unit tests:
   1. Add comments to ignore this kind of error: Error return value is not checked (errcheck)