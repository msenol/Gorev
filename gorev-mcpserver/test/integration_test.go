package test

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	mcphandlers "github.com/msenol/gorev/internal/mcp"
	testinghelpers "github.com/msenol/gorev/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTestConfigWithEmbeddedMigrations returns a test config with embedded migrations
func getTestConfigWithEmbeddedMigrations(t *testing.T) *testinghelpers.TestDatabaseConfig {
	migrationsFS, err := getEmbeddedMigrationsFS()
	require.NoError(t, err, "Failed to get embedded migrations FS")

	return &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsFS:    migrationsFS,
		CreateTemplates: true,
		InitializeI18n:  true,
	}
}

// Helper function to extract text from MCP result
func extractText(t *testing.T, result *mcp.CallToolResult) string {
	require.Len(t, result.Content, 1)

	// Content is []interface{}, and each element should be a TextContent
	textContent, ok := result.Content[0].(mcp.TextContent)
	if ok {
		return textContent.Text
	}

	// If it's a map (serialized form)
	contentMap, ok := result.Content[0].(map[string]interface{})
	require.True(t, ok, "Content should be TextContent or map")
	text, ok := contentMap["text"].(string)
	require.True(t, ok, "Content should have text field")
	return text
}

func TestGorevOlusturVeListele(t *testing.T) {
	// Create test environment using embedded migrations
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test projesi oluştur
	proje, err := isYonetici.ProjeOlustur("Test Projesi", "Test açıklaması")
	require.NoError(t, err)

	// Projeyi aktif yap
	err = isYonetici.VeriYonetici().AktifProjeAyarla(proje.ID)
	require.NoError(t, err)

	// Test: Template listesini kontrol et
	templates, err := isYonetici.VeriYonetici().TemplateListele("")
	require.NoError(t, err)
	t.Logf("Available templates: %d", len(templates))
	for _, tmpl := range templates {
		t.Logf("Template: %s", tmpl.ID)
	}

	// Test: Template ile görev oluştur (gorev_olustur artık deprecated)
	// Bug Raporu template'ini bul
	var firstBugTemplate *gorev.GorevTemplate
	for _, tmpl := range templates {
		if tmpl.Name == "Bug Raporu" {
			firstBugTemplate = tmpl
			break
		}
	}
	if firstBugTemplate == nil {
		t.Fatal("Bug Raporu template not found")
	}

	params := map[string]interface{}{
		constants.ParamTemplateID: firstBugTemplate.ID,
		constants.ParamValues: map[string]interface{}{
			"title":       "Test Bug Görevi",
			"description": "Bu bir test bug raporu",
			"modul":       "test-integration",
			"ortam":       "development",
			"adimlar":     "1. Test çalıştır",
			"beklenen":    "Başarı",
			"mevcut":      "Hata",
			"priority":    "yuksek",
		},
	}

	result, err := handlers.TemplatedenGorevOlustur(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "✓ Template kullanılarak görev oluşturuldu")
	assert.Contains(t, text, "Test Bug Görevi")

	// Test: Görevleri listele (sıralama ve filtreleme ile)
	listParams := map[string]interface{}{
		"sirala": "son_tarih_asc",
	}
	listResult, err := handlers.GorevListele(listParams)
	require.NoError(t, err)
	assert.False(t, listResult.IsError)
	listText := extractText(t, listResult)
	// Check for tasks that were created - the first test creates only one task
	assert.Contains(t, listText, "Test Bug Görevi")
	assert.Contains(t, listText, "Y") // Compact format for "yuksek" priority
}

func TestGorevDurumGuncelle(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Önce bir görev oluştur
	gorevObj, err := isYonetici.GorevOlustur("Durum test görevi", "", "orta", "", "", nil)
	require.NoError(t, err)

	// Durumu güncelle
	updateParams := map[string]interface{}{
		"id":    gorevObj.ID,
		"durum": "devam_ediyor",
	}

	result, err := handlers.GorevGuncelle(updateParams)
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "Görev güncellendi")
	assert.Contains(t, text, "devam_ediyor")

	// Güncellemeyi doğrula
	gorevler, err := isYonetici.GorevListele(map[string]interface{}{"durum": "devam_ediyor"})
	require.NoError(t, err)
	assert.Len(t, gorevler, 1)
	assert.Equal(t, "devam_ediyor", gorevler[0].Status)
}

