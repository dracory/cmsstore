package cmsstore

import (
	"context"
	"testing"
)

func TestStoreMediaCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetEntityID("entity-1").
		SetEntityType("page").
		SetTitle("Test Image").
		SetURL("https://example.com/image.jpg").
		SetType("image/jpeg").
		SetStatus(MEDIA_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if media.ID() == "" {
		t.Fatal("expected non-empty ID after create")
	}
	if media.CreatedAt() == "" {
		t.Fatal("expected non-empty CreatedAt after create")
	}
	if media.UpdatedAt() == "" {
		t.Fatal("expected non-empty UpdatedAt after create")
	}
}

func TestStoreMediaFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetEntityID("entity-1").
		SetEntityType("page").
		SetTitle("Test Image").
		SetURL("https://example.com/image.jpg").
		SetType("image/jpeg").
		SetExtension("jpg").
		SetSize("1024").
		SetStatus(MEDIA_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected non-nil media")
	}
	if found.ID() != media.ID() {
		t.Fatalf("expected ID %q, got %q", media.ID(), found.ID())
	}
	if found.Title() != "Test Image" {
		t.Fatalf("expected Title %q, got %q", "Test Image", found.Title())
	}
	if found.URL() != "https://example.com/image.jpg" {
		t.Fatalf("expected URL %q, got %q", "https://example.com/image.jpg", found.URL())
	}
	if found.Type() != "image/jpeg" {
		t.Fatalf("expected Type %q, got %q", "image/jpeg", found.Type())
	}
	if found.Extension() != "jpg" {
		t.Fatalf("expected Extension %q, got %q", "jpg", found.Extension())
	}
	if found.Size() != "1024" {
		t.Fatalf("expected Size %q, got %q", "1024", found.Size())
	}
	if found.EntityID() != "entity-1" {
		t.Fatalf("expected EntityID %q, got %q", "entity-1", found.EntityID())
	}
	if found.EntityType() != "page" {
		t.Fatalf("expected EntityType %q, got %q", "page", found.EntityType())
	}
	if !found.IsImage() {
		t.Fatal("expected IsImage to be true")
	}
}

func TestStoreMediaFindByIDNotFound(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	found, err := store.MediaFindByID(ctx, "nonexistent-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil media for nonexistent ID")
	}
}

func TestStoreMediaFindByHandle(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetTitle("Test Media").
		SetHandle("test-media-handle").
		SetStatus(MEDIA_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.MediaFindByHandle(ctx, "test-media-handle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected non-nil media")
	}
	if found.ID() != media.ID() {
		t.Fatalf("expected ID %q, got %q", media.ID(), found.ID())
	}
	if found.Handle() != "test-media-handle" {
		t.Fatalf("expected Handle %q, got %q", "test-media-handle", found.Handle())
	}
}

func TestStoreMediaList(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_list",
		PageTableName:      "page_table_list",
		SiteTableName:      "site_table_list",
		TemplateTableName:  "template_table_list",
		MediaEnabled:       true,
		MediaTableName:     "media_table_list",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	for i := 0; i < 3; i++ {
		media := NewMedia().
			SetTitle("Media " + string(rune('A'+i))).
			SetStatus(MEDIA_STATUS_ACTIVE)
		err = store.MediaCreate(ctx, media)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	list, err := store.MediaList(ctx, MediaQuery())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 media items, got %d", len(list))
	}
}

