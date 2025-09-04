package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
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
	fileWatcher       *gorev.FileWatcher
	toolHelpers       *ToolHelpers
	debug             bool
}

func YeniHandlers(isYonetici *gorev.IsYonetici) *Handlers {
	var aiContextYonetici *gorev.AIContextYonetici
	var fileWatcher *gorev.FileWatcher

	// Create AI context manager using the same data manager if isYonetici is not nil
	if isYonetici != nil {
		aiContextYonetici = gorev.YeniAIContextYonetici(isYonetici.VeriYonetici())

		// Initialize file watcher with default configuration
		if fw, err := gorev.NewFileWatcher(isYonetici.VeriYonetici(), gorev.DefaultFileWatcherConfig()); err == nil {
			fileWatcher = fw
		} else {
			// Log error but continue without file watcher
			fmt.Printf("Warning: Failed to initialize file watcher: %v\n", err)
		}
	}

	// Initialize tool helpers with shared utilities
	toolHelpers := NewToolHelpers()

	return &Handlers{
		isYonetici:        isYonetici,
		aiContextYonetici: aiContextYonetici,
		fileWatcher:       fileWatcher,
		toolHelpers:       toolHelpers,
		debug:             false,
	}
}

// YeniHandlersWithDebug creates handlers with debug support
func YeniHandlersWithDebug(isYonetici *gorev.IsYonetici, debug bool) *Handlers {
	var aiContextYonetici *gorev.AIContextYonetici
	var fileWatcher *gorev.FileWatcher

	// Create AI context manager using the same data manager if isYonetici is not nil
	if isYonetici != nil {
		aiContextYonetici = gorev.YeniAIContextYonetici(isYonetici.VeriYonetici())

		// Initialize file watcher with default configuration
		if fw, err := gorev.NewFileWatcher(isYonetici.VeriYonetici(), gorev.DefaultFileWatcherConfig()); err == nil {
			fileWatcher = fw
		} else if debug {
			// In debug mode, log the error
			fmt.Printf("DEBUG: Failed to initialize file watcher: %v\n", err)
		}
	}

	// Initialize tool helpers with shared utilities
	toolHelpers := NewToolHelpers()

	return &Handlers{
		isYonetici:        isYonetici,
		aiContextYonetici: aiContextYonetici,
		fileWatcher:       fileWatcher,
		toolHelpers:       toolHelpers,
		debug:             debug,
	}
}

// Close cleans up resources used by handlers
func (h *Handlers) Close() error {
	if h.fileWatcher != nil {
		return h.fileWatcher.Close()
	}
	return nil
}

// gorevResponseSizeEstimate bir görev için tahmini response boyutunu hesaplar
func (h *Handlers) gorevResponseSizeEstimate(gorev *gorev.Gorev) int {
	// Tahmini karakter sayıları
	size := constants.BaseResponseSize // Temel formatlar için
	size += len(gorev.Baslik) + len(gorev.Aciklama)
	size += len(gorev.ID) + len(gorev.ProjeID)

	if gorev.SonTarih != nil {
		size += constants.DateFormatSize // Tarih formatı için
	}

	for _, etiket := range gorev.Etiketler {
		size += len(etiket.Isim) + constants.TagSizeConstant
	}

	// Bağımlılık bilgileri
	if gorev.BagimliGorevSayisi > 0 || gorev.BuGoreveBagimliSayisi > 0 {
		size += constants.DependencyInfoSize
	}

	return size
}

// gorevOzetYazdir bir görevi özet formatta yazdırır (ProjeGorevleri için)
func (h *Handlers) gorevOzetYazdir(g *gorev.Gorev) string {
	// Öncelik kısaltması
	oncelik := PriorityFormat.GetPriorityShort(g.Oncelik)

	metin := fmt.Sprintf("- **%s** (%s)", g.Baslik, oncelik)

	// Inline detaylar
	details := []string{}
	if g.Aciklama != "" && len(g.Aciklama) <= constants.MaxInlineDescriptionLength {
		details = append(details, g.Aciklama)
	} else if g.Aciklama != "" {
		details = append(details, g.Aciklama[:constants.TruncatedDescriptionLength]+"...")
	}

	if g.SonTarih != nil {
		details = append(details, g.SonTarih.Format(constants.DateFormatShort))
	}

	if len(g.Etiketler) > 0 && len(g.Etiketler) <= 2 {
		etiketler := make([]string, len(g.Etiketler))
		for i, e := range g.Etiketler {
			etiketler[i] = e.Isim
		}
		details = append(details, strings.Join(etiketler, ","))
	} else if len(g.Etiketler) > 2 {
		details = append(details, i18n.T("messages.tagCount", map[string]interface{}{"Count": len(g.Etiketler)}))
	}

	if g.TamamlanmamisBagimlilikSayisi > 0 {
		details = append(details, fmt.Sprintf("🔒%d", g.TamamlanmamisBagimlilikSayisi))
	}

	details = append(details, TaskIDFormat.FormatShortID(g.ID))

	if len(details) > 0 {
		metin += " - " + strings.Join(details, " | ")
	}
	metin += "\n"

	return metin
}

// gorevOzetYazdirTamamlandi tamamlanmış bir görevi özet formatta yazdırır
func (h *Handlers) gorevOzetYazdirTamamlandi(g *gorev.Gorev) string {
	// Çok kısa format - sadece başlık ve ID
	return fmt.Sprintf("- ~~%s~~ | %s\n", g.Baslik, TaskIDFormat.FormatShortID(g.ID))
}

// gorevHiyerarsiYazdir bir görevi ve alt görevlerini hiyerarşik olarak yazdırır
func (h *Handlers) gorevHiyerarsiYazdir(gorev *gorev.Gorev, gorevMap map[string]*gorev.Gorev, seviye int, projeGoster bool) string {
	indent := strings.Repeat("  ", seviye)
	prefix := ""
	if seviye > 0 {
		prefix = "└─ "
	}

	durum := StatusFormat.GetStatusSymbol(gorev.Durum)

	// Öncelik kısaltması
	oncelikKisa := PriorityFormat.GetPriorityShort(gorev.Oncelik)

	// Temel satır - öncelik parantez içinde kısaltılmış
	metin := fmt.Sprintf("%s%s[%s] %s (%s)\n", indent, prefix, durum, gorev.Baslik, oncelikKisa)

	// Sadece dolu alanları göster, boş satırlar ekleme
	details := []string{}

	if gorev.Aciklama != "" {
		// Template sistemi için açıklama limiti büyük ölçüde artırıldı
		// Sadece gerçekten çok uzun açıklamaları kısalt
		aciklama := gorev.Aciklama
		if len(aciklama) > constants.MaxDescriptionLength {
			// İlk karakterleri al ve ... ekle
			aciklama = aciklama[:constants.MaxDescriptionTruncateLength] + "..."
		}
		details = append(details, aciklama)
	}

	if projeGoster && gorev.ProjeID != "" {
		proje, _ := h.isYonetici.ProjeGetir(gorev.ProjeID)
		if proje != nil {
			details = append(details, i18n.T("messages.projectLabel", map[string]interface{}{"Name": proje.Isim}))
		}
	}

	if gorev.SonTarih != nil {
		details = append(details, i18n.T("messages.dateLabel", map[string]interface{}{"Date": gorev.SonTarih.Format(constants.DateFormatShort)}))
	}

	if len(gorev.Etiketler) > 0 && len(gorev.Etiketler) <= constants.MaxTagsToDisplay {
		etiketIsimleri := make([]string, len(gorev.Etiketler))
		for i, etiket := range gorev.Etiketler {
			etiketIsimleri[i] = etiket.Isim
		}
		details = append(details, i18n.T("messages.tagLabel", map[string]interface{}{"Tags": strings.Join(etiketIsimleri, ",")}))
	} else if len(gorev.Etiketler) > constants.MaxTagsToDisplay {
		details = append(details, i18n.T("messages.tagCountLabel", map[string]interface{}{"Count": len(gorev.Etiketler)}))
	}

	// Bağımlılık bilgileri - sadece varsa ve sıfırdan büyükse
	if gorev.TamamlanmamisBagimlilikSayisi > 0 {
		details = append(details, i18n.TMarkdownLabel("bekleyen", gorev.TamamlanmamisBagimlilikSayisi))
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

	durum := StatusFormat.GetStatusSymbol(gorev.Durum)

	// Öncelik kısaltması
	oncelikKisa := PriorityFormat.GetPriorityShort(gorev.Oncelik)

	// Temel satır - öncelik parantez içinde kısaltılmış
	metin := fmt.Sprintf("%s%s[%s] %s (%s)\n", indent, prefix, durum, gorev.Baslik, oncelikKisa)

	// Sadece dolu alanları göster, boş satırlar ekleme
	details := []string{}

	if gorev.Aciklama != "" {
		// Template sistemi için açıklama limiti büyük ölçüde artırıldı
		// Sadece gerçekten çok uzun açıklamaları kısalt
		aciklama := gorev.Aciklama
		if len(aciklama) > constants.MaxDescriptionLength {
			// İlk karakterleri al ve ... ekle
			aciklama = aciklama[:constants.MaxDescriptionTruncateLength] + "..."
		}
		details = append(details, aciklama)
	}

	if projeGoster && gorev.ProjeID != "" {
		proje, _ := h.isYonetici.ProjeGetir(gorev.ProjeID)
		if proje != nil {
			details = append(details, i18n.T("messages.projectLabel", map[string]interface{}{"Name": proje.Isim}))
		}
	}

	if gorev.SonTarih != nil {
		details = append(details, i18n.T("messages.dateLabel", map[string]interface{}{"Date": gorev.SonTarih.Format(constants.DateFormatShort)}))
	}

	if len(gorev.Etiketler) > 0 && len(gorev.Etiketler) <= constants.MaxTagsToDisplay {
		etiketIsimleri := make([]string, len(gorev.Etiketler))
		for i, etiket := range gorev.Etiketler {
			etiketIsimleri[i] = etiket.Isim
		}
		details = append(details, i18n.TMarkdownLabel("etiket", strings.Join(etiketIsimleri, ", ")))
	} else if len(gorev.Etiketler) > constants.MaxTagsToDisplay {
		details = append(details, i18n.T("messages.tagCountLabel", map[string]interface{}{"Count": len(gorev.Etiketler)}))
	}

	// Bağımlılık bilgileri - sadece varsa ve sıfırdan büyükse
	if gorev.TamamlanmamisBagimlilikSayisi > 0 {
		details = append(details, i18n.TMarkdownLabel("bekleyen", gorev.TamamlanmamisBagimlilikSayisi))
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
				ornek = time.Now().AddDate(0, 0, 7).Format(constants.DateFormatISO) // One week from now
			case "text":
				ornek = "örnek " + alan.Isim
			}
			ornekler = append(ornekler, fmt.Sprintf("'%s': '%s'", alan.Isim, ornek))
		}
	}
	return strings.Join(ornekler, ", ")
}

