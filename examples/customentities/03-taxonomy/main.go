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

// Example 3: Taxonomy and Categorization
// This example demonstrates:
// - Creating taxonomies (categories, tags)
// - Creating taxonomy terms
// - Assigning entities to taxonomy terms
// - Querying entities by taxonomy

func main() {
	ctx := context.Background()

	// Initialize SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create CMS store with custom entities and taxonomies enabled
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
			TaxonomiesEnabled:  true,
			AutomigrateEnabled: true,
		},
		CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
			{
				Type:            "product",
				TypeLabel:       "Product",
				Group:           "Shop",
				AllowTaxonomies: true,
				Attributes: []cmsstore.CustomAttributeDefinition{
					{Name: "name", Type: "string", Label: "Product Name", Required: true},
					{Name: "price", Type: "float", Label: "Price", Required: true},
					{Name: "description", Type: "string", Label: "Description"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	customStore := store.CustomEntityStore()
	innerStore := customStore.Inner()

	// Example 1: Create a taxonomy for product categories
	fmt.Println("=== Creating Product Categories Taxonomy ===")
	categoryTaxonomy, err := innerStore.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name:        "Product Categories",
		Slug:        "product-categories",
		Description: "Main product categorization",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created taxonomy: %s (ID: %s)\n\n", categoryTaxonomy.GetName(), categoryTaxonomy.ID())

	// Example 2: Create taxonomy terms (categories)
	fmt.Println("=== Creating Category Terms ===")

	electronicsCategory, err := innerStore.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: categoryTaxonomy.ID(),
		Name:       "Electronics",
		Slug:       "electronics",
		SortOrder:  1,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created category: %s (ID: %s)\n", electronicsCategory.GetName(), electronicsCategory.ID())

	clothingCategory, err := innerStore.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: categoryTaxonomy.ID(),
		Name:       "Clothing",
		Slug:       "clothing",
		SortOrder:  2,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created category: %s (ID: %s)\n", clothingCategory.GetName(), clothingCategory.ID())

	booksCategory, err := innerStore.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: categoryTaxonomy.ID(),
		Name:       "Books",
		Slug:       "books",
		SortOrder:  3,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created category: %s (ID: %s)\n\n", booksCategory.GetName(), booksCategory.ID())

	// Example 3: Create a tags taxonomy
	fmt.Println("=== Creating Tags Taxonomy ===")
	tagsTaxonomy, err := innerStore.TaxonomyCreateByOptions(ctx, entitystore.TaxonomyOptions{
		Name:        "Product Tags",
		Slug:        "product-tags",
		Description: "Product tags for filtering",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created taxonomy: %s (ID: %s)\n\n", tagsTaxonomy.GetName(), tagsTaxonomy.ID())

	// Example 4: Create tag terms
	fmt.Println("=== Creating Tag Terms ===")

	newTag, err := innerStore.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tagsTaxonomy.ID(),
		Name:       "New",
		Slug:       "new",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created tag: %s (ID: %s)\n", newTag.GetName(), newTag.ID())

	saleTag, err := innerStore.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tagsTaxonomy.ID(),
		Name:       "Sale",
		Slug:       "sale",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created tag: %s (ID: %s)\n", saleTag.GetName(), saleTag.ID())

	featuredTag, err := innerStore.TaxonomyTermCreateByOptions(ctx, entitystore.TaxonomyTermOptions{
		TaxonomyID: tagsTaxonomy.ID(),
		Name:       "Featured",
		Slug:       "featured",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created tag: %s (ID: %s)\n\n", featuredTag.GetName(), featuredTag.ID())

	// Example 5: Create products with taxonomy assignments
	fmt.Println("=== Creating Products with Taxonomy Assignments ===")

	// Register taxonomy IDs that products can use
	productDef, _ := customStore.GetEntityDefinition("product")
	productDef.TaxonomyIDs = []string{categoryTaxonomy.ID(), tagsTaxonomy.ID()}
	customStore.RegisterEntityType(productDef)

	// Product 1: Laptop in Electronics, tagged as New and Featured
	laptop, err := customStore.Create(ctx, "product", map[string]interface{}{
		"name":        "Gaming Laptop",
		"price":       1299.99,
		"description": "High-performance gaming laptop",
	}, nil, []string{electronicsCategory.ID(), newTag.ID(), featuredTag.ID()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created product: Gaming Laptop (ID: %s)\n", laptop)

	// Product 2: T-Shirt in Clothing, tagged as Sale
	tshirt, err := customStore.Create(ctx, "product", map[string]interface{}{
		"name":        "Cotton T-Shirt",
		"price":       19.99,
		"description": "Comfortable cotton t-shirt",
	}, nil, []string{clothingCategory.ID(), saleTag.ID()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created product: Cotton T-Shirt (ID: %s)\n", tshirt)

	// Product 3: Book in Books, tagged as New
	book, err := customStore.Create(ctx, "product", map[string]interface{}{
		"name":        "Go Programming Guide",
		"price":       39.99,
		"description": "Comprehensive guide to Go programming",
	}, nil, []string{booksCategory.ID(), newTag.ID()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created product: Go Programming Guide (ID: %s)\n\n", book)

	// Example 6: Query taxonomy assignments for a product
	fmt.Println("=== Querying Taxonomy Assignments for Laptop ===")
	laptopTaxonomies, err := customStore.GetTaxonomyAssignments(ctx, laptop)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Laptop has %d taxonomy assignment(s):\n", len(laptopTaxonomies))
	for _, assignment := range laptopTaxonomies {
		term, _ := innerStore.TaxonomyTermFind(ctx, assignment.GetTermID())
		if term != nil {
			fmt.Printf("- Term: %s (Taxonomy ID: %s)\n", term.GetName(), assignment.GetTaxonomyID())
		}
	}
	fmt.Println()

	// Example 7: Find all products in Electronics category
	fmt.Println("=== Finding All Products in Electronics Category ===")
	electronicsProducts, err := innerStore.EntityTaxonomyList(ctx, entitystore.EntityTaxonomyQueryOptions{
		TaxonomyID: categoryTaxonomy.ID(),
		TermID:     electronicsCategory.ID(),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d product(s) in Electronics:\n", len(electronicsProducts))
	for i, assignment := range electronicsProducts {
		product, _ := customStore.FindByID(ctx, assignment.GetEntityID())
		if product != nil {
			fmt.Printf("%d. Product ID: %s (Type: %s)\n", i+1, product.ID(), product.GetType())
		}
	}
	fmt.Println()

	// Example 8: Find all products tagged as "New"
	fmt.Println("=== Finding All Products Tagged as 'New' ===")
	newProducts, err := innerStore.EntityTaxonomyList(ctx, entitystore.EntityTaxonomyQueryOptions{
		TaxonomyID: tagsTaxonomy.ID(),
		TermID:     newTag.ID(),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d product(s) tagged as 'New':\n", len(newProducts))
	for i, assignment := range newProducts {
		product, _ := customStore.FindByID(ctx, assignment.GetEntityID())
		if product != nil {
			fmt.Printf("%d. Product ID: %s\n", i+1, product.ID())
		}
	}
	fmt.Println()

	// Example 9: List all terms in a taxonomy
	fmt.Println("=== Listing All Category Terms ===")
	allCategories, err := innerStore.TaxonomyTermList(ctx, entitystore.TaxonomyTermQueryOptions{
		TaxonomyID: categoryTaxonomy.ID(),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Product Categories taxonomy has %d term(s):\n", len(allCategories))
	for i, term := range allCategories {
		fmt.Printf("%d. %s (slug: %s, order: %d)\n", i+1, term.GetName(), term.GetSlug(), term.GetSortOrder())
	}
	fmt.Println()

	// Example 10: Summary
	fmt.Println("=== Summary ===")
	fmt.Println("Taxonomy structure:")
	fmt.Println("- Product Categories (taxonomy)")
	fmt.Println("  - Electronics (term) → 1 product")
	fmt.Println("  - Clothing (term) → 1 product")
	fmt.Println("  - Books (term) → 1 product")
	fmt.Println("- Product Tags (taxonomy)")
	fmt.Println("  - New (term) → 2 products")
	fmt.Println("  - Sale (term) → 1 product")
	fmt.Println("  - Featured (term) → 1 product")
	fmt.Println("\nTaxonomies enable:")
	fmt.Println("✓ Categorizing entities")
	fmt.Println("✓ Multi-level classification (categories + tags)")
	fmt.Println("✓ Filtering and searching by taxonomy")
	fmt.Println("✓ Organizing content hierarchically")
}
