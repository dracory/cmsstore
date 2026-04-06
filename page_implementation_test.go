package cmsstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewPageDefaults(t *testing.T) {
	page := NewPage()

	if page.ID() == "" {
		t.Error("expected non-empty ID")
	}
	if page.CreatedAt() == "" {
		t.Error("expected non-empty CreatedAt")
	}
	if page.UpdatedAt() == "" {
		t.Error("expected non-empty UpdatedAt")
	}
	if page.Status() != PAGE_STATUS_DRAFT {
		t.Errorf("expected status %q, got %q", PAGE_STATUS_DRAFT, page.Status())
	}
	if page.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", sb.MAX_DATETIME, page.SoftDeletedAt())
	}
	if page.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false")
	}

	metas, err := page.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Error("expected empty metas")
	}

	createdCarbon := page.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if page.CreatedAt() != createdCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAt to match CreatedAtCarbon")
	}

	updatedCarbon := page.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if page.UpdatedAt() != updatedCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAt to match UpdatedAtCarbon")
	}

	softDeletedCarbon := page.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if !softDeletedCarbon.Gte(carbon.Now(carbon.UTC)) {
		t.Error("expected SoftDeletedAtCarbon to be in the future")
	}
}

func TestPageGetterMethods(t *testing.T) {
	page := NewPage()

	// Test default values
	if page.Alias() != "" {
		t.Error("expected empty Alias")
	}
	if page.CanonicalUrl() != "" {
		t.Error("expected empty CanonicalUrl")
	}
	if page.Content() != "" {
		t.Error("expected empty Content")
	}
	if page.Editor() != "" {
		t.Error("expected empty Editor")
	}
	if page.Handle() != "" {
		t.Error("expected empty Handle")
	}
	if page.Memo() != "" {
		t.Error("expected empty Memo")
	}
	if page.MetaDescription() != "" {
		t.Error("expected empty MetaDescription")
	}
	if page.MetaKeywords() != "" {
		t.Error("expected empty MetaKeywords")
	}
	if page.MetaRobots() != "" {
		t.Error("expected empty MetaRobots")
	}
	if page.Name() != "" {
		t.Error("expected empty Name")
	}
	if page.SiteID() != "" {
		t.Error("expected empty SiteID")
	}
	if page.TemplateID() != "" {
		t.Error("expected empty TemplateID")
	}
	if page.Title() != "" {
		t.Error("expected empty Title")
	}
}

func TestPageStatusMethods(t *testing.T) {
	page := NewPage()

	// Test default status (DRAFT)
	if page.IsActive() {
		t.Error("expected IsActive to be false for DRAFT")
	}
	if page.IsInactive() {
		t.Error("expected IsInactive to be false for DRAFT")
	}

	// Test ACTIVE status
	page.SetStatus(PAGE_STATUS_ACTIVE)
	if !page.IsActive() {
		t.Error("expected IsActive to be true for ACTIVE")
	}
	if page.IsInactive() {
		t.Error("expected IsInactive to be false for ACTIVE")
	}

	// Test INACTIVE status
	page.SetStatus(PAGE_STATUS_INACTIVE)
	if page.IsActive() {
		t.Error("expected IsActive to be false for INACTIVE")
	}
	if !page.IsInactive() {
		t.Error("expected IsInactive to be true for INACTIVE")
	}

	// Test other status values
	page.SetStatus("unknown")
	if page.IsActive() {
		t.Error("expected IsActive to be false for unknown")
	}
	if page.IsInactive() {
		t.Error("expected IsInactive to be false for unknown")
	}
}

func TestPageSoftDeleteMethods(t *testing.T) {
	page := NewPage()
	if page.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false by default")
	}

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	page.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	if page.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false with future date")
	}
	if page.SoftDeletedAt() != future.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt %q, got %q", future.ToDateTimeString(carbon.UTC), page.SoftDeletedAt())
	}

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	page.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	if !page.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be true with past date")
	}
	if page.SoftDeletedAt() != past.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt %q, got %q", past.ToDateTimeString(carbon.UTC), page.SoftDeletedAt())
	}
}

