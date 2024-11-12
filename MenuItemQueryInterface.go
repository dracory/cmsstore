package cmsstore

type MenuItemQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) MenuItemQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) MenuItemQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) MenuItemQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) MenuItemQueryInterface

	HasID() bool
	ID() string
	SetID(id string) MenuItemQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) MenuItemQueryInterface

	HasNameLike() bool
	NameLike() string
	SetNameLike(nameLike string) MenuItemQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) MenuItemQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) MenuItemQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) MenuItemQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) MenuItemQueryInterface

	HasSiteID() bool
	SiteID() string
	SetSiteID(siteID string) MenuItemQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(includeSoftDeleted bool) MenuItemQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) MenuItemQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) MenuItemQueryInterface
}
