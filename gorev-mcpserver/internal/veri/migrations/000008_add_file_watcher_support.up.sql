-- Migration to add file watcher support
-- Add task_file_paths table to track file paths associated with tasks

CREATE TABLE task_file_paths (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL REFERENCES gorevler(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_id, file_path)
);

-- Add indexes for efficient lookups
CREATE INDEX idx_task_file_paths_task_id ON task_file_paths(task_id);
CREATE INDEX idx_task_file_paths_file_path ON task_file_paths(file_path);

-- Add file watcher configuration table
CREATE TABLE file_watcher_config (
    id INTEGER PRIMARY KEY DEFAULT 1,
    enabled BOOLEAN DEFAULT TRUE,
    watched_extensions JSON,
    ignore_patterns JSON,
    auto_update_status BOOLEAN DEFAULT TRUE,
    debounce_duration_ms INTEGER DEFAULT 500,
    max_file_size_bytes INTEGER DEFAULT 10485760, -- 10MB
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (id = 1) -- Ensure only one configuration row
);

-- Insert default configuration
INSERT INTO file_watcher_config (
    watched_extensions,
    ignore_patterns
) VALUES (
    '["\.go", "\.js", "\.ts", "\.py", "\.java", "\.cpp", "\.c", "\.h", "\.md", "\.txt", "\.json", "\.yaml", "\.yml"]',
    '["node_modules", "\.git", "\.vscode", "vendor", "build", "dist", "*\.tmp", "*\.log", "*\.swp"]'
);