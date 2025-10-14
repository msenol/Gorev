package gorev

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

type IsYonetici struct {
	veriYonetici VeriYoneticiInterface
}

func YeniIsYonetici(veriYonetici VeriYoneticiInterface) *IsYonetici {
	return &IsYonetici{
		veriYonetici: veriYonetici,
	}
}

// VeriYonetici returns the data manager interface
func (iy *IsYonetici) VeriYonetici() VeriYoneticiInterface {
	return iy.veriYonetici
}

func (iy *IsYonetici) GorevOlustur(baslik, aciklama, oncelik, projeID, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
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
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     sonTarih,
	}

	if err := iy.veriYonetici.GorevKaydet(gorev); err != nil {
		return nil, fmt.Errorf(i18n.TSaveFailed("tr", "task", err))
	}

	if len(etiketIsimleri) > 0 {
		etiketler, err := iy.veriYonetici.EtiketleriGetirVeyaOlustur(etiketIsimleri)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsProcessFailed", map[string]interface{}{"Error": err}))
		}
		if err := iy.veriYonetici.GorevEtiketleriniAyarla(gorev.ID, etiketler); err != nil {
			return nil, fmt.Errorf(i18n.TSetFailed("tr", "task_tags", err))
		}
		gorev.Tags = etiketler
	}

	return gorev, nil
}

func (iy *IsYonetici) GorevListele(filters map[string]interface{}) ([]*Gorev, error) {
	gorevler, err := iy.veriYonetici.GorevListele(filters)
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
			altGorevler, err := iy.veriYonetici.AltGorevleriGetir(gorev.ID)
			if err == nil && len(altGorevler) > 0 {
				gorev.Subtasks = altGorevler
			}
		}
	}

	return gorevler, nil
}

func (iy *IsYonetici) GorevDurumGuncelle(id, durum string) error {
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

	gorev, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return fmt.Errorf(i18n.TEntityNotFound("tr", "task", err))
	}

	// Eğer görev "devam_ediyor" durumuna geçiyorsa, bağımlılıkları kontrol et
	if durum == constants.TaskStatusInProgress && gorev.Status == constants.TaskStatusPending {
		bagimli, tamamlanmamislar, err := iy.GorevBagimliMi(id)
		if err != nil {
			return fmt.Errorf(i18n.TCheckFailed("tr", "dependency", err))
		}

		if !bagimli {
			return fmt.Errorf(i18n.T("error.taskCannotStartDependencies", map[string]interface{}{"Dependencies": tamamlanmamislar}))
		}
	}

	// Eğer görev "tamamlandi" durumuna geçiyorsa, tüm alt görevlerin tamamlandığını kontrol et
	if durum == constants.TaskStatusCompleted && gorev.Status != constants.TaskStatusCompleted {
		altGorevler, err := iy.veriYonetici.AltGorevleriGetir(id)
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

	return iy.veriYonetici.GorevGuncelle(gorev.ID, map[string]interface{}{
		"status":     durum,
		"updated_at": time.Now(),
	})
}

func (iy *IsYonetici) ProjeOlustur(isim, tanim string) (*Proje, error) {
	proje := &Proje{
		ID:         uuid.New().String(),
		Name:       isim,
		Definition: tanim,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := iy.veriYonetici.ProjeKaydet(proje); err != nil {
		return nil, fmt.Errorf(i18n.TSaveFailed("tr", "project", err))
	}

	return proje, nil
}

func (iy *IsYonetici) GorevGetir(id string) (*Gorev, error) {
	gorev, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return nil, fmt.Errorf(i18n.TEntityNotFound("tr", "task", err))
	}
	return gorev, nil
}

func (iy *IsYonetici) ProjeGetir(id string) (*Proje, error) {
	proje, err := iy.veriYonetici.ProjeGetir(id)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.projectNotFound", map[string]interface{}{"Error": err}))
	}
	return proje, nil
}

