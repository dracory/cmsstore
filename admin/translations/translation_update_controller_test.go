package admin

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	_ "modernc.org/sqlite"
)

func initTranslationUpdateHandler() (func(w http.ResponseWriter, r *http.Request) string, cmsstore.StoreInterface, error) {
	store, err := testutils.InitStore(":memory:")

	if err != nil {
		return nil, nil, err
	}

	return NewTranslationUpdateController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Layout: shared.Layout,
	})).Handler, store, nil
}

func Test_TranslationUpdateController_Index_RequiresTranslationID(t *testing.T) {
	handler, _, err := initTranslationUpdateHandler()

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
		`"message":"translation id is required"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_TranslationUpdateController_Index_RequiresValidTranslation(t *testing.T) {
	handler, _, err := initTranslationUpdateHandler()

	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"translation_id": {"trans-123"},
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
		`"message":"translation not found"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_TranslationUpdateController_Index_Success(t *testing.T) {
	handler, store, err := initTranslationUpdateHandler()

	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	seededTranslation, err := testutils.SeedTranslation(store, testutils.SITE_01, testutils.TRANSLATION_01)

	if err != nil {
		t.Fatalf("Failed to seed translation: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"translation_id": {seededTranslation.ID()},
			"view":           {VIEW_SETTINGS},
		},
	})

	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	expecteds := []string{
		`name="translation_name"`,
		`name="translation_memo"`,
		`name="translation_site_id"`,
		`name="translation_status"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

// Test_TranslationUpdateController_Save_Settings_Success tests successful saving of settings view
func Test_TranslationUpdateController_Save_Settings_Success(t *testing.T) {
	handler, store, err := initTranslationUpdateHandler()
	if err != nil {
		t.Fatalf("initTranslationUpdateHandler should succeed, got error: %v", err)
	}

	// Seed initial translation using testutils
	translation, err := testutils.SeedTranslation(store, testutils.SITE_01, testutils.TRANSLATION_01)
	if err != nil {
		t.Fatalf("Seeding translation should succeed, got error: %v", err)
	}
	if translation == nil {
		t.Fatalf("Seeded translation should not be nil")
	}

	newName := "Updated Translation Name"
	newMemo := "Updated Memo Content"
	newSiteID := testutils.SITE_02 // Change site ID
	newStatus := cmsstore.TRANSLATION_STATUS_ACTIVE

	postData := url.Values{}
	postData.Set("translation_name", newName)
	postData.Set("translation_memo", newMemo)
	postData.Set("translation_site_id", newSiteID)
	postData.Set("translation_status", newStatus)
	postData.Set("view", VIEW_SETTINGS) // Ensure view is set for POST

	// Call the handler using test.CallStringEndpoint
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{ // Pass translation_id via GET params as per handler logic
			"translation_id": {translation.ID()},
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
	if !strings.Contains(body, `"text":"translation saved successfully"`) {
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
	updatedTranslation, err := store.TranslationFindByID(context.Background(), translation.ID())
	if err != nil {
		t.Fatalf("Finding updated translation should succeed, got error: %v", err)
	}
	if updatedTranslation == nil {
		t.Fatalf("Updated translation should exist in DB")
	}
	if updatedTranslation.Name() != newName {
		t.Errorf("DB name mismatch: expected '%s', got '%s'", newName, updatedTranslation.Name())
	}
	if updatedTranslation.Memo() != newMemo {
		t.Errorf("DB memo mismatch: expected '%s', got '%s'", newMemo, updatedTranslation.Memo())
	}
	if updatedTranslation.SiteID() != newSiteID {
		t.Errorf("DB site ID mismatch: expected '%s', got '%s'", newSiteID, updatedTranslation.SiteID())
	}
	if updatedTranslation.Status() != newStatus {
		t.Errorf("DB status mismatch: expected '%s', got '%s'", newStatus, updatedTranslation.Status())
	}
}

// Test_TranslationUpdateController_Save_Content_Success tests successful saving of content view
func Test_TranslationUpdateController_Save_Content_Success(t *testing.T) {
	handler, store, err := initTranslationUpdateHandler()
	if err != nil {
		t.Fatalf("initTranslationUpdateHandler should succeed, got error: %v", err)
	}

	translation, err := testutils.SeedTranslation(store, testutils.SITE_01, testutils.TRANSLATION_01)
	if err != nil {
		t.Fatalf("Seeding translation should succeed, got error: %v", err)
	}

	newHandle := "updated-translation-handle"
	newContent := map[string]string{"en": "English content", "es": "Spanish content"}

	postData := url.Values{}
	postData.Set("translation_handle", newHandle)
	postData.Set("translation_content[en]", newContent["en"])
	postData.Set("translation_content[es]", newContent["es"])
	postData.Set("view", VIEW_CONTENT) // Ensure view is set for POST

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{ // Pass translation_id via GET params
			"translation_id": {translation.ID()},
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
	if !strings.Contains(body, `"text":"translation saved successfully"`) {
		t.Errorf("Expected success swal text in body: %s", body)
	}
	if !strings.Contains(body, `window.location.href = "`) {
		t.Errorf("Expected redirect script in body: %s", body)
	}
	if !strings.Contains(body, "view="+VIEW_CONTENT) {
		t.Errorf("Expected redirect back to content view in body: %s", body)
	}

	// Verify database state
	updatedTranslation, err := store.TranslationFindByID(context.Background(), translation.ID())
	if err != nil {
		t.Fatalf("Finding updated translation should succeed, got error: %v", err)
	}
	if updatedTranslation == nil {
		t.Fatalf("Updated translation should exist in DB")
	}
	if updatedTranslation.Handle() != newHandle {
		t.Errorf("DB handle mismatch: expected '%s', got '%s'", newHandle, updatedTranslation.Handle())
	}
	// Content might be stored as JSON
	translationContent, _ := updatedTranslation.Content()
	if translationContent["en"] != newContent["en"] {
		t.Errorf("DB content mismatch: expected '%s', got '%s'", newContent["en"], translationContent["en"])
	}
	// Ensure other fields didn't change
	if updatedTranslation.Name() != translation.Name() {
		t.Errorf("DB name should not change: expected '%s', got '%s'", translation.Name(), updatedTranslation.Name())
	}
	if updatedTranslation.SiteID() != translation.SiteID() {
		t.Errorf("DB site ID should not change: expected '%s', got '%s'", translation.SiteID(), updatedTranslation.SiteID())
	}
}

