-- Drop indexes
DROP INDEX IF EXISTS idx_milestones_name;
DROP INDEX IF EXISTS idx_milestones_project_id;
DROP INDEX IF EXISTS idx_milestones_status;
DROP INDEX IF EXISTS idx_milestones_start_date;
DROP INDEX IF EXISTS idx_milestones_end_date;
DROP INDEX IF EXISTS idx_milestones_created_at;
DROP INDEX IF EXISTS idx_milestones_updated_at;

-- Drop milestones table
DROP TABLE IF EXISTS milestones;
