package gorev

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

type VeriYonetici struct {
	db *sql.DB
}

func YeniVeriYonetici(dbYolu string, migrationsYolu string) (*VeriYonetici, error) {
	db, err := sql.Open("sqlite3", dbYolu)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.dbOpenFailed", map[string]interface{}{"Error": err}))
	}

	vy := &VeriYonetici{db: db}
	if err := vy.migrateDB(migrationsYolu); err != nil {
		return nil, fmt.Errorf(i18n.T("error.migrationFailed", map[string]interface{}{"Error": err}))
	}

	return vy, nil
}

// YeniVeriYoneticiWithEmbeddedMigrations creates a new VeriYonetici with embedded migrations
func YeniVeriYoneticiWithEmbeddedMigrations(dbYolu string, migrationsFS fs.FS) (*VeriYonetici, error) {
	db, err := sql.Open("sqlite3", dbYolu)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.dbOpenFailed", map[string]interface{}{"Error": err}))
	}

	vy := &VeriYonetici{db: db}
	if err := vy.migrateDBWithFS(migrationsFS); err != nil {
		return nil, fmt.Errorf(i18n.T("error.migrationFailed", map[string]interface{}{"Error": err}))
	}

	return vy, nil
}

func (vy *VeriYonetici) migrateDB(migrationsYolu string) error {
	log.Printf("DEBUG: migrateDB called with path: %s", migrationsYolu)

	// Create schema_migrations table if it doesn't exist
	_, err := vy.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			dirty INTEGER NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		log.Printf("ERROR: Failed to create schema_migrations table: %v", err)
		return fmt.Errorf(i18n.T("error.migrationTableFailed", map[string]interface{}{"Error": err}))
	}
	log.Printf("DEBUG: Schema migrations table ready")

	// Read migration files from directory
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsYolu, "*.up.sql"))
	if err != nil {
		log.Printf("ERROR: Failed to read migration files from %s: %v", migrationsYolu, err)
		return fmt.Errorf(i18n.T("error.migrationFileReadFailed", map[string]interface{}{"Error": err}))
	}
	log.Printf("DEBUG: Found %d migration files", len(migrationFiles))

	for _, migrationFile := range migrationFiles {
		// Extract version from filename (e.g., 000001_initial_schema.up.sql -> 1)
		filename := filepath.Base(migrationFile)
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			continue
		}
		versionStr := parts[0]
		version := 0
		if _, parseErr := fmt.Sscanf(versionStr, "%d", &version); parseErr != nil {
			log.Printf("WARNING: Could not parse version from %s: %v", filename, parseErr)
			continue
		}

		// Check if migration is already applied
		var exists int
		err = vy.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&exists)
		if err != nil {
			log.Printf("ERROR: Failed to check migration status for version %d: %v", version, err)
			return fmt.Errorf(i18n.T("error.migrationCheckFailed", map[string]interface{}{"Version": version, "Error": err}))
		}

		if exists > 0 {
			log.Printf("DEBUG: Migration %d already applied, skipping", version)
			continue
		}

		// Read and execute migration file
		migrationSQL, err := os.ReadFile(migrationFile)
		if err != nil {
			log.Printf("ERROR: Failed to read migration file %s: %v", migrationFile, err)
			return fmt.Errorf(i18n.T("error.migrationFileReadFailed", map[string]interface{}{"File": migrationFile, "Error": err}))
		}

		log.Printf("DEBUG: Applying migration %d from %s", version, filename)
		_, err = vy.db.Exec(string(migrationSQL))
		if err != nil {
			log.Printf("ERROR: Failed to execute migration %d: %v", version, err)
			return fmt.Errorf(i18n.T("error.migrationExecuteFailed", map[string]interface{}{"Version": version, "Error": err}))
		}

		// Mark migration as applied
		_, err = vy.db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
		if err != nil {
			log.Printf("ERROR: Failed to record migration %d: %v", version, err)
			return fmt.Errorf(i18n.T("error.migrationRecordFailed", map[string]interface{}{"Version": version, "Error": err}))
		}

		log.Printf("DEBUG: Migration %d applied successfully", version)
	}

	log.Println("SUCCESS: Database migrated successfully")

	// Varsayılan template'leri oluştur
	if err := vy.VarsayilanTemplateleriOlustur(); err != nil {
		log.Printf("WARNING: Failed to create default templates: %v", err)
		// Hata durumunda devam et, kritik değil
	}

	return nil
}

