package cmsstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMenuQueryDefaults(t *testing.T) {
	query := MenuQuery()

	// Test default values
	if query.HasCreatedAtGte() {
		t.Errorf("expected HasCreatedAtGte to be false")
	}
	if query.HasCreatedAtLte() {
		t.Errorf("expected HasCreatedAtLte to be false")
	}
	if query.HasHandle() {
		t.Errorf("expected HasHandle to be false")
	}
	if query.HasID() {
		t.Errorf("expected HasID to be false")
	}
	if query.HasIDIn() {
		t.Errorf("expected HasIDIn to be false")
	}
	if query.HasLimit() {
		t.Errorf("expected HasLimit to be false")
	}
	if query.HasNameLike() {
		t.Errorf("expected HasNameLike to be false")
	}
	if query.HasOffset() {
		t.Errorf("expected HasOffset to be false")
	}
	if query.HasOrderBy() {
		t.Errorf("expected HasOrderBy to be false")
	}
	if query.HasSiteID() {
		t.Errorf("expected HasSiteID to be false")
	}
	if query.HasSoftDeletedIncluded() {
		t.Errorf("expected HasSoftDeletedIncluded to be false")
	}
	if query.HasSortOrder() {
		t.Errorf("expected HasSortOrder to be false")
	}
	if query.HasStatus() {
		t.Errorf("expected HasStatus to be false")
	}
	if query.HasStatusIn() {
		t.Errorf("expected HasStatusIn to be false")
	}
	if query.HasColumns() {
		t.Errorf("expected HasColumns to be false")
	}
	if query.HasCountOnly() {
		t.Errorf("expected HasCountOnly to be false")
	}
	if query.IsCountOnly() {
		t.Errorf("expected IsCountOnly to be false")
	}
	if len(query.Columns()) != 0 {
		t.Errorf("expected empty Columns")
	}
}

func TestMenuQueryColumns(t *testing.T) {
	query := MenuQuery()

	// Test default columns
	if query.HasColumns() {
		t.Errorf("expected HasColumns to be false")
	}
	if len(query.Columns()) != 0 {
		t.Errorf("expected empty Columns")
	}

	// Test SetColumns
	columns := []string{"id", "name", "status"}
	query.SetColumns(columns)
	if !query.HasColumns() {
		t.Errorf("expected HasColumns to be true")
	}
	if len(query.Columns()) != len(columns) {
		t.Errorf("expected Columns length %d, got %d", len(columns), len(query.Columns()))
	}
	for i, col := range columns {
		if query.Columns()[i] != col {
			t.Errorf("expected column %d %q, got %q", i, col, query.Columns()[i])
		}
	}
}

func TestMenuQueryCreatedAtGte(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasCreatedAtGte() {
		t.Errorf("expected HasCreatedAtGte to be false")
	}

	// Test setting value
	query.SetCreatedAtGte("2023-12-25 10:00:00")
	if !query.HasCreatedAtGte() {
		t.Errorf("expected HasCreatedAtGte to be true")
	}
	if query.CreatedAtGte() != "2023-12-25 10:00:00" {
		t.Errorf("expected CreatedAtGte %q, got %q", "2023-12-25 10:00:00", query.CreatedAtGte())
	}
}

func TestMenuQueryCreatedAtLte(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasCreatedAtLte() {
		t.Errorf("expected HasCreatedAtLte to be false")
	}

	// Test setting value
	query.SetCreatedAtLte("2023-12-25 10:00:00")
	if !query.HasCreatedAtLte() {
		t.Errorf("expected HasCreatedAtLte to be true")
	}
	if query.CreatedAtLte() != "2023-12-25 10:00:00" {
		t.Errorf("expected CreatedAtLte %q, got %q", "2023-12-25 10:00:00", query.CreatedAtLte())
	}
}