// GorevOlustur - DEPRECATED: Template kullanımı artık zorunludur
// GorevOlustur was removed in v0.11.1 - use TemplatedenGorevOlustur instead

// GorevListele görevleri listeler
func (h *Handlers) GorevListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	durum, _ := params[constants.ParamDurum].(string)
	sirala, _ := params["sirala"].(string)
	filtre, _ := params["filtre"].(string)
	etiket, _ := params["etiket"].(string)
	tumProjeler, _ := params["tum_projeler"].(bool)

	// Pagination parametreleri
	limit, offset := h.toolHelpers.Validator.ValidatePagination(params)

	// DEBUG: Log parametreleri
	// fmt.Fprintf(os.Stderr, "[GorevListele] Called - durum: %s, limit: %d, offset: %d\n", durum, limit, offset)

	// Create filters map
	filters := make(map[string]interface{})
	if durum != "" {
		filters[constants.ParamDurum] = durum
	}
	if sirala != "" {
		filters["sirala"] = sirala
	}
	if filtre != "" {
		filters["filtre"] = filtre
	}

	gorevler, err := h.isYonetici.GorevListele(filters)
	if err != nil {
		return mcp.NewToolResultError(i18n.TListFailed("task", err)), nil
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
		mesaj := i18n.T("messages.noTasks")
		if aktifProje != nil {
			mesaj = i18n.T("messages.noTasksInProject", map[string]interface{}{"Project": aktifProje.Isim})
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
		metin = i18n.T("messages.taskListCount", map[string]interface{}{"Count": toplamRootGorevSayisi}) + "\n"
	}

	if aktifProje != nil && !tumProjeler {
		metin += i18n.T("messages.projectHeader", map[string]interface{}{"Name": aktifProje.Isim}) + "\n"
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
	maxResponseSize := constants.MaxResponseSize // ~20K karakter güvenli limit

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
			metin += "\n" + i18n.T("messages.sizeWarning", map[string]interface{}{"Count": len(paginatedKokGorevler) - len(gorevlerToShow)}) + "\n"
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
	projeID, result := h.toolHelpers.Validator.ValidateRequiredString(params, "proje_id")
	if result != nil {
		return result, nil
	}

	if err := h.isYonetici.AktifProjeAyarla(projeID); err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("aktif proje ayarlama başarısız", err), nil
	}

	proje, _ := h.isYonetici.ProjeGetir(projeID)
	if proje != nil {
		return mcp.NewToolResultText(
			i18n.T("success.activeProjectSet", map[string]interface{}{"Project": proje.Isim}),
		), nil
	}
	return mcp.NewToolResultText(
		i18n.T("success.activeProjectSet", map[string]interface{}{"Project": projeID}),
	), nil
}

// AktifProjeGoster mevcut aktif projeyi gösterir
func (h *Handlers) AktifProjeGoster(params map[string]interface{}) (*mcp.CallToolResult, error) {
	proje, err := h.isYonetici.AktifProjeGetir()
	if err != nil {
		return mcp.NewToolResultError(i18n.T("error.activeProjectRetrieve", map[string]interface{}{"Error": err.Error()})), nil
	}

	if proje == nil {
		return mcp.NewToolResultText(i18n.T("messages.noActiveProject")), nil
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

	return mcp.NewToolResultText(i18n.T("success.activeProjectRemoved")), nil
}

// GorevGuncelle görev durumunu günceller
func (h *Handlers) GorevGuncelle(params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Use helper for validation
	id, result := h.toolHelpers.Validator.ValidateTaskID(params)
	if result != nil {
		return result, nil
	}

	durum, result := h.toolHelpers.Validator.ValidateTaskStatus(params, true)
	if result != nil {
		return result, nil
	}

	if err := h.isYonetici.GorevDurumGuncelle(id, durum); err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("görev güncellenemedi", err), nil
	}

	return mcp.NewToolResultText(
		h.toolHelpers.Formatter.FormatSuccessMessage("Görev güncellendi", id, durum),
	), nil
}

// ProjeOlustur yeni bir proje oluşturur
func (h *Handlers) ProjeOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	isim, result := h.toolHelpers.Validator.ValidateRequiredString(params, "isim")
	if result != nil {
		return result, nil
	}

	tanim := h.toolHelpers.Validator.ValidateOptionalString(params, "tanim")

	proje, err := h.isYonetici.ProjeOlustur(isim, tanim)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("proje oluşturulamadı", err), nil
	}

	return mcp.NewToolResultText(
		h.toolHelpers.Formatter.FormatSuccessMessage("Proje oluşturuldu", proje.Isim, proje.ID),
	), nil
}

// GorevDetay tek bir görevin detaylı bilgisini markdown formatında döner
func (h *Handlers) GorevDetay(params map[string]interface{}) (*mcp.CallToolResult, error) {
	id, result := h.toolHelpers.Validator.ValidateTaskID(params)
	if result != nil {
		return result, nil
	}

	gorev, err := h.isYonetici.GorevGetir(id)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatNotFoundError("görev", id), nil
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
		gorev.OlusturmaTarih.Format(constants.DateTimeFormatFull),
		gorev.GuncellemeTarih.Format(constants.DateTimeFormatFull),
	)

	if gorev.ProjeID != "" {
		proje, err := h.isYonetici.ProjeGetir(gorev.ProjeID)
		if err == nil {
			metin += "\n" + i18n.TListItem("proje", proje.Isim)
		}
	}
	if gorev.ParentID != "" {
		parent, err := h.isYonetici.GorevGetir(gorev.ParentID)
		if err == nil {
			metin += fmt.Sprintf("\n- **Üst Görev:** %s", parent.Baslik)
		}
	}
	if gorev.SonTarih != nil {
		metin += "\n" + i18n.TListItem("son_tarih", gorev.SonTarih.Format(constants.DateFormatISO))
	}
	if len(gorev.Etiketler) > 0 {
		var etiketIsimleri []string
		for _, e := range gorev.Etiketler {
			etiketIsimleri = append(etiketIsimleri, e.Isim)
		}
		metin += "\n" + i18n.TListItem("etiketler", strings.Join(etiketIsimleri, ", "))
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
						durum := constants.EmojiStatusCompleted
						if kaynakGorev.Durum != constants.TaskStatusCompleted {
							durum = constants.EmojiStatusPending
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
		if err == nil && !bagimli && gorev.Durum == constants.TaskStatusPending {
			metin += fmt.Sprintf("\n> ⚠️ **Uyarı:** Bu görev başlatılamaz! Önce şu görevler tamamlanmalı: %v\n", tamamlanmamislar)
		}
	}

	metin += "\n\n---\n"
	metin += fmt.Sprintf("\n*Son güncelleme: %s*", gorev.GuncellemeTarih.Format(constants.DateFormatDisplay))

	return mcp.NewToolResultText(metin), nil
}

