package cmsstore

import "errors"

func TranslationQuery() TranslationQueryInterface {
	return &translationQuery{
		properties: make(map[string]interface{}),
	}
}

type translationQuery struct {
	properties map[string]interface{}
}

var _ TranslationQueryInterface = (*translationQuery)(nil)

func (q *translationQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("translation query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("translation query. created_at_lte cannot be empty")
	}

	if q.HasHandle() && q.Handle() == "" {
		return errors.New("translation query. handle cannot be empty")
	}

	if q.HasHandleOrID() && q.HandleOrID() == "" {
		return errors.New("translation query. handle_or_id cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("translation query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("translation query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("translation query. limit cannot be negative")
	}

	if q.HasHandle() && q.Handle() == "" {
		return errors.New("translation query. handle cannot be empty")
	}

	if q.HasNameLike() && q.NameLike() == "" {
		return errors.New("translation query. name_like cannot be empty")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("translation query. offset cannot be negative")
	}

	if q.HasSiteID() && q.SiteID() == "" {
		return errors.New("translation query. site_id cannot be empty")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("translation query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("translation query. status_in cannot be empty array")
	}

	return nil
}

func (q *translationQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *translationQuery) SetColumns(columns []string) TranslationQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *translationQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *translationQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *translationQuery) SetCountOnly(countOnly bool) TranslationQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *translationQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *translationQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *translationQuery) SetCreatedAtGte(createdAtGte string) TranslationQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *translationQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *translationQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *translationQuery) SetCreatedAtLte(createdAtLte string) TranslationQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *translationQuery) HasHandle() bool {
	return q.hasProperty("handle")
}

func (q *translationQuery) Handle() string {
	return q.properties["handle"].(string)
}

func (q *translationQuery) SetHandle(handle string) TranslationQueryInterface {
	q.properties["handle"] = handle
	return q
}

func (q *translationQuery) HasHandleOrID() bool {
	return q.hasProperty("handle_or_id")
}

func (q *translationQuery) HandleOrID() string {
	return q.properties["handle_or_id"].(string)
}

func (q *translationQuery) SetHandleOrID(handleOrID string) TranslationQueryInterface {
	q.properties["handle_or_id"] = handleOrID
	return q
}

func (q *translationQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *translationQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *translationQuery) SetID(id string) TranslationQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *translationQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *translationQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *translationQuery) SetIDIn(idIn []string) TranslationQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *translationQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *translationQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *translationQuery) SetLimit(limit int) TranslationQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *translationQuery) HasNameLike() bool {
	return q.hasProperty("name_like")
}

func (q *translationQuery) NameLike() string {
	return q.properties["name_like"].(string)
}

func (q *translationQuery) SetNameLike(nameLike string) TranslationQueryInterface {
	q.properties["name_like"] = nameLike
	return q
}

func (q *translationQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *translationQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *translationQuery) SetOffset(offset int) TranslationQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *translationQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *translationQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *translationQuery) SetOrderBy(orderBy string) TranslationQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *translationQuery) HasSiteID() bool {
	return q.hasProperty("site_id")
}

func (q *translationQuery) SiteID() string {
	return q.properties["site_id"].(string)
}

func (q *translationQuery) SetSiteID(siteID string) TranslationQueryInterface {
	q.properties["site_id"] = siteID
	return q
}

func (q *translationQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

func (q *translationQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

func (q *translationQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) TranslationQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

func (q *translationQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *translationQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *translationQuery) SetSortOrder(sortOrder string) TranslationQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *translationQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *translationQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *translationQuery) SetStatus(status string) TranslationQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *translationQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *translationQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *translationQuery) SetStatusIn(statusIn []string) TranslationQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *translationQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
