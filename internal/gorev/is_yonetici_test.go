package gorev

import (
	"errors"
	"strings"
	"testing"
)

// MockVeriYonetici is a mock implementation of VeriYonetici for testing
type MockVeriYonetici struct {
	gorevler map[string]*Gorev
	projeler map[string]*Proje

	// Control behavior
	shouldFailGorevKaydet    bool
	shouldFailGorevGetir     bool
	shouldFailGorevGuncelle  bool
	shouldFailGorevSil       bool
	shouldFailProjeKaydet    bool
	shouldFailProjeGetir     bool
	shouldFailGorevleriGetir bool
	shouldFailProjeleriGetir bool
}

func NewMockVeriYonetici() *MockVeriYonetici {
	return &MockVeriYonetici{
		gorevler: make(map[string]*Gorev),
		projeler: make(map[string]*Proje),
	}
}

func (m *MockVeriYonetici) GorevKaydet(gorev *Gorev) error {
	if m.shouldFailGorevKaydet {
		return errors.New("mock error: gorev kaydet failed")
	}
	m.gorevler[gorev.ID] = gorev
	return nil
}

func (m *MockVeriYonetici) GorevGetir(id string) (*Gorev, error) {
	if m.shouldFailGorevGetir {
		return nil, errors.New("mock error: gorev getir failed")
	}
	gorev, ok := m.gorevler[id]
	if !ok {
		return nil, errors.New("görev bulunamadı")
	}
	return gorev, nil
}

func (m *MockVeriYonetici) GorevleriGetir(durum string) ([]*Gorev, error) {
	if m.shouldFailGorevleriGetir {
		return nil, errors.New("mock error: gorevleri getir failed")
	}
	var result []*Gorev
	for _, gorev := range m.gorevler {
		if durum == "" || gorev.Durum == durum {
			result = append(result, gorev)
		}
	}
	return result, nil
}

func (m *MockVeriYonetici) GorevGuncelle(gorev *Gorev) error {
	if m.shouldFailGorevGuncelle {
		return errors.New("mock error: gorev guncelle failed")
	}
	if _, ok := m.gorevler[gorev.ID]; !ok {
		return errors.New("görev bulunamadı")
	}
	m.gorevler[gorev.ID] = gorev
	return nil
}

func (m *MockVeriYonetici) GorevSil(id string) error {
	if m.shouldFailGorevSil {
		return errors.New("mock error: gorev sil failed")
	}
	if _, ok := m.gorevler[id]; !ok {
		return errors.New("görev bulunamadı")
	}
	delete(m.gorevler, id)
	return nil
}

func (m *MockVeriYonetici) ProjeKaydet(proje *Proje) error {
	if m.shouldFailProjeKaydet {
		return errors.New("mock error: proje kaydet failed")
	}
	m.projeler[proje.ID] = proje
	return nil
}

func (m *MockVeriYonetici) ProjeGetir(id string) (*Proje, error) {
	if m.shouldFailProjeGetir {
		return nil, errors.New("mock error: proje getir failed")
	}
	proje, ok := m.projeler[id]
	if !ok {
		return nil, errors.New("proje bulunamadı")
	}
	return proje, nil
}

func (m *MockVeriYonetici) ProjeleriGetir() ([]*Proje, error) {
	if m.shouldFailProjeleriGetir {
		return nil, errors.New("mock error: projeleri getir failed")
	}
	var result []*Proje
	for _, proje := range m.projeler {
		result = append(result, proje)
	}
	return result, nil
}

func (m *MockVeriYonetici) ProjeGorevleriGetir(projeID string) ([]*Gorev, error) {
	var result []*Gorev
	for _, gorev := range m.gorevler {
		if gorev.ProjeID == projeID {
			result = append(result, gorev)
		}
	}
	return result, nil
}

func (m *MockVeriYonetici) Kapat() error {
	return nil
}

// Tests

func TestYeniIsYonetici(t *testing.T) {
	mockVY := NewMockVeriYonetici()
	iy := YeniIsYonetici(mockVY)

	if iy == nil {
		t.Fatal("YeniIsYonetici returned nil")
	}
	if iy.veriYonetici == nil {
		t.Error("veriYonetici not properly set")
	}
}

