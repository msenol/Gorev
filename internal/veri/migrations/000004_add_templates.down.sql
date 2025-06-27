-- Template index'lerini kaldır
DROP INDEX IF EXISTS idx_template_aktif;
DROP INDEX IF EXISTS idx_template_kategori;

-- Template tablosunu kaldır
DROP TABLE IF EXISTS gorev_templateleri;