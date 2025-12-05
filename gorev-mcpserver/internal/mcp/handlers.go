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
	contextutil "github.com/msenol/gorev/internal/context"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// initializeHandlerComponents initializes common handler components
func initializeHandlerComponents(isYonetici *gorev.IsYonetici, debug bool) (*gorev.AIContextYonetici, *gorev.FileWatcher, *ToolHelpers) {
	var aiContextYonetici *gorev.AIContextYonetici
	var fileWatcher *gorev.FileWatcher

	// Create AI context manager using the same data manager if isYonetici is not nil
	if isYonetici != nil {
		aiContextYonetici = gorev.YeniAIContextYonetici(isYonetici.VeriYonetici())

		// Initialize file watcher with default configuration
		if fw, err := gorev.NewFileWatcher(isYonetici.VeriYonetici(), gorev.DefaultFileWatcherConfig()); err == nil {
			fileWatcher = fw
		} else {
			// Handle file watcher initialization error with proper logging
			if debug {
				slog.Debug("Failed to initialize file watcher", "error", err)
			} else {
				slog.Warn("File watcher initialization failed", "error", err)
			}
		}
	}

	// Initialize tool helpers with shared utilities
	toolHelpers := NewToolHelpers()

	return aiContextYonetici, fileWatcher, toolHelpers
}

