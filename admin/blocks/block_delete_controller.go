package admin

import (
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/gouniverse/router"
)

// == CONTROLLER ==============================================================

type blockDeleteController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*blockDeleteController)(nil)

// == CONSTRUCTOR =============================================================

type blockDeleteControllerData struct {
	request        *http.Request
	blockID        string
	block          cmsstore.BlockInterface
	successMessage string
}

func NewBlockDeleteController(ui UiInterface) *blockDeleteController {
	return &blockDeleteController{
		ui: ui,
	}
}

func (controller blockDeleteController) Handler(w http.ResponseWriter, r *http.Request) string {
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

func (controller *blockDeleteController) modal(data blockDeleteControllerData) hb.TagInterface {
	submitUrl := shared.URLR(data.request, shared.PathBlocksBlockDelete, map[string]string{
		"block_id": data.blockID,
	})

	modalID := "ModalBlockDelete"
	modalBackdropClass := "ModalBackdrop"

	formGroupBlockId := hb.Input().
		Type(hb.TYPE_HIDDEN).
		Name("block_id").
		Value(data.blockID)

	buttonDelete := hb.Button().
		HTML("Delete").
		Class("btn btn-primary float-end").
		HxInclude("#Modal" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalBlockDelete").
		HxTarget("body").
		HxSwap("beforeend")

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Delete Block").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalBlockDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

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
						Child(hb.Paragraph().Text("Are you sure you want to delete this block?").Style(`margin-bottom:20px;color:red;`)).
						Child(hb.Paragraph().Text("This action cannot be undone.")).
						Child(formGroupBlockId)).
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

func (controller *blockDeleteController) prepareDataAndValidate(r *http.Request) (data blockDeleteControllerData, errorMessage string) {
	data.request = r
	data.blockID = req.GetStringTrimmed(r, "block_id")

	if data.blockID == "" {
		return data, "block id is required"
	}

	block, err := controller.ui.Store().BlockFindByID(r.Context(), data.blockID)

	if err != nil {
		controller.ui.Logger().Error("Error. At blockDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if block == nil {
		return data, "Block not found"
	}

	data.block = block

	if r.Method != "POST" {
		return data, ""
	}

	err = controller.ui.Store().BlockSoftDelete(r.Context(), block)

	if err != nil {
		controller.ui.Logger().Error("Error. At blockDeleteController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.successMessage = "block deleted successfully."

	return data, ""

}
