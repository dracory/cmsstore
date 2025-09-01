package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dracory/cmsstore"
)

// handleBlocksEndpoint handles HTTP requests for the /api/blocks endpoint
func (api *RestAPI) handleBlocksEndpoint(w http.ResponseWriter, r *http.Request, pathParts []string) {
	switch r.Method {
	case http.MethodPost:
		// Create a new block
		api.handleBlockCreate(w, r)
	case http.MethodGet:
		// Get block(s)
		if len(pathParts) > 0 && pathParts[0] != "" {
			// Get a specific block by ID
			api.handleBlockGet(w, r, pathParts[0])
		} else {
			// List all blocks
			api.handleBlockList(w, r)
		}
	case http.MethodPut:
		// Update a block
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handleBlockUpdate(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Block ID required for update"}`, http.StatusBadRequest)
		}
	case http.MethodDelete:
		// Delete a block
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handleBlockDelete(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Block ID required for deletion"}`, http.StatusBadRequest)
		}
	default:
		http.Error(w, `{"success":false,"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handleBlockCreate handles HTTP requests to create a block
func (api *RestAPI) handleBlockCreate(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to read request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Parse the request body
	var blockData map[string]interface{}
	if err := json.Unmarshal(body, &blockData); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to parse request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	name, ok := blockData["name"].(string)
	if !ok || name == "" {
		http.Error(w, `{"success":false,"error":"Name is required"}`, http.StatusBadRequest)
		return
	}

	content, ok := blockData["content"].(string)
	if !ok {
		content = "" // Default to empty content if not provided
	}

	// Create the block
	block := cmsstore.NewBlock()
	block.SetName(name)
	block.SetContent(content)

	// Set site ID - required field
	siteID, ok := blockData["site_id"].(string)
	if !ok || siteID == "" {
		http.Error(w, `{"success":false,"error":"Site ID is required"}`, http.StatusBadRequest)
		return
	}
	block.SetSiteID(siteID)

	// Set page ID - optional field
	if pageID, ok := blockData["page_id"].(string); ok {
		block.SetPageID(pageID)
	}

	// Set template ID - optional field
	if templateID, ok := blockData["template_id"].(string); ok {
		block.SetTemplateID(templateID)
	} else {
		block.SetTemplateID("") // Set empty string if not provided
	}

	// Set parent ID - optional field
	if parentID, ok := blockData["parent_id"].(string); ok {
		block.SetParentID(parentID)
	} else {
		block.SetParentID("") // Set empty string if not provided
	}

	// Set sequence - optional field with default value 0
	if sequence, ok := blockData["sequence"].(float64); ok {
		block.SetSequenceInt(int(sequence))
	} else {
		block.SetSequenceInt(0) // Default to 0 if not provided
	}

	// Save the block
	if err := api.store.BlockCreate(r.Context(), block); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save block: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the created block
	response := map[string]interface{}{
		"success": true,
		"id":      block.ID(),
		"name":    block.Name(),
		"content": block.Content(),
		"site_id": block.SiteID(),
		"page_id": block.PageID(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleBlockGet handles HTTP requests to get a block by ID
func (api *RestAPI) handleBlockGet(w http.ResponseWriter, r *http.Request, blockID string) {
	// Get the block from the store
	block, err := api.store.BlockFindByID(r.Context(), blockID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find block: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if block == nil {
		http.Error(w, `{"success":false,"error":"Block not found"}`, http.StatusNotFound)
		return
	}

	// Return the block
	response := map[string]interface{}{
		"success": true,
		"block": map[string]interface{}{
			"id":      block.ID(),
			"name":    block.Name(),
			"content": block.Content(),
			"site_id": block.SiteID(),
			"page_id": block.PageID(),
		},
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleBlockList handles HTTP requests to list all blocks
func (api *RestAPI) handleBlockList(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	siteID := r.URL.Query().Get("site_id")

	// Create block query
	query := cmsstore.BlockQuery()
	if siteID != "" {
		query = query.SetSiteID(siteID)
	}

	// Get blocks from the store
	blocks, err := api.store.BlockList(r.Context(), query)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to list blocks: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Convert blocks to response format
	blocksList := make([]map[string]interface{}, 0, len(blocks))
	for _, block := range blocks {
		blockData := map[string]interface{}{
			"id":      block.ID(),
			"name":    block.Name(),
			"content": block.Content(),
			"site_id": block.SiteID(),
		}
		blocksList = append(blocksList, blockData)
	}

	// Return the blocks list
	response := map[string]interface{}{
		"success": true,
		"blocks":  blocksList,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleBlockUpdate handles HTTP requests to update a block
func (api *RestAPI) handleBlockUpdate(w http.ResponseWriter, r *http.Request, blockID string) {
	// Get the existing block
	block, err := api.store.BlockFindByID(r.Context(), blockID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find block: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if block == nil {
		http.Error(w, `{"success":false,"error":"Block not found"}`, http.StatusNotFound)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to read request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Parse the request body
	var updates map[string]interface{}
	if err := json.Unmarshal(body, &updates); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to parse request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok && name != "" {
		block.SetName(name)
	}
	if content, ok := updates["content"].(string); ok {
		block.SetContent(content)
	}

	// Save the updated block
	if err := api.store.BlockUpdate(r.Context(), block); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save block: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the updated block
	response := map[string]interface{}{
		"success": true,
		"id":      block.ID(),
		"name":    block.Name(),
		"content": block.Content(),
		"site_id": block.SiteID(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleBlockDelete handles HTTP requests to delete a block
func (api *RestAPI) handleBlockDelete(w http.ResponseWriter, r *http.Request, blockID string) {
	// Delete the block
	if err := api.store.BlockSoftDeleteByID(r.Context(), blockID); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to delete block: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Block deleted successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
