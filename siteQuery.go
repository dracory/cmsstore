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

func (q *siteQuery) HasCountOnly() bool {
	return q.hasParameter("count_only")
}

func (q *siteQuery) IsCountOnly() bool {
	if !q.HasCountOnly() {
		return false
	}
	return q.parameters["count_only"].(bool)
}

func (q *siteQuery) SetCountOnly(isCountOnly bool) SiteQueryInterface {
	q.parameters["count_only"] = isCountOnly
	return q
}

func (q *siteQuery) HasCreatedAtGte() bool {
	return q.hasParameter("created_at_gte")
}

func (q *siteQuery) CreatedAtGte() string {
	return q.parameters["created_at_gte"].(string)
}

func (q *siteQuery) SetCreatedAtGte(createdAtGte string) SiteQueryInterface {
	q.parameters["created_at_gte"] = createdAtGte
	return q
}

func (q *siteQuery) HasCreatedAtLte() bool {
	return q.hasParameter("created_at_lte")
}

func (q *siteQuery) CreatedAtLte() string {
	return q.parameters["created_at_lte"].(string)
}

func (q *siteQuery) SetCreatedAtLte(createdAtLte string) SiteQueryInterface {
	q.parameters["created_at_lte"] = createdAtLte
	return q
}

func (q *siteQuery) HasDomainName() bool {
	return q.hasParameter("domain_name")
}

func (q *siteQuery) DomainName() string {
	return q.parameters["domain_name"].(string)
}

func (q *siteQuery) SetDomainName(domainName string) SiteQueryInterface {
	q.parameters["domain_name"] = domainName
	return q
}

func (q *siteQuery) HasHandle() bool {
	return q.hasParameter("handle")
}

func (q *siteQuery) Handle() string {
	return q.parameters["handle"].(string)
}

func (q *siteQuery) SetHandle(handle string) SiteQueryInterface {
	q.parameters["handle"] = handle
	return q
}

func (q *siteQuery) HasID() bool {
	return q.hasParameter("id")
}

func (q *siteQuery) ID() string {
	return q.parameters["id"].(string)
}

func (q *siteQuery) SetID(id string) SiteQueryInterface {
	q.parameters["id"] = id
	return q
}

func (q *siteQuery) HasIDIn() bool {
	return q.hasParameter("id_in")
}

func (q *siteQuery) IDIn() []string {
	return q.parameters["id_in"].([]string)
}

func (q *siteQuery) SetIDIn(idIn []string) SiteQueryInterface {
	q.parameters["id_in"] = idIn
	return q
}

func (q *siteQuery) HasLimit() bool {
	return q.hasParameter("limit")
}

func (q *siteQuery) Limit() int {
	return q.parameters["limit"].(int)
}

func (q *siteQuery) SetLimit(limit int) SiteQueryInterface {
	q.parameters["limit"] = limit
	return q
}

func (q *siteQuery) HasNameLike() bool {
	return q.hasParameter("name_like")
}

func (q *siteQuery) NameLike() string {
	return q.parameters["name_like"].(string)
}

func (q *siteQuery) SetNameLike(nameLike string) SiteQueryInterface {
	q.parameters["name_like"] = nameLike
	return q
}

func (q *siteQuery) HasOffset() bool {
	return q.hasParameter("offset")
}

func (q *siteQuery) Offset() int {
	return q.parameters["offset"].(int)
}

func (q *siteQuery) SetOffset(offset int) SiteQueryInterface {
	q.parameters["offset"] = offset
	return q
}

func (q *siteQuery) HasOrderBy() bool {
	return q.hasParameter("order_by")
}

func (q *siteQuery) OrderBy() string {
	return q.parameters["order_by"].(string)
}

func (q *siteQuery) SetOrderBy(orderBy string) SiteQueryInterface {
	q.parameters["order_by"] = orderBy
	return q
}

func (q *siteQuery) HasSoftDeletedIncluded() bool {
	return q.hasParameter("soft_deleted_included")
}

func (q *siteQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.parameters["soft_deleted_included"].(bool)
}

func (q *siteQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) SiteQueryInterface {
	q.parameters["soft_deleted_included"] = softDeletedIncluded
	return q
}

func (q *siteQuery) HasSortOrder() bool {
	return q.hasParameter("sort_order")
}

func (q *siteQuery) SortOrder() string {
	return q.parameters["sort_order"].(string)
}

func (q *siteQuery) SetSortOrder(sortOrder string) SiteQueryInterface {
	q.parameters["sort_order"] = sortOrder
	return q
}

func (q *siteQuery) HasStatus() bool {
	return q.hasParameter("status")
}

func (q *siteQuery) Status() string {
	return q.parameters["status"].(string)
}

func (q *siteQuery) SetStatus(status string) SiteQueryInterface {
	q.parameters["status"] = status
	return q
}

func (q *siteQuery) HasStatusIn() bool {
	return q.hasParameter("status_in")
}

func (q *siteQuery) StatusIn() []string {
	return q.parameters["status_in"].([]string)
}

func (q *siteQuery) SetStatusIn(statusIn []string) SiteQueryInterface {
	q.parameters["status_in"] = statusIn
	return q
}

// == PRIVATE METHODS ========================================================

func (q *siteQuery) hasParameter(key string) bool {
	return q.parameters[key] != nil
}
