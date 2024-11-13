package cmsstore

import "github.com/golang-module/carbon/v2"

type MenuItemInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	SetCreatedAt(createdAt string) MenuItemInterface
	CreatedAtCarbon() carbon.Carbon

	Handle() string
	SetHandle(handle string) MenuItemInterface

	ID() string
	SetID(id string) MenuItemInterface

	Memo() string
	SetMemo(memo string) MenuItemInterface

	MenuID() string
	SetMenuID(menuID string) MenuItemInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) MenuItemInterface

	PageID() string
	SetPageID(pageID string) MenuItemInterface

	ParentID() string
	SetParentID(parentID string) MenuItemInterface

	Sequence() string
	SequenceInt() int
	SetSequence(sequence string) MenuItemInterface
	SetSequenceInt(sequence int) MenuItemInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) MenuItemInterface
	SoftDeletedAtCarbon() carbon.Carbon

	Status() string
	SetStatus(status string) MenuItemInterface

	Target() string
	SetTarget(target string) MenuItemInterface

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) MenuItemInterface
	UpdatedAtCarbon() carbon.Carbon

	URL() string
	SetURL(url string) MenuItemInterface

	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}
