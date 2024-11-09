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
}