// GorevDuzenle görevi düzenler
func (h *Handlers) GorevDuzenle(params map[string]interface{}) (*mcp.CallToolResult, error) {
	id, ok := params[constants.ParamID].(string)
	if !ok || id == "" {
		return mcp.NewToolResultError("id parametresi gerekli"), nil
	}

	// En az bir düzenleme alanı olmalı
	baslik, baslikVar := params[constants.ParamBaslik].(string)
	aciklama, aciklamaVar := params[constants.ParamAciklama].(string)
	oncelik, oncelikVar := params[constants.ParamOncelik].(string)
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
	id, result := h.toolHelpers.Validator.ValidateTaskID(params)
	if result != nil {
		return result, nil
	}

	// Onay kontrolü
	onay := h.toolHelpers.Validator.ValidateBool(params, "onay")
	if !onay {
		return mcp.NewToolResultError("görevi silmek için 'onay' parametresi true olmalıdır"), nil
	}

	gorev, err := h.isYonetici.GorevGetir(id)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatNotFoundError("görev", id), nil
	}

	gorevBaslik := gorev.Baslik

	err = h.isYonetici.GorevSil(id)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("görev silinemedi", err), nil
	}

	return mcp.NewToolResultText(h.toolHelpers.Formatter.FormatSuccessMessage("Görev silindi", gorevBaslik, id)), nil
}

// GorevBulkTransition changes status for multiple tasks
func (h *Handlers) GorevBulkTransition(params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Validate task IDs
	taskIDsRaw, ok := params["task_ids"]
	if !ok {
		return mcp.NewToolResultError("task_ids parametresi gerekli"), nil
	}

	taskIDsInterface, ok := taskIDsRaw.([]interface{})
	if !ok {
		return mcp.NewToolResultError("task_ids array formatında olmalı"), nil
	}

	taskIDs := make([]string, len(taskIDsInterface))
	for i, idInterface := range taskIDsInterface {
		if id, ok := idInterface.(string); ok && id != "" {
			taskIDs[i] = id
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("geçersiz task ID index %d", i)), nil
		}
	}

	// Validate new status
	newStatus, result := h.toolHelpers.Validator.ValidateTaskStatus(params, true)
	if result != nil {
		return result, nil
	}

	// Optional parameters
	force := h.toolHelpers.Validator.ValidateBool(params, "force")
	checkDependencies := h.toolHelpers.Validator.ValidateBool(params, "check_dependencies")
	dryRun := h.toolHelpers.Validator.ValidateBool(params, "dry_run")

	// Create batch processor
	batchProcessor := gorev.NewBatchProcessor(h.isYonetici.VeriYonetici())
	if h.aiContextYonetici != nil {
		batchProcessor.SetAIContextManager(h.aiContextYonetici)
	}

	// Execute bulk transition
	request := gorev.BulkStatusTransitionRequest{
		TaskIDs:           taskIDs,
		NewStatus:         newStatus,
		Force:             force,
		CheckDependencies: checkDependencies,
		DryRun:            dryRun,
	}

	result_batch, err := batchProcessor.BulkStatusTransition(request)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("toplu durum değişikliği başarısız", err), nil
	}

	// Format response
	var response strings.Builder

	if dryRun {
		response.WriteString("🔍 **Kuru Çalıştırma Sonucu**\n\n")
	} else {
		response.WriteString(constants.EmojiStatusCompleted + " **Toplu Durum Değişikliği Tamamlandı**\n\n")
	}

	response.WriteString(i18n.TMarkdownLabel("hedef_durum", newStatus) + "\n")
	response.WriteString(i18n.TCount("islenen_gorev", result_batch.TotalProcessed) + "\n")
	response.WriteString(i18n.TCount("basarili", len(result_batch.Successful)) + "\n")
	response.WriteString(i18n.TCount("basarisiz", len(result_batch.Failed)) + "\n")
	response.WriteString(i18n.TCount("uyari", len(result_batch.Warnings)) + "\n")
	response.WriteString(i18n.TDuration("sure", result_batch.ExecutionTime) + "\n\n")

	if len(result_batch.Successful) > 0 {
		response.WriteString("**" + constants.EmojiStatusCompleted + " Başarılı Görevler:**\n")
		for _, taskID := range result_batch.Successful {
			response.WriteString(fmt.Sprintf("- %s\n", TaskIDFormat.FormatShortID(taskID)))
		}
		response.WriteString("\n")
	}

	if len(result_batch.Failed) > 0 {
		response.WriteString("**" + constants.EmojiStatusCancelled + " Başarısız Görevler:**\n")
		for _, failure := range result_batch.Failed {
			response.WriteString(fmt.Sprintf("- %s: %s\n", TaskIDFormat.FormatShortID(failure.TaskID), failure.Error))
		}
		response.WriteString("\n")
	}

	if len(result_batch.Warnings) > 0 {
		response.WriteString("**" + constants.EmojiPriorityAlert + " Uyarılar:**\n")
	}

	return mcp.NewToolResultText(response.String()), nil
}

// GorevBulkTag adds, removes, or replaces tags for multiple tasks
func (h *Handlers) GorevBulkTag(params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Validate task IDs
	taskIDsRaw, ok := params["task_ids"]
	if !ok {
		return mcp.NewToolResultError("task_ids parametresi gerekli"), nil
	}

	taskIDsInterface, ok := taskIDsRaw.([]interface{})
	if !ok {
		return mcp.NewToolResultError("task_ids array formatında olmalı"), nil
	}

	taskIDs := make([]string, len(taskIDsInterface))
	for i, idInterface := range taskIDsInterface {
		if id, ok := idInterface.(string); ok && id != "" {
			taskIDs[i] = id
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("geçersiz task ID index %d", i)), nil
		}
	}

	// Validate tags
	tagsRaw, ok := params["tags"]
	if !ok {
		return mcp.NewToolResultError("tags parametresi gerekli"), nil
	}

	tagsInterface, ok := tagsRaw.([]interface{})
	if !ok {
		return mcp.NewToolResultError("tags array formatında olmalı"), nil
	}

	tags := make([]string, len(tagsInterface))
	for i, tagInterface := range tagsInterface {
		if tag, ok := tagInterface.(string); ok && tag != "" {
			tags[i] = strings.TrimSpace(tag)
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("geçersiz tag index %d", i)), nil
		}
	}

	// Validate operation
	operation, result := h.toolHelpers.Validator.ValidateEnum(params, "operation", []string{"add", "remove", "replace"}, true)
	if result != nil {
		return result, nil
	}

	// Optional parameters
	dryRun := h.toolHelpers.Validator.ValidateBool(params, "dry_run")

	// Create batch processor
	batchProcessor := gorev.NewBatchProcessor(h.isYonetici.VeriYonetici())
	if h.aiContextYonetici != nil {
		batchProcessor.SetAIContextManager(h.aiContextYonetici)
	}

	// Execute bulk tag operation
	request := gorev.BulkTagOperationRequest{
		TaskIDs:   taskIDs,
		Tags:      tags,
		Operation: operation,
		DryRun:    dryRun,
	}

	result_batch, err := batchProcessor.BulkTagOperation(request)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("toplu etiket işlemi başarısız", err), nil
	}

	// Format response
	var response strings.Builder

	if dryRun {
		response.WriteString("🔍 **Kuru Çalıştırma Sonucu**\n\n")
	} else {
		response.WriteString(constants.EmojiStatusCompleted + " **Toplu Etiket İşlemi Tamamlandı**\n\n")
	}

	response.WriteString(fmt.Sprintf("**İşlem:** %s\n", operation))
	response.WriteString(i18n.TMarkdownLabel("etiketler", strings.Join(tags, ", ")) + "\n")
	response.WriteString(fmt.Sprintf("**İşlenen Görev:** %d\n", result_batch.TotalProcessed))
	response.WriteString(fmt.Sprintf("**Başarılı:** %d\n", len(result_batch.Successful)))
	response.WriteString(fmt.Sprintf("**Başarısız:** %d\n", len(result_batch.Failed)))
	response.WriteString(fmt.Sprintf("**Uyarı:** %d\n", len(result_batch.Warnings)))
	response.WriteString(fmt.Sprintf("**Süre:** %v\n\n", result_batch.ExecutionTime))

	if len(result_batch.Successful) > 0 {
		response.WriteString("**" + constants.EmojiStatusCompleted + " Başarılı Görevler:**\n")
		for _, taskID := range result_batch.Successful {
			response.WriteString(fmt.Sprintf("- %s\n", TaskIDFormat.FormatShortID(taskID)))
		}
		response.WriteString("\n")
	}

	if len(result_batch.Failed) > 0 {
		response.WriteString("**" + constants.EmojiStatusCancelled + " Başarısız Görevler:**\n")
		for _, failure := range result_batch.Failed {
			response.WriteString(fmt.Sprintf("- %s: %s\n", TaskIDFormat.FormatShortID(failure.TaskID), failure.Error))
		}
		response.WriteString("\n")
	}

	if len(result_batch.Warnings) > 0 {
		response.WriteString("**" + constants.EmojiPriorityAlert + " Uyarılar:**\n")
		for _, warning := range result_batch.Warnings {
			response.WriteString(fmt.Sprintf("- %s: %s\n", TaskIDFormat.FormatShortID(warning.TaskID), warning.Message))
		}
	}

	return mcp.NewToolResultText(response.String()), nil
}

