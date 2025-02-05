package admin

import (
	"net/http"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/router"
	"github.com/gouniverse/utils"
)

// == CONTROLLER ==============================================================

type translationDeleteController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*translationDeleteController)(nil)

// == CONSTRUCTOR =============================================================

type translationDeleteControllerData struct {
	request        *http.Request
	translationID  string
	translation    cmsstore.TranslationInterface
	successMessage string
}

func NewTranslationDeleteController(ui UiInterface) *translationDeleteController {
	return &translationDeleteController{
		ui: ui,
	}
}

func (controller translationDeleteController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		return hb.Swal(hb.SwalOptions{
			Icon: "error",
			Text: errorMessage,
		}).ToHTML()
	}

	if data.successMessage != "" {
		return hb.Wrap().
			Child(hb.Swal(hb.SwalOptions{
				Icon: "success",
				Text: data.successMessage,
			})).
			Child(hb.Script("setTimeout(() => {window.location.href = window.location.href}, 2000)")).
			ToHTML()
	}

	return controller.
		modal(data).
		ToHTML()
}

func (controller *translationDeleteController) modal(data translationDeleteControllerData) hb.TagInterface {
	submitUrl := shared.URLR(data.request, shared.PathTranslationsTranslationDelete, map[string]string{
		"translation_id": data.translationID,
	})

	modalID := "ModalTranslationDelete"
	modalBackdropClass := "ModalBackdrop"

	formGroupTranslationId := hb.Input().
		Type(hb.TYPE_HIDDEN).
		Name("translation_id").
		Value(data.translationID)

	buttonDelete := hb.Button().
		HTML("Delete").
		Class("btn btn-primary float-end").
		HxInclude("#Modal" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalTranslationDelete").
		HxTarget("body").
		HxSwap("beforeend")

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Delete Translation").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalTranslationDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Child(hb.Script(jsCloseFn)).
		Child(bs.ModalDialog().
			Child(bs.ModalContent().
				Child(
					bs.ModalHeader().
						Child(modalHeading).
						Child(modalClose)).
				Child(
					bs.ModalBody().
						Child(hb.Paragraph().Text("Are you sure you want to delete this translation?").Style(`margin-bottom:20px;color:red;`)).
						Child(hb.Paragraph().Text("This action cannot be undone.")).
						Child(formGroupTranslationId)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(
						hb.Button().HTML("Close").
							Class("btn btn-secondary float-start").
							Data("bs-dismiss", "modal").
							OnClick(modalCloseScript)).
					Child(buttonDelete)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().
		Children([]hb.TagInterface{
			modal,
			backdrop,
		})
}

func (controller *translationDeleteController) prepareDataAndValidate(r *http.Request) (data translationDeleteControllerData, errorMessage string) {
	data.request = r
	data.translationID = utils.Req(r, "translation_id", "")

	if data.translationID == "" {
		return data, "translation id is required"
	}

	translation, err := controller.ui.Store().TranslationFindByID(r.Context(), data.translationID)

	if err != nil {
		controller.ui.Logger().Error("Error. At translationDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if translation == nil {
		return data, "Translation not found"
	}

	data.translation = translation

	if r.Method != "POST" {
		return data, ""
	}

	err = controller.ui.Store().TranslationSoftDelete(r.Context(), translation)

	if err != nil {
		controller.ui.Logger().Error("Error. At translationDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.successMessage = "translation deleted successfully."

	return data, ""

}
