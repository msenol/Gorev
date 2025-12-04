-- Rollback: Remove workspace_id columns from all tables
-- Note: SQLite doesn't support DROP COLUMN directly in older versions
-- We need to recreate tables without the workspace_id column

-- Drop indexes first
DROP INDEX IF EXISTS idx_gorevler_workspace_id;
DROP INDEX IF EXISTS idx_projeler_workspace_id;
DROP INDEX IF EXISTS idx_baglantilar_workspace_id;
DROP INDEX IF EXISTS idx_etiketler_workspace_id;
DROP INDEX IF EXISTS idx_gorev_templateleri_workspace_id;
DROP INDEX IF EXISTS idx_ai_interactions_workspace_id;
DROP INDEX IF EXISTS idx_ai_context_workspace_id;
DROP INDEX IF EXISTS idx_aktif_proje_workspace_id;
DROP INDEX IF EXISTS idx_filter_profiles_workspace_id;
DROP INDEX IF EXISTS idx_search_history_workspace_id;
DROP INDEX IF EXISTS idx_gorevler_workspace_status;
DROP INDEX IF EXISTS idx_gorevler_workspace_project;
DROP INDEX IF EXISTS idx_projeler_workspace_name;

-- For SQLite 3.35.0+ we can use ALTER TABLE DROP COLUMN
-- For older versions, this migration cannot be rolled back cleanly

-- Drop workspace_id columns (requires SQLite 3.35.0+)
ALTER TABLE gorevler DROP COLUMN workspace_id;
ALTER TABLE projeler DROP COLUMN workspace_id;
ALTER TABLE baglantilar DROP COLUMN workspace_id;
ALTER TABLE etiketler DROP COLUMN workspace_id;
ALTER TABLE gorev_templateleri DROP COLUMN workspace_id;
ALTER TABLE ai_interactions DROP COLUMN workspace_id;
ALTER TABLE ai_context DROP COLUMN workspace_id;
ALTER TABLE aktif_proje DROP COLUMN workspace_id;
ALTER TABLE filter_profiles DROP COLUMN workspace_id;
ALTER TABLE search_history DROP COLUMN workspace_id;