// GorevSuggestions provides intelligent suggestions for task management
func (h *Handlers) GorevSuggestions(params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Get session ID from AI context if available
	sessionID := ""
	activeTaskID := ""

	if h.aiContextYonetici != nil {
		if activeTask, err := h.aiContextYonetici.GetActiveTask(); err == nil && activeTask != nil {
			activeTaskID = activeTask.ID
		}
	}

	// Get optional parameters
	limit := h.toolHelpers.Validator.ValidateNumber(params, "limit", constants.DefaultSuggestionLimit)

	// Get suggestion types filter
	var types []string
	if typesRaw, ok := params["types"].([]interface{}); ok {
		for _, typeInterface := range typesRaw {
			if typeStr, ok := typeInterface.(string); ok {
				types = append(types, typeStr)
			}
		}
	}

	// Create suggestion engine
	suggestionEngine := gorev.NewSuggestionEngine(h.isYonetici.VeriYonetici())
	if h.aiContextYonetici != nil {
		suggestionEngine.SetAIContextManager(h.aiContextYonetici)
	}

	// Generate suggestions
	request := gorev.SuggestionRequest{
		SessionID:    sessionID,
		ActiveTaskID: activeTaskID,
		Limit:        limit,
		Types:        types,
	}

	response, err := suggestionEngine.GetSuggestions(request)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("öneri oluşturma başarısız", err), nil
	}

	// Format response
	var output strings.Builder

	output.WriteString("🎯 **Akıllı Öneriler**\n\n")
	output.WriteString(fmt.Sprintf("**Toplam:** %d öneri\n", response.TotalCount))
	output.WriteString(fmt.Sprintf("**Süre:** %v\n\n", response.ExecutionTime))

	if len(response.Suggestions) == 0 {
		output.WriteString("ℹ️ Şu anda öneri yok.\n")
		return mcp.NewToolResultText(output.String()), nil
	}

	// Group suggestions by type
	suggestionGroups := make(map[string][]gorev.Suggestion)
	for _, suggestion := range response.Suggestions {
		suggestionGroups[suggestion.Type] = append(suggestionGroups[suggestion.Type], suggestion)
	}

	// Display suggestions by type
	typeNames := map[string]string{
		"next_action":   "🚀 Sonraki Aksiyonlar",
		"similar_task":  "🔍 Benzer Görevler",
		"template":      "📋 Template Önerileri",
		"deadline_risk": "⚠️ Son Tarih Uyarıları",
	}

	typeOrder := []string{"deadline_risk", "next_action", "similar_task", "template"}

	for _, suggestionType := range typeOrder {
		suggestions, exists := suggestionGroups[suggestionType]
		if !exists || len(suggestions) == 0 {
			continue
		}

		output.WriteString(fmt.Sprintf("## %s\n\n", typeNames[suggestionType]))

		for i, suggestion := range suggestions {
			// Priority emoji using constant helper
			priorityEmoji := constants.GetSuggestionPriorityEmoji(suggestion.Priority)

			output.WriteString(fmt.Sprintf("### %d. %s %s\n", i+1, priorityEmoji, suggestion.Title))
			output.WriteString(fmt.Sprintf("**Açıklama:** %s\n", suggestion.Description))
			output.WriteString(fmt.Sprintf("**Önerilen Aksiyon:** `%s`\n", suggestion.Action))
			output.WriteString(fmt.Sprintf("**Güven Skoru:** %.1f%%\n", suggestion.Confidence*100))

			if suggestion.TaskID != "" {
				output.WriteString(fmt.Sprintf("**İlgili Görev:** %s\n", TaskIDFormat.FormatShortID(suggestion.TaskID)))
			}

			output.WriteString("\n")
		}
	}

	output.WriteString("---\n")
	output.WriteString("💡 **İpucu:** Önerilen aksiyonları doğrudan kopyalayıp kullanabilirsiniz.\n")

	return mcp.NewToolResultText(output.String()), nil
}

