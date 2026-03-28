# Custom Entity Admin Controllers

This package provides admin UI controllers for managing custom entities in the CMS.

## Features

- **List View**: Paginated table view with search and filtering
- **Create Form**: Modal form for creating new entities
- **Edit Form**: Modal form for editing existing entities
- **Delete Confirmation**: Safe deletion with confirmation dialog
- **HTMX Integration**: Seamless modal interactions without page reloads
- **Bootstrap 5 UI**: Consistent styling with the rest of the admin

## Controllers

### EntityListController
Displays a paginated list of entities with actions for create, edit, and delete.

### EntityCreateController
Handles entity creation with a modal form.

### EntityEditController
Handles entity editing with a pre-populated modal form.

### EntityDeleteController
Handles entity deletion with a confirmation dialog.

## Usage

### 1. Define Your Custom Entity

```go
productDef := cmsstore.CustomEntityDefinition{
    Type:      "product",
    TypeLabel: "Product",
    Group:     "Shop",
    Attributes: []cmsstore.CustomAttributeDefinition{
        {Name: "title", Type: "string", Label: "Title", Required: true},
        {Name: "price", Type: "float", Label: "Price", Required: true},
        {Name: "stock", Type: "int", Label: "Stock Quantity"},
        {Name: "description", Type: "string", Label: "Description"},
    },
}
```

### 2. Initialize the Store with Custom Entities

```go
store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
    DB:                    db,
    BlockTableName:        "cms_block",
    PageTableName:         "cms_page",
    SiteTableName:         "cms_site",
    TemplateTableName:     "cms_template",
    AutomigrateEnabled:    true,
    CustomEntitiesEnabled: true,
    CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{productDef},
})
```

### 3. Create Controllers

```go
import "github.com/dracory/cmsstore/admin/customentity"

// Create UI interface implementation
type AdminUI struct {
    store  cmsstore.StoreInterface
    layout hb.TagInterface
    logger any
}

func (ui *AdminUI) Store() cmsstore.StoreInterface { return ui.store }
func (ui *AdminUI) Layout() hb.TagInterface { return ui.layout }
func (ui *AdminUI) Logger() any { return ui.logger }

// Initialize controllers
ui := &AdminUI{store: store, layout: myLayout, logger: myLogger}

listController := customentity.NewEntityListController(ui, productDef)
createController := customentity.NewEntityCreateController(ui, productDef)
editController := customentity.NewEntityEditController(ui, productDef)
deleteController := customentity.NewEntityDeleteController(ui, productDef)
```

### 4. Register Routes

```go
// List/Index
http.HandleFunc("/admin/custom-entity/product", listController.Handler)

// Create
http.HandleFunc("/admin/custom-entity/product/create", createController.Handler)

// Edit
http.HandleFunc("/admin/custom-entity/product/edit", editController.Handler)

// Delete
http.HandleFunc("/admin/custom-entity/product/delete", deleteController.Handler)
```

### 5. Add Navigation Link

Add a link to your admin navigation:

```html
<a href="/admin/custom-entity/product" class="nav-link">
    <i class="bi bi-box"></i> Products
</a>
```

## Route Pattern

For each custom entity type, the following routes are created:

- `GET  /admin/custom-entity/{type}` - List view
- `GET  /admin/custom-entity/{type}/create` - Create form (modal)
- `POST /admin/custom-entity/{type}/create` - Process creation
- `GET  /admin/custom-entity/{type}/edit?id={id}` - Edit form (modal)
- `POST /admin/custom-entity/{type}/edit?id={id}` - Process update
- `GET  /admin/custom-entity/{type}/delete?id={id}` - Delete confirmation (modal)
- `POST /admin/custom-entity/{type}/delete?id={id}` - Process deletion

## Attribute Types

Supported attribute types and their form field mappings:

| Type | Form Field | Description |
|------|------------|-------------|
| `string` | Text input | Single-line text |
| `int` | Number input | Integer values |
| `float` | Number input (with decimals) | Decimal values |
| `bool` | Select dropdown | True/false values |
| `json` | Textarea | JSON data |

## Customization

### Custom Field Types

You can extend the `getFieldType()` method in the controllers to support additional field types.

### Custom Validation

Add validation logic in the `handleSubmit()` methods before calling `Create()` or `Update()`.

### Custom Actions

Add additional buttons in the table by modifying the `buildActions()` method in `EntityListController`.

## Example: Complete Setup

See `example_admin_setup.go` for a complete working example of setting up custom entity admin controllers.

## Requirements

- Bootstrap 5 (for styling)
- HTMX (for modal interactions)
- Bootstrap Icons (for icons)

These should be included in your admin layout template.

## Notes

- All modals use HTMX for seamless interactions
- Forms use POST method for submissions
- Entities are soft-deleted (moved to trash)
- Pagination shows up to 20 entities per page
- Entity IDs are displayed as code for easy copying
