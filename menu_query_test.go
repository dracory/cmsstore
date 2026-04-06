package cmsstore

import (
	"strings"
	"testing"
)

func TestMenuQueryDefaults(t *testing.T) {
	query := MenuQuery()

	// Test default values
	if query.HasCreatedAtGte() {
		t.Errorf("Expected HasCreatedAtGte to be false")
	}
	if query.HasCreatedAtLte() {
		t.Errorf("Expected HasCreatedAtLte to be false")
	}
	if query.HasHandle() {
		t.Errorf("Expected HasHandle to be false")
	}
	if query.HasID() {
		t.Errorf("Expected HasID to be false")
	}
	if query.HasIDIn() {
		t.Errorf("Expected HasIDIn to be false")
	}
	if query.HasLimit() {
		t.Errorf("Expected HasLimit to be false")
	}
	if query.HasNameLike() {
		t.Errorf("Expected HasNameLike to be false")
	}
	if query.HasOffset() {
		t.Errorf("Expected HasOffset to be false")
	}
	if query.HasOrderBy() {
		t.Errorf("Expected HasOrderBy to be false")
	}
	if query.HasSiteID() {
		t.Errorf("Expected HasSiteID to be false")
	}
	if query.HasSoftDeletedIncluded() {
		t.Errorf("Expected HasSoftDeletedIncluded to be false")
	}
	if query.HasSortOrder() {
		t.Errorf("Expected HasSortOrder to be false")
	}
	if query.HasStatus() {
		t.Errorf("Expected HasStatus to be false")
	}
	if query.HasStatusIn() {
		t.Errorf("Expected HasStatusIn to be false")
	}
	if query.HasColumns() {
		t.Errorf("Expected HasColumns to be false")
	}
	if query.HasCountOnly() {
		t.Errorf("Expected HasCountOnly to be false")
	}
	if query.IsCountOnly() {
		t.Errorf("Expected IsCountOnly to be false")
	}
	if len(query.Columns()) != 0 {
		t.Errorf("Expected empty Columns")
	}
}

func TestMenuQueryColumns(t *testing.T) {
	query := MenuQuery()

	// Test default columns
	if query.HasColumns() {
		t.Errorf("Expected HasColumns to be false")
	}
	if len(query.Columns()) != 0 {
		t.Errorf("Expected empty Columns")
	}

	// Test SetColumns
	columns := []string{"id", "name", "status"}
	query.SetColumns(columns)
	if !query.HasColumns() {
		t.Errorf("Expected HasColumns to be true")
	}
	if len(query.Columns()) != len(columns) {
		t.Errorf("Expected Columns length %d, got %d", len(columns), len(query.Columns()))
	}
	for i, col := range columns {
		if query.Columns()[i] != col {
			t.Errorf("Expected column %d %q, got %q", i, col, query.Columns()[i])
		}
	}
}

func TestMenuQueryCreatedAtGte(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasCreatedAtGte() {
		t.Errorf("Expected HasCreatedAtGte to be false")
	}

	// Test setting value
	query.SetCreatedAtGte("2023-12-25 10:00:00")
	if !query.HasCreatedAtGte() {
		t.Errorf("Expected HasCreatedAtGte to be true")
	}
	if query.CreatedAtGte() != "2023-12-25 10:00:00" {
		t.Errorf("Expected CreatedAtGte %q, got %q", "2023-12-25 10:00:00", query.CreatedAtGte())
	}
}

func TestMenuQueryCreatedAtLte(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasCreatedAtLte() {
		t.Errorf("Expected HasCreatedAtLte to be false")
	}

	// Test setting value
	query.SetCreatedAtLte("2023-12-25 10:00:00")
	if !query.HasCreatedAtLte() {
		t.Errorf("Expected HasCreatedAtLte to be true")
	}
	if query.CreatedAtLte() != "2023-12-25 10:00:00" {
		t.Errorf("Expected CreatedAtLte %q, got %q", "2023-12-25 10:00:00", query.CreatedAtLte())
	}
}

func TestMenuQueryHandle(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasHandle() {
		t.Errorf("Expected HasHandle to be false")
	}

	// Test setting value
	handle := "test-handle"
	query.SetHandle(handle)
	if !query.HasHandle() {
		t.Errorf("Expected HasHandle to be true")
	}
	if query.Handle() != handle {
		t.Errorf("Expected Handle %q, got %q", handle, query.Handle())
	}
}