// GorevIntelligentCreate creates a task with AI-enhanced features
func (h *Handlers) GorevIntelligentCreate(params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Validate required parameters
	title, result := h.toolHelpers.Validator.ValidateRequiredString(params, constants.ParamBaslik)
	if result != nil {
		return result, nil
	}

	// Optional parameters
	description := h.toolHelpers.Validator.ValidateOptionalString(params, constants.ParamAciklama)
	autoSplit := h.toolHelpers.Validator.ValidateBool(params, "auto_split")
	estimateTime := h.toolHelpers.Validator.ValidateBool(params, "estimate_time")
	smartPriority := h.toolHelpers.Validator.ValidateBool(params, "smart_priority")
	suggestTemplate := h.toolHelpers.Validator.ValidateBool(params, "suggest_template")

	// Get project ID if specified
	projeID := h.toolHelpers.Validator.ValidateOptionalString(params, "proje_id")

	// Use active project if no project specified
	if projeID == "" {
		if aktifProje, err := h.isYonetici.AktifProjeGetir(); err == nil && aktifProje != nil {
			projeID = aktifProje.ID
		}
	}

	// Create intelligent task creator
	creator := gorev.NewIntelligentTaskCreator(h.isYonetici.VeriYonetici())

	// Prepare request
	request := gorev.TaskCreationRequest{
		Title:           title,
		Description:     description,
		AutoSplit:       autoSplit,
		EstimateTime:    estimateTime,
		SmartPriority:   smartPriority,
		SuggestTemplate: suggestTemplate,
	}

	// Create task with AI features
	response, err := creator.CreateIntelligentTask(request)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("akıllı görev oluşturma başarısız", err), nil
	}

	// Set project if specified
	if projeID != "" && response.MainTask != nil {
		response.MainTask.ProjeID = projeID
		updateParams := map[string]interface{}{
			"proje_id": projeID,
		}
		if err := h.isYonetici.VeriYonetici().GorevGuncelle(response.MainTask.ID, updateParams); err != nil {
			// Log but don't fail
			slog.Warn("Failed to set project for intelligent task", "error", err)
		}
	}

	// Record interaction with AI context
	if h.aiContextYonetici != nil && response.MainTask != nil {
		h.aiContextYonetici.RecordInteraction(response.MainTask.ID, "intelligent_create", map[string]interface{}{
			"auto_split":       autoSplit,
			"estimate_time":    estimateTime,
			"smart_priority":   smartPriority,
			"suggest_template": suggestTemplate,
			"subtasks_created": len(response.Subtasks),
		})
	}

	// Format response
	var output strings.Builder

	output.WriteString("🧠 **Akıllı Görev Oluşturuldu**\n\n")

	// Main task info
	output.WriteString(fmt.Sprintf("### 📋 Ana Görev\n"))
	output.WriteString(fmt.Sprintf("**Başlık:** %s\n", response.MainTask.Baslik))
	output.WriteString(fmt.Sprintf("**ID:** %s\n", response.MainTask.ID))

	if response.SuggestedPriority != "" {
		priorityEmoji := h.toolHelpers.Formatter.GetPriorityEmoji(response.SuggestedPriority)
		output.WriteString(fmt.Sprintf("**Akıllı Öncelik:** %s %s\n", priorityEmoji, response.SuggestedPriority))
	}

	if response.EstimatedHours > 0 {
		output.WriteString(fmt.Sprintf("**Tahmini Süre:** %.1f saat\n", response.EstimatedHours))
	}

	if projeID != "" {
		if proje, err := h.isYonetici.ProjeGetir(projeID); err == nil {
			output.WriteString(fmt.Sprintf("**Proje:** %s\n", proje.Isim))
		}
	}

	output.WriteString("\n")

	// Subtasks
	if len(response.Subtasks) > 0 {
		output.WriteString(fmt.Sprintf("### 🌳 Otomatik Alt Görevler (%d)\n", len(response.Subtasks)))
		for i, subtask := range response.Subtasks {
			output.WriteString(fmt.Sprintf("%d. %s (`%s`)\n", i+1, subtask.Baslik, TaskIDFormat.FormatShortID(subtask.ID)))
		}
		output.WriteString("\n")
	}

	// Template recommendation
	if response.RecommendedTemplate != "" {
		output.WriteString(fmt.Sprintf("### 📋 Önerilen Template\n"))
		output.WriteString(fmt.Sprintf("**Template:** %s (güven: %.1f%%)\n",
			response.RecommendedTemplate, response.Confidence.TemplateConfidence*100))
		output.WriteString(fmt.Sprintf("**Kullanım:** `template_listele` ile detayları görün\n\n"))
	}

	// Similar tasks
	if len(response.SimilarTasks) > 0 {
		output.WriteString(fmt.Sprintf("### 🔍 Benzer Görevler (%d)\n", len(response.SimilarTasks)))
		for i, similar := range response.SimilarTasks {
			if i >= constants.MaxSuggestionsToShow { // Show top suggestions
				break
			}
			output.WriteString(fmt.Sprintf("%d. %s (%.1f%% benzer - %s)\n",
				i+1, similar.Task.Baslik, similar.SimilarityScore*100, similar.Reason))
		}
		output.WriteString("\n")
	}

	// AI Insights
	if len(response.Insights) > 0 {
		output.WriteString("### 🎯 AI Analiz Sonuçları\n")
		for _, insight := range response.Insights {
			output.WriteString(fmt.Sprintf("- %s\n", insight))
		}
		output.WriteString("\n")
	}

	// Performance info
	output.WriteString("### 📊 Performans\n")
	output.WriteString(fmt.Sprintf("**İşlem Süresi:** %v\n", response.ExecutionTime))
	output.WriteString(fmt.Sprintf("**Güven Skorları:**\n"))
	if response.SuggestedPriority != "" {
		output.WriteString(fmt.Sprintf("  - Öncelik: %.1f%%\n", response.Confidence.PriorityConfidence*100))
	}
	if response.EstimatedHours > 0 {
		output.WriteString(fmt.Sprintf("  - Süre tahmini: %.1f%%\n", response.Confidence.TimeConfidence*100))
	}
	if len(response.Subtasks) > 0 {
		output.WriteString(fmt.Sprintf("  - Alt görev analizi: %.1f%%\n", response.Confidence.SubtaskConfidence*100))
	}

	output.WriteString("\n---\n")
	output.WriteString("💡 **İpucu:** `gorev_detay id='" + response.MainTask.ID + "'` ile detayları görün\n")

	return mcp.NewToolResultText(output.String()), nil
}

