package page_update

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/test"
)

func Test_AjaxLoadMiddlewares_Success(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionLoadMiddlewares},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"success"`) {
		t.Fatalf("Expected success status, got: %s", body)
	}
}

func Test_AjaxSaveMiddlewares_Success(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, _, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"action":  {actionSaveMiddlewares},
		},
		JSONData: map[string]any{
			"page_id":             seededPage.ID(),
			"middlewares_before":  []string{"auth", "cors"},
			"middlewares_after":   []string{"cache"},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	if !strings.Contains(body, `"status":"success"`) {
		t.Fatalf("Expected success status, got: %s", body)
	}
}
