package cmsstore

import (
	"context"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Test Translation")
	err = translation.SetContent(map[string]string{"en": "Hello"})
	require.NoError(t, err)

	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	require.NoError(t, err)
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
	require.NoError(t, err)

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Test Translation")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	require.NoError(t, err)

	// Test find by full ID
	found, err := store.TranslationFindByID(ctx, translation.ID())
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, translation.ID(), found.ID())

	// Test find by shortened ID
	shortID := ShortenID(translation.ID())
	foundShort, err := store.TranslationFindByID(ctx, shortID)
	require.NoError(t, err)
	require.NotNil(t, foundShort)
	require.Equal(t, translation.ID(), foundShort.ID())
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
	require.NoError(t, err)

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Test Translation")
	translation.SetHandle("test-handle")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	require.NoError(t, err)

	found, err := store.TranslationFindByHandle(ctx, "test-handle")
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, translation.ID(), found.ID())
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
	require.NoError(t, err)

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Test Translation")
	translation.SetHandle("test-handle")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	require.NoError(t, err)

	// Find by handle
	found, err := store.TranslationFindByHandleOrID(ctx, "test-handle", "en")
	require.NoError(t, err)
	require.NotNil(t, found)

	// Find by ID
	found, err = store.TranslationFindByHandleOrID(ctx, translation.ID(), "en")
	require.NoError(t, err)
	require.NotNil(t, found)
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
	require.NoError(t, err)

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	translation.SetName("Original Name")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	require.NoError(t, err)

	translation.SetName("Updated Name")
	err = store.TranslationUpdate(ctx, translation)
	require.NoError(t, err)

	found, _ := store.TranslationFindByID(ctx, translation.ID())
	require.Equal(t, "Updated Name", found.Name())
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
	require.NoError(t, err)

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	require.NoError(t, err)

	err = store.TranslationDelete(ctx, translation)
	require.NoError(t, err)

	found, _ := store.TranslationFindByID(ctx, translation.ID())
	require.Nil(t, found)
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
	require.NoError(t, err)

	translation := NewTranslation()
	translation.SetSiteID("test-site")
	ctx := context.Background()
	err = store.TranslationCreate(ctx, translation)
	require.NoError(t, err)

	err = store.TranslationSoftDeleteByID(ctx, translation.ID())
	require.NoError(t, err)

	found, _ := store.TranslationFindByID(ctx, translation.ID())
	require.Nil(t, found, "Should not find soft deleted translation with default query")

	// Should find it when included
	list, err := store.TranslationList(ctx, TranslationQuery().
		SetID(translation.ID()).
		SetSoftDeletedIncluded(true))
	require.NoError(t, err)
	require.Len(t, list, 1)
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
	require.NoError(t, err)

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		t := NewTranslation()
		t.SetSiteID("test-site")
		t.SetName("Translation " + carbon.Now().ToIso8601String())
		_ = store.TranslationCreate(ctx, t)
	}

	list, err := store.TranslationList(ctx, TranslationQuery().SetLimit(3))
	require.NoError(t, err)
	require.Len(t, list, 3)
}

func TestStoreTranslationCount(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		t := NewTranslation()
		t.SetSiteID("test-site")
		_ = store.TranslationCreate(ctx, t)
	}

	count, err := store.TranslationCount(ctx, TranslationQuery())
	require.NoError(t, err)
	require.Equal(t, int64(3), count)
}

func TestTranslationQueryMethods(t *testing.T) {
	q := TranslationQuery()
	
	q.SetColumns([]string{"id", "name"})
	require.Equal(t, []string{"id", "name"}, q.Columns())
	
	q.SetCountOnly(true)
	require.True(t, q.IsCountOnly())
	
	q.SetCreatedAtGte("2023-01-01")
	require.Equal(t, "2023-01-01", q.CreatedAtGte())
	
	q.SetCreatedAtLte("2023-12-31")
	require.Equal(t, "2023-12-31", q.CreatedAtLte())
	
	q.SetHandle("handle")
	require.Equal(t, "handle", q.Handle())
	
	q.SetHandleOrID("handleorid")
	require.Equal(t, "handleorid", q.HandleOrID())
	
	q.SetID("id")
	require.Equal(t, "id", q.ID())
	
	q.SetIDIn([]string{"id1", "id2"})
	require.Equal(t, []string{"id1", "id2"}, q.IDIn())
	
	q.SetLimit(10)
	require.Equal(t, 10, q.Limit())
	
	q.SetNameLike("%test%")
	require.Equal(t, "%test%", q.NameLike())
	
	q.SetOffset(5)
	require.Equal(t, 5, q.Offset())
	
	q.SetOrderBy("name")
	require.Equal(t, "name", q.OrderBy())
	
	q.SetSiteID("siteid")
	require.Equal(t, "siteid", q.SiteID())
	
	q.SetSoftDeletedIncluded(true)
	require.True(t, q.SoftDeletedIncluded())
	
	q.SetSortOrder(sb.ASC)
	require.Equal(t, sb.ASC, q.SortOrder())
	
	q.SetStatus(TRANSLATION_STATUS_ACTIVE)
	require.Equal(t, TRANSLATION_STATUS_ACTIVE, q.Status())
	
	q.SetStatusIn([]string{TRANSLATION_STATUS_ACTIVE, TRANSLATION_STATUS_DRAFT})
	require.Equal(t, []string{TRANSLATION_STATUS_ACTIVE, TRANSLATION_STATUS_DRAFT}, q.StatusIn())
	
	// Validate success
	require.NoError(t, q.Validate())
	
	// Validate failure cases
	errQ := TranslationQuery().SetCreatedAtGte("")
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetCreatedAtLte("")
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetHandle("")
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetHandleOrID("")
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetID("")
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetIDIn([]string{})
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetLimit(-1)
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetNameLike("")
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetOffset(-1)
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetSiteID("")
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetStatus("")
	require.Error(t, errQ.Validate())
	
	errQ = TranslationQuery().SetStatusIn([]string{})
	require.Error(t, errQ.Validate())
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
	require.NoError(t, err)
	
	require.Equal(t, "en", store.TranslationLanguageDefault())
	require.Equal(t, langs, store.TranslationLanguages())
}
