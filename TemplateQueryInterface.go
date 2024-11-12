package cmsstore

type TemplateQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) TemplateQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) TemplateQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) TemplateQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) TemplateQueryInterface

	HasHandle() bool
	Handle() string
	SetHandle(handle string) TemplateQueryInterface

	HasID() bool
	ID() string
	SetID(id string) TemplateQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) TemplateQueryInterface

	HasNameLike() bool
	NameLike() string
	SetNameLike(nameLike string) TemplateQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) TemplateQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) TemplateQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) TemplateQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) TemplateQueryInterface

	HasSiteID() bool
	SiteID() string
	SetSiteID(siteID string) TemplateQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(includeSoftDeleted bool) TemplateQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) TemplateQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) TemplateQueryInterface
}
