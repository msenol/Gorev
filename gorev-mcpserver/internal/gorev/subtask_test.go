package gorev

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAltGorevOperations(t *testing.T) {
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer vy.Kapat()

	// Proje oluştur
	proje := &Proje{
		ID:         uuid.New().String(),
		Name:       "Test Projesi",
		Definition: "Test açıklaması",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = vy.ProjeKaydet(proje)
	require.NoError(t, err)

	// Ana görev oluştur
	anaGorev := &Gorev{
		ID:          uuid.New().String(),
		Title:       "Ana Görev",
		Description: "Ana görev açıklaması",
		Status:      "beklemede",
		Priority:    "yuksek",
		ProjeID:     proje.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = vy.GorevKaydet(anaGorev)
	require.NoError(t, err)

	t.Run("Create Subtask", func(t *testing.T) {
		altGorev := &Gorev{
			ID:          uuid.New().String(),
			Title:       "Alt Görev 1",
			Description: "Alt görev açıklaması",
			Status:      "beklemede",
			Priority:    "orta",
			ProjeID:     proje.ID,
			ParentID:    anaGorev.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = vy.GorevKaydet(altGorev)
		assert.NoError(t, err)

		// Alt görevi getir ve kontrol et
		gorev, err := vy.GorevGetir(altGorev.ID)
		assert.NoError(t, err)
		assert.Equal(t, anaGorev.ID, gorev.ParentID)
	})

	t.Run("Get Direct Subtasks", func(t *testing.T) {
		// İkinci alt görev oluştur
		altGorev2 := &Gorev{
			ID:          uuid.New().String(),
			Title:       "Alt Görev 2",
			Description: "İkinci alt görev",
			Status:      "beklemede",
			Priority:    "dusuk",
			ProjeID:     proje.ID,
			ParentID:    anaGorev.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = vy.GorevKaydet(altGorev2)
		require.NoError(t, err)

		// Alt görevleri getir
		altGorevler, err := vy.AltGorevleriGetir(anaGorev.ID)
		assert.NoError(t, err)
		assert.Len(t, altGorevler, 2)
	})

	t.Run("Get Parent Tasks", func(t *testing.T) {
		// Alt görevin alt görevi oluştur (3. seviye)
		altGorevler, _ := vy.AltGorevleriGetir(anaGorev.ID)
		altAltGorev := &Gorev{
			ID:          uuid.New().String(),
			Title:       "Alt Alt Görev",
			Description: "3. seviye görev",
			Status:      "beklemede",
			Priority:    "orta",
			ProjeID:     proje.ID,
			ParentID:    altGorevler[0].ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = vy.GorevKaydet(altAltGorev)
		require.NoError(t, err)

		// Üst görevleri getir
		ustGorevler, err := vy.UstGorevleriGetir(altAltGorev.ID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(ustGorevler), 1)
	})

	t.Run("Update Parent ID", func(t *testing.T) {
		// Yeni ana görev oluştur
		yeniAnaGorev := &Gorev{
			ID:          uuid.New().String(),
			Title:       "Yeni Ana Görev",
			Description: "Taşıma için yeni ana görev",
			Status:      "beklemede",
			Priority:    "orta",
			ProjeID:     proje.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = vy.GorevKaydet(yeniAnaGorev)
		require.NoError(t, err)

		// Alt görevi taşı
		altGorevler, _ := vy.AltGorevleriGetir(anaGorev.ID)
		err = vy.ParentIDGuncelle(altGorevler[0].ID, yeniAnaGorev.ID)
		assert.NoError(t, err)

		// Kontrol et
		tasinanGorev, err := vy.GorevGetir(altGorevler[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, yeniAnaGorev.ID, tasinanGorev.ParentID)
	})

	t.Run("Circular Dependency Check", func(t *testing.T) {
		// Kendisine parent olamaz
		daireVar, err := vy.DaireBagimliligiKontrolEt(anaGorev.ID, anaGorev.ID)
		assert.NoError(t, err)
		assert.True(t, daireVar)

		// Alt görev ana göreve parent olamaz
		altGorevler, _ := vy.AltGorevleriGetir(anaGorev.ID)
		if len(altGorevler) > 0 {
			daireVar, err = vy.DaireBagimliligiKontrolEt(anaGorev.ID, altGorevler[0].ID)
			assert.NoError(t, err)
			assert.True(t, daireVar)
		}
	})

	t.Run("Task Hierarchy Info", func(t *testing.T) {
		hiyerarsi, err := vy.GorevHiyerarsiGetir(anaGorev.ID)
		assert.NoError(t, err)
		assert.NotNil(t, hiyerarsi)
		assert.Equal(t, anaGorev.ID, hiyerarsi.Gorev.ID)
		assert.GreaterOrEqual(t, hiyerarsi.TotalSubtasks, 1)
	})
}

func TestIsYonetici_AltGorevOperations(t *testing.T) {
	mockVY := NewMockVeriYonetici()
	iy := YeniIsYonetici(mockVY)

	// Test projesi ekle
	proje := &Proje{
		ID:        "proje-1",
		Name:      "Test Projesi",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockVY.projeler[proje.ID] = proje

	// Ana görev ekle
	anaGorev := &Gorev{
		ID:        "gorev-1",
		Title:     "Ana Görev",
		Status:    "beklemede",
		Priority:  "yuksek",
		ProjeID:   proje.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockVY.gorevler[anaGorev.ID] = anaGorev

	t.Run("Create Subtask", func(t *testing.T) {
		altGorev, err := iy.AltGorevOlustur(
			anaGorev.ID,
			"Alt Görev",
			"Alt görev açıklaması",
			"orta",
			"",
			[]string{"test", "altgorev"},
		)

		assert.NoError(t, err)
		assert.NotNil(t, altGorev)
		assert.Equal(t, anaGorev.ID, altGorev.ParentID)
		assert.Equal(t, proje.ID, altGorev.ProjeID)
		assert.Equal(t, "Alt Görev", altGorev.Title)
	})

	t.Run("Create Subtask - Parent Not Found", func(t *testing.T) {
		altGorev, err := iy.AltGorevOlustur(
			"olmayan-gorev",
			"Alt Görev",
			"Açıklama",
			"orta",
			"",
			nil,
		)

		assert.Error(t, err)
		assert.Nil(t, altGorev)
		assert.Contains(t, err.Error(), "parentTaskNotFound")
	})

	t.Run("Change Parent Task", func(t *testing.T) {
		// Yeni ana görev ekle
		yeniAnaGorev := &Gorev{
			ID:        "gorev-2",
			Title:     "Yeni Ana Görev",
			Status:    "beklemede",
			Priority:  "orta",
			ProjeID:   proje.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockVY.gorevler[yeniAnaGorev.ID] = yeniAnaGorev

		// Alt görev ekle
		altGorev := &Gorev{
			ID:        "altgorev-1",
			Title:     "Alt Görev",
			Status:    "beklemede",
			Priority:  "orta",
			ProjeID:   proje.ID,
			ParentID:  anaGorev.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockVY.gorevler[altGorev.ID] = altGorev

		err := iy.GorevUstDegistir(altGorev.ID, yeniAnaGorev.ID)
		assert.NoError(t, err)
		assert.Equal(t, yeniAnaGorev.ID, mockVY.gorevler[altGorev.ID].ParentID)
	})

	t.Run("Delete Task with Subtasks Should Fail", func(t *testing.T) {
		// Alt görev ekle
		altGorev := &Gorev{
			ID:        "altgorev-2",
			Title:     "Silinmeyi engelleyen alt görev",
			Status:    "beklemede",
			Priority:  "orta",
			ProjeID:   proje.ID,
			ParentID:  anaGorev.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockVY.gorevler[altGorev.ID] = altGorev

		err := iy.GorevSil(anaGorev.ID)
		assert.Error(t, err)
		// Check for i18n key or translated text
		errMsg := err.Error()
		if !strings.Contains(errMsg, "taskHasSubtasksCannotDelete") && !strings.Contains(errMsg, "bu görev silinemez") {
			t.Errorf("Expected subtask deletion error, got: %s", errMsg)
		}
	})

	t.Run("Complete Task with Incomplete Subtasks Should Fail", func(t *testing.T) {
		// Alt görevi beklemede durumunda ekle
		altGorev := &Gorev{
			ID:        "altgorev-3",
			Title:     "Tamamlanmamış alt görev",
			Status:    "beklemede",
			Priority:  "orta",
			ProjeID:   proje.ID,
			ParentID:  anaGorev.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockVY.gorevler[altGorev.ID] = altGorev

		err := iy.GorevDurumGuncelle(anaGorev.ID, "tamamlandi")
		assert.Error(t, err)
		// Check for i18n key or translated text
		errMsg := err.Error()
		if !strings.Contains(errMsg, "taskCannotCompleteSubtasks") && !strings.Contains(errMsg, "bu görev tamamlanamaz") {
			t.Errorf("Expected subtask completion error, got: %s", errMsg)
		}
	})
}
