package cmsstore

import (
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
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

func (o *menuItem) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *menuItem) SetCreatedAt(createdAt string) MenuItemInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *menuItem) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *menuItem) Content() string {
	return o.Get(COLUMN_CONTENT)
}

func (o *menuItem) SetContent(content string) MenuItemInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

func (o *menuItem) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

func (o *menuItem) SetEditor(editor string) MenuItemInterface {
	o.Set(COLUMN_EDITOR, editor)
	return o
}

func (o *menuItem) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *menuItem) SetID(id string) MenuItemInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the menuItem
//
// A handle is a human friendly unique identifier for the menuItem, unlike the ID
func (o *menuItem) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the menuItem
//
// A handle is a human friendly unique identifier for the menuItem, unlike the ID
func (o *menuItem) SetHandle(handle string) MenuItemInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *menuItem) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *menuItem) SetMemo(memo string) MenuItemInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *menuItem) MenuID() string {
	return o.Get(COLUMN_MENU_ID)
}

func (o *menuItem) SetMenuID(siteID string) MenuItemInterface {
	o.Set(COLUMN_MENU_ID, siteID)
	return o
}

func (o *menuItem) Metas() (map[string]string, error) {
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

func (o *menuItem) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *menuItem) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

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

func (o *menuItem) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *menuItem) SetName(name string) MenuItemInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *menuItem) PageID() string {
	return o.Get(COLUMN_PAGE_ID)
}

func (o *menuItem) SetPageID(siteID string) MenuItemInterface {
	o.Set(COLUMN_PAGE_ID, siteID)
	return o
}

func (o *menuItem) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *menuItem) SetParentID(parentID string) MenuItemInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *menuItem) Sequence() string {
	return o.Get(COLUMN_SEQUENCE)
}

func (o *menuItem) SequenceInt() int {
	return cast.ToInt(o.Sequence())
}

func (o *menuItem) SetSequence(sequence string) MenuItemInterface {
	o.Set(COLUMN_SEQUENCE, sequence)
	return o
}

func (o *menuItem) SetSequenceInt(sequence int) MenuItemInterface {
	o.SetSequence(cast.ToString(sequence))
	return o
}

func (o *menuItem) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *menuItem) SetSoftDeletedAt(softDeletedAt string) MenuItemInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *menuItem) SoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *menuItem) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *menuItem) SetStatus(status string) MenuItemInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *menuItem) Target() string {
	return o.Get(COLUMN_TARGET)
}

func (o *menuItem) SetTarget(target string) MenuItemInterface {
	o.Set(COLUMN_TARGET, target)
	return o
}

func (o *menuItem) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *menuItem) SetUpdatedAt(updatedAt string) MenuItemInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *menuItem) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}

func (o *menuItem) URL() string {
	return o.Get(COLUMN_URL)
}

func (o *menuItem) SetURL(url string) MenuItemInterface {
	o.Set(COLUMN_URL, url)
	return o
}
