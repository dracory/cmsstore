package customentity

import (
	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
)

// UiInterface defines the interface for custom entity admin UI
type UiInterface interface {
	Store() cmsstore.StoreInterface
	Layout() hb.TagInterface
	Logger() any
}
