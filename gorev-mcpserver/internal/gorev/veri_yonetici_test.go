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
			err := vy.GorevGuncelle(tc.gorev)
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