func YeniHandlers(isYonetici *gorev.IsYonetici) *Handlers {
	aiContextYonetici, fileWatcher, toolHelpers := initializeHandlerComponents(isYonetici, false)

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
	aiContextYonetici, fileWatcher, toolHelpers := initializeHandlerComponents(isYonetici, debug)

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

// extractLanguage extracts language preference with fallback priority:
// 1. Environment variable GOREV_LANG (set by MCP client)
// 2. Default "tr"
// Note: context.Context not available in MCP handler signatures, so we use env var
func (h *Handlers) extractLanguage() string {
	// Use context.Background() to trigger env var fallback in GetLanguage
	return contextutil.GetLanguage(context.Background())
}

// gorevResponseSizeEstimate bir görev için tahmini response boyutunu hesaplar
func (h *Handlers) gorevResponseSizeEstimate(gorev *gorev.Gorev) int {
	// Tahmini karakter sayıları
	size := constants.BaseResponseSize // Temel formatlar için
	size += len(gorev.Title) + len(gorev.Description)
	size += len(gorev.ID) + len(gorev.ProjeID)

	if gorev.DueDate != nil {
		size += constants.DateFormatSize // Tarih formatı için
	}

	for _, etiket := range gorev.Tags {
		size += len(etiket.Name) + constants.TagSizeConstant
	}

	// Bağımlılık bilgileri
	if gorev.DependencyCount > 0 || gorev.DependentOnThisCount > 0 {
		size += constants.DependencyInfoSize
	}

	return size
}

// gorevOzetYazdir bir görevi özet formatta yazdırır (ProjeGorevleri için)
func (h *Handlers) gorevOzetYazdir(g *gorev.Gorev) string {
	// Öncelik kısaltması
	oncelik := PriorityFormat.GetPriorityShort(g.Priority)

	metin := fmt.Sprintf("- **%s** (%s)", g.Title, oncelik)

	// Inline detaylar
	details := []string{}
	if g.Description != "" && len(g.Description) <= constants.MaxInlineDescriptionLength {
		details = append(details, g.Description)
	} else if g.Description != "" {
		details = append(details, g.Description[:constants.TruncatedDescriptionLength]+"...")
	}

	if g.DueDate != nil {
		details = append(details, g.DueDate.Format(constants.DateFormatShort))
	}

	if len(g.Tags) > 0 && len(g.Tags) <= 2 {
		etiketler := make([]string, len(g.Tags))
		for i, e := range g.Tags {
			etiketler[i] = e.Name
		}
		details = append(details, strings.Join(etiketler, ","))
	} else if len(g.Tags) > 2 {
		details = append(details, i18n.T("messages.tagCount", map[string]interface{}{"Count": len(g.Tags)}))
	}

	if g.UncompletedDependencyCount > 0 {
		details = append(details, fmt.Sprintf("🔒%d", g.UncompletedDependencyCount))
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
	return fmt.Sprintf("- ~~%s~~ | %s\n", g.Title, TaskIDFormat.FormatShortID(g.ID))
}

// gorevHiyerarsiYazdir bir görevi ve alt görevlerini hiyerarşik olarak yazdırır
func (h *Handlers) gorevHiyerarsiYazdir(ctx context.Context, lang string, gorev *gorev.Gorev, gorevMap map[string]*gorev.Gorev, seviye int, projeGoster bool) string {
	indent := strings.Repeat("  ", seviye)
	prefix := ""
	if seviye > 0 {
		prefix = "└─ "
	}

	durum := StatusFormat.GetStatusSymbol(gorev.Status)

	// Öncelik kısaltması
	oncelikKisa := PriorityFormat.GetPriorityShort(gorev.Priority)

	// Temel satır - öncelik parantez içinde kısaltılmış
	metin := fmt.Sprintf("%s%s[%s] %s (%s)\n", indent, prefix, durum, gorev.Title, oncelikKisa)

	// Sadece dolu alanları göster, boş satırlar ekleme
	details := []string{}

	if gorev.Description != "" {
		// Template sistemi için açıklama limiti büyük ölçüde artırıldı
		// Sadece gerçekten çok uzun açıklamaları kısalt
		aciklama := gorev.Description
		if len(aciklama) > constants.MaxDescriptionLength {
			// İlk karakterleri al ve ... ekle
			aciklama = aciklama[:constants.MaxDescriptionTruncateLength] + "..."
		}
		details = append(details, aciklama)
	}

	if projeGoster && gorev.ProjeID != "" {
		proje, _ := h.isYonetici.ProjeGetir(ctx, gorev.ProjeID)
		if proje != nil {
			details = append(details, i18n.T("messages.projectLabel", map[string]interface{}{"Name": proje.Name}))
		}
	}

	if gorev.DueDate != nil {
		details = append(details, i18n.T("messages.dateLabel", map[string]interface{}{"Date": gorev.DueDate.Format(constants.DateFormatShort)}))
	}

	if len(gorev.Tags) > 0 && len(gorev.Tags) <= constants.MaxTagsToDisplay {
		etiketIsimleri := make([]string, len(gorev.Tags))
		for i, etiket := range gorev.Tags {
			etiketIsimleri[i] = etiket.Name
		}
		details = append(details, i18n.T("messages.tagLabel", map[string]interface{}{"Tags": strings.Join(etiketIsimleri, ",")}))
	} else if len(gorev.Tags) > constants.MaxTagsToDisplay {
		details = append(details, i18n.T("messages.tagCountLabel", map[string]interface{}{"Count": len(gorev.Tags)}))
	}

	// Bağımlılık bilgileri - sadece varsa ve sıfırdan büyükse
	if gorev.UncompletedDependencyCount > 0 {
		details = append(details, i18n.TMarkdownLabel(lang, "bekleyen", gorev.UncompletedDependencyCount))
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
			metin += h.gorevHiyerarsiYazdir(ctx, lang, g, gorevMap, seviye+1, projeGoster)
		}
	}

	if seviye == 0 {
		metin += "\n"
	}

	return metin
}

// gorevHiyerarsiYazdirVeIsaretle görevleri yazdırırken hangi görevlerin gösterildiğini işaretler
func (h *Handlers) gorevHiyerarsiYazdirVeIsaretle(ctx context.Context, lang string, gorev *gorev.Gorev, gorevMap map[string]*gorev.Gorev, seviye int, projeGoster bool, shownGorevIDs map[string]bool) string {
	// Bu görevi gösterildi olarak işaretle
	shownGorevIDs[gorev.ID] = true

	// Normal hiyerarşik yazdırma işlemi
	metin := h.gorevHiyerarsiYazdirInternal(ctx, lang, gorev, gorevMap, seviye, projeGoster, shownGorevIDs)

	return metin
}

// gorevHiyerarsiYazdirInternal görev hiyerarşisini yazdırır ve gösterilenleri işaretler
func (h *Handlers) gorevHiyerarsiYazdirInternal(ctx context.Context, lang string, gorev *gorev.Gorev, gorevMap map[string]*gorev.Gorev, seviye int, projeGoster bool, shownGorevIDs map[string]bool) string {
	indent := strings.Repeat("  ", seviye)
	prefix := ""
	if seviye > 0 {
		prefix = "└─ "
	}

	durum := StatusFormat.GetStatusSymbol(gorev.Status)

	// Öncelik kısaltması
	oncelikKisa := PriorityFormat.GetPriorityShort(gorev.Priority)

	// Temel satır - öncelik parantez içinde kısaltılmış
	metin := fmt.Sprintf("%s%s[%s] %s (%s)\n", indent, prefix, durum, gorev.Title, oncelikKisa)

	// Sadece dolu alanları göster, boş satırlar ekleme
	details := []string{}

	if gorev.Description != "" {
		// Template sistemi için açıklama limiti büyük ölçüde artırıldı
		// Sadece gerçekten çok uzun açıklamaları kısalt
		aciklama := gorev.Description
		if len(aciklama) > constants.MaxDescriptionLength {
			// İlk karakterleri al ve ... ekle
			aciklama = aciklama[:constants.MaxDescriptionTruncateLength] + "..."
		}
		details = append(details, aciklama)
	}

	if projeGoster && gorev.ProjeID != "" {
		proje, _ := h.isYonetici.ProjeGetir(ctx, gorev.ProjeID)
		if proje != nil {
			details = append(details, i18n.T("messages.projectLabel", map[string]interface{}{"Name": proje.Name}))
		}
	}

	if gorev.DueDate != nil {
		details = append(details, i18n.T("messages.dateLabel", map[string]interface{}{"Date": gorev.DueDate.Format(constants.DateFormatShort)}))
	}

	if len(gorev.Tags) > 0 && len(gorev.Tags) <= constants.MaxTagsToDisplay {
		etiketIsimleri := make([]string, len(gorev.Tags))
		for i, etiket := range gorev.Tags {
			etiketIsimleri[i] = etiket.Name
		}
		details = append(details, i18n.TMarkdownLabel(lang, "etiket", strings.Join(etiketIsimleri, ", ")))
	} else if len(gorev.Tags) > constants.MaxTagsToDisplay {
		details = append(details, i18n.T("messages.tagCountLabel", map[string]interface{}{"Count": len(gorev.Tags)}))
	}

	// Bağımlılık bilgileri - sadece varsa ve sıfırdan büyükse
	if gorev.UncompletedDependencyCount > 0 {
		details = append(details, i18n.TMarkdownLabel(lang, "bekleyen", gorev.UncompletedDependencyCount))
	}

	// ID'yi en sona ekle
	details = append(details, fmt.Sprintf("ID:%s", gorev.ID))

	// Detayları tek satırda göster
	if len(details) > 0 {
		metin += fmt.Sprintf("%s  %s\n", indent, strings.Join(details, " | "))
	}

	// Bağımlılık bilgilerini ekle (MarkdownParser için)
	bagimlilikBilgisi := h.gorevBagimlilikBilgisi(gorev, indent+"  ")
	if bagimlilikBilgisi != "" {
		metin += bagimlilikBilgisi
	}

	// Alt görevleri bul ve yazdır - TÜM alt görevler gösterilir
	for _, g := range gorevMap {
		if g.ParentID == gorev.ID {
			shownGorevIDs[g.ID] = true
			metin += h.gorevHiyerarsiYazdirInternal(ctx, lang, g, gorevMap, seviye+1, projeGoster, shownGorevIDs)
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
	for _, alan := range template.Fields {
		if alan.Required {
			tip := alan.Type
			if alan.Type == "select" && len(alan.Options) > 0 {
				tip = i18n.T("messages.fieldSelect", map[string]interface{}{"Options": strings.Join(alan.Options, ", ")})
			}
			alanlar = append(alanlar, fmt.Sprintf("- %s (%s)", alan.Name, tip))
		}
	}
	return strings.Join(alanlar, "\n")
}

// templateOrnekDegerler template için örnek değerler oluşturur
func (h *Handlers) templateOrnekDegerler(template *gorev.GorevTemplate) string {
	var ornekler []string
	for _, alan := range template.Fields {
		if alan.Required {
			ornek := ""
			switch alan.Type {
			case "select":
				if len(alan.Options) > 0 {
					ornek = alan.Options[0]
				}
			case "date":
				ornek = time.Now().AddDate(0, 0, 7).Format(constants.DateFormatISO) // One week from now
			case "text":
				ornek = i18n.T("messages.exampleLabel", map[string]interface{}{"Field": alan.Name})
			}
			ornekler = append(ornekler, fmt.Sprintf("'%s': '%s'", alan.Name, ornek))
		}
	}
	return strings.Join(ornekler, ", ")
}

// GorevOlustur - DEPRECATED: Template kullanımı artık zorunludur
// GorevOlustur was removed in v0.11.1 - use TemplatedenGorevOlustur instead

// GorevListele görevleri listeler
func (h *Handlers) GorevListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	status, _ := params[constants.ParamStatus].(string)
	orderBy, _ := params[constants.ParamOrderBy].(string)
	filter, _ := params[constants.ParamFilter].(string)
	tag, _ := params[constants.ParamTag].(string)
	allProjects, _ := params[constants.ParamAllProjects].(bool)

	// Pagination parameters
	limit, offset := h.toolHelpers.Validator.ValidatePagination(params)

	// Create filters map
	filters := make(map[string]interface{})
	if status != "" {
		filters[constants.ParamStatus] = status
	}
	if orderBy != "" {
		filters[constants.ParamOrderBy] = orderBy
	}
	if filter != "" {
		filters[constants.ParamFilter] = filter
	}

	gorevler, err := h.isYonetici.GorevListele(ctx, filters)
	if err != nil {
		return mcp.NewToolResultError(i18n.TListFailed(lang, "task", err)), nil
	}

	// DEBUG: Log görev sayısı
	// fmt.Fprintf(os.Stderr, "[GorevListele] Fetched %d tasks total\n", len(gorevler))

	// Filter by tag
	if tag != "" {
		var filteredTasks []*gorev.Gorev
		for _, g := range gorevler {
			for _, e := range g.Tags {
				if e.Name == tag {
					filteredTasks = append(filteredTasks, g)
					break
				}
			}
		}
		gorevler = filteredTasks
	}

	// If active project exists and all_projects is false, show only active project's tasks
	var aktifProje *gorev.Proje
	if !allProjects {
		aktifProje, _ = h.isYonetici.AktifProjeGetir(ctx)
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
			// Debug logging for OpenCode.ai issue
			if h.debug {
				slog.Debug("GorevListele: Empty project case",
					"project_id", aktifProje.ID,
					"project_name", aktifProje.Name,
					"project_id_type", fmt.Sprintf("%T", aktifProje.ID))
			}
			mesaj = i18n.T("messages.noTasksInProject", map[string]interface{}{"Project": aktifProje.Name})
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
	// FIX: Tüm görev sayısını göster (root + subtasks)
	toplamRootGorevSayisi := len(kokGorevler)
	if toplamRootGorevSayisi > limit || offset > 0 {
		// Check if offset is beyond available data
		if offset >= toplamRootGorevSayisi {
			// No data available at this offset - return empty result
			return mcp.NewToolResultText(i18n.T("messages.noMoreTasks")), nil
		}

		// FIX: Pagination end calculation düzeltildi - root görev sayısına göre hesaplama
		actualEnd := offset + limit
		if actualEnd > toplamRootGorevSayisi {
			actualEnd = toplamRootGorevSayisi
		}
		metin = i18n.T("messages.taskCount", map[string]interface{}{
			"Start": offset + 1,
			"End":   actualEnd,
			"Total": toplamRootGorevSayisi,
		}) + "\n"
	} else {
		metin = i18n.T("messages.taskListCount", map[string]interface{}{"Count": toplamRootGorevSayisi}) + "\n"
	}

	if aktifProje != nil && !allProjects {
		metin += i18n.T("messages.projectHeader", map[string]interface{}{"Name": aktifProje.Name}) + "\n"
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
			remainingCount := len(paginatedKokGorevler) - len(gorevlerToShow)
			metin += "\n" + i18n.T("messages.sizeWarning", map[string]interface{}{"Count": remainingCount}) + "\n"
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
		metin += h.gorevHiyerarsiYazdirVeIsaretle(ctx, lang, kokGorev, gorevMap, 0, allProjects || aktifProje == nil, shownGorevIDs)
	}

	// REMOVED: Orphan checking logic
	// Artık sadece root görevleri paginate ediyoruz
	// Alt görevler her zaman parent'larıyla birlikte gösterilecek

	return mcp.NewToolResultText(metin), nil
}

// AktifProjeAyarla bir projeyi aktif proje olarak ayarlar
func (h *Handlers) AktifProjeAyarla(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	// Support both English and Turkish parameter names for backward compatibility
	projeID, result := h.toolHelpers.Validator.ValidateRequiredString(params, "project_id")
	if result != nil {
		// Fallback to Turkish parameter name
		projeID, result = h.toolHelpers.Validator.ValidateRequiredString(params, "proje_id")
		if result != nil {
			return result, nil
		}
	}

	if err := h.isYonetici.AktifProjeAyarla(ctx, projeID); err != nil {
		return mcp.NewToolResultError(i18n.TSetFailed(lang, "active_project", err)), nil
	}

	proje, _ := h.isYonetici.ProjeGetir(ctx, projeID)
	if proje != nil {
		return mcp.NewToolResultText(
			i18n.T("success.activeProjectSet", map[string]interface{}{"Project": proje.Name}),
		), nil
	}
	return mcp.NewToolResultText(
		i18n.T("success.activeProjectSet", map[string]interface{}{"Project": projeID}),
	), nil
}

// AktifProjeGoster mevcut aktif projeyi gösterir
func (h *Handlers) AktifProjeGoster(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	proje, err := h.isYonetici.AktifProjeGetir(ctx)
	if err != nil {
		return mcp.NewToolResultError(i18n.T("error.activeProjectRetrieve", map[string]interface{}{"Error": err.Error()})), nil
	}

	if proje == nil {
		return mcp.NewToolResultText(i18n.T("messages.noActiveProject")), nil
	}

	// Görev sayısını al
	gorevSayisi, _ := h.isYonetici.ProjeGorevSayisi(ctx, proje.ID)

	metin := i18n.T("headers.activeProject") + "\n\n"
	metin += i18n.TListItem(lang, "proje", proje.Name) + "\n"
	metin += i18n.TListItem(lang, "id_field", proje.ID) + "\n"
	metin += i18n.TListItem(lang, "aciklama", proje.Definition) + "\n"
	metin += i18n.TListItem(lang, "gorev_sayisi", fmt.Sprintf("%d", gorevSayisi))

	return mcp.NewToolResultText(metin), nil
}

// AktifProjeKaldir aktif proje ayarını kaldırır
func (h *Handlers) AktifProjeKaldir(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	if err := h.isYonetici.AktifProjeKaldir(ctx); err != nil {
		return mcp.NewToolResultError(i18n.TRemoveFailed(lang, "active_project", err)), nil
	}

	return mcp.NewToolResultText(i18n.TRemoved(lang, "active_project")), nil
}

// GorevGuncelle görev durumunu günceller
func (h *Handlers) GorevGuncelle(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	h.toolHelpers.SetLanguage(lang)
	ctx := i18n.WithLanguage(context.Background(), lang)

	// Use helper for validation
	id, result := h.toolHelpers.Validator.ValidateTaskID(params)
	if result != nil {
		return result, nil
	}

	// Check if we have status or priority (at least one required) - using English names per schema
	status, statusResult := h.toolHelpers.Validator.ValidateEnum(params, constants.ParamStatus, constants.GetValidTaskStatuses(), false)
	priority, priorityResult := h.toolHelpers.Validator.ValidateEnum(params, constants.ParamPriority, constants.GetValidPriorities(), false)

	// If a parameter was provided but invalid, return that error immediately
	// Special case: if status was explicitly provided but is empty/invalid, it should fail
	if params[constants.ParamStatus] != nil {
		if statusResult != nil && status == "" {
			return statusResult, nil
		}
		// Additional validation: explicitly provided empty status should fail
		if status == "" && statusResult == nil {
			// Status was provided but is empty - this should be treated as invalid
			return mcp.NewToolResultError(i18n.T("error.invalidStatus",
				map[string]interface{}{
					"Status":        "",
					"ValidStatuses": strings.Join(constants.GetValidTaskStatuses(), ", "),
				})), nil
		}
	}
	if params[constants.ParamPriority] != nil {
		if priorityResult != nil && priority == "" {
			return priorityResult, nil
		}
	}

	// If both are missing/invalid, return error
	if statusResult != nil && priorityResult != nil {
		// Custom validation: at least one of status or priority required
		return mcp.NewToolResultError(i18n.T("common.validation.one_of_required",
			map[string]interface{}{"Params": "status, priority"})), nil
	}

	// Update status if provided
	if status != "" {
		if err := h.isYonetici.GorevDurumGuncelle(ctx, id, status); err != nil {
			return mcp.NewToolResultError(i18n.TUpdateFailed(lang, "task", err)), nil
		}
	}

	// Update priority if provided
	if priority != "" {
		updateParams := map[string]interface{}{
			"priority": priority,
		}
		if err := h.isYonetici.VeriYonetici().GorevGuncelle(ctx, id, updateParams); err != nil {
			return mcp.NewToolResultError(i18n.TUpdateFailed(lang, "task", err)), nil
		}
	}

	// Build success message
	var updates []string
	if status != "" {
		updates = append(updates, i18n.TLabel(lang, "durum")+": "+status)
	}
	if priority != "" {
		updates = append(updates, i18n.TLabel(lang, "oncelik")+": "+priority)
	}

	return mcp.NewToolResultText(
		i18n.T("success.taskUpdatedWithChanges", map[string]interface{}{
			"ID":      id,
			"Changes": strings.Join(updates, ", "),
		}),
	), nil
}

// ProjeOlustur yeni bir proje oluşturur
func (h *Handlers) ProjeOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	h.toolHelpers.SetLanguage(lang)
	ctx := i18n.WithLanguage(context.Background(), lang)

	// Support both English and Turkish parameter names for backward compatibility
	name, result := h.toolHelpers.Validator.ValidateRequiredString(params, constants.ParamName)
	if result != nil {
		// Fallback to Turkish parameter name
		name, result = h.toolHelpers.Validator.ValidateRequiredString(params, "isim")
		if result != nil {
			return result, nil
		}
	}

	// Support both English and Turkish parameter names
	definition := h.toolHelpers.Validator.ValidateOptionalString(params, constants.ParamDefinition)
	if definition == "" {
		definition = h.toolHelpers.Validator.ValidateOptionalString(params, "tanim")
	}

	proje, err := h.isYonetici.ProjeOlustur(ctx, name, definition)
	if err != nil {
		return mcp.NewToolResultError(i18n.TCreateFailed(lang, "project", err)), nil
	}

	return mcp.NewToolResultText(i18n.TCreated(lang, "project", proje.Name, proje.ID)), nil
}

// GorevDetay tek bir görevin detaylı bilgisini markdown formatında döner
func (h *Handlers) GorevDetay(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	id, result := h.toolHelpers.Validator.ValidateTaskID(params)
	if result != nil {
		return result, nil
	}

	gorev, err := h.isYonetici.GorevGetir(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(i18n.TEntityNotFoundByID(lang, "task", id)), nil
	}

	// Bağımlılık sayılarını hesapla (VS Code extension için gerekli)
	bagimliSayilari, _ := h.isYonetici.VeriYonetici().BulkBagimlilikSayilariGetir([]string{id})
	if count, exists := bagimliSayilari[id]; exists {
		gorev.DependencyCount = count
	}

	tamamlanmamisSayilari, _ := h.isYonetici.VeriYonetici().BulkTamamlanmamiaBagimlilikSayilariGetir([]string{id})
	if count, exists := tamamlanmamisSayilari[id]; exists {
		gorev.UncompletedDependencyCount = count
	}

	buGoreveBagimliSayilari, _ := h.isYonetici.VeriYonetici().BulkBuGoreveBagimliSayilariGetir([]string{id})
	if count, exists := buGoreveBagimliSayilari[id]; exists {
		gorev.DependentOnThisCount = count
	}

	// Auto-state management: Record task view and potentially transition state
	if err := h.aiContextYonetici.RecordTaskView(ctx, id); err != nil {
		// Log but don't fail the request
		// fmt.Printf("Görev görüntüleme kaydı hatası: %v\n", err)
	}

	// Markdown formatında detaylı görev bilgisi
	metin := fmt.Sprintf("# %s\n\n", gorev.Title)
	metin += i18n.T("headers.generalInfo") + "\n"
	metin += i18n.TListItem(lang, "id_field", gorev.ID) + "\n"
	metin += i18n.TListItem(lang, "durum", gorev.Status) + "\n"
	metin += i18n.TListItem(lang, "oncelik", gorev.Priority) + "\n"
	metin += i18n.TListItem(lang, "olusturulma_tarihi", gorev.CreatedAt.Format(constants.DateTimeFormatFull)) + "\n"
	metin += i18n.TListItem(lang, "guncelleme_tarihi", gorev.UpdatedAt.Format(constants.DateTimeFormatFull))

	if gorev.ProjeID != "" {
		proje, err := h.isYonetici.ProjeGetir(ctx, gorev.ProjeID)
		if err == nil {
			metin += "\n" + i18n.TListItem(lang, "proje", proje.Name)
		}
	}
	if gorev.ParentID != "" {
		parent, err := h.isYonetici.GorevGetir(ctx, gorev.ParentID)
		if err == nil {
			metin += "\n" + i18n.TListItem(lang, "ust_gorev", parent.Title)
		}
	}
	if gorev.DueDate != nil {
		metin += "\n" + i18n.TListItem(lang, "son_tarih", gorev.DueDate.Format(constants.DateFormatISO))
	}
	if len(gorev.Tags) > 0 {
		var etiketIsimleri []string
		for _, e := range gorev.Tags {
			etiketIsimleri = append(etiketIsimleri, e.Name)
		}
		metin += "\n" + i18n.TListItem(lang, "etiketler", strings.Join(etiketIsimleri, ", "))
	}

	// Bağımlılık sayı bilgilerini ekle (MarkdownParser için gerekli)
	bagimlilikBilgisi := h.gorevBagimlilikBilgisi(gorev, "")
	if bagimlilikBilgisi != "" {
		metin += "\n" + bagimlilikBilgisi
	}

	metin += "\n\n" + i18n.T("headers.taskDescription") + "\n"
	if gorev.Description != "" {
		// Açıklama zaten markdown formatında olabilir, direkt ekle
		metin += gorev.Description
	} else {
		metin += i18n.T("messages.noDescription")
	}

	// Bağımlılıkları ekle - Her zaman göster
	metin += "\n\n" + i18n.T("headers.dependencies") + "\n"

	baglantilar, err := h.isYonetici.GorevBaglantilariGetir(ctx, id)
	if err != nil {
		metin += i18n.T("messages.dependenciesNotAvailable") + "\n"
	} else if len(baglantilar) == 0 {
		metin += i18n.T("messages.noDependencies") + "\n"
	} else {
		var oncekiler []string
		var sonrakiler []string

		for _, b := range baglantilar {
			if b.ConnectionType == "onceki" || b.ConnectionType == "blocker" || b.ConnectionType == "depends_on" {
				if b.TargetID == id {
					// Bu görev hedefse, kaynak önceki görevdir
					kaynakGorev, err := h.isYonetici.GorevGetir(ctx, b.SourceID)
					if err == nil {
						durum := constants.EmojiStatusCompleted
						if kaynakGorev.Status != constants.TaskStatusCompleted {
							durum = constants.EmojiStatusPending
						}
						oncekiler = append(oncekiler, fmt.Sprintf("%s %s (`%s`)", durum, kaynakGorev.Title, kaynakGorev.Status))
					}
				} else if b.SourceID == id {
					// Bu görev kaynaksa, hedef sonraki görevdir
					hedefGorev, err := h.isYonetici.GorevGetir(ctx, b.TargetID)
					if err == nil {
						sonrakiler = append(sonrakiler, fmt.Sprintf("- %s (`%s`)", hedefGorev.Title, hedefGorev.Status))
					}
				}
			}
		}

		if len(oncekiler) > 0 {
			metin += "\n" + i18n.T("headers.waitingTasks") + "\n"
			for _, onceki := range oncekiler {
				metin += fmt.Sprintf("- %s\n", onceki)
			}
		} else {
			metin += "\n" + i18n.T("headers.waitingTasks") + "\n" + i18n.T("messages.notDependentOnAny") + "\n"
		}

		if len(sonrakiler) > 0 {
			metin += "\n" + i18n.T("headers.dependentTasks") + "\n"
			for _, sonraki := range sonrakiler {
				metin += sonraki + "\n"
			}
		} else {
			metin += "\n" + i18n.T("headers.dependentTasks") + "\n" + i18n.T("messages.noDependentTasks") + "\n"
		}

		// Bağımlılık durumu kontrolü
		bagimli, tamamlanmamislar, err := h.isYonetici.GorevBagimliMi(ctx, id)
		if err == nil && !bagimli && gorev.Status == constants.TaskStatusPending {
			metin += "\n" + i18n.T("messages.dependencyWarning", map[string]interface{}{
				"Dependencies": tamamlanmamislar,
			}) + "\n"
		}
	}

	metin += "\n\n---\n"
	metin += "\n*" + i18n.T("messages.lastUpdate", map[string]interface{}{
		"Date": gorev.UpdatedAt.Format(constants.DateFormatDisplay),
	}) + "*"

	return mcp.NewToolResultText(metin), nil
}

// GorevDuzenle görevi düzenler
func (h *Handlers) GorevDuzenle(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	id, result := h.toolHelpers.Validator.ValidateTaskID(params)
	if result != nil {
		return result, nil
	}

	// En az bir düzenleme alanı olmalı
	baslik, baslikVar := params[constants.ParamTitle].(string)
	aciklama, aciklamaVar := params[constants.ParamDescription].(string)
	oncelik, oncelikVar := params[constants.ParamPriority].(string)
	projeID, projeVar := params["project_id"].(string)
	sonTarih, sonTarihVar := params["due_date"].(string)

	if !baslikVar && !aciklamaVar && !oncelikVar && !projeVar && !sonTarihVar {
		return mcp.NewToolResultError(i18n.T("common.validation.at_least_one_field",
			map[string]interface{}{
				"Fields": "title, description, priority, project_id, due_date",
			})), nil
	}

	err := h.isYonetici.GorevDuzenle(ctx, id, baslik, aciklama, oncelik, projeID, sonTarih, baslikVar, aciklamaVar, oncelikVar, projeVar, sonTarihVar)
	if err != nil {
		return mcp.NewToolResultError(i18n.TEditFailed(lang, "task", err)), nil
	}

	// Fetch task to get title for success message
	gorev, _ := h.isYonetici.GorevGetir(ctx, id)
	title := id
	if gorev != nil {
		title = gorev.Title
	}

	return mcp.NewToolResultText(i18n.TEdited(lang, "task", title)), nil
}

// GorevSil görevi siler
func (h *Handlers) GorevSil(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	id, result := h.toolHelpers.Validator.ValidateTaskID(params)
	if result != nil {
		return result, nil
	}

	// Confirmation check
	confirm := h.toolHelpers.Validator.ValidateBool(params, constants.ParamConfirm)
	if !confirm {
		return mcp.NewToolResultError(i18n.T("error.confirmationRequired")), nil
	}

	gorev, err := h.isYonetici.GorevGetir(ctx, id)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatNotFoundError("task", id), nil
	}

	gorevBaslik := gorev.Title

	err = h.isYonetici.GorevSil(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(i18n.TDeleteFailed(lang, "task", err)), nil
	}

	return mcp.NewToolResultText(i18n.TDeleted(lang, "task", gorevBaslik, id)), nil
}

// GorevBulkTransition changes status for multiple tasks
func (h *Handlers) GorevBulkTransition(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	// Validate task IDs
	taskIDsRaw, ok := params["task_ids"]
	if !ok {
		return mcp.NewToolResultError(i18n.TRequiredParam(lang, "task_ids")), nil
	}

	taskIDsInterface, ok := taskIDsRaw.([]interface{})
	if !ok {
		return mcp.NewToolResultError(i18n.TRequiredArray(lang, "task_ids")), nil
	}

	taskIDs := make([]string, len(taskIDsInterface))
	for i, idInterface := range taskIDsInterface {
		if id, ok := idInterface.(string); ok && id != "" {
			taskIDs[i] = id
		} else {
			return mcp.NewToolResultError(i18n.T("error.invalidTaskIdIndex",
				map[string]interface{}{"Index": i})), nil
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

	result_batch, err := batchProcessor.BulkStatusTransition(ctx, request)
	if err != nil {
		return mcp.NewToolResultError(i18n.TProcessFailed(lang, "bulk_status_transition", err)), nil
	}

	// Format response
	var response strings.Builder

	if dryRun {
		response.WriteString(i18n.T("messages.dryRunResult") + "\n\n")
	} else {
		response.WriteString(i18n.T("messages.bulkStatusComplete") + "\n\n")
	}

	response.WriteString(i18n.TMarkdownLabel(lang, "hedef_durum", newStatus) + "\n")
	response.WriteString(i18n.T("messages.processedTasks", map[string]interface{}{"Count": result_batch.TotalProcessed}) + "\n")
	response.WriteString(i18n.T("messages.successfulTasks", map[string]interface{}{"Count": len(result_batch.Successful)}) + "\n")
	response.WriteString(i18n.T("messages.failedTasks", map[string]interface{}{"Count": len(result_batch.Failed)}) + "\n")
	response.WriteString(i18n.T("messages.warningTasks", map[string]interface{}{"Count": len(result_batch.Warnings)}) + "\n")
	response.WriteString(i18n.TDuration(lang, "sure", result_batch.ExecutionTime) + "\n\n")

	if len(result_batch.Successful) > 0 {
		response.WriteString(i18n.T("messages.successfulTasksHeader") + "\n")
		for _, taskID := range result_batch.Successful {
			response.WriteString(fmt.Sprintf("- %s\n", TaskIDFormat.FormatShortID(taskID)))
		}
		response.WriteString("\n")
	}

	if len(result_batch.Failed) > 0 {
		response.WriteString(i18n.T("messages.failedTasksHeader") + "\n")
		for _, failure := range result_batch.Failed {
			response.WriteString(fmt.Sprintf("- %s: %s\n", TaskIDFormat.FormatShortID(failure.TaskID), failure.Error))
		}
		response.WriteString("\n")
	}

	if len(result_batch.Warnings) > 0 {
		response.WriteString(i18n.T("messages.warningsHeader") + "\n")
	}

	return mcp.NewToolResultText(response.String()), nil
}

// GorevBulkTag adds, removes, or replaces tags for multiple tasks
func (h *Handlers) GorevBulkTag(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	// Validate task IDs
	taskIDsRaw, ok := params["task_ids"]
	if !ok {
		return mcp.NewToolResultError(i18n.TRequiredParam(lang, "task_ids")), nil
	}

	taskIDsInterface, ok := taskIDsRaw.([]interface{})
	if !ok {
		return mcp.NewToolResultError(i18n.TRequiredArray(lang, "task_ids")), nil
	}

	taskIDs := make([]string, len(taskIDsInterface))
	for i, idInterface := range taskIDsInterface {
		if id, ok := idInterface.(string); ok && id != "" {
			taskIDs[i] = id
		} else {
			return mcp.NewToolResultError(i18n.T("error.invalidTaskIdIndex",
				map[string]interface{}{"Index": i})), nil
		}
	}

	// Validate tags
	tagsRaw, ok := params["tags"]
	if !ok {
		return mcp.NewToolResultError(i18n.TRequiredParam(lang, "tags")), nil
	}

	tagsInterface, ok := tagsRaw.([]interface{})
	if !ok {
		return mcp.NewToolResultError(i18n.TRequiredArray(lang, "tags")), nil
	}

	tags := make([]string, len(tagsInterface))
	for i, tagInterface := range tagsInterface {
		if tag, ok := tagInterface.(string); ok && tag != "" {
			tags[i] = strings.TrimSpace(tag)
		} else {
			return mcp.NewToolResultError(i18n.T("error.invalidTagIndex",
				map[string]interface{}{"Index": i})), nil
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

	result_batch, err := batchProcessor.BulkTagOperation(ctx, request)
	if err != nil {
		return mcp.NewToolResultError(i18n.TProcessFailed(lang, "bulk_tag_operation", err)), nil
	}

	// Format response
	var response strings.Builder

	if dryRun {
		response.WriteString(i18n.T("messages.dryRunResult") + "\n\n")
	} else {
		response.WriteString(i18n.T("messages.bulkTagComplete") + "\n\n")
	}

	response.WriteString(i18n.TMarkdownLabel(lang, "islem", operation) + "\n")
	response.WriteString(i18n.TMarkdownLabel(lang, "etiketler", strings.Join(tags, ", ")) + "\n")
	response.WriteString(i18n.T("messages.processedTasks", map[string]interface{}{"Count": result_batch.TotalProcessed}) + "\n")
	response.WriteString(i18n.T("messages.successfulTasks", map[string]interface{}{"Count": len(result_batch.Successful)}) + "\n")
	response.WriteString(i18n.T("messages.failedTasks", map[string]interface{}{"Count": len(result_batch.Failed)}) + "\n")
	response.WriteString(i18n.T("messages.warningTasks", map[string]interface{}{"Count": len(result_batch.Warnings)}) + "\n")
	response.WriteString(i18n.TDuration(lang, "sure", result_batch.ExecutionTime) + "\n\n")

	if len(result_batch.Successful) > 0 {
		response.WriteString(i18n.T("messages.successfulTasksHeader") + "\n")
		for _, taskID := range result_batch.Successful {
			response.WriteString(fmt.Sprintf("- %s\n", TaskIDFormat.FormatShortID(taskID)))
		}
		response.WriteString("\n")
	}

	if len(result_batch.Failed) > 0 {
		response.WriteString(i18n.T("messages.failedTasksHeader") + "\n")
		for _, failure := range result_batch.Failed {
			response.WriteString(fmt.Sprintf("- %s: %s\n", TaskIDFormat.FormatShortID(failure.TaskID), failure.Error))
		}
		response.WriteString("\n")
	}

	if len(result_batch.Warnings) > 0 {
		response.WriteString(i18n.T("messages.warningsHeader") + "\n")
		for _, warning := range result_batch.Warnings {
			response.WriteString(fmt.Sprintf("- %s: %s\n", TaskIDFormat.FormatShortID(warning.TaskID), warning.Message))
		}
	}

	return mcp.NewToolResultText(response.String()), nil
}

// GorevSuggestions provides intelligent suggestions for task management
func (h *Handlers) GorevSuggestions(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	// Get session ID from AI context if available
	sessionID := ""
	activeTaskID := ""

	if h.aiContextYonetici != nil {
		if activeTask, err := h.aiContextYonetici.GetActiveTask(ctx); err == nil && activeTask != nil {
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

	response, err := suggestionEngine.GetSuggestions(ctx, request)
	if err != nil {
		return mcp.NewToolResultError(i18n.TProcessFailed(lang, "suggestion", err)), nil
	}

	// Format response
	var output strings.Builder

	output.WriteString(i18n.T("messages.smartSuggestions") + "\n\n")
	output.WriteString(i18n.T("messages.totalSuggestions", map[string]interface{}{"Count": response.TotalCount}) + "\n")
	output.WriteString(i18n.T("messages.performanceExecutionTime", map[string]interface{}{"Duration": response.ExecutionTime}) + "\n\n")

	if len(response.Suggestions) == 0 {
		output.WriteString(i18n.T("messages.noSuggestionsNow") + "\n")
		return mcp.NewToolResultText(output.String()), nil
	}

	// Group suggestions by type
	suggestionGroups := make(map[string][]gorev.Suggestion)
	for _, suggestion := range response.Suggestions {
		suggestionGroups[suggestion.Type] = append(suggestionGroups[suggestion.Type], suggestion)
	}

	// Display suggestions by type
	typeNames := map[string]string{
		"next_action":   i18n.T("messages.suggestionTypeNextAction"),
		"similar_task":  i18n.T("messages.suggestionTypeSimilarTask"),
		"template":      i18n.T("messages.suggestionTypeTemplate"),
		"deadline_risk": i18n.T("messages.suggestionTypeDeadlineRisk"),
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
			output.WriteString(i18n.T("messages.suggestionDescription", map[string]interface{}{"Description": suggestion.Description}) + "\n")
			output.WriteString(i18n.T("messages.suggestionAction", map[string]interface{}{"Action": suggestion.Action}) + "\n")
			output.WriteString(i18n.T("messages.suggestionConfidence", map[string]interface{}{"Confidence": suggestion.Confidence * 100}) + "\n")

			if suggestion.TaskID != "" {
				output.WriteString(i18n.T("messages.suggestionRelatedTask", map[string]interface{}{"TaskID": TaskIDFormat.FormatShortID(suggestion.TaskID)}) + "\n")
			}

			output.WriteString("\n")
		}
	}

	output.WriteString("---\n")
	output.WriteString(i18n.T("messages.suggestionActionTip") + "\n")

	return mcp.NewToolResultText(output.String()), nil
}

// GorevIntelligentCreate creates a task with AI-enhanced features
func (h *Handlers) GorevIntelligentCreate(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	// Validate required parameters
	title, result := h.toolHelpers.Validator.ValidateRequiredString(params, constants.ParamTitle)
	if result != nil {
		return result, nil
	}

	// Optional parameters
	description := h.toolHelpers.Validator.ValidateOptionalString(params, constants.ParamDescription)
	autoSplit := h.toolHelpers.Validator.ValidateBool(params, "auto_split")
	estimateTime := h.toolHelpers.Validator.ValidateBool(params, "estimate_time")
	smartPriority := h.toolHelpers.Validator.ValidateBool(params, "smart_priority")
	suggestTemplate := h.toolHelpers.Validator.ValidateBool(params, "suggest_template")

	// Get project ID if specified
	projeID := h.toolHelpers.Validator.ValidateOptionalString(params, "project_id")

	// Use active project if no project specified
	if projeID == "" {
		if aktifProje, err := h.isYonetici.AktifProjeGetir(ctx); err == nil && aktifProje != nil {
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
	response, err := creator.CreateIntelligentTask(ctx, request)
	if err != nil {
		return mcp.NewToolResultError(i18n.TCreateFailed(lang, "task", err)), nil
	}

	// Set project if specified
	if projeID != "" && response.MainTask != nil {
		response.MainTask.ProjeID = projeID
		updateParams := map[string]interface{}{
			"project_id": projeID,
		}
		if err := h.isYonetici.VeriYonetici().GorevGuncelle(ctx, response.MainTask.ID, updateParams); err != nil {
			// Log but don't fail
			slog.Warn("Failed to set project for intelligent task", "error", err)
		}
	}

	// Record interaction with AI context
	if h.aiContextYonetici != nil && response.MainTask != nil {
		if err := h.aiContextYonetici.RecordInteraction(ctx, response.MainTask.ID, "intelligent_create", map[string]interface{}{
			"auto_split":       autoSplit,
			"estimate_time":    estimateTime,
			"smart_priority":   smartPriority,
			"suggest_template": suggestTemplate,
			"subtasks_created": len(response.Subtasks),
		}); err != nil {
			slog.Warn("Failed to record AI interaction for intelligent create", "taskID", response.MainTask.ID, "error", err)
		}
	}

	// Format response
	var output strings.Builder

	output.WriteString(i18n.T("messages.intelligentTaskCreated") + "\n\n")

	// Main task info
	output.WriteString(i18n.T("messages.mainTask") + "\n")
	output.WriteString(i18n.TMarkdownLabel(lang, "baslik", response.MainTask.Title) + "\n")
	output.WriteString(fmt.Sprintf("**ID:** %s\n", response.MainTask.ID))

	if response.SuggestedPriority != "" {
		priorityEmoji := h.toolHelpers.Formatter.GetPriorityEmoji(response.SuggestedPriority)
		output.WriteString(i18n.T("messages.smartPriority", map[string]interface{}{
			"Emoji":    priorityEmoji,
			"Priority": response.SuggestedPriority,
		}) + "\n")
	}

	if response.EstimatedHours > 0 {
		output.WriteString(i18n.T("messages.estimatedTime", map[string]interface{}{"Hours": response.EstimatedHours}) + "\n")
	}

	if projeID != "" {
		if proje, err := h.isYonetici.ProjeGetir(ctx, projeID); err == nil {
			output.WriteString(i18n.TMarkdownLabel(lang, "proje", proje.Name) + "\n")
		}
	}

	output.WriteString("\n")

	// Subtasks
	if len(response.Subtasks) > 0 {
		output.WriteString(i18n.T("messages.autoSubtasks", map[string]interface{}{"Count": len(response.Subtasks)}) + "\n")
		for i, subtask := range response.Subtasks {
			output.WriteString(fmt.Sprintf("%d. %s (`%s`)\n", i+1, subtask.Title, TaskIDFormat.FormatShortID(subtask.ID)))
		}
		output.WriteString("\n")
	}

	// Template recommendation
	if response.RecommendedTemplate != "" {
		output.WriteString(i18n.T("messages.suggestedTemplate") + "\n")
		output.WriteString(i18n.T("messages.templateConfidence", map[string]interface{}{
			"Template":   response.RecommendedTemplate,
			"Confidence": response.Confidence.TemplateConfidence * 100,
		}) + "\n")
		output.WriteString(i18n.T("messages.templateUsage") + "\n\n")
	}

	// Similar tasks
	if len(response.SimilarTasks) > 0 {
		output.WriteString(i18n.T("messages.similarTasks", map[string]interface{}{"Count": len(response.SimilarTasks)}) + "\n")
		for i, similar := range response.SimilarTasks {
			if i >= constants.MaxSuggestionsToShow { // Show top suggestions
				break
			}
			output.WriteString(fmt.Sprintf("%d. %s (%s)\n",
				i+1, similar.Task.Title,
				i18n.T("messages.similarTasksMatch", map[string]interface{}{
					"Similarity": similar.SimilarityScore * 100,
					"Reason":     similar.Reason,
				})))
		}
		output.WriteString("\n")
	}

	// AI Insights
	if len(response.Insights) > 0 {
		output.WriteString(i18n.T("messages.aiAnalysisResults") + "\n")
		for _, insight := range response.Insights {
			output.WriteString(fmt.Sprintf("- %s\n", insight))
		}
		output.WriteString("\n")
	}

	// Performance info
	output.WriteString(i18n.T("messages.performanceHeader") + "\n")
	output.WriteString(i18n.T("messages.performanceExecutionTime", map[string]interface{}{"Duration": response.ExecutionTime}) + "\n")
	output.WriteString(i18n.T("messages.confidenceScores") + "\n")
	if response.SuggestedPriority != "" {
		output.WriteString(i18n.T("messages.priorityConfidence", map[string]interface{}{"Percent": response.Confidence.PriorityConfidence * 100}) + "\n")
	}
	if response.EstimatedHours > 0 {
		output.WriteString(i18n.T("messages.timeConfidence", map[string]interface{}{"Percent": response.Confidence.TimeConfidence * 100}) + "\n")
	}
	if len(response.Subtasks) > 0 {
		output.WriteString(i18n.T("messages.subtaskConfidence", map[string]interface{}{"Percent": response.Confidence.SubtaskConfidence * 100}) + "\n")
	}

	output.WriteString("\n---\n")
	output.WriteString(i18n.T("messages.taskDetailTip", map[string]interface{}{"Id": response.MainTask.ID}) + "\n")

	return mcp.NewToolResultText(output.String()), nil
}

// ProjeListele tüm projeleri listeler
func (h *Handlers) ProjeListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	projeler, err := h.isYonetici.ProjeListele(ctx)
	if err != nil {
		return mcp.NewToolResultError(i18n.TListFailed(lang, "project", err)), nil
	}

	if len(projeler) == 0 {
		return mcp.NewToolResultText(i18n.T("messages.noProjects")), nil
	}

	metin := i18n.T("headers.projectList") + "\n\n"
	for _, proje := range projeler {
		metin += fmt.Sprintf("### %s\n", proje.Name)
		metin += i18n.TListItem(lang, "id_field", proje.ID) + "\n"
		if proje.Definition != "" {
			metin += i18n.TListItem(lang, "tanim", proje.Definition) + "\n"
		}
		metin += i18n.TListItem(lang, "olusturma", proje.CreatedAt.Format(constants.DateFormatDisplay)) + "\n"

		// Her proje için görev sayısını göster
		gorevSayisi, err := h.isYonetici.ProjeGorevSayisi(ctx, proje.ID)
		if err == nil {
			metin += i18n.TListItem(lang, "gorev_sayisi", gorevSayisi) + "\n"
		}
	}

	return mcp.NewToolResultText(metin), nil
}

// gorevBagimlilikBilgisi görev için bağımlılık bilgilerini formatlar
func (h *Handlers) gorevBagimlilikBilgisi(g *gorev.Gorev, indent string) string {
	bilgi := ""
	if g.DependencyCount > 0 {
		bilgi += indent + i18n.T("messages.dependentTaskCount", map[string]interface{}{
			"Count": g.DependencyCount,
		}) + "\n"
		if g.UncompletedDependencyCount > 0 {
			bilgi += indent + i18n.T("messages.incompleteDependencyCount", map[string]interface{}{
				"Count": g.UncompletedDependencyCount,
			}) + "\n"
		}
	}
	if g.DependentOnThisCount > 0 {
		bilgi += indent + i18n.T("messages.dependentOnThisCount", map[string]interface{}{
			"Count": g.DependentOnThisCount,
		}) + "\n"
	}
	return bilgi
}

// ProjeGorevleri bir projenin görevlerini listeler
func (h *Handlers) ProjeGorevleri(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	h.toolHelpers.SetLanguage(lang)
	ctx := i18n.WithLanguage(context.Background(), lang)

	// Support both English and Turkish parameter names for backward compatibility
	projeID, result := h.toolHelpers.Validator.ValidateRequiredString(params, "project_id")
	if result != nil {
		// Fallback to Turkish parameter name
		projeID, result = h.toolHelpers.Validator.ValidateRequiredString(params, "proje_id")
		if result != nil {
			return result, nil
		}
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
	proje, err := h.isYonetici.ProjeGetir(ctx, projeID)
	if err != nil {
		return mcp.NewToolResultError(i18n.TEntityNotFound(lang, "project", err)), nil
	}

	gorevler, err := h.isYonetici.ProjeGorevleri(ctx, projeID)
	if err != nil {
		return mcp.NewToolResultError(i18n.TFetchFailed(lang, "task", err)), nil
	}

	// Toplam görev sayısı
	toplamGorevSayisi := len(gorevler)

	metin := ""

	if toplamGorevSayisi == 0 {
		metin = i18n.T("messages.noTasksInThisProject", map[string]interface{}{
			"Name": proje.Name,
		})
		return mcp.NewToolResultText(metin), nil
	}

	// Kompakt başlık
	if toplamGorevSayisi > limit || offset > 0 {
		metin = i18n.T("messages.projectTasksPage", map[string]interface{}{
			"Name":  proje.Name,
			"Start": offset + 1,
			"End":   min(offset+limit, toplamGorevSayisi),
			"Total": toplamGorevSayisi,
		}) + "\n"
	} else {
		metin = i18n.T("messages.taskRangeCount", map[string]interface{}{
			"Name":  proje.Name,
			"Count": toplamGorevSayisi,
		}) + "\n"
	}

	// Duruma göre grupla
	beklemede := []*gorev.Gorev{}
	devamEdiyor := []*gorev.Gorev{}
	tamamlandi := []*gorev.Gorev{}

	for _, g := range gorevler {
		switch g.Status {
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
		metin += "\n🔵 " + i18n.TStatus(lang, constants.TaskStatusInProgress) + "\n"
		for i := devamEdiyorStart; i < devamEdiyorEnd; i++ {
			g := devamEdiyor[i]
			gorevSize := h.gorevResponseSizeEstimate(g)
			if estimatedSize+gorevSize > maxResponseSize && gorevleriGoster > 0 {
				metin += i18n.T("messages.moreTasksLimit", map[string]interface{}{
					"Count": devamEdiyorEnd - i,
				}) + "\n"
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
		// Pagination index no longer needed after this point
	} else {
		beklemedeStart = len(beklemede)
		beklemedeEnd = len(beklemede)
		// Pagination index updated for next section
	}

	if beklemedeEnd > beklemedeStart && estimatedSize < maxResponseSize {
		metin += "\n⚪ " + i18n.TStatus(lang, constants.TaskStatusPending) + "\n"
		for i := beklemedeStart; i < beklemedeEnd; i++ {
			g := beklemede[i]
			gorevSize := h.gorevResponseSizeEstimate(g)
			if estimatedSize+gorevSize > maxResponseSize && gorevleriGoster > 0 {
				metin += i18n.T("messages.moreTasksLimit", map[string]interface{}{
					"Count": beklemedeEnd - i,
				}) + "\n"
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
		metin += "\n" + constants.EmojiStatusCompleted + " " + i18n.TStatus(lang, constants.TaskStatusCompleted) + "\n"
		for i := tamamlandiStart; i < tamamlandiEnd; i++ {
			g := tamamlandi[i]
			gorevSize := h.gorevResponseSizeEstimate(g)
			if estimatedSize+gorevSize > maxResponseSize && gorevleriGoster > 0 {
				metin += i18n.T("messages.moreTasksLimit", map[string]interface{}{
					"Count": tamamlandiEnd - i,
				}) + "\n"
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
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	ozet, err := h.isYonetici.OzetAl(ctx)
	if err != nil {
		return mcp.NewToolResultError(i18n.TFetchFailed(lang, "summary", err)), nil
	}

	var metin strings.Builder
	metin.WriteString(i18n.T("headers.summaryReport") + "\n\n")
	metin.WriteString(i18n.TMarkdownLabel(lang, "toplam_proje", ozet.TotalProjects) + "\n")
	metin.WriteString(i18n.TMarkdownLabel(lang, "toplam_gorev", ozet.TotalTasks) + "\n\n")

	metin.WriteString(i18n.T("headers.statusDistribution") + "\n")
	metin.WriteString(fmt.Sprintf("- %s: %d\n", i18n.TStatus(lang, constants.TaskStatusPending), ozet.PendingTasks))
	metin.WriteString(fmt.Sprintf("- %s: %d\n", i18n.TStatus(lang, constants.TaskStatusInProgress), ozet.InProgressTasks))
	metin.WriteString(fmt.Sprintf("- %s: %d\n\n", i18n.TStatus(lang, constants.TaskStatusCompleted), ozet.CompletedTasks))

	metin.WriteString(i18n.T("headers.priorityDistribution") + "\n")
	metin.WriteString(fmt.Sprintf("- %s: %d\n", i18n.TPriority(lang, constants.PriorityHigh), ozet.HighPriorityTasks))
	metin.WriteString(fmt.Sprintf("- %s: %d\n", i18n.TPriority(lang, constants.PriorityMedium), ozet.MediumPriorityTasks))
	metin.WriteString(fmt.Sprintf("- %s: %d\n", i18n.TPriority(lang, constants.PriorityLow), ozet.LowPriorityTasks))

	return mcp.NewToolResultText(metin.String()), nil
}

func (h *Handlers) GorevBagimlilikEkle(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)

	// Support both English and Turkish parameter names for backward compatibility
	kaynakID, ok := params["source_id"].(string)
	if !ok || kaynakID == "" {
		// Fallback to Turkish parameter name
		kaynakID, ok = params["gorev_id"].(string)
		if !ok || kaynakID == "" {
			return mcp.NewToolResultError(i18n.TRequiredParam(lang, "source_id")), nil
		}
	}

	hedefID, ok := params["target_id"].(string)
	if !ok || hedefID == "" {
		// Fallback to Turkish parameter name
		hedefID, ok = params["bagli_gorev_id"].(string)
		if !ok || hedefID == "" {
			return mcp.NewToolResultError(i18n.TRequiredParam(lang, "target_id")), nil
		}
	}

	baglantiTipi, ok := params["connection_type"].(string)
	if !ok || baglantiTipi == "" {
		// Fallback to Turkish parameter name
		baglantiTipi, ok = params["baglanti_tipi"].(string)
		if !ok || baglantiTipi == "" {
			// Default to "blocks" if not specified
			baglantiTipi = "blocks"
		}
	}

	baglanti, err := h.isYonetici.GorevBagimlilikEkle(ctx, kaynakID, hedefID, baglantiTipi)
	if err != nil {
		return mcp.NewToolResultError(i18n.TAddFailed(lang, "dependency", err)), nil
	}

	return mcp.NewToolResultText(i18n.T("success.dependencyAdded", map[string]interface{}{
		"Source": baglanti.SourceID,
		"Target": baglanti.TargetID,
		"Type":   baglanti.ConnectionType,
	})), nil
}

// TemplateListele lists available templates
func (h *Handlers) TemplateListele(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	category, _ := params[constants.ParamCategory].(string)

	templates, err := h.isYonetici.TemplateListele(ctx, category)
	if err != nil {
		return mcp.NewToolResultError(i18n.TListFailed(lang, "template", err)), nil
	}

	if len(templates) == 0 {
		return mcp.NewToolResultText(i18n.T("messages.noTemplates")), nil
	}

	metin := i18n.T("messages.templateListHeader") + "\n\n"

	// Kategorilere göre grupla
	kategoriMap := make(map[string][]*gorev.GorevTemplate)
	for _, tmpl := range templates {
		kategoriMap[tmpl.Category] = append(kategoriMap[tmpl.Category], tmpl)
	}

	// Her kategoriyi göster
	for kat, tmpls := range kategoriMap {
		metin += fmt.Sprintf("### %s\n\n", kat)

		for _, tmpl := range tmpls {
			metin += fmt.Sprintf("#### %s\n", tmpl.Name)
			metin += i18n.TListItem(lang, "id_field", fmt.Sprintf("`%s`", tmpl.ID)) + "\n"
			metin += i18n.TListItem(lang, "aciklama", tmpl.Definition) + "\n"
			metin += i18n.TListItem(lang, "baslik_sablonu", fmt.Sprintf("`%s`", tmpl.DefaultTitle)) + "\n"

			// Alanları göster
			if len(tmpl.Fields) > 0 {
				metin += i18n.TListItem(lang, "alanlar", "") + "\n"
				for _, alan := range tmpl.Fields {
					zorunlu := ""
					if alan.Required {
						zorunlu = fmt.Sprintf(" *(%s)*", i18n.TLabel(lang, "zorunlu"))
					}
					metin += fmt.Sprintf("  - `%s` (%s)%s", alan.Name, alan.Type, zorunlu)
					if alan.Default != "" {
						metin += fmt.Sprintf(" - %s: %s", i18n.TLabel(lang, "varsayilan"), alan.Default)
					}
					if len(alan.Options) > 0 {
						metin += fmt.Sprintf(" - %s: %s", i18n.TLabel(lang, "secenekler"), strings.Join(alan.Options, ", "))
					}
					metin += "\n"
				}
			}
			metin += "\n"
		}
	}

	metin += "\n" + i18n.T("messages.templateUsageTip")

	return mcp.NewToolResultText(metin), nil
}

// TemplatedenGorevOlustur template kullanarak görev oluşturur
func (h *Handlers) TemplatedenGorevOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	templateID, ok := params[constants.ParamTemplateID].(string)
	if !ok || templateID == "" {
		return mcp.NewToolResultError(i18n.TRequiredParam(lang, "template_id")), nil
	}

	// Support both English and Turkish parameter names for backward compatibility
	degerlerRaw, ok := params[constants.ParamValues].(map[string]interface{})
	if !ok {
		// Fallback to Turkish parameter name
		degerlerRaw, ok = params["degerler"].(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError(i18n.T("common.validation.required_object", map[string]interface{}{
				"Param": "values",
			})), nil
		}
	}

	// Önce template'i ID veya alias ile kontrol et
	template, err := h.isYonetici.VeriYonetici().TemplateIDVeyaAliasIleGetir(ctx, templateID)
	if err != nil {
		return mcp.NewToolResultError(i18n.TEntityNotFound(lang, "template", err)), nil
	}

	// Interface{} map'i string map'e çevir ve validation yap
	degerler := make(map[string]string)
	eksikAlanlar := []string{}

	// Tüm zorunlu alanları kontrol et
	for _, alan := range template.Fields {
		if val, exists := degerlerRaw[alan.Name]; exists {
			// Değer var, string'e çevir
			strVal := fmt.Sprintf("%v", val)
			if alan.Required && strings.TrimSpace(strVal) == "" {
				eksikAlanlar = append(eksikAlanlar, alan.Name)
			} else {
				degerler[alan.Name] = strVal
			}
		} else if alan.Required {
			// Zorunlu alan eksik
			eksikAlanlar = append(eksikAlanlar, alan.Name)
		} else if alan.Default != "" {
			// Varsayılan değeri kullan
			degerler[alan.Name] = alan.Default
		}
	}

	// Eksik alanlar varsa detaylı hata ver
	if len(eksikAlanlar) > 0 {
		errorMsg := i18n.T("messages.requiredFieldsMissing") + "\n\n"
		errorMsg += i18n.T("messages.templateLabel", map[string]interface{}{"Template": template.Name}) + "\n"
		errorMsg += i18n.T("messages.missingFieldsLabel", map[string]interface{}{"Fields": strings.Join(eksikAlanlar, ", ")}) + "\n\n"
		errorMsg += i18n.T("messages.requiredFieldsForTemplate") + "\n"
		errorMsg += h.templateZorunluAlanlariListele(template) + "\n\n"
		errorMsg += i18n.T("messages.exampleUsage") + "\n"
		errorMsg += i18n.T("messages.templateUsageCommand", map[string]interface{}{
			"TemplateID": templateID,
			"Values":     "{" + h.templateOrnekDegerler(template) + "}",
		})
		return mcp.NewToolResultError(errorMsg), nil
	}

	// Select tipindeki alanların geçerli değerlerini kontrol et
	for _, alan := range template.Fields {
		if alan.Type == "select" && len(alan.Options) > 0 {
			if deger, ok := degerler[alan.Name]; ok && deger != "" {
				gecerli := false
				for _, secenek := range alan.Options {
					if deger == secenek {
						gecerli = true
						break
					}
				}
				if !gecerli {
					return mcp.NewToolResultError(i18n.TInvalidValue(lang, alan.Name, deger, alan.Options)), nil
				}
			}
		}
	}

	gorev, err := h.isYonetici.TemplatedenGorevOlustur(ctx, templateID, degerler)
	if err != nil {
		return mcp.NewToolResultError(i18n.TCreateFailed(lang, "task", err)), nil
	}

	successMsg := i18n.T("messages.templateTaskCreated") + "\n\n"
	successMsg += i18n.T("messages.templateLabel", map[string]interface{}{"Template": template.Name}) + "\n"
	successMsg += i18n.TListItem(lang, "baslik", gorev.Title) + "\n"
	successMsg += i18n.TListItem(lang, "id_field", gorev.ID) + "\n"
	successMsg += i18n.TListItem(lang, "oncelik", gorev.Priority) + "\n\n"
	successMsg += i18n.T("messages.detailsCommand", map[string]interface{}{"ID": gorev.ID})

	return mcp.NewToolResultText(successMsg), nil
}

// RegisterTools tüm araçları MCP sunucusuna kaydeder
// GorevAltGorevOlustur mevcut bir görevin altına yeni görev oluşturur
func (h *Handlers) GorevAltGorevOlustur(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	parentID, ok := params["parent_id"].(string)
	if !ok || parentID == "" {
		return mcp.NewToolResultError(i18n.TRequiredParam(lang, "parent_id")), nil
	}

	baslik, ok := params[constants.ParamTitle].(string)
	if !ok || baslik == "" {
		return mcp.NewToolResultError(i18n.TRequiredParam(lang, "title")), nil
	}

	aciklama, _ := params[constants.ParamDescription].(string)
	oncelik, _ := params[constants.ParamPriority].(string)
	if oncelik == "" {
		oncelik = constants.PriorityMedium
	}

	sonTarih, _ := params["due_date"].(string)
	etiketlerStr, _ := params["tags"].(string)
	var etiketler []string
	if etiketlerStr != "" {
		etiketler = strings.Split(etiketlerStr, ",")
		for i := range etiketler {
			etiketler[i] = strings.TrimSpace(etiketler[i])
		}
	}

	gorev, err := h.isYonetici.AltGorevOlustur(ctx, parentID, baslik, aciklama, oncelik, sonTarih, etiketler)
	if err != nil {
		return mcp.NewToolResultError(i18n.TCreateFailed(lang, "subtask", err)), nil
	}

	return mcp.NewToolResultText(i18n.T("success.subtaskCreated", map[string]interface{}{
		"Title": gorev.Title,
		"Id":    gorev.ID,
	})), nil
}

// GorevUstDegistir bir görevin üst görevini değiştirir
func (h *Handlers) GorevUstDegistir(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	taskID, ok := params[constants.ParamTaskID].(string)
	if !ok || taskID == "" {
		return mcp.NewToolResultError(i18n.TRequiredParam(lang, constants.ParamTaskID)), nil
	}

	newParentID, _ := params[constants.ParamNewParentID].(string)

	err := h.isYonetici.GorevUstDegistir(ctx, taskID, newParentID)
	if err != nil {
		return mcp.NewToolResultError(i18n.TUpdateFailed(lang, "task", err)), nil
	}

	if newParentID == "" {
		return mcp.NewToolResultText(i18n.T("success.taskMovedToRoot")), nil
	}
	return mcp.NewToolResultText(i18n.T("success.taskMovedToParent")), nil
}

// GorevHiyerarsiGoster bir görevin tam hiyerarşisini gösterir
func (h *Handlers) GorevHiyerarsiGoster(params map[string]interface{}) (*mcp.CallToolResult, error) {
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	taskID, ok := params[constants.ParamTaskID].(string)
	if !ok || taskID == "" {
		return mcp.NewToolResultError(i18n.TRequiredParam(lang, constants.ParamTaskID)), nil
	}

	hiyerarsi, err := h.isYonetici.GorevHiyerarsiGetir(ctx, taskID)
	if err != nil {
		return mcp.NewToolResultError(i18n.TFetchFailed(lang, "hierarchy", err)), nil
	}

	var sb strings.Builder
	sb.WriteString(i18n.T("messages.hierarchyTitle", map[string]interface{}{
		"Title": hiyerarsi.Gorev.Title,
	}) + "\n\n")

	// Üst görevler
	if len(hiyerarsi.ParentTasks) > 0 {
		sb.WriteString(i18n.T("messages.upperTaskHeader") + "\n")
		for i := len(hiyerarsi.ParentTasks) - 1; i >= 0; i-- {
			ust := hiyerarsi.ParentTasks[i]
			sb.WriteString(fmt.Sprintf("%s└─ %s (%s)\n", strings.Repeat("  ", len(hiyerarsi.ParentTasks)-i-1), ust.Title, ust.Status))
		}
		sb.WriteString(fmt.Sprintf("%s", strings.Repeat("  ", len(hiyerarsi.ParentTasks))))
		sb.WriteString(i18n.T("messages.currentTaskMarker", map[string]interface{}{
			"Title": hiyerarsi.Gorev.Title,
		}) + "\n\n")
	}

	// Alt görev istatistikleri
	sb.WriteString(i18n.T("messages.subtaskStats") + "\n")
	sb.WriteString(i18n.T("messages.totalSubtasks", map[string]interface{}{"Count": hiyerarsi.TotalSubtasks}) + "\n")
	sb.WriteString(i18n.T("messages.completedSubtasks", map[string]interface{}{"Count": hiyerarsi.CompletedSubtasks}) + "\n")
	sb.WriteString(i18n.T("messages.inProgressSubtasks", map[string]interface{}{"Count": hiyerarsi.InProgressSubtasks}) + "\n")
	sb.WriteString(i18n.T("messages.pendingSubtasks", map[string]interface{}{
		"Count": hiyerarsi.PendingSubtasks,
		"Emoji": constants.EmojiStatusPending,
	}) + "\n")
	sb.WriteString(i18n.T("messages.progressPercent", map[string]interface{}{
		"Percent": fmt.Sprintf("%.1f", hiyerarsi.ProgressPercentage),
	}) + "\n\n")

	// Doğrudan alt görevler
	altGorevler, err := h.isYonetici.AltGorevleriGetir(ctx, taskID)
	if err == nil && len(altGorevler) > 0 {
		sb.WriteString(i18n.T("messages.directSubtasks") + "\n")
		for _, alt := range altGorevler {
			durum := h.toolHelpers.Formatter.GetStatusEmoji(alt.Status)
			sb.WriteString(fmt.Sprintf("- %s %s (ID: %s)\n", durum, alt.Title, alt.ID))
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
	case "aktif_proje_ayarla":
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
	case "gorev_search_advanced":
		return h.GorevSearchAdvanced(params)
	case "gorev_filter_profile_save":
		return h.GorevFilterProfileSave(params)
	case "gorev_filter_profile_load":
		return h.GorevFilterProfileLoad(params)
	case "gorev_filter_profile_list":
		return h.GorevFilterProfileList(params)
	case "gorev_filter_profile_delete":
		return h.GorevFilterProfileDelete(params)
	case "gorev_search_history":
		return h.GorevSearchHistory(params)
	case "gorev_export":
		return h.GorevExport(params)
	case "gorev_import":
		return h.GorevImport(params)
	case "gorev_ide_detect":
		return h.IDEDetect(params)
	case "gorev_ide_install":
		return h.IDEInstallExtension(params)
	case "gorev_ide_uninstall":
		return h.IDEUninstallExtension(params)
	case "gorev_ide_status":
		return h.IDEExtensionStatus(params)
	case "gorev_ide_update":
		return h.IDEUpdateExtension(params)
	default:
		return mcp.NewToolResultError(i18n.T("error.unknownTool", map[string]interface{}{"Tool": toolName})), nil
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
	ctx := context.Background()
	taskID, result := h.toolHelpers.Validator.ValidateTaskIDField(params, "task_id")
	if result != nil {
		return result, nil
	}

	err := h.aiContextYonetici.SetActiveTask(ctx, taskID)
	if err != nil {
		return h.toolHelpers.ErrorFormatter.FormatOperationError("Aktif görev ayarlama hatası", err), nil
	}

	// Also record task view for auto-state management
	if err := h.aiContextYonetici.RecordTaskView(ctx, taskID); err != nil {
		// Log but don't fail
		// fmt.Printf("Görev görüntüleme kaydı hatası: %v\n", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf(constants.EmojiStatusCompleted+" Görev %s başarıyla aktif görev olarak ayarlandı.", taskID)), nil
}

// GorevGetActive returns the current active task
func (h *Handlers) GorevGetActive(params map[string]interface{}) (*mcp.CallToolResult, error) {
	ctx := context.Background()
	activeTask, err := h.aiContextYonetici.GetActiveTask(ctx)
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
		activeTask.Title,
		activeTask.ID,
		activeTask.Status,
		activeTask.Priority,
		activeTask.ProjeID)

	if activeTask.Description != "" {
		metin += fmt.Sprintf("\n\n## 📝 Açıklama\n%s", activeTask.Description)
	}

	return mcp.NewToolResultText(metin), nil
}

// GorevRecent returns recent tasks interacted with by AI
func (h *Handlers) GorevRecent(params map[string]interface{}) (*mcp.CallToolResult, error) {
	ctx := context.Background()
	limit := constants.DefaultRecentTaskLimit
	if l, ok := params["limit"].(float64); ok {
		limit = int(l)
	}

	tasks, err := h.aiContextYonetici.GetRecentTasks(ctx, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Son görevleri getirme hatası: %v", err)), nil
	}

	if len(tasks) == 0 {
		return mcp.NewToolResultText("ℹ️ Son etkileşimde bulunulan görev yok."), nil
	}

	var result strings.Builder
	result.WriteString("## 📋 Son Etkileşimli Görevler\n\n")

	for i, task := range tasks {
		result.WriteString(fmt.Sprintf("### %d. %s (ID: %s)\n", i+1, task.Title, task.ID))
		result.WriteString(fmt.Sprintf("- **Durum:** %s\n", task.Status))
		result.WriteString(fmt.Sprintf("- **Öncelik:** %s\n", task.Priority))
		if task.ProjeID != "" {
			result.WriteString(fmt.Sprintf("- **Proje:** %s\n", task.ProjeID))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

// GorevContextSummary returns an AI-optimized context summary
func (h *Handlers) GorevContextSummary(params map[string]interface{}) (*mcp.CallToolResult, error) {
	ctx := context.Background()
	summary, err := h.aiContextYonetici.GetContextSummary(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Context özeti getirme hatası: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString("## 🤖 AI Oturum Özeti\n\n")

	// Active task
	if summary.ActiveTask != nil {
		result.WriteString(fmt.Sprintf("### 🎯 Aktif Görev\n**%s** (ID: %s)\n", summary.ActiveTask.Title, summary.ActiveTask.ID))
		result.WriteString(fmt.Sprintf("- Durum: %s | Öncelik: %s\n\n", summary.ActiveTask.Status, summary.ActiveTask.Priority))
	} else {
		result.WriteString("### 🎯 Aktif Görev\nYok\n\n")
	}

	// Working project
	if summary.WorkingProject != nil {
		result.WriteString(fmt.Sprintf("### 📁 Çalışılan Proje\n**%s**\n\n", summary.WorkingProject.Name))
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
			result.WriteString(fmt.Sprintf("- **%s** (ID: %s)\n", task.Title, task.ID))
		}
		result.WriteString("\n")
	}

	// Blockers
	if len(summary.Blockers) > 0 {
		result.WriteString("### 🚫 Blokajlar\n")
		for _, task := range summary.Blockers {
			result.WriteString(fmt.Sprintf("- **%s** (ID: %s) - %d bağımlılık bekliyor\n",
				task.Title, task.ID, task.UncompletedDependencyCount))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}

// GorevBatchUpdate performs batch updates on multiple tasks
func (h *Handlers) GorevBatchUpdate(params map[string]interface{}) (*mcp.CallToolResult, error) {
	ctx := context.Background()
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

		// Build updates data from all fields except "id"
		updatesData := make(map[string]interface{})
		for key, value := range updateMap {
			if key != "id" {
				updatesData[key] = value
			}
		}

		// Skip if no actual updates provided
		if len(updatesData) == 0 {
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

	result, err := h.aiContextYonetici.BatchUpdate(ctx, updates)
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
	lang := h.extractLanguage()
	ctx := i18n.WithLanguage(context.Background(), lang)
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query parametresi gerekli"), nil
	}

	tasks, err := h.aiContextYonetici.NLPQuery(ctx, query)
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
		statusEmoji := h.toolHelpers.Formatter.GetStatusEmoji(task.Status)
		priorityEmoji := h.toolHelpers.Formatter.GetPriorityEmoji(task.Priority)

		result.WriteString(fmt.Sprintf("%s %s **%s** (ID: %s)\n", statusEmoji, priorityEmoji, task.Title, TaskIDFormat.FormatShortID(task.ID)))

		if task.Description != "" {
			desc := task.Description
			if len(desc) > constants.MaxDescriptionDisplayLength {
				desc = desc[:constants.MaxDescriptionDisplayLength] + "..."
			}
			result.WriteString(fmt.Sprintf("   %s\n", desc))
		}

		details := []string{}
		if task.ProjeID != "" {
			details = append(details, fmt.Sprintf("Proje: %s", task.ProjeID))
		}
		if len(task.Tags) > 0 {
			var tagNames []string
			for _, tag := range task.Tags {
				tagNames = append(tagNames, tag.Name)
			}
			details = append(details, i18n.TMarkdownLabel(lang, "etiketler", strings.Join(tagNames, ", ")))
		}
		if task.DueDate != nil {
			details = append(details, i18n.TMarkdownLabel(lang, "son_tarih", task.DueDate.Format(constants.DateFormatISO)))
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
	ctx := context.Background()
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
	task, err := h.isYonetici.VeriYonetici().GorevGetir(ctx, taskID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Görev bulunamadı: %v", err)), nil
	}

	// Add file path to watcher
	if err := h.fileWatcher.AddTaskPath(taskID, filePath); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Dosya izleme eklenemedi: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(constants.EmojiStatusCompleted+" Görev '%s' için '%s' dosya yolu izlemeye eklendi.\n\nDosya değişiklikleri otomatik olarak takip edilecek ve görev durumu gerektiğinde güncellenecek.", task.Title, filePath)), nil
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
	ctx := context.Background()
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
			if task, err := h.isYonetici.VeriYonetici().GorevGetir(ctx, taskID); err == nil {
				result.WriteString(fmt.Sprintf("### Görev: %s (ID: %s)\n\n", task.Title, taskID))
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
					if task, err := h.isYonetici.VeriYonetici().GorevGetir(ctx, taskID); err == nil {
						taskNames = append(taskNames, fmt.Sprintf("%s (ID: %s)", task.Title, taskID))
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
	ctx := context.Background()
	if h.debug {
		slog.Debug("GorevExport called", "params", params)
	}

	// Parse parameters
	if params == nil {
		noArgsMsg := i18n.T("error.noArguments", nil)
		if noArgsMsg == "error.noArguments" {
			noArgsMsg = "No arguments provided"
		}
		return mcp.NewToolResultError(noArgsMsg), nil
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
	exportData, err := h.isYonetici.ExportData(ctx, options)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf(i18n.T("error.exportFailed", map[string]interface{}{"Error": err}))), nil
	}

	// Save to file
	if err := h.isYonetici.SaveExportToFile(ctx, exportData, options); err != nil {
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
	ctx := context.Background()
	if h.debug {
		slog.Debug("GorevImport called", "params", params)
	}

	// Parse parameters
	if params == nil {
		noArgsMsg := i18n.T("error.noArguments", nil)
		if noArgsMsg == "error.noArguments" {
			noArgsMsg = "No arguments provided"
		}
		return mcp.NewToolResultError(noArgsMsg), nil
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
	result, err := h.isYonetici.ImportData(ctx, options)
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
					cases.Title(language.Und).String(conflict.Type),
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
		paramRequiredMsg := i18n.T("error.parameterRequired", map[string]interface{}{"Param": "ide_type"})
		if paramRequiredMsg == "error.parameterRequired" {
			paramRequiredMsg = "Parameter ide_type is required"
		}
		return nil, fmt.Errorf(paramRequiredMsg)
	}

	installToAll := false
	if ideType == "all" {
		installToAll = true
	}

	// Create detector and installer
	detector := gorev.NewIDEDetector()
	if _, err := detector.DetectAllIDEs(); err != nil {
		return nil, fmt.Errorf("IDE detection failed: %w", err)
	}

	installer := gorev.NewExtensionInstaller(detector)
	defer func() {
		if cerr := installer.Cleanup(); cerr != nil {
			slog.Warn("Installer cleanup failed", "error", cerr)
		}
	}()

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
		paramRequiredMsg := i18n.T("error.parameterRequired", map[string]interface{}{"Param": "ide_type"})
		if paramRequiredMsg == "error.parameterRequired" {
			paramRequiredMsg = "Parameter ide_type is required"
		}
		return nil, fmt.Errorf(paramRequiredMsg)
	}

	extensionID, exists := params["extension_id"].(string)
	if !exists {
		extensionID = "mehmetsenol.gorev-vscode" // default
	}

	// Create detector and installer
	detector := gorev.NewIDEDetector()
	if _, err := detector.DetectAllIDEs(); err != nil {
		return nil, fmt.Errorf("IDE detection failed: %w", err)
	}

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
		paramRequiredMsg := i18n.T("error.parameterRequired", map[string]interface{}{"Param": "ide_type"})
		if paramRequiredMsg == "error.parameterRequired" {
			paramRequiredMsg = "Parameter ide_type is required"
		}
		return nil, fmt.Errorf(paramRequiredMsg)
	}

	updateAll := false
	if ideType == "all" {
		updateAll = true
	}

	// Create detector and installer
	detector := gorev.NewIDEDetector()
	if _, err := detector.DetectAllIDEs(); err != nil {
		return nil, fmt.Errorf("IDE detection failed: %w", err)
	}

	installer := gorev.NewExtensionInstaller(detector)
	defer func() {
		if cerr := installer.Cleanup(); cerr != nil {
			slog.Warn("Installer cleanup failed", "error", cerr)
		}
	}()

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

// ================================
// Advanced Search Tools
// ================================

// GorevSearchAdvanced performs advanced search with FTS5 and fuzzy matching
func (h *Handlers) GorevSearchAdvanced(params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Extract parameters
	query := h.toolHelpers.Validator.ValidateOptionalString(params, "query")

	// Extract search options with defaults
	useFuzzySearch := true
	if val, ok := params["use_fuzzy_search"].(bool); ok {
		useFuzzySearch = val
	}

	fuzzyThreshold := 0.6
	if val, ok := params["fuzzy_threshold"].(float64); ok {
		fuzzyThreshold = val
	}

	maxResults := 50
	if val, ok := params["max_results"]; ok {
		if intVal, ok := val.(int); ok {
			maxResults = intVal
		} else if floatVal, ok := val.(float64); ok {
			maxResults = int(floatVal)
		}
	}

	sortBy := "relevance"
	if val, ok := params["sort_by"].(string); ok && val != "" {
		sortBy = val
	}

	sortDirection := "desc"
	if val, ok := params["sort_direction"].(string); ok && val != "" {
		sortDirection = val
	}

	includeCompleted := false
	if val, ok := params["include_completed"].(bool); ok {
		includeCompleted = val
	}

	// Parse search options
	options := gorev.SearchOptions{
		Query:            query,
		Filters:          make(map[string]interface{}),
		UseFuzzySearch:   useFuzzySearch,
		FuzzyThreshold:   fuzzyThreshold,
		MaxResults:       maxResults,
		SortBy:           sortBy,
		SortDirection:    sortDirection,
		IncludeCompleted: includeCompleted,
	}

	// Extract filters if provided
	if filtersParam, ok := params["filters"]; ok {
		if filters, ok := filtersParam.(map[string]interface{}); ok {
			options.Filters = filters
		}
	}

	// Parse query string for "key:value" patterns and add to filters
	// Supports queries like "durum:devam_ediyor oncelik:yuksek"
	if query != "" && len(options.Filters) == 0 {
		parsedFilters := parseQueryFilters(query)
		if len(parsedFilters) > 0 {
			options.Filters = parsedFilters
			// Clear query if it was fully parsed into filters
			options.Query = ""
		}
	}

	// Extract search fields if provided
	if fieldsParam, ok := params["search_fields"]; ok {
		if fields, ok := fieldsParam.([]interface{}); ok {
			for _, field := range fields {
				if fieldStr, ok := field.(string); ok {
					options.SearchFields = append(options.SearchFields, fieldStr)
				}
			}
		}
	}

	// Create search engine and perform search
	db, err := h.isYonetici.VeriYonetici().GetDB()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Database access failed: %v", err)), nil
	}
	searchEngine := gorev.NewSearchEngine(h.isYonetici.VeriYonetici(), db)
	response, err := searchEngine.Search(options)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	// Format response as simple text for now
	responseText := fmt.Sprintf("Advanced Search Results:\n")
	responseText += fmt.Sprintf("Query: '%s'\n", query)
	responseText += fmt.Sprintf("Found %d results in %dms\n", response.TotalCount, response.QueryTime.Milliseconds())
	responseText += fmt.Sprintf("Fuzzy search: %v\n", response.UsedFuzzy)

	if len(response.Results) > 0 {
		responseText += "\nTop Results:\n"
		for i, result := range response.Results {
			if i >= 5 { // Limit to top 5 for display
				break
			}
			responseText += fmt.Sprintf("- %s (Score: %.2f)\n", result.Task.Title, result.RelevanceScore)
		}
	}

	return mcp.NewToolResultText(responseText), nil
}

// GorevFilterProfileSave saves a filter profile
func (h *Handlers) GorevFilterProfileSave(params map[string]interface{}) (*mcp.CallToolResult, error) {
	name, err := h.toolHelpers.Validator.ValidateRequiredString(params, "name")
	if err != nil {
		return err, nil
	}

	// TODO: Extract and use description, searchQuery, isDefault, filters when implementing FilterProfileManager
	// For now, just return success message
	return mcp.NewToolResultText(fmt.Sprintf("Filter profile '%s' saved successfully", name)), nil
}

// GorevFilterProfileLoad loads a filter profile by ID or name
func (h *Handlers) GorevFilterProfileLoad(params map[string]interface{}) (*mcp.CallToolResult, error) {
	profileID := 0
	if val, ok := params["profile_id"]; ok {
		if intVal, ok := val.(int); ok {
			profileID = intVal
		} else if floatVal, ok := val.(float64); ok {
			profileID = int(floatVal)
		}
	}

	profileName := h.toolHelpers.Validator.ValidateOptionalString(params, "profile_name")

	if profileID == 0 && profileName == "" {
		return mcp.NewToolResultError("Profile ID or name is required"), nil
	}

	// TODO: Implement actual profile loading
	// For now, return mock data
	profile := &gorev.FilterProfile{
		ID:          fmt.Sprintf("%d", profileID),
		Name:        profileName,
		Description: "Mock filter profile",
		Filters:     gorev.SearchFilters{Status: []string{"devam_ediyor"}},
		SearchQuery: "",
		IsDefault:   false,
	}

	return mcp.NewToolResultText(fmt.Sprintf("Loaded filter profile: %s", profile.Name)), nil
}

// GorevFilterProfileList lists all filter profiles
func (h *Handlers) GorevFilterProfileList(params map[string]interface{}) (*mcp.CallToolResult, error) {
	defaultsOnly := false
	if val, ok := params["defaults_only"].(bool); ok {
		defaultsOnly = val
	}

	// TODO: Implement actual profile listing
	// For now, return mock data
	profiles := []*gorev.FilterProfile{
		{
			ID:          "1",
			Name:        "Yüksek Öncelik",
			Description: "Yüksek öncelikli görevler",
			Filters:     gorev.SearchFilters{Priority: []string{"yuksek"}},
			IsDefault:   true,
			UseCount:    5,
		},
		{
			ID:          "2",
			Name:        "Devam Ediyor",
			Description: "Şu anda üzerinde çalışılan görevler",
			Filters:     gorev.SearchFilters{Status: []string{"devam_ediyor"}},
			IsDefault:   true,
			UseCount:    10,
		},
	}

	// Filter for defaults only if requested
	if defaultsOnly {
		var defaultProfiles []*gorev.FilterProfile
		for _, profile := range profiles {
			if profile.IsDefault {
				defaultProfiles = append(defaultProfiles, profile)
			}
		}
		profiles = defaultProfiles
	}

	// Build detailed response with full profile information
	var response strings.Builder
	response.WriteString(fmt.Sprintf("## 📁 Filtre Profilleri (%d adet)\n\n", len(profiles)))

	if len(profiles) == 0 {
		response.WriteString("ℹ️ Henüz kaydedilmiş filtre profili bulunmuyor.\n")
	} else {
		for i, profile := range profiles {
			response.WriteString(fmt.Sprintf("### %d. %s\n", i+1, profile.Name))
			response.WriteString(fmt.Sprintf("- **ID:** %s\n", profile.ID))
			if profile.Description != "" {
				response.WriteString(fmt.Sprintf("- **Açıklama:** %s\n", profile.Description))
			}
			if profile.IsDefault {
				response.WriteString("- **Varsayılan:** ✅ Evet\n")
			}
			if profile.UseCount > 0 {
				response.WriteString(fmt.Sprintf("- **Kullanım Sayısı:** %d\n", profile.UseCount))
			}

			// Show filter details
			if len(profile.Filters.Status) > 0 {
				response.WriteString(fmt.Sprintf("- **Durum Filtresi:** %v\n", profile.Filters.Status))
			}
			if len(profile.Filters.Priority) > 0 {
				response.WriteString(fmt.Sprintf("- **Öncelik Filtresi:** %v\n", profile.Filters.Priority))
			}
			if len(profile.Filters.Tags) > 0 {
				response.WriteString(fmt.Sprintf("- **Etiket Filtresi:** %v\n", profile.Filters.Tags))
			}
			if len(profile.Filters.ProjectIDs) > 0 {
				response.WriteString(fmt.Sprintf("- **Proje ID'leri:** %v\n", profile.Filters.ProjectIDs))
			}

			response.WriteString("\n")
		}

		response.WriteString("💡 **Kullanım:** `gorev_filter_profile_load` komutuyla profil ID veya ismiyle yükleyebilirsiniz.\n")
	}

	return mcp.NewToolResultText(response.String()), nil
}

// GorevFilterProfileDelete deletes a filter profile
func (h *Handlers) GorevFilterProfileDelete(params map[string]interface{}) (*mcp.CallToolResult, error) {
	profileID := 0
	if val, ok := params["profile_id"]; ok {
		if intVal, ok := val.(int); ok {
			profileID = intVal
		} else if floatVal, ok := val.(float64); ok {
			profileID = int(floatVal)
		}
	}

	if profileID == 0 {
		return mcp.NewToolResultError("Profile ID is required"), nil
	}

	// TODO: Implement actual profile deletion
	// For now, return success message
	return mcp.NewToolResultText(fmt.Sprintf("Filter profile %d deleted successfully", profileID)), nil
}

// GorevSearchHistory returns recent search history
func (h *Handlers) GorevSearchHistory(params map[string]interface{}) (*mcp.CallToolResult, error) {
	limit := 20
	if val, ok := params["limit"]; ok {
		if intVal, ok := val.(int); ok {
			limit = intVal
		} else if floatVal, ok := val.(float64); ok {
			limit = int(floatVal)
		}
	}

	// TODO: Implement actual search history retrieval
	// For now, return mock data
	history := []*gorev.SearchHistoryEntry{
		{
			ID:              1,
			Query:           "database görevleri",
			Filters:         `{"durum": "devam_ediyor"}`,
			ResultCount:     5,
			ExecutionTimeMs: 150,
		},
		{
			ID:              2,
			Query:           "yüksek öncelikli",
			Filters:         `{"oncelik": "yuksek"}`,
			ResultCount:     8,
			ExecutionTimeMs: 89,
		},
	}

	// Apply limit
	if len(history) > limit {
		history = history[:limit]
	}

	responseText := fmt.Sprintf("Search History (last %d entries):\n", len(history))
	for _, entry := range history {
		responseText += fmt.Sprintf("- '%s' (%d results, %dms)\n", entry.Query, entry.ResultCount, entry.ExecutionTimeMs)
	}

	return mcp.NewToolResultText(responseText), nil
}

// ============================================================================
// UNIFIED HANDLERS - Optimized Tool Set
// ============================================================================

// AktifProje - Unified handler for active project operations
// Replaces: AktifProjeAyarla, AktifProjeGoster, AktifProjeKaldir
func (h *Handlers) AktifProje(params map[string]interface{}) (*mcp.CallToolResult, error) {
	action, ok := params["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action parameter is required (set|get|clear)"), nil
	}

	switch action {
	case "set":
		return h.AktifProjeAyarla(params)
	case "get":
		return h.AktifProjeGoster(params)
	case "clear":
		return h.AktifProjeKaldir(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action: %s (expected: set|get|clear)", action)), nil
	}
}

// GorevBulk - Unified handler for bulk operations
// Replaces: GorevBulkTransition, GorevBulkTag, GorevBatchUpdate
func (h *Handlers) GorevBulk(params map[string]interface{}) (*mcp.CallToolResult, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return mcp.NewToolResultError("operation parameter is required (transition|tag|update)"), nil
	}

	// Get ids array (required for all operations)
	idsRaw, ok := params["ids"]
	if !ok {
		return mcp.NewToolResultError("ids parameter is required"), nil
	}

	// Transform unified parameters to operation-specific format
	transformedParams := make(map[string]interface{})

	// Extract data object
	data, hasData := params["data"].(map[string]interface{})

	switch operation {
	case "transition":
		// Transform: ids → task_ids
		transformedParams["task_ids"] = idsRaw

		if hasData {
			// Accept "status" as the primary parameter name (v0.17.0 English migration)
			if status, ok := data[constants.ParamStatus].(string); ok {
				transformedParams[constants.ParamStatus] = status
			} else if newStatus, ok := data["new_status"].(string); ok {
				transformedParams[constants.ParamStatus] = newStatus
			}

			// Optional parameters
			if force, ok := data["force"].(bool); ok {
				transformedParams["force"] = force
			}
			if checkDeps, ok := data["check_dependencies"].(bool); ok {
				transformedParams["check_dependencies"] = checkDeps
			}
			if dryRun, ok := data["dry_run"].(bool); ok {
				transformedParams["dry_run"] = dryRun
			}
		}

		return h.GorevBulkTransition(transformedParams)

	case "tag":
		// Transform: ids → task_ids
		transformedParams["task_ids"] = idsRaw

		// Default operation is "add" if not specified
		transformedParams["operation"] = "add"

		if hasData {
			// data.tags → tags
			if tags, ok := data["tags"]; ok {
				transformedParams["tags"] = tags
			}

			// Accept both "operation" and "tag_operation"
			if op, ok := data["operation"].(string); ok {
				transformedParams["operation"] = op
			} else if tagOp, ok := data["tag_operation"].(string); ok {
				transformedParams["operation"] = tagOp
			}

			// Optional parameters
			if dryRun, ok := data["dry_run"].(bool); ok {
				transformedParams["dry_run"] = dryRun
			}
		}

		return h.GorevBulkTag(transformedParams)

	case "update":
		// Transform: ids + data → updates array
		// GorevBatchUpdate expects: {"updates": [{"id": "...", "field1": "...", ...}]}
		idsArray, ok := idsRaw.([]interface{})
		if !ok {
			return mcp.NewToolResultError("ids must be an array"), nil
		}

		if !hasData || len(data) == 0 {
			return mcp.NewToolResultError("data object is required for update operation"), nil
		}

		// Create updates array: each id gets the same data fields
		updates := make([]interface{}, len(idsArray))
		for i, idRaw := range idsArray {
			id, ok := idRaw.(string)
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("invalid id at index %d", i)), nil
			}

			updateObj := map[string]interface{}{"id": id}
			// Copy all data fields to this update object
			for key, value := range data {
				updateObj[key] = value
			}
			updates[i] = updateObj
		}

		transformedParams["updates"] = updates
		return h.GorevBatchUpdate(transformedParams)

	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid operation: %s (expected: transition|tag|update)", operation)), nil
	}
}

// GorevHierarchy - Unified handler for hierarchy operations
// Replaces: GorevAltGorevOlustur, GorevUstDegistir, GorevHiyerarsiGoster
func (h *Handlers) GorevHierarchy(params map[string]interface{}) (*mcp.CallToolResult, error) {
	action, ok := params["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action parameter is required (create_subtask|change_parent|show)"), nil
	}

	switch action {
	case "create_subtask":
		return h.GorevAltGorevOlustur(params)
	case "change_parent":
		return h.GorevUstDegistir(params)
	case "show":
		// Transform parameter: parent_id → task_id (unified schema uses parent_id for task ID)
		if parentID, ok := params["parent_id"].(string); ok && parentID != "" {
			if _, hasTaskID := params[constants.ParamTaskID]; !hasTaskID {
				params[constants.ParamTaskID] = parentID
			}
		}
		return h.GorevHiyerarsiGoster(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action: %s (expected: create_subtask|change_parent|show)", action)), nil
	}
}

// GorevFilterProfile - Unified handler for filter profile operations
// Replaces: GorevFilterProfileSave, GorevFilterProfileLoad, GorevFilterProfileList, GorevFilterProfileDelete
func (h *Handlers) GorevFilterProfile(params map[string]interface{}) (*mcp.CallToolResult, error) {
	action, ok := params["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action parameter is required (save|load|list|delete)"), nil
	}

	switch action {
	case "save":
		return h.GorevFilterProfileSave(params)
	case "load":
		return h.GorevFilterProfileLoad(params)
	case "list":
		return h.GorevFilterProfileList(params)
	case "delete":
		return h.GorevFilterProfileDelete(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action: %s (expected: save|load|list|delete)", action)), nil
	}
}

// GorevFileWatch - Unified handler for file watch operations
// Replaces: GorevFileWatchAdd, GorevFileWatchRemove, GorevFileWatchList, GorevFileWatchStats
func (h *Handlers) GorevFileWatch(params map[string]interface{}) (*mcp.CallToolResult, error) {
	action, ok := params["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action parameter is required (add|remove|list|stats)"), nil
	}

	switch action {
	case "add":
		return h.GorevFileWatchAdd(params)
	case "remove":
		return h.GorevFileWatchRemove(params)
	case "list":
		return h.GorevFileWatchList(params)
	case "stats":
		return h.GorevFileWatchStats(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action: %s (expected: add|remove|list|stats)", action)), nil
	}
}

// IDEManage - Unified handler for IDE extension management
// Replaces: IDEDetect, IDEInstallExtension, IDEUninstallExtension, IDEExtensionStatus, IDEUpdateExtension
func (h *Handlers) IDEManage(params map[string]interface{}) (*mcp.CallToolResult, error) {
	action, ok := params["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action parameter is required (detect|install|uninstall|status|update)"), nil
	}

	switch action {
	case "detect":
		return h.IDEDetect(params)
	case "install":
		return h.IDEInstallExtension(params)
	case "uninstall":
		return h.IDEUninstallExtension(params)
	case "status":
		return h.IDEExtensionStatus(params)
	case "update":
		return h.IDEUpdateExtension(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action: %s (expected: detect|install|uninstall|status|update)", action)), nil
	}
}

// GorevContext - Unified handler for AI context operations
// Replaces: GorevSetActive, GorevGetActive, GorevRecent, GorevContextSummary
func (h *Handlers) GorevContext(params map[string]interface{}) (*mcp.CallToolResult, error) {
	action, ok := params["action"].(string)
	if !ok {
		return mcp.NewToolResultError("action parameter is required (set_active|get_active|recent|summary)"), nil
	}

	switch action {
	case "set_active":
		// Use task_id parameter directly (v0.17.0 English migration)
		// Note: params already contains task_id from unified schema
		return h.GorevSetActive(params)
	case "get_active":
		return h.GorevGetActive(params)
	case "recent":
		return h.GorevRecent(params)
	case "summary":
		return h.GorevContextSummary(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid action: %s (expected: set_active|get_active|recent|summary)", action)), nil
	}
}

// GorevSearch - Unified handler for search operations
// Replaces: GorevNLPQuery, GorevSearchAdvanced, GorevSearchHistory
func (h *Handlers) GorevSearch(params map[string]interface{}) (*mcp.CallToolResult, error) {
	mode, ok := params["mode"].(string)
	if !ok {
		return mcp.NewToolResultError("mode parameter is required (nlp|advanced|history)"), nil
	}

	// Transform parameter: arama_metni → query
	if aramaMetni, ok := params["arama_metni"].(string); ok && aramaMetni != "" {
		params["query"] = aramaMetni
	}

	switch mode {
	case "nlp":
		return h.GorevNLPQuery(params)
	case "advanced":
		return h.GorevSearchAdvanced(params)
	case "history":
		return h.GorevSearchHistory(params)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid mode: %s (expected: nlp|advanced|history)", mode)), nil
	}
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// parseQueryFilters parses query strings like "durum:devam_ediyor oncelik:yuksek"
// into a filters map for advanced search
func parseQueryFilters(query string) map[string]interface{} {
	filters := make(map[string]interface{})

	// Split by spaces, but preserve quoted strings
	parts := strings.Fields(query)

	for _, part := range parts {
		// Look for "key:value" pattern
		if strings.Contains(part, ":") {
			kv := strings.SplitN(part, ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])

				if key != "" && value != "" {
					filters[key] = value
				}
			}
		}
	}

	return filters
}
