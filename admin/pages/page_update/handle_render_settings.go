package page_update

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
)

func renderSettingsTab(r *http.Request, page cmsstore.PageInterface) hb.TagInterface {
	htmlContent, err := pageUpdateFiles.ReadFile("page_settings.html")
	if err != nil {
		slog.Error("Failed to read page settings HTML template", "error", err)
		return hb.Div().HTML("Error loading settings component")
	}

	jsContent, err := pageUpdateFiles.ReadFile("page_settings.js")
	if err != nil {
		slog.Error("Failed to read page settings JavaScript file", "error", err)
		return hb.Div().HTML("Error loading settings component")
	}

	initScript := hb.Script(`
		const pageID = '` + page.ID() + `';
		const urlSettingsLoad = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionLoadSettings}) + `';
		const urlSettingsSave = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID(), "action": actionSaveSettings}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	return hb.Div().
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)
}
