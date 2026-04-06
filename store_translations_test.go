package cmsstore

import (
	"context"
	"slices"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	_ "modernc.org/sqlite"
)

func TestStoreTranslationCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Test Translation")
	err = translation.SetContent(map[string]string{"en": "Hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStoreTranslationFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Test Translation")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test find by full ID
	found, err := store.TranslationFindByID(ctx, translation.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("found MUST NOT be nil")
	}
	if found.ID() != translation.ID() {
		t.Fatalf("Expected ID %s, got %s", translation.ID(), found.ID())
	}

	// Test find by shortened ID
	shortID := ShortenID(translation.ID())
	foundShort, err := store.TranslationFindByID(ctx, shortID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if foundShort == nil {
		t.Fatal("foundShort MUST NOT be nil")
	}
	if foundShort.ID() != translation.ID() {
		t.Fatalf("Expected ID %s, got %s", translation.ID(), foundShort.ID())
	}
}

func TestStoreTranslationFindByHandle(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Test Translation")
	translation.SetHandle("test-handle")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := store.TranslationFindByHandle(ctx, "test-handle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("found MUST NOT be nil")
	}
	if found.ID() != translation.ID() {
		t.Fatalf("Expected ID %s, got %s", translation.ID(), found.ID())
	}
}

func TestStoreTranslationFindByHandleOrID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Test Translation")
	translation.SetHandle("test-handle")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find by handle
	found, err := store.TranslationFindByHandleOrID(ctx, "test-handle", "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("found MUST NOT be nil")
	}

	// Find by ID
	found, err = store.TranslationFindByHandleOrID(ctx, translation.ID(), "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("found MUST NOT be nil")
	}
}

func TestStoreTranslationErrorPaths(t *testing.T) {
	ctx := context.Background()

	// Test with nil DB
	store := &storeImplementation{db: nil}

	_, err := store.TranslationCount(ctx, TranslationQuery())
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationCreate(ctx, NewTranslation())
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationDelete(ctx, NewTranslation())
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationDeleteByID(ctx, "id")
	if err == nil {
		t.Error("Expected error")
	}

	_, err = store.TranslationFindByHandle(ctx, "handle")
	if err == nil {
		t.Error("Expected error")
	}

	_, err = store.TranslationFindByID(ctx, "id")
	if err == nil {
		t.Error("Expected error")
	}

	_, err = store.TranslationFindByHandleOrID(ctx, "handle-or-id", "en")
	if err == nil {
		t.Error("Expected error")
	}

	_, err = store.TranslationList(ctx, TranslationQuery())
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationSoftDelete(ctx, NewTranslation())
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationSoftDeleteByID(ctx, "id")
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationUpdate(ctx, NewTranslation())
	if err == nil {
		t.Error("Expected error")
	}

	// Test with nil entity
	store.db = initDB(":memory:")
	err = store.TranslationCreate(ctx, nil)
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationDelete(ctx, nil)
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationSoftDelete(ctx, nil)
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationUpdate(ctx, nil)
	if err == nil {
		t.Error("Expected error")
	}

	// Test with empty ID/handle
	_, err = store.TranslationFindByHandle(ctx, "")
	if err == nil {
		t.Error("Expected error")
	}

	_, err = store.TranslationFindByID(ctx, "")
	if err == nil {
		t.Error("Expected error")
	}

	_, err = store.TranslationFindByHandleOrID(ctx, "", "en")
	if err == nil {
		t.Error("Expected error")
	}

	err = store.TranslationDeleteByID(ctx, "")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestStoreTranslationUpdate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Original Name")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	translation.SetName("Updated Name")
	err = store.TranslationUpdate(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := store.TranslationFindByID(ctx, translation.ID())
	if found == nil {
		t.Fatal("found MUST NOT be nil")
	}
	if found.Name() != "Updated Name" {
		t.Fatalf("Expected name 'Updated Name', got %s", found.Name())
	}
}

func TestStoreTranslationDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.TranslationDelete(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := store.TranslationFindByID(ctx, translation.ID())
	if found != nil {
		t.Fatal("Expected found to be nil")
	}
}

func TestStoreTranslationSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.TranslationSoftDeleteByID(ctx, translation.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := store.TranslationFindByID(ctx, translation.ID())
	if found != nil {
		t.Fatal("Should not find soft deleted translation with default query")
	}

	// Should find it when included
	list, err := store.TranslationList(ctx, TranslationQuery().
		SetID(translation.ID()).
		SetSoftDeletedIncluded(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(list))
	}
}

func TestStoreTranslationList(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		t := NewTranslation()
		t.SetSiteID("test-site")
		t.SetName("Translation " + carbon.Now().ToIso8601String())
		_ = store.TranslationCreate(ctx, t)
	}

	list, err := store.TranslationList(ctx, TranslationQuery().SetLimit(3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(list))
	}
}

func TestStoreTranslationCount(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table_trans_count",
		PageTableName:        "page_table_trans_count",
		SiteTableName:        "site_table_trans_count",
		TemplateTableName:    "template_table_trans_count",
		TranslationsEnabled:  true,
		TranslationTableName: "translation_table_trans_count",
		AutomigrateEnabled:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		t := NewTranslation()
		t.SetSiteID("test-site")
		_ = store.TranslationCreate(ctx, t)
	}

	count, err := store.TranslationCount(ctx, TranslationQuery())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != int64(3) {
		t.Fatalf("Expected count 3, got %d", count)
	}
}