func TestPageMetasMethods(t *testing.T) {
	page := NewPage()

	// Test empty metas
	metas, err := page.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Error("expected empty metas")
	}

	// Test Meta lookup on empty metas
	if page.Meta("nonexistent") != "" {
		t.Error("expected empty Meta for nonexistent key")
	}

	// Test SetMetas
	err = page.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metas, err = page.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metas["layout"] != "main" {
		t.Errorf("expected layout %q, got %q", "main", metas["layout"])
	}
	if metas["theme"] != "dark" {
		t.Errorf("expected theme %q, got %q", "dark", metas["theme"])
	}

	// Test Meta lookup
	if page.Meta("layout") != "main" {
		t.Errorf("expected layout %q", "main")
	}
	if page.Meta("theme") != "dark" {
		t.Errorf("expected theme %q", "dark")
	}
	if page.Meta("nonexistent") != "" {
		t.Error("expected empty Meta for nonexistent key")
	}

	// Test SetMeta
	err = page.SetMeta("newkey", "newvalue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Meta("newkey") != "newvalue" {
		t.Errorf("expected newkey %q", "newvalue")
	}

	// Test UpsertMetas
	err = page.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Meta("layout") != "sidebar" { // Updated
		t.Errorf("expected updated layout %q", "sidebar")
	}
	if page.Meta("theme") != "dark" { // Preserved
		t.Errorf("expected preserved theme %q", "dark")
	}
	if page.Meta("newkey") != "newvalue" { // Preserved
		t.Errorf("expected preserved newkey %q", "newvalue")
	}
	if page.Meta("color") != "blue" { // Added
		t.Errorf("expected added color %q", "blue")
	}
}

func TestPageMiddlewaresMethods(t *testing.T) {
	page := NewPage()

	// Test default middlewares
	if len(page.MiddlewaresBefore()) != 0 {
		t.Error("expected empty MiddlewaresBefore")
	}
	if len(page.MiddlewaresAfter()) != 0 {
		t.Error("expected empty MiddlewaresAfter")
	}

	// Test SetMiddlewaresBefore
	before := []string{"auth", "csrf"}
	page.SetMiddlewaresBefore(before)
	if len(page.MiddlewaresBefore()) != len(before) {
		t.Errorf("expected MiddlewaresBefore length %d, got %d", len(before), len(page.MiddlewaresBefore()))
	}
	for i, m := range before {
		if page.MiddlewaresBefore()[i] != m {
			t.Errorf("expected MiddlewaresBefore[%d] %q, got %q", i, m, page.MiddlewaresBefore()[i])
		}
	}

	// Test SetMiddlewaresAfter
	after := []string{"log", "cache"}
	page.SetMiddlewaresAfter(after)
	if len(page.MiddlewaresAfter()) != len(after) {
		t.Errorf("expected MiddlewaresAfter length %d, got %d", len(after), len(page.MiddlewaresAfter()))
	}
	for i, m := range after {
		if page.MiddlewaresAfter()[i] != m {
			t.Errorf("expected MiddlewaresAfter[%d] %q, got %q", i, m, page.MiddlewaresAfter()[i])
		}
	}

	// Test empty middlewares
	page.SetMiddlewaresBefore([]string{})
	page.SetMiddlewaresAfter([]string{})
	if len(page.MiddlewaresBefore()) != 0 {
		t.Error("expected empty MiddlewaresBefore")
	}
	if len(page.MiddlewaresAfter()) != 0 {
		t.Error("expected empty MiddlewaresAfter")
	}
}

func TestPageCreatedAtMethods(t *testing.T) {
	page := NewPage()

	// Test default CreatedAt
	createdAt := page.CreatedAt()
	if createdAt == "" {
		t.Error("expected non-empty CreatedAt")
	}

	createdAtCarbon := page.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if createdAt != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAt to match CreatedAtCarbon")
	}

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	page.SetCreatedAt(testDate)
	if page.CreatedAt() != testDate {
		t.Errorf("expected CreatedAt %q, got %q", testDate, page.CreatedAt())
	}

	createdAtCarbon = page.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if testDate != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAtCarbon to match test date")
	}
}

