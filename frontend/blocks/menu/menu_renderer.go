package menu

import (
	"context"
	"strings"

	"github.com/dracory/cmsstore"
)

// MenuRenderer provides comprehensive menu rendering functionality
type MenuRenderer struct {
	store FrontendStore
}

// NewMenuRenderer creates a new menu renderer
func NewMenuRenderer(store FrontendStore) *MenuRenderer {
	return &MenuRenderer{
		store: store,
	}
}

// RenderMenuHTML renders menu items as HTML based on the specified style
func (r *MenuRenderer) RenderMenuHTML(ctx context.Context, menuItems []cmsstore.MenuItemInterface, style, cssClass string, startLevel, maxDepth int) (string, error) {
	if len(menuItems) == 0 {
		return "<!-- No menu items to render -->", nil
	}

	tree := r.buildMenuTree(menuItems, "")

	if startLevel > 0 {
		tree = r.filterMenuTreeByLevel(tree, startLevel)
	}

	if maxDepth > 0 {
		tree = r.limitMenuTreeDepth(tree, maxDepth)
	}

	switch style {
	case cmsstore.BLOCK_MENU_STYLE_HORIZONTAL:
		return r.renderMenuHorizontal(ctx, tree, cssClass), nil
	case cmsstore.BLOCK_MENU_STYLE_VERTICAL:
		return r.renderMenuVertical(ctx, tree, cssClass), nil
	case cmsstore.BLOCK_MENU_STYLE_DROPDOWN:
		return r.renderMenuDropdown(ctx, tree, cssClass), nil
	case cmsstore.BLOCK_MENU_STYLE_BREADCRUMB:
		return r.renderMenuBreadcrumb(ctx, tree, cssClass), nil
	default:
		return r.renderMenuVertical(ctx, tree, cssClass), nil
	}
}

// menuTreeNode represents a node in the menu tree
type menuTreeNode struct {
	Item     cmsstore.MenuItemInterface
	Children []*menuTreeNode
}

// buildMenuTree builds a hierarchical tree structure from flat menu items
func (r *MenuRenderer) buildMenuTree(items []cmsstore.MenuItemInterface, parentID string) []*menuTreeNode {
	var nodes []*menuTreeNode

	for _, item := range items {
		if item.ParentID() == parentID {
			node := &menuTreeNode{
				Item:     item,
				Children: r.buildMenuTree(items, item.ID()),
			}
			nodes = append(nodes, node)
		}
	}

	return nodes
}

// filterMenuTreeByLevel filters the tree to start at a specific level
func (r *MenuRenderer) filterMenuTreeByLevel(tree []*menuTreeNode, level int) []*menuTreeNode {
	if level <= 0 {
		return tree
	}

	var result []*menuTreeNode
	for _, node := range tree {
		filtered := r.filterMenuTreeByLevel(node.Children, level-1)
		result = append(result, filtered...)
	}
	return result
}

// limitMenuTreeDepth limits the depth of the menu tree
func (r *MenuRenderer) limitMenuTreeDepth(tree []*menuTreeNode, maxDepth int) []*menuTreeNode {
	if maxDepth <= 0 {
		return nil
	}

	var result []*menuTreeNode
	for _, node := range tree {
		newNode := &menuTreeNode{
			Item:     node.Item,
			Children: r.limitMenuTreeDepth(node.Children, maxDepth-1),
		}
		result = append(result, newNode)
	}
	return result
}

// renderMenuHorizontal renders a horizontal menu
func (r *MenuRenderer) renderMenuHorizontal(ctx context.Context, tree []*menuTreeNode, cssClass string) string {
	if len(tree) == 0 {
		return ""
	}

	html := `<ul`
	if cssClass != "" {
		html += ` class="` + cssClass + `"`
	}
	html += `>`

	for _, node := range tree {
		html += r.renderMenuItemHTML(ctx, node, false)
	}

	html += `</ul>`
	return html
}

// renderMenuVertical renders a vertical menu
func (r *MenuRenderer) renderMenuVertical(ctx context.Context, tree []*menuTreeNode, cssClass string) string {
	if len(tree) == 0 {
		return ""
	}

	html := `<ul`
	if cssClass != "" {
		html += ` class="` + cssClass + `"`
	}
	html += `>`

	for _, node := range tree {
		html += r.renderMenuItemHTML(ctx, node, true)
	}

	html += `</ul>`
	return html
}

// renderMenuDropdown renders a dropdown menu
func (r *MenuRenderer) renderMenuDropdown(ctx context.Context, tree []*menuTreeNode, cssClass string) string {
	if len(tree) == 0 {
		return ""
	}

	html := `<ul`
	if cssClass != "" {
		html += ` class="` + cssClass + `"`
	}
	html += `>`

	for _, node := range tree {
		html += r.renderMenuItemHTML(ctx, node, true)
	}

	html += `</ul>`
	return html
}

// renderMenuBreadcrumb renders a breadcrumb menu
func (r *MenuRenderer) renderMenuBreadcrumb(ctx context.Context, tree []*menuTreeNode, cssClass string) string {
	if len(tree) == 0 {
		return ""
	}

	html := `<ol`
	if cssClass != "" {
		html += ` class="` + cssClass + `"`
	}
	html += `>`

	for _, node := range tree {
		html += r.renderMenuItemHTML(ctx, node, false)
	}

	html += `</ol>`
	return html
}

// renderMenuItemHTML renders a single menu item with its children
func (r *MenuRenderer) renderMenuItemHTML(ctx context.Context, node *menuTreeNode, renderChildren bool) string {
	url := r.resolveMenuItemURL(ctx, node.Item)
	target := node.Item.Target()

	html := `<li>`

	if url != "" {
		html += `<a href="` + url + `"`
		if target != "" {
			html += ` target="` + target + `"`
		}
		html += `>` + node.Item.Name() + `</a>`
	} else {
		html += node.Item.Name()
	}

	if renderChildren && len(node.Children) > 0 {
		html += `<ul>`
		for _, child := range node.Children {
			html += r.renderMenuItemHTML(ctx, child, true)
		}
		html += `</ul>`
	}

	html += `</li>`
	return html
}

// resolveMenuItemURL resolves the URL for a menu item
func (r *MenuRenderer) resolveMenuItemURL(ctx context.Context, item cmsstore.MenuItemInterface) string {
	if item.URL() != "" {
		return item.URL()
	}

	if item.PageID() != "" {
		page, err := r.store.PageFindByID(ctx, item.PageID())
		if err != nil || page == nil {
			r.store.Logger().Debug("resolveMenuItemURL: Page not found", "pageID", item.PageID(), "error", err)
			return ""
		}
		return "/" + strings.TrimPrefix(page.Alias(), "/")
	}

	return ""
}
