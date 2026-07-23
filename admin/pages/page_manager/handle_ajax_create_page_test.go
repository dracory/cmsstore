package page_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
)

func Test_HandleAjaxCreatePage_Success(t *testing.T) {
	store := initStore(t)

	site, err := testutils.SeedSite(store, testutils.SITE_01)
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionCreatePage},
		},
		JSONData: map[string]any{
			"name":    "Test Page",
			"site_id": site.ID(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "success") {
		t.Errorf("Expected success response, got: %s", body)
	}
}

func Test_HandleAjaxCreatePage_MissingName(t *testing.T) {
	store := initStore(t)

	site, err := testutils.SeedSite(store, testutils.SITE_01)
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionCreatePage},
		},
		JSONData: map[string]any{
			"name":    "",
			"site_id": site.ID(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "Name is required") {
		t.Errorf("Expected 'Name is required', got: %s", body)
	}
}

func Test_HandleAjaxCreatePage_MissingSiteID(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionCreatePage},
		},
		JSONData: map[string]any{
			"name":    "Test Page",
			"site_id": "",
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "Site is required") {
		t.Errorf("Expected 'Site is required', got: %s", body)
	}
}

func Test_HandleAjaxCreatePage_InvalidBody(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionCreatePage},
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
