-- Etiketler için yeni bir tablo oluştur
CREATE TABLE etiketler (
    id TEXT PRIMARY KEY,
    isim TEXT NOT NULL UNIQUE
);

-- Görevler ve etiketler arasında bir birleştirme tablosu oluştur
CREATE TABLE gorev_etiketleri (
    gorev_id TEXT NOT NULL,
    etiket_id TEXT NOT NULL,
    PRIMARY KEY (gorev_id, etiket_id),
    FOREIGN KEY (gorev_id) REFERENCES gorevler(id) ON DELETE CASCADE,
    FOREIGN KEY (etiket_id) REFERENCES etiketler(id) ON DELETE CASCADE
);

-- Etiket ismine göre hızlı arama için bir indeks oluştur
CREATE INDEX idx_etiket_isim ON etiketler(isim);
