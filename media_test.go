package cmsstore

import (
	"encoding/json"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNewMediaDefaults(t *testing.T) {
	media := NewMedia()

	if media.ID() == "" {
		t.Error("expected non-empty ID")
	}
	if media.CreatedAt() == "" {
		t.Error("expected non-empty CreatedAt")
	}
	if media.UpdatedAt() == "" {
		t.Error("expected non-empty UpdatedAt")
	}
	if media.Status() != MEDIA_STATUS_DRAFT {
		t.Errorf("expected status %q, got %q", MEDIA_STATUS_DRAFT, media.Status())
	}
	if media.SoftDeletedAt() != MAX_DATETIME {
		t.Errorf("expected SoftDeletedAt %q, got %q", MAX_DATETIME, media.SoftDeletedAt())
	}
	if media.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false")
	}

	metas, err := media.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Error("expected empty metas")
	}

	createdCarbon := media.CreatedAtCarbon()
	if createdCarbon == nil {
		t.Fatal("expected non-nil CreatedAtCarbon")
	}
	if media.CreatedAt() != createdCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected CreatedAt to match CreatedAtCarbon")
	}

	updatedCarbon := media.UpdatedAtCarbon()
	if updatedCarbon == nil {
		t.Fatal("expected non-nil UpdatedAtCarbon")
	}
	if media.UpdatedAt() != updatedCarbon.ToDateTimeString(carbon.UTC) {
		t.Error("expected UpdatedAt to match UpdatedAtCarbon")
	}

	softDeletedCarbon := media.SoftDeletedAtCarbon()
	if softDeletedCarbon == nil {
		t.Fatal("expected non-nil SoftDeletedAtCarbon")
	}
	if !softDeletedCarbon.Gte(carbon.Now(carbon.UTC)) {
		t.Error("expected SoftDeletedAtCarbon to be in the future")
	}
}

func TestMediaGetterMethods(t *testing.T) {
	media := NewMedia()

	if media.EntityID() != "" {
		t.Error("expected empty EntityID")
	}
	if media.EntityType() != "" {
		t.Error("expected empty EntityType")
	}
	if media.Title() != "" {
		t.Error("expected empty Title")
	}
	if media.Description() != "" {
		t.Error("expected empty Description")
	}
	if media.Memo() != "" {
		t.Error("expected empty Memo")
	}
	if media.URL() != "" {
		t.Error("expected empty URL")
	}
	if media.Type() != "" {
		t.Error("expected empty Type")
	}
	if media.Size() != "0" {
		t.Errorf("expected Size %q, got %q", "0", media.Size())
	}
	if media.Extension() != "" {
		t.Error("expected empty Extension")
	}
	if media.Sequence() != "0" {
		t.Errorf("expected Sequence %q, got %q", "0", media.Sequence())
	}
	if media.SequenceInt() != 0 {
		t.Errorf("expected SequenceInt 0, got %d", media.SequenceInt())
	}
	if media.Handle() != "" {
		t.Error("expected empty Handle")
	}
	if media.SiteID() != "" {
		t.Error("expected empty SiteID")
	}
}

func TestMediaStatusMethods(t *testing.T) {
	media := NewMedia()

	if media.IsActive() {
		t.Error("expected IsActive to be false for DRAFT")
	}
	if media.IsInactive() {
		t.Error("expected IsInactive to be false for DRAFT")
	}
	if !media.IsDraft() {
		t.Error("expected IsDraft to be true for DRAFT")
	}

	media.SetStatus(MEDIA_STATUS_ACTIVE)
	if !media.IsActive() {
		t.Error("expected IsActive to be true for ACTIVE")
	}
	if media.IsInactive() {
		t.Error("expected IsInactive to be false for ACTIVE")
	}
	if media.IsDraft() {
		t.Error("expected IsDraft to be false for ACTIVE")
	}

	media.SetStatus(MEDIA_STATUS_INACTIVE)
	if media.IsActive() {
		t.Error("expected IsActive to be false for INACTIVE")
	}
	if !media.IsInactive() {
		t.Error("expected IsInactive to be true for INACTIVE")
	}
	if media.IsDraft() {
		t.Error("expected IsDraft to be false for INACTIVE")
	}
}

func TestMediaSoftDeleteMethods(t *testing.T) {
	media := NewMedia()
	if media.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false by default")
	}

	future := carbon.Now(carbon.UTC).AddHour()
	media.SetSoftDeletedAt(future.ToDateTimeString(carbon.UTC))
	if media.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be false with future date")
	}

	past := carbon.Now(carbon.UTC).SubHour()
	media.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))
	if !media.IsSoftDeleted() {
		t.Error("expected IsSoftDeleted to be true with past date")
	}
}

