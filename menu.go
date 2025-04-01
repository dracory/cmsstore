// Package cmsstore contains the implementation of the CMS store functionality.
package cmsstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

// This file defines the menu structure and related operations.

// == TYPE ===================================================================

// menu represents a menu item in the CMS store.
// It embeds the DataObject from the gouniverse/dataobject package to provide
// common data object functionalities.
type menu struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*menu)(nil)
// The menu type implements the MenuInterface, which defines the methods
// required for menu operations.
var _ MenuInterface = (*menu)(nil)

// == CONSTRUCTORS ==========================================================

// NewMenu creates a new menu instance with default values.
// It initializes the menu with a unique ID, draft status, and current timestamps.
func NewMenu() MenuInterface {
	o := &menu{}
	o.SetHandle("")
	o.SetID(uid.HumanUid())
	o.SetMemo("")
	o.SetMetas(map[string]string{})
	o.SetName("")
	o.SetStatus(MENU_STATUS_DRAFT)
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)
	return o
}

// NewMenuFromExistingData creates a new menu instance from existing data.
// It hydrates the menu with the provided data map.
func NewMenuFromExistingData(data map[string]string) *menu {
	o := &menu{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

// IsActive checks if the menu is in active status.
func (o *menu) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

// IsInactive checks if the menu is in inactive status.
func (o *menu) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

// IsSoftDeleted checks if the menu is soft deleted.
func (o *menu) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

// CreatedAt returns the creation timestamp of the menu.
func (o *menu) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp of the menu.
func (o *menu) SetCreatedAt(createdAt string) MenuInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// CreatedAtCarbon returns the creation timestamp of the menu as a Carbon instance.
func (o *menu) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

// ID returns the unique identifier of the menu.
func (o *menu) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the menu.
func (o *menu) SetID(id string) MenuInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the human-friendly unique identifier of the menu.
//
// A handle is a human-friendly unique identifier for the menu, unlike the ID.
func (o *menu) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the human-friendly unique identifier of the menu.
//
// A handle is a human-friendly unique identifier for the menu, unlike the ID.
func (o *menu) SetHandle(handle string) MenuInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

// Memo returns the memo associated with the menu.
func (o *menu) Memo() string {
	return o.Get(COLUMN_MEMO)
}

// SetMemo sets the memo associated with the menu.
func (o *menu) SetMemo(memo string) MenuInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// Metas returns the metadata associated with the menu.
func (o *menu) Metas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
}

// Meta returns the value of a specific metadata key.
func (o *menu) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// SetMeta sets the value of a specific metadata key.
func (o *menu) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metadata as a JSON string.
// Warning: it overwrites any existing metadata.
func (o *menu) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

// UpsertMetas updates or inserts metadata.
func (o *menu) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

// Name returns the name of the menu.
func (o *menu) Name() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the name of the menu.
func (o *menu) SetName(name string) MenuInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// SiteID returns the site identifier associated with the menu.
func (o *menu) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

// SetSiteID sets the site identifier associated with the menu.
func (o *menu) SetSiteID(siteID string) MenuInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

// SoftDeletedAt returns the soft deletion timestamp of the menu.
func (o *menu) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// SetSoftDeletedAt sets the soft deletion timestamp of the menu.
func (o *menu) SetSoftDeletedAt(softDeletedAt string) MenuInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

// SoftDeletedAtCarbon returns the soft deletion timestamp of the menu as a Carbon instance.
func (o *menu) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

// Status returns the status of the menu.
func (o *menu) Status() string {
	return o.Get(COLUMN_STATUS)
}

// SetStatus sets the status of the menu.
func (o *menu) SetStatus(status string) MenuInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

// UpdatedAt returns the last update timestamp of the menu.
func (o *menu) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the last update timestamp of the menu.
func (o *menu) SetUpdatedAt(updatedAt string) MenuInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// UpdatedAtCarbon returns the last update timestamp of the menu as a Carbon instance.
func (o *menu) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
