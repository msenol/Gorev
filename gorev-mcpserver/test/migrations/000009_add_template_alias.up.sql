-- Template'lere alias kolonu ekleme
ALTER TABLE gorev_templateleri ADD COLUMN alias TEXT;

-- Alias kolonu için unique index (NULL değerler dahil edilmez)
CREATE UNIQUE INDEX IF NOT EXISTS idx_template_alias ON gorev_templateleri(alias) WHERE alias IS NOT NULL;

-- Mevcut template'lere alias değerleri ata (spesifik isim eşleştirmesi ile)
UPDATE gorev_templateleri SET alias = 'bug' WHERE isim = 'Bug Raporu' AND alias IS NULL;
UPDATE gorev_templateleri SET alias = 'bug2' WHERE isim = 'Bug Raporu v2' AND alias IS NULL;
UPDATE gorev_templateleri SET alias = 'feature' WHERE isim = 'Özellik İsteği' AND alias IS NULL;
UPDATE gorev_templateleri SET alias = 'research' WHERE isim = 'Araştırma Görevi' AND alias IS NULL;
UPDATE gorev_templateleri SET alias = 'spike' WHERE isim = 'Spike Araştırma' AND alias IS NULL;
UPDATE gorev_templateleri SET alias = 'security' WHERE isim = 'Güvenlik Düzeltmesi' AND alias IS NULL;
UPDATE gorev_templateleri SET alias = 'performance' WHERE isim = 'Performans Sorunu' AND alias IS NULL;
UPDATE gorev_templateleri SET alias = 'refactor' WHERE isim = 'Refactoring' AND alias IS NULL;
UPDATE gorev_templateleri SET alias = 'debt' WHERE isim = 'Teknik Borç' AND alias IS NULL;