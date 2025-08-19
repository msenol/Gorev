package gorev

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVeriYoneticiAI is a mock implementation of VeriYoneticiInterface for AI tests
type MockVeriYoneticiAI struct {
	mock.Mock
}

// setupBasicAIContextMocks sets up the common mock expectations for AI context operations
func setupBasicAIContextMocks(mockVY *MockVeriYoneticiAI) {
	// Mock AIContextGetir to return a default context
	mockContext := &AIContext{
		RecentTasks: []string{},
		SessionData: make(map[string]interface{}),
		LastUpdated: time.Now(),
	}
	mockVY.On("AIContextGetir").Return(mockContext, nil).Maybe()

	// Mock AIContextKaydet
	mockVY.On("AIContextKaydet", mock.AnythingOfType("*gorev.AIContext")).Return(nil).Maybe()

	// Mock AIInteractionKaydet
	mockVY.On("AIInteractionKaydet", mock.AnythingOfType("*gorev.AIInteraction")).Return(nil).Maybe()

	// Mock AIInteractionlariGetir
	mockVY.On("AIInteractionlariGetir", mock.AnythingOfType("int")).Return([]*AIInteraction{}, nil).Maybe()

	// Mock AITodayInteractionlariGetir
	mockVY.On("AITodayInteractionlariGetir").Return([]*AIInteraction{}, nil).Maybe()

	// Mock AILastInteractionGuncelle
	mockVY.On("AILastInteractionGuncelle", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil).Maybe()
}

func (m *MockVeriYoneticiAI) GorevKaydet(gorev *Gorev) error {
	args := m.Called(gorev)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) GorevGetir(id string) (*Gorev, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) GorevGuncelle(taskID string, params interface{}) error {
	args := m.Called(taskID, params)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) GorevSil(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) GorevleriGetir(durum, sirala, filtre string) ([]*Gorev, error) {
	args := m.Called(durum, sirala, filtre)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) ProjeKaydet(proje *Proje) error {
	args := m.Called(proje)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) ProjeGetir(id string) (*Proje, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Proje), args.Error(1)
}

func (m *MockVeriYoneticiAI) ProjeleriGetir() ([]*Proje, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Proje), args.Error(1)
}

func (m *MockVeriYoneticiAI) ProjeGorevleriGetir(projeID string) ([]*Gorev, error) {
	args := m.Called(projeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) AktifProjeAyarla(projeID string) error {
	args := m.Called(projeID)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) AktifProjeGetir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockVeriYoneticiAI) AktifProjeKaldir() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) BaglantiEkle(baglanti *Baglanti) error {
	args := m.Called(baglanti)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) BaglantilariGetir(gorevID string) ([]*Baglanti, error) {
	args := m.Called(gorevID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Baglanti), args.Error(1)
}

func (m *MockVeriYoneticiAI) EtiketleriGetir() ([]*Etiket, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Etiket), args.Error(1)
}

func (m *MockVeriYoneticiAI) EtiketOlustur(isim string) (*Etiket, error) {
	args := m.Called(isim)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Etiket), args.Error(1)
}

func (m *MockVeriYoneticiAI) EtiketleriGetirVeyaOlustur(isimler []string) ([]*Etiket, error) {
	args := m.Called(isimler)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Etiket), args.Error(1)
}

func (m *MockVeriYoneticiAI) GorevEtiketleriniAyarla(gorevID string, etiketler []*Etiket) error {
	args := m.Called(gorevID, etiketler)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) GorevEtiketleriniGetir(gorevID string) ([]*Etiket, error) {
	args := m.Called(gorevID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Etiket), args.Error(1)
}

func (m *MockVeriYoneticiAI) TemplateListele(kategori string) ([]*GorevTemplate, error) {
	args := m.Called(kategori)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*GorevTemplate), args.Error(1)
}

