package cmsstore

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)

	BlockCreate(block BlockInterface) error
	BlockCount(options BlockQueryInterface) (int64, error)
	BlockDelete(block BlockInterface) error
	BlockDeleteByID(id string) error
	BlockFindByHandle(blockHandle string) (BlockInterface, error)
	BlockFindByID(blockID string) (BlockInterface, error)
	BlockList(query BlockQueryInterface) ([]BlockInterface, error)
	BlockSoftDelete(block BlockInterface) error
	BlockSoftDeleteByID(id string) error
	BlockUpdate(block BlockInterface) error

	MenusEnabled() bool

	MenuCreate(menu MenuInterface) error
	MenuCount(options MenuQueryInterface) (int64, error)
	MenuDelete(menu MenuInterface) error
	MenuDeleteByID(id string) error
	MenuFindByHandle(menuHandle string) (MenuInterface, error)
	MenuFindByID(menuID string) (MenuInterface, error)
	MenuList(query MenuQueryInterface) ([]MenuInterface, error)
	MenuSoftDelete(menu MenuInterface) error
	MenuSoftDeleteByID(id string) error
	MenuUpdate(menu MenuInterface) error

	MenuItemCreate(menuItem MenuItemInterface) error
	MenuItemCount(options MenuItemQueryInterface) (int64, error)
	MenuItemDelete(menuItem MenuItemInterface) error
	MenuItemDeleteByID(id string) error
	MenuItemFindByID(menuItemID string) (MenuItemInterface, error)
	MenuItemList(query MenuItemQueryInterface) ([]MenuItemInterface, error)
	MenuItemSoftDelete(menuItem MenuItemInterface) error
	MenuItemSoftDeleteByID(id string) error
	MenuItemUpdate(menuItem MenuItemInterface) error

	PageCreate(page PageInterface) error
	PageCount(options PageQueryInterface) (int64, error)
	PageDelete(page PageInterface) error
	PageDeleteByID(id string) error
	PageFindByHandle(pageHandle string) (PageInterface, error)
	PageFindByID(pageID string) (PageInterface, error)
	PageList(query PageQueryInterface) ([]PageInterface, error)
	PageSoftDelete(page PageInterface) error
	PageSoftDeleteByID(id string) error
	PageUpdate(page PageInterface) error

	SiteCreate(site SiteInterface) error
	SiteCount(options SiteQueryInterface) (int64, error)
	SiteDelete(site SiteInterface) error
	SiteDeleteByID(id string) error
	SiteFindByDomainName(siteDomainName string) (SiteInterface, error)
	SiteFindByHandle(siteHandle string) (SiteInterface, error)
	SiteFindByID(siteID string) (SiteInterface, error)
	SiteList(query SiteQueryInterface) ([]SiteInterface, error)
	SiteSoftDelete(site SiteInterface) error
	SiteSoftDeleteByID(id string) error
	SiteUpdate(site SiteInterface) error

	TemplateCreate(template TemplateInterface) error
	TemplateCount(options TemplateQueryInterface) (int64, error)
	TemplateDelete(template TemplateInterface) error
	TemplateDeleteByID(id string) error
	TemplateFindByHandle(templateHandle string) (TemplateInterface, error)
	TemplateFindByID(templateID string) (TemplateInterface, error)
	TemplateList(query TemplateQueryInterface) ([]TemplateInterface, error)
	TemplateSoftDelete(template TemplateInterface) error
	TemplateSoftDeleteByID(id string) error
	TemplateUpdate(template TemplateInterface) error

	TranslationsEnabled() bool

	TranslationCreate(translation TranslationInterface) error
	TranslationCount(options TranslationQueryInterface) (int64, error)
	TranslationDelete(translation TranslationInterface) error
	TranslationDeleteByID(id string) error
	TranslationFindByHandle(translationHandle string) (TranslationInterface, error)
	TranslationFindByID(translationID string) (TranslationInterface, error)
	TranslationList(query TranslationQueryInterface) ([]TranslationInterface, error)
	TranslationSoftDelete(translation TranslationInterface) error
	TranslationSoftDeleteByID(id string) error
	TranslationUpdate(translation TranslationInterface) error
	TranslationLanguageDefault() string
	TranslationLanguages() map[string]string
}