func TestIsYonetici_GorevOlustur(t *testing.T) {
	testCases := []struct {
		name             string
		baslik           string
		aciklama         string
		oncelik          string
		shouldFailKaydet bool
		wantErr          bool
	}{
		{
			name:     "valid task creation",
			baslik:   "Test Görevi",
			aciklama: "Test açıklaması",
			oncelik:  "orta",
			wantErr:  false,
		},
		{
			name:     "empty title",
			baslik:   "",
			aciklama: "Açıklama",
			oncelik:  "yuksek",
			wantErr:  false, // Business logic doesn't validate empty titles
		},
		{
			name:             "database error",
			baslik:           "Test",
			aciklama:         "Test",
			oncelik:          "orta",
			shouldFailKaydet: true,
			wantErr:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVY := NewMockVeriYonetici()
			mockVY.shouldFailGorevKaydet = tc.shouldFailKaydet
			iy := YeniIsYonetici(mockVY)

			gorev, err := iy.GorevOlustur(tc.baslik, tc.aciklama, tc.oncelik)
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

			// Verify the created task
			if gorev.Baslik != tc.baslik {
				t.Errorf("expected Baslik %s, got %s", tc.baslik, gorev.Baslik)
			}
			if gorev.Aciklama != tc.aciklama {
				t.Errorf("expected Aciklama %s, got %s", tc.aciklama, gorev.Aciklama)
			}
			if gorev.Oncelik != tc.oncelik {
				t.Errorf("expected Oncelik %s, got %s", tc.oncelik, gorev.Oncelik)
			}
			if gorev.Durum != "beklemede" {
				t.Errorf("expected Durum 'beklemede', got %s", gorev.Durum)
			}
			if gorev.ID == "" {
				t.Error("ID should not be empty")
			}

			// Verify it was saved
			if _, ok := mockVY.gorevler[gorev.ID]; !ok {
				t.Error("task was not saved to database")
			}
		})
	}
}

func TestIsYonetici_GorevListele(t *testing.T) {
	mockVY := NewMockVeriYonetici()
	iy := YeniIsYonetici(mockVY)

	// Add test data
	testGorevler := []*Gorev{
		{ID: "1", Baslik: "Görev 1", Durum: "beklemede"},
		{ID: "2", Baslik: "Görev 2", Durum: "devam-ediyor"},
		{ID: "3", Baslik: "Görev 3", Durum: "tamamlandi"},
		{ID: "4", Baslik: "Görev 4", Durum: "beklemede"},
	}
	for _, g := range testGorevler {
		mockVY.gorevler[g.ID] = g
	}

	testCases := []struct {
		name          string
		durum         string
		expectedCount int
		shouldFail    bool
		wantErr       bool
	}{
		{
			name:          "list all tasks",
			durum:         "",
			expectedCount: 4,
		},
		{
			name:          "list beklemede tasks",
			durum:         "beklemede",
			expectedCount: 2,
		},
		{
			name:          "list devam-ediyor tasks",
			durum:         "devam-ediyor",
			expectedCount: 1,
		},
		{
			name:       "database error",
			durum:      "",
			shouldFail: true,
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVY.shouldFailGorevleriGetir = tc.shouldFail

			gorevler, err := iy.GorevListele(tc.durum)
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

			if len(gorevler) != tc.expectedCount {
				t.Errorf("expected %d tasks, got %d", tc.expectedCount, len(gorevler))
			}
		})
	}
}

