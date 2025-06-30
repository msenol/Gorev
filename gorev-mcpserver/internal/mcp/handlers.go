package mcp

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/gorev"
)

type Handlers struct {
	isYonetici *gorev.IsYonetici
}

func YeniHandlers(isYonetici *gorev.IsYonetici) *Handlers {
	return &Handlers{
		isYonetici: isYonetici,
	}
}

// gorevHiyerarsiYazdir bir görevi ve alt görevlerini hiyerarşik olarak yazdırır
func (h *Handlers) gorevHiyerarsiYazdir(gorev *gorev.Gorev, gorevMap map[string]*gorev.Gorev, seviye int, projeGoster bool) string {
	indent := strings.Repeat("  ", seviye)
	prefix := ""
	if seviye > 0 {
		prefix = "└─ "
	}

	durum := ""
	switch gorev.Durum {
	case "tamamlandi":
		durum = "✓"
	case "devam_ediyor":
		durum = "🔄"
	case "beklemede":
		durum = "⏳"
	}

	metin := fmt.Sprintf("%s%s[%s] %s (%s öncelik)\n", indent, prefix, durum, gorev.Baslik, gorev.Oncelik)

	if gorev.Aciklama != "" {
		metin += fmt.Sprintf("%s  %s\n", indent, gorev.Aciklama)
	}

	if projeGoster && gorev.ProjeID != "" {
		proje, _ := h.isYonetici.ProjeGetir(gorev.ProjeID)
		if proje != nil {
			metin += fmt.Sprintf("%s  Proje: %s\n", indent, proje.Isim)
		}
		metin += fmt.Sprintf("%s  ProjeID: %s\n", indent, gorev.ProjeID)
	}

	if gorev.SonTarih != nil {
		metin += fmt.Sprintf("%s  Son tarih: %s\n", indent, gorev.SonTarih.Format("2006-01-02"))
	}

	if len(gorev.Etiketler) > 0 {
		etiketIsimleri := make([]string, len(gorev.Etiketler))
		for i, etiket := range gorev.Etiketler {
			etiketIsimleri[i] = etiket.Isim
		}
		metin += fmt.Sprintf("%s  Etiketler: %s\n", indent, strings.Join(etiketIsimleri, ", "))
	}

	metin += fmt.Sprintf("%s  ID: %s\n", indent, gorev.ID)

	// Alt görevleri bul ve yazdır
	for _, g := range gorevMap {
		if g.ParentID == gorev.ID {
			metin += h.gorevHiyerarsiYazdir(g, gorevMap, seviye+1, projeGoster)
		}
	}

	if seviye == 0 {
		metin += "\n"
	}

	return metin
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

	projeID, _ := params["proje_id"].(string)
	sonTarih, _ := params["son_tarih"].(string)
	etiketlerStr, _ := params["etiketler"].(string)
	etiketler := strings.Split(etiketlerStr, ",")

	// Eğer proje_id verilmemişse, aktif projeyi kullan
	if projeID == "" {
		aktifProje, err := h.isYonetici.AktifProjeGetir()
		if err == nil && aktifProje != nil {
			projeID = aktifProje.ID
		}
	}

	gorev, err := h.isYonetici.GorevOlustur(baslik, aciklama, oncelik, projeID, sonTarih, etiketler)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görev oluşturulamadı: %v", err)), nil
	}

	mesaj := fmt.Sprintf("✓ Görev oluşturuldu: %s (ID: %s)", gorev.Baslik, gorev.ID)
	if projeID != "" {
		proje, _ := h.isYonetici.ProjeGetir(projeID)
		if proje != nil {
			mesaj += fmt.Sprintf("\n  Proje: %s", proje.Isim)
		}
	}

	return mcp.NewToolResultText(mesaj), nil
}

