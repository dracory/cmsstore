package page_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
)

func Test_HandleAjaxDeletePage_Success(t *testing.T) {
	store := initStore(t)

	site, err := testutils.SeedSite(store, testutils.SITE_01)
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	page, err := testutils.SeedPage(store, site.ID(), testutils.PAGE_01)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionDeletePage},
		},
		JSONData: map[string]any{
			"page_id": page.ID(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "success") {
		t.Errorf("Expected success response, got: %s", body)
	}
}

func Test_HandleAjaxDeletePage_MissingID(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionDeletePage},
		},
		JSONData: map[string]any{
			"page_id": "",
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "Page ID is required") {
		t.Errorf("Expected 'Page ID is required', got: %s", body)
	}
}

func Test_HandleAjaxDeletePage_NotFound(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionDeletePage},
		},
		JSONData: map[string]any{
			"page_id": "nonexistent-id",
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "Page not found") {
		t.Errorf("Expected 'Page not found', got: %s", body)
	}
}

func Test_HandleAjaxDeletePage_InvalidBody(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionDeletePage},
		},
		Body:        `invalid json`,
		ContentType: "application/json",
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "Invalid request body") {
		t.Errorf("Expected 'Invalid request body', got: %s", body)
	}
}
