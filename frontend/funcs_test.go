package frontend

import (
	"testing"
)

// TestContentFindIdsByPatternPrefix tests the contentFindIdsByPatternPrefix function
func TestContentFindIdsByPatternPrefix(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		prefix   string
		expected []string
	}{
		{
			name:     "Empty content",
			content:  "",
			prefix:   "BLOCK",
			expected: []string{},
		},
		{
			name:     "No matches",
			content:  "Hello world",
			prefix:   "BLOCK",
			expected: []string{},
		},
		{
			name:     "Single BLOCK match",
			content:  "[[BLOCK_abc123]]",
			prefix:   "BLOCK",
			expected: []string{"abc123"},
		},
		{
			name:     "Multiple BLOCK matches",
			content:  "[[BLOCK_abc123]] and [[BLOCK_def456]]",
			prefix:   "BLOCK",
			expected: []string{"abc123", "def456"},
		},
		{
			name:     "BLOCK with spaces variant",
			content:  "[[ BLOCK_abc123 ]]",
			prefix:   "BLOCK",
			expected: []string{}, // Only matches without spaces
		},
		{
			name:     "TRANSLATION prefix",
			content:  "[[TRANSLATION_xyz789]]",
			prefix:   "TRANSLATION",
			expected: []string{"xyz789"},
		},
		{
			name:     "PAGE_URL prefix",
			content:  "[[PAGE_URL_page123]]",
			prefix:   "PAGE_URL",
			expected: []string{"page123"},
		},
		{
			name:     "Mixed prefixes",
			content:  "[[BLOCK_abc]] [[TRANSLATION_def]] [[PAGE_URL_ghi]]",
			prefix:   "BLOCK",
			expected: []string{"abc"},
		},
		{
			name:     "ID with underscores",
			content:  "[[BLOCK_abc_def_ghi]]",
			prefix:   "BLOCK",
			expected: []string{"abc_def_ghi"},
		},
		{
			name:     "ID with hyphens",
			content:  "[[BLOCK_abc-def-ghi]]",
			prefix:   "BLOCK",
			expected: []string{"abc-def-ghi"},
		},
		{
			name:     "ID with numbers",
			content:  "[[BLOCK_123456]]",
			prefix:   "BLOCK",
			expected: []string{"123456"},
		},
		{
			name:     "ID with mixed alphanumeric",
			content:  "[[BLOCK_abc123-def456_ghi789]]",
			prefix:   "BLOCK",
			expected: []string{"abc123-def456_ghi789"},
		},
		{
			name:     "Empty ID after prefix",
			content:  "[[BLOCK_]]",
			prefix:   "BLOCK",
			expected: []string{}, // Empty ID should be skipped
		},
		{
			name:     "Partial match should not match",
			content:  "BLOCK_abc123",
			prefix:   "BLOCK",
			expected: []string{},
		},
		{
			name:     "Duplicate IDs",
			content:  "[[BLOCK_abc]] [[BLOCK_abc]]",
			prefix:   "BLOCK",
			expected: []string{"abc", "abc"},
		},
		{
			name:     "Content with HTML",
			content:  `<div>[[BLOCK_header]]</div><p>[[BLOCK_content]]</p>`,
			prefix:   "BLOCK",
			expected: []string{"header", "content"},
		},
		{
			name:     "Case sensitivity",
			content:  "[[block_lowercase]] [[BLOCK_UPPERCASE]]",
			prefix:   "BLOCK",
			expected: []string{"UPPERCASE"}, // Only matches exact prefix case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contentFindIdsByPatternPrefix(tt.content, tt.prefix)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d IDs, got %d", len(tt.expected), len(result))
				return
			}

			for i, expectedID := range tt.expected {
				if result[i] != expectedID {
					t.Errorf("Expected ID %d to be %q, got %q", i, expectedID, result[i])
				}
			}
		})
	}
}

// TestIsJSON tests the isJSON function
func TestIsJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Valid JSON object",
			input:    `{"key": "value"}`,
			expected: true,
		},
		{
			name:     "Valid JSON array",
			input:    `[1, 2, 3]`,
			expected: false, // isJSON only validates objects, not arrays
		},
		{
			name:     "Valid JSON nested object",
			input:    `{"outer": {"inner": "value"}}`,
			expected: true,
		},
		{
			name:     "Valid JSON with numbers",
			input:    `{"number": 123, "float": 45.67}`,
			expected: true,
		},
		{
			name:     "Valid JSON with booleans",
			input:    `{"bool1": true, "bool2": false}`,
			expected: true,
		},
		{
			name:     "Valid JSON with null",
			input:    `{"value": null}`,
			expected: true,
		},
		{
			name:     "Valid JSON string",
			input:    `"just a string"`,
			expected: false, // isJSON only validates objects, not strings
		},
		{
			name:     "Valid JSON number",
			input:    `123`,
			expected: false, // isJSON only validates objects, not numbers
		},
		{
			name:     "Valid JSON boolean",
			input:    `true`,
			expected: false, // isJSON only validates objects, not booleans
		},
		{
			name:     "Invalid JSON - missing closing brace",
			input:    `{"key": "value"`,
			expected: false,
		},
		{
			name:     "Invalid JSON - unquoted key",
			input:    `{key: "value"}`,
			expected: false,
		},
		{
			name:     "Invalid JSON - trailing comma",
			input:    `{"key": "value",}`,
			expected: false,
		},
		{
			name:     "Plain text",
			input:    "Hello world",
			expected: false,
		},
		{
			name:     "HTML string",
			input:    `<div>Hello</div>`,
			expected: false,
		},
		{
			name:     "JSON-like but not valid",
			input:    `{key: value}`,
			expected: false,
		},
		{
			name:     "Empty object",
			input:    `{}`,
			expected: true,
		},
		{
			name:     "Empty array",
			input:    `[]`,
			expected: false, // isJSON only validates objects, not arrays
		},
		{
			name:     "JSON with escaped characters",
			input:    `{"key": "value \"with quotes\""}`,
			expected: true,
		},
		{
			name:     "JSON with unicode",
			input:    `{"key": "value with \u00e9 unicode"}`,
			expected: true,
		},
		{
			name:     "Whitespace only",
			input:    "   ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSON(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for input %q", tt.expected, result, tt.input)
			}
		})
	}
}
