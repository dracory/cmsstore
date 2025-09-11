package admin

import (
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

// == CONTROLLER ==============================================================

type menuDeleteController struct {
	ui UiInterface
}

// == CONSTRUCTOR =============================================================

type menuDeleteControllerData struct {
	request        *http.Request
	menuID         string
	menu           cmsstore.MenuInterface
	successMessage string
}

func NewMenuDeleteController(ui UiInterface) *menuDeleteController {
	return &menuDeleteController{
		ui: ui,
	}
}

func (controller menuDeleteController) Handler(w http.ResponseWriter, r *http.Request) string {
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

func (controller *menuDeleteController) modal(data menuDeleteControllerData) hb.TagInterface {
	submitUrl := shared.URLR(data.request, shared.PathMenusMenuDelete, map[string]string{
		"menu_id": data.menuID,
	})

	modalID := "ModalMenuDelete"
	modalBackdropClass := "ModalBackdrop"

	formGroupMenuId := hb.Input().
		Type(hb.TYPE_HIDDEN).
		Name("menu_id").
		Value(data.menuID)

	buttonDelete := hb.Button().
		HTML("Delete").
		Class("btn btn-primary float-end").
		HxInclude("#Modal" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalMenuDelete").
		HxTarget("body").
		HxSwap("beforeend")

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Delete Menu").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalMenuDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

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
						Child(hb.Paragraph().Text("Are you sure you want to delete this menu?").Style(`margin-bottom:20px;color:red;`)).
						Child(hb.Paragraph().Text("This action cannot be undone.")).
						Child(formGroupMenuId)).
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

func (controller *menuDeleteController) prepareDataAndValidate(r *http.Request) (data menuDeleteControllerData, errorMessage string) {
	data.request = r
	data.menuID = req.GetStringTrimmed(r, "menu_id")

	if data.menuID == "" {
		return data, "menu id is required"
	}

	menu, err := controller.ui.Store().MenuFindByID(r.Context(), data.menuID)

	if err != nil {
		controller.ui.Logger().Error("Error. At menuDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if menu == nil {
		return data, "Menu not found"
	}

	data.menu = menu

	if r.Method != "POST" {
		return data, ""
	}

	err = controller.ui.Store().MenuSoftDelete(r.Context(), menu)

	if err != nil {
		controller.ui.Logger().Error("Error. At menuDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.successMessage = "menu deleted successfully."

	return data, ""

}
