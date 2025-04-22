package admin

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/base/test"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/cmsstore/testutils"
)

func initUI(store cmsstore.StoreInterface) UiInterface {
	return UI(shared.UiConfig{
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}) string {
			return "" // Placeholder layout function
		},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})
}

func initHandler() (func(w http.ResponseWriter, r *http.Request) string, cmsstore.StoreInterface, error) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		return nil, nil, err
	}

	return NewTemplateUpdateController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(nil, nil)),
		Layout: shared.Layout,
	})).Handler, store, nil
}

func Test_TemplateUpdateController_Index_RequiresTemplateID(t *testing.T) {
	handler, _, err := initHandler()

	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{})

	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`"status":"error"`,
		`"message":"template id is required"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_TemplateUpdateController_Index_RequiresValidTemplate(t *testing.T) {
	handler, _, err := initHandler()

	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"template_id": {"tpl-123"},
		},
	})

	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`"status":"error"`,
		`"message":"template not found"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_TemplateUpdateController_Index_Success(t *testing.T) {
	handler, store, err := initHandler()

	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	_, err = testutils.SeedTemplate(store, testutils.SITE_01, testutils.TEMPLATE_01)

	if err != nil {
		t.Fatalf("Failed to seed template: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"template_id": {testutils.TEMPLATE_01},
			"view":        {VIEW_SETTINGS},
		},
	})

	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`name="template_name"`,
		`name="template_memo"`,
		`name="template_site_id"`,
		`name="template_status"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

// Test_TemplateUpdateController_Save_Settings_Success tests successful saving of settings view
func Test_TemplateUpdateController_Save_Settings_Success(t *testing.T) {
	handler, store, err := initHandler()
	if err != nil {
		t.Fatalf("initHandler should succeed, got error: %v", err)
	}

	// Seed initial template using testutils
	template, err := testutils.SeedTemplate(store, testutils.SITE_01, testutils.TEMPLATE_01)
	if err != nil {
		t.Fatalf("Seeding template should succeed, got error: %v", err)
	}
	if template == nil {
		t.Fatalf("Seeded template should not be nil")
	}

	newName := "Updated Template Name"
	newMemo := "Updated Memo Content"
	newSiteID := testutils.SITE_02 // Change site ID
	newStatus := cmsstore.TEMPLATE_STATUS_ACTIVE

	postData := url.Values{}
	postData.Set("template_name", newName)
	postData.Set("template_memo", newMemo)
	postData.Set("template_site_id", newSiteID)
	postData.Set("template_status", newStatus)
	postData.Set("view", VIEW_SETTINGS) // Ensure view is set for POST

	// Call the handler using test.CallStringEndpoint
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{ // Pass template_id via GET params as per handler logic
			"template_id": {template.ID()},
		},
		PostValues: postData, // Send form data via POST
	})

	if err != nil {
		t.Fatalf("CallStringEndpoint should succeed, got error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK (%d), got: %d", http.StatusOK, response.StatusCode)
	}

	// Check for success message in the response body (rendered by the form builder)
	if !strings.Contains(body, `"icon":"success"`) {
		t.Errorf("Expected success swal icon in body: %s", body)
	}
	if !strings.Contains(body, `"text":"template saved successfully"`) {
		t.Errorf("Expected success swal text in body: %s", body)
	}
	// Check for redirect script
	if !strings.Contains(body, `window.location.href = "`) {
		t.Errorf("Expected redirect script in body: %s", body)
	}
	if !strings.Contains(body, "view="+VIEW_SETTINGS) {
		t.Errorf("Expected redirect back to settings view in body: %s", body)
	}

	// Verify database state directly using the store
	updatedTemplate, err := store.TemplateFindByID(context.Background(), template.ID())
	if err != nil {
		t.Fatalf("Finding updated template should succeed, got error: %v", err)
	}
	if updatedTemplate == nil {
		t.Fatalf("Updated template should exist in DB")
	}
	if updatedTemplate.Name() != newName {
		t.Errorf("DB name mismatch: expected '%s', got '%s'", newName, updatedTemplate.Name())
	}
	if updatedTemplate.Memo() != newMemo {
		t.Errorf("DB memo mismatch: expected '%s', got '%s'", newMemo, updatedTemplate.Memo())
	}
	if updatedTemplate.SiteID() != newSiteID {
		t.Errorf("DB site ID mismatch: expected '%s', got '%s'", newSiteID, updatedTemplate.SiteID())
	}
	if updatedTemplate.Status() != newStatus {
		t.Errorf("DB status mismatch: expected '%s', got '%s'", newStatus, updatedTemplate.Status())
	}
}

