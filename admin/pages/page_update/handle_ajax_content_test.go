package page_update

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/test"
)

func Test_AjaxLoadContent_Success(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionLoadContent},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"success"`) {
		t.Fatalf("Expected success status, got: %s", body)
	}
	if !strings.Contains(body, `"title"`) {
		t.Fatalf("Expected title in response, got: %s", body)
	}
}

func Test_AjaxLoadContent_MissingPageID(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionLoadContent},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"error"`) {
		t.Fatalf("Expected error status, got: %s", body)
	}
	if !strings.Contains(body, "Page ID is required") {
		t.Fatalf("Expected Page ID is required, got: %s", body)
	}
}

func Test_AjaxSaveContent_Success(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionSaveContent},
		},
		JSONData: map[string]any{
			"page_id":      seededPage.ID(),
			"page_title":   "Updated Title",
			"page_content": "Updated content",
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"success"`) {
		t.Fatalf("Expected success status, got: %s", body)
	}

	page, _ := store.PageFindByID(nil, seededPage.ID())
	if page.Title() != "Updated Title" {
		t.Fatalf("Expected title to be 'Updated Title', got: %s", page.Title())
	}
}

func Test_AjaxSaveContent_MissingTitle(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionSaveContent},
		},
		JSONData: map[string]any{
			"page_id":      seededPage.ID(),
			"page_title":   "",
			"page_content": "content",
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"error"`) {
		t.Fatalf("Expected error status, got: %s", body)
	}
	if !strings.Contains(body, "Title is required") {
		t.Fatalf("Expected Title is required, got: %s", body)
	}
}
