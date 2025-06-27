-- Görev template'leri tablosu
CREATE TABLE IF NOT EXISTS gorev_templateleri (
    id TEXT PRIMARY KEY,
    isim TEXT NOT NULL UNIQUE,
    tanim TEXT NOT NULL,
    varsayilan_baslik TEXT NOT NULL,
    aciklama_template TEXT NOT NULL,
    alanlar TEXT NOT NULL, -- JSON array of TemplateAlan
    ornek_degerler TEXT NOT NULL, -- JSON object
    kategori TEXT NOT NULL,
    aktif BOOLEAN NOT NULL DEFAULT 1,
    olusturma_tarih TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    guncelleme_tarih TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Kategori index'i
CREATE INDEX IF NOT EXISTS idx_template_kategori ON gorev_templateleri(kategori);

-- Aktif template'ler index'i
CREATE INDEX IF NOT EXISTS idx_template_aktif ON gorev_templateleri(aktif);