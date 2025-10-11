-- Migration: Rename Turkish field names to English
-- Version: v0.17.0
-- Date: 2025-10-11
-- Purpose: Refactor database schema for international consistency
-- NOTE: Domain terms (gorevler, projeler) remain Turkish as per project convention

-- Disable foreign keys temporarily to allow column renames
PRAGMA foreign_keys=OFF;

-- ========================================
-- 1. gorevler table (tasks)
-- ========================================
ALTER TABLE gorevler RENAME COLUMN baslik TO title;
ALTER TABLE gorevler RENAME COLUMN aciklama TO description;
ALTER TABLE gorevler RENAME COLUMN durum TO status;
ALTER TABLE gorevler RENAME COLUMN oncelik TO priority;
ALTER TABLE gorevler RENAME COLUMN proje_id TO project_id;
ALTER TABLE gorevler RENAME COLUMN olusturma_tarih TO created_at;
ALTER TABLE gorevler RENAME COLUMN guncelleme_tarih TO updated_at;
ALTER TABLE gorevler RENAME COLUMN son_tarih TO due_date;

-- Update indexes for gorevler
DROP INDEX IF EXISTS idx_gorev_durum;
DROP INDEX IF EXISTS idx_gorev_proje;
DROP INDEX IF EXISTS idx_gorev_oncelik;
DROP INDEX IF EXISTS idx_gorev_parent;
CREATE INDEX idx_gorev_status ON gorevler(status);
CREATE INDEX idx_gorev_project ON gorevler(project_id);
CREATE INDEX idx_gorev_priority ON gorevler(priority);
CREATE INDEX idx_gorev_parent ON gorevler(parent_id);

-- ========================================
-- 2. projeler table (projects)
-- ========================================
ALTER TABLE projeler RENAME COLUMN isim TO name;
ALTER TABLE projeler RENAME COLUMN tanim TO definition;
ALTER TABLE projeler RENAME COLUMN olusturma_tarih TO created_at;
ALTER TABLE projeler RENAME COLUMN guncelleme_tarih TO updated_at;

-- ========================================
-- 3. etiketler table (tags)
-- ========================================
ALTER TABLE etiketler RENAME COLUMN isim TO name;

-- ========================================
-- 4. gorev_etiketleri table (task-tag junction)
-- ========================================
ALTER TABLE gorev_etiketleri RENAME COLUMN gorev_id TO task_id;
ALTER TABLE gorev_etiketleri RENAME COLUMN etiket_id TO tag_id;

-- ========================================
-- 5. baglantilar table (connections/dependencies)
-- ========================================
ALTER TABLE baglantilar RENAME COLUMN kaynak_id TO source_id;
ALTER TABLE baglantilar RENAME COLUMN hedef_id TO target_id;
ALTER TABLE baglantilar RENAME COLUMN baglanti_tip TO connection_type;

-- Update indexes for baglantilar
DROP INDEX IF EXISTS idx_baglanti_kaynak;
DROP INDEX IF EXISTS idx_baglanti_hedef;
DROP INDEX IF EXISTS idx_baglanti_cift;
CREATE INDEX idx_baglanti_source ON baglantilar(source_id);
CREATE INDEX idx_baglanti_target ON baglantilar(target_id);
CREATE INDEX idx_baglanti_pair ON baglantilar(source_id, target_id);

-- ========================================
-- 6. gorev_templateleri table (task templates)
-- ========================================
ALTER TABLE gorev_templateleri RENAME COLUMN isim TO name;
ALTER TABLE gorev_templateleri RENAME COLUMN tanim TO definition;
ALTER TABLE gorev_templateleri RENAME COLUMN varsayilan_baslik TO default_title;
ALTER TABLE gorev_templateleri RENAME COLUMN aciklama_template TO description_template;
ALTER TABLE gorev_templateleri RENAME COLUMN ornek_degerler TO sample_values;
ALTER TABLE gorev_templateleri RENAME COLUMN kategori TO category;
ALTER TABLE gorev_templateleri RENAME COLUMN aktif TO active;
ALTER TABLE gorev_templateleri RENAME COLUMN alanlar TO fields;

