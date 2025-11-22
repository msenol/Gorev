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

	// Görev oluştur
	gorev := &Gorev{
		Title:       baslik,
		Description: aciklama,
		Priority:    oncelik,
		Status:      constants.TaskStatusPending,
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
			EN: nil, // No English translation for legacy templates
		},
		// Research Template (old version, Turkish only)
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
			EN: nil, // No English translation for legacy templates
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
