-- Create milestones table
CREATE TABLE IF NOT EXISTS milestones (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    project_id INTEGER NOT NULL,
    start_date INTEGER,
    end_date INTEGER,
    status INTEGER NOT NULL DEFAULT 2,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    -- Add CHECK constraints for status validation
    CHECK (status IN (1, 2)),

    -- Foreign key constraint
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Create indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_milestones_name ON milestones(name);
CREATE INDEX IF NOT EXISTS idx_milestones_project_id ON milestones(project_id);
CREATE INDEX IF NOT EXISTS idx_milestones_status ON milestones(status);
CREATE INDEX IF NOT EXISTS idx_milestones_start_date ON milestones(start_date);
CREATE INDEX IF NOT EXISTS idx_milestones_end_date ON milestones(end_date);
CREATE INDEX IF NOT EXISTS idx_milestones_created_at ON milestones(created_at);
CREATE INDEX IF NOT EXISTS idx_milestones_updated_at ON milestones(updated_at);
