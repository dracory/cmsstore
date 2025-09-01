# API Documentation

## Data Model

### DataObject Base
All entities in the CMS extend the `dataobject.DataObject` from `github.com/dracory/dataobject`. This provides:
- Common data management functionality
- Data storage and retrieval using `Get`/`Set` methods
- Data change tracking with `DataChanged()` and `MarkAsNotDirty()`
- Data hydration from maps with `Hydrate()`

### Common Entity Features
All entities (Page, Menu, Block, Site, Template, Translation) share:
- Soft delete support
- Metadata management
- Timestamp tracking (created_at, updated_at)
- Status management
- Type-safe getters/setters
- Common methods from DataObject:
  - `Data()` - Returns all data as map[string]string
  - `DataChanged()` - Returns changed data
  - `MarkAsNotDirty()` - Resets change tracking

## Query Interfaces

Each entity type has its own Query interface that provides:
- Fluent interface for building queries
- Common filtering capabilities:
  - ID/IDIn for single/multiple record selection
  - Limit/Offset for pagination
  - OrderBy/SortOrder for sorting
  - Status filtering
  - Soft delete handling

### Site Query Interface
The `SiteQueryInterface` provides methods for managing site-related operations:
- Site creation and management
- Domain handling
- Site configuration
- Filtering by domain name, handle, and status

### Page Query Interface
The `PageQueryInterface` handles page-related operations:
- Page CRUD operations
- URL management
- Content organization
- Template association
- Middleware management

### Menu Query Interface
The `MenuQueryInterface` manages menu-related functionality:
- Menu structure management
- Menu item organization
- Navigation handling
- Parent-child relationships

### Block Query Interface
The `BlockQueryInterface` handles content blocks:
- Block creation and management
- Content rendering
- Dynamic content insertion
- Block type management

### Template Query Interface
The `TemplateQueryInterface` manages templates:
- Template management
- Theme handling
- Layout organization
- Site-specific templating

### Translation Query Interface
The `TranslationQueryInterface` handles translations:
- Content localization
- Language management
- Translation storage
- Multi-language support

## Common Patterns

### Entity Creation
```go
// Creating new instance
page := NewPage()
page.SetTitle("My Page")
page.SetContent("Content")

// Loading existing data
existingPage := NewPageFromExistingData(data)
```

### Query Building
```go
query := PageQuery().
    SetLimit(10).
    SetOffset(0).
    SetOrderBy("created_at").
    SetSortOrder("DESC")
```

### Data Management
- CRUD operations
- Error handling
- Data validation
- Transaction management
- Soft delete support
- Metadata handling
- Versioning support (where applicable)

## Entity Interfaces

Each entity implements its specific interface (e.g., PageInterface, MenuInterface) which provides:
- Type-safe property access
- Entity-specific operations
- Status management methods
- Timestamp handling
- Metadata management
- Relationship management

See individual interface documentation for detailed usage examples and implementation details. 