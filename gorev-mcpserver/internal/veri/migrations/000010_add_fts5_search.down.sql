-- Rollback FTS5 search functionality

-- Drop triggers first
DROP TRIGGER IF EXISTS gorevler_fts_tag_delete;
DROP TRIGGER IF EXISTS gorevler_fts_tag_update;
DROP TRIGGER IF EXISTS gorevler_fts_delete;
DROP TRIGGER IF EXISTS gorevler_fts_update;
DROP TRIGGER IF EXISTS gorevler_fts_insert;

-- Drop indexes
DROP INDEX IF EXISTS idx_search_history_created;
DROP INDEX IF EXISTS idx_search_history_query;
DROP INDEX IF EXISTS idx_filter_profiles_last_used;
DROP INDEX IF EXISTS idx_filter_profiles_default;
DROP INDEX IF EXISTS idx_filter_profiles_name;

-- Drop tables
DROP TABLE IF EXISTS search_history;
DROP TABLE IF EXISTS filter_profiles;
DROP TABLE IF EXISTS gorevler_fts;