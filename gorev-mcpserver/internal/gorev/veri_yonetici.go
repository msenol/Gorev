package gorev

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type VeriYonetici struct {
	db *sql.DB
}

func YeniVeriYonetici(dbYolu string, migrationsYolu string) (*VeriYonetici, error) {
	db, err := sql.Open("sqlite3", dbYolu)
	if err != nil {
		return nil, fmt.Errorf("veritabanı açılamadı: %w", err)
	}

	vy := &VeriYonetici{db: db}
	if err := vy.migrateDB(migrationsYolu); err != nil {
		return nil, fmt.Errorf("veritabanı migrate edilemedi: %w", err)
	}

	return vy, nil
}

func (vy *VeriYonetici) migrateDB(migrationsYolu string) error {
	driver, err := sqlite3.WithInstance(vy.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("migration driver oluşturulamadı: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsYolu,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migration instance oluşturulamadı: %w", err)
	}

	// Hata ayıklama için versiyonları logla
	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("Migration öncesi versiyon alınamadı: %v", err)
	} else {
		log.Printf("Migration öncesi veritabanı versiyonu: %d, dirty: %v", version, dirty)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration işlemi başarısız: %w", err)
	}

	version, dirty, err = m.Version()
	if err != nil {
		log.Printf("Migration sonrası versiyon alınamadı: %v", err)
	} else {
		log.Printf("Migration sonrası veritabanı versiyonu: %d, dirty: %v", version, dirty)
	}

	log.Println("Veritabanı başarıyla migrate edildi.")

	// Varsayılan template'leri oluştur
	if err := vy.VarsayilanTemplateleriOlustur(); err != nil {
		log.Printf("Varsayılan template'ler oluşturulurken uyarı: %v", err)
		// Hata durumunda devam et, kritik değil
	}

	return nil
}

func (vy *VeriYonetici) Kapat() error {
	return vy.db.Close()
}

