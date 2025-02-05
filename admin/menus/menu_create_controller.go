package admin

import (
	"net/http"
	"strings"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/router"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

// == CONTROLLER ==============================================================

type menuCreateController struct {
	ui UiInterface
}

type menuCreateControllerData struct {
	request        *http.Request
	siteList       []cmsstore.SiteInterface
	siteID         string
	name           string
	successMessage string
}

var _ router.HTMLControllerInterface = (*menuCreateController)(nil)

// == CONSTRUCTOR =============================================================

func NewMenuCreateController(ui UiInterface) *menuCreateController {
	return &menuCreateController{
		ui: ui,
	}
}

func (controller menuCreateController) Handler(w http.ResponseWriter, r *http.Request) string {
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

func (controller *menuCreateController) modal(data menuCreateControllerData) hb.TagInterface {
	submitUrl := shared.URLR(data.request, shared.PathMenusMenuCreate, nil)

	form := form.NewForm(form.FormOptions{
		ID: "FormMenuCreate",
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label:    "Menu name",
				Name:     "menu_name",
				Type:     form.FORM_FIELD_TYPE_STRING,
				Value:    data.name,
				Required: true,
			}),
			form.NewField(form.FieldOptions{
				Label:    "Site",
				Name:     "site_id",
				Type:     form.FORM_FIELD_TYPE_SELECT,
				Value:    data.siteID,
				Required: true,
				Options: append([]form.FieldOption{
					{
						Value: "Select site",
						Key:   "",
					},
				},
					lo.Map(data.siteList, func(site cmsstore.SiteInterface, index int) form.FieldOption {
						return form.FieldOption{
							Value: site.Name(),
							Key:   site.ID(),
						}
					})...),
			}),
		},
	})

	modalID := "ModalPageCreate"
	modalBackdropClass := "ModalBackdrop"

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("New Menu").Style(`margin:0px;`)

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
						Child(form.Build())).
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

func (controller *menuCreateController) prepareDataAndValidate(r *http.Request) (data menuCreateControllerData, errorMessage string) {
	data.request = r
	data.name = strings.TrimSpace(utils.Req(r, "menu_name", ""))
	data.siteID = strings.TrimSpace(utils.Req(r, "site_id", ""))

	var err error

	data.siteList, err = controller.ui.Store().SiteList(r.Context(), cmsstore.SiteQuery().SetOrderBy(cmsstore.COLUMN_NAME).SetSortOrder(sb.ASC))

	if err != nil {
		controller.ui.Logger().Error("At pageCreateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if r.Method != http.MethodPost {
		return data, ""
	}

	if data.siteID == "" {
		return data, "site id is required"
	}

	if data.name == "" {
		return data, "menu name is required"
	}

	menu := cmsstore.NewMenu()
	menu.SetSiteID(data.siteID)
	menu.SetName(data.name)

	err = controller.ui.Store().MenuCreate(r.Context(), menu)

	if err != nil {
		controller.ui.Logger().Error("At menuCreateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.successMessage = "menu created successfully."

	return data, ""

}
