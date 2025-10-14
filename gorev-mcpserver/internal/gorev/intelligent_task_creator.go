package gorev

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// IntelligentTaskCreator provides advanced task creation capabilities
type IntelligentTaskCreator struct {
	veriYonetici     VeriYoneticiInterface
	suggestionEngine *SuggestionEngine
}

// NewIntelligentTaskCreator creates a new intelligent task creator
func NewIntelligentTaskCreator(vy VeriYoneticiInterface) *IntelligentTaskCreator {
	return &IntelligentTaskCreator{
		veriYonetici:     vy,
		suggestionEngine: NewSuggestionEngine(vy),
	}
}

// TaskCreationRequest represents an intelligent task creation request
type TaskCreationRequest struct {
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	AutoSplit       bool              `json:"auto_split,omitempty"`
	EstimateTime    bool              `json:"estimate_time,omitempty"`
	SmartPriority   bool              `json:"smart_priority,omitempty"`
	SuggestTemplate bool              `json:"suggest_template,omitempty"`
	Context         map[string]string `json:"context,omitempty"`
}

// TaskCreationResponse contains the created task and intelligent suggestions
type TaskCreationResponse struct {
	MainTask            *Gorev            `json:"main_task"`
	Subtasks            []*Gorev          `json:"subtasks,omitempty"`
	EstimatedHours      float64           `json:"estimated_hours,omitempty"`
	SuggestedPriority   string            `json:"suggested_priority,omitempty"`
	RecommendedTemplate string            `json:"recommended_template,omitempty"`
	SimilarTasks        []SimilarTaskInfo `json:"similar_tasks,omitempty"`
	Insights            []string          `json:"insights"`
	Confidence          TaskAnalysisScore `json:"confidence"`
	ExecutionTime       time.Duration     `json:"execution_time"`
}

// SimilarTaskInfo contains information about similar tasks
type SimilarTaskInfo struct {
	Task            *Gorev  `json:"task"`
	SimilarityScore float64 `json:"similarity_score"`
	Reason          string  `json:"reason"`
}

// TaskAnalysisScore contains confidence scores for different analyses
type TaskAnalysisScore struct {
	PriorityConfidence float64 `json:"priority_confidence"`
	TimeConfidence     float64 `json:"time_confidence"`
	TemplateConfidence float64 `json:"template_confidence"`
	SubtaskConfidence  float64 `json:"subtask_confidence"`
}

// CreateIntelligentTask creates a task with AI-enhanced features
func (itc *IntelligentTaskCreator) CreateIntelligentTask(ctx context.Context, request TaskCreationRequest) (*TaskCreationResponse, error) {
	startTime := time.Now()

	log.Printf("Creating intelligent task: title=%s, autoSplit=%t", request.Title, request.AutoSplit)

	response := &TaskCreationResponse{
		Insights: []string{},
	}

	// Determine smart priority
	suggestedPriority := constants.PriorityMedium // default
	priorityConfidence := constants.DefaultConfidence
	if request.SmartPriority {
		suggestedPriority, priorityConfidence = itc.determinePriority(request.Title, request.Description)
		response.SuggestedPriority = suggestedPriority
		response.Insights = append(response.Insights,
			fmt.Sprintf("Öncelik analizi: %s (güven: %.1f%%)", suggestedPriority, priorityConfidence*100))
	}

	// Estimate time based on historical data
	estimatedHours := 0.0
	timeConfidence := 0.0
	if request.EstimateTime {
		estimatedHours, timeConfidence = itc.estimateTaskTime(ctx, request.Title, request.Description)
		response.EstimatedHours = estimatedHours
		if estimatedHours > 0 {
			response.Insights = append(response.Insights,
				fmt.Sprintf("Tahmini süre: %.1f saat (güven: %.1f%%)", estimatedHours, timeConfidence*100))
		}
	}

	// Find similar tasks
	if len(request.Title) > 0 {
		similarTasks := itc.findSimilarTasks(ctx, request.Title, request.Description, 5)
		response.SimilarTasks = similarTasks
		if len(similarTasks) > 0 {
			response.Insights = append(response.Insights,
				fmt.Sprintf("%d benzer görev bulundu", len(similarTasks)))
		}
	}

	// Suggest template
	templateConfidence := 0.0
	if request.SuggestTemplate {
		recommendedTemplate, confidence := itc.suggestTemplate(request.Title, request.Description)
		if recommendedTemplate != "" {
			response.RecommendedTemplate = recommendedTemplate
			templateConfidence = confidence
			response.Insights = append(response.Insights,
				fmt.Sprintf("Önerilen template: %s (güven: %.1f%%)", recommendedTemplate, confidence*100))
		}
	}

	// Create main task
	mainTaskID, err := itc.veriYonetici.GorevOlustur(ctx, map[string]interface{}{
		"title":       request.Title,
		"description": request.Description,
		"priority":    suggestedPriority,
		"project_id":  "", // project will be set by caller if needed
		"due_date":    "", // due date will be set by caller if needed
	})
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.mainTaskCreateFailed", map[string]interface{}{"Error": err}))
	}

	// Get the created task for response
	mainTask, err := itc.veriYonetici.GorevGetir(ctx, mainTaskID)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.createdTaskFetchFailed", map[string]interface{}{"Error": err}))
	}

	response.MainTask = mainTask

	// Auto-split into subtasks if requested
	subtaskConfidence := 0.0
	if request.AutoSplit && len(request.Description) > 50 {
		subtasks, confidence := itc.generateSubtasks(mainTask.ID, request.Description)
		response.Subtasks = subtasks
		subtaskConfidence = confidence
		if len(subtasks) > 0 {
			response.Insights = append(response.Insights,
				fmt.Sprintf("%d alt görev otomatik oluşturuldu", len(subtasks)))
		}
	}

	// Set confidence scores
	response.Confidence = TaskAnalysisScore{
		PriorityConfidence: priorityConfidence,
		TimeConfidence:     timeConfidence,
		TemplateConfidence: templateConfidence,
		SubtaskConfidence:  subtaskConfidence,
	}

	response.ExecutionTime = time.Since(startTime)

	log.Printf("Intelligent task created: taskId=%s, subtasks=%d, insights=%d, duration=%v", mainTask.ID, len(response.Subtasks), len(response.Insights), response.ExecutionTime)

	return response, nil
}

