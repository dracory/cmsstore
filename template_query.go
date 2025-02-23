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
	if !q.hasProperty(propertyKeyColumns) {
		return []string{}
	}

	return q.properties[propertyKeyColumns].([]string)
}

func (q *templateQuery) SetColumns(columns []string) TemplateQueryInterface {
	q.properties[propertyKeyColumns] = columns
	return q
}

func (q *templateQuery) HasCountOnly() bool {
	return q.hasProperty(propertyKeyCountOnly)
}

func (q *templateQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties[propertyKeyCountOnly].(bool)
	}

	return false
}

func (q *templateQuery) SetCountOnly(countOnly bool) TemplateQueryInterface {
	q.properties[propertyKeyCountOnly] = countOnly
	return q
}

func (q *templateQuery) HasCreatedAtGte() bool {
	return q.hasProperty(propertyKeyCreatedAtGte)
}

func (q *templateQuery) CreatedAtGte() string {
	return q.properties[propertyKeyCreatedAtGte].(string)
}

func (q *templateQuery) SetCreatedAtGte(createdAtGte string) TemplateQueryInterface {
	q.properties[propertyKeyCreatedAtGte] = createdAtGte
	return q
}

func (q *templateQuery) HasCreatedAtLte() bool {
	return q.hasProperty(propertyKeyCreatedAtLte)
}

func (q *templateQuery) CreatedAtLte() string {
	return q.properties[propertyKeyCreatedAtLte].(string)
}

func (q *templateQuery) SetCreatedAtLte(createdAtLte string) TemplateQueryInterface {
	q.properties[propertyKeyCreatedAtLte] = createdAtLte
	return q
}

func (q *templateQuery) HasHandle() bool {
	return q.hasProperty(propertyKeyHandle)
}

func (q *templateQuery) Handle() string {
	return q.properties[propertyKeyHandle].(string)
}

func (q *templateQuery) SetHandle(handle string) TemplateQueryInterface {
	q.properties[propertyKeyHandle] = handle
	return q
}

func (q *templateQuery) HasID() bool {
	return q.hasProperty(propertyKeyId)
}

func (q *templateQuery) ID() string {
	return q.properties[propertyKeyId].(string)
}

func (q *templateQuery) SetID(id string) TemplateQueryInterface {
	q.properties[propertyKeyId] = id
	return q
}

func (q *templateQuery) HasIDIn() bool {
	return q.hasProperty(propertyKeyIdIn)
}

func (q *templateQuery) IDIn() []string {
	return q.properties[propertyKeyIdIn].([]string)
}

func (q *templateQuery) SetIDIn(idIn []string) TemplateQueryInterface {
	q.properties[propertyKeyIdIn] = idIn
	return q
}

func (q *templateQuery) HasLimit() bool {
	return q.hasProperty(propertyKeyLimit)
}

func (q *templateQuery) Limit() int {
	return q.properties[propertyKeyLimit].(int)
}

func (q *templateQuery) SetLimit(limit int) TemplateQueryInterface {
	q.properties[propertyKeyLimit] = limit
	return q
}

func (q *templateQuery) HasNameLike() bool {
	return q.hasProperty(propertyKeyNameLike)
}

func (q *templateQuery) NameLike() string {
	return q.properties[propertyKeyNameLike].(string)
}

func (q *templateQuery) SetNameLike(nameLike string) TemplateQueryInterface {
	q.properties[propertyKeyNameLike] = nameLike
	return q
}

func (q *templateQuery) HasOffset() bool {
	return q.hasProperty(propertyKeyOffset)
}

func (q *templateQuery) Offset() int {
	return q.properties[propertyKeyOffset].(int)
}

func (q *templateQuery) SetOffset(offset int) TemplateQueryInterface {
	q.properties[propertyKeyOffset] = offset
	return q
}

func (q *templateQuery) HasOrderBy() bool {
	return q.hasProperty(propertyKeyOrderBy)
}

func (q *templateQuery) OrderBy() string {
	return q.properties[propertyKeyOrderBy].(string)
}

func (q *templateQuery) SetOrderBy(orderBy string) TemplateQueryInterface {
	q.properties[propertyKeyOrderBy] = orderBy
	return q
}

func (q *templateQuery) HasSiteID() bool {
	return q.hasProperty(propertyKeySiteID)
}

func (q *templateQuery) SiteID() string {
	return q.properties[propertyKeySiteID].(string)
}

func (q *templateQuery) SetSiteID(siteID string) TemplateQueryInterface {
	q.properties[propertyKeySiteID] = siteID
	return q
}

func (q *templateQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty(propertyKeySoftDeleteIncluded)
}

func (q *templateQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties[propertyKeySoftDeleteIncluded].(bool)
}

func (q *templateQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) TemplateQueryInterface {
	q.properties[propertyKeySoftDeleteIncluded] = softDeleteIncluded
	return q
}

func (q *templateQuery) HasSortOrder() bool {
	return q.hasProperty(propertyKeySortOrder)
}

func (q *templateQuery) SortOrder() string {
	return q.properties[propertyKeySortOrder].(string)
}

func (q *templateQuery) SetSortOrder(sortOrder string) TemplateQueryInterface {
	q.properties[propertyKeySortOrder] = sortOrder
	return q
}

func (q *templateQuery) HasStatus() bool {
	return q.hasProperty(propertyKeyStatus)
}

func (q *templateQuery) Status() string {
	return q.properties[propertyKeyStatus].(string)
}

func (q *templateQuery) SetStatus(status string) TemplateQueryInterface {
	q.properties[propertyKeyStatus] = status
	return q
}

func (q *templateQuery) HasStatusIn() bool {
	return q.hasProperty(propertyKeyStatusIn)
}

func (q *templateQuery) StatusIn() []string {
	return q.properties[propertyKeyStatusIn].([]string)
}

func (q *templateQuery) SetStatusIn(statusIn []string) TemplateQueryInterface {
	q.properties[propertyKeyStatusIn] = statusIn
	return q
}

func (q *templateQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
