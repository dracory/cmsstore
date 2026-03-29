package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dracory/cmsstore"
	"github.com/dracory/entitystore"
	_ "modernc.org/sqlite"
)

// Example 1: Basic Custom Entity Usage
// This example demonstrates:
// - Initializing the CMS store with custom entities enabled
// - Defining a simple custom entity type (Product)
// - Creating, retrieving, updating, and listing entities

func main() {
	ctx := context.Background()

	// Initialize SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create CMS store with custom entities enabled
	store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                    db,
		BlockTableName:        "cms_block",
		PageTableName:         "cms_page",
		SiteTableName:         "cms_site",
		TemplateTableName:     "cms_template",
		AutomigrateEnabled:    true,
		CustomEntitiesEnabled: true,
		CustomEntityStoreOptions: cmsstore.CustomEntityStoreOptions{
			EntityTableName:    "cms_custom_entity",
			AttributeTableName: "cms_custom_attribute",
		},
		CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
			{
				Type:      "product",
				TypeLabel: "Product",
				Group:     "Shop",
				Attributes: []cmsstore.CustomAttributeDefinition{
					{Name: "title", Type: "string", Label: "Product Title", Required: true},
					{Name: "description", Type: "string", Label: "Description"},
					{Name: "price", Type: "float", Label: "Price (USD)", Required: true},
					{Name: "stock", Type: "int", Label: "Stock Quantity", Required: true},
					{Name: "active", Type: "bool", Label: "Active"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	customStore := store.CustomEntityStore()

	// Example 1: Create a product
	fmt.Println("=== Creating Product ===")
	productID, err := customStore.Create(ctx, "product", map[string]interface{}{
		"title":       "Laptop Pro 15",
		"description": "High-performance laptop with 16GB RAM",
		"price":       1299.99,
		"stock":       25,
		"active":      true,
	}, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created product with ID: %s\n\n", productID)

	// Example 2: Retrieve the product
	fmt.Println("=== Retrieving Product ===")
	product, err := customStore.FindByID(ctx, productID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Product ID: %s\n", product.ID())
	fmt.Printf("Type: %s\n", product.EntityType())
	fmt.Printf("Created: %s\n\n", product.CreatedAt())

	// Example 3: Update the product
	fmt.Println("=== Updating Product ===")
	err = customStore.Update(ctx, product, map[string]interface{}{
		"price": 1199.99, // Price reduction
		"stock": 20,      // Stock sold
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Product updated successfully\n")

	// Example 4: Create more products
	fmt.Println("=== Creating More Products ===")
	products := []map[string]interface{}{
		{
			"title":       "Wireless Mouse",
			"description": "Ergonomic wireless mouse",
			"price":       29.99,
			"stock":       100,
			"active":      true,
		},
		{
			"title":       "Mechanical Keyboard",
			"description": "RGB mechanical keyboard",
			"price":       89.99,
			"stock":       50,
			"active":      true,
		},
		{
			"title":       "USB-C Hub",
			"description": "7-in-1 USB-C hub",
			"price":       49.99,
			"stock":       0, // Out of stock
			"active":      false,
		},
	}

	for _, attrs := range products {
		id, err := customStore.Create(ctx, "product", attrs, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created product: %s (ID: %s)\n", attrs["title"], id)
	}
	fmt.Println()

	// Example 5: List all products
	fmt.Println("=== Listing All Products ===")
	allProducts, err := customStore.List(ctx, entitystore.EntityQueryOptions{
		EntityType: "product",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total products: %d\n\n", len(allProducts))

	for i, p := range allProducts {
		fmt.Printf("%d. Product ID: %s (Type: %s)\n",
			i+1,
			p.ID(),
			p.EntityType(),
		)
	}
	fmt.Println()

	// Example 6: Count products
	fmt.Println("=== Counting Products ===")
	count, err := customStore.Count(ctx, entitystore.EntityQueryOptions{
		EntityType: "product",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total product count: %d\n\n", count)

	// Example 7: Delete a product (soft delete)
	fmt.Println("=== Deleting Product ===")
	err = customStore.Delete(ctx, productID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Product %s moved to trash\n\n", productID)

	// Verify deletion
	deletedProduct, err := customStore.FindByID(ctx, productID)
	if err != nil {
		fmt.Printf("Product not found (as expected): %v\n", err)
	} else if deletedProduct == nil {
		fmt.Println("Product successfully deleted")
	}

	// Count after deletion
	count, _ = customStore.Count(ctx, entitystore.EntityQueryOptions{
		EntityType: "product",
	})
	fmt.Printf("Remaining products: %d\n", count)
}
