-- Drop the hierarchy view
DROP VIEW IF EXISTS gorev_hiyerarsi;

-- Drop the parent index
DROP INDEX IF EXISTS idx_gorev_parent;

-- Remove parent_id column (this will fail if there are any subtasks)
ALTER TABLE gorevler DROP COLUMN parent_id;