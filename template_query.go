package cmsstore

import "errors"

func TemplateQuery() TemplateQueryInterface {
	return &templateQuery{
		properties: make(map[string]interface{}),
	}
}

type templateQuery struct {
	properties map[string]interface{}
}

var _ TemplateQueryInterface = (*templateQuery)(nil)

func (q *templateQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("template query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("template query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("template query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("template query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("template query. limit cannot be negative")
	}

	if q.HasHandle() && q.Handle() == "" {
		return errors.New("template query. handle cannot be empty")
	}

	if q.HasNameLike() && q.NameLike() == "" {
		return errors.New("template query. name_like cannot be empty")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("template query. offset cannot be negative")
	}

	if q.HasSiteID() && q.SiteID() == "" {
		return errors.New("template query. site_id cannot be empty")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("template query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("template query. status_in cannot be empty array")
	}

	return nil
}

func (q *templateQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *templateQuery) SetColumns(columns []string) TemplateQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *templateQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *templateQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *templateQuery) SetCountOnly(countOnly bool) TemplateQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *templateQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *templateQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *templateQuery) SetCreatedAtGte(createdAtGte string) TemplateQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *templateQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *templateQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *templateQuery) SetCreatedAtLte(createdAtLte string) TemplateQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *templateQuery) HasHandle() bool {
	return q.hasProperty("handle")
}

func (q *templateQuery) Handle() string {
	return q.properties["handle"].(string)
}

func (q *templateQuery) SetHandle(handle string) TemplateQueryInterface {
	q.properties["handle"] = handle
	return q
}

func (q *templateQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *templateQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *templateQuery) SetID(id string) TemplateQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *templateQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *templateQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *templateQuery) SetIDIn(idIn []string) TemplateQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *templateQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *templateQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *templateQuery) SetLimit(limit int) TemplateQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *templateQuery) HasNameLike() bool {
	return q.hasProperty("name_like")
}

func (q *templateQuery) NameLike() string {
	return q.properties["name_like"].(string)
}

func (q *templateQuery) SetNameLike(nameLike string) TemplateQueryInterface {
	q.properties["name_like"] = nameLike
	return q
}

func (q *templateQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *templateQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *templateQuery) SetOffset(offset int) TemplateQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *templateQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *templateQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *templateQuery) SetOrderBy(orderBy string) TemplateQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *templateQuery) HasSiteID() bool {
	return q.hasProperty("site_id")
}

func (q *templateQuery) SiteID() string {
	return q.properties["site_id"].(string)
}

func (q *templateQuery) SetSiteID(siteID string) TemplateQueryInterface {
	q.properties["site_id"] = siteID
	return q
}

func (q *templateQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

func (q *templateQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

func (q *templateQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) TemplateQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

func (q *templateQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *templateQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *templateQuery) SetSortOrder(sortOrder string) TemplateQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *templateQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *templateQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *templateQuery) SetStatus(status string) TemplateQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *templateQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *templateQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *templateQuery) SetStatusIn(statusIn []string) TemplateQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *templateQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
