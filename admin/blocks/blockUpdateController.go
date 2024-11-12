package admin

import (
	"net/http"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/router"
	"github.com/gouniverse/utils"
)

const VIEW_SETTINGS = "settings"
const VIEW_CONTENT = "content"

const codemirrorCss = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.css"
const codemirrorJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.js"
const codemirrorXmlJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/xml/xml.min.js"
const codemirrorHtmlmixedJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/htmlmixed/htmlmixed.min.js"
const codemirrorJavascriptJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/javascript/javascript.js"
const codemirrorCssJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/css/css.js"
const codemirrorClikeJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/clike/clike.min.js"
const codemirrorPhpJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/php/php.min.js"
const codemirrorFormattingJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/2.36.0/formatting.min.js"
const codemirrorMatchBracketsJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.22.0/addon/edit/matchbrackets.min.js"

// == CONTROLLER ==============================================================

type blockUpdateController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*blockUpdateController)(nil)

// == CONSTRUCTOR =============================================================

func NewBlockUpdateController(ui UiInterface) *blockUpdateController {
	return &blockUpdateController{
		ui: ui,
	}
}

func (controller *blockUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		return api.Error(errorMessage).ToString()
	}

	if r.Method == http.MethodPost {
		return controller.form(data).ToHTML()
	}

	html := controller.page(data)

	options := struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		Styles: []string{
			`.CodeMirror {
				border: 1px solid #eee;
				height: auto;
			}
			`,
		},
		StyleURLs: []string{
			codemirrorCss,
		},
		Scripts: []string{},
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
			cdn.Htmx_2_0_0(),
			cdn.Jquery_3_7_1(),
			codemirrorJs,
			codemirrorXmlJs,
			codemirrorHtmlmixedJs,
			codemirrorJavascriptJs,
			codemirrorCssJs,
			codemirrorClikeJs,
			codemirrorPhpJs,
			codemirrorFormattingJs,
			codemirrorMatchBracketsJs,
		},
	}

	return controller.ui.Layout(w, r, "Edit Block | CMS", html.ToHTML(), options)
}

func (controller blockUpdateController) page(data blockUpdateControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), controller.ui.Endpoint())
	breadcrumbs := controller.ui.AdminBreadcrumbs(controller.ui.Endpoint(), []shared.Breadcrumb{
		{
			Name: "Block Manager",
			URL:  shared.URL(controller.ui.Endpoint(), shared.PathBlocksBlockManager, nil),
		},
		{
			Name: "Edit Block",
			URL:  shared.URL(controller.ui.Endpoint(), shared.PathBlocksBlockUpdate, map[string]string{"block_id": data.blockID}),
		},
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormBlockUpdate").
		HxPost(shared.URL(controller.ui.Endpoint(), shared.PathBlocksBlockUpdate, map[string]string{"block_id": data.blockID})).
		HxTarget("#FormBlockUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URL(controller.ui.Endpoint(), shared.PathBlocksBlockManager, nil))

	badgeStatus := hb.Div().
		Class("badge fs-6 ms-3").
		ClassIf(data.block.Status() == cmsstore.TEMPLATE_STATUS_ACTIVE, "bg-success").
		ClassIf(data.block.Status() == cmsstore.TEMPLATE_STATUS_INACTIVE, "bg-secondary").
		ClassIf(data.block.Status() == cmsstore.TEMPLATE_STATUS_DRAFT, "bg-warning").
		Text(data.block.Status())

	pageTitle := hb.Heading1().
		Text("CMS. Edit Block:").
		Text(" ").
		Text(data.block.Name()).
		Child(hb.Sup().Child(badgeStatus)).
		Child(buttonSave).
		Child(buttonCancel)

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTMLIf(data.view == VIEW_CONTENT, "Block Content").
					HTMLIf(data.view == VIEW_SETTINGS, "Block Settings").
					Style("margin-bottom:0;display:inline-block;")).
				Child(buttonSave),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(controller.form(data)))

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_CONTENT, "active").
				Href(shared.URL(controller.ui.Endpoint(), shared.PathBlocksBlockUpdate, map[string]string{
					"block_id": data.blockID,
					"view":     VIEW_CONTENT,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SETTINGS, "active").
				Href(shared.URL(controller.ui.Endpoint(), shared.PathBlocksBlockUpdate, map[string]string{
					"block_id": data.blockID,
					"view":     VIEW_SETTINGS,
				})).
				HTML("Settings")))

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader).
		Child(hb.HR()).
		Child(pageTitle).
		Child(tabs).
		Child(card).
		Child(hb.HR().Class("mt-4")).
		Child(hb.Div().
			Class("text-info mb-2").
			Text("To use this block in your website use the following shortcode:").
			Child(hb.BR())).
		Child(hb.PRE().
			Child(hb.Code().
				Text(`<!-- START: Block: ` + data.block.Name() + ` -->`).
				Text("\n").
				Text(`[[BLOCK_` + data.blockID + `]]`).
				Text("\n").
				Text(`<!-- END: Block: ` + data.block.Name() + ` -->`)))
}

func (controller blockUpdateController) form(data blockUpdateControllerData) hb.TagInterface {

	fieldsContent := controller.fieldsContent(data)
	fieldsSettings := controller.fieldsSettings(data)

	formpageUpdate := form.NewForm(form.FormOptions{
		ID: "FormBlockUpdate",
	})

	if data.view == VIEW_SETTINGS {
		formpageUpdate.SetFields(fieldsSettings)
	}

	if data.view == VIEW_CONTENT {
		formpageUpdate.SetFields(fieldsContent)
	}

	if data.formErrorMessage != "" {
		formpageUpdate.AddField(&form.Field{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
		})
	}

	if data.formSuccessMessage != "" {
		formpageUpdate.AddField(&form.Field{
			Type: form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{
				Icon:              "success",
				Text:              data.formSuccessMessage,
				Position:          "top-end",
				Timer:             1500,
				ShowConfirmButton: false,
				ShowCancelButton:  false,
			}).ToHTML(),
		})
	}

	return formpageUpdate.Build()
}

