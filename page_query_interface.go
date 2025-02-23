package cmsstore

// PageQueryInterface defines methods for querying pages in the CMS store.
type PageQueryInterface interface {
	// Validate checks if the query parameters are valid.
	Validate() error

	// Columns returns the list of columns to be queried.
	Columns() []string
	// SetColumns sets the list of columns to be queried.
	SetColumns(columns []string) PageQueryInterface

	// HasAlias checks if an alias is set.
	HasAlias() bool
	// Alias returns the alias if set.
	Alias() string
	// SetAlias sets the alias.
	SetAlias(alias string) PageQueryInterface

	// HasAliasLike checks if an alias pattern is set.
	HasAliasLike() bool
	// AliasLike returns the alias pattern if set.
	AliasLike() string
	// SetAliasLike sets the alias pattern.
	SetAliasLike(nameLike string) PageQueryInterface

	// HasCreatedAtGte checks if a 'created at' greater-than-or-equal-to filter is set.
	HasCreatedAtGte() bool
	// CreatedAtGte returns the 'created at' greater-than-or-equal-to filter if set.
	CreatedAtGte() string
	// SetCreatedAtGte sets the 'created at' greater-than-or-equal-to filter.
	SetCreatedAtGte(createdAtGte string) PageQueryInterface

	// HasCreatedAtLte checks if a 'created at' less-than-or-equal-to filter is set.
	HasCreatedAtLte() bool
	// CreatedAtLte returns the 'created at' less-than-or-equal-to filter if set.
	CreatedAtLte() string
	// SetCreatedAtLte sets the 'created at' less-than-or-equal-to filter.
	SetCreatedAtLte(createdAtLte string) PageQueryInterface

	// HasCountOnly checks if count-only mode is set.
	HasCountOnly() bool
	// IsCountOnly returns the count-only mode setting.
	IsCountOnly() bool
	// SetCountOnly sets the count-only mode.
	SetCountOnly(countOnly bool) PageQueryInterface

	// HasHandle checks if a handle is set.
	HasHandle() bool
	// Handle returns the handle if set.
	Handle() string
	// SetHandle sets the handle.
	SetHandle(handle string) PageQueryInterface

	// HasID checks if an ID is set.
	HasID() bool
	// ID returns the ID if set.
	ID() string
	// SetID sets the ID.
	SetID(id string) PageQueryInterface

	// HasIDIn checks if an ID list is set.
	HasIDIn() bool
	// IDIn returns the ID list if set.
	IDIn() []string
	// SetIDIn sets the ID list.
	SetIDIn(idIn []string) PageQueryInterface

	// HasLimit checks if a limit is set.
	HasLimit() bool
	// Limit returns the limit if set.
	Limit() int
	// SetLimit sets the limit.
	SetLimit(limit int) PageQueryInterface

	// HasNameLike checks if a name pattern is set.
	HasNameLike() bool
	// NameLike returns the name pattern if set.
	NameLike() string
	// SetNameLike sets the name pattern.
	SetNameLike(nameLike string) PageQueryInterface

	// HasOffset checks if an offset is set.
	HasOffset() bool
	// Offset returns the offset if set.
	Offset() int
	// SetOffset sets the offset.
	SetOffset(offset int) PageQueryInterface

	// HasOrderBy checks if an order-by clause is set.
	HasOrderBy() bool
	// OrderBy returns the order-by clause if set.
	OrderBy() string
	// SetOrderBy sets the order-by clause.
	SetOrderBy(orderBy string) PageQueryInterface

	// HasSiteID checks if a site ID is set.
	HasSiteID() bool
	// SiteID returns the site ID if set.
	SiteID() string
	// SetSiteID sets the site ID.
	SetSiteID(siteID string) PageQueryInterface

	// HasSortOrder checks if a sort order is set.
	HasSortOrder() bool
	// SortOrder returns the sort order if set.
	SortOrder() string
	// SetSortOrder sets the sort order.
	SetSortOrder(sortOrder string) PageQueryInterface

	// HasSoftDeletedIncluded checks if soft-deleted records are included.
	HasSoftDeletedIncluded() bool
	// SoftDeletedIncluded returns whether soft-deleted records are included.
	SoftDeletedIncluded() bool
	// SetSoftDeletedIncluded sets whether to include soft-deleted records.
	SetSoftDeletedIncluded(softDeleteIncluded bool) PageQueryInterface

	// HasStatus checks if a status is set.
	HasStatus() bool
	// Status returns the status if set.
	Status() string
	// SetStatus sets the status.
	SetStatus(status string) PageQueryInterface

	// HasStatusIn checks if a status list is set.
	HasStatusIn() bool
	// StatusIn returns the status list if set.
	StatusIn() []string
	// SetStatusIn sets the status list.
	SetStatusIn(statusIn []string) PageQueryInterface

	// HasTemplateID checks if a template ID is set.
	HasTemplateID() bool
	// TemplateID returns the template ID if set.
	TemplateID() string
	// SetTemplateID sets the template ID.
	SetTemplateID(templateID string) PageQueryInterface
}
