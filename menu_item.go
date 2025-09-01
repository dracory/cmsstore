package cmsstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == TYPE ===================================================================

type menuItem struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*menuItem)(nil)
var _ MenuItemInterface = (*menuItem)(nil)

// == CONSTRUCTORS ==========================================================

// NewMenuItem creates a new menu item with default values.
func NewMenuItem() MenuItemInterface {
	o := &menuItem{}
	o.SetID(uid.HumanUid())
	o.SetMemo("")
	o.SetMetas(map[string]string{})
	o.SetName("")
	o.SetPageID("")
	o.SetStatus(MENU_ITEM_STATUS_DRAFT)
	o.SetTarget("")
	o.SetURL("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)
	return o
}

func NewMenuItemFromExistingData(data map[string]string) *menuItem {
	o := &menuItem{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *menuItem) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *menuItem) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *menuItem) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

// CreatedAt returns the creation timestamp of the menu item.
func (o *menuItem) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp of the menu item.
func (o *menuItem) SetCreatedAt(createdAt string) MenuItemInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// CreatedAtCarbon returns the creation timestamp of the menu item as a Carbon object.
func (o *menuItem) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

// Content returns the content of the menu item.
func (o *menuItem) Content() string {
	return o.Get(COLUMN_CONTENT)
}

// SetContent sets the content of the menu item.
func (o *menuItem) SetContent(content string) MenuItemInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

// Editor returns the editor of the menu item.
func (o *menuItem) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

// SetEditor sets the editor of the menu item.
func (o *menuItem) SetEditor(editor string) MenuItemInterface {
	o.Set(COLUMN_EDITOR, editor)
	return o
}

// ID returns the ID of the menu item.
func (o *menuItem) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the ID of the menu item.
func (o *menuItem) SetID(id string) MenuItemInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the menu item.
//
// A handle is a human-friendly unique identifier for the menu item, unlike the ID.
func (o *menuItem) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the menu item.
//
// A handle is a human-friendly unique identifier for the menu item, unlike the ID.
func (o *menuItem) SetHandle(handle string) MenuItemInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

// Memo returns the memo of the menu item.
func (o *menuItem) Memo() string {
	return o.Get(COLUMN_MEMO)
}

// SetMemo sets the memo of the menu item.
func (o *menuItem) SetMemo(memo string) MenuItemInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// MenuID returns the ID of the menu associated with the menu item.
func (o *menuItem) MenuID() string {
	return o.Get(COLUMN_MENU_ID)
}

// SetMenuID sets the ID of the menu associated with the menu item.
func (o *menuItem) SetMenuID(siteID string) MenuItemInterface {
	o.Set(COLUMN_MENU_ID, siteID)
	return o
}

// Metas returns the metas of the menu item as a map.
//
// Metas are additional metadata stored as JSON.
func (o *menuItem) Metas() (map[string]string, error) {
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

// Meta returns the value of a specific meta for the menu item.
func (o *menuItem) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// SetMeta sets the value of a specific meta for the menu item.
func (o *menuItem) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas sets the metas of the menu item.
//
// Warning: This method overwrites any existing metas with the provided map.
func (o *menuItem) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))

	return nil
}

// UpsertMetas merges the provided metas with existing metas.
func (o *menuItem) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

// Name returns the name of the menu item.
func (o *menuItem) Name() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the name of the menu item.
func (o *menuItem) SetName(name string) MenuItemInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// PageID returns the ID of the page associated with the menu item.
func (o *menuItem) PageID() string {
	return o.Get(COLUMN_PAGE_ID)
}

// SetPageID sets the ID of the page associated with the menu item.
func (o *menuItem) SetPageID(siteID string) MenuItemInterface {
	o.Set(COLUMN_PAGE_ID, siteID)
	return o
}

// ParentID returns the ID of the parent menu item.
func (o *menuItem) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

// SetParentID sets the ID of the parent menu item.
func (o *menuItem) SetParentID(parentID string) MenuItemInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

// Sequence returns the sequence of the menu item.
func (o *menuItem) Sequence() string {
	return o.Get(COLUMN_SEQUENCE)
}

// SequenceInt returns the sequence of the menu item as an integer.
func (o *menuItem) SequenceInt() int {
	return cast.ToInt(o.Sequence())
}

// SetSequence sets the sequence of the menu item.
func (o *menuItem) SetSequence(sequence string) MenuItemInterface {
	o.Set(COLUMN_SEQUENCE, sequence)
	return o
}

// SetSequenceInt sets the sequence of the menu item as an integer.
func (o *menuItem) SetSequenceInt(sequence int) MenuItemInterface {
	o.SetSequence(cast.ToString(sequence))
	return o
}

// SoftDeletedAt returns the soft deletion timestamp of the menu item.
func (o *menuItem) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// SetSoftDeletedAt sets the soft deletion timestamp of the menu item.
func (o *menuItem) SetSoftDeletedAt(softDeletedAt string) MenuItemInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *menuItem) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

// Status returns the status of the menu item.
func (o *menuItem) Status() string {
	return o.Get(COLUMN_STATUS)
}

// SetStatus sets the status of the menu item.
func (o *menuItem) SetStatus(status string) MenuItemInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

// Target returns the target attribute of the menu item.
func (o *menuItem) Target() string {
	return o.Get(COLUMN_TARGET)
}

// SetTarget sets the target attribute of the menu item.
func (o *menuItem) SetTarget(target string) MenuItemInterface {
	o.Set(COLUMN_TARGET, target)
	return o
}

// UpdatedAt returns the last update timestamp of the menu item.
func (o *menuItem) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the last update timestamp of the menu item.
func (o *menuItem) SetUpdatedAt(updatedAt string) MenuItemInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *menuItem) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}

// URL returns the URL of the menu item.
func (o *menuItem) URL() string {
	return o.Get(COLUMN_URL)
}

// SetURL sets the URL of the menu item.
func (o *menuItem) SetURL(url string) MenuItemInterface {
	o.Set(COLUMN_URL, url)
	return o
}