func TestPageUpdatedAtMethods(t *testing.T) {
	page := NewPage()

	// Test default UpdatedAt
	updatedAt := page.UpdatedAt()
	if updatedAt == "" {
		t.Error("expected non-empty UpdatedAt")
	}

	updatedAtCarbon := page.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if updatedAt != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAt to match UpdatedAtCarbon")
	}

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	page.SetUpdatedAt(testDate)
	if page.UpdatedAt() != testDate {
		t.Errorf("expected UpdatedAt %q, got %q", testDate, page.UpdatedAt())
	}

	updatedAtCarbon = page.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if testDate != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAtCarbon to match test date")
	}
}

func TestPageSoftDeletedAtMethods(t *testing.T) {
	page := NewPage()

	// Test default SoftDeletedAt
	softDeletedAt := page.SoftDeletedAt()
	if softDeletedAt != sb.MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", sb.MAX_DATETIME, softDeletedAt)
	}

	softDeletedAtCarbon := page.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if softDeletedAt != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected SoftDeletedAt to match SoftDeletedAtCarbon")
	}

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	page.SetSoftDeletedAt(testDate)
	if page.SoftDeletedAt() != testDate {
		t.Errorf("expected SoftDeletedAt %q, got %q", testDate, page.SoftDeletedAt())
	}

	softDeletedAtCarbon = page.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if testDate != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected SoftDeletedAtCarbon to match test date")
	}
}

func TestPageIDMethods(t *testing.T) {
	page := NewPage()

	// Test default ID
	id := page.ID()
	if id == "" {
		t.Error("expected non-empty ID")
	}

	// Test SetID
	newID := "test-page-id-123"
	page.SetID(newID)
	if page.ID() != newID {
		t.Errorf("expected ID %q, got %q", newID, page.ID())
	}
}

func TestPageAliasMethods(t *testing.T) {
	page := NewPage()

	// Test default alias
	if page.Alias() != "" {
		t.Error("expected empty Alias")
	}

	// Test SetAlias
	alias := "test-page-alias"
	page.SetAlias(alias)
	if page.Alias() != alias {
		t.Errorf("expected Alias %q, got %q", alias, page.Alias())
	}
}

func TestPageCanonicalUrlMethods(t *testing.T) {
	page := NewPage()

	// Test default canonical URL
	if page.CanonicalUrl() != "" {
		t.Error("expected empty CanonicalUrl")
	}

	// Test SetCanonicalUrl
	canonicalUrl := "https://example.com/canonical"
	page.SetCanonicalUrl(canonicalUrl)
	if page.CanonicalUrl() != canonicalUrl {
		t.Errorf("expected CanonicalUrl %q, got %q", canonicalUrl, page.CanonicalUrl())
	}
}

func TestPageContentMethods(t *testing.T) {
	page := NewPage()

	// Test default content
	if page.Content() != "" {
		t.Error("expected empty Content")
	}

	// Test SetContent
	content := "This is page content"
	page.SetContent(content)
	if page.Content() != content {
		t.Errorf("expected Content %q, got %q", content, page.Content())
	}
}

func TestPageEditorMethods(t *testing.T) {
	page := NewPage()

	// Test default editor
	if page.Editor() != "" {
		t.Error("expected empty Editor")
	}

	// Test SetEditor
	editor := "test-editor"
	page.SetEditor(editor)
	if page.Editor() != editor {
		t.Errorf("expected Editor %q, got %q", editor, page.Editor())
	}
}

func TestPageHandleMethods(t *testing.T) {
	page := NewPage()

	// Test default handle
	if page.Handle() != "" {
		t.Error("expected empty Handle")
	}

	// Test SetHandle
	handle := "test-page-handle"
	page.SetHandle(handle)
	if page.Handle() != handle {
		t.Errorf("expected Handle %q, got %q", handle, page.Handle())
	}
}

