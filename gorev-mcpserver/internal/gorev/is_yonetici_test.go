package gorev

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// setupTestI18n initializes the i18n system for tests
func setupTestI18n() {
	// Initialize i18n with Turkish (default) for tests
	_ = i18n.Initialize(constants.DefaultTestLanguage)
}

// MockVeriYonetici is a mock implementation of VeriYonetici for testing
type MockVeriYonetici struct {
	gorevler      map[string]*Gorev
	projeler      map[string]*Proje
	baglantilar   []*Baglanti
	aktifProjeID  string
	aktifProjeSet bool

	// AI Context specific data for tests
	aiContext         *AIContext
	interactions      []*AIInteraction
	todayInteractions []*AIInteraction
	allTasks          []*Gorev
	tags              map[string]*Etiket

	// Control behavior
	shouldFailGorevKaydet    bool
	shouldFailGorevGetir     bool
	shouldFailGorevGuncelle  bool
	shouldFailGorevSil       bool
	shouldFailProjeKaydet    bool
	shouldFailProjeGetir     bool
	shouldFailGorevleriGetir bool
	shouldFailProjeleriGetir bool
	shouldFailGorevListele   bool
	shouldFailAktifProje     bool
	shouldFailBaglantiEkle   bool
	shouldReturnError        bool
	errorToReturn            error
}

func NewMockVeriYonetici() *MockVeriYonetici {
	return &MockVeriYonetici{
		gorevler:          make(map[string]*Gorev),
		projeler:          make(map[string]*Proje),
		baglantilar:       make([]*Baglanti, 0),
		tags:              make(map[string]*Etiket),
		aiContext:         &AIContext{RecentTasks: []string{}, SessionData: make(map[string]interface{})},
		interactions:      []*AIInteraction{},
		todayInteractions: []*AIInteraction{},
		shouldReturnError: false,
	}
}
func (m *MockVeriYonetici) AktifProjeAyarla(projeID string) error {
	if m.shouldFailAktifProje {
		return errors.New("mock error")
	}
	m.aktifProjeID = projeID
	m.aktifProjeSet = true
	return nil
}

func (m *MockVeriYonetici) AktifProjeGetir() (string, error) {
	if m.shouldFailAktifProje {
		return "", errors.New("mock error")
	}
	if !m.aktifProjeSet {
		return "", errors.New("aktif proje bulunamadı")
	}
	return m.aktifProjeID, nil
}

func (m *MockVeriYonetici) AktifProjeKaldir() error {
	if m.shouldFailAktifProje {
		return errors.New("mock error")
	}
	m.aktifProjeID = ""
	m.aktifProjeSet = false
	return nil
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
		return nil, errors.New(i18n.TEntityNotFound("task", errors.New("not found")))
	}
	return gorev, nil
}

func (m *MockVeriYonetici) GorevleriGetir(durum, sirala, filtre string) ([]*Gorev, error) {
	if m.shouldFailGorevleriGetir {
		return nil, errors.New("mock error: gorevleri getir failed")
	}
	if m.shouldReturnError {
		return nil, m.errorToReturn
	}
	
	// If allTasks is populated (for tests), use that
	if len(m.allTasks) > 0 {
		return m.allTasks, nil
	}
	
	// Otherwise, use the map
	var result []*Gorev
	for _, gorev := range m.gorevler {
		if durum == "" || gorev.Durum == durum {
			result = append(result, gorev)
		}
	}
	return result, nil
}

