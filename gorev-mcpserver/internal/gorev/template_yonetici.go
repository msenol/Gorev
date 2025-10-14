package gorev

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// TemplateOlustur yeni bir görev template'i oluşturur
func (vy *VeriYonetici) TemplateOlustur(template *GorevTemplate) error {
	template.ID = uuid.New().String()

	// Alanları JSON'a çevir
	alanlarJSON, err := json.Marshal(template.Fields)
	if err != nil {
		return fmt.Errorf(i18n.T("error.fieldsJsonFailed", map[string]interface{}{"Error": err}))
	}

	// Örnek değerleri JSON'a çevir
	ornekDegerlerJSON, err := json.Marshal(template.SampleValues)
	if err != nil {
		return fmt.Errorf(i18n.T("error.exampleValuesJsonFailed", map[string]interface{}{"Error": err}))
	}

	sorgu := `INSERT INTO gorev_templateleri 
		(id, name, definition, alias, default_title, description_template, fields, sample_values, category, active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = vy.db.Exec(sorgu, template.ID, template.Name, template.Definition, template.Alias,
		template.DefaultTitle, template.DescriptionTemplate,
		string(alanlarJSON), string(ornekDegerlerJSON), template.Category, template.Active)

	if err != nil {
		return fmt.Errorf(i18n.TCreateFailed("tr", "template", err))
	}

	return nil
}

// TemplateListele tüm active template'leri listeler
func (vy *VeriYonetici) TemplateListele(category string) ([]*GorevTemplate, error) {
	var sorgu string
	var args []interface{}

	if category != "" {
		sorgu = `SELECT id, name, definition, alias, default_title, description_template, 
				fields, sample_values, category, active 
				FROM gorev_templateleri WHERE active = 1 AND category = ? ORDER BY name`
		args = append(args, category)
	} else {
		sorgu = `SELECT id, name, definition, alias, default_title, description_template, 
				fields, sample_values, category, active 
				FROM gorev_templateleri WHERE active = 1 ORDER BY category, name`
	}

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, fmt.Errorf(i18n.TListFailed("tr", "template", err))
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: rows.Close failed: %v\n", err)
		}
	}()

	var templates []*GorevTemplate
	for rows.Next() {
		template := &GorevTemplate{}
		var alanlarJSON, ornekDegerlerJSON string

		err := rows.Scan(&template.ID, &template.Name, &template.Definition, &template.Alias,
			&template.DefaultTitle, &template.DescriptionTemplate,
			&alanlarJSON, &ornekDegerlerJSON, &template.Category, &template.Active)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.templateReadFailed", map[string]interface{}{"Error": err}))
		}

		// Alanları parse et
		if err := json.Unmarshal([]byte(alanlarJSON), &template.Fields); err != nil {
			return nil, fmt.Errorf(i18n.TParseFailed("tr", "fields", err))
		}

		// Örnek değerleri parse et
		if err := json.Unmarshal([]byte(ornekDegerlerJSON), &template.SampleValues); err != nil {
			return nil, fmt.Errorf(i18n.T("error.exampleValuesParseFailed", map[string]interface{}{"Error": err}))
		}

		templates = append(templates, template)
	}

	return templates, nil
}

// TemplateGetir belirli bir template'i getirir
func (vy *VeriYonetici) TemplateGetir(templateID string) (*GorevTemplate, error) {
	template := &GorevTemplate{}
	var alanlarJSON, ornekDegerlerJSON string

	sorgu := `SELECT id, name, definition, alias, default_title, description_template, 
			fields, sample_values, category, active 
			FROM gorev_templateleri WHERE id = ?`

	err := vy.db.QueryRow(sorgu, templateID).Scan(
		&template.ID, &template.Name, &template.Definition, &template.Alias,
		&template.DefaultTitle, &template.DescriptionTemplate,
		&alanlarJSON, &ornekDegerlerJSON, &template.Category, &template.Active)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(i18n.T("error.templateNotFoundId", map[string]interface{}{"Id": templateID}))
		}
		return nil, fmt.Errorf(i18n.TFetchFailed("tr", "template", err))
	}

	// Alanları parse et
	if err := json.Unmarshal([]byte(alanlarJSON), &template.Fields); err != nil {
		return nil, fmt.Errorf(i18n.TParseFailed("tr", "fields", err))
	}

	// Örnek değerleri parse et
	if err := json.Unmarshal([]byte(ornekDegerlerJSON), &template.SampleValues); err != nil {
		return nil, fmt.Errorf(i18n.T("error.exampleValuesParseFailed", map[string]interface{}{"Error": err}))
	}

	return template, nil
}

// TemplateAliasIleGetir alias ile template getirir
func (vy *VeriYonetici) TemplateAliasIleGetir(alias string) (*GorevTemplate, error) {
	template := &GorevTemplate{}
	var alanlarJSON, ornekDegerlerJSON string

	sorgu := `SELECT id, name, definition, alias, default_title, description_template, 
			fields, sample_values, category, active 
			FROM gorev_templateleri WHERE alias = ? AND active = 1`

	err := vy.db.QueryRow(sorgu, alias).Scan(
		&template.ID, &template.Name, &template.Definition, &template.Alias,
		&template.DefaultTitle, &template.DescriptionTemplate,
		&alanlarJSON, &ornekDegerlerJSON, &template.Category, &template.Active)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(i18n.T("error.templateNotFoundAlias", map[string]interface{}{"Alias": alias}))
		}
		return nil, fmt.Errorf(i18n.TFetchFailed("tr", "template", err))
	}

	// Alanları parse et
	if err := json.Unmarshal([]byte(alanlarJSON), &template.Fields); err != nil {
		return nil, fmt.Errorf(i18n.TParseFailed("tr", "fields", err))
	}

	// Örnek değerleri parse et
	if err := json.Unmarshal([]byte(ornekDegerlerJSON), &template.SampleValues); err != nil {
		return nil, fmt.Errorf(i18n.T("error.exampleValuesParseFailed", map[string]interface{}{"Error": err}))
	}

	return template, nil
}

// TemplateIDVeyaAliasIleGetir ID veya alias ile template getirir
func (vy *VeriYonetici) TemplateIDVeyaAliasIleGetir(idOrAlias string) (*GorevTemplate, error) {
	// Önce ID olarak dene
	template, err := vy.TemplateGetir(idOrAlias)
	if err == nil {
		return template, nil
	}

	// Sonra alias olarak dene
	return vy.TemplateAliasIleGetir(idOrAlias)
}

// TemplatedenGorevOlustur template kullanarak görev oluşturur
func (vy *VeriYonetici) TemplatedenGorevOlustur(templateID string, degerler map[string]string) (*Gorev, error) {
	// Template'i ID veya alias ile getir
	template, err := vy.TemplateIDVeyaAliasIleGetir(templateID)
	if err != nil {
		return nil, err
	}

	// Zorunlu alanları kontrol et
	for _, alan := range template.Fields {
		if alan.Required {
			if _, ok := degerler[alan.Name]; !ok {
				return nil, fmt.Errorf(i18n.T("error.requiredFieldMissing", map[string]interface{}{"Field": alan.Name}))
			}
		}
	}

	// Başlık oluştur
	baslik := template.DefaultTitle
	for key, value := range degerler {
		baslik = strings.ReplaceAll(baslik, "{{"+key+"}}", value)
	}

	// Açıklama oluştur
	aciklama := template.DescriptionTemplate
	for key, value := range degerler {
		aciklama = strings.ReplaceAll(aciklama, "{{"+key+"}}", value)
	}

	// Varsayılan değerleri uygula
	oncelik := constants.PriorityMedium
	if val, ok := degerler["priority"]; ok {
		oncelik = val
	}

	var sonTarih *time.Time
	if val, ok := degerler["due_date"]; ok {
		if t, err := time.Parse(constants.DateFormatISO, val); err == nil {
			sonTarih = &t
		}
	}

	// Etiketleri ayır
	var etiketler []string
	if val, ok := degerler["etiketler"]; ok {
		etiketler = strings.Split(val, ",")
		for i := range etiketler {
			etiketler[i] = strings.TrimSpace(etiketler[i])
		}
	}

	// Görev oluştur
	gorev := &Gorev{
		Title:       baslik,
		Description: aciklama,
		Priority:    oncelik,
		Status:      constants.TaskStatusPending,
	}

	// ProjeID'yi ayarla
	if val, ok := degerler["proje_id"]; ok && val != "" {
		gorev.ProjeID = val
	} else {
		// Aktif projeyi kullan
		aktifProjeID, err := vy.AktifProjeGetir()
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.activeProjectFetchFailed", map[string]interface{}{"Error": err}))
		}
		if aktifProjeID == "" {
			return nil, fmt.Errorf(i18n.T("error.noActiveProjectSet"))
		}
		gorev.ProjeID = aktifProjeID
	}

	// ID ve tarihler ayarla
	gorev.ID = uuid.New().String()
	gorev.CreatedAt = time.Now()
	gorev.UpdatedAt = time.Now()
	gorev.DueDate = sonTarih

	// Görevi kaydet
	err = vy.GorevKaydet(gorev)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.taskSaveFromTemplateFailed", map[string]interface{}{"Error": err}))
	}

	// Etiketleri ayarla
	if len(etiketler) > 0 {
		etiketNesneleri, err := vy.EtiketleriGetirVeyaOlustur(etiketler)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsCreateFromTemplateFailed", map[string]interface{}{"Error": err}))
		}

		err = vy.GorevEtiketleriniAyarla(gorev.ID, etiketNesneleri)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.taskTagsSetFromTemplateFailed", map[string]interface{}{"Error": err}))
		}

		gorev.Tags = etiketNesneleri
	}

	return gorev, nil
}

// VarsayilanTemplateleriOlustur varsayılan template'leri oluşturur
func (vy *VeriYonetici) VarsayilanTemplateleriOlustur() error {
	templates := []*GorevTemplate{
		{
			Name:         "Bug Raporu",
			Definition:   "Yazılım hatası bildirimi için detaylı template",
			Alias:        "bug",
			DefaultTitle: "🐛 [{{modul}}] {{title}}",
			DescriptionTemplate: `## 🐛 Hata Açıklaması
{{description}}

## 📍 Nerede Oluşuyor?
**Modül/Bileşen:** {{modul}}
**Ortam:** {{ortam}}

## 🔄 Tekrar Üretme Adımları
{{adimlar}}

## ✅ Beklenen Davranış
{{beklenen}}

## ❌ Mevcut Davranış
{{mevcut}}

## 📸 Ekran Görüntüleri/Loglar
{{ekler}}

## 🔧 Olası Çözüm
{{cozum}}

## 📊 Öncelik: {{priority}}
## 🏷️ Etiketler: {{etiketler}}`,
			Fields: []TemplateAlan{
				{Name: "title", Type: "text", Required: true},
				{Name: "description", Type: "text", Required: true},
				{Name: "modul", Type: "text", Required: true},
				{Name: "ortam", Type: "select", Required: true, Options: constants.ValidEnvironments},
				{Name: "adimlar", Type: "text", Required: true},
				{Name: "beklenen", Type: "text", Required: true},
				{Name: "mevcut", Type: "text", Required: true},
				{Name: "ekler", Type: "text", Required: false},
				{Name: "cozum", Type: "text", Required: false},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
				{Name: "etiketler", Type: "text", Required: false, Default: "bug"},
			},
			Category: "Teknik",
			Active:   true,
		},
		{
			Name:         "Özellik İsteği",
			Definition:   "Yeni özellik veya geliştirme isteği için template",
			Alias:        "feature",
			DefaultTitle: "✨ {{title}}",
			DescriptionTemplate: `## ✨ Özellik Açıklaması
{{description}}

## 🎯 Amaç ve Faydalar
{{amac}}

## 👥 Hedef Kullanıcılar
{{kullanicilar}}

## 📋 Kabul Kriterleri
{{kriterler}}

## 🎨 UI/UX Düşünceleri
{{ui_ux}}

## 🔗 İlgili Özellikler/Modüller
{{ilgili}}

## 📊 Tahmini Efor
{{efor}}

## 🏷️ Etiketler: {{etiketler}}`,
			Fields: []TemplateAlan{
				{Name: "title", Type: "text", Required: true},
				{Name: "description", Type: "text", Required: true},
				{Name: "amac", Type: "text", Required: true},
				{Name: "kullanicilar", Type: "text", Required: true},
				{Name: "kriterler", Type: "text", Required: true},
				{Name: "ui_ux", Type: "text", Required: false},
				{Name: "ilgili", Type: "text", Required: false},
				{Name: "efor", Type: "select", Required: false, Options: constants.ValidEffortLevels},
				{Name: "due_date", Type: "date", Required: false},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
				{Name: "etiketler", Type: "text", Required: false, Default: "özellik"},
			},
			Category: "Özellik",
			Active:   true,
		},
		{
			Name:         "Teknik Borç",
			Definition:   "Refaktöring veya teknik iyileştirme için template",
			Alias:        "debt",
			DefaultTitle: "🔧 [{{alan}}] {{title}}",
			DescriptionTemplate: `## 🔧 Teknik Borç Açıklaması
{{description}}

## 📍 Etkilenen Alan
**Alan/Modül:** {{alan}}
**Dosyalar:** {{dosyalar}}

## ❓ Neden Gerekli?
{{neden}}

## 📊 Mevcut Durum Analizi
{{analiz}}

## 🎯 Önerilen Çözüm
{{cozum}}

## ⚠️ Riskler
{{riskler}}

## 📈 Beklenen İyileştirmeler
{{iyilestirmeler}}

## ⏱️ Tahmini Süre: {{sure}}
## 🏷️ Etiketler: {{etiketler}}`,
			Fields: []TemplateAlan{
				{Name: "title", Type: "text", Required: true},
				{Name: "description", Type: "text", Required: true},
				{Name: "alan", Type: "text", Required: true},
				{Name: "dosyalar", Type: "text", Required: false},
				{Name: "neden", Type: "text", Required: true},
				{Name: "analiz", Type: "text", Required: true},
				{Name: "cozum", Type: "text", Required: true},
				{Name: "riskler", Type: "text", Required: false},
				{Name: "iyilestirmeler", Type: "text", Required: true},
				{Name: "sure", Type: "select", Required: false, Options: []string{"1 gün", "2-3 gün", "1 hafta", "2+ hafta"}},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
				{Name: "etiketler", Type: "text", Required: false, Default: "teknik-borç,refaktöring"},
			},
			Category: "Teknik",
			Active:   true,
		},
		{
			Name:         "Araştırma Görevi",
			Definition:   "Teknoloji veya çözüm araştırması için template",
			Alias:        "research",
			DefaultTitle: "🔍 {{konu}} Araştırması",
			DescriptionTemplate: `## 🔍 Araştırma Konusu
{{konu}}

## 🎯 Araştırma Amacı
{{amac}}

## ❓ Cevaplanması Gereken Sorular
{{sorular}}

## 📚 Araştırılacak Kaynaklar
{{kaynaklar}}

## 🔄 Alternatifler
{{alternatifler}}

## ⚖️ Değerlendirme Kriterleri
{{kriterler}}

## 📅 Bitiş Tarihi: {{son_tarih}}
## 🏷️ Etiketler: {{etiketler}}`,
			Fields: []TemplateAlan{
				{Name: "konu", Type: "text", Required: true},
				{Name: "amac", Type: "text", Required: true},
				{Name: "sorular", Type: "text", Required: true},
				{Name: "kaynaklar", Type: "text", Required: false},
				{Name: "alternatifler", Type: "text", Required: false},
				{Name: "kriterler", Type: "text", Required: true},
				{Name: "due_date", Type: "date", Required: false},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
				{Name: "etiketler", Type: "text", Required: false, Default: "araştırma"},
			},
			Category: "Araştırma",
			Active:   true,
		},
		// Yeni template'ler - Template zorunluluğu için eklendi
		{
			Name:         "Bug Raporu v2",
			Definition:   "Gelişmiş bug raporu - detaylı adımlar ve environment bilgisi",
			Alias:        "bug2",
			DefaultTitle: "🐛 [{{severity}}] {{modul}}: {{title}}",
			DescriptionTemplate: `## 🐛 Hata Özeti
{{description}}

## 🔄 Tekrar Üretme Adımları
{{steps_to_reproduce}}

## ✅ Beklenen Davranış
{{expected_behavior}}

## ❌ Gerçekleşen Davranış
{{actual_behavior}}

## 💻 Ortam Bilgileri
- **İşletim Sistemi:** {{os_version}}
- **Tarayıcı/Client:** {{client_info}}
- **Server Version:** {{server_version}}
- **Database:** {{db_info}}

## 🚨 Hata Derecesi
**Severity:** {{severity}}
**Etkilenen Kullanıcı Sayısı:** {{affected_users}}

## 📸 Ekler
{{attachments}}

## 🔧 Geçici Çözüm
{{workaround}}`,
			Fields: []TemplateAlan{
				{Name: "title", Type: "text", Required: true},
				{Name: "description", Type: "text", Required: true},
				{Name: "modul", Type: "text", Required: true},
				{Name: "steps_to_reproduce", Type: "text", Required: true},
				{Name: "expected_behavior", Type: "text", Required: true},
				{Name: "actual_behavior", Type: "text", Required: true},
				{Name: "os_version", Type: "text", Required: true},
				{Name: "client_info", Type: "text", Required: true},
				{Name: "server_version", Type: "text", Required: true},
				{Name: "db_info", Type: "text", Required: false},
				{Name: "severity", Type: "select", Required: true, Options: []string{"critical", "high", "medium", "low"}},
				{Name: "affected_users", Type: "text", Required: true},
				{Name: "attachments", Type: "text", Required: false},
				{Name: "workaround", Type: "text", Required: false},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityHigh, Options: constants.GetValidPriorities()},
				{Name: "etiketler", Type: "text", Required: false, Default: "bug,production"},
			},
			Category: "Bug",
			Active:   true,
		},
		{
			Name:         "Spike Araştırma",
			Definition:   "Time-boxed teknik araştırma ve proof-of-concept çalışmaları",
			Alias:        "spike",
			DefaultTitle: "🔬 [SPIKE] {{research_question}}",
			DescriptionTemplate: `## 🔬 Araştırma Sorusu
{{research_question}}

## 🎯 Başarı Kriterleri
{{success_criteria}}

## ⏰ Time Box
**Maksimum Süre:** {{time_box}}
**Karar Tarihi:** {{decision_deadline}}

## 🔍 Araştırma Planı
{{research_plan}}

## 📊 Beklenen Çıktılar
{{expected_outputs}}

## ⚡ Riskler ve Varsayımlar
{{risks_assumptions}}

## 🏷️ Etiketler: {{etiketler}}`,
			Fields: []TemplateAlan{
				{Name: "research_question", Type: "text", Required: true},
				{Name: "success_criteria", Type: "text", Required: true},
				{Name: "time_box", Type: "select", Required: true, Options: []string{"4 saat", "1 gün", "2 gün", "3 gün", "1 hafta"}},
				{Name: "decision_deadline", Type: "date", Required: true},
				{Name: "research_plan", Type: "text", Required: true},
				{Name: "expected_outputs", Type: "text", Required: true},
				{Name: "risks_assumptions", Type: "text", Required: false},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityHigh, Options: constants.GetValidPriorities()},
				{Name: "etiketler", Type: "text", Required: false, Default: "spike,research,poc"},
			},
			Category: "Araştırma",
			Active:   true,
		},
		{
			Name:         "Performans Sorunu",
			Definition:   "Performans problemleri ve optimizasyon görevleri",
			Alias:        "performance",
			DefaultTitle: "⚡ [PERF] {{metric_affected}}: {{title}}",
			DescriptionTemplate: `## ⚡ Performans Sorunu
{{description}}

## 📊 Etkilenen Metrik
**Metrik:** {{metric_affected}}
**Mevcut Değer:** {{current_value}}
**Hedef Değer:** {{target_value}}
**Kabul Edilebilir Değer:** {{acceptable_value}}

## 📏 Ölçüm Yöntemi
{{measurement_method}}

## 👥 Kullanıcı Etkisi
{{user_impact}}

## 🔍 Kök Neden Analizi
{{root_cause}}

## 💡 Önerilen Çözümler
{{proposed_solutions}}

## ⚠️ Trade-offs
{{tradeoffs}}

## 🏷️ Etiketler: {{etiketler}}`,
			Fields: []TemplateAlan{
				{Name: "title", Type: "text", Required: true},
				{Name: "description", Type: "text", Required: true},
				{Name: "metric_affected", Type: "select", Required: true, Options: []string{"response_time", "throughput", "cpu_usage", "memory_usage", "database_query", "page_load", "api_latency"}},
				{Name: "current_value", Type: "text", Required: true},
				{Name: "target_value", Type: "text", Required: true},
				{Name: "acceptable_value", Type: "text", Required: false},
				{Name: "measurement_method", Type: "text", Required: true},
				{Name: "user_impact", Type: "text", Required: true},
				{Name: "root_cause", Type: "text", Required: false},
				{Name: "proposed_solutions", Type: "text", Required: true},
				{Name: "tradeoffs", Type: "text", Required: false},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityHigh, Options: constants.GetValidPriorities()},
				{Name: "etiketler", Type: "text", Required: false, Default: "performance,optimization"},
			},
			Category: "Teknik",
			Active:   true,
		},
		{
			Name:         "Güvenlik Düzeltmesi",
			Definition:   "Güvenlik açıkları ve düzeltmeleri için özel template",
			Alias:        "security",
			DefaultTitle: "🔒 [SEC-{{severity}}] {{vulnerability_type}}: {{title}}",
			DescriptionTemplate: `## 🔒 Güvenlik Açığı
{{description}}

## 🎯 Açık Tipi
**Category:** {{vulnerability_type}}
**CVSS Score:** {{cvss_score}}
**Severity:** {{severity}}

## 🔍 Etkilenen Bileşenler
{{affected_components}}

## 💥 Potansiyel Etki
{{potential_impact}}

## 🛡️ Azaltma Adımları
{{mitigation_steps}}

## ✅ Test Gereksinimleri
{{testing_requirements}}

## 📋 Güvenlik Kontrol Listesi
- [ ] Güvenlik testi yapıldı
- [ ] Penetrasyon testi gerekli mi?
- [ ] Security review tamamlandı
- [ ] Dokümantasyon güncellendi

## 🚨 Disclosure Timeline
{{disclosure_timeline}}

## 🏷️ Etiketler: {{etiketler}}`,
			Fields: []TemplateAlan{
				{Name: "title", Type: "text", Required: true},
				{Name: "description", Type: "text", Required: true},
				{Name: "vulnerability_type", Type: "select", Required: true, Options: []string{"SQL Injection", "XSS", "CSRF", "Authentication", "Authorization", "Data Exposure", "Misconfiguration", "Dependency", "Other"}},
				{Name: "cvss_score", Type: "text", Required: false},
				{Name: "severity", Type: "select", Required: true, Options: []string{"critical", "high", "medium", "low"}},
				{Name: "affected_components", Type: "text", Required: true},
				{Name: "potential_impact", Type: "text", Required: true},
				{Name: "mitigation_steps", Type: "text", Required: true},
				{Name: "testing_requirements", Type: "text", Required: true},
				{Name: "disclosure_timeline", Type: "text", Required: false},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityHigh, Options: []string{constants.PriorityHigh}}, // Güvenlik her zaman yüksek
				{Name: "etiketler", Type: "text", Required: false, Default: "security,vulnerability"},
			},
			Category: "Güvenlik",
			Active:   true,
		},
		{
			Name:         "Refactoring",
			Definition:   "Kod kalitesi ve mimari iyileştirmeler",
			Alias:        "refactor",
			DefaultTitle: "♻️ [REFACTOR] {{code_smell}}: {{title}}",
			DescriptionTemplate: `## ♻️ Refactoring Özeti
{{description}}

## 🦨 Code Smell Tipi
{{code_smell_type}}

## 📁 Etkilenen Dosyalar
{{affected_files}}

## 🎯 Refactoring Stratejisi
{{refactoring_strategy}}

## ✅ Başarı Kriterleri
- [ ] Tüm testler geçiyor
- [ ] Kod coverage düşmedi
- [ ] Performance etkilenmedi
- [ ] API uyumluluğu korundu

## ⚠️ Risk Değerlendirmesi
**Risk Seviyesi:** {{risk_level}}
**Etki Alanı:** {{impact_scope}}

## 🔄 Rollback Planı
{{rollback_plan}}

## 📊 Metrikler
- **Mevcut Cyclomatic Complexity:** {{current_complexity}}
- **Hedef Complexity:** {{target_complexity}}
- **Mevcut Code Coverage:** {{current_coverage}}

## 🏷️ Etiketler: {{etiketler}}`,
			Fields: []TemplateAlan{
				{Name: "title", Type: "text", Required: true},
				{Name: "description", Type: "text", Required: true},
				{Name: "code_smell", Type: "select", Required: true, Options: []string{"Long Method", "Large Class", "Duplicate Code", "Dead Code", "Complex Conditionals", "Feature Envy", "Data Clumps", "Primitive Obsession", "Switch Statements", "Parallel Inheritance", "Lazy Class", "Speculative Generality", "Message Chains", "Middle Man", "Other"}},
				{Name: "code_smell_type", Type: "text", Required: true},
				{Name: "affected_files", Type: "text", Required: true},
				{Name: "refactoring_strategy", Type: "text", Required: true},
				{Name: "risk_level", Type: "select", Required: true, Options: []string{"low", "medium", "high"}},
				{Name: "impact_scope", Type: "text", Required: true},
				{Name: "rollback_plan", Type: "text", Required: true},
				{Name: "current_complexity", Type: "text", Required: false},
				{Name: "target_complexity", Type: "text", Required: false},
				{Name: "current_coverage", Type: "text", Required: false},
				{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
				{Name: "etiketler", Type: "text", Required: false, Default: "refactoring,code-quality"},
			},
			Category: "Teknik",
			Active:   true,
		},
	}

	for _, template := range templates {
		// Generate UUID for template
		template.ID = uuid.New().String()

		// Check if template with this name already exists
		existingTemplates, err := vy.TemplateListele("")
		if err != nil {
			return fmt.Errorf(i18n.TListFailed("tr", "template", err))
		}

		exists := false
		for _, existing := range existingTemplates {
			if existing.Name == template.Name {
				exists = true
				break
			}
		}

		if exists {
			// Template already exists, skip creation
			continue
		}

		if err := vy.TemplateOlustur(template); err != nil {
			return fmt.Errorf(i18n.T("error.defaultTemplateCreateFailed", map[string]interface{}{"Template": template.Name, "Error": err}))
		}
	}

	return nil
}
