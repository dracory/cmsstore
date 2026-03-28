package customentity

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dracory/cmsstore"
	"github.com/dracory/entitystore"
	"github.com/dracory/hb"
)

// EntityListController handles the list view for a custom entity type
type EntityListController struct {
	ui         UiInterface
	definition cmsstore.CustomEntityDefinition
}

// NewEntityListController creates a new list controller
func NewEntityListController(ui UiInterface, definition cmsstore.CustomEntityDefinition) *EntityListController {
	return &EntityListController{
		ui:         ui,
		definition: definition,
	}
}

// Handler renders the entity list page
func (c *EntityListController) Handler(w http.ResponseWriter, r *http.Request) string {
	ctx := context.Background()

	if !c.ui.Store().CustomEntitiesEnabled() {
		return c.errorPage("Custom entities are not enabled")
	}

	customStore := c.ui.Store().CustomEntityStore()
	if customStore == nil {
		return c.errorPage("Custom entity store is not available")
	}

	// Pagination
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	perPage := 20
	offset := (page - 1) * perPage

	// Fetch entities
	entities, err := customStore.List(ctx, entitystore.EntityQueryOptions{
		EntityType: c.definition.Type,
		Limit:      uint64(perPage),
		Offset:     uint64(offset),
		SortOrder:  "desc",
	})
	if err != nil {
		return c.errorPage("Error loading entities: " + err.Error())
	}

	// Count total
	totalCount, _ := customStore.Count(ctx, entitystore.EntityQueryOptions{
		EntityType: c.definition.Type,
	})

	// Build page content
	content := c.buildContent(ctx, customStore, entities, int(totalCount), perPage, page)

	// Return with layout
	return c.wrapInLayout(content)
}

func (c *EntityListController) buildContent(ctx context.Context, customStore *cmsstore.CustomEntityStore, entities []entitystore.EntityInterface, total, perPage, page int) hb.TagInterface {
	container := hb.Div().Class("container-fluid mt-4")

	// Header
	header := hb.Div().Class("d-flex justify-content-between align-items-center mb-4")
	header.Child(hb.Heading1().HTML(c.definition.TypeLabel + " Manager"))
	header.Child(c.createButton())
	container.Child(header)

	// Breadcrumbs
	container.Child(c.breadcrumbs())

	// Table or empty state
	if len(entities) == 0 {
		container.Child(hb.Div().Class("alert alert-info").
			HTML("No " + c.definition.TypeLabel + "s found. Click 'New " + c.definition.TypeLabel + "' to create one."))
	} else {
		container.Child(c.buildTable(ctx, customStore, entities))
		container.Child(c.buildPagination(total, perPage, page))
	}

	return container
}

func (c *EntityListController) breadcrumbs() hb.TagInterface {
	nav := hb.Nav().Attr("aria-label", "breadcrumb")
	ol := hb.OL().Class("breadcrumb")
	ol.Child(hb.LI().Class("breadcrumb-item").
		Child(hb.Hyperlink().HTML("Home").Href("/admin")))
	ol.Child(hb.LI().Class("breadcrumb-item").
		Child(hb.Hyperlink().HTML("Custom Entities").Href("/admin/custom-entities")))
	ol.Child(hb.LI().Class("breadcrumb-item active").
		Attr("aria-current", "page").
		HTML(c.definition.TypeLabel))
	nav.Child(ol)
	return nav
}

func (c *EntityListController) createButton() hb.TagInterface {
	url := "/admin/custom-entity/" + c.definition.Type + "/create"
	return hb.Button().
		Class("btn btn-primary").
		Child(hb.I().Class("bi bi-plus-circle me-2")).
		Child(hb.Span().HTML("New " + c.definition.TypeLabel)).
		HxGet(url).
		HxTarget("body").
		HxSwap("beforeend")
}

