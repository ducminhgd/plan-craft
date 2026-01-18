-- Add cost column to project_resources table
ALTER TABLE project_resources ADD COLUMN cost REAL NOT NULL DEFAULT 0;

-- Create index for cost column
CREATE INDEX IF NOT EXISTS idx_project_resources_cost ON project_resources(cost);
