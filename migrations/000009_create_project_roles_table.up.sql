-- Create project_roles table
CREATE TABLE IF NOT EXISTS project_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    level INTEGER NOT NULL DEFAULT 2,
    headcount INTEGER NOT NULL DEFAULT 1,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    -- Add CHECK constraints for validation
    CHECK (level IN (1, 2, 3, 4)),
    CHECK (headcount >= 0),

    -- Foreign key constraints
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Create unique index on project_id, name, and level
-- Each role of each level can only have 1 row per project
CREATE UNIQUE INDEX IF NOT EXISTS idx_project_role_name_level ON project_roles(project_id, name, level);

-- Create indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_project_roles_project_id ON project_roles(project_id);
CREATE INDEX IF NOT EXISTS idx_project_roles_name ON project_roles(name);
CREATE INDEX IF NOT EXISTS idx_project_roles_level ON project_roles(level);
CREATE INDEX IF NOT EXISTS idx_project_roles_created_at ON project_roles(created_at);
CREATE INDEX IF NOT EXISTS idx_project_roles_updated_at ON project_roles(updated_at);
