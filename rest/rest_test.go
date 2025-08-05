package rest_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/rest" // Import the package to be tested
	"github.com/gouniverse/utils"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// initTestDB creates and returns a new in-memory SQLite database connection.
func initTestDB(t *testing.T, filepath string) (*sql.DB, func()) {
	t.Helper()
	if filepath != ":memory:" && utils.FileExists(filepath) {
		err := os.Remove(filepath)
		if err != nil {
			t.Fatalf("failed to remove existing db file: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	if err = db.Ping(); err != nil {
		db.Close()
		t.Fatalf("failed to ping database: %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}
	return db, cleanup
}

// initTestStore initializes a cmsstore.StoreInterface with the given database connection.
func initTestStore(t *testing.T, db *sql.DB) cmsstore.StoreInterface {
	t.Helper()
	store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                db,
		BlockTableName:    "rest_test_block",
		PageTableName:     "rest_test_page",
		SiteTableName:     "rest_test_site",
		TemplateTableName: "rest_test_template",

		TranslationsEnabled:        true,
		TranslationLanguageDefault: "en",
		TranslationLanguages: map[string]string{
			"en": "English",
			"fr": "French",
		},
		TranslationTableName: "rest_test_translation",

		MenusEnabled:      true,
		MenuTableName:     "rest_test_menu",
		MenuItemTableName: "rest_test_menu_item",

		VersioningEnabled:   true,
		VersioningTableName: "rest_test_version",

		AutomigrateEnabled: true,
		DbDriverName:       "sqlite3",
	})

	if err != nil {
		t.Fatalf("cmsstore.NewStore failed: %v", err)
	}
	return store
}

// CreateTestSite creates a new test site and returns it along with a cleanup function
func CreateTestSite(t *testing.T, store cmsstore.StoreInterface) (cmsstore.SiteInterface, func()) {
	t.Helper()
	
	site := cmsstore.NewSite()
	site.SetName("Test Site - " + t.Name())
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	return site, func() {
		_ = store.SiteDeleteByID(context.Background(), site.ID())
	}
}

// setupTestAPI sets up the RestAPI with an in-memory store and returns an httptest.Server URL and a cleanup func.
func setupTestAPI(t *testing.T) (serverURL string, store cmsstore.StoreInterface, cleanup func()) {
	t.Helper()

	db, dbCleanup := initTestDB(t, ":memory:")
	testStore := initTestStore(t, db)

	api := rest.NewRestAPI(testStore)
	testServer := httptest.NewServer(api.Handler())

	cleanupFunc := func() {
		testServer.Close()
		dbCleanup()
	}

	return testServer.URL, testStore, cleanupFunc
}

func TestRestAPI_Routing(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	client := &http.Client{}

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string // Substring to find in body
	}{
		{
			name:       "invalid api path - too short",
			method:     http.MethodGet,
			path:       "/api",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid API path",
		},
		{
			name:       "not an api request",
			method:     http.MethodGet,
			path:       "/nonapi/pages",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Not an API request",
		},
		{
			name:       "unknown resource",
			method:     http.MethodGet,
			path:       "/api/widgets",
			wantStatus: http.StatusNotFound,
			wantBody:   "Unknown resource",
		},
		{
			name:       "pages endpoint - method not allowed",
			method:     http.MethodPatch,
			path:       "/api/pages",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed",
		},
		{
			name:       "menus endpoint - method not allowed",
			method:     http.MethodDelete,
			path:       "/api/menus",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, serverURL+tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, resp.StatusCode, string(body))
			}
			if tt.wantBody != "" && !strings.Contains(string(body), tt.wantBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.wantBody, string(body))
			}
		})
	}
}