// ProjeListele tüm projeleri listeler
func (h *Handlers) ProjeListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	projeler, err := h.isYonetici.ProjeListele()
	if err != nil {
		return mcp.NewToolResultError(i18n.TListFailed("project", err)), nil
	}

	if len(projeler) == 0 {
		return mcp.NewToolResultText(i18n.T("messages.noProjects")), nil
	}

	metin := i18n.T("headers.projectList") + "\n\n"
	for _, proje := range projeler {
		metin += fmt.Sprintf("### %s\n", proje.Isim)
		metin += fmt.Sprintf("- **ID:** %s\n", proje.ID)
		if proje.Tanim != "" {
			metin += fmt.Sprintf("- **Tanım:** %s\n", proje.Tanim)
		}
		metin += fmt.Sprintf("- **Oluşturma:** %s\n", proje.OlusturmaTarih.Format(constants.DateFormatDisplay))

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
	limit := constants.DefaultTaskLimit // Varsayılan limit
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
		case constants.TaskStatusPending:
			beklemede = append(beklemede, g)
		case constants.TaskStatusInProgress:
			devamEdiyor = append(devamEdiyor, g)
		case constants.TaskStatusCompleted:
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
	maxResponseSize := constants.MaxResponseSize
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
		metin += "\n" + constants.EmojiStatusCompleted + " Tamamlandı\n"
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
		return mcp.NewToolResultError(i18n.TListFailed("template", err)), nil
	}

	if len(templates) == 0 {
		return mcp.NewToolResultText(i18n.T("messages.noTemplates")), nil
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
	templateID, ok := params[constants.ParamTemplateID].(string)
	if !ok || templateID == "" {
		return mcp.NewToolResultError("template_id parametresi gerekli"), nil
	}

	degerlerRaw, ok := params[constants.ParamDegerler].(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("degerler parametresi gerekli ve obje tipinde olmalı"), nil
	}

	// Önce template'i ID veya alias ile kontrol et
	template, err := h.isYonetici.VeriYonetici().TemplateIDVeyaAliasIleGetir(templateID)
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

	baslik, ok := params[constants.ParamBaslik].(string)
	if !ok || baslik == "" {
		return mcp.NewToolResultError("başlık parametresi gerekli"), nil
	}

	aciklama, _ := params[constants.ParamAciklama].(string)
	oncelik, _ := params[constants.ParamOncelik].(string)
	if oncelik == "" {
		oncelik = constants.PriorityMedium
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
	sb.WriteString(fmt.Sprintf("- **Beklemede:** %d %s\n", hiyerarsi.BeklemedeAlt, constants.EmojiStatusPending))
	sb.WriteString(fmt.Sprintf("- **İlerleme:** %.1f%%\n\n", hiyerarsi.IlerlemeYuzdesi))

	// Doğrudan alt görevler
	altGorevler, err := h.isYonetici.AltGorevleriGetir(gorevID)
	if err == nil && len(altGorevler) > 0 {
		sb.WriteString("## 🌳 Doğrudan Alt Görevler\n")
		for _, alt := range altGorevler {
			durum := h.toolHelpers.Formatter.GetStatusEmoji(alt.Durum)
			sb.WriteString(fmt.Sprintf("- %s %s (ID: %s)\n", durum, alt.Baslik, alt.ID))
		}
	}

	return mcp.NewToolResultText(sb.String()), nil
}

// CallTool çağrı yapmak için yardımcı metod
func (h *Handlers) CallTool(toolName string, params map[string]interface{}) (*mcp.CallToolResult, error) {
	switch toolName {
	// gorev_olustur was removed in v0.11.1, use templateden_gorev_olustur instead
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
	case "gorev_bulk_transition":
		return h.GorevBulkTransition(params)
	case "gorev_bulk_tag":
		return h.GorevBulkTag(params)
	case "gorev_suggestions":
		return h.GorevSuggestions(params)
	case "gorev_nlp_query":
		return h.GorevNLPQuery(params)
	case "gorev_file_watch_add":
		return h.GorevFileWatchAdd(params)
	case "gorev_file_watch_remove":
		return h.GorevFileWatchRemove(params)
	case "gorev_file_watch_list":
		return h.GorevFileWatchList(params)
	case "gorev_file_watch_stats":
		return h.GorevFileWatchStats(params)
	case "gorev_export":
		return h.GorevExport(params)
	case "gorev_import":
		return h.GorevImport(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("bilinmeyen araç: %s", toolName)), nil
	}
}

func (h *Handlers) RegisterTools(s *server.MCPServer) {
	// Delegate tool registration to the dedicated ToolRegistry
	registry := NewToolRegistry(h)
	registry.RegisterAllTools(s)
}

// AI Context Management Handlers

// GorevSetActive sets the active task for the AI session
func (h *Handlers) GorevSetActive(params map[string]interface{}) (*mcp.CallToolResult, error) {
	taskID, result := h.toolHelpers.Validator.ValidateTaskIDField(params, "task_id")
	if result != nil {
		return result, nil
	}

	err := h.aiContextYonetici.SetActiveTask(taskID)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("Aktif görev ayarlama hatası", err), nil
	}

	// Also record task view for auto-state management
	if err := h.aiContextYonetici.RecordTaskView(taskID); err != nil {
		// Log but don't fail
		// fmt.Printf("Görev görüntüleme kaydı hatası: %v\n", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf(constants.EmojiStatusCompleted+" Görev %s başarıyla aktif görev olarak ayarlandı.", taskID)), nil
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
	limit := constants.DefaultRecentTaskLimit
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
		response.WriteString("### " + constants.EmojiStatusCompleted + " Başarılı Güncellemeler\n")
		for _, id := range result.Successful {
			response.WriteString(fmt.Sprintf("- %s\n", id))
		}
		response.WriteString("\n")
	}

	if len(result.Failed) > 0 {
		response.WriteString("### ❌ Başarısız Güncellemeler\n")
		for _, fail := range result.Failed {
			response.WriteString(fmt.Sprintf("- %s: %s\n", fail.TaskID, fail.Error))
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
		statusEmoji := h.toolHelpers.Formatter.GetStatusEmoji(task.Durum)
		priorityEmoji := h.toolHelpers.Formatter.GetPriorityEmoji(task.Oncelik)

		result.WriteString(fmt.Sprintf("%s %s **%s** (ID: %s)\n", statusEmoji, priorityEmoji, task.Baslik, TaskIDFormat.FormatShortID(task.ID)))

		if task.Aciklama != "" {
			desc := task.Aciklama
			if len(desc) > constants.MaxDescriptionDisplayLength {
				desc = desc[:constants.MaxDescriptionDisplayLength] + "..."
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
			details = append(details, i18n.TMarkdownLabel("etiketler", strings.Join(tagNames, ", ")))
		}
		if task.SonTarih != nil {
			details = append(details, i18n.TMarkdownLabel("son_tarih", task.SonTarih.Format(constants.DateFormatISO)))
		}

		if len(details) > 0 {
			result.WriteString(fmt.Sprintf("   %s\n", strings.Join(details, " | ")))
		}

		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

// File Watcher Handlers

// GorevFileWatchAdd adds a file path to be watched for a specific task
func (h *Handlers) GorevFileWatchAdd(params map[string]interface{}) (*mcp.CallToolResult, error) {
	if h.fileWatcher == nil {
		return mcp.NewToolResultError("Dosya izleme sistemi başlatılamadı"), nil
	}

	// Validate task_id
	taskID, exists := params["task_id"].(string)
	if !exists || taskID == "" {
		return mcp.NewToolResultError("task_id gerekli"), nil
	}

	// Validate file_path
	filePath, exists := params["file_path"].(string)
	if !exists || filePath == "" {
		return mcp.NewToolResultError("file_path gerekli"), nil
	}

	// Verify task exists
	task, err := h.isYonetici.VeriYonetici().GorevGetir(taskID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Görev bulunamadı: %v", err)), nil
	}

	// Add file path to watcher
	if err := h.fileWatcher.AddTaskPath(taskID, filePath); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Dosya izleme eklenemedi: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(constants.EmojiStatusCompleted+" Görev '%s' için '%s' dosya yolu izlemeye eklendi.\n\nDosya değişiklikleri otomatik olarak takip edilecek ve görev durumu gerektiğinde güncellenecek.", task.Baslik, filePath)), nil
}

// GorevFileWatchRemove removes a file path from being watched for a specific task
func (h *Handlers) GorevFileWatchRemove(params map[string]interface{}) (*mcp.CallToolResult, error) {
	if h.fileWatcher == nil {
		return mcp.NewToolResultError("Dosya izleme sistemi başlatılamadı"), nil
	}

	// Validate task_id
	taskID, exists := params["task_id"].(string)
	if !exists || taskID == "" {
		return mcp.NewToolResultError("task_id gerekli"), nil
	}

	// Validate file_path
	filePath, exists := params["file_path"].(string)
	if !exists || filePath == "" {
		return mcp.NewToolResultError("file_path gerekli"), nil
	}

	// Remove file path from watcher
	if err := h.fileWatcher.RemoveTaskPath(taskID, filePath); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Dosya izleme kaldırılamadı: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(constants.EmojiStatusCompleted+" Görev ID %s için '%s' dosya yolu izlemeden kaldırıldı.", taskID, filePath)), nil
}

// GorevFileWatchList lists all watched file paths and their associated tasks
func (h *Handlers) GorevFileWatchList(params map[string]interface{}) (*mcp.CallToolResult, error) {
	if h.fileWatcher == nil {
		return mcp.NewToolResultError("Dosya izleme sistemi başlatılamadı"), nil
	}

	var result strings.Builder
	result.WriteString("## 📁 Dosya İzleme Durumu\n\n")

	// Check if filtering by specific task
	if taskID, exists := params["task_id"].(string); exists && taskID != "" {
		paths := h.fileWatcher.GetTaskPaths(taskID)
		if len(paths) == 0 {
			result.WriteString(fmt.Sprintf("ℹ️ Görev ID %s için izlenen dosya yolu bulunamadı.\n", taskID))
		} else {
			// Get task info for display
			if task, err := h.isYonetici.VeriYonetici().GorevGetir(taskID); err == nil {
				result.WriteString(fmt.Sprintf("### Görev: %s (ID: %s)\n\n", task.Baslik, taskID))
			} else {
				result.WriteString(fmt.Sprintf("### Görev ID: %s\n\n", taskID))
			}

			result.WriteString("İzlenen dosya yolları:\n")
			for _, path := range paths {
				result.WriteString(fmt.Sprintf("- `%s`\n", path))
			}
		}
	} else {
		// List all watched paths
		watchedPaths := h.fileWatcher.GetWatchedPaths()
		if len(watchedPaths) == 0 {
			result.WriteString("ℹ️ Şu anda hiçbir dosya yolu izlenmiyor.\n")
		} else {
			result.WriteString(fmt.Sprintf("**%d dosya yolu izleniyor:**\n\n", len(watchedPaths)))

			for path, taskIDs := range watchedPaths {
				result.WriteString(fmt.Sprintf("📁 `%s`\n", path))
				result.WriteString("   İlişkili görevler: ")

				taskNames := make([]string, 0, len(taskIDs))
				for _, taskID := range taskIDs {
					if task, err := h.isYonetici.VeriYonetici().GorevGetir(taskID); err == nil {
						taskNames = append(taskNames, fmt.Sprintf("%s (ID: %s)", task.Baslik, taskID))
					} else {
						taskNames = append(taskNames, fmt.Sprintf("ID: %s", taskID))
					}
				}
				result.WriteString(strings.Join(taskNames, ", "))
				result.WriteString("\n\n")
			}
		}
	}

	return mcp.NewToolResultText(result.String()), nil
}

// GorevFileWatchStats shows file watcher system statistics
func (h *Handlers) GorevFileWatchStats(params map[string]interface{}) (*mcp.CallToolResult, error) {
	if h.fileWatcher == nil {
		return mcp.NewToolResultError("Dosya izleme sistemi başlatılamadı"), nil
	}

	stats := h.fileWatcher.GetStats()

	var result strings.Builder
	result.WriteString("## 📊 Dosya İzleme Sistemi İstatistikleri\n\n")

	result.WriteString(fmt.Sprintf("**İzlenen dosya yolu sayısı:** %v\n", stats["watched_paths_count"]))
	result.WriteString(fmt.Sprintf("**İzlenen görev sayısı:** %v\n", stats["watched_tasks_count"]))

	if config, exists := stats["config"].(gorev.FileWatcherConfig); exists {
		result.WriteString("\n### Sistem Konfigürasyonu\n\n")
		result.WriteString(fmt.Sprintf("**Otomatik durum güncellemesi:** %v\n", config.AutoUpdateStatus))
		result.WriteString(fmt.Sprintf("**Debounce süresi:** %v\n", config.DebounceDuration))
		result.WriteString(fmt.Sprintf("**Max dosya boyutu:** %d bytes\n", config.MaxFileSize))

		result.WriteString("\n**İzlenen dosya uzantıları:**\n")
		for _, ext := range config.WatchedExtensions {
			result.WriteString(fmt.Sprintf("- %s\n", ext))
		}

		result.WriteString("\n**Yoksayılan desenler:**\n")
		for _, pattern := range config.IgnorePatterns {
			result.WriteString(fmt.Sprintf("- %s\n", pattern))
		}
	}

	result.WriteString("\n💡 **İpucu:** Dosya izleme sistemi, ilişkili dosyalarda değişiklik olduğunda görev durumunu otomatik olarak 'devam_ediyor' durumuna geçirir.")

	return mcp.NewToolResultText(result.String()), nil
}

// GorevExport exports task data to a file
func (h *Handlers) GorevExport(params map[string]interface{}) (*mcp.CallToolResult, error) {
	if h.debug {
		slog.Debug("GorevExport called", "params", params)
	}

	// Parse parameters
	if params == nil {
		return mcp.NewToolResultError(i18n.T("error.noArguments", nil)), nil
	}

	// Extract parameters
	outputPath, _ := params["output_path"].(string)
	if outputPath == "" {
		return mcp.NewToolResultError(i18n.T("error.outputPathRequired", nil)), nil
	}

	format, _ := params["format"].(string)
	if format == "" {
		format = "json"
	}

	includeCompleted := true
	if val, ok := params["include_completed"].(bool); ok {
		includeCompleted = val
	}

	includeDependencies := true
	if val, ok := params["include_dependencies"].(bool); ok {
		includeDependencies = val
	}

	includeTemplates := false
	if val, ok := params["include_templates"].(bool); ok {
		includeTemplates = val
	}

	includeAIContext := false
	if val, ok := params["include_ai_context"].(bool); ok {
		includeAIContext = val
	}

	// Parse project filter
	var projectFilter []string
	if val, ok := params["project_filter"]; ok {
		if filterInterface, ok := val.([]interface{}); ok {
			for _, pid := range filterInterface {
				if pidStr, ok := pid.(string); ok {
					projectFilter = append(projectFilter, pidStr)
				}
			}
		}
	}

	// Parse date range
	var dateRange *gorev.DateRange
	if val, ok := params["date_range"]; ok {
		if rangeMap, ok := val.(map[string]interface{}); ok {
			dateRange = &gorev.DateRange{}
			if from, ok := rangeMap["from"].(string); ok && from != "" {
				if parsedFrom, err := time.Parse(time.RFC3339, from); err == nil {
					dateRange.From = &parsedFrom
				}
			}
			if to, ok := rangeMap["to"].(string); ok && to != "" {
				if parsedTo, err := time.Parse(time.RFC3339, to); err == nil {
					dateRange.To = &parsedTo
				}
			}
		}
	}

	// Create export options
	options := gorev.ExportOptions{
		Format:              format,
		OutputPath:          outputPath,
		DateRange:           dateRange,
		ProjectFilter:       projectFilter,
		IncludeCompleted:    includeCompleted,
		IncludeDependencies: includeDependencies,
		IncludeMetadata:     true,
		IncludeAIContext:    includeAIContext,
		IncludeTemplates:    includeTemplates,
	}

	// Export data
	exportData, err := h.isYonetici.ExportData(options)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf(i18n.T("error.exportFailed", map[string]interface{}{"Error": err}))), nil
	}

	// Save to file
	if err := h.isYonetici.SaveExportToFile(exportData, options); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf(i18n.T("error.exportSaveFailed", map[string]interface{}{"Error": err}))), nil
	}

	// Create summary
	summary := fmt.Sprintf(i18n.T("export.success", map[string]interface{}{
		"Format":    format,
		"Path":      outputPath,
		"Tasks":     exportData.Metadata.TotalTasks,
		"Projects":  exportData.Metadata.TotalProjects,
		"Tags":      len(exportData.Tags),
		"Templates": len(exportData.Templates),
	}))

	return mcp.NewToolResultText(summary), nil
}

