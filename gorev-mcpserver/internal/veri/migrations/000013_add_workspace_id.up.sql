-- Migration: Add workspace_id column to all tables for centralized mode support
-- This enables multi-tenant database isolation when running in centralized mode

-- Add workspace_id to gorevler (tasks) table
ALTER TABLE gorevler ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to projeler (projects) table
ALTER TABLE projeler ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to baglantilar (dependencies) table
ALTER TABLE baglantilar ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to etiketler (tags) table
ALTER TABLE etiketler ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to gorev_templateleri (task templates) table
ALTER TABLE gorev_templateleri ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to ai_interactions table
ALTER TABLE ai_interactions ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to ai_context table
ALTER TABLE ai_context ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to aktif_proje (active project) table
ALTER TABLE aktif_proje ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to filter_profiles table
ALTER TABLE filter_profiles ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Add workspace_id to search_history table
ALTER TABLE search_history ADD COLUMN workspace_id TEXT NOT NULL DEFAULT 'default';

-- Create indexes for efficient workspace filtering
CREATE INDEX IF NOT EXISTS idx_gorevler_workspace_id ON gorevler(workspace_id);
CREATE INDEX IF NOT EXISTS idx_projeler_workspace_id ON projeler(workspace_id);
CREATE INDEX IF NOT EXISTS idx_baglantilar_workspace_id ON baglantilar(workspace_id);
CREATE INDEX IF NOT EXISTS idx_etiketler_workspace_id ON etiketler(workspace_id);
CREATE INDEX IF NOT EXISTS idx_gorev_templateleri_workspace_id ON gorev_templateleri(workspace_id);
CREATE INDEX IF NOT EXISTS idx_ai_interactions_workspace_id ON ai_interactions(workspace_id);
CREATE INDEX IF NOT EXISTS idx_ai_context_workspace_id ON ai_context(workspace_id);
CREATE INDEX IF NOT EXISTS idx_aktif_proje_workspace_id ON aktif_proje(workspace_id);
CREATE INDEX IF NOT EXISTS idx_filter_profiles_workspace_id ON filter_profiles(workspace_id);
CREATE INDEX IF NOT EXISTS idx_search_history_workspace_id ON search_history(workspace_id);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_gorevler_workspace_status ON gorevler(workspace_id, status);
CREATE INDEX IF NOT EXISTS idx_gorevler_workspace_project ON gorevler(workspace_id, project_id);
CREATE INDEX IF NOT EXISTS idx_projeler_workspace_name ON projeler(workspace_id, name);
