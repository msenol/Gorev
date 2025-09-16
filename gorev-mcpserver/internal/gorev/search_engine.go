package gorev

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/msenol/gorev/internal/i18n"
)

// SearchEngine handles advanced search functionality with FTS5 and fuzzy matching
type SearchEngine struct {
	veriYonetici VeriYoneticiInterface
	db           *sql.DB
	nlpProcessor *NLPProcessor
}

// SearchOptions contains options for advanced search
type SearchOptions struct {
	Query            string                 `json:"query"`
	Filters          map[string]interface{} `json:"filters"`
	UseFuzzySearch   bool                   `json:"use_fuzzy_search"`
	FuzzyThreshold   float64                `json:"fuzzy_threshold"`
	MaxResults       int                    `json:"max_results"`
	SortBy           string                 `json:"sort_by"`
	SortDirection    string                 `json:"sort_direction"`
	IncludeCompleted bool                   `json:"include_completed"`
	SearchFields     []string               `json:"search_fields"`
}

// SearchResult contains search result with relevance scoring
type SearchResult struct {
	Task           *Gorev   `json:"task"`
	RelevanceScore float64  `json:"relevance_score"`
	MatchType      string   `json:"match_type"` // "exact", "fts", "fuzzy"
	MatchedFields  []string `json:"matched_fields"`
}

// SearchResponse contains search results with metadata
type SearchResponse struct {
	Results     []SearchResult `json:"results"`
	TotalCount  int            `json:"total_count"`
	QueryTime   time.Duration  `json:"query_time"`
	UsedFuzzy   bool           `json:"used_fuzzy"`
	Suggestions []string       `json:"suggestions"`
}

// FilterProfile represents a saved search/filter configuration
type FilterProfile struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Filters     map[string]interface{} `json:"filters"`
	SearchQuery string                 `json:"search_query"`
	IsDefault   bool                   `json:"is_default"`
	CreatedAt   time.Time              `json:"created_at"`
	LastUsedAt  *time.Time             `json:"last_used_at"`
	UseCount    int                    `json:"use_count"`
}

// SearchHistoryEntry represents a search history record
type SearchHistoryEntry struct {
	ID              int       `json:"id"`
	Query           string    `json:"query"`
	Filters         string    `json:"filters"`
	ResultCount     int       `json:"result_count"`
	ExecutionTimeMs int       `json:"execution_time_ms"`
	CreatedAt       time.Time `json:"created_at"`
}

// NewSearchEngine creates a new search engine instance
func NewSearchEngine(vy VeriYoneticiInterface, db *sql.DB) *SearchEngine {
	return &SearchEngine{
		veriYonetici: vy,
		db:           db,
		nlpProcessor: NewNLPProcessor(),
	}
}

// Search performs advanced search with FTS5 and fuzzy matching support
func (se *SearchEngine) Search(options SearchOptions) (*SearchResponse, error) {
	startTime := time.Now()

	// Set default options
	if options.MaxResults <= 0 {
		options.MaxResults = 50
	}
	if options.FuzzyThreshold <= 0 {
		options.FuzzyThreshold = 0.6
	}
	if len(options.SearchFields) == 0 {
		options.SearchFields = []string{"baslik", "aciklama", "etiketler", "proje_adi"}
	}

	var results []SearchResult
	usedFuzzy := false
	originalQuery := options.Query

	// Process query with NLP for enhanced filtering and suggestions
	var nlpIntent *QueryIntent
	if strings.TrimSpace(options.Query) != "" {
		if intent, err := se.nlpProcessor.ProcessQuery(options.Query); err == nil {
			nlpIntent = intent

			// Merge NLP-detected filters with provided filters
			for key, value := range intent.Filters {
				if _, exists := options.Filters[key]; !exists {
					options.Filters[key] = value
				}
			}

			// Extract clean search terms (remove filter expressions)
			cleanQuery := se.extractCleanSearchTerms(options.Query, intent)
			if cleanQuery != "" {
				options.Query = cleanQuery
			}
		}
	}

	// Try FTS5 search first if query is provided
	if strings.TrimSpace(options.Query) != "" {
		ftsResults, err := se.performFTSSearch(options)
		if err != nil {
			log.Printf("%s", i18n.T("error.ftsSearchFailed", map[string]interface{}{"Error": err}))
		} else {
			results = append(results, ftsResults...)
		}

		// If FTS didn't return enough results and fuzzy is enabled, try fuzzy search
		if len(results) < options.MaxResults/2 && options.UseFuzzySearch {
			fuzzyResults, err := se.performFuzzySearch(options)
			if err != nil {
				log.Printf("%s", i18n.T("error.fuzzySearchFailed", map[string]interface{}{"Error": err}))
			} else {
				results = append(results, fuzzyResults...)
				usedFuzzy = true
			}
		}
	}

	// Apply filters
	if len(options.Filters) > 0 {
		if len(results) == 0 {
			// No text search, just filter all tasks
			allTasks, err := se.veriYonetici.GorevListele(map[string]interface{}{
				"limit": 1000,
			})
			if err != nil {
				return nil, fmt.Errorf(i18n.T("error.tasksRetrieveFailed", map[string]interface{}{"Error": err}))
			}

			for _, task := range allTasks {
				results = append(results, SearchResult{
					Task:           task,
					RelevanceScore: 1.0,
					MatchType:      "filter",
					MatchedFields:  []string{},
				})
			}
		}

		results = se.applyFilters(results, options.Filters)
	}

	// Remove duplicates and sort by relevance
	results = se.removeDuplicates(results)
	se.sortResults(results, options.SortBy, options.SortDirection)

	// Limit results
	if len(results) > options.MaxResults {
		results = results[:options.MaxResults]
	}

	// Record search history
	se.recordSearchHistory(options, len(results), time.Since(startTime))

	// Generate suggestions based on NLP analysis and search history
	suggestions := se.generateSuggestions(originalQuery, nlpIntent, results)

	return &SearchResponse{
		Results:     results,
		TotalCount:  len(results),
		QueryTime:   time.Since(startTime),
		UsedFuzzy:   usedFuzzy,
		Suggestions: suggestions,
	}, nil
}

