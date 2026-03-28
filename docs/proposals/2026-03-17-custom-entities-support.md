# [Draft] Implement Custom Entities Support

## Status
**[Draft]** - Planned, never implemented, but **ready to implement** using existing libraries

## Summary
- **Problem**: The CMS only supports predefined entity types (pages, blocks, menus, templates, translations) with no way to extend for custom business logic
- **Solution**: Integrate `dracory/entitystore` for EAV storage and `dracory/crud` for admin interface, eliminating need to build from scratch
- **Opportunity**: Existing libraries provide 90% of required functionality

## Available Libraries

### 1. dracory/entitystore - EAV Storage Layer

**Purpose:** Schemaless SQL storage using Entity-Attribute-Value pattern

**Key Features:**
- **No schema changes** - Add entity types without migrations
- **EAV pattern** - Keeps relational structure, avoids JSON blobs
- **Type safety** - `SetString()`, `SetInt()`, `SetFloat()`, `SetInterface()`
- **Soft deletes** - Built-in trash bin with `EntityTrash()`
- **SQL reporting** - Full SQL access for complex queries
- **99% functionality** out of the box

**Usage Pattern:**
```go
import "github.com/dracory/entitystore"

// Initialize
entityStore, err := entitystore.NewStore(entitystore.NewStoreOptions{
    DB:                 db,
    EntityTableName:    "cms_custom_entity",
    AttributeTableName: "cms_custom_attribute",
    AutomigrateEnabled: true,
})

// Create entity
person := entityStore.EntityCreateWithType("user")
person.SetString("first_name", "John")
person.SetString("last_name", "Doe")
person.SetInt("age", 32)

// Retrieve
person := entityStore.EntityFindByID(entityID)
name := person.GetString("first_name", "")
age := person.GetInt("age", 0)
```

**Available Methods:**
- Store: `EntityCreate`, `EntityFindByID`, `EntityList`, `EntityTrash`, `EntityCount`
- Entity: `GetString`, `GetInt`, `GetFloat`, `SetString`, `SetInt`, `SetFloat`
- Attributes: `AttributeSetString`, `AttributeFind`

### 2. dracory/crud - Admin Interface

**Purpose:** Server-side rendered CRUD interface for Go web applications

**Key Features:**
- **Bootstrap 5** - Matches existing CMS admin UI
- **Vue.js 3** - Interactive forms with reactivity
- **Server-side rendered** - No API layer needed
- **Zero custom CSS** - Uses Bootstrap components

**What It Generates:**
- List view with pagination, search, sorting
- Create/Edit forms with validation
- Soft delete (trash) functionality
- Bulk operations

## Implementation Strategy

### Phase 1: Storage Integration (1 week)

Instead of building custom EAV tables, wrap `entitystore`:

```go
// custom_entity_store.go
package cmsstore

import (
    "github.com/dracory/entitystore"
)

type CustomEntityStore struct {
    inner *entitystore.Store
}

func NewCustomEntityStore(db *sql.DB) *CustomEntityStore {
    inner, _ := entitystore.NewStore(entitystore.NewStoreOptions{
        DB:                 db,
        EntityTableName:    "cms_custom_entity",
        AttributeTableName: "cms_custom_attribute",
        AutomigrateEnabled: true,
    })
    return &CustomEntityStore{inner: inner}
}

// Delegate to entitystore with CMS-specific logic
func (s *CustomEntityStore) Create(entityType string, attrs map[string]interface{}) (string, error) {
    entity := s.inner.EntityCreateWithType(entityType)
    for k, v := range attrs {
        switch val := v.(type) {
        case string:
            entity.SetString(k, val)
        case int:
            entity.SetInt(k, int64(val))
        case float64:
            entity.SetFloat(k, val)
        }
    }
    return entity.ID(), nil
}
```

**Benefits:**
- No custom schema to maintain
- Battle-tested EAV implementation
- Soft deletes, trash bin included
- Performance optimizations already done

### Phase 2: Admin Integration (1 week)

Instead of building dynamic forms, use `dracory/crud` with configuration:

```go
// admin/customentity/customentity_controller.go
package customentity

import (
    "github.com/dracory/crud"
)

type CustomEntityController struct {
    crud *crud.Crud
    definitions []CustomEntityDefinition
}

func (c *CustomEntityController) RegisterEntityType(def CustomEntityDefinition) {
    // Generate CRUD for this entity type
    c.crud.Register(c.crud.Config{
        EntityType: def.Type,
        TableName:  "cms_custom_entity",
        Columns:    c.buildColumns(def.Attributes),
        FormFields: c.buildFormFields(def.Attributes),
    })
}

func (c *CustomEntityController) buildColumns(attrs []CustomAttributeDefinition) []crud.Column {
    // Map CustomAttributeStructure to crud.Column
}

func (c *CustomEntityController) buildFormFields(attrs []CustomAttributeDefinition) []crud.Field {
    // Map CustomAttributeStructure to crud.Field
}
```

### Phase 3: CMS Integration (3 days)

Add to existing CMS store:

