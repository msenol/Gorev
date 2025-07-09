package gorev

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVeriYoneticiAI is a mock implementation of VeriYoneticiInterface for AI tests
type MockVeriYoneticiAI struct {
	mock.Mock
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

func (m *MockVeriYoneticiAI) GorevGuncelle(gorev *Gorev) error {
	args := m.Called(gorev)
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

func (m *MockVeriYoneticiAI) VarsayilanTemplateleriOlustur() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockVeriYoneticiAI) Kapat() error {
	args := m.Called()
	return args.Error(0)
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

			// Setup mock expectations
			mockVeriYonetici.On("GorevGetir", tt.taskID).Return(tt.mockTask, tt.mockError)

			if tt.mockTask != nil && tt.mockTask.Durum == "beklemede" {
				// The task will be updated with new status
				mockVeriYonetici.On("GorevGuncelle", mock.MatchedBy(func(g *Gorev) bool {
					return g.ID == tt.taskID && g.Durum == "devam_ediyor"
				})).Return(nil)
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

	// Since GetContext returns a mock implementation, we test the basic flow
	task, err := acy.GetActiveTask()
	assert.NoError(t, err)
	assert.Nil(t, task) // Context returns empty ActiveTaskID
}

// TestGetRecentTasks tests the GetRecentTasks functionality
func TestGetRecentTasks(t *testing.T) {
	mockVeriYonetici := new(MockVeriYoneticiAI)
	acy := YeniAIContextYonetici(mockVeriYonetici)

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

			// Setup mock expectations
			mockVeriYonetici.On("GorevGetir", tt.taskID).Return(tt.mockTask, tt.mockError)

			if tt.shouldUpdate && tt.mockTask != nil {
				mockVeriYonetici.On("GorevGuncelle", mock.MatchedBy(func(g *Gorev) bool {
					return g.ID == tt.taskID && g.Durum == "devam_ediyor"
				})).Return(nil)
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
	mockVeriYonetici.On("GorevGuncelle", mock.MatchedBy(func(g *Gorev) bool {
		return g.ID == "task-1" && g.Durum == "devam_ediyor"
	})).Return(nil)

	mockVeriYonetici.On("GorevGetir", "task-2").Return(&Gorev{ID: "task-2", Durum: "beklemede"}, nil)
	mockVeriYonetici.On("GorevGuncelle", mock.MatchedBy(func(g *Gorev) bool {
		return g.ID == "task-2" && g.Durum == "tamamlandi"
	})).Return(nil)

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
	assert.Equal(t, "task-not-found", result.Failed[0].ID)

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

// TestHelperFunctions tests helper functions
func TestHelperFunctions(t *testing.T) {
	// Test contains function
	assert.True(t, contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, contains([]string{"a", "b", "c"}, "d"))
	assert.False(t, contains([]string{}, "a"))
}