func (blockUpdateController) fieldsContent(data blockUpdateControllerData) []form.FieldInterface {
	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Content (HTML)",
			Name:  "block_content",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formContent,
		}),
		form.NewField(form.FieldOptions{
			Label:    "Block ID",
			Name:     "block_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.blockID,
			Readonly: true,
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_CONTENT,
			Readonly: true,
		}),
	}

	contentScript := hb.Script(`
function codeMirrorSelector() {
	return 'textarea[name="block_content"]';
}
function getCodeMirrorEditor() {
	return document.querySelector(codeMirrorSelector());
}
setTimeout(function () {
    console.log(getCodeMirrorEditor());
	if (getCodeMirrorEditor()) {
		var editor = CodeMirror.fromTextArea(getCodeMirrorEditor(), {
			lineNumbers: true,
			matchBrackets: true,
			mode: "application/x-httpd-php",
			indentUnit: 4,
			indentWithTabs: true,
			enterMode: "keep", tabMode: "shift"
		});
		$(document).on('mouseup', codeMirrorSelector(), function() {
			getCodeMirrorEditor().value = editor.getValue();
		});
		$(document).on('change', codeMirrorSelector(), function() {
			getCodeMirrorEditor().value = editor.getValue();
		});
		setInterval(()=>{
			getCodeMirrorEditor().value = editor.getValue();
		}, 1000)
	}
}, 500);
		`).ToHTML()

	fieldsContent = append(fieldsContent, &form.Field{
		Type:  form.FORM_FIELD_TYPE_RAW,
		Value: contentScript,
	})

	return fieldsContent
}

func (controller blockUpdateController) fieldsSettings(data blockUpdateControllerData) []form.FieldInterface {
	fieldsSettings := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Status",
			Name:  "block_status",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formStatus,
			Help:  "The status of this webpage. Published pages will be displayed on the webblock.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "Draft",
					Key:   cmsstore.BLOCK_STATUS_DRAFT,
				},
				{
					Value: "Published",
					Key:   cmsstore.BLOCK_STATUS_ACTIVE,
				},
				{
					Value: "Unpublished",
					Key:   cmsstore.BLOCK_STATUS_INACTIVE,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Block Name (Internal)",
			Name:  "block_name",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formName,
			Help:  "The name of the block as displayed in the admin panel. This is not vsible to the block vistors",
		}),
		form.NewField(form.FieldOptions{
			Label: "Admin Notes (Internal)",
			Name:  "block_memo",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formMemo,
			Help:  "Admin notes for this block. These notes will not be visible to the public.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "Webblock ID",
			Name:     "block_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.blockID,
			Readonly: true,
			Help:     "The reference number (ID) of the webblock. This is used to identify the webblock in the system and should not be changed.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.view,
			Readonly: true,
		}),
	}

	return fieldsSettings
}

func (controller blockUpdateController) saveBlock(r *http.Request, data blockUpdateControllerData) (d blockUpdateControllerData, errorMessage string) {
	data.formContent = utils.Req(r, "block_content", "")
	data.formMemo = utils.Req(r, "block_memo", "")
	data.formName = utils.Req(r, "block_name", "")
	data.formStatus = utils.Req(r, "block_status", "")
	data.formTitle = utils.Req(r, "block_title", "")

	if data.view == VIEW_SETTINGS {
		if data.formStatus == "" {
			data.formErrorMessage = "Status is required"
			return data, ""
		}
	}

	if data.view == VIEW_SETTINGS {
		data.block.SetMemo(data.formMemo)
		data.block.SetName(data.formName)
		data.block.SetStatus(data.formStatus)
	}

	if data.view == VIEW_CONTENT {
		data.block.SetContent(data.formContent)
	}

	err := controller.ui.Store().BlockUpdate(data.block)

	if err != nil {
		//config.LogStore.ErrorWithContext("At blockUpdateController > prepareDataAndValidate", err.Error())
		data.formErrorMessage = "System error. Saving block failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "block saved successfully"

	return data, ""
}

func (controller blockUpdateController) prepareDataAndValidate(r *http.Request) (data blockUpdateControllerData, errorMessage string) {
	data.action = utils.Req(r, "action", "")
	data.blockID = utils.Req(r, "block_id", "")
	data.view = utils.Req(r, "view", "")

	if data.view == "" {
		data.view = VIEW_CONTENT
	}

	if data.blockID == "" {
		return data, "block id is required"
	}

	var err error
	data.block, err = controller.ui.Store().BlockFindByID(data.blockID)

	if err != nil {
		controller.ui.Logger().Error("At blockUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if data.block == nil {
		return data, "block not found"
	}

	data.formContent = data.block.Content()
	data.formName = data.block.Name()
	data.formMemo = data.block.Memo()
	data.formStatus = data.block.Status()

	if r.Method != http.MethodPost {
		return data, ""
	}

	return controller.saveBlock(r, data)
}

type blockUpdateControllerData struct {
	action  string
	blockID string
	block   cmsstore.BlockInterface
	view    string

	formErrorMessage   string
	formSuccessMessage string
	formContent        string
	formName           string
	formMemo           string
	formStatus         string
	formTitle          string
}
