package gorev

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

type IsYonetici struct {
	veriYonetici VeriYoneticiInterface
	workspaceID  string // Workspace ID for centralized mode filtering
}

func YeniIsYonetici(veriYonetici VeriYoneticiInterface) *IsYonetici {
	return &IsYonetici{
		veriYonetici: veriYonetici,
		workspaceID:  "", // Empty means no workspace filtering (local mode)
	}
}

// YeniIsYoneticiWithWorkspaceID creates a new IsYonetici with workspace ID filtering
// Used in centralized mode where all data is in a single database
func YeniIsYoneticiWithWorkspaceID(veriYonetici VeriYoneticiInterface, workspaceID string) *IsYonetici {
	return &IsYonetici{
		veriYonetici: veriYonetici,
		workspaceID:  workspaceID,
	}
}

// GetWorkspaceID returns the workspace ID for this manager
func (iy *IsYonetici) GetWorkspaceID() string {
	return iy.workspaceID
}

// addWorkspaceFilter adds workspace_id to filters if in centralized mode
func (iy *IsYonetici) addWorkspaceFilter(filters map[string]interface{}) map[string]interface{} {
	if iy.workspaceID == "" {
		return filters
	}
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["workspace_id"] = iy.workspaceID
	return filters
}

// VeriYonetici returns the data manager interface
func (iy *IsYonetici) VeriYonetici() VeriYoneticiInterface {
	return iy.veriYonetici
}

func (iy *IsYonetici) GorevOlustur(ctx context.Context, baslik, aciklama, oncelik, projeID, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	var sonTarih *time.Time
	if sonTarihStr != "" {
		t, err := time.Parse("2006-01-02", sonTarihStr)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.invalidDateFormat", map[string]interface{}{"Error": err}))
		}
		sonTarih = &t
	}

	gorev := &Gorev{
		ID:          uuid.New().String(),
		Title:       baslik,
		Description: aciklama,
		Priority:    oncelik,
		Status:      constants.TaskStatusPending,
		ProjeID:     projeID,
		WorkspaceID: iy.workspaceID, // Set workspace ID for centralized mode
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     sonTarih,
	}

	if err := iy.veriYonetici.GorevKaydet(ctx, gorev); err != nil {
		return nil, fmt.Errorf(i18n.TSaveFailed(i18n.FromContext(ctx), "task", err))
	}

	if len(etiketIsimleri) > 0 {
		etiketler, err := iy.veriYonetici.EtiketleriGetirVeyaOlustur(ctx, etiketIsimleri)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsProcessFailed", map[string]interface{}{"Error": err}))
		}
		if err := iy.veriYonetici.GorevEtiketleriniAyarla(ctx, gorev.ID, etiketler); err != nil {
			return nil, fmt.Errorf(i18n.TSetFailed(i18n.FromContext(ctx), "task_tags", err))
		}
		gorev.Tags = etiketler
	}

	return gorev, nil
}

