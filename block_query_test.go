package cmsstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockQueryDefaults(t *testing.T) {
	query := BlockQuery()

	// Test default values
	require.False(t, query.HasCreatedAtGte())
	require.False(t, query.HasCreatedAtLte())
	require.False(t, query.HasHandle())
	require.False(t, query.HasID())
	require.False(t, query.HasIDIn())
	require.False(t, query.HasLimit())
	require.False(t, query.HasNameLike())
	require.False(t, query.HasOffset())
	require.False(t, query.HasOrderBy())
	require.False(t, query.HasPageID())
	require.False(t, query.HasParentID())
	require.False(t, query.HasSequence())
	require.False(t, query.HasSiteID())
	require.False(t, query.HasSoftDeleted())
	require.False(t, query.HasSortOrder())
	require.False(t, query.HasStatus())
	require.False(t, query.HasStatusIn())
	require.False(t, query.HasTemplateID())
	require.False(t, query.HasColumns())
	require.False(t, query.IsCountOnly())
	require.Empty(t, query.Columns())
}

func TestBlockQueryColumns(t *testing.T) {
	query := BlockQuery()

	// Test default columns
	require.False(t, query.HasColumns())
	require.Empty(t, query.Columns())

	// Test SetColumns
	columns := []string{"id", "name", "status"}
	query.SetColumns(columns)
	require.True(t, query.HasColumns())
	require.Equal(t, columns, query.Columns())
}

func TestBlockQueryCreatedAtGte(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasCreatedAtGte())

	// Test setting value
	query.SetCreatedAtGte("2023-12-25 10:00:00")
	require.True(t, query.HasCreatedAtGte())
	require.Equal(t, "2023-12-25 10:00:00", query.CreatedAtGte())
}

func TestBlockQueryCreatedAtLte(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasCreatedAtLte())

	// Test setting value
	query.SetCreatedAtLte("2023-12-25 10:00:00")
	require.True(t, query.HasCreatedAtLte())
	require.Equal(t, "2023-12-25 10:00:00", query.CreatedAtLte())
}

func TestBlockQueryHandle(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasHandle())

	// Test setting value
	handle := "test-handle"
	query.SetHandle(handle)
	require.True(t, query.HasHandle())
	require.Equal(t, handle, query.Handle())
}

func TestBlockQueryID(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasID())

	// Test setting value
	id := "test-id"
	query.SetID(id)
	require.True(t, query.HasID())
	require.Equal(t, id, query.ID())
}

func TestBlockQueryIDIn(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasIDIn())

	// Test setting value
	ids := []string{"id1", "id2", "id3"}
	query.SetIDIn(ids)
	require.True(t, query.HasIDIn())
	require.Equal(t, ids, query.IDIn())
}

func TestBlockQueryLimit(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasLimit())

	// Test setting value
	limit := 10
	query.SetLimit(limit)
	require.True(t, query.HasLimit())
	require.Equal(t, limit, query.Limit())
}

func TestBlockQueryNameLike(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasNameLike())

	// Test setting value
	nameLike := "test-name"
	query.SetNameLike(nameLike)
	require.True(t, query.HasNameLike())
	require.Equal(t, nameLike, query.NameLike())
}

func TestBlockQueryOffset(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasOffset())

	// Test setting value
	offset := 5
	query.SetOffset(offset)
	require.True(t, query.HasOffset())
	require.Equal(t, offset, query.Offset())
}

func TestBlockQueryOrderBy(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasOrderBy())

	// Test setting value
	orderBy := "name"
	query.SetOrderBy(orderBy)
	require.True(t, query.HasOrderBy())
	require.Equal(t, orderBy, query.OrderBy())
}

func TestBlockQueryPageID(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasPageID())

	// Test setting value
	pageID := "page-123"
	query.SetPageID(pageID)
	require.True(t, query.HasPageID())
	require.Equal(t, pageID, query.PageID())
}

func TestBlockQueryParentID(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasParentID())

	// Test setting value
	parentID := "parent-123"
	query.SetParentID(parentID)
	require.True(t, query.HasParentID())
	require.Equal(t, parentID, query.ParentID())
}

func TestBlockQuerySequence(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasSequence())

	// Test setting value
	sequence := 1
	query.SetSequence(sequence)
	require.True(t, query.HasSequence())
	require.Equal(t, sequence, query.Sequence())
}

func TestBlockQuerySiteID(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasSiteID())

	// Test setting value
	siteID := "site-123"
	query.SetSiteID(siteID)
	require.True(t, query.HasSiteID())
	require.Equal(t, siteID, query.SiteID())
}

func TestBlockQuerySoftDeleteIncluded(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasSoftDeleted())
	require.False(t, query.SoftDeleteIncluded())

	// Test setting value
	query.SetSoftDeleteIncluded(true)
	require.True(t, query.HasSoftDeleted())
	require.True(t, query.SoftDeleteIncluded())
}

func TestBlockQuerySortOrder(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasSortOrder())

	// Test setting value
	sortOrder := "asc"
	query.SetSortOrder(sortOrder)
	require.True(t, query.HasSortOrder())
	require.Equal(t, sortOrder, query.SortOrder())
}

func TestBlockQueryStatus(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasStatus())

	// Test setting value
	status := "active"
	query.SetStatus(status)
	require.True(t, query.HasStatus())
	require.Equal(t, status, query.Status())
}

func TestBlockQueryStatusIn(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasStatusIn())

	// Test setting value
	statuses := []string{"active", "inactive"}
	query.SetStatusIn(statuses)
	require.True(t, query.HasStatusIn())
	require.Equal(t, statuses, query.StatusIn())
}

func TestBlockQueryTemplateID(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.HasTemplateID())

	// Test setting value
	templateID := "template-123"
	query.SetTemplateID(templateID)
	require.True(t, query.HasTemplateID())
	require.Equal(t, templateID, query.TemplateID())
}

func TestBlockQueryCountOnly(t *testing.T) {
	query := BlockQuery()

	// Test default
	require.False(t, query.IsCountOnly())

	// Test setting value
	query.SetCountOnly(true)
	require.True(t, query.IsCountOnly())
}

func TestBlockQueryValidation(t *testing.T) {
	query := BlockQuery()

	// Test valid query
	err := query.Validate()
	require.NoError(t, err)

	// Test invalid created_at_gte
	query.SetCreatedAtGte("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "created_at_gte cannot be empty")

	// Test invalid created_at_lte
	query = BlockQuery()
	query.SetCreatedAtLte("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "created_at_lte cannot be empty")

	// Test invalid id
	query = BlockQuery()
	query.SetID("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "id cannot be empty")

	// Test invalid id_in
	query = BlockQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "id_in cannot be empty array")

	// Test invalid limit
	query = BlockQuery()
	query.SetLimit(-1)
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "limit cannot be negative")

	// Test invalid handle
	query = BlockQuery()
	query.SetHandle("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "handle cannot be empty")

	// Test invalid name_like
	query = BlockQuery()
	query.SetNameLike("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "name_like cannot be empty")

	// Test invalid status
	query = BlockQuery()
	query.SetStatus("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "status cannot be empty")

	// Test invalid offset
	query = BlockQuery()
	query.SetOffset(-1)
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset cannot be negative")
}