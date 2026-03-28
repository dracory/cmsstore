# [Draft] Implement Custom Entities Support

## Status
**[Draft]** - Partially planned, never implemented

## Summary
- **Problem**: The CMS only supports predefined entity types with no way to extend for custom business logic
- **Solution**: Implement a flexible custom entities system for developer-defined entity types

## Current State (As-Is)

The custom entities feature was **planned but never implemented**. Evidence exists in commented code:

**Commented Code Locations:**
1. **`admin/development/main.go`** (lines 177-328): Complete commented `entityList()` function with examples:
   - User entities (first_name, last_name, email, image_url)
   - Shop products (title, description, price, image_url)
   - Shop orders with belongs_to relationships
   - Shop order line items
   - Make-a-wish entities

2. **`admin/shared/admin_header.go`** (lines 37-44): Commented navigation links for:
   - Widgets manager
   - Settings manager
   - (Implied entity manager was planned)

3. **`StoreInterface` (`interfaces.go`)**: No custom entity methods exist

**What Exists:**
- ❌ No `CustomEntityStructure` type
- ❌ No `CustomAttributeStructure` type
- ❌ No custom entity database schema
- ❌ No store methods for custom entities
- ❌ No admin interface for custom entities
- ❌ No custom entity query interface

**The Commented Design (from main.go):**
```go
// CustomEntityStructure - planned but never implemented
type CustomEntityStructure struct {
    Group         string                    // "Users", "Shop"
    Type          string                    // "user", "shop_product"
    TypeLabel     string                    // "User", "Product"
    AttributeList []CustomAttributeStructure
}

type CustomAttributeStructure struct {
    Name             string  // Field name
    Type             string  // "string", "int"
    FormControlLabel string  // Display label
    FormControlType  string  // "input", "textarea"
    FormControlHelp  string  // Help text
    BelongsToType    string  // Optional relationship
}
```

## Proposed Design (To-Be)

### 1. Entity Configuration

```go
type CustomEntityStructure struct {
    Group         string
    Type          string
    TypeLabel     string
    AttributeList []CustomAttributeStructure
}

type CustomAttributeStructure struct {
    Name             string
    Type             string  // string, int, float, text, textarea, select
    FormControlLabel string
    FormControlType  string
    FormControlHelp  string
    BelongsToType    string  // Relationship
    Options          []string
}
```

### 2. Store Interface Extensions

```go
type StoreInterface interface {
    // ... existing methods ...
    
    // Custom Entity Methods
    CustomEntityCreate(ctx context.Context, entityType string, entity map[string]interface{}) error
    CustomEntityFindByID(ctx context.Context, entityType string, entityID string) (map[string]interface{}, error)
    CustomEntityList(ctx context.Context, entityType string, query CustomEntityQueryInterface) ([]map[string]interface{}, error)
    CustomEntityUpdate(ctx context.Context, entityType string, entityID string, updates map[string]interface{}) error
    CustomEntityDelete(ctx context.Context, entityType string, entityID string) error
    CustomEntityQuery(entityType string) CustomEntityQueryInterface
}
```

### 3. Database Schema (EAV Pattern)

```sql
-- Entity definitions
CREATE TABLE custom_entity_definitions (
    id VARCHAR(255) PRIMARY KEY,
    group_name VARCHAR(255),
    type VARCHAR(255) UNIQUE,
    type_label VARCHAR(255),
    attributes JSON,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Entity data (EAV pattern)
CREATE TABLE custom_entity_data (
    id VARCHAR(255) PRIMARY KEY,
    entity_type VARCHAR(255),
    attribute_name VARCHAR(255),
    attribute_value TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Entity relationships
CREATE TABLE custom_entity_relationships (
    id VARCHAR(255) PRIMARY KEY,
    entity_type VARCHAR(255),
    entity_id VARCHAR(255),
    related_type VARCHAR(255),
    related_id VARCHAR(255),
    relationship_type VARCHAR(50),
    created_at TIMESTAMP
);
```

### 4. Admin Interface Integration

- Dynamic CRUD forms based on attribute definitions
- Navigation grouped by Group field
- Versioning support
- Relationship dropdowns

### 5. Configuration Integration

```go
type Config struct {
    // ... existing config ...
    CustomEntityList []CustomEntityStructure
}

// Usage
config := Config{
    CustomEntityList: []CustomEntityStructure{
        {
            Group:     "Shop",
            Type:      "product",
            TypeLabel: "Product",
            AttributeList: []CustomAttributeStructure{
                {
                    Name:             "title",
                    Type:             "string",
                    FormControlLabel: "Title",
                    FormControlType:  "input",
                },
            },
        },
    },
}
```

## Implementation Status

| Feature | Status | Notes |
|---------|--------|-------|
| Entity structure types | ❌ Not implemented | Commented in main.go |
| Attribute structure types | ❌ Not implemented | Commented in main.go |
| Database schema | ❌ Not implemented | No custom_entity tables |
| Store interface methods | ❌ Not implemented | No CustomEntity* methods |
| Query interface | ❌ Not implemented | No CustomEntityQueryInterface |
| Admin forms | ❌ Not implemented | No dynamic form generation |
| Admin navigation | ❌ Not implemented | No entity groups in sidebar |
| Versioning integration | ❌ Not implemented | Not integrated |
| REST API | ❌ Not implemented | No endpoints |

## Files to Create/Modify (If Implementing)

1. New: `custom_entity.go` - Core types and interfaces
2. New: `custom_entity_query.go` - Query interface
3. New: `store_custom_entities.go` - Store implementation
4. New: `admin/customentities/` - Admin controllers
5. Modify: `interfaces.go` - Add to StoreInterface
6. Modify: `admin/shared/admin_header.go` - Uncomment/add navigation
7. Modify: `options.go` - Add CustomEntityList to Config

## Historical Context

The custom entities feature appears to have been:
1. **Planned early** in the CMS architecture
2. **Partially implemented** in development examples
3. **Never fully integrated** into the production system
4. **Commented out** when the CMS was refactored

The commented code in `main.go` shows a complete working design that was likely functional in an earlier version but removed during a major refactor.

## Risks and Mitigations

1. **EAV Performance**
   - Risk: EAV pattern can be slow
   - Mitigation: Implement caching, consider hybrid approach

2. **Schema Complexity**
   - Risk: Flexible schema makes validation harder
   - Mitigation: Strong validation in application layer

3. **Admin Complexity**
   - Risk: Dynamic forms may not cover all use cases
   - Mitigation: Extension points, custom form templates

4. **Migration**
   - Risk: Breaking existing data structures
   - Mitigation: Migration tools, backward compatibility

## Future Considerations

The commented code suggests these features were also planned:
- Custom validation rules
- Bulk operations
- Export/Import (CSV, JSON)
- Webhooks on entity changes
- GraphQL API support