func (iy *IsYonetici) GorevListele(ctx context.Context, filters map[string]interface{}) ([]*Gorev, error) {
	// Add workspace filter for centralized mode
	filters = iy.addWorkspaceFilter(filters)

	gorevler, err := iy.veriYonetici.GorevListele(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Performans optimizasyonu: Tüm görevlerin ID'lerini topla
	if len(gorevler) == 0 {
		return gorevler, nil
	}

	gorevIDs := make([]string, len(gorevler))
	for i, gorev := range gorevler {
		gorevIDs[i] = gorev.ID
	}

	// Tek sorguda tüm görevlerin bağımlılık sayılarını hesapla (N+1 sorgu problemi çözüldü)
	bagimliSayilari, err := iy.veriYonetici.BulkBagimlilikSayilariGetir(gorevIDs)
	if err != nil {
		// Hata durumunda bile devam et, sadece bağımlılık sayıları 0 olarak kalır
		bagimliSayilari = make(map[string]int)
	}

	// Tek sorguda tüm görevlerin tamamlanmamış bağımlılık sayılarını hesapla
	tamamlanmamisSayilari, err := iy.veriYonetici.BulkTamamlanmamiaBagimlilikSayilariGetir(gorevIDs)
	if err != nil {
		// Hata durumunda bile devam et, sadece tamamlanmamış bağımlılık sayıları 0 olarak kalır
		tamamlanmamisSayilari = make(map[string]int)
	}

	// Bu göreve bağımlı olanları hesaplamak için bulk query kullan
	buGoreveBagimliSayilari, err := iy.veriYonetici.BulkBuGoreveBagimliSayilariGetir(gorevIDs)
	if err != nil {
		// Hata durumunda bile devam et, sadece bağımlı sayıları 0 olarak kalır
		buGoreveBagimliSayilari = make(map[string]int)
	}

	// Her görev için hesaplanan değerleri ata
	for _, gorev := range gorevler {
		if count, exists := bagimliSayilari[gorev.ID]; exists {
			gorev.DependencyCount = count
		}

		if count, exists := tamamlanmamisSayilari[gorev.ID]; exists {
			gorev.UncompletedDependencyCount = count
		}

		if count, exists := buGoreveBagimliSayilari[gorev.ID]; exists {
			gorev.DependentOnThisCount = count
		}

		// Alt görevleri getir (parent_id null olan görevler için)
		if gorev.ParentID == "" {
			altGorevler, err := iy.veriYonetici.AltGorevleriGetir(ctx, gorev.ID)
			if err == nil && len(altGorevler) > 0 {
				gorev.Subtasks = altGorevler
			}
		}
	}

	return gorevler, nil
}

func (iy *IsYonetici) GorevDurumGuncelle(ctx context.Context, id, durum string) error {
	// Validate status values
	validStatuses := constants.GetValidTaskStatuses()
	isValidStatus := false
	for _, validStatus := range validStatuses {
		if durum == validStatus {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return fmt.Errorf(i18n.T("error.invalidStatus", map[string]interface{}{"Status": durum, "ValidStatuses": validStatuses}))
	}

	gorev, err := iy.veriYonetici.GorevGetir(ctx, id)
	if err != nil {
		return fmt.Errorf(i18n.TEntityNotFound(i18n.FromContext(ctx), "task", err))
	}

	// Eğer görev "devam_ediyor" durumuna geçiyorsa, bağımlılıkları kontrol et
	if durum == constants.TaskStatusInProgress && gorev.Status == constants.TaskStatusPending {
		bagimli, tamamlanmamislar, err := iy.GorevBagimliMi(ctx, id)
		if err != nil {
			return fmt.Errorf(i18n.TCheckFailed(i18n.FromContext(ctx), "dependency", err))
		}

		if !bagimli {
			return fmt.Errorf(i18n.T("error.taskCannotStartDependencies", map[string]interface{}{"Dependencies": tamamlanmamislar}))
		}
	}

	// Eğer görev "tamamlandi" durumuna geçiyorsa, tüm alt görevlerin tamamlandığını kontrol et
	if durum == constants.TaskStatusCompleted && gorev.Status != constants.TaskStatusCompleted {
		altGorevler, err := iy.veriYonetici.AltGorevleriGetir(ctx, id)
		if err != nil {
			return fmt.Errorf(i18n.T("error.subtasksCheckFailed", map[string]interface{}{"Error": err}))
		}

		for _, altGorev := range altGorevler {
			if altGorev.Status != constants.TaskStatusCompleted {
				return fmt.Errorf(i18n.T("error.taskCannotCompleteSubtasks"))
			}
		}
	}

	gorev.Status = durum
	gorev.UpdatedAt = time.Now()

	return iy.veriYonetici.GorevGuncelle(ctx, gorev.ID, map[string]interface{}{
		"status":     durum,
		"updated_at": time.Now(),
	})
}

func (iy *IsYonetici) ProjeOlustur(ctx context.Context, isim, tanim string) (*Proje, error) {
	proje := &Proje{
		ID:          uuid.New().String(),
		Name:        isim,
		Definition:  tanim,
		WorkspaceID: iy.workspaceID, // Set workspace ID for centralized mode
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := iy.veriYonetici.ProjeKaydet(ctx, proje); err != nil {
		return nil, fmt.Errorf(i18n.TSaveFailed(i18n.FromContext(ctx), "project", err))
	}

	return proje, nil
}

func (iy *IsYonetici) GorevGetir(ctx context.Context, id string) (*Gorev, error) {
	gorev, err := iy.veriYonetici.GorevGetir(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(i18n.TEntityNotFound(i18n.FromContext(ctx), "task", err))
	}
	return gorev, nil
}

func (iy *IsYonetici) ProjeGetir(ctx context.Context, id string) (*Proje, error) {
	proje, err := iy.veriYonetici.ProjeGetir(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.projectNotFound", map[string]interface{}{"Error": err}))
	}
	return proje, nil
}

func (iy *IsYonetici) GorevDuzenle(ctx context.Context, id, baslik, aciklama, oncelik, projeID, sonTarihStr string, baslikVar, aciklamaVar, oncelikVar, projeVar, sonTarihVar bool) error {
	// Önce mevcut görevi al
	gorev, err := iy.veriYonetici.GorevGetir(ctx, id)
	if err != nil {
		return fmt.Errorf(i18n.TEntityNotFound(i18n.FromContext(ctx), "task", err))
	}

	// Sadece belirtilen alanları güncelle
	if baslikVar && baslik != "" {
		gorev.Title = baslik
	}
	if aciklamaVar {
		gorev.Description = aciklama
	}
	if oncelikVar && oncelik != "" {
		gorev.Priority = oncelik
	}
	if projeVar {
		// Proje değiştiriliyorsa, tüm alt görevleri de taşı
		if gorev.ProjeID != projeID {
			altGorevler, err := iy.veriYonetici.TumAltGorevleriGetir(ctx, id)
			if err != nil {
				return fmt.Errorf(i18n.T("error.subtasksFetchFailed", map[string]interface{}{"Error": err}))
			}

			// Tüm alt görevlerin projesini güncelle
			for _, altGorev := range altGorevler {
				if err := iy.veriYonetici.GorevGuncelle(ctx, altGorev.ID, map[string]interface{}{
					"project_id": projeID,
					"updated_at": time.Now(),
				}); err != nil {
					return fmt.Errorf(i18n.T("error.subtaskUpdateFailed", map[string]interface{}{"Error": err}))
				}
			}
		}
		gorev.ProjeID = projeID
	}
	if sonTarihVar {
		if sonTarihStr == "" {
			gorev.DueDate = nil
		} else {
			t, err := time.Parse("2006-01-02", sonTarihStr)
			if err != nil {
				return fmt.Errorf(i18n.T("error.invalidDateFormat", map[string]interface{}{"Error": err}))
			}
			gorev.DueDate = &t
		}
	}

	updateParams := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if baslikVar && baslik != "" {
		updateParams["title"] = baslik
	}
	if aciklamaVar {
		updateParams["description"] = aciklama
	}
	if oncelikVar && oncelik != "" {
		updateParams["priority"] = oncelik
	}
	if projeVar {
		updateParams["project_id"] = projeID
	}
	if sonTarihVar {
		if sonTarihStr == "" {
			updateParams["due_date"] = nil
		} else {
			t, _ := time.Parse("2006-01-02", sonTarihStr)
			updateParams["due_date"] = t
		}
	}

	return iy.veriYonetici.GorevGuncelle(ctx, gorev.ID, updateParams)
}

func (iy *IsYonetici) GorevSil(ctx context.Context, id string) error {
	// Önce görevin var olduğunu kontrol et
	_, err := iy.veriYonetici.GorevGetir(ctx, id)
	if err != nil {
		return fmt.Errorf(i18n.TEntityNotFound(i18n.FromContext(ctx), "task", err))
	}

	// Alt görevleri kontrol et
	altGorevler, err := iy.veriYonetici.AltGorevleriGetir(ctx, id)
	if err != nil {
		return fmt.Errorf(i18n.T("error.subtasksCheckFailed", map[string]interface{}{"Error": err}))
	}

	if len(altGorevler) > 0 {
		return fmt.Errorf(i18n.T("error.taskHasSubtasksCannotDelete", map[string]interface{}{"Count": len(altGorevler)}))
	}

	return iy.veriYonetici.GorevSil(ctx, id)
}

func (iy *IsYonetici) ProjeListele(ctx context.Context) ([]*Proje, error) {
	return iy.veriYonetici.ProjeleriGetir(ctx)
}

func (iy *IsYonetici) ProjeGorevleri(ctx context.Context, projeID string) ([]*Gorev, error) {
	gorevler, err := iy.veriYonetici.ProjeGorevleriGetir(ctx, projeID)
	if err != nil {
		return nil, err
	}

	// Her görev için bağımlılık sayılarını hesapla
	for _, gorev := range gorevler {
		// Bu görevin bağımlılıklarını al (bu görev başka görevlere bağımlı)
		baglantilar, err := iy.veriYonetici.BaglantilariGetir(ctx, gorev.ID)
		if err == nil && len(baglantilar) > 0 {
			gorev.DependencyCount = len(baglantilar)

			// Tamamlanmamış bağımlılıkları say
			tamamlanmamisSayisi := 0
			for _, baglanti := range baglantilar {
				hedefGorev, err := iy.veriYonetici.GorevGetir(ctx, baglanti.TargetID)
				if err == nil && hedefGorev.Status != constants.TaskStatusCompleted {
					tamamlanmamisSayisi++
				}
			}
			gorev.UncompletedDependencyCount = tamamlanmamisSayisi
		}

		// Bu göreve bağımlı olan görevleri bul
		buGoreveBagimliSayisi := 0
		tumGorevler, _ := iy.veriYonetici.GorevleriGetir(ctx, "", "", "")
		for _, digerGorev := range tumGorevler {
			digerBaglantilar, err := iy.veriYonetici.BaglantilariGetir(ctx, digerGorev.ID)
			if err == nil {
				for _, baglanti := range digerBaglantilar {
					if baglanti.TargetID == gorev.ID {
						buGoreveBagimliSayisi++
						break
					}
				}
			}
		}
		gorev.DependentOnThisCount = buGoreveBagimliSayisi
	}

	return gorevler, nil
}

func (iy *IsYonetici) ProjeGorevSayisi(ctx context.Context, projeID string) (int, error) {
	gorevler, err := iy.veriYonetici.ProjeGorevleriGetir(ctx, projeID)
	if err != nil {
		return 0, err
	}
	return len(gorevler), nil
}

func (iy *IsYonetici) AktifProjeAyarla(ctx context.Context, projeID string) error {
	return iy.veriYonetici.AktifProjeAyarla(ctx, projeID)
}

func (iy *IsYonetici) AktifProjeGetir(ctx context.Context) (*Proje, error) {
	projeID, err := iy.veriYonetici.AktifProjeGetir(ctx)
	if err != nil {
		return nil, err
	}
	if projeID == "" {
		return nil, nil // Aktif proje yok
	}
	return iy.veriYonetici.ProjeGetir(ctx, projeID)
}

func (iy *IsYonetici) AktifProjeKaldir(ctx context.Context) error {
	return iy.veriYonetici.AktifProjeKaldir(ctx)
}

func (iy *IsYonetici) OzetAl(ctx context.Context) (*Ozet, error) {
	gorevler, err := iy.veriYonetici.GorevleriGetir(ctx, "", "", "")
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.tasksFetchFailed", map[string]interface{}{"Error": err}))
	}

	projeler, err := iy.veriYonetici.ProjeleriGetir(ctx)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.projectListFailed", map[string]interface{}{"Error": err}))
	}

	ozet := &Ozet{
		TotalProjects: len(projeler),
		TotalTasks:    len(gorevler),
	}

	for _, gorev := range gorevler {
		switch gorev.Status {
		case constants.TaskStatusPending:
			ozet.PendingTasks++
		case constants.TaskStatusInProgress:
			ozet.InProgressTasks++
		case constants.TaskStatusCompleted:
			ozet.CompletedTasks++
		}

		switch gorev.Priority {
		case constants.PriorityHigh:
			ozet.HighPriorityTasks++
		case constants.PriorityMedium:
			ozet.MediumPriorityTasks++
		case constants.PriorityLow:
			ozet.LowPriorityTasks++
		}
	}

	return ozet, nil
}

func (iy *IsYonetici) GorevBagimlilikEkle(ctx context.Context, kaynakID, hedefID, baglantiTipi string) (*Baglanti, error) {
	// Görevlerin var olup olmadığını kontrol et
	_, err := iy.veriYonetici.GorevGetir(ctx, kaynakID)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.sourceTaskNotFound", map[string]interface{}{"Error": err}))
	}
	_, err = iy.veriYonetici.GorevGetir(ctx, hedefID)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.targetTaskNotFound", map[string]interface{}{"Error": err}))
	}

	baglanti := &Baglanti{
		ID:             uuid.New().String(),
		SourceID:       kaynakID,
		TargetID:       hedefID,
		ConnectionType: baglantiTipi,
	}

	if err := iy.veriYonetici.BaglantiEkle(ctx, baglanti); err != nil {
		return nil, fmt.Errorf(i18n.TAddFailed(i18n.FromContext(ctx), "link", err))
	}

	return baglanti, nil
}