func TestProjeOlustur(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Proje oluştur
	params := map[string]interface{}{
		"isim":  "Test Projesi",
		"tanim": "Test amaçlı proje",
	}

	result, err := handlers.ProjeOlustur(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "✓ Proje oluşturuldu")
	assert.Contains(t, text, "Test Projesi")
}

func TestOzetGoster(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test verisi oluştur
	_, err := isYonetici.ProjeOlustur("Proje 1", "")
	require.NoError(t, err)

	_, err = isYonetici.GorevOlustur("Görev 1", "", "yuksek", "", "", nil)
	require.NoError(t, err)

	gorev2, err := isYonetici.GorevOlustur("Görev 2", "", "orta", "", "", nil)
	require.NoError(t, err)

	err = isYonetici.GorevDurumGuncelle(gorev2.ID, "tamamlandi")
	require.NoError(t, err)

	// Özet al
	result, err := handlers.OzetGoster(map[string]interface{}{})
	require.NoError(t, err)

	assert.False(t, result.IsError)
	ozetText := extractText(t, result)
	assert.Contains(t, ozetText, "Toplam Proje:** 1")
	assert.Contains(t, ozetText, "Toplam Görev:** 2")
	assert.Contains(t, ozetText, "Beklemede: 1")
	assert.Contains(t, ozetText, "Tamamlandı: 1")
	assert.Contains(t, ozetText, "Yüksek: 1")
	assert.Contains(t, ozetText, "Orta: 1")
}

func TestHataYonetimi(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Varsayılan template'leri oluştur (idempotent)
	err := isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur()
	require.NoError(t, err)

	// Test projesi oluştur
	proje, err := isYonetici.ProjeOlustur("Test Projesi", "Test açıklaması")
	require.NoError(t, err)

	// Projeyi aktif yap
	err = isYonetici.VeriYonetici().AktifProjeAyarla(proje.ID)
	require.NoError(t, err)

	// Get available templates
	templates, err := isYonetici.VeriYonetici().TemplateListele("")
	require.NoError(t, err)
	if len(templates) == 0 {
		t.Fatal("No templates available for testing")
	}

	// Find the "Bug Raporu" template specifically (templates are sorted by category, name)
	var bugTemplate *gorev.GorevTemplate
	for _, tmpl := range templates {
		if tmpl.Name == "Bug Raporu" {
			bugTemplate = tmpl
			break
		}
	}
	if bugTemplate == nil {
		t.Fatal("Bug Raporu template not found")
	}

	templateParams := map[string]interface{}{
		constants.ParamTemplateID: bugTemplate.ID,
		constants.ParamValues: map[string]interface{}{
			"title":       "Integration Test Bug",
			"description": "Template ile oluşturulan test bug raporu",
			"modul":       "integration-test",
			"ortam":       "development",
			"adimlar":     "1. Test çalıştır 2. Hatayı gözle",
			"beklenen":    "Test başarılı olması",
			"mevcut":      "Test hata verdi",
			"priority":    "orta",
		},
	}

	result, err := handlers.TemplatedenGorevOlustur(templateParams)
	require.NoError(t, err)
	assert.False(t, result.IsError, "Template creation should succeed: %s", extractText(t, result))
	text := extractText(t, result)
	assert.Contains(t, text, "Integration Test Bug")

	// Test: Geçersiz ID ile güncelleme
	updateParams := map[string]interface{}{
		"id":    "gecersiz-id",
		"durum": "tamamlandi",
	}

	result, err = handlers.GorevGuncelle(updateParams)
	require.NoError(t, err)
	assert.True(t, result.IsError)
	text2 := extractText(t, result)
	assert.Contains(t, text2, "güncellenemedi")
}

