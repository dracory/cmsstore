package page_update

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore"
)

func handleAjaxLoadContent(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	page, err := store.PageFindByID(r.Context(), pageID)
	if err != nil || page == nil {
		return api.Error("Page not found").ToString()
	}

	return api.SuccessWithData("Content loaded successfully", map[string]any{
		"title":   page.Title(),
		"content": page.Content(),
		"editor":  page.Editor(),
	}).ToString()
}

func handleAjaxSaveContent(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")

	var reqData struct {
		PageID  string `json:"page_id"`
		Title   string `json:"page_title"`
		Content string `json:"page_content"`
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

	if reqData.Title == "" {
		return api.Error("Title is required").ToString()
	}

	page, err := store.PageFindByID(r.Context(), reqData.PageID)
	if err != nil || page == nil {
		return api.Error("Page not found").ToString()
	}

	page.SetTitle(reqData.Title)
	page.SetContent(reqData.Content)

	if err := store.PageUpdate(r.Context(), page); err != nil {
		slog.Error("Failed to save page content", "error", err)
		return api.Error("Failed to save page content").ToString()
	}

	return api.Success("Page saved successfully").ToString()
}
