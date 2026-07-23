package page_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore"
)

func handleAjaxDeletePage(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	var reqData struct {
		PageID string `json:"page_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.PageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	page, err := store.PageFindByID(ctx, reqData.PageID)
	if err != nil {
		slog.Error("Failed to find page for delete", "error", err)
		return api.Error("Page not found").ToString()
	}

	if page == nil {
		return api.Error("Page not found").ToString()
	}

	if err := store.PageDelete(ctx, page); err != nil {
		slog.Error("Failed to delete page", "error", err)
		return api.Error("Failed to delete page").ToString()
	}

	return api.SuccessWithData("Page deleted successfully", map[string]any{}).ToString()
}
