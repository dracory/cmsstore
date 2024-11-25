package cmsstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

// == TYPE ===================================================================

type menu struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*menu)(nil)
var _ MenuInterface = (*menu)(nil)

// == CONSTRUCTORS ==========================================================

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

func NewMenuFromExistingData(data map[string]string) *menu {
	o := &menu{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *menu) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *menu) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *menu) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *menu) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *menu) SetCreatedAt(createdAt string) MenuInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *menu) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *menu) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *menu) SetID(id string) MenuInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the menu
//
// A handle is a human friendly unique identifier for the menu, unlike the ID
func (o *menu) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the menu
//
// A handle is a human friendly unique identifier for the menu, unlike the ID
func (o *menu) SetHandle(handle string) MenuInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *menu) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *menu) SetMemo(memo string) MenuInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

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

func (o *menu) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *menu) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

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

func (o *menu) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *menu) SetName(name string) MenuInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *menu) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

func (o *menu) SetSiteID(siteID string) MenuInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

func (o *menu) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *menu) SetSoftDeletedAt(softDeletedAt string) MenuInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *menu) SoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *menu) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *menu) SetStatus(status string) MenuInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *menu) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *menu) SetUpdatedAt(updatedAt string) MenuInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *menu) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
