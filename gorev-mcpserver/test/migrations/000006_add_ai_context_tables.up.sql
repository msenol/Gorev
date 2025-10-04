-- Add AI context management tables for tracking AI interactions and session context

-- Table for tracking AI interactions with tasks
CREATE TABLE IF NOT EXISTS ai_interactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gorev_id INTEGER NOT NULL,
    action_type TEXT NOT NULL CHECK (action_type IN ('viewed', 'created', 'updated', 'completed', 'set_active', 'bulk_operation')),
    context TEXT, -- JSON field for storing additional context
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (gorev_id) REFERENCES gorevler(id) ON DELETE CASCADE
);

-- Table for maintaining AI session context
CREATE TABLE IF NOT EXISTS ai_context (
    id INTEGER PRIMARY KEY CHECK (id = 1), -- Single row table
    active_task_id INTEGER,
    recent_tasks TEXT, -- JSON array of recent task IDs
    session_data TEXT, -- JSON field for additional session context
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (active_task_id) REFERENCES gorevler(id) ON DELETE SET NULL
);

-- Add new columns to gorevler table for AI tracking
ALTER TABLE gorevler ADD COLUMN last_ai_interaction DATETIME;
ALTER TABLE gorevler ADD COLUMN estimated_hours INTEGER;
ALTER TABLE gorevler ADD COLUMN actual_hours INTEGER;

-- Create indexes for better query performance
CREATE INDEX idx_ai_interactions_gorev_id ON ai_interactions(gorev_id);
CREATE INDEX idx_ai_interactions_timestamp ON ai_interactions(timestamp);
CREATE INDEX idx_ai_interactions_action_type ON ai_interactions(action_type);
CREATE INDEX idx_gorevler_last_ai_interaction ON gorevler(last_ai_interaction);

-- Initialize the ai_context table with a single row
INSERT INTO ai_context (id, active_task_id, recent_tasks, session_data) 
VALUES (1, NULL, '[]', '{}');