func (vy *VeriYonetici) GorevKaydet(gorev *Gorev) error {
	sorgu := `INSERT INTO gorevler (id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih, son_tarih)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := vy.db.Exec(sorgu,
		gorev.ID,
		gorev.Baslik,
		gorev.Aciklama,
		gorev.Durum,
		gorev.Oncelik,
		gorev.ProjeID,
		gorev.OlusturmaTarih,
		gorev.GuncellemeTarih,
		gorev.SonTarih,
	)

	return err
}

func (vy *VeriYonetici) GorevGetir(id string) (*Gorev, error) {
	sorgu := `SELECT id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih, son_tarih
	          FROM gorevler WHERE id = ?`

	gorev := &Gorev{}
	var projeID sql.NullString

	err := vy.db.QueryRow(sorgu, id).Scan(
		&gorev.ID,
		&gorev.Baslik,
		&gorev.Aciklama,
		&gorev.Durum,
		&gorev.Oncelik,
		&projeID,
		&gorev.OlusturmaTarih,
		&gorev.GuncellemeTarih,
		&gorev.SonTarih,
	)

	if err != nil {
		return nil, err
	}

	if projeID.Valid {
		gorev.ProjeID = projeID.String
	}

	// Etiketleri getir
	etiketler, err := vy.gorevEtiketleriniGetir(gorev.ID)
	if err != nil {
		log.Printf("görev etiketleri getirilemedi: %v", err)
		// Etiket getirme başarısız olsa bile görevi döndür
		gorev.Etiketler = []*Etiket{}
	} else {
		gorev.Etiketler = etiketler
	}

	return gorev, nil
}

func (vy *VeriYonetici) GorevleriGetir(durum, sirala, filtre string) ([]*Gorev, error) {
	sorgu := `SELECT id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih, son_tarih
	          FROM gorevler`
	args := []interface{}{}
	whereClauses := []string{}

	if durum != "" {
		whereClauses = append(whereClauses, "durum = ?")
		args = append(args, durum)
	}

	if filtre == "acil" {
		whereClauses = append(whereClauses, "son_tarih IS NOT NULL AND son_tarih >= date('now') AND son_tarih < date('now', '+7 days')")
	} else if filtre == "gecmis" {
		whereClauses = append(whereClauses, "son_tarih IS NOT NULL AND son_tarih < date('now')")
	}

	if len(whereClauses) > 0 {
		sorgu += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	switch sirala {
	case "son_tarih_asc":
		sorgu += " ORDER BY son_tarih ASC"
	case "son_tarih_desc":
		sorgu += " ORDER BY son_tarih DESC"
	default:
		sorgu += " ORDER BY olusturma_tarih DESC"
	}

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gorevler []*Gorev
	for rows.Next() {
		gorev := &Gorev{}
		var projeID sql.NullString

		err := rows.Scan(
			&gorev.ID,
			&gorev.Baslik,
			&gorev.Aciklama,
			&gorev.Durum,
			&gorev.Oncelik,
			&projeID,
			&gorev.OlusturmaTarih,
			&gorev.GuncellemeTarih,
			&gorev.SonTarih,
		)
		if err != nil {
			return nil, err
		}

		if projeID.Valid {
			gorev.ProjeID = projeID.String
		}

		// Etiketleri getir
		etiketler, err := vy.gorevEtiketleriniGetir(gorev.ID)
		if err != nil {
			// Hata durumunda logla ve devam et, görevi etiketsiz döndür
			log.Printf("görev etiketleri getirilemedi: %v", err)
		}
		gorev.Etiketler = etiketler

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

func (vy *VeriYonetici) gorevEtiketleriniGetir(gorevID string) ([]*Etiket, error) {
	// First check if etiketler table exists
	var tableExists int
	err := vy.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='etiketler'").Scan(&tableExists)
	if err != nil || tableExists == 0 {
		// Table doesn't exist, return empty slice instead of error
		return []*Etiket{}, nil
	}

	sorgu := `SELECT e.id, e.isim FROM etiketler e
	          JOIN gorev_etiketleri ge ON e.id = ge.etiket_id
	          WHERE ge.gorev_id = ?`
	rows, err := vy.db.Query(sorgu, gorevID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var etiketler []*Etiket
	for rows.Next() {
		e := &Etiket{}
		if err := rows.Scan(&e.ID, &e.Isim); err != nil {
			return nil, err
		}
		etiketler = append(etiketler, e)
	}
	return etiketler, nil
}

func (vy *VeriYonetici) EtiketleriGetirVeyaOlustur(isimler []string) ([]*Etiket, error) {
	etiketler := make([]*Etiket, 0, len(isimler))
	tx, err := vy.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction başlatılamadı: %w", err)
	}
	defer tx.Rollback() // Hata durumunda geri al

	stmtSelect, err := tx.Prepare("SELECT id, isim FROM etiketler WHERE isim = ?")
	if err != nil {
		return nil, fmt.Errorf("select statement hazırlanamadı: %w", err)
	}
	defer stmtSelect.Close()

	stmtInsert, err := tx.Prepare("INSERT INTO etiketler (id, isim) VALUES (?, ?)")
	if err != nil {
		return nil, fmt.Errorf("insert statement hazırlanamadı: %w", err)
	}
	defer stmtInsert.Close()

	for _, isim := range isimler {
		if strings.TrimSpace(isim) == "" {
			continue
		}
		etiket := &Etiket{Isim: strings.TrimSpace(isim)}
		err := stmtSelect.QueryRow(etiket.Isim).Scan(&etiket.ID, &etiket.Isim)
		if err == sql.ErrNoRows {
			// Etiket yok, oluştur
			etiket.ID = uuid.New().String()
			if _, err := stmtInsert.Exec(etiket.ID, etiket.Isim); err != nil {
				return nil, fmt.Errorf("yeni etiket oluşturulamadı '%s': %w", etiket.Isim, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("etiket sorgulanamadı '%s': %w", etiket.Isim, err)
		}
		etiketler = append(etiketler, etiket)
	}

	return etiketler, tx.Commit()
}

func (vy *VeriYonetici) GorevEtiketleriniAyarla(gorevID string, etiketler []*Etiket) error {
	tx, err := vy.db.Begin()
	if err != nil {
		return fmt.Errorf("transaction başlatılamadı: %w", err)
	}
	defer tx.Rollback()

	// Mevcut bağlantıları sil
	if _, err := tx.Exec("DELETE FROM gorev_etiketleri WHERE gorev_id = ?", gorevID); err != nil {
		return fmt.Errorf("mevcut etiketler silinemedi: %w", err)
	}

	// Yeni bağlantıları ekle
	stmt, err := tx.Prepare("INSERT INTO gorev_etiketleri (gorev_id, etiket_id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("insert statement hazırlanamadı: %w", err)
	}
	defer stmt.Close()

	for _, etiket := range etiketler {
		if _, err := stmt.Exec(gorevID, etiket.ID); err != nil {
			return fmt.Errorf("görev etiketi eklenemedi '%s': %w", etiket.Isim, err)
		}
	}

	return tx.Commit()
}

func (vy *VeriYonetici) GorevGuncelle(gorev *Gorev) error {
	sorgu := `UPDATE gorevler SET baslik = ?, aciklama = ?, durum = ?, oncelik = ?, 
	          proje_id = ?, guncelleme_tarih = ?, son_tarih = ? WHERE id = ?`

	_, err := vy.db.Exec(sorgu,
		gorev.Baslik,
		gorev.Aciklama,
		gorev.Durum,
		gorev.Oncelik,
		gorev.ProjeID,
		gorev.GuncellemeTarih,
		gorev.SonTarih,
		gorev.ID,
	)

	return err
}

func (vy *VeriYonetici) ProjeKaydet(proje *Proje) error {
	sorgu := `INSERT INTO projeler (id, isim, tanim, olusturma_tarih, guncelleme_tarih)
	          VALUES (?, ?, ?, ?, ?)`

	_, err := vy.db.Exec(sorgu,
		proje.ID,
		proje.Isim,
		proje.Tanim,
		proje.OlusturmaTarih,
		proje.GuncellemeTarih,
	)

	return err
}

func (vy *VeriYonetici) ProjeGetir(id string) (*Proje, error) {
	sorgu := `SELECT id, isim, tanim, olusturma_tarih, guncelleme_tarih
	          FROM projeler WHERE id = ?`

	proje := &Proje{}
	err := vy.db.QueryRow(sorgu, id).Scan(
		&proje.ID,
		&proje.Isim,
		&proje.Tanim,
		&proje.OlusturmaTarih,
		&proje.GuncellemeTarih,
	)

	if err != nil {
		return nil, err
	}

	return proje, nil
}

func (vy *VeriYonetici) ProjeleriGetir() ([]*Proje, error) {
	sorgu := `SELECT id, isim, tanim, olusturma_tarih, guncelleme_tarih
	          FROM projeler ORDER BY olusturma_tarih DESC`

	rows, err := vy.db.Query(sorgu)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projeler []*Proje
	for rows.Next() {
		proje := &Proje{}
		err := rows.Scan(
			&proje.ID,
			&proje.Isim,
			&proje.Tanim,
			&proje.OlusturmaTarih,
			&proje.GuncellemeTarih,
		)
		if err != nil {
			return nil, err
		}
		projeler = append(projeler, proje)
	}

	return projeler, nil
}

func (vy *VeriYonetici) GorevSil(id string) error {
	sorgu := `DELETE FROM gorevler WHERE id = ?`

	result, err := vy.db.Exec(sorgu, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("görev bulunamadı")
	}

	return nil
}

func (vy *VeriYonetici) ProjeGorevleriGetir(projeID string) ([]*Gorev, error) {
	sorgu := `SELECT id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih
	          FROM gorevler WHERE proje_id = ? ORDER BY olusturma_tarih DESC`

	rows, err := vy.db.Query(sorgu, projeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gorevler []*Gorev
	for rows.Next() {
		gorev := &Gorev{}
		var pID sql.NullString

		err := rows.Scan(
			&gorev.ID,
			&gorev.Baslik,
			&gorev.Aciklama,
			&gorev.Durum,
			&gorev.Oncelik,
			&pID,
			&gorev.OlusturmaTarih,
			&gorev.GuncellemeTarih,
		)
		if err != nil {
			return nil, err
		}

		if pID.Valid {
			gorev.ProjeID = pID.String
		}

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

func (vy *VeriYonetici) BaglantiEkle(baglanti *Baglanti) error {
	sorgu := `INSERT INTO baglantilar (id, kaynak_id, hedef_id, baglanti_tip) VALUES (?, ?, ?, ?)`
	_, err := vy.db.Exec(sorgu, baglanti.ID, baglanti.KaynakID, baglanti.HedefID, baglanti.BaglantiTip)
	return err
}

func (vy *VeriYonetici) BaglantilariGetir(gorevID string) ([]*Baglanti, error) {
	sorgu := `SELECT id, kaynak_id, hedef_id, baglanti_tip FROM baglantilar WHERE kaynak_id = ? OR hedef_id = ?`
	rows, err := vy.db.Query(sorgu, gorevID, gorevID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var baglantilar []*Baglanti
	for rows.Next() {
		b := &Baglanti{}
		if err := rows.Scan(&b.ID, &b.KaynakID, &b.HedefID, &b.BaglantiTip); err != nil {
			return nil, err
		}
		baglantilar = append(baglantilar, b)
	}
	return baglantilar, nil
}
