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

type translation struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*translation)(nil)
var _ TranslationInterface = (*translation)(nil)

// == CONSTRUCTORS ==========================================================

func NewTranslation() TranslationInterface {
	o := &translation{}
	o.SetContent(map[string]string{})
	o.SetHandle("")
	o.SetID(uid.HumanUid())
	o.SetMemo("")
	o.SetMetas(map[string]string{})
	o.SetName("")
	o.SetStatus(TEMPLATE_STATUS_DRAFT)
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)
	return o
}

func NewTranslationFromExistingData(data map[string]string) *translation {
	o := &translation{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *translation) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *translation) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *translation) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *translation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *translation) SetCreatedAt(createdAt string) TranslationInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *translation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *translation) Content() (languageCodeContentMap map[string]string, err error) {
	languageCodeContentStr := o.Get(COLUMN_CONTENT)

	if languageCodeContentStr == "" {
		languageCodeContentStr = "{}"
	}

	languageCodeContentJSON, errJson := utils.FromJSON(languageCodeContentStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(languageCodeContentJSON.(map[string]any)), nil
}

func (o *translation) SetContent(languageCodeContentMap map[string]string) error {
	mapString, err := utils.ToJSON(languageCodeContentMap)

	if err != nil {
		return err
	}

	o.Set(COLUMN_CONTENT, mapString)

	return nil
}

func (o *translation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *translation) SetID(id string) TranslationInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the translation
//
// A handle is a human friendly unique identifier for the translation, unlike the ID
func (o *translation) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the translation
//
// A handle is a human friendly unique identifier for the translation, unlike the ID
func (o *translation) SetHandle(handle string) TranslationInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *translation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *translation) SetMemo(memo string) TranslationInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *translation) Metas() (map[string]string, error) {
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

func (o *translation) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *translation) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *translation) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

func (o *translation) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *translation) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *translation) SetName(name string) TranslationInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *translation) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

func (o *translation) SetSiteID(siteID string) TranslationInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

func (o *translation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *translation) SetSoftDeletedAt(softDeletedAt string) TranslationInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *translation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *translation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *translation) SetStatus(status string) TranslationInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *translation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *translation) SetUpdatedAt(updatedAt string) TranslationInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *translation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