func TestRestAPI_PageCreate(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()
	client := &http.Client{}

	testSite := cmsstore.NewSite().SetID("site1").SetName("Test Site")
	if err := store.SiteCreate(context.Background(), testSite); err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	tests := []struct {
		name          string
		payload       map[string]interface{}
		wantStatus    int
		wantInBody    []string
		checkStore    bool
		expectedTitle string
	}{
		{
			name: "valid page creation",
			payload: map[string]interface{}{
				"title":   "My New Page",
				"content": "Page content here.",
				"status":  "published",
				"site_id": "site1",
			},
			wantStatus:    http.StatusOK,
			wantInBody:    []string{`"success":true`, `"title":"My New Page"`},
			checkStore:    true,
			expectedTitle: "My New Page",
		},
		{
			name: "missing title",
			payload: map[string]interface{}{
				"content": "Page content here.",
				"status":  "draft",
				"site_id": "site1",
			},
			wantStatus: http.StatusBadRequest,
			wantInBody: []string{`"success":false`, "Title is required"},
		},
		{
			name: "missing site_id",
			payload: map[string]interface{}{
				"title":   "Page with no site",
				"content": "Content.",
			},
			wantStatus: http.StatusBadRequest,
			wantInBody: []string{`"success":false`, "Site ID is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, serverURL+"/api/pages", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			bodyStr := string(body)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, resp.StatusCode, bodyStr)
			}

			for _, sub := range tt.wantInBody {
				if !strings.Contains(bodyStr, sub) {
					t.Errorf("Expected body to contain %q, got %q", sub, bodyStr)
				}
			}

			if tt.checkStore && resp.StatusCode == http.StatusOK {
				var respJSON map[string]interface{}
				if err := json.Unmarshal(body, &respJSON); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				pageID, ok := respJSON["id"].(string)
				if !ok || pageID == "" {
					t.Fatal("Response did not contain a valid page ID")
				}

				createdPage, errStore := store.PageFindByID(context.Background(), pageID)
				if errStore != nil || createdPage == nil {
					t.Fatalf("Created page not found in store: %v", errStore)
				}
				if createdPage.Title() != tt.expectedTitle {
					t.Errorf("Expected page title in store to be %q, got %q", tt.expectedTitle, createdPage.Title())
				}
			}
		})
	}
}

func TestRestAPI_PageGetListUpdateDelete(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()
	client := &http.Client{}

	testSite := cmsstore.NewSite().SetID("site-for-pages").SetName("Site For Pages")
	if err := store.SiteCreate(context.Background(), testSite); err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	createPayload := map[string]interface{}{"title": "Initial Page", "content": "Initial content.", "status": "draft", "site_id": "site-for-pages"}
	payloadBytes, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest(http.MethodPost, serverURL+"/api/pages", bytes.NewBuffer(payloadBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := client.Do(createReq)
	var createdPageResp map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createdPageResp)
	pageID := createdPageResp["id"].(string)
	createResp.Body.Close()

	t.Run("Get Existing Page", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, serverURL+"/api/pages/"+pageID, nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		var pageResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&pageResp)
		if pageResp["id"] != pageID || pageResp["title"] != "Initial Page" {
			t.Errorf("Unexpected page data: %+v", pageResp)
		}
	})

	t.Run("List Pages", func(t *testing.T) {
		listReq, _ := http.NewRequest(http.MethodGet, serverURL+"/api/pages", nil)
		listResp, _ := client.Do(listReq)
		defer listResp.Body.Close()
		var listResult struct {
			Pages []map[string]interface{} `json:"pages"`
		}
		json.NewDecoder(listResp.Body).Decode(&listResult)
		if len(listResult.Pages) < 1 {
			t.Errorf("Expected at least 1 page, got %d", len(listResult.Pages))
		}
	})

	t.Run("Update Page", func(t *testing.T) {
		updatePl := map[string]interface{}{"title": "Updated Page", "status": "published"}
		plBytes, _ := json.Marshal(updatePl)
		req, _ := http.NewRequest(http.MethodPut, serverURL+"/api/pages/"+pageID, bytes.NewBuffer(plBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		var updatedResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&updatedResp)
		if updatedResp["title"] != "Updated Page" {
			t.Error("Title not updated")
		}
	})

	t.Run("Delete Page", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, serverURL+"/api/pages/"+pageID, nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		p, _ := store.PageFindByID(context.Background(), pageID)
		if p != nil {
			t.Error("Page not soft-deleted")
		}
	})

	t.Run("Get Non-Existent Page", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, serverURL+"/api/pages/nonexistent", nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, resp.StatusCode, string(body))
		}
	})

	t.Run("Update Non-Existent Page", func(t *testing.T) {
		updatePl := map[string]interface{}{"title": "No Such Page"}
		plBytes, _ := json.Marshal(updatePl)
		req, _ := http.NewRequest(http.MethodPut, serverURL+"/api/pages/nonexistent", bytes.NewBuffer(plBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, resp.StatusCode, string(body))
		}
	})

	t.Run("Delete Non-Existent Page", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, serverURL+"/api/pages/nonexistent", nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		// Based on current store logic, deleting a non-existent page results in an error from PageSoftDeleteByID,
		// which the handler turns into a 500.
		if resp.StatusCode != http.StatusInternalServerError {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusInternalServerError, resp.StatusCode, string(body))
		}
	})
}

