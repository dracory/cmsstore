package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
)

type UiInterface interface {
	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger() *slog.Logger
	Store() cmsstore.StoreInterface
}
