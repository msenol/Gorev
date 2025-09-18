-- Indexes for baglantilar table to improve dependency query performance
-- This fixes the N+1 query performance issue in GorevListele

-- Index on kaynak_id for finding dependencies of a task
CREATE INDEX IF NOT EXISTS idx_baglantilar_kaynak_id ON baglantilar(kaynak_id);

-- Index on hedef_id for finding reverse dependencies (what depends on this task)
CREATE INDEX IF NOT EXISTS idx_baglantilar_hedef_id ON baglantilar(hedef_id);

-- Composite index for faster joins and unique constraint checks
CREATE INDEX IF NOT EXISTS idx_baglantilar_kaynak_hedef ON baglantilar(kaynak_id, hedef_id);