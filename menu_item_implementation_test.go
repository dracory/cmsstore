package cmsstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/require"
)

func TestNewMenuItemDefaults(t *testing.T) {
	menuItem := NewMenuItem()

	require.NotEmpty(t, menuItem.ID())
	require.NotEmpty(t, menuItem.CreatedAt())
	require.NotEmpty(t, menuItem.UpdatedAt())
	require.Equal(t, MENU_ITEM_STATUS_DRAFT, menuItem.Status())
	require.Equal(t, sb.MAX_DATETIME, menuItem.SoftDeletedAt())
	require.False(t, menuItem.IsSoftDeleted())

	metas, err := menuItem.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	createdCarbon := menuItem.CreatedAtCarbon()
	require.NotNil(t, createdCarbon)
	require.Equal(t, menuItem.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

	updatedCarbon := menuItem.UpdatedAtCarbon()
	require.NotNil(t, updatedCarbon)
	require.Equal(t, menuItem.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

	softDeletedCarbon := menuItem.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedCarbon)
	require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)))
}

func TestMenuItemGetterMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default values
	require.Equal(t, "", menuItem.Handle())
	require.Equal(t, "", menuItem.Memo())
	require.Equal(t, "", menuItem.MenuID())
	require.Equal(t, "", menuItem.Name())
	require.Equal(t, "", menuItem.PageID())
	require.Equal(t, "", menuItem.ParentID())
	require.Equal(t, "0", menuItem.Sequence())
	require.Equal(t, 0, menuItem.SequenceInt())
	require.Equal(t, "", menuItem.Target())
	require.Equal(t, "", menuItem.URL())
}

func TestMenuItemStatusMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default status (DRAFT)
	require.False(t, menuItem.IsActive())
	require.False(t, menuItem.IsInactive())

	// Test ACTIVE status
	menuItem.SetStatus(MENU_ITEM_STATUS_ACTIVE)
	require.True(t, menuItem.IsActive())
	require.False(t, menuItem.IsInactive())

	// Test INACTIVE status
	menuItem.SetStatus(MENU_ITEM_STATUS_INACTIVE)
	require.False(t, menuItem.IsActive())
	require.True(t, menuItem.IsInactive())

	// Test other status values
	menuItem.SetStatus("unknown")
	require.False(t, menuItem.IsActive())
	require.False(t, menuItem.IsInactive())
}

func TestMenuItemSoftDeleteMethods(t *testing.T) {
	menuItem := NewMenuItem()
	require.False(t, menuItem.IsSoftDeleted())

	// Test with future date
	future := carbon.Now(carbon.UTC).AddHour()
	menuItem.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	require.False(t, menuItem.IsSoftDeleted())
	require.Equal(t, future.ToDateTimeString(carbon.UTC), menuItem.SoftDeletedAt())

	// Test with past date
	past := carbon.Now(carbon.UTC).SubHour()
	menuItem.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	require.True(t, menuItem.IsSoftDeleted())
	require.Equal(t, past.ToDateTimeString(carbon.UTC), menuItem.SoftDeletedAt())
}

func TestMenuItemMetasMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test empty metas
	metas, err := menuItem.Metas()
	require.NoError(t, err)
	require.Empty(t, metas)

	// Test Meta lookup on empty metas
	require.Equal(t, "", menuItem.Meta("nonexistent"))

	// Test SetMetas
	err = menuItem.SetMetas(map[string]string{"layout": "main", "theme": "dark"})
	require.NoError(t, err)

	metas, err = menuItem.Metas()
	require.NoError(t, err)
	require.Equal(t, "main", metas["layout"])
	require.Equal(t, "dark", metas["theme"])

	// Test Meta lookup
	require.Equal(t, "main", menuItem.Meta("layout"))
	require.Equal(t, "dark", menuItem.Meta("theme"))
	require.Equal(t, "", menuItem.Meta("nonexistent"))

	// Test SetMeta
	err = menuItem.SetMeta("newkey", "newvalue")
	require.NoError(t, err)
	require.Equal(t, "newvalue", menuItem.Meta("newkey"))

	// Test UpsertMetas
	err = menuItem.UpsertMetas(map[string]string{"layout": "sidebar", "color": "blue"})
	require.NoError(t, err)
	require.Equal(t, "sidebar", menuItem.Meta("layout")) // Updated
	require.Equal(t, "dark", menuItem.Meta("theme"))     // Preserved
	require.Equal(t, "newvalue", menuItem.Meta("newkey")) // Preserved
	require.Equal(t, "blue", menuItem.Meta("color"))      // Added
}

func TestMenuItemSequenceMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default sequence
	require.Equal(t, "0", menuItem.Sequence())
	require.Equal(t, 0, menuItem.SequenceInt())

	// Test SetSequenceInt
	menuItem.SetSequenceInt(42)
	require.Equal(t, "42", menuItem.Sequence())
	require.Equal(t, 42, menuItem.SequenceInt())

	// Test SetSequence
	menuItem.SetSequence("123")
	require.Equal(t, "123", menuItem.Sequence())
	require.Equal(t, 123, menuItem.SequenceInt())

	// Test invalid sequence
	menuItem.SetSequence("invalid")
	require.Equal(t, "invalid", menuItem.Sequence())
	require.Equal(t, 0, menuItem.SequenceInt()) // Should default to 0 for invalid
}

func TestMenuItemCreatedAtMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default CreatedAt
	createdAt := menuItem.CreatedAt()
	require.NotEmpty(t, createdAt)

	createdAtCarbon := menuItem.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, createdAt, createdAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetCreatedAt
	testDate := "2023-12-25 10:30:00"
	menuItem.SetCreatedAt(testDate)
	require.Equal(t, testDate, menuItem.CreatedAt())

	createdAtCarbon = menuItem.CreatedAtCarbon()
	require.NotNil(t, createdAtCarbon)
	require.Equal(t, testDate, createdAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestMenuItemUpdatedAtMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default UpdatedAt
	updatedAt := menuItem.UpdatedAt()
	require.NotEmpty(t, updatedAt)

	updatedAtCarbon := menuItem.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, updatedAt, updatedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetUpdatedAt
	testDate := "2023-12-25 15:45:00"
	menuItem.SetUpdatedAt(testDate)
	require.Equal(t, testDate, menuItem.UpdatedAt())

	updatedAtCarbon = menuItem.UpdatedAtCarbon()
	require.NotNil(t, updatedAtCarbon)
	require.Equal(t, testDate, updatedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestMenuItemSoftDeletedAtMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default SoftDeletedAt
	softDeletedAt := menuItem.SoftDeletedAt()
	require.Equal(t, sb.MAX_DATETIME, softDeletedAt)

	softDeletedAtCarbon := menuItem.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, softDeletedAt, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))

	// Test SetSoftDeletedAt
	testDate := "2023-12-25 20:00:00"
	menuItem.SetSoftDeletedAt(testDate)
	require.Equal(t, testDate, menuItem.SoftDeletedAt())

	softDeletedAtCarbon = menuItem.SoftDeletedAtCarbon()
	require.NotNil(t, softDeletedAtCarbon)
	require.Equal(t, testDate, softDeletedAtCarbon.ToDateTimeString(carbon.UTC))
}

func TestMenuItemIDMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default ID
	id := menuItem.ID()
	require.NotEmpty(t, id)

	// Test SetID
	newID := "test-menu-item-id-123"
	menuItem.SetID(newID)
	require.Equal(t, newID, menuItem.ID())
}


func TestMenuItemHandleMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default handle
	require.Equal(t, "", menuItem.Handle())

	// Test SetHandle
	handle := "test-menu-item-handle"
	menuItem.SetHandle(handle)
	require.Equal(t, handle, menuItem.Handle())
}

func TestMenuItemMemoMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default memo
	require.Equal(t, "", menuItem.Memo())

	// Test SetMemo
	memo := "This is a menu item memo"
	menuItem.SetMemo(memo)
	require.Equal(t, memo, menuItem.Memo())
}

func TestMenuItemMenuIDMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default menu ID
	require.Equal(t, "", menuItem.MenuID())

	// Test SetMenuID
	menuID := "test-menu-id"
	menuItem.SetMenuID(menuID)
	require.Equal(t, menuID, menuItem.MenuID())
}

func TestMenuItemNameMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default name
	require.Equal(t, "", menuItem.Name())

	// Test SetName
	name := "Test Menu Item Name"
	menuItem.SetName(name)
	require.Equal(t, name, menuItem.Name())
}

func TestMenuItemPageIDMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default page ID
	require.Equal(t, "", menuItem.PageID())

	// Test SetPageID
	pageID := "test-page-id"
	menuItem.SetPageID(pageID)
	require.Equal(t, pageID, menuItem.PageID())
}

func TestMenuItemParentIDMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default parent ID
	require.Equal(t, "", menuItem.ParentID())

	// Test SetParentID
	parentID := "test-parent-id"
	menuItem.SetParentID(parentID)
	require.Equal(t, parentID, menuItem.ParentID())
}

func TestMenuItemTargetMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default target
	require.Equal(t, "", menuItem.Target())

	// Test SetTarget
	target := "_blank"
	menuItem.SetTarget(target)
	require.Equal(t, target, menuItem.Target())
}

func TestMenuItemURLMethods(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default URL
	require.Equal(t, "", menuItem.URL())

	// Test SetURL
	url := "https://example.com"
	menuItem.SetURL(url)
	require.Equal(t, url, menuItem.URL())
}

func TestMenuItemStatusSettersAndGetters(t *testing.T) {
	menuItem := NewMenuItem()

	// Test default status
	require.Equal(t, MENU_ITEM_STATUS_DRAFT, menuItem.Status())

	// Test SetStatus
	menuItem.SetStatus(MENU_ITEM_STATUS_ACTIVE)
	require.Equal(t, MENU_ITEM_STATUS_ACTIVE, menuItem.Status())

	menuItem.SetStatus(MENU_ITEM_STATUS_INACTIVE)
	require.Equal(t, MENU_ITEM_STATUS_INACTIVE, menuItem.Status())

	menuItem.SetStatus("custom-status")
	require.Equal(t, "custom-status", menuItem.Status())
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
	require.NoError(t, err)
	require.NotEmpty(t, versionedJSON)

	// Parse the JSON to verify it contains expected fields
	var versionedData map[string]string
	err = json.Unmarshal([]byte(versionedJSON), &versionedData)
	require.NoError(t, err)

	// Check that expected fields are present
	require.Equal(t, "test-handle", versionedData[COLUMN_HANDLE])
	require.Equal(t, "test-memo", versionedData[COLUMN_MEMO])
	require.Equal(t, "test-menu", versionedData[COLUMN_MENU_ID])
	require.Equal(t, "Test Menu Item", versionedData[COLUMN_NAME])
	require.Equal(t, "test-page", versionedData[COLUMN_PAGE_ID])
	require.Equal(t, "test-parent", versionedData[COLUMN_PARENT_ID])
	require.Equal(t, "1", versionedData[COLUMN_SEQUENCE])
	require.Equal(t, "_blank", versionedData[COLUMN_TARGET])
	require.Equal(t, "https://example.com", versionedData[COLUMN_URL])
	require.Equal(t, MENU_ITEM_STATUS_ACTIVE, versionedData[COLUMN_STATUS])

	// Check that timestamps and soft delete fields are excluded
	_, hasCreatedAt := versionedData[COLUMN_CREATED_AT]
	_, hasUpdatedAt := versionedData[COLUMN_UPDATED_AT]
	_, hasSoftDeletedAt := versionedData[COLUMN_SOFT_DELETED_AT]
	require.False(t, hasCreatedAt)
	require.False(t, hasUpdatedAt)
	require.False(t, hasSoftDeletedAt)
}