package gorev

import (
	"context"
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
func (vy *VeriYonetici) TemplateOlustur(ctx context.Context, template *GorevTemplate) error {
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
		(id, name, definition, alias, default_title, description_template, fields, sample_values, category, active, language_code, base_template_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Ensure required fields have defaults
	languageCode := template.LanguageCode
	if languageCode == "" {
		languageCode = "tr"
	}

	baseTemplateID := template.BaseTemplateID
	if baseTemplateID == nil {
		baseTemplateID = &template.ID
	}

	_, err = vy.db.Exec(sorgu, template.ID, template.Name, template.Definition, template.Alias,
		template.DefaultTitle, template.DescriptionTemplate,
		string(alanlarJSON), string(ornekDegerlerJSON), template.Category, template.Active,
		languageCode, baseTemplateID)

	if err != nil {
		return fmt.Errorf(i18n.TCreateFailed(i18n.FromContext(ctx), "template", err))
	}

	return nil
}

// TemplateListele tüm active template'leri listeler (language-aware)
func (vy *VeriYonetici) TemplateListele(ctx context.Context, category string) ([]*GorevTemplate, error) {
	lang := i18n.FromContext(ctx)
	if lang == "" {
		lang = "tr"
	}

	var sorgu string
	var args []interface{}

	if category != "" {
		sorgu = `SELECT id, name, definition, alias, default_title, description_template,
				fields, sample_values, category, active, language_code, base_template_id
				FROM gorev_templateleri WHERE active = 1 AND category = ? AND language_code = ? ORDER BY name`
		args = append(args, category, lang)
	} else {
		sorgu = `SELECT id, name, definition, alias, default_title, description_template,
				fields, sample_values, category, active, language_code, base_template_id
				FROM gorev_templateleri WHERE active = 1 AND language_code = ? ORDER BY category, name`
		args = append(args, lang)
	}

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, fmt.Errorf(i18n.TListFailed(i18n.FromContext(ctx), "template", err))
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
			&alanlarJSON, &ornekDegerlerJSON, &template.Category, &template.Active,
			&template.LanguageCode, &template.BaseTemplateID)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.templateReadFailed", map[string]interface{}{"Error": err}))
		}

		// Alanları parse et
		if err := json.Unmarshal([]byte(alanlarJSON), &template.Fields); err != nil {
			return nil, fmt.Errorf(i18n.TParseFailed(i18n.FromContext(ctx), "fields", err))
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
func (vy *VeriYonetici) TemplateGetir(ctx context.Context, templateID string) (*GorevTemplate, error) {
	template := &GorevTemplate{}
	var alanlarJSON, ornekDegerlerJSON string

	sorgu := `SELECT id, name, definition, alias, default_title, description_template,
			fields, sample_values, category, active, language_code, base_template_id
			FROM gorev_templateleri WHERE id = ?`

	err := vy.db.QueryRow(sorgu, templateID).Scan(
		&template.ID, &template.Name, &template.Definition, &template.Alias,
		&template.DefaultTitle, &template.DescriptionTemplate,
		&alanlarJSON, &ornekDegerlerJSON, &template.Category, &template.Active,
		&template.LanguageCode, &template.BaseTemplateID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(i18n.T("error.templateNotFoundId", map[string]interface{}{"Id": templateID}))
		}
		return nil, fmt.Errorf(i18n.TFetchFailed(i18n.FromContext(ctx), "template", err))
	}

	// Alanları parse et
	if err := json.Unmarshal([]byte(alanlarJSON), &template.Fields); err != nil {
		return nil, fmt.Errorf(i18n.TParseFailed(i18n.FromContext(ctx), "fields", err))
	}

	// Örnek değerleri parse et
	if err := json.Unmarshal([]byte(ornekDegerlerJSON), &template.SampleValues); err != nil {
		return nil, fmt.Errorf(i18n.T("error.exampleValuesParseFailed", map[string]interface{}{"Error": err}))
	}

	return template, nil
}

// TemplateAliasIleGetir alias ile template getirir (language-aware)
func (vy *VeriYonetici) TemplateAliasIleGetir(ctx context.Context, alias string) (*GorevTemplate, error) {
	lang := i18n.FromContext(ctx)
	if lang == "" {
		lang = "tr"
	}

	template := &GorevTemplate{}
	var alanlarJSON, ornekDegerlerJSON string

	sorgu := `SELECT id, name, definition, alias, default_title, description_template,
			fields, sample_values, category, active, language_code, base_template_id
			FROM gorev_templateleri WHERE alias = ? AND active = 1 AND language_code = ?`

	err := vy.db.QueryRow(sorgu, alias, lang).Scan(
		&template.ID, &template.Name, &template.Definition, &template.Alias,
		&template.DefaultTitle, &template.DescriptionTemplate,
		&alanlarJSON, &ornekDegerlerJSON, &template.Category, &template.Active,
		&template.LanguageCode, &template.BaseTemplateID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(i18n.T("error.templateNotFoundAlias", map[string]interface{}{"Alias": alias}))
		}
		return nil, fmt.Errorf(i18n.TFetchFailed(i18n.FromContext(ctx), "template", err))
	}

	// Alanları parse et
	if err := json.Unmarshal([]byte(alanlarJSON), &template.Fields); err != nil {
		return nil, fmt.Errorf(i18n.TParseFailed(i18n.FromContext(ctx), "fields", err))
	}

	// Örnek değerleri parse et
	if err := json.Unmarshal([]byte(ornekDegerlerJSON), &template.SampleValues); err != nil {
		return nil, fmt.Errorf(i18n.T("error.exampleValuesParseFailed", map[string]interface{}{"Error": err}))
	}

	return template, nil
}

