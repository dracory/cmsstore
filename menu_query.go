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
	if !q.hasProperty(propertyKeyColumns) {
		return []string{}
	}

	return q.properties[propertyKeyColumns].([]string)
}

func (q *menuQuery) SetColumns(columns []string) MenuQueryInterface {
	q.properties[propertyKeyColumns] = columns
	return q
}

func (q *menuQuery) HasCountOnly() bool {
	return q.hasProperty(propertyKeyCountOnly)
}

func (q *menuQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties[propertyKeyCountOnly].(bool)
	}

	return false
}

func (q *menuQuery) SetCountOnly(countOnly bool) MenuQueryInterface {
	q.properties[propertyKeyCountOnly] = countOnly
	return q
}

func (q *menuQuery) HasCreatedAtGte() bool {
	return q.hasProperty(propertyKeyCreatedAtGte)
}

func (q *menuQuery) CreatedAtGte() string {
	return q.properties[propertyKeyCreatedAtGte].(string)
}

func (q *menuQuery) SetCreatedAtGte(createdAtGte string) MenuQueryInterface {
	q.properties[propertyKeyCreatedAtGte] = createdAtGte
	return q
}

func (q *menuQuery) HasCreatedAtLte() bool {
	return q.hasProperty(propertyKeyCreatedAtLte)
}

func (q *menuQuery) CreatedAtLte() string {
	return q.properties[propertyKeyCreatedAtLte].(string)
}

func (q *menuQuery) SetCreatedAtLte(createdAtLte string) MenuQueryInterface {
	q.properties[propertyKeyCreatedAtLte] = createdAtLte
	return q
}

func (q *menuQuery) HasHandle() bool {
	return q.hasProperty(propertyKeyHandle)
}

func (q *menuQuery) Handle() string {
	return q.properties[propertyKeyHandle].(string)
}

func (q *menuQuery) SetHandle(handle string) MenuQueryInterface {
	q.properties[propertyKeyHandle] = handle
	return q
}

func (q *menuQuery) HasID() bool {
	return q.hasProperty(propertyKeyId)
}

func (q *menuQuery) ID() string {
	return q.properties[propertyKeyId].(string)
}

func (q *menuQuery) SetID(id string) MenuQueryInterface {
	q.properties[propertyKeyId] = id
	return q
}

func (q *menuQuery) HasIDIn() bool {
	return q.hasProperty(propertyKeyIdIn)
}

func (q *menuQuery) IDIn() []string {
	return q.properties[propertyKeyIdIn].([]string)
}

func (q *menuQuery) SetIDIn(idIn []string) MenuQueryInterface {
	q.properties[propertyKeyIdIn] = idIn
	return q
}

func (q *menuQuery) HasLimit() bool {
	return q.hasProperty(propertyKeyLimit)
}

func (q *menuQuery) Limit() int {
	return q.properties[propertyKeyLimit].(int)
}

func (q *menuQuery) SetLimit(limit int) MenuQueryInterface {
	q.properties[propertyKeyLimit] = limit
	return q
}

func (q *menuQuery) HasNameLike() bool {
	return q.hasProperty(propertyKeyNameLike)
}

func (q *menuQuery) NameLike() string {
	return q.properties[propertyKeyNameLike].(string)
}

func (q *menuQuery) SetNameLike(nameLike string) MenuQueryInterface {
	q.properties[propertyKeyNameLike] = nameLike
	return q
}

func (q *menuQuery) HasOffset() bool {
	return q.hasProperty(propertyKeyOffset)
}

func (q *menuQuery) Offset() int {
	return q.properties[propertyKeyOffset].(int)
}

func (q *menuQuery) SetOffset(offset int) MenuQueryInterface {
	q.properties[propertyKeyOffset] = offset
	return q
}

func (q *menuQuery) HasOrderBy() bool {
	return q.hasProperty(propertyKeyOrderBy)
}

func (q *menuQuery) OrderBy() string {
	return q.properties[propertyKeyOrderBy].(string)
}

func (q *menuQuery) SetOrderBy(orderBy string) MenuQueryInterface {
	q.properties[propertyKeyOrderBy] = orderBy
	return q
}

func (q *menuQuery) HasSiteID() bool {
	return q.hasProperty(propertyKeySiteID)
}

func (q *menuQuery) SiteID() string {
	return q.properties[propertyKeySiteID].(string)
}

func (q *menuQuery) SetSiteID(siteID string) MenuQueryInterface {
	q.properties[propertyKeySiteID] = siteID
	return q
}

func (q *menuQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty(propertyKeySoftDeleteIncluded)
}

func (q *menuQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties[propertyKeySoftDeleteIncluded].(bool)
}

func (q *menuQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) MenuQueryInterface {
	q.properties[propertyKeySoftDeleteIncluded] = softDeleteIncluded
	return q
}

func (q *menuQuery) HasSortOrder() bool {
	return q.hasProperty(propertyKeySortOrder)
}

func (q *menuQuery) SortOrder() string {
	return q.properties[propertyKeySortOrder].(string)
}

func (q *menuQuery) SetSortOrder(sortOrder string) MenuQueryInterface {
	q.properties[propertyKeySortOrder] = sortOrder
	return q
}

func (q *menuQuery) HasStatus() bool {
	return q.hasProperty(propertyKeyStatus)
}

func (q *menuQuery) Status() string {
	return q.properties[propertyKeyStatus].(string)
}

func (q *menuQuery) SetStatus(status string) MenuQueryInterface {
	q.properties[propertyKeyStatus] = status
	return q
}

func (q *menuQuery) HasStatusIn() bool {
	return q.hasProperty(propertyKeyStatusIn)
}

func (q *menuQuery) StatusIn() []string {
	return q.properties[propertyKeyStatusIn].([]string)
}

func (q *menuQuery) SetStatusIn(statusIn []string) MenuQueryInterface {
	q.properties[propertyKeyStatusIn] = statusIn
	return q
}

func (q *menuQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
