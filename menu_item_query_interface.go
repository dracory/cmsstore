// Package cmsstore defines the interface for querying menu items.
package cmsstore

// MenuItemQueryInterface defines the interface for querying menu items.
type MenuItemQueryInterface interface {
	// Validate validates the query parameters.
	Validate() error

	// Columns returns the columns to be returned in the query.
	Columns() []string
	// SetColumns sets the columns to be returned in the query.
	SetColumns(columns []string) MenuItemQueryInterface

	// HasCountOnly returns true if the query is for counting only.
	HasCountOnly() bool
	// IsCountOnly returns true if the query is for counting only.
	IsCountOnly() bool
	// SetCountOnly sets the query to be for counting only.
	SetCountOnly(countOnly bool) MenuItemQueryInterface

	// HasCreatedAtGte returns true if the query has a CreatedAtGte filter.
	HasCreatedAtGte() bool
	// CreatedAtGte returns the CreatedAtGte filter.
	CreatedAtGte() string
	// SetCreatedAtGte sets the CreatedAtGte filter.
	SetCreatedAtGte(createdAtGte string) MenuItemQueryInterface

	// HasCreatedAtLte returns true if the query has a CreatedAtLte filter.
	HasCreatedAtLte() bool
	// CreatedAtLte returns the CreatedAtLte filter.
	CreatedAtLte() string
	// SetCreatedAtLte sets the CreatedAtLte filter.
	SetCreatedAtLte(createdAtLte string) MenuItemQueryInterface

	// HasID returns true if the query has an ID filter.
	HasID() bool
	// ID returns the ID filter.
	ID() string
	// SetID sets the ID filter.
	SetID(id string) MenuItemQueryInterface

	// HasIDIn returns true if the query has an IDIn filter.
	HasIDIn() bool
	// IDIn returns the IDIn filter.
	IDIn() []string
	// SetIDIn sets the IDIn filter.
	SetIDIn(idIn []string) MenuItemQueryInterface

	// HasMenuID returns true if the query has a MenuID filter.
	HasMenuID() bool
	// MenuID returns the MenuID filter.
	MenuID() string
	// SetMenuID sets the MenuID filter.
	SetMenuID(menuID string) MenuItemQueryInterface

	// HasNameLike returns true if the query has a NameLike filter.
	HasNameLike() bool
	// NameLike returns the NameLike filter.
	NameLike() string
	// SetNameLike sets the NameLike filter.
	SetNameLike(nameLike string) MenuItemQueryInterface

	// HasOffset returns true if the query has an Offset.
	HasOffset() bool
	// Offset returns the Offset.
	Offset() int
	// SetOffset sets the Offset.
	SetOffset(offset int) MenuItemQueryInterface

	// HasLimit returns true if the query has a Limit.
	HasLimit() bool
	// Limit returns the Limit.
	Limit() int
	// SetLimit sets the Limit.
	SetLimit(limit int) MenuItemQueryInterface

	// HasSortOrder returns true if the query has a SortOrder.
	HasSortOrder() bool
	// SortOrder returns the SortOrder.
	SortOrder() string
	// SetSortOrder sets the SortOrder.
	SetSortOrder(sortOrder string) MenuItemQueryInterface

	// HasOrderBy returns true if the query has an OrderBy.
	HasOrderBy() bool
	// OrderBy returns the OrderBy.
	OrderBy() string
	// SetOrderBy sets the OrderBy.
	SetOrderBy(orderBy string) MenuItemQueryInterface

	// HasSiteID returns true if the query has a SiteID filter.
	HasSiteID() bool
	// SiteID returns the SiteID filter.
	SiteID() string
	// SetSiteID sets the SiteID filter.
	SetSiteID(siteID string) MenuItemQueryInterface

	// HasSoftDeletedIncluded returns true if soft deleted items should be included.
	HasSoftDeletedIncluded() bool
	// SoftDeletedIncluded returns true if soft deleted items should be included.
	SoftDeletedIncluded() bool
	// SetSoftDeletedIncluded sets whether soft deleted items should be included.
	SetSoftDeletedIncluded(includeSoftDeleted bool) MenuItemQueryInterface

	// HasStatus returns true if the query has a Status filter.
	HasStatus() bool
	// Status returns the Status filter.
	Status() string
	// SetStatus sets the Status filter.
	SetStatus(status string) MenuItemQueryInterface

	// HasStatusIn returns true if the query has a StatusIn filter.
	HasStatusIn() bool
	// StatusIn returns the StatusIn filter.
	StatusIn() []string
	// SetStatusIn sets the StatusIn filter.
	SetStatusIn(statusIn []string) MenuItemQueryInterface
}
