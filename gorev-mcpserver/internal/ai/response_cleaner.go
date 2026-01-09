package ai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// extractJSON extracts clean JSON from AI response
// Handles markdown code blocks and plain JSON
// Returns cleaned JSON bytes or error
func extractJSON(rawContent string) ([]byte, error) {
	if rawContent == "" {
		return nil, fmt.Errorf("empty response")
	}

	// Try direct parsing first (for clean responses)
	var result interface{}
	if err := json.Unmarshal([]byte(rawContent), &result); err == nil {
		return []byte(rawContent), nil
	}

	// Extract from markdown code blocks
	// Pattern: ```json...``` or ```...```
	// Note: Using string concatenation to avoid backtick conflicts in raw strings
	codeBlockRegex := regexp.MustCompile("(?s)" + "```" + `(?:json)?\s*(.*?)\s*` + "```")
	matches := codeBlockRegex.FindStringSubmatch(rawContent)

	if len(matches) > 1 {
		extracted := strings.TrimSpace(matches[1])
		unmarshalErr := json.Unmarshal([]byte(extracted), &result)
		if unmarshalErr == nil {
			return []byte(extracted), nil
		}
		return nil, fmt.Errorf("extracted content is not valid JSON: %w", unmarshalErr)
	}

	return nil, fmt.Errorf("no valid JSON found in response (tried direct parse and code block extraction)")
}
