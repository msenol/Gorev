package mcp

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/gorev"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Handlers struct {
	isYonetici        *gorev.IsYonetici
	aiContextYonetici *gorev.AIContextYonetici
}

func YeniHandlers(isYonetici *gorev.IsYonetici) *Handlers {
	// Create AI context manager using the same data manager
	aiContextYonetici := gorev.YeniAIContextYonetici(isYonetici.VeriYonetici())

	return &Handlers{
		isYonetici:        isYonetici,
		aiContextYonetici: aiContextYonetici,
	}
}

// gorevResponseSizeEstimate bir görev için tahmini response boyutunu hesaplar
func (h *Handlers) gorevResponseSizeEstimate(gorev *gorev.Gorev) int {
	// Tahmini karakter sayıları
	size := 100 // Temel formatlar için
	size += len(gorev.Baslik) + len(gorev.Aciklama)
	size += len(gorev.ID) + len(gorev.ProjeID)

	if gorev.SonTarih != nil {
		size += 30 // Tarih formatı için
	}

	for _, etiket := range gorev.Etiketler {
		size += len(etiket.Isim) + 5
	}

	// Bağımlılık bilgileri
	if gorev.BagimliGorevSayisi > 0 || gorev.BuGoreveBagimliSayisi > 0 {
		size += 100
	}

	return size
}

// gorevOzetYazdir bir görevi özet formatta yazdırır (ProjeGorevleri için)
func (h *Handlers) gorevOzetYazdir(g *gorev.Gorev) string {
	// Öncelik kısaltması
	oncelik := ""
	switch g.Oncelik {
	case "yuksek":
		oncelik = "Y"
	case "orta":
		oncelik = "O"
	case "dusuk":
		oncelik = "D"
	}

	metin := fmt.Sprintf("- **%s** (%s)", g.Baslik, oncelik)

	// Inline detaylar
	details := []string{}
	if g.Aciklama != "" && len(g.Aciklama) <= 50 {
		details = append(details, g.Aciklama)
	} else if g.Aciklama != "" {
		details = append(details, g.Aciklama[:47]+"...")
	}

	if g.SonTarih != nil {
		details = append(details, g.SonTarih.Format("02/01"))
	}

	if len(g.Etiketler) > 0 && len(g.Etiketler) <= 2 {
		etiketler := make([]string, len(g.Etiketler))
		for i, e := range g.Etiketler {
			etiketler[i] = e.Isim
		}
		details = append(details, strings.Join(etiketler, ","))
	} else if len(g.Etiketler) > 2 {
		details = append(details, fmt.Sprintf("%d etiket", len(g.Etiketler)))
	}

	if g.TamamlanmamisBagimlilikSayisi > 0 {
		details = append(details, fmt.Sprintf("🔒%d", g.TamamlanmamisBagimlilikSayisi))
	}

	details = append(details, g.ID[:8])

	if len(details) > 0 {
		metin += " - " + strings.Join(details, " | ")
	}
	metin += "\n"

	return metin
}

