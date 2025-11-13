package cmsstore

import (
	"encoding/json"
	"strings"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// page represents a page in the CMS system.
type pageImplementation struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*page)(nil)
var _ PageInterface = (*pageImplementation)(nil)

// == CONSTRUCTORS ==========================================================

// NewPage creates a new page with default values.
func NewPage() PageInterface {
	o := &pageImplementation{}
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
func NewPageFromExistingData(data map[string]string) *pageImplementation {
	o := &pageImplementation{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

// IsActive checks if the page is active.
func (o *pageImplementation) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

// IsInactive checks if the page is inactive.
func (o *pageImplementation) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

// IsSoftDeleted checks if the page is soft deleted.
func (o *pageImplementation) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// MarshalToVersioning marshals the page data to a versioned JSON string, excluding timestamps and soft delete information.
func (o *pageImplementation) MarshalToVersioning() (string, error) {
	versionedData := map[string]string{}

	for k, v := range o.Data() {
		if k == COLUMN_CREATED_AT ||
			k == COLUMN_UPDATED_AT ||
			k == COLUMN_SOFT_DELETED_AT {
			continue
		}
		versionedData[k] = v
	}

	b, err := json.Marshal(versionedData)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// == SETTERS AND GETTERS =====================================================

// Alias returns the alias of the page.
func (o *pageImplementation) Alias() string {
	return o.Get(COLUMN_ALIAS)
}

// SetAlias sets the alias of the page.
func (o *pageImplementation) SetAlias(alias string) PageInterface {
	o.Set(COLUMN_ALIAS, alias)
	return o
}

// CreatedAt returns the creation timestamp of the page.
func (o *pageImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp of the page.
func (o *pageImplementation) SetCreatedAt(createdAt string) PageInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// CreatedAtCarbon returns the creation timestamp of the page as a Carbon object.
func (o *pageImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

// CanonicalUrl returns the canonical URL of the page.
func (o *pageImplementation) CanonicalUrl() string {
	return o.Get(COLUMN_CANONICAL_URL)
}

// SetCanonicalUrl sets the canonical URL of the page.
func (o *pageImplementation) SetCanonicalUrl(canonicalUrl string) PageInterface {
	o.Set(COLUMN_CANONICAL_URL, canonicalUrl)
	return o
}

// Content returns the content of the page.
func (o *pageImplementation) Content() string {
	return o.Get(COLUMN_CONTENT)
}

// SetContent sets the content of the page.
func (o *pageImplementation) SetContent(content string) PageInterface {
	o.Set(COLUMN_CONTENT, content)
	return o
}

// Editor returns the editor of the page.
func (o *pageImplementation) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

// SetEditor sets the editor of the page.
func (o *pageImplementation) SetEditor(editor string) PageInterface {
	o.Set(COLUMN_EDITOR, editor)
	return o
}

// ID returns the ID of the page.
func (o *pageImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the ID of the page.
func (o *pageImplementation) SetID(id string) PageInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the page.
//
// A handle is a human-friendly unique identifier for the page, unlike the ID.
func (o *pageImplementation) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the page.
//
// A handle is a human-friendly unique identifier for the page, unlike the ID.
func (o *pageImplementation) SetHandle(handle string) PageInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

// Memo returns the memo of the page.
func (o *pageImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

// SetMemo sets the memo of the page.
func (o *pageImplementation) SetMemo(memo string) PageInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// MetaDescription returns the meta description of the page.
func (o *pageImplementation) MetaDescription() string {
	return o.Get(COLUMN_META_DESCRIPTION)
}

// SetMetaDescription sets the meta description of the page.
func (o *pageImplementation) SetMetaDescription(metaDescription string) PageInterface {
	o.Set(COLUMN_META_DESCRIPTION, metaDescription)
	return o
}

// MetaKeywords returns the meta keywords of the page.
func (o *pageImplementation) MetaKeywords() string {
	return o.Get(COLUMN_META_KEYWORDS)
}

// SetMetaKeywords sets the meta keywords of the page.
func (o *pageImplementation) SetMetaKeywords(metaKeywords string) PageInterface {
	o.Set(COLUMN_META_KEYWORDS, metaKeywords)
	return o
}

// MetaRobots returns the meta robots of the page.
func (o *pageImplementation) MetaRobots() string {
	return o.Get(COLUMN_META_ROBOTS)
}

// SetMetaRobots sets the meta robots of the page.
func (o *pageImplementation) SetMetaRobots(metaRobots string) PageInterface {
	o.Set(COLUMN_META_ROBOTS, metaRobots)
	return o
}

// Metas returns the metas of the page as a map.
func (o *pageImplementation) Metas() (map[string]string, error) {
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

// Meta returns the value of a specific meta key.
func (o *pageImplementation) Meta(name string) string {
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
func (o *pageImplementation) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as a JSON string.
// Warning: it overwrites any existing metas.
func (o *pageImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))

	return nil
}

// UpsertMetas updates or inserts metas into the existing metas.
func (o *pageImplementation) UpsertMetas(metas map[string]string) error {
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
func (o *pageImplementation) MiddlewaresBefore() []string {
	s := o.Get(COLUMN_MIDDLEWARES_BEFORE)
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

// SetMiddlewaresBefore sets the middlewares that run before the page.
func (o *pageImplementation) SetMiddlewaresBefore(middlewaresBefore []string) PageInterface {
	s := strings.Join(middlewaresBefore, ",")
	o.Set(COLUMN_MIDDLEWARES_BEFORE, s)
	return o
}

// MiddlewaresAfter returns the middlewares that run after the page.
func (o *pageImplementation) MiddlewaresAfter() []string {
	s := o.Get(COLUMN_MIDDLEWARES_AFTER)
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

// SetMiddlewaresAfter sets the middlewares that run after the page.
func (o *pageImplementation) SetMiddlewaresAfter(middlewaresAfter []string) PageInterface {
	s := strings.Join(middlewaresAfter, ",")
	o.Set(COLUMN_MIDDLEWARES_AFTER, s)
	return o
}

// Name returns the name of the page.
func (o *pageImplementation) Name() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the name of the page.
func (o *pageImplementation) SetName(name string) PageInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// SiteID returns the site ID of the page.
func (o *pageImplementation) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

// SetSiteID sets the site ID of the page.
func (o *pageImplementation) SetSiteID(siteID string) PageInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

// SoftDeletedAt returns the soft delete timestamp of the page.
func (o *pageImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// SetSoftDeletedAt sets the soft delete timestamp of the page.
func (o *pageImplementation) SetSoftDeletedAt(softDeletedAt string) PageInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

// SoftDeletedAtCarbon returns the soft delete timestamp of the page as a Carbon object.
func (o *pageImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

// Status returns the status of the page.
func (o *pageImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

// SetStatus sets the status of the page.
func (o *pageImplementation) SetStatus(status string) PageInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

// Title returns the title of the page.
func (o *pageImplementation) Title() string {
	return o.Get(COLUMN_TITLE)
}

// SetTitle sets the title of the page.
func (o *pageImplementation) SetTitle(title string) PageInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

// TemplateID returns the template ID of the page.
func (o *pageImplementation) TemplateID() string {
	return o.Get(COLUMN_TEMPLATE_ID)
}

// SetTemplateID sets the template ID of the page.
func (o *pageImplementation) SetTemplateID(templateID string) PageInterface {
	o.Set(COLUMN_TEMPLATE_ID, templateID)
	return o
}

// UpdatedAt returns the update timestamp of the page.
func (o *pageImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the update timestamp of the page.
func (o *pageImplementation) SetUpdatedAt(updatedAt string) PageInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// UpdatedAtCarbon returns the update timestamp of the page as a Carbon object.
func (o *pageImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
