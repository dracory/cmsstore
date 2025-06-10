package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gouniverse/cmsstore"
)

func TestTranslationEndpoints(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	// Create a test site for our tests
	testSite := cmsstore.NewSite()
	testSite.SetName("Test Site for Translations")
	testSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), testSite)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	// Create a test translation for our tests
	testTranslation := cmsstore.NewTranslation()
	testTranslation.SetName("welcome_message") // Key is stored in Name

	// Set content as a map of locale to text
	contentMap := map[string]string{
		"en": "Welcome to our site",
		"fr": "Bienvenue sur notre site",
	}
	err = testTranslation.SetContent(contentMap)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	// Store key and locale in metadata
	err = testTranslation.SetMeta("key", "welcome_message")
	if err != nil {
		t.Fatalf("Failed to set translation key metadata: %v", err)
	}

	testTranslation.SetSiteID(testSite.ID())
	testTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), testTranslation)
	if err != nil {
		t.Fatalf("Failed to create test translation: %v", err)
	}

	t.Run("List Translations", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/translations?site_id=" + testSite.ID())
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if success, ok := result["success"].(bool); !ok || !success {
			t.Errorf("Expected success to be true, got %v", result["success"])
		}

		translations, ok := result["translations"].([]interface{})
		if !ok {
			t.Fatalf("Expected translations to be an array, got %T", result["translations"])
		}

		if len(translations) < 1 {
			t.Errorf("Expected at least one translation, got %d", len(translations))
		}
	})

	t.Run("List Translations with Locale Filter", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/translations?site_id=" + testSite.ID() + "&locale=fr")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if success, ok := result["success"].(bool); !ok || !success {
			t.Errorf("Expected success to be true, got %v", result["success"])
		}

		translations, ok := result["translations"].([]interface{})
		if !ok {
			t.Fatalf("Expected translations to be an array, got %T", result["translations"])
		}

		// Since we're filtering by locale, we should still see our test translation
		if len(translations) < 1 {
			t.Errorf("Expected at least one translation with locale 'fr', got %d", len(translations))
		}
	})

	t.Run("Get Translation", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/translations/" + testTranslation.ID())
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if success, ok := result["success"].(bool); !ok || !success {
			t.Errorf("Expected success to be true, got %v", result["success"])
		}

		if id, ok := result["id"].(string); !ok || id != testTranslation.ID() {
			t.Errorf("Expected id to be %s, got %v", testTranslation.ID(), id)
		}

		if key, ok := result["key"].(string); !ok || key != "welcome_message" {
			t.Errorf("Expected key to be 'welcome_message', got %v", key)
		}

		// Check that content is properly returned
		content, ok := result["content"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected content to be a map, got %T", result["content"])
		}

		if enText, ok := content["en"].(string); !ok || enText != "Welcome to our site" {
			t.Errorf("Expected English text to be 'Welcome to our site', got %v", enText)
		}

		if frText, ok := content["fr"].(string); !ok || frText != "Bienvenue sur notre site" {
			t.Errorf("Expected French text to be 'Bienvenue sur notre site', got %v", frText)
		}
	})

	t.Run("Create Translation", func(t *testing.T) {
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
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if success, ok := result["success"].(bool); !ok || !success {
			t.Errorf("Expected success to be true, got %v", result["success"])
		}

		if _, ok := result["id"].(string); !ok {
			t.Errorf("Expected id to be a string")
		}

		if key, ok := result["key"].(string); !ok || key != "login_button" {
			t.Errorf("Expected key to be 'login_button', got %v", key)
		}

		// Check that content is properly returned
		content, ok := result["content"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected content to be a map, got %T", result["content"])
		}

		if enText, ok := content["en"].(string); !ok || enText != "Login" {
			t.Errorf("Expected English text to be 'Login', got %v", enText)
		}

		if frText, ok := content["fr"].(string); !ok || frText != "Connexion" {
			t.Errorf("Expected French text to be 'Connexion', got %v", frText)
		}
	})

	t.Run("Update Translation", func(t *testing.T) {
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

		req, err := http.NewRequest(http.MethodPut, serverURL+"/api/translations/"+testTranslation.ID(), bytes.NewBuffer(jsonData))
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
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if success, ok := result["success"].(bool); !ok || !success {
			t.Errorf("Expected success to be true, got %v", result["success"])
		}

		// Check that content is properly updated
		content, ok := result["content"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected content to be a map, got %T", result["content"])
		}

		if enText, ok := content["en"].(string); !ok || enText != "Welcome to our updated site" {
			t.Errorf("Expected English text to be 'Welcome to our updated site', got %v", enText)
		}

		if frText, ok := content["fr"].(string); !ok || frText != "Bienvenue sur notre site mis à jour" {
			t.Errorf("Expected French text to be 'Bienvenue sur notre site mis à jour', got %v", frText)
		}

		if esText, ok := content["es"].(string); !ok || esText != "Bienvenido a nuestro sitio actualizado" {
			t.Errorf("Expected Spanish text to be 'Bienvenido a nuestro sitio actualizado', got %v", esText)
		}
	})

	t.Run("Delete Translation", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, serverURL+"/api/translations/"+testTranslation.ID(), nil)
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
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if success, ok := result["success"].(bool); !ok || !success {
			t.Errorf("Expected success to be true, got %v", result["success"])
		}

		// Verify the translation was soft deleted
		translation, err := store.TranslationFindByID(context.Background(), testTranslation.ID())
		if err != nil {
			t.Fatalf("Failed to find translation: %v", err)
		}
		if translation == nil {
			t.Fatalf("Translation should still exist after soft delete")
		}
		if !translation.IsSoftDeleted() {
			t.Errorf("Translation should be marked as soft deleted")
		}
	})
}
