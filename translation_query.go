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
	if !q.hasProperty(propertyKeyColumns) {
		return []string{}
	}

	return q.properties[propertyKeyColumns].([]string)
}

func (q *translationQuery) SetColumns(columns []string) TranslationQueryInterface {
	q.properties[propertyKeyColumns] = columns
	return q
}

func (q *translationQuery) HasCountOnly() bool {
	return q.hasProperty(propertyKeyCountOnly)
}

func (q *translationQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties[propertyKeyCountOnly].(bool)
	}

	return false
}

func (q *translationQuery) SetCountOnly(countOnly bool) TranslationQueryInterface {
	q.properties[propertyKeyCountOnly] = countOnly
	return q
}

func (q *translationQuery) HasCreatedAtGte() bool {
	return q.hasProperty(propertyKeyCreatedAtGte)
}

func (q *translationQuery) CreatedAtGte() string {
	return q.properties[propertyKeyCreatedAtGte].(string)
}

func (q *translationQuery) SetCreatedAtGte(createdAtGte string) TranslationQueryInterface {
	q.properties[propertyKeyCreatedAtGte] = createdAtGte
	return q
}

func (q *translationQuery) HasCreatedAtLte() bool {
	return q.hasProperty(propertyKeyCreatedAtLte)
}

func (q *translationQuery) CreatedAtLte() string {
	return q.properties[propertyKeyCreatedAtLte].(string)
}

func (q *translationQuery) SetCreatedAtLte(createdAtLte string) TranslationQueryInterface {
	q.properties[propertyKeyCreatedAtLte] = createdAtLte
	return q
}

func (q *translationQuery) HasHandle() bool {
	return q.hasProperty(propertyKeyHandle)
}

func (q *translationQuery) Handle() string {
	return q.properties[propertyKeyHandle].(string)
}

func (q *translationQuery) SetHandle(handle string) TranslationQueryInterface {
	q.properties[propertyKeyHandle] = handle
	return q
}

func (q *translationQuery) HasHandleOrID() bool {
	return q.hasProperty(propertyKeyHandleOrID)
}

func (q *translationQuery) HandleOrID() string {
	return q.properties[propertyKeyHandleOrID].(string)
}

func (q *translationQuery) SetHandleOrID(handleOrID string) TranslationQueryInterface {
	q.properties[propertyKeyHandleOrID] = handleOrID
	return q
}

func (q *translationQuery) HasID() bool {
	return q.hasProperty(propertyKeyId)
}

func (q *translationQuery) ID() string {
	return q.properties[propertyKeyId].(string)
}

func (q *translationQuery) SetID(id string) TranslationQueryInterface {
	q.properties[propertyKeyId] = id
	return q
}

func (q *translationQuery) HasIDIn() bool {
	return q.hasProperty(propertyKeyIdIn)
}

func (q *translationQuery) IDIn() []string {
	return q.properties[propertyKeyIdIn].([]string)
}

func (q *translationQuery) SetIDIn(idIn []string) TranslationQueryInterface {
	q.properties[propertyKeyIdIn] = idIn
	return q
}

func (q *translationQuery) HasLimit() bool {
	return q.hasProperty(propertyKeyLimit)
}

func (q *translationQuery) Limit() int {
	return q.properties[propertyKeyLimit].(int)
}

func (q *translationQuery) SetLimit(limit int) TranslationQueryInterface {
	q.properties[propertyKeyLimit] = limit
	return q
}

func (q *translationQuery) HasNameLike() bool {
	return q.hasProperty(propertyKeyNameLike)
}

func (q *translationQuery) NameLike() string {
	return q.properties[propertyKeyNameLike].(string)
}

func (q *translationQuery) SetNameLike(nameLike string) TranslationQueryInterface {
	q.properties[propertyKeyNameLike] = nameLike
	return q
}

func (q *translationQuery) HasOffset() bool {
	return q.hasProperty(propertyKeyOffset)
}

func (q *translationQuery) Offset() int {
	return q.properties[propertyKeyOffset].(int)
}

func (q *translationQuery) SetOffset(offset int) TranslationQueryInterface {
	q.properties[propertyKeyOffset] = offset
	return q
}

func (q *translationQuery) HasOrderBy() bool {
	return q.hasProperty(propertyKeyOrderBy)
}

func (q *translationQuery) OrderBy() string {
	return q.properties[propertyKeyOrderBy].(string)
}

func (q *translationQuery) SetOrderBy(orderBy string) TranslationQueryInterface {
	q.properties[propertyKeyOrderBy] = orderBy
	return q
}

func (q *translationQuery) HasSiteID() bool {
	return q.hasProperty(propertyKeySiteID)
}

func (q *translationQuery) SiteID() string {
	return q.properties[propertyKeySiteID].(string)
}

func (q *translationQuery) SetSiteID(siteID string) TranslationQueryInterface {
	q.properties[propertyKeySiteID] = siteID
	return q
}

func (q *translationQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty(propertyKeySoftDeleteIncluded)
}

func (q *translationQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties[propertyKeySoftDeleteIncluded].(bool)
}

func (q *translationQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) TranslationQueryInterface {
	q.properties[propertyKeySoftDeleteIncluded] = softDeleteIncluded
	return q
}

func (q *translationQuery) HasSortOrder() bool {
	return q.hasProperty(propertyKeySortOrder)
}

func (q *translationQuery) SortOrder() string {
	return q.properties[propertyKeySortOrder].(string)
}

func (q *translationQuery) SetSortOrder(sortOrder string) TranslationQueryInterface {
	q.properties[propertyKeySortOrder] = sortOrder
	return q
}

func (q *translationQuery) HasStatus() bool {
	return q.hasProperty(propertyKeyStatus)
}

func (q *translationQuery) Status() string {
	return q.properties[propertyKeyStatus].(string)
}

func (q *translationQuery) SetStatus(status string) TranslationQueryInterface {
	q.properties[propertyKeyStatus] = status
	return q
}

func (q *translationQuery) HasStatusIn() bool {
	return q.hasProperty(propertyKeyStatusIn)
}

func (q *translationQuery) StatusIn() []string {
	return q.properties[propertyKeyStatusIn].([]string)
}

func (q *translationQuery) SetStatusIn(statusIn []string) TranslationQueryInterface {
	q.properties[propertyKeyStatusIn] = statusIn
	return q
}

func (q *translationQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
