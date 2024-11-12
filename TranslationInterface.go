package cmsstore

import "github.com/golang-module/carbon/v2"

type TranslationInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	ID() string
	SetID(id string) TranslationInterface

	CreatedAt() string
	SetCreatedAt(createdAt string) TranslationInterface
	CreatedAtCarbon() carbon.Carbon

	Content() (languageCodeContent map[string]string, err error)
	SetContent(languageCodeContent map[string]string) error

	Handle() string
	SetHandle(handle string) TranslationInterface

	Memo() string
	SetMemo(memo string) TranslationInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) TranslationInterface

	SiteID() string
	SetSiteID(siteID string) TranslationInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) TranslationInterface
	SoftDeletedAtCarbon() carbon.Carbon

	Status() string
	SetStatus(status string) TranslationInterface

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) TranslationInterface
	UpdatedAtCarbon() carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}
