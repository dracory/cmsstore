package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dracory/cmsstore"
)

func (m *MCP) toolSiteGet(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	site, err := m.store.SiteFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if site == nil {
		return "", errors.New("site not found")
	}

	domains, _ := site.DomainNames()

	respBytes, err := json.Marshal(map[string]any{
		"id":          site.ID(),
		"name":        site.Name(),
		"handle":      site.Handle(),
		"domainNames": domains,
		"status":      site.Status(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolSiteList(ctx context.Context, args map[string]any) (string, error) {
	q := cmsstore.SiteQuery()

	if v, ok := args["status"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetStatus(v)
	}
	if v, ok := args["name_like"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetNameLike(v)
	}
	if v, ok := args["domain_name"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetDomainName(v)
	}
	if v, ok := args["handle"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetHandle(v)
	}
	if v, ok := argBool(args, "include_soft_deleted"); ok {
		q.SetSoftDeletedIncluded(v)
	}
	if v, ok := argInt(args, "limit"); ok {
		q.SetLimit(int(v))
	}
	if v, ok := argInt(args, "offset"); ok {
		q.SetOffset(int(v))
	}
	if v, ok := args["order_by"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetOrderBy(v)
	}
	if v, ok := args["sort_order"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetSortOrder(v)
	}

	sites, err := m.store.SiteList(ctx, q)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(sites))
	for _, s := range sites {
		if s == nil {
			continue
		}
		domains, _ := s.DomainNames()
		items = append(items, map[string]any{
			"id":          s.ID(),
			"name":        s.Name(),
			"handle":      s.Handle(),
			"status":      s.Status(),
			"domainNames": domains,
			"created_at":  s.CreatedAt(),
			"updated_at":  s.UpdatedAt(),
		})
	}

	respBytes, err := json.Marshal(map[string]any{
		"items": items,
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolSiteUpsert(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	name, _ := args["name"].(string)
	handle, _ := args["handle"].(string)
	status, _ := args["status"].(string)

	if strings.TrimSpace(name) == "" {
		return "", errors.New("missing required parameter: name")
	}

	var site cmsstore.SiteInterface
	var err error

	// If ID is provided, try to find existing site
	if strings.TrimSpace(id) != "" {
		site, err = m.store.SiteFindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if site == nil {
			return "", errors.New("site not found")
		}
	} else {
		// Create new site
		site = cmsstore.NewSite()
	}

	// Set/update fields
	site.SetName(name)
	if handle != "" {
		site.SetHandle(handle)
	}
	if status != "" {
		site.SetStatus(status)
	}
	if v := strings.TrimSpace(argString(args, "memo")); v != "" {
		site.SetMemo(v)
	}

	// Handle domain names
	if domainsInterface, ok := args["domain_names"].([]any); ok {
		domains := make([]string, 0, len(domainsInterface))
		for _, domain := range domainsInterface {
			if domainStr, ok := domain.(string); ok {
				domains = append(domains, domainStr)
			}
		}
		if len(domains) > 0 {
			site.SetDomainNames(domains)
		}
	}

	// Save site
	if strings.TrimSpace(id) != "" {
		// Update existing site
		if err := m.store.SiteUpdate(ctx, site); err != nil {
			return "", err
		}
	} else {
		// Create new site
		if err := m.store.SiteCreate(ctx, site); err != nil {
			return "", err
		}
	}

	// Create versioning record if versioning is enabled
	if m.store.VersioningEnabled() {
		if err := m.createSiteVersioning(ctx, site); err != nil {
			// Log error but don't fail the operation
			// In a production environment, you might want to handle this differently
		}
	}

	domains, _ := site.DomainNames()

	respBytes, err := json.Marshal(map[string]any{
		"id":          site.ID(),
		"name":        site.Name(),
		"handle":      site.Handle(),
		"status":      site.Status(),
		"domainNames": domains,
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolSiteDelete(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	site, err := m.store.SiteFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if site == nil {
		return "", errors.New("site not found")
	}

	if err := m.store.SiteSoftDeleteByID(ctx, id); err != nil {
		if err := m.store.SiteDelete(ctx, site); err != nil {
			return "", err
		}
	}

	respBytes, err := json.Marshal(map[string]any{"id": id})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

// createSiteVersioning creates a versioning record for a site if versioning is enabled
func (m *MCP) createSiteVersioning(ctx context.Context, site cmsstore.SiteInterface) error {
	if !m.store.VersioningEnabled() {
		return nil
	}

	if site == nil {
		return errors.New("site is nil")
	}

	// Get last versioning to check if content has changed
	lastVersioningList, err := m.store.VersioningList(ctx, cmsstore.NewVersioningQuery().
		SetEntityType(cmsstore.VERSIONING_TYPE_SITE).
		SetEntityID(site.ID()).
		SetOrderBy("created_at").
		SetSortOrder("DESC").
		SetLimit(1))

	if err != nil {
		return err
	}

	// Marshal site content for versioning
	domains, _ := site.DomainNames()
	siteData := map[string]any{
		"id":          site.ID(),
		"name":        site.Name(),
		"handle":      site.Handle(),
		"status":      site.Status(),
		"domainNames": domains,
		"memo":        site.Memo(),
	}
	content, err := json.Marshal(siteData)
	if err != nil {
		return err
	}

	// Check if last versioning has the same content
	if len(lastVersioningList) > 0 {
		lastVersioning := lastVersioningList[0]
		if lastVersioning.Content() == string(content) {
			return nil // No change needed
		}
	}

	// Create new versioning record
	return m.store.VersioningCreate(ctx, cmsstore.NewVersioning().
		SetEntityID(site.ID()).
		SetEntityType(cmsstore.VERSIONING_TYPE_SITE).
		SetContent(string(content)))
}
