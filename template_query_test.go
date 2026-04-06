package cmsstore

import (
	"slices"
	"strings"
	"testing"
)

func TestTemplateQueryDefaults(t *testing.T) {
	query := TemplateQuery()

	// Test default values
	if query.HasCreatedAtGte() {
		t.Error("Expected HasCreatedAtGte to be false")
	}
	if query.HasCreatedAtLte() {
		t.Error("Expected HasCreatedAtLte to be false")
	}
	if query.HasHandle() {
		t.Error("Expected HasHandle to be false")
	}
	if query.HasID() {
		t.Error("Expected HasID to be false")
	}
	if query.HasIDIn() {
		t.Error("Expected HasIDIn to be false")
	}
	if query.HasLimit() {
		t.Error("Expected HasLimit to be false")
	}
	if query.HasNameLike() {
		t.Error("Expected HasNameLike to be false")
	}
	if query.HasOffset() {
		t.Error("Expected HasOffset to be false")
	}
	if query.HasOrderBy() {
		t.Error("Expected HasOrderBy to be false")
	}
	if query.HasSiteID() {
		t.Error("Expected HasSiteID to be false")
	}
	if query.HasSoftDeletedIncluded() {
		t.Error("Expected HasSoftDeletedIncluded to be false")
	}
	if query.HasSortOrder() {
		t.Error("Expected HasSortOrder to be false")
	}
	if query.HasStatus() {
		t.Error("Expected HasStatus to be false")
	}
	if query.HasStatusIn() {
		t.Error("Expected HasStatusIn to be false")
	}
	if query.HasColumns() {
		t.Error("Expected HasColumns to be false")
	}
	if query.HasCountOnly() {
		t.Error("Expected HasCountOnly to be false")
	}
	if query.IsCountOnly() {
		t.Error("Expected IsCountOnly to be false")
	}
	if len(query.Columns()) != 0 {
		t.Errorf("Expected empty Columns, got %v", query.Columns())
	}
}

func TestTemplateQueryColumns(t *testing.T) {
	query := TemplateQuery()

	// Test default columns
	if query.HasColumns() {
		t.Error("Expected HasColumns to be false")
	}
	if len(query.Columns()) != 0 {
		t.Errorf("Expected empty Columns, got %v", query.Columns())
	}

	// Test SetColumns
	columns := []string{"id", "name", "status"}
	query.SetColumns(columns)
	if !query.HasColumns() {
		t.Error("Expected HasColumns to be true")
	}
	if !slices.Equal(columns, query.Columns()) {
		t.Errorf("Expected Columns %v, got %v", columns, query.Columns())
	}
}

func TestTemplateQueryCreatedAtGte(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasCreatedAtGte() {
		t.Error("Expected HasCreatedAtGte to be false")
	}

	// Test setting value
	query.SetCreatedAtGte("2023-12-25 10:00:00")
	if !query.HasCreatedAtGte() {
		t.Error("Expected HasCreatedAtGte to be true")
	}
	if query.CreatedAtGte() != "2023-12-25 10:00:00" {
		t.Errorf("Expected CreatedAtGte %s, got %s", "2023-12-25 10:00:00", query.CreatedAtGte())
	}
}

func TestTemplateQueryCreatedAtLte(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasCreatedAtLte() {
		t.Error("Expected HasCreatedAtLte to be false")
	}

	// Test setting value
	query.SetCreatedAtLte("2023-12-25 10:00:00")
	if !query.HasCreatedAtLte() {
		t.Error("Expected HasCreatedAtLte to be true")
	}
	if query.CreatedAtLte() != "2023-12-25 10:00:00" {
		t.Errorf("Expected CreatedAtLte %s, got %s", "2023-12-25 10:00:00", query.CreatedAtLte())
	}
}

func TestTemplateQueryHandle(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasHandle() {
		t.Error("Expected HasHandle to be false")
	}

	// Test setting value
	handle := "test-handle"
	query.SetHandle(handle)
	if !query.HasHandle() {
		t.Error("Expected HasHandle to be true")
	}
	if query.Handle() != handle {
		t.Errorf("Expected Handle %s, got %s", handle, query.Handle())
	}
}

