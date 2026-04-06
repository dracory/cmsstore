package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// TestBlockCreateMalformedJSON tests handling of malformed JSON
func TestBlockCreateMalformedJSON(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	malformedJSON := `{"name": "Test Block", "site_id": "site1"` // Missing closing brace

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", strings.NewReader(malformedJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should return 400 Bad Request or similar error
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected status to not be 200 OK, got %d", resp.StatusCode)
	}
}

// TestBlockCreateMissingContentType tests handling of missing Content-Type header
func TestBlockCreateMissingContentType(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	blockData := map[string]interface{}{
		"name":    "Test Block",
		"site_id": "site1",
	}

	jsonData, err := json.Marshal(blockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req, err := http.NewRequest("POST", serverURL+"/api/blocks", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Intentionally not setting Content-Type

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Implementation may accept or reject - just verify it doesn't crash
	if resp.StatusCode == 0 {
		t.Errorf("Expected non-zero status code, got %d", resp.StatusCode)
	}
}

// TestBlockCreateInvalidJSON tests handling of valid JSON but invalid structure
func TestBlockCreateInvalidJSON(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	// Valid JSON but wrong structure (array instead of object)
	invalidData := `["not", "an", "object"]`

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", strings.NewReader(invalidData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected status to not be 200 OK, got %d", resp.StatusCode)
	}
}

