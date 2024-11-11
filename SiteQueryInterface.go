package cmsstore

type SiteQueryOptions struct {
	ID           string
	IDIn         []string
	Handle       string
	NameLike     string
	Status       string
	StatusIn     []string
	CreatedAtGte string
	CreatedAtLte string
	Offset       int
	Limit        int
	SortOrder    string
	OrderBy      string
	CountOnly    bool
	WithDeleted  bool
}

type SiteQueryInterface interface {
	Validate() error

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) SiteQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) SiteQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) SiteQueryInterface

	HasDomainName() bool
	DomainName() string
	SetDomainName(domainName string) SiteQueryInterface

	HasHandle() bool
	Handle() string
	SetHandle(handle string) SiteQueryInterface

	HasID() bool
	ID() string
	SetID(id string) SiteQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) SiteQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) SiteQueryInterface

	HasNameLike() bool
	NameLike() string
	SetNameLike(nameLike string) SiteQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) SiteQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) SiteQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) SiteQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) SiteQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) SiteQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) SiteQueryInterface
}
