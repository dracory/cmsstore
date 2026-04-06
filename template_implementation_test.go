package cmsstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewTemplateDefaults(t *testing.T) {
	template := NewTemplate()

	// Test default values
	if len(template.ID()) == 0 {
		t.Error("Expected ID to be non-empty")
	}
	if len(template.CreatedAt()) == 0 {
		t.Error("Expected CreatedAt to be non-empty")
	}
	if len(template.UpdatedAt()) == 0 {
		t.Error("Expected UpdatedAt to be non-empty")
	}
	if template.Status() != TEMPLATE_STATUS_DRAFT {
		t.Errorf("Expected Status %s, got %s", TEMPLATE_STATUS_DRAFT, template.Status())
	}
	if template.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("Expected SoftDeletedAt %s, got %s", sb.MAX_DATETIME, template.SoftDeletedAt())
	}
	if template.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false")
	}

	metas, err := template.Metas()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("Expected empty metas, got %v", metas)
	}

	createdCarbon := template.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Error("Expected CreatedAtCarbon to be non-nil")
	}
	if template.CreatedAt() != createdCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected CreatedAt %s, got %s", template.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))
	}

	updatedCarbon := template.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Error("Expected UpdatedAtCarbon to be non-nil")
	}
	if template.UpdatedAt() != updatedCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected UpdatedAt %s, got %s", template.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))
	}

	softDeletedCarbon := template.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Error("Expected SoftDeletedAtCarbon to be non-nil")
	}
	if !softDeletedCarbon.Gte(carbon.Now(carbon.UTC)) {
		t.Error("Expected SoftDeletedAtCarbon to be greater than or equal to now")
	}
}

func TestTemplateGetterMethods(t *testing.T) {
	template := NewTemplate()

	// Test default values
	if template.Content() != "" {
		t.Errorf("Expected empty Content, got %s", template.Content())
	}
	if template.Editor() != "" {
		t.Errorf("Expected empty Editor, got %s", template.Editor())
	}
	if template.Handle() != "" {
		t.Errorf("Expected empty Handle, got %s", template.Handle())
	}
	if template.Memo() != "" {
		t.Errorf("Expected empty Memo, got %s", template.Memo())
	}
	if template.Name() != "" {
		t.Errorf("Expected empty Name, got %s", template.Name())
	}
	if template.SiteID() != "" {
		t.Errorf("Expected empty SiteID, got %s", template.SiteID())
	}
}

func TestTemplateStatusMethods(t *testing.T) {
	template := NewTemplate()

	// Test default status (DRAFT)
	if template.IsActive() {
		t.Error("Expected IsActive to be false")
	}
	if template.IsInactive() {
		t.Error("Expected IsInactive to be false")
	}

	// Test ACTIVE status
	template.SetStatus(TEMPLATE_STATUS_ACTIVE)
	if !template.IsActive() {
		t.Error("Expected IsActive to be true")
	}
	if template.IsInactive() {
		t.Error("Expected IsInactive to be false")
	}

	// Test INACTIVE status
	template.SetStatus(TEMPLATE_STATUS_INACTIVE)
	if template.IsActive() {
		t.Error("Expected IsActive to be false")
	}
	if !template.IsInactive() {
		t.Error("Expected IsInactive to be true")
	}

	// Test other status values
	template.SetStatus("unknown")
	if template.IsActive() {
		t.Error("Expected IsActive to be false")
	}
	if template.IsInactive() {
		t.Error("Expected IsInactive to be false")
	}
}

func TestTemplateSoftDeleteMethods(t *testing.T) {
	template := NewTemplate()
	if template.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false")
	}

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	template.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	if template.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false")
	}
	if template.SoftDeletedAt() != future.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", future.ToDateTimeString(carbon.UTC), template.SoftDeletedAt())
	}

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	template.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	if !template.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be true")
	}
	if template.SoftDeletedAt() != past.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", past.ToDateTimeString(carbon.UTC), template.SoftDeletedAt())
	}
}

