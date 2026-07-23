package page_update

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/test"
)

func Test_AjaxLoadSettings_Success(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionLoadSettings},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"success"`) {
		t.Fatalf("Expected success status, got: %s", body)
	}
	if !strings.Contains(body, `"sites"`) {
		t.Fatalf("Expected sites in response, got: %s", body)
	}
}

func Test_AjaxSaveSettings_Success(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionSaveSettings},
		},
		JSONData: map[string]any{
			"page_id":         seededPage.ID(),
			"page_status":     "active",
			"page_name":       "Updated Name",
			"page_editor":     "codemirror",
			"page_template_id": "",
			"page_site_id":    "",
			"page_memo":       "Test memo",
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"success"`) {
		t.Fatalf("Expected success status, got: %s", body)
	}
}

func Test_AjaxSaveSettings_MissingStatus(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionSaveSettings},
		},
		JSONData: map[string]any{
			"page_id":     seededPage.ID(),
			"page_status": "",
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"error"`) {
		t.Fatalf("Expected error status, got: %s", body)
	}
	if !strings.Contains(body, "Status is required") {
		t.Fatalf("Expected Status is required, got: %s", body)
	}
}
