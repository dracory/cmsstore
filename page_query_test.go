package cmsstore

import (
	"slices"
	"strings"
	"testing"
)

func TestPageQueryDefaults(t *testing.T) {
	query := PageQuery()

	if query.HasColumns() {
		t.Error("Expected HasColumns to be false")
	}
	if query.HasAlias() {
		t.Error("Expected HasAlias to be false")
	}
	if query.HasAliasLike() {
		t.Error("Expected HasAliasLike to be false")
	}
	if query.HasCreatedAtGte() {
		t.Error("Expected HasCreatedAtGte to be false")
	}
	if query.HasCreatedAtLte() {
		t.Error("Expected HasCreatedAtLte to be false")
	}
	if query.HasCountOnly() {
		t.Error("Expected HasCountOnly to be false")
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
	if query.HasTemplateID() {
		t.Error("Expected HasTemplateID to be false")
	}
	if query.IsCountOnly() {
		t.Error("Expected IsCountOnly to be false")
	}
	if len(query.Columns()) != 0 {
		t.Errorf("Expected empty Columns, got %v", query.Columns())
	}
}

func TestPageQueryColumns(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryAlias(t *testing.T) {
	query := PageQuery()

	// Test default
	if query.HasAlias() {
		t.Error("Expected HasAlias to be false")
	}

	// Test setting value
	alias := "test-alias"
	query.SetAlias(alias)
	if !query.HasAlias() {
		t.Error("Expected HasAlias to be true")
	}
	if query.Alias() != alias {
		t.Errorf("Expected Alias %s, got %s", alias, query.Alias())
	}
}

func TestPageQueryAliasLike(t *testing.T) {
	query := PageQuery()

	// Test default
	if query.HasAliasLike() {
		t.Error("Expected HasAliasLike to be false")
	}

	// Test setting value
	aliasLike := "test-alias"
	query.SetAliasLike(aliasLike)
	if !query.HasAliasLike() {
		t.Error("Expected HasAliasLike to be true")
	}
	if query.AliasLike() != aliasLike {
		t.Errorf("Expected AliasLike %s, got %s", aliasLike, query.AliasLike())
	}
}

func TestPageQueryCreatedAtGte(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryCreatedAtLte(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryHandle(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryID(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryIDIn(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryLimit(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryNameLike(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryOffset(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryOrderBy(t *testing.T) {
	query := PageQuery()

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

func TestPageQuerySiteID(t *testing.T) {
	query := PageQuery()

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

func TestPageQuerySoftDeletedIncluded(t *testing.T) {
	query := PageQuery()

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

func TestPageQuerySortOrder(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryStatus(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryStatusIn(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryTemplateID(t *testing.T) {
	query := PageQuery()

	// Test default
	if query.HasTemplateID() {
		t.Error("Expected HasTemplateID to be false")
	}

	// Test setting value
	templateID := "template-123"
	query.SetTemplateID(templateID)
	if !query.HasTemplateID() {
		t.Error("Expected HasTemplateID to be true")
	}
	if query.TemplateID() != templateID {
		t.Errorf("Expected TemplateID %s, got %s", templateID, query.TemplateID())
	}
}

func TestPageQueryCountOnly(t *testing.T) {
	query := PageQuery()

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

func TestPageQueryValidation(t *testing.T) {
	query := PageQuery()

	// Test valid query
	err := query.Validate()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test invalid alias_like
	query.SetAliasLike("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty alias_like")
	}
	if !strings.Contains(err.Error(), "alias_like cannot be empty") {
		t.Errorf("Expected error message to contain 'alias_like cannot be empty', got %s", err.Error())
	}

	// Test invalid created_at_gte
	query = PageQuery()
	query.SetCreatedAtGte("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty created_at_gte")
	}
	if !strings.Contains(err.Error(), "created_at_gte cannot be empty") {
		t.Errorf("Expected error message to contain 'created_at_gte cannot be empty', got %s", err.Error())
	}

	// Test invalid created_at_lte
	query = PageQuery()
	query.SetCreatedAtLte("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty created_at_lte")
	}
	if !strings.Contains(err.Error(), "created_at_lte cannot be empty") {
		t.Errorf("Expected error message to contain 'created_at_lte cannot be empty', got %s", err.Error())
	}

	// Test invalid id
	query = PageQuery()
	query.SetID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id")
	}
	if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("Expected error message to contain 'id cannot be empty', got %s", err.Error())
	}

	// Test invalid id_in
	query = PageQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id_in")
	}
	if !strings.Contains(err.Error(), "id_in cannot be empty array") {
		t.Errorf("Expected error message to contain 'id_in cannot be empty array', got %s", err.Error())
	}

	// Test invalid limit
	query = PageQuery()
	query.SetLimit(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative limit")
	}
	if !strings.Contains(err.Error(), "limit cannot be negative") {
		t.Errorf("Expected error message to contain 'limit cannot be negative', got %s", err.Error())
	}

	// Test invalid handle
	query = PageQuery()
	query.SetHandle("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty handle")
	}
	if !strings.Contains(err.Error(), "handle cannot be empty") {
		t.Errorf("Expected error message to contain 'handle cannot be empty', got %s", err.Error())
	}

	// Test invalid name_like
	query = PageQuery()
	query.SetNameLike("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty name_like")
	}
	if !strings.Contains(err.Error(), "name_like cannot be empty") {
		t.Errorf("Expected error message to contain 'name_like cannot be empty', got %s", err.Error())
	}

	// Test invalid offset
	query = PageQuery()
	query.SetOffset(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative offset")
	}
	if !strings.Contains(err.Error(), "offset cannot be negative") {
		t.Errorf("Expected error message to contain 'offset cannot be negative', got %s", err.Error())
	}

	// Test invalid order_by
	query = PageQuery()
	query.SetOrderBy("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty order_by")
	}
	if !strings.Contains(err.Error(), "order_by cannot be empty") {
		t.Errorf("Expected error message to contain 'order_by cannot be empty', got %s", err.Error())
	}

	// Test invalid status
	query = PageQuery()
	query.SetStatus("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status")
	}
	if !strings.Contains(err.Error(), "status cannot be empty") {
		t.Errorf("Expected error message to contain 'status cannot be empty', got %s", err.Error())
	}

	// Test invalid status_in
	query = PageQuery()
	query.SetStatusIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status_in")
	}
	if !strings.Contains(err.Error(), "status_in cannot be empty array") {
		t.Errorf("Expected error message to contain 'status_in cannot be empty array', got %s", err.Error())
	}

	// Test invalid template_id
	query = PageQuery()
	query.SetTemplateID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty template_id")
	}
	if !strings.Contains(err.Error(), "template_id cannot be empty") {
		t.Errorf("Expected error message to contain 'template_id cannot be empty', got %s", err.Error())
	}
}