func TestMenuQueryHandle(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasHandle() {
		t.Errorf("expected HasHandle to be false")
	}

	// Test setting value
	handle := "test-handle"
	query.SetHandle(handle)
	if !query.HasHandle() {
		t.Errorf("expected HasHandle to be true")
	}
	if query.Handle() != handle {
		t.Errorf("expected Handle %q, got %q", handle, query.Handle())
	}
}

func TestMenuQueryID(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasID() {
		t.Errorf("expected HasID to be false")
	}

	// Test setting value
	id := "test-id"
	query.SetID(id)
	if !query.HasID() {
		t.Errorf("expected HasID to be true")
	}
	if query.ID() != id {
		t.Errorf("expected ID %q, got %q", id, query.ID())
	}
}

func TestMenuQueryIDIn(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasIDIn() {
		t.Errorf("expected HasIDIn to be false")
	}

	// Test setting value
	ids := []string{"id1", "id2", "id3"}
	query.SetIDIn(ids)
	if !query.HasIDIn() {
		t.Errorf("expected HasIDIn to be true")
	}
	if len(query.IDIn()) != len(ids) {
		t.Errorf("expected IDIn length %d, got %d", len(ids), len(query.IDIn()))
	}
	for i, id := range ids {
		if query.IDIn()[i] != id {
			t.Errorf("expected IDIn[%d] %q, got %q", i, id, query.IDIn()[i])
		}
	}
}

func TestMenuQueryLimit(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasLimit() {
		t.Errorf("expected HasLimit to be false")
	}

	// Test setting value
	limit := 10
	query.SetLimit(limit)
	if !query.HasLimit() {
		t.Errorf("expected HasLimit to be true")
	}
	if query.Limit() != limit {
		t.Errorf("expected Limit %d, got %d", limit, query.Limit())
	}
}

func TestMenuQueryNameLike(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasNameLike() {
		t.Errorf("expected HasNameLike to be false")
	}

	// Test setting value
	nameLike := "test-name"
	query.SetNameLike(nameLike)
	if !query.HasNameLike() {
		t.Errorf("expected HasNameLike to be true")
	}
	if query.NameLike() != nameLike {
		t.Errorf("expected NameLike %q, got %q", nameLike, query.NameLike())
	}
}

func TestMenuQueryOffset(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasOffset() {
		t.Errorf("expected HasOffset to be false")
	}

	// Test setting value
	offset := 5
	query.SetOffset(offset)
	if !query.HasOffset() {
		t.Errorf("expected HasOffset to be true")
	}
	if query.Offset() != offset {
		t.Errorf("expected Offset %d, got %d", offset, query.Offset())
	}
}

func TestMenuQueryOrderBy(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasOrderBy() {
		t.Errorf("expected HasOrderBy to be false")
	}

	// Test setting value
	orderBy := "name"
	query.SetOrderBy(orderBy)
	if !query.HasOrderBy() {
		t.Errorf("expected HasOrderBy to be true")
	}
	if query.OrderBy() != orderBy {
		t.Errorf("expected OrderBy %q, got %q", orderBy, query.OrderBy())
	}
}

func TestMenuQuerySiteID(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasSiteID() {
		t.Errorf("expected HasSiteID to be false")
	}

	// Test setting value
	siteID := "site-123"
	query.SetSiteID(siteID)
	if !query.HasSiteID() {
		t.Errorf("expected HasSiteID to be true")
	}
	if query.SiteID() != siteID {
		t.Errorf("expected SiteID %q, got %q", siteID, query.SiteID())
	}
}

func TestMenuQuerySoftDeletedIncluded(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasSoftDeletedIncluded() {
		t.Errorf("expected HasSoftDeletedIncluded to be false")
	}
	if query.SoftDeletedIncluded() {
		t.Errorf("expected SoftDeletedIncluded to be false")
	}

	// Test setting value
	query.SetSoftDeletedIncluded(true)
	if !query.HasSoftDeletedIncluded() {
		t.Errorf("expected HasSoftDeletedIncluded to be true")
	}
	if !query.SoftDeletedIncluded() {
		t.Errorf("expected SoftDeletedIncluded to be true")
	}
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
