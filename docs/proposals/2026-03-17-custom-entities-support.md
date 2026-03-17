# [Draft] Implement Custom Entities Support

## Summary
- **Problem**: The CMS currently only supports predefined entity types (pages, blocks, menus, templates, translations) with no way to extend functionality for custom business logic or domain-specific data.
- **Solution**: Implement a flexible custom entities system that allows developers to define and manage their own entity types with configurable fields, relationships, and admin interfaces.

## Background

During codebase analysis, I discovered that custom entities functionality was partially implemented but commented out in several places:

1. **Admin Header** (`admin/shared/admin_header.go`): Contains commented code for entity manager navigation
2. **Development Main** (`admin/development/main.go`): Has extensive commented examples of custom entity definitions including:
   - User entities with first_name, last_name, email, image_url
   - Shop products with title, description, price, image_url
   - Shop orders with user_id, total
   - Shop order line items with order_id, product_id, quantity, subtotal
   - Make-a-wish entities with wish, referral

3. **Store Interface**: The `StoreInterface` in `interfaces.go` shows no custom entity methods, indicating this functionality was never fully integrated.

The existing architecture provides a solid foundation with:
- Versioning support for all entities
- Admin interface patterns
- Query interfaces
- Middleware support

## Detailed Design

### 1. Custom Entity Configuration Structure

```go
type CustomEntityStructure struct {
    Group         string                    // Group for admin navigation (e.g., "Shop", "Users")
    Type          string                    // Unique entity type identifier
    TypeLabel     string                    // Human-readable label
    AttributeList []CustomAttributeStructure // Field definitions
}

type CustomAttributeStructure struct {
    Name             string                 // Field name
    Type             string                 // Field type (string, int, float, text, textarea, select, etc.)
    FormControlLabel string                 // Label for admin form
    FormControlType  string                 // Form control type (input, textarea, select, etc.)
    FormControlHelp  string                 // Help text for admin form
    BelongsToType    string                 // Optional relationship to another entity type
    Options          []string               // For select fields
}
```

### 2. Store Interface Extensions

```go
type StoreInterface interface {
    // ... existing methods ...
    
    // Custom Entity Methods
    CustomEntityCreate(ctx context.Context, entityType string, entity map[string]interface{}) error
    CustomEntityCount(ctx context.Context, entityType string, query CustomEntityQueryInterface) (int64, error)
    CustomEntityDelete(ctx context.Context, entityType string, entityID string) error
    CustomEntityFindByID(ctx context.Context, entityType string, entityID string) (map[string]interface{}, error)
    CustomEntityList(ctx context.Context, entityType string, query CustomEntityQueryInterface) ([]map[string]interface{}, error)
    CustomEntitySoftDelete(ctx context.Context, entityType string, entityID string) error
    CustomEntityUpdate(ctx context.Context, entityType string, entityID string, updates map[string]interface{}) error
    
    // Custom Entity Query Interface
    CustomEntityQuery(entityType string) CustomEntityQueryInterface
}
```

### 3. Admin Interface Integration

The admin interface will automatically generate CRUD operations for each custom entity:

- **Navigation**: Entities grouped by their Group field in the admin sidebar
- **List View**: Table with configurable columns, search, and pagination
- **Form View**: Dynamic forms based on attribute definitions
- **Versioning**: Full version history support
- **Relationships**: Support for belongs_to relationships with dropdowns

### 4. Database Schema

Custom entities will use a flexible schema:

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

-- Entity data (EAV pattern for flexibility)
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
    relationship_type VARCHAR(50), -- belongs_to, has_many, etc.
    created_at TIMESTAMP
);
```

### 5. Configuration Integration

Extend the existing configuration system to support custom entities:

```go
type Config struct {
    // ... existing config ...
    CustomEntityList []CustomEntityStructure
}