// migrateDBWithFS migrates database using embedded filesystem
func (vy *VeriYonetici) migrateDBWithFS(migrationsFS fs.FS) error {
	log.Printf("DEBUG: migrateDBWithFS called for platform: %s", runtime.GOOS)

	// Extract embedded migrations to temporary directory
	tempDir, err := os.MkdirTemp("", "gorev-migrations-*")
	if err != nil {
		log.Printf("ERROR: Failed to create temp dir: %v", err)
		return fmt.Errorf(i18n.T("error.migrationDirCreateFailed", map[string]interface{}{"Error": err}))
	}
	log.Printf("DEBUG: Created temp dir: %s", tempDir)
	defer os.RemoveAll(tempDir)

	// Copy all migration files from embedded FS to temp directory
	fileCount := 0
	err = fs.WalkDir(migrationsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("ERROR: WalkDir error for %s: %v", path, err)
			return err
		}
		if d.IsDir() {
			log.Printf("DEBUG: Skipping directory: %s", path)
			return nil
		}

		// Read file from embedded FS
		content, err := fs.ReadFile(migrationsFS, path)
		if err != nil {
			log.Printf("ERROR: Failed to read embedded file %s: %v", path, err)
			return err
		}

		// Write to temp directory
		destPath := filepath.Join(tempDir, path)
		log.Printf("DEBUG: Writing file %s to %s", path, destPath)
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			log.Printf("ERROR: Failed to write file %s: %v", destPath, err)
			return fmt.Errorf(i18n.T("error.migrationFileWriteFailed", map[string]interface{}{"Error": err}))
		}
		fileCount++
		return nil
	})
	if err != nil {
		log.Printf("ERROR: Failed to extract migrations: %v", err)
		return fmt.Errorf(i18n.T("error.migrationExtractFailed", map[string]interface{}{"Error": err}))
	}
	log.Printf("DEBUG: Extracted %d migration files to %s", fileCount, tempDir)

	// Now use regular file-based migration with the temp directory
	log.Printf("DEBUG: Using temp directory for migrations: %s", tempDir)
	return vy.migrateDB(tempDir)
}

// GorevListele retrieves tasks based on filters
func (vy *VeriYonetici) GorevListele(filters map[string]interface{}) ([]*Gorev, error) {
	// Convert filters to old format for compatibility
	durum := ""
	sirala := ""
	filtre := ""

	if v, ok := filters["durum"]; ok {
		if s, ok := v.(string); ok {
			durum = s
		}
	}
	if v, ok := filters["sirala"]; ok {
		if s, ok := v.(string); ok {
			sirala = s
		}
	}
	if v, ok := filters["filtre"]; ok {
		if s, ok := v.(string); ok {
			filtre = s
		}
	}

	return vy.GorevleriGetir(durum, sirala, filtre)
}

// GorevOlustur creates a new task
func (vy *VeriYonetici) GorevOlustur(params map[string]interface{}) (string, error) {
	gorev := &Gorev{
		ID:              uuid.New().String(),
		Durum:           constants.TaskStatusPending,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}

	if v, ok := params["baslik"]; ok {
		if s, ok := v.(string); ok {
			gorev.Baslik = s
		}
	}
	if v, ok := params["aciklama"]; ok {
		if s, ok := v.(string); ok {
			gorev.Aciklama = s
		}
	}
	if v, ok := params["oncelik"]; ok {
		if s, ok := v.(string); ok {
			gorev.Oncelik = s
		}
	}
	if v, ok := params["proje_id"]; ok {
		if s, ok := v.(string); ok {
			gorev.ProjeID = s
		}
	}
	if v, ok := params["parent_id"]; ok {
		if s, ok := v.(string); ok {
			gorev.ParentID = s
		}
	}

	if err := vy.GorevKaydet(gorev); err != nil {
		return "", err
	}

	return gorev.ID, nil
}

