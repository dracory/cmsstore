package shared

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
)

type UiInterface interface {
	Endpoint() string
	AdminBreadcrumbs(endpoint string, breadcrumbs []Breadcrumb) hb.TagInterface
	// AdminHeader() hb.TagInterface
	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger() *slog.Logger
	Store() cmsstore.StoreInterface
}
