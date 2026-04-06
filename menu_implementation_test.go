package cmsstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewMenuDefaults(t *testing.T) {
	menu := NewMenu()

	if menu.ID() == "" {
		t.Error("expected non-empty ID")
	}
	if menu.CreatedAt() == "" {
		t.Error("expected non-empty CreatedAt")
	}
	if menu.UpdatedAt() == "" {
		t.Error("expected non-empty UpdatedAt")
	}
	if menu.Status() != MENU_STATUS_DRAFT {
		t.Errorf("expected status %q, got %q", MENU_STATUS_DRAFT, menu.Status())
	}
	if menu.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", sb.MAX_DATETIME, menu.SoftDeletedAt())
	}
	if menu.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false")
	}

	metas, err := menu.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Error("expected empty metas")
	}

	createdCarbon := menu.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if menu.CreatedAt() != createdCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAt to match CreatedAtCarbon")
	}

	updatedCarbon := menu.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if menu.UpdatedAt() != updatedCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAt to match UpdatedAtCarbon")
	}

	softDeletedCarbon := menu.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if !softDeletedCarbon.Gte(carbon.Now(carbon.UTC)) {
		t.Error("expected SoftDeletedAtCarbon to be in the future")
	}
}

func TestMenuGetterMethods(t *testing.T) {
	menu := NewMenu()

	// Test default values
	if menu.Handle() != "" {
		t.Error("expected empty Handle")
	}
	if menu.Memo() != "" {
		t.Error("expected empty Memo")
	}
	if menu.Name() != "" {
		t.Error("expected empty Name")
	}
	if menu.SiteID() != "" {
		t.Error("expected empty SiteID")
	}
}

func TestMenuStatusMethods(t *testing.T) {
	menu := NewMenu()

	// Test default status (DRAFT)
	if menu.IsActive() {
		t.Error("expected IsActive to be false for DRAFT")
	}
	if menu.IsInactive() {
		t.Error("expected IsInactive to be false for DRAFT")
	}

	// Test ACTIVE status
	menu.SetStatus(MENU_STATUS_ACTIVE)
	if !menu.IsActive() {
		t.Error("expected IsActive to be true for ACTIVE")
	}
	if menu.IsInactive() {
		t.Error("expected IsInactive to be false for ACTIVE")
	}

	// Test INACTIVE status
	menu.SetStatus(MENU_STATUS_INACTIVE)
	if menu.IsActive() {
		t.Error("expected IsActive to be false for INACTIVE")
	}
	if !menu.IsInactive() {
		t.Error("expected IsInactive to be true for INACTIVE")
	}

	// Test other status values
	menu.SetStatus("unknown")
	if menu.IsActive() {
		t.Error("expected IsActive to be false for unknown")
	}
	if menu.IsInactive() {
		t.Error("expected IsInactive to be false for unknown")
	}
}

func TestMenuSoftDeleteMethods(t *testing.T) {
	menu := NewMenu()
	if menu.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false by default")
	}

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	menu.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	if menu.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false with future date")
	}
	if menu.SoftDeletedAt() != future.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt %q, got %q", future.ToDateTimeString(carbon.UTC), menu.SoftDeletedAt())
	}

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	menu.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	if !menu.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be true with past date")
	}
	if menu.SoftDeletedAt() != past.ToDateTimeString(carbon.UTC) {
		t.Errorf("expected SoftDeletedAt %q, got %q", past.ToDateTimeString(carbon.UTC), menu.SoftDeletedAt())
	}
}

func TestMenuMetasMethods(t *testing.T) {
	menu := NewMenu()

	// Test empty metas
	metas, err := menu.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Error("expected empty metas")
	}

	// Test Meta lookup on empty metas
	if menu.Meta("nonexistent") != "" {
		t.Error("expected empty Meta for nonexistent key")
	}

	// Test SetMetas
	err = menu.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metas, err = menu.Metas()
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
	if menu.Meta("layout") != "main" {
		t.Errorf("expected layout %q", "main")
	}
	if menu.Meta("theme") != "dark" {
		t.Errorf("expected theme %q", "dark")
	}
	if menu.Meta("nonexistent") != "" {
		t.Error("expected empty Meta for nonexistent key")
	}

	// Test SetMeta
	err = menu.SetMeta("newkey", "newvalue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if menu.Meta("newkey") != "newvalue" {
		t.Errorf("expected newkey %q", "newvalue")
	}

	// Test UpsertMetas
	err = menu.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if menu.Meta("layout") != "sidebar" { // Updated
		t.Errorf("expected updated layout %q", "sidebar")
	}
	if menu.Meta("theme") != "dark" { // Preserved
		t.Errorf("expected preserved theme %q", "dark")
	}
	if menu.Meta("newkey") != "newvalue" { // Preserved
		t.Errorf("expected preserved newkey %q", "newvalue")
	}
	if menu.Meta("color") != "blue" { // Added
		t.Errorf("expected added color %q", "blue")
	}
}

func TestMenuCreatedAtMethods(t *testing.T) {
	menu := NewMenu()

	// Test default CreatedAt
	createdAt := menu.CreatedAt()
	if createdAt == "" {
		t.Error("expected non-empty CreatedAt")
	}

	createdAtCarbon := menu.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if createdAt != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAt to match CreatedAtCarbon")
	}

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	menu.SetCreatedAt(testDate)
	if menu.CreatedAt() != testDate {
		t.Errorf("expected CreatedAt %q, got %q", testDate, menu.CreatedAt())
	}

	createdAtCarbon = menu.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if testDate != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAtCarbon to match test date")
	}
}