func (iy *IsYonetici) GorevBaglantilariGetir(ctx context.Context, gorevID string) ([]*Baglanti, error) {
	return iy.veriYonetici.BaglantilariGetir(ctx, gorevID)
}

// GorevBagimliMi görevi başlatmak için tüm bağımlılıkların tamamlanıp tamamlanmadığını kontrol eder
func (iy *IsYonetici) GorevBagimliMi(ctx context.Context, gorevID string) (bool, []string, error) {
	baglantilar, err := iy.veriYonetici.BaglantilariGetir(ctx, gorevID)
	if err != nil {
		return false, nil, fmt.Errorf(i18n.T("error.dependencyFetchFailed", map[string]interface{}{"Error": err}))
	}

	var tamamlanmamisBagimliliklar []string

	for _, baglanti := range baglantilar {
		// Bu görev hedef konumundaysa ve bağlantı tipi "onceki" ise
		if baglanti.TargetID == gorevID && baglanti.ConnectionType == "onceki" {
			// Kaynak görevin durumunu kontrol et
			kaynakGorev, err := iy.veriYonetici.GorevGetir(ctx, baglanti.SourceID)
			if err != nil {
				return false, nil, fmt.Errorf(i18n.T("error.dependentTaskNotFound", map[string]interface{}{"Error": err}))
			}

			// Eğer bağımlı görev tamamlanmamışsa
			if kaynakGorev.Status != constants.TaskStatusCompleted {
				tamamlanmamisBagimliliklar = append(tamamlanmamisBagimliliklar, kaynakGorev.Title)
			}
		}
	}

	return len(tamamlanmamisBagimliliklar) == 0, tamamlanmamisBagimliliklar, nil
}

