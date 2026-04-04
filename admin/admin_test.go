package admin

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

func newAdminForTest(t *testing.T, options AdminOptions) *admin {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	if options.Store == nil {
		options.Store = store
	}
	if options.Logger == nil {
		options.Logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	a, err := New(options)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}
	return a
}

func TestAdminHandle_HomePath(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Dashboard") {
		t.Errorf("Expected body to contain 'Dashboard'")
	}
}

func TestAdminHandle_InvalidPathDefaultsToHome(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=nonexistent", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Dashboard") {
		t.Errorf("Expected body to contain 'Dashboard'")
	}
}

func TestAdminHandle_PageManagerPath(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=/pages/page-manager", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Pages") {
		t.Errorf("Expected body to contain 'Pages'")
	}
}

func TestAdminHandle_MediaLinkRenderedWhenConfigured(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{
		MediaManagerURL: "/admin/media",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "href=\"/admin/media\"") {
		t.Errorf("Expected body to contain media link")
	}
	if !strings.Contains(w.Body.String(), "target=\"_blank\"") {
		t.Errorf("Expected body to contain target blank")
	}
	if !strings.Contains(w.Body.String(), ">Media</a>") {
		t.Errorf("Expected body to contain Media link")
	}
}

func TestAdminHandle_MediaLinkNotRenderedWhenNotConfigured(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if strings.Contains(w.Body.String(), ">Media</a>") {
		t.Errorf("Expected body to NOT contain Media link")
	}
}

func TestAdminHandle_CMSContainerPaddingRendered(t *testing.T) {
	a := newAdminForTest(t, AdminOptions{
		PaddingTopPx:    9,
		PaddingRightPx:  8,
		PaddingBottomPx: 7,
		PaddingLeftPx:   6,
	})

	req := httptest.NewRequest(http.MethodGet, "/admin?path=home", nil)
	w := httptest.NewRecorder()
	a.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "class=\"cms\"") {
		t.Errorf("Expected body to contain cms class")
	}
	if !strings.Contains(w.Body.String(), "padding: 9px 8px 7px 6px;") {
		t.Errorf("Expected body to contain padding style")
	}
}
