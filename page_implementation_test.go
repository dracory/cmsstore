package cmsstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/require"
)

func TestNewPageDefaults(t *testing.T) {
	page := NewPage()

	require.NotEmpty(t, page.ID())
	require.NotEmpty(t, page.CreatedAt())
	require.NotEmpty(t, page.UpdatedAt())
	require.Equal(t, PAGE_STATUS_DRAFT, page.Status())
	require.Equal(t, sb.MAX_DATETIME, page.SoftDeletedAt())
	require.False(t, page.IsSoftDeleted())

	metas, err := page.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	createdCarbon := page.CreatedAtCarbon()
	require.NotNil(t, createdCarbon)
	require.Equal(t, page.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

	updatedCarbon := page.UpdatedAtCarbon()
	require.NotNil(t, updatedCarbon)
	require.Equal(t, page.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

	softDeletedCarbon := page.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedCarbon)
	require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestPageGetterMethods(t *testing.T) {
	page := NewPage()

	// Test default values
	require.Equal(t, "", page.Alias())
	require.Equal(t, "", page.CanonicalUrl())
	require.Equal(t, "", page.Content())
	require.Equal(t, "", page.Editor())
	require.Equal(t, "", page.Handle())
	require.Equal(t, "", page.Memo())
	require.Equal(t, "", page.MetaDescription())
	require.Equal(t, "", page.MetaKeywords())
	require.Equal(t, "", page.MetaRobots())
	require.Equal(t, "", page.Name())
	require.Equal(t, "", page.SiteID())
	require.Equal(t, "", page.TemplateID())
	require.Equal(t, "", page.Title())
}

func TestPageStatusMethods(t *testing.T) {
	page := NewPage()

	// Test default status (DRAFT)
	require.False(t, page.IsActive())
	require.False(t, page.IsInactive())

	// Test ACTIVE status
	page.SetStatus(PAGE_STATUS_ACTIVE)
	require.True(t, page.IsActive())
	require.False(t, page.IsInactive())

	// Test INACTIVE status
	page.SetStatus(PAGE_STATUS_INACTIVE)
	require.False(t, page.IsActive())
	require.True(t, page.IsInactive())

	// Test other status values
	page.SetStatus("unknown")
	require.False(t, page.IsActive())
	require.False(t, page.IsInactive())
}

func TestPageSoftDeleteMethods(t *testing.T) {
	page := NewPage()
	require.False(t, page.IsSoftDeleted())

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	page.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	require.False(t, page.IsSoftDeleted())
	require.Equal(t, future.ToDateTimeString(carbon.UTC), page.SoftDeletedAt())

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	page.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	require.True(t, page.IsSoftDeleted())
	require.Equal(t, past.ToDateTimeString(carbon.UTC), page.SoftDeletedAt())
}

func TestPageMetasMethods(t *testing.T) {
	page := NewPage()

	// Test empty metas
	metas, err := page.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	// Test Meta lookup on empty metas
	require.Equal(t, "", page.Meta("nonexistent"))

	// Test SetMetas
	err = page.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	require.NoError(t, err)

	metas, err = page.Metas()
	require.NoError(t, err)
	require.Equal(t, "main", metas["layout"])
	require.Equal(t, "dark", metas["theme"])

	// Test Meta lookup
	require.Equal(t, "main", page.Meta("layout"))
	require.Equal(t, "dark", page.Meta("theme"))
	require.Equal(t, "", page.Meta("nonexistent"))

	// Test SetMeta
	err = page.SetMeta("newkey", "newvalue")
	require.NoError(t, err)
	require.Equal(t, "newvalue", page.Meta("newkey"))

	// Test UpsertMetas
	err = page.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	require.NoError(t, err)
	require.Equal(t, "sidebar", page.Meta("layout")) // Updated
	require.Equal(t, "dark", page.Meta("theme"))     // Preserved
	require.Equal(t, "newvalue", page.Meta("newkey")) // Preserved
	require.Equal(t, "blue", page.Meta("color"))      // Added
}

func TestPageMiddlewaresMethods(t *testing.T) {
	page := NewPage()

	// Test default middlewares
	require.Empty(t, page.MiddlewaresBefore())
	require.Empty(t, page.MiddlewaresAfter())

	// Test SetMiddlewaresBefore
	before := []string{"auth", "csrf"}
	page.SetMiddlewaresBefore(before)
	require.Equal(t, before, page.MiddlewaresBefore())

	// Test SetMiddlewaresAfter
	after := []string{"log", "cache"}
	page.SetMiddlewaresAfter(after)
	require.Equal(t, after, page.MiddlewaresAfter())

	// Test empty middlewares
	page.SetMiddlewaresBefore([]string{})
	page.SetMiddlewaresAfter([]string{})
	require.Empty(t, page.MiddlewaresBefore())
	require.Empty(t, page.MiddlewaresAfter())
}

func TestPageCreatedAtMethods(t *testing.T) {
	page := NewPage()

	// Test default CreatedAt
	createdAt := page.CreatedAt()
	require.NotEmpty(t, createdAt)

	createdAtCarbon := page.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, createdAt, createdAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	page.SetCreatedAt(testDate)
	require.Equal(t, testDate, page.CreatedAt())

	createdAtCarbon = page.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, testDate, createdAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestPageUpdatedAtMethods(t *testing.T) {
	page := NewPage()

	// Test default UpdatedAt
	updatedAt := page.UpdatedAt()
	require.NotEmpty(t, updatedAt)

	updatedAtCarbon := page.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, updatedAt, updatedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	page.SetUpdatedAt(testDate)
	require.Equal(t, testDate, page.UpdatedAt())

	updatedAtCarbon = page.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, testDate, updatedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestPageSoftDeletedAtMethods(t *testing.T) {
	page := NewPage()

	// Test default SoftDeletedAt
	softDeletedAt := page.SoftDeletedAt()
	require.Equal(t, sb.MAX_DATETIME, softDeletedAt)

	softDeletedAtCarbon := page.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, softDeletedAt, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	page.SetSoftDeletedAt(testDate)
	require.Equal(t, testDate, page.SoftDeletedAt())

	softDeletedAtCarbon = page.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, testDate, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestPageIDMethods(t *testing.T) {
	page := NewPage()

	// Test default ID
	id := page.ID()
	require.NotEmpty(t, id)

	// Test SetID
	newID := "test-page-id-123"
	page.SetID(newID)
	require.Equal(t, newID, page.ID())
}

func TestPageAliasMethods(t *testing.T) {
	page := NewPage()

	// Test default alias
	require.Equal(t, "", page.Alias())

	// Test SetAlias
	alias := "test-page-alias"
	page.SetAlias(alias)
	require.Equal(t, alias, page.Alias())
}

func TestPageCanonicalUrlMethods(t *testing.T) {
	page := NewPage()

	// Test default canonical URL
	require.Equal(t, "", page.CanonicalUrl())

	// Test SetCanonicalUrl
	canonicalUrl := "https://example.com/canonical"
	page.SetCanonicalUrl(canonicalUrl)
	require.Equal(t, canonicalUrl, page.CanonicalUrl())
}

func TestPageContentMethods(t *testing.T) {
	page := NewPage()

	// Test default content
	require.Equal(t, "", page.Content())

	// Test SetContent
	content := "This is page content"
	page.SetContent(content)
	require.Equal(t, content, page.Content())
}

func TestPageEditorMethods(t *testing.T) {
	page := NewPage()

	// Test default editor
	require.Equal(t, "", page.Editor())

	// Test SetEditor
	editor := "test-editor"
	page.SetEditor(editor)
	require.Equal(t, editor, page.Editor())
}

func TestPageHandleMethods(t *testing.T) {
	page := NewPage()

	// Test default handle
	require.Equal(t, "", page.Handle())

	// Test SetHandle
	handle := "test-page-handle"
	page.SetHandle(handle)
	require.Equal(t, handle, page.Handle())
}

func TestPageMemoMethods(t *testing.T) {
	page := NewPage()

	// Test default memo
	require.Equal(t, "", page.Memo())

	// Test SetMemo
	memo := "This is a page memo"
	page.SetMemo(memo)
	require.Equal(t, memo, page.Memo())
}

func TestPageMetaDescriptionMethods(t *testing.T) {
	page := NewPage()

	// Test default meta description
	require.Equal(t, "", page.MetaDescription())

	// Test SetMetaDescription
	metaDescription := "This is a meta description"
	page.SetMetaDescription(metaDescription)
	require.Equal(t, metaDescription, page.MetaDescription())
}

func TestPageMetaKeywordsMethods(t *testing.T) {
	page := NewPage()

	// Test default meta keywords
	require.Equal(t, "", page.MetaKeywords())

	// Test SetMetaKeywords
	metaKeywords := "keyword1, keyword2, keyword3"
	page.SetMetaKeywords(metaKeywords)
	require.Equal(t, metaKeywords, page.MetaKeywords())
}

func TestPageMetaRobotsMethods(t *testing.T) {
	page := NewPage()

	// Test default meta robots
	require.Equal(t, "", page.MetaRobots())

	// Test SetMetaRobots
	metaRobots := "noindex, nofollow"
	page.SetMetaRobots(metaRobots)
	require.Equal(t, metaRobots, page.MetaRobots())
}

func TestPageNameMethods(t *testing.T) {
	page := NewPage()

	// Test default name
	require.Equal(t, "", page.Name())

	// Test SetName
	name := "Test Page Name"
	page.SetName(name)
	require.Equal(t, name, page.Name())
}

func TestPageSiteIDMethods(t *testing.T) {
	page := NewPage()

	// Test default site ID
	require.Equal(t, "", page.SiteID())

	// Test SetSiteID
	siteID := "test-site-id"
	page.SetSiteID(siteID)
	require.Equal(t, siteID, page.SiteID())
}

func TestPageTemplateIDMethods(t *testing.T) {
	page := NewPage()

	// Test default template ID
	require.Equal(t, "", page.TemplateID())

	// Test SetTemplateID
	templateID := "test-template-id"
	page.SetTemplateID(templateID)
	require.Equal(t, templateID, page.TemplateID())
}

func TestPageTitleMethods(t *testing.T) {
	page := NewPage()

	// Test default title
	require.Equal(t, "", page.Title())

	// Test SetTitle
	title := "Test Page Title"
	page.SetTitle(title)
	require.Equal(t, title, page.Title())
}

func TestPageStatusSettersAndGetters(t *testing.T) {
	page := NewPage()

	// Test default status
	require.Equal(t, PAGE_STATUS_DRAFT, page.Status())

	// Test SetStatus
	page.SetStatus(PAGE_STATUS_ACTIVE)
	require.Equal(t, PAGE_STATUS_ACTIVE, page.Status())

	page.SetStatus(PAGE_STATUS_INACTIVE)
	require.Equal(t, PAGE_STATUS_INACTIVE, page.Status())

	page.SetStatus("custom-status")
	require.Equal(t, "custom-status", page.Status())
}

func TestPageMarshalToVersioning(t *testing.T) {
	page := NewPage()
	page.SetAlias("test-alias")
	page.SetCanonicalUrl("https://example.com/canonical")
	page.SetContent("Test content")
	page.SetEditor("test-editor")
	page.SetHandle("test-handle")
	page.SetMemo("test-memo")
	page.SetMetaDescription("Test meta description")
	page.SetMetaKeywords("test, keywords")
	page.SetMetaRobots("index, follow")
	page.SetName("Test Page")
	page.SetSiteID("test-site")
	page.SetTemplateID("test-template")
	page.SetTitle("Test Page Title")
	page.SetStatus(PAGE_STATUS_ACTIVE)

	versionedJSON, err := page.MarshalToVersioning()
	require.NoError(t, err)
	require.NotEmpty(t, versionedJSON)

	// Parse the JSON to verify it contains expected fields
	var versionedData map[string]string
	err = json.Unmarshal([]byte(versionedJSON), &versionedData)
	require.NoError(t, err)

	// Check that expected fields are present
	require.Equal(t, "test-alias", versionedData[COLUMN_ALIAS])
	require.Equal(t, "https://example.com/canonical", versionedData[COLUMN_CANONICAL_URL])
	require.Equal(t, "Test content", versionedData[COLUMN_CONTENT])
	require.Equal(t, "test-editor", versionedData[COLUMN_EDITOR])
	require.Equal(t, "test-handle", versionedData[COLUMN_HANDLE])
	require.Equal(t, "test-memo", versionedData[COLUMN_MEMO])
	require.Equal(t, "Test meta description", versionedData[COLUMN_META_DESCRIPTION])
	require.Equal(t, "test, keywords", versionedData[COLUMN_META_KEYWORDS])
	require.Equal(t, "index, follow", versionedData[COLUMN_META_ROBOTS])
	require.Equal(t, "Test Page", versionedData[COLUMN_NAME])
	require.Equal(t, "test-site", versionedData[COLUMN_SITE_ID])
	require.Equal(t, "test-template", versionedData[COLUMN_TEMPLATE_ID])
	require.Equal(t, "Test Page Title", versionedData[COLUMN_TITLE])
	require.Equal(t, PAGE_STATUS_ACTIVE, versionedData[COLUMN_STATUS])

	// Check that timestamps and soft delete fields are excluded
	_, hasCreatedAt := versionedData[COLUMN_CREATED_AT]
	_, hasUpdatedAt := versionedData[COLUMN_UPDATED_AT]
	_, hasSoftDeletedAt := versionedData[COLUMN_SOFT_DELETED_AT]
	require.False(t, hasCreatedAt)
	require.False(t, hasUpdatedAt)
	require.False(t, hasSoftDeletedAt)
}