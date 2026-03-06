package cmsstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMenuQueryDefaults(t *testing.T) {
	query := MenuQuery()

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
	require.False(t, query.HasSiteID())
	require.False(t, query.HasSoftDeletedIncluded())
	require.False(t, query.HasSortOrder())
	require.False(t, query.HasStatus())
	require.False(t, query.HasStatusIn())
	require.False(t, query.HasColumns())
	require.False(t, query.HasCountOnly())
	require.False(t, query.IsCountOnly())
	require.Empty(t, query.Columns())
}

func TestMenuQueryColumns(t *testing.T) {
	query := MenuQuery()

	// Test default columns
	require.False(t, query.HasColumns())
	require.Empty(t, query.Columns())

	// Test SetColumns
	columns := []string{"id", "name", "status"}
	query.SetColumns(columns)
	require.True(t, query.HasColumns())
	require.Equal(t, columns, query.Columns())
}

func TestMenuQueryCreatedAtGte(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasCreatedAtGte())

	// Test setting value
	query.SetCreatedAtGte("2023-12-25 10:00:00")
	require.True(t, query.HasCreatedAtGte())
	require.Equal(t, "2023-12-25 10:00:00", query.CreatedAtGte())
}

func TestMenuQueryCreatedAtLte(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasCreatedAtLte())

	// Test setting value
	query.SetCreatedAtLte("2023-12-25 10:00:00")
	require.True(t, query.HasCreatedAtLte())
	require.Equal(t, "2023-12-25 10:00:00", query.CreatedAtLte())
}

func TestMenuQueryHandle(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasHandle())

	// Test setting value
	handle := "test-handle"
	query.SetHandle(handle)
	require.True(t, query.HasHandle())
	require.Equal(t, handle, query.Handle())
}

func TestMenuQueryID(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasID())

	// Test setting value
	id := "test-id"
	query.SetID(id)
	require.True(t, query.HasID())
	require.Equal(t, id, query.ID())
}

func TestMenuQueryIDIn(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasIDIn())

	// Test setting value
	ids := []string{"id1", "id2", "id3"}
	query.SetIDIn(ids)
	require.True(t, query.HasIDIn())
	require.Equal(t, ids, query.IDIn())
}

func TestMenuQueryLimit(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasLimit())

	// Test setting value
	limit := 10
	query.SetLimit(limit)
	require.True(t, query.HasLimit())
	require.Equal(t, limit, query.Limit())
}

func TestMenuQueryNameLike(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasNameLike())

	// Test setting value
	nameLike := "test-name"
	query.SetNameLike(nameLike)
	require.True(t, query.HasNameLike())
	require.Equal(t, nameLike, query.NameLike())
}

func TestMenuQueryOffset(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasOffset())

	// Test setting value
	offset := 5
	query.SetOffset(offset)
	require.True(t, query.HasOffset())
	require.Equal(t, offset, query.Offset())
}

func TestMenuQueryOrderBy(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasOrderBy())

	// Test setting value
	orderBy := "name"
	query.SetOrderBy(orderBy)
	require.True(t, query.HasOrderBy())
	require.Equal(t, orderBy, query.OrderBy())
}

func TestMenuQuerySiteID(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasSiteID())

	// Test setting value
	siteID := "site-123"
	query.SetSiteID(siteID)
	require.True(t, query.HasSiteID())
	require.Equal(t, siteID, query.SiteID())
}

func TestMenuQuerySoftDeletedIncluded(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasSoftDeletedIncluded())
	require.False(t, query.SoftDeletedIncluded())

	// Test setting value
	query.SetSoftDeletedIncluded(true)
	require.True(t, query.HasSoftDeletedIncluded())
	require.True(t, query.SoftDeletedIncluded())
}

func TestMenuQuerySortOrder(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasSortOrder())

	// Test setting value
	sortOrder := "asc"
	query.SetSortOrder(sortOrder)
	require.True(t, query.HasSortOrder())
	require.Equal(t, sortOrder, query.SortOrder())
}

func TestMenuQueryStatus(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasStatus())

	// Test setting value
	status := "active"
	query.SetStatus(status)
	require.True(t, query.HasStatus())
	require.Equal(t, status, query.Status())
}

func TestMenuQueryStatusIn(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasStatusIn())

	// Test setting value
	statuses := []string{"active", "inactive"}
	query.SetStatusIn(statuses)
	require.True(t, query.HasStatusIn())
	require.Equal(t, statuses, query.StatusIn())
}

func TestMenuQueryCountOnly(t *testing.T) {
	query := MenuQuery()

	// Test default
	require.False(t, query.HasCountOnly())
	require.False(t, query.IsCountOnly())

	// Test setting value
	query.SetCountOnly(true)
	require.True(t, query.HasCountOnly())
	require.True(t, query.IsCountOnly())
}

func TestMenuQueryValidation(t *testing.T) {
	query := MenuQuery()

	// Test valid query
	err := query.Validate()
	require.NoError(t, err)

	// Test invalid created_at_gte
	query.SetCreatedAtGte("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "created_at_gte cannot be empty")

	// Test invalid created_at_lte
	query = MenuQuery()
	query.SetCreatedAtLte("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "created_at_lte cannot be empty")

	// Test invalid id
	query = MenuQuery()
	query.SetID("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "id cannot be empty")

	// Test invalid id_in
	query = MenuQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "id_in cannot be empty array")

	// Test invalid limit
	query = MenuQuery()
	query.SetLimit(-1)
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "limit cannot be negative")

	// Test invalid handle
	query = MenuQuery()
	query.SetHandle("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "handle cannot be empty")

	// Test invalid name_like
	query = MenuQuery()
	query.SetNameLike("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "name_like cannot be empty")

	// Test invalid offset
	query = MenuQuery()
	query.SetOffset(-1)
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset cannot be negative")

	// Test invalid site_id
	query = MenuQuery()
	query.SetSiteID("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "site_id cannot be empty")

	// Test invalid status
	query = MenuQuery()
	query.SetStatus("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "status cannot be empty")

	// Test invalid status_in
	query = MenuQuery()
	query.SetStatusIn([]string{})
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "status_in cannot be empty array")
}