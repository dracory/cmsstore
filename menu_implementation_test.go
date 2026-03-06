package cmsstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/require"
)

func TestNewMenuDefaults(t *testing.T) {
	menu := NewMenu()

	require.NotEmpty(t, menu.ID())
	require.NotEmpty(t, menu.CreatedAt())
	require.NotEmpty(t, menu.UpdatedAt())
	require.Equal(t, MENU_STATUS_DRAFT, menu.Status())
	require.Equal(t, sb.MAX_DATETIME, menu.SoftDeletedAt())
	require.False(t, menu.IsSoftDeleted())

	metas, err := menu.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	createdCarbon := menu.CreatedAtCarbon()
	require.NotNil(t, createdCarbon)
	require.Equal(t, menu.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

	updatedCarbon := menu.UpdatedAtCarbon()
	require.NotNil(t, updatedCarbon)
	require.Equal(t, menu.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

	softDeletedCarbon := menu.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedCarbon)
	require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestMenuGetterMethods(t *testing.T) {
	menu := NewMenu()

	// Test default values
	require.Equal(t, "", menu.Handle())
	require.Equal(t, "", menu.Memo())
	require.Equal(t, "", menu.Name())
	require.Equal(t, "", menu.SiteID())
}

func TestMenuStatusMethods(t *testing.T) {
	menu := NewMenu()

	// Test default status (DRAFT)
	require.False(t, menu.IsActive())
	require.False(t, menu.IsInactive())

	// Test ACTIVE status
	menu.SetStatus(MENU_STATUS_ACTIVE)
	require.True(t, menu.IsActive())
	require.False(t, menu.IsInactive())

	// Test INACTIVE status
	menu.SetStatus(MENU_STATUS_INACTIVE)
	require.False(t, menu.IsActive())
	require.True(t, menu.IsInactive())

	// Test other status values
	menu.SetStatus("unknown")
	require.False(t, menu.IsActive())
	require.False(t, menu.IsInactive())
}

func TestMenuSoftDeleteMethods(t *testing.T) {
	menu := NewMenu()
	require.False(t, menu.IsSoftDeleted())

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	menu.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	require.False(t, menu.IsSoftDeleted())
	require.Equal(t, future.ToDateTimeString(carbon.UTC), menu.SoftDeletedAt())

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	menu.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	require.True(t, menu.IsSoftDeleted())
	require.Equal(t, past.ToDateTimeString(carbon.UTC), menu.SoftDeletedAt())
}

func TestMenuMetasMethods(t *testing.T) {
	menu := NewMenu()

	// Test empty metas
	metas, err := menu.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	// Test Meta lookup on empty metas
	require.Equal(t, "", menu.Meta("nonexistent"))

	// Test SetMetas
	err = menu.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	require.NoError(t, err)

	metas, err = menu.Metas()
	require.NoError(t, err)
	require.Equal(t, "main", metas["layout"])
	require.Equal(t, "dark", metas["theme"])

	// Test Meta lookup
	require.Equal(t, "main", menu.Meta("layout"))
	require.Equal(t, "dark", menu.Meta("theme"))
	require.Equal(t, "", menu.Meta("nonexistent"))

	// Test SetMeta
	err = menu.SetMeta("newkey", "newvalue")
	require.NoError(t, err)
	require.Equal(t, "newvalue", menu.Meta("newkey"))

	// Test UpsertMetas
	err = menu.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	require.NoError(t, err)
	require.Equal(t, "sidebar", menu.Meta("layout")) // Updated
	require.Equal(t, "dark", menu.Meta("theme"))     // Preserved
	require.Equal(t, "newvalue", menu.Meta("newkey")) // Preserved
	require.Equal(t, "blue", menu.Meta("color"))      // Added
}

func TestMenuCreatedAtMethods(t *testing.T) {
	menu := NewMenu()

	// Test default CreatedAt
	createdAt := menu.CreatedAt()
	require.NotEmpty(t, createdAt)

	createdAtCarbon := menu.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, createdAt, createdAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	menu.SetCreatedAt(testDate)
	require.Equal(t, testDate, menu.CreatedAt())

	createdAtCarbon = menu.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, testDate, createdAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestMenuUpdatedAtMethods(t *testing.T) {
	menu := NewMenu()

	// Test default UpdatedAt
	updatedAt := menu.UpdatedAt()
	require.NotEmpty(t, updatedAt)

	updatedAtCarbon := menu.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, updatedAt, updatedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	menu.SetUpdatedAt(testDate)
	require.Equal(t, testDate, menu.UpdatedAt())

	updatedAtCarbon = menu.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, testDate, updatedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestMenuSoftDeletedAtMethods(t *testing.T) {
	menu := NewMenu()

	// Test default SoftDeletedAt
	softDeletedAt := menu.SoftDeletedAt()
	require.Equal(t, sb.MAX_DATETIME, softDeletedAt)

	softDeletedAtCarbon := menu.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, softDeletedAt, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	menu.SetSoftDeletedAt(testDate)
	require.Equal(t, testDate, menu.SoftDeletedAt())

	softDeletedAtCarbon = menu.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, testDate, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestMenuIDMethods(t *testing.T) {
	menu := NewMenu()

	// Test default ID
	id := menu.ID()
	require.NotEmpty(t, id)

	// Test SetID
	newID := "test-menu-id-123"
	menu.SetID(newID)
	require.Equal(t, newID, menu.ID())
}

func TestMenuHandleMethods(t *testing.T) {
	menu := NewMenu()

	// Test default handle
	require.Equal(t, "", menu.Handle())

	// Test SetHandle
	handle := "test-menu-handle"
	menu.SetHandle(handle)
	require.Equal(t, handle, menu.Handle())
}

func TestMenuMemoMethods(t *testing.T) {
	menu := NewMenu()

	// Test default memo
	require.Equal(t, "", menu.Memo())

	// Test SetMemo
	memo := "This is a menu memo"
	menu.SetMemo(memo)
	require.Equal(t, memo, menu.Memo())
}

func TestMenuNameMethods(t *testing.T) {
	menu := NewMenu()

	// Test default name
	require.Equal(t, "", menu.Name())

	// Test SetName
	name := "Test Menu Name"
	menu.SetName(name)
	require.Equal(t, name, menu.Name())
}

func TestMenuSiteIDMethods(t *testing.T) {
	menu := NewMenu()

	// Test default site ID
	require.Equal(t, "", menu.SiteID())

	// Test SetSiteID
	siteID := "test-site-id"
	menu.SetSiteID(siteID)
	require.Equal(t, siteID, menu.SiteID())
}

func TestMenuStatusSettersAndGetters(t *testing.T) {
	menu := NewMenu()

	// Test default status
	require.Equal(t, MENU_STATUS_DRAFT, menu.Status())

	// Test SetStatus
	menu.SetStatus(MENU_STATUS_ACTIVE)
	require.Equal(t, MENU_STATUS_ACTIVE, menu.Status())

	menu.SetStatus(MENU_STATUS_INACTIVE)
	require.Equal(t, MENU_STATUS_INACTIVE, menu.Status())

	menu.SetStatus("custom-status")
	require.Equal(t, "custom-status", menu.Status())
}

func TestMenuMarshalToVersioning(t *testing.T) {
	menu := NewMenu()
	menu.SetHandle("test-handle")
	menu.SetMemo("test-memo")
	menu.SetName("Test Menu")
	menu.SetSiteID("test-site")
	menu.SetStatus(MENU_STATUS_ACTIVE)

	versionedJSON, err := menu.MarshalToVersioning()
	require.NoError(t, err)
	require.NotEmpty(t, versionedJSON)

	// Parse the JSON to verify it contains expected fields
	var versionedData map[string]string
	err = json.Unmarshal([]byte(versionedJSON), &versionedData)
	require.NoError(t, err)

	// Check that expected fields are present
	require.Equal(t, "test-handle", versionedData[COLUMN_HANDLE])
	require.Equal(t, "test-memo", versionedData[COLUMN_MEMO])
	require.Equal(t, "Test Menu", versionedData[COLUMN_NAME])
	require.Equal(t, "test-site", versionedData[COLUMN_SITE_ID])
	require.Equal(t, MENU_STATUS_ACTIVE, versionedData[COLUMN_STATUS])

	// Check that timestamps and soft delete fields are excluded
	_, hasCreatedAt := versionedData[COLUMN_CREATED_AT]
	_, hasUpdatedAt := versionedData[COLUMN_UPDATED_AT]
	_, hasSoftDeletedAt := versionedData[COLUMN_SOFT_DELETED_AT]
	require.False(t, hasCreatedAt)
	require.False(t, hasUpdatedAt)
	require.False(t, hasSoftDeletedAt)
}