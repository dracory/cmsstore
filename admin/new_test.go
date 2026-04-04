package admin

import (
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

func TestNewAdmin_ValidOptions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	options := AdminOptions{
		Store:           store,
		Logger:          slog.New(slog.NewTextHandler(os.Stderr, nil)),
		MediaManagerURL: "/admin/media",
		PaddingTopPx:    11,
		PaddingRightPx:  12,
		PaddingBottomPx: 13,
		PaddingLeftPx:   14,
	}
	a, err := New(options)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}
	if a == nil {
		t.Errorf("Expected admin to be created, got nil")
	}
	if a.mediaManagerURL != "/admin/media" {
		t.Errorf("Expected mediaManagerURL '/admin/media', got '%s'", a.mediaManagerURL)
	}
	if a.paddingTopPx != 11 {
		t.Errorf("Expected paddingTopPx 11, got %d", a.paddingTopPx)
	}
	if a.paddingRightPx != 12 {
		t.Errorf("Expected paddingRightPx 12, got %d", a.paddingRightPx)
	}
	if a.paddingBottomPx != 13 {
		t.Errorf("Expected paddingBottomPx 13, got %d", a.paddingBottomPx)
	}
	if a.paddingLeftPx != 14 {
		t.Errorf("Expected paddingLeftPx 14, got %d", a.paddingLeftPx)
	}
}

func TestNewAdmin_MissingStore(t *testing.T) {
	options := AdminOptions{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	a, err := New(options)
	if err == nil {
		t.Errorf("Expected error when store is missing")
	}
	if a != nil {
		t.Errorf("Expected nil admin when store is missing")
	}
	if !strings.Contains(err.Error(), shared.ERROR_STORE_IS_NIL) {
		t.Errorf("Expected error to contain '%s', got '%s'", shared.ERROR_STORE_IS_NIL, err.Error())
	}
}

func TestNewAdmin_MissingLogger(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	options := AdminOptions{
		Store: store,
	}
	a, err := New(options)
	if err == nil {
		t.Errorf("Expected error when logger is missing")
	}
	if a != nil {
		t.Errorf("Expected nil admin when logger is missing")
	}
	if !strings.Contains(err.Error(), shared.ERROR_LOGGER_IS_NIL) {
		t.Errorf("Expected error to contain '%s', got '%s'", shared.ERROR_LOGGER_IS_NIL, err.Error())
	}
}