func TestIsYonetici_GorevDurumGuncelle(t *testing.T) {
	testCases := []struct {
		name             string
		gorevID          string
		yeniDurum        string
		shouldFailGetir  bool
		shouldFailUpdate bool
		wantErr          bool
		expectedError    string
	}{
		{
			name:      "update existing task",
			gorevID:   "existing-task",
			yeniDurum: "devam-ediyor",
			wantErr:   false,
		},
		{
			name:          "non-existing task",
			gorevID:       "non-existing",
			yeniDurum:     "tamamlandi",
			wantErr:       true,
			expectedError: "görev bulunamadı",
		},
		{
			name:            "database getir error",
			gorevID:         "existing-task",
			yeniDurum:       "tamamlandi",
			shouldFailGetir: true,
			wantErr:         true,
		},
		{
			name:             "database update error",
			gorevID:          "existing-task",
			yeniDurum:        "tamamlandi",
			shouldFailUpdate: true,
			wantErr:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVY := NewMockVeriYonetici()
			iy := YeniIsYonetici(mockVY)

			// Add test task
			if tc.gorevID == "existing-task" {
				mockVY.gorevler["existing-task"] = &Gorev{
					ID:      "existing-task",
					Baslik:  "Test Task",
					Durum:   "beklemede",
					Oncelik: "orta",
				}
			}

			mockVY.shouldFailGorevGetir = tc.shouldFailGetir
			mockVY.shouldFailGorevGuncelle = tc.shouldFailUpdate

			err := iy.GorevDurumGuncelle(tc.gorevID, tc.yeniDurum)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				} else if tc.expectedError != "" && !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("expected error containing '%s', got '%s'", tc.expectedError, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify update
			gorev := mockVY.gorevler["existing-task"]
			if gorev.Durum != tc.yeniDurum {
				t.Errorf("expected Durum %s, got %s", tc.yeniDurum, gorev.Durum)
			}
		})
	}
}

func TestIsYonetici_ProjeOlustur(t *testing.T) {
	testCases := []struct {
		name             string
		isim             string
		tanim            string
		shouldFailKaydet bool
		wantErr          bool
	}{
		{
			name:    "valid project creation",
			isim:    "Test Projesi",
			tanim:   "Test proje açıklaması",
			wantErr: false,
		},
		{
			name:    "empty name",
			isim:    "",
			tanim:   "Açıklama",
			wantErr: false, // Business logic doesn't validate empty names
		},
		{
			name:             "database error",
			isim:             "Test",
			tanim:            "Test",
			shouldFailKaydet: true,
			wantErr:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVY := NewMockVeriYonetici()
			mockVY.shouldFailProjeKaydet = tc.shouldFailKaydet
			iy := YeniIsYonetici(mockVY)

			proje, err := iy.ProjeOlustur(tc.isim, tc.tanim)
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

			// Verify the created project
			if proje.Isim != tc.isim {
				t.Errorf("expected Isim %s, got %s", tc.isim, proje.Isim)
			}
			if proje.Tanim != tc.tanim {
				t.Errorf("expected Tanim %s, got %s", tc.tanim, proje.Tanim)
			}
			if proje.ID == "" {
				t.Error("ID should not be empty")
			}

			// Verify it was saved
			if _, ok := mockVY.projeler[proje.ID]; !ok {
				t.Error("project was not saved to database")
			}
		})
	}
}

