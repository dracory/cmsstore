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

	// 6. Register routes (wrap handlers to write HTML response)
	http.HandleFunc("/admin/custom-entity/product", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(productListCtrl.Handler(w, r)))
	})
	http.HandleFunc("/admin/custom-entity/product/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(productCreateCtrl.Handler(w, r)))
	})
	http.HandleFunc("/admin/custom-entity/product/edit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(productEditCtrl.Handler(w, r)))
	})
	http.HandleFunc("/admin/custom-entity/product/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(productDeleteCtrl.Handler(w, r)))
	})

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
// Note: This is a simplified example. In production, use your actual layout implementation.
func createAdminLayout() hb.TagInterface {
	// Return a simple container div that will wrap the content
	// In a real implementation, this would include full HTML structure with head, body, etc.
	return hb.Div().Class("admin-layout")
}

// Note: Navigation creation removed from example as it's not used.
// In production, integrate custom entity links into your existing admin navigation.

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

		// Wrap handlers to write HTML response
		http.HandleFunc("/admin/custom-entity/"+entityType, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(listCtrl.Handler(w, r)))
		})
		http.HandleFunc("/admin/custom-entity/"+entityType+"/create", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(createCtrl.Handler(w, r)))
		})
		http.HandleFunc("/admin/custom-entity/"+entityType+"/edit", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(editCtrl.Handler(w, r)))
		})
		http.HandleFunc("/admin/custom-entity/"+entityType+"/delete", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(deleteCtrl.Handler(w, r)))
		})
	}

	http.ListenAndServe(":8080", nil)
}