func TestPageMemoMethods(t *testing.T) {
	page := NewPage()

	// Test default memo
	if page.Memo() != "" {
		t.Error("expected empty Memo")
	}

	// Test SetMemo
	memo := "This is a page memo"
	page.SetMemo(memo)
	if page.Memo() != memo {
		t.Errorf("expected Memo %q, got %q", memo, page.Memo())
	}
}

func TestPageMetaDescriptionMethods(t *testing.T) {
	page := NewPage()

	// Test default meta description
	if page.MetaDescription() != "" {
		t.Error("expected empty MetaDescription")
	}

	// Test SetMetaDescription
	metaDescription := "This is a meta description"
	page.SetMetaDescription(metaDescription)
	if page.MetaDescription() != metaDescription {
		t.Errorf("expected MetaDescription %q, got %q", metaDescription, page.MetaDescription())
	}
}

func TestPageMetaKeywordsMethods(t *testing.T) {
	page := NewPage()

	// Test default meta keywords
	if page.MetaKeywords() != "" {
		t.Error("expected empty MetaKeywords")
	}

	// Test SetMetaKeywords
	metaKeywords := "keyword1, keyword2, keyword3"
	page.SetMetaKeywords(metaKeywords)
	if page.MetaKeywords() != metaKeywords {
		t.Errorf("expected MetaKeywords %q, got %q", metaKeywords, page.MetaKeywords())
	}
}

func TestPageMetaRobotsMethods(t *testing.T) {
	page := NewPage()

	// Test default meta robots
	if page.MetaRobots() != "" {
		t.Error("expected empty MetaRobots")
	}

	// Test SetMetaRobots
	metaRobots := "noindex, nofollow"
	page.SetMetaRobots(metaRobots)
	if page.MetaRobots() != metaRobots {
		t.Errorf("expected MetaRobots %q, got %q", metaRobots, page.MetaRobots())
	}
}

func TestPageNameMethods(t *testing.T) {
	page := NewPage()

	// Test default name
	if page.Name() != "" {
		t.Error("expected empty Name")
	}

	// Test SetName
	name := "Test Page Name"
	page.SetName(name)
	if page.Name() != name {
		t.Errorf("expected Name %q, got %q", name, page.Name())
	}
}

func TestPageSiteIDMethods(t *testing.T) {
	page := NewPage()

	// Test default site ID
	if page.SiteID() != "" {
		t.Error("expected empty SiteID")
	}

	// Test SetSiteID
	siteID := "test-site-id"
	page.SetSiteID(siteID)
	if page.SiteID() != siteID {
		t.Errorf("expected SiteID %q, got %q", siteID, page.SiteID())
	}
}

func TestPageTemplateIDMethods(t *testing.T) {
	page := NewPage()

	// Test default template ID
	if page.TemplateID() != "" {
		t.Error("expected empty TemplateID")
	}

	// Test SetTemplateID
	templateID := "test-template-id"
	page.SetTemplateID(templateID)
	if page.TemplateID() != templateID {
		t.Errorf("expected TemplateID %q, got %q", templateID, page.TemplateID())
	}
}

func TestPageTitleMethods(t *testing.T) {
	page := NewPage()

	// Test default title
	if page.Title() != "" {
		t.Error("expected empty Title")
	}

	// Test SetTitle
	title := "Test Page Title"
	page.SetTitle(title)
	if page.Title() != title {
		t.Errorf("expected Title %q, got %q", title, page.Title())
	}
}

