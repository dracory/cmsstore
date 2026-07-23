package page_update

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/test"
)

func Test_AjaxLoadSEO_Success(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionLoadSEO},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"success"`) {
		t.Fatalf("Expected success status, got: %s", body)
	}
}

func Test_AjaxSaveSEO_Success(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionSaveSEO},
		},
		JSONData: map[string]any{
			"page_id":               seededPage.ID(),
			"page_alias":            "/test-alias",
			"page_meta_description": "Test description",
			"page_meta_keywords":    "test, keywords",
			"page_meta_robots":      "INDEX, FOLLOW",
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"success"`) {
		t.Fatalf("Expected success status, got: %s", body)
	}
}
