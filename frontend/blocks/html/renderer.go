package html

import (
	"context"

	"github.com/dracory/cmsstore"
)

// BlockRenderer renders HTML blocks by returning their content
type BlockRenderer struct{}

// NewBlockRenderer creates a new HTML block renderer
func NewBlockRenderer() *BlockRenderer {
	return &BlockRenderer{}
}

// Render renders an HTML block by returning its content
func (r *BlockRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	if block == nil {
		return "<!-- Block is nil -->", nil
	}
	content := block.Content()
	if content == "" {
		return "<!-- Empty block content -->", nil
	}
	return content, nil
}