// performFTSSearch executes FTS5 full-text search
func (se *SearchEngine) performFTSSearch(options SearchOptions) ([]SearchResult, error) {
	query := se.prepareFTSQuery(options.Query)
	if query == "" {
		return nil, nil
	}

	sqlQuery := `
		SELECT g.id, g.baslik, g.aciklama, g.durum, g.oncelik, g.son_tarih,
		       g.olusturma_tarih, g.guncelleme_tarih, g.proje_id, g.parent_id,
		       COALESCE(p.ad, '') as proje_adi,
		       fts.rank
		FROM gorevler_fts fts
		JOIN gorevler g ON fts.id = g.id
		LEFT JOIN projeler p ON g.proje_id = p.id
		WHERE fts MATCH ?
		ORDER BY fts.rank
		LIMIT ?
	`

	rows, err := se.db.Query(sqlQuery, query, options.MaxResults)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.ftsSearchFailed", map[string]interface{}{"Error": err}))
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var task Gorev
		var projeAdi string
		var rank float64

		err := rows.Scan(
			&task.ID, &task.Baslik, &task.Aciklama, &task.Durum, &task.Oncelik,
			&task.SonTarih, &task.OlusturmaTarih, &task.GuncellemeTarih,
			&task.ProjeID, &task.ParentID, &projeAdi, &rank,
		)
		if err != nil {
			log.Printf("%s", i18n.T("error.scanResultFailed", map[string]interface{}{"Error": err}))
			continue
		}

		// Load tags - first get the task details which includes tags
		taskDetail, err := se.veriYonetici.GorevDetay(task.ID)
		if err == nil && taskDetail != nil {
			task.Etiketler = taskDetail.Etiketler
		}

		// Calculate relevance score based on FTS rank
		relevanceScore := se.calculateFTSRelevance(rank, options.Query, &task)

		// Determine matched fields
		matchedFields := se.getMatchedFields(options.Query, &task)

		results = append(results, SearchResult{
			Task:           &task,
			RelevanceScore: relevanceScore,
			MatchType:      "fts",
			MatchedFields:  matchedFields,
		})
	}

	return results, nil
}

