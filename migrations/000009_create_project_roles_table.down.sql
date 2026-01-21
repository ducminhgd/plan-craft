-- Drop indexes
DROP INDEX IF EXISTS idx_project_roles_updated_at;
DROP INDEX IF EXISTS idx_project_roles_created_at;
DROP INDEX IF EXISTS idx_project_roles_level;
DROP INDEX IF EXISTS idx_project_roles_name;
DROP INDEX IF EXISTS idx_project_roles_project_id;
DROP INDEX IF EXISTS idx_project_role_name_level;

-- Drop project_roles table
DROP TABLE IF EXISTS project_roles;