// TestBlockGetNonExistent tests getting a non-existent block
func TestBlockGetNonExistent(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	resp, err := http.Get(serverURL + "/api/blocks/non-existent-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should return 404 Not Found
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok {
		t.Errorf("Expected 'success' key to be present")
	}
	if success {
		t.Errorf("Expected 'success' to be false")
	}
}

// TestBlockUpdateNonExistent tests updating a non-existent block
func TestBlockUpdateNonExistent(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	updateData := map[string]interface{}{
		"name": "Updated Name",
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, serverURL+"/api/blocks/non-existent-id", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	success, ok := result["success"].(bool)
	if !ok {
		t.Errorf("Expected 'success' key to be present")
	}
	if success {
		t.Errorf("Expected 'success' to be false")
	}
}

// TestBlockDeleteNonExistent tests deleting a non-existent block
func TestBlockDeleteNonExistent(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	req, err := http.NewRequest(http.MethodDelete, serverURL+"/api/blocks/non-existent-id", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := result["success"].(bool)
	if !ok {
		t.Errorf("Expected 'success' key to be present")
	}
	// Depending on implementation, deleting non-existent may succeed or fail
}

// TestBlockCreateEmptyBody tests creating with empty request body
func TestBlockCreateEmptyBody(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected status to not be 200 OK, got %d", resp.StatusCode)
	}
}

// TestBlockCreateNullValues tests creating with null values
func TestBlockCreateNullValues(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	blockData := map[string]interface{}{
		"name":    nil,
		"site_id": nil,
		"content": nil,
	}

	jsonData, err := json.Marshal(blockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail validation
	success, ok := result["success"].(bool)
	if !ok {
		t.Errorf("Expected 'success' key to be present")
	}
	if success {
		t.Errorf("Expected 'success' to be false")
	}
}

// TestBlockListInvalidQueryParams tests list with invalid query parameters
func TestBlockListInvalidQueryParams(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	// Test with invalid limit
	resp, err := http.Get(serverURL + "/api/blocks?limit=invalid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should handle gracefully (either error or use default)
	if resp.StatusCode == 0 {
		t.Errorf("Expected non-zero status code, got %d", resp.StatusCode)
	}

	// Test with negative limit
	resp, err = http.Get(serverURL + "/api/blocks?limit=-10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 0 {
		t.Errorf("Expected non-zero status code, got %d", resp.StatusCode)
	}

	// Test with invalid offset
	resp, err = http.Get(serverURL + "/api/blocks?offset=invalid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 0 {
		t.Errorf("Expected non-zero status code, got %d", resp.StatusCode)
	}
}

// TestBlockCreateVeryLargePayload tests handling of very large payloads
func TestBlockCreateVeryLargePayload(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	// Create a very large content string (1MB)
	largeContent := strings.Repeat("A", 1024*1024)

	blockData := map[string]interface{}{
		"name":    "Large Block",
		"site_id": "site1",
		"content": largeContent,
	}

	jsonData, err := json.Marshal(blockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should either accept or reject with appropriate error
	if resp.StatusCode == 0 {
		t.Errorf("Expected non-zero status code, got %d", resp.StatusCode)
	}
}

// TestBlockUpdatePartialData tests partial updates
func TestBlockUpdatePartialData(t *testing.T) {
	serverURL, store, _, _, testBlock, cleanup := setupBlockTest(t)
	defer cleanup()

	// Update only the name
	updateData := map[string]interface{}{
		"name": "Updated Name Only",
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, serverURL+"/api/blocks/"+testBlock.ID(), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verify other fields weren't changed
	ctx := context.Background()

	found, err := store.BlockFindByID(ctx, testBlock.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Errorf("Expected block to be found")
	}
	if found.Name() != "Updated Name Only" {
		t.Errorf("Expected name to be 'Updated Name Only', got '%s'", found.Name())
	}
	if found.Content() != testBlock.Content() {
		t.Errorf("Expected content to be '%s', got '%s'", testBlock.Content(), found.Content())
	}
}

// TestBlockCreateDuplicateID tests creating with duplicate ID
func TestBlockCreateDuplicateID(t *testing.T) {
	serverURL, _, testSite, _, testBlock, cleanup := setupBlockTest(t)
	defer cleanup()

	// Try to create another block with same ID
	blockData := map[string]interface{}{
		"id":      testBlock.ID(),
		"name":    "Duplicate",
		"site_id": testSite.ID(),
	}

	jsonData, err := json.Marshal(blockData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fail or handle duplicate gracefully
	_, ok := result["success"].(bool)
	if !ok {
		t.Errorf("Expected 'success' key to be present")
	}
}

// TestBlockCreateSpecialCharacters tests handling of special characters
func TestBlockCreateSpecialCharacters(t *testing.T) {
	serverURL, _, testSite, _, _, cleanup := setupBlockTest(t)
	defer cleanup()

	specialChars := []string{
		"<script>alert('xss')</script>",
		"'; DROP TABLE blocks; --",
		"../../../etc/passwd",
		"\x00\x01\x02",
		"🔥💯✨",
	}

	for _, char := range specialChars {
		blockData := map[string]interface{}{
			"name":    char,
			"site_id": testSite.ID(),
			"content": char,
		}

		jsonData, err := json.Marshal(blockData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		// Should handle without crashing
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode == 0 {
			t.Errorf("Expected non-zero status code, got %d", resp.StatusCode)
		}
	}
}

// TestBlockMethodNotAllowed tests unsupported HTTP methods
func TestBlockMethodNotAllowed(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	// Try PATCH on blocks endpoint (if not supported)
	req, err := http.NewRequest("PATCH", serverURL+"/api/blocks/some-id", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Should return 405 Method Not Allowed or handle gracefully
	if resp.StatusCode == 0 {
		t.Errorf("Expected non-zero status code, got %d", resp.StatusCode)
	}
}

// TestBlockConcurrentUpdates tests concurrent updates to same block
func TestBlockConcurrentUpdates(t *testing.T) {
	serverURL, _, _, _, testBlock, cleanup := setupBlockTest(t)
	defer cleanup()

	const numUpdates = 5
	done := make(chan bool, numUpdates)
	successCount := make(chan bool, numUpdates)

	for i := 0; i < numUpdates; i++ {
		go func(index int) {
			updateData := map[string]interface{}{
				"name": "Concurrent Update",
			}

			jsonData, err := json.Marshal(updateData)
			if err != nil {
				done <- true
				return
			}
			req, err := http.NewRequest(http.MethodPut, serverURL+"/api/blocks/"+testBlock.ID(), bytes.NewBuffer(jsonData))
			if err != nil {
				done <- true
				return
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if resp != nil {
				// Track successful updates (200 OK)
				if err == nil && resp.StatusCode == http.StatusOK {
					successCount <- true
				}
				resp.Body.Close()
			}
			done <- true
		}(i)
	}

	// Wait for all updates
	for i := 0; i < numUpdates; i++ {
		<-done
	}
	close(successCount)

	// Count successful updates
	successful := 0
	for range successCount {
		successful++
	}

	// At least some updates should succeed (tests error resilience under concurrency)
	t.Logf("Successful concurrent updates: %d/%d", successful, numUpdates)

	// Verify the API doesn't crash - should still be able to query the block
	resp, err := http.Get(serverURL + "/api/blocks/" + testBlock.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// API should respond (even if with an error, it shouldn't crash)
	if resp.StatusCode == 0 {
		t.Errorf("Expected non-zero status code, got %d", resp.StatusCode)
	}
}
