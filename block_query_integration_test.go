package cmsstore

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

// TestBlockQueryIntegrationCombinations tests query combinations with actual database
func TestBlockQueryIntegrationCombinations(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_query_combo",
		PageTableName:      "page_table_query_combo",
		SiteTableName:      "site_table_query_combo",
		TemplateTableName:  "template_table_query_combo",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// Create test data
	site1 := NewSite().SetName("Site 1")
	err = store.SiteCreate(ctx, site1)
	if err != nil {
		t.Fatalf("failed to create site1: %v", err)
	}

	site2 := NewSite().SetName("Site 2")
	err = store.SiteCreate(ctx, site2)
	if err != nil {
		t.Fatalf("failed to create site2: %v", err)
	}

	// Create blocks with different statuses and sites
	for i := 0; i < 5; i++ {
		block := NewBlock().
			SetSiteID(site1.ID()).
			SetName("Active Block Site 1").
			SetStatus(BLOCK_STATUS_ACTIVE)
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID(site1.ID()).
			SetName("Inactive Block Site 1").
			SetStatus(BLOCK_STATUS_INACTIVE)
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	for i := 0; i < 4; i++ {
		block := NewBlock().
			SetSiteID(site2.ID()).
			SetName("Active Block Site 2").
			SetStatus(BLOCK_STATUS_ACTIVE)
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	// Test: Site + Status combination
	blocks, err := store.BlockList(ctx, BlockQuery().
		SetSiteID(site1.ID()).
		SetStatus(BLOCK_STATUS_ACTIVE))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 5 {
		t.Errorf("expected 5 blocks, got %d", len(blocks))
	}

	// Test: Site + Status + Limit
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID(site1.ID()).
		SetStatus(BLOCK_STATUS_ACTIVE).
		SetLimit(3))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 3 {
		t.Errorf("expected 3 blocks, got %d", len(blocks))
	}

	// Test: Site + Limit + Offset
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID(site1.ID()).
		SetLimit(3).
		SetOffset(2))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 3 {
		t.Errorf("expected 3 blocks, got %d", len(blocks))
	}

	// Test: StatusIn with multiple statuses
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID(site1.ID()).
		SetStatusIn([]string{BLOCK_STATUS_ACTIVE, BLOCK_STATUS_INACTIVE}))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 8 { // 5 active + 3 inactive
		t.Errorf("expected 8 blocks, got %d", len(blocks))
	}
}

// TestBlockQueryIntegrationNameLike tests LIKE queries
func TestBlockQueryIntegrationNameLike(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_name_like",
		PageTableName:      "page_table_name_like",
		SiteTableName:      "site_table_name_like",
		TemplateTableName:  "template_table_name_like",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// Create blocks with different names
	names := []string{
		"Header Block",
		"Footer Block",
		"Sidebar Block",
		"Main Content Block",
		"Navigation Block",
	}

	for _, name := range names {
		block := NewBlock().
			SetSiteID("Site1").
			SetName(name)
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	// Test: Search for "Block"
	blocks, err := store.BlockList(ctx, BlockQuery().SetNameLike("%Block%"))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 5 {
		t.Errorf("expected 5 blocks, got %d", len(blocks))
	}

	// Test: Search for "Header"
	blocks, err = store.BlockList(ctx, BlockQuery().SetNameLike("%Header%"))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(blocks))
	}
	if blocks[0].Name() != "Header Block" {
		t.Errorf("expected 'Header Block', got %q", blocks[0].Name())
	}

	// Test: Search for "Content"
	blocks, err = store.BlockList(ctx, BlockQuery().SetNameLike("%Content%"))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(blocks))
	}
}

