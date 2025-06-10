package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gouniverse/cmsstore"
)

// handlePagesEndpoint handles HTTP requests for the /api/pages endpoint
func (api *RestAPI) handlePagesEndpoint(w http.ResponseWriter, r *http.Request, pathParts []string) {
	switch r.Method {
	case http.MethodPost:
		// Create a new page
		api.handlePageCreate(w, r)
	case http.MethodGet:
		// Get page(s)
		if len(pathParts) > 0 && pathParts[0] != "" {
			// Get a specific page by ID
			api.handlePageGet(w, r, pathParts[0])
		} else {
			// List all pages
			api.handlePageList(w, r)
		}
	case http.MethodPut:
		// Update a page
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handlePageUpdate(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Page ID required for update"}`, http.StatusBadRequest)
		}
	case http.MethodDelete:
		// Delete a page
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handlePageDelete(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Page ID required for deletion"}`, http.StatusBadRequest)
		}
	default:
		http.Error(w, `{"success":false,"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handlePageCreate handles HTTP requests to create a page
func (api *RestAPI) handlePageCreate(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to read request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Parse the request body
	var pageData map[string]interface{}
	if err := json.Unmarshal(body, &pageData); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to parse request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	title, ok := pageData["title"].(string)
	if !ok || title == "" {
		http.Error(w, `{"success":false,"error":"Title is required"}`, http.StatusBadRequest)
		return
	}

	content, _ := pageData["content"].(string) // Content is optional
	status, _ := pageData["status"].(string)   // Status is optional

	// Create the page
	page := cmsstore.NewPage()
	page.SetTitle(title)
	page.SetContent(content)

	// Set site ID - required field
	siteID, ok := pageData["site_id"].(string)
	if !ok || siteID == "" {
		http.Error(w, `{"success":false,"error":"Site ID is required"}`, http.StatusBadRequest)
		return
	}
	page.SetSiteID(siteID)

	if status != "" {
		page.SetStatus(status)
	} else {
		page.SetStatus("draft") // Default status
	}

	// Save the page
	if err := api.store.PageCreate(r.Context(), page); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save page: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the created page
	response := map[string]interface{}{
		"success": true,
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handlePageGet handles HTTP requests to get a page by ID
func (api *RestAPI) handlePageGet(w http.ResponseWriter, r *http.Request, pageID string) {
	// Get the page from the store
	page, err := api.store.PageFindByID(r.Context(), pageID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find page: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if page == nil {
		http.Error(w, `{"success":false,"error":"Page not found"}`, http.StatusNotFound)
		return
	}

	// Return the page
	response := map[string]interface{}{
		"success": true,
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handlePageList handles HTTP requests to list all pages
func (api *RestAPI) handlePageList(w http.ResponseWriter, r *http.Request) {
	// Get all pages from the store
	pages, err := api.store.PageList(r.Context(), cmsstore.PageQuery())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to list pages: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Convert pages to response format
	pagesList := make([]map[string]interface{}, 0, len(pages))
	for _, page := range pages {
		pageData := map[string]interface{}{
			"id":      page.ID(),
			"title":   page.Title(),
			"content": page.Content(),
			"status":  page.Status(),
		}
		pagesList = append(pagesList, pageData)
	}

	// Return the pages list
	response := map[string]interface{}{
		"success": true,
		"pages":   pagesList,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handlePageUpdate handles HTTP requests to update a page
func (api *RestAPI) handlePageUpdate(w http.ResponseWriter, r *http.Request, pageID string) {
	// Get the existing page
	page, err := api.store.PageFindByID(r.Context(), pageID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find page: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if page == nil {
		http.Error(w, `{"success":false,"error":"Page not found"}`, http.StatusNotFound)
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
	if title, ok := updates["title"].(string); ok && title != "" {
		page.SetTitle(title)
	}
	if content, ok := updates["content"].(string); ok {
		page.SetContent(content)
	}
	if status, ok := updates["status"].(string); ok && status != "" {
		page.SetStatus(status)
	}

	// Save the updated page
	if err := api.store.PageUpdate(r.Context(), page); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save page: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the updated page
	response := map[string]interface{}{
		"success": true,
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handlePageDelete handles HTTP requests to delete a page
func (api *RestAPI) handlePageDelete(w http.ResponseWriter, r *http.Request, pageID string) {
	// Delete the page
	if err := api.store.PageSoftDeleteByID(r.Context(), pageID); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to delete page: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Page deleted successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
