package cmsstore

type BlockQueryOptions struct {
	ID           string
	IDIn         []string
	NameLike     string
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

type BlockQueryInterface interface {
	ID() string
	SetID(id string) (BlockQueryInterface, error)

	IDIn() []string
	SetIDIn(idIn []string) (BlockQueryInterface, error)

	NameLike() string
	SetNameLike(nameLike string) (BlockQueryInterface, error)

	Status() string
	SetStatus(status string) (BlockQueryInterface, error)

	StatusIn() []string
	SetStatusIn(statusIn []string) (BlockQueryInterface, error)

	Handle() string
	SetHandle(handle string) (BlockQueryInterface, error)

	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) (BlockQueryInterface, error)

	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) (BlockQueryInterface, error)

	Offset() int
	SetOffset(offset int) (BlockQueryInterface, error)

	Limit() int
	SetLimit(limit int) (BlockQueryInterface, error)

	SortOrder() string
	SetSortOrder(sortOrder string) (BlockQueryInterface, error)

	OrderBy() string
	SetOrderBy(orderBy string) (BlockQueryInterface, error)

	CountOnly() bool
	SetCountOnly(countOnly bool) BlockQueryInterface

	WithSoftDeleted() bool
	SetWithSoftDeleted(withDeleted bool) BlockQueryInterface
}
