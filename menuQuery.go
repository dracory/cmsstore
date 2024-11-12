package cmsstore

import "errors"

func MenuQuery() MenuQueryInterface {
	return &menuQuery{
		properties: make(map[string]interface{}),
	}
}

type menuQuery struct {
	properties map[string]interface{}
}

var _ MenuQueryInterface = (*menuQuery)(nil)

func (q *menuQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("menu query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("menu query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("menu query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("menu query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("menu query. limit cannot be negative")
	}

	if q.HasHandle() && q.Handle() == "" {
		return errors.New("menu query. handle cannot be empty")
	}

	if q.HasNameLike() && q.NameLike() == "" {
		return errors.New("menu query. name_like cannot be empty")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("menu query. offset cannot be negative")
	}

	if q.HasSiteID() && q.SiteID() == "" {
		return errors.New("menu query. site_id cannot be empty")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("menu query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("menu query. status_in cannot be empty array")
	}

	return nil
}

func (q *menuQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *menuQuery) SetColumns(columns []string) MenuQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *menuQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *menuQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *menuQuery) SetCountOnly(countOnly bool) MenuQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *menuQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *menuQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *menuQuery) SetCreatedAtGte(createdAtGte string) MenuQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *menuQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *menuQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *menuQuery) SetCreatedAtLte(createdAtLte string) MenuQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *menuQuery) HasHandle() bool {
	return q.hasProperty("handle")
}

func (q *menuQuery) Handle() string {
	return q.properties["handle"].(string)
}

func (q *menuQuery) SetHandle(handle string) MenuQueryInterface {
	q.properties["handle"] = handle
	return q
}

func (q *menuQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *menuQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *menuQuery) SetID(id string) MenuQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *menuQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *menuQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *menuQuery) SetIDIn(idIn []string) MenuQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *menuQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *menuQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *menuQuery) SetLimit(limit int) MenuQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *menuQuery) HasNameLike() bool {
	return q.hasProperty("name_like")
}

func (q *menuQuery) NameLike() string {
	return q.properties["name_like"].(string)
}

func (q *menuQuery) SetNameLike(nameLike string) MenuQueryInterface {
	q.properties["name_like"] = nameLike
	return q
}

func (q *menuQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *menuQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *menuQuery) SetOffset(offset int) MenuQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *menuQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *menuQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *menuQuery) SetOrderBy(orderBy string) MenuQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *menuQuery) HasSiteID() bool {
	return q.hasProperty("site_id")
}

func (q *menuQuery) SiteID() string {
	return q.properties["site_id"].(string)
}

func (q *menuQuery) SetSiteID(siteID string) MenuQueryInterface {
	q.properties["site_id"] = siteID
	return q
}

func (q *menuQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

func (q *menuQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

func (q *menuQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) MenuQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

func (q *menuQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *menuQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *menuQuery) SetSortOrder(sortOrder string) MenuQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *menuQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *menuQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *menuQuery) SetStatus(status string) MenuQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *menuQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *menuQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *menuQuery) SetStatusIn(statusIn []string) MenuQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *menuQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