// TemplateListele kullanılabilir template'leri listeler
func (iy *IsYonetici) TemplateListele(ctx context.Context, kategori string) ([]*GorevTemplate, error) {
	return iy.veriYonetici.TemplateListele(ctx, kategori)
}

// TemplatedenGorevOlustur template kullanarak görev oluşturur
func (iy *IsYonetici) TemplatedenGorevOlustur(ctx context.Context, templateID string, degerler map[string]string) (*Gorev, error) {
	// Inject workspace_id for centralized mode
	if iy.workspaceID != "" {
		if degerler == nil {
			degerler = make(map[string]string)
		}
		degerler["_workspace_id"] = iy.workspaceID
	}
	return iy.veriYonetici.TemplatedenGorevOlustur(ctx, templateID, degerler)
}

// AltGorevOlustur mevcut bir görevin altına yeni görev oluşturur
func (iy *IsYonetici) AltGorevOlustur(ctx context.Context, parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	// Parent görevi kontrol et
	parent, err := iy.veriYonetici.GorevGetir(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("parentTaskNotFound"))
	}

	var sonTarih *time.Time
	if sonTarihStr != "" {
		t, err := time.Parse("2006-01-02", sonTarihStr)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.invalidDateFormat", map[string]interface{}{"Error": err}))
		}
		sonTarih = &t
	}

	gorev := &Gorev{
		ID:          uuid.New().String(),
		Title:       baslik,
		Description: aciklama,
		Priority:    oncelik,
		Status:      constants.TaskStatusPending,
		ProjeID:     parent.ProjeID, // Alt görev aynı projede olmalı
		ParentID:    parentID,
		WorkspaceID: iy.workspaceID, // Set workspace ID for centralized mode
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     sonTarih,
	}

	if err := iy.veriYonetici.GorevKaydet(ctx, gorev); err != nil {
		return nil, fmt.Errorf(i18n.TSaveFailed(i18n.FromContext(ctx), "subtask", err))
	}

	if len(etiketIsimleri) > 0 {
		etiketler, err := iy.veriYonetici.EtiketleriGetirVeyaOlustur(ctx, etiketIsimleri)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsProcessFailed", map[string]interface{}{"Error": err}))
		}
		if err := iy.veriYonetici.GorevEtiketleriniAyarla(ctx, gorev.ID, etiketler); err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsSetFailed", map[string]interface{}{"Error": err}))
		}
		gorev.Tags = etiketler
	}

	return gorev, nil
}