func TestMenuQueryID(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasID() {
		t.Errorf("Expected HasID to be false")
	}

	// Test setting value
	id := "test-id"
	query.SetID(id)
	if !query.HasID() {
		t.Errorf("Expected HasID to be true")
	}
	if query.ID() != id {
		t.Errorf("Expected ID %q, got %q", id, query.ID())
	}
}

func TestMenuQueryIDIn(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasIDIn() {
		t.Errorf("Expected HasIDIn to be false")
	}

	// Test setting value
	ids := []string{"id1", "id2", "id3"}
	query.SetIDIn(ids)
	if !query.HasIDIn() {
		t.Errorf("Expected HasIDIn to be true")
	}
	if len(query.IDIn()) != len(ids) {
		t.Errorf("Expected IDIn length %d, got %d", len(ids), len(query.IDIn()))
	}
	for i, id := range ids {
		if query.IDIn()[i] != id {
			t.Errorf("Expected IDIn[%d] %q, got %q", i, id, query.IDIn()[i])
		}
	}
}

func TestMenuQueryLimit(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasLimit() {
		t.Errorf("Expected HasLimit to be false")
	}

	// Test setting value
	limit := 10
	query.SetLimit(limit)
	if !query.HasLimit() {
		t.Errorf("Expected HasLimit to be true")
	}
	if query.Limit() != limit {
		t.Errorf("Expected Limit %d, got %d", limit, query.Limit())
	}
}

func TestMenuQueryNameLike(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasNameLike() {
		t.Errorf("Expected HasNameLike to be false")
	}

	// Test setting value
	nameLike := "test-name"
	query.SetNameLike(nameLike)
	if !query.HasNameLike() {
		t.Errorf("Expected HasNameLike to be true")
	}
	if query.NameLike() != nameLike {
		t.Errorf("Expected NameLike %q, got %q", nameLike, query.NameLike())
	}
}

func TestMenuQueryOffset(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasOffset() {
		t.Errorf("Expected HasOffset to be false")
	}

	// Test setting value
	offset := 5
	query.SetOffset(offset)
	if !query.HasOffset() {
		t.Errorf("Expected HasOffset to be true")
	}
	if query.Offset() != offset {
		t.Errorf("Expected Offset %d, got %d", offset, query.Offset())
	}
}

func TestMenuQueryOrderBy(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasOrderBy() {
		t.Errorf("Expected HasOrderBy to be false")
	}

	// Test setting value
	orderBy := "name"
	query.SetOrderBy(orderBy)
	if !query.HasOrderBy() {
		t.Errorf("Expected HasOrderBy to be true")
	}
	if query.OrderBy() != orderBy {
		t.Errorf("Expected OrderBy %q, got %q", orderBy, query.OrderBy())
	}
}

func TestMenuQuerySiteID(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasSiteID() {
		t.Errorf("Expected HasSiteID to be false")
	}

	// Test setting value
	siteID := "site-123"
	query.SetSiteID(siteID)
	if !query.HasSiteID() {
		t.Errorf("Expected HasSiteID to be true")
	}
	if query.SiteID() != siteID {
		t.Errorf("Expected SiteID %q, got %q", siteID, query.SiteID())
	}
}

func TestMenuQuerySoftDeletedIncluded(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasSoftDeletedIncluded() {
		t.Errorf("Expected HasSoftDeletedIncluded to be false")
	}
	if query.SoftDeletedIncluded() {
		t.Errorf("Expected SoftDeletedIncluded to be false")
	}

	// Test setting value
	query.SetSoftDeletedIncluded(true)
	if !query.HasSoftDeletedIncluded() {
		t.Errorf("Expected HasSoftDeletedIncluded to be true")
	}
	if !query.SoftDeletedIncluded() {
		t.Errorf("Expected SoftDeletedIncluded to be true")
	}
}

func TestMenuQuerySortOrder(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasSortOrder() {
		t.Errorf("Expected HasSortOrder to be false")
	}

	// Test setting value
	sortOrder := "asc"
	query.SetSortOrder(sortOrder)
	if !query.HasSortOrder() {
		t.Errorf("Expected HasSortOrder to be true")
	}
	if query.SortOrder() != sortOrder {
		t.Errorf("Expected SortOrder %q, got %q", sortOrder, query.SortOrder())
	}
}