func (m *MockVeriYonetici) GorevGuncelle(taskID string, params interface{}) error {
	if m.shouldFailGorevGuncelle {
		return errors.New("mock error: gorev guncelle failed")
	}
	gorev, ok := m.gorevler[taskID]
	if !ok {
		return errors.New("görev bulunamadı")
	}
	
	// Apply updates from params map
	if updateParams, ok := params.(map[string]interface{}); ok {
		for key, value := range updateParams {
			switch key {
			case "baslik":
				if val, ok := value.(string); ok {
					gorev.Baslik = val
				}
			case "aciklama":
				if val, ok := value.(string); ok {
					gorev.Aciklama = val
				}
			case "durum":
				if val, ok := value.(string); ok {
					gorev.Durum = val
				}
			case "oncelik":
				if val, ok := value.(string); ok {
					gorev.Oncelik = val
				}
			case "proje_id":
				if val, ok := value.(string); ok {
					gorev.ProjeID = val
				}
			case "parent_id":
				if val, ok := value.(string); ok {
					gorev.ParentID = val
				}
			case "guncelleme_tarih":
				if val, ok := value.(time.Time); ok {
					gorev.GuncellemeTarih = val
				}
			}
		}
	}
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

func (m *MockVeriYonetici) BulkBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error) {
	// Simple mock implementation
	result := make(map[string]int)
	for _, id := range gorevIDs {
		// Count total dependencies - tasks that this task depends on
		count := 0
		for _, b := range m.baglantilar {
			if b.HedefID == id && b.BaglantiTip == "onceki" { // This task depends on the source task
				count++
			}
		}
		result[id] = count
	}
	return result, nil
}

func (m *MockVeriYonetici) BulkTamamlanmamiaBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error) {
	// Simple mock implementation
	result := make(map[string]int)
	for _, id := range gorevIDs {
		// Count uncompleted dependencies - tasks that this task depends on
		count := 0
		for _, b := range m.baglantilar {
			if b.HedefID == id && b.BaglantiTip == "onceki" { // This task depends on the source task
				if kaynakGorev, exists := m.gorevler[b.KaynakID]; exists {
					if kaynakGorev.Durum != "tamamlandi" {
						count++
					}
				}
			}
		}
		result[id] = count
	}
	return result, nil
}

func (m *MockVeriYonetici) Kapat() error {
	return nil
}

func (m *MockVeriYonetici) EtiketleriGetir() ([]*Etiket, error) {
	if m.shouldReturnError {
		return nil, m.errorToReturn
	}
	var result []*Etiket
	for _, tag := range m.tags {
		result = append(result, tag)
	}
	return result, nil
}

func (m *MockVeriYonetici) EtiketOlustur(isim string) (*Etiket, error) {
	if m.shouldReturnError {
		return nil, m.errorToReturn
	}
	tag := &Etiket{ID: "tag-" + isim, Isim: isim}
	m.tags[tag.ID] = tag
	return tag, nil
}

func (m *MockVeriYonetici) EtiketleriGetirVeyaOlustur(isimler []string) ([]*Etiket, error) {
	if m.shouldReturnError {
		return nil, m.errorToReturn
	}
	var result []*Etiket
	for _, isim := range isimler {
		// Try to find existing tag
		var found *Etiket
		for _, tag := range m.tags {
			if tag.Isim == isim {
				found = tag
				break
			}
		}
		if found != nil {
			result = append(result, found)
		} else {
			// Create new tag
			tag, _ := m.EtiketOlustur(isim)
			result = append(result, tag)
		}
	}
	return result, nil
}

func (m *MockVeriYonetici) GorevEtiketleriniGetir(gorevID string) ([]*Etiket, error) {
	if m.shouldReturnError {
		return nil, m.errorToReturn
	}
	gorev, exists := m.gorevler[gorevID]
	if !exists {
		return nil, errors.New("task not found")
	}
	return gorev.Etiketler, nil
}

func (m *MockVeriYonetici) GorevEtiketleriniAyarla(gorevID string, etiketler []*Etiket) error {
	if m.shouldReturnError {
		return m.errorToReturn
	}
	gorev, exists := m.gorevler[gorevID]
	if !exists {
		return errors.New("task not found")
	}
	gorev.Etiketler = etiketler
	return nil
}

func (m *MockVeriYonetici) BaglantiEkle(baglanti *Baglanti) error {
	if m.shouldFailBaglantiEkle {
		return errors.New("mock error: baglanti ekle failed")
	}
	m.baglantilar = append(m.baglantilar, baglanti)
	return nil
}

func (m *MockVeriYonetici) BaglantilariGetir(gorevID string) ([]*Baglanti, error) {
	var result []*Baglanti
	for _, b := range m.baglantilar {
		if b.KaynakID == gorevID || b.HedefID == gorevID {
			result = append(result, b)
		}
	}
	return result, nil
}