func TestTemplateQueryID(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasID() {
		t.Error("Expected HasID to be false")
	}

	// Test setting value
	id := "test-id"
	query.SetID(id)
	if !query.HasID() {
		t.Error("Expected HasID to be true")
	}
	if query.ID() != id {
		t.Errorf("Expected ID %s, got %s", id, query.ID())
	}
}

func TestTemplateQueryIDIn(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasIDIn() {
		t.Error("Expected HasIDIn to be false")
	}

	// Test setting value
	ids := []string{"id1", "id2", "id3"}
	query.SetIDIn(ids)
	if !query.HasIDIn() {
		t.Error("Expected HasIDIn to be true")
	}
	if !slices.Equal(ids, query.IDIn()) {
		t.Errorf("Expected IDIn %v, got %v", ids, query.IDIn())
	}
}

func TestTemplateQueryLimit(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasLimit() {
		t.Error("Expected HasLimit to be false")
	}

	// Test setting value
	limit := 10
	query.SetLimit(limit)
	if !query.HasLimit() {
		t.Error("Expected HasLimit to be true")
	}
	if query.Limit() != limit {
		t.Errorf("Expected Limit %d, got %d", limit, query.Limit())
	}
}

func TestTemplateQueryNameLike(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasNameLike() {
		t.Error("Expected HasNameLike to be false")
	}

	// Test setting value
	nameLike := "test-name"
	query.SetNameLike(nameLike)
	if !query.HasNameLike() {
		t.Error("Expected HasNameLike to be true")
	}
	if query.NameLike() != nameLike {
		t.Errorf("Expected NameLike %s, got %s", nameLike, query.NameLike())
	}
}

func TestTemplateQueryOffset(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasOffset() {
		t.Error("Expected HasOffset to be false")
	}

	// Test setting value
	offset := 5
	query.SetOffset(offset)
	if !query.HasOffset() {
		t.Error("Expected HasOffset to be true")
	}
	if query.Offset() != offset {
		t.Errorf("Expected Offset %d, got %d", offset, query.Offset())
	}
}

func TestTemplateQueryOrderBy(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasOrderBy() {
		t.Error("Expected HasOrderBy to be false")
	}

	// Test setting value
	orderBy := "name"
	query.SetOrderBy(orderBy)
	if !query.HasOrderBy() {
		t.Error("Expected HasOrderBy to be true")
	}
	if query.OrderBy() != orderBy {
		t.Errorf("Expected OrderBy %s, got %s", orderBy, query.OrderBy())
	}
}

func TestTemplateQuerySiteID(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasSiteID() {
		t.Error("Expected HasSiteID to be false")
	}

	// Test setting value
	siteID := "site-123"
	query.SetSiteID(siteID)
	if !query.HasSiteID() {
		t.Error("Expected HasSiteID to be true")
	}
	if query.SiteID() != siteID {
		t.Errorf("Expected SiteID %s, got %s", siteID, query.SiteID())
	}
}

func TestTemplateQuerySoftDeletedIncluded(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasSoftDeletedIncluded() {
		t.Error("Expected HasSoftDeletedIncluded to be false")
	}
	if query.SoftDeletedIncluded() {
		t.Error("Expected SoftDeletedIncluded to be false")
	}

	// Test setting value
	query.SetSoftDeletedIncluded(true)
	if !query.HasSoftDeletedIncluded() {
		t.Error("Expected HasSoftDeletedIncluded to be true")
	}
	if !query.SoftDeletedIncluded() {
		t.Error("Expected SoftDeletedIncluded to be true")
	}
}

func TestTemplateQuerySortOrder(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasSortOrder() {
		t.Error("Expected HasSortOrder to be false")
	}

	// Test setting value
	sortOrder := "asc"
	query.SetSortOrder(sortOrder)
	if !query.HasSortOrder() {
		t.Error("Expected HasSortOrder to be true")
	}
	if query.SortOrder() != sortOrder {
		t.Errorf("Expected SortOrder %s, got %s", sortOrder, query.SortOrder())
	}
}

