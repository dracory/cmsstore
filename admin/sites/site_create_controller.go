package admin

import (
	"net/http"

	// "project/config"
	// "project/controllers/admin/cms/shared"
	// "project/internal/helpers"
	// "project/pkg/cmsstore"

	"github.com/dracory/bs"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

// == CONTROLLER ==============================================================

type siteCreateController struct {
	ui UiInterface
}

type siteCreateControllerData struct {
	request        *http.Request
	name           string
	successMessage string
}

// == CONSTRUCTOR =============================================================

func NewSiteCreateController(ui UiInterface) *siteCreateController {
	return &siteCreateController{
		ui: ui,
	}
}

func (controller siteCreateController) Handler(w http.ResponseWriter, r *http.Request) string {
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

func (controller *siteCreateController) modal(data siteCreateControllerData) hb.TagInterface {
	submitUrl := shared.URLR(data.request, shared.PathSitesSiteCreate, nil)

	formGroupName := bs.FormGroup().
		Class("mb-3").
		Child(bs.FormLabel("Website name")).
		Child(bs.FormInput().Name("site_name").Value(data.name))

	modalID := "ModalPageCreate"
	modalBackdropClass := "ModalBackdrop"

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("New Site").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalPageCreate').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

	buttonSend := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Create & Edit").
		Class("btn btn-primary float-end").
		HxInclude("#" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalpageCreate").
		HxTarget("body").
		HxSwap("beforeend")

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Close").
		Class("btn btn-secondary float-start").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

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
						Child(formGroupName)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(buttonCancel).
					Child(buttonSend)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})
}

func (controller *siteCreateController) prepareDataAndValidate(r *http.Request) (data siteCreateControllerData, errorMessage string) {
	data.request = r
	data.name = req.GetStringTrimmed(r, "site_name")

	if r.Method != http.MethodPost {
		return data, ""
	}

	if data.name == "" {
		return data, "site name is required"
	}

	site := cmsstore.NewSite()
	site.SetName(data.name)

	err := controller.ui.Store().SiteCreate(r.Context(), site)

	if err != nil {
		controller.ui.Logger().Error("At siteCreateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.successMessage = "site created successfully."

	return data, ""

}
