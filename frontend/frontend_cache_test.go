package frontend

import (
	"context"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
)

// TestFrontendCache_CacheDisabled tests cache behavior when cache is disabled
func TestFrontendCache_CacheDisabled(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: false,
	})

	fe := f.(*frontend)

	// CacheHas should return false when cache is disabled
	if fe.CacheHas("test_key") {
		t.Error("CacheHas should return false when cache is disabled")
	}

	// CacheGet should return nil when cache is disabled
	if fe.CacheGet("test_key") != nil {
		t.Error("CacheGet should return nil when cache is disabled")
	}

	// CacheSet should not panic when cache is disabled
	fe.CacheSet("test_key", "test_value", 60)
}

// TestFrontendCache_CacheEnabled tests cache behavior when cache is enabled
func TestFrontendCache_CacheEnabled(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: true,
	})

	fe := f.(*frontend)

	// Give cache a moment to initialize
	time.Sleep(100 * time.Millisecond)

	testKey := "test_key"
	testValue := "test_value"

	// Initially, cache should not have the key
	if fe.CacheHas(testKey) {
		t.Error("Cache should not have the key initially")
	}

	// Set a value
	fe.CacheSet(testKey, testValue, 60)

	// Now cache should have the key
	if !fe.CacheHas(testKey) {
		t.Error("Cache should have the key after setting")
	}

	// Get the value
	retrievedValue := fe.CacheGet(testKey)
	if retrievedValue == nil {
		t.Error("CacheGet should return the value")
	}

	if retrievedValue.(string) != testValue {
		t.Errorf("Expected %q, got %q", testValue, retrievedValue.(string))
	}
}

// TestFrontendCache_CacheGetNil tests CacheGet with nil cache
func TestFrontendCache_CacheGetNil(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: true,
	})

	fe := f.(*frontend)

	// Manually set cache to nil to test the nil check
	fe.cache = nil

	// CacheGet should return nil when cache is nil
	if fe.CacheGet("test_key") != nil {
		t.Error("CacheGet should return nil when cache is nil")
	}
}

// TestFrontendCache_CacheHasNil tests CacheHas with nil cache
func TestFrontendCache_CacheHasNil(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: true,
	})

	fe := f.(*frontend)

	// Manually set cache to nil to test the nil check
	fe.cache = nil

	// CacheHas should return false when cache is nil
	if fe.CacheHas("test_key") {
		t.Error("CacheHas should return false when cache is nil")
	}
}

// TestFrontendCache_CacheSetNil tests CacheSet with nil cache
func TestFrontendCache_CacheSetNil(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: true,
	})

	fe := f.(*frontend)

	// Manually set cache to nil to test the nil check
	fe.cache = nil

	// CacheSet should not panic when cache is nil
	fe.CacheSet("test_key", "test_value", 60)
}

// TestFrontendCache_DifferentTypes tests caching different types of values
func TestFrontendCache_DifferentTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: true,
	})

	fe := f.(*frontend)

	// Give cache a moment to initialize
	time.Sleep(100 * time.Millisecond)

	tests := []struct {
		key   string
		value any
	}{
		{"string_key", "string_value"},
		{"int_key", 123},
		{"float_key", 45.67},
		{"bool_key", true},
		{"slice_key", []string{"a", "b", "c"}},
		{"map_key", map[string]string{"a": "1", "b": "2"}},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			fe.CacheSet(tt.key, tt.value, 60)

			if !fe.CacheHas(tt.key) {
				t.Errorf("Cache should have key %q", tt.key)
			}

			retrieved := fe.CacheGet(tt.key)
			if retrieved == nil {
				t.Errorf("CacheGet should return value for key %q", tt.key)
			}
		})
	}
}

// TestFrontendCache_Overwrite tests overwriting a cached value
func TestFrontendCache_Overwrite(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: true,
	})

	fe := f.(*frontend)

	// Give cache a moment to initialize
	time.Sleep(100 * time.Millisecond)

	testKey := "test_key"
	initialValue := "initial_value"
	updatedValue := "updated_value"

	fe.CacheSet(testKey, initialValue, 60)

	retrieved := fe.CacheGet(testKey)
	if retrieved.(string) != initialValue {
		t.Errorf("Expected %q, got %q", initialValue, retrieved.(string))
	}

	// Overwrite the value
	fe.CacheSet(testKey, updatedValue, 60)

	retrieved = fe.CacheGet(testKey)
	if retrieved.(string) != updatedValue {
		t.Errorf("Expected %q, got %q", updatedValue, retrieved.(string))
	}
}

// TestFrontendCache_DefaultExpire tests that default cache expire seconds is set correctly
func TestFrontendCache_DefaultExpire(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: true,
	})

	fe := f.(*frontend)

	// Default should be 10 minutes (600 seconds)
	if fe.cacheExpireSeconds != 600 {
		t.Errorf("Expected default cache expire seconds to be 600, got %d", fe.cacheExpireSeconds)
	}
}

// TestFrontendCache_CustomExpire tests that custom cache expire seconds is set correctly
func TestFrontendCache_CustomExpire(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	customExpire := 300
	f := New(Config{
		Store:               store,
		CacheEnabled:        true,
		CacheExpireSeconds:  customExpire,
	})

	fe := f.(*frontend)

	if fe.cacheExpireSeconds != customExpire {
		t.Errorf("Expected cache expire seconds to be %d, got %d", customExpire, fe.cacheExpireSeconds)
	}
}

// TestFrontendCache_ZeroExpireDefaultsTo600 tests that zero expire seconds defaults to 600
func TestFrontendCache_ZeroExpireDefaultsTo600(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:               store,
		CacheEnabled:        true,
		CacheExpireSeconds:  0,
	})

	fe := f.(*frontend)

	// Should default to 600 seconds
	if fe.cacheExpireSeconds != 600 {
		t.Errorf("Expected cache expire seconds to default to 600, got %d", fe.cacheExpireSeconds)
	}
}

// TestFrontendCache_NegativeExpireDefaultsTo600 tests that negative expire seconds defaults to 600
func TestFrontendCache_NegativeExpireDefaultsTo600(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:               store,
		CacheEnabled:        true,
		CacheExpireSeconds:  -100,
	})

	fe := f.(*frontend)

	// Should default to 600 seconds
	if fe.cacheExpireSeconds != 600 {
		t.Errorf("Expected cache expire seconds to default to 600, got %d", fe.cacheExpireSeconds)
	}
}

// TestFrontendCache_WarmUpCache tests the warmUpCache function
func TestFrontendCache_WarmUpCache(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Create a site to test cache warming
	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:        store,
		CacheEnabled: true,
	})

	fe := f.(*frontend)

	// Give cache a moment to initialize and warm up
	time.Sleep(200 * time.Millisecond)

	// The warmUpCache function runs in a goroutine and calls fetchActiveSites
	// We can't easily test the exact behavior, but we can verify it doesn't panic
	// and that the cache is functional after initialization
	
	// Test that cache is working after warm-up
	testKey := "test_after_warmup"
	fe.CacheSet(testKey, "test_value", 60)
	
	if !fe.CacheHas(testKey) {
		t.Error("Cache should be functional after warm-up")
	}
}