func TestTemplateQueryStatus(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasStatus() {
		t.Error("Expected HasStatus to be false")
	}

	// Test setting value
	status := "active"
	query.SetStatus(status)
	if !query.HasStatus() {
		t.Error("Expected HasStatus to be true")
	}
	if query.Status() != status {
		t.Errorf("Expected Status %s, got %s", status, query.Status())
	}
}

func TestTemplateQueryStatusIn(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasStatusIn() {
		t.Error("Expected HasStatusIn to be false")
	}

	// Test setting value
	statuses := []string{"active", "inactive"}
	query.SetStatusIn(statuses)
	if !query.HasStatusIn() {
		t.Error("Expected HasStatusIn to be true")
	}
	if !slices.Equal(statuses, query.StatusIn()) {
		t.Errorf("Expected StatusIn %v, got %v", statuses, query.StatusIn())
	}
}

func TestTemplateQueryCountOnly(t *testing.T) {
	query := TemplateQuery()

	// Test default
	if query.HasCountOnly() {
		t.Error("Expected HasCountOnly to be false")
	}
	if query.IsCountOnly() {
		t.Error("Expected IsCountOnly to be false")
	}

	// Test setting value
	query.SetCountOnly(true)
	if !query.HasCountOnly() {
		t.Error("Expected HasCountOnly to be true")
	}
	if !query.IsCountOnly() {
		t.Error("Expected IsCountOnly to be true")
	}
}

func TestTemplateQueryValidation(t *testing.T) {
	query := TemplateQuery()

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
		t.Errorf("Expected error message to contain 'created_at_gte cannot be empty', got %s", err.Error())
	}

	// Test invalid created_at_lte
	query = TemplateQuery()
	query.SetCreatedAtLte("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty created_at_lte")
	}
	if !strings.Contains(err.Error(), "created_at_lte cannot be empty") {
		t.Errorf("Expected error message to contain 'created_at_lte cannot be empty', got %s", err.Error())
	}

	// Test invalid id
	query = TemplateQuery()
	query.SetID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id")
	}
	if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("Expected error message to contain 'id cannot be empty', got %s", err.Error())
	}

	// Test invalid id_in
	query = TemplateQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id_in")
	}
	if !strings.Contains(err.Error(), "id_in cannot be empty array") {
		t.Errorf("Expected error message to contain 'id_in cannot be empty array', got %s", err.Error())
	}

	// Test invalid limit
	query = TemplateQuery()
	query.SetLimit(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative limit")
	}
	if !strings.Contains(err.Error(), "limit cannot be negative") {
		t.Errorf("Expected error message to contain 'limit cannot be negative', got %s", err.Error())
	}

	// Test invalid handle
	query = TemplateQuery()
	query.SetHandle("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty handle")
	}
	if !strings.Contains(err.Error(), "handle cannot be empty") {
		t.Errorf("Expected error message to contain 'handle cannot be empty', got %s", err.Error())
	}

	// Test invalid name_like
	query = TemplateQuery()
	query.SetNameLike("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty name_like")
	}
	if !strings.Contains(err.Error(), "name_like cannot be empty") {
		t.Errorf("Expected error message to contain 'name_like cannot be empty', got %s", err.Error())
	}

	// Test invalid offset
	query = TemplateQuery()
	query.SetOffset(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative offset")
	}
	if !strings.Contains(err.Error(), "offset cannot be negative") {
		t.Errorf("Expected error message to contain 'offset cannot be negative', got %s", err.Error())
	}

	// Test invalid site_id
	query = TemplateQuery()
	query.SetSiteID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty site_id")
	}
	if !strings.Contains(err.Error(), "site_id cannot be empty") {
		t.Errorf("Expected error message to contain 'site_id cannot be empty', got %s", err.Error())
	}

	// Test invalid status
	query = TemplateQuery()
	query.SetStatus("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status")
	}
	if !strings.Contains(err.Error(), "status cannot be empty") {
		t.Errorf("Expected error message to contain 'status cannot be empty', got %s", err.Error())
	}

	// Test invalid status_in
	query = TemplateQuery()
	query.SetStatusIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status_in")
	}
	if !strings.Contains(err.Error(), "status_in cannot be empty array") {
		t.Errorf("Expected error message to contain 'status_in cannot be empty array', got %s", err.Error())
	}
}