// GorevDetay retrieves detailed task information
func (vy *VeriYonetici) GorevDetay(taskID string) (*Gorev, error) {
	return vy.GorevGetir(taskID)
}

// GorevBagimlilikGetir retrieves task dependencies
func (vy *VeriYonetici) GorevBagimlilikGetir(taskID string) ([]*Gorev, error) {
	// Get all dependencies for the task
	baglantilari, err := vy.BaglantilariGetir(taskID)
	if err != nil {
		return nil, err
	}

	var bagimliGorevler []*Gorev
	for _, baglanti := range baglantilari {
		if baglanti.HedefID == taskID {
			// This task depends on the source task
			gorev, err := vy.GorevGetir(baglanti.KaynakID)
			if err == nil {
				bagimliGorevler = append(bagimliGorevler, gorev)
			}
		}
	}

	return bagimliGorevler, nil
}

// AltGorevOlustur creates a subtask under a parent task
func (vy *VeriYonetici) AltGorevOlustur(parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	var sonTarih *time.Time
	if sonTarihStr != "" {
		t, err := time.Parse("2006-01-02", sonTarihStr)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.invalidDateFormat", map[string]interface{}{"Error": err}))
		}
		sonTarih = &t
	}

	gorev := &Gorev{
		ID:              uuid.New().String(),
		Baslik:          baslik,
		Aciklama:        aciklama,
		Oncelik:         oncelik,
		Durum:           constants.TaskStatusPending,
		ParentID:        parentID,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
		SonTarih:        sonTarih,
	}

	if err := vy.GorevKaydet(gorev); err != nil {
		return nil, fmt.Errorf(i18n.TSaveFailed("task", err))
	}

	if len(etiketIsimleri) > 0 {
		etiketler, err := vy.EtiketleriGetirVeyaOlustur(etiketIsimleri)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsProcessFailed", map[string]interface{}{"Error": err}))
		}
		if err := vy.GorevEtiketleriniAyarla(gorev.ID, etiketler); err != nil {
			return nil, fmt.Errorf(i18n.TSetFailed("task_tags", err))
		}
		gorev.Etiketler = etiketler
	}

	return gorev, nil
}

func (vy *VeriYonetici) Kapat() error {
	return vy.db.Close()
}

// GetDB returns the underlying database connection for advanced operations
func (vy *VeriYonetici) GetDB() (*sql.DB, error) {
	if vy.db == nil {
		return nil, fmt.Errorf(i18n.T("error.dbConnectionClosed"))
	}
	return vy.db, nil
}

// ProjeOlustur creates a new project with minimal data for testing
func (vy *VeriYonetici) ProjeOlustur(isim, aciklama string, etiketler ...string) (*Proje, error) {
	proje := &Proje{
		ID:             fmt.Sprintf("proj_%d", time.Now().UnixNano()),
		Isim:           isim,
		Tanim:          aciklama,
		OlusturmaTarih: time.Now(),
	}

	err := vy.ProjeKaydet(proje)
	if err != nil {
		return nil, err
	}

	return proje, nil
}

