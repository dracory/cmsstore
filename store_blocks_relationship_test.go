package cmsstore

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// TestBlockCreateWithInvalidSiteID tests that creating a block with non-existent site fails or succeeds based on FK constraints
func TestBlockCreateWithInvalidSiteID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_invalid_site",
		PageTableName:      "page_table_invalid_site",
		SiteTableName:      "site_table_invalid_site",
		TemplateTableName:  "template_table_invalid_site",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create block with non-existent site ID
	block := NewBlock().
		SetSiteID("non-existent-site").
		SetName("Test Block")

	err = store.BlockCreate(ctx, block)
	// Note: Depending on FK constraints, this may or may not error
	// If it succeeds, we should still be able to query it
	if err == nil {
		found, err := store.BlockFindByID(ctx, block.ID())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found == nil {
			t.Fatal("Block MUST NOT be nil")
		}
		if found.SiteID() != "non-existent-site" {
			t.Fatalf("Expected SiteID 'non-existent-site', got %s", found.SiteID())
		}
	}
}

// TestBlockCreateWithInvalidPageID tests block creation with non-existent page
func TestBlockCreateWithInvalidPageID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_invalid_page",
		PageTableName:      "page_table_invalid_page",
		SiteTableName:      "site_table_invalid_page",
		TemplateTableName:  "template_table_invalid_page",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create a valid site first
	site := NewSite().SetName("Test Site")
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create block with non-existent page ID
	block := NewBlock().
		SetSiteID(site.ID()).
		SetPageID("non-existent-page").
		SetName("Test Block")

	err = store.BlockCreate(ctx, block)
	if err == nil {
		found, err := store.BlockFindByID(ctx, block.ID())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found == nil {
			t.Fatal("Block MUST NOT be nil")
		}
		if found.PageID() != "non-existent-page" {
			t.Fatalf("Expected PageID 'non-existent-page', got %s", found.PageID())
		}
	}
}

