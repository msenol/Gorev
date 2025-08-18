package constants

// Task status constants to eliminate hardcoded strings throughout the codebase
const (
	// TaskStatusPending represents a task that is waiting to be started
	TaskStatusPending = "beklemede"

	// TaskStatusInProgress represents a task that is currently being worked on
	TaskStatusInProgress = "devam_ediyor"

	// TaskStatusCompleted represents a task that has been completed
	TaskStatusCompleted = "tamamlandi"

	// TaskStatusCancelled represents a task that has been cancelled
	TaskStatusCancelled = "iptal"
)

// Priority constants to eliminate hardcoded strings throughout the codebase
const (
	// PriorityLow represents low priority tasks
	PriorityLow = "dusuk"

	// PriorityMedium represents medium priority tasks
	PriorityMedium = "orta"

	// PriorityHigh represents high priority tasks
	PriorityHigh = "yuksek"
)

// Dependency type constants
const (
	// DependencyTypeBlocker represents a blocking dependency
	DependencyTypeBlocker = "blocker"

	// DependencyTypeDependsOn represents a depends-on dependency
	DependencyTypeDependsOn = "depends_on"
)

// Common task limits and defaults
const (
	// DefaultTaskLimit is the default number of tasks to return in listings
	DefaultTaskLimit = 50

	// MaxTaskLimit is the maximum number of tasks allowed in a single request
	MaxTaskLimit = 200

	// DefaultRecentTaskLimit is the default number of recent tasks to return
	DefaultRecentTaskLimit = 5

	// DefaultPaginationOffset is the default pagination offset
	DefaultPaginationOffset = 0

	// MaxInlineDescriptionLength is the maximum length for inline task descriptions
	MaxInlineDescriptionLength = 50

	// TagSizeConstant is the constant added to tag name length for calculations
	TagSizeConstant = 5

	// TruncatedDescriptionLength is the length for truncated descriptions
	TruncatedDescriptionLength = 47

	// DefaultSuggestionLimit is the default number of suggestions to return
	DefaultSuggestionLimit = 10

	// BaseResponseSize is the base size estimate for formatting responses
	BaseResponseSize = 100

	// DependencyInfoSize is the size estimate for dependency information
	DependencyInfoSize = 100

	// MaxDescriptionDisplayLength is the maximum length for displaying descriptions
	MaxDescriptionDisplayLength = 100

	// MaxRecentTasks is the maximum number of recent tasks to keep in AI context
	MaxRecentTasks = 10

	// MaxSummaryItems is the maximum number of items to show in AI context summaries
	MaxSummaryItems = 5

	// MaxInteractionHistory is the maximum number of AI interactions to retrieve
	MaxInteractionHistory = 50

	// DateFormatSize is the estimated size for date formatting
	DateFormatSize = 30

	// ShortIDLength is the length for truncated task IDs
	ShortIDLength = 8

	// MaxResponseSize is the maximum size for MCP responses
	MaxResponseSize = 20000

	// MaxDescriptionLength is the maximum length before truncating descriptions
	MaxDescriptionLength = 2000

	// MaxDescriptionTruncateLength is the length to truncate to (with ... added)
	MaxDescriptionTruncateLength = 1997

	// LastCreatedCount is the number for "last created" queries
	LastCreatedCount = 1

	// RecentlyCreatedCount is the number for "recently created" queries
	RecentlyCreatedCount = 5

	// MaxCommonWordsDisplay is the maximum common words to display before "and X more"
	MaxCommonWordsDisplay = 3

	// MaxEstimatedHours is the maximum hours for task estimation
	MaxEstimatedHours = 40
)

// Date format constants to eliminate hardcoded date formats
const (
	// DateFormatISO is the standard ISO date format YYYY-MM-DD
	DateFormatISO = "2006-01-02"

	// DateTimeFormatFull is the full datetime format
	DateTimeFormatFull = "2006-01-02 15:04:05"

	// DateFormatDisplay is the human-readable date format
	DateFormatDisplay = "02 Jan 2006, 15:04"

	// DateFormatShort is the short date format for compact display
	DateFormatShort = "02/01"
)

// Confidence score constants for AI/ML operations
const (
	// ConfidenceVeryHigh for extremely reliable suggestions
	ConfidenceVeryHigh = 0.9

	// ConfidenceHigh for reliable suggestions
	ConfidenceHigh = 0.8

	// ConfidenceMedium for moderately reliable suggestions
	ConfidenceMedium = 0.7

	// ConfidenceLow for less reliable suggestions
	ConfidenceLow = 0.6

	// ConfidenceVeryLow for minimal confidence suggestions
	ConfidenceVeryLow = 0.3

	// ConfidenceThreshold minimum confidence for suggestions
	ConfidenceThreshold = 0.3

	// ConfidenceWeightHigh for high-importance NLP factors
	ConfidenceWeightHigh = 0.4

	// ConfidenceWeightMedium for medium-importance NLP factors
	ConfidenceWeightMedium = 0.3

	// ConfidenceWeightLow for low-importance NLP factors
	ConfidenceWeightLow = 0.2

	// ConfidenceWeightMinimal for minimal confidence adjustments
	ConfidenceWeightMinimal = 0.1

	// DefaultConfidence when no specific confidence can be calculated
	DefaultConfidence = 0.5
)

