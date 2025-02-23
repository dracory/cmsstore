package cmsstore

import "errors"

// MenuQuery returns a new instance of MenuQueryInterface.
func MenuQuery() MenuQueryInterface {
	return &menuQuery{
		properties: make(map[string]interface{}),
	}
}

// menuQuery is a struct that implements MenuQueryInterface.
type menuQuery struct {
	properties map[string]interface{}
}

// Ensuring menuQuery implements MenuQueryInterface.
var _ MenuQueryInterface = (*menuQuery)(nil)

// Validate checks the validity of the menuQuery struct properties.
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

// Columns returns the list of columns to be queried.
func (q *menuQuery) Columns() []string {
	if !q.hasProperty(propertyKeyColumns) {
		return []string{}
	}

	return q.properties[propertyKeyColumns].([]string)
}

// SetColumns sets the list of columns to be queried.
func (q *menuQuery) SetColumns(columns []string) MenuQueryInterface {
	q.properties[propertyKeyColumns] = columns
	return q
}

// HasCountOnly checks if CountOnly property is set.
func (q *menuQuery) HasCountOnly() bool {
	return q.hasProperty(propertyKeyCountOnly)
}

// IsCountOnly returns the value of CountOnly property.
func (q *menuQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties[propertyKeyCountOnly].(bool)
	}

	return false
}

// SetCountOnly sets the value of CountOnly property.
func (q *menuQuery) SetCountOnly(countOnly bool) MenuQueryInterface {
	q.properties[propertyKeyCountOnly] = countOnly
	return q
}

// HasCreatedAtGte checks if CreatedAtGte property is set.
func (q *menuQuery) HasCreatedAtGte() bool {
	return q.hasProperty(propertyKeyCreatedAtGte)
}

// CreatedAtGte returns the value of CreatedAtGte property.
func (q *menuQuery) CreatedAtGte() string {
	return q.properties[propertyKeyCreatedAtGte].(string)
}

// SetCreatedAtGte sets the value of CreatedAtGte property.
func (q *menuQuery) SetCreatedAtGte(createdAtGte string) MenuQueryInterface {
	q.properties[propertyKeyCreatedAtGte] = createdAtGte
	return q
}

// HasCreatedAtLte checks if CreatedAtLte property is set.
func (q *menuQuery) HasCreatedAtLte() bool {
	return q.hasProperty(propertyKeyCreatedAtLte)
}

// CreatedAtLte returns the value of CreatedAtLte property.
func (q *menuQuery) CreatedAtLte() string {
	return q.properties[propertyKeyCreatedAtLte].(string)
}

// SetCreatedAtLte sets the value of CreatedAtLte property.
func (q *menuQuery) SetCreatedAtLte(createdAtLte string) MenuQueryInterface {
	q.properties[propertyKeyCreatedAtLte] = createdAtLte
	return q
}

// HasHandle checks if Handle property is set.
func (q *menuQuery) HasHandle() bool {
	return q.hasProperty(propertyKeyHandle)
}

// Handle returns the value of Handle property.
func (q *menuQuery) Handle() string {
	return q.properties[propertyKeyHandle].(string)
}

// SetHandle sets the value of Handle property.
func (q *menuQuery) SetHandle(handle string) MenuQueryInterface {
	q.properties[propertyKeyHandle] = handle
	return q
}

// HasID checks if ID property is set.
func (q *menuQuery) HasID() bool {
	return q.hasProperty(propertyKeyId)
}

// ID returns the value of ID property.
func (q *menuQuery) ID() string {
	return q.properties[propertyKeyId].(string)
}

// SetID sets the value of ID property.
func (q *menuQuery) SetID(id string) MenuQueryInterface {
	q.properties[propertyKeyId] = id
	return q
}

// HasIDIn checks if IDIn property is set.
func (q *menuQuery) HasIDIn() bool {
	return q.hasProperty(propertyKeyIdIn)
}

// IDIn returns the value of IDIn property.
func (q *menuQuery) IDIn() []string {
	return q.properties[propertyKeyIdIn].([]string)
}

// SetIDIn sets the value of IDIn property.
func (q *menuQuery) SetIDIn(idIn []string) MenuQueryInterface {
	q.properties[propertyKeyIdIn] = idIn
	return q
}

// HasLimit checks if Limit property is set.
func (q *menuQuery) HasLimit() bool {
	return q.hasProperty(propertyKeyLimit)
}

// Limit returns the value of Limit property.
func (q *menuQuery) Limit() int {
	return q.properties[propertyKeyLimit].(int)
}

