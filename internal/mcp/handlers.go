package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yourusername/gorev/internal/gorev"
)

type Handlers struct {
	isYonetici *gorev.IsYonetici
}

func YeniHandlers(isYonetici *gorev.IsYonetici) *Handlers {
	return &Handlers{
		isYonetici: isYonetici,
	}
}

// GorevOlustur yeni bir görev oluşturur
func (h *Handlers) GorevOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	baslik, ok := params["baslik"].(string)
	if !ok || baslik == "" {
		return mcp.NewToolResultError("başlık parametresi gerekli"), nil
	}

	aciklama, _ := params["aciklama"].(string)
	oncelik, _ := params["oncelik"].(string)
	if oncelik == "" {
		oncelik = "orta"
	}

	gorev, err := h.isYonetici.GorevOlustur(baslik, aciklama, oncelik)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görev oluşturulamadı: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("✓ Görev oluşturuldu: %s (ID: %s)", gorev.Baslik, gorev.ID),
	), nil
}

// GorevListele görevleri listeler
func (h *Handlers) GorevListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	durum, _ := params["durum"].(string)

	gorevler, err := h.isYonetici.GorevListele(durum)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görevler listelenemedi: %v", err)), nil
	}

	if len(gorevler) == 0 {
		return mcp.NewToolResultText("Henüz görev bulunmuyor."), nil
	}

	metin := "## Görev Listesi\n\n"
	for _, gorev := range gorevler {
		metin += fmt.Sprintf("- [%s] %s (%s öncelik)\n", gorev.Durum, gorev.Baslik, gorev.Oncelik)
		if gorev.Aciklama != "" {
			metin += fmt.Sprintf("  %s\n", gorev.Aciklama)
		}
		metin += fmt.Sprintf("  ID: %s\n\n", gorev.ID)
	}

	return mcp.NewToolResultText(metin), nil
}

// GorevGuncelle görev durumunu günceller
func (h *Handlers) GorevGuncelle(params map[string]interface{}) (*mcp.CallToolResult, error) {
	id, ok := params["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("id parametresi gerekli"), nil
	}

	durum, ok := params["durum"].(string)
	if !ok || durum == "" {
		return mcp.NewToolResultError("durum parametresi gerekli"), nil
	}

	if err := h.isYonetici.GorevDurumGuncelle(id, durum); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görev güncellenemedi: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("✓ Görev güncellendi: %s → %s", id, durum),
	), nil
}

// ProjeOlustur yeni bir proje oluşturur
func (h *Handlers) ProjeOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	isim, ok := params["isim"].(string)
	if !ok || isim == "" {
		return mcp.NewToolResultError("isim parametresi gerekli"), nil
	}

	tanim, _ := params["tanim"].(string)

	proje, err := h.isYonetici.ProjeOlustur(isim, tanim)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("proje oluşturulamadı: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf("✓ Proje oluşturuldu: %s (ID: %s)", proje.Isim, proje.ID),
	), nil
}

// GorevDetay tek bir görevin detaylı bilgisini markdown formatında döner
func (h *Handlers) GorevDetay(params map[string]interface{}) (*mcp.CallToolResult, error) {
	id, ok := params["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("id parametresi gerekli"), nil
	}

	gorev, err := h.isYonetici.GorevDetayAl(id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görev bulunamadı: %v", err)), nil
	}

	// Markdown formatında detaylı görev bilgisi
	metin := fmt.Sprintf(`# %s

## 📋 Genel Bilgiler
- **ID:** %s
- **Durum:** %s
- **Öncelik:** %s
- **Oluşturma Tarihi:** %s
- **Son Güncelleme:** %s`,
		gorev.Baslik,
		gorev.ID,
		gorev.Durum,
		gorev.Oncelik,
		gorev.OlusturmaTarih.Format("2006-01-02 15:04:05"),
		gorev.GuncellemeTarih.Format("2006-01-02 15:04:05"),
	)

	if gorev.ProjeID != "" {
		proje, err := h.isYonetici.ProjeDetayAl(gorev.ProjeID)
		if err == nil {
			metin += fmt.Sprintf("\n- **Proje:** %s", proje.Isim)
		}
	}

	metin += "\n\n## 📝 Açıklama\n"
	if gorev.Aciklama != "" {
		// Açıklama zaten markdown formatında olabilir, direkt ekle
		metin += gorev.Aciklama
	} else {
		metin += "*Açıklama girilmemiş*"
	}

	metin += "\n\n---\n"
	metin += fmt.Sprintf("\n*Son güncelleme: %s*", gorev.GuncellemeTarih.Format("02 Jan 2006, 15:04"))

	return mcp.NewToolResultText(metin), nil
}

