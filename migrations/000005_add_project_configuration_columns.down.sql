-- SQLite doesn't support DROP COLUMN directly before version 3.35.0
-- We need to recreate the table without these columns

-- Create temporary table with old schema
CREATE TABLE projects_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    client_id INTEGER NOT NULL,
    start_date INTEGER,
    end_date INTEGER,
    status INTEGER NOT NULL DEFAULT 2,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    CHECK (status IN (1, 2)),

    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE RESTRICT
);

-- Copy data to backup table
INSERT INTO projects_backup (id, name, description, client_id, start_date, end_date, status, created_at, updated_at)
SELECT id, name, description, client_id, start_date, end_date, status, created_at, updated_at
FROM projects;

-- Drop the original table
DROP TABLE projects;

-- Rename backup to original
ALTER TABLE projects_backup RENAME TO projects;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_projects_client_id ON projects(client_id);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_start_date ON projects(start_date);
CREATE INDEX IF NOT EXISTS idx_projects_end_date ON projects(end_date);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at);
CREATE INDEX IF NOT EXISTS idx_projects_updated_at ON projects(updated_at);