// GorevListele görevleri listeler
func (h *Handlers) GorevListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	durum, _ := params["durum"].(string)
	sirala, _ := params["sirala"].(string)
	filtre, _ := params["filtre"].(string)
	etiket, _ := params["etiket"].(string)
	tumProjeler, _ := params["tum_projeler"].(bool)

	gorevler, err := h.isYonetici.GorevListele(durum, sirala, filtre)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görevler listelenemedi: %v", err)), nil
	}

	// Etikete göre filtrele
	if etiket != "" {
		var filtreliGorevler []*gorev.Gorev
		for _, g := range gorevler {
			for _, e := range g.Etiketler {
				if e.Isim == etiket {
					filtreliGorevler = append(filtreliGorevler, g)
					break
				}
			}
		}
		gorevler = filtreliGorevler
	}

	// Aktif proje varsa ve tum_projeler false ise, sadece aktif projenin görevlerini göster
	var aktifProje *gorev.Proje
	if !tumProjeler {
		aktifProje, _ = h.isYonetici.AktifProjeGetir()
		if aktifProje != nil {
			// Görevleri filtrele
			var filtreliGorevler []*gorev.Gorev
			for _, g := range gorevler {
				if g.ProjeID == aktifProje.ID {
					filtreliGorevler = append(filtreliGorevler, g)
				}
			}
			gorevler = filtreliGorevler
		}
	}

	if len(gorevler) == 0 {
		mesaj := "Henüz görev bulunmuyor."
		if aktifProje != nil {
			mesaj = fmt.Sprintf("%s projesinde henüz görev bulunmuyor.", aktifProje.Isim)
		}
		return mcp.NewToolResultText(mesaj), nil
	}

	metin := "## Görev Listesi"
	if aktifProje != nil && !tumProjeler {
		metin += fmt.Sprintf(" - %s", aktifProje.Isim)
	}
	metin += "\n\n"

	// Görevleri hiyerarşik olarak organize et
	gorevMap := make(map[string]*gorev.Gorev)
	kokGorevler := []*gorev.Gorev{}

	for _, g := range gorevler {
		gorevMap[g.ID] = g
		if g.ParentID == "" {
			kokGorevler = append(kokGorevler, g)
		}
	}

	// Kök görevlerden başlayarak hiyerarşiyi oluştur
	for _, kokGorev := range kokGorevler {
		metin += h.gorevHiyerarsiYazdir(kokGorev, gorevMap, 0, aktifProje == nil)
	}

	// Parent'ı olmayan ama parent_id'si dolu olanları da göster (parent görünmeyen görevler)
	for _, g := range gorevler {
		if g.ParentID != "" {
			if _, parentVar := gorevMap[g.ParentID]; !parentVar {
				metin += h.gorevHiyerarsiYazdir(g, gorevMap, 0, aktifProje == nil)
			}
		}
	}

	return mcp.NewToolResultText(metin), nil
}

// AktifProjeAyarla bir projeyi aktif proje olarak ayarlar
func (h *Handlers) AktifProjeAyarla(params map[string]interface{}) (*mcp.CallToolResult, error) {
	projeID, ok := params["proje_id"].(string)
	if !ok || projeID == "" {
		return mcp.NewToolResultError("proje_id parametresi gerekli"), nil
	}

	if err := h.isYonetici.AktifProjeAyarla(projeID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("aktif proje ayarlanamadı: %v", err)), nil
	}

	proje, _ := h.isYonetici.ProjeGetir(projeID)
	if proje != nil {
		return mcp.NewToolResultText(
			fmt.Sprintf("✓ Aktif proje ayarlandı: %s", proje.Isim),
		), nil
	}
	return mcp.NewToolResultText(
		fmt.Sprintf("✓ Aktif proje ayarlandı: %s", projeID),
	), nil
}