// gorevOzetYazdirTamamlandi tamamlanmış bir görevi özet formatta yazdırır
func (h *Handlers) gorevOzetYazdirTamamlandi(g *gorev.Gorev) string {
	// Çok kısa format - sadece başlık ve ID
	return fmt.Sprintf("- ~~%s~~ | %s\n", g.Baslik, g.ID[:8])
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

	// Öncelik kısaltması
	oncelikKisa := ""
	switch gorev.Oncelik {
	case "yuksek":
		oncelikKisa = "Y"
	case "orta":
		oncelikKisa = "O"
	case "dusuk":
		oncelikKisa = "D"
	default:
		oncelikKisa = gorev.Oncelik
	}

	// Temel satır - öncelik parantez içinde kısaltılmış
	metin := fmt.Sprintf("%s%s[%s] %s (%s)\n", indent, prefix, durum, gorev.Baslik, oncelikKisa)

	// Sadece dolu alanları göster, boş satırlar ekleme
	details := []string{}

	if gorev.Aciklama != "" {
		// Template sistemi için açıklama limiti büyük ölçüde artırıldı
		// Sadece gerçekten çok uzun açıklamaları kısalt (2000+ karakter)
		aciklama := gorev.Aciklama
		if len(aciklama) > 2000 {
			// İlk 1997 karakteri al ve ... ekle
			aciklama = aciklama[:1997] + "..."
		}
		details = append(details, aciklama)
	}

	if projeGoster && gorev.ProjeID != "" {
		proje, _ := h.isYonetici.ProjeGetir(gorev.ProjeID)
		if proje != nil {
			details = append(details, fmt.Sprintf("Proje: %s", proje.Isim))
		}
	}

	if gorev.SonTarih != nil {
		details = append(details, fmt.Sprintf("Tarih: %s", gorev.SonTarih.Format("02/01")))
	}

	if len(gorev.Etiketler) > 0 && len(gorev.Etiketler) <= 3 {
		etiketIsimleri := make([]string, len(gorev.Etiketler))
		for i, etiket := range gorev.Etiketler {
			etiketIsimleri[i] = etiket.Isim
		}
		details = append(details, fmt.Sprintf("Etiket: %s", strings.Join(etiketIsimleri, ",")))
	} else if len(gorev.Etiketler) > 3 {
		details = append(details, fmt.Sprintf("Etiket: %d adet", len(gorev.Etiketler)))
	}

	// Bağımlılık bilgileri - sadece varsa ve sıfırdan büyükse
	if gorev.TamamlanmamisBagimlilikSayisi > 0 {
		details = append(details, fmt.Sprintf("Bekleyen: %d", gorev.TamamlanmamisBagimlilikSayisi))
	}

	// ID'yi en sona ekle
	details = append(details, fmt.Sprintf("ID:%s", gorev.ID))

	// Detayları tek satırda göster
	if len(details) > 0 {
		metin += fmt.Sprintf("%s  %s\n", indent, strings.Join(details, " | "))
	}

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

// gorevHiyerarsiYazdirVeIsaretle görevleri yazdırırken hangi görevlerin gösterildiğini işaretler
func (h *Handlers) gorevHiyerarsiYazdirVeIsaretle(gorev *gorev.Gorev, gorevMap map[string]*gorev.Gorev, seviye int, projeGoster bool, shownGorevIDs map[string]bool) string {
	// Bu görevi gösterildi olarak işaretle
	shownGorevIDs[gorev.ID] = true

	// Normal hiyerarşik yazdırma işlemi
	metin := h.gorevHiyerarsiYazdirInternal(gorev, gorevMap, seviye, projeGoster, shownGorevIDs)

	return metin
}

// gorevHiyerarsiYazdirInternal görev hiyerarşisini yazdırır ve gösterilenleri işaretler
func (h *Handlers) gorevHiyerarsiYazdirInternal(gorev *gorev.Gorev, gorevMap map[string]*gorev.Gorev, seviye int, projeGoster bool, shownGorevIDs map[string]bool) string {
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

	// Öncelik kısaltması
	oncelikKisa := ""
	switch gorev.Oncelik {
	case "yuksek":
		oncelikKisa = "Y"
	case "orta":
		oncelikKisa = "O"
	case "dusuk":
		oncelikKisa = "D"
	default:
		oncelikKisa = gorev.Oncelik
	}

	// Temel satır - öncelik parantez içinde kısaltılmış
	metin := fmt.Sprintf("%s%s[%s] %s (%s)\n", indent, prefix, durum, gorev.Baslik, oncelikKisa)

	// Sadece dolu alanları göster, boş satırlar ekleme
	details := []string{}

	if gorev.Aciklama != "" {
		// Template sistemi için açıklama limiti büyük ölçüde artırıldı
		// Sadece gerçekten çok uzun açıklamaları kısalt (2000+ karakter)
		aciklama := gorev.Aciklama
		if len(aciklama) > 2000 {
			// İlk 1997 karakteri al ve ... ekle
			aciklama = aciklama[:1997] + "..."
		}
		details = append(details, aciklama)
	}

	if projeGoster && gorev.ProjeID != "" {
		proje, _ := h.isYonetici.ProjeGetir(gorev.ProjeID)
		if proje != nil {
			details = append(details, fmt.Sprintf("Proje: %s", proje.Isim))
		}
	}

	if gorev.SonTarih != nil {
		details = append(details, fmt.Sprintf("Tarih: %s", gorev.SonTarih.Format("02/01")))
	}

	if len(gorev.Etiketler) > 0 && len(gorev.Etiketler) <= 3 {
		etiketIsimleri := make([]string, len(gorev.Etiketler))
		for i, etiket := range gorev.Etiketler {
			etiketIsimleri[i] = etiket.Isim
		}
		details = append(details, fmt.Sprintf("Etiket: %s", strings.Join(etiketIsimleri, ", ")))
	} else if len(gorev.Etiketler) > 3 {
		details = append(details, fmt.Sprintf("Etiket: %d adet", len(gorev.Etiketler)))
	}

	// Bağımlılık bilgileri - sadece varsa ve sıfırdan büyükse
	if gorev.TamamlanmamisBagimlilikSayisi > 0 {
		details = append(details, fmt.Sprintf("Bekleyen: %d", gorev.TamamlanmamisBagimlilikSayisi))
	}

	// ID'yi en sona ekle
	details = append(details, fmt.Sprintf("ID:%s", gorev.ID))

	// Detayları tek satırda göster
	if len(details) > 0 {
		metin += fmt.Sprintf("%s  %s\n", indent, strings.Join(details, " | "))
	}

	// Alt görevleri bul ve yazdır - TÜM alt görevler gösterilir
	for _, g := range gorevMap {
		if g.ParentID == gorev.ID {
			shownGorevIDs[g.ID] = true
			metin += h.gorevHiyerarsiYazdirInternal(g, gorevMap, seviye+1, projeGoster, shownGorevIDs)
		}
	}

	if seviye == 0 {
		metin += "\n"
	}

	return metin
}

// templateZorunluAlanlariListele template'in zorunlu alanlarını listeler
func (h *Handlers) templateZorunluAlanlariListele(template *gorev.GorevTemplate) string {
	var alanlar []string
	for _, alan := range template.Alanlar {
		if alan.Zorunlu {
			tip := alan.Tip
			if alan.Tip == "select" && len(alan.Secenekler) > 0 {
				tip = fmt.Sprintf("select [%s]", strings.Join(alan.Secenekler, ", "))
			}
			alanlar = append(alanlar, fmt.Sprintf("- %s (%s)", alan.Isim, tip))
		}
	}
	return strings.Join(alanlar, "\n")
}

// templateOrnekDegerler template için örnek değerler oluşturur
func (h *Handlers) templateOrnekDegerler(template *gorev.GorevTemplate) string {
	var ornekler []string
	for _, alan := range template.Alanlar {
		if alan.Zorunlu {
			ornek := ""
			switch alan.Tip {
			case "select":
				if len(alan.Secenekler) > 0 {
					ornek = alan.Secenekler[0]
				}
			case "date":
				ornek = "2025-01-15"
			case "text":
				ornek = "örnek " + alan.Isim
			}
			ornekler = append(ornekler, fmt.Sprintf("'%s': '%s'", alan.Isim, ornek))
		}
	}
	return strings.Join(ornekler, ", ")
}

// GorevOlustur - DEPRECATED: Template kullanımı artık zorunludur
func (h *Handlers) GorevOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(`❌ gorev_olustur artık kullanılmıyor!

Template kullanımı zorunludur. Lütfen şu adımları takip edin:

1. Önce mevcut template'leri listeleyin:
   template_listele

2. Uygun template'i seçin ve görev oluşturun:
   templateden_gorev_olustur template_id='bug_report_v2' baslik='...' ...

Mevcut template kategorileri:
• 🐛 Bug: bug_report, bug_report_v2
• ✨ Feature: feature_request
• 🔬 Araştırma: research_task, spike_research
• ⚡ Performans: performance_issue
• 🔒 Güvenlik: security_fix
• ♻️ Teknik: technical_debt, refactoring

Detaylı bilgi için: template_listele kategori='Bug'`), nil
}