// GorevDuzenle görevi düzenler
func (h *Handlers) GorevDuzenle(params map[string]interface{}) (*mcp.CallToolResult, error) {
	id, ok := params["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("id parametresi gerekli"), nil
	}

	// En az bir düzenleme alanı olmalı
	baslik, baslikVar := params["baslik"].(string)
	aciklama, aciklamaVar := params["aciklama"].(string)
	oncelik, oncelikVar := params["oncelik"].(string)
	projeID, projeVar := params["proje_id"].(string)

	if !baslikVar && !aciklamaVar && !oncelikVar && !projeVar {
		return mcp.NewToolResultError("en az bir düzenleme alanı belirtilmeli (baslik, aciklama, oncelik veya proje_id)"), nil
	}

	err := h.isYonetici.GorevDuzenle(id, baslik, aciklama, oncelik, projeID, baslikVar, aciklamaVar, oncelikVar, projeVar)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görev düzenlenemedi: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✓ Görev düzenlendi: %s", id)), nil
}

// GorevSil görevi siler
func (h *Handlers) GorevSil(params map[string]interface{}) (*mcp.CallToolResult, error) {
	id, ok := params["id"].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("id parametresi gerekli"), nil
	}

	// Onay kontrolü
	onay, onayVar := params["onay"].(bool)
	if !onayVar || !onay {
		return mcp.NewToolResultError("görevi silmek için 'onay' parametresi true olmalıdır"), nil
	}

	gorev, err := h.isYonetici.GorevDetayAl(id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görev bulunamadı: %v", err)), nil
	}

	gorevBaslik := gorev.Baslik

	err = h.isYonetici.GorevSil(id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görev silinemedi: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✓ Görev silindi: %s (ID: %s)", gorevBaslik, id)), nil
}

// ProjeListele tüm projeleri listeler
func (h *Handlers) ProjeListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	projeler, err := h.isYonetici.ProjeListele()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("projeler listelenemedi: %v", err)), nil
	}

	if len(projeler) == 0 {
		return mcp.NewToolResultText("Henüz proje bulunmuyor."), nil
	}

	metin := "## Proje Listesi\n\n"
	for _, proje := range projeler {
		metin += fmt.Sprintf("### %s\n", proje.Isim)
		metin += fmt.Sprintf("- **ID:** %s\n", proje.ID)
		if proje.Tanim != "" {
			metin += fmt.Sprintf("- **Tanım:** %s\n", proje.Tanim)
		}
		metin += fmt.Sprintf("- **Oluşturma:** %s\n", proje.OlusturmaTarih.Format("02 Jan 2006, 15:04"))

		// Her proje için görev sayısını göster
		gorevSayisi, err := h.isYonetici.ProjeGorevSayisi(proje.ID)
		if err == nil {
			metin += fmt.Sprintf("- **Görev Sayısı:** %d\n", gorevSayisi)
		}
		metin += "\n"
	}

	return mcp.NewToolResultText(metin), nil
}

// ProjeGorevleri bir projenin görevlerini listeler
func (h *Handlers) ProjeGorevleri(params map[string]interface{}) (*mcp.CallToolResult, error) {
	projeID, ok := params["proje_id"].(string)
	if !ok || projeID == "" {
		return mcp.NewToolResultError("proje_id parametresi gerekli"), nil
	}

	// Önce projenin var olduğunu kontrol et
	proje, err := h.isYonetici.ProjeDetayAl(projeID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("proje bulunamadı: %v", err)), nil
	}

	gorevler, err := h.isYonetici.ProjeGorevleri(projeID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görevler alınamadı: %v", err)), nil
	}

	metin := fmt.Sprintf("## %s - Görevler\n\n", proje.Isim)

	if len(gorevler) == 0 {
		metin += "*Bu projede henüz görev bulunmuyor.*"
		return mcp.NewToolResultText(metin), nil
	}

	// Duruma göre grupla
	beklemede := []*gorev.Gorev{}
	devamEdiyor := []*gorev.Gorev{}
	tamamlandi := []*gorev.Gorev{}

	for _, g := range gorevler {
		switch g.Durum {
		case "beklemede":
			beklemede = append(beklemede, g)
		case "devam_ediyor":
			devamEdiyor = append(devamEdiyor, g)
		case "tamamlandi":
			tamamlandi = append(tamamlandi, g)
		}
	}

	// Devam eden görevler
	if len(devamEdiyor) > 0 {
		metin += "### 🔵 Devam Ediyor\n"
		for _, g := range devamEdiyor {
			metin += fmt.Sprintf("- **%s** (%s öncelik)\n", g.Baslik, g.Oncelik)
			if g.Aciklama != "" {
				metin += fmt.Sprintf("  %s\n", g.Aciklama)
			}
			metin += fmt.Sprintf("  `ID: %s`\n", g.ID)
		}
		metin += "\n"
	}

	// Bekleyen görevler
	if len(beklemede) > 0 {
		metin += "### ⚪ Beklemede\n"
		for _, g := range beklemede {
			metin += fmt.Sprintf("- **%s** (%s öncelik)\n", g.Baslik, g.Oncelik)
			if g.Aciklama != "" {
				metin += fmt.Sprintf("  %s\n", g.Aciklama)
			}
			metin += fmt.Sprintf("  `ID: %s`\n", g.ID)
		}
		metin += "\n"
	}

	// Tamamlanan görevler
	if len(tamamlandi) > 0 {
		metin += "### ✅ Tamamlandı\n"
		for _, g := range tamamlandi {
			metin += fmt.Sprintf("- ~~%s~~ (%s öncelik)\n", g.Baslik, g.Oncelik)
			metin += fmt.Sprintf("  `ID: %s`\n", g.ID)
		}
	}

	return mcp.NewToolResultText(metin), nil
}

