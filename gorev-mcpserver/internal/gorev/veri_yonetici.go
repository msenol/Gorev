package gorev

import (
	"context"
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
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
	_ "modernc.org/sqlite"
)

// EventEmitter interface for emitting database change events
// Import from websocket package to avoid circular dependency
type EventEmitter interface {
	EmitTaskCreated(workspaceID, taskID string, data map[string]interface{})
	EmitTaskUpdated(workspaceID, taskID string, data map[string]interface{})
	EmitTaskDeleted(workspaceID, taskID string)
	EmitProjectCreated(workspaceID, projectID string, data map[string]interface{})
	EmitProjectUpdated(workspaceID, projectID string, data map[string]interface{})
	EmitProjectDeleted(workspaceID, projectID string)
	EmitTemplateChanged(workspaceID string)
	EmitWorkspaceSync(workspaceID string)
}

type VeriYonetici struct {
	db           *sql.DB
	eventEmitter EventEmitter
	workspaceID  string // Workspace ID for event emission
}

// configureSQLiteForConcurrency configures SQLite for better concurrent access
func configureSQLiteForConcurrency(db *sql.DB) error {
	// Enable WAL mode for better concurrent access
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return fmt.Errorf(i18n.T("error.walModeFailed", map[string]interface{}{"Error": err}))
	}

	// Set busy timeout to avoid immediate failures (10 seconds)
	if _, err := db.Exec("PRAGMA busy_timeout=10000"); err != nil {
		return fmt.Errorf(i18n.T("error.busyTimeoutFailed", map[string]interface{}{"Error": err}))
	}

	// Set connection pooling for better concurrency
	// Lower max open connections to reduce contention
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	return nil
}

// retryOnBusy retries a database operation if it fails with SQLITE_BUSY error
func retryOnBusy(operation func() error, maxRetries int) error {
	var err error
	for i := 0; i <= maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}

		// Check if error is SQLITE_BUSY
		if strings.Contains(err.Error(), "database is locked") ||
			strings.Contains(err.Error(), "SQLITE_BUSY") {
			// Exponential backoff with cap: 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1000ms (capped)
			backoff := time.Duration(10<<uint(i)) * time.Millisecond
			if backoff > time.Second {
				backoff = time.Second // Cap at 1 second
			}
			time.Sleep(backoff)
			continue
		}

		// Not a busy error, return immediately
		return err
	}
	return fmt.Errorf(i18n.T("error.operationRetryFailed", map[string]interface{}{"Retries": maxRetries, "Error": err}))
}

func YeniVeriYonetici(dbYolu string, migrationsYolu string) (*VeriYonetici, error) {
	return YeniVeriYoneticiWithEventEmitter(dbYolu, migrationsYolu, nil, "")
}

// YeniVeriYoneticiWithEventEmitter creates a new VeriYonetici with optional event emitter
func YeniVeriYoneticiWithEventEmitter(dbYolu string, migrationsYolu string, eventEmitter EventEmitter, workspaceID string) (*VeriYonetici, error) {
	db, err := sql.Open("sqlite", dbYolu)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.dbOpenFailed", map[string]interface{}{"Error": err}))
	}

	// Configure SQLite for better concurrent access
	if err := configureSQLiteForConcurrency(db); err != nil {
		db.Close()
		return nil, fmt.Errorf(i18n.T("error.dbConfigureFailed", map[string]interface{}{"Error": err}))
	}

	vy := &VeriYonetici{
		db:           db,
		eventEmitter: eventEmitter,
		workspaceID:  workspaceID,
	}
	if err := vy.migrateDB(migrationsYolu); err != nil {
		return nil, fmt.Errorf(i18n.T("error.migrationFailed", map[string]interface{}{"Error": err}))
	}

	return vy, nil
}

// YeniVeriYoneticiWithEmbeddedMigrations creates a new VeriYonetici with embedded migrations
func YeniVeriYoneticiWithEmbeddedMigrations(dbYolu string, migrationsFS fs.FS) (*VeriYonetici, error) {
	return YeniVeriYoneticiWithEmbeddedMigrationsAndEventEmitter(dbYolu, migrationsFS, nil, "")
}