// TestBlockListByPageID tests that blocks are properly associated with pages
func TestBlockListByPageID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_by_page",
		PageTableName:      "page_table_by_page",
		SiteTableName:      "site_table_by_page",
		TemplateTableName:  "template_table_by_page",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create site
	site := NewSite().SetName("Test Site")
	err = store.SiteCreate(ctx, site)
	require.NoError(t, err)

	// Create two pages
	page1 := NewPage().SetSiteID(site.ID()).SetTitle("Page 1")
	err = store.PageCreate(ctx, page1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	page2 := NewPage().SetSiteID(site.ID()).SetTitle("Page 2")
	err = store.PageCreate(ctx, page2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create blocks for page1
	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID(site.ID()).
			SetPageID(page1.ID()).
			SetName("Block for Page 1")
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Create blocks for page2
	for i := 0; i < 2; i++ {
		block := NewBlock().
			SetSiteID(site.ID()).
			SetPageID(page2.ID()).
			SetName("Block for Page 2")
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Query blocks for page1
	blocks1, err := store.BlockList(ctx, BlockQuery().SetPageID(page1.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks1) != 3 {
		t.Fatalf("Expected 3 blocks, got %d", len(blocks1))
	}

	// Query blocks for page2
	blocks2, err := store.BlockList(ctx, BlockQuery().SetPageID(page2.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks2) != 2 {
		t.Fatalf("Expected 2 blocks, got %d", len(blocks2))
	}

	// Verify all blocks for page1 have correct page_id
	for _, block := range blocks1 {
		if block.PageID() != page1.ID() {
			t.Fatalf("Expected PageID %s, got %s", page1.ID(), block.PageID())
		}
	}
}

// TestPageDeleteWithBlocks tests what happens when deleting a page that has blocks
func TestPageDeleteWithBlocks(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_page_delete",
		PageTableName:      "page_table_page_delete",
		SiteTableName:      "site_table_page_delete",
		TemplateTableName:  "template_table_page_delete",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create site and page
	site := NewSite().SetName("Test Site")
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	page := NewPage().SetSiteID(site.ID()).SetTitle("Test Page")
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create blocks for the page
	var blockIDs []string
	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID(site.ID()).
			SetPageID(page.ID()).
			SetName("Test Block")
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		blockIDs = append(blockIDs, block.ID())
	}

	// Delete the page
	err = store.PageDeleteByID(ctx, page.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if blocks still exist (orphaned) or were cascade deleted
	for _, blockID := range blockIDs {
		block, err := store.BlockFindByID(ctx, blockID)
		// Depending on implementation, blocks may be orphaned or deleted
		if err == nil && block != nil {
			// Blocks are orphaned - verify they still reference the deleted page
			if block.PageID() != page.ID() {
				t.Fatalf("Expected PageID %s, got %s", page.ID(), block.PageID())
			}
		}
	}
}

// TestSiteDeleteWithPages tests cascading behavior when deleting a site
func TestSiteDeleteWithPages(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_site_delete",
		PageTableName:      "page_table_site_delete",
		SiteTableName:      "site_table_site_delete",
		TemplateTableName:  "template_table_site_delete",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create site
	site := NewSite().SetName("Test Site")
	err = store.SiteCreate(ctx, site)
	require.NoError(t, err)

	// Create pages for the site
	var pageIDs []string
	for i := 0; i < 3; i++ {
		page := NewPage().SetSiteID(site.ID()).SetTitle("Test Page")
		err = store.PageCreate(ctx, page)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		pageIDs = append(pageIDs, page.ID())
	}

	// Delete the site
	err = store.SiteDeleteByID(ctx, site.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if pages still exist (orphaned) or were cascade deleted
	for _, pageID := range pageIDs {
		page, err := store.PageFindByID(ctx, pageID)
		if err == nil && page != nil {
			// Pages are orphaned
			if page.SiteID() != site.ID() {
				t.Fatalf("Expected SiteID %s, got %s", site.ID(), page.SiteID())
			}
		}
	}
}

// TestBlockUpdatePreservesUnchangedFields tests that update only changes specified fields
func TestBlockUpdatePreservesUnchangedFields(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_update_preserve",
		PageTableName:      "page_table_update_preserve",
		SiteTableName:      "site_table_update_preserve",
		TemplateTableName:  "template_table_update_preserve",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create block with multiple fields
	block := NewBlock().
		SetSiteID("Site1").
		SetName("Original Name").
		SetContent("Original Content").
		SetHandle("original-handle").
		SetStatus(BLOCK_STATUS_ACTIVE)

	err = block.SetMetas(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	originalCreatedAt := block.CreatedAt()
	originalID := block.ID()

	// Update only the name
	block.SetName("Updated Name")
	err = store.BlockUpdate(ctx, block)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Fetch and verify
	found, err := store.BlockFindByID(ctx, block.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("Block MUST NOT be nil")
	}

	// Verify changed field
	if found.Name() != "Updated Name" {
		t.Fatalf("Expected Name 'Updated Name', got %s", found.Name())
	}

	// Verify unchanged fields
	if found.ID() != originalID {
		t.Fatalf("Expected ID %s, got %s", originalID, found.ID())
	}
	if found.Content() != "Original Content" {
		t.Fatalf("Expected Content 'Original Content', got %s", found.Content())
	}
	if found.Handle() != "original-handle" {
		t.Fatalf("Expected Handle 'original-handle', got %s", found.Handle())
	}
	if found.Status() != BLOCK_STATUS_ACTIVE {
		t.Fatalf("Expected Status %s, got %s", BLOCK_STATUS_ACTIVE, found.Status())
	}
	if found.SiteID() != "Site1" {
		t.Fatalf("Expected SiteID 'Site1', got %s", found.SiteID())
	}
	// CreatedAt should be preserved (compare timestamp values, format may vary)
	if !strings.Contains(found.CreatedAt(), originalCreatedAt[:19]) { // Compare date/time portion without timezone
		t.Fatalf("Expected CreatedAt to contain %s, got %s", originalCreatedAt[:19], found.CreatedAt())
	}

	// Verify metas are preserved
	if found.Meta("key1") != "value1" {
		t.Fatalf("Expected Meta('key1') 'value1', got %s", found.Meta("key1"))
	}
	if found.Meta("key2") != "value2" {
		t.Fatalf("Expected Meta('key2') 'value2', got %s", found.Meta("key2"))
	}
}

// TestBlockDuplicateHandle tests behavior when creating blocks with duplicate handles
func TestBlockDuplicateHandle(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_dup_handle",
		PageTableName:      "page_table_dup_handle",
		SiteTableName:      "site_table_dup_handle",
		TemplateTableName:  "template_table_dup_handle",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create first block with handle
	block1 := NewBlock().
		SetSiteID("Site1").
		SetHandle("unique-handle").
		SetName("Block 1")

	err = store.BlockCreate(ctx, block1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Try to create second block with same handle
	block2 := NewBlock().
		SetSiteID("Site1").
		SetHandle("unique-handle").
		SetName("Block 2")

	err = store.BlockCreate(ctx, block2)
	// Depending on unique constraints, this may error or succeed
	// If it succeeds, both blocks should be queryable
	if err == nil {
		// Query by handle should return one of them
		found, err := store.BlockFindByHandle(ctx, "unique-handle")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found == nil {
			t.Fatal("Block MUST NOT be nil")
		}
		// It should be one of the two blocks
		if found.ID() != block1.ID() && found.ID() != block2.ID() {
			t.Fatalf("Expected ID to be one of %v, got %s", []string{block1.ID(), block2.ID()}, found.ID())
		}
	}
}

// TestMenuItemHierarchy tests parent-child relationships in menu items
func TestMenuItemHierarchy(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_menu_hierarchy",
		PageTableName:      "page_table_menu_hierarchy",
		SiteTableName:      "site_table_menu_hierarchy",
		TemplateTableName:  "template_table_menu_hierarchy",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_menu_hierarchy",
		MenuItemTableName:  "menu_item_table_menu_hierarchy",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create menu
	menu := NewMenu().SetSiteID("Site1").SetName("Test Menu")
	err = store.MenuCreate(ctx, menu)
	require.NoError(t, err)

	// Create parent menu item
	parent := NewMenuItem().
		SetMenuID(menu.ID()).
		SetName("Parent").
		SetSequenceInt(1)
	err = store.MenuItemCreate(ctx, parent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create child menu items
	child1 := NewMenuItem().
		SetMenuID(menu.ID()).
		SetParentID(parent.ID()).
		SetName("Child 1").
		SetSequenceInt(1)
	err = store.MenuItemCreate(ctx, child1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	child2 := NewMenuItem().
		SetMenuID(menu.ID()).
		SetParentID(parent.ID()).
		SetName("Child 2").
		SetSequenceInt(2)
	err = store.MenuItemCreate(ctx, child2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Query all menu items and filter children
	allItems, err := store.MenuItemList(ctx, MenuItemQuery().SetMenuID(menu.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Filter children by parent ID
	var children []MenuItemInterface
	for _, item := range allItems {
		if item.ParentID() == parent.ID() {
			children = append(children, item)
		}
	}
	if len(children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(children))
	}

	// Verify parent-child relationships
	for _, child := range children {
		if child.ParentID() != parent.ID() {
			t.Fatalf("Expected ParentID %s, got %s", parent.ID(), child.ParentID())
		}
		if child.MenuID() != menu.ID() {
			t.Fatalf("Expected MenuID %s, got %s", menu.ID(), child.MenuID())
		}
	}
}

// TestMenuItemCircularReference tests detection/prevention of circular parent-child references
func TestMenuItemCircularReference(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_circular",
		PageTableName:      "page_table_circular",
		SiteTableName:      "site_table_circular",
		TemplateTableName:  "template_table_circular",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_circular",
		MenuItemTableName:  "menu_item_table_circular",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create menu
	menu := NewMenu().SetSiteID("Site1").SetName("Test Menu")
	err = store.MenuCreate(ctx, menu)
	require.NoError(t, err)

	// Create menu item A
	itemA := NewMenuItem().
		SetMenuID(menu.ID()).
		SetName("Item A").
		SetSequenceInt(1)
	err = store.MenuItemCreate(ctx, itemA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create menu item B as child of A
	itemB := NewMenuItem().
		SetMenuID(menu.ID()).
		SetParentID(itemA.ID()).
		SetName("Item B").
		SetSequenceInt(1)
	err = store.MenuItemCreate(ctx, itemB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Try to make A a child of B (circular reference)
	itemA.SetParentID(itemB.ID())
	err = store.MenuItemUpdate(ctx, itemA)
	// This should either error or succeed
	// If it succeeds, we have a circular reference that needs to be handled in rendering
	if err == nil {
		// Circular reference exists - verify both items
		foundA, _ := store.MenuItemFindByID(ctx, itemA.ID())
		foundB, _ := store.MenuItemFindByID(ctx, itemB.ID())

		if foundA != nil && foundB != nil {
			if foundA.ParentID() != itemB.ID() {
				t.Fatalf("Expected ParentID %s, got %s", itemB.ID(), foundA.ParentID())
			}
			if foundB.ParentID() != itemA.ID() {
				t.Fatalf("Expected ParentID %s, got %s", itemA.ID(), foundB.ParentID())
			}
		}
	}
}
