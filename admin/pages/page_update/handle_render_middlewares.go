package page_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
)

func renderMiddlewaresTab(r *http.Request, page cmsstore.PageInterface) hb.TagInterface {
	htmlContent, err := pageUpdateFiles.ReadFile("page_middlewares.html")
	if err != nil {
		slog.Error("Failed to read page middlewares HTML template", "error", err)
		return hb.Div().HTML("Error loading middlewares component")
	}

	jsContent, err := pageUpdateFiles.ReadFile("page_middlewares.js")
	if err != nil {
		slog.Error("Failed to read page middlewares JavaScript file", "error", err)
		return hb.Div().HTML("Error loading middlewares component")
	}

	initScript := hb.Script(`
		const pageID = '` + page.ID() + `';
		const urlMiddlewaresLoad = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionLoadMiddlewares}) + `';
		const urlMiddlewaresSave = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionSaveMiddlewares}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	return hb.Div().
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)
}