// Usage example
config := Config{
    CustomEntityList: []CustomEntityStructure{
        {
            Group:     "Users",
            Type:      "user",
            TypeLabel: "User",
            AttributeList: []CustomAttributeStructure{
                {
                    Name:             "first_name",
                    Type:             "string",
                    FormControlLabel: "First Name",
                    FormControlType:  "input",
                    FormControlHelp:  "The first name of the user",
                },
                // ... more attributes
            },
        },
    },
}
```

## Implementation Plan

### Phase 1: Core Infrastructure (2 weeks)
1. **Database Schema**: Implement flexible storage for custom entities
2. **Store Interface**: Add custom entity methods to StoreInterface
3. **Query Interface**: Create CustomEntityQueryInterface with filtering, sorting, pagination
4. **Versioning Integration**: Ensure custom entities work with existing versioning system

### Phase 2: Admin Interface (2 weeks)
1. **Dynamic Forms**: Generate admin forms based on entity definitions
2. **List Views**: Create configurable table views with search and pagination
3. **Navigation**: Add entity groups to admin sidebar
4. **CRUD Operations**: Implement create, read, update, delete functionality

### Phase 3: Advanced Features (1 week)
1. **Relationships**: Support for belongs_to relationships with dropdowns
2. **Validation**: Field-level validation based on attribute types
3. **Permissions**: Basic permission system for entity access
4. **API Endpoints**: REST API for custom entities

### Phase 4: Testing and Documentation (1 week)
1. **Unit Tests**: Comprehensive test coverage for all new functionality
2. **Integration Tests**: Test admin interface and API endpoints
3. **Documentation**: Update README and create custom entity examples
4. **Migration**: Provide migration path for existing systems

## Alternatives Considered

### 1. Plugin System
- **Pros**: More flexible, allows custom business logic
- **Cons**: Complex implementation, harder to maintain consistency
- **Rejected**: Too complex for initial implementation

### 2. Separate Entity Service
- **Pros**: Isolated concerns, easier to scale
- **Cons**: Adds complexity, requires API integration
- **Rejected**: Overkill for current needs

### 3. Database-per-Entity
- **Pros**: Better performance, cleaner schema
- **Cons**: Harder to implement dynamic forms, more complex migrations
- **Rejected**: EAV pattern provides better flexibility for initial implementation

## Risks and Mitigations

### 1. Performance with EAV Pattern
- **Risk**: EAV can be slow for complex queries
- **Mitigation**: Implement caching, consider hybrid approach for high-volume entities

### 2. Schema Complexity
- **Risk**: Flexible schema makes data integrity harder to enforce
- **Mitigation**: Strong validation in application layer, clear documentation

### 3. Admin Interface Complexity
- **Risk**: Dynamic forms may not cover all use cases
- **Mitigation**: Provide extension points, allow custom form templates

### 4. Migration Complexity
- **Risk**: Existing systems may have custom data structures
- **Mitigation**: Provide migration tools, maintain backward compatibility

## Benefits

1. **Extensibility**: Developers can add custom entities without modifying core code
2. **Rapid Development**: Pre-built admin interface reduces development time
3. **Consistency**: All entities follow same patterns for versioning, permissions, etc.
4. **Flexibility**: Supports various business domains (e-commerce, user management, etc.)
5. **Maintainability**: Centralized entity management reduces code duplication

## Success Criteria

1. **Functional**: Custom entities can be defined, created, updated, and deleted through admin interface
2. **Performance**: Entity operations complete within acceptable time limits (< 1 second for CRUD operations)
3. **Usability**: Developers can add new entity types in under 30 minutes
4. **Integration**: Custom entities work seamlessly with existing versioning and middleware systems
5. **Documentation**: Complete examples and API documentation available

## Future Enhancements

1. **Custom Validation Rules**: Allow developers to define custom validation logic
2. **Custom Actions**: Support for custom admin actions and bulk operations
3. **Advanced Relationships**: Support for has_many, many_to_many relationships
4. **Custom Views**: Allow custom list and form templates
5. **Export/Import**: CSV and JSON import/export functionality
6. **Webhooks**: Trigger events on entity changes
7. **GraphQL API**: Alternative to REST for complex queries