// TemplateIDVeyaAliasIleGetir ID veya alias ile template getirir
func (vy *VeriYonetici) TemplateIDVeyaAliasIleGetir(ctx context.Context, idOrAlias string) (*GorevTemplate, error) {
	// Önce ID olarak dene
	template, err := vy.TemplateGetir(ctx, idOrAlias)
	if err == nil {
		return template, nil
	}

	// Sonra alias olarak dene
	return vy.TemplateAliasIleGetir(ctx, idOrAlias)
}

// TemplatedenGorevOlustur template kullanarak görev oluşturur
func (vy *VeriYonetici) TemplatedenGorevOlustur(ctx context.Context, templateID string, degerler map[string]string) (*Gorev, error) {
	// Template'i ID veya alias ile getir
	template, err := vy.TemplateIDVeyaAliasIleGetir(ctx, templateID)
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
	if val, ok := degerler["tags"]; ok {
		etiketler = strings.Split(val, ",")
		for i := range etiketler {
			etiketler[i] = strings.TrimSpace(etiketler[i])
		}
	}

	// Get workspace_id if injected by IsYonetici
	workspaceID := ""
	if val, ok := degerler["_workspace_id"]; ok {
		workspaceID = val
		delete(degerler, "_workspace_id") // Clean up internal key
	}

	// Görev oluştur
	gorev := &Gorev{
		Title:       baslik,
		Description: aciklama,
		Priority:    oncelik,
		Status:      constants.TaskStatusPending,
		WorkspaceID: workspaceID,
	}

	// ProjeID'yi ayarla
	if val, ok := degerler["project_id"]; ok && val != "" {
		gorev.ProjeID = val
	} else {
		// Aktif projeyi kullan
		aktifProjeID, err := vy.AktifProjeGetir(ctx)
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
	err = vy.GorevKaydet(ctx, gorev)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.taskSaveFromTemplateFailed", map[string]interface{}{"Error": err}))
	}

	// Etiketleri ayarla
	if len(etiketler) > 0 {
		etiketNesneleri, err := vy.EtiketleriGetirVeyaOlustur(ctx, etiketler)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.tagsCreateFromTemplateFailed", map[string]interface{}{"Error": err}))
		}

		err = vy.GorevEtiketleriniAyarla(ctx, gorev.ID, etiketNesneleri)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.taskTagsSetFromTemplateFailed", map[string]interface{}{"Error": err}))
		}

		gorev.Tags = etiketNesneleri
	}

	return gorev, nil
}

