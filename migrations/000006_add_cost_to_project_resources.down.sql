-- Drop cost index
DROP INDEX IF EXISTS idx_project_resources_cost;

-- Remove cost column from project_resources table
-- Note: SQLite does not support DROP COLUMN directly, so we need to recreate the table
CREATE TABLE project_resources_backup AS SELECT
    id, project_id, human_resource_id, role, allocation, start_date, end_date, notes, status, created_at, updated_at
FROM project_resources;

DROP TABLE project_resources;

CREATE TABLE project_resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    human_resource_id INTEGER NOT NULL,
    role TEXT,
    allocation REAL NOT NULL DEFAULT 100,
    start_date INTEGER,
    end_date INTEGER,
    notes TEXT,
    status INTEGER NOT NULL DEFAULT 2,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    CHECK (status IN (1, 2)),
    CHECK (allocation >= 0 AND allocation <= 100),

    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (human_resource_id) REFERENCES human_resources(id) ON DELETE RESTRICT
);

INSERT INTO project_resources SELECT * FROM project_resources_backup;

DROP TABLE project_resources_backup;

-- Recreate indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_project_human_resource ON project_resources(project_id, human_resource_id);
CREATE INDEX IF NOT EXISTS idx_project_resources_project_id ON project_resources(project_id);
CREATE INDEX IF NOT EXISTS idx_project_resources_human_resource_id ON project_resources(human_resource_id);
CREATE INDEX IF NOT EXISTS idx_project_resources_role ON project_resources(role);
CREATE INDEX IF NOT EXISTS idx_project_resources_status ON project_resources(status);
CREATE INDEX IF NOT EXISTS idx_project_resources_start_date ON project_resources(start_date);
CREATE INDEX IF NOT EXISTS idx_project_resources_end_date ON project_resources(end_date);
CREATE INDEX IF NOT EXISTS idx_project_resources_created_at ON project_resources(created_at);
CREATE INDEX IF NOT EXISTS idx_project_resources_updated_at ON project_resources(updated_at);
