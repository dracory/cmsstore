package page_update

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/dracory/blockeditor"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	_ "modernc.org/sqlite"
)

type testUi struct {
	store                  cmsstore.StoreInterface
	logger                 *slog.Logger
	blockEditorDefinitions []blockeditor.BlockDefinition
	layoutFunc             func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
}

func (t testUi) Store() cmsstore.StoreInterface { return t.store }
func (t testUi) Logger() *slog.Logger           { return t.logger }
func (t testUi) BlockEditorDefinitions() []blockeditor.BlockDefinition {
	return t.blockEditorDefinitions
}
func (t testUi) Layout(w http.ResponseWriter, r *http.Request, title, html string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	if t.layoutFunc != nil {
		return t.layoutFunc(w, r, title, html, options)
	}
	return html
}

func newTestUi(store cmsstore.StoreInterface) uiInterface {
	return testUi{
		store:  store,
		logger: slog.New(slog.NewTextHandler(nil, nil)),
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
	return NewPageUpdateController(newTestUi(store)).Handler
}

func callHandler(method string, handler func(w http.ResponseWriter, r *http.Request) string, opts test.NewRequestOptions) (string, error) {
	body, _, err := test.CallStringEndpoint(method, handler, opts)
	return body, err
}
