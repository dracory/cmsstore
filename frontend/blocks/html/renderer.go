package html

import (
	"context"

	"github.com/dracory/cmsstore"
)

// HTMLRenderer provides HTML block rendering functionality
type HTMLRenderer struct{}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{}
}

// Render renders an HTML block by returning its content
func (r *HTMLRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	if block == nil {
		return "<!-- Block is nil -->", nil
	}
	content := block.Content()
	if content == "" {
		return "<!-- Empty block content -->", nil
	}
	return content, nil
}