// performFuzzySearch executes fuzzy string matching
func (se *SearchEngine) performFuzzySearch(options SearchOptions) ([]SearchResult, error) {
	// Get all tasks for fuzzy matching
	allTasks, err := se.veriYonetici.GorevListele(map[string]interface{}{
		"limit": 1000,
	})
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.tasksRetrieveFailed", map[string]interface{}{"Error": err}))
	}

	var results []SearchResult
	query := strings.ToLower(strings.TrimSpace(options.Query))

	for _, task := range allTasks {
		// Calculate fuzzy scores for different fields
		titleScore := se.calculateLevenshteinSimilarity(query, strings.ToLower(task.Baslik))
		descScore := se.calculateLevenshteinSimilarity(query, strings.ToLower(task.Aciklama))

		// Get tag names for fuzzy matching
		taskDetail, _ := se.veriYonetici.GorevDetay(task.ID)
		var tagNames []string
		if taskDetail != nil && taskDetail.Etiketler != nil {
			for _, etiket := range taskDetail.Etiketler {
				tagNames = append(tagNames, etiket.Isim)
			}
			task.Etiketler = taskDetail.Etiketler
		}
		tagText := strings.ToLower(strings.Join(tagNames, " "))
		tagScore := se.calculateLevenshteinSimilarity(query, tagText)

		// Take the highest score
		maxScore := titleScore
		matchedField := "baslik"

		if descScore > maxScore {
			maxScore = descScore
			matchedField = "aciklama"
		}
		if tagScore > maxScore {
			maxScore = tagScore
			matchedField = "etiketler"
		}

		// Only include if above threshold
		if maxScore >= options.FuzzyThreshold {
			// task.Etiketler already set above

			results = append(results, SearchResult{
				Task:           task,
				RelevanceScore: maxScore,
				MatchType:      "fuzzy",
				MatchedFields:  []string{matchedField},
			})
		}
	}

	return results, nil
}

// calculateLevenshteinSimilarity calculates similarity between two strings using Levenshtein distance
func (se *SearchEngine) calculateLevenshteinSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Convert to runes for proper UTF-8 handling
	r1 := []rune(s1)
	r2 := []rune(s2)

	distance := se.levenshteinDistance(r1, r2)
	maxLen := len(r1)
	if len(r2) > maxLen {
		maxLen = len(r2)
	}

	// Convert distance to similarity (0.0 to 1.0)
	similarity := 1.0 - float64(distance)/float64(maxLen)
	if similarity < 0 {
		similarity = 0
	}

	return similarity
}

// levenshteinDistance calculates the Levenshtein distance between two rune slices
func (se *SearchEngine) levenshteinDistance(s1, s2 []rune) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// prepareFTSQuery prepares the query for FTS5 search
func (se *SearchEngine) prepareFTSQuery(query string) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return ""
	}

	// Split into words and add wildcard matching
	words := strings.Fields(query)
	var ftsTerms []string

	for _, word := range words {
		// Remove special characters that might break FTS
		cleaned := strings.Trim(word, `"'*()[]{}`)
		if len(cleaned) > 0 {
			// Add wildcard for partial matching
			ftsTerms = append(ftsTerms, cleaned+"*")
		}
	}

	// Join with OR for broad matching
	return strings.Join(ftsTerms, " OR ")
}

// calculateFTSRelevance calculates relevance score for FTS results
func (se *SearchEngine) calculateFTSRelevance(rank float64, query string, task *Gorev) float64 {
	// FTS rank is negative (closer to 0 = better)
	baseScore := 1.0 / (1.0 + (-rank))

	// Boost score for exact matches in title
	queryLower := strings.ToLower(query)
	titleLower := strings.ToLower(task.Baslik)

	if strings.Contains(titleLower, queryLower) {
		baseScore *= 1.5
	}

	// Boost for high priority tasks
	if task.Oncelik == "yuksek" {
		baseScore *= 1.2
	}

	// Boost for in-progress tasks
	if task.Durum == "devam_ediyor" {
		baseScore *= 1.3
	}

	return baseScore
}

// getMatchedFields determines which fields matched the search query
func (se *SearchEngine) getMatchedFields(query string, task *Gorev) []string {
	var matched []string
	queryLower := strings.ToLower(query)

	if strings.Contains(strings.ToLower(task.Baslik), queryLower) {
		matched = append(matched, "baslik")
	}
	if strings.Contains(strings.ToLower(task.Aciklama), queryLower) {
		matched = append(matched, "aciklama")
	}

	// Check tags
	for _, tag := range task.Etiketler {
		if strings.Contains(strings.ToLower(tag.Isim), queryLower) {
			matched = append(matched, "etiketler")
			break
		}
	}

	return matched
}