// VarsayilanTemplateleriOlustur varsayılan template'leri TR/EN çifti olarak oluşturur
func (vy *VeriYonetici) VarsayilanTemplateleriOlustur(ctx context.Context) error {
	// Define all template groups with base IDs
	templateGroups := []struct {
		BaseTemplateID string
		TR             *GorevTemplate
		EN             *GorevTemplate
	}{
		// Bug Report Template
		{
			BaseTemplateID: "bug-report",
			TR: &GorevTemplate{
				Name:         "Bug Raporu",
				Definition:   "Yazılım hatası bildirimi için detaylı template",
				Alias:        "bug",
				DefaultTitle: "🐛 [{{module}}] {{title}}",
				DescriptionTemplate: `## 🐛 Hata Açıklaması
{{description}}

## 📍 Nerede Oluşuyor?
**Modül/Bileşen:** {{module}}
**Ortam:** {{environment}}

## 🔄 Tekrar Üretme Adımları
{{steps}}

## ✅ Beklenen Davranış
{{expected}}

## ❌ Mevcut Davranış
{{actual}}

## 📸 Ekran Görüntüleri/Loglar
{{attachments}}

## 🔧 Olası Çözüm
{{solution}}

## 📊 Öncelik: {{priority}}
## 🏷️ Tags: {{tags}}`,
				Fields: []TemplateAlan{
					{Name: "title", Type: "text", Required: true},
					{Name: "description", Type: "text", Required: true},
					{Name: "module", Type: "text", Required: true},
					{Name: "environment", Type: "select", Required: true, Options: constants.ValidEnvironments},
					{Name: "steps", Type: "text", Required: true},
					{Name: "expected", Type: "text", Required: true},
					{Name: "actual", Type: "text", Required: true},
					{Name: "attachments", Type: "text", Required: false},
					{Name: "solution", Type: "text", Required: false},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "bug"},
				},
				Category: "Teknik",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Bug Report",
				Definition:   "Detailed template for software bug reporting",
				Alias:        "bug",
				DefaultTitle: "🐛 [{{module}}] {{title}}",
				DescriptionTemplate: `## 🐛 Bug Description
{{description}}

## 📍 Where Does It Occur?
**Module/Component:** {{module}}
**Environment:** {{environment}}

## 🔄 Reproduction Steps
{{steps}}

## ✅ Expected Behavior
{{expected}}

## ❌ Actual Behavior
{{actual}}

## 📸 Screenshots/Logs
{{attachments}}

## 🔧 Possible Solution
{{solution}}

## 📊 Priority: {{priority}}
## 🏷️ Tags: {{tags}}`,
				Fields: []TemplateAlan{
					{Name: "title", Type: "text", Required: true},
					{Name: "description", Type: "text", Required: true},
					{Name: "module", Type: "text", Required: true},
					{Name: "environment", Type: "select", Required: true, Options: constants.ValidEnvironments},
					{Name: "steps", Type: "text", Required: true},
					{Name: "expected", Type: "text", Required: true},
					{Name: "actual", Type: "text", Required: true},
					{Name: "attachments", Type: "text", Required: false},
					{Name: "solution", Type: "text", Required: false},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "bug"},
				},
				Category: "Technical",
				Active:   true,
			},
		},
		// Feature Request Template
		{
			BaseTemplateID: "feature-request",
			TR: &GorevTemplate{
				Name:         "Özellik İsteği",
				Definition:   "Yeni özellik veya geliştirme isteği için template",
				Alias:        "feature",
				DefaultTitle: "✨ {{title}}",
				DescriptionTemplate: `## ✨ Özellik Açıklaması
{{description}}

## 🎯 Amaç ve Faydalar
{{purpose}}

## 👥 Hedef Kullanıcılar
{{users}}

## 📋 Kabul Kriterleri
{{criteria}}

## 🎨 UI/UX Düşünceleri
{{ui_ux}}

## 🔗 İlgili Özellikler/Modüller
{{related}}

## 📊 Tahmini Efor
{{effort}}

## 🏷️ Tags: {{tags}}`,
				Fields: []TemplateAlan{
					{Name: "title", Type: "text", Required: true},
					{Name: "description", Type: "text", Required: true},
					{Name: "purpose", Type: "text", Required: true},
					{Name: "users", Type: "text", Required: true},
					{Name: "criteria", Type: "text", Required: true},
					{Name: "ui_ux", Type: "text", Required: false},
					{Name: "related", Type: "text", Required: false},
					{Name: "effort", Type: "select", Required: false, Options: constants.ValidEffortLevels},
					{Name: "due_date", Type: "date", Required: false},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "özellik"},
				},
				Category: "Özellik",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Feature Request",
				Definition:   "Template for new feature or enhancement requests",
				Alias:        "feature",
				DefaultTitle: "✨ {{title}}",
				DescriptionTemplate: `## ✨ Feature Description
{{description}}

## 🎯 Purpose and Benefits
{{purpose}}

## 👥 Target Users
{{users}}

## 📋 Acceptance Criteria
{{criteria}}

## 🎨 UI/UX Thoughts
{{ui_ux}}

## 🔗 Related Features/Modules
{{related}}

## 📊 Estimated Effort
{{effort}}

## 🏷️ Tags: {{tags}}`,
				Fields: []TemplateAlan{
					{Name: "title", Type: "text", Required: true},
					{Name: "description", Type: "text", Required: true},
					{Name: "purpose", Type: "text", Required: true},
					{Name: "users", Type: "text", Required: true},
					{Name: "criteria", Type: "text", Required: true},
					{Name: "ui_ux", Type: "text", Required: false},
					{Name: "related", Type: "text", Required: false},
					{Name: "effort", Type: "select", Required: false, Options: constants.ValidEffortLevels},
					{Name: "due_date", Type: "date", Required: false},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "feature"},
				},
				Category: "Feature",
				Active:   true,
			},
		},
		// Technical Debt Template (old version, Turkish only, no English translation needed for legacy)
		{
			BaseTemplateID: "technical-debt",
			TR: &GorevTemplate{
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
## 🏷️ Tags: {{tags}}`,
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
					{Name: "tags", Type: "text", Required: false, Default: "teknik-borç,refaktöring"},
				},
				Category: "Teknik",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Technical Debt",
				Definition:   "Template for refactoring or technical improvements",
				Alias:        "debt",
				DefaultTitle: "🔧 [{{alan}}] {{title}}",
				DescriptionTemplate: `## 🔧 Technical Debt Description
{{description}}

## 📍 Affected Area
**Area/Module:** {{alan}}
**Files:** {{dosyalar}}

## ❓ Why Is It Needed?
{{neden}}

## 📊 Current State Analysis
{{analiz}}

## 🎯 Proposed Solution
{{cozum}}

## ⚠️ Risks
{{riskler}}

## 📈 Expected Improvements
{{iyilestirmeler}}

## ⏱️ Estimated Time: {{sure}}
## 🏷️ Tags: {{tags}}`,
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
					{Name: "sure", Type: "select", Required: false, Options: []string{"1 day", "2-3 days", "1 week", "2+ weeks"}},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "tech-debt,refactoring"},
				},
				Category: "Technical",
				Active:   true,
			},
		},
		// Research Template
		{
			BaseTemplateID: "research-task",
			TR: &GorevTemplate{
				Name:         "Araştırma Görevi",
				Definition:   "Teknoloji veya çözüm araştırması için template",
				Alias:        "research",
				DefaultTitle: "🔍 {{topic}} Araştırması",
				DescriptionTemplate: `## 🔍 Araştırma Konusu
{{topic}}

## 🎯 Araştırma Amacı
{{purpose}}

## ❓ Cevaplanması Gereken Sorular
{{questions}}

## 📚 Araştırılacak Kaynaklar
{{sources}}

## 🔄 Alternatifler
{{alternatives}}

## ⚖️ Değerlendirme Kriterleri
{{criteria}}

## 📅 Bitiş Tarihi: {{due_date}}
## 🏷️ Tags: {{tags}}`,
				Fields: []TemplateAlan{
					{Name: "topic", Type: "text", Required: true},
					{Name: "purpose", Type: "text", Required: true},
					{Name: "questions", Type: "text", Required: true},
					{Name: "sources", Type: "text", Required: false},
					{Name: "alternatives", Type: "text", Required: false},
					{Name: "criteria", Type: "text", Required: true},
					{Name: "due_date", Type: "date", Required: false},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "araştırma"},
				},
				Category: "Araştırma",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Research Task",
				Definition:   "Template for technology or solution research",
				Alias:        "research",
				DefaultTitle: "🔍 {{topic}} Research",
				DescriptionTemplate: `## 🔍 Research Topic
{{topic}}

## 🎯 Research Purpose
{{purpose}}

## ❓ Questions to Answer
{{questions}}

## 📚 Sources to Research
{{sources}}

## 🔄 Alternatives
{{alternatives}}

## ⚖️ Evaluation Criteria
{{criteria}}

## 📅 Due Date: {{due_date}}
## 🏷️ Tags: {{tags}}`,
				Fields: []TemplateAlan{
					{Name: "topic", Type: "text", Required: true},
					{Name: "purpose", Type: "text", Required: true},
					{Name: "questions", Type: "text", Required: true},
					{Name: "sources", Type: "text", Required: false},
					{Name: "alternatives", Type: "text", Required: false},
					{Name: "criteria", Type: "text", Required: true},
					{Name: "due_date", Type: "date", Required: false},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityMedium, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "research"},
				},
				Category: "Research",
				Active:   true,
			},
		},
		// Bug Report v2 Template (Enhanced)
		{
			BaseTemplateID: "bug-report-v2",
			TR: &GorevTemplate{
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
					{Name: "tags", Type: "text", Required: false, Default: "bug,production"},
				},
				Category: "Bug",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Bug Report v2",
				Definition:   "Enhanced bug report - detailed steps and environment info",
				Alias:        "bug2",
				DefaultTitle: "🐛 [{{severity}}] {{modul}}: {{title}}",
				DescriptionTemplate: `## 🐛 Bug Summary
{{description}}

## 🔄 Steps to Reproduce
{{steps_to_reproduce}}

## ✅ Expected Behavior
{{expected_behavior}}

## ❌ Actual Behavior
{{actual_behavior}}

## 💻 Environment Info
- **Operating System:** {{os_version}}
- **Browser/Client:** {{client_info}}
- **Server Version:** {{server_version}}
- **Database:** {{db_info}}

## 🚨 Bug Severity
**Severity:** {{severity}}
**Affected Users:** {{affected_users}}

## 📸 Attachments
{{attachments}}

## 🔧 Workaround
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
					{Name: "tags", Type: "text", Required: false, Default: "bug,production"},
				},
				Category: "Bug",
				Active:   true,
			},
		},
		// Spike Research Template
		{
			BaseTemplateID: "spike-research",
			TR: &GorevTemplate{
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

## 🏷️ Tags: {{tags}}`,
				Fields: []TemplateAlan{
					{Name: "research_question", Type: "text", Required: true},
					{Name: "success_criteria", Type: "text", Required: true},
					{Name: "time_box", Type: "select", Required: true, Options: []string{"4 saat", "1 gün", "2 gün", "3 gün", "1 hafta"}},
					{Name: "decision_deadline", Type: "date", Required: true},
					{Name: "research_plan", Type: "text", Required: true},
					{Name: "expected_outputs", Type: "text", Required: true},
					{Name: "risks_assumptions", Type: "text", Required: false},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityHigh, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "spike,research,poc"},
				},
				Category: "Araştırma",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Spike Research",
				Definition:   "Time-boxed technical research and proof-of-concept work",
				Alias:        "spike",
				DefaultTitle: "🔬 [SPIKE] {{research_question}}",
				DescriptionTemplate: `## 🔬 Research Question
{{research_question}}

## 🎯 Success Criteria
{{success_criteria}}

## ⏰ Time Box
**Maximum Duration:** {{time_box}}
**Decision Deadline:** {{decision_deadline}}

## 🔍 Research Plan
{{research_plan}}

## 📊 Expected Outputs
{{expected_outputs}}

## ⚡ Risks and Assumptions
{{risks_assumptions}}

## 🏷️ Tags: {{tags}}`,
				Fields: []TemplateAlan{
					{Name: "research_question", Type: "text", Required: true},
					{Name: "success_criteria", Type: "text", Required: true},
					{Name: "time_box", Type: "select", Required: true, Options: []string{"4 hours", "1 day", "2 days", "3 days", "1 week"}},
					{Name: "decision_deadline", Type: "date", Required: true},
					{Name: "research_plan", Type: "text", Required: true},
					{Name: "expected_outputs", Type: "text", Required: true},
					{Name: "risks_assumptions", Type: "text", Required: false},
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityHigh, Options: constants.GetValidPriorities()},
					{Name: "tags", Type: "text", Required: false, Default: "spike,research,poc"},
				},
				Category: "Research",
				Active:   true,
			},
		},
		// Performance Issue Template
		{
			BaseTemplateID: "performance-issue",
			TR: &GorevTemplate{
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

## 🏷️ Tags: {{tags}}`,
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
					{Name: "tags", Type: "text", Required: false, Default: "performance,optimization"},
				},
				Category: "Teknik",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Performance Issue",
				Definition:   "Performance problems and optimization tasks",
				Alias:        "performance",
				DefaultTitle: "⚡ [PERF] {{metric_affected}}: {{title}}",
				DescriptionTemplate: `## ⚡ Performance Issue
{{description}}

## 📊 Affected Metric
**Metric:** {{metric_affected}}
**Current Value:** {{current_value}}
**Target Value:** {{target_value}}
**Acceptable Value:** {{acceptable_value}}

## 📏 Measurement Method
{{measurement_method}}

## 👥 User Impact
{{user_impact}}

## 🔍 Root Cause Analysis
{{root_cause}}

## 💡 Proposed Solutions
{{proposed_solutions}}

## ⚠️ Trade-offs
{{tradeoffs}}

## 🏷️ Tags: {{tags}}`,
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
					{Name: "tags", Type: "text", Required: false, Default: "performance,optimization"},
				},
				Category: "Technical",
				Active:   true,
			},
		},
		// Security Fix Template
		{
			BaseTemplateID: "security-fix",
			TR: &GorevTemplate{
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

## 🏷️ Tags: {{tags}}`,
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
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityHigh, Options: []string{constants.PriorityHigh}},
					{Name: "tags", Type: "text", Required: false, Default: "security,vulnerability"},
				},
				Category: "Güvenlik",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Security Fix",
				Definition:   "Special template for security vulnerabilities and fixes",
				Alias:        "security",
				DefaultTitle: "🔒 [SEC-{{severity}}] {{vulnerability_type}}: {{title}}",
				DescriptionTemplate: `## 🔒 Security Vulnerability
{{description}}

## 🎯 Vulnerability Type
**Category:** {{vulnerability_type}}
**CVSS Score:** {{cvss_score}}
**Severity:** {{severity}}

## 🔍 Affected Components
{{affected_components}}

## 💥 Potential Impact
{{potential_impact}}

## 🛡️ Mitigation Steps
{{mitigation_steps}}

## ✅ Testing Requirements
{{testing_requirements}}

## 📋 Security Checklist
- [ ] Security testing completed
- [ ] Penetration testing required?
- [ ] Security review completed
- [ ] Documentation updated

## 🚨 Disclosure Timeline
{{disclosure_timeline}}

## 🏷️ Tags: {{tags}}`,
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
					{Name: "priority", Type: "select", Required: true, Default: constants.PriorityHigh, Options: []string{constants.PriorityHigh}},
					{Name: "tags", Type: "text", Required: false, Default: "security,vulnerability"},
				},
				Category: "Security",
				Active:   true,
			},
		},
		// Refactoring Template
		{
			BaseTemplateID: "refactoring",
			TR: &GorevTemplate{
				Name:         "Yeniden Düzenleme",
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

## 🏷️ Tags: {{tags}}`,
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
					{Name: "tags", Type: "text", Required: false, Default: "refactoring,code-quality"},
				},
				Category: "Teknik",
				Active:   true,
			},
			EN: &GorevTemplate{
				Name:         "Refactoring",
				Definition:   "Code quality and architectural improvements",
				Alias:        "refactor",
				DefaultTitle: "♻️ [REFACTOR] {{code_smell}}: {{title}}",
				DescriptionTemplate: `## ♻️ Refactoring Summary
{{description}}

## 🦨 Code Smell Type
{{code_smell_type}}

## 📁 Affected Files
{{affected_files}}

## 🎯 Refactoring Strategy
{{refactoring_strategy}}

## ✅ Success Criteria
- [ ] All tests passing
- [ ] Code coverage not decreased
- [ ] Performance not affected
- [ ] API compatibility maintained

## ⚠️ Risk Assessment
**Risk Level:** {{risk_level}}
**Impact Scope:** {{impact_scope}}

## 🔄 Rollback Plan
{{rollback_plan}}

## 📊 Metrics
- **Current Cyclomatic Complexity:** {{current_complexity}}
- **Target Complexity:** {{target_complexity}}
- **Current Code Coverage:** {{current_coverage}}

## 🏷️ Tags: {{tags}}`,
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
					{Name: "tags", Type: "text", Required: false, Default: "refactoring,code-quality"},
				},
				Category: "Technical",
				Active:   true,
			},
		},
	}

	// Create templates for each language
	for _, group := range templateGroups {
		// Set base_template_id for TR version
		group.TR.BaseTemplateID = &group.BaseTemplateID
		group.TR.LanguageCode = "tr"
		group.TR.ID = uuid.New().String()

		// Check if Turkish version exists
		ctxTR := i18n.WithLanguage(ctx, "tr")
		existingTR, err := vy.TemplateAliasIleGetir(ctxTR, group.TR.Alias)
		if err != nil || existingTR == nil {
			// Create Turkish version
			if err := vy.TemplateOlustur(ctxTR, group.TR); err != nil {
				return fmt.Errorf(i18n.T("error.defaultTemplateCreateFailed", map[string]interface{}{"Template": group.TR.Name, "Error": err}))
			}
		}

		// Create English version if defined
		if group.EN != nil {
			group.EN.BaseTemplateID = &group.BaseTemplateID
			group.EN.LanguageCode = "en"
			group.EN.ID = uuid.New().String()

			ctxEN := i18n.WithLanguage(ctx, "en")
			existingEN, err := vy.TemplateAliasIleGetir(ctxEN, group.EN.Alias)
			if err != nil || existingEN == nil {
				// Create English version
				if err := vy.TemplateOlustur(ctxEN, group.EN); err != nil {
					return fmt.Errorf(i18n.T("error.defaultTemplateCreateFailed", map[string]interface{}{"Template": group.EN.Name, "Error": err}))
				}
			}
		}
	}

	return nil
}