func TestTemplateMetasMethods(t *testing.T) {
	template := NewTemplate()

	// Test empty metas
	metas, err := template.Metas()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("Expected empty metas, got %v", metas)
	}

	// Test Meta lookup on empty metas
	if template.Meta("nonexistent") != "" {
		t.Errorf("Expected empty Meta, got %s", template.Meta("nonexistent"))
	}

	// Test SetMetas
	err = template.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	metas, err = template.Metas()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if metas["layout"] != "main" {
		t.Errorf("Expected layout 'main', got %s", metas["layout"])
	}
	if metas["theme"] != "dark" {
		t.Errorf("Expected theme 'dark', got %s", metas["theme"])
	}

	// Test Meta lookup
	if template.Meta("layout") != "main" {
		t.Errorf("Expected layout 'main', got %s", template.Meta("layout"))
	}
	if template.Meta("theme") != "dark" {
		t.Errorf("Expected theme 'dark', got %s", template.Meta("theme"))
	}
	if template.Meta("nonexistent") != "" {
		t.Errorf("Expected empty Meta, got %s", template.Meta("nonexistent"))
	}

	// Test SetMeta
	err = template.SetMeta("newkey", "newvalue")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if template.Meta("newkey") != "newvalue" {
		t.Errorf("Expected newkey 'newvalue', got %s", template.Meta("newkey"))
	}

	// Test UpsertMetas
	err = template.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if template.Meta("layout") != "sidebar" { // Updated
		t.Errorf("Expected layout 'sidebar', got %s", template.Meta("layout"))
	}
	if template.Meta("theme") != "dark" { // Preserved
		t.Errorf("Expected theme 'dark', got %s", template.Meta("theme"))
	}
	if template.Meta("newkey") != "newvalue" { // Preserved
		t.Errorf("Expected newkey 'newvalue', got %s", template.Meta("newkey"))
	}
	if template.Meta("color") != "blue" { // Added
		t.Errorf("Expected color 'blue', got %s", template.Meta("color"))
	}
}

func TestTemplateCreatedAtMethods(t *testing.T) {
	template := NewTemplate()

	// Test default CreatedAt
	createdAt := template.CreatedAt()
	if len(createdAt) == 0 {
		t.Error("Expected CreatedAt to be non-empty")
	}

	createdAtCarbon := template.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Error("Expected CreatedAtCarbon to be non-nil")
	}
	if createdAt != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected CreatedAt %s, got %s", createdAt, createdAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	template.SetCreatedAt(testDate)
	if template.CreatedAt() != testDate {
		t.Errorf("Expected CreatedAt %s, got %s", testDate, template.CreatedAt())
	}

	createdAtCarbon = template.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Error("Expected CreatedAtCarbon to be non-nil")
	}
	if testDate != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected CreatedAt %s, got %s", testDate, createdAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestTemplateUpdatedAtMethods(t *testing.T) {
	template := NewTemplate()

	// Test default UpdatedAt
	updatedAt := template.UpdatedAt()
	if len(updatedAt) == 0 {
		t.Error("Expected UpdatedAt to be non-empty")
	}

	updatedAtCarbon := template.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Error("Expected UpdatedAtCarbon to be non-nil")
	}
	if updatedAt != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected UpdatedAt %s, got %s", updatedAt, updatedAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	template.SetUpdatedAt(testDate)
	if template.UpdatedAt() != testDate {
		t.Errorf("Expected UpdatedAt %s, got %s", testDate, template.UpdatedAt())
	}

	updatedAtCarbon = template.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Error("Expected UpdatedAtCarbon to be non-nil")
	}
	if testDate != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected UpdatedAt %s, got %s", testDate, updatedAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestTemplateSoftDeletedAtMethods(t *testing.T) {
	template := NewTemplate()

	// Test default SoftDeletedAt
	softDeletedAt := template.SoftDeletedAt()
	if softDeletedAt != sb.MAX_DATETIME {
		t.Errorf("Expected SoftDeletedAt %s, got %s", sb.MAX_DATETIME, softDeletedAt)
	}

	softDeletedAtCarbon := template.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Error("Expected SoftDeletedAtCarbon to be non-nil")
	}
	if softDeletedAt != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", softDeletedAt, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	template.SetSoftDeletedAt(testDate)
	if template.SoftDeletedAt() != testDate {
		t.Errorf("Expected SoftDeletedAt %s, got %s", testDate, template.SoftDeletedAt())
	}

	softDeletedAtCarbon = template.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Error("Expected SoftDeletedAtCarbon to be non-nil")
	}
	if testDate != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", testDate, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestTemplateIDMethods(t *testing.T) {
	template := NewTemplate()

	// Test default ID
	id := template.ID()
	if len(id) == 0 {
		t.Error("Expected ID to be non-empty")
	}

	// Test SetID
	newID := "test-template-id-123"
	template.SetID(newID)
	if template.ID() != newID {
		t.Errorf("Expected ID %s, got %s", newID, template.ID())
	}
}

