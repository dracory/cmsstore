package frontend

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/ui"
)

// TestNew_Basic tests creating a frontend with basic configuration
func TestNew_Basic(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	logger := slog.Default()

	config := Config{
		Store:        store,
		Logger:       logger,
		CacheEnabled: false,
	}

	f := New(config)

	if f == nil {
		t.Error("New should return a non-nil frontend")
	}

	// Verify it implements the interface
	var _ FrontendInterface = f
}

// TestNew_WithCache tests creating a frontend with cache enabled
func TestNew_WithCache(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store:        store,
		CacheEnabled: true,
	}

	f := New(config)

	if f == nil {
		t.Error("New should return a non-nil frontend")
	}

	// Give cache a moment to initialize
	time.Sleep(100 * time.Millisecond)

	fe := f.(*frontend)
	if fe.cache == nil {
		t.Error("Cache should be initialized when CacheEnabled is true")
	}
}

// TestNew_WithoutCache tests creating a frontend with cache disabled
func TestNew_WithoutCache(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store:        store,
		CacheEnabled: false,
	}

	f := New(config)

	if f == nil {
		t.Error("New should return a non-nil frontend")
	}

	fe := f.(*frontend)
	if fe.cache != nil {
		t.Error("Cache should not be initialized when CacheEnabled is false")
	}
}

// TestNew_DefaultCacheExpire tests that default cache expire is set correctly
func TestNew_DefaultCacheExpire(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store:        store,
		CacheEnabled: true,
	}

	f := New(config)
	fe := f.(*frontend)

	if fe.cacheExpireSeconds != 600 {
		t.Errorf("Expected default cache expire to be 600, got %d", fe.cacheExpireSeconds)
	}
}

// TestNew_CustomCacheExpire tests that custom cache expire is set correctly
func TestNew_CustomCacheExpire(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	customExpire := 300
	config := Config{
		Store:              store,
		CacheEnabled:       true,
		CacheExpireSeconds: customExpire,
	}

	f := New(config)
	fe := f.(*frontend)

	if fe.cacheExpireSeconds != customExpire {
		t.Errorf("Expected cache expire to be %d, got %d", customExpire, fe.cacheExpireSeconds)
	}
}

// TestNew_ZeroCacheExpireDefaultsTo600 tests that zero cache expire defaults to 600
func TestNew_ZeroCacheExpireDefaultsTo600(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store:              store,
		CacheEnabled:       true,
		CacheExpireSeconds: 0,
	}

	f := New(config)
	fe := f.(*frontend)

	if fe.cacheExpireSeconds != 600 {
		t.Errorf("Expected cache expire to default to 600, got %d", fe.cacheExpireSeconds)
	}
}

// TestNew_NegativeCacheExpireDefaultsTo600 tests that negative cache expire defaults to 600
func TestNew_NegativeCacheExpireDefaultsTo600(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store:              store,
		CacheEnabled:       true,
		CacheExpireSeconds: -100,
	}

	f := New(config)
	fe := f.(*frontend)

	if fe.cacheExpireSeconds != 600 {
		t.Errorf("Expected cache expire to default to 600, got %d", fe.cacheExpireSeconds)
	}
}

// TestNew_WithBlockEditorRenderer tests creating a frontend with block editor renderer
func TestNew_WithBlockEditorRenderer(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	blockEditorRenderer := func(blocks []ui.BlockInterface) string {
		return "rendered blocks"
	}

	config := Config{
		Store:               store,
		BlockEditorRenderer: blockEditorRenderer,
	}

	f := New(config)
	fe := f.(*frontend)

	if fe.blockEditorRenderer == nil {
		t.Error("BlockEditorRenderer should be set")
	}
}

// TestNew_WithShortcodes tests creating a frontend with shortcodes
func TestNew_WithShortcodes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	shortcodes := []cmsstore.ShortcodeInterface{}

	config := Config{
		Store:      store,
		Shortcodes: shortcodes,
	}

	f := New(config)
	fe := f.(*frontend)

	if fe.shortcodes == nil {
		t.Error("Shortcodes should be set")
	}
}

// TestNew_WithPageNotFoundHandler tests creating a frontend with custom page not found handler
func TestNew_WithPageNotFoundHandler(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handlerCalled := false
	pageNotFoundHandler := func(w http.ResponseWriter, r *http.Request, alias string) (handled bool, result string) {
		handlerCalled = true
		return true, "Custom 404: " + alias
	}

	config := Config{
		Store:               store,
		PageNotFoundHandler: pageNotFoundHandler,
	}

	f := New(config)
	fe := f.(*frontend)

	if fe.pageNotFoundHandler == nil {
		t.Error("PageNotFoundHandler should be set")
	}

	// Test the handler
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	recorder := httptest.NewRecorder()

	handled, result := fe.pageNotFoundHandler(recorder, req, "test-page")

	if !handled {
		t.Error("Handler should return handled=true")
	}

	if !handlerCalled {
		t.Error("Handler should have been called")
	}

	if result != "Custom 404: test-page" {
		t.Errorf("Expected 'Custom 404: test-page', got %q", result)
	}
}

// TestNew_BlockRegistry tests that BlockRegistry is initialized
func TestNew_BlockRegistry(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store: store,
	}

	f := New(config)

	registry := f.(*frontend).BlockRegistry()
	if registry == nil {
		t.Error("BlockRegistry should be initialized")
	}
}

// TestNew_WithoutStore tests that New handles nil store gracefully
func TestNew_WithoutStore(t *testing.T) {
	config := Config{
		Store: nil,
	}

	f := New(config)

	// Should still return a frontend instance
	if f == nil {
		t.Error("New should return a frontend even with nil store")
	}
}

// TestNew_InterfaceImplementation tests that the returned frontend implements all required interfaces
func TestNew_InterfaceImplementation(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store: store,
	}

	f := New(config)

	// Test FrontendInterface
	var fi FrontendInterface = f
	if fi == nil {
		t.Error("Should implement FrontendInterface")
	}

	// Test that Handler and StringHandler can be called
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	// These should not panic
	fi.Handler(recorder, req)
	result := fi.StringHandler(recorder, req)

	// Result might be an error message, but shouldn't panic
	_ = result
}

// TestNew_MenuStoreInterface tests that frontend implements menu.FrontendStore interface
func TestNew_MenuStoreInterface(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store: store,
	}

	f := New(config)
	fe := f.(*frontend)

	// Test menu interface methods
	_, err = fe.MenuFindByID(context.Background(), "test")
	if err != nil {
		// Expected to fail with non-existent ID, but shouldn't panic
	}

	_, err = fe.MenuItemList(context.Background(), nil)
	if err != nil {
		// Expected to fail with nil query, but shouldn't panic
	}

	enabled := fe.MenusEnabled()
	_ = enabled // Should not panic

	_, err = fe.PageFindByID(context.Background(), "test")
	if err != nil {
		// Expected to fail with non-existent ID, but shouldn't panic
	}
}

// TestNew_Logger tests that logger is set correctly
func TestNew_Logger(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	logger := slog.Default()

	config := Config{
		Store:  store,
		Logger: logger,
	}

	f := New(config)
	fe := f.(*frontend)

	if fe.logger != logger {
		t.Error("Logger should be set to the provided logger")
	}
}

// TestNew_LoggerNil tests that nil logger is handled gracefully
func TestNew_LoggerNil(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	config := Config{
		Store:  store,
		Logger: nil,
	}

	f := New(config)
	fe := f.(*frontend)

	// Should accept nil logger without panicking
	if fe.logger != nil {
		t.Error("Logger should be nil when not provided")
	}
}
