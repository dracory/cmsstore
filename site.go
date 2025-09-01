package cmsstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == TYPE ===================================================================

// site represents a site in the CMS store.
type site struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// Ensure that site implements SiteInterface.
var _ SiteInterface = (*site)(nil)

// == CONSTRUCTORS ==========================================================

// NewSite creates a new site with default values.
func NewSite() SiteInterface {
	o := &site{}
	o.SetDomainNames([]string{})
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

// NewSiteFromExistingData creates a new site from existing data.
func NewSiteFromExistingData(data map[string]string) *site {
	o := &site{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

// IsActive checks if the site is active.
func (o *site) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

// IsInactive checks if the site is inactive.
func (o *site) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

// IsSoftDeleted checks if the site is soft-deleted.
func (o *site) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

// CreatedAt returns the creation timestamp of the site.
func (o *site) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the creation timestamp of the site.
func (o *site) SetCreatedAt(createdAt string) SiteInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

// CreatedAtCarbon returns the creation timestamp of the site as a Carbon object.
func (o *site) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

// DomainNames returns the domain names associated with the site.
func (o *site) DomainNames() ([]string, error) {
	domainNamesStr := o.Get(COLUMN_DOMAIN_NAMES)

	if domainNamesStr == "" {
		domainNamesStr = "[]"
	}

	domainNamesJson := []string{}
	errJson := json.Unmarshal([]byte(domainNamesStr), &domainNamesJson)
	if errJson != nil {
		return []string{}, errJson
	}

	if domainNamesJson == nil {
		return []string{}, nil
	}

	return domainNamesJson, nil
}

// SetDomainNames sets the domain names associated with the site.
func (o *site) SetDomainNames(domainNames []string) (SiteInterface, error) {
	domainNamesBytes, err := json.Marshal(domainNames)
	if err != nil {
		return o, err
	}
	o.Set(COLUMN_DOMAIN_NAMES, string(domainNamesBytes))
	return o, nil
}

// ID returns the unique identifier of the site.
func (o *site) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the site.
func (o *site) SetID(id string) SiteInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the site.
//
// A handle is a human-friendly unique identifier for the site, unlike the ID.
func (o *site) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the site.
//
// A handle is a human-friendly unique identifier for the site, unlike the ID.
func (o *site) SetHandle(handle string) SiteInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

// Memo returns the memo associated with the site.
func (o *site) Memo() string {
	return o.Get(COLUMN_MEMO)
}

// SetMemo sets the memo associated with the site.
func (o *site) SetMemo(memo string) SiteInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// Metas returns the metadata associated with the site.
func (o *site) Metas() (map[string]string, error) {
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

// Meta returns a specific metadata field.
func (o *site) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

// SetMeta sets a specific metadata field.
func (o *site) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metadata as a JSON string.
// Warning: it overwrites any existing metadata.
func (o *site) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))

	return nil
}

// UpsertMetas updates or inserts metadata.
func (o *site) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

// Name returns the name of the site.
func (o *site) Name() string {
	return o.Get(COLUMN_NAME)
}

// SetName sets the name of the site.
func (o *site) SetName(name string) SiteInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

// SoftDeletedAt returns the soft-deletion timestamp of the site.
func (o *site) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// SetSoftDeletedAt sets the soft-deletion timestamp of the site.
func (o *site) SetSoftDeletedAt(softDeletedAt string) SiteInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

// SoftDeletedAtCarbon returns the soft-deletion timestamp of the site as a Carbon object.
func (o *site) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

// Status returns the status of the site.
func (o *site) Status() string {
	return o.Get(COLUMN_STATUS)
}

// SetStatus sets the status of the site.
func (o *site) SetStatus(status string) SiteInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

// UpdatedAt returns the last updated timestamp of the site.
func (o *site) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// SetUpdatedAt sets the last updated timestamp of the site.
func (o *site) SetUpdatedAt(updatedAt string) SiteInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

// UpdatedAtCarbon returns the last updated timestamp of the site as a Carbon object.
func (o *site) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