// applyFilters applies search filters to results
func (se *SearchEngine) applyFilters(results []SearchResult, filters map[string]interface{}) []SearchResult {
	var filtered []SearchResult

	for _, result := range results {
		if se.matchesFilters(result.Task, filters) {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

// matchesFilters checks if a task matches the given filters
func (se *SearchEngine) matchesFilters(task *Gorev, filters map[string]interface{}) bool {
	for key, value := range filters {
		switch key {
		case "durum":
			if valueStr, ok := value.(string); ok && task.Durum != valueStr {
				return false
			}
		case "oncelik":
			if valueStr, ok := value.(string); ok && task.Oncelik != valueStr {
				return false
			}
		case "proje_id":
			if valueStr, ok := value.(string); ok {
				if task.ProjeID == "" && valueStr != "" {
					return false
				}
				if task.ProjeID != "" && task.ProjeID != valueStr {
					return false
				}
			}
		case "son_tarih":
			if valueStr, ok := value.(string); ok {
				if !se.matchesDateFilter(task.SonTarih, valueStr) {
					return false
				}
			}
		case "etiket":
			if valueStr, ok := value.(string); ok {
				found := false
				for _, tag := range task.Etiketler {
					if tag.Isim == valueStr {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
	}

	return true
}

// matchesDateFilter checks if a date matches the given filter
func (se *SearchEngine) matchesDateFilter(taskDate *time.Time, filter string) bool {
	if taskDate == nil {
		return filter == "no_date"
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch filter {
	case "today":
		return taskDate.Format("2006-01-02") == today.Format("2006-01-02")
	case "tomorrow":
		tomorrow := today.AddDate(0, 0, 1)
		return taskDate.Format("2006-01-02") == tomorrow.Format("2006-01-02")
	case "this_week":
		weekStart := today.AddDate(0, 0, -int(today.Weekday()))
		weekEnd := weekStart.AddDate(0, 0, 7)
		return taskDate.After(weekStart) && taskDate.Before(weekEnd)
	case "overdue":
		return taskDate.Before(today)
	case "upcoming":
		return taskDate.After(today)
	}

	return true
}

// sortResults sorts search results by the specified criteria
func (se *SearchEngine) sortResults(results []SearchResult, sortBy, direction string) {
	if sortBy == "" {
		sortBy = "relevance"
	}
	if direction == "" {
		direction = "desc"
	}

	sort.Slice(results, func(i, j int) bool {
		var less bool

		switch sortBy {
		case "relevance":
			less = results[i].RelevanceScore < results[j].RelevanceScore
		case "created":
			less = results[i].Task.OlusturmaTarih.Before(results[j].Task.OlusturmaTarih)
		case "updated":
			less = results[i].Task.GuncellemeTarih.Before(results[j].Task.GuncellemeTarih)
		case "due_date":
			if results[i].Task.SonTarih == nil && results[j].Task.SonTarih == nil {
				less = false
			} else if results[i].Task.SonTarih == nil {
				less = false
			} else if results[j].Task.SonTarih == nil {
				less = true
			} else {
				less = results[i].Task.SonTarih.Before(*results[j].Task.SonTarih)
			}
		case "priority":
			priOrder := map[string]int{
				"yuksek": 3,
				"orta":   2,
				"dusuk":  1,
			}
			less = priOrder[results[i].Task.Oncelik] < priOrder[results[j].Task.Oncelik]
		default:
			less = results[i].RelevanceScore < results[j].RelevanceScore
		}

		if direction == "desc" {
			return !less
		}
		return less
	})
}

// removeDuplicates removes duplicate tasks from search results
func (se *SearchEngine) removeDuplicates(results []SearchResult) []SearchResult {
	seen := make(map[string]bool)
	var unique []SearchResult

	for _, result := range results {
		if !seen[result.Task.ID] {
			seen[result.Task.ID] = true
			unique = append(unique, result)
		}
	}

	return unique
}

// recordSearchHistory records search queries for analytics
func (se *SearchEngine) recordSearchHistory(options SearchOptions, resultCount int, executionTime time.Duration) {
	if se.db == nil {
		return
	}

	filtersJSON, _ := json.Marshal(options.Filters)

	_, err := se.db.Exec(`
		INSERT INTO search_history (query, filters, result_count, execution_time_ms)
		VALUES (?, ?, ?, ?)
	`, options.Query, string(filtersJSON), resultCount, int(executionTime.Milliseconds()))

	if err != nil {
		log.Printf("%s", i18n.T("error.searchHistoryFailed", map[string]interface{}{"Error": err}))
	}
}

// generateSuggestions generates smart search suggestions based on NLP analysis and context
func (se *SearchEngine) generateSuggestions(originalQuery string, nlpIntent *QueryIntent, results []SearchResult) []string {
	var suggestions []string

	// If no results found, suggest alternative queries
	if len(results) == 0 && originalQuery != "" {
		suggestions = append(suggestions, se.generateNoResultsSuggestions(originalQuery)...)
	}

	// Add NLP-based filter suggestions
	if nlpIntent != nil {
		suggestions = append(suggestions, se.generateNLPSuggestions(nlpIntent)...)
	}

	// Add common search patterns based on existing tags and projects
	suggestions = append(suggestions, se.generateCommonPatternSuggestions()...)

	// Add time-based suggestions
	suggestions = append(suggestions, se.generateTimeBasedSuggestions()...)

	// Remove duplicates and limit
	suggestions = se.removeDuplicateStrings(suggestions)
	if len(suggestions) > 8 {
		suggestions = suggestions[:8]
	}

	return suggestions
}

// extractCleanSearchTerms removes filter expressions from query to get clean search terms
func (se *SearchEngine) extractCleanSearchTerms(query string, intent *QueryIntent) string {
	if intent == nil {
		return query
	}

	cleaned := query

	// Remove tag expressions (etiket:value, tag:value)
	tagRegex := regexp.MustCompile(`\b(etiket|tag):\w+\b`)
	cleaned = tagRegex.ReplaceAllString(cleaned, "")

	// Remove status expressions (durum:value, status:value)
	statusRegex := regexp.MustCompile(`\b(durum|status):\w+\b`)
	cleaned = statusRegex.ReplaceAllString(cleaned, "")

	// Remove priority expressions (öncelik:value, priority:value)
	priorityRegex := regexp.MustCompile(`\b(öncelik|priority):\w+\b`)
	cleaned = priorityRegex.ReplaceAllString(cleaned, "")

	// Clean up extra spaces
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(cleaned), " ")

	return cleaned
}

// generateNoResultsSuggestions provides alternative queries when no results found
func (se *SearchEngine) generateNoResultsSuggestions(query string) []string {
	var suggestions []string

	// Suggest fuzzy search if not already enabled
	suggestions = append(suggestions, fmt.Sprintf("Fuzzy search: %s", query))

	// Suggest broader terms
	words := strings.Fields(strings.ToLower(query))
	if len(words) > 1 {
		// Suggest searching for individual words
		for _, word := range words {
			if len(word) > 3 {
				suggestions = append(suggestions, word)
			}
		}
	}

	// Suggest common alternatives for Turkish/English
	alternatives := map[string][]string{
		"bug":      {"hata", "sorun", "problem"},
		"hata":     {"bug", "sorun", "problem"},
		"feature":  {"özellik", "yeni", "ekleme"},
		"özellik":  {"feature", "yeni", "ekleme"},
		"test":     {"test", "deneme", "kontrol"},
		"fix":      {"düzelt", "tamir", "çözüm"},
		"düzelt":   {"fix", "tamir", "çözüm"},
		"update":   {"güncelle", "revize", "değiştir"},
		"güncelle": {"update", "revize", "değiştir"},
	}

	for _, word := range words {
		if alts, exists := alternatives[word]; exists {
			suggestions = append(suggestions, alts...)
		}
	}

	return suggestions
}

// generateNLPSuggestions creates suggestions based on NLP analysis
func (se *SearchEngine) generateNLPSuggestions(intent *QueryIntent) []string {
	var suggestions []string

	// Suggest related filters based on detected intent
	if intent.Action == "list" {
		suggestions = append(suggestions, []string{
			"durum:devam_ediyor",
			"öncelik:yuksek",
			"son_tarih:bugün",
			"son_tarih:yarın",
		}...)
	}

	// Add temporal suggestions based on time range detection
	if intent.TimeRange != nil {
		if intent.TimeRange.Relative != "" {
			switch intent.TimeRange.Relative {
			case "bugün", "today":
				suggestions = append(suggestions, "son_tarih:bugün", "oluşturma:bugün")
			case "yarın", "tomorrow":
				suggestions = append(suggestions, "son_tarih:yarın")
			case "bu hafta", "this week":
				suggestions = append(suggestions, "son_tarih:bu_hafta")
			}
		}
	}

	return suggestions
}

// generateCommonPatternSuggestions suggests common search patterns
func (se *SearchEngine) generateCommonPatternSuggestions() []string {
	// These would ideally be based on actual tag/project analysis
	return []string{
		"etiket:bug",
		"etiket:feature",
		"etiket:urgent",
		"durum:beklemede",
		"durum:tamamlandi",
		"öncelik:yuksek",
		"bugün tamamlanacak",
		"gecikmis görevler",
	}
}

// generateTimeBasedSuggestions provides time-based search suggestions
func (se *SearchEngine) generateTimeBasedSuggestions() []string {
	return []string{
		"bugün",
		"yarın",
		"bu hafta",
		"geçen hafta",
		"bu ay",
		"geciken",
		"yaklaşan deadline",
	}
}

// removeDuplicateStrings removes duplicate strings from slice
func (se *SearchEngine) removeDuplicateStrings(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] && strings.TrimSpace(item) != "" {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}
