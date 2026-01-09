package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dracory/cmsstore"
)

func (m *MCP) toolBlockGet(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	block, err := m.store.BlockFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if block == nil {
		return "", errors.New("block not found")
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      block.ID(),
		"type":    block.Type(),
		"content": block.Content(),
		"status":  block.Status(),
		"site_id": block.SiteID(),
		"page_id": block.PageID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolBlockList(ctx context.Context, args map[string]any) (string, error) {
	q := cmsstore.BlockQuery()

	if v, ok := args["site_id"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetSiteID(v)
	}
	if v, ok := args["page_id"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetPageID(v)
	}
	if v, ok := args["status"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetStatus(v)
	}
	if v, ok := args["handle"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetHandle(v)
	}
	if v, ok := argBool(args, "include_soft_deleted"); ok {
		q.SetSoftDeleteIncluded(v)
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

	blocks, err := m.store.BlockList(ctx, q)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(blocks))
	for _, block := range blocks {
		if block == nil {
			continue
		}
		items = append(items, map[string]any{
			"id":          block.ID(),
			"type":        block.Type(),
			"name":        block.Name(),
			"handle":      block.Handle(),
			"status":      block.Status(),
			"site_id":     block.SiteID(),
			"page_id":     block.PageID(),
			"template_id": block.TemplateID(),
			"parent_id":   block.ParentID(),
			"sequence":    block.Sequence(),
			"created_at":  block.CreatedAt(),
			"updated_at":  block.UpdatedAt(),
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

func (m *MCP) toolBlockUpsert(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	blockType, _ := args["type"].(string)
	content, _ := args["content"].(string)
	status, _ := args["status"].(string)
	siteID, _ := args["site_id"].(string)

	if strings.TrimSpace(blockType) == "" {
		return "", errors.New("missing required parameter: type")
	}

	var block cmsstore.BlockInterface
	var err error

	// If ID is provided, try to find existing block
	if strings.TrimSpace(id) != "" {
		block, err = m.store.BlockFindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if block == nil {
			return "", errors.New("block not found")
		}
	} else {
		// Create new block
		block = cmsstore.NewBlock()

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
		block.SetSiteID(siteID)
	}

	// Set/update fields
	block.SetType(blockType)
	if content != "" {
		block.SetContent(content)
	}
	if status != "" {
		block.SetStatus(status)
	}
	if v := strings.TrimSpace(argString(args, "site_id")); v != "" {
		block.SetSiteID(v)
	}
	if v := strings.TrimSpace(argString(args, "page_id")); v != "" {
		block.SetPageID(v)
	}
	if v := strings.TrimSpace(argString(args, "template_id")); v != "" {
		block.SetTemplateID(v)
	}
	if v := strings.TrimSpace(argString(args, "parent_id")); v != "" {
		block.SetParentID(v)
	}
	if v := strings.TrimSpace(argString(args, "name")); v != "" {
		block.SetName(v)
	}
	if v := strings.TrimSpace(argString(args, "handle")); v != "" {
		block.SetHandle(v)
	}
	if v := strings.TrimSpace(argString(args, "editor")); v != "" {
		block.SetEditor(v)
	}
	if v := strings.TrimSpace(argString(args, "memo")); v != "" {
		block.SetMemo(v)
	}
	if v, ok := argInt(args, "sequence"); ok {
		block.SetSequenceInt(int(v))
	}

	// Save block
	if strings.TrimSpace(id) != "" {
		// Update existing block
		if err := m.store.BlockUpdate(ctx, block); err != nil {
			return "", err
		}
	} else {
		// Create new block
		if err := m.store.BlockCreate(ctx, block); err != nil {
			return "", err
		}
	}

	// Create versioning record if versioning is enabled
	if m.store.VersioningEnabled() {
		if err := m.createBlockVersioning(ctx, block); err != nil {
			// Log error but don't fail the operation
			// In a production environment, you might want to handle this differently
		}
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      block.ID(),
		"type":    block.Type(),
		"content": block.Content(),
		"status":  block.Status(),
		"site_id": block.SiteID(),
		"page_id": block.PageID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolBlockDelete(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	block, err := m.store.BlockFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if block == nil {
		return "", errors.New("block not found")
	}

	if err := m.store.BlockSoftDeleteByID(ctx, id); err != nil {
		if err := m.store.BlockDelete(ctx, block); err != nil {
			return "", err
		}
	}

	respBytes, err := json.Marshal(map[string]any{"id": id})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

// createBlockVersioning creates a versioning record for a block if versioning is enabled
func (m *MCP) createBlockVersioning(ctx context.Context, block cmsstore.BlockInterface) error {
	if !m.store.VersioningEnabled() {
		return nil
	}

	if block == nil {
		return errors.New("block is nil")
	}

	// Get last versioning to check if content has changed
	lastVersioningList, err := m.store.VersioningList(ctx, cmsstore.NewVersioningQuery().
		SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK).
		SetEntityID(block.ID()).
		SetOrderBy("created_at").
		SetSortOrder("DESC").
		SetLimit(1))

	if err != nil {
		return err
	}

	// Marshal block content for versioning
	blockData := map[string]any{
		"id":      block.ID(),
		"type":    block.Type(),
		"content": block.Content(),
		"name":    block.Name(),
		"handle":  block.Handle(),
		"status":  block.Status(),
		"site_id": block.SiteID(),
		"page_id": block.PageID(),
		"memo":    block.Memo(),
	}
	content, err := json.Marshal(blockData)
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
		SetEntityID(block.ID()).
		SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK).
		SetContent(string(content)))
}