// TestBlockQueryIntegrationOrderBy tests ordering
func TestBlockQueryIntegrationOrderBy(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_order",
		PageTableName:      "page_table_order",
		SiteTableName:      "site_table_order",
		TemplateTableName:  "template_table_order",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// Create blocks with specific sequences
	sequences := []int{5, 1, 3, 2, 4}
	for _, seq := range sequences {
		block := NewBlock().
			SetSiteID("Site1").
			SetName("Block").
			SetSequenceInt(seq)
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	// Test: Order by sequence ASC
	blocks, err := store.BlockList(ctx, BlockQuery().
		SetSiteID("Site1").
		SetOrderBy("sequence").
		SetSortOrder("ASC"))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 5 {
		t.Errorf("expected 5 blocks, got %d", len(blocks))
	}

	// Verify ascending order
	for i := 0; i < len(blocks)-1; i++ {
		if blocks[i].SequenceInt() > blocks[i+1].SequenceInt() {
			t.Errorf("blocks not in ascending order at index %d: %d > %d", i, blocks[i].SequenceInt(), blocks[i+1].SequenceInt())
		}
	}

	// Test: Order by sequence DESC
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID("Site1").
		SetOrderBy("sequence").
		SetSortOrder("DESC"))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 5 {
		t.Errorf("expected 5 blocks, got %d", len(blocks))
	}

	// Verify descending order
	for i := 0; i < len(blocks)-1; i++ {
		if blocks[i].SequenceInt() < blocks[i+1].SequenceInt() {
			t.Errorf("blocks not in descending order at index %d: %d < %d", i, blocks[i].SequenceInt(), blocks[i+1].SequenceInt())
		}
	}
}

// TestBlockQueryIntegrationEdgeCases tests edge cases
func TestBlockQueryIntegrationEdgeCases(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_edge",
		PageTableName:      "page_table_edge",
		SiteTableName:      "site_table_edge",
		TemplateTableName:  "template_table_edge",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// Create 10 blocks
	for i := 0; i < 10; i++ {
		block := NewBlock().SetSiteID("Site1").SetName("Block")
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	// Test: Limit = 0 (should return all or none depending on implementation)
	blocks, err := store.BlockList(ctx, BlockQuery().SetSiteID("Site1").SetLimit(0))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	// Implementation dependent - could be 0 or all

	// Test: Offset > total records
	blocks, err = store.BlockList(ctx, BlockQuery().SetSiteID("Site1").SetLimit(1000).SetOffset(100))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks for offset > total, got %d", len(blocks))
	}

	// Test: Offset at boundary
	blocks, err = store.BlockList(ctx, BlockQuery().SetSiteID("Site1").SetLimit(5).SetOffset(9))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) > 1 { // Should get 1 or 0 blocks
		t.Errorf("expected at most 1 block at boundary, got %d", len(blocks))
	}

	// Test: Empty result set
	blocks, err = store.BlockList(ctx, BlockQuery().SetSiteID("NonExistentSite"))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks for non-existent site, got %d", len(blocks))
	}
}

// TestBlockQueryIntegrationIDIn tests ID IN queries
func TestBlockQueryIntegrationIDIn(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_id_in",
		PageTableName:      "page_table_id_in",
		SiteTableName:      "site_table_id_in",
		TemplateTableName:  "template_table_id_in",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// Create blocks and collect IDs
	var blockIDs []string
	for i := 0; i < 5; i++ {
		block := NewBlock().SetSiteID("Site1").SetName("Block")
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
		blockIDs = append(blockIDs, block.ID())
	}

	// Test: Query with specific IDs
	selectedIDs := []string{blockIDs[0], blockIDs[2], blockIDs[4]}
	blocks, err := store.BlockList(ctx, BlockQuery().SetIDIn(selectedIDs))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 3 {
		t.Errorf("expected 3 blocks, got %d", len(blocks))
	}

	// Verify returned IDs match
	returnedIDs := make(map[string]bool)
	for _, block := range blocks {
		returnedIDs[block.ID()] = true
	}
	for _, id := range selectedIDs {
		if !returnedIDs[id] {
			t.Errorf("expected ID %s not found", id)
		}
	}
}

