package admin

import (
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

func initHandler(store cmsstore.StoreInterface) (func(w http.ResponseWriter, r *http.Request) string, error) {
	return NewPageUpdateController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(nil, nil)),
		Layout: shared.Layout,
	})).Handler, nil
}

func Test_PageUpdateController_PageIdIsRequired(t *testing.T) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	handler, err := initHandler(store)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`"status":"error"`,
		`"message":"page ID is required"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_PageUpdateController_PageIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	handler, err := initHandler(store)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {"invalid"},
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`"status":"error"`,
		`"message":"page not found"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_PageUpdateController_ViewContent_IsDefault(t *testing.T) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	handler, err := initHandler(store)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	// seededSite, err := testutils.SeedSite(store, testutils.SITE_01)

	// if err != nil {
	// 	t.Fatalf("Expected no error, got: %s", err)
	// }

	seededPage, err := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`name="page_title"`,
		`name="page_content"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_PageUpdateController_ViewSettings(t *testing.T) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	handler, err := initHandler(store)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	// seededSite, err := testutils.SeedSite(store, testutils.SITE_01)

	// if err != nil {
	// 	t.Fatalf("Expected no error, got: %s", err)
	// }

	seededPage, err := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"view":    {VIEW_SETTINGS},
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`name="page_editor"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_PageUpdateController_ViewMiddlewares(t *testing.T) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	handler, err := initHandler(store)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	// seededSite, err := testutils.SeedSite(store, testutils.SITE_01)

	// if err != nil {
	// 	t.Fatalf("Expected no error, got: %s", err)
	// }

	seededPage, err := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"view":    {VIEW_MIDDLEWARES},
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`value="middlewares"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_PageUpdateController_ViewSEO(t *testing.T) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	handler, err := initHandler(store)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	// seededSite, err := testutils.SeedSite(store, testutils.SITE_01)

	// if err != nil {
	// 	t.Fatalf("Expected no error, got: %s", err)
	// }

	seededPage, err := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)

	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"view":    {VIEW_SEO},
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`name="page_alias"`,
		`name="page_meta_keywords"`,
		`name="page_meta_description"`,
		`name="page_meta_robots"`,
		`name="page_canonical_url"`,
		`name="page_meta_description"`,
		`name="page_id"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}
