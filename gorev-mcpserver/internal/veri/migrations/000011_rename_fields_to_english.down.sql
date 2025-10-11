-- Rollback Migration: Revert English field names back to Turkish
-- Version: v0.17.0
-- Date: 2025-10-11
-- Purpose: Rollback schema changes if needed

-- ========================================
-- 1. gorevler table (tasks)
-- ========================================
ALTER TABLE gorevler RENAME COLUMN title TO baslik;
ALTER TABLE gorevler RENAME COLUMN description TO aciklama;
ALTER TABLE gorevler RENAME COLUMN status TO durum;
ALTER TABLE gorevler RENAME COLUMN priority TO oncelik;
ALTER TABLE gorevler RENAME COLUMN created_at TO olusturma_tarih;
ALTER TABLE gorevler RENAME COLUMN updated_at TO guncelleme_tarih;
ALTER TABLE gorevler RENAME COLUMN due_date TO son_tarih;

-- Restore original indexes
DROP INDEX IF EXISTS idx_gorev_status;
DROP INDEX IF EXISTS idx_gorev_proje;
DROP INDEX IF EXISTS idx_gorev_priority;
DROP INDEX IF EXISTS idx_gorev_parent;
CREATE INDEX idx_gorev_durum ON gorevler(durum);
CREATE INDEX idx_gorev_proje ON gorevler(proje_id);
CREATE INDEX idx_gorev_oncelik ON gorevler(oncelik);
CREATE INDEX idx_gorev_parent ON gorevler(parent_id);

-- ========================================
-- 2. projeler table (projects)
-- ========================================
ALTER TABLE projeler RENAME COLUMN name TO isim;
ALTER TABLE projeler RENAME COLUMN definition TO tanim;
ALTER TABLE projeler RENAME COLUMN created_at TO olusturma_tarih;
ALTER TABLE projeler RENAME COLUMN updated_at TO guncelleme_tarih;

-- ========================================
-- 3. etiketler table (tags)
-- ========================================
ALTER TABLE etiketler RENAME COLUMN name TO isim;

-- ========================================
-- 4. baglantilar table (connections/dependencies)
-- ========================================
ALTER TABLE baglantilar RENAME COLUMN source_id TO kaynak_id;
ALTER TABLE baglantilar RENAME COLUMN target_id TO hedef_id;
ALTER TABLE baglantilar RENAME COLUMN connection_type TO baglanti_tip;

-- Restore original indexes
DROP INDEX IF EXISTS idx_baglanti_source;
DROP INDEX IF EXISTS idx_baglanti_target;
DROP INDEX IF EXISTS idx_baglanti_pair;
CREATE INDEX idx_baglanti_kaynak ON baglantilar(kaynak_id);
CREATE INDEX idx_baglanti_hedef ON baglantilar(hedef_id);
CREATE INDEX idx_baglanti_cift ON baglantilar(kaynak_id, hedef_id);

-- ========================================
-- 5. gorev_templateleri table (task templates)
-- ========================================
ALTER TABLE gorev_templateleri RENAME COLUMN name TO isim;
ALTER TABLE gorev_templateleri RENAME COLUMN definition TO tanim;
ALTER TABLE gorev_templateleri RENAME COLUMN default_title TO varsayilan_baslik;
ALTER TABLE gorev_templateleri RENAME COLUMN description_template TO aciklama_template;
ALTER TABLE gorev_templateleri RENAME COLUMN sample_values TO ornek_degerler;
ALTER TABLE gorev_templateleri RENAME COLUMN category TO kategori;
ALTER TABLE gorev_templateleri RENAME COLUMN active TO aktif;
ALTER TABLE gorev_templateleri RENAME COLUMN fields TO alanlar;

-- ========================================
-- 6. gorevler_fts table (FTS5 virtual table)
-- ========================================
DROP TABLE IF EXISTS gorevler_fts;
CREATE VIRTUAL TABLE gorevler_fts USING fts5(
    task_id UNINDEXED,
    baslik,
    aciklama,
    etiketler,
    content='gorevler',
    content_rowid='rowid'
);

-- Rebuild FTS index with original column names
INSERT INTO gorevler_fts(task_id, baslik, aciklama, etiketler)
SELECT
    g.id,
    g.baslik,
    COALESCE(g.aciklama, ''),
    COALESCE((
        SELECT GROUP_CONCAT(e.isim, ' ')
        FROM gorev_etiketleri ge
        JOIN etiketler e ON ge.etiket_id = e.id
        WHERE ge.gorev_id = g.id
    ), '')
FROM gorevler g;

-- Recreate FTS triggers with original column names
DROP TRIGGER IF EXISTS gorevler_ai;
DROP TRIGGER IF EXISTS gorevler_ad;
DROP TRIGGER IF EXISTS gorevler_au;

CREATE TRIGGER gorevler_ai AFTER INSERT ON gorevler BEGIN
    INSERT INTO gorevler_fts(task_id, baslik, aciklama, etiketler)
    VALUES (
        new.id,
        new.baslik,
        COALESCE(new.aciklama, ''),
        COALESCE((
            SELECT GROUP_CONCAT(e.isim, ' ')
            FROM gorev_etiketleri ge
            JOIN etiketler e ON ge.etiket_id = e.id
            WHERE ge.gorev_id = new.id
        ), '')
    );
END;

CREATE TRIGGER gorevler_ad AFTER DELETE ON gorevler BEGIN
    DELETE FROM gorevler_fts WHERE task_id = old.id;
END;

CREATE TRIGGER gorevler_au AFTER UPDATE ON gorevler BEGIN
    UPDATE gorevler_fts
    SET
        baslik = new.baslik,
        aciklama = COALESCE(new.aciklama, ''),
        etiketler = COALESCE((
            SELECT GROUP_CONCAT(e.isim, ' ')
            FROM gorev_etiketleri ge
            JOIN etiketler e ON ge.etiket_id = e.id
            WHERE ge.gorev_id = new.id
        ), '')
    WHERE task_id = old.id;
END;

-- ========================================
-- 7. gorev_hiyerarsi VIEW
-- ========================================
DROP VIEW IF EXISTS gorev_hiyerarsi;
CREATE VIEW gorev_hiyerarsi AS
WITH RECURSIVE hierarchy AS (
    SELECT
        g.id,
        g.baslik,
        g.aciklama,
        g.durum,
        g.oncelik,
        g.proje_id,
        g.parent_id,
        g.olusturma_tarih,
        g.guncelleme_tarih,
        g.son_tarih,
        0 as seviye,
        g.id as kok_id
    FROM gorevler g
    WHERE g.parent_id IS NULL

    UNION ALL

    SELECT
        g.id,
        g.baslik,
        g.aciklama,
        g.durum,
        g.oncelik,
        g.proje_id,
        g.parent_id,
        g.olusturma_tarih,
        g.guncelleme_tarih,
        g.son_tarih,
        h.seviye + 1,
        h.kok_id
    FROM gorevler g
    INNER JOIN hierarchy h ON g.parent_id = h.id
)
SELECT * FROM hierarchy;

-- Rollback completed successfully
-- All English field names reverted to Turkish
