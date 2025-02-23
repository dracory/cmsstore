package cmsstore

import (
	"strings"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

// == TYPE ===================================================================

// page represents a page in the CMS system.
type page struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*page)(nil)
var _ PageInterface = (*page)(nil)

// == CONSTRUCTORS ==========================================================

// NewPage creates a new page with default values.
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
	o.SetMiddlewaresAfter([]string{})
	o.SetMiddlewaresBefore([]string{})
	o.SetName("")
	o.SetStatus(PAGE_STATUS_DRAFT)
	o.SetTemplateID("")
	o.SetTitle("")
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(sb.MAX_DATETIME)
	return o
}

// NewPageFromExistingData creates a new page from existing data.
func NewPageFromExistingData(data map[string]string) *page {
	o := &page{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

// IsActive checks if the page is active.
func (o *page) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

// IsInactive checks if the page is inactive.
func (o *page) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

// IsSoftDeleted checks if the page is soft deleted.
func (o *page) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// MarshalToVersioning marshals the page data to a versioned JSON string, excluding timestamps and soft delete information.
func (o *page) MarshalToVersioning() (string, error) {
	versionedData := map[string]string{}

	for k, v := range o.Data() {
		if k == COLUMN_CREATED_AT ||
			k == COLUMN_UPDATED_AT ||
			k == COLUMN_SOFT_DELETED_AT {
			continue
		}
		versionedData[k] = v
	}

	return utils.ToJSON(versionedData)
}

// == SETTERS AND GETTERS =====================================================

// Alias returns the alias of the page.
func (o *page) Alias() string {
	return o.Get(COLUMN_ALIAS)
}

// SetAlias sets the alias of the page.
func (o *page) SetAlias(alias string) PageInterface {
	o.Set(COLUMN_ALIAS, alias)
	return o
}

// CreatedAt returns the creation timestamp of the page.
func (o *page) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp of the page.
func (o *page) SetCreatedAt(createdAt string) PageInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// CreatedAtCarbon returns the creation timestamp of the page as a Carbon object.
func (o *page) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

// CanonicalUrl returns the canonical URL of the page.
func (o *page) CanonicalUrl() string {
	return o.Get(COLUMN_CANONICAL_URL)
}

// SetCanonicalUrl sets the canonical URL of the page.
func (o *page) SetCanonicalUrl(canonicalUrl string) PageInterface {
	o.Set(COLUMN_CANONICAL_URL, canonicalUrl)
	return o
}

// Content returns the content of the page.
func (o *page) Content() string {
	return o.Get(COLUMN_CONTENT)
}

// SetContent sets the content of the page.
func (o *page) SetContent(content string) PageInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

// Editor returns the editor of the page.
func (o *page) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

// SetEditor sets the editor of the page.
func (o *page) SetEditor(editor string) PageInterface {
	o.Set(COLUMN_EDITOR, editor)
	return o
}

// ID returns the ID of the page.
func (o *page) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the ID of the page.
func (o *page) SetID(id string) PageInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the page.
//
// A handle is a human-friendly unique identifier for the page, unlike the ID.
func (o *page) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the page.
//
// A handle is a human-friendly unique identifier for the page, unlike the ID.
func (o *page) SetHandle(handle string) PageInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

// Memo returns the memo of the page.
func (o *page) Memo() string {
	return o.Get(COLUMN_MEMO)
}

// SetMemo sets the memo of the page.
func (o *page) SetMemo(memo string) PageInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// MetaDescription returns the meta description of the page.
func (o *page) MetaDescription() string {
	return o.Get(COLUMN_META_DESCRIPTION)
}

// SetMetaDescription sets the meta description of the page.
func (o *page) SetMetaDescription(metaDescription string) PageInterface {
	o.Set(COLUMN_META_DESCRIPTION, metaDescription)
	return o
}

// MetaKeywords returns the meta keywords of the page.
func (o *page) MetaKeywords() string {
	return o.Get(COLUMN_META_KEYWORDS)
}

// SetMetaKeywords sets the meta keywords of the page.
func (o *page) SetMetaKeywords(metaKeywords string) PageInterface {
	o.Set(COLUMN_META_KEYWORDS, metaKeywords)
	return o
}

// MetaRobots returns the meta robots of the page.
func (o *page) MetaRobots() string {
	return o.Get(COLUMN_META_ROBOTS)
}

// SetMetaRobots sets the meta robots of the page.
func (o *page) SetMetaRobots(metaRobots string) PageInterface {
	o.Set(COLUMN_META_ROBOTS, metaRobots)
	return o
}

// Metas returns the metas of the page as a map.
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

// Meta returns the value of a specific meta key.
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

// SetMeta sets the value of a specific meta key.
func (o *page) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as a JSON string.
// Warning: it overwrites any existing metas.
func (o *page) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

// UpsertMetas updates or inserts metas into the existing metas.
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

// MiddlewaresBefore returns the middlewares that run before the page.
func (o *page) MiddlewaresBefore() []string {
	s := o.Get(COLUMN_MIDDLEWARES_BEFORE)
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

// SetMiddlewaresBefore sets the middlewares that run before the page.
func (o *page) SetMiddlewaresBefore(middlewaresBefore []string) PageInterface {
	s := strings.Join(middlewaresBefore, ",")
	o.Set(COLUMN_MIDDLEWARES_BEFORE, s)
	return o
}

// MiddlewaresAfter returns the middlewares that run after the page.
func (o *page) MiddlewaresAfter() []string {
	s := o.Get(COLUMN_MIDDLEWARES_AFTER)
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

// SetMiddlewaresAfter sets the middlewares that run after the page.
func (o *page) SetMiddlewaresAfter(middlewaresAfter []string) PageInterface {
	s := strings.Join(middlewaresAfter, ",")
	o.Set(COLUMN_MIDDLEWARES_AFTER, s)
	return o
}

// Name returns the name of the page.
func (o *page) Name() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the name of the page.
func (o *page) SetName(name string) PageInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// SiteID returns the site ID of the page.
func (o *page) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

// SetSiteID sets the site ID of the page.
func (o *page) SetSiteID(siteID string) PageInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

// SoftDeletedAt returns the soft delete timestamp of the page.
func (o *page) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// SetSoftDeletedAt sets the soft delete timestamp of the page.
func (o *page) SetSoftDeletedAt(softDeletedAt string) PageInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

// SoftDeletedAtCarbon returns the soft delete timestamp of the page as a Carbon object.
func (o *page) SoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

// Status returns the status of the page.
func (o *page) Status() string {
	return o.Get(COLUMN_STATUS)
}

// SetStatus sets the status of the page.
func (o *page) SetStatus(status string) PageInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

// Title returns the title of the page.
func (o *page) Title() string {
	return o.Get(COLUMN_TITLE)
}

// SetTitle sets the title of the page.
func (o *page) SetTitle(title string) PageInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

// TemplateID returns the template ID of the page.
func (o *page) TemplateID() string {
	return o.Get(COLUMN_TEMPLATE_ID)
}

// SetTemplateID sets the template ID of the page.
func (o *page) SetTemplateID(templateID string) PageInterface {
	o.Set(COLUMN_TEMPLATE_ID, templateID)
	return o
}

// UpdatedAt returns the update timestamp of the page.
func (o *page) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the update timestamp of the page.
func (o *page) SetUpdatedAt(updatedAt string) PageInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// UpdatedAtCarbon returns the update timestamp of the page as a Carbon object.
func (o *page) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
