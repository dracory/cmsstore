package cmsstore

import "errors"

func BlockQuery() BlockQueryInterface {
	return &blockQuery{
		properties: make(map[string]interface{}),
	}
}

var _ BlockQueryInterface = (*blockQuery)(nil)

type blockQuery struct {
	properties map[string]interface{}
}

func (q *blockQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("block query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("block query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("block query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("block query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("block query. limit cannot be negative")
	}

	if q.HasHandle() && q.Handle() == "" {
		return errors.New("block query. handle cannot be empty")
	}

	if q.HasNameLike() && q.NameLike() == "" {
		return errors.New("block query. name_like cannot be empty")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("block query. status cannot be empty")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("block query. offset cannot be negative")
	}

	return nil
}

func (q *blockQuery) Columns() []string {
	if !q.hasProperty(propertyKeyColumns) {
		return []string{}
	}

	return q.properties[propertyKeyColumns].([]string)
}

func (q *blockQuery) SetColumns(columns []string) BlockQueryInterface {
	q.properties[propertyKeyColumns] = columns
	return q
}

func (q *blockQuery) HasCreatedAtGte() bool {
	return q.hasProperty(propertyKeyCreatedAtGte)
}

func (q *blockQuery) HasCreatedAtLte() bool {
	return q.hasProperty(propertyKeyCreatedAtLte)
}

func (q *blockQuery) HasHandle() bool {
	return q.hasProperty(propertyKeyHandle)
}

func (q *blockQuery) HasID() bool {
	return q.hasProperty(propertyKeyId)
}

func (q *blockQuery) HasIDIn() bool {
	return q.hasProperty(propertyKeyIdIn)
}

func (q *blockQuery) HasLimit() bool {
	return q.hasProperty(propertyKeyLimit)
}

func (q *blockQuery) HasNameLike() bool {
	return q.hasProperty(propertyKeyNameLike)
}

func (q *blockQuery) HasOffset() bool {
	return q.hasProperty(propertyKeyOffset)
}

func (q *blockQuery) HasOrderBy() bool {
	return q.hasProperty(propertyKeyOrderBy)
}

func (q *blockQuery) HasPageID() bool {
	return q.hasProperty(propertyKeyPageID)
}

func (q *blockQuery) HasParentID() bool {
	return q.hasProperty(propertyKeyParentID)
}

func (q *blockQuery) HasSequence() bool {
	return q.hasProperty(propertyKeySequence)
}

func (q *blockQuery) HasSiteID() bool {
	return q.hasProperty(propertyKeySiteID)
}

func (q *blockQuery) HasSoftDeleted() bool {
	return q.hasProperty(propertyKeySoftDeleteIncluded)
}

func (q *blockQuery) HasSortOrder() bool {
	return q.hasProperty(propertyKeySortOrder)
}

func (q *blockQuery) HasStatus() bool {
	return q.hasProperty(propertyKeyStatus)
}

func (q *blockQuery) HasStatusIn() bool {
	return q.hasProperty(propertyKeyStatusIn)
}

func (q *blockQuery) HasTemplateID() bool {
	return q.hasProperty(propertyKeyTemplateID)
}

func (q *blockQuery) IsCountOnly() bool {
	if q.hasProperty(propertyKeyCountOnly) {
		return q.properties[propertyKeyCountOnly].(bool)
	}

	return false
}

func (q *blockQuery) CreatedAtGte() string {
	return q.properties[propertyKeyCreatedAtGte].(string)
}

func (q *blockQuery) CreatedAtLte() string {
	return q.properties[propertyKeyCreatedAtLte].(string)
}

func (q *blockQuery) Handle() string {
	return q.properties[propertyKeyHandle].(string)
}

func (q *blockQuery) ID() string {
	return q.properties[propertyKeyId].(string)
}

func (q *blockQuery) IDIn() []string {
	return q.properties[propertyKeyIdIn].([]string)
}

func (q *blockQuery) Limit() int {
	return q.properties[propertyKeyLimit].(int)
}

func (q *blockQuery) NameLike() string {
	return q.properties[propertyKeyNameLike].(string)
}

func (q *blockQuery) Offset() int {
	return q.properties[propertyKeyOffset].(int)
}

func (q *blockQuery) OrderBy() string {
	return q.properties[propertyKeyOrderBy].(string)
}

func (q *blockQuery) PageID() string {
	return q.properties[propertyKeyPageID].(string)
}

func (q *blockQuery) ParentID() string {
	return q.properties[propertyKeyParentID].(string)
}

func (q *blockQuery) Sequence() int {
	return q.properties[propertyKeySequence].(int)
}

func (q *blockQuery) SiteID() string {
	return q.properties[propertyKeySiteID].(string)
}

func (q *blockQuery) SoftDeleteIncluded() bool {
	if !q.hasProperty(propertyKeySoftDeleteIncluded) {
		return false
	}

	return q.properties[propertyKeySoftDeleteIncluded].(bool)
}

func (q *blockQuery) SortOrder() string {
	return q.properties[propertyKeySortOrder].(string)
}

func (q *blockQuery) Status() string {
	return q.properties[propertyKeyStatus].(string)
}

func (q *blockQuery) StatusIn() []string {
	return q.properties[propertyKeyStatusIn].([]string)
}

func (q *blockQuery) TemplateID() string {
	return q.properties[propertyKeyTemplateID].(string)
}

func (q *blockQuery) SetCountOnly(countOnly bool) BlockQueryInterface {
	q.properties[propertyKeyCountOnly] = countOnly
	return q
}

func (q *blockQuery) SetID(id string) BlockQueryInterface {
	q.properties[propertyKeyId] = id
	return q
}

func (q *blockQuery) SetIDIn(idIn []string) BlockQueryInterface {
	q.properties[propertyKeyIdIn] = idIn
	return q
}

func (q *blockQuery) SetHandle(handle string) BlockQueryInterface {
	q.properties[propertyKeyHandle] = handle
	return q
}

func (q *blockQuery) SetLimit(limit int) BlockQueryInterface {
	q.properties[propertyKeyLimit] = limit
	return q
}

func (q *blockQuery) SetNameLike(nameLike string) BlockQueryInterface {
	q.properties[propertyKeyNameLike] = nameLike
	return q
}

func (q *blockQuery) SetOffset(offset int) BlockQueryInterface {
	q.properties[propertyKeyOffset] = offset
	return q
}

func (q *blockQuery) SetOrderBy(orderBy string) BlockQueryInterface {
	q.properties[propertyKeyOrderBy] = orderBy
	return q
}

func (q *blockQuery) SetPageID(pageID string) BlockQueryInterface {
	q.properties[propertyKeyPageID] = pageID
	return q
}

func (q *blockQuery) SetParentID(parentID string) BlockQueryInterface {
	q.properties[propertyKeyParentID] = parentID
	return q
}

func (q *blockQuery) SetSequence(sequence int) BlockQueryInterface {
	q.properties[propertyKeySequence] = sequence
	return q
}

func (q *blockQuery) SetSiteID(siteID string) BlockQueryInterface {
	q.properties[propertyKeySiteID] = siteID
	return q
}

func (q *blockQuery) SetSoftDeleteIncluded(SoftDeleteIncluded bool) BlockQueryInterface {
	q.properties[propertyKeySoftDeleteIncluded] = SoftDeleteIncluded
	return q
}

func (q *blockQuery) SetSortOrder(sortOrder string) BlockQueryInterface {
	q.properties[propertyKeySortOrder] = sortOrder
	return q
}

func (q *blockQuery) SetStatus(status string) BlockQueryInterface {
	q.properties[propertyKeyStatus] = status
	return q
}

func (q *blockQuery) SetStatusIn(statusIn []string) BlockQueryInterface {
	q.properties[propertyKeyStatusIn] = statusIn
	return q
}

func (q *blockQuery) SetTemplateID(templateID string) BlockQueryInterface {
	q.properties[propertyKeyTemplateID] = templateID
	return q
}

func (q *blockQuery) hasProperty(key string) bool {
	_, ok := q.properties[key]
	return ok
}
