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
3. For the Database management features:
   1. When the application is started, there is no DB to connect to, so the application will use the memory database. And, on the menu bar, display "draft" instead of the database name.
   2. If the user Save the current database as another one (the Save As feature), for both cases from memory to a file or from a file to another file, the application will wire up with new database and update the menu bar with the new database name.
   3. If the user open an existing database, the application will wire up with the new database, and update the menu bar with the new database name, and reload the current page.
   4. If the user close the application when the database is in draft mode, the application will ask if the user want to save the current database before closing.