// GorevImport imports task data from a file
func (h *Handlers) GorevImport(params map[string]interface{}) (*mcp.CallToolResult, error) {
	if h.debug {
		slog.Debug("GorevImport called", "params", params)
	}

	// Parse parameters
	if params == nil {
		return mcp.NewToolResultError(i18n.T("error.noArguments", nil)), nil
	}

	// Extract parameters
	filePath, _ := params["file_path"].(string)
	if filePath == "" {
		return mcp.NewToolResultError(i18n.T("error.filePathRequired", nil)), nil
	}

	importMode, _ := params["import_mode"].(string)
	if importMode == "" {
		importMode = "merge"
	}

	conflictResolution, _ := params["conflict_resolution"].(string)
	if conflictResolution == "" {
		conflictResolution = "skip"
	}

	preserveIDs := false
	if val, ok := params["preserve_ids"].(bool); ok {
		preserveIDs = val
	}

	dryRun := false
	if val, ok := params["dry_run"].(bool); ok {
		dryRun = val
	}

	// Parse project mapping
	var projectMapping map[string]string
	if val, ok := params["project_mapping"]; ok {
		if mappingInterface, ok := val.(map[string]interface{}); ok {
			projectMapping = make(map[string]string)
			for k, v := range mappingInterface {
				if vStr, ok := v.(string); ok {
					projectMapping[k] = vStr
				}
			}
		}
	}

	// Create import options
	options := gorev.ImportOptions{
		FilePath:           filePath,
		ImportMode:         importMode,
		ConflictResolution: conflictResolution,
		PreserveIDs:        preserveIDs,
		ProjectMapping:     projectMapping,
		DryRun:             dryRun,
	}

	// Import data
	result, err := h.isYonetici.ImportData(options)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf(i18n.T("error.importFailed", map[string]interface{}{"Error": err}))), nil
	}

	// Create summary
	var summary strings.Builder
	if dryRun {
		summary.WriteString(fmt.Sprintf("🔍 **%s**\n\n", i18n.T("import.dryRunResults", nil)))
	} else {
		summary.WriteString(fmt.Sprintf("✅ **%s**\n\n", i18n.T("import.success", nil)))
	}

	summary.WriteString(fmt.Sprintf("📊 **%s**\n", i18n.T("import.statistics", nil)))
	summary.WriteString(fmt.Sprintf("- %s: %d\n", i18n.T("import.importedTasks", nil), result.ImportedTasks))
	summary.WriteString(fmt.Sprintf("- %s: %d\n", i18n.T("import.importedProjects", nil), result.ImportedProjects))
	summary.WriteString(fmt.Sprintf("- %s: %d\n", i18n.T("import.importedTags", nil), result.ImportedTags))
	summary.WriteString(fmt.Sprintf("- %s: %d\n", i18n.T("import.importedTemplates", nil), result.ImportedTemplates))

	if len(result.Conflicts) > 0 {
		summary.WriteString(fmt.Sprintf("\n⚠️ **%s (%d)**\n", i18n.T("import.conflicts", nil), len(result.Conflicts)))
		for i, conflict := range result.Conflicts {
			if i < 5 { // Show max 5 conflicts
				summary.WriteString(fmt.Sprintf("- %s %s: %s\n",
					strings.Title(conflict.Type),
					i18n.T("import.conflictResolution", nil),
					conflict.Resolution))
			}
		}
		if len(result.Conflicts) > 5 {
			summary.WriteString(fmt.Sprintf("... %s %d %s\n",
				i18n.T("import.and", nil),
				len(result.Conflicts)-5,
				i18n.T("import.moreConflicts", nil)))
		}
	}

	if len(result.Errors) > 0 {
		summary.WriteString(fmt.Sprintf("\n❌ **%s**\n", i18n.T("import.errors", nil)))
		for _, errMsg := range result.Errors {
			summary.WriteString(fmt.Sprintf("- %s\n", errMsg))
		}
	}

	if len(result.Warnings) > 0 {
		summary.WriteString(fmt.Sprintf("\n⚠️ **%s**\n", i18n.T("import.warnings", nil)))
		for _, warning := range result.Warnings {
			summary.WriteString(fmt.Sprintf("- %s\n", warning))
		}
	}

	return mcp.NewToolResultText(summary.String()), nil
}

