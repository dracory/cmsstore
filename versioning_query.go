package cmsstore

import "errors"

// newVersioningQuery creates a new versioning query
func newVersioningQuery() *versioningQuery {
	return &versioningQuery{
		properties: map[string]any{},
	}
}

// versioningQuery implements VersioningQueryInterface
type versioningQuery struct {
	properties map[string]any
}

// Validate validates the query parameters
func (q *versioningQuery) Validate() error {
	if q.HasEntityID() && q.EntityID() == "" {
		return errors.New("versioning query. entity_id cannot be empty")
	}

	if q.HasEntityType() && q.EntityType() == "" {
		return errors.New("versioning query. entity_type cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("versioning query. id cannot be empty")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("versioning query. limit cannot be negative")
	}

	if q.HasLimit() && q.Limit() < 1 {
		return errors.New("versioning query. limit cannot be less than 1")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("versioning query. offset cannot be negative")
	}

	return nil
}

// Columns returns the columns to select
func (q *versioningQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

// SetColumns sets the columns to select
func (q *versioningQuery) SetColumns(columns []string) VersioningQueryInterface {
	q.properties["columns"] = columns
	return q
}

// IsCountOnly returns true if only count is requested
func (q *versioningQuery) IsCountOnly() bool {
	if !q.hasProperty("count_only") {
		return false
	}

	return q.properties["count_only"].(bool)
}

// HasCountOnly returns true if count_only is set
func (q *versioningQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

// SetCountOnly sets whether to return only count
func (q *versioningQuery) SetCountOnly(countOnly bool) VersioningQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

// HasEntityID returns true if entity_id is set
func (q *versioningQuery) HasEntityID() bool {
	return q.hasProperty("entity_id")
}

// EntityID returns the entity ID
func (q *versioningQuery) EntityID() string {
	if !q.hasProperty("entity_id") {
		return ""
	}

	return q.properties["entity_id"].(string)
}

// SetEntityID sets the entity ID
func (q *versioningQuery) SetEntityID(entityID string) VersioningQueryInterface {
	q.properties["entity_id"] = entityID
	return q
}

// HasEntityType returns true if entity_type is set
func (q *versioningQuery) HasEntityType() bool {
	return q.hasProperty("entity_type")
}

// EntityType returns the entity type
func (q *versioningQuery) EntityType() string {
	if !q.hasProperty("entity_type") {
		return ""
	}

	return q.properties["entity_type"].(string)
}

// SetEntityType sets the entity type
func (q *versioningQuery) SetEntityType(entityType string) VersioningQueryInterface {
	q.properties["entity_type"] = entityType
	return q
}

// HasID returns true if id is set
func (q *versioningQuery) HasID() bool {
	return q.hasProperty("id")
}

// ID returns the versioning ID
func (q *versioningQuery) ID() string {
	if !q.hasProperty("id") {
		return ""
	}

	return q.properties["id"].(string)
}

// SetID sets the versioning ID
func (q *versioningQuery) SetID(id string) VersioningQueryInterface {
	q.properties["id"] = id
	return q
}

// HasLimit returns true if limit is set
func (q *versioningQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

// Limit returns the query limit
func (q *versioningQuery) Limit() int {
	if !q.hasProperty("limit") {
		return 0
	}

	return q.properties["limit"].(int)
}

// SetLimit sets the query limit
func (q *versioningQuery) SetLimit(limit int) VersioningQueryInterface {
	q.properties["limit"] = limit
	return q
}

// HasOffset returns true if offset is set
func (q *versioningQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

// Offset returns the query offset
func (q *versioningQuery) Offset() int64 {
	if !q.hasProperty("offset") {
		return 0
	}

	return q.properties["offset"].(int64)
}

// SetOffset sets the query offset
func (q *versioningQuery) SetOffset(offset int64) VersioningQueryInterface {
	q.properties["offset"] = offset
	return q
}

// HasOrderBy returns true if order_by is set
func (q *versioningQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

// OrderBy returns the order by field
func (q *versioningQuery) OrderBy() string {
	if !q.hasProperty("order_by") {
		return ""
	}

	return q.properties["order_by"].(string)
}

// SetOrderBy sets the order by field
func (q *versioningQuery) SetOrderBy(orderBy string) VersioningQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

// HasSortOrder returns true if sort_order is set
func (q *versioningQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

// SortOrder returns the sort order (ASC or DESC)
func (q *versioningQuery) SortOrder() string {
	if !q.hasProperty("sort_order") {
		return ""
	}

	return q.properties["sort_order"].(string)
}

// SetSortOrder sets the sort order (ASC or DESC)
func (q *versioningQuery) SetSortOrder(sortOrder string) VersioningQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

// HasSoftDeletedIncluded returns true if soft_deleted_included is set
func (q *versioningQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

// SoftDeletedIncluded returns true if soft deleted versionings should be included
func (q *versioningQuery) SoftDeletedIncluded() bool {
	if q.hasProperty("soft_deleted_included") {
		return q.properties["soft_deleted_included"].(bool)
	}

	return false
}

// SetSoftDeletedIncluded sets whether to include soft deleted versionings
func (q *versioningQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) VersioningQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

// hasProperty returns true if the property exists in the map
func (q *versioningQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
