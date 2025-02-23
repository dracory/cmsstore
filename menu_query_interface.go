package cmsstore

// MenuQueryInterface defines the methods required for querying menus.
type MenuQueryInterface interface {
	// Validate checks if the query parameters are valid.
	Validate() error

	// Columns returns the list of columns to be selected in the query.
	Columns() []string
	// SetColumns sets the list of columns to be selected in the query.
	SetColumns(columns []string) MenuQueryInterface

	// HasCountOnly checks if the query is set to return only the count.
	HasCountOnly() bool
	// IsCountOnly returns true if the query is set to return only the count.
	IsCountOnly() bool
	// SetCountOnly sets the query to return only the count.
	SetCountOnly(countOnly bool) MenuQueryInterface

	// HasCreatedAtGte checks if the query has a 'created_at' greater than or equal to condition.
	HasCreatedAtGte() bool
	// CreatedAtGte returns the 'created_at' greater than or equal to condition.
	CreatedAtGte() string
	// SetCreatedAtGte sets the 'created_at' greater than or equal to condition.
	SetCreatedAtGte(createdAtGte string) MenuQueryInterface

	// HasCreatedAtLte checks if the query has a 'created_at' less than or equal to condition.
	HasCreatedAtLte() bool
	// CreatedAtLte returns the 'created_at' less than or equal to condition.
	CreatedAtLte() string
	// SetCreatedAtLte sets the 'created_at' less than or equal to condition.
	SetCreatedAtLte(createdAtLte string) MenuQueryInterface

	// HasHandle checks if the query has a 'handle' condition.
	HasHandle() bool
	// Handle returns the 'handle' condition.
	Handle() string
	// SetHandle sets the 'handle' condition.
	SetHandle(handle string) MenuQueryInterface

	// HasID checks if the query has an 'id' condition.
	HasID() bool
	// ID returns the 'id' condition.
	ID() string
	// SetID sets the 'id' condition.
	SetID(id string) MenuQueryInterface

	// HasIDIn checks if the query has an 'id' in condition.
	HasIDIn() bool
	// IDIn returns the 'id' in condition.
	IDIn() []string
	// SetIDIn sets the 'id' in condition.
	SetIDIn(idIn []string) MenuQueryInterface

	// HasNameLike checks if the query has a 'name' like condition.
	HasNameLike() bool
	// NameLike returns the 'name' like condition.
	NameLike() string
	// SetNameLike sets the 'name' like condition.
	SetNameLike(nameLike string) MenuQueryInterface

	// HasOffset checks if the query has an offset condition.
	HasOffset() bool
	// Offset returns the offset condition.
	Offset() int
	// SetOffset sets the offset condition.
	SetOffset(offset int) MenuQueryInterface

	// HasLimit checks if the query has a limit condition.
	HasLimit() bool
	// Limit returns the limit condition.
	Limit() int
	// SetLimit sets the limit condition.
	SetLimit(limit int) MenuQueryInterface

	// HasSortOrder checks if the query has a sort order condition.
	HasSortOrder() bool
	// SortOrder returns the sort order condition.
	SortOrder() string
	// SetSortOrder sets the sort order condition.
	SetSortOrder(sortOrder string) MenuQueryInterface

	// HasOrderBy checks if the query has an order by condition.
	HasOrderBy() bool
	// OrderBy returns the order by condition.
	OrderBy() string
	// SetOrderBy sets the order by condition.
	SetOrderBy(orderBy string) MenuQueryInterface

	// HasSiteID checks if the query has a 'site_id' condition.
	HasSiteID() bool
	// SiteID returns the 'site_id' condition.
	SiteID() string
	// SetSiteID sets the 'site_id' condition.
	SetSiteID(siteID string) MenuQueryInterface

	// HasSoftDeletedIncluded checks if the query includes soft deleted records.
	HasSoftDeletedIncluded() bool
	// SoftDeletedIncluded returns true if the query includes soft deleted records.
	SoftDeletedIncluded() bool
	// SetSoftDeletedIncluded sets whether the query should include soft deleted records.
	SetSoftDeletedIncluded(includeSoftDeleted bool) MenuQueryInterface

	// HasStatus checks if the query has a 'status' condition.
	HasStatus() bool
	// Status returns the 'status' condition.
	Status() string
	// SetStatus sets the 'status' condition.
	SetStatus(status string) MenuQueryInterface

	// HasStatusIn checks if the query has a 'status' in condition.
	HasStatusIn() bool
	// StatusIn returns the 'status' in condition.
	StatusIn() []string
	// SetStatusIn sets the 'status' in condition.
	SetStatusIn(statusIn []string) MenuQueryInterface
}
