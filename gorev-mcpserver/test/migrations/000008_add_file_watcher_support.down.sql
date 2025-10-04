-- Rollback migration for file watcher support

DROP TABLE IF EXISTS task_file_paths;
DROP TABLE IF EXISTS file_watcher_config;