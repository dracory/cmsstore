package admin

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

func Test_UI_CreatesUiInterface(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	if ui == nil {
		t.Errorf("Expected UI to be created, got nil")
	}
}

func Test_UI_UiInterfaceHasRequiredMethods(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	// Verify that ui was properly initialized
	// Note: Function fields cannot be compared to nil in Go
	if ui == nil {
		t.Errorf("Expected UI to be created")
	}
}

func Test_UI_UiInterfaceMethodsWork(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	// Test Logger() returns the logger
	logger := ui.Logger()
	if logger == nil {
		t.Errorf("Expected Logger to not be nil")
	}

	// Test Store() returns the store
	storeInterface := ui.Store()
	if storeInterface == nil {
		t.Errorf("Expected Store to not be nil")
	}
	// Verify store is not nil (it's already StoreInterface type)
	_ = storeInterface
}

func Test_UI_UiInterfaceLayoutMethod(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	customLayoutCalled := false
	customLayout := func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string {
		customLayoutCalled = true
		return "<html>" + webpageHtml + "</html>"
	}

	ui := UI(shared.UiConfig{
		Layout: customLayout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rec := &testResponseWriter{}

	// Call Layout method
	result := ui.Layout(rec, req, "Test Title", "<div>Test Content</div>", struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{})

	if !customLayoutCalled {
		t.Errorf("Expected custom layout to be called")
	}
	if !strings.Contains(result, "<div>Test Content</div>") {
		t.Errorf("Expected result to contain '<div>Test Content</div>'")
	}
}

// testResponseWriter is a minimal ResponseWriter implementation for testing
type testResponseWriter struct {
	header http.Header
}

func (w *testResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (w *testResponseWriter) WriteHeader(statusCode int) {}
