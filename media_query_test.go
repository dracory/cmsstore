package cmsstore

import (
	"slices"
	"strings"
	"testing"
)

func TestMediaQueryDefaults(t *testing.T) {
	query := MediaQuery()

	if query.HasColumns() {
		t.Error("Expected HasColumns to be false")
	}
	if query.HasID() {
		t.Error("Expected HasID to be false")
	}
	if query.HasIDIn() {
		t.Error("Expected HasIDIn to be false")
	}
	if query.HasEntityID() {
		t.Error("Expected HasEntityID to be false")
	}
	if query.HasEntityType() {
		t.Error("Expected HasEntityType to be false")
	}
	if query.HasSiteID() {
		t.Error("Expected HasSiteID to be false")
	}
	if query.HasHandle() {
		t.Error("Expected HasHandle to be false")
	}
	if query.HasExtension() {
		t.Error("Expected HasExtension to be false")
	}
	if query.HasType() {
		t.Error("Expected HasType to be false")
	}
	if query.HasStatus() {
		t.Error("Expected HasStatus to be false")
	}
	if query.HasStatusIn() {
		t.Error("Expected HasStatusIn to be false")
	}
	if query.HasNameLike() {
		t.Error("Expected HasNameLike to be false")
	}
	if query.HasCountOnly() {
		t.Error("Expected HasCountOnly to be false")
	}
	if query.HasLimit() {
		t.Error("Expected HasLimit to be false")
	}
	if query.HasOffset() {
		t.Error("Expected HasOffset to be false")
	}
	if query.HasSortOrder() {
		t.Error("Expected HasSortOrder to be false")
	}
	if query.HasOrderBy() {
		t.Error("Expected HasOrderBy to be false")
	}
	if query.HasSoftDeletedIncluded() {
		t.Error("Expected HasSoftDeletedIncluded to be false")
	}
	if query.IsCountOnly() {
		t.Error("Expected IsCountOnly to be false")
	}
	if len(query.Columns()) != 0 {
		t.Errorf("Expected empty Columns, got %v", query.Columns())
	}
}

func TestMediaQueryColumns(t *testing.T) {
	query := MediaQuery()

	columns := []string{"id", "title", "status"}
	query.SetColumns(columns)
	if !query.HasColumns() {
		t.Error("Expected HasColumns to be true")
	}
	if !slices.Equal(columns, query.Columns()) {
		t.Errorf("Expected Columns %v, got %v", columns, query.Columns())
	}
}

func TestMediaQueryID(t *testing.T) {
	query := MediaQuery()

	query.SetID("test-id")
	if !query.HasID() {
		t.Error("Expected HasID to be true")
	}
	if query.ID() != "test-id" {
		t.Errorf("Expected ID %s, got %s", "test-id", query.ID())
	}
}

func TestMediaQueryIDIn(t *testing.T) {
	query := MediaQuery()

	ids := []string{"id1", "id2", "id3"}
	query.SetIDIn(ids)
	if !query.HasIDIn() {
		t.Error("Expected HasIDIn to be true")
	}
	if !slices.Equal(ids, query.IDIn()) {
		t.Errorf("Expected IDIn %v, got %v", ids, query.IDIn())
	}
}

func TestMediaQueryEntityID(t *testing.T) {
	query := MediaQuery()

	query.SetEntityID("entity-123")
	if !query.HasEntityID() {
		t.Error("Expected HasEntityID to be true")
	}
	if query.EntityID() != "entity-123" {
		t.Errorf("Expected EntityID %s, got %s", "entity-123", query.EntityID())
	}
}

func TestMediaQueryEntityType(t *testing.T) {
	query := MediaQuery()

	query.SetEntityType("page")
	if !query.HasEntityType() {
		t.Error("Expected HasEntityType to be true")
	}
	if query.EntityType() != "page" {
		t.Errorf("Expected EntityType %s, got %s", "page", query.EntityType())
	}
}

func TestMediaQuerySiteID(t *testing.T) {
	query := MediaQuery()

	query.SetSiteID("site-123")
	if !query.HasSiteID() {
		t.Error("Expected HasSiteID to be true")
	}
	if query.SiteID() != "site-123" {
		t.Errorf("Expected SiteID %s, got %s", "site-123", query.SiteID())
	}
}

func TestMediaQueryHandle(t *testing.T) {
	query := MediaQuery()

	query.SetHandle("test-handle")
	if !query.HasHandle() {
		t.Error("Expected HasHandle to be true")
	}
	if query.Handle() != "test-handle" {
		t.Errorf("Expected Handle %s, got %s", "test-handle", query.Handle())
	}
}

func TestMediaQueryExtension(t *testing.T) {
	query := MediaQuery()

	query.SetExtension("jpg")
	if !query.HasExtension() {
		t.Error("Expected HasExtension to be true")
	}
	if query.Extension() != "jpg" {
		t.Errorf("Expected Extension %s, got %s", "jpg", query.Extension())
	}
}

func TestMediaQueryType(t *testing.T) {
	query := MediaQuery()

	query.SetType("image/jpeg")
	if !query.HasType() {
		t.Error("Expected HasType to be true")
	}
	if query.Type() != "image/jpeg" {
		t.Errorf("Expected Type %s, got %s", "image/jpeg", query.Type())
	}
}

func TestMediaQueryStatus(t *testing.T) {
	query := MediaQuery()

	query.SetStatus("active")
	if !query.HasStatus() {
		t.Error("Expected HasStatus to be true")
	}
	if query.Status() != "active" {
		t.Errorf("Expected Status %s, got %s", "active", query.Status())
	}
}