func TestStoreMediaListByEntityID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	media1 := NewMedia().
		SetEntityID("page-1").
		SetEntityType("page").
		SetTitle("Media 1").
		SetSequenceInt(2).
		SetStatus(MEDIA_STATUS_ACTIVE)
	err = store.MediaCreate(ctx, media1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media2 := NewMedia().
		SetEntityID("page-1").
		SetEntityType("page").
		SetTitle("Media 2").
		SetSequenceInt(1).
		SetStatus(MEDIA_STATUS_ACTIVE)
	err = store.MediaCreate(ctx, media2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media3 := NewMedia().
		SetEntityID("block-1").
		SetEntityType("block").
		SetTitle("Media 3").
		SetStatus(MEDIA_STATUS_ACTIVE)
	err = store.MediaCreate(ctx, media3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, err := store.MediaListByEntityID(ctx, "page-1", "page")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 media items for page-1, got %d", len(list))
	}

	if list[0].SequenceInt() != 1 {
		t.Fatalf("expected first item SequenceInt 1, got %d", list[0].SequenceInt())
	}
	if list[1].SequenceInt() != 2 {
		t.Fatalf("expected second item SequenceInt 2, got %d", list[1].SequenceInt())
	}
}

func TestStoreMediaCount(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_count",
		PageTableName:      "page_table_count",
		SiteTableName:      "site_table_count",
		TemplateTableName:  "template_table_count",
		MediaEnabled:       true,
		MediaTableName:     "media_table_count",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		media := NewMedia().
			SetTitle("Media " + string(rune('A'+i))).
			SetStatus(MEDIA_STATUS_ACTIVE)
		err = store.MediaCreate(ctx, media)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	count, err := store.MediaCount(ctx, MediaQuery())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected count 5, got %d", count)
	}

	count, err = store.MediaCount(ctx, MediaQuery().SetStatus(MEDIA_STATUS_ACTIVE))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected count 5 for active, got %d", count)
	}

	count, err = store.MediaCount(ctx, MediaQuery().SetStatus(MEDIA_STATUS_DRAFT))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0 for draft, got %d", count)
	}
}

func TestStoreMediaUpdate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetTitle("Original Title").
		SetStatus(MEDIA_STATUS_DRAFT)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media.SetTitle("Updated Title")
	media.SetStatus(MEDIA_STATUS_ACTIVE)
	err = store.MediaUpdate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected non-nil media")
	}
	if found.Title() != "Updated Title" {
		t.Fatalf("expected Title %q, got %q", "Updated Title", found.Title())
	}
	if found.Status() != MEDIA_STATUS_ACTIVE {
		t.Fatalf("expected Status %q, got %q", MEDIA_STATUS_ACTIVE, found.Status())
	}
}

func TestStoreMediaSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetTitle("Test Media").
		SetStatus(MEDIA_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.MediaSoftDelete(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil media after soft delete")
	}

	list, err := store.MediaList(ctx, MediaQuery().
		SetID(media.ID()).
		SetSoftDeletedIncluded(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 media with soft delete included, got %d", len(list))
	}
	if !list[0].IsSoftDeleted() {
		t.Fatal("expected media to be soft deleted")
	}
}

func TestStoreMediaSoftDeleteByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetTitle("Test Media").
		SetStatus(MEDIA_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.MediaSoftDeleteByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil media after soft delete")
	}
}

func TestStoreMediaDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetTitle("Test Media").
		SetStatus(MEDIA_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.MediaDelete(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil media after delete")
	}
}

func TestStoreMediaDeleteByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetTitle("Test Media").
		SetStatus(MEDIA_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.MediaDeleteByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil media after delete")
	}
}

func TestStoreMediaMetas(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media := NewMedia().
		SetTitle("Test Media").
		SetStatus(MEDIA_STATUS_ACTIVE)

	err = media.SetMeta("alt", "Alternative text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = media.SetMeta("caption", "Test caption")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	err = store.MediaCreate(ctx, media)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.MediaFindByID(ctx, media.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected non-nil media")
	}
	if found.Meta("alt") != "Alternative text" {
		t.Fatalf("expected Meta('alt') %q, got %q", "Alternative text", found.Meta("alt"))
	}
	if found.Meta("caption") != "Test caption" {
		t.Fatalf("expected Meta('caption') %q, got %q", "Test caption", found.Meta("caption"))
	}
}

func TestStoreMediaEnabledFlag(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "media_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !store.MediaEnabled() {
		t.Fatal("expected MediaEnabled to be true")
	}

	store2, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_2",
		PageTableName:      "page_table_2",
		SiteTableName:      "site_table_2",
		TemplateTableName:  "template_table_2",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store2.MediaEnabled() {
		t.Fatal("expected MediaEnabled to be false")
	}
}

func TestStoreMediaValidationMissingTableName(t *testing.T) {
	db := initDB(":memory:")

	_, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MediaEnabled:       true,
		MediaTableName:     "",
		AutomigrateEnabled: true,
	})
	if err == nil {
		t.Fatal("expected error for missing MediaTableName")
	}
}
