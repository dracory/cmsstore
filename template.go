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

type template struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*template)(nil)
var _ TemplateInterface = (*template)(nil)

// == CONSTRUCTORS ==========================================================

func NewTemplate() TemplateInterface {
	o := &template{}
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

func NewTemplateFromExistingData(data map[string]string) *template {
	o := &template{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *template) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *template) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *template) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *template) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *template) SetCreatedAt(createdAt string) TemplateInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *template) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *template) Content() string {
	return o.Get(COLUMN_CONTENT)
}

func (o *template) SetContent(content string) TemplateInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

func (o *template) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

func (o *template) SetEditor(editor string) TemplateInterface {
	o.Set(COLUMN_EDITOR, editor)
	return o
}

func (o *template) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *template) SetID(id string) TemplateInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the template
//
// A handle is a human friendly unique identifier for the template, unlike the ID
func (o *template) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the template
//
// A handle is a human friendly unique identifier for the template, unlike the ID
func (o *template) SetHandle(handle string) TemplateInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *template) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *template) SetMemo(memo string) TemplateInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *template) Metas() (map[string]string, error) {
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

func (o *template) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *template) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *template) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

func (o *template) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *template) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *template) SetName(name string) TemplateInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *template) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

func (o *template) SetSiteID(siteID string) TemplateInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

func (o *template) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *template) SetSoftDeletedAt(softDeletedAt string) TemplateInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *template) SoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *template) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *template) SetStatus(status string) TemplateInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *template) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *template) SetUpdatedAt(updatedAt string) TemplateInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *template) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
