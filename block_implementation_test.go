package cmsstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewBlockDefaults(t *testing.T) {
	block := NewBlock()

	if block.ID() == "" {
		t.Errorf("expected non-empty ID")
	}
	if block.CreatedAt() == "" {
		t.Errorf("expected non-empty CreatedAt")
	}
	if block.UpdatedAt() == "" {
		t.Errorf("expected non-empty UpdatedAt")
	}
	if block.Status() != BLOCK_STATUS_DRAFT {
		t.Errorf("expected status %q, got %q", BLOCK_STATUS_DRAFT, block.Status())
	}
	if block.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", sb.MAX_DATETIME, block.SoftDeletedAt())
	}
	if block.IsSoftDeleted() {
		t.Errorf("expected IsSoftDeleted to be false")
	}

	metas, err := block.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("expected empty metas")
	}

	createdCarbon := block.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Fatalf("expected non-nil CreatedAtCarbon")
	}
	if block.CreatedAt() != createdCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected CreatedAt to match CreatedAtCarbon")
	}

	updatedCarbon := block.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Fatalf("expected non-nil UpdatedAtCarbon")
	}
	if block.UpdatedAt() != updatedCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected UpdatedAt to match UpdatedAtCarbon")
	}

	softDeletedCarbon := block.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Fatalf("expected non-nil SoftDeletedAtCarbon")
	}
	if !softDeletedCarbon.Gte(carbon.Now(carbon.UTC)) {
		t.Errorf("expected SoftDeletedAtCarbon to be in the future")
	}
}

func TestBlockGetterMethods(t *testing.T) {
	block := NewBlock()

	// Test default values
	if block.Content() != "" {
		t.Errorf("expected empty Content")
	}
	if block.Editor() != "" {
		t.Errorf("expected empty Editor")
	}
	if block.Handle() != "" {
		t.Errorf("expected empty Handle")
	}
	if block.Memo() != "" {
		t.Errorf("expected empty Memo")
	}
	if block.Name() != "" {
		t.Errorf("expected empty Name")
	}
	if block.PageID() != "" {
		t.Errorf("expected empty PageID")
	}
	if block.ParentID() != "" {
		t.Errorf("expected empty ParentID")
	}
	if block.Sequence() != "0" {
		t.Errorf("expected Sequence %q, got %q", "0", block.Sequence())
	}
	if block.SequenceInt() != 0 {
		t.Errorf("expected SequenceInt 0, got %d", block.SequenceInt())
	}
	if block.SiteID() != "" {
		t.Errorf("expected empty SiteID")
	}
	if block.TemplateID() != "" {
		t.Errorf("expected empty TemplateID")
	}
	if block.Type() != BLOCK_TYPE_HTML {
		t.Errorf("expected Type %q, got %q", BLOCK_TYPE_HTML, block.Type())
	}
}

func TestBlockStatusMethods(t *testing.T) {
	block := NewBlock()

	// Test default status (DRAFT)
	if block.IsActive() {
		t.Errorf("expected IsActive to be false for DRAFT")
	}
	if block.IsInactive() {
		t.Errorf("expected IsInactive to be false for DRAFT")
	}

	// Test ACTIVE status
	block.SetStatus(BLOCK_STATUS_ACTIVE)
	if !block.IsActive() {
		t.Errorf("expected IsActive to be true for ACTIVE")
	}
	if block.IsInactive() {
		t.Errorf("expected IsInactive to be false for ACTIVE")
	}

	// Test INACTIVE status
	block.SetStatus(BLOCK_STATUS_INACTIVE)
	if block.IsActive() {
		t.Errorf("expected IsActive to be false for INACTIVE")
	}
	if !block.IsInactive() {
		t.Errorf("expected IsInactive to be true for INACTIVE")
	}

	// Test other status values
	block.SetStatus("unknown")
	if block.IsActive() {
		t.Errorf("expected IsActive to be false for unknown")
	}
	if block.IsInactive() {
		t.Errorf("expected IsInactive to be false for unknown")
	}
}

func TestBlockSoftDeleteMethods(t *testing.T) {
	block := NewBlock()
	if block.IsSoftDeleted() {
		t.Errorf("expected IsSoftDeleted to be false by default")
	}

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	block.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	if block.IsSoftDeleted() {
		t.Errorf("expected IsSoftDeleted to be false with future date")
	}
	if block.SoftDeletedAt() != future.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt %q, got %q", future.ToDateTimeString(carbon.UTC), block.SoftDeletedAt())
	}

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	block.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	if !block.IsSoftDeleted() {
		t.Errorf("expected IsSoftDeleted to be true with past date")
	}
	if block.SoftDeletedAt() != past.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt %q, got %q", past.ToDateTimeString(carbon.UTC), block.SoftDeletedAt())
	}
}

