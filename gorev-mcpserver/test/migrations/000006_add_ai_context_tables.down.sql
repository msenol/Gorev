-- Rollback AI context management changes

-- Drop indexes
DROP INDEX IF EXISTS idx_gorevler_last_ai_interaction;
DROP INDEX IF EXISTS idx_ai_interactions_action_type;
DROP INDEX IF EXISTS idx_ai_interactions_timestamp;
DROP INDEX IF EXISTS idx_ai_interactions_gorev_id;

-- SQLite doesn't support DROP COLUMN directly, so we need to recreate the table
-- First, create a temporary table with the original schema
CREATE TABLE gorevler_temp AS 
SELECT id, baslik, aciklama, durum, oncelik, proje_id, parent_id, 
       olusturma_tarih, guncelleme_tarih, son_tarih
FROM gorevler;

-- Drop the original table
DROP TABLE gorevler;

-- Recreate the original table without the AI columns
CREATE TABLE gorevler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    baslik TEXT NOT NULL,
    aciklama TEXT,
    durum TEXT DEFAULT 'beklemede' CHECK (durum IN ('beklemede', 'devam_ediyor', 'tamamlandi')),
    oncelik TEXT DEFAULT 'orta' CHECK (oncelik IN ('yuksek', 'orta', 'dusuk', 'düşük')),
    proje_id INTEGER,
    parent_id INTEGER,
    olusturma_tarih DATETIME DEFAULT CURRENT_TIMESTAMP,
    guncelleme_tarih DATETIME DEFAULT CURRENT_TIMESTAMP,
    son_tarih DATE,
    FOREIGN KEY (proje_id) REFERENCES projeler(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES gorevler(id) ON DELETE CASCADE
);

-- Copy data back
INSERT INTO gorevler SELECT * FROM gorevler_temp;

-- Drop the temporary table
DROP TABLE gorevler_temp;

-- Recreate indexes and triggers that were on the original table
CREATE INDEX idx_gorevler_proje_id ON gorevler(proje_id);
CREATE INDEX idx_gorevler_durum ON gorevler(durum);
CREATE INDEX idx_gorevler_oncelik ON gorevler(oncelik);
CREATE INDEX idx_gorevler_parent_id ON gorevler(parent_id);
CREATE INDEX idx_gorevler_son_tarih ON gorevler(son_tarih);

-- Drop AI-specific tables
DROP TABLE IF EXISTS ai_context;
DROP TABLE IF EXISTS ai_interactions;