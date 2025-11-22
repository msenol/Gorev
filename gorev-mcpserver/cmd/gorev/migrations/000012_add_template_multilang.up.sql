-- Migration: Add multi-language support to templates
-- Version: v0.18.0
-- Date: 2025-11-22
-- Purpose: Enable templates to support multiple languages (TR, EN, etc.)

-- ========================================
-- 1. Add language support columns to gorev_templateleri
-- ========================================
ALTER TABLE gorev_templateleri
ADD COLUMN language_code VARCHAR(2) NOT NULL DEFAULT 'tr';

ALTER TABLE gorev_templateleri
ADD COLUMN base_template_id TEXT;

-- ========================================
-- 2. Remove old UNIQUE constraint on alias
-- ========================================
-- Drop the old unique index on alias created in migration 009
DROP INDEX IF EXISTS idx_template_alias;

-- ========================================
-- 3. Update existing templates to multi-language format
-- ========================================
-- Set language_code to 'tr' for existing templates (they're in Turkish)
-- Set base_template_id to their own ID for grouping
UPDATE gorev_templateleri
SET language_code = 'tr', base_template_id = id
WHERE language_code IS NULL OR base_template_id IS NULL;

-- ========================================
-- 4. Create new composite unique index
-- ========================================
-- Create new unique index on (alias, language_code) to allow same alias in different languages
CREATE UNIQUE INDEX idx_template_alias_lang ON gorev_templateleri(alias, language_code);

-- ========================================
-- 5. Create indexes for language and base template
-- ========================================
CREATE INDEX idx_template_language ON gorev_templateleri(language_code);
CREATE INDEX idx_template_base_id ON gorev_templateleri(base_template_id);
CREATE INDEX idx_template_language_active ON gorev_templateleri(language_code, active);

-- ========================================
-- 6. Migration verification
-- ========================================
-- Count templates before and after (for verification)
-- This will be logged during migration:
-- SELECT COUNT(*) as total_templates FROM gorev_templateleri WHERE language_code = 'tr';

-- Migration completed successfully
-- Templates now support multi-language with unique (alias, language_code) constraint
-- Existing templates are preserved and updated to support multiple languages
