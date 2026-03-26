package frontend

import (
	"context"

	"github.com/dracory/cmsstore"
)

// HTMLBlockRenderer renders HTML blocks by returning their content
type HTMLBlockRenderer struct{}

// NewHTMLBlockRenderer creates a new HTML block renderer
func NewHTMLBlockRenderer() *HTMLBlockRenderer {
	return &HTMLBlockRenderer{}
}

// Render renders an HTML block by returning its content
func (r *HTMLBlockRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	return block.Content(), nil
}
