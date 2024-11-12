package shared

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/blockeditor"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
)

type UiConfig struct {
	BlockEditorDefinitions []blockeditor.BlockDefinition
	// AdminHeader            hb.TagInterface
	AdminBreadcrumbs func(endpoint string, breadcrumbs []Breadcrumb) hb.TagInterface
	Endpoint         string
	Layout           func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger *slog.Logger
	Store  cmsstore.StoreInterface
}
