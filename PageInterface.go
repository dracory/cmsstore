package cmsstore

import "github.com/golang-module/carbon/v2"

type PageInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	ID() string
	SetID(id string) PageInterface

	Alias() string
	SetAlias(alias string) PageInterface

	CreatedAt() string
	SetCreatedAt(createdAt string) PageInterface
	CreatedAtCarbon() carbon.Carbon

	CanonicalUrl() string
	SetCanonicalUrl(canonicalUrl string) PageInterface

	Content() string
	SetContent(content string) PageInterface

	Editor() string
	SetEditor(editor string) PageInterface

	Handle() string
	SetHandle(handle string) PageInterface

	Memo() string
	SetMemo(memo string) PageInterface

	MetaDescription() string
	SetMetaDescription(metaDescription string) PageInterface

	MetaKeywords() string
	SetMetaKeywords(metaKeywords string) PageInterface

	MetaRobots() string
	SetMetaRobots(metaRobots string) PageInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) PageInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) PageInterface
	SoftDeletedAtCarbon() carbon.Carbon

	Status() string
	SetStatus(status string) PageInterface

	Title() string
	SetTitle(title string) PageInterface

	TemplateID() string
	SetTemplateID(templateID string) PageInterface

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) PageInterface
	UpdatedAtCarbon() carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsDeleted() bool
}
