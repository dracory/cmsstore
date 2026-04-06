package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
)

// CreateTestSite is now defined in test_utils.go
// createTestTranslation creates a test translation and returns it
func createTestTranslation(t *testing.T, store cmsstore.StoreInterface, siteID string) cmsstore.TranslationInterface {
	translation := cmsstore.NewTranslation()
	translation.SetName("welcome_message")
	translation.SetHandle("welcome_message")
	translation.SetSiteID(siteID)
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	contentMap := map[string]string{
		"en": "Welcome to our site",
		"fr": "Bienvenue sur notre site",
	}
	err := translation.SetContent(contentMap)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = translation.SetMeta("key", "welcome_message")
	if err != nil {
		t.Fatalf("Failed to set translation key metadata: %v", err)
	}

	err = translation.SetMeta("description", "Welcome message for the homepage")
	if err != nil {
		t.Fatalf("Failed to set translation description metadata: %v", err)
	}

	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create test translation: %v", err)
	}

	return translation
}

// TestListTranslations tests the GET /api/translations endpoint
func TestListTranslations(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	testSite, cleanupSite := CreateTestSite(t, store)
	defer cleanupSite()

	// Create a test translation
	_ = createTestTranslation(t, store, testSite.ID())

	// Test the endpoint
	resp, err := http.Get(serverURL + "/api/translations?site_id=" + testSite.ID())
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success to be true, got %v", success)
	}

	translations, ok := result["translations"].([]interface{})
	if !ok {
		t.Errorf("Expected translations to be an array, got %v", result["translations"])
	}
	if len(translations) < 1 {
		t.Errorf("Expected at least one translation, got %d", len(translations))
	}
}

// TestListTranslationsWithLocaleFilter tests the GET /api/translations endpoint with locale filter
func TestListTranslationsWithLocaleFilter(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	testSite, cleanupSite := CreateTestSite(t, store)
	defer cleanupSite()

	// Create a test translation
	translation := createTestTranslation(t, store, testSite.ID())

	// Set the locale we'll filter by
	err := translation.SetMeta("locale", "fr")
	if err != nil {
		t.Fatalf("Failed to set translation locale metadata: %v", err)
	}

	// Update the translation with the new metadata
	err = store.TranslationUpdate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to update test translation with locale metadata: %v", err)
	}

	// Test the endpoint with locale filter
	resp, err := http.Get(serverURL + "/api/translations?site_id=" + testSite.ID() + "&locale=fr")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success to be true, got %v", success)
	}

	translations, ok := result["translations"].([]interface{})
	if !ok {
		t.Errorf("Expected translations to be an array, got %v", result["translations"])
	}
	if len(translations) < 1 {
		t.Errorf("Expected at least one translation with locale 'fr', got %d", len(translations))
	}
}

// TestGetTranslation tests the GET /api/translations/:id endpoint
func TestGetTranslation(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	testSite, cleanupSite := CreateTestSite(t, store)
	defer cleanupSite()

	// Create a test translation
	translation := createTestTranslation(t, store, testSite.ID())

	// Test the endpoint
	resp, err := http.Get(serverURL + "/api/translations/" + translation.ID())
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success to be true, got %v", success)
	}

	id, ok := result["id"].(string)
	if !ok {
		t.Errorf("Expected id to be a string, got %v", result["id"])
	}
	if id != translation.ID() {
		t.Errorf("Expected translation ID %s, got %s", translation.ID(), id)
	}

	key, ok := result["key"].(string)
	if !ok {
		t.Errorf("Expected key to be a string, got %v", result["key"])
	}
	if key != "welcome_message" {
		t.Errorf("Expected translation key 'welcome_message', got %s", key)
	}

	// Check that content is properly returned
	content, ok := result["content"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected content to be a map, got %v", result["content"])
	}

	enText, ok := content["en"].(string)
	if !ok {
		t.Errorf("Expected English text to be present, got %v", content["en"])
	}
	if enText != "Welcome to our site" {
		t.Errorf("Expected English text 'Welcome to our site', got %s", enText)
	}

	frText, ok := content["fr"].(string)
	if !ok {
		t.Errorf("Expected French text to be present, got %v", content["fr"])
	}
	if frText != "Bienvenue sur notre site" {
		t.Errorf("Expected French text 'Bienvenue sur notre site', got %s", frText)
	}
}

