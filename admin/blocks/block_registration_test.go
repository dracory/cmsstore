package admin

import (
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/blocks/breadcrumbs"
	"github.com/dracory/cmsstore/blocks/navbar"
	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

// TestBlockRegistration tests that navbar and breadcrumbs blocks are properly registered
func TestBlockRegistration(t *testing.T) {
	// Initialize store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Manually register the blocks (simulating what admin UI does)
	cmsstore.RegisterBlockType(navbar.NewNavbarBlockType(store))
	cmsstore.RegisterBlockType(breadcrumbs.NewBreadcrumbsBlockType(store))

	// Check if blocks are registered in the global registry
	navbarBlock := cmsstore.GetBlockType(cmsstore.BLOCK_TYPE_NAVBAR)
	breadcrumbsBlock := cmsstore.GetBlockType(cmsstore.BLOCK_TYPE_BREADCRUMBS)

	// Test navbar block availability
	if navbarBlock == nil {
		t.Error("Navbar block should be registered but was not found")
	} else {
		if navbarBlock.TypeKey() != cmsstore.BLOCK_TYPE_NAVBAR {
			t.Errorf("Expected navbar TypeKey to be %s, got %s", cmsstore.BLOCK_TYPE_NAVBAR, navbarBlock.TypeKey())
		}
		if navbarBlock.TypeLabel() != "Navbar" {
			t.Errorf("Expected navbar TypeLabel to be 'Navbar', got '%s'", navbarBlock.TypeLabel())
		}
	}

	// Test breadcrumbs block availability
	if breadcrumbsBlock == nil {
		t.Error("Breadcrumbs block should be registered but was not found")
	} else {
		if breadcrumbsBlock.TypeKey() != cmsstore.BLOCK_TYPE_BREADCRUMBS {
			t.Errorf("Expected breadcrumbs TypeKey to be %s, got %s", cmsstore.BLOCK_TYPE_BREADCRUMBS, breadcrumbsBlock.TypeKey())
		}
		if breadcrumbsBlock.TypeLabel() != "Breadcrumbs" {
			t.Errorf("Expected breadcrumbs TypeLabel to be 'Breadcrumbs', got '%s'", breadcrumbsBlock.TypeLabel())
		}
	}
}
