package gorev

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYeniVeriYonetici(t *testing.T) {
	testCases := []struct {
		name    string
		dbYolu  string
		wantErr bool
	}{
		{
			name:    "valid memory database",
			dbYolu:  ":memory:",
			wantErr: false,
		},
		{
			name:    "invalid database path",
			dbYolu:  "/invalid\x00path/db.sqlite",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vy, err := YeniVeriYonetici(tc.dbYolu, "file://../../internal/veri/migrations")
			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer vy.Kapat()
		})
	}
}

func TestVeriYonetici_GorevKaydet(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	testCases := []struct {
		name    string
		gorev   *Gorev
		wantErr bool
	}{
		{
			name: "valid task",
			gorev: &Gorev{
				ID:              "test-1",
				Baslik:          "Test Görevi",
				Aciklama:        "Test açıklaması",
				Durum:           "beklemede",
				Oncelik:         "orta",
				OlusturmaTarih:  time.Now(),
				GuncellemeTarih: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "task with project",
			gorev: &Gorev{
				ID:              "test-2",
				Baslik:          "Proje Görevi",
				Aciklama:        "Proje ile ilişkili görev",
				Durum:           "devam-ediyor",
				Oncelik:         "yuksek",
				ProjeID:         "proje-1",
				OlusturmaTarih:  time.Now(),
				GuncellemeTarih: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "duplicate ID",
			gorev: &Gorev{
				ID:              "test-1",
				Baslik:          "Duplicate",
				Aciklama:        "Should fail",
				Durum:           "beklemede",
				Oncelik:         "orta",
				OlusturmaTarih:  time.Now(),
				GuncellemeTarih: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty ID",
			gorev: &Gorev{
				ID:              "",
				Baslik:          "No ID",
				Aciklama:        "Empty ID allowed in SQLite",
				Durum:           "beklemede",
				Oncelik:         "orta",
				OlusturmaTarih:  time.Now(),
				GuncellemeTarih: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := vy.GorevKaydet(tc.gorev)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestVeriYonetici_GorevGetir(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Insert test data
	testGorev := &Gorev{
		ID:              "test-get-1",
		Baslik:          "Test Getir",
		Aciklama:        "Getir testi",
		Durum:           "beklemede",
		Oncelik:         "orta",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	if err := vy.GorevKaydet(testGorev); err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	testCases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "existing task",
			id:      "test-get-1",
			wantErr: false,
		},
		{
			name:    "non-existing task",
			id:      "non-existing",
			wantErr: true,
		},
		{
			name:    "empty ID",
			id:      "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gorev, err := vy.GorevGetir(tc.id)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if gorev.ID != tc.id {
				t.Errorf("expected ID %s, got %s", tc.id, gorev.ID)
			}
			if gorev.Baslik != testGorev.Baslik {
				t.Errorf("expected Baslik %s, got %s", testGorev.Baslik, gorev.Baslik)
			}
		})
	}
}

func TestVeriYonetici_GorevleriGetir(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Insert test data
	testGorevler := []*Gorev{
		{
			ID:              "test-list-1",
			Baslik:          "Bekleyen Görev",
			Durum:           "beklemede",
			Oncelik:         "orta",
			OlusturmaTarih:  time.Now().Add(-2 * time.Hour),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "test-list-2",
			Baslik:          "Devam Eden Görev",
			Durum:           "devam-ediyor",
			Oncelik:         "yuksek",
			OlusturmaTarih:  time.Now().Add(-1 * time.Hour),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "test-list-3",
			Baslik:          "Tamamlanan Görev",
			Durum:           "tamamlandi",
			Oncelik:         "dusuk",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
	}

	for _, gorev := range testGorevler {
		if err := vy.GorevKaydet(gorev); err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	testCases := []struct {
		name          string
		durum         string
		expectedCount int
	}{
		{
			name:          "all tasks",
			durum:         "",
			expectedCount: 3,
		},
		{
			name:          "beklemede tasks",
			durum:         "beklemede",
			expectedCount: 1,
		},
		{
			name:          "devam-ediyor tasks",
			durum:         "devam-ediyor",
			expectedCount: 1,
		},
		{
			name:          "tamamlandi tasks",
			durum:         "tamamlandi",
			expectedCount: 1,
		},
		{
			name:          "non-existing status",
			durum:         "iptal",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gorevler, err := vy.GorevleriGetir(tc.durum, "", "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(gorevler) != tc.expectedCount {
				t.Errorf("expected %d tasks, got %d", tc.expectedCount, len(gorevler))
			}

			// Verify order (newest first)
			if len(gorevler) > 1 {
				for i := 0; i < len(gorevler)-1; i++ {
					if gorevler[i].OlusturmaTarih.Before(gorevler[i+1].OlusturmaTarih) {
						t.Error("tasks not ordered by creation date (newest first)")
					}
				}
			}
		})
	}
}

func TestVeriYonetici_GorevGuncelle(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Insert test data
	originalGorev := &Gorev{
		ID:              "test-update-1",
		Baslik:          "Original Title",
		Aciklama:        "Original Description",
		Durum:           "beklemede",
		Oncelik:         "orta",
		OlusturmaTarih:  time.Now().Add(-1 * time.Hour),
		GuncellemeTarih: time.Now().Add(-1 * time.Hour),
	}
	if err := vy.GorevKaydet(originalGorev); err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	testCases := []struct {
		name    string
		gorev   *Gorev
		wantErr bool
	}{
		{
			name: "update all fields",
			gorev: &Gorev{
				ID:              "test-update-1",
				Baslik:          "Updated Title",
				Aciklama:        "Updated Description",
				Durum:           "devam-ediyor",
				Oncelik:         "yuksek",
				ProjeID:         "proje-1",
				GuncellemeTarih: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "update non-existing task",
			gorev: &Gorev{
				ID:              "non-existing",
				Baslik:          "Should not update",
				Durum:           "beklemede",
				Oncelik:         "orta",
				GuncellemeTarih: time.Now(),
			},
			wantErr: false, // SQL UPDATE doesn't error on non-existing rows
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert gorev struct to map for GorevGuncelle
			params := map[string]interface{}{
				"baslik":           tc.gorev.Baslik,
				"aciklama":         tc.gorev.Aciklama,
				"durum":            tc.gorev.Durum,
				"oncelik":          tc.gorev.Oncelik,
				"proje_id":         tc.gorev.ProjeID,
				"guncelleme_tarih": tc.gorev.GuncellemeTarih,
			}

			err := vy.GorevGuncelle(tc.gorev.ID, params)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify update
			if tc.gorev.ID == "test-update-1" {
				updated, err := vy.GorevGetir(tc.gorev.ID)
				if err != nil {
					t.Fatalf("failed to get updated task: %v", err)
				}
				if updated.Baslik != tc.gorev.Baslik {
					t.Errorf("expected Baslik %s, got %s", tc.gorev.Baslik, updated.Baslik)
				}
				if updated.Durum != tc.gorev.Durum {
					t.Errorf("expected Durum %s, got %s", tc.gorev.Durum, updated.Durum)
				}
			}
		})
	}
}

func TestVeriYonetici_GorevSil(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Insert test data
	testGorev := &Gorev{
		ID:              "test-delete-1",
		Baslik:          "To be deleted",
		Durum:           "beklemede",
		Oncelik:         "orta",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	if err := vy.GorevKaydet(testGorev); err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	testCases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "delete existing task",
			id:      "test-delete-1",
			wantErr: false,
		},
		{
			name:    "delete non-existing task",
			id:      "non-existing",
			wantErr: true,
		},
		{
			name:    "delete with empty ID",
			id:      "",
			wantErr: true, // Will fail because no rows affected
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := vy.GorevSil(tc.id)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify deletion
			_, err = vy.GorevGetir(tc.id)
			if err == nil {
				t.Error("task still exists after deletion")
			}
		})
	}
}

func TestVeriYonetici_ProjeKaydet(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	testCases := []struct {
		name    string
		proje   *Proje
		wantErr bool
	}{
		{
			name: "valid project",
			proje: &Proje{
				ID:              "proje-1",
				Isim:            "Test Projesi",
				Tanim:           "Test projesi açıklaması",
				OlusturmaTarih:  time.Now(),
				GuncellemeTarih: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "duplicate project ID",
			proje: &Proje{
				ID:              "proje-1",
				Isim:            "Duplicate",
				OlusturmaTarih:  time.Now(),
				GuncellemeTarih: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty project ID",
			proje: &Proje{
				ID:              "",
				Isim:            "No ID",
				OlusturmaTarih:  time.Now(),
				GuncellemeTarih: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := vy.ProjeKaydet(tc.proje)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestVeriYonetici_ProjeGetir(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Insert test data
	testProje := &Proje{
		ID:              "proje-get-1",
		Isim:            "Test Projesi",
		Tanim:           "Test açıklaması",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	if err := vy.ProjeKaydet(testProje); err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	testCases := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "existing project",
			id:      "proje-get-1",
			wantErr: false,
		},
		{
			name:    "non-existing project",
			id:      "non-existing",
			wantErr: true,
		},
		{
			name:    "empty ID",
			id:      "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			proje, err := vy.ProjeGetir(tc.id)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if proje.ID != tc.id {
				t.Errorf("expected ID %s, got %s", tc.id, proje.ID)
			}
			if proje.Isim != testProje.Isim {
				t.Errorf("expected Isim %s, got %s", testProje.Isim, proje.Isim)
			}
		})
	}
}

func TestVeriYonetici_ProjeleriGetir(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Insert test data
	testProjeler := []*Proje{
		{
			ID:              "proje-list-1",
			Isim:            "Proje 1",
			Tanim:           "İlk proje",
			OlusturmaTarih:  time.Now().Add(-2 * time.Hour),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "proje-list-2",
			Isim:            "Proje 2",
			Tanim:           "İkinci proje",
			OlusturmaTarih:  time.Now().Add(-1 * time.Hour),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "proje-list-3",
			Isim:            "Proje 3",
			Tanim:           "Üçüncü proje",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
	}

	for _, proje := range testProjeler {
		if err := vy.ProjeKaydet(proje); err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	projeler, err := vy.ProjeleriGetir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projeler) != 3 {
		t.Errorf("expected 3 projects, got %d", len(projeler))
	}

	// Verify order (newest first)
	for i := 0; i < len(projeler)-1; i++ {
		if projeler[i].OlusturmaTarih.Before(projeler[i+1].OlusturmaTarih) {
			t.Error("projects not ordered by creation date (newest first)")
		}
	}
}

func TestVeriYonetici_ProjeGorevleriGetir(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Insert test project
	testProje := &Proje{
		ID:              "proje-tasks-1",
		Isim:            "Test Projesi",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	if err := vy.ProjeKaydet(testProje); err != nil {
		t.Fatalf("failed to insert test project: %v", err)
	}

	// Insert tasks for the project
	testGorevler := []*Gorev{
		{
			ID:              "task-1",
			Baslik:          "Proje Görevi 1",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         "proje-tasks-1",
			OlusturmaTarih:  time.Now().Add(-2 * time.Hour),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "task-2",
			Baslik:          "Proje Görevi 2",
			Durum:           "devam-ediyor",
			Oncelik:         "yuksek",
			ProjeID:         "proje-tasks-1",
			OlusturmaTarih:  time.Now().Add(-1 * time.Hour),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "task-3",
			Baslik:          "Başka Proje Görevi",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         "other-project",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "task-4",
			Baslik:          "Projesi Olmayan Görev",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         "",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
	}

	for _, gorev := range testGorevler {
		if err := vy.GorevKaydet(gorev); err != nil {
			t.Fatalf("failed to insert test task: %v", err)
		}
	}

	testCases := []struct {
		name          string
		projeID       string
		expectedCount int
	}{
		{
			name:          "existing project with tasks",
			projeID:       "proje-tasks-1",
			expectedCount: 2,
		},
		{
			name:          "non-existing project",
			projeID:       "non-existing",
			expectedCount: 0,
		},
		{
			name:          "empty project ID",
			projeID:       "",
			expectedCount: 1, // Will return the task with empty ProjeID
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gorevler, err := vy.ProjeGorevleriGetir(tc.projeID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(gorevler) != tc.expectedCount {
				t.Errorf("expected %d tasks, got %d", tc.expectedCount, len(gorevler))
			}

			// Verify all tasks belong to the project
			for _, gorev := range gorevler {
				if gorev.ProjeID != tc.projeID {
					t.Errorf("task %s has wrong project ID: expected %s, got %s",
						gorev.ID, tc.projeID, gorev.ProjeID)
				}
			}

			// Verify order (newest first)
			if len(gorevler) > 1 {
				for i := 0; i < len(gorevler)-1; i++ {
					if gorevler[i].OlusturmaTarih.Before(gorevler[i+1].OlusturmaTarih) {
						t.Error("tasks not ordered by creation date (newest first)")
					}
				}
			}
		})
	}
}

func TestVeriYonetici_Kapat(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}

	// Close should work
	err = vy.Kapat()
	if err != nil {
		t.Errorf("unexpected error closing database: %v", err)
	}

	// Operations after close should fail
	_, err = vy.GorevleriGetir("", "", "")
	if err == nil {
		t.Error("expected error after closing database, but got nil")
	}
}

func TestVeriYonetici_ConcurrentAccess(t *testing.T) {
	// Use a temporary file database for concurrent access testing
	// as :memory: databases don't support true concurrency in SQLite
	tempDB := fmt.Sprintf("/tmp/test_concurrent_%d.db", time.Now().UnixNano())
	defer func() {
		// Clean up temp file
		os.Remove(tempDB)
	}()

	vy, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Test concurrent writes
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			gorev := &Gorev{
				ID:              fmt.Sprintf("concurrent-%d", id),
				Baslik:          fmt.Sprintf("Concurrent Task %d", id),
				Durum:           "beklemede",
				Oncelik:         "orta",
				OlusturmaTarih:  time.Now(),
				GuncellemeTarih: time.Now(),
			}
			done <- vy.GorevKaydet(gorev)
		}(i)
	}

	// Wait for all goroutines and check errors
	var errors []error
	successCount := 0
	for i := 0; i < 10; i++ {
		if err := <-done; err != nil {
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	// Allow some concurrent access failures, but at least 50% should succeed
	if successCount < 5 {
		t.Errorf("Too many concurrent access failures. Success: %d/10, Errors: %v", successCount, errors)
	}

	// Verify tasks were created
	gorevler, err := vy.GorevleriGetir("", "", "")
	if err != nil {
		t.Fatalf("failed to get tasks: %v", err)
	}
	if len(gorevler) < 1 {
		t.Errorf("expected at least 1 task, got %d", len(gorevler))
	}
}

func TestVeriYonetici_SQLInjection(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Try SQL injection in various fields
	maliciousGorev := &Gorev{
		ID:              "test-injection",
		Baslik:          "'; DROP TABLE gorevler; --",
		Aciklama:        "' OR '1'='1",
		Durum:           "beklemede",
		Oncelik:         "orta",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}

	// Should save without executing the injection
	err = vy.GorevKaydet(maliciousGorev)
	if err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Verify table still exists
	var count int
	err = vy.db.QueryRow("SELECT COUNT(*) FROM gorevler").Scan(&count)
	if err != nil {
		t.Error("gorevler table was dropped - SQL injection succeeded!")
	}

	// Verify the malicious string was stored as data, not executed
	retrieved, err := vy.GorevGetir("test-injection")
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}
	if retrieved.Baslik != maliciousGorev.Baslik {
		t.Error("malicious string was not stored correctly")
	}

	// Try injection in filter parameter
	_, err = vy.GorevleriGetir("'; DROP TABLE gorevler; --", "", "")
	if err != nil {
		t.Errorf("query failed: %v", err)
	}

	// Verify table still exists
	err = vy.db.QueryRow("SELECT COUNT(*) FROM gorevler").Scan(&count)
	if err != nil {
		t.Error("gorevler table was dropped - SQL injection succeeded!")
	}
}

func TestVeriYonetici_NullHandling(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Insert task without project (NULL proje_id)
	gorevWithoutProject := &Gorev{
		ID:              "no-project",
		Baslik:          "Task without project",
		Aciklama:        "", // Empty description
		Durum:           "beklemede",
		Oncelik:         "orta",
		ProjeID:         "", // Should be stored as NULL
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}

	err = vy.GorevKaydet(gorevWithoutProject)
	if err != nil {
		t.Fatalf("failed to save task: %v", err)
	}

	// Retrieve and verify NULL handling
	retrieved, err := vy.GorevGetir("no-project")
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}

	if retrieved.ProjeID != "" {
		t.Errorf("expected empty ProjeID for NULL value, got %s", retrieved.ProjeID)
	}

	// Verify in list query
	gorevler, err := vy.GorevleriGetir("", "", "")
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}

	found := false
	for _, g := range gorevler {
		if g.ID == "no-project" {
			found = true
			if g.ProjeID != "" {
				t.Errorf("expected empty ProjeID in list, got %s", g.ProjeID)
			}
			break
		}
	}
	if !found {
		t.Error("task not found in list")
	}
}

// Helper function to compare times (ignoring nanoseconds)
func timesEqual(t1, t2 time.Time) bool {
	return t1.Unix() == t2.Unix()
}

func TestVeriYonetici_Etiketleme(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer vy.Kapat()

	// Görev oluştur
	gorev := &Gorev{
		ID:     "etiket-test-gorev",
		Baslik: "Etiket Testi",
	}
	err = vy.GorevKaydet(gorev)
	require.NoError(t, err)

	// 1. Yeni etiketler oluştur ve getir
	isimler := []string{"bug", "acil", "  yeni-ozellik  "}
	etiketler, err := vy.EtiketleriGetirVeyaOlustur(isimler)
	require.NoError(t, err)
	require.Len(t, etiketler, 3)
	assert.Equal(t, "bug", etiketler[0].Isim)
	assert.Equal(t, "acil", etiketler[1].Isim)
	assert.Equal(t, "yeni-ozellik", etiketler[2].Isim) // Boşlukların temizlendiğini doğrula

	// 2. Göreve etiketleri ata
	err = vy.GorevEtiketleriniAyarla(gorev.ID, etiketler)
	require.NoError(t, err)

	// 3. Görevi getir ve etiketleri doğrula
	getirilenGorev, err := vy.GorevGetir(gorev.ID)
	require.NoError(t, err)
	require.NotNil(t, getirilenGorev)
	require.Len(t, getirilenGorev.Etiketler, 3)

	// Etiket isimlerini bir map'e koyarak kontrol et
	etiketMap := make(map[string]bool)
	for _, e := range getirilenGorev.Etiketler {
		etiketMap[e.Isim] = true
	}
	assert.True(t, etiketMap["bug"])
	assert.True(t, etiketMap["acil"])
	assert.True(t, etiketMap["yeni-ozellik"])

	// 4. Etiketleri güncelle (birini çıkar, birini ekle)
	yeniIsimler := []string{"acil", "dokumantasyon"}
	yeniEtiketler, err := vy.EtiketleriGetirVeyaOlustur(yeniIsimler)
	require.NoError(t, err)
	err = vy.GorevEtiketleriniAyarla(gorev.ID, yeniEtiketler)
	require.NoError(t, err)

	// 5. Güncellenmiş görevi getir ve etiketleri doğrula
	getirilenGorev, err = vy.GorevGetir(gorev.ID)
	require.NoError(t, err)
	require.NotNil(t, getirilenGorev)
	require.Len(t, getirilenGorev.Etiketler, 2)

	yeniEtiketMap := make(map[string]bool)
	for _, e := range getirilenGorev.Etiketler {
		yeniEtiketMap[e.Isim] = true
	}
	assert.False(t, yeniEtiketMap["bug"])
	assert.True(t, yeniEtiketMap["acil"])
	assert.True(t, yeniEtiketMap["dokumantasyon"])
}

// TestVeriYonetici_BulkDependencyCounts tests the new bulk dependency count methods
func TestVeriYonetici_BulkDependencyCounts(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Create test tasks
	tasks := []*Gorev{
		{
			ID:              "task-1",
			Baslik:          "Task 1",
			Durum:           "beklemede",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "task-2",
			Baslik:          "Task 2",
			Durum:           "tamamlandi",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "task-3",
			Baslik:          "Task 3",
			Durum:           "beklemede",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "task-4",
			Baslik:          "Task 4",
			Durum:           "devam_ediyor",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
	}

	for _, task := range tasks {
		if err := vy.GorevKaydet(task); err != nil {
			t.Fatalf("failed to save task %s: %v", task.ID, err)
		}
	}

	// Create dependencies:
	// task-1 depends on task-2 (completed)
	// task-3 depends on task-1 (not completed)
	// task-4 depends on task-2 (completed) and task-3 (not completed)
	dependencies := []*Baglanti{
		{
			ID:          "dep-1",
			KaynakID:    "task-2",
			HedefID:     "task-1",
			BaglantiTip: "onceki",
		},
		{
			ID:          "dep-2",
			KaynakID:    "task-1",
			HedefID:     "task-3",
			BaglantiTip: "onceki",
		},
		{
			ID:          "dep-3",
			KaynakID:    "task-2",
			HedefID:     "task-4",
			BaglantiTip: "onceki",
		},
		{
			ID:          "dep-4",
			KaynakID:    "task-3",
			HedefID:     "task-4",
			BaglantiTip: "onceki",
		},
	}

	for _, dep := range dependencies {
		if err := vy.BaglantiEkle(dep); err != nil {
			t.Fatalf("failed to add dependency %s: %v", dep.ID, err)
		}
	}

	taskIDs := []string{"task-1", "task-2", "task-3", "task-4"}

	// Test BulkBagimlilikSayilariGetir (tasks that each task depends on)
	t.Run("BulkBagimlilikSayilariGetir", func(t *testing.T) {
		counts, err := vy.BulkBagimlilikSayilariGetir(taskIDs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := map[string]int{
			"task-1": 1, // depends on task-2
			"task-2": 0, // depends on nothing
			"task-3": 1, // depends on task-1
			"task-4": 2, // depends on task-2 and task-3
		}

		for taskID, expectedCount := range expected {
			if counts[taskID] != expectedCount {
				t.Errorf("task %s: expected %d dependencies, got %d", taskID, expectedCount, counts[taskID])
			}
		}
	})

	// Test BulkTamamlanmamiaBagimlilikSayilariGetir (incomplete dependencies)
	t.Run("BulkTamamlanmamiaBagimlilikSayilariGetir", func(t *testing.T) {
		counts, err := vy.BulkTamamlanmamiaBagimlilikSayilariGetir(taskIDs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := map[string]int{
			"task-1": 0, // depends on task-2 (completed)
			"task-2": 0, // depends on nothing
			"task-3": 1, // depends on task-1 (not completed)
			"task-4": 1, // depends on task-2 (completed) and task-3 (not completed) = 1 incomplete
		}

		for taskID, expectedCount := range expected {
			if counts[taskID] != expectedCount {
				t.Errorf("task %s: expected %d incomplete dependencies, got %d", taskID, expectedCount, counts[taskID])
			}
		}
	})

	// Test BulkBuGoreveBagimliSayilariGetir (tasks that depend on each task)
	t.Run("BulkBuGoreveBagimliSayilariGetir", func(t *testing.T) {
		counts, err := vy.BulkBuGoreveBagimliSayilariGetir(taskIDs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := map[string]int{
			"task-1": 1, // task-3 depends on task-1
			"task-2": 2, // task-1 and task-4 depend on task-2
			"task-3": 1, // task-4 depends on task-3
			"task-4": 0, // nothing depends on task-4
		}

		for taskID, expectedCount := range expected {
			if counts[taskID] != expectedCount {
				t.Errorf("task %s: expected %d dependent tasks, got %d", taskID, expectedCount, counts[taskID])
			}
		}
	})

	// Test with empty slice
	t.Run("EmptySlice", func(t *testing.T) {
		counts, err := vy.BulkBagimlilikSayilariGetir([]string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(counts) != 0 {
			t.Errorf("expected empty map, got %v", counts)
		}
	})

	// Test with non-existent task IDs
	t.Run("NonExistentTasks", func(t *testing.T) {
		counts, err := vy.BulkBagimlilikSayilariGetir([]string{"non-existent-1", "non-existent-2"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should return empty counts for non-existent tasks
		if len(counts) != 0 {
			t.Errorf("expected empty map for non-existent tasks, got %v", counts)
		}
	})
}

// ==================== PERFORMANCE AND CONCURRENCY TESTS ====================

// Benchmark functions for performance testing
func BenchmarkVeriYonetici_GorevKaydet(b *testing.B) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		b.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gorev := &Gorev{
			ID:              fmt.Sprintf("benchmark-%d", i),
			Baslik:          fmt.Sprintf("Benchmark Task %d", i),
			Aciklama:        "This is a benchmark task for performance testing",
			Durum:           "beklemede",
			Oncelik:         "orta",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		err = vy.GorevKaydet(gorev)
		if err != nil {
			b.Fatalf("failed to save task: %v", err)
		}
	}
}

func BenchmarkVeriYonetici_GorevleriGetir(b *testing.B) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		b.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Pre-populate with some tasks
	for i := 0; i < 100; i++ {
		gorev := &Gorev{
			ID:              fmt.Sprintf("benchmark-%d", i),
			Baslik:          fmt.Sprintf("Benchmark Task %d", i),
			Durum:           "beklemede",
			Oncelik:         "orta",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		vy.GorevKaydet(gorev)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = vy.GorevleriGetir("", "", "")
		if err != nil {
			b.Fatalf("failed to get tasks: %v", err)
		}
	}
}

func BenchmarkVeriYonetici_ProjeKaydet(b *testing.B) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		b.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proje := &Proje{
			ID:              fmt.Sprintf("benchmark-proje-%d", i),
			Isim:            fmt.Sprintf("Benchmark Project %d", i),
			Tanim:           fmt.Sprintf("Benchmark project description %d", i),
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		err = vy.ProjeKaydet(proje)
		if err != nil {
			b.Fatalf("failed to save project: %v", err)
		}
	}
}

func BenchmarkVeriYonetici_GorevGuncelle(b *testing.B) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		b.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Create a task to update
	gorev := &Gorev{
		ID:              "benchmark-update",
		Baslik:          "Benchmark Update Task",
		Durum:           "beklemede",
		Oncelik:         "orta",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	err = vy.GorevKaydet(gorev)
	if err != nil {
		b.Fatalf("failed to create task: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taskID := fmt.Sprintf("benchmark-update-%d", i)
		gorev := &Gorev{
			ID:              taskID,
			Baslik:          "Benchmark Update Task",
			Durum:           "beklemede",
			Oncelik:         "orta",
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		// Create task first
		err = vy.GorevKaydet(gorev)
		if err != nil {
			b.Fatalf("failed to create task: %v", err)
		}

		// Then update it
		params := map[string]interface{}{
			"durum": "devam_ediyor",
		}
		err = vy.GorevGuncelle(taskID, params)
		if err != nil {
			b.Fatalf("failed to update task: %v", err)
		}
	}
}

// Advanced concurrency tests
func TestVeriYonetici_HighConcurrencyAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping high concurrency test in short mode")
	}

	// Use temporary file for true concurrency testing
	tempDB := fmt.Sprintf("/tmp/test_high_concurrent_%d.db", time.Now().UnixNano())
	defer os.Remove(tempDB)

	vy, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Test with high concurrency (100 goroutines)
	concurrentOps := 100
	done := make(chan error, concurrentOps)
	startTime := time.Now()

	for i := 0; i < concurrentOps; i++ {
		go func(id int) {
			// Perform a mix of read and write operations
			if id%2 == 0 {
				// Write operation
				gorev := &Gorev{
					ID:              fmt.Sprintf("high-concurrent-%d", id),
					Baslik:          fmt.Sprintf("High Concurrent Task %d", id),
					Durum:           "beklemede",
					Oncelik:         "orta",
					OlusturmaTarih:  time.Now(),
					GuncellemeTarih: time.Now(),
				}
				done <- vy.GorevKaydet(gorev)
			} else {
				// Read operation
				_, err := vy.GorevleriGetir("", "", "")
				done <- err
			}
		}(i)
	}

	// Wait for all operations
	var errors []error
	successCount := 0
	for i := 0; i < concurrentOps; i++ {
		if err := <-done; err != nil {
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	duration := time.Since(startTime)
	t.Logf("Completed %d operations in %v (%.2f ops/sec)", concurrentOps, duration, float64(concurrentOps)/duration.Seconds())

	// Allow some failures, but at least 80% should succeed
	if successCount < concurrentOps*80/100 {
		t.Errorf("Too many high concurrency failures. Success: %d/%d, Errors: %v", successCount, concurrentOps, errors)
	}

	// Verify data integrity
	gorevler, err := vy.GorevleriGetir("", "", "")
	if err != nil {
		t.Fatalf("failed to verify data: %v", err)
	}
	t.Logf("Final task count: %d", len(gorevler))
}

func TestVeriYonetici_MixedOperationsConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping mixed operations concurrency test in short mode")
	}

	// Use temporary file for true concurrency testing
	tempDB := fmt.Sprintf("/tmp/test_mixed_concurrent_%d.db", time.Now().UnixNano())
	defer os.Remove(tempDB)

	vy, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Create some initial data
	for i := 0; i < 20; i++ {
		proje := &Proje{
			ID:              fmt.Sprintf("proje-%d", i),
			Isim:            fmt.Sprintf("Proje %d", i),
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		vy.ProjeKaydet(proje)
	}

	// Test concurrent mixed operations
	workers := 10
	operationsPerWorker := 20
	done := make(chan error, workers*operationsPerWorker)

	for worker := 0; worker < workers; worker++ {
		go func(wid int) {
			for op := 0; op < operationsPerWorker; op++ {
				taskID := fmt.Sprintf("mixed-task-%d-%d", wid, op)

				switch op % 5 {
				case 0: // Create task
					gorev := &Gorev{
						ID:              taskID,
						Baslik:          fmt.Sprintf("Mixed Task %d-%d", wid, op),
						Durum:           "beklemede",
						Oncelik:         "orta",
						ProjeID:         fmt.Sprintf("proje-%d", wid%20),
						OlusturmaTarih:  time.Now(),
						GuncellemeTarih: time.Now(),
					}
					done <- vy.GorevKaydet(gorev)

				case 1: // Update task
					gorev, err := vy.GorevGetir(taskID)
					if err == nil {
						params := map[string]interface{}{
							"durum": "devam_ediyor",
						}
						done <- vy.GorevGuncelle(gorev.ID, params)
					} else {
						done <- nil // Task doesn't exist yet, that's ok
					}

				case 2: // List tasks
					_, err := vy.GorevleriGetir("", "", "")
					done <- err

				case 3: // List projects
					_, err := vy.ProjeleriGetir()
					done <- err

				case 4: // Get task details
					_, err := vy.GorevGetir(taskID)
					done <- err
				}
			}
		}(worker)
	}

	// Wait for all operations
	var errors []error
	successCount := 0
	totalOps := workers * operationsPerWorker

	for i := 0; i < totalOps; i++ {
		if err := <-done; err != nil {
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	// At least 90% should succeed
	if successCount < totalOps*90/100 {
		t.Errorf("Too many mixed operations failures. Success: %d/%d, Errors: %v", successCount, totalOps, errors)
	}

	t.Logf("Mixed operations: %d/%d successful", successCount, totalOps)
}

func TestVeriYonetici_ConnectionPoolStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping connection pool stress test in short mode")
	}

	// Use temporary file for connection pool testing
	tempDB := fmt.Sprintf("/tmp/test_pool_stress_%d.db", time.Now().UnixNano())
	defer os.Remove(tempDB)

	vy, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Rapid connection cycling test
	iterations := 1000
	done := make(chan error, iterations)

	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		go func(id int) {
			// Each goroutine performs rapid successive operations
			for j := 0; j < 10; j++ {
				gorev := &Gorev{
					ID:              fmt.Sprintf("pool-stress-%d-%d", id, j),
					Baslik:          fmt.Sprintf("Pool Stress Task %d-%d", id, j),
					Durum:           "beklemede",
					Oncelik:         "orta",
					OlusturmaTarih:  time.Now(),
					GuncellemeTarih: time.Now(),
				}

				err := vy.GorevKaydet(gorev)
				if err != nil {
					done <- err
					return
				}

				// Immediate read after write
				_, err = vy.GorevGetir(gorev.ID)
				if err != nil {
					done <- err
					return
				}
			}
			done <- nil
		}(i)
	}

	// Wait for completion
	var errors []error
	successCount := 0

	for i := 0; i < iterations; i++ {
		if err := <-done; err != nil {
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	duration := time.Since(startTime)
	totalOps := iterations * 10

	t.Logf("Connection pool stress test: %d total operations in %v (%.2f ops/sec)",
		totalOps, duration, float64(totalOps)/duration.Seconds())
	t.Logf("Success rate: %d/%d (%.1f%%)", successCount, totalOps, float64(successCount)*100/float64(totalOps))

	// Verify no database corruption
	gorevler, err := vy.GorevleriGetir("", "", "")
	if err != nil {
		t.Fatalf("database potentially corrupted after stress test: %v", err)
	}
	t.Logf("Final database integrity check: %d tasks stored", len(gorevler))

	// Should have very high success rate (at least 95%)
	if successCount < totalOps*95/100 {
		t.Errorf("Connection pool stress test had too many failures. Success: %d/%d", successCount, totalOps)
	}
}

// TestVeriYonetici_LongRunningConcurrency tests database behavior under prolonged concurrent load
func TestVeriYonetici_LongRunningConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long running concurrency test")
	}

	tempDB := fmt.Sprintf("/tmp/test_long_running_%d.db", time.Now().UnixNano())
	defer os.Remove(tempDB)

	vy, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("failed to create VeriYonetici: %v", err)
	}
	defer vy.Kapat()

	// Run concurrent operations for 5 seconds
	duration := 5 * time.Second
	done := make(chan bool)
	errorChan := make(chan error, 100)
	operations := make(chan int, 1000)

	// Start workers
	workerCount := 5
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			opCount := 0
			for {
				select {
				case <-done:
					operations <- opCount
					return
				default:
					// Perform random operation
					taskID := fmt.Sprintf("long-running-%d-%d", workerID, opCount)
					gorev := &Gorev{
						ID:              taskID,
						Baslik:          fmt.Sprintf("Long Running Task %d-%d", workerID, opCount),
						Durum:           "beklemede",
						Oncelik:         "orta",
						OlusturmaTarih:  time.Now(),
						GuncellemeTarih: time.Now(),
					}

					err := vy.GorevKaydet(gorev)
					if err != nil {
						errorChan <- fmt.Errorf("worker %d: %v", workerID, err)
					} else {
						opCount++
					}

					// Small delay to simulate realistic load
					time.Sleep(time.Millisecond * time.Duration(10+workerID))
				}
			}
		}(i)
	}

	// Let it run for the specified duration
	time.Sleep(duration)

	// Signal workers to stop
	for i := 0; i < workerCount; i++ {
		done <- true
	}

	// Collect results
	totalOps := 0
	for i := 0; i < workerCount; i++ {
		totalOps += <-operations
	}

	// Collect any errors
	close(errorChan)
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	t.Logf("Long running concurrency test completed in %v", duration)
	t.Logf("Total operations: %d (%.2f ops/sec)", totalOps, float64(totalOps)/duration.Seconds())
	t.Logf("Errors encountered: %d", len(errors))

	if len(errors) > 0 {
		t.Errorf("Long running concurrency test had %d errors: %v", len(errors), errors)
	}
}
