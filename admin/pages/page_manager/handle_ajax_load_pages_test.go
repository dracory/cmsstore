package page_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
)

func Test_HandleAjaxLoadPages_Success(t *testing.T) {
	store := initStore(t)

	site, err := testutils.SeedSite(store, testutils.SITE_01)
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	_, err = testutils.SeedPage(store, site.ID(), testutils.PAGE_01)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionLoadPages},
		},
		JSONData: map[string]any{
			"page":     0,
			"per_page": 20,
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "success") {
		t.Errorf("Expected success response, got: %s", body)
	}
	if !strings.Contains(body, testutils.PAGE_01) {
		t.Errorf("Expected response to contain page name %s, got: %s", testutils.PAGE_01, body)
	}
}

func Test_HandleAjaxLoadPages_Empty(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionLoadPages},
		},
		JSONData: map[string]any{
			"page":     0,
			"per_page": 20,
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "success") {
		t.Errorf("Expected success response with empty list, got: %s", body)
	}
}

func Test_HandleAjaxLoadPages_InvalidBody(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionLoadPages},
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