// Template mock methods
func (m *MockVeriYonetici) TemplateOlustur(template *GorevTemplate) error {
	return nil
}

func (m *MockVeriYonetici) TemplateListele(kategori string) ([]*GorevTemplate, error) {
	return []*GorevTemplate{}, nil
}

func (m *MockVeriYonetici) TemplateGetir(templateID string) (*GorevTemplate, error) {
	return &GorevTemplate{}, nil
}

func (m *MockVeriYonetici) TemplateAliasIleGetir(alias string) (*GorevTemplate, error) {
	return &GorevTemplate{}, nil
}

func (m *MockVeriYonetici) TemplateIDVeyaAliasIleGetir(idOrAlias string) (*GorevTemplate, error) {
	return &GorevTemplate{}, nil
}

func (m *MockVeriYonetici) TemplatedenGorevOlustur(templateID string, degerler map[string]string) (*Gorev, error) {
	return &Gorev{}, nil
}

func (m *MockVeriYonetici) VarsayilanTemplateleriOlustur() error {
	return nil
}

func (m *MockVeriYonetici) AltGorevleriGetir(parentID string) ([]*Gorev, error) {
	var result []*Gorev
	for _, gorev := range m.gorevler {
		if gorev.ParentID == parentID {
			result = append(result, gorev)
		}
	}
	return result, nil
}

func (m *MockVeriYonetici) TumAltGorevleriGetir(parentID string) ([]*Gorev, error) {
	// Simplified implementation for testing
	return m.AltGorevleriGetir(parentID)
}

func (m *MockVeriYonetici) UstGorevleriGetir(gorevID string) ([]*Gorev, error) {
	var result []*Gorev
	gorev, ok := m.gorevler[gorevID]
	if !ok || gorev.ParentID == "" {
		return result, nil
	}

	parent, ok := m.gorevler[gorev.ParentID]
	if ok {
		result = append(result, parent)
	}
	return result, nil
}

func (m *MockVeriYonetici) GorevHiyerarsiGetir(gorevID string) (*GorevHiyerarsi, error) {
	gorev, err := m.GorevGetir(gorevID)
	if err != nil {
		return nil, err
	}

	return &GorevHiyerarsi{
		Gorev:           gorev,
		UstGorevler:     []*Gorev{},
		ToplamAltGorev:  0,
		TamamlananAlt:   0,
		DevamEdenAlt:    0,
		BeklemedeAlt:    0,
		IlerlemeYuzdesi: 0,
	}, nil
}

func (m *MockVeriYonetici) ParentIDGuncelle(gorevID, yeniParentID string) error {
	gorev, ok := m.gorevler[gorevID]
	if !ok {
		return errors.New("görev bulunamadı")
	}
	gorev.ParentID = yeniParentID
	return nil
}

func (m *MockVeriYonetici) DaireBagimliligiKontrolEt(gorevID, hedefParentID string) (bool, error) {
	// Simple check for testing - just check if they're the same
	return gorevID == hedefParentID, nil
}

// AI Context Management methods
func (m *MockVeriYonetici) AIContextGetir() (*AIContext, error) {
	if m.shouldReturnError {
		return nil, m.errorToReturn
	}
	return m.aiContext, nil
}

func (m *MockVeriYonetici) AIContextKaydet(context *AIContext) error {
	if m.shouldReturnError {
		return m.errorToReturn
	}
	m.aiContext = context
	return nil
}

func (m *MockVeriYonetici) AIInteractionKaydet(interaction *AIInteraction) error {
	if m.shouldReturnError {
		return m.errorToReturn
	}
	m.interactions = append(m.interactions, interaction)
	return nil
}

func (m *MockVeriYonetici) AIInteractionlariGetir(limit int) ([]*AIInteraction, error) {
	if m.shouldReturnError {
		return nil, m.errorToReturn
	}
	if limit <= 0 || limit >= len(m.interactions) {
		return m.interactions, nil
	}
	return m.interactions[:limit], nil
}