// AktifProjeGoster mevcut aktif projeyi gösterir
func (h *Handlers) AktifProjeGoster(params map[string]interface{}) (*mcp.CallToolResult, error) {
	proje, err := h.isYonetici.AktifProjeGetir()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("aktif proje getirilemedi: %v", err)), nil
	}

	if proje == nil {
		return mcp.NewToolResultText("Henüz aktif proje ayarlanmamış."), nil
	}

	// Görev sayısını al
	gorevSayisi, _ := h.isYonetici.ProjeGorevSayisi(proje.ID)

	metin := fmt.Sprintf(`## Aktif Proje

**Proje:** %s
**ID:** %s
**Açıklama:** %s
**Görev Sayısı:** %d`,
		proje.Isim,
		proje.ID,
		proje.Tanim,
		gorevSayisi,
	)

	return mcp.NewToolResultText(metin), nil
}

// AktifProjeKaldir aktif proje ayarını kaldırır
func (h *Handlers) AktifProjeKaldir(params map[string]interface{}) (*mcp.CallToolResult, error) {
	if err := h.isYonetici.AktifProjeKaldir(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("aktif proje kaldırılamadı: %v", err)), nil
	}

	return mcp.NewToolResultText("✓ Aktif proje ayarı kaldırıldı."), nil
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

	gorev, err := h.isYonetici.GorevGetir(id)
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
		proje, err := h.isYonetici.ProjeGetir(gorev.ProjeID)
		if err == nil {
			metin += fmt.Sprintf("\n- **Proje:** %s", proje.Isim)
		}
	}
	if gorev.SonTarih != nil {
		metin += fmt.Sprintf("\n- **Son Teslim Tarihi:** %s", gorev.SonTarih.Format("2006-01-02"))
	}
	if len(gorev.Etiketler) > 0 {
		var etiketIsimleri []string
		for _, e := range gorev.Etiketler {
			etiketIsimleri = append(etiketIsimleri, e.Isim)
		}
		metin += fmt.Sprintf("\n- **Etiketler:** %s", strings.Join(etiketIsimleri, ", "))
	}

	metin += "\n\n## 📝 Açıklama\n"
	if gorev.Aciklama != "" {
		// Açıklama zaten markdown formatında olabilir, direkt ekle
		metin += gorev.Aciklama
	} else {
		metin += "*Açıklama girilmemiş*"
	}

	// Bağımlılıkları ekle
	baglantilar, err := h.isYonetici.GorevBaglantilariGetir(id)
	if err == nil && len(baglantilar) > 0 {
		metin += "\n\n## 🔗 Bağımlılıklar\n"

		var oncekiler []string
		var sonrakiler []string

		for _, b := range baglantilar {
			if b.BaglantiTip == "onceki" {
				if b.HedefID == id {
					// Bu görev hedefse, kaynak önceki görevdir
					kaynakGorev, err := h.isYonetici.GorevGetir(b.KaynakID)
					if err == nil {
						durum := "✅"
						if kaynakGorev.Durum != "tamamlandi" {
							durum = "⏳"
						}
						oncekiler = append(oncekiler, fmt.Sprintf("%s %s (`%s`)", durum, kaynakGorev.Baslik, kaynakGorev.Durum))
					}
				} else if b.KaynakID == id {
					// Bu görev kaynaksa, hedef sonraki görevdir
					hedefGorev, err := h.isYonetici.GorevGetir(b.HedefID)
					if err == nil {
						sonrakiler = append(sonrakiler, fmt.Sprintf("- %s (`%s`)", hedefGorev.Baslik, hedefGorev.Durum))
					}
				}
			}
		}

		if len(oncekiler) > 0 {
			metin += "\n### 📋 Bu görev için beklenen görevler:\n"
			for _, onceki := range oncekiler {
				metin += fmt.Sprintf("- %s\n", onceki)
			}
		}

		if len(sonrakiler) > 0 {
			metin += "\n### 🎯 Bu göreve bağımlı görevler:\n"
			for _, sonraki := range sonrakiler {
				metin += sonraki + "\n"
			}
		}

		// Bağımlılık durumu kontrolü
		bagimli, tamamlanmamislar, err := h.isYonetici.GorevBagimliMi(id)
		if err == nil && !bagimli && gorev.Durum == "beklemede" {
			metin += fmt.Sprintf("\n> ⚠️ **Uyarı:** Bu görev başlatılamaz! Önce şu görevler tamamlanmalı: %v\n", tamamlanmamislar)
		}
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
	sonTarih, sonTarihVar := params["son_tarih"].(string)

	if !baslikVar && !aciklamaVar && !oncelikVar && !projeVar && !sonTarihVar {
		return mcp.NewToolResultError("en az bir düzenleme alanı belirtilmeli (baslik, aciklama, oncelik, proje_id veya son_tarih)"), nil
	}

	err := h.isYonetici.GorevDuzenle(id, baslik, aciklama, oncelik, projeID, sonTarih, baslikVar, aciklamaVar, oncelikVar, projeVar, sonTarihVar)
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

	gorev, err := h.isYonetici.GorevGetir(id)
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
	proje, err := h.isYonetici.ProjeGetir(projeID)
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
			metin += fmt.Sprintf("  ProjeID: %s\n", g.ProjeID)
			if g.SonTarih != nil {
				metin += fmt.Sprintf("  Son tarih: %s\n", g.SonTarih.Format("2006-01-02"))
			}
			if len(g.Etiketler) > 0 {
				etiketIsimleri := make([]string, len(g.Etiketler))
				for i, etiket := range g.Etiketler {
					etiketIsimleri[i] = etiket.Isim
				}
				metin += fmt.Sprintf("  Etiketler: %s\n", strings.Join(etiketIsimleri, ", "))
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
			metin += fmt.Sprintf("  ProjeID: %s\n", g.ProjeID)
			if g.SonTarih != nil {
				metin += fmt.Sprintf("  Son tarih: %s\n", g.SonTarih.Format("2006-01-02"))
			}
			if len(g.Etiketler) > 0 {
				etiketIsimleri := make([]string, len(g.Etiketler))
				for i, etiket := range g.Etiketler {
					etiketIsimleri[i] = etiket.Isim
				}
				metin += fmt.Sprintf("  Etiketler: %s\n", strings.Join(etiketIsimleri, ", "))
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
			metin += fmt.Sprintf("  ProjeID: %s\n", g.ProjeID)
			if len(g.Etiketler) > 0 {
				etiketIsimleri := make([]string, len(g.Etiketler))
				for i, etiket := range g.Etiketler {
					etiketIsimleri[i] = etiket.Isim
				}
				metin += fmt.Sprintf("  Etiketler: %s\n", strings.Join(etiketIsimleri, ", "))
			}
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

func (h *Handlers) GorevBagimlilikEkle(params map[string]interface{}) (*mcp.CallToolResult, error) {
	kaynakID, ok := params["kaynak_id"].(string)
	if !ok || kaynakID == "" {
		return mcp.NewToolResultError("kaynak_id parametresi gerekli"), nil
	}

	hedefID, ok := params["hedef_id"].(string)
	if !ok || hedefID == "" {
		return mcp.NewToolResultError("hedef_id parametresi gerekli"), nil
	}

	baglantiTipi, ok := params["baglanti_tipi"].(string)
	if !ok || baglantiTipi == "" {
		return mcp.NewToolResultError("baglanti_tipi parametresi gerekli"), nil
	}

	baglanti, err := h.isYonetici.GorevBagimlilikEkle(kaynakID, hedefID, baglantiTipi)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("bağımlılık eklenemedi: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✓ Bağımlılık eklendi: %s -> %s (%s)", baglanti.KaynakID, baglanti.HedefID, baglanti.BaglantiTip)), nil
}

// TemplateListele kullanılabilir template'leri listeler
func (h *Handlers) TemplateListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	kategori, _ := params["kategori"].(string)

	templates, err := h.isYonetici.TemplateListele(kategori)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("template'ler listelenemedi: %v", err)), nil
	}

	if len(templates) == 0 {
		return mcp.NewToolResultText("Henüz template bulunmuyor."), nil
	}

	metin := "## 📋 Görev Template'leri\n\n"

	// Kategorilere göre grupla
	kategoriMap := make(map[string][]*gorev.GorevTemplate)
	for _, tmpl := range templates {
		kategoriMap[tmpl.Kategori] = append(kategoriMap[tmpl.Kategori], tmpl)
	}

	// Her kategoriyi göster
	for kat, tmpls := range kategoriMap {
		metin += fmt.Sprintf("### %s\n\n", kat)

		for _, tmpl := range tmpls {
			metin += fmt.Sprintf("#### %s\n", tmpl.Isim)
			metin += fmt.Sprintf("- **ID:** `%s`\n", tmpl.ID)
			metin += fmt.Sprintf("- **Açıklama:** %s\n", tmpl.Tanim)
			metin += fmt.Sprintf("- **Başlık Şablonu:** `%s`\n", tmpl.VarsayilanBaslik)

			// Alanları göster
			if len(tmpl.Alanlar) > 0 {
				metin += "- **Alanlar:**\n"
				for _, alan := range tmpl.Alanlar {
					zorunlu := ""
					if alan.Zorunlu {
						zorunlu = " *(zorunlu)*"
					}
					metin += fmt.Sprintf("  - `%s` (%s)%s", alan.Isim, alan.Tip, zorunlu)
					if alan.Varsayilan != "" {
						metin += fmt.Sprintf(" - varsayılan: %s", alan.Varsayilan)
					}
					if len(alan.Secenekler) > 0 {
						metin += fmt.Sprintf(" - seçenekler: %s", strings.Join(alan.Secenekler, ", "))
					}
					metin += "\n"
				}
			}
			metin += "\n"
		}
	}

	metin += "\n💡 **Kullanım:** `templateden_gorev_olustur` komutunu template ID'si ve alan değerleriyle kullanın."

	return mcp.NewToolResultText(metin), nil
}

