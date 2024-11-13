package admin

import (
	"net/http"

	// "project/config"
	// "project/controllers/admin/cms/shared"
	// "project/internal/helpers"
	// "project/pkg/cmsstore"
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

type pageCreateController struct {
	ui UiInterface
}

type pageCreateControllerData struct {
	request        *http.Request
	siteList       []cmsstore.SiteInterface
	siteID         string
	name           string
	successMessage string
}

var _ router.HTMLControllerInterface = (*pageDeleteController)(nil)

// == CONSTRUCTOR =============================================================

func NewPageCreateController(ui UiInterface) *pageCreateController {
	return &pageCreateController{
		ui: ui,
	}
}

func (controller pageCreateController) Handler(w http.ResponseWriter, r *http.Request) string {
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

func (controller *pageCreateController) modal(data pageCreateControllerData) hb.TagInterface {
	submitUrl := shared.URLR(data.request, shared.PathPagesPageCreate, nil)

	form := form.NewForm(form.FormOptions{
		ID: "FormPageCreate",
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label:    "Page name",
				Name:     "page_name",
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

	modalHeading := hb.Heading5().HTML("New Page").Style(`margin:0px;`)

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

func (controller *pageCreateController) prepareDataAndValidate(r *http.Request) (data pageCreateControllerData, errorMessage string) {
	data.request = r
	data.name = strings.TrimSpace(utils.Req(r, "page_name", ""))
	data.siteID = strings.TrimSpace(utils.Req(r, "site_id", ""))

	var err error

	data.siteList, err = controller.ui.Store().SiteList(cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC))

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
		return data, "page name is required"
	}

	page := cmsstore.NewPage()
	page.SetSiteID(data.siteID)
	page.SetName(data.name)

	err = controller.ui.Store().PageCreate(page)

	if err != nil {
		controller.ui.Logger().Error("At pageCreateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.successMessage = "page created successfully."

	return data, ""

}