func TestRestAPI_MenuCreate(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()
	client := &http.Client{}

	testSite := cmsstore.NewSite().SetID("site-for-menus").SetName("Site For Menus")
	if err := store.SiteCreate(context.Background(), testSite); err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	tests := []struct {
		name         string
		payload      map[string]interface{}
		wantStatus   int
		wantInBody   []string
		checkStore   bool
		expectedName string
	}{
		{
			name: "valid menu creation",
			payload: map[string]interface{}{
				"name":    "Main Menu",
				"site_id": "site-for-menus",
			},
			wantStatus:   http.StatusOK,
			wantInBody:   []string{`"success":true`, `"name":"Main Menu"`},
			checkStore:   true,
			expectedName: "Main Menu",
		},
		{
			name: "missing name",
			payload: map[string]interface{}{
				"site_id": "site-for-menus",
			},
			wantStatus: http.StatusBadRequest,
			wantInBody: []string{`"success":false`, "Name is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, serverURL+"/api/menus", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := client.Do(req)
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, resp.StatusCode, string(body))
			}
			for _, sub := range tt.wantInBody {
				if !strings.Contains(string(body), sub) {
					t.Errorf("Expected body to contain %q, got %q", sub, string(body))
				}
			}
			if tt.checkStore && resp.StatusCode == http.StatusOK {
				var respJSON map[string]interface{}
				json.Unmarshal(body, &respJSON)
				menuID := respJSON["id"].(string)
				m, _ := store.MenuFindByID(context.Background(), menuID)
				if m == nil || m.Name() != tt.expectedName {
					t.Errorf("Menu not found in store or name mismatch. Expected: %s, Got: %v", tt.expectedName, m)
				}
			}
		})
	}
}

func TestRestAPI_MenuGetList(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()
	client := &http.Client{}

	testSite := cmsstore.NewSite().SetID("s1").SetName("S1")
	store.SiteCreate(context.Background(), testSite)

	// Create a menu
	createPl := map[string]interface{}{"name": "Nav Menu", "site_id": "s1"}
	plBytes, _ := json.Marshal(createPl)
	createReq, _ := http.NewRequest(http.MethodPost, serverURL+"/api/menus", bytes.NewBuffer(plBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := client.Do(createReq)
	var createdMenuResp map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createdMenuResp)
	menuID := createdMenuResp["id"].(string)
	createResp.Body.Close()

	t.Run("Get Existing Menu", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, serverURL+"/api/menus/"+menuID, nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		var menuResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&menuResp)
		if menuResp["id"] != menuID || menuResp["name"] != "Nav Menu" {
			t.Errorf("Unexpected menu data: %+v", menuResp)
		}
	})

	t.Run("List Menus", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, serverURL+"/api/menus", nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		var listResp struct {
			Menus []map[string]interface{} `json:"menus"`
		}
		json.NewDecoder(resp.Body).Decode(&listResp)
		if len(listResp.Menus) < 1 {
			t.Errorf("Expected at least 1 menu, got %d", len(listResp.Menus))
		}
	})

	t.Run("Get Non-Existent Menu", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, serverURL+"/api/menus/nonexistent", nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNotFound, resp.StatusCode, string(body))
		}
	})
}
