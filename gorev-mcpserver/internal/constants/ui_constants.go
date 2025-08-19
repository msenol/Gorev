package constants

// Status emoji constants to eliminate hardcoded emojis throughout the codebase
const (
	// Task status emojis
	EmojiStatusPending    = "⏳" // beklemede
	EmojiStatusInProgress = "🔄" // devam_ediyor
	EmojiStatusCompleted  = "✅" // tamamlandi
	EmojiStatusCancelled  = "❌" // iptal
	EmojiStatusUnknown    = "⚪" // bilinmeyen durum
)

// Priority emoji constants
const (
	// Task priority emojis
	EmojiPriorityHigh    = "🔴" // yuksek
	EmojiPriorityMedium  = "🟡" // orta
	EmojiPriorityLow     = "🟢" // dusuk
	EmojiPriorityUnknown = "⚪" // bilinmeyen öncelik

	// Alternative priority emojis
	EmojiPriorityHighAlt = "🔥"  // yuksek (alternatif)
	EmojiPriorityAlert   = "⚠️" // uyarı
	
	// Suggestion priority emojis (for AI suggestions)
	EmojiSuggestionHigh   = "🔥"  // yuksek oneri
	EmojiSuggestionMedium = "⚡"  // orta oneri
	EmojiSuggestionLow    = "ℹ️"  // dusuk oneri
)

// Message prefix constants
const (
	// Success message prefix
	PrefixSuccess = "✓ "

	// Error message prefix
	PrefixError = "✗ "

	// Warning message prefix
	PrefixWarning = "⚠ "

	// Info message prefix
	PrefixInfo = "ℹ "

	// Loading/Processing prefix
	PrefixLoading = "⌛ "
)

// Formatting symbols
const (
	// Markdown formatting
	MarkdownBold   = "**"
	MarkdownItalic = "*"
	MarkdownCode   = "`"

	// List bullets
	BulletPoint = "•"
	BulletArrow = "→"
	BulletCheck = "✓"
	BulletCross = "✗"

	// Separators
	SeparatorDash  = " - "
	SeparatorColon = ": "
	SeparatorPipe  = " | "
	SeparatorComma = ", "
)

// Icon constants for UI elements
const (
	// Common icons
	IconTask     = "📋"
	IconProject  = "📁"
	IconTag      = "🏷️"
	IconDate     = "📅"
	IconTime     = "⏰"
	IconUser     = "👤"
	IconSettings = "⚙️"
	IconHelp     = "❓"
	IconSearch   = "🔍"
	IconFilter   = "🔽"
	IconSort     = "🔄"
	IconAdd      = "➕"
	IconEdit     = "✏️"
	IconDelete   = "🗑️"
	IconSave     = "💾"
	IconCancel   = "❌"
	IconRefresh  = "🔄"
	IconExport   = "📤"
	IconImport   = "📥"
)

// Progress indicators
const (
	// Progress bars
	ProgressEmpty = "░"
	ProgressFull  = "█"
	ProgressHalf  = "▌"

	// Spinner characters
	SpinnerChars = "|/-\\"
)

// Template emoji constants for different template types
const (
	// Template category emojis
	EmojiTemplateBug        = "🐛"
	EmojiTemplateFeature    = "✨"
	EmojiTemplateTask       = "📋"
	EmojiTemplateMeeting    = "👥"
	EmojiTemplateResearch   = "🔬"
	EmojiTemplateSecurity   = "🔒"
	EmojiTemplateRefactor   = "🔧"
	EmojiTemplateDoc        = "📚"
	EmojiTemplateTest       = "🧪"
	EmojiTemplateDeployment = "🚀"
)

// Status indicator combinations
var (
	// Status with emoji combinations
	StatusDisplayMap = map[string]string{
		TaskStatusPending:    EmojiStatusPending + " " + "Beklemede",
		TaskStatusInProgress: EmojiStatusInProgress + " " + "Devam Ediyor",
		TaskStatusCompleted:  EmojiStatusCompleted + " " + "Tamamlandı",
		TaskStatusCancelled:  EmojiStatusCancelled + " " + "İptal",
	}

	// Priority with emoji combinations
	PriorityDisplayMap = map[string]string{
		PriorityHigh:   EmojiPriorityHigh + " " + "Yüksek",
		PriorityMedium: EmojiPriorityMedium + " " + "Orta",
		PriorityLow:    EmojiPriorityLow + " " + "Düşük",
	}
)

// Helper functions for UI display
func GetStatusEmoji(status string) string {
	switch status {
	case TaskStatusPending:
		return EmojiStatusPending
	case TaskStatusInProgress:
		return EmojiStatusInProgress
	case TaskStatusCompleted:
		return EmojiStatusCompleted
	case TaskStatusCancelled:
		return EmojiStatusCancelled
	default:
		return EmojiStatusUnknown
	}
}

func GetPriorityEmoji(priority string) string {
	switch priority {
	case PriorityHigh:
		return EmojiPriorityHigh
	case PriorityMedium:
		return EmojiPriorityMedium
	case PriorityLow:
		return EmojiPriorityLow
	default:
		return EmojiPriorityUnknown
	}
}

func GetStatusDisplay(status string) string {
	if display, exists := StatusDisplayMap[status]; exists {
		return display
	}
	return EmojiStatusUnknown + " Bilinmeyen"
}

func GetPriorityDisplay(priority string) string {
	if display, exists := PriorityDisplayMap[priority]; exists {
		return display
	}
	return EmojiPriorityUnknown + " Bilinmeyen"
}

func GetSuggestionPriorityEmoji(priority string) string {
	switch priority {
	case "high":
		return EmojiSuggestionHigh
	case "medium":
		return EmojiSuggestionMedium
	case "low":
		return EmojiSuggestionLow
	default:
		return EmojiPriorityUnknown
	}
}
