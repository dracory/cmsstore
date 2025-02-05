package cmsstore

type MenuQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) MenuQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) MenuQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) MenuQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) MenuQueryInterface

	HasHandle() bool
	Handle() string
	SetHandle(handle string) MenuQueryInterface

	HasID() bool
	ID() string
	SetID(id string) MenuQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) MenuQueryInterface

	HasNameLike() bool
	NameLike() string
	SetNameLike(nameLike string) MenuQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) MenuQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) MenuQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) MenuQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) MenuQueryInterface

	HasSiteID() bool
	SiteID() string
	SetSiteID(siteID string) MenuQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(includeSoftDeleted bool) MenuQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) MenuQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) MenuQueryInterface
}
