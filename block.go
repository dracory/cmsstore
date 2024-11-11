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

type block struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*block)(nil)
var _ BlockInterface = (*block)(nil)

// == CONSTRUCTORS ==========================================================

func NewBlock() BlockInterface {
	o := &block{}
	o.SetContent("")
	o.SetEditor("")
	o.SetHandle("")
	o.SetID(uid.HumanUid())
	o.SetMemo("")
	o.SetMetas(map[string]string{})
	o.SetName("")
	o.SetStatus(BLOCK_STATUS_DRAFT)
	o.SetType("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)
	return o
}

func NewBlockFromExistingData(data map[string]string) *block {
	o := &block{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *block) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *block) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *block) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *block) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *block) SetCreatedAt(createdAt string) BlockInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *block) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *block) Content() string {
	return o.Get(COLUMN_CONTENT)
}

func (o *block) SetContent(content string) BlockInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

func (o *block) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

func (o *block) SetEditor(editor string) BlockInterface {
	o.Set(COLUMN_EDITOR, editor)
	return o
}

func (o *block) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *block) SetID(id string) BlockInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the block
//
// A handle is a human friendly unique identifier for the block, unlike the ID
func (o *block) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the block
//
// A handle is a human friendly unique identifier for the block, unlike the ID
func (o *block) SetHandle(handle string) BlockInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *block) Memo() string {
	return o.Get("memo")
}

func (o *block) SetMemo(memo string) BlockInterface {
	o.Set("memo", memo)
	return o
}

func (o *block) Metas() (map[string]string, error) {
	metasStr := o.Get("metas")

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
}

func (o *block) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *block) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *block) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

func (o *block) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *block) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *block) SetName(name string) BlockInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *block) PageID() string {
	return o.Get(COLUMN_PAGE_ID)
}

func (o *block) SetPageID(pageID string) BlockInterface {
	o.Set(COLUMN_PAGE_ID, pageID)
	return o
}

func (o *block) ParentID() string {
	return o.Get(COLUMN_PARENT_ID)
}

func (o *block) SetParentID(parentID string) BlockInterface {
	o.Set(COLUMN_PARENT_ID, parentID)
	return o
}

func (o *block) Sequence() string {
	return o.Get(COLUMN_SEQUENCE)
}

func (o *block) SetSequence(sequence string) BlockInterface {
	o.Set(COLUMN_SEQUENCE, sequence)
	return o
}

func (o *block) SequenceInt() int {
	return cast.ToInt(o.Sequence())
}

func (o *block) SetSequenceInt(sequence int) BlockInterface {
	o.Set(COLUMN_SEQUENCE, cast.ToString(sequence))
	return o
}

func (o *block) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

func (o *block) SetSiteID(siteID string) BlockInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

func (o *block) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *block) SetSoftDeletedAt(softDeletedAt string) BlockInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *block) SoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *block) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *block) SetStatus(status string) BlockInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *block) TemplateID() string {
	return o.Get(COLUMN_TEMPLATE_ID)
}

func (o *block) SetTemplateID(templateID string) BlockInterface {
	o.Set(COLUMN_TEMPLATE_ID, templateID)
	return o
}

// Type returns the type of the block, i.e. "text"
func (o *block) Type() string {
	return o.Get(COLUMN_TYPE)
}

// SetType sets the type of the block, i.e. "text"
func (o *block) SetType(blockType string) BlockInterface {
	o.Set(COLUMN_TYPE, blockType)
	return o
}

func (o *block) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *block) SetUpdatedAt(updatedAt string) BlockInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *block) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
