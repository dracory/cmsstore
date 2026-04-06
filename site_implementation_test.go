package cmsstore

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewSiteDefaults(t *testing.T) {
	site := NewSite()

	// Test default values
	if len(site.ID()) == 0 {
		t.Error("Expected ID to be non-empty")
	}
	if len(site.CreatedAt()) == 0 {
		t.Error("Expected CreatedAt to be non-empty")
	}
	if len(site.UpdatedAt()) == 0 {
		t.Error("Expected UpdatedAt to be non-empty")
	}
	if site.Status() != TEMPLATE_STATUS_DRAFT {
		t.Errorf("Expected Status %s, got %s", TEMPLATE_STATUS_DRAFT, site.Status())
	}
	if site.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("Expected SoftDeletedAt %s, got %s", sb.MAX_DATETIME, site.SoftDeletedAt())
	}
	if site.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false")
	}

	metas, err := site.Metas()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("Expected empty metas, got %v", metas)
	}

	createdCarbon := site.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Error("Expected CreatedAtCarbon to be non-nil")
	}
	if site.CreatedAt() != createdCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected CreatedAt %s, got %s", site.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))
	}

	updatedCarbon := site.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Error("Expected UpdatedAtCarbon to be non-nil")
	}
	if site.UpdatedAt() != updatedCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected UpdatedAt %s, got %s", site.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))
	}

	softDeletedCarbon := site.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Error("Expected SoftDeletedAtCarbon to be non-nil")
	}
	if !softDeletedCarbon.Gte(carbon.Now(carbon.UTC)) {
		t.Error("Expected SoftDeletedAtCarbon to be greater than or equal to now")
	}
}

func TestSiteGetterMethods(t *testing.T) {
	site := NewSite()

	// Test default values
	if site.Handle() != "" {
		t.Errorf("Expected empty Handle, got %s", site.Handle())
	}
	if site.Memo() != "" {
		t.Errorf("Expected empty Memo, got %s", site.Memo())
	}
	if site.Name() != "" {
		t.Errorf("Expected empty Name, got %s", site.Name())
	}
}

func TestSiteStatusMethods(t *testing.T) {
	site := NewSite()

	// Test default status (DRAFT)
	if site.IsActive() {
		t.Error("Expected IsActive to be false")
	}
	if site.IsInactive() {
		t.Error("Expected IsInactive to be false")
	}

	// Test ACTIVE status
	site.SetStatus(PAGE_STATUS_ACTIVE)
	if !site.IsActive() {
		t.Error("Expected IsActive to be true")
	}
	if site.IsInactive() {
		t.Error("Expected IsInactive to be false")
	}

	// Test INACTIVE status
	site.SetStatus(PAGE_STATUS_INACTIVE)
	if site.IsActive() {
		t.Error("Expected IsActive to be false")
	}
	if !site.IsInactive() {
		t.Error("Expected IsInactive to be true")
	}

	// Test other status values
	site.SetStatus("unknown")
	if site.IsActive() {
		t.Error("Expected IsActive to be false")
	}
	if site.IsInactive() {
		t.Error("Expected IsInactive to be false")
	}
}

func TestSiteSoftDeleteMethods(t *testing.T) {
	site := NewSite()
	if site.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false")
	}

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	site.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	if site.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be false")
	}
	if site.SoftDeletedAt() != future.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", future.ToDateTimeString(carbon.UTC), site.SoftDeletedAt())
	}

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	site.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	if !site.IsSoftDeleted() {
		t.Error("Expected IsSoftDeleted to be true")
	}
	if site.SoftDeletedAt() != past.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", past.ToDateTimeString(carbon.UTC), site.SoftDeletedAt())
	}
}

