package page_update

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore"
)

func handleAjaxLoadSettings(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	page, err := store.PageFindByID(r.Context(), pageID)
	if err != nil || page == nil {
		return api.Error("Page not found").ToString()
	}

	sites, err := store.SiteList(r.Context(), cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(cmsstore.SORT_ORDER_ASC).
		SetOffset(0).
		SetLimit(100))
	if err != nil {
		slog.Error("Failed to load sites", "error", err)
		sites = []cmsstore.SiteInterface{}
	}

	templates, err := store.TemplateList(r.Context(), cmsstore.TemplateQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(cmsstore.SORT_ORDER_ASC).
		SetOffset(0).
		SetLimit(100))
	if err != nil {
		slog.Error("Failed to load templates", "error", err)
		templates = []cmsstore.TemplateInterface{}
	}

	siteList := []map[string]any{}
	for _, site := range sites {
		siteList = append(siteList, map[string]any{
			"id":     site.ID(),
			"name":   site.Name(),
			"status": site.Status(),
		})
	}

	templateList := []map[string]any{}
	for _, tpl := range templates {
		templateList = append(templateList, map[string]any{
			"id":   tpl.ID(),
			"name": tpl.Name(),
		})
	}

	return api.SuccessWithData("Settings loaded successfully", map[string]any{
		"status":      page.Status(),
		"template_id": page.TemplateID(),
		"editor":      page.Editor(),
		"name":        page.Name(),
		"site_id":     page.SiteID(),
		"memo":        page.Memo(),
		"sites":       siteList,
		"templates":   templateList,
	}).ToString()
}

func handleAjaxSaveSettings(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")

	var reqData struct {
		PageID     string `json:"page_id"`
		Status     string `json:"page_status"`
		TemplateID string `json:"page_template_id"`
		Editor     string `json:"page_editor"`
		Name       string `json:"page_name"`
		SiteID     string `json:"page_site_id"`
		Memo       string `json:"page_memo"`
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

	if reqData.Status == "" {
		return api.Error("Status is required").ToString()
	}

	page, err := store.PageFindByID(r.Context(), reqData.PageID)
	if err != nil || page == nil {
		return api.Error("Page not found").ToString()
	}

	page.SetStatus(reqData.Status)
	page.SetTemplateID(reqData.TemplateID)
	page.SetEditor(reqData.Editor)
	page.SetName(reqData.Name)
	page.SetSiteID(reqData.SiteID)
	page.SetMemo(reqData.Memo)

	if err := store.PageUpdate(r.Context(), page); err != nil {
		slog.Error("Failed to save page settings", "error", err)
		return api.Error("Failed to save page settings").ToString()
	}

	return api.Success("Page saved successfully").ToString()
}
