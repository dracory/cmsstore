package cmsstore

import (
	"context"
	"sync"
	"testing"

	_ "modernc.org/sqlite"
)

// TestConcurrentBlockCreate tests concurrent block creation
func TestConcurrentBlockCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create site first (this also verifies tables exist after migration)
	site := NewSite().SetName("Test Site")
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Concurrently create blocks
	const numGoroutines = 10
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)
	blockIDs := make(chan string, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			block := NewBlock().
				SetSiteID(site.ID()).
				SetName("Concurrent Block")

			if err := store.BlockCreate(ctx, block); err != nil {
				errors <- err
				return
			}
			blockIDs <- block.ID()
		}(i)
	}

	wg.Wait()
	close(errors)
	close(blockIDs)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent create error: %v", err)
	}

	// Verify all blocks were created
	var ids []string
	for id := range blockIDs {
		ids = append(ids, id)
	}
	if len(ids) != numGoroutines {
		t.Fatalf("Expected %d blocks, got %d", numGoroutines, len(ids))
	}

	// Verify all IDs are unique
	uniqueIDs := make(map[string]bool)
	for _, id := range ids {
		if uniqueIDs[id] {
			t.Fatalf("Duplicate ID found: %s", id)
		}
		uniqueIDs[id] = true
	}
}

// TestConcurrentBlockUpdate tests concurrent updates to the same block
func TestConcurrentBlockUpdate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create a block
	block := NewBlock().
		SetSiteID("Site1").
		SetName("Original Name").
		SetContent("Original Content")

	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	blockID := block.ID()

	// Concurrently update the same block
	const numGoroutines = 5
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Fetch the block
			b, err := store.BlockFindByID(ctx, blockID)
			if err != nil {
				errors <- err
				return
			}

			// Update it
			b.SetContent("Updated Content")
			if err := store.BlockUpdate(ctx, b); err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors (some may occur due to race conditions)
	errorCount := 0
	for err := range errors {
		errorCount++
		t.Logf("Concurrent update error (expected): %v", err)
	}

	// Verify final state
	finalBlock, err := store.BlockFindByID(ctx, blockID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if finalBlock == nil {
		t.Fatal("finalBlock MUST NOT be nil")
	}
	if finalBlock.Content() != "Updated Content" {
		t.Fatalf("Expected Content 'Updated Content', got %s", finalBlock.Content())
	}
}

// TestConcurrentReadWrite tests concurrent reads and writes
func TestConcurrentReadWrite(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create initial blocks
	var blockIDs []string
	for i := 0; i < 5; i++ {
		block := NewBlock().
			SetSiteID("Site1").
			SetName("Block")
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		blockIDs = append(blockIDs, block.ID())
	}

	// Concurrent readers and writers
	const numReaders = 10
	const numWriters = 5
	var wg sync.WaitGroup

	// Start readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_, _ = store.BlockList(ctx, BlockQuery().SetSiteID("Site1"))
			}
		}()
	}

	// Start writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			block := NewBlock().
				SetSiteID("Site1").
				SetName("New Block")
			_ = store.BlockCreate(ctx, block)
		}(i)
	}

	wg.Wait()

	// Verify data integrity
	blocks, err := store.BlockList(ctx, BlockQuery().SetSiteID("Site1"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) < 5 {
		t.Fatalf("Expected at least 5 blocks, got %d", len(blocks))
	}
}

// TestConcurrentSoftDelete tests concurrent soft deletes
func TestConcurrentSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create a block
	block := NewBlock().SetSiteID("Site1").SetName("Test Block")
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	blockID := block.ID()

	// Try to soft delete concurrently
	const numGoroutines = 5
	var wg sync.WaitGroup
	successCount := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := store.BlockSoftDeleteByID(ctx, blockID)
			successCount <- (err == nil)
		}()
	}

	wg.Wait()
	close(successCount)

	// Count successes
	successes := 0
	for success := range successCount {
		if success {
			successes++
		}
	}

	// At least one should succeed
	if successes < 1 {
		t.Fatalf("Expected at least 1 success, got %d", successes)
	}

	// Verify block is soft deleted
	found, err := store.BlockFindByID(ctx, blockID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("Expected found to be nil")
	}

	// Should find with soft delete included
	blocks, err := store.BlockList(ctx, BlockQuery().SetID(blockID).SetSoftDeleteIncluded(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(blocks))
	}
	if !blocks[0].IsSoftDeleted() {
		t.Fatal("Expected block to be soft deleted")
	}
}

// TestConcurrentPageAndBlockCreation tests creating pages and blocks concurrently
func TestConcurrentPageAndBlockCreation(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create site
	site := NewSite().SetName("Test Site")
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const numPages = 5
	var wg sync.WaitGroup
	pageIDs := make(chan string, numPages)

	// Create pages concurrently
	for i := 0; i < numPages; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			page := NewPage().
				SetSiteID(site.ID()).
				SetTitle("Page")

			if err := store.PageCreate(ctx, page); err != nil {
				t.Errorf("Page create error: %v", err)
				return
			}

			pageIDs <- page.ID()

			// Create blocks for this page
			for j := 0; j < 3; j++ {
				block := NewBlock().
					SetSiteID(site.ID()).
					SetPageID(page.ID()).
					SetName("Block")

				if err := store.BlockCreate(ctx, block); err != nil {
					t.Errorf("Block create error: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()
	close(pageIDs)

	// Verify pages were created
	var ids []string
	for id := range pageIDs {
		ids = append(ids, id)
	}
	if len(ids) != numPages {
		t.Fatalf("Expected %d pages, got %d", numPages, len(ids))
	}

	// Verify blocks were created (should be 3 per page)
	blocks, err := store.BlockList(ctx, BlockQuery().SetSiteID(site.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != numPages*3 {
		t.Fatalf("Expected %d blocks, got %d", numPages*3, len(blocks))
	}
}

// TestConcurrentCountOperations tests concurrent count operations
func TestConcurrentCountOperations(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_count_ops",
		PageTableName:      "page_table_count_ops",
		SiteTableName:      "site_table_count_ops",
		TemplateTableName:  "template_table_count_ops",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create initial blocks
	for i := 0; i < 10; i++ {
		block := NewBlock().SetSiteID("Site1").SetName("Block")
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Concurrent count operations
	const numGoroutines = 20
	var wg sync.WaitGroup
	counts := make(chan int64, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count, err := store.BlockCount(ctx, BlockQuery().SetSiteID("Site1"))
			if err != nil {
				t.Errorf("Count error: %v", err)
				return
			}
			counts <- count
		}()
	}

	wg.Wait()
	close(counts)

	// All counts should be consistent
	var countValues []int64
	for count := range counts {
		countValues = append(countValues, count)
	}

	if len(countValues) != numGoroutines {
		t.Fatalf("Expected %d count values, got %d", numGoroutines, len(countValues))
	}
	// All counts should be the same (10)
	for _, count := range countValues {
		if count != int64(10) {
			t.Fatalf("Expected count 10, got %d", count)
		}
	}
}
