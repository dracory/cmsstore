package page_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
)

func renderSEOTab(r *http.Request, page cmsstore.PageInterface) hb.TagInterface {
	htmlContent, err := pageUpdateFiles.ReadFile("page_seo.html")
	if err != nil {
		slog.Error("Failed to read page SEO HTML template", "error", err)
		return hb.Div().HTML("Error loading SEO component")
	}

	jsContent, err := pageUpdateFiles.ReadFile("page_seo.js")
	if err != nil {
		slog.Error("Failed to read page SEO JavaScript file", "error", err)
		return hb.Div().HTML("Error loading SEO component")
	}

	initScript := hb.Script(`
		const pageID = '` + page.ID() + `';
		const urlSEOLoad = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionLoadSEO}) + `';
		const urlSEOSave = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionSaveSEO}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	return hb.Div().
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)
}