func TestMediaQueryStatusIn(t *testing.T) {
	query := MediaQuery()

	statuses := []string{"active", "inactive"}
	query.SetStatusIn(statuses)
	if !query.HasStatusIn() {
		t.Error("Expected HasStatusIn to be true")
	}
	if !slices.Equal(statuses, query.StatusIn()) {
		t.Errorf("Expected StatusIn %v, got %v", statuses, query.StatusIn())
	}
}

func TestMediaQueryNameLike(t *testing.T) {
	query := MediaQuery()

	query.SetNameLike("test-name")
	if !query.HasNameLike() {
		t.Error("Expected HasNameLike to be true")
	}
	if query.NameLike() != "test-name" {
		t.Errorf("Expected NameLike %s, got %s", "test-name", query.NameLike())
	}
}

func TestMediaQueryCountOnly(t *testing.T) {
	query := MediaQuery()

	query.SetCountOnly(true)
	if !query.HasCountOnly() {
		t.Error("Expected HasCountOnly to be true")
	}
	if !query.IsCountOnly() {
		t.Error("Expected IsCountOnly to be true")
	}
}

func TestMediaQueryLimit(t *testing.T) {
	query := MediaQuery()

	query.SetLimit(10)
	if !query.HasLimit() {
		t.Error("Expected HasLimit to be true")
	}
	if query.Limit() != 10 {
		t.Errorf("Expected Limit %d, got %d", 10, query.Limit())
	}
}

func TestMediaQueryOffset(t *testing.T) {
	query := MediaQuery()

	query.SetOffset(5)
	if !query.HasOffset() {
		t.Error("Expected HasOffset to be true")
	}
	if query.Offset() != 5 {
		t.Errorf("Expected Offset %d, got %d", 5, query.Offset())
	}
}

func TestMediaQuerySortOrder(t *testing.T) {
	query := MediaQuery()

	query.SetSortOrder("asc")
	if !query.HasSortOrder() {
		t.Error("Expected HasSortOrder to be true")
	}
	if query.SortOrder() != "asc" {
		t.Errorf("Expected SortOrder %s, got %s", "asc", query.SortOrder())
	}
}

func TestMediaQueryOrderBy(t *testing.T) {
	query := MediaQuery()

	query.SetOrderBy("title")
	if !query.HasOrderBy() {
		t.Error("Expected HasOrderBy to be true")
	}
	if query.OrderBy() != "title" {
		t.Errorf("Expected OrderBy %s, got %s", "title", query.OrderBy())
	}
}

func TestMediaQuerySoftDeletedIncluded(t *testing.T) {
	query := MediaQuery()

	query.SetSoftDeletedIncluded(true)
	if !query.HasSoftDeletedIncluded() {
		t.Error("Expected HasSoftDeletedIncluded to be true")
	}
	if !query.SoftDeletedIncluded() {
		t.Error("Expected SoftDeletedIncluded to be true")
	}
}

func TestMediaQueryValidation(t *testing.T) {
	query := MediaQuery()

	err := query.Validate()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	query.SetID("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id")
	}
	if !strings.Contains(err.Error(), "id cannot be empty") {
		t.Errorf("Expected error message to contain 'id cannot be empty', got %s", err.Error())
	}

	query = MediaQuery()
	query.SetIDIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty id_in")
	}
	if !strings.Contains(err.Error(), "id_in cannot be empty array") {
		t.Errorf("Expected error message to contain 'id_in cannot be empty array', got %s", err.Error())
	}

	query = MediaQuery()
	query.SetHandle("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty handle")
	}
	if !strings.Contains(err.Error(), "handle cannot be empty") {
		t.Errorf("Expected error message to contain 'handle cannot be empty', got %s", err.Error())
	}

	query = MediaQuery()
	query.SetNameLike("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty name_like")
	}
	if !strings.Contains(err.Error(), "name_like cannot be empty") {
		t.Errorf("Expected error message to contain 'name_like cannot be empty', got %s", err.Error())
	}

	query = MediaQuery()
	query.SetLimit(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative limit")
	}
	if !strings.Contains(err.Error(), "limit cannot be negative") {
		t.Errorf("Expected error message to contain 'limit cannot be negative', got %s", err.Error())
	}

	query = MediaQuery()
	query.SetOffset(-1)
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for negative offset")
	}
	if !strings.Contains(err.Error(), "offset cannot be negative") {
		t.Errorf("Expected error message to contain 'offset cannot be negative', got %s", err.Error())
	}

	query = MediaQuery()
	query.SetOrderBy("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty order_by")
	}
	if !strings.Contains(err.Error(), "order_by cannot be empty") {
		t.Errorf("Expected error message to contain 'order_by cannot be empty', got %s", err.Error())
	}

	query = MediaQuery()
	query.SetStatus("")
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status")
	}
	if !strings.Contains(err.Error(), "status cannot be empty") {
		t.Errorf("Expected error message to contain 'status cannot be empty', got %s", err.Error())
	}

	query = MediaQuery()
	query.SetStatusIn([]string{})
	err = query.Validate()
	if err == nil {
		t.Error("Expected error for empty status_in")
	}
	if !strings.Contains(err.Error(), "status_in cannot be empty array") {
		t.Errorf("Expected error message to contain 'status_in cannot be empty array', got %s", err.Error())
	}
}
