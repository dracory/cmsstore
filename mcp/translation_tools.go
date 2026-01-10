package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dracory/cmsstore"
)

func (m *MCP) toolTranslationGet(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	translation, err := m.store.TranslationFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if translation == nil {
		return "", errors.New("translation not found")
	}

	content, err := translation.Content()
	if err != nil {
		return "", err
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      translation.ID(),
		"name":    translation.Name(),
		"handle":  translation.Handle(),
		"content": content,
		"status":  translation.Status(),
		"site_id": translation.SiteID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolTranslationList(ctx context.Context, args map[string]any) (string, error) {
	q := cmsstore.TranslationQuery()

	if v, ok := args["site_id"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetSiteID(v)
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

	translations, err := m.store.TranslationList(ctx, q)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(translations))
	for _, translation := range translations {
		if translation == nil {
			continue
		}
		items = append(items, map[string]any{
			"id":         translation.ID(),
			"name":       translation.Name(),
			"handle":     translation.Handle(),
			"status":     translation.Status(),
			"site_id":    translation.SiteID(),
			"created_at": translation.CreatedAt(),
			"updated_at": translation.UpdatedAt(),
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

func (m *MCP) toolTranslationUpsert(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	name, _ := args["name"].(string)
	handle, _ := args["handle"].(string)
	status, _ := args["status"].(string)
	siteID, _ := args["site_id"].(string)

	if strings.TrimSpace(name) == "" {
		return "", errors.New("missing required parameter: name")
	}

	var translation cmsstore.TranslationInterface
	var err error

	// If ID is provided, try to find existing translation
	if strings.TrimSpace(id) != "" {
		translation, err = m.store.TranslationFindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if translation == nil {
			return "", errors.New("translation not found")
		}
	} else {
		// Create new translation
		translation = cmsstore.NewTranslation()

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
		translation.SetSiteID(siteID)
	}

	// Set/update fields
	translation.SetName(name)
	if handle != "" {
		translation.SetHandle(handle)
	}
	if status != "" {
		translation.SetStatus(status)
	}
	if v := strings.TrimSpace(argString(args, "site_id")); v != "" {
		translation.SetSiteID(v)
	}
	if v := strings.TrimSpace(argString(args, "memo")); v != "" {
		translation.SetMemo(v)
	}

	// Handle content (language-specific content)
	if contentMap, ok := args["content"].(map[string]any); ok {
		// Convert map[string]any to map[string]string
		contentStrMap := make(map[string]string)
		for k, v := range contentMap {
			if str, ok := v.(string); ok {
				contentStrMap[k] = str
			}
		}
		if err := translation.SetContent(contentStrMap); err != nil {
			return "", err
		}
	}

	// Save translation
	if strings.TrimSpace(id) != "" {
		// Update existing translation
		if err := m.store.TranslationUpdate(ctx, translation); err != nil {
			return "", err
		}
	} else {
		// Create new translation
		if err := m.store.TranslationCreate(ctx, translation); err != nil {
			return "", err
		}
	}

	content, err := translation.Content()
	if err != nil {
		return "", err
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      translation.ID(),
		"name":    translation.Name(),
		"handle":  translation.Handle(),
		"content": content,
		"status":  translation.Status(),
		"site_id": translation.SiteID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolTranslationDelete(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	translation, err := m.store.TranslationFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if translation == nil {
		return "", errors.New("translation not found")
	}

	if err := m.store.TranslationSoftDeleteByID(ctx, id); err != nil {
		if err := m.store.TranslationDelete(ctx, translation); err != nil {
			return "", err
		}
	}

	respBytes, err := json.Marshal(map[string]any{"id": id})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}