func TestSiteMetasMethods(t *testing.T) {
	site := NewSite()

	// Test empty metas
	metas, err := site.Metas()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("Expected empty metas, got %v", metas)
	}

	// Test Meta lookup on empty metas
	if site.Meta("nonexistent") != "" {
		t.Errorf("Expected empty Meta, got %s", site.Meta("nonexistent"))
	}

	// Test SetMetas
	err = site.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	metas, err = site.Metas()
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
	if site.Meta("layout") != "main" {
		t.Errorf("Expected layout 'main', got %s", site.Meta("layout"))
	}
	if site.Meta("theme") != "dark" {
		t.Errorf("Expected theme 'dark', got %s", site.Meta("theme"))
	}
	if site.Meta("nonexistent") != "" {
		t.Errorf("Expected empty Meta, got %s", site.Meta("nonexistent"))
	}

	// Test SetMeta
	err = site.SetMeta("newkey", "newvalue")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if site.Meta("newkey") != "newvalue" {
		t.Errorf("Expected newkey 'newvalue', got %s", site.Meta("newkey"))
	}

	// Test UpsertMetas
	err = site.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if site.Meta("layout") != "sidebar" { // Updated
		t.Errorf("Expected layout 'sidebar', got %s", site.Meta("layout"))
	}
	if site.Meta("theme") != "dark" { // Preserved
		t.Errorf("Expected theme 'dark', got %s", site.Meta("theme"))
	}
	if site.Meta("newkey") != "newvalue" { // Preserved
		t.Errorf("Expected newkey 'newvalue', got %s", site.Meta("newkey"))
	}
	if site.Meta("color") != "blue" { // Added
		t.Errorf("Expected color 'blue', got %s", site.Meta("color"))
	}
}

func TestSiteDomainNamesMethods(t *testing.T) {
	site := NewSite()

	// Test default domain names
	domainNames, err := site.DomainNames()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(domainNames) != 0 {
		t.Errorf("Expected empty domainNames, got %v", domainNames)
	}

	// Test SetDomainNames
	domains := []string{"example.com", "www.example.com"}
	site, err = site.SetDomainNames(domains)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	domainNames, err = site.DomainNames()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !slices.Equal(domainNames, domains) {
		t.Errorf("Expected domainNames %v, got %v", domains, domainNames)
	}

	// Test empty domain names
	site, err = site.SetDomainNames([]string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	domainNames, err = site.DomainNames()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(domainNames) != 0 {
		t.Errorf("Expected empty domainNames, got %v", domainNames)
	}
}

func TestSiteCreatedAtMethods(t *testing.T) {
	site := NewSite()

	// Test default CreatedAt
	createdAt := site.CreatedAt()
	if len(createdAt) == 0 {
		t.Error("Expected CreatedAt to be non-empty")
	}

	createdAtCarbon := site.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Error("Expected CreatedAtCarbon to be non-nil")
	}
	if createdAt != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected CreatedAt %s, got %s", createdAt, createdAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	site.SetCreatedAt(testDate)
	if site.CreatedAt() != testDate {
		t.Errorf("Expected CreatedAt %s, got %s", testDate, site.CreatedAt())
	}

	createdAtCarbon = site.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Error("Expected CreatedAtCarbon to be non-nil")
	}
	if testDate != createdAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected CreatedAt %s, got %s", testDate, createdAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestSiteUpdatedAtMethods(t *testing.T) {
	site := NewSite()

	// Test default UpdatedAt
	updatedAt := site.UpdatedAt()
	if len(updatedAt) == 0 {
		t.Error("Expected UpdatedAt to be non-empty")
	}

	updatedAtCarbon := site.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Error("Expected UpdatedAtCarbon to be non-nil")
	}
	if updatedAt != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected UpdatedAt %s, got %s", updatedAt, updatedAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	site.SetUpdatedAt(testDate)
	if site.UpdatedAt() != testDate {
		t.Errorf("Expected UpdatedAt %s, got %s", testDate, site.UpdatedAt())
	}

	updatedAtCarbon = site.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Error("Expected UpdatedAtCarbon to be non-nil")
	}
	if testDate != updatedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected UpdatedAt %s, got %s", testDate, updatedAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestSiteSoftDeletedAtMethods(t *testing.T) {
	site := NewSite()

	// Test default SoftDeletedAt
	softDeletedAt := site.SoftDeletedAt()
	if softDeletedAt != sb.MAX_DATETIME {
		t.Errorf("Expected SoftDeletedAt %s, got %s", sb.MAX_DATETIME, softDeletedAt)
	}

	softDeletedAtCarbon := site.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Error("Expected SoftDeletedAtCarbon to be non-nil")
	}
	if softDeletedAt != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", softDeletedAt, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	site.SetSoftDeletedAt(testDate)
	if site.SoftDeletedAt() != testDate {
		t.Errorf("Expected SoftDeletedAt %s, got %s", testDate, site.SoftDeletedAt())
	}

	softDeletedAtCarbon = site.SoftDeletedAtCarbon()
	if softDeletedAtCarbon == nil {
		t.Error("Expected SoftDeletedAtCarbon to be non-nil")
	}
	if testDate != softDeletedAtCarbon.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", testDate, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestSiteIDMethods(t *testing.T) {
	site := NewSite()

	// Test default ID
	id := site.ID()
	if len(id) == 0 {
		t.Error("Expected ID to be non-empty")
	}

	// Test SetID
	newID := "test-site-id-123"
	site.SetID(newID)
	if site.ID() != newID {
		t.Errorf("Expected ID %s, got %s", newID, site.ID())
	}
}

