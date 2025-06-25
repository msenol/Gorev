package gorev

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type VeriYonetici struct {
	db *sql.DB
}

func YeniVeriYonetici(dbYolu string) (*VeriYonetici, error) {
	db, err := sql.Open("sqlite3", dbYolu)
	if err != nil {
		return nil, fmt.Errorf("veritabanı açılamadı: %w", err)
	}

	vy := &VeriYonetici{db: db}
	if err := vy.tablolariOlustur(); err != nil {
		return nil, fmt.Errorf("tablolar oluşturulamadı: %w", err)
	}

	return vy, nil
}

func (vy *VeriYonetici) Kapat() error {
	return vy.db.Close()
}

func (vy *VeriYonetici) tablolariOlustur() error {
	sorgular := []string{
		`CREATE TABLE IF NOT EXISTS projeler (
			id TEXT PRIMARY KEY,
			isim TEXT NOT NULL,
			tanim TEXT,
			olusturma_tarih DATETIME NOT NULL,
			guncelleme_tarih DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS gorevler (
			id TEXT PRIMARY KEY,
			baslik TEXT NOT NULL,
			aciklama TEXT,
			durum TEXT NOT NULL DEFAULT 'beklemede',
			oncelik TEXT NOT NULL DEFAULT 'orta',
			proje_id TEXT,
			olusturma_tarih DATETIME NOT NULL,
			guncelleme_tarih DATETIME NOT NULL,
			FOREIGN KEY (proje_id) REFERENCES projeler(id)
		)`,
		`CREATE TABLE IF NOT EXISTS baglantilar (
			id TEXT PRIMARY KEY,
			kaynak_id TEXT NOT NULL,
			hedef_id TEXT NOT NULL,
			baglanti_tip TEXT NOT NULL,
			FOREIGN KEY (kaynak_id) REFERENCES gorevler(id),
			FOREIGN KEY (hedef_id) REFERENCES gorevler(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_gorev_durum ON gorevler(durum)`,
		`CREATE INDEX IF NOT EXISTS idx_gorev_proje ON gorevler(proje_id)`,
	}

	for _, sorgu := range sorgular {
		if _, err := vy.db.Exec(sorgu); err != nil {
			return fmt.Errorf("sorgu çalıştırılamadı: %w", err)
		}
	}

	return nil
}

func (vy *VeriYonetici) GorevKaydet(gorev *Gorev) error {
	sorgu := `INSERT INTO gorevler (id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := vy.db.Exec(sorgu,
		gorev.ID,
		gorev.Baslik,
		gorev.Aciklama,
		gorev.Durum,
		gorev.Oncelik,
		gorev.ProjeID,
		gorev.OlusturmaTarih,
		gorev.GuncellemeTarih,
	)

	return err
}

func (vy *VeriYonetici) GorevGetir(id string) (*Gorev, error) {
	sorgu := `SELECT id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih
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
	)

	if err != nil {
		return nil, err
	}

	if projeID.Valid {
		gorev.ProjeID = projeID.String
	}

	return gorev, nil
}

func (vy *VeriYonetici) GorevleriGetir(durum string) ([]*Gorev, error) {
	sorgu := `SELECT id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih
	          FROM gorevler`
	args := []interface{}{}

	if durum != "" {
		sorgu += " WHERE durum = ?"
		args = append(args, durum)
	}

	sorgu += " ORDER BY olusturma_tarih DESC"

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
		)
		if err != nil {
			return nil, err
		}

		if projeID.Valid {
			gorev.ProjeID = projeID.String
		}

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

func (vy *VeriYonetici) GorevGuncelle(gorev *Gorev) error {
	sorgu := `UPDATE gorevler SET baslik = ?, aciklama = ?, durum = ?, oncelik = ?, 
	          proje_id = ?, guncelleme_tarih = ? WHERE id = ?`

	_, err := vy.db.Exec(sorgu,
		gorev.Baslik,
		gorev.Aciklama,
		gorev.Durum,
		gorev.Oncelik,
		gorev.ProjeID,
		gorev.GuncellemeTarih,
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
