package cmsstore

import "errors"

// MenuItemQuery returns a new instance of MenuItemQueryInterface.
// It initializes the properties map to store query parameters.
func MenuItemQuery() MenuItemQueryInterface {
	return &menuItemQuery{
		properties: make(map[string]interface{}),
	}
}

// menuItemQuery is a struct that implements the MenuItemQueryInterface.
// It uses a map to store various query parameters.
type menuItemQuery struct {
	properties map[string]interface{}
}

// Ensure menuItemQuery implements MenuItemQueryInterface.
var _ MenuItemQueryInterface = (*menuItemQuery)(nil)

// Validate checks the validity of the query parameters.
// It returns an error if any parameter is invalid.
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

// Columns returns the list of columns to be selected in the query.
// If no columns are specified, it returns an empty slice.
func (q *menuItemQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

// SetColumns sets the list of columns to be selected in the query.
func (q *menuItemQuery) SetColumns(columns []string) MenuItemQueryInterface {
	q.properties["columns"] = columns
	return q
}

// HasCountOnly checks if the count_only parameter is set.
func (q *menuItemQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

// IsCountOnly returns the value of the count_only parameter.
// If count_only is not set, it returns false.
func (q *menuItemQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

// SetCountOnly sets the count_only parameter.
func (q *menuItemQuery) SetCountOnly(countOnly bool) MenuItemQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

// HasCreatedAtGte checks if the created_at_gte parameter is set.
func (q *menuItemQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

// CreatedAtGte returns the value of the created_at_gte parameter.
func (q *menuItemQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

// SetCreatedAtGte sets the created_at_gte parameter.
func (q *menuItemQuery) SetCreatedAtGte(createdAtGte string) MenuItemQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

// HasCreatedAtLte checks if the created_at_lte parameter is set.
func (q *menuItemQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

// CreatedAtLte returns the value of the created_at_lte parameter.
func (q *menuItemQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

// SetCreatedAtLte sets the created_at_lte parameter.
func (q *menuItemQuery) SetCreatedAtLte(createdAtLte string) MenuItemQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

// HasID checks if the id parameter is set.
func (q *menuItemQuery) HasID() bool {
	return q.hasProperty("id")
}

// ID returns the value of the id parameter.
func (q *menuItemQuery) ID() string {
	return q.properties["id"].(string)
}

// SetID sets the id parameter.
func (q *menuItemQuery) SetID(id string) MenuItemQueryInterface {
	q.properties["id"] = id
	return q
}

// HasIDIn checks if the id_in parameter is set.
func (q *menuItemQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

// IDIn returns the value of the id_in parameter.
func (q *menuItemQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

// SetIDIn sets the id_in parameter.
func (q *menuItemQuery) SetIDIn(idIn []string) MenuItemQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

// HasLimit checks if the limit parameter is set.
func (q *menuItemQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

// Limit returns the value of the limit parameter.
func (q *menuItemQuery) Limit() int {
	return q.properties["limit"].(int)
}

// SetLimit sets the limit parameter.
func (q *menuItemQuery) SetLimit(limit int) MenuItemQueryInterface {
	q.properties["limit"] = limit
	return q
}

// HasMenuID checks if the menu_id parameter is set.
func (q *menuItemQuery) HasMenuID() bool {
	return q.hasProperty("menu_id")
}

// MenuID returns the value of the menu_id parameter.
func (q *menuItemQuery) MenuID() string {
	return q.properties["menu_id"].(string)
}

// SetMenuID sets the menu_id parameter.
func (q *menuItemQuery) SetMenuID(menuID string) MenuItemQueryInterface {
	q.properties["menu_id"] = menuID
	return q
}

// HasNameLike checks if the name_like parameter is set.
func (q *menuItemQuery) HasNameLike() bool {
	return q.hasProperty("name_like")
}

// NameLike returns the value of the name_like parameter.
func (q *menuItemQuery) NameLike() string {
	return q.properties["name_like"].(string)
}

// SetNameLike sets the name_like parameter.
func (q *menuItemQuery) SetNameLike(nameLike string) MenuItemQueryInterface {
	q.properties["name_like"] = nameLike
	return q
}

// HasOffset checks if the offset parameter is set.
func (q *menuItemQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

// Offset returns the value of the offset parameter.
func (q *menuItemQuery) Offset() int {
	return q.properties["offset"].(int)
}

// SetOffset sets the offset parameter.
func (q *menuItemQuery) SetOffset(offset int) MenuItemQueryInterface {
	q.properties["offset"] = offset
	return q
}

// HasOrderBy checks if the order_by parameter is set.
func (q *menuItemQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

// OrderBy returns the value of the order_by parameter.
func (q *menuItemQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

// SetOrderBy sets the order_by parameter.
func (q *menuItemQuery) SetOrderBy(orderBy string) MenuItemQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

// HasSiteID checks if the site_id parameter is set.
func (q *menuItemQuery) HasSiteID() bool {
	return q.hasProperty("site_id")
}

// SiteID returns the value of the site_id parameter.
func (q *menuItemQuery) SiteID() string {
	return q.properties["site_id"].(string)
}

// SetSiteID sets the site_id parameter.
func (q *menuItemQuery) SetSiteID(siteID string) MenuItemQueryInterface {
	q.properties["site_id"] = siteID
	return q
}

// HasSoftDeletedIncluded checks if the soft_delete_included parameter is set.
func (q *menuItemQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

// SoftDeletedIncluded returns the value of the soft_delete_included parameter.
// If soft_delete_included is not set, it returns false.
func (q *menuItemQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

// SetSoftDeletedIncluded sets the soft_delete_included parameter.
func (q *menuItemQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) MenuItemQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

// HasSortOrder checks if the sort_order parameter is set.
func (q *menuItemQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

// SortOrder returns the value of the sort_order parameter.
func (q *menuItemQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

// SetSortOrder sets the sort_order parameter.
func (q *menuItemQuery) SetSortOrder(sortOrder string) MenuItemQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

// HasStatus checks if the status parameter is set.
func (q *menuItemQuery) HasStatus() bool {
	return q.hasProperty("status")
}

// Status returns the value of the status parameter.
func (q *menuItemQuery) Status() string {
	return q.properties["status"].(string)
}

// SetStatus sets the status parameter.
func (q *menuItemQuery) SetStatus(status string) MenuItemQueryInterface {
	q.properties["status"] = status
	return q
}

// HasStatusIn checks if the status_in parameter is set.
func (q *menuItemQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

// StatusIn returns the value of the status_in parameter.
func (q *menuItemQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

// SetStatusIn sets the status_in parameter.
func (q *menuItemQuery) SetStatusIn(statusIn []string) MenuItemQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

// hasProperty checks if a given property key is set in the properties map.
func (q *menuItemQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
