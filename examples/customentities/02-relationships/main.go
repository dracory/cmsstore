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

// Example 2: Entity Relationships
// This example demonstrates:
// - Creating entities with relationships (belongs_to, has_many)
// - Linking authors to blog posts
// - Querying relationships

func main() {
	ctx := context.Background()

	// Initialize SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create CMS store with custom entities and relationships enabled
	store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                    db,
		BlockTableName:        "cms_block",
		PageTableName:         "cms_page",
		SiteTableName:         "cms_site",
		TemplateTableName:     "cms_template",
		AutomigrateEnabled:    true,
		CustomEntitiesEnabled: true,
		CustomEntityStoreOptions: cmsstore.CustomEntityStoreOptions{
			EntityTableName:      "cms_custom_entity",
			AttributeTableName:   "cms_custom_attribute",
			RelationshipsEnabled: true,
			AutomigrateEnabled:   true,
		},
		CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
			{
				Type:                 "author",
				TypeLabel:            "Author",
				Group:                "Blog",
				AllowRelationships:   true,
				AllowedRelationTypes: []string{"has_many"},
				Attributes: []cmsstore.CustomAttributeDefinition{
					{Name: "name", Type: "string", Label: "Author Name", Required: true},
					{Name: "email", Type: "string", Label: "Email", Required: true},
					{Name: "bio", Type: "string", Label: "Biography"},
				},
			},
			{
				Type:                 "post",
				TypeLabel:            "Blog Post",
				Group:                "Blog",
				AllowRelationships:   true,
				AllowedRelationTypes: []string{"belongs_to"},
				Attributes: []cmsstore.CustomAttributeDefinition{
					{Name: "title", Type: "string", Label: "Post Title", Required: true},
					{Name: "content", Type: "string", Label: "Content", Required: true},
					{Name: "published", Type: "bool", Label: "Published"},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	customStore := store.CustomEntityStore()

	// Example 1: Create an author
	fmt.Println("=== Creating Author ===")
	authorID, err := customStore.Create(ctx, "author", map[string]interface{}{
		"name":  "Jane Doe",
		"email": "jane@example.com",
		"bio":   "Technical writer and software developer",
	}, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created author with ID: %s\n\n", authorID)

	// Example 2: Create blog posts with relationship to author
	fmt.Println("=== Creating Blog Posts with Author Relationship ===")

	post1ID, err := customStore.Create(ctx, "post", map[string]interface{}{
		"title":     "Introduction to Go",
		"content":   "Go is a statically typed, compiled programming language...",
		"published": true,
	}, []cmsstore.RelationshipDefinition{
		{
			Type:     "belongs_to",
			TargetID: authorID,
			Metadata: map[string]interface{}{
				"role": "author",
			},
		},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created post 1 with ID: %s\n", post1ID)

	post2ID, err := customStore.Create(ctx, "post", map[string]interface{}{
		"title":     "Advanced Go Patterns",
		"content":   "This article covers advanced design patterns in Go...",
		"published": true,
	}, []cmsstore.RelationshipDefinition{
		{
			Type:     "belongs_to",
			TargetID: authorID,
			Metadata: map[string]interface{}{
				"role": "author",
			},
		},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created post 2 with ID: %s\n\n", post2ID)

	// Example 3: Query relationships for a post
	fmt.Println("=== Querying Post Relationships ===")
	post1Relationships, err := customStore.GetRelationships(ctx, post1ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Post 1 has %d relationship(s)\n", len(post1Relationships))
	for _, rel := range post1Relationships {
		fmt.Printf("- Relationship Type: %s\n", rel.RelationshipType())
		fmt.Printf("  Related Entity ID: %s\n", rel.RelatedEntityID())
		fmt.Printf("  Metadata: %s\n", rel.Metadata())
	}
	fmt.Println()

	// Example 4: Find all posts by author (using relationship query)
	fmt.Println("=== Finding All Posts by Author ===")

	// Query relationships where the author is the related entity
	authorRelationships, err := customStore.Inner().RelationshipList(ctx, entitystore.RelationshipQueryOptions{
		RelatedEntityID:  authorID,
		RelationshipType: "belongs_to",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Author has %d post(s):\n", len(authorRelationships))
	for i, rel := range authorRelationships {
		// Get the post entity
		post, err := customStore.FindByID(ctx, rel.EntityID())
		if err != nil {
			continue
		}
		fmt.Printf("%d. Post ID: %s (Type: %s)\n", i+1, post.ID(), post.EntityType())
	}
	fmt.Println()

	// Example 5: Create a comment entity with relationship to post
	fmt.Println("=== Adding Comment Entity Type ===")

	// Register comment entity type
	err = customStore.RegisterEntityType(cmsstore.CustomEntityDefinition{
		Type:                 "comment",
		TypeLabel:            "Comment",
		Group:                "Blog",
		AllowRelationships:   true,
		AllowedRelationTypes: []string{"belongs_to"},
		Attributes: []cmsstore.CustomAttributeDefinition{
			{Name: "author_name", Type: "string", Label: "Commenter Name", Required: true},
			{Name: "content", Type: "string", Label: "Comment Content", Required: true},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create comments for post 1
	comment1ID, err := customStore.Create(ctx, "comment", map[string]interface{}{
		"author_name": "John Smith",
		"content":     "Great article! Very helpful.",
	}, []cmsstore.RelationshipDefinition{
		{
			Type:     "belongs_to",
			TargetID: post1ID,
			Metadata: map[string]interface{}{
				"relationship": "comment_on_post",
			},
		},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created comment with ID: %s\n", comment1ID)

	comment2ID, err := customStore.Create(ctx, "comment", map[string]interface{}{
		"author_name": "Alice Johnson",
		"content":     "Thanks for sharing this!",
	}, []cmsstore.RelationshipDefinition{
		{
			Type:     "belongs_to",
			TargetID: post1ID,
			Metadata: map[string]interface{}{
				"relationship": "comment_on_post",
			},
		},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created comment with ID: %s\n\n", comment2ID)

	// Example 6: Query all comments for a post
	fmt.Println("=== Finding All Comments for Post ===")

	postComments, err := customStore.Inner().RelationshipList(ctx, entitystore.RelationshipQueryOptions{
		RelatedEntityID:  post1ID,
		RelationshipType: "belongs_to",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Post 1 has %d comment(s):\n", len(postComments))
	for i, rel := range postComments {
		comment, err := customStore.FindByID(ctx, rel.EntityID())
		if err != nil {
			continue
		}
		// Only show comments (not the author relationship)
		if comment.EntityType() == "comment" {
			fmt.Printf("%d. Comment ID: %s\n", i+1, comment.ID())
		}
	}
	fmt.Println()

	// Example 7: Summary
	fmt.Println("=== Summary ===")
	fmt.Println("Relationship structure:")
	fmt.Println("- 1 Author")
	fmt.Println("  - has 2 Posts (via belongs_to relationships)")
	fmt.Println("    - Post 1 has 2 Comments (via belongs_to relationships)")
	fmt.Println("    - Post 2 has 0 Comments")
	fmt.Println("\nRelationships enable:")
	fmt.Println("✓ Linking entities together")
	fmt.Println("✓ Querying related entities")
	fmt.Println("✓ Storing metadata on relationships")
	fmt.Println("✓ Building complex data structures")
}
