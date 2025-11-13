// Package cmsstore contains the implementation of the CMS store functionality.
package cmsstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// This file defines the menu structure and related operations.

// == TYPE ===================================================================

// menuImplementation represents a menu item in the CMS store.
// It embeds the DataObject from the gouniverse/dataobject package to provide
// common data object functionalities.
type menuImplementation struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*menu)(nil)
// The menu type implements the MenuInterface, which defines the methods
// required for menu operations.
var _ MenuInterface = (*menuImplementation)(nil)

// == CONSTRUCTORS ==========================================================

// NewMenu creates a new menu instance with default values.
// It initializes the menu with a unique ID, draft status, and current timestamps.
func NewMenu() MenuInterface {
	o := &menuImplementation{}
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
func NewMenuFromExistingData(data map[string]string) *menuImplementation {
	o := &menuImplementation{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

// IsActive checks if the menu is in active status.
func (o *menuImplementation) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

// IsInactive checks if the menu is in inactive status.
func (o *menuImplementation) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

// IsSoftDeleted checks if the menu is soft deleted.
func (o *menuImplementation) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

// CreatedAt returns the creation timestamp of the menu.
func (o *menuImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp of the menu.
func (o *menuImplementation) SetCreatedAt(createdAt string) MenuInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// CreatedAtCarbon returns the creation timestamp of the menu as a Carbon instance.
func (o *menuImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

// ID returns the unique identifier of the menu.
func (o *menuImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the menu.
func (o *menuImplementation) SetID(id string) MenuInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the human-friendly unique identifier of the menu.
//
// A handle is a human-friendly unique identifier for the menu, unlike the ID.
func (o *menuImplementation) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the human-friendly unique identifier of the menu.
//
// A handle is a human-friendly unique identifier for the menu, unlike the ID.
func (o *menuImplementation) SetHandle(handle string) MenuInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

// Memo returns the memo associated with the menu.
func (o *menuImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

// SetMemo sets the memo associated with the menu.
func (o *menuImplementation) SetMemo(memo string) MenuInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// Metas returns the metadata associated with the menu.
func (o *menuImplementation) Metas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson := map[string]string{}
	errJson := json.Unmarshal([]byte(metasStr), &metasJson)
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return metasJson, nil
}

// Meta returns the value of a specific metadata key.
func (o *menuImplementation) Meta(name string) string {
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
func (o *menuImplementation) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metadata as a JSON string.
// Warning: it overwrites any existing metadata.
func (o *menuImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))

	return nil
}

// UpsertMetas updates or inserts metadata.
func (o *menuImplementation) UpsertMetas(metas map[string]string) error {
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
func (o *menuImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the name of the menu.
func (o *menuImplementation) SetName(name string) MenuInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// SiteID returns the site identifier associated with the menu.
func (o *menuImplementation) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

// SetSiteID sets the site identifier associated with the menu.
func (o *menuImplementation) SetSiteID(siteID string) MenuInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

// SoftDeletedAt returns the soft deletion timestamp of the menu.
func (o *menuImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// SetSoftDeletedAt sets the soft deletion timestamp of the menu.
func (o *menuImplementation) SetSoftDeletedAt(softDeletedAt string) MenuInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

// SoftDeletedAtCarbon returns the soft deletion timestamp of the menu as a Carbon instance.
func (o *menuImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

// Status returns the status of the menu.
func (o *menuImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

// SetStatus sets the status of the menu.
func (o *menuImplementation) SetStatus(status string) MenuInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

// UpdatedAt returns the last update timestamp of the menu.
func (o *menuImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the last update timestamp of the menu.
func (o *menuImplementation) SetUpdatedAt(updatedAt string) MenuInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// UpdatedAtCarbon returns the last update timestamp of the menu as a Carbon instance.
func (o *menuImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
