-- Drop indexes
DROP INDEX IF EXISTS idx_projects_name;
DROP INDEX IF EXISTS idx_projects_client_id;
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_start_date;
DROP INDEX IF EXISTS idx_projects_end_date;
DROP INDEX IF EXISTS idx_projects_created_at;
DROP INDEX IF EXISTS idx_projects_updated_at;

-- Drop projects table
DROP TABLE IF EXISTS projects;
