package admin

import (
	"net/http"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/router"
	"github.com/gouniverse/utils"
)

// == CONTROLLER ==============================================================

type templateDeleteController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*templateDeleteController)(nil)

// == CONSTRUCTOR =============================================================

type templateDeleteControllerData struct {
	templateID     string
	template       cmsstore.TemplateInterface
	successMessage string
}

func NewTemplateDeleteController(ui UiInterface) *templateDeleteController {
	return &templateDeleteController{
		ui: ui,
	}
}

func (controller templateDeleteController) Handler(w http.ResponseWriter, r *http.Request) string {
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

func (controller *templateDeleteController) modal(data templateDeleteControllerData) hb.TagInterface {
	submitUrl := controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathTemplateDelete(), map[string]string{
		"template_id": data.templateID,
	})

	modalID := "ModalTemplateDelete"
	modalBackdropClass := "ModalBackdrop"

	formGroupTemplateId := hb.Input().
		Type(hb.TYPE_HIDDEN).
		Name("template_id").
		Value(data.templateID)

	buttonDelete := hb.Button().
		HTML("Delete").
		Class("btn btn-primary float-end").
		HxInclude("#Modal" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalTemplateDelete").
		HxTarget("body").
		HxSwap("beforeend")

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Delete Template").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalTemplateDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

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
						Child(hb.Paragraph().Text("Are you sure you want to delete this template?").Style(`margin-bottom:20px;color:red;`)).
						Child(hb.Paragraph().Text("This action cannot be undone.")).
						Child(formGroupTemplateId)).
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

func (controller *templateDeleteController) prepareDataAndValidate(r *http.Request) (data templateDeleteControllerData, errorMessage string) {
	data.templateID = utils.Req(r, "template_id", "")

	if data.templateID == "" {
		return data, "template id is required"
	}

	template, err := controller.ui.Store().TemplateFindByID(data.templateID)

	if err != nil {
		controller.ui.Logger().Error("Error. At templateDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if template == nil {
		return data, "Template not found"
	}

	data.template = template

	if r.Method != "POST" {
		return data, ""
	}

	err = controller.ui.Store().TemplateSoftDelete(template)

	if err != nil {
		controller.ui.Logger().Error("Error. At templateDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.successMessage = "template deleted successfully."

	return data, ""

}