func TestGorevDetay(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test verisi oluştur
	proje, err := isYonetici.ProjeOlustur("Test Projesi", "Proje açıklaması")
	require.NoError(t, err)

	gorev1, err := isYonetici.GorevOlustur("Detaylı Test Görevi", "## Açıklama\n\nBu bir **markdown** açıklamadır.", "yuksek", proje.ID, "2025-12-31", []string{"bug", "acil"})
	require.NoError(t, err)

	gorev2, err := isYonetici.GorevOlustur("Bağlı Görev", "", "orta", proje.ID, "", nil)
	require.NoError(t, err)

	// Bağımlılık ekle (gorev1 önce tamamlanmalı, sonra gorev2 başlayabilir)
	_, err = isYonetici.GorevBagimlilikEkle(gorev1.ID, gorev2.ID, "onceki")
	require.NoError(t, err)

	// Detay al
	params := map[string]interface{}{
		"id": gorev1.ID,
	}

	result, err := handlers.GorevDetay(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)

	detayText := extractText(t, result)
	assert.Contains(t, detayText, "# Detaylı Test Görevi")
	assert.Contains(t, detayText, "**Proje:** Test Projesi")
	assert.Contains(t, detayText, "**Son Tarih:** 2025-12-31")
	// Etiketler farklı sırada olabilir
	assert.Contains(t, detayText, "**Etiketler:**")
	assert.Contains(t, detayText, "bug")
	assert.Contains(t, detayText, "acil")
	assert.Contains(t, detayText, "## 🔗 Bağımlılıklar")
	assert.Contains(t, detayText, "### 🎯 Bu göreve bağımlı görevler:")
	assert.Contains(t, detayText, "- Bağlı Görev")
}

func TestGorevDuzenle(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test görevi oluştur
	gorevObj, err := isYonetici.GorevOlustur("Eski Başlık", "Eski açıklama", "orta", "", "", nil)
	require.NoError(t, err)

	// Başlık ve açıklama güncelle
	params := map[string]interface{}{
		"id":          gorevObj.ID,
		"title":       "Yeni Başlık",
		"description": "## Yeni Açıklama\n\nMarkdown destekli",
		"priority":    "yuksek",
	}

	result, err := handlers.GorevDuzenle(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "✓ Görev düzenlendi")

	// Değişiklikleri doğrula
	guncelGorev, err := isYonetici.GorevGetir(gorevObj.ID)
	require.NoError(t, err)
	assert.Equal(t, "Yeni Başlık", guncelGorev.Title)
	assert.Equal(t, "## Yeni Açıklama\n\nMarkdown destekli", guncelGorev.Description)
	assert.Equal(t, "yuksek", guncelGorev.Priority)
}

func TestGorevSil(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test görevi oluştur
	gorevObj, err := isYonetici.GorevOlustur("Silinecek Görev", "", "orta", "", "", nil)
	require.NoError(t, err)

	// Onaysız silme denemesi
	params := map[string]interface{}{
		"id":   gorevObj.ID,
		"onay": false,
	}

	result, err := handlers.GorevSil(params)
	require.NoError(t, err)
	assert.True(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "onay' parametresi true olmalıdır")

	// Onaylı silme
	params["onay"] = true
	result, err = handlers.GorevSil(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text = extractText(t, result)
	assert.Contains(t, text, "✓ Görev silindi: Silinecek Görev")

	// Silinen görevi arama
	_, err = isYonetici.GorevGetir(gorevObj.ID)
	assert.Error(t, err)
}

func TestProjeListele(t *testing.T) {
	// Create test environment using standardized helpers
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPathIntegration,
		CreateTemplates: true,
		InitializeI18n:  true,
	}
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test projeleri oluştur
	proje1, err := isYonetici.ProjeOlustur("Proje 1", "İlk proje")
	require.NoError(t, err)

	_, err = isYonetici.ProjeOlustur("Proje 2", "İkinci proje")
	require.NoError(t, err)

	// Proje 1'e görevler ekle
	gorev1, err := isYonetici.GorevOlustur("Görev 1", "", "yuksek", "", "", nil)
	require.NoError(t, err)
	err = isYonetici.GorevDuzenle(gorev1.ID, "", "", "", proje1.ID, "", false, false, false, true, false)
	require.NoError(t, err)

	gorev2, err := isYonetici.GorevOlustur("Görev 2", "", "orta", "", "", nil)
	require.NoError(t, err)
	err = isYonetici.GorevDuzenle(gorev2.ID, "", "", "", proje1.ID, "", false, false, false, true, false)
	require.NoError(t, err)

	// Projeleri listele
	result, err := handlers.ProjeListele(map[string]interface{}{})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	listText := extractText(t, result)
	assert.Contains(t, listText, "## Proje Listesi")
	assert.Contains(t, listText, "### Proje 1")
	assert.Contains(t, listText, "**Tanım:** İlk proje")
	assert.Contains(t, listText, "**Görev Sayısı:** 2")
	assert.Contains(t, listText, "### Proje 2")
	assert.Contains(t, listText, "**Görev Sayısı:** 0")
}

