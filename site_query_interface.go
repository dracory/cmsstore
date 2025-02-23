package cmsstore

// SiteQueryInterface defines the methods required for querying site data.
type SiteQueryInterface interface {
	// Validate checks if the query parameters are valid.
	Validate() error

	// Columns returns the list of columns to be selected in the query.
	Columns() []string
	// SetColumns sets the list of columns to be selected in the query.
	SetColumns(columns []string) SiteQueryInterface

	// HasCountOnly checks if the query is set to return only the count of records.
	HasCountOnly() bool
	// IsCountOnly returns true if the query is set to return only the count of records.
	IsCountOnly() bool
	// SetCountOnly sets the query to return only the count of records.
	SetCountOnly(countOnly bool) SiteQueryInterface

	// HasCreatedAtGte checks if the query has a 'created_at' greater than or equal to condition.
	HasCreatedAtGte() bool
	// CreatedAtGte returns the 'created_at' greater than or equal to condition.
	CreatedAtGte() string
	// SetCreatedAtGte sets the 'created_at' greater than or equal to condition.
	SetCreatedAtGte(createdAtGte string) SiteQueryInterface

	// HasCreatedAtLte checks if the query has a 'created_at' less than or equal to condition.
	HasCreatedAtLte() bool
	// CreatedAtLte returns the 'created_at' less than or equal to condition.
	CreatedAtLte() string
	// SetCreatedAtLte sets the 'created_at' less than or equal to condition.
	SetCreatedAtLte(createdAtLte string) SiteQueryInterface

	// HasDomainName checks if the query has a domain name condition.
	HasDomainName() bool
	// DomainName returns the domain name condition.
	DomainName() string
	// SetDomainName sets the domain name condition.
	SetDomainName(domainName string) SiteQueryInterface

	// HasHandle checks if the query has a handle condition.
	HasHandle() bool
	// Handle returns the handle condition.
	Handle() string
	// SetHandle sets the handle condition.
	SetHandle(handle string) SiteQueryInterface

	// HasID checks if the query has an ID condition.
	HasID() bool
	// ID returns the ID condition.
	ID() string
	// SetID sets the ID condition.
	SetID(id string) SiteQueryInterface

	// HasIDIn checks if the query has an ID in condition.
	HasIDIn() bool
	// IDIn returns the ID in condition.
	IDIn() []string
	// SetIDIn sets the ID in condition.
	SetIDIn(idIn []string) SiteQueryInterface

	// HasLimit checks if the query has a limit condition.
	HasLimit() bool
	// Limit returns the limit condition.
	Limit() int
	// SetLimit sets the limit condition.
	SetLimit(limit int) SiteQueryInterface

	// HasNameLike checks if the query has a name like condition.
	HasNameLike() bool
	// NameLike returns the name like condition.
	NameLike() string
	// SetNameLike sets the name like condition.
	SetNameLike(nameLike string) SiteQueryInterface

	// HasOffset checks if the query has an offset condition.
	HasOffset() bool
	// Offset returns the offset condition.
	Offset() int
	// SetOffset sets the offset condition.
	SetOffset(offset int) SiteQueryInterface

	// HasSortOrder checks if the query has a sort order condition.
	HasSortOrder() bool
	// SortOrder returns the sort order condition.
	SortOrder() string
	// SetSortOrder sets the sort order condition.
	SetSortOrder(sortOrder string) SiteQueryInterface

	// HasOrderBy checks if the query has an order by condition.
	HasOrderBy() bool
	// OrderBy returns the order by condition.
	OrderBy() string
	// SetOrderBy sets the order by condition.
	SetOrderBy(orderBy string) SiteQueryInterface

	// HasSoftDeletedIncluded checks if the query includes soft deleted records.
	HasSoftDeletedIncluded() bool
	// SoftDeletedIncluded returns true if the query includes soft deleted records.
	SoftDeletedIncluded() bool
	// SetSoftDeletedIncluded sets the query to include soft deleted records.
	SetSoftDeletedIncluded(softDeletedIncluded bool) SiteQueryInterface

	// HasStatus checks if the query has a status condition.
	HasStatus() bool
	// Status returns the status condition.
	Status() string
	// SetStatus sets the status condition.
	SetStatus(status string) SiteQueryInterface

	// HasStatusIn checks if the query has a status in condition.
	HasStatusIn() bool
	// StatusIn returns the status in condition.
	StatusIn() []string
	// SetStatusIn sets the status in condition.
	SetStatusIn(statusIn []string) SiteQueryInterface
}
