package page_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore"
	"github.com/samber/lo"
)

func handleAjaxLoadPages(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	var reqData struct {
		Page      int    `json:"page"`
		PerPage   int    `json:"per_page"`
		SortOrder string `json:"sort_order"`
		SortBy    string `json:"sort_by"`
		Status    string `json:"status"`
		Search    string `json:"search"`
		SiteID    string `json:"site_id"`
		DateFrom  string `json:"date_from"`
		DateTo    string `json:"date_to"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	ctx := r.Context()

	// Load sites for the create modal dropdown
	sitesForDropdown, err := store.SiteList(ctx, cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(cmsstore.SORT_ORDER_ASC).
		SetOffset(0).
		SetLimit(100))
	if err != nil {
		slog.Error("Failed to load sites for dropdown", "error", err)
		sitesForDropdown = []cmsstore.SiteInterface{}
	}

	page := reqData.Page
	perPage := reqData.PerPage
	if perPage == 0 {
		perPage = 20
	}
	sortOrder := reqData.SortOrder
	if sortOrder == "" {
		sortOrder = cmsstore.SORT_ORDER_DESC
	}
	sortBy := reqData.SortBy
	if sortBy == "" {
		sortBy = cmsstore.COLUMN_CREATED_AT
	}

	query := cmsstore.PageQuery().
		SetLimit(perPage).
		SetOffset(page * perPage).
		SetOrderBy(sortBy).
		SetSortOrder(sortOrder)

	if reqData.Status != "" {
		query = query.SetStatus(reqData.Status)
	}

	if reqData.SiteID != "" {
		query = query.SetSiteID(reqData.SiteID)
	}

	if reqData.Search != "" {
		query = query.SetNameLike(reqData.Search)
	}

	if reqData.DateFrom != "" {
		query = query.SetCreatedAtGte(reqData.DateFrom + " 00:00:00")
	}

	if reqData.DateTo != "" {
		query = query.SetCreatedAtLte(reqData.DateTo + " 23:59:59")
	}

	pages, err := store.PageList(ctx, query)
	if err != nil {
		slog.Error("Failed to load pages", "error", err)
		return api.Error("Failed to load pages").ToString()
	}

	// Load sites for display
	sites, err := store.SiteList(ctx, cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(cmsstore.SORT_ORDER_ASC).
		SetOffset(0).
		SetLimit(100))
	if err != nil {
		slog.Error("Failed to load sites", "error", err)
		sites = []cmsstore.SiteInterface{}
	}

	siteMap := lo.Associate(sites, func(site cmsstore.SiteInterface) (string, string) {
		return site.ID(), site.Name()
	})

	// Build a map of site ID -> first domain name for live URL construction
	siteDomainMap := map[string]string{}
	for _, site := range sites {
		domainNames, err := site.DomainNames()
		if err == nil && len(domainNames) > 0 {
			siteDomainMap[site.ID()] = domainNames[0]
		}
	}

	pageList := []map[string]any{}
	for _, p := range pages {
		siteName := lo.IfF(siteMap != nil, func() string {
			if name, ok := siteMap[p.SiteID()]; ok {
				return name
			}
			return ""
		}).Else("")

		liveURL := ""
		if domain, ok := siteDomainMap[p.SiteID()]; ok {
			alias := p.Alias()
			if alias != "" {
				if !strings.HasPrefix(alias, "/") {
					alias = "/" + alias
				}
				if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
					if strings.HasPrefix(domain, "localhost") || strings.HasSuffix(domain, ".local") {
						domain = "http://" + domain
					} else {
						domain = "https://" + domain
					}
				}
				liveURL = strings.TrimSuffix(domain, "/") + alias
			}
		}

		pageList = append(pageList, map[string]any{
			"id":         p.ID(),
			"name":       p.Name(),
			"alias":      p.Alias(),
			"status":     p.Status(),
			"site_id":    p.SiteID(),
			"site_name":  siteName,
			"live_url":   liveURL,
			"created_at": p.CreatedAt(),
			"updated_at": p.UpdatedAt(),
		})
	}

	count, err := store.PageCount(ctx, query)
	if err != nil {
		slog.Error("Failed to get pages count", "error", err)
		return api.Error("Failed to get pages count").ToString()
	}

	return api.SuccessWithData("Pages loaded successfully", map[string]any{
		"pages": pageList,
		"total": count,
		"sites": siteListForResponse(sitesForDropdown),
	}).ToString()
}

func siteListForResponse(sites []cmsstore.SiteInterface) []map[string]any {
	result := []map[string]any{}
	for _, s := range sites {
		result = append(result, map[string]any{
			"id":   s.ID(),
			"name": s.Name() + " (" + string(s.Status()) + ")",
		})
	}
	return result
}
