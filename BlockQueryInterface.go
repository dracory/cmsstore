package cmsstore

// type BlockQueryOptions struct {
// 	ID           string
// 	IDIn         []string
// 	NameLike     string
// 	Status       string
// 	StatusIn     []string
// 	Handle       string
// 	CreatedAtGte string
// 	CreatedAtLte string
// 	Offset       int
// 	Limit        int
// 	SortOrder    string
// 	OrderBy      string
// 	CountOnly    bool
// 	WithDeleted  bool
// }

type BlockQueryInterface interface {
	Validate() error

	IsCountOnly() bool

	HasCreatedAtGte() bool
	HasCreatedAtLte() bool
	HasHandle() bool
	HasID() bool
	HasIDIn() bool
	HasLimit() bool
	HasNameLike() bool
	HasOffset() bool
	HasOrderBy() bool
	HasPageID() bool
	HasParentID() bool
	HasSequence() bool
	HasSiteID() bool
	HasSoftDeleted() bool
	HasSortOrder() bool
	HasStatus() bool
	HasStatusIn() bool
	HasTemplateID() bool

	CreatedAtGte() string
	CreatedAtLte() string
	Handle() string
	ID() string
	IDIn() []string
	Limit() int
	NameLike() string
	Offset() int
	OrderBy() string
	PageID() string
	ParentID() string
	Sequence() int
	SiteID() string
	SoftDeleteIncluded() bool
	SortOrder() string
	Status() string
	StatusIn() []string
	TemplateID() string

	SetCountOnly(countOnly bool) BlockQueryInterface
	SetID(id string) BlockQueryInterface
	SetIDIn(idIn []string) BlockQueryInterface
	SetHandle(handle string) BlockQueryInterface
	SetLimit(limit int) BlockQueryInterface
	SetNameLike(nameLike string) BlockQueryInterface
	SetOffset(offset int) BlockQueryInterface
	SetOrderBy(orderBy string) BlockQueryInterface
	SetPageID(pageID string) BlockQueryInterface
	SetParentID(parentID string) BlockQueryInterface
	SetSequence(sequence int) BlockQueryInterface
	SetSiteID(websiteID string) BlockQueryInterface
	SetSoftDeleteIncluded(withSoftDeleted bool) BlockQueryInterface
	SetSortOrder(sortOrder string) BlockQueryInterface
	SetStatus(status string) BlockQueryInterface
	SetStatusIn(statusIn []string) BlockQueryInterface
	SetTemplateID(templateID string) BlockQueryInterface

	// ID() string
	// SetID(id string) (BlockQueryInterface, error)

	// IDIn() []string
	// SetIDIn(idIn []string) (BlockQueryInterface, error)

	// NameLike() string
	// SetNameLike(nameLike string) (BlockQueryInterface, error)

	// Status() string
	// SetStatus(status string) (BlockQueryInterface, error)

	// StatusIn() []string
	// SetStatusIn(statusIn []string) (BlockQueryInterface, error)

	// Handle() string
	// SetHandle(handle string) (BlockQueryInterface, error)

	// CreatedAtGte() string
	// SetCreatedAtGte(createdAtGte string) (BlockQueryInterface, error)

	// CreatedAtLte() string
	// SetCreatedAtLte(createdAtLte string) (BlockQueryInterface, error)

	// Offset() int
	// SetOffset(offset int) (BlockQueryInterface, error)

	// Limit() int
	// SetLimit(limit int) (BlockQueryInterface, error)

	// SortOrder() string
	// SetSortOrder(sortOrder string) (BlockQueryInterface, error)

	// OrderBy() string
	// SetOrderBy(orderBy string) (BlockQueryInterface, error)

	// CountOnly() bool
	// SetCountOnly(countOnly bool) BlockQueryInterface

	// WithSoftDeleted() bool
	// SetWithSoftDeleted(withDeleted bool) BlockQueryInterface
}