// GorevOlusturBasit creates a task with individual parameters for testing
func (vy *VeriYonetici) GorevOlusturBasit(baslik, aciklama, projeID, oncelik, sonTarih, parentID, etiketler string) (*Gorev, error) {
	gorev := &Gorev{
		ID:              fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Baslik:          baslik,
		Aciklama:        aciklama,
		Durum:           "beklemede",
		Oncelik:         oncelik,
		ProjeID:         projeID,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}

	if sonTarih != "" {
		if tarih, err := time.Parse("2006-01-02", sonTarih); err == nil {
			gorev.SonTarih = &tarih
		}
	}

	if parentID != "" {
		gorev.ParentID = parentID
	}

	err := vy.GorevKaydet(gorev)
	if err != nil {
		return nil, err
	}

	return gorev, nil
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
	defer func() {
		if err := rows.Close(); err != nil {
			// best-effort close
		}
	}()

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
	defer func() {
		if err := rows.Close(); err != nil {
			// best-effort close
		}
	}()

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
		return nil, fmt.Errorf(i18n.T("error.transactionFailed", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = tx.Rollback() }() // Hata durumunda geri al

	stmtSelect, err := tx.Prepare("SELECT id, isim FROM etiketler WHERE isim = ?")
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.selectPrepFailed", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = stmtSelect.Close() }()

	stmtInsert, err := tx.Prepare("INSERT INTO etiketler (id, isim) VALUES (?, ?)")
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.insertPrepFailed", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = stmtInsert.Close() }()

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
		return fmt.Errorf(i18n.T("error.transactionFailed", map[string]interface{}{"Error": err}))
	}
	defer func() {
		// Silently ignore rollback errors as transaction may already be committed
		_ = tx.Rollback()
	}()

	// Mevcut bağlantıları sil
	if _, err := tx.Exec("DELETE FROM gorev_etiketleri WHERE gorev_id = ?", gorevID); err != nil {
		return fmt.Errorf(i18n.T("error.currentTagsRemoveFailed", map[string]interface{}{"Error": err}))
	}

	// Yeni bağlantıları ekle
	stmt, err := tx.Prepare("INSERT INTO gorev_etiketleri (gorev_id, etiket_id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf(i18n.T("error.insertPrepFailed", map[string]interface{}{"Error": err}))
	}
	defer func() {
		_ = stmt.Close()
	}()

	for _, etiket := range etiketler {
		if _, err := stmt.Exec(gorevID, etiket.ID); err != nil {
			return fmt.Errorf(i18n.T("error.taskTagAddFailed", map[string]interface{}{"Tag": etiket.Isim, "Error": err}))
		}
	}

	return tx.Commit()
}

func (vy *VeriYonetici) GorevGuncelle(taskID string, params interface{}) error {
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid params type, expected map[string]interface{}")
	}

	if len(paramsMap) == 0 {
		return nil // No updates to perform
	}

	// Build dynamic UPDATE query
	var setParts []string
	var args []interface{}

	for key, value := range paramsMap {
		setParts = append(setParts, key+" = ?")
		args = append(args, value)
	}

	sorgu := fmt.Sprintf("UPDATE gorevler SET %s WHERE id = ?", strings.Join(setParts, ", "))
	args = append(args, taskID)

	_, err := vy.db.Exec(sorgu, args...)
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
	defer func() {
		_ = rows.Close()
	}()

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
		return fmt.Errorf(i18n.TEntityNotFound("task", errors.New("not found")))
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
	defer func() {
		_ = rows.Close()
	}()

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
	defer func() {
		_ = rows.Close()
	}()

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
	defer func() {
		_ = rows.Close()
	}()

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
	defer func() {
		_ = rows.Close()
	}()

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
	defer func() {
		_ = rows.Close()
	}()

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
	} else if gorev.Durum == constants.TaskStatusCompleted {
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
			return fmt.Errorf(i18n.T("error.circularDependency"))
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
		SELECT hedef_id, COUNT(*) as bagli_sayi
		FROM baglantilar 
		WHERE hedef_id IN (%s)
		GROUP BY hedef_id
	`, strings.Join(placeholders, ","))

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

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
		SELECT b.hedef_id, COUNT(*) as tamamlanmamis_sayi
		FROM baglantilar b
		INNER JOIN gorevler g ON b.kaynak_id = g.id
		WHERE b.hedef_id IN (%s) AND g.durum != 'tamamlandi'
		GROUP BY b.hedef_id
	`, strings.Join(placeholders, ","))

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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

