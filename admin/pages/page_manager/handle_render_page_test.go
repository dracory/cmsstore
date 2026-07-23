package page_manager

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dracory/test"
)

func Test_HandleRenderPage_Success(t *testing.T) {
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
	if !strings.Contains(body, "New Page") {
		t.Errorf("Expected body to contain 'New Page'")
	}
}

func Test_HandleRenderPage_ContainsVueApp(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, _, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if !strings.Contains(body, "vue.global.js") {
		t.Errorf("Expected body to contain Vue CDN script")
	}
	if !strings.Contains(body, "CmsPagesApp") {
		t.Errorf("Expected body to contain Vue app initialization")
	}
}