func TestProjeGorevleri(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test projesi ve görevleri oluştur
	proje, err := isYonetici.ProjeOlustur("Test Projesi", "")
	require.NoError(t, err)

	// Farklı durumlarda görevler oluştur
	gorev1, err := isYonetici.GorevOlustur("Devam Eden Görev", "Açıklama 1", "yuksek", "", "", nil)
	require.NoError(t, err)
	err = isYonetici.GorevDuzenle(gorev1.ID, "", "", "", proje.ID, "", false, false, false, true, false)
	require.NoError(t, err)
	err = isYonetici.GorevDurumGuncelle(gorev1.ID, "devam_ediyor")
	require.NoError(t, err)

	gorev2, err := isYonetici.GorevOlustur("Bekleyen Görev", "Açıklama 2", "orta", "", "", nil)
	require.NoError(t, err)
	err = isYonetici.GorevDuzenle(gorev2.ID, "", "", "", proje.ID, "", false, false, false, true, false)
	require.NoError(t, err)

	gorev3, err := isYonetici.GorevOlustur("Tamamlanan Görev", "", "dusuk", "", "", nil)
	require.NoError(t, err)
	err = isYonetici.GorevDuzenle(gorev3.ID, "", "", "", proje.ID, "", false, false, false, true, false)
	require.NoError(t, err)
	err = isYonetici.GorevDurumGuncelle(gorev3.ID, "tamamlandi")
	require.NoError(t, err)

	// Proje görevlerini listele
	params := map[string]interface{}{
		"proje_id": proje.ID,
	}

	result, err := handlers.ProjeGorevleri(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)

	gorevlerText := extractText(t, result)
	assert.Contains(t, gorevlerText, "Test Projesi (3 görev)")
	assert.Contains(t, gorevlerText, "🔵 Devam Ediyor")
	assert.Contains(t, gorevlerText, "**Devam Eden Görev** (Y)")
	assert.Contains(t, gorevlerText, "⚪ Beklemede")
	assert.Contains(t, gorevlerText, "**Bekleyen Görev** (O)")
	assert.Contains(t, gorevlerText, "✅ Tamamlandı")
	assert.Contains(t, gorevlerText, "~~Tamamlanan Görev~~")
}

func TestGorevBagimlilikEkle(t *testing.T) {
	// Create test environment using standardized helpers
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// İki test görevi oluştur
	gorev1, err := isYonetici.GorevOlustur("Kaynak Görev", "", "orta", "", "", nil)
	require.NoError(t, err)
	gorev2, err := isYonetici.GorevOlustur("Hedef Görev", "", "yuksek", "", "", nil)
	require.NoError(t, err)

	// Bağımlılık ekle
	params := map[string]interface{}{
		"kaynak_id":     gorev1.ID,
		"hedef_id":      gorev2.ID,
		"baglanti_tipi": "engelliyor",
	}

	result, err := handlers.GorevBagimlilikEkle(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "✓ Bağımlılık eklendi")

	// Bağımlılığı doğrula
	baglantilar, err := isYonetici.GorevBaglantilariGetir(gorev1.ID)
	require.NoError(t, err)
	assert.Len(t, baglantilar, 1)
	assert.Equal(t, gorev1.ID, baglantilar[0].SourceID)
	assert.Equal(t, gorev2.ID, baglantilar[0].TargetID)
	assert.Equal(t, "engelliyor", baglantilar[0].ConnectionType)
}