// analyzeTaskContent performs basic content analysis
func (itc *IntelligentTaskCreator) analyzeTaskContent(title, description string) map[string]interface{} {
	content := strings.ToLower(title + " " + description)

	analysis := map[string]interface{}{
		"length":                len(content),
		"word_count":            len(strings.Fields(content)),
		"has_numbers":           regexp.MustCompile(`\d+`).MatchString(content),
		"urgency_keywords":      itc.containsUrgencyKeywords(content),
		"complexity_indicators": itc.detectComplexityIndicators(content),
	}

	return analysis
}

// determinePriority analyzes content to suggest task priority
func (itc *IntelligentTaskCreator) determinePriority(title, description string) (string, float64) {
	content := strings.ToLower(title + " " + description)

	// High priority indicators
	highPriorityKeywords := []string{
		"acil", "urgent", "kritik", "critical", "hemen", "immediately", "asap",
		"bug", "hata", "error", "çöktü", "crash", "down", "broken",
		"security", "güvenlik", "vulnerability", "zafiyet",
		"deadline", "son tarih", "due", "release",
	}

	// Medium priority indicators
	mediumPriorityKeywords := []string{
		"important", "önemli", "should", "need", "gerek", "lazım",
		"improvement", "iyileştirme", "enhance", "optimize",
		"feature", "özellik", "functionality", "fonksiyon",
	}

	// Low priority indicators
	lowPriorityKeywords := []string{
		"nice to have", "optional", "opsiyonel", "when time permits",
		"documentation", "dokuman", "cleanup", "temizlik",
		"refactor", "refaktör", "research", "araştırma",
	}

	highScore := itc.countKeywords(content, highPriorityKeywords)
	mediumScore := itc.countKeywords(content, mediumPriorityKeywords)
	lowScore := itc.countKeywords(content, lowPriorityKeywords)

	// Calculate priority based on keyword scores
	totalKeywords := highScore + mediumScore + lowScore

	if totalKeywords == 0 {
		return constants.PriorityMedium, constants.ConfidenceVeryLow // Default with low confidence
	}

	confidence := float64(totalKeywords) / constants.ConfidenceNormalizer // Normalize confidence
	if confidence > constants.ConfidenceVeryHigh {
		confidence = constants.ConfidenceVeryHigh
	}

	if highScore > 0 && highScore >= mediumScore && highScore >= lowScore {
		return constants.PriorityHigh, confidence + constants.ConfidenceWeightMedium
	} else if lowScore > 0 && lowScore > highScore && lowScore >= mediumScore {
		return constants.PriorityLow, confidence + constants.ConfidenceWeightLow
	} else {
		return constants.PriorityMedium, confidence + constants.ConfidenceWeightMinimal
	}
}

