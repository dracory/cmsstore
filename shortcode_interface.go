package cmsstore

import "net/http"

// ShortcodeInterface defines the methods that a shortcode must implement.
// A shortcode is a reusable snippet of code or content that can be rendered
// within a web page or template.
type ShortcodeInterface interface {
	// Alias returns a unique identifier for the shortcode.
	Alias() string

	// Description provides a brief explanation of the shortcode's purpose and usage.
	Description() string

	// Render generates the final output of the shortcode based on the provided
	// HTTP request, shortcode name, and a map of attributes.
	// r: HTTP request containing the context in which the shortcode is rendered.
	// s: The name of the shortcode.
	// m: A map of attributes and their values that configure the shortcode behavior.
	Render(r *http.Request, s string, m map[string]string) string
}
