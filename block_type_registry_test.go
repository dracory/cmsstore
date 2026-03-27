package cmsstore

import (
	"context"
	"net/http"
	"testing"
)

// TestBlockTypeRegistryOrigin tests that the registry correctly tracks block type origins
type testSystemBlock struct{}

func (t *testSystemBlock) TypeKey() string   { return "test_system" }
func (t *testSystemBlock) TypeLabel() string { return "Test System Block" }
func (t *testSystemBlock) Render(ctx context.Context, block BlockInterface) (string, error) {
	return "", nil
}
func (t *testSystemBlock) GetAdminFields(block BlockInterface, r *http.Request) interface{} {
	return nil
}
func (t *testSystemBlock) SaveAdminFields(r *http.Request, block BlockInterface) error {
	return nil
}

type testCustomBlock struct{}

func (t *testCustomBlock) TypeKey() string   { return "test_custom" }
func (t *testCustomBlock) TypeLabel() string { return "Test Custom Block" }
func (t *testCustomBlock) Render(ctx context.Context, block BlockInterface) (string, error) {
	return "", nil
}
func (t *testCustomBlock) GetAdminFields(block BlockInterface, r *http.Request) interface{} {
	return nil
}
func (t *testCustomBlock) SaveAdminFields(r *http.Request, block BlockInterface) error {
	return nil
}

func TestBlockTypeRegistryOrigin(t *testing.T) {
	// Clear registry for clean test
	globalBlockTypeRegistry = &BlockTypeRegistry{
		types:   make(map[string]BlockType),
		origins: make(map[string]string),
	}

	systemBlock := &testSystemBlock{}
	customBlock := &testCustomBlock{}

	// Register system block
	RegisterSystemBlockType(systemBlock)

	// Register custom block
	RegisterCustomBlockType(customBlock)

	// Test GetBlockTypeOrigin
	systemOrigin := GetBlockTypeOrigin("test_system")
	if systemOrigin != BLOCK_ORIGIN_SYSTEM {
		t.Errorf("Expected system origin '%s', got '%s'", BLOCK_ORIGIN_SYSTEM, systemOrigin)
	}

	customOrigin := GetBlockTypeOrigin("test_custom")
	if customOrigin != BLOCK_ORIGIN_CUSTOM {
		t.Errorf("Expected custom origin '%s', got '%s'", BLOCK_ORIGIN_CUSTOM, customOrigin)
	}

	// Test unknown block type returns custom as default
	unknownOrigin := GetBlockTypeOrigin("unknown")
	if unknownOrigin != BLOCK_ORIGIN_CUSTOM {
		t.Errorf("Expected default origin '%s' for unknown type, got '%s'", BLOCK_ORIGIN_CUSTOM, unknownOrigin)
	}

	// Test GetSystemBlockTypes
	systemTypes := GetSystemBlockTypes()
	if len(systemTypes) != 1 {
		t.Errorf("Expected 1 system type, got %d", len(systemTypes))
	}
	if _, ok := systemTypes["test_system"]; !ok {
		t.Error("Expected test_system in system types")
	}

	// Test GetCustomBlockTypes
	customTypes := GetCustomBlockTypes()
	if len(customTypes) != 1 {
		t.Errorf("Expected 1 custom type, got %d", len(customTypes))
	}
	if _, ok := customTypes["test_custom"]; !ok {
		t.Error("Expected test_custom in custom types")
	}
}

func TestBlockTypeRegistryGetAll(t *testing.T) {
	// Clear registry for clean test
	globalBlockTypeRegistry = &BlockTypeRegistry{
		types:   make(map[string]BlockType),
		origins: make(map[string]string),
	}

	RegisterSystemBlockType(&testSystemBlock{})
	RegisterCustomBlockType(&testCustomBlock{})

	allTypes := GetAllBlockTypes()
	if len(allTypes) != 2 {
		t.Errorf("Expected 2 total types, got %d", len(allTypes))
	}
}

func TestBlockTypeRegistryGet(t *testing.T) {
	// Clear registry for clean test
	globalBlockTypeRegistry = &BlockTypeRegistry{
		types:   make(map[string]BlockType),
		origins: make(map[string]string),
	}

	RegisterSystemBlockType(&testSystemBlock{})

	blockType := GetBlockType("test_system")
	if blockType == nil {
		t.Error("Expected to find test_system block type")
	}
	if blockType.TypeKey() != "test_system" {
		t.Errorf("Expected TypeKey 'test_system', got '%s'", blockType.TypeKey())
	}
	if blockType.TypeLabel() != "Test System Block" {
		t.Errorf("Expected TypeLabel 'Test System Block', got '%s'", blockType.TypeLabel())
	}

	// Test non-existent type
	nonExistent := GetBlockType("non_existent")
	if nonExistent != nil {
		t.Error("Expected nil for non-existent block type")
	}
}