func TestTranslationQueryMethods(t *testing.T) {
	q := TranslationQuery()

	q.SetColumns([]string{"id", "name"})
	if !slices.Equal(q.Columns(), []string{"id", "name"}) {
		t.Fatalf("Expected columns [id name], got %v", q.Columns())
	}

	q.SetCountOnly(true)
	if !q.IsCountOnly() {
		t.Error("Expected IsCountOnly to be true")
	}

	q.SetCreatedAtGte("2023-01-01")
	if q.CreatedAtGte() != "2023-01-01" {
		t.Errorf("Expected CreatedAtGte '2023-01-01', got %s", q.CreatedAtGte())
	}

	q.SetCreatedAtLte("2023-12-31")
	if q.CreatedAtLte() != "2023-12-31" {
		t.Errorf("Expected CreatedAtLte '2023-12-31', got %s", q.CreatedAtLte())
	}

	q.SetHandle("handle")
	if q.Handle() != "handle" {
		t.Errorf("Expected Handle 'handle', got %s", q.Handle())
	}

	q.SetHandleOrID("handleorid")
	if q.HandleOrID() != "handleorid" {
		t.Errorf("Expected HandleOrID 'handleorid', got %s", q.HandleOrID())
	}

	q.SetID("id")
	if q.ID() != "id" {
		t.Errorf("Expected ID 'id', got %s", q.ID())
	}

	q.SetIDIn([]string{"id1", "id2"})
	if !slices.Equal(q.IDIn(), []string{"id1", "id2"}) {
		t.Fatalf("Expected IDIn [id1 id2], got %v", q.IDIn())
	}

	q.SetLimit(10)
	if q.Limit() != 10 {
		t.Errorf("Expected Limit 10, got %d", q.Limit())
	}

	q.SetNameLike("%test%")
	if q.NameLike() != "%test%" {
		t.Errorf("Expected NameLike '%%test%%', got %s", q.NameLike())
	}

	q.SetOffset(5)
	if q.Offset() != 5 {
		t.Errorf("Expected Offset 5, got %d", q.Offset())
	}

	q.SetOrderBy("name")
	if q.OrderBy() != "name" {
		t.Errorf("Expected OrderBy 'name', got %s", q.OrderBy())
	}

	q.SetSiteID("siteid")
	if q.SiteID() != "siteid" {
		t.Errorf("Expected SiteID 'siteid', got %s", q.SiteID())
	}

	q.SetSoftDeletedIncluded(true)
	if !q.SoftDeletedIncluded() {
		t.Error("Expected SoftDeletedIncluded to be true")
	}

	q.SetSortOrder(sb.ASC)
	if q.SortOrder() != sb.ASC {
		t.Errorf("Expected SortOrder sb.ASC, got %v", q.SortOrder())
	}

	q.SetStatus(TRANSLATION_STATUS_ACTIVE)
	if q.Status() != TRANSLATION_STATUS_ACTIVE {
		t.Errorf("Expected Status %s, got %s", TRANSLATION_STATUS_ACTIVE, q.Status())
	}

	q.SetStatusIn([]string{TRANSLATION_STATUS_ACTIVE, TRANSLATION_STATUS_DRAFT})
	if !slices.Equal(q.StatusIn(), []string{TRANSLATION_STATUS_ACTIVE, TRANSLATION_STATUS_DRAFT}) {
		t.Fatalf("Expected StatusIn [%s %s], got %v", TRANSLATION_STATUS_ACTIVE, TRANSLATION_STATUS_DRAFT, q.StatusIn())
	}

	// Validate success
	if err := q.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate failure cases
	errQ := TranslationQuery().SetCreatedAtGte("")
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetCreatedAtLte("")
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetHandle("")
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetHandleOrID("")
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetID("")
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetIDIn([]string{})
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetLimit(-1)
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetNameLike("")
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetOffset(-1)
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetSiteID("")
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetStatus("")
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}

	errQ = TranslationQuery().SetStatusIn([]string{})
	if errQ.Validate() == nil {
		t.Error("Expected error")
	}
}

func TestStoreTranslationLanguages(t *testing.T) {
	langs := map[string]string{"en": "English", "fr": "French"}
	store, err := NewStore(NewStoreOptions{
		DB:                         initDB(":memory:"),
		BlockTableName:             "block_table",
		PageTableName:              "page_table",
		SiteTableName:              "site_table",
		TemplateTableName:          "template_table",
		TranslationsEnabled:        true,
		TranslationTableName:       "translation_table",
		TranslationLanguageDefault: "en",
		TranslationLanguages:       langs,
		AutomigrateEnabled:         true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.TranslationLanguageDefault() != "en" {
		t.Fatalf("Expected default language 'en', got %s", store.TranslationLanguageDefault())
	}
	gotLangs := store.TranslationLanguages()
	if len(gotLangs) != len(langs) {
		t.Fatalf("Expected %d languages, got %d", len(langs), len(gotLangs))
	}
	for k, v := range langs {
		if gotLangs[k] != v {
			t.Fatalf("Expected language %s=%s, got %s", k, v, gotLangs[k])
		}
	}
}