// estimateTaskTime estimates completion time based on similar tasks
func (itc *IntelligentTaskCreator) estimateTaskTime(ctx context.Context, title, description string) (float64, float64) {
	// Find similar completed tasks
	similarTasks := itc.findSimilarTasks(ctx, title, description, 10)

	if len(similarTasks) == 0 {
		// Fallback estimation based on content analysis
		return itc.estimateTimeFromContent(title, description), constants.ConfidenceVeryLow
	}

	// Calculate average time from similar tasks
	var totalHours float64
	var validSamples int

	for _, similar := range similarTasks {
		// Use estimated duration as base, since Gorev struct doesn't have ActualHours
		if estimatedDays := itc.extractDurationFromDescription(similar.Task.Description); estimatedDays > 0 {
			// Weight by similarity score (convert days to hours)
			totalHours += float64(estimatedDays) * 8.0 * similar.SimilarityScore
			validSamples++
		}
	}

	if validSamples == 0 {
		return itc.estimateTimeFromContent(title, description), 0.2
	}

	avgHours := totalHours / float64(validSamples)
	confidence := float64(validSamples) / 10.0
	if confidence > 0.9 {
		confidence = 0.9
	}

	return avgHours, confidence
}

// estimateTimeFromContent provides fallback time estimation
func (itc *IntelligentTaskCreator) estimateTimeFromContent(title, description string) float64 {
	content := title + " " + description
	wordCount := len(strings.Fields(content))

	// Simple heuristic based on content length and complexity
	baseHours := float64(wordCount) / constants.WordsPerHourEstimate // Time estimation

	// Adjust based on complexity indicators
	complexity := itc.detectComplexityLevel(content)
	multiplier := map[string]float64{
		"low":    0.5,
		"medium": 1.0,
		"high":   2.0,
	}[complexity]

	estimatedHours := baseHours * multiplier

	// Reasonable bounds
	if estimatedHours < constants.MinimumTaskTimeHours {
		estimatedHours = constants.MinimumTaskTimeHours
	} else if estimatedHours > constants.MaxEstimatedHours {
		estimatedHours = constants.MaxEstimatedHours
	}

	return estimatedHours
}

// findSimilarTasks finds tasks similar to the given title and description
func (itc *IntelligentTaskCreator) findSimilarTasks(ctx context.Context, title, description string, limit int) []SimilarTaskInfo {
	allTasks, err := itc.veriYonetici.GorevListele(ctx, map[string]interface{}{})
	if err != nil {
		return nil
	}

	targetWords := extractKeywords(title + " " + description)
	var similarities []SimilarTaskInfo

	for _, task := range allTasks {
		taskWords := extractKeywords(task.Title + " " + task.Description)
		similarity := calculateSimilarity(targetWords, taskWords)

		if similarity > 0.2 { // 20% similarity threshold
			reason := itc.generateSimilarityReason(targetWords, taskWords)
			similarities = append(similarities, SimilarTaskInfo{
				Task:            task,
				SimilarityScore: similarity,
				Reason:          reason,
			})
		}
	}

	// Sort by similarity score
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].SimilarityScore > similarities[j].SimilarityScore
	})

	// Limit results
	if len(similarities) > limit {
		similarities = similarities[:limit]
	}

	return similarities
}

// suggestTemplate recommends a template based on content analysis
func (itc *IntelligentTaskCreator) suggestTemplate(title, description string) (string, float64) {
	content := strings.ToLower(title + " " + description)

	// Template matching patterns
	templatePatterns := map[string][]string{
		"bug_report": {
			"bug", "hata", "error", "exception", "crash", "çöktü", "broken",
			"not working", "çalışmıyor", "issue", "problem",
		},
		"feature_request": {
			"feature", "özellik", "functionality", "fonksiyon", "add", "ekle",
			"implement", "uygula", "new", "yeni", "enhancement",
		},
		"research_task": {
			"research", "araştırma", "investigate", "incele", "analyze", "analiz",
			"study", "çalışma", "explore", "keşfet", "evaluate",
		},
		"performance_issue": {
			"performance", "performans", "slow", "yavaş", "optimization", "optimizasyon",
			"speed", "hız", "memory", "bellek", "cpu", "latency",
		},
		"security_fix": {
			"security", "güvenlik", "vulnerability", "zafiyet", "auth", "authentication",
			"authorization", "yetki", "encryption", "şifreleme",
		},
	}

	bestTemplate := ""
	bestScore := 0.0

	for template, keywords := range templatePatterns {
		score := float64(itc.countKeywords(content, keywords))
		if score > bestScore {
			bestScore = score
			bestTemplate = template
		}
	}

	if bestScore > 0 {
		confidence := bestScore / 10.0
		if confidence > 1.0 {
			confidence = 1.0
		}
		return bestTemplate, confidence
	}

	return "", 0.0
}

