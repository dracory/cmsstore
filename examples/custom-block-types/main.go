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

// Custom Block Types Example
// This example shows how to extend the CMS with custom block types
// using custom entities to store additional block variations.
// Run: go run main.go

func main() {
	ctx := context.Background()

	// Initialize database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create CMS store with custom block types
	store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                    db,
		BlockTableName:        "cms_block",
		PageTableName:         "cms_page",
		SiteTableName:         "cms_site",
		TemplateTableName:     "cms_template",
		AutomigrateEnabled:    true,
		CustomEntitiesEnabled: true,
		CustomEntityStoreOptions: cmsstore.CustomEntityStoreOptions{
			AutomigrateEnabled: true,
		},
		CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
			// Hero Block - Large banner with image and CTA
			{
				Type:      "hero_block",
				TypeLabel: "Hero Block",
				Group:     "Content Blocks",
				Attributes: []cmsstore.CustomAttributeDefinition{
					{Name: "title", Type: "string", Label: "Title", Required: true},
					{Name: "subtitle", Type: "string", Label: "Subtitle"},
					{Name: "background_image", Type: "string", Label: "Background Image URL"},
					{Name: "cta_text", Type: "string", Label: "CTA Button Text"},
					{Name: "cta_link", Type: "string", Label: "CTA Button Link"},
					{Name: "alignment", Type: "string", Label: "Text Alignment"},
				},
			},
			// Feature Block - Icon, title, description
			{
				Type:      "feature_block",
				TypeLabel: "Feature Block",
				Group:     "Content Blocks",
				Attributes: []cmsstore.CustomAttributeDefinition{
					{Name: "icon", Type: "string", Label: "Icon Name", Required: true},
					{Name: "title", Type: "string", Label: "Title", Required: true},
					{Name: "description", Type: "string", Label: "Description"},
					{Name: "link", Type: "string", Label: "Learn More Link"},
				},
			},
			// Testimonial Block - Quote with author
			{
				Type:      "testimonial_block",
				TypeLabel: "Testimonial Block",
				Group:     "Content Blocks",
				Attributes: []cmsstore.CustomAttributeDefinition{
					{Name: "quote", Type: "string", Label: "Quote", Required: true},
					{Name: "author_name", Type: "string", Label: "Author Name", Required: true},
					{Name: "author_title", Type: "string", Label: "Author Title"},
					{Name: "author_image", Type: "string", Label: "Author Image URL"},
					{Name: "rating", Type: "int", Label: "Rating (1-5)"},
				},
			},
			// FAQ Block - Question and answer
			{
				Type:      "faq_block",
				TypeLabel: "FAQ Block",
				Group:     "Content Blocks",
				Attributes: []cmsstore.CustomAttributeDefinition{
					{Name: "question", Type: "string", Label: "Question", Required: true},
					{Name: "answer", Type: "string", Label: "Answer", Required: true},
					{Name: "category", Type: "string", Label: "Category"},
					{Name: "order", Type: "int", Label: "Display Order"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	customStore := store.CustomEntityStore()

	fmt.Println("=== Custom Block Types Example ===")
	fmt.Println("Demonstrating how to extend CMS with custom block types")

	// Example 1: Create a Hero Block
	fmt.Println("=== Creating Hero Block ===")
	heroID, err := customStore.Create(ctx, "hero_block", map[string]interface{}{
		"title":            "Welcome to Our Platform",
		"subtitle":         "Build amazing websites with custom blocks",
		"background_image": "/images/hero-bg.jpg",
		"cta_text":         "Get Started",
		"cta_link":         "/signup",
		"alignment":        "center",
	}, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Created hero block: %s\n", heroID)

	// Example 2: Create Feature Blocks
	fmt.Println("\n=== Creating Feature Blocks ===")
	features := []map[string]interface{}{
		{
			"icon":        "zap",
			"title":       "Fast Performance",
			"description": "Lightning-fast load times and optimized delivery",
			"link":        "/features/performance",
		},
		{
			"icon":        "shield",
			"title":       "Secure by Default",
			"description": "Enterprise-grade security built in",
			"link":        "/features/security",
		},
		{
			"icon":        "code",
			"title":       "Developer Friendly",
			"description": "Clean APIs and extensive documentation",
			"link":        "/features/developer",
		},
	}

	for _, feature := range features {
		id, err := customStore.Create(ctx, "feature_block", feature, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("✓ Created feature: %s\n", feature["title"])
		_ = id
	}

	// Example 3: Create Testimonial Blocks
	fmt.Println("\n=== Creating Testimonial Blocks ===")
	testimonials := []map[string]interface{}{
		{
			"quote":        "This platform transformed how we build websites. Highly recommended!",
			"author_name":  "Sarah Johnson",
			"author_title": "CTO, TechCorp",
			"author_image": "/images/sarah.jpg",
			"rating":       5,
		},
		{
			"quote":        "The custom blocks feature is a game-changer for our content team.",
			"author_name":  "Mike Chen",
			"author_title": "Marketing Director",
			"author_image": "/images/mike.jpg",
			"rating":       5,
		},
	}

	for _, testimonial := range testimonials {
		id, err := customStore.Create(ctx, "testimonial_block", testimonial, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("✓ Created testimonial from: %s\n", testimonial["author_name"])
		_ = id
	}

	// Example 4: Create FAQ Blocks
	fmt.Println("\n=== Creating FAQ Blocks ===")
	faqs := []map[string]interface{}{
		{
			"question": "How do custom blocks work?",
			"answer":   "Custom blocks extend the CMS with new content types without database migrations.",
			"category": "Getting Started",
			"order":    1,
		},
		{
			"question": "Can I create my own block types?",
			"answer":   "Yes! Define any block type with custom attributes to fit your needs.",
			"category": "Getting Started",
			"order":    2,
		},
		{
			"question": "Are custom blocks searchable?",
			"answer":   "Yes, all custom block content is stored in the database and fully searchable.",
			"category": "Features",
			"order":    3,
		},
	}

	for _, faq := range faqs {
		id, err := customStore.Create(ctx, "faq_block", faq, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("✓ Created FAQ: %s\n", faq["question"])
		_ = id
	}

	// Example 5: List all blocks by type
	fmt.Println("\n=== Listing Blocks by Type ===")

	blockTypes := []string{"hero_block", "feature_block", "testimonial_block", "faq_block"}
	for _, blockType := range blockTypes {
		blocks, err := customStore.List(ctx, entitystore.EntityQueryOptions{
			EntityType: blockType,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-20s: %d block(s)\n", blockType, len(blocks))
	}

	// Example 6: Count all custom blocks
	fmt.Println("\n=== Total Custom Blocks ===")
	totalCount := int64(0)
	for _, blockType := range blockTypes {
		count, err := customStore.Count(ctx, entitystore.EntityQueryOptions{
			EntityType: blockType,
		})
		if err != nil {
			log.Fatal(err)
		}
		totalCount += count
	}
	fmt.Printf("Total custom blocks created: %d\n", totalCount)

	// Example 7: Retrieve and display a specific block
	fmt.Println("\n=== Retrieving Hero Block ===")
	hero, err := customStore.FindByID(ctx, heroID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Block Type: %s\n", hero.GetType())
	fmt.Printf("Block ID: %s\n", hero.ID())
	fmt.Printf("Created: %s\n", hero.GetCreatedAt())

	// Summary
	fmt.Println("\n=== Summary ===")
	fmt.Println("This example demonstrated:")
	fmt.Println("✓ Defining custom block types")
	fmt.Println("✓ Creating hero blocks with CTAs")
	fmt.Println("✓ Creating feature blocks with icons")
	fmt.Println("✓ Creating testimonial blocks with ratings")
	fmt.Println("✓ Creating FAQ blocks with categories")
	fmt.Println("✓ Listing and counting blocks by type")
	fmt.Println()
	fmt.Println("Use Cases:")
	fmt.Println("• Landing page builders")
	fmt.Println("• Content management systems")
	fmt.Println("• Marketing websites")
	fmt.Println("• Documentation sites")
	fmt.Println()
	fmt.Println("Benefits:")
	fmt.Println("• No code changes for new block types")
	fmt.Println("• Flexible content structure")
	fmt.Println("• Reusable across pages")
	fmt.Println("• Version control friendly")
}
