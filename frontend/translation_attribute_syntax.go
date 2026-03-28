package frontend

import (
	"html"
	"net/http"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

// Package-level compiled regex for performance
// Matches: <translation ...attributes... /> OR [[translation ...attributes...]]
var translationAttributeAngleBrackets = regexp.MustCompile(`<translation\s+([^>]+?)\s*/>`)
var translationAttributeSquareBrackets = regexp.MustCompile(`\[\[translation\s+([^\]]+?)\s*\]\]`)

// applyTranslationAttributeSyntax processes translation references with attributes.
// Supports two syntaxes:
// 1. <translation id="..." attr="value" /> - Primary syntax
// 2. [[translation id='...' attr='value']] - Alternative for HTML attribute contexts
// This is called after legacy [[TRANSLATION_id]] processing.
func (frontend *frontend) applyTranslationAttributeSyntax(req *http.Request, content string, language string) (string, error) {
	// Find all <translation ... /> tags (angle bracket syntax)
	angleMatches := translationAttributeAngleBrackets.FindAllStringSubmatch(content, -1)

	// Find all [[translation ...]] tags (square bracket syntax)
	squareMatches := translationAttributeSquareBrackets.FindAllStringSubmatch(content, -1)

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
		return content, nil // No translation references found
	}

	// Process each match
	for _, m := range allMatches {
		fullTag := m.fullTag
		attrString := m.attrs

		// Parse attributes
		attrs := parseAttributes(attrString)

		// Get translation ID (required)
		translationID := attrs["id"]
		if translationID == "" {
			frontend.logger.Warn("Translation attribute syntax: missing id attribute", "tag", fullTag)
			content = strings.Replace(content, fullTag, "<!-- Translation reference missing id -->", 1)
			continue
		}

		// Get fallback language (optional)
		fallbackLang := attrs["fallback"]

		// Fetch translation from database
		translation, err := frontend.store.TranslationFindByHandleOrID(req.Context(), translationID, language)
		if err != nil {
			frontend.logger.Error("Translation attribute syntax: error fetching translation", "id", translationID, "error", err)
			content = strings.Replace(content, fullTag, "<!-- Translation error: "+translationID+" -->", 1)
			continue
		}

		if translation == nil {
			frontend.logger.Warn("Translation attribute syntax: translation not found", "id", translationID)
			content = strings.Replace(content, fullTag, "<!-- Translation not found: "+translationID+" -->", 1)
			continue
		}

		// Get translation content map
		translationMap, err := translation.Content()
		if err != nil {
			frontend.logger.Error("Translation attribute syntax: error parsing translation content", "id", translationID, "error", err)
			content = strings.Replace(content, fullTag, "<!-- Translation parse error: "+translationID+" -->", 1)
			continue
		}

		// Get text for current language
		text := lo.ValueOr(translationMap, language, "")

		// Fallback handling
		if text == "" && fallbackLang != "" {
			text = lo.ValueOr(translationMap, fallbackLang, "")
		}

		// If still empty, use empty string
		if text == "" {
			frontend.logger.Warn("Translation attribute syntax: no translation found for language", "id", translationID, "language", language, "fallback", fallbackLang)
			text = ""
		}

		// Variable interpolation - replace {{key}} with attribute values
		// Remove system attributes (id, fallback) before interpolation
		for key, value := range attrs {
			if key != "id" && key != "fallback" {
				// Security: HTML escape interpolation values to prevent XSS
				escapedValue := html.EscapeString(value)
				placeholder := "{{" + key + "}}"
				text = strings.ReplaceAll(text, placeholder, escapedValue)
			}
		}

		// Replace the tag with translated text
		content = strings.Replace(content, fullTag, text, 1)
	}

	return content, nil
}
