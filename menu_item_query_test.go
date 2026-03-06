package cmsstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMenuItemQueryDefaults(t *testing.T) {
	query := MenuItemQuery()

	// Test default values
	require.False(t, query.HasCreatedAtGte())
	require.False(t, query.HasCreatedAtLte())
	require.False(t, query.HasID())
	require.False(t, query.HasIDIn())
	require.False(t, query.HasLimit())
	require.False(t, query.HasNameLike())
	require.False(t, query.HasOffset())
	require.False(t, query.HasOrderBy())
	require.False(t, query.HasMenuID())
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

func TestMenuItemQueryColumns(t *testing.T) {
	query := MenuItemQuery()

	// Test default columns
	require.False(t, query.HasColumns())
	require.Empty(t, query.Columns())

	// Test SetColumns
	columns := []string{"id", "name", "status"}
	query.SetColumns(columns)
	require.True(t, query.HasColumns())
	require.Equal(t, columns, query.Columns())
}

func TestMenuItemQueryCreatedAtGte(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasCreatedAtGte())

	// Test setting value
	query.SetCreatedAtGte("2023-12-25 10:00:00")
	require.True(t, query.HasCreatedAtGte())
	require.Equal(t, "2023-12-25 10:00:00", query.CreatedAtGte())
}

func TestMenuItemQueryCreatedAtLte(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasCreatedAtLte())

	// Test setting value
	query.SetCreatedAtLte("2023-12-25 10:00:00")
	require.True(t, query.HasCreatedAtLte())
	require.Equal(t, "2023-12-25 10:00:00", query.CreatedAtLte())
}

func TestMenuItemQueryID(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasID())

	// Test setting value
	id := "test-id"
	query.SetID(id)
	require.True(t, query.HasID())
	require.Equal(t, id, query.ID())
}

func TestMenuItemQueryIDIn(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasIDIn())

	// Test setting value
	ids := []string{"id1", "id2", "id3"}
	query.SetIDIn(ids)
	require.True(t, query.HasIDIn())
	require.Equal(t, ids, query.IDIn())
}

func TestMenuItemQueryLimit(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasLimit())

	// Test setting value
	limit := 10
	query.SetLimit(limit)
	require.True(t, query.HasLimit())
	require.Equal(t, limit, query.Limit())
}

func TestMenuItemQueryNameLike(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasNameLike())

	// Test setting value
	nameLike := "test-name"
	query.SetNameLike(nameLike)
	require.True(t, query.HasNameLike())
	require.Equal(t, nameLike, query.NameLike())
}

func TestMenuItemQueryOffset(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasOffset())

	// Test setting value
	offset := 5
	query.SetOffset(offset)
	require.True(t, query.HasOffset())
	require.Equal(t, offset, query.Offset())
}

func TestMenuItemQueryOrderBy(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasOrderBy())

	// Test setting value
	orderBy := "name"
	query.SetOrderBy(orderBy)
	require.True(t, query.HasOrderBy())
	require.Equal(t, orderBy, query.OrderBy())
}

func TestMenuItemQueryMenuID(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasMenuID())

	// Test setting value
	menuID := "menu-123"
	query.SetMenuID(menuID)
	require.True(t, query.HasMenuID())
	require.Equal(t, menuID, query.MenuID())
}

func TestMenuItemQuerySiteID(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasSiteID())

	// Test setting value
	siteID := "site-123"
	query.SetSiteID(siteID)
	require.True(t, query.HasSiteID())
	require.Equal(t, siteID, query.SiteID())
}

func TestMenuItemQuerySoftDeletedIncluded(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasSoftDeletedIncluded())
	require.False(t, query.SoftDeletedIncluded())

	// Test setting value
	query.SetSoftDeletedIncluded(true)
	require.True(t, query.HasSoftDeletedIncluded())
	require.True(t, query.SoftDeletedIncluded())
}

func TestMenuItemQuerySortOrder(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasSortOrder())

	// Test setting value
	sortOrder := "asc"
	query.SetSortOrder(sortOrder)
	require.True(t, query.HasSortOrder())
	require.Equal(t, sortOrder, query.SortOrder())
}

func TestMenuItemQueryStatus(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasStatus())

	// Test setting value
	status := "active"
	query.SetStatus(status)
	require.True(t, query.HasStatus())
	require.Equal(t, status, query.Status())
}

func TestMenuItemQueryStatusIn(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasStatusIn())

	// Test setting value
	statuses := []string{"active", "inactive"}
	query.SetStatusIn(statuses)
	require.True(t, query.HasStatusIn())
	require.Equal(t, statuses, query.StatusIn())
}

func TestMenuItemQueryCountOnly(t *testing.T) {
	query := MenuItemQuery()

	// Test default
	require.False(t, query.HasCountOnly())
	require.False(t, query.IsCountOnly())

	// Test setting value
	query.SetCountOnly(true)
	require.True(t, query.HasCountOnly())
	require.True(t, query.IsCountOnly())
}

func TestMenuItemQueryValidation(t *testing.T) {
	query := MenuItemQuery()

	// Test valid query
	err := query.Validate()
	require.NoError(t, err)

	// Test invalid created_at_gte
	query.SetCreatedAtGte("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "created_at_gte cannot be empty")

	// Test invalid created_at_lte
	query = MenuItemQuery()
	query.SetCreatedAtLte("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "created_at_lte cannot be empty")

	// Test invalid id
	query = MenuItemQuery()
	query.SetID("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "id cannot be empty")

	// Test invalid id_in
	query = MenuItemQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "id_in cannot be empty array")

	// Test invalid limit
	query = MenuItemQuery()
	query.SetLimit(-1)
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "limit cannot be negative")

	// Test invalid name_like
	query = MenuItemQuery()
	query.SetNameLike("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "name_like cannot be empty")

	// Test invalid offset
	query = MenuItemQuery()
	query.SetOffset(-1)
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset cannot be negative")

	// Test invalid site_id
	query = MenuItemQuery()
	query.SetSiteID("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "site_id cannot be empty")

	// Test invalid status
	query = MenuItemQuery()
	query.SetStatus("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "status cannot be empty")

	// Test invalid status_in
	query = MenuItemQuery()
	query.SetStatusIn([]string{})
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "status_in cannot be empty array")
}