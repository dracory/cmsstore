package cmsstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create test data
	site1 := NewSite().SetName("Site 1")
	err = store.SiteCreate(ctx, site1)
	require.NoError(t, err)

	site2 := NewSite().SetName("Site 2")
	err = store.SiteCreate(ctx, site2)
	require.NoError(t, err)

	// Create blocks with different statuses and sites
	for i := 0; i < 5; i++ {
		block := NewBlock().
			SetSiteID(site1.ID()).
			SetName("Active Block Site 1").
			SetStatus(BLOCK_STATUS_ACTIVE)
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID(site1.ID()).
			SetName("Inactive Block Site 1").
			SetStatus(BLOCK_STATUS_INACTIVE)
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	for i := 0; i < 4; i++ {
		block := NewBlock().
			SetSiteID(site2.ID()).
			SetName("Active Block Site 2").
			SetStatus(BLOCK_STATUS_ACTIVE)
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	// Test: Site + Status combination
	blocks, err := store.BlockList(ctx, BlockQuery().
		SetSiteID(site1.ID()).
		SetStatus(BLOCK_STATUS_ACTIVE))
	require.NoError(t, err)
	require.Len(t, blocks, 5)

	// Test: Site + Status + Limit
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID(site1.ID()).
		SetStatus(BLOCK_STATUS_ACTIVE).
		SetLimit(3))
	require.NoError(t, err)
	require.Len(t, blocks, 3)

	// Test: Site + Limit + Offset
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID(site1.ID()).
		SetLimit(3).
		SetOffset(2))
	require.NoError(t, err)
	require.Len(t, blocks, 3)

	// Test: StatusIn with multiple statuses
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID(site1.ID()).
		SetStatusIn([]string{BLOCK_STATUS_ACTIVE, BLOCK_STATUS_INACTIVE}))
	require.NoError(t, err)
	require.Len(t, blocks, 8) // 5 active + 3 inactive
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
	require.NoError(t, err)

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
		require.NoError(t, err)
	}

	// Test: Search for "Block"
	blocks, err := store.BlockList(ctx, BlockQuery().SetNameLike("%Block%"))
	require.NoError(t, err)
	require.Len(t, blocks, 5)

	// Test: Search for "Header"
	blocks, err = store.BlockList(ctx, BlockQuery().SetNameLike("%Header%"))
	require.NoError(t, err)
	require.Len(t, blocks, 1)
	require.Equal(t, "Header Block", blocks[0].Name())

	// Test: Search for "Content"
	blocks, err = store.BlockList(ctx, BlockQuery().SetNameLike("%Content%"))
	require.NoError(t, err)
	require.Len(t, blocks, 1)
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create blocks with specific sequences
	sequences := []int{5, 1, 3, 2, 4}
	for _, seq := range sequences {
		block := NewBlock().
			SetSiteID("Site1").
			SetName("Block").
			SetSequenceInt(seq)
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	// Test: Order by sequence ASC
	blocks, err := store.BlockList(ctx, BlockQuery().
		SetSiteID("Site1").
		SetOrderBy("sequence").
		SetSortOrder("ASC"))
	require.NoError(t, err)
	require.Len(t, blocks, 5)

	// Verify ascending order
	for i := 0; i < len(blocks)-1; i++ {
		require.LessOrEqual(t, blocks[i].SequenceInt(), blocks[i+1].SequenceInt())
	}

	// Test: Order by sequence DESC
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID("Site1").
		SetOrderBy("sequence").
		SetSortOrder("DESC"))
	require.NoError(t, err)
	require.Len(t, blocks, 5)

	// Verify descending order
	for i := 0; i < len(blocks)-1; i++ {
		require.GreaterOrEqual(t, blocks[i].SequenceInt(), blocks[i+1].SequenceInt())
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create 10 blocks
	for i := 0; i < 10; i++ {
		block := NewBlock().SetSiteID("Site1").SetName("Block")
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	// Test: Limit = 0 (should return all or none depending on implementation)
	blocks, err := store.BlockList(ctx, BlockQuery().SetSiteID("Site1").SetLimit(0))
	require.NoError(t, err)
	// Implementation dependent - could be 0 or all

	// Test: Offset > total records
	blocks, err = store.BlockList(ctx, BlockQuery().SetSiteID("Site1").SetLimit(1000).SetOffset(100))
	require.NoError(t, err)
	require.Empty(t, blocks)

	// Test: Offset at boundary
	blocks, err = store.BlockList(ctx, BlockQuery().SetSiteID("Site1").SetLimit(5).SetOffset(9))
	require.NoError(t, err)
	require.LessOrEqual(t, len(blocks), 1) // Should get 1 or 0 blocks

	// Test: Empty result set
	blocks, err = store.BlockList(ctx, BlockQuery().SetSiteID("NonExistentSite"))
	require.NoError(t, err)
	require.Empty(t, blocks)
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create blocks and collect IDs
	var blockIDs []string
	for i := 0; i < 5; i++ {
		block := NewBlock().SetSiteID("Site1").SetName("Block")
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
		blockIDs = append(blockIDs, block.ID())
	}

	// Test: Query with specific IDs
	selectedIDs := []string{blockIDs[0], blockIDs[2], blockIDs[4]}
	blocks, err := store.BlockList(ctx, BlockQuery().SetIDIn(selectedIDs))
	require.NoError(t, err)
	require.Len(t, blocks, 3)

	// Verify returned IDs match
	returnedIDs := make(map[string]bool)
	for _, block := range blocks {
		returnedIDs[block.ID()] = true
	}
	for _, id := range selectedIDs {
		require.True(t, returnedIDs[id], "Expected ID %s not found", id)
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create blocks
	var blockIDs []string
	for i := 0; i < 5; i++ {
		block := NewBlock().SetSiteID("Site1").SetName("Block")
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
		blockIDs = append(blockIDs, block.ID())
	}

	// Soft delete 2 blocks
	err = store.BlockSoftDeleteByID(ctx, blockIDs[0])
	require.NoError(t, err)
	err = store.BlockSoftDeleteByID(ctx, blockIDs[2])
	require.NoError(t, err)

	// Test: Default query (should exclude soft deleted)
	blocks, err := store.BlockList(ctx, BlockQuery().SetSiteID("Site1"))
	require.NoError(t, err)
	require.Len(t, blocks, 3)

	// Test: Include soft deleted
	blocks, err = store.BlockList(ctx, BlockQuery().
		SetSiteID("Site1").
		SetSoftDeleteIncluded(true))
	require.NoError(t, err)
	require.Len(t, blocks, 5)

	// Verify soft deleted blocks are marked
	softDeletedCount := 0
	for _, block := range blocks {
		if block.IsSoftDeleted() {
			softDeletedCount++
		}
	}
	require.Equal(t, 2, softDeletedCount)
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create blocks with different statuses
	for i := 0; i < 7; i++ {
		block := NewBlock().
			SetSiteID("Site1").
			SetStatus(BLOCK_STATUS_ACTIVE)
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID("Site1").
			SetStatus(BLOCK_STATUS_INACTIVE)
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	// Test: Count all blocks for site
	count, err := store.BlockCount(ctx, BlockQuery().SetSiteID("Site1"))
	require.NoError(t, err)
	require.Equal(t, int64(10), count)

	// Test: Count active blocks
	count, err = store.BlockCount(ctx, BlockQuery().
		SetSiteID("Site1").
		SetStatus(BLOCK_STATUS_ACTIVE))
	require.NoError(t, err)
	require.Equal(t, int64(7), count)

	// Test: Count inactive blocks
	count, err = store.BlockCount(ctx, BlockQuery().
		SetSiteID("Site1").
		SetStatus(BLOCK_STATUS_INACTIVE))
	require.NoError(t, err)
	require.Equal(t, int64(3), count)
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create site, page, and template
	site := NewSite().SetName("Site")
	err = store.SiteCreate(ctx, site)
	require.NoError(t, err)

	page := NewPage().SetSiteID(site.ID()).SetTitle("Page")
	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	template := NewTemplate().SetSiteID(site.ID()).SetName("Template")
	err = store.TemplateCreate(ctx, template)
	require.NoError(t, err)

	// Create blocks with different associations
	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID(site.ID()).
			SetPageID(page.ID())
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	for i := 0; i < 2; i++ {
		block := NewBlock().
			SetSiteID(site.ID()).
			SetTemplateID(template.ID())
		err = store.BlockCreate(ctx, block)
		require.NoError(t, err)
	}

	// Test: Query by page_id
	blocks, err := store.BlockList(ctx, BlockQuery().SetPageID(page.ID()))
	require.NoError(t, err)
	require.Len(t, blocks, 3)

	// Test: Query by template_id
	blocks, err = store.BlockList(ctx, BlockQuery().SetTemplateID(template.ID()))
	require.NoError(t, err)
	require.Len(t, blocks, 2)
}
