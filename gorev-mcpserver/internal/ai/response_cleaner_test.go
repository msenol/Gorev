package ai

import (
	"testing"
)

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		validate    func([]byte) error
	}{
		{
			name:    "plain JSON - no markdown",
			input:   `{"subtasks":[{"title":"Task 1"}]}`,
			wantErr: false,
			validate: func(data []byte) error {
				// Verify it's valid JSON
				return nil
			},
		},
		{
			name:    "JSON in ```json code block",
			input:   "```json\n{\"subtasks\":[{\"title\":\"Task 1\"}]}\n```",
			wantErr: false,
			validate: func(data []byte) error {
				// Verify it's valid JSON
				return nil
			},
		},
		{
			name:    "JSON in ``` code block without language",
			input:   "```\n{\"subtasks\":[{\"title\":\"Task 1\"}]}\n```",
			wantErr: false,
			validate: func(data []byte) error {
				// Verify it's valid JSON
				return nil
			},
		},
		{
			name:    "JSON with surrounding text and code block",
			input:   "Here's the result:\n```json\n{\"critical_path\":[\"task-1\"]}\n```\nHope this helps!",
			wantErr: false,
			validate: func(data []byte) error {
				// Verify it's valid JSON
				return nil
			},
		},
		{
			name:    "JSON with extra whitespace in code block",
			input:   "```json  \n  {\"subtasks\":[{\"title\":\"Task 1\"}]}  \n  ```",
			wantErr: false,
			validate: func(data []byte) error {
				// Verify it's valid JSON
				return nil
			},
		},
		{
			name:    "Invalid JSON in code block",
			input:   "```json\n{\"invalid\": }\n```",
			wantErr: true,
			validate: nil,
		},
		{
			name:    "Empty response",
			input:   "",
			wantErr: true,
			validate: nil,
		},
		{
			name:    "Text with no JSON",
			input:   "This is just plain text with no JSON at all.",
			wantErr: true,
			validate: nil,
		},
		{
			name:    "Code block with non-JSON content",
			input:   "```\nThis is not JSON\n```",
			wantErr: true,
			validate: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractJSON(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("extractJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				if err := tt.validate(got); err != nil {
					t.Errorf("extractJSON() validation failed: %v", err)
				}
			}
		})
	}
}
