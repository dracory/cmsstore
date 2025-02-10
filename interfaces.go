package cmsstore

import (
	"context"

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

	MiddlewaresBefore() []string
	SetMiddlewaresBefore(middlewaresBefore []string) PageInterface

	MiddlewaresAfter() []string
	SetMiddlewaresAfter(middlewaresAfter []string) PageInterface

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

type StoreInterface interface {
	AutoMigrate(ctx context.Context, opts ...Option) error
	EnableDebug(debug bool)

	BlockCreate(ctx context.Context, block BlockInterface) error
	BlockCount(ctx context.Context, options BlockQueryInterface) (int64, error)
	BlockDelete(ctx context.Context, block BlockInterface) error
	BlockDeleteByID(ctx context.Context, id string) error
	BlockFindByHandle(ctx context.Context, blockHandle string) (BlockInterface, error)
	BlockFindByID(ctx context.Context, blockID string) (BlockInterface, error)
	BlockList(ctx context.Context, query BlockQueryInterface) ([]BlockInterface, error)
	BlockSoftDelete(ctx context.Context, block BlockInterface) error
	BlockSoftDeleteByID(ctx context.Context, id string) error
	BlockUpdate(ctx context.Context, block BlockInterface) error

	MenusEnabled() bool

	MenuCreate(ctx context.Context, menu MenuInterface) error
	MenuCount(ctx context.Context, options MenuQueryInterface) (int64, error)
	MenuDelete(ctx context.Context, menu MenuInterface) error
	MenuDeleteByID(ctx context.Context, id string) error
	MenuFindByHandle(ctx context.Context, menuHandle string) (MenuInterface, error)
	MenuFindByID(ctx context.Context, menuID string) (MenuInterface, error)
	MenuList(ctx context.Context, query MenuQueryInterface) ([]MenuInterface, error)
	MenuSoftDelete(ctx context.Context, menu MenuInterface) error
	MenuSoftDeleteByID(ctx context.Context, id string) error
	MenuUpdate(ctx context.Context, menu MenuInterface) error

	MenuItemCreate(ctx context.Context, menuItem MenuItemInterface) error
	MenuItemCount(ctx context.Context, options MenuItemQueryInterface) (int64, error)
	MenuItemDelete(ctx context.Context, menuItem MenuItemInterface) error
	MenuItemDeleteByID(ctx context.Context, id string) error
	MenuItemFindByID(ctx context.Context, menuItemID string) (MenuItemInterface, error)
	MenuItemList(ctx context.Context, query MenuItemQueryInterface) ([]MenuItemInterface, error)
	MenuItemSoftDelete(ctx context.Context, menuItem MenuItemInterface) error
	MenuItemSoftDeleteByID(ctx context.Context, id string) error
	MenuItemUpdate(ctx context.Context, menuItem MenuItemInterface) error

	PageCreate(ctx context.Context, page PageInterface) error
	PageCount(ctx context.Context, options PageQueryInterface) (int64, error)
	PageDelete(ctx context.Context, page PageInterface) error
	PageDeleteByID(ctx context.Context, id string) error
	PageFindByHandle(ctx context.Context, pageHandle string) (PageInterface, error)
	PageFindByID(ctx context.Context, pageID string) (PageInterface, error)
	PageList(ctx context.Context, query PageQueryInterface) ([]PageInterface, error)
	PageSoftDelete(ctx context.Context, page PageInterface) error
	PageSoftDeleteByID(ctx context.Context, id string) error
	PageUpdate(ctx context.Context, page PageInterface) error

	SiteCreate(ctx context.Context, site SiteInterface) error
	SiteCount(ctx context.Context, options SiteQueryInterface) (int64, error)
	SiteDelete(ctx context.Context, site SiteInterface) error
	SiteDeleteByID(ctx context.Context, id string) error
	SiteFindByDomainName(ctx context.Context, siteDomainName string) (SiteInterface, error)
	SiteFindByHandle(ctx context.Context, siteHandle string) (SiteInterface, error)
	SiteFindByID(ctx context.Context, siteID string) (SiteInterface, error)
	SiteList(ctx context.Context, query SiteQueryInterface) ([]SiteInterface, error)
	SiteSoftDelete(ctx context.Context, site SiteInterface) error
	SiteSoftDeleteByID(ctx context.Context, id string) error
	SiteUpdate(ctx context.Context, site SiteInterface) error

	TemplateCreate(ctx context.Context, template TemplateInterface) error
	TemplateCount(ctx context.Context, options TemplateQueryInterface) (int64, error)
	TemplateDelete(ctx context.Context, template TemplateInterface) error
	TemplateDeleteByID(ctx context.Context, id string) error
	TemplateFindByHandle(ctx context.Context, templateHandle string) (TemplateInterface, error)
	TemplateFindByID(ctx context.Context, templateID string) (TemplateInterface, error)
	TemplateList(ctx context.Context, query TemplateQueryInterface) ([]TemplateInterface, error)
	TemplateSoftDelete(ctx context.Context, template TemplateInterface) error
	TemplateSoftDeleteByID(ctx context.Context, id string) error
	TemplateUpdate(ctx context.Context, template TemplateInterface) error

	TranslationsEnabled() bool

	TranslationCreate(ctx context.Context, translation TranslationInterface) error
	TranslationCount(ctx context.Context, options TranslationQueryInterface) (int64, error)
	TranslationDelete(ctx context.Context, translation TranslationInterface) error
	TranslationDeleteByID(ctx context.Context, id string) error
	TranslationFindByHandle(ctx context.Context, translationHandle string) (TranslationInterface, error)
	TranslationFindByHandleOrID(ctx context.Context, translationHandleOrID string, language string) (TranslationInterface, error)
	TranslationFindByID(ctx context.Context, translationID string) (TranslationInterface, error)
	TranslationList(ctx context.Context, query TranslationQueryInterface) ([]TranslationInterface, error)
	TranslationSoftDelete(ctx context.Context, translation TranslationInterface) error
	TranslationSoftDeleteByID(ctx context.Context, id string) error
	TranslationUpdate(ctx context.Context, translation TranslationInterface) error
	TranslationLanguageDefault() string
	TranslationLanguages() map[string]string

	// Versioning
	VersioningEnabled() bool
	VersioningCreate(ctx context.Context, versioning VersioningInterface) error
	// VersioningCount(options VersioningQueryInterface) (int64, error)
	VersioningDelete(ctx context.Context, versioning VersioningInterface) error
	VersioningDeleteByID(ctx context.Context, id string) error
	VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error)
	VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error)
	VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error
	VersioningSoftDeleteByID(ctx context.Context, id string) error
	VersioningUpdate(ctx context.Context, versioning VersioningInterface) error

	Shortcodes() []ShortcodeInterface
	AddShortcode(shortcode ShortcodeInterface)
	AddShortcodes(shortcodes []ShortcodeInterface)
	SetShortcodes(shortcodes []ShortcodeInterface)

	Middlewares() []MiddlewareInterface
	AddMiddleware(middleware MiddlewareInterface)
	AddMiddlewares(middlewares []MiddlewareInterface)
	SetMiddlewares(middlewares []MiddlewareInterface)
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
