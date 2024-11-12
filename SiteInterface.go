package cmsstore

import "github.com/golang-module/carbon/v2"

type SiteInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	SetCreatedAt(createdAt string) SiteInterface
	CreatedAtCarbon() carbon.Carbon

	DomainNames() ([]string, error)
	SetDomainNames(domainNames []string) (SiteInterface, error)

	Handle() string
	SetHandle(handle string) SiteInterface

	ID() string
	SetID(id string) SiteInterface

	Memo() string
	SetMemo(memo string) SiteInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) SiteInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) SiteInterface
	SoftDeletedAtCarbon() carbon.Carbon

	Status() string
	SetStatus(status string) SiteInterface

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) SiteInterface
	UpdatedAtCarbon() carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}
