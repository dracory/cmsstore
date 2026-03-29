# Custom Entities Examples

This directory contains comprehensive examples demonstrating how to use the custom entities feature in `cmsstore`.

## Overview

Custom entities allow you to extend the CMS with your own entity types without modifying the database schema. Built on top of `dracory/entitystore`, this feature provides:

- **Schema-less storage** - Add new entity types without migrations
- **EAV pattern** - Flexible attribute storage with type safety
- **Relationships** - Link entities together (belongs_to, has_many)
- **Taxonomies** - Categorize and tag entities
- **Soft deletes** - Built-in trash/restore functionality

## Examples

### 01_basic_usage.go

**Basic Custom Entity Operations**

Demonstrates:
- Initializing CMS store with custom entities enabled
- Defining a simple entity type (Product)
- Creating, retrieving, updating, and listing entities
- Soft delete functionality
- Counting entities

**Run:**
```bash
cd examples/customentities
go run 01_basic_usage.go
```

**Key Concepts:**
- Entity type definitions with attributes
- Required vs optional attributes
- Type-safe attribute storage (string, int, float, bool)
- CRUD operations

---

### 02_relationships.go

**Entity Relationships**

Demonstrates:
- Creating entities with relationships
- Linking authors to blog posts
- Creating comments on posts
- Querying relationships
- Building hierarchical data structures

**Run:**
```bash
cd examples/customentities
go run 02_relationships.go
```

**Key Concepts:**
- Relationship types (belongs_to, has_many)
- Relationship metadata
- Querying related entities
- Multi-level relationships (Author → Post → Comment)

---

### 03_taxonomy.go

**Taxonomy and Categorization**

Demonstrates:
- Creating taxonomies (categories, tags)
- Creating taxonomy terms
- Assigning entities to multiple taxonomies
- Querying entities by taxonomy
- Hierarchical categorization

**Run:**
```bash
cd examples/customentities
go run 03_taxonomy.go
```

**Key Concepts:**
- Taxonomy vs taxonomy terms
- Multiple taxonomies per entity type
- Filtering by category/tag
- Taxonomy metadata (slug, sort order)

---

## Quick Start

### 1. Basic Setup

```go
import (
    "github.com/dracory/cmsstore"
    "github.com/dracory/entitystore"
)

// Create store with custom entities enabled
store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
    DB:                    db,
    CustomEntitiesEnabled: true,
    CustomEntityStoreOptions: cmsstore.CustomEntityStoreOptions{
        EntityTableName:    "cms_custom_entity",
        AttributeTableName: "cms_custom_attribute",
        AutomigrateEnabled: true,
    },
    CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
        {
            Type:      "product",
            TypeLabel: "Product",
            Group:     "Shop",
            Attributes: []cmsstore.CustomAttributeDefinition{
                {Name: "title", Type: "string", Label: "Title", Required: true},
                {Name: "price", Type: "float", Label: "Price", Required: true},
            },
        },
    },
})
```

### 2. Create an Entity

```go
customStore := store.CustomEntityStore()

productID, err := customStore.Create(ctx, "product", map[string]interface{}{
    "title": "Laptop",
    "price": 999.99,
}, nil, nil)
```

### 3. Retrieve an Entity

```go
product, err := customStore.FindByID(ctx, productID)
fmt.Printf("Product ID: %s\n", product.ID())
fmt.Printf("Type: %s\n", product.EntityType())
```

### 4. List Entities

```go
products, err := customStore.List(ctx, entitystore.EntityQueryOptions{
    EntityType: "product",
})

for _, p := range products {
    fmt.Printf("Product: %s\n", p.ID())
}
```

## Features

### Attribute Types

Supported attribute types:
- `string` - Text values
- `int` - Integer numbers
- `float` - Decimal numbers
- `bool` - Boolean values (stored as 0/1)

### Relationships (Optional)

Enable relationships between entities:

```go
CustomEntityStoreOptions: cmsstore.CustomEntityStoreOptions{
    RelationshipsEnabled: true,
}
```

Create entities with relationships:

```go
postID, err := customStore.Create(ctx, "post", attrs, 
    []cmsstore.RelationshipDefinition{
        {
            Type:     "belongs_to",
            TargetID: authorID,
            Metadata: map[string]interface{}{
                "role": "author",
            },
        },
    }, nil)
```