func TestMenuQueryStatus(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasStatus() {
		t.Errorf("Expected HasStatus to be false")
	}

	// Test setting value
	status := "active"
	query.SetStatus(status)
	if !query.HasStatus() {
		t.Errorf("Expected HasStatus to be true")
	}
	if query.Status() != status {
		t.Errorf("Expected Status %q, got %q", status, query.Status())
	}
}

func TestMenuQueryStatusIn(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasStatusIn() {
		t.Errorf("Expected HasStatusIn to be false")
	}

	// Test setting value
	statuses := []string{"active", "inactive"}
	query.SetStatusIn(statuses)
	if !query.HasStatusIn() {
		t.Errorf("Expected HasStatusIn to be true")
	}
	if len(query.StatusIn()) != len(statuses) {
		t.Errorf("Expected StatusIn length %d, got %d", len(statuses), len(query.StatusIn()))
	}
	for i, status := range statuses {
		if query.StatusIn()[i] != status {
			t.Errorf("Expected StatusIn[%d] %q, got %q", i, status, query.StatusIn()[i])
		}
	}
}

func TestMenuQueryCountOnly(t *testing.T) {
	query := MenuQuery()

	// Test default
	if query.HasCountOnly() {
		t.Errorf("Expected HasCountOnly to be false")
	}
	if query.IsCountOnly() {
		t.Errorf("Expected IsCountOnly to be false")
	}

	// Test setting value
	query.SetCountOnly(true)
	if !query.HasCountOnly() {
		t.Errorf("Expected HasCountOnly to be true")
	}
	if !query.IsCountOnly() {
		t.Errorf("Expected IsCountOnly to be true")
	}
}

func TestMenuQueryValidation(t *testing.T) {
	query := MenuQuery()

	// Test valid query
	err := query.Validate()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test invalid created_at_gte
	query.SetCreatedAtGte("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty created_at_gte")
	}
	if !strings.Contains(err.Error(), "created_at_gte cannot be empty") {
		t.Errorf("Expected error to contain 'created_at_gte cannot be empty'")
	}

	// Test invalid created_at_lte
	query = MenuQuery()
	query.SetCreatedAtLte("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty created_at_lte")
	}
	if !strings.Contains(err.Error(), "created_at_lte cannot be empty") {
		t.Errorf("Expected error to contain 'created_at_lte cannot be empty'")
	}

	// Test invalid id
	query = MenuQuery()
	query.SetID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id")
	}
	if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("Expected error to contain 'id cannot be empty'")
	}

	// Test invalid id_in
	query = MenuQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id_in")
	}
	if !strings.Contains(err.Error(), "id_in cannot be empty array") {
		t.Errorf("Expected error to contain 'id_in cannot be empty array'")
	}

	// Test invalid limit
	query = MenuQuery()
	query.SetLimit(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative limit")
	}
	if !strings.Contains(err.Error(), "limit cannot be negative") {
		t.Errorf("Expected error to contain 'limit cannot be negative'")
	}

	// Test invalid handle
	query = MenuQuery()
	query.SetHandle("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty handle")
	}
	if !strings.Contains(err.Error(), "handle cannot be empty") {
		t.Errorf("Expected error to contain 'handle cannot be empty'")
	}

	// Test invalid name_like
	query = MenuQuery()
	query.SetNameLike("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty name_like")
	}
	if !strings.Contains(err.Error(), "name_like cannot be empty") {
		t.Errorf("Expected error to contain 'name_like cannot be empty'")
	}

	// Test invalid offset
	query = MenuQuery()
	query.SetOffset(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative offset")
	}
	if !strings.Contains(err.Error(), "offset cannot be negative") {
		t.Errorf("Expected error to contain 'offset cannot be negative'")
	}

	// Test invalid site_id
	query = MenuQuery()
	query.SetSiteID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty site_id")
	}
	if !strings.Contains(err.Error(), "site_id cannot be empty") {
		t.Errorf("Expected error to contain 'site_id cannot be empty'")
	}

	// Test invalid status
	query = MenuQuery()
	query.SetStatus("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status")
	}
	if !strings.Contains(err.Error(), "status cannot be empty") {
		t.Errorf("Expected error to contain 'status cannot be empty'")
	}

	// Test invalid status_in
	query = MenuQuery()
	query.SetStatusIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status_in")
	}
	if !strings.Contains(err.Error(), "status_in cannot be empty array") {
		t.Errorf("Expected error to contain 'status_in cannot be empty array'")
	}
}
