-- Drop indexes
DROP INDEX IF EXISTS idx_human_resources_updated_at;
DROP INDEX IF EXISTS idx_human_resources_created_at;
DROP INDEX IF EXISTS idx_human_resources_status;
DROP INDEX IF EXISTS idx_human_resources_level;
DROP INDEX IF EXISTS idx_human_resources_title;
DROP INDEX IF EXISTS idx_human_resources_name;

-- Drop table
DROP TABLE IF EXISTS human_resources;