func (iy *IsYonetici) GorevDuzenle(id, baslik, aciklama, oncelik, projeID, sonTarihStr string, baslikVar, aciklamaVar, oncelikVar, projeVar, sonTarihVar bool) error {
	// Önce mevcut görevi al
	gorev, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return fmt.Errorf(i18n.TEntityNotFound("tr", "task", err))
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
			altGorevler, err := iy.veriYonetici.TumAltGorevleriGetir(id)
			if err != nil {
				return fmt.Errorf(i18n.T("error.subtasksFetchFailed", map[string]interface{}{"Error": err}))
			}

			// Tüm alt görevlerin projesini güncelle
			for _, altGorev := range altGorevler {
				if err := iy.veriYonetici.GorevGuncelle(altGorev.ID, map[string]interface{}{
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

	return iy.veriYonetici.GorevGuncelle(gorev.ID, updateParams)
}

func (iy *IsYonetici) GorevSil(id string) error {
	// Önce görevin var olduğunu kontrol et
	_, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return fmt.Errorf(i18n.TEntityNotFound("tr", "task", err))
	}

	// Alt görevleri kontrol et
	altGorevler, err := iy.veriYonetici.AltGorevleriGetir(id)
	if err != nil {
		return fmt.Errorf(i18n.T("error.subtasksCheckFailed", map[string]interface{}{"Error": err}))
	}

	if len(altGorevler) > 0 {
		return fmt.Errorf(i18n.T("error.taskHasSubtasksCannotDelete", map[string]interface{}{"Count": len(altGorevler)}))
	}

	return iy.veriYonetici.GorevSil(id)
}

func (iy *IsYonetici) ProjeListele() ([]*Proje, error) {
	return iy.veriYonetici.ProjeleriGetir()
}

func (iy *IsYonetici) ProjeGorevleri(projeID string) ([]*Gorev, error) {
	gorevler, err := iy.veriYonetici.ProjeGorevleriGetir(projeID)
	if err != nil {
		return nil, err
	}

	// Her görev için bağımlılık sayılarını hesapla
	for _, gorev := range gorevler {
		// Bu görevin bağımlılıklarını al (bu görev başka görevlere bağımlı)
		baglantilar, err := iy.veriYonetici.BaglantilariGetir(gorev.ID)
		if err == nil && len(baglantilar) > 0 {
			gorev.DependencyCount = len(baglantilar)

			// Tamamlanmamış bağımlılıkları say
			tamamlanmamisSayisi := 0
			for _, baglanti := range baglantilar {
				hedefGorev, err := iy.veriYonetici.GorevGetir(baglanti.TargetID)
				if err == nil && hedefGorev.Status != constants.TaskStatusCompleted {
					tamamlanmamisSayisi++
				}
			}
			gorev.UncompletedDependencyCount = tamamlanmamisSayisi
		}

		// Bu göreve bağımlı olan görevleri bul
		buGoreveBagimliSayisi := 0
		tumGorevler, _ := iy.veriYonetici.GorevleriGetir("", "", "")
		for _, digerGorev := range tumGorevler {
			digerBaglantilar, err := iy.veriYonetici.BaglantilariGetir(digerGorev.ID)
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

func (iy *IsYonetici) ProjeGorevSayisi(projeID string) (int, error) {
	gorevler, err := iy.veriYonetici.ProjeGorevleriGetir(projeID)
	if err != nil {
		return 0, err
	}
	return len(gorevler), nil
}

func (iy *IsYonetici) AktifProjeAyarla(projeID string) error {
	return iy.veriYonetici.AktifProjeAyarla(projeID)
}

func (iy *IsYonetici) AktifProjeGetir() (*Proje, error) {
	projeID, err := iy.veriYonetici.AktifProjeGetir()
	if err != nil {
		return nil, err
	}
	if projeID == "" {
		return nil, nil // Aktif proje yok
	}
	return iy.veriYonetici.ProjeGetir(projeID)
}

func (iy *IsYonetici) AktifProjeKaldir() error {
	return iy.veriYonetici.AktifProjeKaldir()
}

func (iy *IsYonetici) OzetAl() (*Ozet, error) {
	gorevler, err := iy.veriYonetici.GorevleriGetir("", "", "")
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.tasksFetchFailed", map[string]interface{}{"Error": err}))
	}

	projeler, err := iy.veriYonetici.ProjeleriGetir()
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

func (iy *IsYonetici) GorevBagimlilikEkle(kaynakID, hedefID, baglantiTipi string) (*Baglanti, error) {
	// Görevlerin var olup olmadığını kontrol et
	_, err := iy.veriYonetici.GorevGetir(kaynakID)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.sourceTaskNotFound", map[string]interface{}{"Error": err}))
	}
	_, err = iy.veriYonetici.GorevGetir(hedefID)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.targetTaskNotFound", map[string]interface{}{"Error": err}))
	}

	baglanti := &Baglanti{
		ID:             uuid.New().String(),
		SourceID:       kaynakID,
		TargetID:       hedefID,
		ConnectionType: baglantiTipi,
	}

	if err := iy.veriYonetici.BaglantiEkle(baglanti); err != nil {
		return nil, fmt.Errorf(i18n.TAddFailed("tr", "link", err))
	}

	return baglanti, nil
}

func (iy *IsYonetici) GorevBaglantilariGetir(gorevID string) ([]*Baglanti, error) {
	return iy.veriYonetici.BaglantilariGetir(gorevID)
}

// GorevBagimliMi görevi başlatmak için tüm bağımlılıkların tamamlanıp tamamlanmadığını kontrol eder
func (iy *IsYonetici) GorevBagimliMi(gorevID string) (bool, []string, error) {
	baglantilar, err := iy.veriYonetici.BaglantilariGetir(gorevID)
	if err != nil {
		return false, nil, fmt.Errorf(i18n.T("error.dependencyFetchFailed", map[string]interface{}{"Error": err}))
	}

	var tamamlanmamisBagimliliklar []string

	for _, baglanti := range baglantilar {
		// Bu görev hedef konumundaysa ve bağlantı tipi "onceki" ise
		if baglanti.TargetID == gorevID && baglanti.ConnectionType == "onceki" {
			// Kaynak görevin durumunu kontrol et
			kaynakGorev, err := iy.veriYonetici.GorevGetir(baglanti.SourceID)
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
func (iy *IsYonetici) TemplateListele(kategori string) ([]*GorevTemplate, error) {
	return iy.veriYonetici.TemplateListele(kategori)
}

// TemplatedenGorevOlustur template kullanarak görev oluşturur
func (iy *IsYonetici) TemplatedenGorevOlustur(templateID string, degerler map[string]string) (*Gorev, error) {
	return iy.veriYonetici.TemplatedenGorevOlustur(templateID, degerler)
}

// AltGorevOlustur mevcut bir görevin altına yeni görev oluşturur
func (iy *IsYonetici) AltGorevOlustur(parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	// Parent görevi kontrol et
	parent, err := iy.veriYonetici.GorevGetir(parentID)
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
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     sonTarih,
	}

	if err := iy.veriYonetici.GorevKaydet(gorev); err != nil {
		return nil, fmt.Errorf(i18n.TSaveFailed("tr", "subtask", err))
	}

	if len(etiketIsimleri) > 0 {
		etiketler, err := iy.veriYonetici.EtiketleriGetirVeyaOlustur(etiketIsimleri)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsProcessFailed", map[string]interface{}{"Error": err}))
		}
		if err := iy.veriYonetici.GorevEtiketleriniAyarla(gorev.ID, etiketler); err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsSetFailed", map[string]interface{}{"Error": err}))
		}
		gorev.Tags = etiketler
	}

	return gorev, nil
}

// GorevUstDegistir bir görevin üst görevini değiştirir
func (iy *IsYonetici) GorevUstDegistir(gorevID, yeniParentID string) error {
	// Görevi kontrol et
	gorev, err := iy.veriYonetici.GorevGetir(gorevID)
	if err != nil {
		return fmt.Errorf(i18n.TEntityNotFound("tr", "task", err))
	}

	// Yeni parent varsa kontrol et
	if yeniParentID != "" {
		parent, err := iy.veriYonetici.GorevGetir(yeniParentID)
		if err != nil {
			return fmt.Errorf(i18n.T("parentTaskNotFound"))
		}

		// Circular dependency kontrolü
		circular, err := iy.veriYonetici.DaireBagimliligiKontrolEt(gorevID, yeniParentID)
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

	return iy.veriYonetici.ParentIDGuncelle(gorevID, yeniParentID)
}

// GorevHiyerarsiGetir bir görevin tam hiyerarşi bilgilerini getirir
func (iy *IsYonetici) GorevHiyerarsiGetir(gorevID string) (*GorevHiyerarsi, error) {
	return iy.veriYonetici.GorevHiyerarsiGetir(gorevID)
}

// AltGorevleriGetir bir görevin doğrudan alt görevlerini getirir
func (iy *IsYonetici) AltGorevleriGetir(parentID string) ([]*Gorev, error) {
	return iy.veriYonetici.AltGorevleriGetir(parentID)
}

// TumAltGorevleriGetir bir görevin tüm alt görev ağacını getirir
func (iy *IsYonetici) TumAltGorevleriGetir(parentID string) ([]*Gorev, error) {
	return iy.veriYonetici.TumAltGorevleriGetir(parentID)
}