// generateSubtasks automatically creates subtasks from description
func (itc *IntelligentTaskCreator) generateSubtasks(parentID, description string) ([]*Gorev, float64) {
	if len(description) < 50 {
		return nil, 0.0
	}

	// Look for structured content that indicates subtasks
	subtaskPatterns := []string{
		`\d+\.\s+(.+)`,             // "1. Do something"
		`[-*]\s+(.+)`,              // "- Do something" or "* Do something"
		`\n(.+?)[:]\s*\n`,          // "Step: description"
		`(?i)step\s+\d+[:]\s*(.+)`, // "Step 1: description"
		`(?i)task\s+\d+[:]\s*(.+)`, // "Task 1: description"
	}

	var subtasks []*Gorev
	var foundItems []string

	for _, pattern := range subtaskPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(description, -1)

		for _, match := range matches {
			if len(match) > 1 && len(strings.TrimSpace(match[1])) > 5 {
				item := strings.TrimSpace(match[1])
				// Avoid duplicates
				if !contains(foundItems, item) {
					foundItems = append(foundItems, item)
				}
			}
		}

		if len(foundItems) >= constants.MaxSubtasksAuto { // Limit to reasonable number
			break
		}
	}

	// Create subtasks
	for i, item := range foundItems {
		if i >= constants.MaxSubtasksAuto { // Maximum subtasks
			break
		}

		subtask, err := itc.veriYonetici.AltGorevOlustur(context.Background(), 
			parentID,
			item,
			"",                       // No description for auto-generated subtasks
			constants.PriorityMedium, // Default priority
			"",                       // No due date
			nil,                      // No tags
		)

		if err != nil {
			log.Printf("Failed to create subtask: error=%v, title=%s", err, item)
			continue
		}

		subtasks = append(subtasks, subtask)
	}

	confidence := 0.0
	if len(foundItems) > 0 {
		confidence = float64(len(subtasks)) / float64(len(foundItems))
	}

	return subtasks, confidence
}

// Helper functions

func (itc *IntelligentTaskCreator) containsUrgencyKeywords(content string) bool {
	urgencyKeywords := []string{
		"urgent", "acil", "asap", "immediately", "hemen", "critical", "kritik",
	}
	return itc.countKeywords(content, urgencyKeywords) > 0
}

func (itc *IntelligentTaskCreator) detectComplexityIndicators(content string) []string {
	var indicators []string

	complexityPatterns := map[string]string{
		"multiple_steps":  `\d+\.\s+|\n[-*]\s+`,
		"technical_terms": `api|database|algorithm|integration|deployment`,
		"time_mentions":   `\d+\s*(hour|day|week|month)s?`,
		"dependencies":    `depend|require|need|after|before`,
	}

	for indicator, pattern := range complexityPatterns {
		if matched, _ := regexp.MatchString(pattern, content); matched {
			indicators = append(indicators, indicator)
		}
	}

	return indicators
}

func (itc *IntelligentTaskCreator) detectComplexityLevel(content string) string {
	indicators := itc.detectComplexityIndicators(content)
	wordCount := len(strings.Fields(content))

	if len(indicators) >= constants.ComplexityIndicatorMinimum || wordCount > constants.ComplexityWordCountHigh {
		return "high"
	} else if len(indicators) >= 1 || wordCount > constants.ComplexityWordCountMedium {
		return "medium"
	} else {
		return "low"
	}
}

func (itc *IntelligentTaskCreator) countKeywords(content string, keywords []string) int {
	count := 0
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			count++
		}
	}
	return count
}

func (itc *IntelligentTaskCreator) generateSimilarityReason(words1, words2 []string) string {
	commonWords := []string{}
	wordSet1 := make(map[string]bool)

	for _, word := range words1 {
		wordSet1[word] = true
	}

	for _, word := range words2 {
		if wordSet1[word] && len(word) > 3 { // Only meaningful words
			commonWords = append(commonWords, word)
		}
	}

	if len(commonWords) == 0 {
		return "Genel benzerlik"
	}

	if len(commonWords) <= constants.MaxCommonWordsDisplay {
		return fmt.Sprintf("Ortak kelimeler: %s", strings.Join(commonWords, ", "))
	} else {
		return fmt.Sprintf("Ortak kelimeler: %s ve %d tane daha",
			strings.Join(commonWords[:constants.MaxCommonWordsDisplay], ", "), len(commonWords)-constants.MaxCommonWordsDisplay)
	}
}

// extractDurationFromDescription extracts estimated duration in days from task description
func (itc *IntelligentTaskCreator) extractDurationFromDescription(description string) int {
	// Simple pattern matching for common duration expressions
	// In a real implementation, this could use more sophisticated NLP
	desc := strings.ToLower(description)

	if strings.Contains(desc, "1 gün") || strings.Contains(desc, "1 day") {
		return 1
	}
	if strings.Contains(desc, "2 gün") || strings.Contains(desc, "2 day") {
		return 2
	}
	if strings.Contains(desc, "3 gün") || strings.Contains(desc, "3 day") {
		return 3
	}
	if strings.Contains(desc, "1 hafta") || strings.Contains(desc, "1 week") {
		return 5
	}
	if strings.Contains(desc, "hızlı") || strings.Contains(desc, "quick") || strings.Contains(desc, "basit") {
		return 1
	}

	return 0 // Default if no duration found
}
