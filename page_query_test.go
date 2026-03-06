package cmsstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPageQueryDefaults(t *testing.T) {
	query := PageQuery()

	// Test default values
	require.False(t, query.HasColumns())
	require.False(t, query.HasAlias())
	require.False(t, query.HasAliasLike())
	require.False(t, query.HasCreatedAtGte())
	require.False(t, query.HasCreatedAtLte())
	require.False(t, query.HasCountOnly())
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
	require.False(t, query.HasTemplateID())
	require.False(t, query.IsCountOnly())
	require.Empty(t, query.Columns())
}

func TestPageQueryColumns(t *testing.T) {
	query := PageQuery()

	// Test default columns
	require.False(t, query.HasColumns())
	require.Empty(t, query.Columns())

	// Test SetColumns
	columns := []string{"id", "name", "status"}
	query.SetColumns(columns)
	require.True(t, query.HasColumns())
	require.Equal(t, columns, query.Columns())
}

func TestPageQueryAlias(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasAlias())

	// Test setting value
	alias := "test-alias"
	query.SetAlias(alias)
	require.True(t, query.HasAlias())
	require.Equal(t, alias, query.Alias())
}

func TestPageQueryAliasLike(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasAliasLike())

	// Test setting value
	aliasLike := "test-alias"
	query.SetAliasLike(aliasLike)
	require.True(t, query.HasAliasLike())
	require.Equal(t, aliasLike, query.AliasLike())
}

func TestPageQueryCreatedAtGte(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasCreatedAtGte())

	// Test setting value
	query.SetCreatedAtGte("2023-12-25 10:00:00")
	require.True(t, query.HasCreatedAtGte())
	require.Equal(t, "2023-12-25 10:00:00", query.CreatedAtGte())
}

func TestPageQueryCreatedAtLte(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasCreatedAtLte())

	// Test setting value
	query.SetCreatedAtLte("2023-12-25 10:00:00")
	require.True(t, query.HasCreatedAtLte())
	require.Equal(t, "2023-12-25 10:00:00", query.CreatedAtLte())
}

func TestPageQueryHandle(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasHandle())

	// Test setting value
	handle := "test-handle"
	query.SetHandle(handle)
	require.True(t, query.HasHandle())
	require.Equal(t, handle, query.Handle())
}

func TestPageQueryID(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasID())

	// Test setting value
	id := "test-id"
	query.SetID(id)
	require.True(t, query.HasID())
	require.Equal(t, id, query.ID())
}

func TestPageQueryIDIn(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasIDIn())

	// Test setting value
	ids := []string{"id1", "id2", "id3"}
	query.SetIDIn(ids)
	require.True(t, query.HasIDIn())
	require.Equal(t, ids, query.IDIn())
}

func TestPageQueryLimit(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasLimit())

	// Test setting value
	limit := 10
	query.SetLimit(limit)
	require.True(t, query.HasLimit())
	require.Equal(t, limit, query.Limit())
}

func TestPageQueryNameLike(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasNameLike())

	// Test setting value
	nameLike := "test-name"
	query.SetNameLike(nameLike)
	require.True(t, query.HasNameLike())
	require.Equal(t, nameLike, query.NameLike())
}

func TestPageQueryOffset(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasOffset())

	// Test setting value
	offset := 5
	query.SetOffset(offset)
	require.True(t, query.HasOffset())
	require.Equal(t, offset, query.Offset())
}

func TestPageQueryOrderBy(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasOrderBy())

	// Test setting value
	orderBy := "name"
	query.SetOrderBy(orderBy)
	require.True(t, query.HasOrderBy())
	require.Equal(t, orderBy, query.OrderBy())
}

func TestPageQuerySiteID(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasSiteID())

	// Test setting value
	siteID := "site-123"
	query.SetSiteID(siteID)
	require.True(t, query.HasSiteID())
	require.Equal(t, siteID, query.SiteID())
}

func TestPageQuerySoftDeletedIncluded(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasSoftDeletedIncluded())
	require.False(t, query.SoftDeletedIncluded())

	// Test setting value
	query.SetSoftDeletedIncluded(true)
	require.True(t, query.HasSoftDeletedIncluded())
	require.True(t, query.SoftDeletedIncluded())
}

func TestPageQuerySortOrder(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasSortOrder())

	// Test setting value
	sortOrder := "asc"
	query.SetSortOrder(sortOrder)
	require.True(t, query.HasSortOrder())
	require.Equal(t, sortOrder, query.SortOrder())
}

func TestPageQueryStatus(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasStatus())

	// Test setting value
	status := "active"
	query.SetStatus(status)
	require.True(t, query.HasStatus())
	require.Equal(t, status, query.Status())
}

func TestPageQueryStatusIn(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasStatusIn())

	// Test setting value
	statuses := []string{"active", "inactive"}
	query.SetStatusIn(statuses)
	require.True(t, query.HasStatusIn())
	require.Equal(t, statuses, query.StatusIn())
}

func TestPageQueryTemplateID(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasTemplateID())

	// Test setting value
	templateID := "template-123"
	query.SetTemplateID(templateID)
	require.True(t, query.HasTemplateID())
	require.Equal(t, templateID, query.TemplateID())
}

func TestPageQueryCountOnly(t *testing.T) {
	query := PageQuery()

	// Test default
	require.False(t, query.HasCountOnly())
	require.False(t, query.IsCountOnly())

	// Test setting value
	query.SetCountOnly(true)
	require.True(t, query.HasCountOnly())
	require.True(t, query.IsCountOnly())
}

func TestPageQueryValidation(t *testing.T) {
	query := PageQuery()

	// Test valid query
	err := query.Validate()
	require.NoError(t, err)

	// Test invalid alias_like
	query.SetAliasLike("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "alias_like cannot be empty")

	// Test invalid created_at_gte
	query = PageQuery()
	query.SetCreatedAtGte("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "created_at_gte cannot be empty")

	// Test invalid created_at_lte
	query = PageQuery()
	query.SetCreatedAtLte("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "created_at_lte cannot be empty")

	// Test invalid id
	query = PageQuery()
	query.SetID("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "id cannot be empty")

	// Test invalid id_in
	query = PageQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "id_in cannot be empty array")

	// Test invalid limit
	query = PageQuery()
	query.SetLimit(-1)
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "limit cannot be negative")

	// Test invalid handle
	query = PageQuery()
	query.SetHandle("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "handle cannot be empty")

	// Test invalid name_like
	query = PageQuery()
	query.SetNameLike("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "name_like cannot be empty")

	// Test invalid offset
	query = PageQuery()
	query.SetOffset(-1)
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset cannot be negative")

	// Test invalid order_by
	query = PageQuery()
	query.SetOrderBy("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "order_by cannot be empty")

	// Test invalid status
	query = PageQuery()
	query.SetStatus("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "status cannot be empty")

	// Test invalid status_in
	query = PageQuery()
	query.SetStatusIn([]string{})
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "status_in cannot be empty array")

	// Test invalid template_id
	query = PageQuery()
	query.SetTemplateID("")
	err = query.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "template_id cannot be empty")
}