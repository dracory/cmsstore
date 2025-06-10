package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gouniverse/cmsstore"
)

// handleTemplatesEndpoint handles HTTP requests for the /api/templates endpoint
func (api *RestAPI) handleTemplatesEndpoint(w http.ResponseWriter, r *http.Request, pathParts []string) {
	switch r.Method {
	case http.MethodPost:
		// Create a new template
		api.handleTemplateCreate(w, r)
	case http.MethodGet:
		// Get template(s)
		if len(pathParts) > 0 && pathParts[0] != "" {
			// Get a specific template by ID
			api.handleTemplateGet(w, r, pathParts[0])
		} else {
			// List all templates
			api.handleTemplateList(w, r)
		}
	case http.MethodPut:
		// Update a template
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handleTemplateUpdate(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Template ID required for update"}`, http.StatusBadRequest)
		}
	case http.MethodDelete:
		// Delete a template
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handleTemplateDelete(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Template ID required for deletion"}`, http.StatusBadRequest)
		}
	default:
		http.Error(w, `{"success":false,"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handleTemplateCreate handles HTTP requests to create a template
func (api *RestAPI) handleTemplateCreate(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to read request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Parse the request body
	var templateData map[string]interface{}
	if err := json.Unmarshal(body, &templateData); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to parse request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	name, ok := templateData["name"].(string)
	if !ok || name == "" {
		http.Error(w, `{"success":false,"error":"Name is required"}`, http.StatusBadRequest)
		return
	}

	content, ok := templateData["content"].(string)
	if !ok {
		content = "" // Default to empty content if not provided
	}

	// Create the template
	template := cmsstore.NewTemplate()
	template.SetName(name)
	template.SetContent(content)

	// Set site ID - required field
	siteID, ok := templateData["site_id"].(string)
	if !ok || siteID == "" {
		http.Error(w, `{"success":false,"error":"Site ID is required"}`, http.StatusBadRequest)
		return
	}
	template.SetSiteID(siteID)

	// Save the template
	if err := api.store.TemplateCreate(r.Context(), template); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save template: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the created template
	response := map[string]interface{}{
		"success": true,
		"id":      template.ID(),
		"name":    template.Name(),
		"content": template.Content(),
		"site_id": template.SiteID(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleTemplateGet handles HTTP requests to get a template by ID
func (api *RestAPI) handleTemplateGet(w http.ResponseWriter, r *http.Request, templateID string) {
	// Get the template from the store
	template, err := api.store.TemplateFindByID(r.Context(), templateID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find template: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if template == nil {
		http.Error(w, `{"success":false,"error":"Template not found"}`, http.StatusNotFound)
		return
	}

	// Return the template
	response := map[string]interface{}{
		"success": true,
		"id":      template.ID(),
		"name":    template.Name(),
		"content": template.Content(),
		"site_id": template.SiteID(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleTemplateList handles HTTP requests to list all templates
func (api *RestAPI) handleTemplateList(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	siteID := r.URL.Query().Get("site_id")

	// Create template query
	query := cmsstore.TemplateQuery()
	if siteID != "" {
		query = query.SetSiteID(siteID)
	}

	// Get templates from the store
	templates, err := api.store.TemplateList(r.Context(), query)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to list templates: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Convert templates to response format
	templatesList := make([]map[string]interface{}, 0, len(templates))
	for _, template := range templates {
		templateData := map[string]interface{}{
			"id":      template.ID(),
			"name":    template.Name(),
			"content": template.Content(),
			"site_id": template.SiteID(),
		}
		templatesList = append(templatesList, templateData)
	}

	// Return the templates list
	response := map[string]interface{}{
		"success":   true,
		"templates": templatesList,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleTemplateUpdate handles HTTP requests to update a template
func (api *RestAPI) handleTemplateUpdate(w http.ResponseWriter, r *http.Request, templateID string) {
	// Get the existing template
	template, err := api.store.TemplateFindByID(r.Context(), templateID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find template: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if template == nil {
		http.Error(w, `{"success":false,"error":"Template not found"}`, http.StatusNotFound)
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
		template.SetName(name)
	}
	if content, ok := updates["content"].(string); ok {
		template.SetContent(content)
	}

	// Save the updated template
	if err := api.store.TemplateUpdate(r.Context(), template); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save template: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the updated template
	response := map[string]interface{}{
		"success": true,
		"id":      template.ID(),
		"name":    template.Name(),
		"content": template.Content(),
		"site_id": template.SiteID(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleTemplateDelete handles HTTP requests to delete a template
func (api *RestAPI) handleTemplateDelete(w http.ResponseWriter, r *http.Request, templateID string) {
	// Delete the template
	if err := api.store.TemplateSoftDeleteByID(r.Context(), templateID); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to delete template: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Template deleted successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
