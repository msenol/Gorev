CREATE TABLE projeler (
    id TEXT PRIMARY KEY,
    isim TEXT NOT NULL,
    tanim TEXT,
    olusturma_tarih DATETIME NOT NULL,
    guncelleme_tarih DATETIME NOT NULL
);

CREATE TABLE gorevler (
    id TEXT PRIMARY KEY,
    baslik TEXT NOT NULL,
    aciklama TEXT,
    durum TEXT NOT NULL DEFAULT 'beklemede',
    oncelik TEXT NOT NULL DEFAULT 'orta',
    proje_id TEXT,
    olusturma_tarih DATETIME NOT NULL,
    guncelleme_tarih DATETIME NOT NULL,
    FOREIGN KEY (proje_id) REFERENCES projeler(id)
);

CREATE TABLE baglantilar (
    id TEXT PRIMARY KEY,
    kaynak_id TEXT NOT NULL,
    hedef_id TEXT NOT NULL,
    baglanti_tip TEXT NOT NULL,
    FOREIGN KEY (kaynak_id) REFERENCES gorevler(id),
    FOREIGN KEY (hedef_id) REFERENCES gorevler(id)
);

CREATE INDEX idx_gorev_durum ON gorevler(durum);
CREATE INDEX idx_gorev_proje ON gorevler(proje_id);

CREATE TABLE aktif_proje (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    proje_id TEXT NOT NULL,
    FOREIGN KEY (proje_id) REFERENCES projeler(id)
);
