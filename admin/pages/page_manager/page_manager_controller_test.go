package page_manager

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/test"
)

func Test_Handler_RenderPage(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Page Manager") {
		t.Errorf("Expected body to contain 'Page Manager'")
	}
}

func Test_Handler_AjaxGetRejected(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionLoadPages},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "Method not allowed") {
		t.Errorf("Expected 'Method not allowed' for GET on ajax action, got: %s", body)
	}
}

func Test_Handler_NilStore(t *testing.T) {
	handler := initHandler(nil)

	body, _, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "Store not available") {
		t.Errorf("Expected 'Store not available', got: %s", body)
	}
}
