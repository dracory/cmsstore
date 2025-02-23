package cmsstore

import "errors"

// site_query.go provides a query builder for site-related data.
// It includes methods to set and validate query parameters.

// == CONSTRUCTOR ============================================================
// SiteQuery returns a new instance of SiteQueryInterface.

// SiteQuery returns a new instance of SiteQueryInterface with an initialized parameters map.
func SiteQuery() SiteQueryInterface {
	return &siteQuery{
		parameters: make(map[string]any),
	}
}

// == TYPE =====================================================================
// siteQuery is the internal implementation of SiteQueryInterface.

type siteQuery struct {
	parameters map[string]any
}

// == INTERFACE VERIFICATION =================================================
// Verify that siteQuery implements SiteQueryInterface.

var _ SiteQueryInterface = (*siteQuery)(nil)

// == INTERFACE IMPLEMENTATION ===============================================
// Validate checks if the query parameters are valid.

// Validate checks if the query parameters are valid.
// It returns an error if any parameter is invalid.
func (q *siteQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("site query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("site query. created_at_lte cannot be empty")
	}

	if q.HasHandle() && q.Handle() == "" {
		return errors.New("site query. handle cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("site query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("site query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("site query. limit cannot be negative")
	}

	if q.HasNameLike() && q.NameLike() == "" {
		return errors.New("site query. name_like cannot be empty")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("site query. offset cannot be negative")
	}

	if q.HasOrderBy() && q.OrderBy() == "" {
		return errors.New("site query. order_by cannot be empty")
	}

	if q.HasSortOrder() && q.SortOrder() == "" {
		return errors.New("site query. sort_order cannot be empty")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("site query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("site query. status_in cannot be empty array")
	}

	// if q.HasUpdatedAtGte() && q.UpdatedAtGte() == "" {
	// 	return errors.New("site query. updated_at_gte cannot be empty")
	// }

	return nil
}

// Columns returns the list of columns set in the query.
func (q *siteQuery) Columns() []string {
	if !q.hasParameter(propertyKeyColumns) {
		return []string{}
	}

	return q.parameters[propertyKeyColumns].([]string)
}

// SetColumns sets the list of columns for the query.
func (q *siteQuery) SetColumns(columns []string) SiteQueryInterface {
	q.parameters[propertyKeyColumns] = columns
	return q
}

// HasCountOnly checks if the count_only parameter is set.
func (q *siteQuery) HasCountOnly() bool {
	return q.hasParameter(propertyKeyCountOnly)
}

// IsCountOnly returns the value of the count_only parameter.
func (q *siteQuery) IsCountOnly() bool {
	if !q.HasCountOnly() {
		return false
	}
	return q.parameters[propertyKeyCountOnly].(bool)
}

// SetCountOnly sets the count_only parameter.
func (q *siteQuery) SetCountOnly(isCountOnly bool) SiteQueryInterface {
	q.parameters[propertyKeyCountOnly] = isCountOnly
	return q
}

// HasCreatedAtGte checks if the created_at_gte parameter is set.
func (q *siteQuery) HasCreatedAtGte() bool {
	return q.hasParameter(propertyKeyCreatedAtGte)
}

// CreatedAtGte returns the value of the created_at_gte parameter.
func (q *siteQuery) CreatedAtGte() string {
	return q.parameters[propertyKeyCreatedAtGte].(string)
}

// SetCreatedAtGte sets the created_at_gte parameter.
func (q *siteQuery) SetCreatedAtGte(createdAtGte string) SiteQueryInterface {
	q.parameters[propertyKeyCreatedAtGte] = createdAtGte
	return q
}

// HasCreatedAtLte checks if the created_at_lte parameter is set.
func (q *siteQuery) HasCreatedAtLte() bool {
	return q.hasParameter(propertyKeyCreatedAtLte)
}

// CreatedAtLte returns the value of the created_at_lte parameter.
func (q *siteQuery) CreatedAtLte() string {
	return q.parameters[propertyKeyCreatedAtLte].(string)
}

// SetCreatedAtLte sets the created_at_lte parameter.
func (q *siteQuery) SetCreatedAtLte(createdAtLte string) SiteQueryInterface {
	q.parameters[propertyKeyCreatedAtLte] = createdAtLte
	return q
}

// HasDomainName checks if the domain_name parameter is set.
func (q *siteQuery) HasDomainName() bool {
	return q.hasParameter(propertyKeyDomainName)
}

// DomainName returns the value of the domain_name parameter.
func (q *siteQuery) DomainName() string {
	return q.parameters[propertyKeyDomainName].(string)
}

// SetDomainName sets the domain_name parameter.
func (q *siteQuery) SetDomainName(domainName string) SiteQueryInterface {
	q.parameters[propertyKeyDomainName] = domainName
	return q
}

// HasHandle checks if the handle parameter is set.
func (q *siteQuery) HasHandle() bool {
	return q.hasParameter(propertyKeyHandle)
}

// Handle returns the value of the handle parameter.
func (q *siteQuery) Handle() string {
	return q.parameters[propertyKeyHandle].(string)
}

// SetHandle sets the handle parameter.
func (q *siteQuery) SetHandle(handle string) SiteQueryInterface {
	q.parameters[propertyKeyHandle] = handle
	return q
}

// HasID checks if the id parameter is set.
func (q *siteQuery) HasID() bool {
	return q.hasParameter(propertyKeyId)
}

// ID returns the value of the id parameter.
func (q *siteQuery) ID() string {
	return q.parameters[propertyKeyId].(string)
}

// SetID sets the id parameter.
func (q *siteQuery) SetID(id string) SiteQueryInterface {
	q.parameters[propertyKeyId] = id
	return q
}

// HasIDIn checks if the id_in parameter is set.
func (q *siteQuery) HasIDIn() bool {
	return q.hasParameter(propertyKeyIdIn)
}

// IDIn returns the value of the id_in parameter.
func (q *siteQuery) IDIn() []string {
	return q.parameters[propertyKeyIdIn].([]string)
}

// SetIDIn sets the id_in parameter.
func (q *siteQuery) SetIDIn(idIn []string) SiteQueryInterface {
	q.parameters[propertyKeyIdIn] = idIn
	return q
}

// HasLimit checks if the limit parameter is set.
func (q *siteQuery) HasLimit() bool {
	return q.hasParameter(propertyKeyLimit)
}

// Limit returns the value of the limit parameter.
func (q *siteQuery) Limit() int {
	return q.parameters[propertyKeyLimit].(int)
}

// SetLimit sets the limit parameter.
func (q *siteQuery) SetLimit(limit int) SiteQueryInterface {
	q.parameters[propertyKeyLimit] = limit
	return q
}

// HasNameLike checks if the name_like parameter is set.
func (q *siteQuery) HasNameLike() bool {
	return q.hasParameter(propertyKeyNameLike)
}

// NameLike returns the value of the name_like parameter.
func (q *siteQuery) NameLike() string {
	return q.parameters[propertyKeyNameLike].(string)
}

// SetNameLike sets the name_like parameter.
func (q *siteQuery) SetNameLike(nameLike string) SiteQueryInterface {
	q.parameters[propertyKeyNameLike] = nameLike
	return q
}

// HasOffset checks if the offset parameter is set.
func (q *siteQuery) HasOffset() bool {
	return q.hasParameter(propertyKeyOffset)
}

// Offset returns the value of the offset parameter.
func (q *siteQuery) Offset() int {
	return q.parameters[propertyKeyOffset].(int)
}

// SetOffset sets the offset parameter.
func (q *siteQuery) SetOffset(offset int) SiteQueryInterface {
	q.parameters[propertyKeyOffset] = offset
	return q
}

// HasOrderBy checks if the order_by parameter is set.
func (q *siteQuery) HasOrderBy() bool {
	return q.hasParameter(propertyKeyOrderBy)
}

// OrderBy returns the value of the order_by parameter.
func (q *siteQuery) OrderBy() string {
	return q.parameters[propertyKeyOrderBy].(string)
}

// SetOrderBy sets the order_by parameter.
func (q *siteQuery) SetOrderBy(orderBy string) SiteQueryInterface {
	q.parameters[propertyKeyOrderBy] = orderBy
	return q
}

// HasSoftDeletedIncluded checks if the soft_deleted_included parameter is set.
func (q *siteQuery) HasSoftDeletedIncluded() bool {
	return q.hasParameter(propertyKeySoftDeleteIncluded)
}

// SoftDeletedIncluded returns the value of the soft_deleted_included parameter.
func (q *siteQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.parameters[propertyKeySoftDeleteIncluded].(bool)
}

// SetSoftDeletedIncluded sets the soft_deleted_included parameter.
func (q *siteQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) SiteQueryInterface {
	q.parameters[propertyKeySoftDeleteIncluded] = softDeletedIncluded
	return q
}

// HasSortOrder checks if the sort_order parameter is set.
func (q *siteQuery) HasSortOrder() bool {
	return q.hasParameter(propertyKeySortOrder)
}

// SortOrder returns the value of the sort_order parameter.
func (q *siteQuery) SortOrder() string {
	return q.parameters[propertyKeySortOrder].(string)
}

// SetSortOrder sets the sort_order parameter.
func (q *siteQuery) SetSortOrder(sortOrder string) SiteQueryInterface {
	q.parameters[propertyKeySortOrder] = sortOrder
	return q
}

// HasStatus checks if the status parameter is set.
func (q *siteQuery) HasStatus() bool {
	return q.hasParameter(propertyKeyStatus)
}

// Status returns the value of the status parameter.
func (q *siteQuery) Status() string {
	return q.parameters[propertyKeyStatus].(string)
}

// SetStatus sets the status parameter.
func (q *siteQuery) SetStatus(status string) SiteQueryInterface {
	q.parameters[propertyKeyStatus] = status
	return q
}

// HasStatusIn checks if the status_in parameter is set.
func (q *siteQuery) HasStatusIn() bool {
	return q.hasParameter(propertyKeyStatusIn)
}

// StatusIn returns the value of the status_in parameter.
func (q *siteQuery) StatusIn() []string {
	return q.parameters[propertyKeyStatusIn].([]string)
}

// SetStatusIn sets the status_in parameter.
func (q *siteQuery) SetStatusIn(statusIn []string) SiteQueryInterface {
	q.parameters[propertyKeyStatusIn] = statusIn
	return q
}

// == PRIVATE METHODS ========================================================

// hasParameter checks if a parameter is set.
func (q *siteQuery) hasParameter(key string) bool {
	return q.parameters[key] != nil
}
