package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dracory/cmsstore"
	"github.com/dromara/carbon/v2"
)

func (m *MCP) toolTemplateGet(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	template, err := m.store.TemplateFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if template == nil {
		return "", errors.New("template not found")
	}

	metas, err := template.Metas()
	if err != nil {
		return "", err
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":         cmsstore.ShortenID(template.ID()),
		"name":       template.Name(),
		"handle":     template.Handle(),
		"status":     template.Status(),
		"editor":     template.Editor(),
		"content":    template.Content(),
		"site_id":    cmsstore.ShortenID(template.SiteID()),
		"memo":       template.Memo(),
		"created_at": template.CreatedAtCarbon().ToDateTimeString(carbon.UTC),
		"updated_at": template.UpdatedAtCarbon().ToDateTimeString(carbon.UTC),
		// "soft_deleted_at": template.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC), // commented out to avoid confusing LLMs since list operations exclude soft deleted items by default
		"metas": metas,
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolTemplateList(ctx context.Context, args map[string]any) (string, error) {
	q := cmsstore.TemplateQuery()

	if v, ok := args["site_id"].(string); ok && strings.TrimSpace(v) != "" {
		// Unshorten the site_id before setting it in the query
		// The database stores full IDs, so we need to unshorten shortened IDs
		siteID := cmsstore.UnshortenID(v)
		q.SetSiteID(siteID)
	}
	if v, ok := args["status"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetStatus(v)
	}
	if v, ok := args["name_like"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetNameLike(v)
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

	templates, err := m.store.TemplateList(ctx, q)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(templates))
	for _, template := range templates {
		if template == nil {
			continue
		}
		metas, err := template.Metas()
		if err != nil {
			return "", err
		}
		items = append(items, map[string]any{
			"id":         cmsstore.ShortenID(template.ID()),
			"name":       template.Name(),
			"handle":     template.Handle(),
			"status":     template.Status(),
			"site_id":    cmsstore.ShortenID(template.SiteID()),
			"editor":     template.Editor(),
			"memo":       template.Memo(),
			"created_at": template.CreatedAtCarbon().ToDateTimeString(carbon.UTC),
			"updated_at": template.UpdatedAtCarbon().ToDateTimeString(carbon.UTC),
			// "soft_deleted_at": template.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC), // commented out to avoid confusing LLMs since list operations exclude soft deleted items by default
			"metas": metas,
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

func (m *MCP) toolTemplateUpsert(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	name, _ := args["name"].(string)
	content, _ := args["content"].(string)
	status, _ := args["status"].(string)
	siteID, _ := args["site_id"].(string)

	if strings.TrimSpace(name) == "" {
		return "", errors.New("missing required parameter: name")
	}

	var template cmsstore.TemplateInterface
	var err error

	// If ID is provided, try to find existing template
	if strings.TrimSpace(id) != "" {
		template, err = m.store.TemplateFindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if template == nil {
			return "", errors.New("template not found")
		}
	} else {
		// Create new template
		template = cmsstore.NewTemplate()

		// Set site ID if not provided
		if siteID == "" {
			sites, err := m.store.SiteList(ctx, cmsstore.SiteQuery())
			if err != nil {
				return "", err
			}
			if len(sites) > 0 {
				siteID = sites[0].ID()
			} else {
				defaultSite := cmsstore.NewSite()
				defaultSite.SetName("Default Site")
				defaultSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
				if err := m.store.SiteCreate(ctx, defaultSite); err != nil {
					return "", err
				}
				siteID = defaultSite.ID()
			}
		}
		template.SetSiteID(siteID)
	}

	// Set/update fields
	template.SetName(name)
	if content != "" {
		template.SetContent(content)
	}
	if status != "" {
		template.SetStatus(status)
	}
	if v := strings.TrimSpace(argString(args, "site_id")); v != "" {
		template.SetSiteID(v)
	}
	if v := strings.TrimSpace(argString(args, "handle")); v != "" {
		template.SetHandle(v)
	}
	if v := strings.TrimSpace(argString(args, "editor")); v != "" {
		template.SetEditor(v)
	}
	if v := strings.TrimSpace(argString(args, "memo")); v != "" {
		template.SetMemo(v)
	}

	// Save template
	if strings.TrimSpace(id) != "" {
		// Update existing template
		if err := m.store.TemplateUpdate(ctx, template); err != nil {
			return "", err
		}
	} else {
		// Create new template
		if err := m.store.TemplateCreate(ctx, template); err != nil {
			return "", err
		}
	}

	metas, err := template.Metas()
	if err != nil {
		return "", err
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":         cmsstore.ShortenID(template.ID()),
		"name":       template.Name(),
		"content":    template.Content(),
		"status":     template.Status(),
		"site_id":    cmsstore.ShortenID(template.SiteID()),
		"handle":     template.Handle(),
		"editor":     template.Editor(),
		"memo":       template.Memo(),
		"created_at": template.CreatedAtCarbon().ToDateTimeString(carbon.UTC),
		"updated_at": template.UpdatedAtCarbon().ToDateTimeString(carbon.UTC),
		// "soft_deleted_at": template.SoftDeletedAtCarbon().ToDateTimeString(carbon.UTC), // commented out to avoid confusing LLMs since list operations exclude soft deleted items by default
		"metas": metas,
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolTemplateDelete(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	template, err := m.store.TemplateFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if template == nil {
		return "", errors.New("template not found")
	}

	if err := m.store.TemplateSoftDeleteByID(ctx, id); err != nil {
		if err := m.store.TemplateDelete(ctx, template); err != nil {
			return "", err
		}
	}

	respBytes, err := json.Marshal(map[string]any{"id": cmsstore.ShortenID(id)})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}