// TemplatedenGorevOlustur template kullanarak görev oluşturur
func (h *Handlers) TemplatedenGorevOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	templateID, ok := params["template_id"].(string)
	if !ok || templateID == "" {
		return mcp.NewToolResultError("template_id parametresi gerekli"), nil
	}

	degerlerRaw, ok := params["degerler"].(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("degerler parametresi gerekli ve obje tipinde olmalı"), nil
	}

	// Interface{} map'i string map'e çevir
	degerler := make(map[string]string)
	for k, v := range degerlerRaw {
		degerler[k] = fmt.Sprintf("%v", v)
	}

	gorev, err := h.isYonetici.TemplatedenGorevOlustur(templateID, degerler)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("template'den görev oluşturulamadı: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✓ Template kullanılarak görev oluşturuldu: %s (ID: %s)", gorev.Baslik, gorev.ID)), nil
}

// RegisterTools tüm araçları MCP sunucusuna kaydeder
// GorevAltGorevOlustur mevcut bir görevin altına yeni görev oluşturur
func (h *Handlers) GorevAltGorevOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	parentID, ok := params["parent_id"].(string)
	if !ok || parentID == "" {
		return mcp.NewToolResultError("parent_id parametresi gerekli"), nil
	}

	baslik, ok := params["baslik"].(string)
	if !ok || baslik == "" {
		return mcp.NewToolResultError("başlık parametresi gerekli"), nil
	}

	aciklama, _ := params["aciklama"].(string)
	oncelik, _ := params["oncelik"].(string)
	if oncelik == "" {
		oncelik = "orta"
	}

	sonTarih, _ := params["son_tarih"].(string)
	etiketlerStr, _ := params["etiketler"].(string)
	var etiketler []string
	if etiketlerStr != "" {
		etiketler = strings.Split(etiketlerStr, ",")
		for i := range etiketler {
			etiketler[i] = strings.TrimSpace(etiketler[i])
		}
	}

	gorev, err := h.isYonetici.AltGorevOlustur(parentID, baslik, aciklama, oncelik, sonTarih, etiketler)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("alt görev oluşturulamadı: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("✓ Alt görev oluşturuldu: %s (ID: %s)", gorev.Baslik, gorev.ID)), nil
}

