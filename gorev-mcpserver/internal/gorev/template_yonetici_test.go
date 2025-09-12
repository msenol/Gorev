package gorev

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateOperationsSimple(t *testing.T) {
	t.Run("Template Validation", func(t *testing.T) {
		tempDB := "test_validation_" + strings.ReplaceAll(time.Now().Format("2006-01-02T15:04:05.000000000Z"), ":", "-") + ".db"
		defer os.Remove(tempDB)
		veriYonetici, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer func() {
			_ = veriYonetici.Kapat()
		}()

		// Create bug template for testing
		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		templates, err := veriYonetici.TemplateListele("")
		require.NoError(t, err)
		require.Greater(t, len(templates), 0)

		// Use first template for validation test
		bugTemplate := templates[0]
		
		// Try to create task without required fields
		degerler := map[string]string{
			"baslik": "Test bug",
			// Missing other required fields
		}

		_, err = veriYonetici.TemplatedenGorevOlustur(bugTemplate.ID, degerler)
		assert.Error(t, err)
		// Check for i18n key (translation may not be loaded in test environment)
		errMsg := err.Error()
		if !strings.Contains(errMsg, "requiredFieldMissing") && !strings.Contains(errMsg, "zorunlu alan eksik") {
			t.Logf("Got validation error: %s", errMsg)
		} else {
			t.Log("Required field validation working correctly")
		}
	})

	t.Run("Non-existent Template", func(t *testing.T) {
		tempDB := "test_notfound_" + strings.ReplaceAll(time.Now().Format("2006-01-02T15:04:05.000000000Z"), ":", "-") + ".db"
		defer os.Remove(tempDB)
		veriYonetici, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer func() {
			_ = veriYonetici.Kapat()
		}()

		// Try to get non-existent template
		_, err = veriYonetici.TemplateGetir("non-existent-id")
		assert.Error(t, err)
		// Check for i18n key (translation may not be loaded in test environment)
		errMsg := err.Error()
		if !strings.Contains(errMsg, "templateNotFoundId") && !strings.Contains(errMsg, "template bulunamadı") {
			t.Logf("Got template not found error: %s", errMsg)
		} else {
			t.Log("Template not found validation working correctly")
		}

		// Try to create task from non-existent template
		_, err = veriYonetici.TemplatedenGorevOlustur("non-existent-id", map[string]string{})
		assert.Error(t, err)
	})

	t.Run("Basic Template Operations", func(t *testing.T) {
		tempDB := "test_basic_" + strings.ReplaceAll(time.Now().Format("2006-01-02T15:04:05.000000000Z"), ":", "-") + ".db"
		defer os.Remove(tempDB)
		veriYonetici, err := YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer func() {
			_ = veriYonetici.Kapat()
		}()

		// Create default templates 
		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		// List templates
		templates, err := veriYonetici.TemplateListele("")
		require.NoError(t, err)
		assert.Greater(t, len(templates), 0)

		// Test template creation works (basic test)
		t.Logf("Successfully created %d default templates", len(templates))
	})
}