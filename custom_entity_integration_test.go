package cmsstore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/entitystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestCustomEntityIntegration(t *testing.T) {
	// Initialize database
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	require.NoError(t, err)
	defer db.Close()

	// Create CMS store with custom entities enabled
	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "cms_block",
		PageTableName:      "cms_page",
		SiteTableName:      "cms_site",
		TemplateTableName:  "cms_template",
		AutomigrateEnabled: true,

		CustomEntitiesEnabled: true,
		CustomEntityStoreOptions: CustomEntityStoreOptions{
			RelationshipsEnabled: false,
			TaxonomiesEnabled:    false,
		},
		CustomEntityDefinitions: []CustomEntityDefinition{
			{
				Type:      "product",
				TypeLabel: "Product",
				Group:     "Shop",
				Attributes: []CustomAttributeDefinition{
					{Name: "title", Type: "string", Label: "Title", Required: true},
					{Name: "price", Type: "float", Label: "Price", Required: true},
					{Name: "stock", Type: "int", Label: "Stock Quantity"},
				},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, store)

	// Verify custom entities are enabled
	assert.True(t, store.CustomEntitiesEnabled())

	// Get custom entity store
	customStore := store.CustomEntityStore()
	require.NotNil(t, customStore)

	// Verify entity type is registered
	def, ok := customStore.GetEntityDefinition("product")
	assert.True(t, ok)
	assert.Equal(t, "product", def.Type)
	assert.Equal(t, "Product", def.TypeLabel)
	assert.Equal(t, "Shop", def.Group)
	assert.Len(t, def.Attributes, 3)

	ctx := context.Background()

	// Test 1: Create a product
	t.Run("Create", func(t *testing.T) {
		attrs := map[string]interface{}{
			"title": "Test Laptop",
			"price": 999.99,
			"stock": 10,
		}

		productID, err := customStore.Create(ctx, "product", attrs, nil, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, productID)

		// Verify entity was created
		entity, err := customStore.FindByID(ctx, productID)
		assert.NoError(t, err)
		assert.NotNil(t, entity)
		assert.Equal(t, "product", entity.GetType())

		// Verify attributes
		titleAttr, err := customStore.Inner().AttributeFind(ctx, entity.ID(), "title")
		assert.NoError(t, err)
		assert.NotNil(t, titleAttr)
		assert.Equal(t, "Test Laptop", titleAttr.GetValue())

		priceAttr, err := customStore.Inner().AttributeFind(ctx, entity.ID(), "price")
		assert.NoError(t, err)
		assert.NotNil(t, priceAttr)
		priceValue, _ := priceAttr.GetFloat()
		assert.Equal(t, 999.99, priceValue)
	})

	// Test 2: Validation - missing required attribute
	t.Run("ValidationRequired", func(t *testing.T) {
		attrs := map[string]interface{}{
			"title": "Incomplete Product",
			// Missing required "price"
		}

		_, err := customStore.Create(ctx, "product", attrs, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required attribute 'price' is missing")
	})

	// Test 3: Unregistered entity type
	t.Run("UnregisteredType", func(t *testing.T) {
		attrs := map[string]interface{}{
			"name": "Test",
		}

		_, err := customStore.Create(ctx, "unknown_type", attrs, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "entity type 'unknown_type' is not registered")
	})

	// Test 4: Update entity
	t.Run("Update", func(t *testing.T) {
		// Create initial entity
		attrs := map[string]interface{}{
			"title": "Original Title",
			"price": 50.00,
		}

		productID, err := customStore.Create(ctx, "product", attrs, nil, nil)
		require.NoError(t, err)

		// Find entity
		entity, err := customStore.FindByID(ctx, productID)
		require.NoError(t, err)

		// Update attributes
		updateAttrs := map[string]interface{}{
			"title": "Updated Title",
			"price": 75.00,
		}

		err = customStore.Update(ctx, entity, updateAttrs)
		assert.NoError(t, err)

		// Verify update
		titleAttr, _ := customStore.Inner().AttributeFind(ctx, productID, "title")
		assert.Equal(t, "Updated Title", titleAttr.GetValue())

		priceAttr, _ := customStore.Inner().AttributeFind(ctx, productID, "price")
		priceValue, _ := priceAttr.GetFloat()
		assert.Equal(t, 75.00, priceValue)
	})

	// Test 5: Delete entity
	t.Run("Delete", func(t *testing.T) {
		// Create entity
		attrs := map[string]interface{}{
			"title": "To Delete",
			"price": 100.00,
		}

		productID, err := customStore.Create(ctx, "product", attrs, nil, nil)
		require.NoError(t, err)

		// Delete entity
		err = customStore.Delete(ctx, productID)
		assert.NoError(t, err)

		// Verify soft delete (entity should not be found)
		entity, err := customStore.FindByID(ctx, productID)
		assert.NoError(t, err)
		assert.Nil(t, entity)
	})

	// Test 6: Count entities
	t.Run("Count", func(t *testing.T) {
		// Create multiple entities
		for i := 1; i <= 3; i++ {
			attrs := map[string]interface{}{
				"title": "Product " + string(rune('0'+i)),
				"price": float64(i * 100),
			}
			_, err := customStore.Create(ctx, "product", attrs, nil, nil)
			require.NoError(t, err)
		}

		// Count entities
		count, err := customStore.Count(ctx, entitystore.EntityQueryOptions{
			EntityType: "product",
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(3))
	})

	// Test 7: Get all definitions
	t.Run("GetAllDefinitions", func(t *testing.T) {
		defs := customStore.GetAllDefinitions()
		assert.Len(t, defs, 1)
		assert.Equal(t, "product", defs[0].Type)
	})
}

func TestCustomEntityWithRelationships(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	require.NoError(t, err)
	defer db.Close()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "cms_block",
		PageTableName:      "cms_page",
		SiteTableName:      "cms_site",
		TemplateTableName:  "cms_template",
		AutomigrateEnabled: true,

		CustomEntitiesEnabled: true,
		CustomEntityStoreOptions: CustomEntityStoreOptions{
			RelationshipsEnabled: true,
		},
		CustomEntityDefinitions: []CustomEntityDefinition{
			{
				Type:               "author",
				TypeLabel:          "Author",
				AllowRelationships: true,
				Attributes: []CustomAttributeDefinition{
					{Name: "name", Type: "string", Label: "Name", Required: true},
				},
			},
			{
				Type:                 "book",
				TypeLabel:            "Book",
				AllowRelationships:   true,
				AllowedRelationTypes: []string{"belongs_to"},
				Attributes: []CustomAttributeDefinition{
					{Name: "title", Type: "string", Label: "Title", Required: true},
				},
			},
		},
	})
	require.NoError(t, err)

	ctx := context.Background()
	customStore := store.CustomEntityStore()

	// Create author
	authorAttrs := map[string]interface{}{"name": "John Doe"}
	authorID, err := customStore.Create(ctx, "author", authorAttrs, nil, nil)
	require.NoError(t, err)

	// Create book with relationship
	bookAttrs := map[string]interface{}{"title": "My Book"}
	relationships := []RelationshipDefinition{
		{
			TargetID: authorID,
			Type:     "belongs_to",
			Metadata: map[string]interface{}{"role": "author"},
		},
	}

	bookID, err := customStore.Create(ctx, "book", bookAttrs, relationships, nil)
	assert.NoError(t, err)

	// Verify relationship
	rels, err := customStore.GetRelationships(ctx, bookID)
	assert.NoError(t, err)
	assert.Len(t, rels, 1)
	assert.Equal(t, authorID, rels[0].GetRelatedEntityID())
	assert.Equal(t, "belongs_to", rels[0].GetRelationshipType())
}