// Limit and threshold constants for various operations
const (
	// MaxTagsToDisplay maximum tags to show before summarizing
	MaxTagsToDisplay = 3

	// MaxSuggestionsToShow maximum suggestions to display
	MaxSuggestionsToShow = 3

	// SimilarityThreshold minimum similarity score for task matching
	SimilarityThreshold = 0.3

	// MinWordLength minimum word length for keyword extraction
	MinWordLength = 2

	// ConfidenceNormalizer for normalizing confidence scores
	ConfidenceNormalizer = 10.0

	// MaxSubtasksAuto maximum subtasks to auto-generate
	MaxSubtasksAuto = 5

	// WordsPerHourEstimate words per hour for time estimation
	WordsPerHourEstimate = 20.0

	// ComplexityWordCountHigh word count threshold for high complexity
	ComplexityWordCountHigh = 100

	// ComplexityWordCountMedium word count threshold for medium complexity
	ComplexityWordCountMedium = 30

	// ComplexityIndicatorMinimum minimum indicators for complexity detection
	ComplexityIndicatorMinimum = 3

	// MinimumTaskTimeHours minimum hours for task estimation
	MinimumTaskTimeHours = 0.5
)

// Priority scoring constants for suggestion ranking
const (
	// PriorityScoreHigh numeric value for high priority ranking
	PriorityScoreHigh = 3

	// PriorityScoreMedium numeric value for medium priority ranking
	PriorityScoreMedium = 2

	// PriorityScoreLow numeric value for low priority ranking
	PriorityScoreLow = 1
)

// GetValidTaskStatuses returns all valid task status values
func GetValidTaskStatuses() []string {
	return []string{
		TaskStatusPending,
		TaskStatusInProgress,
		TaskStatusCompleted,
		TaskStatusCancelled,
	}
}

// GetValidPriorities returns all valid priority values
func GetValidPriorities() []string {
	return []string{
		PriorityLow,
		PriorityMedium,
		PriorityHigh,
	}
}

// GetValidDependencyTypes returns all valid dependency type values
func GetValidDependencyTypes() []string {
	return []string{
		DependencyTypeBlocker,
		DependencyTypeDependsOn,
	}
}

// IsValidTaskStatus checks if a given status is valid
func IsValidTaskStatus(status string) bool {
	for _, validStatus := range GetValidTaskStatuses() {
		if status == validStatus {
			return true
		}
	}
	return false
}

// IsValidPriority checks if a given priority is valid
func IsValidPriority(priority string) bool {
	for _, validPriority := range GetValidPriorities() {
		if priority == validPriority {
			return true
		}
	}
	return false
}

// IsValidDependencyType checks if a given dependency type is valid
func IsValidDependencyType(depType string) bool {
	for _, validType := range GetValidDependencyTypes() {
		if depType == validType {
			return true
		}
	}
	return false
}

// Template-related constants for form options
var (
	// ValidEnvironments for deployment environments in templates
	ValidEnvironments = []string{"development", "staging", "production"}

	// ValidEffortLevels for effort estimation in templates
	ValidEffortLevels = []string{"küçük", "orta", "büyük"}
)

// Error format constants to eliminate duplicated error patterns
const (
	// ErrorFormatValidation for validation errors
	ErrorFormatValidation = "validation failed: %v"

	// ErrorFormatTool for tool operation errors
	ErrorFormatTool = "tool error: %v"

	// ErrorFormatDatabase for database operation errors
	ErrorFormatDatabase = "database error: %v"

	// ErrorFormatMigration for migration errors
	ErrorFormatMigration = "migration failed: %v"

	// ErrorFormatTemplate for template errors
	ErrorFormatTemplate = "template error: %v"

	// ErrorFormatParamRequired for required parameter errors
	ErrorFormatParamRequired = "required parameter missing: %s"

	// ErrorFormatParamInvalid for invalid parameter errors
	ErrorFormatParamInvalid = "invalid parameter value: %s = %v"

	// ErrorFormatTaskNotFound for task not found errors
	ErrorFormatTaskNotFound = "task not found: %s"

	// ErrorFormatProjectNotFound for project not found errors
	ErrorFormatProjectNotFound = "project not found: %s"

	// ErrorFormatUnauthorized for authorization errors
	ErrorFormatUnauthorized = "unauthorized operation: %s"

	// ErrorFormatTimeout for timeout errors
	ErrorFormatTimeout = "operation timeout: %s"

	// ErrorFormatUnexpected for unexpected errors
	ErrorFormatUnexpected = "unexpected error in %s: %v"
)

// Common error prefix patterns
const (
	// ErrorPrefixFailed for "Failed to..." patterns
	ErrorPrefixFailed = "Failed to %s"

	// ErrorPrefixCannot for "Cannot..." patterns
	ErrorPrefixCannot = "Cannot %s"

	// ErrorPrefixUnable for "Unable to..." patterns
	ErrorPrefixUnable = "Unable to %s"
)