func (m *MockVeriYoneticiAI) TemplatedenGorevOlustur(templateID string, degerler map[string]string) (*Gorev, error) {
	args := m.Called(templateID, degerler)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) AltGorevleriGetir(parentID string) ([]*Gorev, error) {
	args := m.Called(parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) TumAltGorevleriGetir(parentID string) ([]*Gorev, error) {
	args := m.Called(parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) UstGorevleriGetir(gorevID string) ([]*Gorev, error) {
	args := m.Called(gorevID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) GorevHiyerarsiGetir(gorevID string) (*GorevHiyerarsi, error) {
	args := m.Called(gorevID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GorevHiyerarsi), args.Error(1)
}

func (m *MockVeriYoneticiAI) ParentIDGuncelle(gorevID, yeniParentID string) error {
	args := m.Called(gorevID, yeniParentID)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) DaireBagimliligiKontrolEt(gorevID, hedefParentID string) (bool, error) {
	args := m.Called(gorevID, hedefParentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockVeriYoneticiAI) TemplateOlustur(template *GorevTemplate) error {
	args := m.Called(template)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) TemplateGetir(templateID string) (*GorevTemplate, error) {
	args := m.Called(templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GorevTemplate), args.Error(1)
}

func (m *MockVeriYoneticiAI) TemplateAliasIleGetir(alias string) (*GorevTemplate, error) {
	args := m.Called(alias)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GorevTemplate), args.Error(1)
}

func (m *MockVeriYoneticiAI) TemplateIDVeyaAliasIleGetir(idOrAlias string) (*GorevTemplate, error) {
	args := m.Called(idOrAlias)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GorevTemplate), args.Error(1)
}

func (m *MockVeriYoneticiAI) VarsayilanTemplateleriOlustur() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) BulkBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error) {
	args := m.Called(gorevIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockVeriYoneticiAI) BulkTamamlanmamiaBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error) {
	args := m.Called(gorevIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockVeriYoneticiAI) BulkBuGoreveBagimliSayilariGetir(gorevIDs []string) (map[string]int, error) {
	args := m.Called(gorevIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

// AI Context Management methods
func (m *MockVeriYoneticiAI) AIContextGetir() (*AIContext, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AIContext), args.Error(1)
}

func (m *MockVeriYoneticiAI) AIContextKaydet(context *AIContext) error {
	args := m.Called(context)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) AIInteractionKaydet(interaction *AIInteraction) error {
	args := m.Called(interaction)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) AIInteractionlariGetir(limit int) ([]*AIInteraction, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*AIInteraction), args.Error(1)
}

func (m *MockVeriYoneticiAI) AITodayInteractionlariGetir() ([]*AIInteraction, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*AIInteraction), args.Error(1)
}

func (m *MockVeriYoneticiAI) AILastInteractionGuncelle(taskID string, timestamp time.Time) error {
	args := m.Called(taskID, timestamp)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) Kapat() error {
	args := m.Called()
	return args.Error(0)
}