// GorevUstDegistir bir görevin üst görevini değiştirir
func (h *Handlers) GorevUstDegistir(params map[string]interface{}) (*mcp.CallToolResult, error) {
	gorevID, ok := params["gorev_id"].(string)
	if !ok || gorevID == "" {
		return mcp.NewToolResultError("gorev_id parametresi gerekli"), nil
	}

	yeniParentID, _ := params["yeni_parent_id"].(string)

	err := h.isYonetici.GorevUstDegistir(gorevID, yeniParentID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("üst görev değiştirilemedi: %v", err)), nil
	}

	if yeniParentID == "" {
		return mcp.NewToolResultText(fmt.Sprintf("✓ Görev kök seviyeye taşındı")), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("✓ Görev yeni üst göreve taşındı")), nil
}

// GorevHiyerarsiGoster bir görevin tam hiyerarşisini gösterir
func (h *Handlers) GorevHiyerarsiGoster(params map[string]interface{}) (*mcp.CallToolResult, error) {
	gorevID, ok := params["gorev_id"].(string)
	if !ok || gorevID == "" {
		return mcp.NewToolResultError("gorev_id parametresi gerekli"), nil
	}

	hiyerarsi, err := h.isYonetici.GorevHiyerarsiGetir(gorevID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("hiyerarşi alınamadı: %v", err)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# 📊 Görev Hiyerarşisi: %s\n\n", hiyerarsi.Gorev.Baslik))

	// Üst görevler
	if len(hiyerarsi.UstGorevler) > 0 {
		sb.WriteString("## 📍 Üst Görevler\n")
		for i := len(hiyerarsi.UstGorevler) - 1; i >= 0; i-- {
			ust := hiyerarsi.UstGorevler[i]
			sb.WriteString(fmt.Sprintf("%s└─ %s (%s)\n", strings.Repeat("  ", len(hiyerarsi.UstGorevler)-i-1), ust.Baslik, ust.Durum))
		}
		sb.WriteString(fmt.Sprintf("%s└─ **%s** (Mevcut Görev)\n\n", strings.Repeat("  ", len(hiyerarsi.UstGorevler)), hiyerarsi.Gorev.Baslik))
	}

	// Alt görev istatistikleri
	sb.WriteString("## 📈 Alt Görev İstatistikleri\n")
	sb.WriteString(fmt.Sprintf("- **Toplam Alt Görev:** %d\n", hiyerarsi.ToplamAltGorev))
	sb.WriteString(fmt.Sprintf("- **Tamamlanan:** %d ✓\n", hiyerarsi.TamamlananAlt))
	sb.WriteString(fmt.Sprintf("- **Devam Eden:** %d 🔄\n", hiyerarsi.DevamEdenAlt))
	sb.WriteString(fmt.Sprintf("- **Beklemede:** %d ⏳\n", hiyerarsi.BeklemedeAlt))
	sb.WriteString(fmt.Sprintf("- **İlerleme:** %.1f%%\n\n", hiyerarsi.IlerlemeYuzdesi))

	// Doğrudan alt görevler
	altGorevler, err := h.isYonetici.AltGorevleriGetir(gorevID)
	if err == nil && len(altGorevler) > 0 {
		sb.WriteString("## 🌳 Doğrudan Alt Görevler\n")
		for _, alt := range altGorevler {
			durum := ""
			switch alt.Durum {
			case "tamamlandi":
				durum = "✓"
			case "devam_ediyor":
				durum = "🔄"
			case "beklemede":
				durum = "⏳"
			}
			sb.WriteString(fmt.Sprintf("- %s %s (ID: %s)\n", durum, alt.Baslik, alt.ID))
		}
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func (h *Handlers) RegisterTools(s *server.MCPServer) {
	// Görev oluştur
	s.AddTool(mcp.Tool{
		Name:        "gorev_olustur",
		Description: "Kullanıcının doğal dil isteğinden bir görev oluşturur. Başlık, açıklama ve öncelik gibi bilgileri akıllıca çıkarır. Örneğin, kullanıcı 'çok acil olarak sunucu çökmesini düzeltmem lazım' derse, başlığı 'Sunucu çökmesini düzelt' ve önceliği 'yuksek' olarak ayarla. Eğer bir proje aktif ise görevi o projeye ata.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": "Görevin başlığı. Kullanıcının isteğindeki ana eylemden çıkarılmalıdır.",
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": "Görevin detaylı açıklaması. Kullanıcının isteğindeki ek bağlam veya detayları içerir.",
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": "Öncelik seviyesi. 'acil', 'önemli' gibi kelimelerden 'yuksek', 'düşük öncelikli' gibi ifadelerden 'dusuk' olarak çıkarım yapılmalıdır. Varsayılan 'orta'dır.",
					"enum":        []string{"dusuk", "orta", "yuksek"},
				},
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin atanacağı projenin ID'si. Kullanıcı belirtmezse ve aktif bir proje varsa, o kullanılır.",
				},
				"son_tarih": map[string]interface{}{
					"type":        "string",
					"description": "Görevin son teslim tarihi (YYYY-AA-GG formatında).",
				},
				"etiketler": map[string]interface{}{
					"type":        "string",
					"description": "Virgülle ayrılmış etiket listesi (örn: 'bug,acil,onemli').",
				},
			},
			Required: []string{"baslik"},
		},
	}, h.GorevOlustur)

	// Görev listele
	s.AddTool(mcp.Tool{
		Name:        "gorev_listele",
		Description: "Görevleri durum, proje, son teslim tarihi gibi kriterlere göre filtreleyerek ve sıralayarak listeler.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"durum": map[string]interface{}{
					"type":        "string",
					"description": "Filtrelenecek görev durumu (beklemede, devam_ediyor, tamamlandi).",
					"enum":        []string{"beklemede", "devam_ediyor", "tamamlandi"},
				},
				"sirala": map[string]interface{}{
					"type":        "string",
					"description": "Sıralama ölçütü ('son_tarih_asc', 'son_tarih_desc'). Varsayılan oluşturma tarihine göredir.",
					"enum":        []string{"son_tarih_asc", "son_tarih_desc"},
				},
				"filtre": map[string]interface{}{
					"type":        "string",
					"description": "Özel filtreler ('acil' - son 7 gün, 'gecmis' - tarihi geçmiş).",
					"enum":        []string{"acil", "gecmis"},
				},
				"etiket": map[string]interface{}{
					"type":        "string",
					"description": "Belirtilen etikete sahip görevleri filtreler.",
				},
				"tum_projeler": map[string]interface{}{
					"type":        "boolean",
					"description": "Tüm projelerdeki görevleri gösterir. Varsayılan olarak sadece aktif projenin görevleri listelenir.",
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
		Description: "Mevcut bir görevin başlık, açıklama, öncelik veya proje bilgilerini günceller. Kullanıcının isteğinden hangi alanların güncelleneceğini anlar. Örneğin, '123 ID'li görevin başlığını 'Yeni Başlık' yap' komutunu işler.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Düzenlenecek görevin ID'si.",
				},
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": "Görev için yeni başlık (opsiyonel).",
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": "Görev için yeni açıklama (opsiyonel).",
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": "Görev için yeni öncelik seviyesi (opsiyonel).",
					"enum":        []string{"dusuk", "orta", "yuksek"},
				},
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin atanacağı yeni projenin ID'si (opsiyonel).",
				},
				"son_tarih": map[string]interface{}{
					"type":        "string",
					"description": "Görevin yeni son teslim tarihi (YYYY-AA-GG formatında, boş string tarihi kaldırır).",
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

	// Proje aktif yap
	s.AddTool(mcp.Tool{
		Name:        "proje_aktif_yap",
		Description: "Bir projeyi aktif proje olarak ayarla",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": "Aktif yapılacak proje ID",
				},
			},
			Required: []string{"proje_id"},
		},
	}, h.AktifProjeAyarla)

	// Aktif proje göster
	s.AddTool(mcp.Tool{
		Name:        "aktif_proje_goster",
		Description: "Mevcut aktif projeyi göster",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.AktifProjeGoster)

	// Aktif proje kaldır
	s.AddTool(mcp.Tool{
		Name:        "aktif_proje_kaldir",
		Description: "Aktif proje ayarını kaldır",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.AktifProjeKaldir)

	// Görev bağımlılık ekle
	s.AddTool(mcp.Tool{
		Name:        "gorev_bagimlilik_ekle",
		Description: "İki görev arasına bir bağımlılık ekler",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"kaynak_id": map[string]interface{}{
					"type":        "string",
					"description": "Bağımlılığın kaynağı olan görev ID",
				},
				"hedef_id": map[string]interface{}{
					"type":        "string",
					"description": "Bağımlılığın hedefi olan görev ID",
				},
				"baglanti_tipi": map[string]interface{}{
					"type":        "string",
					"description": "Bağımlılık tipi (örn: 'engelliyor', 'ilişkili')",
				},
			},
			Required: []string{"kaynak_id", "hedef_id", "baglanti_tipi"},
		},
	}, h.GorevBagimlilikEkle)

	// Template listele
	s.AddTool(mcp.Tool{
		Name:        "template_listele",
		Description: "Kullanılabilir görev template'lerini listeler",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"kategori": map[string]interface{}{
					"type":        "string",
					"description": "Filtrelenecek template kategorisi (Teknik, Özellik, Araştırma vb.)",
				},
			},
		},
	}, h.TemplateListele)

	// Template'den görev oluştur
	s.AddTool(mcp.Tool{
		Name:        "templateden_gorev_olustur",
		Description: "Seçilen template'i kullanarak özelleştirilmiş bir görev oluşturur",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"template_id": map[string]interface{}{
					"type":        "string",
					"description": "Kullanılacak template'in ID'si",
				},
				"degerler": map[string]interface{}{
					"type":        "object",
					"description": "Template alanları için değerler (key-value çiftleri)",
				},
			},
			Required: []string{"template_id", "degerler"},
		},
	}, h.TemplatedenGorevOlustur)

	// Alt görev oluştur
	s.AddTool(mcp.Tool{
		Name:        "gorev_altgorev_olustur",
		Description: "Mevcut bir görevin altına yeni görev oluşturur",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"parent_id": map[string]interface{}{
					"type":        "string",
					"description": "Üst görevin ID'si",
				},
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": "Alt görevin başlığı",
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": "Alt görevin açıklaması",
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": "Öncelik seviyesi (yuksek, orta, dusuk)",
				},
				"son_tarih": map[string]interface{}{
					"type":        "string",
					"description": "Son tarih (YYYY-AA-GG formatında)",
				},
				"etiketler": map[string]interface{}{
					"type":        "string",
					"description": "Virgülle ayrılmış etiket listesi",
				},
			},
			Required: []string{"parent_id", "baslik"},
		},
	}, h.GorevAltGorevOlustur)

	// Görev üst değiştir
	s.AddTool(mcp.Tool{
		Name:        "gorev_ust_degistir",
		Description: "Bir görevin üst görevini değiştirir veya kök göreve taşır",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"gorev_id": map[string]interface{}{
					"type":        "string",
					"description": "Taşınacak görevin ID'si",
				},
				"yeni_parent_id": map[string]interface{}{
					"type":        "string",
					"description": "Yeni üst görevin ID'si (boş string kök göreve taşır)",
				},
			},
			Required: []string{"gorev_id"},
		},
	}, h.GorevUstDegistir)

	// Görev hiyerarşi göster
	s.AddTool(mcp.Tool{
		Name:        "gorev_hiyerarsi_goster",
		Description: "Bir görevin tam hiyerarşisini ve alt görev istatistiklerini gösterir",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"gorev_id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin ID'si",
				},
			},
			Required: []string{"gorev_id"},
		},
	}, h.GorevHiyerarsiGoster)
}
