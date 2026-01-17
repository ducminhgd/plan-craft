-- Drop indexes
DROP INDEX IF EXISTS idx_project_human_resource;
DROP INDEX IF EXISTS idx_project_resources_project_id;
DROP INDEX IF EXISTS idx_project_resources_human_resource_id;
DROP INDEX IF EXISTS idx_project_resources_role;
DROP INDEX IF EXISTS idx_project_resources_status;
DROP INDEX IF EXISTS idx_project_resources_start_date;
DROP INDEX IF EXISTS idx_project_resources_end_date;
DROP INDEX IF EXISTS idx_project_resources_created_at;
DROP INDEX IF EXISTS idx_project_resources_updated_at;

-- Drop project_resources table
DROP TABLE IF EXISTS project_resources;