### Taxonomies (Optional)

Enable taxonomies for categorization:

```go
CustomEntityStoreOptions: cmsstore.CustomEntityStoreOptions{
    TaxonomiesEnabled: true,
}
```

Create and assign taxonomies:

```go
// Create taxonomy
taxonomy, err := innerStore.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
    Name: "Categories",
    Slug: "categories",
})

// Create term
term, err := innerStore.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
    TaxonomyID: taxonomy.ID(),
    Name:       "Electronics",
    Slug:       "electronics",
})

// Assign to entity
productID, err := customStore.Create(ctx, "product", attrs, nil, []string{term.ID()})
```

## Database Tables

When custom entities are enabled, the following tables are created:

### Core Tables
- `cms_custom_entity` - Entity records
- `cms_custom_attribute` - EAV attribute storage
- `cms_custom_entity_trash` - Soft-deleted entities
- `cms_custom_attribute_trash` - Soft-deleted attributes

### Relationship Tables (if enabled)
- `cms_custom_relationship` - Entity relationships
- `cms_custom_relationship_trash` - Soft-deleted relationships

### Taxonomy Tables (if enabled)
- `cms_custom_taxonomy` - Taxonomy definitions
- `cms_custom_taxonomy_term` - Taxonomy terms
- `cms_custom_entity_taxonomy` - Entity-term assignments
- Plus corresponding trash tables

## API Reference

### CustomEntityStore Methods

```go
// Create a new entity
Create(ctx, entityType, attrs, relationships, taxonomyTermIDs) (string, error)

// Find entity by ID
FindByID(ctx, entityID) (entitystore.EntityInterface, error)

// List entities with filtering
List(ctx, options entitystore.EntityQueryOptions) ([]entitystore.EntityInterface, error)

// Update entity attributes
Update(ctx, entity, attrs) error

// Soft delete entity
Delete(ctx, entityID) error

// Count entities
Count(ctx, options entitystore.EntityQueryOptions) (int64, error)

// Get entity relationships
GetRelationships(ctx, entityID) ([]entitystore.RelationshipInterface, error)

// Get taxonomy assignments
GetTaxonomyAssignments(ctx, entityID) ([]entitystore.EntityTaxonomyInterface, error)

// Access underlying entitystore for advanced operations
Inner() entitystore.StoreInterface
```

### Entity Definition

```go
type CustomEntityDefinition struct {
    Type                 string                      // Entity type identifier
    TypeLabel            string                      // Display name
    Group                string                      // Grouping for admin UI
    Attributes           []CustomAttributeDefinition // Attribute definitions
    AllowRelationships   bool                        // Enable relationships
    AllowedRelationTypes []string                    // Allowed relationship types
    AllowTaxonomies      bool                        // Enable taxonomies
    TaxonomyIDs          []string                    // Allowed taxonomy IDs
}

type CustomAttributeDefinition struct {
    Name     string // Attribute name
    Type     string // Attribute type (string, int, float, bool)
    Label    string // Display label
    Required bool   // Is required?
}
```

## Best Practices

1. **Define entity types at initialization** - Register all entity types when creating the store
2. **Use required attributes** - Mark essential attributes as required for validation
3. **Enable features selectively** - Only enable relationships/taxonomies if needed
4. **Use meaningful slugs** - Taxonomy slugs should be URL-friendly
5. **Leverage metadata** - Store additional context on relationships
6. **Query efficiently** - Use EntityQueryOptions for filtering

## Troubleshooting

### "entity type not registered"
Make sure the entity type is defined in `CustomEntityDefinitions` when creating the store.

### "required attribute missing"
Ensure all required attributes are provided when creating entities.

### "relationships not enabled"
Set `RelationshipsEnabled: true` in `CustomEntityStoreOptions`.

### "taxonomies not enabled"
Set `TaxonomiesEnabled: true` in `CustomEntityStoreOptions`.

## Additional Resources

- [entitystore documentation](https://github.com/dracory/entitystore)
- [Custom Entities Proposal](../../docs/proposals/2026-03-17-custom-entities-support.md)
- [CMS Store Documentation](../../docs/CUSTOM_ENTITIES.md)

## License

Same as parent project.