func TestIsYonetici_GorevDuzenle(t *testing.T) {
	testCases := []struct {
		name             string
		gorevID          string
		baslik           string
		aciklama         string
		oncelik          string
		projeID          string
		baslikVar        bool
		aciklamaVar      bool
		oncelikVar       bool
		projeVar         bool
		shouldFailGetir  bool
		shouldFailUpdate bool
		wantErr          bool
	}{
		{
			name:      "update only title",
			gorevID:   "existing-task",
			baslik:    "Yeni Başlık",
			baslikVar: true,
			wantErr:   false,
		},
		{
			name:        "update only description",
			gorevID:     "existing-task",
			aciklama:    "Yeni Açıklama",
			aciklamaVar: true,
			wantErr:     false,
		},
		{
			name:        "update all fields",
			gorevID:     "existing-task",
			baslik:      "Yeni Başlık",
			aciklama:    "Yeni Açıklama",
			oncelik:     "yuksek",
			projeID:     "proje-1",
			baslikVar:   true,
			aciklamaVar: true,
			oncelikVar:  true,
			projeVar:    true,
			wantErr:     false,
		},
		{
			name:      "non-existing task",
			gorevID:   "non-existing",
			baslik:    "Test",
			baslikVar: true,
			wantErr:   true,
		},
		{
			name:            "database getir error",
			gorevID:         "existing-task",
			baslik:          "Test",
			baslikVar:       true,
			shouldFailGetir: true,
			wantErr:         true,
		},
		{
			name:             "database update error",
			gorevID:          "existing-task",
			baslik:           "Test",
			baslikVar:        true,
			shouldFailUpdate: true,
			wantErr:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVY := NewMockVeriYonetici()
			iy := YeniIsYonetici(mockVY)

			// Add test task
			originalTask := &Gorev{
				ID:       "existing-task",
				Baslik:   "Original Title",
				Aciklama: "Original Description",
				Durum:    "beklemede",
				Oncelik:  "orta",
				ProjeID:  "",
			}
			mockVY.gorevler["existing-task"] = originalTask

			mockVY.shouldFailGorevGetir = tc.shouldFailGetir
			mockVY.shouldFailGorevGuncelle = tc.shouldFailUpdate

			err := iy.GorevDuzenle(tc.gorevID, tc.baslik, tc.aciklama, tc.oncelik, tc.projeID,
				tc.baslikVar, tc.aciklamaVar, tc.oncelikVar, tc.projeVar)

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

			// Verify updates
			gorev := mockVY.gorevler["existing-task"]
			if tc.baslikVar && tc.baslik != "" {
				if gorev.Baslik != tc.baslik {
					t.Errorf("expected Baslik %s, got %s", tc.baslik, gorev.Baslik)
				}
			} else {
				if gorev.Baslik != originalTask.Baslik {
					t.Error("Baslik should not have changed")
				}
			}

			if tc.aciklamaVar {
				if gorev.Aciklama != tc.aciklama {
					t.Errorf("expected Aciklama %s, got %s", tc.aciklama, gorev.Aciklama)
				}
			} else {
				if gorev.Aciklama != originalTask.Aciklama {
					t.Error("Aciklama should not have changed")
				}
			}
		})
	}
}

func TestIsYonetici_GorevSil(t *testing.T) {
	testCases := []struct {
		name            string
		gorevID         string
		shouldFailGetir bool
		shouldFailSil   bool
		wantErr         bool
	}{
		{
			name:    "delete existing task",
			gorevID: "existing-task",
			wantErr: false,
		},
		{
			name:    "delete non-existing task",
			gorevID: "non-existing",
			wantErr: true,
		},
		{
			name:            "database getir error",
			gorevID:         "existing-task",
			shouldFailGetir: true,
			wantErr:         true,
		},
		{
			name:          "database sil error",
			gorevID:       "existing-task",
			shouldFailSil: true,
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVY := NewMockVeriYonetici()
			iy := YeniIsYonetici(mockVY)

			// Add test task
			if tc.gorevID == "existing-task" {
				mockVY.gorevler["existing-task"] = &Gorev{
					ID:     "existing-task",
					Baslik: "Test Task",
				}
			}

			mockVY.shouldFailGorevGetir = tc.shouldFailGetir
			mockVY.shouldFailGorevSil = tc.shouldFailSil

			err := iy.GorevSil(tc.gorevID)
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
			if _, ok := mockVY.gorevler["existing-task"]; ok {
				t.Error("task should have been deleted")
			}
		})
	}
}

