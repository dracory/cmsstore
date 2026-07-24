package page_update

import (
	"embed"
	"net/http"
	"slices"

	"github.com/dracory/api"
	"github.com/dracory/blockeditor"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/req"
)

//go:embed *.html
//go:embed *.js
var pageUpdateFiles embed.FS

// uiInterface extends shared.UiInterface with BlockEditorDefinitions
type uiInterface interface {
	shared.UiInterface
	BlockEditorDefinitions() []blockeditor.BlockDefinition
}

// == CONTROLLER ==============================================================

type pageUpdateController struct {
	ui uiInterface
}

// == CONSTRUCTOR =============================================================

func NewPageUpdateController(ui uiInterface) *pageUpdateController {
	return &pageUpdateController{ui: ui}
}

func (controller *pageUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")
	pageID := req.GetStringTrimmed(r, "page_id")
	view := req.GetStringTrimmedOr(r, "view", viewContent)

	store := controller.ui.Store()
	if store == nil {
		return api.Error("Store not available").ToString()
	}

	// AJAX actions that require POST
	postActions := []string{
		actionSaveContent, actionSaveSEO, actionSaveSettings, actionSaveMiddlewares,
		actionUploadMedia, actionDeleteMedia, actionAddMedia,
	}

	// AJAX actions (any method)
	ajaxActions := []string{
		actionLoadContent, actionSaveContent,
		actionLoadSEO, actionSaveSEO,
		actionLoadSettings, actionSaveSettings,
		actionLoadMiddlewares, actionSaveMiddlewares,
		actionBlockeditor,
		actionLoadMedia, actionUploadMedia, actionSaveMedia, actionDeleteMedia, actionAddMedia,
	}

	if slices.Contains(postActions, action) && r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	if slices.Contains(ajaxActions, action) {
		switch action {
		case actionLoadContent:
			return handleAjaxLoadContent(store, w, r)
		case actionSaveContent:
			return handleAjaxSaveContent(store, w, r)
		case actionLoadSEO:
			return handleAjaxLoadSEO(store, w, r)
		case actionSaveSEO:
			return handleAjaxSaveSEO(store, w, r)
		case actionLoadSettings:
			return handleAjaxLoadSettings(store, w, r)
		case actionSaveSettings:
			return handleAjaxSaveSettings(store, w, r)
		case actionLoadMiddlewares:
			return handleAjaxLoadMiddlewares(store, w, r)
		case actionSaveMiddlewares:
			return handleAjaxSaveMiddlewares(store, w, r)
		case actionBlockeditor:
			return blockeditor.Handle(w, r, controller.ui.BlockEditorDefinitions())
		case actionLoadMedia:
			return handleAjaxLoadMedia(store, w, r)
		case actionUploadMedia:
			return handleAjaxUploadMedia(store, w, r)
		case actionSaveMedia:
			return handleAjaxSaveMedia(store, w, r)
		case actionDeleteMedia:
			return handleAjaxDeleteMedia(store, w, r)
		case actionAddMedia:
			return handleAjaxAddMedia(store, w, r)
		}
	}

	// Page rendering
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	page, err := store.PageFindByID(r.Context(), pageID)
	if err != nil {
		return api.Error("Page not found").ToString()
	}

	if page == nil {
		return api.Error("Page not found").ToString()
	}

	return handleRenderPage(controller.ui, store, page, view, w, r)
}