// GorevUstDegistir bir görevin üst görevini değiştirir
func (iy *IsYonetici) GorevUstDegistir(ctx context.Context, gorevID, yeniParentID string) error {
	// Görevi kontrol et
	gorev, err := iy.veriYonetici.GorevGetir(ctx, gorevID)
	if err != nil {
		return fmt.Errorf(i18n.TEntityNotFound(i18n.FromContext(ctx), "task", err))
	}

	// Yeni parent varsa kontrol et
	if yeniParentID != "" {
		parent, err := iy.veriYonetici.GorevGetir(ctx, yeniParentID)
		if err != nil {
			return fmt.Errorf(i18n.T("parentTaskNotFound"))
		}

		// Circular dependency kontrolü
		circular, err := iy.veriYonetici.DaireBagimliligiKontrolEt(ctx, gorevID, yeniParentID)
		if err != nil {
			return fmt.Errorf(i18n.T("error.circularDependencyCheckFailed", map[string]interface{}{"Error": err}))
		}
		if circular {
			return fmt.Errorf(i18n.T("circularDependency"))
		}

		// Alt görev ve üst görev aynı projede olmalı
		if gorev.ProjeID != parent.ProjeID {
			return fmt.Errorf(i18n.T("error.subtaskProjectMismatch"))
		}
	}

	return iy.veriYonetici.ParentIDGuncelle(ctx, gorevID, yeniParentID)
}

// GorevHiyerarsiGetir bir görevin tam hiyerarşi bilgilerini getirir
func (iy *IsYonetici) GorevHiyerarsiGetir(ctx context.Context, gorevID string) (*GorevHiyerarsi, error) {
	return iy.veriYonetici.GorevHiyerarsiGetir(ctx, gorevID)
}

// AltGorevleriGetir bir görevin doğrudan alt görevlerini getirir
func (iy *IsYonetici) AltGorevleriGetir(ctx context.Context, parentID string) ([]*Gorev, error) {
	return iy.veriYonetici.AltGorevleriGetir(ctx, parentID)
}

// TumAltGorevleriGetir bir görevin tüm alt görev ağacını getirir
func (iy *IsYonetici) TumAltGorevleriGetir(ctx context.Context, parentID string) ([]*Gorev, error) {
	return iy.veriYonetici.TumAltGorevleriGetir(ctx, parentID)
}