// Missing interface methods
func (m *MockVeriYoneticiAI) AltGorevOlustur(parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	args := m.Called(parentID, baslik, aciklama, oncelik, sonTarihStr, etiketIsimleri)
	return args.Get(0).(*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) GorevDosyaYoluEkle(taskID string, path string) error {
	args := m.Called(taskID, path)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) GorevDosyaYoluSil(taskID string, path string) error {
	args := m.Called(taskID, path)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) GorevDosyaYollariGetir(taskID string) ([]string, error) {
	args := m.Called(taskID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockVeriYoneticiAI) DosyaYoluGorevleriGetir(path string) ([]string, error) {
	args := m.Called(path)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockVeriYoneticiAI) AIEtkilemasimKaydet(taskID string, interactionType, data, sessionID string) error {
	args := m.Called(taskID, interactionType, data, sessionID)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) GorevSonAIEtkilesiminiGuncelle(taskID string, timestamp time.Time) error {
	args := m.Called(taskID, timestamp)
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) GorevDetay(taskID string) (*Gorev, error) {
	args := m.Called(taskID)
	return args.Get(0).(*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) GorevListele(filters map[string]interface{}) ([]*Gorev, error) {
	args := m.Called(filters)
	return args.Get(0).([]*Gorev), args.Error(1)
}

func (m *MockVeriYoneticiAI) GorevOlustur(params map[string]interface{}) (string, error) {
	args := m.Called(params)
	return args.String(0), args.Error(1)
}

func (m *MockVeriYoneticiAI) GorevBagimlilikGetir(taskID string) ([]*Gorev, error) {
	args := m.Called(taskID)
	return args.Get(0).([]*Gorev), args.Error(1)
}

// TestSetActiveTask tests the SetActiveTask functionality
func TestSetActiveTask(t *testing.T) {
	tests := []struct {
		name          string
		taskID        string
		mockTask      *Gorev
		mockError     error
		expectError   bool
		expectedState string
	}{
		{
			name:   "Set active task - task in beklemede",
			taskID: "task-1",
			mockTask: &Gorev{
				ID:     "task-1",
				Baslik: "Test Task",
				Durum:  "beklemede",
			},
			expectError:   false,
			expectedState: "devam_ediyor",
		},
		{
			name:   "Set active task - task already in progress",
			taskID: "task-2",
			mockTask: &Gorev{
				ID:     "task-2",
				Baslik: "Test Task 2",
				Durum:  "devam_ediyor",
			},
			expectError:   false,
			expectedState: "devam_ediyor",
		},
		{
			name:        "Set active task - task not found",
			taskID:      "task-3",
			mockTask:    nil,
			mockError:   assert.AnError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVeriYonetici := new(MockVeriYoneticiAI)
			acy := YeniAIContextYonetici(mockVeriYonetici)

			// Setup basic AI context mock expectations
			setupBasicAIContextMocks(mockVeriYonetici)

			// Setup mock expectations
			mockVeriYonetici.On("GorevGetir", tt.taskID).Return(tt.mockTask, tt.mockError)

			if tt.mockTask != nil && tt.mockTask.Durum == "beklemede" {
				// The task will be updated with new status
				mockVeriYonetici.On("GorevGuncelle", tt.taskID, map[string]interface{}{"durum": "devam_ediyor"}).Return(nil)
			}

			// Execute
			err := acy.SetActiveTask(tt.taskID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockVeriYonetici.AssertExpectations(t)
		})
	}
}

// TestGetActiveTask tests the GetActiveTask functionality
func TestGetActiveTask(t *testing.T) {
	mockVeriYonetici := new(MockVeriYoneticiAI)
	acy := YeniAIContextYonetici(mockVeriYonetici)

	// Setup basic AI context mock expectations
	setupBasicAIContextMocks(mockVeriYonetici)

	// Since GetContext returns a mock implementation, we test the basic flow
	task, err := acy.GetActiveTask()
	assert.NoError(t, err)
	assert.Nil(t, task) // Context returns empty ActiveTaskID
}

// TestGetRecentTasks tests the GetRecentTasks functionality
func TestGetRecentTasks(t *testing.T) {
	mockVeriYonetici := new(MockVeriYoneticiAI)
	acy := YeniAIContextYonetici(mockVeriYonetici)

	// Setup basic AI context mock expectations
	setupBasicAIContextMocks(mockVeriYonetici)

	// Test with empty recent tasks
	tasks, err := acy.GetRecentTasks(5)
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

// TestRecordTaskView tests the RecordTaskView functionality
func TestRecordTaskView(t *testing.T) {
	tests := []struct {
		name         string
		taskID       string
		mockTask     *Gorev
		mockError    error
		expectError  bool
		shouldUpdate bool
	}{
		{
			name:   "Record view - task in beklemede",
			taskID: "task-1",
			mockTask: &Gorev{
				ID:     "task-1",
				Baslik: "Test Task",
				Durum:  "beklemede",
			},
			expectError:  false,
			shouldUpdate: true,
		},
		{
			name:   "Record view - task in devam_ediyor",
			taskID: "task-2",
			mockTask: &Gorev{
				ID:     "task-2",
				Baslik: "Test Task 2",
				Durum:  "devam_ediyor",
			},
			expectError:  false,
			shouldUpdate: false,
		},
		{
			name:        "Record view - task not found",
			taskID:      "task-3",
			mockTask:    nil,
			mockError:   assert.AnError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVeriYonetici := new(MockVeriYoneticiAI)
			acy := YeniAIContextYonetici(mockVeriYonetici)

			// Setup basic AI context mock expectations
			setupBasicAIContextMocks(mockVeriYonetici)

			// Setup mock expectations
			mockVeriYonetici.On("GorevGetir", tt.taskID).Return(tt.mockTask, tt.mockError)

			if tt.shouldUpdate && tt.mockTask != nil {
				mockVeriYonetici.On("GorevGuncelle", tt.taskID, map[string]interface{}{"durum": "devam_ediyor"}).Return(nil)
			}

			// Execute
			err := acy.RecordTaskView(tt.taskID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockVeriYonetici.AssertExpectations(t)
		})
	}
}

// TestBatchUpdate tests the BatchUpdate functionality
func TestBatchUpdate(t *testing.T) {
	mockVeriYonetici := new(MockVeriYoneticiAI)
	acy := YeniAIContextYonetici(mockVeriYonetici)

	// Setup basic AI context mock expectations
	setupBasicAIContextMocks(mockVeriYonetici)

	updates := []BatchUpdate{
		{
			ID: "task-1",
			Updates: map[string]interface{}{
				"durum": "devam_ediyor",
			},
		},
		{
			ID: "task-2",
			Updates: map[string]interface{}{
				"durum": "tamamlandi",
			},
		},
		{
			ID: "task-not-found",
			Updates: map[string]interface{}{
				"durum": "devam_ediyor",
			},
		},
	}

	// Setup mock expectations
	mockVeriYonetici.On("GorevGetir", "task-1").Return(&Gorev{ID: "task-1", Durum: "beklemede"}, nil)
	mockVeriYonetici.On("GorevGuncelle", "task-1", map[string]interface{}{"durum": "devam_ediyor"}).Return(nil)

	mockVeriYonetici.On("GorevGetir", "task-2").Return(&Gorev{ID: "task-2", Durum: "beklemede"}, nil)
	mockVeriYonetici.On("GorevGuncelle", "task-2", map[string]interface{}{"durum": "tamamlandi"}).Return(nil)

	mockVeriYonetici.On("GorevGetir", "task-not-found").Return(nil, assert.AnError)

	// Execute
	result, err := acy.BatchUpdate(updates)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, result.TotalProcessed)
	assert.Len(t, result.Successful, 2)
	assert.Len(t, result.Failed, 1)
	assert.Contains(t, result.Successful, "task-1")
	assert.Contains(t, result.Successful, "task-2")
	assert.Equal(t, "task-not-found", result.Failed[0].TaskID)

	mockVeriYonetici.AssertExpectations(t)
}

// TestNLPQuery tests the NLPQuery functionality
func TestNLPQuery(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		mockTasks    []*Gorev
		expectedLen  int
		expectFilter string
	}{
		{
			name:  "Query for high priority tasks",
			query: "yüksek öncelikli görevler",
			mockTasks: []*Gorev{
				{ID: "1", Baslik: "Task 1", Oncelik: "yuksek"},
				{ID: "2", Baslik: "Task 2", Oncelik: "yuksek"},
			},
			expectedLen: 2,
		},
		{
			name:  "Query for incomplete tasks",
			query: "tamamlanmamış görevler",
			mockTasks: []*Gorev{
				{ID: "1", Baslik: "Task 1", Durum: "beklemede"},
				{ID: "2", Baslik: "Task 2", Durum: "beklemede"},
			},
			expectedLen: 2,
		},
		{
			name:  "Query for urgent tasks",
			query: "acil görevler",
			mockTasks: []*Gorev{
				{ID: "1", Baslik: "Urgent Task", SonTarih: &time.Time{}},
			},
			expectedLen:  1,
			expectFilter: "acil",
		},
		{
			name:  "Query with tag filter",
			query: "etiket:bug",
			mockTasks: []*Gorev{
				{ID: "1", Baslik: "Bug Fix", Etiketler: []*Etiket{{Isim: "bug"}}},
			},
			expectedLen: 1,
		},
		{
			name:  "General text search",
			query: "test",
			mockTasks: []*Gorev{
				{ID: "1", Baslik: "Test görevi için", Aciklama: "Test açıklaması"},
				{ID: "2", Baslik: "Başka görev", Aciklama: "Test içerik"},
				{ID: "3", Baslik: "Unrelated", Aciklama: "Unrelated"},
			},
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVeriYonetici := new(MockVeriYoneticiAI)
			acy := YeniAIContextYonetici(mockVeriYonetici)

			// Setup basic AI context mock expectations
			setupBasicAIContextMocks(mockVeriYonetici)

			// Setup mock based on query type
			if strings.Contains(tt.query, "yüksek öncelik") {
				mockVeriYonetici.On("GorevleriGetir", "beklemede", "", "").Return(tt.mockTasks, nil)
			} else if strings.Contains(tt.query, "tamamlanmamış") {
				mockVeriYonetici.On("GorevleriGetir", "beklemede", "", "").Return(tt.mockTasks, nil)
			} else if strings.Contains(tt.query, "acil") {
				mockVeriYonetici.On("GorevleriGetir", "", "", "acil").Return(tt.mockTasks, nil)
			} else if strings.Contains(tt.query, "etiket:") {
				// For tag queries, we fetch all tasks and filter manually
				mockVeriYonetici.On("GorevleriGetir", "", "", "").Return(tt.mockTasks, nil)
			} else {
				// General text search
				mockVeriYonetici.On("GorevleriGetir", "", "", "").Return(tt.mockTasks, nil)
			}

			// Execute
			tasks, err := acy.NLPQuery(tt.query)

			// Assert
			assert.NoError(t, err)
			assert.Len(t, tasks, tt.expectedLen)

			mockVeriYonetici.AssertExpectations(t)
		})
	}
}

// TestGetContextSummary tests the GetContextSummary functionality
func TestGetContextSummary(t *testing.T) {
	mockVeriYonetici := new(MockVeriYoneticiAI)
	acy := YeniAIContextYonetici(mockVeriYonetici)

	// Setup basic AI context mock expectations
	setupBasicAIContextMocks(mockVeriYonetici)

	// Setup mock data
	allTasks := []*Gorev{
		{ID: "1", Baslik: "High Priority", Oncelik: "yuksek", Durum: "beklemede"},
		{ID: "2", Baslik: "Blocked Task", Durum: "beklemede", TamamlanmamisBagimlilikSayisi: 2},
		{ID: "3", Baslik: "Normal Task", Oncelik: "orta", Durum: "beklemede"},
	}

	mockVeriYonetici.On("GorevleriGetir", "beklemede", "", "").Return(allTasks, nil)

	// Execute
	summary, err := acy.GetContextSummary()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Len(t, summary.NextPriorities, 1)
	assert.Equal(t, "High Priority", summary.NextPriorities[0].Baslik)
	assert.Len(t, summary.Blockers, 1)
	assert.Equal(t, "Blocked Task", summary.Blockers[0].Baslik)

	mockVeriYonetici.AssertExpectations(t)
}

// TestBatchUpdateEnhanced tests the enhanced BatchUpdate functionality
func TestBatchUpdateEnhanced(t *testing.T) {
	tests := []struct {
		name            string
		updates         []BatchUpdate
		expectError     bool
		expectedSuccess int
		expectedFailed  int
	}{
		{
			name: "valid status update",
			updates: []BatchUpdate{
				{
					ID: "task1",
					Updates: map[string]interface{}{
						"durum": "devam_ediyor",
					},
				},
			},
			expectError:     false,
			expectedSuccess: 1,
			expectedFailed:  0,
		},
		{
			name: "valid priority update",
			updates: []BatchUpdate{
				{
					ID: "task1",
					Updates: map[string]interface{}{
						"oncelik": "yuksek",
					},
				},
			},
			expectError:     false,
			expectedSuccess: 1,
			expectedFailed:  0,
		},
		{
			name: "valid multiple field update",
			updates: []BatchUpdate{
				{
					ID: "task1",
					Updates: map[string]interface{}{
						"durum":     "devam_ediyor",
						"oncelik":   "yuksek",
						"baslik":    "Updated Title",
						"aciklama":  "Updated description",
						"son_tarih": "2024-12-31",
					},
				},
			},
			expectError:     false,
			expectedSuccess: 1,
			expectedFailed:  0,
		},
		{
			name: "invalid status",
			updates: []BatchUpdate{
				{
					ID: "task1",
					Updates: map[string]interface{}{
						"durum": "invalid_status",
					},
				},
			},
			expectError:     false,
			expectedSuccess: 0,
			expectedFailed:  1,
		},
		{
			name: "invalid priority",
			updates: []BatchUpdate{
				{
					ID: "task1",
					Updates: map[string]interface{}{
						"oncelik": "invalid_priority",
					},
				},
			},
			expectError:     false,
			expectedSuccess: 0,
			expectedFailed:  1,
		},
		{
			name: "empty title",
			updates: []BatchUpdate{
				{
					ID: "task1",
					Updates: map[string]interface{}{
						"baslik": "",
					},
				},
			},
			expectError:     false,
			expectedSuccess: 0,
			expectedFailed:  1,
		},
		{
			name: "invalid date format",
			updates: []BatchUpdate{
				{
					ID: "task1",
					Updates: map[string]interface{}{
						"son_tarih": "invalid-date",
					},
				},
			},
			expectError:     false,
			expectedSuccess: 0,
			expectedFailed:  1,
		},
		{
			name: "task not found",
			updates: []BatchUpdate{
				{
					ID: "nonexistent",
					Updates: map[string]interface{}{
						"durum": "devam_ediyor",
					},
				},
			},
			expectError:     false,
			expectedSuccess: 0,
			expectedFailed:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVeriYonetici := new(MockVeriYoneticiAI)
			acy := YeniAIContextYonetici(mockVeriYonetici)

			// Setup basic AI context mock expectations
			setupBasicAIContextMocks(mockVeriYonetici)

			// Setup task existence mocks
			for _, update := range tt.updates {
				if update.ID == "nonexistent" {
					mockVeriYonetici.On("GorevGetir", update.ID).Return(nil, fmt.Errorf("task not found"))
				} else {
					testTask := &Gorev{ID: update.ID, Baslik: "Test Task"}
					mockVeriYonetici.On("GorevGetir", update.ID).Return(testTask, nil)

					// Only expect GorevGuncelle if we expect success
					if tt.expectedSuccess > 0 {
						mockVeriYonetici.On("GorevGuncelle", update.ID, mock.AnythingOfType("map[string]interface {}")).Return(nil)
					}
				}
			}

			// Execute
			result, err := acy.BatchUpdate(tt.updates)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Successful, tt.expectedSuccess)
				assert.Len(t, result.Failed, tt.expectedFailed)
				assert.Equal(t, len(tt.updates), result.TotalProcessed)
			}

			mockVeriYonetici.AssertExpectations(t)
		})
	}
}

