package admin

import (
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_UI_CreatesUiInterface(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	assert.NotNil(t, ui)
}

func Test_UI_UiInterfaceHasRequiredMethods(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	// Verify that ui implements UiInterface (excluding Endpoint which is not in the interface)
	assert.NotNil(t, ui.Layout)
	assert.NotNil(t, ui.Logger)
	assert.NotNil(t, ui.Store)
	assert.NotNil(t, ui.TranslationCreate)
	assert.NotNil(t, ui.TranslationManager)
	assert.NotNil(t, ui.TranslationDelete)
	assert.NotNil(t, ui.TranslationUpdate)
	assert.NotNil(t, ui.TranslationVersioning)
}

func Test_UI_UiInterfaceMethodsWork(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	// Test Logger() returns the logger
	logger := ui.Logger()
	assert.NotNil(t, logger)

	// Test Store() returns the store
	storeInterface := ui.Store()
	assert.NotNil(t, storeInterface)
	assert.Implements(t, (*cmsstore.StoreInterface)(nil), storeInterface)
}

func Test_UI_UiInterfaceLayoutMethod(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

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

	assert.True(t, customLayoutCalled)
	assert.Contains(t, result, "<div>Test Content</div>")
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