func (m *MockVeriYonetici) AITodayInteractionlariGetir() ([]*AIInteraction, error) {
	if m.shouldReturnError {
		return nil, m.errorToReturn
	}
	return m.todayInteractions, nil
}

func (m *MockVeriYonetici) AILastInteractionGuncelle(taskID string, timestamp time.Time) error {
	return nil
}

func (m *MockVeriYonetici) AltGorevOlustur(parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	return nil, nil
}

func (m *MockVeriYonetici) GorevDosyaYoluEkle(gorevID, dosyaYolu string) error {
	return nil
}

func (m *MockVeriYonetici) GorevDosyaYoluSil(gorevID, dosyaYolu string) error {
	return nil
}

func (m *MockVeriYonetici) GorevDosyaYollariGetir(gorevID string) ([]string, error) {
	return nil, nil
}

func (m *MockVeriYonetici) DosyaYoluGorevleriGetir(dosyaYolu string) ([]string, error) {
	return nil, nil
}

func (m *MockVeriYonetici) AIEtkilemasimKaydet(taskID string, interactionType, data, sessionID string) error {
	return nil
}

func (m *MockVeriYonetici) GorevSonAIEtkilesiminiGuncelle(gorevID string, timestamp time.Time) error {
	return nil
}

func (m *MockVeriYonetici) GorevDetay(id string) (*Gorev, error) {
	return m.GorevGetir(id)
}

func (m *MockVeriYonetici) GorevListele(filters map[string]interface{}) ([]*Gorev, error) {
	if m.shouldFailGorevListele {
		return nil, errors.New("mock error: gorev listele failed")
	}

	var result []*Gorev
	durum := ""
	if v, ok := filters["durum"]; ok {
		if s, ok := v.(string); ok {
			durum = s
		}
	}

	for _, gorev := range m.gorevler {
		// Apply durum filter if specified
		if durum == "" || gorev.Durum == durum {
			result = append(result, gorev)
		}
	}
	return result, nil
}

func (m *MockVeriYonetici) GorevOlustur(params map[string]interface{}) (string, error) {
	return "test-id", nil
}

func (m *MockVeriYonetici) GorevBagimlilikGetir(gorevID string) ([]*Gorev, error) {
	return nil, nil
}


