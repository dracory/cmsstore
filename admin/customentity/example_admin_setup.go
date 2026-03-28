package customentity

import (
	"database/sql"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
)

// ExampleAdminSetup demonstrates how to set up custom entity admin controllers
func ExampleAdminSetup() {
	// 1. Initialize database
	db, _ := sql.Open("sqlite", "cms.db")
	defer db.Close()

	// 2. Define custom entity types
	productDef := cmsstore.CustomEntityDefinition{
		Type:      "product",
		TypeLabel: "Product",
		Group:     "Shop",
		Icon:      "bi-box",
		Attributes: []cmsstore.CustomAttributeDefinition{
			{
				Name:     "title",
				Type:     "string",
				Label:    "Product Title",
				Required: true,
				Help:     "Enter the product name",
			},
			{
				Name:     "price",
				Type:     "float",
				Label:    "Price",
				Required: true,
				Help:     "Enter the price in USD",
			},
			{
				Name:     "stock",
				Type:     "int",
				Label:    "Stock Quantity",
				Required: false,
				Help:     "Number of items in stock",
			},
			{
				Name:     "description",
				Type:     "string",
				Label:    "Description",
				Required: false,
			},
		},
		AllowRelationships: true,
		AllowTaxonomies:    true,
	}

	// 3. Create CMS store with custom entities
	store, _ := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                 db,
		BlockTableName:     "cms_block",
		PageTableName:      "cms_page",
		SiteTableName:      "cms_site",
		TemplateTableName:  "cms_template",
		AutomigrateEnabled: true,

		CustomEntitiesEnabled: true,
		CustomEntityStoreOptions: cmsstore.CustomEntityStoreOptions{
			RelationshipsEnabled: true,
			TaxonomiesEnabled:    true,
		},
		CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
			productDef,
		},
	})

	// 4. Create UI interface implementation
	ui := &ExampleUI{
		store:  store,
		layout: createAdminLayout(),
	}

	// 5. Initialize controllers for each entity type
	productListCtrl := NewEntityListController(ui, productDef)
	productCreateCtrl := NewEntityCreateController(ui, productDef)
	productEditCtrl := NewEntityEditController(ui, productDef)
	productDeleteCtrl := NewEntityDeleteController(ui, productDef)

	// 6. Register routes
	http.HandleFunc("/admin/custom-entity/product", productListCtrl.Handler)
	http.HandleFunc("/admin/custom-entity/product/create", productCreateCtrl.Handler)
	http.HandleFunc("/admin/custom-entity/product/edit", productEditCtrl.Handler)
	http.HandleFunc("/admin/custom-entity/product/delete", productDeleteCtrl.Handler)

	// 7. Start server
	http.ListenAndServe(":8080", nil)
}

// ExampleUI implements the UiInterface for custom entity controllers
type ExampleUI struct {
	store  cmsstore.StoreInterface
	layout hb.TagInterface
	logger interface{}
}

func (ui *ExampleUI) Store() cmsstore.StoreInterface {
	return ui.store
}

func (ui *ExampleUI) Layout() hb.TagInterface {
	return ui.layout
}

func (ui *ExampleUI) Logger() any {
	return ui.logger
}

// createAdminLayout creates a basic admin layout
func createAdminLayout() hb.TagInterface {
	return hb.HTML().
		Children([]hb.TagInterface{
			hb.Head().Children([]hb.TagInterface{
				hb.Meta().Charset("utf-8"),
				hb.Meta().Name("viewport").Content("width=device-width, initial-scale=1"),
				hb.Title().HTML("CMS Admin"),
				// Bootstrap CSS
				hb.Link().Rel("stylesheet").
					Href("https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"),
				// Bootstrap Icons
				hb.Link().Rel("stylesheet").
					Href("https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.0/font/bootstrap-icons.css"),
				// HTMX
				hb.Script("").Src("https://unpkg.com/htmx.org@1.9.10"),
			}),
			hb.Body().Children([]hb.TagInterface{
				// Navigation
				createNavigation(),
				// Main content area (will be filled by controllers)
				hb.Div().ID("main-content"),
				// Bootstrap JS
				hb.Script("").
					Src("https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"),
			}),
		})
}

