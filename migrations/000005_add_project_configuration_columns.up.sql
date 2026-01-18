-- Add configuration columns to projects table
ALTER TABLE projects ADD COLUMN hours_per_day INTEGER NOT NULL DEFAULT 8;
ALTER TABLE projects ADD COLUMN days_per_week INTEGER NOT NULL DEFAULT 5;
ALTER TABLE projects ADD COLUMN working_days_per_week TEXT NOT NULL DEFAULT '[1,2,3,4,5]';
ALTER TABLE projects ADD COLUMN timezone TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN currency TEXT NOT NULL DEFAULT '';

-- Add CHECK constraints for configuration validation
-- Note: SQLite doesn't support ALTER TABLE ADD CONSTRAINT, so we validate in application layer
