package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dracory/cmsstore"
	_ "modernc.org/sqlite"
)

// Basic Example: Getting Started with CMS Store
// This is the simplest possible example to get started with the CMS.
// Demonstrates basic CMS features: Sites, Pages, and Blocks.
// Run: go run main.go

func main() {
	ctx := context.Background()

	// Step 1: Initialize database
	fmt.Println("=== Initializing Database ===")
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("✓ Database initialized")

	// Step 2: Create CMS store
	fmt.Println("\n=== Creating CMS Store ===")
	store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                 db,
		BlockTableName:     "cms_block",
		PageTableName:      "cms_page",
		SiteTableName:      "cms_site",
		TemplateTableName:  "cms_template",
		AutomigrateEnabled: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ CMS store created")
	fmt.Println("✓ Tables auto-migrated")

	// Step 3: Create a site
	fmt.Println("\n=== Creating a Site ===")
	site := cmsstore.NewSite()
	site.SetName("My Website")
	site.SetHandle("my-website")
	site.SetStatus("active")

	err = store.SiteCreate(ctx, site)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Created site: %s (ID: %s)\n", site.Name(), site.ID())

	// Step 4: Create pages
	fmt.Println("\n=== Creating Pages ===")

	homePage := cmsstore.NewPage()
	homePage.SetSiteID(site.ID())
	homePage.SetTitle("Home")
	homePage.SetAlias("/")
	homePage.SetHandle("home")
	homePage.SetStatus("active")

	err = store.PageCreate(ctx, homePage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Created page: %s (Alias: %s)\n", homePage.Title(), homePage.Alias())

	aboutPage := cmsstore.NewPage()
	aboutPage.SetSiteID(site.ID())
	aboutPage.SetTitle("About Us")
	aboutPage.SetAlias("/about")
	aboutPage.SetHandle("about")
	aboutPage.SetStatus("active")

	err = store.PageCreate(ctx, aboutPage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Created page: %s (Alias: %s)\n", aboutPage.Title(), aboutPage.Alias())

	contactPage := cmsstore.NewPage()
	contactPage.SetSiteID(site.ID())
	contactPage.SetTitle("Contact")
	contactPage.SetAlias("/contact")
	contactPage.SetHandle("contact")
	contactPage.SetStatus("active")

	err = store.PageCreate(ctx, contactPage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Created page: %s (Alias: %s)\n", contactPage.Title(), contactPage.Alias())

	// Step 5: Create blocks
	fmt.Println("\n=== Creating Content Blocks ===")

	headerBlock := cmsstore.NewBlock()
	headerBlock.SetSiteID(site.ID())
	headerBlock.SetName("Header")
	headerBlock.SetHandle("header")
	headerBlock.SetType("html")
	headerBlock.SetContent("<header><h1>Welcome to My Website</h1></header>")
	headerBlock.SetStatus("active")

	err = store.BlockCreate(ctx, headerBlock)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Created block: %s (Type: %s)\n", headerBlock.Name(), headerBlock.Type())

	footerBlock := cmsstore.NewBlock()
	footerBlock.SetSiteID(site.ID())
	footerBlock.SetName("Footer")
	footerBlock.SetHandle("footer")
	footerBlock.SetType("html")
	footerBlock.SetContent("<footer><p>&copy; 2026 My Website</p></footer>")
	footerBlock.SetStatus("active")

	err = store.BlockCreate(ctx, footerBlock)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Created block: %s (Type: %s)\n", footerBlock.Name(), footerBlock.Type())

	// Step 6: List all pages
	fmt.Println("\n=== Listing All Pages ===")
	pages, err := store.PageList(ctx, cmsstore.PageQuery().SetSiteID(site.ID()))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total pages: %d\n", len(pages))
	for i, p := range pages {
		fmt.Printf("%d. %s - %s (Status: %s)\n", i+1, p.Title(), p.Alias(), p.Status())
	}

	// Step 7: List all blocks
	fmt.Println("\n=== Listing All Blocks ===")
	blocks, err := store.BlockList(ctx, cmsstore.BlockQuery().SetSiteID(site.ID()))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total blocks: %d\n", len(blocks))
	for i, b := range blocks {
		fmt.Printf("%d. %s (Handle: %s, Type: %s)\n", i+1, b.Name(), b.Handle(), b.Type())
	}

	// Step 8: Update a page
	fmt.Println("\n=== Updating a Page ===")
	homePage.SetTitle("Home - Welcome!")
	err = store.PageUpdate(ctx, homePage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Updated page title to: %s\n", homePage.Title())

	// Step 9: Find a page by handle
	fmt.Println("\n=== Finding Page by Handle ===")
	foundPage, err := store.PageFindByHandle(ctx, "about")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✓ Found page: %s (ID: %s)\n", foundPage.Title(), foundPage.ID())

	// Step 10: Count pages
	fmt.Println("\n=== Counting Pages ===")
	pageCount, err := store.PageCount(ctx, cmsstore.PageQuery().SetSiteID(site.ID()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total page count: %d\n", pageCount)

	// Summary
	fmt.Println("\n=== Summary ===")
	fmt.Println("This example demonstrated:")
	fmt.Println("✓ Creating a CMS store")
	fmt.Println("✓ Creating a site")
	fmt.Println("✓ Creating pages with aliases")
	fmt.Println("✓ Creating content blocks")
	fmt.Println("✓ Listing pages and blocks")
	fmt.Println("✓ Updating pages")
	fmt.Println("✓ Finding pages by handle")
	fmt.Println("✓ Counting pages")
	fmt.Println()
	fmt.Println("Core CMS Features:")
	fmt.Println("• Sites - Multi-site support")
	fmt.Println("• Pages - Content pages with aliases")
	fmt.Println("• Blocks - Reusable content blocks")
	fmt.Println("• Templates - Page templates (not shown)")
	fmt.Println()
	fmt.Println("Next Steps:")
	fmt.Println("→ See ../customentities for custom entity examples")
	fmt.Println("→ See ../custom-block-types for extending with custom blocks")
}