// TestCreateTranslation tests the POST /api/translations endpoint
func TestCreateTranslation(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	testSite, cleanupSite := CreateTestSite(t, store)
	defer cleanupSite()

	translationData := map[string]interface{}{
		"key":     "login_button",
		"site_id": testSite.ID(),
		"status":  cmsstore.TRANSLATION_STATUS_ACTIVE,
		"content": map[string]string{
			"en": "Login",
			"fr": "Connexion",
		},
	}

	jsonData, err := json.Marshal(translationData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(serverURL+"/api/translations", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success to be true, got %v", success)
	}

	id, ok := result["id"].(string)
	if !ok {
		t.Errorf("Expected id to be a string, got %v", result["id"])
	}
	if id == "" {
		t.Errorf("Expected non-empty id, got %s", id)
	}

	key, ok := result["key"].(string)
	if !ok {
		t.Errorf("Expected key to be a string, got %v", result["key"])
	}
	if key != "login_button" {
		t.Errorf("Expected translation key 'login_button', got %s", key)
	}

	// Check that content is properly returned
	content, ok := result["content"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected content to be a map, got %v", result["content"])
	}

	enText, ok := content["en"].(string)
	if !ok {
		t.Errorf("Expected English text to be present, got %v", content["en"])
	}
	if enText != "Login" {
		t.Errorf("Expected English text 'Login', got %s", enText)
	}

	frText, ok := content["fr"].(string)
	if !ok {
		t.Errorf("Expected French text to be present, got %v", content["fr"])
	}
	if frText != "Connexion" {
		t.Errorf("Expected French text 'Connexion', got %s", frText)
	}
}

// TestUpdateTranslation tests the PUT /api/translations/:id endpoint
func TestUpdateTranslation(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	testSite, cleanupSite := CreateTestSite(t, store)
	defer cleanupSite()

	// Create a test translation to update
	translation := createTestTranslation(t, store, testSite.ID())

	updateData := map[string]interface{}{
		"key": "welcome_message",
		"content": map[string]string{
			"en": "Welcome to our updated site",
			"fr": "Bienvenue sur notre site mis à jour",
			"es": "Bienvenido a nuestro sitio actualizado", // Adding a new language
		},
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, serverURL+"/api/translations/"+translation.ID(), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success to be true, got %v", success)
	}

	// Check that content is properly updated
	content, ok := result["content"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected content to be a map, got %v", result["content"])
	}

	enText, ok := content["en"].(string)
	if !ok {
		t.Errorf("Expected English text to be present, got %v", content["en"])
	}
	if enText != "Welcome to our updated site" {
		t.Errorf("Expected English text 'Welcome to our updated site', got %s", enText)
	}

	frText, ok := content["fr"].(string)
	if !ok {
		t.Errorf("Expected French text to be present, got %v", content["fr"])
	}
	if frText != "Bienvenue sur notre site mis à jour" {
		t.Errorf("Expected French text 'Bienvenue sur notre site mis à jour', got %s", frText)
	}

	esText, ok := content["es"].(string)
	if !ok {
		t.Errorf("Expected Spanish text to be present, got %v", content["es"])
	}
	if esText != "Bienvenido a nuestro sitio actualizado" {
		t.Errorf("Expected Spanish text 'Bienvenido a nuestro sitio actualizado', got %s", esText)
	}
}

// TestDeleteTranslation tests the DELETE /api/translations/:id endpoint
func TestDeleteTranslation(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	testSite, cleanupSite := CreateTestSite(t, store)
	defer cleanupSite()

	// Create a test translation to delete
	translation := createTestTranslation(t, store, testSite.ID())

	req, err := http.NewRequest(http.MethodDelete, serverURL+"/api/translations/"+translation.ID(), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success to be true, got %v", success)
	}

	// Verify the translation was soft deleted by querying with soft-deleted included
	translations, err := store.TranslationList(context.Background(),
		cmsstore.TranslationQuery().
			SetID(translation.ID()).
			SetSoftDeletedIncluded(true).
			SetLimit(1))
	if err != nil {
		t.Fatalf("Failed to find translation: %v", err)
	}
	if len(translations) == 0 {
		t.Errorf("Translation should still exist after soft delete")
	}

	translationAfterDelete := translations[0]
	if !translationAfterDelete.IsSoftDeleted() {
		t.Errorf("Translation should be marked as soft deleted")
	}
}
