package cmsstore

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/dracory/entitystore"
	_ "modernc.org/sqlite"
)

func TestCustomEntityIntegration(t *testing.T) {
	// Initialize database
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
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
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	if store == nil {
		t.Fatal("store is nil")
	}

	// Verify custom entities are enabled
	if !store.CustomEntitiesEnabled() {
		t.Error("expected CustomEntitiesEnabled to be true")
	}

	// Get custom entity store
	customStore := store.CustomEntityStore()
	if customStore == nil {
		t.Fatal("customStore is nil")
	}

	// Verify entity type is registered
	def, ok := customStore.GetEntityDefinition("product")
	if !ok {
		t.Error("expected entity definition 'product' to be registered")
	}
	if def.Type != "product" {
		t.Errorf("expected type 'product', got %q", def.Type)
	}
	if def.TypeLabel != "Product" {
		t.Errorf("expected type label 'Product', got %q", def.TypeLabel)
	}
	if def.Group != "Shop" {
		t.Errorf("expected group 'Shop', got %q", def.Group)
	}
	if len(def.Attributes) != 3 {
		t.Errorf("expected 3 attributes, got %d", len(def.Attributes))
	}

	ctx := context.Background()

	// Test 1: Create a product
	t.Run("Create", func(t *testing.T) {
		attrs := map[string]interface{}{
			"title": "Test Laptop",
			"price": 999.99,
			"stock": 10,
		}

		productID, err := customStore.Create(ctx, "product", attrs, nil, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if productID == "" {
			t.Error("expected non-empty productID")
		}

		// Verify entity was created
		entity, err := customStore.FindByID(ctx, productID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if entity == nil {
			t.Fatal("entity is nil")
		}
		if entity.GetType() != "product" {
			t.Errorf("expected type 'product', got %q", entity.GetType())
		}

		// Verify attributes
		titleAttr, err := customStore.Inner().AttributeFind(ctx, entity.ID(), "title")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if titleAttr == nil {
			t.Fatal("titleAttr is nil")
		}
		if titleAttr.GetValue() != "Test Laptop" {
			t.Errorf("expected title 'Test Laptop', got %q", titleAttr.GetValue())
		}

		priceAttr, err := customStore.Inner().AttributeFind(ctx, entity.ID(), "price")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if priceAttr == nil {
			t.Fatal("priceAttr is nil")
		}
		priceValue, _ := priceAttr.GetFloat()
		if priceValue != 999.99 {
			t.Errorf("expected price 999.99, got %f", priceValue)
		}
	})

	// Test 2: Validation - missing required attribute
	t.Run("ValidationRequired", func(t *testing.T) {
		attrs := map[string]interface{}{
			"title": "Incomplete Product",
			// Missing required "price"
		}

		_, err := customStore.Create(ctx, "product", attrs, nil, nil)
		if err == nil {
			t.Error("expected error for missing required attribute")
		}
		if !strings.Contains(err.Error(), "required attribute 'price' is missing") {
			t.Errorf("expected error to contain 'required attribute 'price' is missing', got %q", err.Error())
		}
	})

	// Test 3: Unregistered entity type
	t.Run("UnregisteredType", func(t *testing.T) {
		attrs := map[string]interface{}{
			"name": "Test",
		}

		_, err := customStore.Create(ctx, "unknown_type", attrs, nil, nil)
		if err == nil {
			t.Error("expected error for unregistered entity type")
		}
		if !strings.Contains(err.Error(), "entity type 'unknown_type' is not registered") {
			t.Errorf("expected error to contain 'entity type 'unknown_type' is not registered', got %q", err.Error())
		}
	})

	// Test 4: Update entity
	t.Run("Update", func(t *testing.T) {
		// Create initial entity
		attrs := map[string]interface{}{
			"title": "Original Title",
			"price": 50.00,
		}

		productID, err := customStore.Create(ctx, "product", attrs, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Find entity
		entity, err := customStore.FindByID(ctx, productID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Update attributes
		updateAttrs := map[string]interface{}{
			"title": "Updated Title",
			"price": 75.00,
		}

		err = customStore.Update(ctx, entity, updateAttrs)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify update
		titleAttr, _ := customStore.Inner().AttributeFind(ctx, productID, "title")
		if titleAttr.GetValue() != "Updated Title" {
			t.Errorf("expected title 'Updated Title', got %q", titleAttr.GetValue())
		}

		priceAttr, _ := customStore.Inner().AttributeFind(ctx, productID, "price")
		priceValue, _ := priceAttr.GetFloat()
		if priceValue != 75.00 {
			t.Errorf("expected price 75.00, got %f", priceValue)
		}
	})

	// Test 5: Delete entity
	t.Run("Delete", func(t *testing.T) {
		// Create entity
		attrs := map[string]interface{}{
			"title": "To Delete",
			"price": 100.00,
		}

		productID, err := customStore.Create(ctx, "product", attrs, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Delete entity
		err = customStore.Delete(ctx, productID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify soft delete (entity should not be found)
		entity, err := customStore.FindByID(ctx, productID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if entity != nil {
			t.Error("expected entity to be nil after delete")
		}
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
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		// Count entities
		count, err := customStore.Count(ctx, entitystore.EntityQueryOptions{
			EntityType: "product",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if count < 3 {
			t.Errorf("expected count >= 3, got %d", count)
		}
	})

	// Test 7: Get all definitions
	t.Run("GetAllDefinitions", func(t *testing.T) {
		defs := customStore.GetAllDefinitions()
		if len(defs) != 1 {
			t.Errorf("expected 1 definition, got %d", len(defs))
		}
		if defs[0].Type != "product" {
			t.Errorf("expected type 'product', got %q", defs[0].Type)
		}
	})
}

func TestCustomEntityWithRelationships(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
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
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()
	customStore := store.CustomEntityStore()

	// Create author
	authorAttrs := map[string]interface{}{"name": "John Doe"}
	authorID, err := customStore.Create(ctx, "author", authorAttrs, nil, nil)
	if err != nil {
		t.Fatalf("failed to create author: %v", err)
	}

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
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify relationship
	rels, err := customStore.GetRelationships(ctx, bookID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(rels) != 1 {
		t.Errorf("expected 1 relationship, got %d", len(rels))
	}
	if rels[0].GetRelatedEntityID() != authorID {
		t.Errorf("expected related entity ID %q, got %q", authorID, rels[0].GetRelatedEntityID())
	}
	if rels[0].GetRelationshipType() != "belongs_to" {
		t.Errorf("expected relationship type 'belongs_to', got %q", rels[0].GetRelationshipType())
	}
}
