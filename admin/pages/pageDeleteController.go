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

type pageDeleteController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*pageCreateController)(nil)

// == CONSTRUCTOR =============================================================

type pageDeleteControllerData struct {
	request        *http.Request
	pageID         string
	page           cmsstore.PageInterface
	successMessage string
}

func NewPageDeleteController(ui UiInterface) *pageDeleteController {
	return &pageDeleteController{
		ui: ui,
	}
}

func (controller pageDeleteController) Handler(w http.ResponseWriter, r *http.Request) string {
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

func (controller *pageDeleteController) modal(data pageDeleteControllerData) hb.TagInterface {
	submitUrl := shared.URL(shared.Endpoint(data.request), shared.PathPagesPageDelete, map[string]string{
		"page_id": data.pageID,
	})

	modalID := "ModalPageDelete"
	modalBackdropClass := "ModalBackdrop"

	formGroupPageId := hb.Input().
		Type(hb.TYPE_HIDDEN).
		Name("page_id").
		Value(data.pageID)

	buttonDelete := hb.Button().
		HTML("Delete").
		Class("btn btn-primary float-end").
		HxInclude("#Modal" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalPageDelete").
		HxTarget("body").
		HxSwap("beforeend")

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Delete Page").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalPageDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

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
						Child(hb.Paragraph().Text("Are you sure you want to delete this page?").Style(`margin-bottom:20px;color:red;`)).
						Child(hb.Paragraph().Text("This action cannot be undone.")).
						Child(formGroupPageId)).
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

func (controller *pageDeleteController) prepareDataAndValidate(r *http.Request) (data pageDeleteControllerData, errorMessage string) {
	data.request = r
	data.pageID = utils.Req(r, "page_id", "")

	if data.pageID == "" {
		return data, "page id is required"
	}

	page, err := controller.ui.Store().PageFindByID(data.pageID)

	if err != nil {
		controller.ui.Logger().Error("Error. At pageDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if page == nil {
		return data, "Page not found"
	}

	data.page = page

	if r.Method != "POST" {
		return data, ""
	}

	err = controller.ui.Store().PageSoftDelete(page)

	if err != nil {
		controller.ui.Logger().Error("Error. At pageDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.successMessage = "page deleted successfully."

	return data, ""

}
