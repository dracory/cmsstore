package frontend

import (
	"net/http"
)

// Mock middleware implementation
type MockMiddleware struct {
	identifier     string
	name           string
	description    string
	middlewareType string
	handler        func(http.Handler) http.Handler
}

func (m *MockMiddleware) Identifier() string                       { return m.identifier }
func (m *MockMiddleware) Name() string                             { return m.name }
func (m *MockMiddleware) Description() string                      { return m.description }
func (m *MockMiddleware) Type() string                             { return m.middlewareType }
func (m *MockMiddleware) Handler() func(http.Handler) http.Handler { return m.handler }
