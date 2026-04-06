package cmsstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewTranslationDefaults(t *testing.T) {
	translation := NewTranslation()

	// Test default values
	if len(translation.ID()) == 0 {
		t.Error("ID should be generated")
	}
	if len(translation.CreatedAt()) == 0 {
		t.Error("CreatedAt should be set")
	}
	if len(translation.UpdatedAt()) == 0 {
		t.Error("UpdatedAt should be set")
	}
	if translation.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("SoftDeletedAt should default to max datetime, got %s", translation.SoftDeletedAt())
	}
	if translation.Status() != TEMPLATE_STATUS_DRAFT {
		t.Errorf("Status should default to draft, got %s", translation.Status())
	}
	if translation.IsSoftDeleted() {
		t.Error("New translation should not be marked as soft deleted")
	}

	content, err := translation.Content()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(content) != 0 {
		t.Errorf("Content should default to an empty map, got %v", content)
	}

	metas, err := translation.Metas()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("Metas should default to an empty map, got %v", metas)
	}

	createdCarbon := translation.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Error("CreatedAtCarbon should not be nil")
	}
	if createdCarbon.IsZero() {
		t.Error("CreatedAtCarbon should be parseable")
	}

	updatedCarbon := translation.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Error("UpdatedAtCarbon should not be nil")
	}
	if updatedCarbon.IsZero() {
		t.Error("UpdatedAtCarbon should be parseable")
	}

	softDeletedCarbon := translation.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Error("SoftDeletedAtCarbon should not be nil")
	}
	if !softDeletedCarbon.Gte(carbon.Now(carbon.UTC)) {
		t.Error("SoftDeletedAt should be in the future by default")
	}
}

func TestTranslationContentRoundTrip(t *testing.T) {
	translation := NewTranslation()

	expectedContent := map[string]string{
		"en": "Hello",
		"fr": "Bonjour",
	}

	err := translation.SetContent(expectedContent)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	content, err := translation.Content()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if content["en"] != expectedContent["en"] || content["fr"] != expectedContent["fr"] {
		t.Errorf("Expected content %v, got %v", expectedContent, content)
	}
}

func TestTranslationMetasUpsertAndMetaLookup(t *testing.T) {
	translation := NewTranslation()

	err := translation.SetMetas(map[string]string{"locale": "en"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if translation.Meta("locale") != "en" {
		t.Errorf("Expected locale 'en', got %s", translation.Meta("locale"))
	}

	err = translation.UpsertMetas(map[string]string{"locale": "fr", "category": "general"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if translation.Meta("locale") != "fr" {
		t.Errorf("Expected locale 'fr', got %s", translation.Meta("locale"))
	}
	if translation.Meta("category") != "general" {
		t.Errorf("Expected category 'general', got %s", translation.Meta("category"))
	}

	if translation.Meta("missing") != "" {
		t.Errorf("Expected empty Meta for missing key, got %s", translation.Meta("missing"))
	}
}

func TestTranslationSoftDeleteBehaviour(t *testing.T) {
	translation := NewTranslation()
	if translation.IsSoftDeleted() {
		t.Error("New translation should not be soft deleted")
	}

	past := carbon.Now(carbon.UTC).SubHour()
	translation.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))

	if !translation.IsSoftDeleted() {
		t.Error("Translation should be marked as soft deleted when past timestamp is set")
	}
	if translation.SoftDeletedAt() != past.ToDateTimeString(carbon.UTC) {
		t.Errorf("Expected SoftDeletedAt %s, got %s", past.ToDateTimeString(carbon.UTC), translation.SoftDeletedAt())
	}
}

func TestNewTranslationFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:      "test-id",
		COLUMN_NAME:    "test-name",
		COLUMN_STATUS:  PAGE_STATUS_ACTIVE,
		COLUMN_CONTENT: "{\"en\":\"Hello\"}",
	}

	translation := NewTranslationFromExistingData(data)

	if translation.ID() != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", translation.ID())
	}
	if translation.Name() != "test-name" {
		t.Errorf("Expected Name 'test-name', got %s", translation.Name())
	}
	if translation.Status() != PAGE_STATUS_ACTIVE {
		t.Errorf("Expected Status %s, got %s", PAGE_STATUS_ACTIVE, translation.Status())
	}

	content, err := translation.Content()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if content["en"] != "Hello" {
		t.Errorf("Expected content['en'] 'Hello', got %s", content["en"])
	}
}

func TestTranslationStatusChecks(t *testing.T) {
	translation := NewTranslation()

	translation.SetStatus(PAGE_STATUS_ACTIVE)
	if !translation.IsActive() {
		t.Error("Expected IsActive to be true")
	}
	if translation.IsInactive() {
		t.Error("Expected IsInactive to be false")
	}

	translation.SetStatus(PAGE_STATUS_INACTIVE)
	if translation.IsActive() {
		t.Error("Expected IsActive to be false")
	}
	if !translation.IsInactive() {
		t.Error("Expected IsInactive to be true")
	}
}

func TestTranslationSetMeta(t *testing.T) {
	translation := NewTranslation()

	err := translation.SetMeta("key1", "value1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if translation.Meta("key1") != "value1" {
		t.Errorf("Expected Meta('key1') 'value1', got %s", translation.Meta("key1"))
	}

	err = translation.SetMeta("key1", "value2")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if translation.Meta("key1") != "value2" {
		t.Errorf("Expected Meta('key1') 'value2', got %s", translation.Meta("key1"))
	}
}

func TestTranslationSettersGetters(t *testing.T) {
	translation := NewTranslation()

	translation.SetHandle("test-handle")
	if translation.Handle() != "test-handle" {
		t.Errorf("Expected Handle 'test-handle', got %s", translation.Handle())
	}

	translation.SetMemo("test-memo")
	if translation.Memo() != "test-memo" {
		t.Errorf("Expected Memo 'test-memo', got %s", translation.Memo())
	}

	translation.SetName("test-name")
	if translation.Name() != "test-name" {
		t.Errorf("Expected Name 'test-name', got %s", translation.Name())
	}

	translation.SetSiteID("test-site")
	if translation.SiteID() != "test-site" {
		t.Errorf("Expected SiteID 'test-site', got %s", translation.SiteID())
	}

	translation.SetUpdatedAt("2023-01-01 12:00:00")
	if translation.UpdatedAt() != "2023-01-01 12:00:00" {
		t.Errorf("Expected UpdatedAt '2023-01-01 12:00:00', got %s", translation.UpdatedAt())
	}
	if translation.UpdatedAtCarbon() == nil {
		t.Error("Expected UpdatedAtCarbon to be non-nil")
	}
}
