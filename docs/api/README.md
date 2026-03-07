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

### Transaction & Versioning Integrity

The CMS store ensures atomic integrity between entity updates and version tracking. When `VersioningEnabled` is set to `true`:

1.  **Automatic Transactions**: Every `Create` and `Update` operation (e.g., `PageCreate`, `BlockUpdate`) is automatically wrapped in a database transaction if one is not already present in the context.
2.  **Version Tracking**: The version record is created *within* the same transaction as the entity data. If versioning fails, the entire transaction (including the entity change) is rolled back.
3.  **Attribution (UserID)**: Version records automatically capture the identity of the editor. If the entity implements an `Editor()` method (as `Page`, `Block`, and `Template` do), the returned ID is stored within the version snapshot as `_userID`.
4.  **Manual Transactions**: You can participate in this transactional flow by providing your own `*sql.Tx` via `cmsstore.WithTransaction(tx)` or by wrapping a transaction in the context using `database.Context(ctx, tx)`. If a transaction is detected in the context, the store will use it instead of starting a new one, allowing you to group multiple CMS operations into a single atomic unit.

## Entity Interfaces

Each entity implements its specific interface (e.g., PageInterface, MenuInterface) which provides:
- Type-safe property access
- Entity-specific operations
- Status management methods
- Timestamp handling
- Metadata management
- Relationship management

See individual interface documentation for detailed usage examples and implementation details. 