// IDEDetect detects all installed IDEs on the system
func (h *Handlers) IDEDetect(params map[string]interface{}) (*mcp.CallToolResult, error) {
	detector := gorev.NewIDEDetector()

	detectedIDEs, err := detector.DetectAllIDEs()
	if err != nil {
		return nil, fmt.Errorf("IDE detection failed: %w", err)
	}

	if len(detectedIDEs) == 0 {
		return mcp.NewToolResultText(i18n.T("display.noIDEsDetected")), nil
	}

	var result strings.Builder
	result.WriteString("# 🔍 Detected IDEs\n\n")

	for ideType, ide := range detectedIDEs {
		result.WriteString(fmt.Sprintf("## %s %s\n", getIDEIcon(ideType), ide.Name))
		result.WriteString(fmt.Sprintf("- **Type**: %s\n", ideType))
		result.WriteString(fmt.Sprintf("- **Executable**: %s\n", ide.ExecutablePath))
		result.WriteString(fmt.Sprintf("- **Config Path**: %s\n", ide.ConfigPath))
		result.WriteString(fmt.Sprintf("- **Extensions Path**: %s\n", ide.ExtensionsPath))
		if ide.Version != "unknown" && ide.Version != "" {
			result.WriteString(fmt.Sprintf("- **Version**: %s\n", ide.Version))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

// IDEInstallExtension installs the Gorev extension to specified IDE(s)
func (h *Handlers) IDEInstallExtension(params map[string]interface{}) (*mcp.CallToolResult, error) {
	ctx := context.Background()

	// Parse parameters
	ideType, exists := params["ide_type"].(string)
	if !exists {
		return nil, fmt.Errorf(i18n.T("error.parameterRequired", map[string]interface{}{"Param": "ide_type"}))
	}

	installToAll := false
	if ideType == "all" {
		installToAll = true
	}

	// Create detector and installer
	detector := gorev.NewIDEDetector()
	detector.DetectAllIDEs()

	installer := gorev.NewExtensionInstaller(detector)
	defer installer.Cleanup()

	// Get latest extension info from GitHub
	extensionInfo, err := installer.GetLatestExtensionInfo(ctx, "msenol", "Gorev")
	if err != nil {
		return nil, fmt.Errorf("failed to get extension info: %w", err)
	}

	var results []gorev.InstallResult
	var err2 error

	if installToAll {
		results, err2 = installer.InstallToAllIDEs(ctx, extensionInfo)
	} else {
		result, err := installer.InstallExtension(ctx, gorev.IDEType(ideType), extensionInfo)
		if result != nil {
			results = append(results, *result)
		}
		err2 = err
	}

	// Format results
	var output strings.Builder
	output.WriteString("# 📦 Extension Installation Results\n\n")

	successCount := 0
	for _, result := range results {
		icon := "❌"
		if result.Success {
			icon = "✅"
			successCount++
		}

		output.WriteString(fmt.Sprintf("## %s %s\n", icon, result.IDE))
		output.WriteString(fmt.Sprintf("- **Extension**: %s\n", result.Extension))
		if result.Version != "" {
			output.WriteString(fmt.Sprintf("- **Version**: %s\n", result.Version))
		}
		output.WriteString(fmt.Sprintf("- **Status**: %s\n\n", result.Message))
	}

	if len(results) > 0 {
		output.WriteString(fmt.Sprintf("**Summary**: %d/%d installations successful\n", successCount, len(results)))
	}

	if err2 != nil && len(results) == 0 {
		return nil, err2
	}

	return mcp.NewToolResultText(output.String()), nil
}

// IDEUninstallExtension removes the Gorev extension from specified IDE
func (h *Handlers) IDEUninstallExtension(params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Parse parameters
	ideType, exists := params["ide_type"].(string)
	if !exists {
		return nil, fmt.Errorf(i18n.T("error.parameterRequired", map[string]interface{}{"Param": "ide_type"}))
	}

	extensionID, exists := params["extension_id"].(string)
	if !exists {
		extensionID = "mehmetsenol.gorev-vscode" // default
	}

	// Create detector and installer
	detector := gorev.NewIDEDetector()
	detector.DetectAllIDEs()

	installer := gorev.NewExtensionInstaller(detector)

	result, err := installer.UninstallExtension(gorev.IDEType(ideType), extensionID)
	if err != nil {
		return nil, err
	}

	var output strings.Builder
	output.WriteString("# 🗑️ Extension Uninstallation Result\n\n")

	icon := "❌"
	if result.Success {
		icon = "✅"
	}

	output.WriteString(fmt.Sprintf("## %s %s\n", icon, result.IDE))
	output.WriteString(fmt.Sprintf("- **Extension**: %s\n", result.Extension))
	output.WriteString(fmt.Sprintf("- **Status**: %s\n", result.Message))

	return mcp.NewToolResultText(output.String()), nil
}

// IDEExtensionStatus checks the installation status of Gorev extension in IDEs
func (h *Handlers) IDEExtensionStatus(params map[string]interface{}) (*mcp.CallToolResult, error) {
	ctx := context.Background()

	// Create detector and installer
	detector := gorev.NewIDEDetector()
	detectedIDEs, err := detector.DetectAllIDEs()
	if err != nil {
		return nil, err
	}

	installer := gorev.NewExtensionInstaller(detector)

	// Get latest version from GitHub
	extensionInfo, err := installer.GetLatestExtensionInfo(ctx, "msenol", "Gorev")
	var latestVersion string
	if err == nil {
		latestVersion = extensionInfo.Version
	}

	var output strings.Builder
	output.WriteString("# 📊 Extension Status Report\n\n")

	if latestVersion != "" {
		output.WriteString(fmt.Sprintf("**Latest Available Version**: v%s\n\n", latestVersion))
	}

	for ideType, ide := range detectedIDEs {
		output.WriteString(fmt.Sprintf("## %s %s\n", getIDEIcon(ideType), ide.Name))

		// Check if extension is installed
		isInstalled, err := detector.IsExtensionInstalled(ideType, "mehmetsenol.gorev-vscode")
		if err != nil {
			output.WriteString(fmt.Sprintf("- **Status**: ❌ Error checking installation: %s\n", err))
		} else if !isInstalled {
			output.WriteString("- **Status**: ❌ Not installed\n")
		} else {
			// Get installed version
			installedVersion, err := detector.GetExtensionVersion(ideType, "mehmetsenol.gorev-vscode")
			if err != nil {
				output.WriteString("- **Status**: ✅ Installed (version unknown)\n")
			} else {
				statusIcon := "✅"
				updateStatus := ""

				if latestVersion != "" && installedVersion != latestVersion {
					statusIcon = "⚠️"
					updateStatus = fmt.Sprintf(" (Update available: v%s)", latestVersion)
				}

				output.WriteString(fmt.Sprintf("- **Status**: %s Installed v%s%s\n", statusIcon, installedVersion, updateStatus))
			}
		}
		output.WriteString("\n")
	}

	return mcp.NewToolResultText(output.String()), nil
}

// IDEUpdateExtension updates the Gorev extension to latest version
func (h *Handlers) IDEUpdateExtension(params map[string]interface{}) (*mcp.CallToolResult, error) {
	ctx := context.Background()

	// Parse parameters
	ideType, exists := params["ide_type"].(string)
	if !exists {
		return nil, fmt.Errorf(i18n.T("error.parameterRequired", map[string]interface{}{"Param": "ide_type"}))
	}

	updateAll := false
	if ideType == "all" {
		updateAll = true
	}

	// Create detector and installer
	detector := gorev.NewIDEDetector()
	detector.DetectAllIDEs()

	installer := gorev.NewExtensionInstaller(detector)
	defer installer.Cleanup()

	// Get latest extension info
	extensionInfo, err := installer.GetLatestExtensionInfo(ctx, "msenol", "Gorev")
	if err != nil {
		return nil, fmt.Errorf("failed to get extension info: %w", err)
	}

	var results []gorev.InstallResult
	var err2 error

	if updateAll {
		// Update all detected IDEs
		allIDEs := detector.GetAllDetectedIDEs()
		for ideTypeKey := range allIDEs {
			result, err := installer.InstallExtension(ctx, ideTypeKey, extensionInfo)
			if result != nil {
				results = append(results, *result)
			}
			if err != nil && err2 == nil {
				err2 = err // Keep first error
			}
		}
	} else {
		result, err := installer.InstallExtension(ctx, gorev.IDEType(ideType), extensionInfo)
		if result != nil {
			results = append(results, *result)
		}
		err2 = err
	}

	// Format results
	var output strings.Builder
	output.WriteString(fmt.Sprintf("# 🔄 Extension Update Results (v%s)\n\n", extensionInfo.Version))

	successCount := 0
	for _, result := range results {
		icon := "❌"
		if result.Success {
			icon = "✅"
			successCount++
		}

		output.WriteString(fmt.Sprintf("## %s %s\n", icon, result.IDE))
		output.WriteString(fmt.Sprintf("- **Extension**: %s\n", result.Extension))
		output.WriteString(fmt.Sprintf("- **Version**: v%s\n", result.Version))
		output.WriteString(fmt.Sprintf("- **Status**: %s\n\n", result.Message))
	}

	if len(results) > 0 {
		output.WriteString(fmt.Sprintf("**Summary**: %d/%d updates successful\n", successCount, len(results)))
	}

	if err2 != nil && len(results) == 0 {
		return nil, err2
	}

	return mcp.NewToolResultText(output.String()), nil
}

// getIDEIcon returns an appropriate emoji icon for the IDE type
func getIDEIcon(ideType gorev.IDEType) string {
	switch ideType {
	case gorev.IDETypeVSCode:
		return "🔵"
	case gorev.IDETypeCursor:
		return "🖱️"
	case gorev.IDETypeWindsurf:
		return "🌊"
	default:
		return "💻"
	}
}
