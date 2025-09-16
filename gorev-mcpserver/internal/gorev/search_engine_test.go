package gorev

import (
	"testing"

	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchEngine_Initialize(t *testing.T) {
	// Initialize i18n
	err := i18n.Initialize("tr")
	require.NoError(t, err)

	veriYonetici, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	db, err := veriYonetici.GetDB()
	require.NoError(t, err)

	searchEngine := NewSearchEngine(veriYonetici, db)
	require.NotNil(t, searchEngine)

	// Test initialization
	err = searchEngine.Initialize()
	assert.NoError(t, err)
}

func TestSearchEngine_PerformSearch(t *testing.T) {
	// Initialize i18n
	err := i18n.Initialize("tr")
	require.NoError(t, err)

	veriYonetici, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	db, err := veriYonetici.GetDB()
	require.NoError(t, err)

	searchEngine := NewSearchEngine(veriYonetici, db)
	require.NoError(t, searchEngine.Initialize())

	// Create test project
	proje, err := veriYonetici.ProjeOlustur("Test Proje", "Test açıklama")
	require.NoError(t, err)

	// Create test tasks
	task1, err := veriYonetici.GorevOlusturBasit("Search Test Task", "Test açıklama database", proje.ID, "yuksek", "2024-12-31", "", "")
	require.NoError(t, err)

	_, err = veriYonetici.GorevOlusturBasit("Another Task", "Different content here", proje.ID, "orta", "2024-12-31", "", "")
	require.NoError(t, err)

	// Test basic search
	results, err := searchEngine.PerformSearch("Search", SearchFilters{})
	assert.NoError(t, err)
	assert.NotNil(t, results)

	// Should find the task with "Search" in title
	assert.GreaterOrEqual(t, len(results.Results), 1)

	// Verify task1 is in results
	found := false
	for _, result := range results.Results {
		if result.Task.ID == task1.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Should find task1 in search results")
}