// GorevListele görevleri listeler
func (h *Handlers) GorevListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	durum, _ := params["durum"].(string)
	sirala, _ := params["sirala"].(string)
	filtre, _ := params["filtre"].(string)
	etiket, _ := params["etiket"].(string)
	tumProjeler, _ := params["tum_projeler"].(bool)

	// Pagination parametreleri
	limit := 50 // Varsayılan limit
	if l, ok := params["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}
	offset := 0
	if o, ok := params["offset"].(float64); ok && o >= 0 {
		offset = int(o)
	}

	// DEBUG: Log parametreleri
	// fmt.Fprintf(os.Stderr, "[GorevListele] Called - durum: %s, limit: %d, offset: %d\n", durum, limit, offset)

	gorevler, err := h.isYonetici.GorevListele(durum, sirala, filtre)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("görevler listelenemedi: %v", err)), nil
	}

	// DEBUG: Log görev sayısı
	// fmt.Fprintf(os.Stderr, "[GorevListele] Fetched %d tasks total\n", len(gorevler))

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

	// Toplam görev sayısı
	toplamGorevSayisi := len(gorevler)

	if toplamGorevSayisi == 0 {
		mesaj := "Henüz görev bulunmuyor."
		if aktifProje != nil {
			mesaj = fmt.Sprintf("%s projesinde henüz görev bulunmuyor.", aktifProje.Isim)
		}
		return mcp.NewToolResultText(mesaj), nil
	}

	// Görevleri hiyerarşik olarak organize et
	gorevMap := make(map[string]*gorev.Gorev)
	kokGorevler := []*gorev.Gorev{}

	for _, g := range gorevler {
		gorevMap[g.ID] = g
		if g.ParentID == "" {
			kokGorevler = append(kokGorevler, g)
		}
	}

	metin := ""

	// Kompakt başlık ve pagination bilgisi
	// NOT: Artık sadece root görev sayısını gösteriyoruz
	toplamRootGorevSayisi := len(kokGorevler)
	if toplamRootGorevSayisi > limit || offset > 0 {
		metin = fmt.Sprintf("Görevler (%d-%d / %d)\n",
			offset+1,
			min(offset+limit, toplamRootGorevSayisi),
			toplamRootGorevSayisi)
	} else {
		metin = fmt.Sprintf("Görevler (%d)\n", toplamRootGorevSayisi)
	}

	if aktifProje != nil && !tumProjeler {
		metin += fmt.Sprintf("Proje: %s\n", aktifProje.Isim)
	}
	metin += "\n"

	// Pagination uygula - SADECE ROOT görevlere
	// Subtask'lar parent'larıyla birlikte gösterilecek
	var paginatedKokGorevler []*gorev.Gorev
	if offset < len(kokGorevler) {
		end := offset + limit
		if end > len(kokGorevler) {
			end = len(kokGorevler)
		}
		paginatedKokGorevler = kokGorevler[offset:end]
	} else {
		paginatedKokGorevler = []*gorev.Gorev{}
	}

	// Response boyutunu tahmin et ve gerekirse daha az görev göster
	estimatedSize := 0
	const maxResponseSize = 20000 // ~20K karakter güvenli limit

	gorevlerToShow := []*gorev.Gorev{}
	for _, kokGorev := range paginatedKokGorevler {
		gorevSize := h.gorevResponseSizeEstimate(kokGorev)
		// Alt görevler için ek boyut tahmin et
		for _, g := range gorevMap {
			if g.ParentID == kokGorev.ID {
				gorevSize += h.gorevResponseSizeEstimate(g)
			}
		}

		if estimatedSize+gorevSize > maxResponseSize && len(gorevlerToShow) > 0 {
			// Boyut aşılacak, daha fazla görev ekleme
			metin += fmt.Sprintf("\n*Not: Response boyut limiti nedeniyle %d görev daha var. 'offset' parametresi ile devam edebilirsiniz.*\n",
				len(paginatedKokGorevler)-len(gorevlerToShow))
			break
		}
		estimatedSize += gorevSize
		gorevlerToShow = append(gorevlerToShow, kokGorev)
	}

	// Hangi görevlerin gösterildiğini takip et
	shownGorevIDs := make(map[string]bool)

	// Kök görevlerden başlayarak hiyerarşiyi oluştur
	// NOT: gorevMap tüm görevleri içerir, böylece paginated bir görevin TÜM alt görevleri gösterilir
	for _, kokGorev := range gorevlerToShow {
		metin += h.gorevHiyerarsiYazdirVeIsaretle(kokGorev, gorevMap, 0, tumProjeler || aktifProje == nil, shownGorevIDs)
	}
	
	// REMOVED: Orphan checking logic
	// Artık sadece root görevleri paginate ediyoruz
	// Alt görevler her zaman parent'larıyla birlikte gösterilecek

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

	// Auto-state management: Record task view and potentially transition state
	if err := h.aiContextYonetici.RecordTaskView(id); err != nil {
		// Log but don't fail the request
		// fmt.Printf("Görev görüntüleme kaydı hatası: %v\n", err)
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

	// Bağımlılıkları ekle - Her zaman göster
	metin += "\n\n## 🔗 Bağımlılıklar\n"

	baglantilar, err := h.isYonetici.GorevBaglantilariGetir(id)
	if err != nil {
		metin += "*Bağımlılık bilgileri alınamadı*\n"
	} else if len(baglantilar) == 0 {
		metin += "*Bu görevin herhangi bir bağımlılığı bulunmuyor*\n"
	} else {
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
		} else {
			metin += "\n### 📋 Bu görev için beklenen görevler:\n*Hiçbir göreve bağımlı değil*\n"
		}

		if len(sonrakiler) > 0 {
			metin += "\n### 🎯 Bu göreve bağımlı görevler:\n"
			for _, sonraki := range sonrakiler {
				metin += sonraki + "\n"
			}
		} else {
			metin += "\n### 🎯 Bu göreve bağımlı görevler:\n*Hiçbir görev bu göreve bağımlı değil*\n"
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
	}

	return mcp.NewToolResultText(metin), nil
}

