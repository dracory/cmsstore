package cmsstore

import "net/http"

type MiddlewareInterface interface {
	// Human friendly name of the middleware
	Name() string

	// Description of the middleware
	Description() string

	// Handler for the middleware, this is the actual middleware
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

func (m *middleware) Name() string {
	if m.hasProperty("name") {
		return m.properties["name"].(string)
	}

	return ""
}

func (m *middleware) SetName(name string) *middleware {
	m.properties["name"] = name
	return m
}

func (m *middleware) Description() string {
	if m.hasProperty("description") {
		return m.properties["description"].(string)
	}

	return ""
}

func (m *middleware) SetDescription(description string) *middleware {
	m.properties["description"] = description
	return m
}

func (m *middleware) Handler() func(next http.Handler) http.Handler {
	if m.hasProperty("handler") {
		return m.properties["handler"].(func(next http.Handler) http.Handler)
	}

	return func(next http.Handler) http.Handler {
		return next
	}
}

func (m *middleware) SetHandler(handler func(next http.Handler) http.Handler) *middleware {
	m.properties["handler"] = handler
	return m
}

func (m *middleware) hasProperty(key string) bool {
	_, ok := m.properties[key]
	return ok
}
