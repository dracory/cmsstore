package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gouniverse/cmsstore"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err, "Failed to set translation content")

	err = translation.SetMeta("key", "welcome_message")
	require.NoError(t, err, "Failed to set translation key metadata")

	err = translation.SetMeta("description", "Welcome message for the homepage")
	require.NoError(t, err, "Failed to set translation description metadata")

	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err, "Failed to create test translation")

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
	require.NoError(t, err, "Failed to make request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode response")

	success, ok := result["success"].(bool)
	require.True(t, ok && success, "Expected success to be true")

	translations, ok := result["translations"].([]interface{})
	require.True(t, ok, "Expected translations to be an array")
	require.GreaterOrEqual(t, len(translations), 1, "Expected at least one translation")
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
	require.NoError(t, err, "Failed to set translation locale metadata")

	// Update the translation with the new metadata
	err = store.TranslationUpdate(context.Background(), translation)
	require.NoError(t, err, "Failed to update test translation with locale metadata")

	// Test the endpoint with locale filter
	resp, err := http.Get(serverURL + "/api/translations?site_id=" + testSite.ID() + "&locale=fr")
	require.NoError(t, err, "Failed to make request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode response")

	success, ok := result["success"].(bool)
	require.True(t, ok && success, "Expected success to be true")

	translations, ok := result["translations"].([]interface{})
	require.True(t, ok, "Expected translations to be an array")
	require.GreaterOrEqual(t, len(translations), 1, "Expected at least one translation with locale 'fr'")
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
	require.NoError(t, err, "Failed to make request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode response")

	success, ok := result["success"].(bool)
	require.True(t, ok && success, "Expected success to be true")

	id, ok := result["id"].(string)
	require.True(t, ok, "Expected id to be a string")
	require.Equal(t, translation.ID(), id, "Unexpected translation ID")

	key, ok := result["key"].(string)
	require.True(t, ok, "Expected key to be a string")
	require.Equal(t, "welcome_message", key, "Unexpected translation key")

	// Check that content is properly returned
	content, ok := result["content"].(map[string]interface{})
	require.True(t, ok, "Expected content to be a map")

	enText, ok := content["en"].(string)
	require.True(t, ok, "Expected English text to be present")
	require.Equal(t, "Welcome to our site", enText, "Unexpected English text")

	frText, ok := content["fr"].(string)
	require.True(t, ok, "Expected French text to be present")
	require.Equal(t, "Bienvenue sur notre site", frText, "Unexpected French text")
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
	require.NoError(t, err, "Failed to marshal JSON")

	resp, err := http.Post(serverURL+"/api/translations", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err, "Failed to make request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode response")

	success, ok := result["success"].(bool)
	require.True(t, ok && success, "Expected success to be true")

	id, ok := result["id"].(string)
	require.True(t, ok, "Expected id to be a string")
	require.NotEmpty(t, id, "Expected non-empty id")

	key, ok := result["key"].(string)
	require.True(t, ok, "Expected key to be a string")
	require.Equal(t, "login_button", key, "Unexpected translation key")

	// Check that content is properly returned
	content, ok := result["content"].(map[string]interface{})
	require.True(t, ok, "Expected content to be a map")

	enText, ok := content["en"].(string)
	require.True(t, ok, "Expected English text to be present")
	require.Equal(t, "Login", enText, "Unexpected English text")

	frText, ok := content["fr"].(string)
	require.True(t, ok, "Expected French text to be present")
	require.Equal(t, "Connexion", frText, "Unexpected French text")
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
	require.NoError(t, err, "Failed to marshal JSON")

	req, err := http.NewRequest(http.MethodPut, serverURL+"/api/translations/"+translation.ID(), bytes.NewBuffer(jsonData))
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to make request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode response")

	success, ok := result["success"].(bool)
	require.True(t, ok && success, "Expected success to be true")

	// Check that content is properly updated
	content, ok := result["content"].(map[string]interface{})
	require.True(t, ok, "Expected content to be a map")

	enText, ok := content["en"].(string)
	require.True(t, ok, "Expected English text to be present")
	require.Equal(t, "Welcome to our updated site", enText, "Unexpected English text")

	frText, ok := content["fr"].(string)
	require.True(t, ok, "Expected French text to be present")
	require.Equal(t, "Bienvenue sur notre site mis à jour", frText, "Unexpected French text")

	esText, ok := content["es"].(string)
	require.True(t, ok, "Expected Spanish text to be present")
	require.Equal(t, "Bienvenido a nuestro sitio actualizado", esText, "Unexpected Spanish text")
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
	require.NoError(t, err, "Failed to create request")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to make request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode response")

	success, ok := result["success"].(bool)
	require.True(t, ok && success, "Expected success to be true")

	// Verify the translation was soft deleted by querying with soft-deleted included
	translations, err := store.TranslationList(context.Background(), 
		cmsstore.TranslationQuery().
			SetID(translation.ID()).
			SetSoftDeletedIncluded(true).
			SetLimit(1))
	require.NoError(t, err, "Failed to find translation")
	require.NotEmpty(t, translations, "Translation should still exist after soft delete")
	
	translationAfterDelete := translations[0]
	require.True(t, translationAfterDelete.IsSoftDeleted(), "Translation should be marked as soft deleted")
}
