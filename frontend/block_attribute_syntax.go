package frontend

import (
	"html"
	"net/http"
	"regexp"
	"strings"

	"github.com/dracory/cmsstore"
)

// Package-level compiled regex for performance
// Matches: <block ...attributes... /> OR [[block ...attributes...]]
var blockAttributeAngleBrackets = regexp.MustCompile(`<block\s+([^>]+?)\s*/>`)
var blockAttributeSquareBrackets = regexp.MustCompile(`\[\[block\s+([^\]]+?)\s*\]\]`)

// Attribute parsing regex - handles double quotes, single quotes, unquoted values, and boolean flags
// Supports hyphens in attribute names (e.g., start-level, max-depth)
var attributePattern = regexp.MustCompile(`([\w-]+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s/>]+)))?`)

// applyBlockAttributeSyntax processes block references with attributes.
// Supports two syntaxes:
// 1. <block id="..." attr="value" /> - Primary syntax
// 2. [[block id='...' attr='value']] - Alternative for HTML attribute contexts
// This is called after legacy [[BLOCK_id]] processing.
func (frontend *frontend) applyBlockAttributeSyntax(req *http.Request, content string) (string, error) {
	// Find all <block ... /> tags (angle bracket syntax)
	angleMatches := blockAttributeAngleBrackets.FindAllStringSubmatch(content, -1)

	// Find all [[block ...]] tags (square bracket syntax)
	squareMatches := blockAttributeSquareBrackets.FindAllStringSubmatch(content, -1)

	// Combine matches - store both the full tag and the attributes string
	type match struct {
		fullTag string
		attrs   string
	}
	allMatches := make([]match, 0, len(angleMatches)+len(squareMatches))
	for _, m := range angleMatches {
		allMatches = append(allMatches, match{fullTag: m[0], attrs: m[1]})
	}
	for _, m := range squareMatches {
		allMatches = append(allMatches, match{fullTag: m[0], attrs: m[1]})
	}

	if len(allMatches) == 0 {
		return content, nil // No block references found
	}

	// Process each match
	for _, m := range allMatches {
		fullTag := m.fullTag
		attrString := m.attrs

		// Parse attributes
		attrs := parseAttributes(attrString)

		// Get block ID (required)
		blockID := attrs["id"]
		if blockID == "" {
			frontend.logger.Warn("Block attribute syntax: missing id attribute", "tag", fullTag)
			content = strings.Replace(content, fullTag, "<!-- Block reference missing id -->", 1)
			continue
		}

		// Fetch block from database
		block, err := frontend.store.BlockFindByID(req.Context(), blockID)
		if err != nil {
			frontend.logger.Error("Block attribute syntax: error fetching block", "id", blockID, "error", err)
			content = strings.Replace(content, fullTag, "<!-- Block error: "+blockID+" -->", 1)
			continue
		}

		if block == nil {
			frontend.logger.Warn("Block attribute syntax: block not found", "id", blockID)
			content = strings.Replace(content, fullTag, "<!-- Block not found: "+blockID+" -->", 1)
			continue
		}

		// Security: Check if block is active
		if !block.IsActive() {
			frontend.logger.Warn("Block attribute syntax: inactive block", "id", blockID)
			content = strings.Replace(content, fullTag, "<!-- Block inactive: "+blockID+" -->", 1)
			continue
		}

		// Get block type from stored value
		blockTypeKey := block.Type()

		// Get block type (from global registry)
		blockType := cmsstore.GetBlockType(blockTypeKey)

		// Extract wrap attribute before filtering (it's handled here, not passed to renderer)
		wrapElement := attrs["wrap"]

		// Remove system attrs and sanitize before passing to renderer
		runtimeAttrs := filterAndSanitizeAttrs(attrs) // remove "id", "wrap", sanitize values

		// Render with attributes
		var htmlOutput string
		if blockType != nil {
			// Render with attributes (or without if empty)
			// Block types validate attributes internally if needed
			if len(runtimeAttrs) > 0 {
				htmlOutput, err = blockType.Render(req.Context(), block, cmsstore.WithAttributes(runtimeAttrs))
			} else {
				htmlOutput, err = blockType.Render(req.Context(), block)
			}
		} else {
			// Fallback to local renderer registry
			renderer := frontend.blockRenderers.GetRenderer(blockTypeKey)
			htmlOutput, err = renderer.Render(req.Context(), block)
		}

		if err != nil {
			frontend.logger.Error("Block attribute syntax: render error", "id", blockID, "error", err)
			htmlOutput = "<!-- Block render error: " + blockID + " -->"
		}

		// Apply wrap element if specified
		if wrapElement != "" {
			// Sanitize the wrap element name (only allow alphanumeric and hyphens)
			wrapElement = html.EscapeString(wrapElement)
			htmlOutput = "<" + wrapElement + ">" + htmlOutput + "</" + wrapElement + ">"
		}

		// Replace the tag with rendered content
		content = strings.Replace(content, fullTag, htmlOutput, 1)
	}

	return content, nil
}

// parseAttributes parses "key=value key2='value2' key3=unquoted" into map
// Handles edge cases:
// - Double quotes: key="value with spaces"
// - Single quotes: key='value with spaces'
// - Unquoted: key=value (no spaces allowed)
// - Boolean flags: key (no value, empty string)
func parseAttributes(s string) map[string]string {
	if s == "" {
		return map[string]string{}
	}

	attrs := make(map[string]string)
	matches := attributePattern.FindAllStringSubmatch(s, -1)

	for _, m := range matches {
		if len(m) < 2 || m[1] == "" {
			continue // Skip invalid matches
		}

		key := strings.TrimSpace(m[1])

		// Value could be in group 2 (double quotes), 3 (single quotes), or 4 (unquoted)
		val := ""
		if len(m) > 2 && m[2] != "" {
			val = m[2] // Double quoted
		} else if len(m) > 3 && m[3] != "" {
			val = m[3] // Single quoted
		} else if len(m) > 4 && m[4] != "" {
			val = m[4] // Unquoted
		}
		// If no value found, it's a boolean flag (empty string)

		attrs[key] = val
	}

	return attrs
}

// filterAndSanitizeAttrs removes system-reserved attributes and sanitizes values
func filterAndSanitizeAttrs(attrs map[string]string) map[string]string {
	systemAttrs := map[string]bool{"id": true, "wrap": true} // 'id' and 'wrap' are system-reserved
	filtered := make(map[string]string)
	for k, v := range attrs {
		if !systemAttrs[k] {
			// Security: HTML escape attribute values to prevent XSS
			filtered[k] = html.EscapeString(v)
		}
	}
	return filtered
}
