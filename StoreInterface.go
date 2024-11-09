package cmsstore

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)
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
}
