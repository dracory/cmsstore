package cmsstore

import (
	"context"
	"net/http"
	"testing"
)

// mockBlockType is a minimal BlockType implementation for testing GetCustomVariables.
type mockBlockType struct {
	vars []BlockCustomVariable
}

func (m *mockBlockType) TypeKey() string   { return "mock" }
func (m *mockBlockType) TypeLabel() string { return "Mock Block" }
func (m *mockBlockType) Render(_ context.Context, _ BlockInterface, _ ...RenderOption) (string, error) {
	return "", nil
}
func (m *mockBlockType) GetAdminFields(_ BlockInterface, _ *http.Request) interface{} { return nil }
func (m *mockBlockType) SaveAdminFields(_ *http.Request, _ BlockInterface) error      { return nil }
func (m *mockBlockType) GetCustomVariables() []BlockCustomVariable                    { return m.vars }

func TestBlockCustomVariable_Fields(t *testing.T) {
	v := BlockCustomVariable{
		Name:        "blog_title",
		Description: "The blog post title",
	}

	if v.Name != "blog_title" {
		t.Errorf("expected Name %q, got %q", "blog_title", v.Name)
	}
	if v.Description != "The blog post title" {
		t.Errorf("expected Description %q, got %q", "The blog post title", v.Description)
	}
}

func TestGetCustomVariables_ReturnsNil(t *testing.T) {
	bt := &mockBlockType{vars: nil}
	if bt.GetCustomVariables() != nil {
		t.Error("expected nil custom variables")
	}
}

func TestGetCustomVariables_ReturnsEmpty(t *testing.T) {
	bt := &mockBlockType{vars: []BlockCustomVariable{}}
	vars := bt.GetCustomVariables()
	if len(vars) != 0 {
		t.Errorf("expected 0 variables, got %d", len(vars))
	}
}

func TestGetCustomVariables_ReturnsVariables(t *testing.T) {
	expected := []BlockCustomVariable{
		{Name: "blog_title", Description: "The blog post title"},
		{Name: "blog_author", Description: "The post author name"},
	}
	bt := &mockBlockType{vars: expected}

	vars := bt.GetCustomVariables()
	if len(vars) != 2 {
		t.Fatalf("expected 2 variables, got %d", len(vars))
	}
	if vars[0].Name != "blog_title" {
		t.Errorf("expected Name %q, got %q", "blog_title", vars[0].Name)
	}
	if vars[0].Description != "The blog post title" {
		t.Errorf("expected Description %q, got %q", "The blog post title", vars[0].Description)
	}
	if vars[1].Name != "blog_author" {
		t.Errorf("expected Name %q, got %q", "blog_author", vars[1].Name)
	}
}
