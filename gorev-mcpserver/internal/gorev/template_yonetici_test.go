package gorev

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateOperations(t *testing.T) {
	// Create temporary database for tests
	tempDB := "test_template_ops.db"
	defer os.Remove(tempDB)

	veriYonetici, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
	require.NoError(t, err)

	t.Run("Create and Retrieve Template", func(t *testing.T) {
		// Create a custom template
		template := &GorevTemplate{
			Isim:             "Test Template",
			Tanim:            "A template for testing",
			VarsayilanBaslik: "Test: {{title}}",
			AciklamaTemplate: "Description: {{description}}\nPriority: {{priority}}",
			Alanlar: []TemplateAlan{
				{Isim: "title", Tip: "text", Zorunlu: true},
				{Isim: "description", Tip: "text", Zorunlu: true},
				{Isim: "priority", Tip: "select", Zorunlu: true, Varsayilan: "medium", Secenekler: []string{"low", "medium", "high"}},
			},
			OrnekDegerler: map[string]string{
				"title":       "Example Title",
				"description": "Example Description",
				"priority":    "high",
			},
			Kategori: "Test",
			Aktif:    true,
		}

		// Create template
		err := veriYonetici.TemplateOlustur(template)
		require.NoError(t, err)
		assert.NotEmpty(t, template.ID)

		// Retrieve template
		retrieved, err := veriYonetici.TemplateGetir(template.ID)
		require.NoError(t, err)
		assert.Equal(t, template.Isim, retrieved.Isim)
		assert.Equal(t, template.Tanim, retrieved.Tanim)
		assert.Equal(t, template.Kategori, retrieved.Kategori)
		assert.Len(t, retrieved.Alanlar, 3)
	})

	t.Run("List Templates by Category", func(t *testing.T) {
		// Initialize default templates
		err := veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		// List all templates
		allTemplates, err := veriYonetici.TemplateListele("")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allTemplates), 9) // At least 9 default templates

		// List by category
		teknikTemplates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(teknikTemplates), 4) // Bug Report, Technical Debt, Performance, Refactoring

		// Verify all templates in Teknik category
		for _, tmpl := range teknikTemplates {
			assert.Equal(t, "Teknik", tmpl.Kategori)
		}
	})

	t.Run("Create Task from Template with Defaults", func(t *testing.T) {
		// Create a template with default values
		template := &GorevTemplate{
			Isim:             "Default Test Template",
			Tanim:            "Template with default values",
			VarsayilanBaslik: "{{type}}: {{subject}}",
			AciklamaTemplate: "Type: {{type}}\nSubject: {{subject}}\nTags: {{tags}}",
			Alanlar: []TemplateAlan{
				{Isim: "type", Tip: "text", Zorunlu: true},
				{Isim: "subject", Tip: "text", Zorunlu: true},
				{Isim: "tags", Tip: "text", Zorunlu: false, Varsayilan: "test,automated"},
			},
			Kategori: "Test",
			Aktif:    true,
		}

		err := veriYonetici.TemplateOlustur(template)
		require.NoError(t, err)

		// Create a test project
		proje := &Proje{
			ID:    "test-project",
			Isim:  "Test Project",
			Tanim: "Test project description",
		}
		err = veriYonetici.ProjeKaydet(proje)
		require.NoError(t, err)

		// Set as active project
		err = veriYonetici.AktifProjeAyarla(proje.ID)
		require.NoError(t, err)

		// Create task from template
		degerler := map[string]string{
			"type":      "Bug",
			"subject":   "Login issue",
			"tags":      "bug,urgent", // Provide tags value for template
			"oncelik":   "yuksek",
			"etiketler": "bug,urgent", // Also set as task tags
		}

		gorev, err := veriYonetici.TemplatedenGorevOlustur(template.ID, degerler)
		require.NoError(t, err)
		assert.Equal(t, "Bug: Login issue", gorev.Baslik)
		assert.Contains(t, gorev.Aciklama, "Type: Bug")
		assert.Contains(t, gorev.Aciklama, "Subject: Login issue")
		assert.Contains(t, gorev.Aciklama, "Tags: bug,urgent")
		assert.Equal(t, "yuksek", gorev.Oncelik)
		assert.Len(t, gorev.Etiketler, 2)
	})

	t.Run("Template Validation", func(t *testing.T) {
		// Get a default template
		templates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)
		require.NotEmpty(t, templates)

		var bugTemplate *GorevTemplate
		for _, tmpl := range templates {
			if tmpl.Isim == "Bug Raporu" {
				bugTemplate = tmpl
				break
			}
		}
		require.NotNil(t, bugTemplate)

		// Try to create task without required fields
		degerler := map[string]string{
			"baslik": "Test bug",
			// Missing other required fields
		}

		_, err = veriYonetici.TemplatedenGorevOlustur(bugTemplate.ID, degerler)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "zorunlu alan eksik")
	})

	t.Run("Non-existent Template", func(t *testing.T) {
		// Try to get non-existent template
		_, err := veriYonetici.TemplateGetir("non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template bulunamadı")

		// Try to create task from non-existent template
		_, err = veriYonetici.TemplatedenGorevOlustur("non-existent-id", map[string]string{})
		assert.Error(t, err)
	})

	t.Run("Template with All Field Types", func(t *testing.T) {
		// Create template with various field types
		template := &GorevTemplate{
			Isim:             "All Fields Template",
			Tanim:            "Template with all field types",
			VarsayilanBaslik: "{{title}}",
			AciklamaTemplate: "Text: {{text_field}}\nNumber: {{number_field}}\nDate: {{date_field}}\nSelect: {{select_field}}",
			Alanlar: []TemplateAlan{
				{Isim: "title", Tip: "text", Zorunlu: true},
				{Isim: "text_field", Tip: "text", Zorunlu: false},
				{Isim: "number_field", Tip: "number", Zorunlu: false},
				{Isim: "date_field", Tip: "date", Zorunlu: false},
				{Isim: "select_field", Tip: "select", Zorunlu: false, Secenekler: []string{"option1", "option2", "option3"}},
			},
			Kategori: "Test",
			Aktif:    true,
		}

		err := veriYonetici.TemplateOlustur(template)
		require.NoError(t, err)

		// Create a test project
		proje := &Proje{
			ID:    "test-project-all-fields",
			Isim:  "Test Project All Fields",
			Tanim: "Test project for all field types",
		}
		err = veriYonetici.ProjeKaydet(proje)
		require.NoError(t, err)

		// Set as active project
		err = veriYonetici.AktifProjeAyarla(proje.ID)
		require.NoError(t, err)

		// Create task with all field types
		degerler := map[string]string{
			"title":        "Test Task",
			"text_field":   "Some text",
			"number_field": "42",
			"date_field":   "2025-12-31",
			"select_field": "option2",
			"oncelik":      "orta",
			"son_tarih":    "2025-12-31",
		}

		gorev, err := veriYonetici.TemplatedenGorevOlustur(template.ID, degerler)
		require.NoError(t, err)
		assert.Equal(t, "Test Task", gorev.Baslik)
		assert.Contains(t, gorev.Aciklama, "Text: Some text")
		assert.Contains(t, gorev.Aciklama, "Number: 42")
		assert.Contains(t, gorev.Aciklama, "Date: 2025-12-31")
		assert.Contains(t, gorev.Aciklama, "Select: option2")
		assert.NotNil(t, gorev.SonTarih)
	})
}

