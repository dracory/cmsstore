package page_update

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore"
)

func handleAjaxLoadMiddlewares(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	page, err := store.PageFindByID(r.Context(), pageID)
	if err != nil || page == nil {
		return api.Error("Page not found").ToString()
	}

	return api.SuccessWithData("Middlewares loaded successfully", map[string]any{
		"before": page.MiddlewaresBefore(),
		"after":  page.MiddlewaresAfter(),
	}).ToString()
}

func handleAjaxSaveMiddlewares(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")

	var reqData struct {
		PageID            string   `json:"page_id"`
		MiddlewaresBefore []string `json:"middlewares_before"`
		MiddlewaresAfter  []string `json:"middlewares_after"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.PageID == "" {
		reqData.PageID = pageID
	}

	if reqData.PageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	page, err := store.PageFindByID(r.Context(), reqData.PageID)
	if err != nil || page == nil {
		return api.Error("Page not found").ToString()
	}

	page.SetMiddlewaresBefore(reqData.MiddlewaresBefore)
	page.SetMiddlewaresAfter(reqData.MiddlewaresAfter)

	if err := store.PageUpdate(r.Context(), page); err != nil {
		slog.Error("Failed to save page middlewares", "error", err)
		return api.Error("Failed to save page middlewares").ToString()
	}

	return api.Success("Middlewares saved successfully").ToString()
}