func TestMediaTypeMethods(t *testing.T) {
	media := NewMedia()

	media.SetType("image/jpeg")
	if !media.IsImage() {
		t.Error("expected IsImage to be true for image/jpeg")
	}
	if media.IsVideo() {
		t.Error("expected IsVideo to be false for image/jpeg")
	}

	media.SetType("video/mp4")
	if media.IsImage() {
		t.Error("expected IsImage to be false for video/mp4")
	}
	if !media.IsVideo() {
		t.Error("expected IsVideo to be true for video/mp4")
	}

	media.SetType("application/pdf")
	if media.IsImage() {
		t.Error("expected IsImage to be false for application/pdf")
	}
	if media.IsVideo() {
		t.Error("expected IsVideo to be false for application/pdf")
	}
}

func TestMediaMetasMethods(t *testing.T) {
	media := NewMedia()

	metas, err := media.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Error("expected empty metas")
	}

	if media.Meta("nonexistent") != "" {
		t.Error("expected empty Meta for nonexistent key")
	}

	err = media.SetMetas(map[string]string{"alt": "test alt", "caption": "test caption"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	metas, err = media.Metas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metas["alt"] != "test alt" {
		t.Errorf("expected alt %q, got %q", "test alt", metas["alt"])
	}
	if metas["caption"] != "test caption" {
		t.Errorf("expected caption %q, got %q", "test caption", metas["caption"])
	}

	if media.Meta("alt") != "test alt" {
		t.Errorf("expected alt %q", "test alt")
	}

	err = media.SetMeta("newkey", "newvalue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if media.Meta("newkey") != "newvalue" {
		t.Errorf("expected newkey %q", "newvalue")
	}

	err = media.UpsertMetas(map[string]string{"alt": "updated alt", "extra": "extra value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if media.Meta("alt") != "updated alt" {
		t.Errorf("expected updated alt %q", "updated alt")
	}
	if media.Meta("caption") != "test caption" {
		t.Errorf("expected preserved caption %q", "test caption")
	}
	if media.Meta("extra") != "extra value" {
		t.Errorf("expected added extra %q", "extra value")
	}
}

func TestMediaSequenceMethods(t *testing.T) {
	media := NewMedia()

	media.SetSequenceInt(42)
	if media.SequenceInt() != 42 {
		t.Errorf("expected SequenceInt 42, got %d", media.SequenceInt())
	}
	if media.Sequence() != "42" {
		t.Errorf("expected Sequence %q, got %q", "42", media.Sequence())
	}

	media.SetSequence("7")
	if media.SequenceInt() != 7 {
		t.Errorf("expected SequenceInt 7, got %d", media.SequenceInt())
	}
	if media.Sequence() != "7" {
		t.Errorf("expected Sequence %q, got %q", "7", media.Sequence())
	}
}

func TestMediaMarshalToVersioning(t *testing.T) {
	media := NewMedia()
	media.SetEntityID("entity-123")
	media.SetEntityType("page")
	media.SetTitle("Test Media")
	media.SetDescription("Test Description")
	media.SetURL("https://example.com/image.jpg")
	media.SetType("image/jpeg")
	media.SetSize("1024")
	media.SetExtension("jpg")
	media.SetStatus(MEDIA_STATUS_ACTIVE)
	media.SetHandle("test-media-handle")
	media.SetSiteID("site-1")

	versionedJSON, err := media.MarshalToVersioning()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if versionedJSON == "" {
		t.Error("expected non-empty versionedJSON")
	}

	var versionedData map[string]string
	err = json.Unmarshal([]byte(versionedJSON), &versionedData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if versionedData[COLUMN_ENTITY_ID] != "entity-123" {
		t.Errorf("expected entity_id %q, got %q", "entity-123", versionedData[COLUMN_ENTITY_ID])
	}
	if versionedData[COLUMN_ENTITY_TYPE] != "page" {
		t.Errorf("expected entity_type %q, got %q", "page", versionedData[COLUMN_ENTITY_TYPE])
	}
	if versionedData[COLUMN_TITLE] != "Test Media" {
		t.Errorf("expected title %q, got %q", "Test Media", versionedData[COLUMN_TITLE])
	}
	if versionedData[COLUMN_MEDIA_URL] != "https://example.com/image.jpg" {
		t.Errorf("expected media_url %q, got %q", "https://example.com/image.jpg", versionedData[COLUMN_MEDIA_URL])
	}
	if versionedData[COLUMN_MEDIA_TYPE] != "image/jpeg" {
		t.Errorf("expected media_type %q, got %q", "image/jpeg", versionedData[COLUMN_MEDIA_TYPE])
	}
	if versionedData[COLUMN_STATUS] != MEDIA_STATUS_ACTIVE {
		t.Errorf("expected status %q, got %q", MEDIA_STATUS_ACTIVE, versionedData[COLUMN_STATUS])
	}

	_, hasCreatedAt := versionedData[COLUMN_CREATED_AT]
	_, hasUpdatedAt := versionedData[COLUMN_UPDATED_AT]
	_, hasSoftDeletedAt := versionedData[COLUMN_SOFT_DELETED_AT]
	if hasCreatedAt {
		t.Error("expected CreatedAt to be excluded")
	}
	if hasUpdatedAt {
		t.Error("expected UpdatedAt to be excluded")
	}
	if hasSoftDeletedAt {
		t.Error("expected SoftDeletedAt to be excluded")
	}
}

func TestMediaFromExistingData(t *testing.T) {
	data := map[string]string{
		"id":              "test-id-123",
		"entity_id":       "entity-456",
		"entity_type":     "block",
		"title":           "Existing Media",
		"media_url":       "https://example.com/file.png",
		"media_type":      "image/png",
		"file_size":       "2048",
		"file_extension":  "png",
		"status":          MEDIA_STATUS_ACTIVE,
		"handle":          "existing-handle",
		"site_id":         "site-2",
		"created_at":      "2023-01-01 10:00:00",
		"updated_at":      "2023-01-02 10:00:00",
		"soft_deleted_at": MAX_DATETIME,
	}

	media := NewMediaFromExistingData(data)

	if media.ID() != "test-id-123" {
		t.Errorf("expected ID %q, got %q", "test-id-123", media.ID())
	}
	if media.EntityID() != "entity-456" {
		t.Errorf("expected EntityID %q, got %q", "entity-456", media.EntityID())
	}
	if media.EntityType() != "block" {
		t.Errorf("expected EntityType %q, got %q", "block", media.EntityType())
	}
	if media.Title() != "Existing Media" {
		t.Errorf("expected Title %q, got %q", "Existing Media", media.Title())
	}
	if media.URL() != "https://example.com/file.png" {
		t.Errorf("expected URL %q, got %q", "https://example.com/file.png", media.URL())
	}
	if media.Type() != "image/png" {
		t.Errorf("expected Type %q, got %q", "image/png", media.Type())
	}
	if media.Size() != "2048" {
		t.Errorf("expected Size %q, got %q", "2048", media.Size())
	}
	if media.Extension() != "png" {
		t.Errorf("expected Extension %q, got %q", "png", media.Extension())
	}
	if media.Status() != MEDIA_STATUS_ACTIVE {
		t.Errorf("expected Status %q, got %q", MEDIA_STATUS_ACTIVE, media.Status())
	}
	if media.Handle() != "existing-handle" {
		t.Errorf("expected Handle %q, got %q", "existing-handle", media.Handle())
	}
	if media.SiteID() != "site-2" {
		t.Errorf("expected SiteID %q, got %q", "site-2", media.SiteID())
	}
	if !media.IsActive() {
		t.Error("expected IsActive to be true")
	}
	if !media.IsImage() {
		t.Error("expected IsImage to be true for image/png")
	}
}

func TestMediaServeURL(t *testing.T) {
	media := NewMedia()
	media.SetID("abc123")
	media.SetExtension(".png")

	got := media.ServeURL()
	expected := "/cms/media/abc123.png"
	if got != expected {
		t.Errorf("expected ServeURL %q, got %q", expected, got)
	}
}

func TestMediaServeURL_NoExtension(t *testing.T) {
	media := NewMedia()
	media.SetID("abc123")
	media.SetExtension("")

	got := media.ServeURL()
	expected := "/cms/media/abc123"
	if got != expected {
		t.Errorf("expected ServeURL %q, got %q", expected, got)
	}
}

func TestMediaServeURL_ExtensionWithoutDot(t *testing.T) {
	media := NewMedia()
	media.SetID("abc123")
	media.SetExtension("jpg")

	got := media.ServeURL()
	expected := "/cms/media/abc123.jpg"
	if got != expected {
		t.Errorf("expected ServeURL %q, got %q", expected, got)
	}
}