// SetLimit sets the value of Limit property.
func (q *menuQuery) SetLimit(limit int) MenuQueryInterface {
	q.properties[propertyKeyLimit] = limit
	return q
}

// HasNameLike checks if NameLike property is set.
func (q *menuQuery) HasNameLike() bool {
	return q.hasProperty(propertyKeyNameLike)
}

// NameLike returns the value of NameLike property.
func (q *menuQuery) NameLike() string {
	return q.properties[propertyKeyNameLike].(string)
}

// SetNameLike sets the value of NameLike property.
func (q *menuQuery) SetNameLike(nameLike string) MenuQueryInterface {
	q.properties[propertyKeyNameLike] = nameLike
	return q
}

// HasOffset checks if Offset property is set.
func (q *menuQuery) HasOffset() bool {
	return q.hasProperty(propertyKeyOffset)
}

// Offset returns the value of Offset property.
func (q *menuQuery) Offset() int {
	return q.properties[propertyKeyOffset].(int)
}

// SetOffset sets the value of Offset property.
func (q *menuQuery) SetOffset(offset int) MenuQueryInterface {
	q.properties[propertyKeyOffset] = offset
	return q
}

// HasOrderBy checks if OrderBy property is set.
func (q *menuQuery) HasOrderBy() bool {
	return q.hasProperty(propertyKeyOrderBy)
}

// OrderBy returns the value of OrderBy property.
func (q *menuQuery) OrderBy() string {
	return q.properties[propertyKeyOrderBy].(string)
}

// SetOrderBy sets the value of OrderBy property.
func (q *menuQuery) SetOrderBy(orderBy string) MenuQueryInterface {
	q.properties[propertyKeyOrderBy] = orderBy
	return q
}

// HasSiteID checks if SiteID property is set.
func (q *menuQuery) HasSiteID() bool {
	return q.hasProperty(propertyKeySiteID)
}

// SiteID returns the value of SiteID property.
func (q *menuQuery) SiteID() string {
	return q.properties[propertyKeySiteID].(string)
}

// SetSiteID sets the value of SiteID property.
func (q *menuQuery) SetSiteID(siteID string) MenuQueryInterface {
	q.properties[propertyKeySiteID] = siteID
	return q
}

// HasSoftDeletedIncluded checks if SoftDeletedIncluded property is set.
func (q *menuQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty(propertyKeySoftDeleteIncluded)
}

// SoftDeletedIncluded returns the value of SoftDeletedIncluded property.
func (q *menuQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties[propertyKeySoftDeleteIncluded].(bool)
}

// SetSoftDeletedIncluded sets the value of SoftDeletedIncluded property.
func (q *menuQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) MenuQueryInterface {
	q.properties[propertyKeySoftDeleteIncluded] = softDeleteIncluded
	return q
}

// HasSortOrder checks if SortOrder property is set.
func (q *menuQuery) HasSortOrder() bool {
	return q.hasProperty(propertyKeySortOrder)
}

// SortOrder returns the value of SortOrder property.
func (q *menuQuery) SortOrder() string {
	return q.properties[propertyKeySortOrder].(string)
}

// SetSortOrder sets the value of SortOrder property.
func (q *menuQuery) SetSortOrder(sortOrder string) MenuQueryInterface {
	q.properties[propertyKeySortOrder] = sortOrder
	return q
}

// HasStatus checks if Status property is set.
func (q *menuQuery) HasStatus() bool {
	return q.hasProperty(propertyKeyStatus)
}

// Status returns the value of Status property.
func (q *menuQuery) Status() string {
	return q.properties[propertyKeyStatus].(string)
}

// SetStatus sets the value of Status property.
func (q *menuQuery) SetStatus(status string) MenuQueryInterface {
	q.properties[propertyKeyStatus] = status
	return q
}

// HasStatusIn checks if StatusIn property is set.
func (q *menuQuery) HasStatusIn() bool {
	return q.hasProperty(propertyKeyStatusIn)
}

// StatusIn returns the value of StatusIn property.
func (q *menuQuery) StatusIn() []string {
	return q.properties[propertyKeyStatusIn].([]string)
}

// SetStatusIn sets the value of StatusIn property.
func (q *menuQuery) SetStatusIn(statusIn []string) MenuQueryInterface {
	q.properties[propertyKeyStatusIn] = statusIn
	return q
}

// hasProperty checks if a property exists in the menuQuery struct.
func (q *menuQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
