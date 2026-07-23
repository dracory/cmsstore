package page_update

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore"
)

func handleAjaxLoadSEO(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	page, err := store.PageFindByID(r.Context(), pageID)
	if err != nil || page == nil {
		return api.Error("Page not found").ToString()
	}

	return api.SuccessWithData("SEO data loaded successfully", map[string]any{
		"alias":            page.Alias(),
		"canonical_url":    page.CanonicalUrl(),
		"meta_description": page.MetaDescription(),
		"meta_keywords":    page.MetaKeywords(),
		"meta_robots":      page.MetaRobots(),
	}).ToString()
}

func handleAjaxSaveSEO(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")

	var reqData struct {
		PageID          string `json:"page_id"`
		Alias           string `json:"page_alias"`
		CanonicalURL    string `json:"page_canonical_url"`
		MetaDescription string `json:"page_meta_description"`
		MetaKeywords    string `json:"page_meta_keywords"`
		MetaRobots      string `json:"page_meta_robots"`
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

	page.SetAlias(reqData.Alias)
	page.SetCanonicalUrl(reqData.CanonicalURL)
	page.SetMetaDescription(reqData.MetaDescription)
	page.SetMetaKeywords(reqData.MetaKeywords)
	page.SetMetaRobots(reqData.MetaRobots)

	if err := store.PageUpdate(r.Context(), page); err != nil {
		slog.Error("Failed to save page SEO", "error", err)
		return api.Error("Failed to save page SEO").ToString()
	}

	return api.Success("Page saved successfully").ToString()
}
