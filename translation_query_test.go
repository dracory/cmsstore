package cmsstore

import (
	"slices"
	"strings"
	"testing"
)

func TestTranslationQueryDefaults(t *testing.T) {
	query := TranslationQuery()

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
	if query.HasHandleOrID() {
		t.Error("Expected HasHandleOrID to be false")
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

func TestTranslationQueryColumns(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryCreatedAtGte(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryCreatedAtLte(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryHandle(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryHandleOrID(t *testing.T) {
	query := TranslationQuery()

	// Test default
	if query.HasHandleOrID() {
		t.Error("Expected HasHandleOrID to be false")
	}

	// Test setting value
	handleOrID := "test-handle-or-id"
	query.SetHandleOrID(handleOrID)
	if !query.HasHandleOrID() {
		t.Error("Expected HasHandleOrID to be true")
	}
	if query.HandleOrID() != handleOrID {
		t.Errorf("Expected HandleOrID %s, got %s", handleOrID, query.HandleOrID())
	}
}

func TestTranslationQueryID(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryIDIn(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryLimit(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryNameLike(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryOffset(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryOrderBy(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQuerySiteID(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQuerySoftDeletedIncluded(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQuerySortOrder(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryStatus(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryStatusIn(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryCountOnly(t *testing.T) {
	query := TranslationQuery()

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

func TestTranslationQueryValidation(t *testing.T) {
	query := TranslationQuery()

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
	} else if !strings.Contains(err.Error(), "created_at_gte cannot be empty") {
		t.Errorf("Expected error message to contain 'created_at_gte cannot be empty', got %s", err.Error())
	}

	// Test invalid created_at_lte
	query = TranslationQuery()
	query.SetCreatedAtLte("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty created_at_lte")
	} else if !strings.Contains(err.Error(), "created_at_lte cannot be empty") {
		t.Errorf("Expected error message to contain 'created_at_lte cannot be empty', got %s", err.Error())
	}

	// Test invalid handle
	query = TranslationQuery()
	query.SetHandle("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty handle")
	} else if !strings.Contains(err.Error(), "handle cannot be empty") {
		t.Errorf("Expected error message to contain 'handle cannot be empty', got %s", err.Error())
	}

	// Test invalid handle_or_id
	query = TranslationQuery()
	query.SetHandleOrID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty handle_or_id")
	} else if !strings.Contains(err.Error(), "handle_or_id cannot be empty") {
		t.Errorf("Expected error message to contain 'handle_or_id cannot be empty', got %s", err.Error())
	}

	// Test invalid id
	query = TranslationQuery()
	query.SetID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id")
	} else if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("Expected error message to contain 'id cannot be empty', got %s", err.Error())
	}

	// Test invalid id_in
	query = TranslationQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id_in")
	} else if !strings.Contains(err.Error(), "id_in cannot be empty array") {
		t.Errorf("Expected error message to contain 'id_in cannot be empty array', got %s", err.Error())
	}

	// Test invalid limit
	query = TranslationQuery()
	query.SetLimit(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative limit")
	} else if !strings.Contains(err.Error(), "limit cannot be negative") {
		t.Errorf("Expected error message to contain 'limit cannot be negative', got %s", err.Error())
	}

	// Test invalid name_like
	query = TranslationQuery()
	query.SetNameLike("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty name_like")
	} else if !strings.Contains(err.Error(), "name_like cannot be empty") {
		t.Errorf("Expected error message to contain 'name_like cannot be empty', got %s", err.Error())
	}

	// Test invalid offset
	query = TranslationQuery()
	query.SetOffset(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative offset")
	} else if !strings.Contains(err.Error(), "offset cannot be negative") {
		t.Errorf("Expected error message to contain 'offset cannot be negative', got %s", err.Error())
	}

	// Test invalid site_id
	query = TranslationQuery()
	query.SetSiteID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty site_id")
	} else if !strings.Contains(err.Error(), "site_id cannot be empty") {
		t.Errorf("Expected error message to contain 'site_id cannot be empty', got %s", err.Error())
	}

	// Test invalid status
	query = TranslationQuery()
	query.SetStatus("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status")
	} else if !strings.Contains(err.Error(), "status cannot be empty") {
		t.Errorf("Expected error message to contain 'status cannot be empty', got %s", err.Error())
	}

	// Test invalid status_in
	query = TranslationQuery()
	query.SetStatusIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status_in")
	} else if !strings.Contains(err.Error(), "status_in cannot be empty array") {
		t.Errorf("Expected error message to contain 'status_in cannot be empty array', got %s", err.Error())
	}
}
