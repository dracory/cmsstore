# Custom Entities Support

## Overview

The CMS now supports custom entities through integration with `dracory/entitystore`, providing a flexible Entity-Attribute-Value (EAV) storage system that allows you to extend the CMS with custom business logic without modifying the core schema.

## Features

- **Schemaless Storage**: Add new entity types without database migrations
- **Type Safety**: Strongly-typed attribute setters (string, int, float, bool)
- **Relationships**: Optional support for entity relationships (belongs_to, has_many, many_many)
- **Taxonomies**: Optional support for categorization and tagging
- **Soft Deletes**: Built-in trash bin functionality
- **SQL Queries**: Full SQL access for complex reporting

## Quick Start

### 1. Enable Custom Entities

```go
store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
    DB:                    db,
    BlockTableName:        "cms_block",
    PageTableName:         "cms_page",
    SiteTableName:         "cms_site",
    TemplateTableName:     "cms_template",
    AutomigrateEnabled:    true,
    
    // Enable custom entities
    CustomEntitiesEnabled: true,
    CustomEntityStoreOptions: cmsstore.CustomEntityStoreOptions{
        RelationshipsEnabled: true,  // Optional: enable relationships
        TaxonomiesEnabled:    true,  // Optional: enable taxonomies
    },
    
    // Register entity types
    CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
        {
            Type:      "product",
            TypeLabel: "Product",
            Group:     "Shop",
            Attributes: []cmsstore.CustomAttributeDefinition{
                {Name: "title", Type: "string", Label: "Title", Required: true},
                {Name: "price", Type: "float", Label: "Price", Required: true},
                {Name: "stock", Type: "int", Label: "Stock Quantity"},
                {Name: "description", Type: "string", Label: "Description"},
            },
            AllowRelationships: true,
            AllowTaxonomies:    true,
        },
    },
})
```

### 2. Create Custom Entities

```go
ctx := context.Background()

// Get the custom entity store
customStore := store.CustomEntityStore()

// Create a product
attrs := map[string]interface{}{
    "title":       "Laptop Computer",
    "price":       999.99,
    "stock":       25,
    "description": "High-performance laptop",
}

productID, err := customStore.Create(ctx, "product", attrs, nil, nil)
if err != nil {
    log.Fatal(err)
}
```

### 3. Retrieve and Update Entities

```go
// Find by ID
entity, err := customStore.FindByID(ctx, productID)
if err != nil {
    log.Fatal(err)
}

// Get attribute values
titleAttr, _ := customStore.Inner().AttributeFind(ctx, entity.ID(), "title")
title := titleAttr.AttributeValue()

// Update attributes
updateAttrs := map[string]interface{}{
    "price": 899.99,
    "stock": 20,
}

err = customStore.Update(ctx, entity, updateAttrs)
```

### 4. List and Query Entities

```go
// List all products
entities, err := customStore.List(ctx, entitystore.EntityQueryOptions{
    EntityType: "product",
    Limit:      10,
    Offset:     0,
})

// Count entities
count, err := customStore.Count(ctx, entitystore.EntityQueryOptions{
    EntityType: "product",
})
```

## Working with Relationships

```go
// Register related entity types
authorDef := cmsstore.CustomEntityDefinition{
    Type:               "author",
    TypeLabel:          "Author",
    AllowRelationships: true,
    Attributes: []cmsstore.CustomAttributeDefinition{
        {Name: "name", Type: "string", Label: "Name", Required: true},
    },
}

bookDef := cmsstore.CustomEntityDefinition{
    Type:                 "book",
    TypeLabel:            "Book",
    AllowRelationships:   true,
    AllowedRelationTypes: []string{"belongs_to"},
    Attributes: []cmsstore.CustomAttributeDefinition{
        {Name: "title", Type: "string", Label: "Title", Required: true},
    },
}

// Create author
authorAttrs := map[string]interface{}{"name": "John Doe"}
authorID, _ := customStore.Create(ctx, "author", authorAttrs, nil, nil)

// Create book with relationship to author
bookAttrs := map[string]interface{}{"title": "My Book"}
relationships := []cmsstore.RelationshipDefinition{
    {
        TargetID: authorID,
        Type:     "belongs_to",
        Metadata: map[string]interface{}{"role": "author"},
    },
}

bookID, _ := customStore.Create(ctx, "book", bookAttrs, relationships, nil)

// Get relationships
rels, _ := customStore.GetRelationships(ctx, bookID)
for _, rel := range rels {
    fmt.Printf("Related to: %s (type: %s)\n", rel.RelatedEntityID(), rel.RelationshipType())
}
```

## Working with Taxonomies

```go
// Create a taxonomy
taxonomy, err := customStore.Inner().TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
    Name:        "Categories",
    Slug:        "categories",
    EntityTypes: []string{"product"},
})

// Create taxonomy terms
electronics, _ := customStore.Inner().TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
    TaxonomyID: taxonomy.ID(),
    Name:       "Electronics",
    Slug:       "electronics",
})

computers, _ := customStore.Inner().TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
    TaxonomyID: taxonomy.ID(),
    ParentID:   electronics.ID(),
    Name:       "Computers",
    Slug:       "computers",
})

// Create product with taxonomy assignment
productAttrs := map[string]interface{}{
    "title": "Laptop",
    "price": 999.99,
}

productID, _ := customStore.Create(ctx, "product", productAttrs, nil, []string{computers.ID()})

// Get taxonomy assignments
assignments, _ := customStore.GetTaxonomyAssignments(ctx, productID)
```