```go
// store.go

type Store struct {
    // ... existing stores ...
    customEntities *CustomEntityStore
}

func (s *Store) CustomEntity() *CustomEntityStore {
    return s.customEntities
}
```

Update `admin/shared/admin_header.go` to add dynamic navigation based on registered entity types.

## Database Schema Comparison

### Original Proposal (Build From Scratch)

```sql
-- 3 tables to create and maintain
CREATE TABLE custom_entity_definitions (...);
CREATE TABLE custom_entity_data (...);
CREATE TABLE custom_entity_relationships (...);
```

### Revised Approach (Use entitystore)

```sql
-- entitystore creates these automatically
CREATE TABLE cms_custom_entity (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(255),           -- entity type (user, product, order)
    status VARCHAR(50),          -- active, inactive, trashed
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    soft_deleted_at TIMESTAMP
);

CREATE TABLE cms_custom_attribute (
    id VARCHAR(255) PRIMARY KEY,
    entity_id VARCHAR(255),      -- FK to entity
    key VARCHAR(255),            -- attribute name
    value TEXT,                  -- serialized value
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE cms_custom_entity_trash;     -- soft delete support
CREATE TABLE cms_custom_attribute_trash;  -- soft delete support
```

**Advantage:** Zero schema maintenance, library handles migrations

## Revised Implementation Status

| Feature | Original Plan | Revised Plan | Library Coverage |
|---------|--------------|--------------|------------------|
| EAV Storage | ❌ Build from scratch | ✅ Use `entitystore` | 100% |
| Store methods | ❌ Custom implementation | ✅ Wrap `entitystore` | 90% |
| Soft deletes | ❌ Custom implementation | ✅ `entitystore` built-in | 100% |
| Admin forms | ❌ Build from scratch | ✅ Use `crud` | 80% |
| Validation | ❌ Custom implementation | ✅ `crud` + app layer | 70% |
| Relationships | ❌ Custom tables | ⚠️ Manual implementation | 0% |
| Versioning | ❌ Not started | ⚠️ Integrate with existing | 0% |

## Implementation Effort Comparison

| Approach | Estimated Time | Risk |
|----------|----------------|------|
| **Original (Build Everything)** | 6-8 weeks | High - EAV is complex |
| **Revised (Use Libraries)** | 2-3 weeks | Low - proven libraries |
| **Time Savings** | **4-5 weeks** | |

## Files to Create (Revised)

1. `custom_entity_store.go` - Wrapper around `entitystore`
2. `custom_entity_definition.go` - CMS entity type configuration
3. `admin/customentity/` - Controllers using `crud`
4. Modify `store.go` - Add `CustomEntity()` method
5. Modify `options.go` - Add `CustomEntityTypes []CustomEntityDefinition`
6. Modify `admin/shared/admin_header.go` - Dynamic entity navigation

## Integration Example

```go
// main.go - CMS initialization
func main() {
    store := cmsstore.NewStore(...)
    
    // Register custom entity types
    store.RegisterCustomEntityType(cmsstore.CustomEntityDefinition{
        Type:      "shop_product",
        TypeLabel: "Product",
        Group:     "Shop",
        Attributes: []cmsstore.CustomAttributeDefinition{
            {Name: "title", Type: "string", Label: "Title"},
            {Name: "price", Type: "float", Label: "Price"},
            {Name: "stock", Type: "int", Label: "Stock Quantity"},
        },
    })
    
    // entitystore handles storage
    // crud generates admin UI
    // CMS provides unified interface
}
```

### 3. Relationship/Taxonomy Support (Add to entitystore)

**Current Gap:** `entitystore` has no native relationship support.

**Proposed Addition to entitystore:**

```go
// entitystore/relationship.go - New file in dracory/entitystore
package entitystore

type Relationship struct {
    ID               string
    EntityID         string
    RelatedEntityID  string
    RelationshipType string // "belongs_to", "has_many", "taxonomy"
    Metadata         map[string]string
}

type Taxonomy struct {
    ID          string
    Name        string
    Slug        string
    ParentID    string // For hierarchical taxonomies
    EntityTypes []string // Which entity types can use this taxonomy
}

type EntityTaxonomy struct {
    EntityID   string
    TaxonomyID string
    TermID     string
}

// Store methods to add
type StoreInterface interface {
    // ... existing methods ...
    
    // Relationships
    RelationshipCreate(rel *Relationship) error
    RelationshipFindByEntity(entityID string) ([]Relationship, error)
    RelationshipDelete(entityID, relatedID, relType string) error
    
    // Taxonomies
    TaxonomyCreate(tax *Taxonomy) error
    TaxonomyFindBySlug(slug string) (*Taxonomy, error)
    TaxonomyList() ([]Taxonomy, error)
    
    // Entity-Taxonomy associations
    EntityTaxonomyAssign(entityID, taxonomyID, termID string) error
    EntityTaxonomyRemove(entityID, taxonomyID string) error
    EntityTaxonomyList(entityID string) ([]EntityTaxonomy, error)
}
```

**Database Schema Addition:**

