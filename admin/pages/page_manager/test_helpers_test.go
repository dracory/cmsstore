package page_manager

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	_ "modernc.org/sqlite"
)

// testUi is a minimal implementation of shared.UiInterface for testing.
type testUi struct {
	store  cmsstore.StoreInterface
	logger *slog.Logger
	layout func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}

func (t testUi) Store() cmsstore.StoreInterface { return t.store }
func (t testUi) Logger() *slog.Logger           { return t.logger }
func (t testUi) Layout(w http.ResponseWriter, r *http.Request, title, html string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	return html
}

func newTestUi(store cmsstore.StoreInterface) shared.UiInterface {
	return testUi{
		store:  store,
		logger: slog.New(slog.NewTextHandler(nil, nil)),
		layout: nil,
	}
}

func initStore(t *testing.T) cmsstore.StoreInterface {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}
	return store
}

func initHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewPageManagerController(newTestUi(store)).Handler
}

// callHandler calls the handler and returns the body string.
func callHandler(method string, handler func(w http.ResponseWriter, r *http.Request) string, opts test.NewRequestOptions) (string, error) {
	body, _, err := test.CallStringEndpoint(method, handler, opts)
	return body, err
}
