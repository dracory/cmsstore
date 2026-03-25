package cmsstore

import (
	"testing"
)

func TestNewBlockDefaultsToHTMLType(t *testing.T) {
	block := NewBlock()

	if block.Type() != BLOCK_TYPE_HTML {
		t.Errorf("Expected new block to default to HTML type, got: %s", block.Type())
	}
}

func TestBlockTypeCanBeSet(t *testing.T) {
	block := NewBlock()

	block.SetType(BLOCK_TYPE_MENU)

	if block.Type() != BLOCK_TYPE_MENU {
		t.Errorf("Expected block type to be MENU, got: %s", block.Type())
	}
}

func TestBlockMetaForMenuConfiguration(t *testing.T) {
	block := NewBlock()
	block.SetType(BLOCK_TYPE_MENU)

	// Set menu configuration via metas
	block.SetMeta(BLOCK_META_MENU_ID, "menu-123")
	block.SetMeta(BLOCK_META_MENU_STYLE, BLOCK_MENU_STYLE_HORIZONTAL)
	block.SetMeta(BLOCK_META_MENU_CSS_CLASS, "nav-primary")
	block.SetMeta(BLOCK_META_MENU_START_LEVEL, "0")
	block.SetMeta(BLOCK_META_MENU_MAX_DEPTH, "3")

	// Verify metas are stored correctly
	if block.Meta(BLOCK_META_MENU_ID) != "menu-123" {
		t.Errorf("Expected menu_id to be 'menu-123', got: %s", block.Meta(BLOCK_META_MENU_ID))
	}

	if block.Meta(BLOCK_META_MENU_STYLE) != BLOCK_MENU_STYLE_HORIZONTAL {
		t.Errorf("Expected menu_style to be 'horizontal', got: %s", block.Meta(BLOCK_META_MENU_STYLE))
	}

	if block.Meta(BLOCK_META_MENU_CSS_CLASS) != "nav-primary" {
		t.Errorf("Expected menu_css_class to be 'nav-primary', got: %s", block.Meta(BLOCK_META_MENU_CSS_CLASS))
	}

	if block.Meta(BLOCK_META_MENU_START_LEVEL) != "0" {
		t.Errorf("Expected menu_start_level to be '0', got: %s", block.Meta(BLOCK_META_MENU_START_LEVEL))
	}

	if block.Meta(BLOCK_META_MENU_MAX_DEPTH) != "3" {
		t.Errorf("Expected menu_max_depth to be '3', got: %s", block.Meta(BLOCK_META_MENU_MAX_DEPTH))
	}
}

func TestBlockTypeConstants(t *testing.T) {
	if BLOCK_TYPE_HTML != "html" {
		t.Errorf("Expected BLOCK_TYPE_HTML to be 'html', got: %s", BLOCK_TYPE_HTML)
	}

	if BLOCK_TYPE_MENU != "menu" {
		t.Errorf("Expected BLOCK_TYPE_MENU to be 'menu', got: %s", BLOCK_TYPE_MENU)
	}
}

func TestBlockMenuStyleConstants(t *testing.T) {
	if BLOCK_MENU_STYLE_HORIZONTAL != "horizontal" {
		t.Errorf("Expected BLOCK_MENU_STYLE_HORIZONTAL to be 'horizontal', got: %s", BLOCK_MENU_STYLE_HORIZONTAL)
	}

	if BLOCK_MENU_STYLE_VERTICAL != "vertical" {
		t.Errorf("Expected BLOCK_MENU_STYLE_VERTICAL to be 'vertical', got: %s", BLOCK_MENU_STYLE_VERTICAL)
	}

	if BLOCK_MENU_STYLE_DROPDOWN != "dropdown" {
		t.Errorf("Expected BLOCK_MENU_STYLE_DROPDOWN to be 'dropdown', got: %s", BLOCK_MENU_STYLE_DROPDOWN)
	}

	if BLOCK_MENU_STYLE_BREADCRUMB != "breadcrumb" {
		t.Errorf("Expected BLOCK_MENU_STYLE_BREADCRUMB to be 'breadcrumb', got: %s", BLOCK_MENU_STYLE_BREADCRUMB)
	}
}
