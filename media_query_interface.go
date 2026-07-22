package cmsstore

type MediaQueryInterface interface {
	Validate() error

	Columns() []string
	HasColumns() bool
	SetColumns(columns []string) MediaQueryInterface

	HasID() bool
	ID() string
	SetID(id string) MediaQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) MediaQueryInterface

	HasEntityID() bool
	EntityID() string
	SetEntityID(entityID string) MediaQueryInterface

	HasEntityType() bool
	EntityType() string
	SetEntityType(entityType string) MediaQueryInterface

	HasSiteID() bool
	SiteID() string
	SetSiteID(siteID string) MediaQueryInterface

	HasHandle() bool
	Handle() string
	SetHandle(handle string) MediaQueryInterface

	HasExtension() bool
	Extension() string
	SetExtension(extension string) MediaQueryInterface

	HasType() bool
	Type() string
	SetType(mediaType string) MediaQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) MediaQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) MediaQueryInterface

	HasNameLike() bool
	NameLike() string
	SetNameLike(nameLike string) MediaQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) MediaQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) MediaQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) MediaQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) MediaQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) MediaQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeleteIncluded bool) MediaQueryInterface
}
