package admin

import (
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/base/test"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/cmsstore/testutils"
)

func initUI(store cmsstore.StoreInterface) UiInterface {
	return UI(shared.UiConfig{
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			return "" // Placeholder layout function
		},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})
}

func initHandler() (func(w http.ResponseWriter, r *http.Request) string, cmsstore.StoreInterface, error) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		return nil, nil, err
	}

	return NewTemplateUpdateController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(nil, nil)),
		Layout: shared.Layout,
	})).Handler, store, nil
}

func Test_TemplateUpdateController_Index_RequiresTemplateID(t *testing.T) {
	handler, _, err := initHandler()

	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{})

	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`"status":"error"`,
		`"message":"template id is required"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_TemplateUpdateController_Index_RequiresValidTemplate(t *testing.T) {
	handler, _, err := initHandler()

	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"template_id": {"tpl-123"},
		},
	})

	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`"status":"error"`,
		`"message":"template not found"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_TemplateUpdateController_Index_Success(t *testing.T) {
	handler, store, err := initHandler()

	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	_, err = testutils.SeedTemplate(store, testutils.SITE_01, testutils.TEMPLATE_01)

	if err != nil {
		t.Fatalf("Failed to seed template: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"template_id": {testutils.TEMPLATE_01},
			"view":        {VIEW_SETTINGS},
		},
	})

	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`name="template_name"`,
		`name="template_memo"`,
		`name="template_site_id"`,
		`name="template_status"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}
