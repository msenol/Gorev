-- Remove indexes from baglantilar table

DROP INDEX IF EXISTS idx_baglantilar_kaynak_id;
DROP INDEX IF EXISTS idx_baglantilar_hedef_id;
DROP INDEX IF EXISTS idx_baglantilar_kaynak_hedef;