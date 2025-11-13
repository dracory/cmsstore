package cmsstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

type templateImplementation struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*template)(nil)
var _ TemplateInterface = (*templateImplementation)(nil)

// == CONSTRUCTORS ==========================================================

func NewTemplate() TemplateInterface {
	o := &templateImplementation{}
	o.SetContent("")
	o.SetEditor("")
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

func NewTemplateFromExistingData(data map[string]string) TemplateInterface {
	o := &templateImplementation{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *templateImplementation) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *templateImplementation) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *templateImplementation) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *templateImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *templateImplementation) SetCreatedAt(createdAt string) TemplateInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *templateImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *templateImplementation) Content() string {
	return o.Get(COLUMN_CONTENT)
}

func (o *templateImplementation) SetContent(content string) TemplateInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

func (o *templateImplementation) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

func (o *templateImplementation) SetEditor(editor string) TemplateInterface {
	o.Set(COLUMN_EDITOR, editor)
	return o
}

func (o *templateImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *templateImplementation) SetID(id string) TemplateInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the template
//
// A handle is a human friendly unique identifier for the template, unlike the ID
func (o *templateImplementation) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the template
//
// A handle is a human friendly unique identifier for the template, unlike the ID
func (o *templateImplementation) SetHandle(handle string) TemplateInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *templateImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *templateImplementation) SetMemo(memo string) TemplateInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *templateImplementation) Metas() (map[string]string, error) {
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

func (o *templateImplementation) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *templateImplementation) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *templateImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))

	return nil
}

func (o *templateImplementation) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *templateImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *templateImplementation) SetName(name string) TemplateInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *templateImplementation) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

func (o *templateImplementation) SetSiteID(siteID string) TemplateInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

func (o *templateImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *templateImplementation) SetSoftDeletedAt(softDeletedAt string) TemplateInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *templateImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *templateImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *templateImplementation) SetStatus(status string) TemplateInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *templateImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *templateImplementation) SetUpdatedAt(updatedAt string) TemplateInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *templateImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
