package cmsstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

type translationImplementation struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*translation)(nil)
var _ TranslationInterface = (*translationImplementation)(nil)

// == CONSTRUCTORS ==========================================================

func NewTranslation() TranslationInterface {
	o := &translationImplementation{}
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

func NewTranslationFromExistingData(data map[string]string) TranslationInterface {
	o := &translationImplementation{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *translationImplementation) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *translationImplementation) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *translationImplementation) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *translationImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *translationImplementation) SetCreatedAt(createdAt string) TranslationInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *translationImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *translationImplementation) Content() (languageCodeContentMap map[string]string, err error) {
	languageCodeContentStr := o.Get(COLUMN_CONTENT)

	if languageCodeContentStr == "" {
		languageCodeContentStr = "{}"
	}

	languageCodeContentJSON := map[string]string{}
	errJson := json.Unmarshal([]byte(languageCodeContentStr), &languageCodeContentJSON)
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return languageCodeContentJSON, nil
}

func (o *translationImplementation) SetContent(languageCodeContentMap map[string]string) error {
	mapString, err := json.Marshal(languageCodeContentMap)

	if err != nil {
		return err
	}

	o.Set(COLUMN_CONTENT, string(mapString))

	return nil
}

func (o *translationImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *translationImplementation) SetID(id string) TranslationInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the translation
//
// A handle is a human friendly unique identifier for the translation, unlike the ID
func (o *translationImplementation) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the translation
//
// A handle is a human friendly unique identifier for the translation, unlike the ID
func (o *translationImplementation) SetHandle(handle string) TranslationInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *translationImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *translationImplementation) SetMemo(memo string) TranslationInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *translationImplementation) Metas() (map[string]string, error) {
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

func (o *translationImplementation) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *translationImplementation) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *translationImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))

	return nil
}

func (o *translationImplementation) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *translationImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *translationImplementation) SetName(name string) TranslationInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *translationImplementation) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

func (o *translationImplementation) SetSiteID(siteID string) TranslationInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

func (o *translationImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *translationImplementation) SetSoftDeletedAt(softDeletedAt string) TranslationInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *translationImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *translationImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *translationImplementation) SetStatus(status string) TranslationInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *translationImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *translationImplementation) SetUpdatedAt(updatedAt string) TranslationInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *translationImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