// OzetGoster sistem özetini gösterir
func (h *Handlers) OzetGoster(params map[string]interface{}) (*mcp.CallToolResult, error) {
	ozet, err := h.isYonetici.OzetAl()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("özet alınamadı: %v", err)), nil
	}

	metin := fmt.Sprintf(`## Özet Rapor

**Toplam Proje:** %d
**Toplam Görev:** %d

### Durum Dağılımı
- Beklemede: %d
- Devam Ediyor: %d
- Tamamlandı: %d

### Öncelik Dağılımı
- Yüksek: %d
- Orta: %d
- Düşük: %d`,
		ozet.ToplamProje,
		ozet.ToplamGorev,
		ozet.BeklemedeGorev,
		ozet.DevamEdenGorev,
		ozet.TamamlananGorev,
		ozet.YuksekOncelik,
		ozet.OrtaOncelik,
		ozet.DusukOncelik,
	)

	return mcp.NewToolResultText(metin), nil
}

// RegisterTools tüm araçları MCP sunucusuna kaydeder
func (h *Handlers) RegisterTools(s *server.MCPServer) {
	// Görev oluştur
	s.AddTool(mcp.Tool{
		Name:        "gorev_olustur",
		Description: "Yeni bir görev oluştur",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": "Görev başlığı",
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": "Görev açıklaması",
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": "Öncelik seviyesi (dusuk, orta, yuksek)",
					"enum":        []string{"dusuk", "orta", "yuksek"},
				},
			},
			Required: []string{"baslik"},
		},
	}, h.GorevOlustur)

	// Görev listele
	s.AddTool(mcp.Tool{
		Name:        "gorev_listele",
		Description: "Görevleri listele",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"durum": map[string]interface{}{
					"type":        "string",
					"description": "Filtrelenecek durum",
					"enum":        []string{"beklemede", "devam_ediyor", "tamamlandi"},
				},
			},
		},
	}, h.GorevListele)

	// Görev güncelle
	s.AddTool(mcp.Tool{
		Name:        "gorev_guncelle",
		Description: "Görev durumunu güncelle",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Görev ID",
				},
				"durum": map[string]interface{}{
					"type":        "string",
					"description": "Yeni durum",
					"enum":        []string{"beklemede", "devam_ediyor", "tamamlandi"},
				},
			},
			Required: []string{"id", "durum"},
		},
	}, h.GorevGuncelle)

	// Görev detay
	s.AddTool(mcp.Tool{
		Name:        "gorev_detay",
		Description: "Bir görevin detaylı bilgilerini markdown formatında göster",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Görev ID",
				},
			},
			Required: []string{"id"},
		},
	}, h.GorevDetay)

	// Görev düzenle
	s.AddTool(mcp.Tool{
		Name:        "gorev_duzenle",
		Description: "Bir görevin başlık, açıklama, öncelik veya proje bilgilerini düzenle",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Görev ID",
				},
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": "Yeni başlık (opsiyonel)",
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": "Yeni açıklama - markdown destekler (opsiyonel)",
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": "Yeni öncelik seviyesi (opsiyonel)",
					"enum":        []string{"dusuk", "orta", "yuksek"},
				},
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": "Yeni proje ID (opsiyonel)",
				},
			},
			Required: []string{"id"},
		},
	}, h.GorevDuzenle)

	// Görev sil
	s.AddTool(mcp.Tool{
		Name:        "gorev_sil",
		Description: "Bir görevi kalıcı olarak sil",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Görev ID",
				},
				"onay": map[string]interface{}{
					"type":        "boolean",
					"description": "Silme işlemini onaylamak için true olmalı",
				},
			},
			Required: []string{"id", "onay"},
		},
	}, h.GorevSil)

	// Proje oluştur
	s.AddTool(mcp.Tool{
		Name:        "proje_olustur",
		Description: "Yeni bir proje oluştur",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"isim": map[string]interface{}{
					"type":        "string",
					"description": "Proje ismi",
				},
				"tanim": map[string]interface{}{
					"type":        "string",
					"description": "Proje tanımı",
				},
			},
			Required: []string{"isim"},
		},
	}, h.ProjeOlustur)

	// Proje listele
	s.AddTool(mcp.Tool{
		Name:        "proje_listele",
		Description: "Tüm projeleri listele",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.ProjeListele)

	// Proje görevleri
	s.AddTool(mcp.Tool{
		Name:        "proje_gorevleri",
		Description: "Bir projenin görevlerini listele",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": "Proje ID",
				},
			},
			Required: []string{"proje_id"},
		},
	}, h.ProjeGorevleri)

	// Özet göster
	s.AddTool(mcp.Tool{
		Name:        "ozet_goster",
		Description: "Proje ve görev özetini göster",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.OzetGoster)
}
