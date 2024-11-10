package cmsstore

type PageQueryOptions struct {
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

type PageQueryInterface interface {
	ID() string
	SetID(id string) (PageQueryInterface, error)

	IDIn() []string
	SetIDIn(idIn []string) (PageQueryInterface, error)

	Handle() string
	SetHandle(handle string) (PageQueryInterface, error)

	AliasLike() string
	SetAliasLike(nameLike string) (PageQueryInterface, error)

	NameLike() string
	SetNameLike(nameLike string) (PageQueryInterface, error)

	Status() string
	SetStatus(status string) (PageQueryInterface, error)

	StatusIn() []string
	SetStatusIn(statusIn []string) (PageQueryInterface, error)

	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) (PageQueryInterface, error)

	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) (PageQueryInterface, error)

	Offset() int
	SetOffset(offset int) (PageQueryInterface, error)

	Limit() int
	SetLimit(limit int) (PageQueryInterface, error)

	SortOrder() string
	SetSortOrder(sortOrder string) (PageQueryInterface, error)

	OrderBy() string
	SetOrderBy(orderBy string) (PageQueryInterface, error)

	CountOnly() bool
	SetCountOnly(countOnly bool) PageQueryInterface

	WithSoftDeleted() bool
	SetWithSoftDeleted(withDeleted bool) PageQueryInterface
}