// BulkBuGoreveBagimliSayilariGetir tüm görevlere bağımlı olan görev sayılarını hesaplar
func (vy *VeriYonetici) BulkBuGoreveBagimliSayilariGetir(gorevIDs []string) (map[string]int, error) {
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

	// Bu göreve bağımlı olan görevleri say (bu görev kaynak olduğu bağlantılar)
	sorgu := fmt.Sprintf(`
		SELECT b.kaynak_id, COUNT(*) as bagimli_sayi
		FROM baglantilar b
		WHERE b.kaynak_id IN (%s)
		GROUP BY b.kaynak_id
	`, strings.Join(placeholders, ","))

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

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

// AI Context Management Methods

// AIContextGetir retrieves the current AI context from the database
func (vy *VeriYonetici) AIContextGetir() (*AIContext, error) {
	var activeTaskID sql.NullString
	var recentTasksJSON, sessionDataJSON string
	var lastUpdated time.Time

	sorgu := `SELECT active_task_id, recent_tasks, session_data, last_updated FROM ai_context WHERE id = 1`
	err := vy.db.QueryRow(sorgu).Scan(&activeTaskID, &recentTasksJSON, &sessionDataJSON, &lastUpdated)

	if err != nil {
		if err == sql.ErrNoRows {
			// Initialize default context if not exists
			defaultContext := &AIContext{
				RecentTasks: []string{},
				SessionData: make(map[string]interface{}),
				LastUpdated: time.Now(),
			}
			// Save the default context
			if err := vy.AIContextKaydet(defaultContext); err != nil {
				return nil, fmt.Errorf(i18n.T("error.contextInitializationFailed", map[string]interface{}{"Error": err}))
			}
			return defaultContext, nil
		}
		return nil, fmt.Errorf(i18n.T("error.contextRetrievalFailed", map[string]interface{}{"Error": err}))
	}

	// Parse JSON fields
	var recentTasks []string
	if err := json.Unmarshal([]byte(recentTasksJSON), &recentTasks); err != nil {
		recentTasks = []string{}
	}

	var sessionData map[string]interface{}
	if err := json.Unmarshal([]byte(sessionDataJSON), &sessionData); err != nil {
		sessionData = make(map[string]interface{})
	}

	context := &AIContext{
		RecentTasks: recentTasks,
		SessionData: sessionData,
		LastUpdated: lastUpdated,
	}

	if activeTaskID.Valid {
		context.ActiveTaskID = activeTaskID.String
	}

	return context, nil
}

// AIContextKaydet saves the AI context to the database
func (vy *VeriYonetici) AIContextKaydet(context *AIContext) error {
	recentTasksJSON, err := json.Marshal(context.RecentTasks)
	if err != nil {
		return fmt.Errorf(i18n.T("error.jsonMarshalFailed", map[string]interface{}{"Field": "recent_tasks", "Error": err}))
	}

	sessionDataJSON, err := json.Marshal(context.SessionData)
	if err != nil {
		return fmt.Errorf(i18n.T("error.jsonMarshalFailed", map[string]interface{}{"Field": "session_data", "Error": err}))
	}

	var activeTaskID *string
	if context.ActiveTaskID != "" {
		activeTaskID = &context.ActiveTaskID
	}

	sorgu := `
		INSERT OR REPLACE INTO ai_context (id, active_task_id, recent_tasks, session_data, last_updated)
		VALUES (1, ?, ?, ?, ?)`

	_, err = vy.db.Exec(sorgu, activeTaskID, string(recentTasksJSON), string(sessionDataJSON), time.Now())
	if err != nil {
		return fmt.Errorf(i18n.T("error.contextSaveFailed", map[string]interface{}{"Error": err}))
	}

	return nil
}

// AIInteractionKaydet records an AI interaction
func (vy *VeriYonetici) AIInteractionKaydet(interaction *AIInteraction) error {
	contextJSON := ""
	if interaction.Context != "" {
		contextJSON = interaction.Context
	}

	sorgu := `
		INSERT INTO ai_interactions (gorev_id, action_type, context, timestamp)
		VALUES (?, ?, ?, ?)`

	_, err := vy.db.Exec(sorgu, interaction.GorevID, interaction.ActionType, contextJSON, time.Now())
	if err != nil {
		return fmt.Errorf(i18n.T("error.interactionSaveFailed", map[string]interface{}{"Error": err}))
	}

	return nil
}

// AIInteractionlariGetir retrieves recent AI interactions
func (vy *VeriYonetici) AIInteractionlariGetir(limit int) ([]*AIInteraction, error) {
	if limit <= 0 {
		limit = 10
	}

	sorgu := `
		SELECT id, gorev_id, action_type, context, timestamp
		FROM ai_interactions
		ORDER BY timestamp DESC
		LIMIT ?`

	rows, err := vy.db.Query(sorgu, limit)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.interactionQueryFailed", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = rows.Close() }()

	var interactions []*AIInteraction
	for rows.Next() {
		var interaction AIInteraction
		var contextStr sql.NullString
		err := rows.Scan(&interaction.ID, &interaction.GorevID, &interaction.ActionType, &contextStr, &interaction.Timestamp)
		if err != nil {
			return nil, err
		}

		if contextStr.Valid {
			interaction.Context = contextStr.String
		}

		interactions = append(interactions, &interaction)
	}

	return interactions, nil
}

// AITodayInteractionlariGetir retrieves today's AI interactions
func (vy *VeriYonetici) AITodayInteractionlariGetir() ([]*AIInteraction, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	sorgu := `
		SELECT id, gorev_id, action_type, context, timestamp
		FROM ai_interactions
		WHERE timestamp >= ? AND timestamp < ?
		ORDER BY timestamp DESC`

	rows, err := vy.db.Query(sorgu, today, tomorrow)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.todayInteractionQueryFailed", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = rows.Close() }()

	var interactions []*AIInteraction
	for rows.Next() {
		var interaction AIInteraction
		var contextStr sql.NullString
		err := rows.Scan(&interaction.ID, &interaction.GorevID, &interaction.ActionType, &contextStr, &interaction.Timestamp)
		if err != nil {
			return nil, err
		}

		if contextStr.Valid {
			interaction.Context = contextStr.String
		}

		interactions = append(interactions, &interaction)
	}

	return interactions, nil
}

// AILastInteractionGuncelle updates the last AI interaction timestamp for a task
func (vy *VeriYonetici) AILastInteractionGuncelle(taskID string, timestamp time.Time) error {
	sorgu := `UPDATE gorevler SET last_ai_interaction = ? WHERE id = ?`

	_, err := vy.db.Exec(sorgu, timestamp, taskID)
	if err != nil {
		return fmt.Errorf(i18n.T("error.lastInteractionUpdateFailed", map[string]interface{}{"TaskID": taskID, "Error": err}))
	}

	return nil
}

// AIEtkilemasimKaydet saves an AI interaction record
func (vy *VeriYonetici) AIEtkilemasimKaydet(taskID string, interactionType, data, sessionID string) error {
	// For now, we'll use the existing AIInteractionKaydet method
	// This is a simplified implementation that matches the interface
	interaction := &AIInteraction{
		GorevID:    taskID,
		ActionType: interactionType,
		Context:    data,
		Timestamp:  time.Now(),
	}
	return vy.AIInteractionKaydet(interaction)
}

// GorevSonAIEtkilesiminiGuncelle updates the last AI interaction timestamp for a task
func (vy *VeriYonetici) GorevSonAIEtkilesiminiGuncelle(taskID string, timestamp time.Time) error {
	return vy.AILastInteractionGuncelle(taskID, timestamp)
}

// GorevDosyaYoluEkle adds a file path to a task
func (vy *VeriYonetici) GorevDosyaYoluEkle(taskID string, path string) error {
	// This would need proper implementation with a file_paths table
	// For now, return nil as a placeholder
	return nil
}

// GorevDosyaYoluSil removes a file path from a task
func (vy *VeriYonetici) GorevDosyaYoluSil(taskID string, path string) error {
	// This would need proper implementation with a file_paths table
	// For now, return nil as a placeholder
	return nil
}

// GorevDosyaYollariGetir gets all file paths for a task
func (vy *VeriYonetici) GorevDosyaYollariGetir(taskID string) ([]string, error) {
	// This would need proper implementation with a file_paths table
	// For now, return empty slice as a placeholder
	return []string{}, nil
}

// DosyaYoluGorevleriGetir gets all tasks associated with a file path
func (vy *VeriYonetici) DosyaYoluGorevleriGetir(path string) ([]string, error) {
	// This would need proper implementation with a file_paths table
	// For now, return empty slice as a placeholder
	return []string{}, nil
}
