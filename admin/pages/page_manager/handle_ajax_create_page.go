package page_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore"
)

func handleAjaxCreatePage(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	var reqData struct {
		Name   string `json:"name"`
		SiteID string `json:"site_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.Name == "" {
		return api.Error("Name is required").ToString()
	}

	if reqData.SiteID == "" {
		return api.Error("Site is required").ToString()
	}

	page := cmsstore.NewPage()
	page.SetName(reqData.Name)
	page.SetSiteID(reqData.SiteID)
	page.SetStatus(cmsstore.PAGE_STATUS_DRAFT)

	if err := store.PageCreate(ctx, page); err != nil {
		slog.Error("Failed to create page", "error", err)
		return api.Error("Failed to create page: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Page created successfully", map[string]any{
		"id": page.ID(),
	}).ToString()
}
