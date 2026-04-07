package cmsstore

import (
	"context"
	"sync"
	"testing"
)

func TestVarsContext_SetAndGet(t *testing.T) {
	vars := NewVarsContext()

	// Test Set and Get
	vars.Set("test_key", "test_value")
	val, ok := vars.Get("test_key")

	if !ok {
		t.Error("Expected key to exist")
	}

	if val != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", val)
	}
}

func TestVarsContext_GetNonExistent(t *testing.T) {
	vars := NewVarsContext()

	val, ok := vars.Get("nonexistent")

	if ok {
		t.Error("Expected key to not exist")
	}

	if val != "" {
		t.Errorf("Expected empty string, got '%s'", val)
	}
}

func TestVarsContext_Overwrite(t *testing.T) {
	vars := NewVarsContext()

	vars.Set("key", "value1")
	vars.Set("key", "value2")

	val, ok := vars.Get("key")

	if !ok {
		t.Error("Expected key to exist")
	}

	if val != "value2" {
		t.Errorf("Expected 'value2', got '%s'", val)
	}
}

func TestVarsContext_All(t *testing.T) {
	vars := NewVarsContext()

	vars.Set("key1", "value1")
	vars.Set("key2", "value2")
	vars.Set("key3", "value3")

	all := vars.All()

	if len(all) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(all))
	}

	if all["key1"] != "value1" {
		t.Errorf("Expected 'value1', got '%s'", all["key1"])
	}

	if all["key2"] != "value2" {
		t.Errorf("Expected 'value2', got '%s'", all["key2"])
	}

	if all["key3"] != "value3" {
		t.Errorf("Expected 'value3', got '%s'", all["key3"])
	}
}

func TestVarsContext_AllIsCopy(t *testing.T) {
	vars := NewVarsContext()

	vars.Set("key1", "value1")

	all := vars.All()
	all["key1"] = "modified"
	all["key2"] = "added"

	// Original should be unchanged
	val, _ := vars.Get("key1")
	if val != "value1" {
		t.Errorf("Expected original to be unchanged, got '%s'", val)
	}

	// New key should not exist in original
	_, ok := vars.Get("key2")
	if ok {
		t.Error("Expected key2 to not exist in original")
	}
}

func TestWithVarsContext(t *testing.T) {
	ctx := context.Background()
	ctx = WithVarsContext(ctx)

	vars := VarsFromContext(ctx)

	if vars == nil {
		t.Error("Expected VarsContext to be in context")
	}
}

func TestVarsFromContext_NoVars(t *testing.T) {
	ctx := context.Background()

	vars := VarsFromContext(ctx)

	if vars != nil {
		t.Error("Expected nil when no VarsContext in context")
	}
}

func TestVarsContext_ContextFlow(t *testing.T) {
	ctx := context.Background()
	ctx = WithVarsContext(ctx)

	// Simulate setting variables in a nested function
	setVars := func(ctx context.Context) {
		if vars := VarsFromContext(ctx); vars != nil {
			vars.Set("nested_key", "nested_value")
		}
	}

	setVars(ctx)

	// Verify variables are accessible in parent context
	vars := VarsFromContext(ctx)
	val, ok := vars.Get("nested_key")

	if !ok {
		t.Error("Expected nested_key to exist")
	}

	if val != "nested_value" {
		t.Errorf("Expected 'nested_value', got '%s'", val)
	}
}

func TestVarsContext_ThreadSafety(t *testing.T) {
	vars := NewVarsContext()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			vars.Set("key", "value")
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			vars.Get("key")
		}()
	}

	// Concurrent All() calls
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			vars.All()
		}()
	}

	wg.Wait()

	// If we get here without race conditions, test passes
}

func TestVarsContext_VariableNamingConventions(t *testing.T) {
	vars := NewVarsContext()

	// Test various naming conventions
	testCases := []struct {
		key   string
		value string
	}{
		{"snake_case", "value1"},
		{"camelCase", "value2"},
		{"PascalCase", "value3"},
		{"blog:title", "value4"},
		{"$prefixed", "value5"},
		{"@userName", "value6"},
		{"#productPrice", "value7"},
		{"kebab-case", "value8"},
		{"UPPER_CASE", "value9"},
		{"number123", "value10"},
	}

	for _, tc := range testCases {
		vars.Set(tc.key, tc.value)
	}

	for _, tc := range testCases {
		val, ok := vars.Get(tc.key)
		if !ok {
			t.Errorf("Expected key '%s' to exist", tc.key)
		}
		if val != tc.value {
			t.Errorf("For key '%s', expected '%s', got '%s'", tc.key, tc.value, val)
		}
	}
}
