-- Alias index'ini kaldır
DROP INDEX IF EXISTS idx_template_alias;

-- Alias kolonunu kaldır
ALTER TABLE gorev_templateleri DROP COLUMN alias;