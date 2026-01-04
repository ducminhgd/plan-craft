-- Drop all triggers
DROP TRIGGER IF EXISTS update_costs_updated_at;
DROP TRIGGER IF EXISTS update_task_assignments_updated_at;
DROP TRIGGER IF EXISTS update_resource_allocations_updated_at;
DROP TRIGGER IF EXISTS update_project_roles_updated_at;
DROP TRIGGER IF EXISTS update_resources_updated_at;
DROP TRIGGER IF EXISTS update_task_dependencies_updated_at;
DROP TRIGGER IF EXISTS update_tasks_updated_at;
DROP TRIGGER IF EXISTS update_milestones_updated_at;
DROP TRIGGER IF EXISTS update_projects_updated_at;
DROP TRIGGER IF EXISTS update_clients_updated_at;

-- Drop all indexes
DROP INDEX IF EXISTS idx_costs_date;
DROP INDEX IF EXISTS idx_costs_is_estimated;
DROP INDEX IF EXISTS idx_costs_project_role;
DROP INDEX IF EXISTS idx_costs_resource;
DROP INDEX IF EXISTS idx_costs_category;
DROP INDEX IF EXISTS idx_costs_type;
DROP INDEX IF EXISTS idx_costs_task;
DROP INDEX IF EXISTS idx_costs_milestone;
DROP INDEX IF EXISTS idx_costs_project;

DROP INDEX IF EXISTS idx_task_assignments_resource;
DROP INDEX IF EXISTS idx_task_assignments_role;
DROP INDEX IF EXISTS idx_task_assignments_task;

DROP INDEX IF EXISTS idx_resource_allocation_active;
DROP INDEX IF EXISTS idx_resource_allocation_dates;
DROP INDEX IF EXISTS idx_resource_allocation_project;
DROP INDEX IF EXISTS idx_resource_allocation_resource;
DROP INDEX IF EXISTS idx_resource_allocation_role;

DROP INDEX IF EXISTS idx_project_roles_resource;
DROP INDEX IF EXISTS idx_project_roles_project;

DROP INDEX IF EXISTS idx_resources_is_active;
DROP INDEX IF EXISTS idx_resources_role;
DROP INDEX IF EXISTS idx_resources_email;

DROP INDEX IF EXISTS idx_task_dependencies_depends;
DROP INDEX IF EXISTS idx_task_dependencies_task;

DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_parent;
DROP INDEX IF EXISTS idx_tasks_milestone;
DROP INDEX IF EXISTS idx_tasks_project;

DROP INDEX IF EXISTS idx_milestones_status;
DROP INDEX IF EXISTS idx_milestones_project;

DROP INDEX IF EXISTS idx_projects_client_id;
DROP INDEX IF EXISTS idx_projects_owner_id;
DROP INDEX IF EXISTS idx_projects_status;
DROP INDEX IF EXISTS idx_projects_name;

DROP INDEX IF EXISTS idx_clients_code;
DROP INDEX IF EXISTS idx_clients_is_active;

-- Drop all tables in reverse dependency order
DROP TABLE IF EXISTS costs;
DROP TABLE IF EXISTS task_assignments;
DROP TABLE IF EXISTS resource_allocations;
DROP TABLE IF EXISTS project_roles;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS task_dependencies;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS milestones;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS clients;
