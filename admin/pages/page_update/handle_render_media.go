package page_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
)

func renderMediaTab(r *http.Request, page cmsstore.PageInterface) hb.TagInterface {
	htmlContent, err := pageUpdateFiles.ReadFile("page_media.html")
	if err != nil {
		slog.Error("Failed to read page media HTML template", "error", err)
		return hb.Div().HTML("Error loading media component")
	}

	jsContent, err := pageUpdateFiles.ReadFile("page_media.js")
	if err != nil {
		slog.Error("Failed to read page media JavaScript file", "error", err)
		return hb.Div().HTML("Error loading media component")
	}

	initScript := hb.Script(`
		const pageID = '` + page.ID() + `';
		const urlMediaLoad = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionLoadMedia}) + `';
		const urlMediaUpload = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionUploadMedia}) + `';
		const urlMediaSave = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionSaveMedia}) + `';
		const urlMediaDelete = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"action": actionDeleteMedia}) + `';
		const urlMediaAdd = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionAddMedia}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	return hb.Div().
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)
}
