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

func TestApplyTranslationAttributeSyntax(t *testing.T) {
	ctx := context.Background()

	// Initialize test store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create test translation
	translation := cmsstore.NewTranslation().
		SetHandle("welcome").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	// Set translation content with interpolation placeholders
	translationContent := map[string]string{
		"en": "Welcome, {{name}}!",
		"es": "¡Bienvenido, {{name}}!",
		"fr": "Bienvenue, {{name}}!",
	}
	err = translation.SetContent(translationContent)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create translation with fallback scenario
	translationRare := cmsstore.NewTranslation().
		SetHandle("rare_term").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	rareContent := map[string]string{
		"en": "Rare technical term",
		// Missing other languages
	}
	err = translationRare.SetContent(rareContent)
	if err != nil {
		t.Fatalf("Failed to set rare translation content: %v", err)
	}

	err = store.TranslationCreate(ctx, translationRare)
	if err != nil {
		t.Fatalf("Failed to create rare translation: %v", err)
	}

	// Create frontend instance
	fe := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: false,
	})
	frontend := fe.(*frontend)

	tests := []struct {
		name     string
		content  string
		language string
		expected string
	}{
		{
			name:     "basic translation - angle brackets",
			content:  `<translation id="welcome" name="John" />`,
			language: "en",
			expected: "Welcome, John!",
		},
		{
			name:     "basic translation - square brackets",
			content:  `[[translation id='welcome' name='Mary']]`,
			language: "en",
			expected: "Welcome, Mary!",
		},
		{
			name:     "translation with different language",
			content:  `<translation id="welcome" name="Juan" />`,
			language: "es",
			expected: "¡Bienvenido, Juan!",
		},
		{
			name:     "translation with fallback",
			content:  `<translation id="rare_term" fallback="en" />`,
			language: "fr",
			expected: "Rare technical term",
		},
		{
			name:     "multiple interpolation variables",
			content:  `<translation id="welcome" name="Alice" />`,
			language: "en",
			expected: "Welcome, Alice!",
		},
		{
			name:     "translation in HTML context",
			content:  `<p><translation id="welcome" name="Bob" /></p>`,
			language: "en",
			expected: "<p>Welcome, Bob!</p>",
		},
		{
			name:     "missing id attribute",
			content:  `<translation name="John" />`,
			language: "en",
			expected: "<!-- Translation reference missing id -->",
		},
		{
			name:     "non-existent translation",
			content:  `<translation id="nonexistent" />`,
			language: "en",
			expected: "<!-- Translation not found: nonexistent -->",
		},
		{
			name:     "XSS prevention in interpolation",
			content:  `<translation id="welcome" name="&lt;script&gt;alert('xss')&lt;/script&gt;" />`,
			language: "en",
			expected: "Welcome, &amp;lt;script&amp;gt;alert(&#39;xss&#39;)&amp;lt;/script&amp;gt;!",
		},
		{
			name:     "multiple translations in content",
			content:  `<translation id="welcome" name="Alice" /> and <translation id="welcome" name="Bob" />`,
			language: "en",
			expected: "Welcome, Alice! and Welcome, Bob!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			result, err := frontend.applyTranslationAttributeSyntax(req, tt.content, tt.language)

			if err != nil {
				t.Errorf("applyTranslationAttributeSyntax() error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("applyTranslationAttributeSyntax() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTranslationAttributeSyntaxIntegration(t *testing.T) {
	ctx := context.Background()

	// Initialize test store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create test translation
	translation := cmsstore.NewTranslation().
		SetHandle("order_summary").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	translationContent := map[string]string{
		"en": "You have {{count}} items totaling {{total}}",
		"es": "Tienes {{count}} artículos por un total de {{total}}",
	}
	err = translation.SetContent(translationContent)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create frontend instance
	fe := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: false,
	})
	frontend := fe.(*frontend)

	// Test full rendering pipeline
	content := `<h1>Order Confirmation</h1>
<p><translation id="order_summary" count="3" total="$125.50" /></p>`

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	options := TemplateRenderHtmlByIDOptions{
		Language: "en",
	}

	result, err := frontend.renderContentToHtml(req, content, options)
	if err != nil {
		t.Fatalf("renderContentToHtml() error = %v", err)
	}

	expected := `<h1>Order Confirmation</h1>
<p>You have 3 items totaling $125.50</p>`

	if result != expected {
		t.Errorf("renderContentToHtml() = %q, want %q", result, expected)
	}
}

func TestTranslationAttributeSyntaxWithLegacy(t *testing.T) {
	ctx := context.Background()

	// Initialize test store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create test translation
	translation := cmsstore.NewTranslation().
		SetHandle("greeting").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	translationContent := map[string]string{
		"en": "Hello, {{name}}!",
	}
	err = translation.SetContent(translationContent)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create legacy translation
	legacyTranslation := cmsstore.NewTranslation().
		SetHandle("goodbye").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	legacyContent := map[string]string{
		"en": "Goodbye!",
	}
	err = legacyTranslation.SetContent(legacyContent)
	if err != nil {
		t.Fatalf("Failed to set legacy translation content: %v", err)
	}

	err = store.TranslationCreate(ctx, legacyTranslation)
	if err != nil {
		t.Fatalf("Failed to create legacy translation: %v", err)
	}

	// Create frontend instance
	fe := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: false,
	})
	frontend := fe.(*frontend)

	// Test both legacy and new syntax together
	content := `[[TRANSLATION_goodbye]] <translation id="greeting" name="World" />`

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	options := TemplateRenderHtmlByIDOptions{
		Language: "en",
	}

	result, err := frontend.renderContentToHtml(req, content, options)
	if err != nil {
		t.Fatalf("renderContentToHtml() error = %v", err)
	}

	expected := "Goodbye! Hello, World!"

	if result != expected {
		t.Errorf("renderContentToHtml() = %q, want %q", result, expected)
	}
}
