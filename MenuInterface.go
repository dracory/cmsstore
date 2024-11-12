package cmsstore

import "github.com/golang-module/carbon/v2"

type MenuInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	SetCreatedAt(createdAt string) MenuInterface
	CreatedAtCarbon() carbon.Carbon

	Handle() string
	SetHandle(handle string) MenuInterface

	ID() string
	SetID(id string) MenuInterface

	Memo() string
	SetMemo(memo string) MenuInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) MenuInterface

	SiteID() string
	SetSiteID(siteID string) MenuInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) MenuInterface
	SoftDeletedAtCarbon() carbon.Carbon

	Status() string
	SetStatus(status string) MenuInterface

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) MenuInterface
	UpdatedAtCarbon() carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}
