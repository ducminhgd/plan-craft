-- Create human_resources table
CREATE TABLE IF NOT EXISTS human_resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    title TEXT NOT NULL,
    level TEXT NOT NULL,
    status INTEGER NOT NULL DEFAULT 2,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,

    -- Add CHECK constraints for status validation
    CHECK (status IN (1, 2))
);

-- Create indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_human_resources_name ON human_resources(name);
CREATE INDEX IF NOT EXISTS idx_human_resources_title ON human_resources(title);
CREATE INDEX IF NOT EXISTS idx_human_resources_level ON human_resources(level);
CREATE INDEX IF NOT EXISTS idx_human_resources_status ON human_resources(status);
CREATE INDEX IF NOT EXISTS idx_human_resources_created_at ON human_resources(created_at);
CREATE INDEX IF NOT EXISTS idx_human_resources_updated_at ON human_resources(updated_at);