func TestSiteHandleMethods(t *testing.T) {
	site := NewSite()

	// Test default handle
	if site.Handle() != "" {
		t.Errorf("Expected empty Handle, got %s", site.Handle())
	}

	// Test SetHandle
	handle := "test-site-handle"
	site.SetHandle(handle)
	if site.Handle() != handle {
		t.Errorf("Expected Handle %s, got %s", handle, site.Handle())
	}
}

func TestSiteMemoMethods(t *testing.T) {
	site := NewSite()

	// Test default memo
	if site.Memo() != "" {
		t.Errorf("Expected empty Memo, got %s", site.Memo())
	}

	// Test SetMemo
	memo := "This is a site memo"
	site.SetMemo(memo)
	if site.Memo() != memo {
		t.Errorf("Expected Memo %s, got %s", memo, site.Memo())
	}
}

func TestSiteNameMethods(t *testing.T) {
	site := NewSite()

	// Test default name
	if site.Name() != "" {
		t.Errorf("Expected empty Name, got %s", site.Name())
	}

	// Test SetName
	name := "Test Site Name"
	site.SetName(name)
	if site.Name() != name {
		t.Errorf("Expected Name %s, got %s", name, site.Name())
	}
}

func TestSiteStatusSettersAndGetters(t *testing.T) {
	site := NewSite()

	// Test default status
	if site.Status() != TEMPLATE_STATUS_DRAFT {
		t.Errorf("Expected Status %s, got %s", TEMPLATE_STATUS_DRAFT, site.Status())
	}

	// Test SetStatus
	site.SetStatus(PAGE_STATUS_ACTIVE)
	if site.Status() != PAGE_STATUS_ACTIVE {
		t.Errorf("Expected Status %s, got %s", PAGE_STATUS_ACTIVE, site.Status())
	}

	site.SetStatus(PAGE_STATUS_INACTIVE)
	if site.Status() != PAGE_STATUS_INACTIVE {
		t.Errorf("Expected Status %s, got %s", PAGE_STATUS_INACTIVE, site.Status())
	}

	site.SetStatus("custom-status")
	if site.Status() != "custom-status" {
		t.Errorf("Expected Status %s, got %s", "custom-status", site.Status())
	}
}

func TestSiteMarshalToVersioning(t *testing.T) {
	site := NewSite()
	site.SetHandle("test-handle")
	site.SetMemo("test-memo")
	site.SetName("Test Site")
	site.SetStatus(PAGE_STATUS_ACTIVE)

	versionedJSON, err := site.MarshalToVersioning()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(versionedJSON) == 0 {
		t.Error("Expected versionedJSON to be non-empty")
	}

	// Parse the JSON to verify it contains expected fields
	var versionedData map[string]string
	err = json.Unmarshal([]byte(versionedJSON), &versionedData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that expected fields are present
	if versionedData[COLUMN_HANDLE] != "test-handle" {
		t.Errorf("Expected handle 'test-handle', got %s", versionedData[COLUMN_HANDLE])
	}
	if versionedData[COLUMN_MEMO] != "test-memo" {
		t.Errorf("Expected memo 'test-memo', got %s", versionedData[COLUMN_MEMO])
	}
	if versionedData[COLUMN_NAME] != "Test Site" {
		t.Errorf("Expected name 'Test Site', got %s", versionedData[COLUMN_NAME])
	}
	if versionedData[COLUMN_STATUS] != PAGE_STATUS_ACTIVE {
		t.Errorf("Expected status %s, got %s", PAGE_STATUS_ACTIVE, versionedData[COLUMN_STATUS])
	}

	// Check that timestamps and soft delete fields are excluded
	_, hasCreatedAt := versionedData[COLUMN_CREATED_AT]
	_, hasUpdatedAt := versionedData[COLUMN_UPDATED_AT]
	_, hasSoftDeletedAt := versionedData[COLUMN_SOFT_DELETED_AT]
	if hasCreatedAt {
		t.Error("Expected hasCreatedAt to be false")
	}
	if hasUpdatedAt {
		t.Error("Expected hasUpdatedAt to be false")
	}
	if hasSoftDeletedAt {
		t.Error("Expected hasSoftDeletedAt to be false")
	}
}
