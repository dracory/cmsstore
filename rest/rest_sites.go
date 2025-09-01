package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dracory/cmsstore"
)

// handleSitesEndpoint handles HTTP requests for the /api/sites endpoint
func (api *RestAPI) handleSitesEndpoint(w http.ResponseWriter, r *http.Request, pathParts []string) {
	switch r.Method {
	case http.MethodPost:
		// Create a new site
		api.handleSiteCreate(w, r)
	case http.MethodGet:
		// Get site(s)
		if len(pathParts) > 0 && pathParts[0] != "" {
			// Get a specific site by ID
			api.handleSiteGet(w, r, pathParts[0])
		} else {
			// List all sites
			api.handleSiteList(w, r)
		}
	case http.MethodPut:
		// Update a site
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handleSiteUpdate(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Site ID required for update"}`, http.StatusBadRequest)
		}
	case http.MethodDelete:
		// Delete a site
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handleSiteDelete(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Site ID required for deletion"}`, http.StatusBadRequest)
		}
	default:
		http.Error(w, `{"success":false,"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handleSiteCreate handles HTTP requests to create a site
func (api *RestAPI) handleSiteCreate(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to read request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Parse the request body
	var siteData map[string]interface{}
	if err := json.Unmarshal(body, &siteData); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to parse request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	name, ok := siteData["name"].(string)
	if !ok || name == "" {
		http.Error(w, `{"success":false,"error":"Name is required"}`, http.StatusBadRequest)
		return
	}

	// Create the site
	site := cmsstore.NewSite()
	site.SetName(name)

	// Handle domain names
	var domainNames []string
	if domains, ok := siteData["domain_names"].([]interface{}); ok {
		for _, domain := range domains {
			if domainStr, ok := domain.(string); ok && domainStr != "" {
				domainNames = append(domainNames, domainStr)
			}
		}
	} else if domain, ok := siteData["domain_name"].(string); ok && domain != "" {
		domainNames = []string{domain}
	}

	if len(domainNames) > 0 {
		site, err = site.SetDomainNames(domainNames)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to set domain names: %v"}`, err), http.StatusBadRequest)
			return
		}
	}

	// Save the site
	if err := api.store.SiteCreate(r.Context(), site); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save site: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Get domain names for response
	domainNames, err = site.DomainNames()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to get domain names: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the created site
	response := map[string]interface{}{
		"success":      true,
		"id":           site.ID(),
		"name":         site.Name(),
		"domain_names": domainNames,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleSiteGet handles HTTP requests to get a site by ID
func (api *RestAPI) handleSiteGet(w http.ResponseWriter, r *http.Request, siteID string) {
	// Get the site from the store
	site, err := api.store.SiteFindByID(r.Context(), siteID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find site: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if site == nil {
		http.Error(w, `{"success":false,"error":"Site not found"}`, http.StatusNotFound)
		return
	}

	// Get domain names
	domainNames, err := site.DomainNames()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to get domain names: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the site
	response := map[string]interface{}{
		"success":      true,
		"id":           site.ID(),
		"name":         site.Name(),
		"domain_names": domainNames,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleSiteList handles HTTP requests to list all sites
func (api *RestAPI) handleSiteList(w http.ResponseWriter, r *http.Request) {
	// Get all sites from the store
	sites, err := api.store.SiteList(r.Context(), cmsstore.SiteQuery())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to list sites: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Convert sites to response format
	sitesList := make([]map[string]interface{}, 0, len(sites))
	for _, site := range sites {
		domainNames, err := site.DomainNames()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to get domain names: %v"}`, err), http.StatusInternalServerError)
			return
		}

		siteData := map[string]interface{}{
			"id":           site.ID(),
			"name":         site.Name(),
			"domain_names": domainNames,
		}
		sitesList = append(sitesList, siteData)
	}

	// Return the sites list
	response := map[string]interface{}{
		"success": true,
		"sites":   sitesList,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleSiteUpdate handles HTTP requests to update a site
func (api *RestAPI) handleSiteUpdate(w http.ResponseWriter, r *http.Request, siteID string) {
	// Get the existing site
	site, err := api.store.SiteFindByID(r.Context(), siteID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find site: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if site == nil {
		http.Error(w, `{"success":false,"error":"Site not found"}`, http.StatusNotFound)
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
		site.SetName(name)
	}

	// Handle domain names update
	var domainNames []string
	updateDomains := false

	if domains, ok := updates["domain_names"].([]interface{}); ok {
		updateDomains = true
		for _, domain := range domains {
			if domainStr, ok := domain.(string); ok && domainStr != "" {
				domainNames = append(domainNames, domainStr)
			}
		}
	} else if domain, ok := updates["domain_name"].(string); ok && domain != "" {
		updateDomains = true
		domainNames = []string{domain}
	}

	if updateDomains {
		site, err = site.SetDomainNames(domainNames)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to set domain names: %v"}`, err), http.StatusBadRequest)
			return
		}
	}

	// Save the updated site
	if err := api.store.SiteUpdate(r.Context(), site); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save site: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Get domain names for response
	domainNames, err = site.DomainNames()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to get domain names: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the updated site
	response := map[string]interface{}{
		"success":      true,
		"id":           site.ID(),
		"name":         site.Name(),
		"domain_names": domainNames,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleSiteDelete handles HTTP requests to delete a site
func (api *RestAPI) handleSiteDelete(w http.ResponseWriter, r *http.Request, siteID string) {
	// Delete the site
	if err := api.store.SiteSoftDeleteByID(r.Context(), siteID); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to delete site: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Site deleted successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
