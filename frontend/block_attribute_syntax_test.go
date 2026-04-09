package frontend

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

func TestParseAttributes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "double quotes",
			input: `id="block123" depth="2"`,
			expected: map[string]string{
				"id":    "block123",
				"depth": "2",
			},
		},
		{
			name:  "single quotes",
			input: `id='block123' depth='2'`,
			expected: map[string]string{
				"id":    "block123",
				"depth": "2",
			},
		},
		{
			name:  "unquoted values",
			input: `id=block123 depth=2`,
			expected: map[string]string{
				"id":    "block123",
				"depth": "2",
			},
		},
		{
			name:  "mixed quotes",
			input: `id="block123" style='horizontal' class=custom`,
			expected: map[string]string{
				"id":    "block123",
				"style": "horizontal",
				"class": "custom",
			},
		},
		{
			name:  "boolean flags",
			input: `id="block123" featured`,
			expected: map[string]string{
				"id":       "block123",
				"featured": "",
			},
		},
		{
			name:  "values with spaces in quotes",
			input: `id="block123" class="my custom class"`,
			expected: map[string]string{
				"id":    "block123",
				"class": "my custom class",
			},
		},
		{
			name:     "empty string",
			input:    ``,
			expected: map[string]string{},
		},
		{
			name:  "hyphenated attribute names",
			input: `id="block123" start-level="1" max-depth="3"`,
			expected: map[string]string{
				"id":          "block123",
				"start-level": "1",
				"max-depth":   "3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAttributes(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("parseAttributes() length = %v, want %v", len(result), len(tt.expected))
			}

			for key, expectedVal := range tt.expected {
				if actualVal, ok := result[key]; !ok {
					t.Errorf("parseAttributes() missing key %q", key)
				} else if actualVal != expectedVal {
					t.Errorf("parseAttributes() key %q = %q, want %q", key, actualVal, expectedVal)
				}
			}
		})
	}
}

func TestFilterAndSanitizeAttrs(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name: "removes id attribute",
			input: map[string]string{
				"id":    "block123",
				"depth": "2",
			},
			expected: map[string]string{
				"depth": "2",
			},
		},
		{
			name: "HTML escapes values",
			input: map[string]string{
				"class": "<script>alert('xss')</script>",
				"style": "color: red;",
			},
			expected: map[string]string{
				"class": "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
				"style": "color: red;",
			},
		},
		{
			name: "escapes quotes and ampersands",
			input: map[string]string{
				"title": `Hello "World" & Friends`,
			},
			expected: map[string]string{
				"title": "Hello &#34;World&#34; &amp; Friends",
			},
		},
		{
			name:     "empty input",
			input:    map[string]string{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterAndSanitizeAttrs(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("filterAndSanitizeAttrs() length = %v, want %v", len(result), len(tt.expected))
			}

			for key, expectedVal := range tt.expected {
				if actualVal, ok := result[key]; !ok {
					t.Errorf("filterAndSanitizeAttrs() missing key %q", key)
				} else if actualVal != expectedVal {
					t.Errorf("filterAndSanitizeAttrs() key %q = %q, want %q", key, actualVal, expectedVal)
				}
			}
		})
	}
}

func TestApplyBlockAttributeSyntax(t *testing.T) {
	ctx := context.Background()

	// Initialize test store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test HTML block
	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetContent("<p>Hello World</p>")
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal(err)
	}

	// Create frontend instance
	fe := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: false,
	})
	f := fe.(*frontend)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "angle bracket syntax without attributes",
			content:  `<block id="` + block.ID() + `" />`,
			expected: `<p>Hello World</p>`,
		},
		{
			name:     "angle bracket syntax with wrap attribute",
			content:  `<block id="` + block.ID() + `" wrap="div" />`,
			expected: `<div><p>Hello World</p></div>`,
		},
		{
			name:     "square bracket syntax without attributes",
			content:  `[[block id='` + block.ID() + `']]`,
			expected: `<p>Hello World</p>`,
		},
		{
			name:     "square bracket syntax with wrap attribute",
			content:  `[[block id='` + block.ID() + `' wrap='section']]`,
			expected: `<section><p>Hello World</p></section>`,
		},
		{
			name:     "multiple blocks",
			content:  `<block id="` + block.ID() + `" /> and <block id="` + block.ID() + `" wrap="div" />`,
			expected: `<p>Hello World</p> and <div><p>Hello World</p></div>`,
		},
		{
			name:     "mixed syntax",
			content:  `<block id="` + block.ID() + `" /> and [[block id='` + block.ID() + `']]`,
			expected: `<p>Hello World</p> and <p>Hello World</p>`,
		},
		{
			name:     "no blocks",
			content:  `<p>Just regular content</p>`,
			expected: `<p>Just regular content</p>`,
		},
		{
			name:     "missing id attribute",
			content:  `<block depth="2" />`,
			expected: `<!-- Block reference missing id -->`,
		},
		{
			name:     "non-existent block",
			content:  `<block id="nonexistent" />`,
			expected: `<!-- Block not found: nonexistent -->`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := cmsstore.WithVarsContext(req.Context())
			result, err := f.applyBlockAttributeSyntax(ctx, req, tt.content)

			if err != nil {
				t.Errorf("applyBlockAttributeSyntax() error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("applyBlockAttributeSyntax() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestApplyBlockAttributeSyntax_InactiveBlock(t *testing.T) {
	ctx := context.Background()

	// Initialize test store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatal(err)
	}

	// Create an inactive block
	block := cmsstore.NewBlock()
	block.SetName("Inactive Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetContent("<p>Should not render</p>")
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_INACTIVE)
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal(err)
	}

	// Create frontend instance
	fe := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: false,
	})
	f := fe.(*frontend)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx = cmsstore.WithVarsContext(req.Context())
	content := `<block id="` + block.ID() + `" />`
	result, err := f.applyBlockAttributeSyntax(ctx, req, content)

	if err != nil {
		t.Errorf("applyBlockAttributeSyntax() error = %v", err)
		return
	}

	expected := `<!-- Block inactive: ` + block.ID() + ` -->`
	if result != expected {
		t.Errorf("applyBlockAttributeSyntax() = %q, want %q", result, expected)
	}
}

func TestApplyBlockAttributeSyntax_XSSPrevention(t *testing.T) {
	ctx := context.Background()

	// Initialize test store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test block
	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetContent("<p>Content</p>")
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal(err)
	}

	// Create frontend instance
	fe := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: false,
	})
	f := fe.(*frontend)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx = cmsstore.WithVarsContext(req.Context())
	// Use properly encoded HTML entities in attribute value (as browsers would parse it)
	content := `<block id="` + block.ID() + `" wrap="&lt;script&gt;alert('xss')&lt;/script&gt;" />`
	result, err := f.applyBlockAttributeSyntax(ctx, req, content)

	if err != nil {
		t.Errorf("applyBlockAttributeSyntax() error = %v", err)
		return
	}

	// The wrap attribute value gets double-escaped (once in HTML, once by our code)
	// This prevents XSS by ensuring script tags are never executed
	if !contains(result, "&amp;lt;script&amp;gt;") {
		t.Errorf("applyBlockAttributeSyntax() did not properly escape HTML, got: %q", result)
	}

	// Should NOT contain unescaped script tags
	if contains(result, "<script>") {
		t.Error("applyBlockAttributeSyntax() did not escape XSS attempt - contains unescaped script tag")
	}
}
