package cmsstore

import "net/http"

// MiddlewareInterface defines the structure of a middleware.
type MiddlewareInterface interface {
	// Identifier is a unique identifier for internal use (e.g., "auth_before").
	// Must be unique, and cannot be changed after creation.
	Identifier() string

	// Name is a human-friendly label for display purposes (e.g., "Authentication Middleware").
	Name() string

	// Description provides details about the middleware’s functionality.
	Description() string

	// Type specifies when the middleware is executed:
	// - "before"  → Runs before rendering page content.
	// - "after"   → Runs after rendering page content.
	Type() string

	// Handler returns the middleware function that processes HTTP requests.
	Handler() func(next http.Handler) http.Handler
}

func Middleware() *middleware {
	m := new(middleware)
	m.properties = map[string]any{}
	return m
}

var _ MiddlewareInterface = (*middleware)(nil)

type middleware struct {
	properties map[string]any
}

// Identifier is a unique identifier for internal use (e.g., "auth_before").
// Must be unique. Cannot be changed after creation.
func (m *middleware) Identifier() string {
	if m.hasProperty("identifier") {
		return m.properties["identifier"].(string)
	}

	return ""
}

// SetIdentifier sets the identifier of the middleware.
func (m *middleware) SetIdentifier(identifier string) *middleware {
	m.properties["identifier"] = identifier
	return m
}

// Name is a human-friendly label for display purposes (e.g., "Authentication Middleware").
func (m *middleware) Name() string {
	if m.hasProperty("name") {
		return m.properties["name"].(string)
	}

	return ""
}

// SetName sets the name of the middleware.
func (m *middleware) SetName(name string) *middleware {
	m.properties["name"] = name
	return m
}

// Description provides details about the middleware’s functionality.
func (m *middleware) Description() string {
	if m.hasProperty("description") {
		return m.properties["description"].(string)
	}

	return ""
}

// SetDescription sets the description of the middleware.
func (m *middleware) SetDescription(description string) *middleware {
	m.properties["description"] = description
	return m
}

// Handler returns the middleware function that processes HTTP requests.
func (m *middleware) Handler() func(next http.Handler) http.Handler {
	if m.hasProperty("handler") {
		return m.properties["handler"].(func(next http.Handler) http.Handler)
	}

	return func(next http.Handler) http.Handler {
		return next
	}
}

// SetHandler sets the middleware function that processes HTTP requests.
func (m *middleware) SetHandler(handler func(next http.Handler) http.Handler) *middleware {
	m.properties["handler"] = handler
	return m
}

func (m *middleware) hasProperty(key string) bool {
	_, ok := m.properties[key]
	return ok
}

// Type specifies when the middleware is executed:
// - "before"  → Runs before rendering page content.
// - "after"   → Runs after rendering page content.
// - "replace" → Modifies or replaces page content.
func (m *middleware) Type() string {
	if m.hasProperty("type") {
		return m.properties["type"].(string)
	}

	return ""
}

// SetType specifies when the middleware is executed:
// - "before"  → Runs before rendering page content.
// - "after"   → Runs after rendering page content.
// - "replace" → Modifies or replaces page content.
func (m *middleware) SetType(middlewareType string) *middleware {
	m.properties["type"] = middlewareType
	return m
}