func TestIsYonetici_OzetAl(t *testing.T) {
	testCases := []struct {
		name                     string
		gorevler                 []*Gorev
		projeler                 []*Proje
		shouldFailGorevleriGetir bool
		shouldFailProjeleriGetir bool
		wantErr                  bool
		expectedOzet             *Ozet
	}{
		{
			name: "calculate summary correctly",
			gorevler: []*Gorev{
				{ID: "1", Durum: "beklemede", Oncelik: "yuksek"},
				{ID: "2", Durum: "beklemede", Oncelik: "orta"},
				{ID: "3", Durum: "devam_ediyor", Oncelik: "orta"},
				{ID: "4", Durum: "tamamlandi", Oncelik: "dusuk"},
				{ID: "5", Durum: "tamamlandi", Oncelik: "yuksek"},
			},
			projeler: []*Proje{
				{ID: "p1", Isim: "Proje 1"},
				{ID: "p2", Isim: "Proje 2"},
			},
			expectedOzet: &Ozet{
				ToplamProje:     2,
				ToplamGorev:     5,
				BeklemedeGorev:  2,
				DevamEdenGorev:  1,
				TamamlananGorev: 2,
				YuksekOncelik:   2,
				OrtaOncelik:     2,
				DusukOncelik:    1,
			},
		},
		{
			name:                     "database gorevler error",
			shouldFailGorevleriGetir: true,
			wantErr:                  true,
		},
		{
			name:                     "database projeler error",
			gorevler:                 []*Gorev{},
			shouldFailProjeleriGetir: true,
			wantErr:                  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVY := NewMockVeriYonetici()
			iy := YeniIsYonetici(mockVY)

			// Add test data
			for _, g := range tc.gorevler {
				mockVY.gorevler[g.ID] = g
			}
			for _, p := range tc.projeler {
				mockVY.projeler[p.ID] = p
			}

			mockVY.shouldFailGorevleriGetir = tc.shouldFailGorevleriGetir
			mockVY.shouldFailProjeleriGetir = tc.shouldFailProjeleriGetir

			ozet, err := iy.OzetAl()
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

			// Verify summary
			if tc.expectedOzet != nil {
				if ozet.ToplamProje != tc.expectedOzet.ToplamProje {
					t.Errorf("expected ToplamProje %d, got %d", tc.expectedOzet.ToplamProje, ozet.ToplamProje)
				}
				if ozet.ToplamGorev != tc.expectedOzet.ToplamGorev {
					t.Errorf("expected ToplamGorev %d, got %d", tc.expectedOzet.ToplamGorev, ozet.ToplamGorev)
				}
				if ozet.BeklemedeGorev != tc.expectedOzet.BeklemedeGorev {
					t.Errorf("expected BeklemedeGorev %d, got %d", tc.expectedOzet.BeklemedeGorev, ozet.BeklemedeGorev)
				}
				if ozet.DevamEdenGorev != tc.expectedOzet.DevamEdenGorev {
					t.Errorf("expected DevamEdenGorev %d, got %d", tc.expectedOzet.DevamEdenGorev, ozet.DevamEdenGorev)
				}
				if ozet.TamamlananGorev != tc.expectedOzet.TamamlananGorev {
					t.Errorf("expected TamamlananGorev %d, got %d", tc.expectedOzet.TamamlananGorev, ozet.TamamlananGorev)
				}
				if ozet.YuksekOncelik != tc.expectedOzet.YuksekOncelik {
					t.Errorf("expected YuksekOncelik %d, got %d", tc.expectedOzet.YuksekOncelik, ozet.YuksekOncelik)
				}
				if ozet.OrtaOncelik != tc.expectedOzet.OrtaOncelik {
					t.Errorf("expected OrtaOncelik %d, got %d", tc.expectedOzet.OrtaOncelik, ozet.OrtaOncelik)
				}
				if ozet.DusukOncelik != tc.expectedOzet.DusukOncelik {
					t.Errorf("expected DusukOncelik %d, got %d", tc.expectedOzet.DusukOncelik, ozet.DusukOncelik)
				}
			}
		})
	}
}

func TestIsYonetici_ProjeGorevSayisi(t *testing.T) {
	mockVY := NewMockVeriYonetici()
	iy := YeniIsYonetici(mockVY)

	// Add test data
	mockVY.gorevler["1"] = &Gorev{ID: "1", ProjeID: "proje-1"}
	mockVY.gorevler["2"] = &Gorev{ID: "2", ProjeID: "proje-1"}
	mockVY.gorevler["3"] = &Gorev{ID: "3", ProjeID: "proje-2"}
	mockVY.gorevler["4"] = &Gorev{ID: "4", ProjeID: ""}

	testCases := []struct {
		name          string
		projeID       string
		expectedCount int
	}{
		{
			name:          "project with tasks",
			projeID:       "proje-1",
			expectedCount: 2,
		},
		{
			name:          "project with one task",
			projeID:       "proje-2",
			expectedCount: 1,
		},
		{
			name:          "project with no tasks",
			projeID:       "proje-3",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count, err := iy.ProjeGorevSayisi(tc.projeID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if count != tc.expectedCount {
				t.Errorf("expected count %d, got %d", tc.expectedCount, count)
			}
		})
	}
}