func TestDefaultTemplates(t *testing.T) {
	// Create temporary database
	tempDB := "test_default_templates.db"
	defer os.Remove(tempDB)

	veriYonetici, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
	require.NoError(t, err)

	// Initialize default templates
	err = veriYonetici.VarsayilanTemplateleriOlustur()
	require.NoError(t, err)

	// Verify all default templates exist
	templates, err := veriYonetici.TemplateListele("")
	require.NoError(t, err)
	assert.Len(t, templates, 9)

	// Check each template
	templateNames := make(map[string]bool)
	for _, tmpl := range templates {
		templateNames[tmpl.Isim] = true
	}

	assert.True(t, templateNames["Bug Raporu"])
	assert.True(t, templateNames["Özellik İsteği"])
	assert.True(t, templateNames["Teknik Borç"])
	assert.True(t, templateNames["Araştırma Görevi"])
	assert.True(t, templateNames["Bug Raporu v2"])
	assert.True(t, templateNames["Spike Araştırma"])
	assert.True(t, templateNames["Performans Sorunu"])
	assert.True(t, templateNames["Güvenlik Düzeltmesi"])
	assert.True(t, templateNames["Refactoring"])

	// Verify template categories
	categories := make(map[string]int)
	for _, tmpl := range templates {
		categories[tmpl.Kategori]++
	}

	assert.Equal(t, 4, categories["Teknik"])    // Bug Raporu, Teknik Borç, Performans Sorunu, Refactoring
	assert.Equal(t, 1, categories["Özellik"])   // Özellik İsteği
	assert.Equal(t, 2, categories["Araştırma"]) // Araştırma Görevi, Spike Araştırma
	assert.Equal(t, 1, categories["Bug"])       // Bug Raporu v2
	assert.Equal(t, 1, categories["Güvenlik"])  // Güvenlik Düzeltmesi

	// Test creating duplicate templates (should not error due to UNIQUE constraint handling)
	err = veriYonetici.VarsayilanTemplateleriOlustur()
	assert.NoError(t, err) // Should handle UNIQUE constraint gracefully

	// Verify no duplicates were created
	templates, err = veriYonetici.TemplateListele("")
	require.NoError(t, err)
	assert.Len(t, templates, 9) // Still only 9 templates
}
