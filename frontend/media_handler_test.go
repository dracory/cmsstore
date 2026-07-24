package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
)

func TestExtractMediaID(t *testing.T) {
	tests := []struct {
		name     string
		urlPath  string
		expected string
	}{
		{"short form png", "/cms/media/1jq6fby3kzj.png", "1jq6fby3kzj"},
		{"short form jpg", "/cms/media/abc123.jpg", "abc123"},
		{"short form no extension", "/cms/media/abc123", "abc123"},
		{"readable form", "/cms/media/1jq6fby3kzj/dulydo.png", "1jq6fby3kzj"},
		{"readable form with underscores", "/cms/media/abc123/my_image_name.jpg", "abc123"},
		{"trailing slash", "/cms/media/1jq6fby3kzj/", "1jq6fby3kzj"},
		{"double slash", "/cms/media//1jq6fby3kzj.png", "1jq6fby3kzj"},
		{"empty path", "/cms/media/", ""},
		{"only extension", "/cms/media/.png", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMediaID(tt.urlPath)
			if got != tt.expected {
				t.Errorf("extractMediaID(%q) = %q, want %q", tt.urlPath, got, tt.expected)
			}
		})
	}
}

func TestIsMediaURL(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/cms/media/abc123.png", true},
		{"/cms/media/abc123", true},
		{"/cms/media/abc123/handle.png", true},
		{"/cms/media/", true},
		{"/page/about", false},
		{"/", false},
		{"/blog/post", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isMediaURL(tt.path)
			if got != tt.expected {
				t.Errorf("isMediaURL(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestMimeTypeFromExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".png", "image/png"},
		{".jpg", "image/jpeg"},
		{".jpeg", "image/jpeg"},
		{".gif", "image/gif"},
		{".webp", "image/webp"},
		{".svg", "image/svg+xml"},
		{".ico", "image/x-icon"},
		{".pdf", "application/pdf"},
		{".unknown", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := mimeTypeFromExtension(tt.ext)
			if got != tt.expected {
				t.Errorf("mimeTypeFromExtension(%q) = %q, want %q", tt.ext, got, tt.expected)
			}
		})
	}
}

func TestMediaHandler_DataURI(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	media := cmsstore.NewMedia().
		SetID("test123").
		SetExtension(".png").
		SetType("image/png").
		SetSize("68").
		SetURL("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==").
		SetStatus(cmsstore.MEDIA_STATUS_ACTIVE)

	if err := store.MediaCreate(nil, media); err != nil {
		t.Fatalf("Failed to create media: %v", err)
	}

	f := New(Config{Store: store})

	req := httptest.NewRequest("GET", "/cms/media/test123.png", nil)
	recorder := httptest.NewRecorder()

	f.Handler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", recorder.Header().Get("Content-Type"))
	}

	if recorder.Header().Get("Cache-Control") != "public, max-age=31536000" {
		t.Errorf("expected Cache-Control header, got %s", recorder.Header().Get("Cache-Control"))
	}

	if recorder.Header().Get("ETag") != "test123" {
		t.Errorf("expected ETag test123, got %s", recorder.Header().Get("ETag"))
	}

	if recorder.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestMediaHandler_NotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{Store: store})

	req := httptest.NewRequest("GET", "/cms/media/nonexistent.png", nil)
	recorder := httptest.NewRecorder()

	f.Handler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestMediaHandler_InactiveMedia(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	media := cmsstore.NewMedia().
		SetID("inactive123").
		SetExtension(".png").
		SetType("image/png").
		SetURL("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==").
		SetStatus(cmsstore.MEDIA_STATUS_DRAFT)

	if err := store.MediaCreate(nil, media); err != nil {
		t.Fatalf("Failed to create media: %v", err)
	}

	f := New(Config{Store: store})

	req := httptest.NewRequest("GET", "/cms/media/inactive123.png", nil)
	recorder := httptest.NewRecorder()

	f.Handler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status %d for inactive media, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestMediaHandler_HTTPRedirect(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	media := cmsstore.NewMedia().
		SetID("redirect123").
		SetExtension(".png").
		SetType("image/png").
		SetURL("https://example.com/image.png").
		SetStatus(cmsstore.MEDIA_STATUS_ACTIVE)

	if err := store.MediaCreate(nil, media); err != nil {
		t.Fatalf("Failed to create media: %v", err)
	}

	f := New(Config{Store: store})

	req := httptest.NewRequest("GET", "/cms/media/redirect123.png", nil)
	recorder := httptest.NewRecorder()

	f.Handler(recorder, req)

	if recorder.Code != http.StatusFound {
		t.Errorf("expected status %d for HTTP redirect, got %d", http.StatusFound, recorder.Code)
	}

	if recorder.Header().Get("Location") != "https://example.com/image.png" {
		t.Errorf("expected redirect to https://example.com/image.png, got %s", recorder.Header().Get("Location"))
	}
}

func TestStringHandler_MediaURL(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	media := cmsstore.NewMedia().
		SetID("strhandler123").
		SetExtension(".png").
		SetType("image/png").
		SetSize("68").
		SetURL("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==").
		SetStatus(cmsstore.MEDIA_STATUS_ACTIVE)

	if err := store.MediaCreate(nil, media); err != nil {
		t.Fatalf("Failed to create media: %v", err)
	}

	f := New(Config{Store: store})

	req := httptest.NewRequest("GET", "/cms/media/strhandler123.png", nil)
	recorder := httptest.NewRecorder()

	f.Handler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", recorder.Header().Get("Content-Type"))
	}
}