// YeniVeriYoneticiWithEmbeddedMigrationsAndEventEmitter creates a new VeriYonetici with embedded migrations and optional event emitter
func YeniVeriYoneticiWithEmbeddedMigrationsAndEventEmitter(dbYolu string, migrationsFS fs.FS, eventEmitter EventEmitter, workspaceID string) (*VeriYonetici, error) {
	db, err := sql.Open("sqlite", dbYolu)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.dbOpenFailed", map[string]interface{}{"Error": err}))
	}

	// Configure SQLite for better concurrent access
	if err := configureSQLiteForConcurrency(db); err != nil {
		db.Close()
		return nil, fmt.Errorf(i18n.T("error.dbConfigureFailed", map[string]interface{}{"Error": err}))
	}

	vy := &VeriYonetici{
		db:           db,
		eventEmitter: eventEmitter,
		workspaceID:  workspaceID,
	}
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

	// Check if database needs migration state repair
	if err := vy.repairMigrationStateIfNeeded(); err != nil {
		log.Printf("WARNING: Migration state repair failed: %v", err)
		// Continue with normal migration process
	}

	// Parse migration path and remove file:// prefix if present
	migrationsPath := migrationsYolu
	if strings.HasPrefix(migrationsPath, "file://") {
		migrationsPath = strings.TrimPrefix(migrationsPath, "file://")
		// Make relative paths absolute
		if !filepath.IsAbs(migrationsPath) {
			if abs, err := filepath.Abs(migrationsPath); err == nil {
				migrationsPath = abs
			}
		}
	}
	log.Printf("DEBUG: Using filesystem path: %s", migrationsPath)

	// Read migration files from directory
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsPath, "*.up.sql"))
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
	if err := vy.VarsayilanTemplateleriOlustur(context.Background()); err != nil {
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
func (vy *VeriYonetici) GorevListele(ctx context.Context, filters map[string]interface{}) ([]*Gorev, error) {
	// Convert filters to old format for compatibility
	status := ""
	sirala := ""
	filtre := ""

	if v, ok := filters["status"]; ok {
		if s, ok := v.(string); ok {
			status = s
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

	return vy.GorevleriGetir(ctx, status, sirala, filtre)
}

// GorevOlustur creates a new task
func (vy *VeriYonetici) GorevOlustur(ctx context.Context, params map[string]interface{}) (string, error) {
	gorev := &Gorev{
		ID:        uuid.New().String(),
		Status:    constants.TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if v, ok := params["title"]; ok {
		if s, ok := v.(string); ok {
			gorev.Title = s
		}
	}
	if v, ok := params["description"]; ok {
		if s, ok := v.(string); ok {
			gorev.Description = s
		}
	}
	if v, ok := params["priority"]; ok {
		if s, ok := v.(string); ok {
			gorev.Priority = s
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

	if err := vy.GorevKaydet(ctx, gorev); err != nil {
		return "", err
	}

	return gorev.ID, nil
}

// GorevDetay retrieves detailed task information
func (vy *VeriYonetici) GorevDetay(ctx context.Context, taskID string) (*Gorev, error) {
	return vy.GorevGetir(ctx, taskID)
}

// GorevBagimlilikGetir retrieves task dependencies
func (vy *VeriYonetici) GorevBagimlilikGetir(ctx context.Context, taskID string) ([]*Gorev, error) {
	// Get all dependencies for the task
	baglantilari, err := vy.BaglantilariGetir(ctx, taskID)
	if err != nil {
		return nil, err
	}

	var bagimliGorevler []*Gorev
	for _, baglanti := range baglantilari {
		if baglanti.TargetID == taskID {
			// This task depends on the source task
			gorev, err := vy.GorevGetir(ctx, baglanti.SourceID)
			if err == nil {
				bagimliGorevler = append(bagimliGorevler, gorev)
			}
		}
	}

	return bagimliGorevler, nil
}

// AltGorevOlustur creates a subtask under a parent task
func (vy *VeriYonetici) AltGorevOlustur(ctx context.Context, parentID, title, description, priority, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	var sonTarih *time.Time
	if sonTarihStr != "" {
		t, err := time.Parse("2006-01-02", sonTarihStr)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.invalidDateFormat", map[string]interface{}{"Error": err}))
		}
		sonTarih = &t
	}

	gorev := &Gorev{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      constants.TaskStatusPending,
		ParentID:    parentID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     sonTarih,
	}

	if err := vy.GorevKaydet(ctx, gorev); err != nil {
		return nil, fmt.Errorf(i18n.TSaveFailed(i18n.FromContext(ctx), "task", err))
	}

	if len(etiketIsimleri) > 0 {
		etiketler, err := vy.EtiketleriGetirVeyaOlustur(ctx, etiketIsimleri)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsProcessFailed", map[string]interface{}{"Error": err}))
		}
		if err := vy.GorevEtiketleriniAyarla(ctx, gorev.ID, etiketler); err != nil {
			return nil, fmt.Errorf(i18n.TSetFailed(i18n.FromContext(ctx), "task_tags", err))
		}
		gorev.Tags = etiketler
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
func (vy *VeriYonetici) ProjeOlustur(ctx context.Context, name, description string, etiketler ...string) (*Proje, error) {
	proje := &Proje{
		ID:         fmt.Sprintf("proj_%d", time.Now().UnixNano()),
		Name:       name,
		Definition: description,
		CreatedAt:  time.Now(),
	}

	err := vy.ProjeKaydet(ctx, proje)
	if err != nil {
		return nil, err
	}

	return proje, nil
}

// GorevOlusturBasit creates a task with individual parameters for testing
func (vy *VeriYonetici) GorevOlusturBasit(ctx context.Context, title, description, projeID, priority, sonTarih, parentID, etiketler string) (*Gorev, error) {
	gorev := &Gorev{
		ID:          fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Title:       title,
		Description: description,
		Status:      "beklemede",
		Priority:    priority,
		ProjeID:     projeID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if sonTarih != "" {
		if tarih, err := time.Parse("2006-01-02", sonTarih); err == nil {
			gorev.DueDate = &tarih
		}
	}

	if parentID != "" {
		gorev.ParentID = parentID
	}

	err := vy.GorevKaydet(ctx, gorev)
	if err != nil {
		return nil, err
	}

	return gorev, nil
}

func (vy *VeriYonetici) GorevKaydet(ctx context.Context, gorev *Gorev) error {
	sorgu := `INSERT INTO gorevler (id, title, description, status, priority, project_id, parent_id, created_at, updated_at, due_date)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Use retry logic for better concurrent write handling
	err := retryOnBusy(func() error {
		_, err := vy.db.Exec(sorgu,
			gorev.ID,
			gorev.Title,
			gorev.Description,
			gorev.Status,
			gorev.Priority,
			sql.NullString{String: gorev.ProjeID, Valid: gorev.ProjeID != ""},
			sql.NullString{String: gorev.ParentID, Valid: gorev.ParentID != ""},
			gorev.CreatedAt,
			gorev.UpdatedAt,
			gorev.DueDate,
		)
		return err
	}, 10) // Retry up to 10 times with exponential backoff (capped at 1s)

	// Emit task created event if operation succeeded
	if err == nil && vy.eventEmitter != nil {
		vy.eventEmitter.EmitTaskCreated(vy.workspaceID, gorev.ID, map[string]interface{}{
			"title":    gorev.Title,
			"status":   gorev.Status,
			"priority": gorev.Priority,
		})
	}

	return err
}

func (vy *VeriYonetici) GorevGetir(ctx context.Context, id string) (*Gorev, error) {
	sorgu := `SELECT id, title, description, status, priority, project_id, parent_id, created_at, updated_at, due_date
	          FROM gorevler WHERE id = ?`

	gorev := &Gorev{}
	var projeID, parentID sql.NullString

	err := vy.db.QueryRow(sorgu, id).Scan(
		&gorev.ID,
		&gorev.Title,
		&gorev.Description,
		&gorev.Status,
		&gorev.Priority,
		&projeID,
		&parentID,
		&gorev.CreatedAt,
		&gorev.UpdatedAt,
		&gorev.DueDate,
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
		gorev.Tags = []*Etiket{}
	} else {
		gorev.Tags = etiketler
	}

	// Proje adını getir (Web UI ve VS Code için)
	if gorev.ProjeID != "" {
		proje, err := vy.ProjeGetir(ctx, gorev.ProjeID)
		if err == nil && proje != nil {
			gorev.ProjeName = proje.Name
		}
	}

	return gorev, nil
}

func (vy *VeriYonetici) GorevleriGetir(ctx context.Context, status, sirala, filtre string) ([]*Gorev, error) {
	sorgu := `SELECT id, title, description, status, priority, project_id, parent_id, created_at, updated_at, due_date
	          FROM gorevler`
	args := []interface{}{}
	whereClauses := []string{}

	if status != "" {
		whereClauses = append(whereClauses, "status = ?")
		args = append(args, status)
	}

	if filtre == "acil" {
		whereClauses = append(whereClauses, "due_date IS NOT NULL AND due_date >= date('now') AND due_date < date('now', '+7 days')")
	} else if filtre == "gecmis" {
		whereClauses = append(whereClauses, "due_date IS NOT NULL AND due_date < date('now')")
	}

	if len(whereClauses) > 0 {
		sorgu += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	switch sirala {
	case "son_tarih_asc":
		sorgu += " ORDER BY due_date ASC"
	case "son_tarih_desc":
		sorgu += " ORDER BY due_date DESC"
	default:
		sorgu += " ORDER BY created_at DESC"
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
			&gorev.Title,
			&gorev.Description,
			&gorev.Status,
			&gorev.Priority,
			&projeID,
			&parentID,
			&gorev.CreatedAt,
			&gorev.UpdatedAt,
			&gorev.DueDate,
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
		gorev.Tags = etiketler

		// Proje adını getir (Web UI ve VS Code için)
		if gorev.ProjeID != "" {
			proje, err := vy.ProjeGetir(ctx, gorev.ProjeID)
			if err == nil && proje != nil {
				gorev.ProjeName = proje.Name
			}
		}

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

func (vy *VeriYonetici) gorevEtiketleriniGetir(gorevID string) ([]*Etiket, error) {
	sorgu := `SELECT e.id, e.name FROM etiketler e
	          JOIN gorev_etiketleri ge ON e.id = ge.tag_id
	          WHERE ge.task_id = ?`
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
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, err
		}
		etiketler = append(etiketler, e)
	}
	return etiketler, nil
}

func (vy *VeriYonetici) EtiketleriGetirVeyaOlustur(ctx context.Context, isimler []string) ([]*Etiket, error) {
	etiketler := make([]*Etiket, 0, len(isimler))
	tx, err := vy.db.Begin()
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.transactionFailed", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = tx.Rollback() }() // Hata durumunda geri al

	stmtSelect, err := tx.Prepare("SELECT id, name FROM etiketler WHERE name = ?")
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.selectPrepFailed", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = stmtSelect.Close() }()

	stmtInsert, err := tx.Prepare("INSERT INTO etiketler (id, name) VALUES (?, ?)")
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.insertPrepFailed", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = stmtInsert.Close() }()

	for _, name := range isimler {
		if strings.TrimSpace(name) == "" {
			continue
		}
		etiket := &Etiket{Name: strings.TrimSpace(name)}
		err := stmtSelect.QueryRow(etiket.Name).Scan(&etiket.ID, &etiket.Name)
		if err == sql.ErrNoRows {
			// Etiket yok, oluştur
			etiket.ID = uuid.New().String()
			if _, err := stmtInsert.Exec(etiket.ID, etiket.Name); err != nil {
				return nil, fmt.Errorf(i18n.T("error.tagCreateFailed", map[string]interface{}{"Tag": etiket.Name, "Error": err}))
			}
		} else if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagQueryFailed", map[string]interface{}{"Tag": etiket.Name, "Error": err}))
		}
		etiketler = append(etiketler, etiket)
	}

	return etiketler, tx.Commit()
}

func (vy *VeriYonetici) GorevEtiketleriniAyarla(ctx context.Context, gorevID string, etiketler []*Etiket) error {
	// Use retry logic for the entire transaction
	return retryOnBusy(func() error {
		tx, err := vy.db.Begin()
		if err != nil {
			return fmt.Errorf(i18n.T("error.transactionFailed", map[string]interface{}{"Error": err}))
		}
		defer func() {
			// Silently ignore rollback errors as transaction may already be committed
			_ = tx.Rollback()
		}()

		// Mevcut bağlantıları sil
		if _, err := tx.Exec("DELETE FROM gorev_etiketleri WHERE task_id = ?", gorevID); err != nil {
			return fmt.Errorf(i18n.T("error.currentTagsRemoveFailed", map[string]interface{}{"Error": err}))
		}

		// Yeni bağlantıları ekle
		stmt, err := tx.Prepare("INSERT INTO gorev_etiketleri (task_id, tag_id) VALUES (?, ?)")
		if err != nil {
			return fmt.Errorf(i18n.T("error.insertPrepFailed", map[string]interface{}{"Error": err}))
		}
		defer func() {
			_ = stmt.Close()
		}()

		for _, etiket := range etiketler {
			if _, err := stmt.Exec(gorevID, etiket.ID); err != nil {
				return fmt.Errorf(i18n.T("error.taskTagAddFailed", map[string]interface{}{"Tag": etiket.Name, "Error": err}))
			}
		}

		return tx.Commit()
	}, 10) // Retry up to 10 times with exponential backoff (capped at 1s)
}

func (vy *VeriYonetici) GorevGuncelle(ctx context.Context, taskID string, params interface{}) error {
	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf(i18n.T("error.invalidParamsType"))
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

	// Emit task updated event if operation succeeded
	if err == nil && vy.eventEmitter != nil {
		vy.eventEmitter.EmitTaskUpdated(vy.workspaceID, taskID, paramsMap)
	}

	return err
}

func (vy *VeriYonetici) ProjeKaydet(ctx context.Context, proje *Proje) error {
	sorgu := `INSERT INTO projeler (id, name, definition, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?)`

	_, err := vy.db.Exec(sorgu,
		proje.ID,
		proje.Name,
		proje.Definition,
		proje.CreatedAt,
		proje.UpdatedAt,
	)

	return err
}

func (vy *VeriYonetici) ProjeGetir(ctx context.Context, id string) (*Proje, error) {
	sorgu := `SELECT id, name, definition, created_at, updated_at
	          FROM projeler WHERE id = ?`

	proje := &Proje{}
	err := vy.db.QueryRow(sorgu, id).Scan(
		&proje.ID,
		&proje.Name,
		&proje.Definition,
		&proje.CreatedAt,
		&proje.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return proje, nil
}

func (vy *VeriYonetici) ProjeleriGetir(ctx context.Context) ([]*Proje, error) {
	sorgu := `SELECT p.id, p.name, p.definition, p.created_at, p.updated_at,
	          COUNT(g.id) as gorev_sayisi
	          FROM projeler p
	          LEFT JOIN gorevler g ON p.id = g.project_id
	          GROUP BY p.id, p.name, p.definition, p.created_at, p.updated_at
	          ORDER BY p.created_at DESC`

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
			&proje.Name,
			&proje.Definition,
			&proje.CreatedAt,
			&proje.UpdatedAt,
			&proje.TaskCount,
		)
		if err != nil {
			return nil, err
		}
		projeler = append(projeler, proje)
	}

	return projeler, nil
}

func (vy *VeriYonetici) GorevSil(ctx context.Context, id string) error {
	// Start a transaction to ensure atomic deletion
	tx, err := vy.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, delete all dependencies where this task is involved
	// This prevents FK constraint violations
	deleteDeps := `DELETE FROM baglantilar WHERE source_id = ? OR target_id = ?`
	_, err = tx.Exec(deleteDeps, id, id)
	if err != nil {
		return fmt.Errorf("failed to delete task dependencies: %w", err)
	}

	// Now delete the task itself
	sorgu := `DELETE FROM gorevler WHERE id = ?`
	result, err := tx.Exec(sorgu, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf(i18n.TEntityNotFound(i18n.FromContext(ctx), "task", errors.New("not found")))
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// Emit task deleted event if operation succeeded
	if vy.eventEmitter != nil {
		vy.eventEmitter.EmitTaskDeleted(vy.workspaceID, id)
	}

	return nil
}

func (vy *VeriYonetici) ProjeGorevleriGetir(ctx context.Context, projeID string) ([]*Gorev, error) {
	var sorgu string
	var rows *sql.Rows
	var err error

	// Handle empty projeID as NULL search
	if projeID == "" {
		sorgu = `SELECT id, title, description, status, priority, project_id, created_at, updated_at
		          FROM gorevler WHERE project_id IS NULL ORDER BY created_at DESC`
		rows, err = vy.db.Query(sorgu)
	} else {
		sorgu = `SELECT id, title, description, status, priority, project_id, created_at, updated_at
		          FROM gorevler WHERE project_id = ? ORDER BY created_at DESC`
		rows, err = vy.db.Query(sorgu, projeID)
	}
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
			&gorev.Title,
			&gorev.Description,
			&gorev.Status,
			&gorev.Priority,
			&pID,
			&gorev.CreatedAt,
			&gorev.UpdatedAt,
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

func (vy *VeriYonetici) BaglantiEkle(ctx context.Context, baglanti *Baglanti) error {
	sorgu := `INSERT INTO baglantilar (id, source_id, target_id, connection_type) VALUES (?, ?, ?, ?)`
	_, err := vy.db.Exec(sorgu, baglanti.ID, baglanti.SourceID, baglanti.TargetID, baglanti.ConnectionType)
	return err
}

// BaglantiSil removes a dependency relationship between two tasks
func (vy *VeriYonetici) BaglantiSil(ctx context.Context, kaynakID, hedefID string) error {
	sorgu := `DELETE FROM baglantilar WHERE source_id = ? AND target_id = ?`
	result, err := vy.db.Exec(sorgu, kaynakID, hedefID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf(i18n.T("error.dependencyNotFound", map[string]interface{}{"Source": kaynakID, "Target": hedefID}))
	}

	return nil
}

func (vy *VeriYonetici) BaglantilariGetir(ctx context.Context, gorevID string) ([]*Baglanti, error) {
	sorgu := `SELECT id, source_id, target_id, connection_type FROM baglantilar WHERE source_id = ? OR target_id = ?`
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
		if err := rows.Scan(&b.ID, &b.SourceID, &b.TargetID, &b.ConnectionType); err != nil {
			return nil, err
		}
		baglantilar = append(baglantilar, b)
	}
	return baglantilar, nil
}

// AltGorevleriGetir belirtilen görevin doğrudan alt görevlerini getirir
func (vy *VeriYonetici) AltGorevleriGetir(ctx context.Context, parentID string) ([]*Gorev, error) {
	sorgu := `SELECT id, title, description, status, priority, project_id, parent_id, created_at, updated_at, due_date
	          FROM gorevler WHERE parent_id = ? ORDER BY created_at`

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
			&gorev.Title,
			&gorev.Description,
			&gorev.Status,
			&gorev.Priority,
			&projeID,
			&parentID,
			&gorev.CreatedAt,
			&gorev.UpdatedAt,
			&gorev.DueDate,
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
		gorev.Tags = etiketler

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

// TumAltGorevleriGetir belirtilen görevin tüm alt görev hiyerarşisini getirir (recursive)
func (vy *VeriYonetici) TumAltGorevleriGetir(ctx context.Context, parentID string) ([]*Gorev, error) {
	sorgu := `
		WITH RECURSIVE alt_gorevler AS (
			SELECT id, title, description, status, priority, project_id, parent_id, 
			       created_at, updated_at, due_date, 1 as level
			FROM gorevler
			WHERE parent_id = ?
			
			UNION ALL
			
			SELECT g.id, g.title, g.description, g.status, g.priority, g.project_id, g.parent_id,
			       g.created_at, g.updated_at, g.due_date, ag.level + 1
			FROM gorevler g
			INNER JOIN alt_gorevler ag ON g.parent_id = ag.id
		)
		SELECT * FROM alt_gorevler ORDER BY level, created_at`

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
			&gorev.Title,
			&gorev.Description,
			&gorev.Status,
			&gorev.Priority,
			&projeID,
			&parentID,
			&gorev.CreatedAt,
			&gorev.UpdatedAt,
			&gorev.DueDate,
			&gorev.Level,
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
		gorev.Tags = etiketler

		gorevler = append(gorevler, gorev)
	}

	return gorevler, nil
}

// UstGorevleriGetir belirtilen görevin tüm üst görev hiyerarşisini getirir
func (vy *VeriYonetici) UstGorevleriGetir(ctx context.Context, gorevID string) ([]*Gorev, error) {
	sorgu := `
		WITH RECURSIVE ust_gorevler AS (
			SELECT g2.id, g2.title, g2.description, g2.status, g2.priority, g2.project_id, g2.parent_id,
			       g2.created_at, g2.updated_at, g2.due_date
			FROM gorevler g1
			JOIN gorevler g2 ON g1.parent_id = g2.id
			WHERE g1.id = ?
			
			UNION ALL
			
			SELECT g.id, g.title, g.description, g.status, g.priority, g.project_id, g.parent_id,
			       g.created_at, g.updated_at, g.due_date
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
			&gorev.Title,
			&gorev.Description,
			&gorev.Status,
			&gorev.Priority,
			&projeID,
			&parentID,
			&gorev.CreatedAt,
			&gorev.UpdatedAt,
			&gorev.DueDate,
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
func (vy *VeriYonetici) GorevHiyerarsiGetir(ctx context.Context, gorevID string) (*GorevHiyerarsi, error) {
	// Ana görevi getir
	gorev, err := vy.GorevGetir(ctx, gorevID)
	if err != nil {
		return nil, err
	}

	// Üst görevleri getir
	ustGorevler, err := vy.UstGorevleriGetir(ctx, gorevID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Alt görev istatistiklerini hesapla
	sorgu := `
		WITH RECURSIVE alt_gorevler AS (
			SELECT id, status
			FROM gorevler
			WHERE parent_id = ?
			
			UNION ALL
			
			SELECT g.id, g.status
			FROM gorevler g
			INNER JOIN alt_gorevler ag ON g.parent_id = ag.id
		)
		SELECT 
			COUNT(*) as toplam,
			COALESCE(SUM(CASE WHEN status = 'tamamlandi' THEN 1 ELSE 0 END), 0) as tamamlanan,
			COALESCE(SUM(CASE WHEN status = 'devam_ediyor' THEN 1 ELSE 0 END), 0) as devam_eden,
			COALESCE(SUM(CASE WHEN status = 'beklemede' THEN 1 ELSE 0 END), 0) as beklemede
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
	} else if gorev.Status == constants.TaskStatusCompleted {
		ilerlemeYuzdesi = 100
	}

	return &GorevHiyerarsi{
		Gorev:              gorev,
		ParentTasks:        ustGorevler,
		TotalSubtasks:      toplam,
		CompletedSubtasks:  tamamlanan,
		InProgressSubtasks: devamEden,
		PendingSubtasks:    beklemede,
		ProgressPercentage: ilerlemeYuzdesi,
	}, nil
}

// ParentIDGuncelle bir görevin parent_id'sini günceller
func (vy *VeriYonetici) ParentIDGuncelle(ctx context.Context, gorevID, yeniParentID string) error {
	// Önce circular dependency kontrolü yap
	if yeniParentID != "" {
		daireVar, err := vy.DaireBagimliligiKontrolEt(ctx, gorevID, yeniParentID)
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
		sorgu = `UPDATE gorevler SET parent_id = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = vy.db.Exec(sorgu, gorevID)
	} else {
		sorgu = `UPDATE gorevler SET parent_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = vy.db.Exec(sorgu, yeniParentID, gorevID)
	}

	return err
}

// DaireBagimliligiKontrolEt bir görevin belirtilen parent'a taşınması durumunda dairesel bağımlılık oluşup oluşmayacağını kontrol eder
func (vy *VeriYonetici) DaireBagimliligiKontrolEt(ctx context.Context, gorevID, hedefParentID string) (bool, error) {
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
		SELECT target_id, COUNT(*) as bagli_sayi
		FROM baglantilar 
		WHERE target_id IN (%s)
		GROUP BY target_id
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
		SELECT b.target_id, COUNT(*) as tamamlanmamis_sayi
		FROM baglantilar b
		INNER JOIN gorevler g ON b.source_id = g.id
		WHERE b.target_id IN (%s) AND g.status != 'tamamlandi'
		GROUP BY b.target_id
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
		SELECT b.source_id, COUNT(*) as bagimli_sayi
		FROM baglantilar b
		WHERE b.source_id IN (%s)
		GROUP BY b.source_id
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
		INSERT INTO ai_interactions (task_id, action_type, context, timestamp)
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
		SELECT id, task_id, action_type, context, timestamp
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
		SELECT id, task_id, action_type, context, timestamp
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
	// Insert file path for task into task_file_paths table
	query := `
		INSERT INTO task_file_paths (task_id, file_path)
		VALUES (?, ?)
		ON CONFLICT(task_id, file_path) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
	`
	_, err := vy.db.Exec(query, taskID, path)
	if err != nil {
		return fmt.Errorf(i18n.T("error.filePathAddFailed", map[string]interface{}{"Error": err}))
	}
	return nil
}

// GorevDosyaYoluSil removes a file path from a task
func (vy *VeriYonetici) GorevDosyaYoluSil(taskID string, path string) error {
	// Delete file path from task_file_paths table
	query := `
		DELETE FROM task_file_paths
		WHERE task_id = ? AND file_path = ?
	`
	result, err := vy.db.Exec(query, taskID, path)
	if err != nil {
		return fmt.Errorf(i18n.T("error.filePathRemoveFailed", map[string]interface{}{"Error": err}))
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(i18n.T("error.filePathNotFound"))
	}

	return nil
}

// GorevDosyaYollariGetir gets all file paths for a task
func (vy *VeriYonetici) GorevDosyaYollariGetir(taskID string) ([]string, error) {
	// Query file paths from task_file_paths table
	query := `
		SELECT file_path
		FROM task_file_paths
		WHERE task_id = ?
		ORDER BY created_at ASC
	`

	rows, err := vy.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.filePathQueryFailed", map[string]interface{}{"Error": err}))
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, fmt.Errorf(i18n.T("error.filePathReadFailed", map[string]interface{}{"Error": err}))
		}
		paths = append(paths, path)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(i18n.T("error.rowIterationError", map[string]interface{}{"Error": err}))
	}

	return paths, nil
}

// DosyaYoluGorevleriGetir gets all tasks associated with a file path
func (vy *VeriYonetici) DosyaYoluGorevleriGetir(path string) ([]string, error) {
	// Query task IDs from task_file_paths table
	query := `
		SELECT task_id
		FROM task_file_paths
		WHERE file_path = ?
		ORDER BY created_at ASC
	`

	rows, err := vy.db.Query(query, path)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.tasksQueryFailed", map[string]interface{}{"Error": err}))
	}
	defer rows.Close()

	var taskIDs []string
	for rows.Next() {
		var taskID string
		if err := rows.Scan(&taskID); err != nil {
			return nil, fmt.Errorf(i18n.T("error.taskIdReadFailed", map[string]interface{}{"Error": err}))
		}
		taskIDs = append(taskIDs, taskID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(i18n.T("error.rowIterationError", map[string]interface{}{"Error": err}))
	}

	return taskIDs, nil
}

// repairMigrationStateIfNeeded checks if existing tables exist but migration state is missing
// and repairs the migration state to prevent "table already exists" errors
func (vy *VeriYonetici) repairMigrationStateIfNeeded() error {
	log.Printf("DEBUG: Checking if migration state repair is needed")

	// Check if there are any migrations recorded
	var migrationCount int
	err := vy.db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&migrationCount)
	if err != nil {
		log.Printf("DEBUG: Could not check migration count: %v", err)
		return nil // Not critical, continue
	}

	// If we have migrations recorded, no repair needed
	if migrationCount > 0 {
		log.Printf("DEBUG: Found %d recorded migrations, no repair needed", migrationCount)
		return nil
	}

	log.Printf("DEBUG: No migrations recorded, checking for existing tables")

	// Check if core tables exist (indicating migrations were run before)
	tables := []string{"projeler", "gorevler", "baglantilar"}
	tablesExist := 0

	for _, table := range tables {
		var exists int
		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
		err := vy.db.QueryRow(query, table).Scan(&exists)
		if err == nil && exists > 0 {
			tablesExist++
			log.Printf("DEBUG: Table %s exists", table)
		}
	}

	// If most core tables exist, repair migration state
	if tablesExist >= 2 {
		log.Printf("DEBUG: Found %d core tables, repairing migration state", tablesExist)
		return vy.repairMigrationState(tablesExist)
	}

	log.Printf("DEBUG: Only %d core tables exist, normal migration will proceed", tablesExist)
	return nil
}

// repairMigrationState adds migration records for already applied migrations
func (vy *VeriYonetici) repairMigrationState(tablesExist int) error {
	log.Printf("DEBUG: Starting migration state repair")

	// Define expected migrations based on table existence
	migrations := []struct {
		version int
		tables  []string
	}{
		{1, []string{"projeler", "gorevler", "baglantilar", "etiketler", "gorev_etiketleri"}},
		{2, []string{"gorevler"}}, // due_date column
		{3, []string{"etiketler", "gorev_etiketleri"}},
		{4, []string{"gorev_templateleri"}},
		{5, []string{"gorevler"}}, // parent_id column
		{6, []string{"ai_interactions", "ai_context", "aktif_proje"}},
		{7, []string{"baglantilar"}}, // indexes
		{8, []string{"file_watches", "file_changes"}},
		{9, []string{"gorev_templateleri"}}, // alias column
		{10, []string{"gorevler_fts", "filter_profiles", "search_history"}},
	}

	for _, migration := range migrations {
		// Check if tables for this migration exist
		allTablesExist := true
		for _, table := range migration.tables {
			var exists int
			query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
			err := vy.db.QueryRow(query, table).Scan(&exists)
			if err != nil || exists == 0 {
				allTablesExist = false
				break
			}
		}

		// Special handling for migration 1 - if core tables exist, mark it as applied
		if migration.version == 1 && tablesExist >= 2 {
			allTablesExist = true
		}

		if allTablesExist {
			// Mark this migration as applied
			_, err := vy.db.Exec("INSERT OR IGNORE INTO schema_migrations (version) VALUES (?)", migration.version)
			if err != nil {
				log.Printf("WARNING: Failed to record migration %d during repair: %v", migration.version, err)
			} else {
				log.Printf("DEBUG: Marked migration %d as applied during repair", migration.version)
			}
		}
	}

	log.Printf("DEBUG: Migration state repair completed")
	return nil
}
