package cmsstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/versionstore"
)

type BlockInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	ID() string
	SetID(id string) BlockInterface

	CreatedAt() string
	SetCreatedAt(createdAt string) BlockInterface
	CreatedAtCarbon() carbon.Carbon

	Content() string
	SetContent(content string) BlockInterface

	Editor() string
	SetEditor(editor string) BlockInterface

	Handle() string
	SetHandle(handle string) BlockInterface

	Memo() string
	SetMemo(memo string) BlockInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) BlockInterface

	PageID() string
	SetPageID(pageID string) BlockInterface

	ParentID() string
	SetParentID(parentID string) BlockInterface

	Sequence() string
	SequenceInt() int
	SetSequenceInt(sequence int) BlockInterface
	SetSequence(sequence string) BlockInterface

	SiteID() string
	SetSiteID(siteID string) BlockInterface

	TemplateID() string
	SetTemplateID(templateID string) BlockInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) BlockInterface
	SoftDeletedAtCarbon() carbon.Carbon

	// Status returns the status of the block, i.e. BLOCK_STATUS_ACTIVE
	Status() string

	// SetStatus sets the status of the block, i.e. BLOCK_STATUS_ACTIVE
	SetStatus(status string) BlockInterface

	// Type returns the type of the block, i.e. "text"
	Type() string

	// SetType sets the type of the block, i.e. "text"
	SetType(blockType string) BlockInterface

	// UpdatedAt returns the last updated time of block
	UpdatedAt() string

	// SetUpdatedAt sets the last updated time of block
	SetUpdatedAt(updatedAt string) BlockInterface

	// UpdatedAtCarbon returns carbon.Carbon of the last updated time of block
	UpdatedAtCarbon() carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}

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

type PageInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Methods

	MarshalToVersioning() (string, error)

	// Setters and Getters

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

	SiteID() string
	SetSiteID(siteID string) PageInterface

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
	IsSoftDeleted() bool
}

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

type TemplateInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	ID() string
	SetID(id string) TemplateInterface

	CreatedAt() string
	SetCreatedAt(createdAt string) TemplateInterface
	CreatedAtCarbon() carbon.Carbon

	Content() string
	SetContent(content string) TemplateInterface

	Editor() string
	SetEditor(editor string) TemplateInterface

	Handle() string
	SetHandle(handle string) TemplateInterface

	Memo() string
	SetMemo(memo string) TemplateInterface

	Meta(key string) string
	SetMeta(key, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Name() string
	SetName(name string) TemplateInterface

	SiteID() string
	SetSiteID(siteID string) TemplateInterface

	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) TemplateInterface
	SoftDeletedAtCarbon() carbon.Carbon

	Status() string
	SetStatus(status string) TemplateInterface

	UpdatedAt() string
	SetUpdatedAt(updatedAt string) TemplateInterface
	UpdatedAtCarbon() carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsSoftDeleted() bool
}

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

type VersioningInterface interface {
	versionstore.VersionInterface
}

func NewVersioning() VersioningInterface {
	return versionstore.NewVersion()
}

type VersioningQueryInterface interface {
	versionstore.VersionQueryInterface
}

func NewVersioningQuery() VersioningQueryInterface {
	return versionstore.NewVersionQuery()
}
