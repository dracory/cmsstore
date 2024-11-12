package cmsstore

type TranslationQueryOptions struct {
	ID           string
	IDIn         []string
	Status       string
	StatusIn     []string
	Handle       string
	CreatedAtGte string
	CreatedAtLte string
	Offset       int
	Limit        int
	SortOrder    string
	OrderBy      string
	CountOnly    bool
	WithDeleted  bool
}

type TranslationQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) TranslationQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) TranslationQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) TranslationQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) TranslationQueryInterface

	HasHandle() bool
	Handle() string
	SetHandle(handle string) TranslationQueryInterface

	HasID() bool
	ID() string
	SetID(id string) TranslationQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) TranslationQueryInterface

	HasNameLike() bool
	NameLike() string
	SetNameLike(nameLike string) TranslationQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) TranslationQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) TranslationQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) TranslationQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) TranslationQueryInterface

	HasSiteID() bool
	SiteID() string
	SetSiteID(siteID string) TranslationQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(includeSoftDeleted bool) TranslationQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) TranslationQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) TranslationQueryInterface
}
