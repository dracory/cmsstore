package cmsstore

import (
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

// == TYPE ===================================================================

type page struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*page)(nil)
var _ PageInterface = (*page)(nil)

// == CONSTRUCTORS ==========================================================

func NewPage() PageInterface {
	o := &page{}
	o.SetAlias("")
	o.SetCanonicalUrl("")
	o.SetContent("")
	o.SetEditor("")
	o.SetHandle("")
	o.SetID(uid.HumanUid())
	o.SetMemo("")
	o.SetMetaDescription("")
	o.SetMetaKeywords("")
	o.SetMetaRobots("")
	o.SetMetas(map[string]string{})
	o.SetName("")
	o.SetStatus(PAGE_STATUS_DRAFT)
	o.SetTemplateID("")
	o.SetTitle("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)
	return o
}

func NewPageFromExistingData(data map[string]string) *page {
	o := &page{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *page) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *page) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *page) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *page) Alias() string {
	return o.Get(COLUMN_ALIAS)
}

func (o *page) SetAlias(alias string) PageInterface {
	o.Set(COLUMN_ALIAS, alias)
	return o
}

func (o *page) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *page) SetCreatedAt(createdAt string) PageInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *page) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *page) CanonicalUrl() string {
	return o.Get(COLUMN_CANONICAL_URL)
}

func (o *page) SetCanonicalUrl(canonicalUrl string) PageInterface {
	o.Set(COLUMN_CANONICAL_URL, canonicalUrl)
	return o
}

func (o *page) Content() string {
	return o.Get(COLUMN_CONTENT)
}

func (o *page) SetContent(content string) PageInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

func (o *page) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

func (o *page) SetEditor(editor string) PageInterface {
	o.Set(COLUMN_EDITOR, editor)
	return o
}

func (o *page) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *page) SetID(id string) PageInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the page
//
// A handle is a human friendly unique identifier for the page, unlike the ID
func (o *page) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the page
//
// A handle is a human friendly unique identifier for the page, unlike the ID
func (o *page) SetHandle(handle string) PageInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *page) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *page) SetMemo(memo string) PageInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *page) MetaDescription() string {
	return o.Get(COLUMN_META_DESCRIPTION)
}

func (o *page) SetMetaDescription(metaDescription string) PageInterface {
	o.Set(COLUMN_META_DESCRIPTION, metaDescription)
	return o
}

func (o *page) MetaKeywords() string {
	return o.Get(COLUMN_META_KEYWORDS)
}

func (o *page) SetMetaKeywords(metaKeywords string) PageInterface {
	o.Set(COLUMN_META_KEYWORDS, metaKeywords)
	return o
}

func (o *page) MetaRobots() string {
	return o.Get(COLUMN_META_ROBOTS)
}

func (o *page) SetMetaRobots(metaRobots string) PageInterface {
	o.Set(COLUMN_META_ROBOTS, metaRobots)
	return o
}

func (o *page) Metas() (map[string]string, error) {
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

func (o *page) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *page) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *page) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

func (o *page) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *page) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *page) SetName(name string) PageInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *page) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

func (o *page) SetSiteID(siteID string) PageInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

func (o *page) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *page) SetSoftDeletedAt(softDeletedAt string) PageInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *page) SoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *page) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *page) SetStatus(status string) PageInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *page) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *page) SetTitle(title string) PageInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *page) TemplateID() string {
	return o.Get(COLUMN_TEMPLATE_ID)
}

func (o *page) SetTemplateID(templateID string) PageInterface {
	o.Set(COLUMN_TEMPLATE_ID, templateID)
	return o
}

func (o *page) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *page) SetUpdatedAt(updatedAt string) PageInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *page) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
