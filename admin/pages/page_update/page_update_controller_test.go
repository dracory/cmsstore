package page_update

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/test"
)

func Test_PageUpdate_PageIdIsRequired(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`"status":"error"`,
		`"message":"Page ID is required"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_PageUpdate_PageNotFound(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

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
		`"message":"Page not found"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_PageUpdate_RenderContent(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
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

	if !strings.Contains(body, "page-content-app") {
		t.Fatalf("Expected to find page-content-app in body, got: %s", body)
	}
}

func Test_PageUpdate_RenderSEO(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"view":    {viewSEO},
		},
	})
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "page-seo-app") {
		t.Fatalf("Expected to find page-seo-app in body, got: %s", body)
	}
}

func Test_PageUpdate_RenderSettings(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"view":    {viewSettings},
		},
	})
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "page-settings-app") {
		t.Fatalf("Expected to find page-settings-app in body, got: %s", body)
	}
}

func Test_PageUpdate_RenderMiddlewares(t *testing.T) {
	store := initStore(t)
	handler := initHandler(store)

	seededPage, err := seedTestPage(store)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
			"view":    {viewMiddlewares},
		},
	})
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "page-middlewares-app") {
		t.Fatalf("Expected to find page-middlewares-app in body, got: %s", body)
	}
}
