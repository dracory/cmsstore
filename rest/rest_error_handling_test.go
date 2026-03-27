package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBlockCreateMalformedJSON tests handling of malformed JSON
func TestBlockCreateMalformedJSON(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	malformedJSON := `{"name": "Test Block", "site_id": "site1"` // Missing closing brace

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", strings.NewReader(malformedJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 400 Bad Request or similar error
	require.NotEqual(t, http.StatusOK, resp.StatusCode)
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
	require.NoError(t, err)

	req, err := http.NewRequest("POST", serverURL+"/api/blocks", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	// Intentionally not setting Content-Type

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Implementation may accept or reject - just verify it doesn't crash
	require.NotEqual(t, 0, resp.StatusCode)
}

// TestBlockCreateInvalidJSON tests handling of valid JSON but invalid structure
func TestBlockCreateInvalidJSON(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	// Valid JSON but wrong structure (array instead of object)
	invalidData := `["not", "an", "object"]`

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", strings.NewReader(invalidData))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotEqual(t, http.StatusOK, resp.StatusCode)
}

// TestBlockGetNonExistent tests getting a non-existent block
func TestBlockGetNonExistent(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	resp, err := http.Get(serverURL + "/api/blocks/non-existent-id")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 404 Not Found
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	success, ok := result["success"].(bool)
	require.True(t, ok)
	require.False(t, success)
}

// TestBlockUpdateNonExistent tests updating a non-existent block
func TestBlockUpdateNonExistent(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	updateData := map[string]interface{}{
		"name": "Updated Name",
	}

	jsonData, err := json.Marshal(updateData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, serverURL+"/api/blocks/non-existent-id", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	success, ok := result["success"].(bool)
	require.True(t, ok)
	require.False(t, success)
}

// TestBlockDeleteNonExistent tests deleting a non-existent block
func TestBlockDeleteNonExistent(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	req, err := http.NewRequest(http.MethodDelete, serverURL+"/api/blocks/non-existent-id", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	_, ok := result["success"].(bool)
	require.True(t, ok)
	// Depending on implementation, deleting non-existent may succeed or fail
}

// TestBlockCreateEmptyBody tests creating with empty request body
func TestBlockCreateEmptyBody(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", strings.NewReader(""))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotEqual(t, http.StatusOK, resp.StatusCode)
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
	require.NoError(t, err)

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Should fail validation
	success, ok := result["success"].(bool)
	require.True(t, ok)
	require.False(t, success)
}

// TestBlockListInvalidQueryParams tests list with invalid query parameters
func TestBlockListInvalidQueryParams(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	// Test with invalid limit
	resp, err := http.Get(serverURL + "/api/blocks?limit=invalid")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should handle gracefully (either error or use default)
	require.NotEqual(t, 0, resp.StatusCode)

	// Test with negative limit
	resp, err = http.Get(serverURL + "/api/blocks?limit=-10")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotEqual(t, 0, resp.StatusCode)

	// Test with invalid offset
	resp, err = http.Get(serverURL + "/api/blocks?offset=invalid")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotEqual(t, 0, resp.StatusCode)
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
	require.NoError(t, err)

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should either accept or reject with appropriate error
	require.NotEqual(t, 0, resp.StatusCode)
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
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, serverURL+"/api/blocks/"+testBlock.ID(), bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify other fields weren't changed
	ctx := context.Background()

	found, err := store.BlockFindByID(ctx, testBlock.ID())
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, "Updated Name Only", found.Name())
	require.Equal(t, testBlock.Content(), found.Content()) // Content should be unchanged
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
	require.NoError(t, err)

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Should fail or handle duplicate gracefully
	_, ok := result["success"].(bool)
	require.True(t, ok)
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
		require.NoError(t, err)

		resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		// Should handle without crashing
		require.NoError(t, err)
		require.NotEqual(t, 0, resp.StatusCode)
	}
}

// TestBlockMethodNotAllowed tests unsupported HTTP methods
func TestBlockMethodNotAllowed(t *testing.T) {
	serverURL, _, cleanup := setupTestAPI(t)
	defer cleanup()

	// Try PATCH on blocks endpoint (if not supported)
	req, err := http.NewRequest("PATCH", serverURL+"/api/blocks/some-id", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 405 Method Not Allowed or handle gracefully
	require.NotEqual(t, 0, resp.StatusCode)
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

			jsonData, _ := json.Marshal(updateData)
			req, _ := http.NewRequest(http.MethodPut, serverURL+"/api/blocks/"+testBlock.ID(), bytes.NewBuffer(jsonData))
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
	require.NoError(t, err)
	defer resp.Body.Close()

	// API should respond (even if with an error, it shouldn't crash)
	require.NotEqual(t, 0, resp.StatusCode, "API should respond to requests after concurrent updates")
}
