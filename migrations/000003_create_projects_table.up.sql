-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    client_id INTEGER NOT NULL,
    start_date INTEGER,
    end_date INTEGER,
    status INTEGER NOT NULL DEFAULT 2,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    -- Add CHECK constraints for status validation
    CHECK (status IN (1, 2)),

    -- Foreign key constraint
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE RESTRICT
);

-- Create indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_projects_client_id ON projects(client_id);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_start_date ON projects(start_date);
CREATE INDEX IF NOT EXISTS idx_projects_end_date ON projects(end_date);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at);
CREATE INDEX IF NOT EXISTS idx_projects_updated_at ON projects(updated_at);