## Entity Definition Configuration

### CustomEntityDefinition

| Field | Type | Description |
|-------|------|-------------|
| `Type` | string | Unique identifier for the entity type (e.g., "product") |
| `TypeLabel` | string | Human-readable label (e.g., "Product") |
| `Group` | string | Group for admin navigation (e.g., "Shop") |
| `Icon` | string | Icon for admin UI (optional) |
| `Attributes` | []CustomAttributeDefinition | Attribute definitions |
| `AllowRelationships` | bool | Enable relationships for this entity type |
| `AllowedRelationTypes` | []string | Restrict allowed relationship types |
| `AllowTaxonomies` | bool | Enable taxonomies for this entity type |
| `TaxonomyIDs` | []string | Restrict allowed taxonomies |

### CustomAttributeDefinition

| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | Attribute name (e.g., "title", "price") |
| `Type` | string | Data type: "string", "int", "float", "bool", "json" |
| `Label` | string | Human-readable label for admin UI |
| `Required` | bool | Whether the attribute is required |
| `DefaultValue` | interface{} | Default value for the attribute |
| `Validation` | string | Validation rules (optional) |
| `Help` | string | Help text for admin UI |

## Database Schema

When custom entities are enabled, the following tables are created:

```sql
-- Core entity table
CREATE TABLE cms_custom_entity (
    id VARCHAR(255) PRIMARY KEY,
    entity_type VARCHAR(255),
    entity_handle VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Attribute storage (EAV pattern)
CREATE TABLE cms_custom_attribute (
    id VARCHAR(255) PRIMARY KEY,
    entity_id VARCHAR(255),
    attribute_key VARCHAR(255),
    attribute_value TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    INDEX idx_entity_key (entity_id, attribute_key)
);

-- Soft delete support
CREATE TABLE cms_custom_entity_trash (...);
CREATE TABLE cms_custom_attribute_trash (...);
```

If relationships are enabled:

```sql
CREATE TABLE cms_custom_relationship (
    id VARCHAR(255) PRIMARY KEY,
    entity_id VARCHAR(255),
    related_entity_id VARCHAR(255),
    relationship_type VARCHAR(50),
    metadata TEXT,
    created_at TIMESTAMP,
    INDEX idx_entity (entity_id),
    INDEX idx_related (related_entity_id)
);
```

If taxonomies are enabled:

```sql
CREATE TABLE cms_custom_taxonomy (...);
CREATE TABLE cms_custom_taxonomy_term (...);
CREATE TABLE cms_custom_entity_taxonomy (...);
```

## Advanced Usage

### Direct Access to EntityStore

For advanced operations, you can access the underlying entitystore:

```go
innerStore := customStore.Inner()

// Use any entitystore method
entity := entitystore.NewEntity()
entity.SetEntityType("custom_type")
innerStore.EntityCreate(ctx, entity)
```

### Custom Validation

Implement custom validation by checking attributes before creation:

```go
func validateProduct(attrs map[string]interface{}) error {
    price, ok := attrs["price"].(float64)
    if !ok || price < 0 {
        return errors.New("price must be a positive number")
    }
    return nil
}

// Before creating
if err := validateProduct(attrs); err != nil {
    return err
}
productID, err := customStore.Create(ctx, "product", attrs, nil, nil)
```

### Soft Delete and Restore

```go
// Soft delete
err := customStore.Delete(ctx, entityID)

// Restore (using inner store)
err = customStore.Inner().EntityRestore(ctx, entityID)
```

## Best Practices

1. **Define Entity Types Early**: Register all entity types during store initialization
2. **Use Required Fields**: Mark essential attributes as required to ensure data integrity
3. **Leverage Relationships**: Use relationships instead of storing IDs as attributes
4. **Index Frequently Queried Attributes**: Consider adding database indexes for performance
5. **Validate Input**: Always validate attribute values before creation/update
6. **Use Transactions**: For complex operations involving multiple entities

## Limitations

- Attribute values are stored as strings (with type conversion helpers)
- Complex queries may require direct SQL access
- No built-in full-text search (use external search engine)
- Relationships are unidirectional (query from both sides if needed)

## Migration from Existing Systems

If you have existing custom tables, you can migrate to custom entities:

```go
// Read from old table
rows, _ := db.Query("SELECT id, name, price FROM old_products")

// Create custom entities
for rows.Next() {
    var id, name string
    var price float64
    rows.Scan(&id, &name, &price)
    
    attrs := map[string]interface{}{
        "title": name,
        "price": price,
    }
    customStore.Create(ctx, "product", attrs, nil, nil)
}
```

## See Also

- [entitystore Documentation](https://github.com/dracory/entitystore)
- [Proposal: Custom Entities Support](../proposals/2026-03-17-custom-entities-support.md)
- [CMS Store README](../README.md)
