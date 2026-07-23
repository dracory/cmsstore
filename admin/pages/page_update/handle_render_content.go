package page_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
)

func renderContentTab(r *http.Request, page cmsstore.PageInterface) hb.TagInterface {
	htmlContent, err := pageUpdateFiles.ReadFile("page_content.html")
	if err != nil {
		slog.Error("Failed to read page content HTML template", "error", err)
		return hb.Div().HTML("Error loading content component")
	}

	jsContent, err := pageUpdateFiles.ReadFile("page_content.js")
	if err != nil {
		slog.Error("Failed to read page content JavaScript file", "error", err)
		return hb.Div().HTML("Error loading content component")
	}

	initScript := hb.Script(`
		const pageID = '` + page.ID() + `';
		const urlContentLoad = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionLoadContent}) + `';
		const urlContentSave = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionSaveContent}) + `';
		const urlBlockEditorHandle = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionBlockeditor}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	return hb.Div().
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)
}