// createNavigation creates the admin navigation menu
func createNavigation() hb.TagInterface {
	nav := hb.Nav().Class("navbar navbar-expand-lg navbar-dark bg-dark")
	container := hb.Div().Class("container-fluid")

	// Brand
	container.Child(hb.Hyperlink().Class("navbar-brand").Href("/admin").HTML("CMS Admin"))

	// Toggle button for mobile
	container.Child(hb.Button().
		Class("navbar-toggler").
		Type("button").
		Data("bs-toggle", "collapse").
		Data("bs-target", "#navbarNav").
		Child(hb.Span().Class("navbar-toggler-icon")))

	// Nav items
	collapse := hb.Div().Class("collapse navbar-collapse").ID("navbarNav")
	ul := hb.UL().Class("navbar-nav")

	// Dashboard
	ul.Child(hb.LI().Class("nav-item").
		Child(hb.Hyperlink().Class("nav-link").Href("/admin").
			Child(hb.I().Class("bi bi-house me-2")).
			Child(hb.Span().HTML("Dashboard"))))

	// Pages
	ul.Child(hb.LI().Class("nav-item").
		Child(hb.Hyperlink().Class("nav-link").Href("/admin/pages").
			Child(hb.I().Class("bi bi-file-text me-2")).
			Child(hb.Span().HTML("Pages"))))

	// Custom Entities Dropdown
	dropdown := hb.LI().Class("nav-item dropdown")
	dropdown.Child(hb.Hyperlink().
		Class("nav-link dropdown-toggle").
		Href("#").
		ID("customEntitiesDropdown").
		Attr("role", "button").
		Data("bs-toggle", "dropdown").
		Child(hb.I().Class("bi bi-grid me-2")).
		Child(hb.Span().HTML("Custom Entities")))

	dropdownMenu := hb.UL().Class("dropdown-menu").Attr("aria-labelledby", "customEntitiesDropdown")
	dropdownMenu.Child(hb.LI().Child(hb.Hyperlink().Class("dropdown-item").
		Href("/admin/custom-entity/product").
		Child(hb.I().Class("bi bi-box me-2")).
		Child(hb.Span().HTML("Products"))))
	dropdown.Child(dropdownMenu)
	ul.Child(dropdown)

	collapse.Child(ul)
	container.Child(collapse)
	nav.Child(container)

	return nav
}

// ExampleMultipleEntityTypes shows how to set up multiple custom entity types
func ExampleMultipleEntityTypes() {
	db, _ := sql.Open("sqlite", "cms.db")
	defer db.Close()

	// Define multiple entity types
	definitions := []cmsstore.CustomEntityDefinition{
		{
			Type:      "product",
			TypeLabel: "Product",
			Group:     "Shop",
			Attributes: []cmsstore.CustomAttributeDefinition{
				{Name: "title", Type: "string", Label: "Title", Required: true},
				{Name: "price", Type: "float", Label: "Price", Required: true},
			},
		},
		{
			Type:      "customer",
			TypeLabel: "Customer",
			Group:     "Shop",
			Attributes: []cmsstore.CustomAttributeDefinition{
				{Name: "name", Type: "string", Label: "Name", Required: true},
				{Name: "email", Type: "string", Label: "Email", Required: true},
			},
		},
		{
			Type:      "order",
			TypeLabel: "Order",
			Group:     "Shop",
			Attributes: []cmsstore.CustomAttributeDefinition{
				{Name: "order_number", Type: "string", Label: "Order Number", Required: true},
				{Name: "total", Type: "float", Label: "Total", Required: true},
			},
		},
	}

	// Create store
	store, _ := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                      db,
		BlockTableName:          "cms_block",
		PageTableName:           "cms_page",
		SiteTableName:           "cms_site",
		TemplateTableName:       "cms_template",
		AutomigrateEnabled:      true,
		CustomEntitiesEnabled:   true,
		CustomEntityDefinitions: definitions,
	})

	ui := &ExampleUI{store: store, layout: createAdminLayout()}

	// Register routes for each entity type
	for _, def := range definitions {
		entityType := def.Type

		listCtrl := NewEntityListController(ui, def)
		createCtrl := NewEntityCreateController(ui, def)
		editCtrl := NewEntityEditController(ui, def)
		deleteCtrl := NewEntityDeleteController(ui, def)

		http.HandleFunc("/admin/custom-entity/"+entityType, listCtrl.Handler)
		http.HandleFunc("/admin/custom-entity/"+entityType+"/create", createCtrl.Handler)
		http.HandleFunc("/admin/custom-entity/"+entityType+"/edit", editCtrl.Handler)
		http.HandleFunc("/admin/custom-entity/"+entityType+"/delete", deleteCtrl.Handler)
	}

	http.ListenAndServe(":8080", nil)
}
