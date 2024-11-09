package cmsstore

type TemplateQueryOptions struct {
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

type TemplateQueryInterface interface {
	ID() string
	SetID(id string) (TemplateQueryInterface, error)

	IDIn() []string
	SetIDIn(idIn []string) (TemplateQueryInterface, error)

	Status() string
	SetStatus(status string) (TemplateQueryInterface, error)

	StatusIn() []string
	SetStatusIn(statusIn []string) (TemplateQueryInterface, error)

	Handle() string
	SetHandle(handle string) (TemplateQueryInterface, error)

	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) (TemplateQueryInterface, error)

	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) (TemplateQueryInterface, error)

	Offset() int
	SetOffset(offset int) (TemplateQueryInterface, error)

	Limit() int
	SetLimit(limit int) (TemplateQueryInterface, error)

	SortOrder() string
	SetSortOrder(sortOrder string) (TemplateQueryInterface, error)

	OrderBy() string
	SetOrderBy(orderBy string) (TemplateQueryInterface, error)

	CountOnly() bool
	SetCountOnly(countOnly bool) TemplateQueryInterface

	WithSoftDeleted() bool
	SetWithSoftDeleted(withDeleted bool) TemplateQueryInterface
}
