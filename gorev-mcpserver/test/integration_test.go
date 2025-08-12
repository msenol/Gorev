package test

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	mcphandlers "github.com/msenol/gorev/internal/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupIntegrationTestI18n initializes the i18n system for integration tests
func setupIntegrationTestI18n() {
	// Initialize i18n with Turkish (default) for tests
	i18n.Initialize("tr")
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
	setupIntegrationTestI18n() // Initialize i18n for tests
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Handler'ları oluştur
	handlers := mcphandlers.YeniHandlers(isYonetici)

	// Varsayılan template'leri oluştur
	err = veriYonetici.VarsayilanTemplateleriOlustur()
	require.NoError(t, err)

	// Test projesi oluştur
	proje, err := isYonetici.ProjeOlustur("Test Projesi", "Test açıklaması")
	require.NoError(t, err)

	// Projeyi aktif yap
	err = veriYonetici.AktifProjeAyarla(proje.ID)
	require.NoError(t, err)

	// Test: Template listesini kontrol et
	templates, err := veriYonetici.TemplateListele("")
	require.NoError(t, err)
	t.Logf("Available templates: %d", len(templates))
	for _, tmpl := range templates {
		t.Logf("Template: %s", tmpl.ID)
	}

	// Test: Template ile görev oluştur (gorev_olustur artık deprecated)
	// İlk template'i kullan (research template)
	templateID := templates[0].ID
	params := map[string]interface{}{
		"template_id": templateID,
		"degerler": map[string]interface{}{
			"konu":      "Test görevi",
			"amac":      "Bu bir test görevidir",
			"sorular":   "Test soruları",
			"kriterler": "Test kriterleri",
			"oncelik":   "yuksek",
		},
	}

	result, err := handlers.TemplatedenGorevOlustur(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "✓ Template kullanılarak görev oluşturuldu")
	assert.Contains(t, text, "Test görevi")

	// Test: Görevleri listele (sıralama ve filtreleme ile)
	listParams := map[string]interface{}{
		"sirala": "son_tarih_asc",
	}
	listResult, err := handlers.GorevListele(listParams)
	require.NoError(t, err)
	assert.False(t, listResult.IsError)
	listText := extractText(t, listResult)
	assert.Contains(t, listText, "Test görevi")
	assert.Contains(t, listText, "Y") // Compact format for "yuksek" priority
}

func TestGorevDurumGuncelle(t *testing.T) {
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

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
	assert.Contains(t, text, "✓ Görev güncellendi")
	assert.Contains(t, text, "devam_ediyor")

	// Güncellemeyi doğrula
	gorevler, err := isYonetici.GorevListele("devam_ediyor", "", "")
	require.NoError(t, err)
	assert.Len(t, gorevler, 1)
	assert.Equal(t, "devam_ediyor", gorevler[0].Durum)
}

func TestProjeOlustur(t *testing.T) {
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

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
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

	// Test verisi oluştur
	_, err = isYonetici.ProjeOlustur("Proje 1", "")
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
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

	// Test: Deprecated GorevOlustur method
	params := map[string]interface{}{
		"aciklama": "Başlıksız görev",
	}

	result, err := handlers.GorevOlustur(params)
	require.NoError(t, err)
	assert.True(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "gorev_olustur artık kullanılmıyor")

	// Test: Geçersiz ID ile güncelleme
	updateParams := map[string]interface{}{
		"id":    "gecersiz-id",
		"durum": "tamamlandi",
	}

	result, err = handlers.GorevGuncelle(updateParams)
	require.NoError(t, err)
	assert.True(t, result.IsError)
	text2 := extractText(t, result)
	assert.Contains(t, text2, "görev güncellenemedi")
}

func TestGorevDetay(t *testing.T) {
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

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
	assert.Contains(t, detayText, "**Son Teslim Tarihi:** 2025-12-31")
	// Etiketler farklı sırada olabilir
	assert.Contains(t, detayText, "**Etiketler:**")
	assert.Contains(t, detayText, "bug")
	assert.Contains(t, detayText, "acil")
	assert.Contains(t, detayText, "## 🔗 Bağımlılıklar")
	assert.Contains(t, detayText, "### 🎯 Bu göreve bağımlı görevler:")
	assert.Contains(t, detayText, "- Bağlı Görev")
}

func TestGorevDuzenle(t *testing.T) {
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

	// Test görevi oluştur
	gorevObj, err := isYonetici.GorevOlustur("Eski Başlık", "Eski açıklama", "orta", "", "", nil)
	require.NoError(t, err)

	// Başlık ve açıklama güncelle
	params := map[string]interface{}{
		"id":       gorevObj.ID,
		"baslik":   "Yeni Başlık",
		"aciklama": "## Yeni Açıklama\n\nMarkdown destekli",
		"oncelik":  "yuksek",
	}

	result, err := handlers.GorevDuzenle(params)
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "✓ Görev düzenlendi")

	// Değişiklikleri doğrula
	guncelGorev, err := isYonetici.GorevGetir(gorevObj.ID)
	require.NoError(t, err)
	assert.Equal(t, "Yeni Başlık", guncelGorev.Baslik)
	assert.Equal(t, "## Yeni Açıklama\n\nMarkdown destekli", guncelGorev.Aciklama)
	assert.Equal(t, "yuksek", guncelGorev.Oncelik)
}

func TestGorevSil(t *testing.T) {
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

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
	setupIntegrationTestI18n() // Initialize i18n for tests
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

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
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

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
	// Test veritabanı oluştur
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	// İş yöneticisi ve handler'ları oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcphandlers.YeniHandlers(isYonetici)

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
	assert.Equal(t, gorev1.ID, baglantilar[0].KaynakID)
	assert.Equal(t, gorev2.ID, baglantilar[0].HedefID)
	assert.Equal(t, "engelliyor", baglantilar[0].BaglantiTip)
}