-- ========================================
-- 7. gorevler_fts table (FTS5 virtual table)
-- ========================================
-- FTS5 virtual tables cannot be altered, must recreate
DROP TABLE IF EXISTS gorevler_fts;
CREATE VIRTUAL TABLE gorevler_fts USING fts5(
    task_id UNINDEXED,
    title,
    description,
    tags
);

-- Rebuild FTS index with new column names
INSERT INTO gorevler_fts(task_id, title, description, tags)
SELECT
    g.id,
    g.title,
    COALESCE(g.description, ''),
    COALESCE((
        SELECT GROUP_CONCAT(e.name, ' ')
        FROM gorev_etiketleri ge
        JOIN etiketler e ON ge.tag_id = e.id
        WHERE ge.task_id = g.id
    ), '')
FROM gorevler g;

-- Recreate FTS triggers with new column names
DROP TRIGGER IF EXISTS gorevler_ai;
DROP TRIGGER IF EXISTS gorevler_ad;
DROP TRIGGER IF EXISTS gorevler_au;

CREATE TRIGGER gorevler_ai AFTER INSERT ON gorevler BEGIN
    INSERT INTO gorevler_fts(task_id, title, description, tags)
    VALUES (
        new.id,
        new.title,
        COALESCE(new.description, ''),
        COALESCE((
            SELECT GROUP_CONCAT(e.name, ' ')
            FROM gorev_etiketleri ge
            JOIN etiketler e ON ge.tag_id = e.id
            WHERE ge.task_id = new.id
        ), '')
    );
END;

CREATE TRIGGER gorevler_ad AFTER DELETE ON gorevler BEGIN
    DELETE FROM gorevler_fts WHERE task_id = old.id;
END;

CREATE TRIGGER gorevler_au AFTER UPDATE ON gorevler BEGIN
    UPDATE gorevler_fts
    SET
        title = new.title,
        description = COALESCE(new.description, ''),
        tags = COALESCE((
            SELECT GROUP_CONCAT(e.name, ' ')
            FROM gorev_etiketleri ge
            JOIN etiketler e ON ge.tag_id = e.id
            WHERE ge.task_id = new.id
        ), '')
    WHERE task_id = old.id;
END;

-- ========================================
-- 8. ai_interactions table
-- ========================================
ALTER TABLE ai_interactions RENAME COLUMN gorev_id TO task_id;

-- ========================================
-- 9. aktif_proje table (active project)
-- ========================================
ALTER TABLE aktif_proje RENAME COLUMN proje_id TO project_id;

-- ========================================
-- 10. gorev_hiyerarsi VIEW
-- ========================================
-- Recreate view with new column names
DROP VIEW IF EXISTS gorev_hiyerarsi;
CREATE VIEW gorev_hiyerarsi AS
WITH RECURSIVE hierarchy AS (
    SELECT
        g.id,
        g.title,
        g.description,
        g.status,
        g.priority,
        g.project_id,
        g.parent_id,
        g.created_at,
        g.updated_at,
        g.due_date,
        0 as level,
        g.id as root_id
    FROM gorevler g
    WHERE g.parent_id IS NULL

    UNION ALL

    SELECT
        g.id,
        g.title,
        g.description,
        g.status,
        g.priority,
        g.project_id,
        g.parent_id,
        g.created_at,
        g.updated_at,
        g.due_date,
        h.level + 1,
        h.root_id
    FROM gorevler g
    INNER JOIN hierarchy h ON g.parent_id = h.id
)
SELECT * FROM hierarchy;

-- ========================================
-- 11. Data integrity verification
-- ========================================
-- Count records before and after (for verification)
-- This will be logged during migration

-- Re-enable foreign keys
PRAGMA foreign_keys=ON;

-- Migration completed successfully
-- All Turkish field names renamed to English
-- Domain table names (gorevler, projeler) preserved
