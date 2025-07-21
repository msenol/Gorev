package gorev

import (
	"database/sql"
	"errors"
	"fmt"
	// "log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/msenol/gorev/internal/i18n"
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
		return fmt.Errorf(i18n.T("error.migrationDriverFailed", map[string]interface{}{"Error": err}))
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsYolu,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf(i18n.T("error.migrationInstanceFailed", map[string]interface{}{"Error": err}))
	}

	// Hata ayıklama için versiyonları logla
	_, _, err = m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		// //log.Printf("Migration öncesi versiyon alınamadı: %v", err)
	} else {
		// //log.Printf("Migration öncesi veritabanı versiyonu: %d, dirty: %v", version, dirty)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration işlemi başarısız: %w", err)
	}

	_, _, err = m.Version()
	if err != nil {
		// //log.Printf("Migration sonrası versiyon alınamadı: %v", err)
	} else {
		// //log.Printf("Migration sonrası veritabanı versiyonu: %d, dirty: %v", version, dirty)
	}

	// //log.Println("Veritabanı başarıyla migrate edildi.")

	// Varsayılan template'leri oluştur
	if err := vy.VarsayilanTemplateleriOlustur(); err != nil {
		// //log.Printf("Varsayılan template'ler oluşturulurken uyarı: %v", err)
		// Hata durumunda devam et, kritik değil
	}

	return nil
}

func (vy *VeriYonetici) Kapat() error {
	return vy.db.Close()
}