```sql
-- Relationships (many-to-many between entities)
CREATE TABLE entity_relationships (
    id VARCHAR(255) PRIMARY KEY,
    entity_id VARCHAR(255),
    related_entity_id VARCHAR(255),
    relationship_type VARCHAR(50),
    metadata JSON,
    created_at TIMESTAMP,
    INDEX idx_entity (entity_id),
    INDEX idx_related (related_entity_id)
);

-- Taxonomy definitions (categories, tags, etc.)
CREATE TABLE taxonomies (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    slug VARCHAR(255) UNIQUE,
    parent_id VARCHAR(255),
    entity_types JSON, -- ["product", "post"]
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Taxonomy terms (actual category/tag values)
CREATE TABLE taxonomy_terms (
    id VARCHAR(255) PRIMARY KEY,
    taxonomy_id VARCHAR(255),
    name VARCHAR(255),
    slug VARCHAR(255),
    parent_id VARCHAR(255),
    sort_order INT,
    created_at TIMESTAMP
);

-- Entity-Term associations
CREATE TABLE entity_taxonomies (
    id VARCHAR(255) PRIMARY KEY,
    entity_id VARCHAR(255),
    taxonomy_id VARCHAR(255),
    term_id VARCHAR(255),
    created_at TIMESTAMP,
    UNIQUE KEY unique_entity_term (entity_id, taxonomy_id, term_id)
);
```

**Benefits of Adding to entitystore:**
- All projects using entitystore get relationships/taxonomy for free
- Single implementation, multiple consumers (cmsstore, other projects)
- Consistent API across projects
- EAV pattern extended naturally to relationships

**Updated Coverage Table:**

| Feature | Original | Revised Plan | Library Coverage |
|---------|----------|--------------|------------------|
| EAV Storage | ❌ Build | ✅ `entitystore` | 100% |
| **Relationships** | ❌ Build | ✅ **Add to `entitystore`** | **0% → 100%** |
| **Taxonomy** | ❌ Build | ✅ **Add to `entitystore`** | **0% → 100%** |
| Soft deletes | ❌ Build | ✅ `entitystore` | 100% |
| Admin forms | ❌ Build | ✅ `crud` | 80% |
| Validation | ❌ Build | ✅ `crud` + app | 70% |
| Versioning | ❌ Not started | ⚠️ CMS integration | 0% |

## Implementation Plan (Revised with entitystore Extensions)

### Phase 0: Extend entitystore (1 week)

Add relationship and taxonomy support to `dracory/entitystore`:

1. Create `relationship.go` with Relationship type and methods
2. Create `taxonomy.go` with Taxonomy, TaxonomyTerm, EntityTaxonomy types
3. Add relationship tables to AutoMigrate
4. Add taxonomy tables to AutoMigrate
5. Write tests for new functionality

**This benefits all projects using entitystore, not just cmsstore.**

### Phase 1: CMS Storage Integration (3 days)

Wrap extended `entitystore`:

```go
// cmsstore/custom_entity_store.go
type CustomEntityStore struct {
    inner *entitystore.Store
}

func (s *CustomEntityStore) CreateWithRelationships(
    entityType string, 
    attrs map[string]interface{},
    relationships []Relationship,
    taxonomyIDs []string,
) (string, error) {
    // Create entity via entitystore
    entity := s.inner.EntityCreateWithType(entityType)
    // ... set attributes ...
    
    // Add relationships via entitystore
    for _, rel := range relationships {
        s.inner.RelationshipCreate(&entitystore.Relationship{
            EntityID: entity.ID(),
            RelatedEntityID: rel.TargetID,
            RelationshipType: rel.Type,
        })
    }
    
    // Assign taxonomies via entitystore
    for _, taxID := range taxonomyIDs {
        s.inner.EntityTaxonomyAssign(entity.ID(), taxID, "")
    }
    
    return entity.ID(), nil
}
```

### Phase 2: Admin Integration (3 days)

Use `dracory/crud` + relationship UI components:

```go
// admin/customentity/controller.go
func (c *Controller) buildRelationshipFields(def CustomEntityDefinition) []crud.Field {
    // Dropdown for belongs_to relationships
    // Multi-select for taxonomies
}
```

### Phase 3: CMS Integration (2 days)

Final integration and testing.

## Effort Comparison

| Phase | Original Plan | Revised Plan | Savings |
|-------|---------------|--------------|---------|
| Storage | 3 weeks (build EAV + relationships) | 1 week (extend entitystore) | 2 weeks |
| Admin | 2 weeks | 3 days | 1 week |
| Integration | 1 week | 2 days | 3 days |
| **Total** | **6 weeks** | **2 weeks** | **4 weeks** |

## Recommendation

**Extend `entitystore` with relationships/taxonomy first.**

This approach:
1. **Reusability** - Other projects using entitystore benefit
2. **Single source** - One implementation to maintain
3. **Natural extension** - EAV + relationships is a complete data layer
4. **Low risk** - Build on proven foundation

**Next Actions:**
1. Design relationship API for entitystore (consider HABTM, belongs_to, has_many patterns)
2. Design taxonomy system (hierarchical categories, flat tags)
3. Implement in entitystore (1 week)
4. Update proposal with actual entitystore API
5. Proceed to cmsstore integration