func TestMenuUpdatedAtMethods(t *testing.T) {
	menu := NewMenu()

	// Test default UpdatedAt
	updatedAt := menu.UpdatedAt()
	if updatedAt == "" {
		t.Error("expected non-empty UpdatedAt")
	}

	updatedAtCarbon := menu.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if updatedAt != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAt to match UpdatedAtCarbon")
	}

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	menu.SetUpdatedAt(testDate)
	if menu.UpdatedAt() != testDate {
		t.Errorf("expected UpdatedAt %q, got %q", testDate, menu.UpdatedAt())
	}

	updatedAtCarbon = menu.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if testDate != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAtCarbon to match test date")
	}
}

func TestMenuSoftDeletedAtMethods(t *testing.T) {
	menu := NewMenu()

	// Test default SoftDeletedAt
	softDeletedAt := menu.SoftDeletedAt()
	if softDeletedAt != sb.MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", sb.MAX_DATETIME, softDeletedAt)
	}

	softDeletedAtCarbon := menu.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if softDeletedAt != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected SoftDeletedAt to match SoftDeletedAtCarbon")
	}

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	menu.SetSoftDeletedAt(testDate)
	if menu.SoftDeletedAt() != testDate {
		t.Errorf("expected SoftDeletedAt %q, got %q", testDate, menu.SoftDeletedAt())
	}

	softDeletedAtCarbon = menu.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if testDate != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected SoftDeletedAtCarbon to match test date")
	}
}

func TestMenuIDMethods(t *testing.T) {
	menu := NewMenu()

	// Test default ID
	id := menu.ID()
	if id == "" {
		t.Error("expected non-empty ID")
	}

	// Test SetID
	newID := "test-menu-id-123"
	menu.SetID(newID)
	if menu.ID() != newID {
		t.Errorf("expected ID %q, got %q", newID, menu.ID())
	}
}

func TestMenuHandleMethods(t *testing.T) {
	menu := NewMenu()

	// Test default handle
	if menu.Handle() != "" {
		t.Error("expected empty Handle")
	}

	// Test SetHandle
	handle := "test-menu-handle"
	menu.SetHandle(handle)
	if menu.Handle() != handle {
		t.Errorf("expected Handle %q, got %q", handle, menu.Handle())
	}
}

func TestMenuMemoMethods(t *testing.T) {
	menu := NewMenu()

	// Test default memo
	if menu.Memo() != "" {
		t.Error("expected empty Memo")
	}

	// Test SetMemo
	memo := "This is a menu memo"
	menu.SetMemo(memo)
	if menu.Memo() != memo {
		t.Errorf("expected Memo %q, got %q", memo, menu.Memo())
	}
}

func TestMenuNameMethods(t *testing.T) {
	menu := NewMenu()

	// Test default name
	if menu.Name() != "" {
		t.Error("expected empty Name")
	}

	// Test SetName
	name := "Test Menu Name"
	menu.SetName(name)
	if menu.Name() != name {
		t.Errorf("expected Name %q, got %q", name, menu.Name())
	}
}

func TestMenuSiteIDMethods(t *testing.T) {
	menu := NewMenu()

	// Test default site ID
	if menu.SiteID() != "" {
		t.Error("expected empty SiteID")
	}

	// Test SetSiteID
	siteID := "test-site-id"
	menu.SetSiteID(siteID)
	if menu.SiteID() != siteID {
		t.Errorf("expected SiteID %q, got %q", siteID, menu.SiteID())
	}
}

func TestMenuStatusSettersAndGetters(t *testing.T) {
	menu := NewMenu()

	// Test default status
	if menu.Status() != MENU_STATUS_DRAFT {
		t.Errorf("expected Status %q, got %q", MENU_STATUS_DRAFT, menu.Status())
	}

	// Test SetStatus
	menu.SetStatus(MENU_STATUS_ACTIVE)
	if menu.Status() != MENU_STATUS_ACTIVE {
		t.Errorf("expected Status %q, got %q", MENU_STATUS_ACTIVE, menu.Status())
	}

	menu.SetStatus(MENU_STATUS_INACTIVE)
	if menu.Status() != MENU_STATUS_INACTIVE {
		t.Errorf("expected Status %q, got %q", MENU_STATUS_INACTIVE, menu.Status())
	}

	menu.SetStatus("custom-status")
	if menu.Status() != "custom-status" {
		t.Errorf("expected Status %q, got %q", "custom-status", menu.Status())
	}
}

func TestMenuMarshalToVersioning(t *testing.T) {
	menu := NewMenu()
	menu.SetHandle("test-handle")
	menu.SetMemo("test-memo")
	menu.SetName("Test Menu")
	menu.SetSiteID("test-site")
	menu.SetStatus(MENU_STATUS_ACTIVE)

	versionedJSON, err := menu.MarshalToVersioning()
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
	if versionedData[COLUMN_NAME] != "Test Menu" {
		t.Errorf("expected name %q, got %q", "Test Menu", versionedData[COLUMN_NAME])
	}
	if versionedData[COLUMN_SITE_ID] != "test-site" {
		t.Errorf("expected site_id %q, got %q", "test-site", versionedData[COLUMN_SITE_ID])
	}
	if versionedData[COLUMN_STATUS] != MENU_STATUS_ACTIVE {
		t.Errorf("expected status %q, got %q", MENU_STATUS_ACTIVE, versionedData[COLUMN_STATUS])
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
