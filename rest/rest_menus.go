package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gouniverse/cmsstore"
)

// handleMenusEndpoint handles HTTP requests for the /api/menus endpoint
func (api *RestAPI) handleMenusEndpoint(w http.ResponseWriter, r *http.Request, pathParts []string) {
	switch r.Method {
	case http.MethodPost:
		// Create a new menu
		api.handleMenuCreate(w, r)
	case http.MethodGet:
		// Get menu(s)
		if len(pathParts) > 0 && pathParts[0] != "" {
			// Get a specific menu by ID
			api.handleMenuGet(w, r, pathParts[0])
		} else {
			// List all menus
			api.handleMenuList(w, r)
		}
	case http.MethodPut:
		// Update a menu
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handleMenuUpdate(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Menu ID required for update"}`, http.StatusBadRequest)
		}
	case http.MethodDelete:
		// Delete a menu
		if len(pathParts) > 0 && pathParts[0] != "" {
			api.handleMenuDelete(w, r, pathParts[0])
		} else {
			http.Error(w, `{"success":false,"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, `{"success":false,"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handleMenuCreate handles HTTP requests to create a menu
func (api *RestAPI) handleMenuCreate(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to read request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Parse the request body
	var menuData map[string]interface{}
	if err := json.Unmarshal(body, &menuData); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to parse request body: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	name, ok := menuData["name"].(string)
	if !ok || name == "" {
		http.Error(w, `{"success":false,"error":"Name is required"}`, http.StatusBadRequest)
		return
	}

	// Create the menu
	menu := cmsstore.NewMenu()
	menu.SetName(name)

	// Set site ID - required field
	siteID, ok := menuData["site_id"].(string)
	if !ok || siteID == "" {
		http.Error(w, `{"success":false,"error":"Site ID is required"}`, http.StatusBadRequest)
		return
	}
	menu.SetSiteID(siteID)

	// Save the menu
	if err := api.store.MenuCreate(r.Context(), menu); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save menu: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the created menu
	response := map[string]interface{}{
		"success": true,
		"id":      menu.ID(),
		"name":    menu.Name(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleMenuGet handles HTTP requests to get a menu by ID
func (api *RestAPI) handleMenuGet(w http.ResponseWriter, r *http.Request, menuID string) {
	// Get the menu from the store
	menu, err := api.store.MenuFindByID(r.Context(), menuID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find menu: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if menu == nil {
		http.Error(w, `{"success":false,"error":"Menu not found"}`, http.StatusNotFound)
		return
	}

	// Return the menu
	response := map[string]interface{}{
		"success": true,
		"id":      menu.ID(),
		"name":    menu.Name(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleMenuList handles HTTP requests to list all menus
func (api *RestAPI) handleMenuList(w http.ResponseWriter, r *http.Request) {
	// Get all menus from the store
	menus, err := api.store.MenuList(r.Context(), cmsstore.MenuQuery())
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to list menus: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Convert menus to response format
	menusList := make([]map[string]interface{}, 0, len(menus))
	for _, menu := range menus {
		menuData := map[string]interface{}{
			"id":   menu.ID(),
			"name": menu.Name(),
		}
		menusList = append(menusList, menuData)
	}

	// Return the menus list
	response := map[string]interface{}{
		"success": true,
		"menus":   menusList,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleMenuUpdate handles HTTP requests to update a menu
func (api *RestAPI) handleMenuUpdate(w http.ResponseWriter, r *http.Request, menuID string) {
	// Get the existing menu
	menu, err := api.store.MenuFindByID(r.Context(), menuID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to find menu: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if menu == nil {
		http.Error(w, `{"success":false,"error":"Menu not found"}`, http.StatusNotFound)
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
		menu.SetName(name)
	}

	// Save the updated menu
	if err := api.store.MenuUpdate(r.Context(), menu); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to save menu: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the updated menu
	response := map[string]interface{}{
		"success": true,
		"id":      menu.ID(),
		"name":    menu.Name(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// handleMenuDelete handles HTTP requests to delete a menu
func (api *RestAPI) handleMenuDelete(w http.ResponseWriter, r *http.Request, menuID string) {
	// Delete the menu
	if err := api.store.MenuSoftDeleteByID(r.Context(), menuID); err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to delete menu: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Menu deleted successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to create response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