func TestBlockMetasMethods(t *testing.T) {
	block := NewBlock()

	// Test empty metas
	metas, err := block.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("expected empty metas")
	}

	// Test Meta lookup on empty metas
	if block.Meta("nonexistent") != "" {
		t.Errorf("expected empty Meta for nonexistent key")
	}

	// Test SetMetas
	err = block.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metas, err = block.Metas()
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
	if block.Meta("layout") != "main" {
		t.Errorf("expected layout %q", "main")
	}
	if block.Meta("theme") != "dark" {
		t.Errorf("expected theme %q", "dark")
	}
	if block.Meta("nonexistent") != "" {
		t.Errorf("expected empty Meta for nonexistent key")
	}

	// Test SetMeta
	err = block.SetMeta("newkey", "newvalue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Meta("newkey") != "newvalue" {
		t.Errorf("expected newkey %q", "newvalue")
	}

	// Test UpsertMetas
	err = block.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Meta("layout") != "sidebar" { // Updated
		t.Errorf("expected updated layout %q", "sidebar")
	}
	if block.Meta("theme") != "dark" { // Preserved
		t.Errorf("expected preserved theme %q", "dark")
	}
	if block.Meta("newkey") != "newvalue" { // Preserved
		t.Errorf("expected preserved newkey %q", "newvalue")
	}
	if block.Meta("color") != "blue" { // Added
		t.Errorf("expected added color %q", "blue")
	}
}

func TestBlockSequenceMethods(t *testing.T) {
	block := NewBlock()

	// Test default sequence
	if block.Sequence() != "0" {
		t.Errorf("expected Sequence %q, got %q", "0", block.Sequence())
	}
	if block.SequenceInt() != 0 {
		t.Errorf("expected SequenceInt 0, got %d", block.SequenceInt())
	}

	// Test SetSequenceInt
	block.SetSequenceInt(42)
	if block.Sequence() != "42" {
		t.Errorf("expected Sequence %q, got %q", "42", block.Sequence())
	}
	if block.SequenceInt() != 42 {
		t.Errorf("expected SequenceInt 42, got %d", block.SequenceInt())
	}

	// Test SetSequence
	block.SetSequence("123")
	if block.Sequence() != "123" {
		t.Errorf("expected Sequence %q, got %q", "123", block.Sequence())
	}
	if block.SequenceInt() != 123 {
		t.Errorf("expected SequenceInt 123, got %d", block.SequenceInt())
	}

	// Test invalid sequence
	block.SetSequence("invalid")
	if block.Sequence() != "invalid" {
		t.Errorf("expected Sequence %q, got %q", "invalid", block.Sequence())
	}
	if block.SequenceInt() != 0 { // Should default to 0 for invalid
		t.Errorf("expected SequenceInt 0 for invalid, got %d", block.SequenceInt())
	}
}

func TestBlockCreatedAtMethods(t *testing.T) {
	block := NewBlock()

	// Test default CreatedAt
	createdAt := block.CreatedAt()
	if createdAt == "" {
		t.Errorf("expected non-empty CreatedAt")
	}

	createdAtCarbon := block.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatalf("expected non-nil CreatedAtCarbon")
	}
	if createdAt != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected CreatedAt to match CreatedAtCarbon")
	}

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	block.SetCreatedAt(testDate)
	if block.CreatedAt() != testDate {
		t.Errorf("expected CreatedAt %q, got %q", testDate, block.CreatedAt())
	}

	createdAtCarbon = block.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatalf("expected non-nil CreatedAtCarbon")
	}
	if testDate != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected CreatedAtCarbon to match test date")
	}
}

func TestBlockUpdatedAtMethods(t *testing.T) {
	block := NewBlock()

	// Test default UpdatedAt
	updatedAt := block.UpdatedAt()
	if updatedAt == "" {
		t.Errorf("expected non-empty UpdatedAt")
	}

	updatedAtCarbon := block.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatalf("expected non-nil UpdatedAtCarbon")
	}
	if updatedAt != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected UpdatedAt to match UpdatedAtCarbon")
	}

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	block.SetUpdatedAt(testDate)
	if block.UpdatedAt() != testDate {
		t.Errorf("expected UpdatedAt %q, got %q", testDate, block.UpdatedAt())
	}

	updatedAtCarbon = block.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatalf("expected non-nil UpdatedAtCarbon")
	}
	if testDate != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected UpdatedAtCarbon to match test date")
	}
}

func TestBlockSoftDeletedAtMethods(t *testing.T) {
	block := NewBlock()

	// Test default SoftDeletedAt
	softDeletedAt := block.SoftDeletedAt()
	if softDeletedAt != sb.MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", sb.MAX_DATETIME, softDeletedAt)
	}

	softDeletedAtCarbon := block.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Fatalf("expected non-nil SoftDeletedAtCarbon")
	}
	if softDeletedAt != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt to match SoftDeletedAtCarbon")
	}

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	block.SetSoftDeletedAt(testDate)
	if block.SoftDeletedAt() != testDate {
		t.Errorf("expected SoftDeletedAt %q, got %q", testDate, block.SoftDeletedAt())
	}

	softDeletedAtCarbon = block.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Fatalf("expected non-nil SoftDeletedAtCarbon")
	}
	if testDate != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAtCarbon to match test date")
	}
}

func TestBlockIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default ID
	id := block.ID()
	if id == "" {
		t.Errorf("expected non-empty ID")
	}

	// Test SetID
	newID := "test-block-id-123"
	block.SetID(newID)
	if block.ID() != newID {
		t.Errorf("expected ID %q, got %q", newID, block.ID())
	}
}

func TestBlockContentMethods(t *testing.T) {
	block := NewBlock()

	// Test default content
	if block.Content() != "" {
		t.Errorf("expected empty Content")
	}

	// Test SetContent
	content := "This is block content"
	block.SetContent(content)
	if block.Content() != content {
		t.Errorf("expected Content %q, got %q", content, block.Content())
	}
}

func TestBlockEditorMethods(t *testing.T) {
	block := NewBlock()

	// Test default editor
	if block.Editor() != "" {
		t.Errorf("expected empty Editor")
	}

	// Test SetEditor
	editor := "test-editor"
	block.SetEditor(editor)
	if block.Editor() != editor {
		t.Errorf("expected Editor %q, got %q", editor, block.Editor())
	}
}

func TestBlockHandleMethods(t *testing.T) {
	block := NewBlock()

	// Test default handle
	if block.Handle() != "" {
		t.Errorf("expected empty Handle")
	}

	// Test SetHandle
	handle := "test-block-handle"
	block.SetHandle(handle)
	if block.Handle() != handle {
		t.Errorf("expected Handle %q, got %q", handle, block.Handle())
	}
}

func TestBlockMemoMethods(t *testing.T) {
	block := NewBlock()

	// Test default memo
	if block.Memo() != "" {
		t.Errorf("expected empty Memo")
	}

	// Test SetMemo
	memo := "This is a block memo"
	block.SetMemo(memo)
	if block.Memo() != memo {
		t.Errorf("expected Memo %q, got %q", memo, block.Memo())
	}
}

func TestBlockNameMethods(t *testing.T) {
	block := NewBlock()

	// Test default name
	if block.Name() != "" {
		t.Errorf("expected empty Name")
	}

	// Test SetName
	name := "Test Block Name"
	block.SetName(name)
	if block.Name() != name {
		t.Errorf("expected Name %q, got %q", name, block.Name())
	}
}

func TestBlockPageIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default page ID
	if block.PageID() != "" {
		t.Errorf("expected empty PageID")
	}

	// Test SetPageID
	pageID := "test-page-id"
	block.SetPageID(pageID)
	if block.PageID() != pageID {
		t.Errorf("expected PageID %q, got %q", pageID, block.PageID())
	}
}

func TestBlockParentIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default parent ID
	if block.ParentID() != "" {
		t.Errorf("expected empty ParentID")
	}

	// Test SetParentID
	parentID := "test-parent-id"
	block.SetParentID(parentID)
	if block.ParentID() != parentID {
		t.Errorf("expected ParentID %q, got %q", parentID, block.ParentID())
	}
}

func TestBlockSiteIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default site ID
	if block.SiteID() != "" {
		t.Errorf("expected empty SiteID")
	}

	// Test SetSiteID
	siteID := "test-site-id"
	block.SetSiteID(siteID)
	if block.SiteID() != siteID {
		t.Errorf("expected SiteID %q, got %q", siteID, block.SiteID())
	}
}

func TestBlockTemplateIDMethods(t *testing.T) {
	block := NewBlock()

	// Test default template ID
	if block.TemplateID() != "" {
		t.Errorf("expected empty TemplateID")
	}

	// Test SetTemplateID
	templateID := "test-template-id"
	block.SetTemplateID(templateID)
	if block.TemplateID() != templateID {
		t.Errorf("expected TemplateID %q, got %q", templateID, block.TemplateID())
	}
}

func TestBlockTypeMethods(t *testing.T) {
	block := NewBlock()

	// Test default type (HTML)
	if block.Type() != BLOCK_TYPE_HTML {
		t.Errorf("expected Type %q, got %q", BLOCK_TYPE_HTML, block.Type())
	}

	// Test SetType
	blockType := "text"
	block.SetType(blockType)
	if block.Type() != blockType {
		t.Errorf("expected Type %q, got %q", blockType, block.Type())
	}
}

func TestBlockStatusSettersAndGetters(t *testing.T) {
	block := NewBlock()

	// Test default status
	if block.Status() != BLOCK_STATUS_DRAFT {
		t.Errorf("expected Status %q, got %q", BLOCK_STATUS_DRAFT, block.Status())
	}

	// Test SetStatus
	block.SetStatus(BLOCK_STATUS_ACTIVE)
	if block.Status() != BLOCK_STATUS_ACTIVE {
		t.Errorf("expected Status %q, got %q", BLOCK_STATUS_ACTIVE, block.Status())
	}

	block.SetStatus(BLOCK_STATUS_INACTIVE)
	if block.Status() != BLOCK_STATUS_INACTIVE {
		t.Errorf("expected Status %q, got %q", BLOCK_STATUS_INACTIVE, block.Status())
	}

	block.SetStatus("custom-status")
	if block.Status() != "custom-status" {
		t.Errorf("expected Status %q, got %q", "custom-status", block.Status())
	}
}