func TestTemplateContentMethods(t *testing.T) {
	template := NewTemplate()

	// Test default content
	if template.Content() != "" {
		t.Errorf("Expected empty Content, got %s", template.Content())
	}

	// Test SetContent
	content := "This is template content"
	template.SetContent(content)
	if template.Content() != content {
		t.Errorf("Expected Content %s, got %s", content, template.Content())
	}
}

func TestTemplateEditorMethods(t *testing.T) {
	template := NewTemplate()

	// Test default editor
	if template.Editor() != "" {
		t.Errorf("Expected empty Editor, got %s", template.Editor())
	}

	// Test SetEditor
	editor := "test-editor"
	template.SetEditor(editor)
	if template.Editor() != editor {
		t.Errorf("Expected Editor %s, got %s", editor, template.Editor())
	}
}

func TestTemplateHandleMethods(t *testing.T) {
	template := NewTemplate()

	// Test default handle
	if template.Handle() != "" {
		t.Errorf("Expected empty Handle, got %s", template.Handle())
	}

	// Test SetHandle
	handle := "test-template-handle"
	template.SetHandle(handle)
	if template.Handle() != handle {
		t.Errorf("Expected Handle %s, got %s", handle, template.Handle())
	}
}

func TestTemplateMemoMethods(t *testing.T) {
	template := NewTemplate()

	// Test default memo
	if template.Memo() != "" {
		t.Errorf("Expected empty Memo, got %s", template.Memo())
	}

	// Test SetMemo
	memo := "This is a template memo"
	template.SetMemo(memo)
	if template.Memo() != memo {
		t.Errorf("Expected Memo %s, got %s", memo, template.Memo())
	}
}

func TestTemplateNameMethods(t *testing.T) {
	template := NewTemplate()

	// Test default name
	if template.Name() != "" {
		t.Errorf("Expected empty Name, got %s", template.Name())
	}

	// Test SetName
	name := "Test Template Name"
	template.SetName(name)
	if template.Name() != name {
		t.Errorf("Expected Name %s, got %s", name, template.Name())
	}
}

func TestTemplateSiteIDMethods(t *testing.T) {
	template := NewTemplate()

	// Test default site ID
	if template.SiteID() != "" {
		t.Errorf("Expected empty SiteID, got %s", template.SiteID())
	}

	// Test SetSiteID
	siteID := "test-site-id"
	template.SetSiteID(siteID)
	if template.SiteID() != siteID {
		t.Errorf("Expected SiteID %s, got %s", siteID, template.SiteID())
	}
}

func TestTemplateStatusSettersAndGetters(t *testing.T) {
	template := NewTemplate()

	// Test default status
	if template.Status() != TEMPLATE_STATUS_DRAFT {
		t.Errorf("Expected Status %s, got %s", TEMPLATE_STATUS_DRAFT, template.Status())
	}

	// Test SetStatus
	template.SetStatus(TEMPLATE_STATUS_ACTIVE)
	if template.Status() != TEMPLATE_STATUS_ACTIVE {
		t.Errorf("Expected Status %s, got %s", TEMPLATE_STATUS_ACTIVE, template.Status())
	}

	template.SetStatus(TEMPLATE_STATUS_INACTIVE)
	if template.Status() != TEMPLATE_STATUS_INACTIVE {
		t.Errorf("Expected Status %s, got %s", TEMPLATE_STATUS_INACTIVE, template.Status())
	}

	template.SetStatus("custom-status")
	if template.Status() != "custom-status" {
		t.Errorf("Expected Status %s, got %s", "custom-status", template.Status())
	}
}
