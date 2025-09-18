-- Add parent_id column to support subtask hierarchy
ALTER TABLE gorevler ADD COLUMN parent_id TEXT REFERENCES gorevler(id) ON DELETE RESTRICT;

-- Create index for efficient subtask queries
CREATE INDEX idx_gorev_parent ON gorevler(parent_id);

-- Create a view for recursive subtask queries using CTE
CREATE VIEW gorev_hiyerarsi AS
WITH RECURSIVE hiyerarsi AS (
    -- Base case: root tasks (no parent)
    SELECT 
        id,
        baslik,
        aciklama,
        durum,
        oncelik,
        proje_id,
        parent_id,
        son_tarih,
        0 as seviye,
        id as kok_id,
        '/' || id as yol
    FROM gorevler
    WHERE parent_id IS NULL
    
    UNION ALL
    
    -- Recursive case: subtasks
    SELECT 
        g.id,
        g.baslik,
        g.aciklama,
        g.durum,
        g.oncelik,
        g.proje_id,
        g.parent_id,
        g.son_tarih,
        h.seviye + 1,
        h.kok_id,
        h.yol || '/' || g.id
    FROM gorevler g
    INNER JOIN hiyerarsi h ON g.parent_id = h.id
)
SELECT * FROM hiyerarsi;