func (m *MockVeriYonetici) BulkBuGoreveBagimliSayilariGetir(gorevIDs []string) (map[string]int, error) {
	result := make(map[string]int)
	for _, id := range gorevIDs {
		// Count how many tasks depend on this task (this task as source)
		count := 0
		for _, b := range m.baglantilar {
			if b.KaynakID == id && b.BaglantiTip == "onceki" { // Other tasks depend on this task
				count++
			}
		}
		result[id] = count
	}
	return result, nil
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
		projeID          string
		shouldFailKaydet bool
		wantErr          bool
	}{
		{
			name:     "valid task creation",
			baslik:   "Test Görevi",
			aciklama: "Test açıklaması",
			oncelik:  "orta",
			projeID:  "proje-1",
			wantErr:  false,
		},
		{
			name:     "empty title",
			baslik:   "",
			aciklama: "Açıklama",
			oncelik:  "yuksek",
			projeID:  "",
			wantErr:  false, // Business logic doesn't validate empty titles
		},
		{
			name:             "database error",
			baslik:           "Test",
			aciklama:         "Test",
			oncelik:          "orta",
			projeID:          "",
			shouldFailKaydet: true,
			wantErr:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockVY := NewMockVeriYonetici()
			mockVY.shouldFailGorevKaydet = tc.shouldFailKaydet
			iy := YeniIsYonetici(mockVY)

			gorev, err := iy.GorevOlustur(tc.baslik, tc.aciklama, tc.oncelik, tc.projeID, "", nil)
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
			mockVY.shouldFailGorevListele = tc.shouldFail

			gorevler, err := iy.GorevListele(map[string]interface{}{"durum": tc.durum})
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
	setupTestI18n() // Initialize i18n for tests
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
			yeniDurum: "devam_ediyor",
			wantErr:   false,
		},
		{
			name:          "non-existing task",
			gorevID:       "non-existing",
			yeniDurum:     "tamamlandi",
			wantErr:       true,
			expectedError: "not found", // Will be compared after i18n initialization
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
				} else if tc.expectedError != "" {
					expectedTranslated := i18n.T(tc.expectedError)
					if !strings.Contains(err.Error(), expectedTranslated) {
						t.Errorf("expected error containing '%s', got '%s'", expectedTranslated, err.Error())
					}
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

			err := iy.GorevDuzenle(tc.gorevID, tc.baslik, tc.aciklama, tc.oncelik, tc.projeID, "",
				tc.baslikVar, tc.aciklamaVar, tc.oncelikVar, tc.projeVar, false)

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

func TestIsYonetici_GorevBagimliMi(t *testing.T) {
	mockVY := NewMockVeriYonetici()
	iy := YeniIsYonetici(mockVY)

	// Test görevleri ekle
	mockVY.gorevler["gorev1"] = &Gorev{ID: "gorev1", Baslik: "Görev 1", Durum: "tamamlandi"}
	mockVY.gorevler["gorev2"] = &Gorev{ID: "gorev2", Baslik: "Görev 2", Durum: "beklemede"}
	mockVY.gorevler["gorev3"] = &Gorev{ID: "gorev3", Baslik: "Görev 3", Durum: "devam_ediyor"}
	mockVY.gorevler["gorev4"] = &Gorev{ID: "gorev4", Baslik: "Görev 4", Durum: "beklemede"}

	// Bağımlılıklar: gorev4, gorev1 ve gorev2'ye bağımlı
	mockVY.baglantilar = []*Baglanti{
		{ID: "b1", KaynakID: "gorev1", HedefID: "gorev4", BaglantiTip: "onceki"},
		{ID: "b2", KaynakID: "gorev2", HedefID: "gorev4", BaglantiTip: "onceki"},
	}

	testCases := []struct {
		name                  string
		gorevID               string
		expectedBagimli       bool
		expectedTamamlanmamis []string
	}{
		{
			name:                  "no dependencies",
			gorevID:               "gorev1",
			expectedBagimli:       true,
			expectedTamamlanmamis: nil,
		},
		{
			name:                  "all dependencies completed",
			gorevID:               "gorev3",
			expectedBagimli:       true,
			expectedTamamlanmamis: nil,
		},
		{
			name:                  "some dependencies not completed",
			gorevID:               "gorev4",
			expectedBagimli:       false,
			expectedTamamlanmamis: []string{"Görev 2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bagimli, tamamlanmamislar, err := iy.GorevBagimliMi(tc.gorevID)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if bagimli != tc.expectedBagimli {
				t.Errorf("expected bagimli=%v, got %v", tc.expectedBagimli, bagimli)
			}

			if len(tamamlanmamislar) != len(tc.expectedTamamlanmamis) {
				t.Errorf("expected %d tamamlanmamis, got %d", len(tc.expectedTamamlanmamis), len(tamamlanmamislar))
			}

			for i, expected := range tc.expectedTamamlanmamis {
				if i < len(tamamlanmamislar) && tamamlanmamislar[i] != expected {
					t.Errorf("expected tamamlanmamis[%d]=%s, got %s", i, expected, tamamlanmamislar[i])
				}
			}
		})
	}
}

func TestIsYonetici_GorevDurumGuncelle_WithDependencies(t *testing.T) {
	setupTestI18n() // Initialize i18n for tests
	mockVY := NewMockVeriYonetici()
	iy := YeniIsYonetici(mockVY)

	// Test görevleri ekle
	mockVY.gorevler["gorev1"] = &Gorev{ID: "gorev1", Baslik: "Görev 1", Durum: "beklemede"}
	mockVY.gorevler["gorev2"] = &Gorev{ID: "gorev2", Baslik: "Görev 2", Durum: "beklemede"}

	// gorev2, gorev1'e bağımlı
	mockVY.baglantilar = []*Baglanti{
		{ID: "b1", KaynakID: "gorev1", HedefID: "gorev2", BaglantiTip: "onceki"},
	}

	// gorev2'yi devam_ediyor yapmaya çalış (gorev1 henüz tamamlanmadı)
	err := iy.GorevDurumGuncelle("gorev2", "devam_ediyor")
	if err == nil {
		t.Error("expected error when trying to start task with incomplete dependencies")
	}
	// Check for the i18n key (translation may not be loaded in test environment)
	if !strings.Contains(err.Error(), "taskCannotStartDependencies") && !strings.Contains(err.Error(), "bu görev başlatılamaz") {
		t.Errorf("unexpected error message: %v", err)
	}

	// gorev1'i tamamla
	err = iy.GorevDurumGuncelle("gorev1", "tamamlandi")
	if err != nil {
		t.Errorf("unexpected error completing gorev1: %v", err)
	}

	// Şimdi gorev2'yi başlatabilmeli
	err = iy.GorevDurumGuncelle("gorev2", "devam_ediyor")
	if err != nil {
		t.Errorf("unexpected error starting gorev2 after dependencies completed: %v", err)
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

func TestIsYonetici_GorevListele_WithDependencyCounts(t *testing.T) {
	mockVeri := NewMockVeriYonetici()
	iy := YeniIsYonetici(mockVeri)

	// Create test tasks
	task1 := &Gorev{ID: "task1", Baslik: "Task 1", Durum: "beklemede", Oncelik: "orta"}
	task2 := &Gorev{ID: "task2", Baslik: "Task 2", Durum: "tamamlandi", Oncelik: "orta"}
	task3 := &Gorev{ID: "task3", Baslik: "Task 3", Durum: "beklemede", Oncelik: "orta"}

	mockVeri.gorevler["task1"] = task1
	mockVeri.gorevler["task2"] = task2
	mockVeri.gorevler["task3"] = task3

	// Create dependencies
	// task1 depends on task2 (completed)
	mockVeri.baglantilar = append(mockVeri.baglantilar, &Baglanti{
		ID:          "dep1",
		KaynakID:    "task2",
		HedefID:     "task1",
		BaglantiTip: "onceki",
	})

	// task3 depends on task1 (not completed)
	mockVeri.baglantilar = append(mockVeri.baglantilar, &Baglanti{
		ID:          "dep2",
		KaynakID:    "task1",
		HedefID:     "task3",
		BaglantiTip: "onceki",
	})

	// Get tasks with dependency counts
	gorevler, err := iy.GorevListele(map[string]interface{}{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify dependency counts
	for _, gorev := range gorevler {
		switch gorev.ID {
		case "task1":
			if gorev.BagimliGorevSayisi != 1 {
				t.Errorf("task1: expected BagimliGorevSayisi 1, got %d", gorev.BagimliGorevSayisi)
			}
			if gorev.TamamlanmamisBagimlilikSayisi != 0 {
				t.Errorf("task1: expected TamamlanmamisBagimlilikSayisi 0, got %d", gorev.TamamlanmamisBagimlilikSayisi)
			}
			if gorev.BuGoreveBagimliSayisi != 1 {
				t.Errorf("task1: expected BuGoreveBagimliSayisi 1, got %d", gorev.BuGoreveBagimliSayisi)
			}
		case "task2":
			if gorev.BagimliGorevSayisi != 0 {
				t.Errorf("task2: expected BagimliGorevSayisi 0, got %d", gorev.BagimliGorevSayisi)
			}
			if gorev.BuGoreveBagimliSayisi != 1 {
				t.Errorf("task2: expected BuGoreveBagimliSayisi 1, got %d", gorev.BuGoreveBagimliSayisi)
			}
		case "task3":
			if gorev.BagimliGorevSayisi != 1 {
				t.Errorf("task3: expected BagimliGorevSayisi 1, got %d", gorev.BagimliGorevSayisi)
			}
			if gorev.TamamlanmamisBagimlilikSayisi != 1 {
				t.Errorf("task3: expected TamamlanmamisBagimlilikSayisi 1, got %d", gorev.TamamlanmamisBagimlilikSayisi)
			}
			if gorev.BuGoreveBagimliSayisi != 0 {
				t.Errorf("task3: expected BuGoreveBagimliSayisi 0, got %d", gorev.BuGoreveBagimliSayisi)
			}
		}
	}
}