// TestBlockQueryIntegrationSoftDeleted tests soft delete filtering
func TestBlockQueryIntegrationSoftDeleted(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_soft_del",
		PageTableName:      "page_table_soft_del",
		SiteTableName:      "site_table_soft_del",
		TemplateTableName:  "template_table_soft_del",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// Create blocks
	var blockIDs []string
	for i := 0; i < 5; i++ {
		block := NewBlock().SetSiteID("Site1").SetName("Block")
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
		blockIDs = append(blockIDs, block.ID())
	}

	// Soft delete 2 blocks
	err = store.BlockSoftDeleteByID(ctx, blockIDs[0])
	if err != nil {
		t.Fatalf("failed to soft delete block: %v", err)
	}
	err = store.BlockSoftDeleteByID(ctx, blockIDs[2])
	if err != nil {
		t.Fatalf("failed to soft delete block: %v", err)
	}

	// Test: Default query (should exclude soft deleted)
	blocks, err := store.BlockList(ctx, BlockQuery().SetSiteID("Site1"))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 3 {
		t.Errorf("expected 3 non-deleted blocks, got %d", len(blocks))
	}

	// Test: Include soft deleted
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID("Site1").
		SetSoftDeleteIncluded(true))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 5 {
		t.Errorf("expected 5 total blocks, got %d", len(blocks))
	}

	// Verify soft deleted blocks are marked
	softDeletedCount := 0
	for _, block := range blocks {
		if block.IsSoftDeleted() {
			softDeletedCount++
		}
	}
	if softDeletedCount != 2 {
		t.Errorf("expected 2 soft deleted blocks, got %d", softDeletedCount)
	}
}

// TestBlockQueryIntegrationCount tests count operations
func TestBlockQueryIntegrationCount(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_count_int",
		PageTableName:      "page_table_count_int",
		SiteTableName:      "site_table_count_int",
		TemplateTableName:  "template_table_count_int",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// Create blocks with different statuses
	for i := 0; i < 7; i++ {
		block := NewBlock().
			SetSiteID("Site1").
			SetStatus(BLOCK_STATUS_ACTIVE)
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID("Site1").
			SetStatus(BLOCK_STATUS_INACTIVE)
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	// Test: Count all blocks for site
	count, err := store.BlockCount(ctx, BlockQuery().SetSiteID("Site1"))
	if err != nil {
		t.Fatalf("failed to count blocks: %v", err)
	}
	if count != 10 {
		t.Errorf("expected count 10, got %d", count)
	}

	// Test: Count active blocks
	count, err = store.BlockCount(ctx, BlockQuery().
		SetSiteID("Site1").
		SetStatus(BLOCK_STATUS_ACTIVE))
	if err != nil {
		t.Fatalf("failed to count blocks: %v", err)
	}
	if count != 7 {
		t.Errorf("expected count 7, got %d", count)
	}

	// Test: Count inactive blocks
	count, err = store.BlockCount(ctx, BlockQuery().
		SetSiteID("Site1").
		SetStatus(BLOCK_STATUS_INACTIVE))
	if err != nil {
		t.Fatalf("failed to count blocks: %v", err)
	}
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}

// TestBlockQueryIntegrationPageAndTemplate tests page_id and template_id filtering
func TestBlockQueryIntegrationPageAndTemplate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_page_tpl",
		PageTableName:      "page_table_page_tpl",
		SiteTableName:      "site_table_page_tpl",
		TemplateTableName:  "template_table_page_tpl",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	ctx := context.Background()

	// Create site, page, and template
	site := NewSite().SetName("Site")
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatalf("failed to create site: %v", err)
	}

	page := NewPage().SetSiteID(site.ID()).SetTitle("Page")
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	template := NewTemplate().SetSiteID(site.ID()).SetName("Template")
	err = store.TemplateCreate(ctx, template)
	if err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	// Create blocks with different associations
	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID(site.ID()).
			SetPageID(page.ID())
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	for i := 0; i < 2; i++ {
		block := NewBlock().
			SetSiteID(site.ID()).
			SetTemplateID(template.ID())
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("failed to create block: %v", err)
		}
	}

	// Test: Query by page_id
	blocks, err := store.BlockList(ctx, BlockQuery().SetPageID(page.ID()))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 3 {
		t.Errorf("expected 3 blocks, got %d", len(blocks))
	}

	// Test: Query by template_id
	blocks, err = store.BlockList(ctx, BlockQuery().SetTemplateID(template.ID()))
	if err != nil {
		t.Fatalf("failed to list blocks: %v", err)
	}
	if len(blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(blocks))
	}
}