func (c *EntityListController) buildTable(ctx context.Context, customStore *cmsstore.CustomEntityStore, entities []entitystore.EntityInterface) hb.TagInterface {
	table := hb.Table().Class("table table-striped table-hover")

	// Header
	thead := hb.Thead()
	headerRow := hb.TR()
	headerRow.Child(hb.TH().HTML("ID"))
	for _, attr := range c.definition.Attributes {
		headerRow.Child(hb.TH().HTML(attr.Label))
	}
	headerRow.Child(hb.TH().HTML("Created"))
	headerRow.Child(hb.TH().HTML("Actions").Style("width: 200px;"))
	thead.Child(headerRow)
	table.Child(thead)

	// Body
	tbody := hb.Tbody()
	for _, entity := range entities {
		row := hb.TR()

		// ID
		row.Child(hb.TD().Child(hb.Code().HTML(entity.ID())).Style("font-size: 0.85em;"))

		// Attributes
		for _, attr := range c.definition.Attributes {
			value := c.getAttributeValue(ctx, customStore, entity.ID(), attr.Name)
			if len(value) > 50 {
				value = value[:47] + "..."
			}
			row.Child(hb.TD().HTML(value))
		}

		// Created
		row.Child(hb.TD().HTML(entity.CreatedAtCarbon().Format("Y-m-d H:i")))

		// Actions
		row.Child(c.buildActions(entity.ID()))

		tbody.Child(row)
	}
	table.Child(tbody)

	return table
}

func (c *EntityListController) getAttributeValue(ctx context.Context, customStore *cmsstore.CustomEntityStore, entityID, attrName string) string {
	attr, err := customStore.Inner().AttributeFind(ctx, entityID, attrName)
	if err != nil || attr == nil {
		return ""
	}
	return attr.AttributeValue()
}

func (c *EntityListController) buildActions(entityID string) hb.TagInterface {
	cell := hb.TD()

	// Edit button
	editUrl := "/admin/custom-entity/" + c.definition.Type + "/edit?id=" + entityID
	cell.Child(hb.Button().
		Class("btn btn-sm btn-primary me-1").
		Child(hb.I().Class("bi bi-pencil-square")).
		HxGet(editUrl).
		HxTarget("body").
		HxSwap("beforeend"))

	// Delete button
	deleteUrl := "/admin/custom-entity/" + c.definition.Type + "/delete?id=" + entityID
	cell.Child(hb.Button().
		Class("btn btn-sm btn-danger").
		Child(hb.I().Class("bi bi-trash")).
		HxGet(deleteUrl).
		HxTarget("body").
		HxSwap("beforeend"))

	return cell
}

func (c *EntityListController) buildPagination(total, perPage, currentPage int) hb.TagInterface {
	if total <= perPage {
		return hb.Div()
	}

	totalPages := (total + perPage - 1) / perPage
	baseUrl := "/admin/custom-entity/" + c.definition.Type

	nav := hb.Nav().Attr("aria-label", "Page navigation")
	ul := hb.UL().Class("pagination justify-content-center")

	// Previous
	if currentPage > 1 {
		ul.Child(hb.LI().Class("page-item").
			Child(hb.Hyperlink().Class("page-link").
				Href(baseUrl + "?page=" + strconv.Itoa(currentPage-1)).
				HTML("Previous")))
	} else {
		ul.Child(hb.LI().Class("page-item disabled").
			Child(hb.Span().Class("page-link").HTML("Previous")))
	}

	// Page numbers
	startPage := currentPage - 2
	if startPage < 1 {
		startPage = 1
	}
	endPage := startPage + 4
	if endPage > totalPages {
		endPage = totalPages
	}

	for i := startPage; i <= endPage; i++ {
		if i == currentPage {
			ul.Child(hb.LI().Class("page-item active").
				Child(hb.Span().Class("page-link").HTML(strconv.Itoa(i))))
		} else {
			ul.Child(hb.LI().Class("page-item").
				Child(hb.Hyperlink().Class("page-link").
					Href(baseUrl + "?page=" + strconv.Itoa(i)).
					HTML(strconv.Itoa(i))))
		}
	}

	// Next
	if currentPage < totalPages {
		ul.Child(hb.LI().Class("page-item").
			Child(hb.Hyperlink().Class("page-link").
				Href(baseUrl + "?page=" + strconv.Itoa(currentPage+1)).
				HTML("Next")))
	} else {
		ul.Child(hb.LI().Class("page-item disabled").
			Child(hb.Span().Class("page-link").HTML("Next")))
	}

	nav.Child(ul)
	return nav
}

func (c *EntityListController) errorPage(message string) string {
	return hb.Swal(hb.SwalOptions{
		Icon: "error",
		Text: message,
	}).ToHTML()
}

func (c *EntityListController) wrapInLayout(content hb.TagInterface) string {
	layout := c.ui.Layout()
	// Use ChildIf or similar pattern based on hb API
	wrapper := hb.Wrap().Child(layout).Child(content)
	return wrapper.ToHTML()
}