// Test_TemplateUpdateController_Save_Content_Success tests successful saving of content view
func Test_TemplateUpdateController_Save_Content_Success(t *testing.T) {
	handler, store, err := initHandler()
	if err != nil {
		t.Fatalf("initHandler should succeed, got error: %v", err)
	}

	template, err := testutils.SeedTemplate(store, testutils.SITE_01, testutils.TEMPLATE_01)
	if err != nil {
		t.Fatalf("Seeding template should succeed, got error: %v", err)
	}

	newContent := "<div>New HTML Content</div>"

	postData := url.Values{}
	postData.Set("template_content", newContent)
	postData.Set("view", VIEW_CONTENT) // Ensure view is set for POST

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{ // Pass template_id via GET params
			"template_id": {template.ID()},
		},
		PostValues: postData, // Send form data via POST
	})

	if err != nil {
		t.Fatalf("CallStringEndpoint should succeed, got error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK (%d), got: %d", http.StatusOK, response.StatusCode)
	}

	// Check for success message and redirect
	if !strings.Contains(body, `"icon":"success"`) {
		t.Errorf("Expected success swal icon in body: %s", body)
	}
	if !strings.Contains(body, `"text":"template saved successfully"`) {
		t.Errorf("Expected success swal text in body: %s", body)
	}
	if !strings.Contains(body, `window.location.href = "`) {
		t.Errorf("Expected redirect script in body: %s", body)
	}
	if !strings.Contains(body, "view="+VIEW_CONTENT) {
		t.Errorf("Expected redirect back to content view in body: %s", body)
	}

	// Verify database state
	updatedTemplate, err := store.TemplateFindByID(context.Background(), template.ID())
	if err != nil {
		t.Fatalf("Finding updated template should succeed, got error: %v", err)
	}
	if updatedTemplate == nil {
		t.Fatalf("Updated template should exist in DB")
	}
	if updatedTemplate.Content() != newContent {
		t.Errorf("DB content mismatch: expected '%s', got '%s'", newContent, updatedTemplate.Content())
	}
	// Ensure other fields didn't change
	if updatedTemplate.Name() != template.Name() {
		t.Errorf("DB name should not change: expected '%s', got '%s'", template.Name(), updatedTemplate.Name())
	}
	if updatedTemplate.SiteID() != template.SiteID() {
		t.Errorf("DB site ID should not change: expected '%s', got '%s'", template.SiteID(), updatedTemplate.SiteID())
	}
}

// Test_TemplateUpdateController_Save_Settings_MissingStatus tests validation failure
func Test_TemplateUpdateController_Save_Settings_MissingStatus(t *testing.T) {
	handler, store, err := initHandler()
	if err != nil {
		t.Fatalf("initHandler should succeed, got error: %v", err)
	}

	// Seed initial template
	template, err := testutils.SeedTemplate(store, testutils.SITE_01, testutils.TEMPLATE_01)
	if err != nil {
		t.Fatalf("Seeding template should succeed, got error: %v", err)
	}

	newName := "Test Name No Status"
	newSiteID := testutils.SITE_01

	postData := url.Values{}
	postData.Set("template_name", newName)
	postData.Set("template_site_id", newSiteID)
	// Missing: postData.Set("template_status", ...)
	postData.Set("view", VIEW_SETTINGS) // Ensure view is set for POST

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{ // Pass template_id via GET params
			"template_id": {template.ID()},
		},
		PostValues: postData, // Send form data via POST
	})

	if err != nil {
		t.Fatalf("CallStringEndpoint should succeed, got error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK (%d), got: %d", http.StatusOK, response.StatusCode)
	}

	// Check for error message in the response body
	if !strings.Contains(body, `"icon":"error"`) {
		t.Errorf("Expected error swal icon in body: %s", body)
	}
	if !strings.Contains(body, `"text":"Status is required"`) {
		t.Errorf("Expected error swal text 'Status is required' in body: %s", body)
	}
	// Ensure no redirect script
	if strings.Contains(body, `window.location.href = "`) {
		t.Errorf("Expected no redirect script on validation error in body: %s", body)
	}

	// Verify database state (should NOT have changed)
	notUpdatedTemplate, err := store.TemplateFindByID(context.Background(), template.ID())
	if err != nil {
		t.Fatalf("Finding template after failed save should succeed, got error: %v", err)
	}
	if notUpdatedTemplate == nil {
		t.Fatalf("Template should still exist")
	}
	if notUpdatedTemplate.Name() != template.Name() {
		t.Errorf("DB name should not have changed: expected '%s', got '%s'", template.Name(), notUpdatedTemplate.Name())
	}
	if notUpdatedTemplate.Status() != template.Status() {
		t.Errorf("DB status should not have changed: expected '%s', got '%s'", template.Status(), notUpdatedTemplate.Status())
	}
}

func Test_TemplateUpdateController_Index_Success_SettingsView(t *testing.T) { // Renamed for clarity
	handler, store, err := initHandler() // Use local initHandler
	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	_, err = testutils.SeedTemplate(store, testutils.SITE_01, testutils.TEMPLATE_01)
	if err != nil {
		t.Fatalf("Failed to seed template: %v", err)
	}

	// Test GET request for settings view
	// Changed to GET as per original test
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"template_id": {testutils.TEMPLATE_01},
			"view":        {VIEW_SETTINGS}, // Explicitly request settings view
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	// Check for form fields specific to the settings view
	expecteds := []string{
		`name="template_name"`,
		`name="template_memo"`,
		`name="template_site_id"`,
		`name="template_status"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_TemplateUpdateController_Index_Success_ContentView(t *testing.T) { // Added test for content view
	handler, store, err := initHandler() // Use local initHandler
	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	_, err = testutils.SeedTemplate(store, testutils.SITE_01, testutils.TEMPLATE_01)
	if err != nil {
		t.Fatalf("Failed to seed template: %v", err)
	}

	// Test GET request for content view (default)
	// Changed to GET as per original test
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"template_id": {testutils.TEMPLATE_01},
			// "view":        {VIEW_CONTENT}, // Default view
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	// Check for form fields specific to the content view
	expecteds := []string{
		`name="template_content"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}
