package gorev

import (
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
		ID:              uuid.New().String(),
		Isim:            "Test Projesi",
		Tanim:           "Test açıklaması",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	err = vy.ProjeKaydet(proje)
	require.NoError(t, err)

	// Ana görev oluştur
	anaGorev := &Gorev{
		ID:              uuid.New().String(),
		Baslik:          "Ana Görev",
		Aciklama:        "Ana görev açıklaması",
		Durum:           "beklemede",
		Oncelik:         "yuksek",
		ProjeID:         proje.ID,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	err = vy.GorevKaydet(anaGorev)
	require.NoError(t, err)

	t.Run("Create Subtask", func(t *testing.T) {
		altGorev := &Gorev{
			ID:              uuid.New().String(),
			Baslik:          "Alt Görev 1",
			Aciklama:        "Alt görev açıklaması",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         proje.ID,
			ParentID:        anaGorev.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
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
			ID:              uuid.New().String(),
			Baslik:          "Alt Görev 2",
			Aciklama:        "İkinci alt görev",
			Durum:           "beklemede",
			Oncelik:         "dusuk",
			ProjeID:         proje.ID,
			ParentID:        anaGorev.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
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
			ID:              uuid.New().String(),
			Baslik:          "Alt Alt Görev",
			Aciklama:        "3. seviye görev",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         proje.ID,
			ParentID:        altGorevler[0].ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
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
			ID:              uuid.New().String(),
			Baslik:          "Yeni Ana Görev",
			Aciklama:        "Taşıma için yeni ana görev",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         proje.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
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
		assert.GreaterOrEqual(t, hiyerarsi.ToplamAltGorev, 1)
	})
}

func TestIsYonetici_AltGorevOperations(t *testing.T) {
	mockVY := NewMockVeriYonetici()
	iy := YeniIsYonetici(mockVY)

	// Test projesi ekle
	proje := &Proje{
		ID:              "proje-1",
		Isim:            "Test Projesi",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	mockVY.projeler[proje.ID] = proje

	// Ana görev ekle
	anaGorev := &Gorev{
		ID:              "gorev-1",
		Baslik:          "Ana Görev",
		Durum:           "beklemede",
		Oncelik:         "yuksek",
		ProjeID:         proje.ID,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
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
		assert.Equal(t, "Alt Görev", altGorev.Baslik)
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
		assert.Contains(t, err.Error(), "üst görev bulunamadı")
	})

	t.Run("Change Parent Task", func(t *testing.T) {
		// Yeni ana görev ekle
		yeniAnaGorev := &Gorev{
			ID:              "gorev-2",
			Baslik:          "Yeni Ana Görev",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         proje.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		mockVY.gorevler[yeniAnaGorev.ID] = yeniAnaGorev

		// Alt görev ekle
		altGorev := &Gorev{
			ID:              "altgorev-1",
			Baslik:          "Alt Görev",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         proje.ID,
			ParentID:        anaGorev.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		mockVY.gorevler[altGorev.ID] = altGorev

		err := iy.GorevUstDegistir(altGorev.ID, yeniAnaGorev.ID)
		assert.NoError(t, err)
		assert.Equal(t, yeniAnaGorev.ID, mockVY.gorevler[altGorev.ID].ParentID)
	})

	t.Run("Delete Task with Subtasks Should Fail", func(t *testing.T) {
		// Alt görev ekle
		altGorev := &Gorev{
			ID:              "altgorev-2",
			Baslik:          "Silinmeyi engelleyen alt görev",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         proje.ID,
			ParentID:        anaGorev.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		mockVY.gorevler[altGorev.ID] = altGorev

		err := iy.GorevSil(anaGorev.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "alt görev silinmeli")
	})

	t.Run("Complete Task with Incomplete Subtasks Should Fail", func(t *testing.T) {
		// Alt görevi beklemede durumunda ekle
		altGorev := &Gorev{
			ID:              "altgorev-3",
			Baslik:          "Tamamlanmamış alt görev",
			Durum:           "beklemede",
			Oncelik:         "orta",
			ProjeID:         proje.ID,
			ParentID:        anaGorev.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		}
		mockVY.gorevler[altGorev.ID] = altGorev

		err := iy.GorevDurumGuncelle(anaGorev.ID, "tamamlandi")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tüm alt görevler tamamlanmalı")
	})
}
