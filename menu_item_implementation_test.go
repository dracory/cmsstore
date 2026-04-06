package cmsstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewMenuItemDefaults(t *testing.T) {
	menuItem := NewMenuItem()

	if menuItem.ID() == "" {
		t.Error("expected non-empty ID")
	}
	if menuItem.CreatedAt() == "" {
		t.Error("expected non-empty CreatedAt")
	}
	if menuItem.UpdatedAt() == "" {
		t.Error("expected non-empty UpdatedAt")
	}
	if menuItem.Status() != MENU_ITEM_STATUS_DRAFT {
		t.Errorf("expected status %q, got %q", MENU_ITEM_STATUS_DRAFT, menuItem.Status())
	}
	if menuItem.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", sb.MAX_DATETIME, menuItem.SoftDeletedAt())
	}
	if menuItem.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false")
	}

	metas, err := menuItem.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Error("expected empty metas")
	}

	createdCarbon := menuItem.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if menuItem.CreatedAt() != createdCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAt to match CreatedAtCarbon")
	}

	updatedCarbon := menuItem.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if menuItem.UpdatedAt() != updatedCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAt to match UpdatedAtCarbon")
	}

	softDeletedCarbon := menuItem.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if !softDeletedCarbon.Gte(carbon.Now(carbon.UTC)) {
		t.Error("expected SoftDeletedAtCarbon to be in the future")
	}
}

func TestMenuItemGetterMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default values
	if menuItem.Handle() != "" {
		t.Error("expected empty Handle")
	}
	if menuItem.Memo() != "" {
		t.Error("expected empty Memo")
	}
	if menuItem.MenuID() != "" {
		t.Error("expected empty MenuID")
	}
	if menuItem.Name() != "" {
		t.Error("expected empty Name")
	}
	if menuItem.PageID() != "" {
		t.Error("expected empty PageID")
	}
	if menuItem.ParentID() != "" {
		t.Error("expected empty ParentID")
	}
	if menuItem.Sequence() != "0" {
		t.Errorf("expected Sequence %q, got %q", "0", menuItem.Sequence())
	}
	if menuItem.SequenceInt() != 0 {
		t.Errorf("expected SequenceInt 0, got %d", menuItem.SequenceInt())
	}
	if menuItem.Target() != "" {
		t.Error("expected empty Target")
	}
	if menuItem.URL() != "" {
		t.Error("expected empty URL")
	}
}

func TestMenuItemStatusMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default status (DRAFT)
	if menuItem.IsActive() {
		t.Error("expected IsActive to be false for DRAFT")
	}
	if menuItem.IsInactive() {
		t.Error("expected IsInactive to be false for DRAFT")
	}

	// Test ACTIVE status
	menuItem.SetStatus(MENU_ITEM_STATUS_ACTIVE)
	if !menuItem.IsActive() {
		t.Error("expected IsActive to be true for ACTIVE")
	}
	if menuItem.IsInactive() {
		t.Error("expected IsInactive to be false for ACTIVE")
	}

	// Test INACTIVE status
	menuItem.SetStatus(MENU_ITEM_STATUS_INACTIVE)
	if menuItem.IsActive() {
		t.Error("expected IsActive to be false for INACTIVE")
	}
	if !menuItem.IsInactive() {
		t.Error("expected IsInactive to be true for INACTIVE")
	}

	// Test other status values
	menuItem.SetStatus("unknown")
	if menuItem.IsActive() {
		t.Error("expected IsActive to be false for unknown")
	}
	if menuItem.IsInactive() {
		t.Error("expected IsInactive to be false for unknown")
	}
}

func TestMenuItemSoftDeleteMethods(t *testing.T) {
	menuItem := NewMenuItem()
	if menuItem.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false by default")
	}

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	menuItem.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	if menuItem.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false with future date")
	}
	if menuItem.SoftDeletedAt() != future.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt %q, got %q", future.ToDateTimeString(carbon.UTC), menuItem.SoftDeletedAt())
	}

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	menuItem.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	if !menuItem.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be true with past date")
	}
	if menuItem.SoftDeletedAt() != past.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt %q, got %q", past.ToDateTimeString(carbon.UTC), menuItem.SoftDeletedAt())
	}
}

func TestMenuItemMetasMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test empty metas
	metas, err := menuItem.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Error("expected empty metas")
	}

	// Test Meta lookup on empty metas
	if menuItem.Meta("nonexistent") != "" {
		t.Error("expected empty Meta for nonexistent key")
	}

	// Test SetMetas
	err = menuItem.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metas, err = menuItem.Metas()
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
	if menuItem.Meta("layout") != "main" {
		t.Errorf("expected layout %q", "main")
	}
	if menuItem.Meta("theme") != "dark" {
		t.Errorf("expected theme %q", "dark")
	}
	if menuItem.Meta("nonexistent") != "" {
		t.Error("expected empty Meta for nonexistent key")
	}

	// Test SetMeta
	err = menuItem.SetMeta("newkey", "newvalue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if menuItem.Meta("newkey") != "newvalue" {
		t.Errorf("expected newkey %q", "newvalue")
	}

	// Test UpsertMetas
	err = menuItem.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if menuItem.Meta("layout") != "sidebar" { // Updated
		t.Errorf("expected updated layout %q", "sidebar")
	}
	if menuItem.Meta("theme") != "dark" { // Preserved
		t.Errorf("expected preserved theme %q", "dark")
	}
	if menuItem.Meta("newkey") != "newvalue" { // Preserved
		t.Errorf("expected preserved newkey %q", "newvalue")
	}
	if menuItem.Meta("color") != "blue" { // Added
		t.Errorf("expected added color %q", "blue")
	}
}

func TestMenuItemSequenceMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default sequence
	if menuItem.Sequence() != "0" {
		t.Errorf("expected Sequence %q, got %q", "0", menuItem.Sequence())
	}
	if menuItem.SequenceInt() != 0 {
		t.Errorf("expected SequenceInt 0, got %d", menuItem.SequenceInt())
	}

	// Test SetSequenceInt
	menuItem.SetSequenceInt(42)
	if menuItem.Sequence() != "42" {
		t.Errorf("expected Sequence %q, got %q", "42", menuItem.Sequence())
	}
	if menuItem.SequenceInt() != 42 {
		t.Errorf("expected SequenceInt 42, got %d", menuItem.SequenceInt())
	}

	// Test SetSequence
	menuItem.SetSequence("123")
	if menuItem.Sequence() != "123" {
		t.Errorf("expected Sequence %q, got %q", "123", menuItem.Sequence())
	}
	if menuItem.SequenceInt() != 123 {
		t.Errorf("expected SequenceInt 123, got %d", menuItem.SequenceInt())
	}

	// Test invalid sequence
	menuItem.SetSequence("invalid")
	if menuItem.Sequence() != "invalid" {
		t.Errorf("expected Sequence %q, got %q", "invalid", menuItem.Sequence())
	}
	if menuItem.SequenceInt() != 0 { // Should default to 0 for invalid
		t.Errorf("expected SequenceInt 0 for invalid, got %d", menuItem.SequenceInt())
	}
}

func TestMenuItemCreatedAtMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default CreatedAt
	createdAt := menuItem.CreatedAt()
	if createdAt == "" {
		t.Error("expected non-empty CreatedAt")
	}

	createdAtCarbon := menuItem.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if createdAt != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAt to match CreatedAtCarbon")
	}

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	menuItem.SetCreatedAt(testDate)
	if menuItem.CreatedAt() != testDate {
		t.Errorf("expected CreatedAt %q, got %q", testDate, menuItem.CreatedAt())
	}

	createdAtCarbon = menuItem.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if testDate != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAtCarbon to match test date")
	}
}

func TestMenuItemUpdatedAtMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default UpdatedAt
	updatedAt := menuItem.UpdatedAt()
	if updatedAt == "" {
		t.Error("expected non-empty UpdatedAt")
	}

	updatedAtCarbon := menuItem.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if updatedAt != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAt to match UpdatedAtCarbon")
	}

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	menuItem.SetUpdatedAt(testDate)
	if menuItem.UpdatedAt() != testDate {
		t.Errorf("expected UpdatedAt %q, got %q", testDate, menuItem.UpdatedAt())
	}

	updatedAtCarbon = menuItem.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if testDate != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAtCarbon to match test date")
	}
}

func TestMenuItemSoftDeletedAtMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default SoftDeletedAt
	softDeletedAt := menuItem.SoftDeletedAt()
	if softDeletedAt != sb.MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", sb.MAX_DATETIME, softDeletedAt)
	}

	softDeletedAtCarbon := menuItem.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if softDeletedAt != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected SoftDeletedAt to match SoftDeletedAtCarbon")
	}

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	menuItem.SetSoftDeletedAt(testDate)
	if menuItem.SoftDeletedAt() != testDate {
		t.Errorf("expected SoftDeletedAt %q, got %q", testDate, menuItem.SoftDeletedAt())
	}

	softDeletedAtCarbon = menuItem.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if testDate != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected SoftDeletedAtCarbon to match test date")
	}
}

func TestMenuItemIDMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default ID
	id := menuItem.ID()
	if id == "" {
		t.Error("expected non-empty ID")
	}

	// Test SetID
	newID := "test-menu-item-id-123"
	menuItem.SetID(newID)
	if menuItem.ID() != newID {
		t.Errorf("expected ID %q, got %q", newID, menuItem.ID())
	}
}

func TestMenuItemHandleMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default handle
	if menuItem.Handle() != "" {
		t.Error("expected empty Handle")
	}

	// Test SetHandle
	handle := "test-menu-item-handle"
	menuItem.SetHandle(handle)
	if menuItem.Handle() != handle {
		t.Errorf("expected Handle %q, got %q", handle, menuItem.Handle())
	}
}

func TestMenuItemMemoMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default memo
	if menuItem.Memo() != "" {
		t.Error("expected empty Memo")
	}

	// Test SetMemo
	memo := "This is a menu item memo"
	menuItem.SetMemo(memo)
	if menuItem.Memo() != memo {
		t.Errorf("expected Memo %q, got %q", memo, menuItem.Memo())
	}
}

func TestMenuItemMenuIDMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default menu ID
	if menuItem.MenuID() != "" {
		t.Error("expected empty MenuID")
	}

	// Test SetMenuID
	menuID := "test-menu-id"
	menuItem.SetMenuID(menuID)
	if menuItem.MenuID() != menuID {
		t.Errorf("expected MenuID %q, got %q", menuID, menuItem.MenuID())
	}
}

func TestMenuItemNameMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default name
	if menuItem.Name() != "" {
		t.Error("expected empty Name")
	}

	// Test SetName
	name := "Test Menu Item Name"
	menuItem.SetName(name)
	if menuItem.Name() != name {
		t.Errorf("expected Name %q, got %q", name, menuItem.Name())
	}
}

func TestMenuItemPageIDMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default page ID
	if menuItem.PageID() != "" {
		t.Error("expected empty PageID")
	}

	// Test SetPageID
	pageID := "test-page-id"
	menuItem.SetPageID(pageID)
	if menuItem.PageID() != pageID {
		t.Errorf("expected PageID %q, got %q", pageID, menuItem.PageID())
	}
}

func TestMenuItemParentIDMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default parent ID
	if menuItem.ParentID() != "" {
		t.Error("expected empty ParentID")
	}

	// Test SetParentID
	parentID := "test-parent-id"
	menuItem.SetParentID(parentID)
	if menuItem.ParentID() != parentID {
		t.Errorf("expected ParentID %q, got %q", parentID, menuItem.ParentID())
	}
}

func TestMenuItemTargetMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default target
	if menuItem.Target() != "" {
		t.Error("expected empty Target")
	}

	// Test SetTarget
	target := "_blank"
	menuItem.SetTarget(target)
	if menuItem.Target() != target {
		t.Errorf("expected Target %q, got %q", target, menuItem.Target())
	}
}

func TestMenuItemURLMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default URL
	if menuItem.URL() != "" {
		t.Error("expected empty URL")
	}

	// Test SetURL
	url := "https://example.com"
	menuItem.SetURL(url)
	if menuItem.URL() != url {
		t.Errorf("expected URL %q, got %q", url, menuItem.URL())
	}
}

func TestMenuItemStatusSettersAndGetters(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default status
	if menuItem.Status() != MENU_ITEM_STATUS_DRAFT {
		t.Errorf("expected Status %q, got %q", MENU_ITEM_STATUS_DRAFT, menuItem.Status())
	}

	// Test SetStatus
	menuItem.SetStatus(MENU_ITEM_STATUS_ACTIVE)
	if menuItem.Status() != MENU_ITEM_STATUS_ACTIVE {
		t.Errorf("expected Status %q, got %q", MENU_ITEM_STATUS_ACTIVE, menuItem.Status())
	}

	menuItem.SetStatus(MENU_ITEM_STATUS_INACTIVE)
	if menuItem.Status() != MENU_ITEM_STATUS_INACTIVE {
		t.Errorf("expected Status %q, got %q", MENU_ITEM_STATUS_INACTIVE, menuItem.Status())
	}

	menuItem.SetStatus("custom-status")
	if menuItem.Status() != "custom-status" {
		t.Errorf("expected Status %q, got %q", "custom-status", menuItem.Status())
	}
}

func TestMenuItemMarshalToVersioning(t *testing.T) {
	menuItem := NewMenuItem()
	menuItem.SetHandle("test-handle")
	menuItem.SetMemo("test-memo")
	menuItem.SetMenuID("test-menu")
	menuItem.SetName("Test Menu Item")
	menuItem.SetPageID("test-page")
	menuItem.SetParentID("test-parent")
	menuItem.SetSequenceInt(1)
	menuItem.SetTarget("_blank")
	menuItem.SetURL("https://example.com")
	menuItem.SetStatus(MENU_ITEM_STATUS_ACTIVE)

	versionedJSON, err := menuItem.MarshalToVersioning()
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
	if versionedData[COLUMN_HANDLE] != "test-handle" {
		t.Errorf("expected handle %q, got %q", "test-handle", versionedData[COLUMN_HANDLE])
	}
	if versionedData[COLUMN_MEMO] != "test-memo" {
		t.Errorf("expected memo %q, got %q", "test-memo", versionedData[COLUMN_MEMO])
	}
	if versionedData[COLUMN_MENU_ID] != "test-menu" {
		t.Errorf("expected menu_id %q, got %q", "test-menu", versionedData[COLUMN_MENU_ID])
	}
	if versionedData[COLUMN_NAME] != "Test Menu Item" {
		t.Errorf("expected name %q, got %q", "Test Menu Item", versionedData[COLUMN_NAME])
	}
	if versionedData[COLUMN_PAGE_ID] != "test-page" {
		t.Errorf("expected page_id %q, got %q", "test-page", versionedData[COLUMN_PAGE_ID])
	}
	if versionedData[COLUMN_PARENT_ID] != "test-parent" {
		t.Errorf("expected parent_id %q, got %q", "test-parent", versionedData[COLUMN_PARENT_ID])
	}
	if versionedData[COLUMN_SEQUENCE] != "1" {
		t.Errorf("expected sequence %q, got %q", "1", versionedData[COLUMN_SEQUENCE])
	}
	if versionedData[COLUMN_TARGET] != "_blank" {
		t.Errorf("expected target %q, got %q", "_blank", versionedData[COLUMN_TARGET])
	}
	if versionedData[COLUMN_URL] != "https://example.com" {
		t.Errorf("expected url %q, got %q", "https://example.com", versionedData[COLUMN_URL])
	}
	if versionedData[COLUMN_STATUS] != MENU_ITEM_STATUS_ACTIVE {
		t.Errorf("expected status %q, got %q", MENU_ITEM_STATUS_ACTIVE, versionedData[COLUMN_STATUS])
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
