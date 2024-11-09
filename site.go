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

type site struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

// var _ dataobject.DataObjectInterface = (*site)(nil)
var _ SiteInterface = (*site)(nil)

// == CONSTRUCTORS ==========================================================

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

func NewSiteFromExistingData(data map[string]string) *site {
	o := &site{}
	o.Hydrate(data)
	return o
}

// == METHODS ===============================================================

func (o *site) IsActive() bool {
	return o.Status() == PAGE_STATUS_ACTIVE
}

func (o *site) IsInactive() bool {
	return o.Status() == PAGE_STATUS_INACTIVE
}

func (o *site) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *site) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *site) SetCreatedAt(createdAt string) SiteInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *site) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *site) DomainNames() ([]string, error) {
	domainNamesStr := o.Get(COLUMN_DOMAIN_NAMES)

	if domainNamesStr == "" {
		domainNamesStr = "[]"
	}

	domainNamesJson, errJson := utils.FromJSON(domainNamesStr, []string{})
	if errJson != nil {
		return []string{}, errJson
	}

	return domainNamesJson.([]string), nil
}

func (o *site) SetDomainNames(domainNames []string) (SiteInterface, error) {
	domainNamesJson, errJson := utils.ToJSON(domainNames)
	if errJson != nil {
		return o, errJson
	}
	o.Set(COLUMN_DOMAIN_NAMES, domainNamesJson)
	return o, nil
}

func (o *site) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *site) SetID(id string) SiteInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// Handle returns the handle of the site
//
// A handle is a human friendly unique identifier for the site, unlike the ID
func (o *site) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

// SetHandle sets the handle of the site
//
// A handle is a human friendly unique identifier for the site, unlike the ID
func (o *site) SetHandle(handle string) SiteInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *site) Memo() string {
	return o.Get("memo")
}

func (o *site) SetMemo(memo string) SiteInterface {
	o.Set("memo", memo)
	return o
}

func (o *site) Metas() (map[string]string, error) {
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

func (o *site) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *site) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, mapString)

	return nil
}

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

func (o *site) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *site) SetName(name string) SiteInterface {
	o.Set(COLUMN_NAME, name)
	return o
}

func (o *site) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *site) SetSoftDeletedAt(softDeletedAt string) SiteInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *site) SoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func (o *site) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *site) SetStatus(status string) SiteInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *site) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *site) SetUpdatedAt(updatedAt string) SiteInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *site) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
