package cmsstore

import "errors"

// == CONSTRUCTOR ============================================================

func SiteQuery() SiteQueryInterface {
	return &siteQuery{
		parameters: make(map[string]any),
	}
}

// ==TYPE =====================================================================

type siteQuery struct {
	parameters map[string]any
}

// == INTERFACE VERIFICATION =================================================

var _ SiteQueryInterface = (*siteQuery)(nil)

// == INTERFACE IMPLEMENTATION ===============================================

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

func (q *siteQuery) Columns() []string {
	if !q.hasParameter(propertyKeyColumns) {
		return []string{}
	}

	return q.parameters[propertyKeyColumns].([]string)
}

func (q *siteQuery) SetColumns(columns []string) SiteQueryInterface {
	q.parameters[propertyKeyColumns] = columns
	return q
}

func (q *siteQuery) HasCountOnly() bool {
	return q.hasParameter(propertyKeyCountOnly)
}

func (q *siteQuery) IsCountOnly() bool {
	if !q.HasCountOnly() {
		return false
	}
	return q.parameters[propertyKeyCountOnly].(bool)
}

func (q *siteQuery) SetCountOnly(isCountOnly bool) SiteQueryInterface {
	q.parameters[propertyKeyCountOnly] = isCountOnly
	return q
}

func (q *siteQuery) HasCreatedAtGte() bool {
	return q.hasParameter(propertyKeyCreatedAtGte)
}

func (q *siteQuery) CreatedAtGte() string {
	return q.parameters[propertyKeyCreatedAtGte].(string)
}

func (q *siteQuery) SetCreatedAtGte(createdAtGte string) SiteQueryInterface {
	q.parameters[propertyKeyCreatedAtGte] = createdAtGte
	return q
}

func (q *siteQuery) HasCreatedAtLte() bool {
	return q.hasParameter(propertyKeyCreatedAtLte)
}

func (q *siteQuery) CreatedAtLte() string {
	return q.parameters[propertyKeyCreatedAtLte].(string)
}

func (q *siteQuery) SetCreatedAtLte(createdAtLte string) SiteQueryInterface {
	q.parameters[propertyKeyCreatedAtLte] = createdAtLte
	return q
}

func (q *siteQuery) HasDomainName() bool {
	return q.hasParameter(propertyKeyDomainName)
}

func (q *siteQuery) DomainName() string {
	return q.parameters[propertyKeyDomainName].(string)
}

func (q *siteQuery) SetDomainName(domainName string) SiteQueryInterface {
	q.parameters[propertyKeyDomainName] = domainName
	return q
}

func (q *siteQuery) HasHandle() bool {
	return q.hasParameter(propertyKeyHandle)
}

func (q *siteQuery) Handle() string {
	return q.parameters[propertyKeyHandle].(string)
}

func (q *siteQuery) SetHandle(handle string) SiteQueryInterface {
	q.parameters[propertyKeyHandle] = handle
	return q
}

func (q *siteQuery) HasID() bool {
	return q.hasParameter(propertyKeyId)
}

func (q *siteQuery) ID() string {
	return q.parameters[propertyKeyId].(string)
}

func (q *siteQuery) SetID(id string) SiteQueryInterface {
	q.parameters[propertyKeyId] = id
	return q
}

func (q *siteQuery) HasIDIn() bool {
	return q.hasParameter(propertyKeyIdIn)
}

func (q *siteQuery) IDIn() []string {
	return q.parameters[propertyKeyIdIn].([]string)
}

func (q *siteQuery) SetIDIn(idIn []string) SiteQueryInterface {
	q.parameters[propertyKeyIdIn] = idIn
	return q
}

func (q *siteQuery) HasLimit() bool {
	return q.hasParameter(propertyKeyLimit)
}

func (q *siteQuery) Limit() int {
	return q.parameters[propertyKeyLimit].(int)
}

func (q *siteQuery) SetLimit(limit int) SiteQueryInterface {
	q.parameters[propertyKeyLimit] = limit
	return q
}

func (q *siteQuery) HasNameLike() bool {
	return q.hasParameter(propertyKeyNameLike)
}

func (q *siteQuery) NameLike() string {
	return q.parameters[propertyKeyNameLike].(string)
}

func (q *siteQuery) SetNameLike(nameLike string) SiteQueryInterface {
	q.parameters[propertyKeyNameLike] = nameLike
	return q
}

func (q *siteQuery) HasOffset() bool {
	return q.hasParameter(propertyKeyOffset)
}

func (q *siteQuery) Offset() int {
	return q.parameters[propertyKeyOffset].(int)
}

func (q *siteQuery) SetOffset(offset int) SiteQueryInterface {
	q.parameters[propertyKeyOffset] = offset
	return q
}

func (q *siteQuery) HasOrderBy() bool {
	return q.hasParameter(propertyKeyOrderBy)
}

func (q *siteQuery) OrderBy() string {
	return q.parameters[propertyKeyOrderBy].(string)
}

func (q *siteQuery) SetOrderBy(orderBy string) SiteQueryInterface {
	q.parameters[propertyKeyOrderBy] = orderBy
	return q
}

func (q *siteQuery) HasSoftDeletedIncluded() bool {
	return q.hasParameter(propertyKeySoftDeleteIncluded)
}

func (q *siteQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.parameters[propertyKeySoftDeleteIncluded].(bool)
}

func (q *siteQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) SiteQueryInterface {
	q.parameters[propertyKeySoftDeleteIncluded] = softDeletedIncluded
	return q
}

func (q *siteQuery) HasSortOrder() bool {
	return q.hasParameter(propertyKeySortOrder)
}

func (q *siteQuery) SortOrder() string {
	return q.parameters[propertyKeySortOrder].(string)
}

func (q *siteQuery) SetSortOrder(sortOrder string) SiteQueryInterface {
	q.parameters[propertyKeySortOrder] = sortOrder
	return q
}

func (q *siteQuery) HasStatus() bool {
	return q.hasParameter(propertyKeyStatus)
}

func (q *siteQuery) Status() string {
	return q.parameters[propertyKeyStatus].(string)
}

func (q *siteQuery) SetStatus(status string) SiteQueryInterface {
	q.parameters[propertyKeyStatus] = status
	return q
}

func (q *siteQuery) HasStatusIn() bool {
	return q.hasParameter(propertyKeyStatusIn)
}

func (q *siteQuery) StatusIn() []string {
	return q.parameters[propertyKeyStatusIn].([]string)
}

func (q *siteQuery) SetStatusIn(statusIn []string) SiteQueryInterface {
	q.parameters[propertyKeyStatusIn] = statusIn
	return q
}

// == PRIVATE METHODS ========================================================

func (q *siteQuery) hasParameter(key string) bool {
	return q.parameters[key] != nil
}
