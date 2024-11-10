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
	ID() string
	SetID(id string) (SiteQueryInterface, error)

	IDIn() []string
	SetIDIn(idIn []string) (SiteQueryInterface, error)

	NameLike() string
	SetNameLike(nameLike string) (SiteQueryInterface, error)

	Status() string
	SetStatus(status string) (SiteQueryInterface, error)

	StatusIn() []string
	SetStatusIn(statusIn []string) (SiteQueryInterface, error)

	Handle() string
	SetHandle(handle string) (SiteQueryInterface, error)

	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) (SiteQueryInterface, error)

	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) (SiteQueryInterface, error)

	Offset() int
	SetOffset(offset int) (SiteQueryInterface, error)

	Limit() int
	SetLimit(limit int) (SiteQueryInterface, error)

	SortOrder() string
	SetSortOrder(sortOrder string) (SiteQueryInterface, error)

	OrderBy() string
	SetOrderBy(orderBy string) (SiteQueryInterface, error)

	CountOnly() bool
	SetCountOnly(countOnly bool) SiteQueryInterface

	WithSoftDeleted() bool
	SetWithSoftDeleted(withDeleted bool) SiteQueryInterface
}
