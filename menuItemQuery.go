package cmsstore

import "errors"

func MenuItemQuery() MenuItemQueryInterface {
	return &menuItemQuery{
		properties: make(map[string]interface{}),
	}
}

type menuItemQuery struct {
	properties map[string]interface{}
}

var _ MenuItemQueryInterface = (*menuItemQuery)(nil)

func (q *menuItemQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("menuItem query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("menuItem query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("menuItem query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("menuItem query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("menuItem query. limit cannot be negative")
	}

	if q.HasNameLike() && q.NameLike() == "" {
		return errors.New("menuItem query. name_like cannot be empty")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("menuItem query. offset cannot be negative")
	}

	if q.HasSiteID() && q.SiteID() == "" {
		return errors.New("menuItem query. site_id cannot be empty")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("menuItem query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("menuItem query. status_in cannot be empty array")
	}

	return nil
}

func (q *menuItemQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *menuItemQuery) SetColumns(columns []string) MenuItemQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *menuItemQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *menuItemQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *menuItemQuery) SetCountOnly(countOnly bool) MenuItemQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *menuItemQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *menuItemQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *menuItemQuery) SetCreatedAtGte(createdAtGte string) MenuItemQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *menuItemQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *menuItemQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *menuItemQuery) SetCreatedAtLte(createdAtLte string) MenuItemQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *menuItemQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *menuItemQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *menuItemQuery) SetID(id string) MenuItemQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *menuItemQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *menuItemQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *menuItemQuery) SetIDIn(idIn []string) MenuItemQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *menuItemQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *menuItemQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *menuItemQuery) SetLimit(limit int) MenuItemQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *menuItemQuery) HasNameLike() bool {
	return q.hasProperty("name_like")
}

func (q *menuItemQuery) NameLike() string {
	return q.properties["name_like"].(string)
}

func (q *menuItemQuery) SetNameLike(nameLike string) MenuItemQueryInterface {
	q.properties["name_like"] = nameLike
	return q
}

func (q *menuItemQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *menuItemQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *menuItemQuery) SetOffset(offset int) MenuItemQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *menuItemQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *menuItemQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *menuItemQuery) SetOrderBy(orderBy string) MenuItemQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *menuItemQuery) HasSiteID() bool {
	return q.hasProperty("site_id")
}

func (q *menuItemQuery) SiteID() string {
	return q.properties["site_id"].(string)
}

func (q *menuItemQuery) SetSiteID(siteID string) MenuItemQueryInterface {
	q.properties["site_id"] = siteID
	return q
}

func (q *menuItemQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

func (q *menuItemQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

func (q *menuItemQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) MenuItemQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

func (q *menuItemQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *menuItemQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *menuItemQuery) SetSortOrder(sortOrder string) MenuItemQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *menuItemQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *menuItemQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *menuItemQuery) SetStatus(status string) MenuItemQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *menuItemQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *menuItemQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *menuItemQuery) SetStatusIn(statusIn []string) MenuItemQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *menuItemQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
