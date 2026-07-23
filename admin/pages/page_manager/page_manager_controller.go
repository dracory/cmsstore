package page_manager

import (
	"embed"
	"net/http"
	"slices"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/req"
)

//go:embed *.html
//go:embed *.js
var pageFiles embed.FS

const (
	actionLoadPages  = "load-pages"
	actionDeletePage = "delete-page"
	actionCreatePage = "create-page"
)

// == CONTROLLER ==============================================================

type pageManagerController struct {
	ui shared.UiInterface
}

// == CONSTRUCTOR =============================================================

func NewPageManagerController(ui shared.UiInterface) *pageManagerController {
	return &pageManagerController{ui: ui}
}

func (controller *pageManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	ajaxActions := []string{actionLoadPages, actionCreatePage, actionDeletePage}

	if slices.Contains(ajaxActions, action) && r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	store := controller.ui.Store()
	if store == nil {
		return api.Error("Store not available").ToString()
	}

	switch action {
	case actionLoadPages:
		return handleAjaxLoadPages(store, w, r)
	case actionCreatePage:
		return handleAjaxCreatePage(store, w, r)
	case actionDeletePage:
		return handleAjaxDeletePage(store, w, r)
	default:
		return handleRenderPage(controller.ui, store, w, r)
	}
}