// Test_TranslationUpdateController_Save_Settings_MissingStatus tests validation failure
func Test_TranslationUpdateController_Save_Settings_MissingStatus(t *testing.T) {
	handler, store, err := initTranslationUpdateHandler()
	if err != nil {
		t.Fatalf("initTranslationUpdateHandler should succeed, got error: %v", err)
	}

	// Seed initial translation
	translation, err := testutils.SeedTranslation(store, testutils.SITE_01, testutils.TRANSLATION_01)
	if err != nil {
		t.Fatalf("Seeding translation should succeed, got error: %v", err)
	}

	newName := "Test Name No Status"
	newSiteID := testutils.SITE_01

	postData := url.Values{}
	postData.Set("translation_name", newName)
	postData.Set("translation_site_id", newSiteID)
	// Missing: postData.Set("translation_status", ...)
	postData.Set("view", VIEW_SETTINGS) // Ensure view is set for POST

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{ // Pass translation_id via GET params
			"translation_id": {translation.ID()},
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
	notUpdatedTranslation, err := store.TranslationFindByID(context.Background(), translation.ID())
	if err != nil {
		t.Fatalf("Finding translation after failed save should succeed, got error: %v", err)
	}
	if notUpdatedTranslation == nil {
		t.Fatalf("Translation should still exist")
	}
	if notUpdatedTranslation.Name() != translation.Name() {
		t.Errorf("DB name should not have changed: expected '%s', got '%s'", translation.Name(), notUpdatedTranslation.Name())
	}
	if notUpdatedTranslation.Status() != translation.Status() {
		t.Errorf("DB status should not have changed: expected '%s', got '%s'", translation.Status(), notUpdatedTranslation.Status())
	}
}

func Test_TranslationUpdateController_Index_Success_SettingsView(t *testing.T) {
	handler, store, err := initTranslationUpdateHandler()
	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	seededTranslation, err := testutils.SeedTranslation(store, testutils.SITE_01, testutils.TRANSLATION_01)
	if err != nil {
		t.Fatalf("Failed to seed translation: %v", err)
	}

	// Test GET request for settings view
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"translation_id": {seededTranslation.ID()},
			"view":           {VIEW_SETTINGS}, // Explicitly request settings view
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
		`name="translation_name"`,
		`name="translation_memo"`,
		`name="translation_site_id"`,
		`name="translation_status"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}

func Test_TranslationUpdateController_Index_Success_ContentView(t *testing.T) {
	handler, store, err := initTranslationUpdateHandler()
	if err != nil {
		t.Fatalf("Failed to initialize controller: %v", err)
	}

	seededTranslation, err := testutils.SeedTranslation(store, testutils.SITE_01, testutils.TRANSLATION_01)
	if err != nil {
		t.Fatalf("Failed to seed translation: %v", err)
	}

	// Test GET request for content view (default)
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"translation_id": {seededTranslation.ID()},
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
		`name="translation_handle"`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatalf("Expected to find %s in the response body, but found: %s", expected, body)
		}
	}
}