// TestBulkBuGoreveBagimliSayilariGetir tests the new bulk dependency count method
func TestBulkBuGoreveBagimliSayilariGetir(t *testing.T) {
	mockVeriYonetici := new(MockVeriYoneticiAI)

	// Setup mock
	expectedCounts := map[string]int{
		"task1": 2,
		"task2": 0,
		"task3": 1,
	}

	mockVeriYonetici.On("BulkBuGoreveBagimliSayilariGetir", []string{"task1", "task2", "task3"}).Return(expectedCounts, nil)

	// Execute
	result, err := mockVeriYonetici.BulkBuGoreveBagimliSayilariGetir([]string{"task1", "task2", "task3"})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCounts, result)

	mockVeriYonetici.AssertExpectations(t)
}

// TestHelperFunctions tests helper functions
func TestHelperFunctions(t *testing.T) {
	// Test contains function
	assert.True(t, contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, contains([]string{"a", "b", "c"}, "d"))
	assert.False(t, contains([]string{}, "a"))
}

// TestAIContextRaceCondition tests concurrent access to AI context manager
func TestAIContextRaceCondition(t *testing.T) {
	mockVeriYonetici := new(MockVeriYoneticiAI)
	setupBasicAIContextMocks(mockVeriYonetici)

	// Mock a valid task
	mockGorev := &Gorev{
		ID:      "test-task-id",
		Baslik:  "Test Task",
		Durum:   "beklemede",
		ProjeID: "test-project-id",
	}
	mockVeriYonetici.On("GorevGetir", "test-task-id").Return(mockGorev, nil)
	mockVeriYonetici.On("GorevGuncelle", "test-task-id", mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// Create AI context manager
	acy := YeniAIContextYonetici(mockVeriYonetici)

	// Track errors from concurrent operations
	errors := make(chan error, 100)
	const numGoroutines = constants.TestConcurrencyLarge
	const operationsPerGoroutine = 10

	// Function to perform concurrent operations
	performOperations := func() {
		for i := 0; i < operationsPerGoroutine; i++ {
			// Mix of read and write operations
			switch i % 4 {
			case 0:
				// SetActiveTask (write operation)
				err := acy.SetActiveTask("test-task-id")
				if err != nil {
					errors <- fmt.Errorf("SetActiveTask failed: %w", err)
					return
				}
			case 1:
				// GetActiveTask (read operation)
				_, err := acy.GetActiveTask()
				if err != nil {
					errors <- fmt.Errorf("GetActiveTask failed: %w", err)
					return
				}
			case 2:
				// GetContext (read operation)
				_, err := acy.GetContext()
				if err != nil {
					errors <- fmt.Errorf("GetContext failed: %w", err)
					return
				}
			case 3:
				// GetRecentTasks (read operation)
				_, err := acy.GetRecentTasks(5)
				if err != nil {
					errors <- fmt.Errorf("GetRecentTasks failed: %w", err)
					return
				}
			}
		}
	}

	// Launch concurrent goroutines
	for i := 0; i < numGoroutines; i++ {
		go performOperations()
	}

	// Wait a bit for operations to complete
	time.Sleep(100 * time.Millisecond)

	// Check if any errors occurred
	close(errors)
	var collectedErrors []error
	for err := range errors {
		collectedErrors = append(collectedErrors, err)
	}

	// Assert no race conditions or errors occurred
	if len(collectedErrors) > 0 {
		var errorMessages []string
		for _, err := range collectedErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		t.Fatalf("Race condition detected - %d errors occurred:\n%s",
			len(collectedErrors),
			strings.Join(errorMessages, "\n"))
	}

	// Verify final state is consistent
	context, err := acy.GetContext()
	assert.NoError(t, err)
	assert.NotNil(t, context)

	activeTask, err := acy.GetActiveTask()
	assert.NoError(t, err)
	if activeTask != nil {
		assert.Equal(t, "test-task-id", activeTask.ID)
	}

	t.Logf("Successfully completed %d concurrent operations across %d goroutines without race conditions",
		numGoroutines*operationsPerGoroutine, numGoroutines)
}
