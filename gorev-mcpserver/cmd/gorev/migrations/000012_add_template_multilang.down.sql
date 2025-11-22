-- Migration Down: Remove multi-language support from templates
-- Reverts changes from 000012_add_template_multilang.up.sql

-- ========================================
-- 1. Drop new indexes
-- ========================================
DROP INDEX IF EXISTS idx_template_alias_lang;
DROP INDEX IF EXISTS idx_template_language;
DROP INDEX IF EXISTS idx_template_base_id;
DROP INDEX IF EXISTS idx_template_language_active;

-- ========================================
-- 2. Remove language support columns
-- ========================================
ALTER TABLE gorev_templateleri
DROP COLUMN IF EXISTS base_template_id;

ALTER TABLE gorev_templateleri
DROP COLUMN IF EXISTS language_code;

-- ========================================
-- 3. Restore old UNIQUE constraint on alias
-- ========================================
-- Recreate the old unique index from migration 009
CREATE UNIQUE INDEX idx_template_alias ON gorev_templateleri(alias) WHERE alias IS NOT NULL;

-- Migration down completed successfully
-- Templates reverted to single-language (Turkish only)
