# Workflow

1. When a new entity is created:
   1. Create the entity in `internal/entities/`.
   2. Create the repository interface in `internal/repositories/`.
   3. Create the service interface in `internal/services/`.
   4. Create the handler in `internal/handlers/`.
   5. Create the router in `internal/routers/`.
   6. Create the migration in `migrations/`.
   7. Update `/internal/infrastructures/database.go` for auto-migration.