// gorevBagimlilikBilgisi görev için bağımlılık bilgilerini formatlar
func (h *Handlers) gorevBagimlilikBilgisi(g *gorev.Gorev, indent string) string {
	bilgi := ""
	if g.BagimliGorevSayisi > 0 {
		bilgi += fmt.Sprintf("%sBağımlı görev sayısı: %d\n", indent, g.BagimliGorevSayisi)
		if g.TamamlanmamisBagimlilikSayisi > 0 {
			bilgi += fmt.Sprintf("%sTamamlanmamış bağımlılık sayısı: %d\n", indent, g.TamamlanmamisBagimlilikSayisi)
		}
	}
	if g.BuGoreveBagimliSayisi > 0 {
		bilgi += fmt.Sprintf("%sBu göreve bağımlı sayısı: %d\n", indent, g.BuGoreveBagimliSayisi)
	}
	return bilgi
}

// ProjeGorevleri bir projenin görevlerini listeler
func (h *Handlers) ProjeGorevleri(params map[string]interface{}) (*mcp.CallToolResult, error) {
	projeID, ok := params["proje_id"].(string)
	if !ok || projeID == "" {
		return mcp.NewToolResultError("proje_id parametresi gerekli"), nil
	}

	// Pagination parametreleri
	limit := 50 // Varsayılan limit
	if l, ok := params["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}
	offset := 0
	if o, ok := params["offset"].(float64); ok && o >= 0 {
		offset = int(o)
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

	// Toplam görev sayısı
	toplamGorevSayisi := len(gorevler)

	metin := ""

	if toplamGorevSayisi == 0 {
		metin = fmt.Sprintf("%s - Görev yok", proje.Isim)
		return mcp.NewToolResultText(metin), nil
	}

	// Kompakt başlık
	if toplamGorevSayisi > limit || offset > 0 {
		metin = fmt.Sprintf("%s (%d-%d / %d)\n",
			proje.Isim,
			offset+1,
			min(offset+limit, toplamGorevSayisi),
			toplamGorevSayisi)
	} else {
		metin = fmt.Sprintf("%s (%d görev)\n", proje.Isim, toplamGorevSayisi)
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

	// Pagination uygula - tüm görevleri tek bir listede topla
	allGorevler := append(append(devamEdiyor, beklemede...), tamamlandi...)

	// Pagination limitleri
	start := offset
	end := offset + limit
	if start > len(allGorevler) {
		start = len(allGorevler)
	}
	if end > len(allGorevler) {
		end = len(allGorevler)
	}

	// Response boyut kontrolü
	estimatedSize := len(metin)
	const maxResponseSize = 20000
	gorevleriGoster := 0

	// Önce devam eden görevleri göster
	devamEdiyorStart := 0
	devamEdiyorEnd := len(devamEdiyor)

	if start < len(devamEdiyor) {
		devamEdiyorStart = start
		if end < len(devamEdiyor) {
			devamEdiyorEnd = end
		}
		start = len(devamEdiyor)
	} else {
		devamEdiyorStart = len(devamEdiyor)
		devamEdiyorEnd = len(devamEdiyor)
		start -= len(devamEdiyor)
	}

	if devamEdiyorEnd > devamEdiyorStart {
		metin += "\n🔵 Devam Ediyor\n"
		for i := devamEdiyorStart; i < devamEdiyorEnd; i++ {
			g := devamEdiyor[i]
			gorevSize := h.gorevResponseSizeEstimate(g)
			if estimatedSize+gorevSize > maxResponseSize && gorevleriGoster > 0 {
				metin += fmt.Sprintf("*... ve %d görev daha (boyut limiti)*\n", devamEdiyorEnd-i)
				break
			}
			metin += h.gorevOzetYazdir(g)
			estimatedSize += gorevSize
			gorevleriGoster++
		}
	}

	// Bekleyen görevleri göster
	beklemedeStart := 0
	beklemedeEnd := len(beklemede)

	if start < len(devamEdiyor)+len(beklemede) {
		if start > len(devamEdiyor) {
			beklemedeStart = start - len(devamEdiyor)
		}
		if end < len(devamEdiyor)+len(beklemede) {
			beklemedeEnd = end - len(devamEdiyor)
			if beklemedeEnd < 0 {
				beklemedeEnd = 0
			}
		}
		start = len(devamEdiyor) + len(beklemede)
	} else {
		beklemedeStart = len(beklemede)
		beklemedeEnd = len(beklemede)
		start -= len(beklemede)
	}

	if beklemedeEnd > beklemedeStart && estimatedSize < maxResponseSize {
		metin += "\n⚪ Beklemede\n"
		for i := beklemedeStart; i < beklemedeEnd; i++ {
			g := beklemede[i]
			gorevSize := h.gorevResponseSizeEstimate(g)
			if estimatedSize+gorevSize > maxResponseSize && gorevleriGoster > 0 {
				metin += fmt.Sprintf("*... ve %d görev daha (boyut limiti)*\n", beklemedeEnd-i)
				break
			}
			metin += h.gorevOzetYazdir(g)
			estimatedSize += gorevSize
			gorevleriGoster++
		}
	}

	// Tamamlanan görevleri göster
	tamamlandiStart := 0
	tamamlandiEnd := len(tamamlandi)

	remainingOffset := offset - len(devamEdiyor) - len(beklemede)
	if remainingOffset > 0 && remainingOffset < len(tamamlandi) {
		tamamlandiStart = remainingOffset
	}

	remainingEnd := end - len(devamEdiyor) - len(beklemede)
	if remainingEnd < len(tamamlandi) && remainingEnd >= 0 {
		tamamlandiEnd = remainingEnd
	}

	if tamamlandiEnd > tamamlandiStart && estimatedSize < maxResponseSize {
		metin += "\n✅ Tamamlandı\n"
		for i := tamamlandiStart; i < tamamlandiEnd; i++ {
			g := tamamlandi[i]
			gorevSize := h.gorevResponseSizeEstimate(g)
			if estimatedSize+gorevSize > maxResponseSize && gorevleriGoster > 0 {
				metin += fmt.Sprintf("*... ve %d görev daha (boyut limiti)*\n", tamamlandiEnd-i)
				break
			}
			metin += h.gorevOzetYazdirTamamlandi(g)
			estimatedSize += gorevSize
			gorevleriGoster++
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

	// Önce template'i kontrol et
	template, err := h.isYonetici.VeriYonetici().TemplateGetir(templateID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("template bulunamadı: %v", err)), nil
	}

	// Interface{} map'i string map'e çevir ve validation yap
	degerler := make(map[string]string)
	eksikAlanlar := []string{}
	
	// Tüm zorunlu alanları kontrol et
	for _, alan := range template.Alanlar {
		if val, exists := degerlerRaw[alan.Isim]; exists {
			// Değer var, string'e çevir
			strVal := fmt.Sprintf("%v", val)
			if alan.Zorunlu && strings.TrimSpace(strVal) == "" {
				eksikAlanlar = append(eksikAlanlar, alan.Isim)
			} else {
				degerler[alan.Isim] = strVal
			}
		} else if alan.Zorunlu {
			// Zorunlu alan eksik
			eksikAlanlar = append(eksikAlanlar, alan.Isim)
		} else if alan.Varsayilan != "" {
			// Varsayılan değeri kullan
			degerler[alan.Isim] = alan.Varsayilan
		}
	}

	// Eksik alanlar varsa detaylı hata ver
	if len(eksikAlanlar) > 0 {
		return mcp.NewToolResultError(fmt.Sprintf(`❌ Zorunlu alanlar eksik!

Template: %s
Eksik alanlar: %s

Bu template için zorunlu alanlar:
%s

Örnek kullanım:
templateden_gorev_olustur template_id='%s' degerler={%s}`,
			template.Isim,
			strings.Join(eksikAlanlar, ", "),
			h.templateZorunluAlanlariListele(template),
			templateID,
			h.templateOrnekDegerler(template))), nil
	}

	// Select tipindeki alanların geçerli değerlerini kontrol et
	for _, alan := range template.Alanlar {
		if alan.Tip == "select" && len(alan.Secenekler) > 0 {
			if deger, ok := degerler[alan.Isim]; ok && deger != "" {
				gecerli := false
				for _, secenek := range alan.Secenekler {
					if deger == secenek {
						gecerli = true
						break
					}
				}
				if !gecerli {
					return mcp.NewToolResultError(fmt.Sprintf("'%s' alanı için geçersiz değer: '%s'. Geçerli değerler: %s", 
						alan.Isim, deger, strings.Join(alan.Secenekler, ", "))), nil
				}
			}
		}
	}

	gorev, err := h.isYonetici.TemplatedenGorevOlustur(templateID, degerler)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("template'den görev oluşturulamadı: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(`✓ Template kullanılarak görev oluşturuldu!

Template: %s
Başlık: %s
ID: %s
Öncelik: %s

Detaylar için: gorev_detay id='%s'`, 
		template.Isim, gorev.Baslik, gorev.ID, gorev.Oncelik, gorev.ID)), nil
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

// CallTool çağrı yapmak için yardımcı metod
func (h *Handlers) CallTool(toolName string, params map[string]interface{}) (*mcp.CallToolResult, error) {
	switch toolName {
	case "gorev_olustur":
		return h.GorevOlustur(params)
	case "gorev_listele":
		return h.GorevListele(params)
	case "gorev_detay":
		return h.GorevDetay(params)
	case "gorev_guncelle":
		return h.GorevGuncelle(params)
	case "gorev_duzenle":
		return h.GorevDuzenle(params)
	case "gorev_sil":
		return h.GorevSil(params)
	case "gorev_bagimlilik_ekle":
		return h.GorevBagimlilikEkle(params)
	case "gorev_altgorev_olustur":
		return h.GorevAltGorevOlustur(params)
	case "gorev_ust_degistir":
		return h.GorevUstDegistir(params)
	case "gorev_hiyerarsi_goster":
		return h.GorevHiyerarsiGoster(params)
	case "proje_olustur":
		return h.ProjeOlustur(params)
	case "proje_listele":
		return h.ProjeListele(params)
	case "proje_gorevleri":
		return h.ProjeGorevleri(params)
	case "proje_aktif_yap":
		return h.AktifProjeAyarla(params)
	case "aktif_proje_goster":
		return h.AktifProjeGoster(params)
	case "aktif_proje_kaldir":
		return h.AktifProjeKaldir(params)
	case "ozet_goster":
		return h.OzetGoster(params)
	case "template_listele":
		return h.TemplateListele(params)
	case "templateden_gorev_olustur":
		return h.TemplatedenGorevOlustur(params)
	case "gorev_set_active":
		return h.GorevSetActive(params)
	case "gorev_get_active":
		return h.GorevGetActive(params)
	case "gorev_recent":
		return h.GorevRecent(params)
	case "gorev_context_summary":
		return h.GorevContextSummary(params)
	case "gorev_batch_update":
		return h.GorevBatchUpdate(params)
	case "gorev_nlp_query":
		return h.GorevNLPQuery(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("bilinmeyen araç: %s", toolName)), nil
	}
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
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Gösterilecek maksimum görev sayısı. Varsayılan: 50",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Atlanacak görev sayısı (pagination için). Varsayılan: 0",
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
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Gösterilecek maksimum görev sayısı. Varsayılan: 50",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Atlanacak görev sayısı (pagination için). Varsayılan: 0",
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

	// AI Context Management Tools
	// Set active task
	s.AddTool(mcp.Tool{
		Name:        "gorev_set_active",
		Description: "AI oturumu için aktif görevi belirler. Görev otomatik olarak 'devam_ediyor' durumuna geçer.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": "Aktif olarak ayarlanacak görevin ID'si",
				},
			},
			Required: []string{"task_id"},
		},
	}, h.GorevSetActive)

	// Get active task
	s.AddTool(mcp.Tool{
		Name:        "gorev_get_active",
		Description: "AI oturumu için şu anda aktif olan görevi getirir",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.GorevGetActive)

	// Get recent tasks
	s.AddTool(mcp.Tool{
		Name:        "gorev_recent",
		Description: "AI'ın son etkileşimde bulunduğu görevleri listeler",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Döndürülecek görev sayısı (varsayılan: 5)",
				},
			},
		},
	}, h.GorevRecent)

	// Get context summary
	s.AddTool(mcp.Tool{
		Name:        "gorev_context_summary",
		Description: "AI için optimize edilmiş oturum özeti getirir (aktif görev, son görevler, öncelikler, blokajlar)",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.GorevContextSummary)

	// Batch update
	s.AddTool(mcp.Tool{
		Name:        "gorev_batch_update",
		Description: "Birden fazla görevi tek seferde günceller",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"updates": map[string]interface{}{
					"type":        "array",
					"description": "Güncelleme listesi [{id: string, updates: {durum?: string, oncelik?: string, ...}}]",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"id": map[string]interface{}{
								"type":        "string",
								"description": "Güncellenecek görevin ID'si",
							},
							"updates": map[string]interface{}{
								"type":        "object",
								"description": "Güncellenecek alanlar",
							},
						},
						"required": []string{"id", "updates"},
					},
				},
			},
			Required: []string{"updates"},
		},
	}, h.GorevBatchUpdate)

	// NLP Query
	s.AddTool(mcp.Tool{
		Name:        "gorev_nlp_query",
		Description: "Doğal dil sorgusu ile görev arama. Örnekler: 'bugün üzerinde çalıştığım görevler', 'yüksek öncelikli görevler', 'database ile ilgili görevler', 'son oluşturduğum görev', 'tamamlanmamış görevler', 'etiket:bug', 'acil görevler'",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Doğal dil sorgusu",
				},
			},
			Required: []string{"query"},
		},
	}, h.GorevNLPQuery)
}

// AI Context Management Handlers

// GorevSetActive sets the active task for the AI session
func (h *Handlers) GorevSetActive(params map[string]interface{}) (*mcp.CallToolResult, error) {
	taskID, ok := params["task_id"].(string)
	if !ok || taskID == "" {
		return mcp.NewToolResultError("task_id parametresi gerekli"), nil
	}

	err := h.aiContextYonetici.SetActiveTask(taskID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Aktif görev ayarlama hatası: %v", err)), nil
	}

	// Also record task view for auto-state management
	if err := h.aiContextYonetici.RecordTaskView(taskID); err != nil {
		// Log but don't fail
		// fmt.Printf("Görev görüntüleme kaydı hatası: %v\n", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("✅ Görev %s başarıyla aktif görev olarak ayarlandı.", taskID)), nil
}

// GorevGetActive returns the current active task
func (h *Handlers) GorevGetActive(params map[string]interface{}) (*mcp.CallToolResult, error) {
	activeTask, err := h.aiContextYonetici.GetActiveTask()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Aktif görev getirme hatası: %v", err)), nil
	}

	if activeTask == nil {
		return mcp.NewToolResultText("ℹ️ Şu anda aktif görev yok."), nil
	}

	// Format active task details
	metin := fmt.Sprintf(`# Aktif Görev: %s

## 📋 Genel Bilgiler
- **ID:** %s
- **Durum:** %s
- **Öncelik:** %s
- **Proje:** %s`,
		activeTask.Baslik,
		activeTask.ID,
		activeTask.Durum,
		activeTask.Oncelik,
		activeTask.ProjeID)

	if activeTask.Aciklama != "" {
		metin += fmt.Sprintf("\n\n## 📝 Açıklama\n%s", activeTask.Aciklama)
	}

	return mcp.NewToolResultText(metin), nil
}

// GorevRecent returns recent tasks interacted with by AI
func (h *Handlers) GorevRecent(params map[string]interface{}) (*mcp.CallToolResult, error) {
	limit := 5
	if l, ok := params["limit"].(float64); ok {
		limit = int(l)
	}

	tasks, err := h.aiContextYonetici.GetRecentTasks(limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Son görevleri getirme hatası: %v", err)), nil
	}

	if len(tasks) == 0 {
		return mcp.NewToolResultText("ℹ️ Son etkileşimde bulunulan görev yok."), nil
	}

	var result strings.Builder
	result.WriteString("## 📋 Son Etkileşimli Görevler\n\n")

	for i, task := range tasks {
		result.WriteString(fmt.Sprintf("### %d. %s (ID: %s)\n", i+1, task.Baslik, task.ID))
		result.WriteString(fmt.Sprintf("- **Durum:** %s\n", task.Durum))
		result.WriteString(fmt.Sprintf("- **Öncelik:** %s\n", task.Oncelik))
		if task.ProjeID != "" {
			result.WriteString(fmt.Sprintf("- **Proje:** %s\n", task.ProjeID))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

// GorevContextSummary returns an AI-optimized context summary
func (h *Handlers) GorevContextSummary(params map[string]interface{}) (*mcp.CallToolResult, error) {
	summary, err := h.aiContextYonetici.GetContextSummary()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Context özeti getirme hatası: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString("## 🤖 AI Oturum Özeti\n\n")

	// Active task
	if summary.ActiveTask != nil {
		result.WriteString(fmt.Sprintf("### 🎯 Aktif Görev\n**%s** (ID: %s)\n", summary.ActiveTask.Baslik, summary.ActiveTask.ID))
		result.WriteString(fmt.Sprintf("- Durum: %s | Öncelik: %s\n\n", summary.ActiveTask.Durum, summary.ActiveTask.Oncelik))
	} else {
		result.WriteString("### 🎯 Aktif Görev\nYok\n\n")
	}

	// Working project
	if summary.WorkingProject != nil {
		result.WriteString(fmt.Sprintf("### 📁 Çalışılan Proje\n**%s**\n\n", summary.WorkingProject.Isim))
	}

	// Session summary
	result.WriteString("### 📊 Oturum İstatistikleri\n")
	result.WriteString(fmt.Sprintf("- Oluşturulan: %d\n", summary.SessionSummary.Created))
	result.WriteString(fmt.Sprintf("- Güncellenen: %d\n", summary.SessionSummary.Updated))
	result.WriteString(fmt.Sprintf("- Tamamlanan: %d\n\n", summary.SessionSummary.Completed))

	// Next priorities
	if len(summary.NextPriorities) > 0 {
		result.WriteString("### 🔥 Öncelikli Görevler\n")
		for _, task := range summary.NextPriorities {
			result.WriteString(fmt.Sprintf("- **%s** (ID: %s)\n", task.Baslik, task.ID))
		}
		result.WriteString("\n")
	}

	// Blockers
	if len(summary.Blockers) > 0 {
		result.WriteString("### 🚫 Blokajlar\n")
		for _, task := range summary.Blockers {
			result.WriteString(fmt.Sprintf("- **%s** (ID: %s) - %d bağımlılık bekliyor\n",
				task.Baslik, task.ID, task.TamamlanmamisBagimlilikSayisi))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

// GorevBatchUpdate performs batch updates on multiple tasks
func (h *Handlers) GorevBatchUpdate(params map[string]interface{}) (*mcp.CallToolResult, error) {
	updatesRaw, ok := params["updates"].([]interface{})
	if !ok {
		return mcp.NewToolResultError("updates parametresi gerekli ve dizi olmalı"), nil
	}

	var updates []gorev.BatchUpdate
	for _, u := range updatesRaw {
		updateMap, ok := u.(map[string]interface{})
		if !ok {
			continue
		}

		id, ok := updateMap["id"].(string)
		if !ok || id == "" {
			continue
		}

		updatesData, ok := updateMap["updates"].(map[string]interface{})
		if !ok {
			continue
		}

		updates = append(updates, gorev.BatchUpdate{
			ID:      id,
			Updates: updatesData,
		})
	}

	if len(updates) == 0 {
		return mcp.NewToolResultError("Geçerli güncelleme bulunamadı"), nil
	}

	result, err := h.aiContextYonetici.BatchUpdate(updates)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Toplu güncelleme hatası: %v", err)), nil
	}

	var response strings.Builder
	response.WriteString("## 📦 Toplu Güncelleme Sonucu\n\n")
	response.WriteString(fmt.Sprintf("**Toplam İşlenen:** %d\n", result.TotalProcessed))
	response.WriteString(fmt.Sprintf("**Başarılı:** %d\n", len(result.Successful)))
	response.WriteString(fmt.Sprintf("**Başarısız:** %d\n\n", len(result.Failed)))

	if len(result.Successful) > 0 {
		response.WriteString("### ✅ Başarılı Güncellemeler\n")
		for _, id := range result.Successful {
			response.WriteString(fmt.Sprintf("- %s\n", id))
		}
		response.WriteString("\n")
	}

	if len(result.Failed) > 0 {
		response.WriteString("### ❌ Başarısız Güncellemeler\n")
		for _, fail := range result.Failed {
			response.WriteString(fmt.Sprintf("- %s: %s\n", fail.ID, fail.Error))
		}
	}

	return mcp.NewToolResultText(response.String()), nil
}

// GorevNLPQuery performs natural language query on tasks
func (h *Handlers) GorevNLPQuery(params map[string]interface{}) (*mcp.CallToolResult, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query parametresi gerekli"), nil
	}

	tasks, err := h.aiContextYonetici.NLPQuery(query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Sorgu hatası: %v", err)), nil
	}

	if len(tasks) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("ℹ️ '%s' sorgusu için sonuç bulunamadı.", query)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("## 🔍 Sorgu Sonuçları: '%s'\n\n", query))
	result.WriteString(fmt.Sprintf("**%d görev bulundu**\n\n", len(tasks)))

	// Use compact format for results
	for _, task := range tasks {
		statusEmoji := "⏳"
		if task.Durum == "devam_ediyor" {
			statusEmoji = "🔄"
		} else if task.Durum == "tamamlandi" {
			statusEmoji = "✅"
		}

		priorityEmoji := "ℹ️"
		if task.Oncelik == "yuksek" {
			priorityEmoji = "🔥"
		} else if task.Oncelik == "orta" {
			priorityEmoji = "⚡"
		}

		result.WriteString(fmt.Sprintf("%s %s **%s** (ID: %s)\n", statusEmoji, priorityEmoji, task.Baslik, task.ID[:8]))

		if task.Aciklama != "" {
			desc := task.Aciklama
			if len(desc) > 100 {
				desc = desc[:100] + "..."
			}
			result.WriteString(fmt.Sprintf("   %s\n", desc))
		}

		details := []string{}
		if task.ProjeID != "" {
			details = append(details, fmt.Sprintf("Proje: %s", task.ProjeID))
		}
		if len(task.Etiketler) > 0 {
			var tagNames []string
			for _, tag := range task.Etiketler {
				tagNames = append(tagNames, tag.Isim)
			}
			details = append(details, fmt.Sprintf("Etiketler: %s", strings.Join(tagNames, ", ")))
		}
		if task.SonTarih != nil {
			details = append(details, fmt.Sprintf("Son tarih: %s", task.SonTarih.Format("2006-01-02")))
		}

		if len(details) > 0 {
			result.WriteString(fmt.Sprintf("   %s\n", strings.Join(details, " | ")))
		}

		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}