func TestPageStatusSettersAndGetters(t *testing.T) {
	page := NewPage()

	// Test default status
	if page.Status() != PAGE_STATUS_DRAFT {
		t.Errorf("expected Status %q, got %q", PAGE_STATUS_DRAFT, page.Status())
	}

	// Test SetStatus
	page.SetStatus(PAGE_STATUS_ACTIVE)
	if page.Status() != PAGE_STATUS_ACTIVE {
		t.Errorf("expected Status %q, got %q", PAGE_STATUS_ACTIVE, page.Status())
	}

	page.SetStatus(PAGE_STATUS_INACTIVE)
	if page.Status() != PAGE_STATUS_INACTIVE {
		t.Errorf("expected Status %q, got %q", PAGE_STATUS_INACTIVE, page.Status())
	}

	page.SetStatus("custom-status")
	if page.Status() != "custom-status" {
		t.Errorf("expected Status %q, got %q", "custom-status", page.Status())
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if versionedJSON == "" {
		t.Error("expected non-empty versionedJSON")
	}

	// Parse the JSON to verify it contains expected fields
	var versionedData map[string]string
	err = json.Unmarshal([]byte(versionedJSON), &versionedData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that expected fields are present
	if versionedData[COLUMN_ALIAS] != "test-alias" {
		t.Errorf("expected alias %q, got %q", "test-alias", versionedData[COLUMN_ALIAS])
	}
	if versionedData[COLUMN_CANONICAL_URL] != "https://example.com/canonical" {
		t.Errorf("expected canonical_url %q, got %q", "https://example.com/canonical", versionedData[COLUMN_CANONICAL_URL])
	}
	if versionedData[COLUMN_CONTENT] != "Test content" {
		t.Errorf("expected content %q, got %q", "Test content", versionedData[COLUMN_CONTENT])
	}
	if versionedData[COLUMN_EDITOR] != "test-editor" {
		t.Errorf("expected editor %q, got %q", "test-editor", versionedData[COLUMN_EDITOR])
	}
	if versionedData[COLUMN_HANDLE] != "test-handle" {
		t.Errorf("expected handle %q, got %q", "test-handle", versionedData[COLUMN_HANDLE])
	}
	if versionedData[COLUMN_MEMO] != "test-memo" {
		t.Errorf("expected memo %q, got %q", "test-memo", versionedData[COLUMN_MEMO])
	}
	if versionedData[COLUMN_META_DESCRIPTION] != "Test meta description" {
		t.Errorf("expected meta_description %q, got %q", "Test meta description", versionedData[COLUMN_META_DESCRIPTION])
	}
	if versionedData[COLUMN_META_KEYWORDS] != "test, keywords" {
		t.Errorf("expected meta_keywords %q, got %q", "test, keywords", versionedData[COLUMN_META_KEYWORDS])
	}
	if versionedData[COLUMN_META_ROBOTS] != "index, follow" {
		t.Errorf("expected meta_robots %q, got %q", "index, follow", versionedData[COLUMN_META_ROBOTS])
	}
	if versionedData[COLUMN_NAME] != "Test Page" {
		t.Errorf("expected name %q, got %q", "Test Page", versionedData[COLUMN_NAME])
	}
	if versionedData[COLUMN_SITE_ID] != "test-site" {
		t.Errorf("expected site_id %q, got %q", "test-site", versionedData[COLUMN_SITE_ID])
	}
	if versionedData[COLUMN_TEMPLATE_ID] != "test-template" {
		t.Errorf("expected template_id %q, got %q", "test-template", versionedData[COLUMN_TEMPLATE_ID])
	}
	if versionedData[COLUMN_TITLE] != "Test Page Title" {
		t.Errorf("expected title %q, got %q", "Test Page Title", versionedData[COLUMN_TITLE])
	}
	if versionedData[COLUMN_STATUS] != PAGE_STATUS_ACTIVE {
		t.Errorf("expected status %q, got %q", PAGE_STATUS_ACTIVE, versionedData[COLUMN_STATUS])
	}

	// Check that timestamps and soft delete fields are excluded
	_, hasCreatedAt := versionedData[COLUMN_CREATED_AT]
	_, hasUpdatedAt := versionedData[COLUMN_UPDATED_AT]
	_, hasSoftDeletedAt := versionedData[COLUMN_SOFT_DELETED_AT]
	if hasCreatedAt {
		t.Error("expected CreatedAt to be excluded")
	}
	if hasUpdatedAt {
		t.Error("expected UpdatedAt to be excluded")
	}
	if hasSoftDeletedAt {
		t.Error("expected SoftDeletedAt to be excluded")
	}
}