func (vy *VeriYonetici) GorevKaydet(gorev *Gorev) error {
	sorgu := `INSERT INTO gorevler (id, baslik, aciklama, durum, oncelik, proje_id, parent_id, olusturma_tarih, guncelleme_tarih, son_tarih)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := vy.db.Exec(sorgu,
		gorev.ID,
		gorev.Baslik,
		gorev.Aciklama,
		gorev.Durum,
		gorev.Oncelik,
		gorev.ProjeID,
		sql.NullString{String: gorev.ParentID, Valid: gorev.ParentID != ""},
		gorev.OlusturmaTarih,
		gorev.GuncellemeTarih,
		gorev.SonTarih,
	)

	return err
}

func (vy *VeriYonetici) GorevGetir(id string) (*Gorev, error) {
	sorgu := `SELECT id, baslik, aciklama, durum, oncelik, proje_id, parent_id, olusturma_tarih, guncelleme_tarih, son_tarih
	          FROM gorevler WHERE id = ?`

	gorev := &Gorev{}
	var projeID, parentID sql.NullString

	err := vy.db.QueryRow(sorgu, id).Scan(
		&gorev.ID,
		&gorev.Baslik,
		&gorev.Aciklama,
		&gorev.Durum,
		&gorev.Oncelik,
		&projeID,
		&parentID,
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

	if parentID.Valid {
		gorev.ParentID = parentID.String
	}

	// Etiketleri getir
	etiketler, err := vy.gorevEtiketleriniGetir(gorev.ID)
	if err != nil {
		//log.Printf("görev etiketleri getirilemedi: %v", err)
		// Etiket getirme başarısız olsa bile görevi döndür
		gorev.Etiketler = []*Etiket{}
	} else {
		gorev.Etiketler = etiketler
	}

	return gorev, nil
}

func (vy *VeriYonetici) GorevleriGetir(durum, sirala, filtre string) ([]*Gorev, error) {
	sorgu := `SELECT id, baslik, aciklama, durum, oncelik, proje_id, parent_id, olusturma_tarih, guncelleme_tarih, son_tarih
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
		var projeID, parentID sql.NullString

		err := rows.Scan(
			&gorev.ID,
			&gorev.Baslik,
			&gorev.Aciklama,
			&gorev.Durum,
			&gorev.Oncelik,
			&projeID,
			&parentID,
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

		if parentID.Valid {
			gorev.ParentID = parentID.String
		}

		// Etiketleri getir
		etiketler, err := vy.gorevEtiketleriniGetir(gorev.ID)
		if err != nil {
			// Hata durumunda logla ve devam et, görevi etiketsiz döndür
			//log.Printf("görev etiketleri getirilemedi: %v", err)
		}
		gorev.Etiketler = etiketler

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

func (vy *VeriYonetici) gorevEtiketleriniGetir(gorevID string) ([]*Etiket, error) {
	sorgu := `SELECT e.id, e.isim FROM etiketler e
	          JOIN gorev_etiketleri ge ON e.id = ge.etiket_id
	          WHERE ge.gorev_id = ?`
	rows, err := vy.db.Query(sorgu, gorevID)
	if err != nil {
		// Muhtemelen tablo yok, boş dön
		return []*Etiket{}, nil
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
				return nil, fmt.Errorf(i18n.T("error.tagCreateFailed", map[string]interface{}{"Tag": etiket.Isim, "Error": err}))
			}
		} else if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagQueryFailed", map[string]interface{}{"Tag": etiket.Isim, "Error": err}))
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
			return fmt.Errorf(i18n.T("error.taskTagAddFailed", map[string]interface{}{"Tag": etiket.Isim, "Error": err}))
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
		return fmt.Errorf(i18n.T("error.taskNotFound"))
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

// AltGorevleriGetir belirtilen görevin doğrudan alt görevlerini getirir
func (vy *VeriYonetici) AltGorevleriGetir(parentID string) ([]*Gorev, error) {
	sorgu := `SELECT id, baslik, aciklama, durum, oncelik, proje_id, parent_id, olusturma_tarih, guncelleme_tarih, son_tarih
	          FROM gorevler WHERE parent_id = ? ORDER BY olusturma_tarih`

	rows, err := vy.db.Query(sorgu, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gorevler []*Gorev
	for rows.Next() {
		gorev := &Gorev{}
		var projeID, parentID sql.NullString

		err := rows.Scan(
			&gorev.ID,
			&gorev.Baslik,
			&gorev.Aciklama,
			&gorev.Durum,
			&gorev.Oncelik,
			&projeID,
			&parentID,
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
		if parentID.Valid {
			gorev.ParentID = parentID.String
		}

		// Etiketleri getir
		etiketler, err := vy.gorevEtiketleriniGetir(gorev.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			//log.Printf("Etiketler getirilemedi: %v", err)
		}
		gorev.Etiketler = etiketler

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

// TumAltGorevleriGetir belirtilen görevin tüm alt görev hiyerarşisini getirir (recursive)
func (vy *VeriYonetici) TumAltGorevleriGetir(parentID string) ([]*Gorev, error) {
	sorgu := `
		WITH RECURSIVE alt_gorevler AS (
			SELECT id, baslik, aciklama, durum, oncelik, proje_id, parent_id, 
			       olusturma_tarih, guncelleme_tarih, son_tarih, 1 as seviye
			FROM gorevler
			WHERE parent_id = ?
			
			UNION ALL
			
			SELECT g.id, g.baslik, g.aciklama, g.durum, g.oncelik, g.proje_id, g.parent_id,
			       g.olusturma_tarih, g.guncelleme_tarih, g.son_tarih, ag.seviye + 1
			FROM gorevler g
			INNER JOIN alt_gorevler ag ON g.parent_id = ag.id
		)
		SELECT * FROM alt_gorevler ORDER BY seviye, olusturma_tarih`

	rows, err := vy.db.Query(sorgu, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gorevler []*Gorev
	for rows.Next() {
		gorev := &Gorev{}
		var projeID, parentID sql.NullString

		err := rows.Scan(
			&gorev.ID,
			&gorev.Baslik,
			&gorev.Aciklama,
			&gorev.Durum,
			&gorev.Oncelik,
			&projeID,
			&parentID,
			&gorev.OlusturmaTarih,
			&gorev.GuncellemeTarih,
			&gorev.SonTarih,
			&gorev.Seviye,
		)
		if err != nil {
			return nil, err
		}

		if projeID.Valid {
			gorev.ProjeID = projeID.String
		}
		if parentID.Valid {
			gorev.ParentID = parentID.String
		}

		// Etiketleri getir
		etiketler, err := vy.gorevEtiketleriniGetir(gorev.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			//log.Printf("Etiketler getirilemedi: %v", err)
		}
		gorev.Etiketler = etiketler

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

// UstGorevleriGetir belirtilen görevin tüm üst görev hiyerarşisini getirir
func (vy *VeriYonetici) UstGorevleriGetir(gorevID string) ([]*Gorev, error) {
	sorgu := `
		WITH RECURSIVE ust_gorevler AS (
			SELECT g2.id, g2.baslik, g2.aciklama, g2.durum, g2.oncelik, g2.proje_id, g2.parent_id,
			       g2.olusturma_tarih, g2.guncelleme_tarih, g2.son_tarih
			FROM gorevler g1
			JOIN gorevler g2 ON g1.parent_id = g2.id
			WHERE g1.id = ?
			
			UNION ALL
			
			SELECT g.id, g.baslik, g.aciklama, g.durum, g.oncelik, g.proje_id, g.parent_id,
			       g.olusturma_tarih, g.guncelleme_tarih, g.son_tarih
			FROM gorevler g
			INNER JOIN ust_gorevler ug ON ug.parent_id = g.id
		)
		SELECT * FROM ust_gorevler`

	rows, err := vy.db.Query(sorgu, gorevID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gorevler []*Gorev
	for rows.Next() {
		gorev := &Gorev{}
		var projeID, parentID sql.NullString

		err := rows.Scan(
			&gorev.ID,
			&gorev.Baslik,
			&gorev.Aciklama,
			&gorev.Durum,
			&gorev.Oncelik,
			&projeID,
			&parentID,
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
		if parentID.Valid {
			gorev.ParentID = parentID.String
		}

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

// GorevHiyerarsiGetir bir görevin tam hiyerarşi bilgilerini getirir
func (vy *VeriYonetici) GorevHiyerarsiGetir(gorevID string) (*GorevHiyerarsi, error) {
	// Ana görevi getir
	gorev, err := vy.GorevGetir(gorevID)
	if err != nil {
		return nil, err
	}

	// Üst görevleri getir
	ustGorevler, err := vy.UstGorevleriGetir(gorevID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Alt görev istatistiklerini hesapla
	sorgu := `
		WITH RECURSIVE alt_gorevler AS (
			SELECT id, durum
			FROM gorevler
			WHERE parent_id = ?
			
			UNION ALL
			
			SELECT g.id, g.durum
			FROM gorevler g
			INNER JOIN alt_gorevler ag ON g.parent_id = ag.id
		)
		SELECT 
			COUNT(*) as toplam,
			COALESCE(SUM(CASE WHEN durum = 'tamamlandi' THEN 1 ELSE 0 END), 0) as tamamlanan,
			COALESCE(SUM(CASE WHEN durum = 'devam_ediyor' THEN 1 ELSE 0 END), 0) as devam_eden,
			COALESCE(SUM(CASE WHEN durum = 'beklemede' THEN 1 ELSE 0 END), 0) as beklemede
		FROM alt_gorevler`

	var toplam, tamamlanan, devamEden, beklemede int
	err = vy.db.QueryRow(sorgu, gorevID).Scan(&toplam, &tamamlanan, &devamEden, &beklemede)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// İlerleme yüzdesini hesapla
	var ilerlemeYuzdesi float64
	if toplam > 0 {
		ilerlemeYuzdesi = (float64(tamamlanan) / float64(toplam)) * 100
	} else if gorev.Durum == "tamamlandi" {
		ilerlemeYuzdesi = 100
	}

	return &GorevHiyerarsi{
		Gorev:           gorev,
		UstGorevler:     ustGorevler,
		ToplamAltGorev:  toplam,
		TamamlananAlt:   tamamlanan,
		DevamEdenAlt:    devamEden,
		BeklemedeAlt:    beklemede,
		IlerlemeYuzdesi: ilerlemeYuzdesi,
	}, nil
}

// ParentIDGuncelle bir görevin parent_id'sini günceller
func (vy *VeriYonetici) ParentIDGuncelle(gorevID, yeniParentID string) error {
	// Önce circular dependency kontrolü yap
	if yeniParentID != "" {
		daireVar, err := vy.DaireBagimliligiKontrolEt(gorevID, yeniParentID)
		if err != nil {
			return err
		}
		if daireVar {
			return fmt.Errorf("dairesel bağımlılık tespit edildi")
		}
	}

	var sorgu string
	var err error

	if yeniParentID == "" {
		sorgu = `UPDATE gorevler SET parent_id = NULL, guncelleme_tarih = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = vy.db.Exec(sorgu, gorevID)
	} else {
		sorgu = `UPDATE gorevler SET parent_id = ?, guncelleme_tarih = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = vy.db.Exec(sorgu, yeniParentID, gorevID)
	}

	return err
}

// DaireBagimliligiKontrolEt bir görevin belirtilen parent'a taşınması durumunda dairesel bağımlılık oluşup oluşmayacağını kontrol eder
func (vy *VeriYonetici) DaireBagimliligiKontrolEt(gorevID, hedefParentID string) (bool, error) {
	// Kendisine parent olamaz
	if gorevID == hedefParentID {
		return true, nil
	}

	// Hedef parent'ın üst hiyerarşisinde gorevID var mı kontrol et
	sorgu := `
		WITH RECURSIVE ust_gorevler AS (
			SELECT id, parent_id
			FROM gorevler
			WHERE id = ?
			
			UNION ALL
			
			SELECT g.id, g.parent_id
			FROM gorevler g
			INNER JOIN ust_gorevler ug ON g.id = ug.parent_id
		)
		SELECT COUNT(*) FROM ust_gorevler WHERE id = ?`

	var count int
	err := vy.db.QueryRow(sorgu, hedefParentID, gorevID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// BulkBagimlilikSayilariGetir tüm görevlerin bağımlılık sayılarını tek sorguda hesaplar
// Bu N+1 sorgu problemini çözer ve performansı büyük ölçüde artırır
func (vy *VeriYonetici) BulkBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error) {
	if len(gorevIDs) == 0 {
		return make(map[string]int), nil
	}

	// Placeholder'ları oluştur
	placeholders := make([]string, len(gorevIDs))
	args := make([]interface{}, len(gorevIDs))
	for i, id := range gorevIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	sorgu := fmt.Sprintf(`
		SELECT kaynak_id, COUNT(*) as bagli_sayi
		FROM baglantilar 
		WHERE kaynak_id IN (%s)
		GROUP BY kaynak_id
	`, strings.Join(placeholders, ","))

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var gorevID string
		var count int
		if err := rows.Scan(&gorevID, &count); err != nil {
			return nil, err
		}
		result[gorevID] = count
	}

	return result, nil
}

// BulkTamamlanmamiaBagimlilikSayilariGetir tüm görevlerin tamamlanmamış bağımlılık sayılarını hesaplar
func (vy *VeriYonetici) BulkTamamlanmamiaBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error) {
	if len(gorevIDs) == 0 {
		return make(map[string]int), nil
	}

	// Placeholder'ları oluştur
	placeholders := make([]string, len(gorevIDs))
	args := make([]interface{}, len(gorevIDs))
	for i, id := range gorevIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	sorgu := fmt.Sprintf(`
		SELECT b.kaynak_id, COUNT(*) as tamamlanmamis_sayi
		FROM baglantilar b
		INNER JOIN gorevler g ON b.hedef_id = g.id
		WHERE b.kaynak_id IN (%s) AND g.durum != 'tamamlandi'
		GROUP BY b.kaynak_id
	`, strings.Join(placeholders, ","))

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var gorevID string
		var count int
		if err := rows.Scan(&gorevID, &count); err != nil {
			return nil, err
		}
		result[gorevID] = count
	}

	return result, nil
}
