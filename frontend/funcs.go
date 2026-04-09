package frontend

import (
	"encoding/json"
	"regexp"
)

// returns the IDs in the content who have the following format [[prefix_id]]
func contentFindIdsByPatternPrefix(content, prefix string) []string {
	ids := []string{}

	re := regexp.MustCompilePOSIX(`\[\[` + prefix + `_([a-zA-Z0-9_-]+)\]\]`)

	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if match[0] == "" {
			continue
		}
		if match[1] == "" {
			continue // no need to add empty IDs
		}
		ids = append(ids, match[1])
	}

	return ids
}

func isJSON(str string) bool {
	var js any
	err := json.Unmarshal([]byte(str), &js)
	if err != nil {
		return false
	}

	// Accept both objects (map) and arrays (slice), but not primitives
	switch js.(type) {
	case map[string]any:
		return true
	case []any:
		return true
	default:
		return false
	}
}
