package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/samber/lo"
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
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), data.request)

	breadcrumbs := shared.AdminBreadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Block Manager",
			URL:  shared.URLR(data.request, shared.PathBlocksBlockManager, nil),
		},
		{
			Name: "Edit Block",
			URL:  shared.URLR(data.request, shared.PathBlocksBlockUpdate, map[string]string{"block_id": data.blockID}),
		},
	}, struct{ SiteList []cmsstore.SiteInterface }{
		SiteList: data.siteList,
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormBlockUpdate").
		HxPost(shared.URLR(data.request, shared.PathBlocksBlockUpdate, map[string]string{"block_id": data.blockID})).
		HxTarget("#FormBlockUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URLR(data.request, shared.PathBlocksBlockManager, nil))

	buttonVersion := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-code-slash").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Version History").
		HxGet(shared.URLR(data.request, shared.PathBlocksBlockVersioning, map[string]string{
			"block_id": data.blockID,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	badgeStatus := hb.Div().
		Class("badge fs-6 ms-3").
		ClassIf(data.block.Status() == cmsstore.BLOCK_STATUS_ACTIVE, "bg-success").
		ClassIf(data.block.Status() == cmsstore.BLOCK_STATUS_INACTIVE, "bg-secondary").
		ClassIf(data.block.Status() == cmsstore.BLOCK_STATUS_DRAFT, "bg-warning").
		Text(lo.If(data.block.Status() == cmsstore.BLOCK_STATUS_ACTIVE, "Published").
			ElseIf(data.block.Status() == cmsstore.BLOCK_STATUS_INACTIVE, "Unpublished").
			ElseIf(data.block.Status() == cmsstore.BLOCK_STATUS_DRAFT, "Draft").
			Else(data.block.Status()))

	pageTitle := hb.Heading1().
		Text("CMS. Edit Block:").
		Text(" ").
		Text(data.block.Name()).
		Child(hb.Sup().Child(badgeStatus)).
		Child(buttonSave).
		Child(buttonVersion).
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
				Href(shared.URLR(data.request, shared.PathBlocksBlockUpdate, map[string]string{
					"block_id": data.blockID,
					"view":     VIEW_CONTENT,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SETTINGS, "active").
				Href(shared.URLR(data.request, shared.PathBlocksBlockUpdate, map[string]string{
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
			Class("text-info mb-3").
			Text("To use this block in your website, click to copy one of the following syntaxes:")).
		Child(controller.codeSnippet(
			"Legacy Syntax",
			`<!-- START: Block: `+data.block.Name()+` -->`+"\n"+
				`[[BLOCK_`+data.blockID+`]]`+"\n"+
				`<!-- END: Block: `+data.block.Name()+` -->`)).
		Child(controller.codeSnippet(
			"New Attribute Syntax (Recommended)",
			`<!-- START: Block with attributes: `+data.block.Name()+` -->`+"\n"+
				`<block id="`+data.blockID+`" />`+"\n"+
				`<!-- END: Block with attributes: `+data.block.Name()+` -->`)).
		Child(controller.codeSnippet(
			"With Runtime Attributes",
			`<block id="`+data.blockID+`" depth="2" style="horizontal" />`)).
		Child(controller.codeSnippet(
			"Alternative Syntax (for HTML attributes)",
			`[[block id='`+data.blockID+`']]`)).
		Child(controller.customVariablesSection(data))
}

// customVariablesSection renders a reference table of custom variables exposed by the block type.
// Returns nil if the block type exposes no custom variables.
func (controller blockUpdateController) customVariablesSection(data blockUpdateControllerData) hb.TagInterface {
	blockType := cmsstore.GetBlockType(data.block.Type())
	if blockType == nil {
		return hb.Div()
	}

	vars := blockType.GetCustomVariables()
	if len(vars) == 0 {
		return hb.Div()
	}

	rows := hb.Tbody()
	for _, v := range vars {
		rows.Child(hb.TR().
			Child(hb.TD().Child(hb.Code().Text("[[" + v.Name + "]]"))).
			Child(hb.TD().Text(v.Description)))
	}

	return hb.Div().
		Class("card mt-4").
		Child(hb.Div().
			Class("card-header").
			Text("Available Custom Variables")).
		Child(hb.Div().
			Class("card-body p-0").
			Child(hb.Div().
				Class("text-muted px-3 pt-3 pb-2").
				Text("This block type sets the following variables. Use them in your page or template content as shown.")).
			Child(hb.Table().
				Class("table table-sm table-bordered mb-0").
				Child(hb.Thead().
					Child(hb.TR().
						Child(hb.TH().Text("Variable")).
						Child(hb.TH().Text("Description")))).
				Child(rows)))
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

	if data.formRedirectURL != "" {
		formpageUpdate.AddField(&form.Field{
			Type: form.FORM_FIELD_TYPE_RAW,
			Value: hb.Script(`window.location.href = "` + data.formRedirectURL + `";`).
				ToHTML(),
		})
	}

	return formpageUpdate.Build()
}

func (controller blockUpdateController) fieldsContent(data blockUpdateControllerData) []form.FieldInterface {
	blockType := data.block.Type()

	if blockType == "" {
		blockType = cmsstore.BLOCK_TYPE_HTML
	}

	var fieldsContent []form.FieldInterface

	// First, check global BlockType registry
	globalBlockType := cmsstore.GetBlockType(blockType)
	if globalBlockType != nil {
		fields := globalBlockType.GetAdminFields(data.block, data.request)
		if formFields, ok := fields.([]form.FieldInterface); ok {
			fieldsContent = formFields
		} else {
			controller.ui.Logger().Error("GetAdminFields returned unexpected type",
				"blockType", blockType,
				"actualType", fmt.Sprintf("%T", fields))
			// Fall through to legacy provider fallback
		}
	}

	if len(fieldsContent) == 0 {
		// Fall back to local admin provider registry
		registry := controller.ui.BlockAdminRegistry()
		provider := registry.GetProvider(blockType)

		if provider != nil {
			fields := provider.GetContentFields(data.block, data.request)
			if formFields, ok := fields.([]form.FieldInterface); ok {
				fieldsContent = formFields
			} else {
				controller.ui.Logger().Error("GetContentFields returned unexpected type",
					"blockType", blockType,
					"actualType", fmt.Sprintf("%T", fields))
				// Fall through to HTML provider fallback
			}
		}

		if len(fieldsContent) == 0 && provider == nil {
			// Fallback to HTML provider if block type not registered
			htmlProvider := registry.GetProvider(cmsstore.BLOCK_TYPE_HTML)
			if htmlProvider != nil {
				fields := htmlProvider.GetContentFields(data.block, data.request)
				if formFields, ok := fields.([]form.FieldInterface); ok {
					fieldsContent = formFields
				} else {
					controller.ui.Logger().Error("HTML provider GetContentFields returned unexpected type",
						"actualType", fmt.Sprintf("%T", fields))
					// Fall through to ultimate fallback
				}
			}

			if len(fieldsContent) == 0 {
				// Ultimate fallback - basic textarea
				fieldsContent = []form.FieldInterface{
					form.NewField(form.FieldOptions{
						Label: "Content",
						Name:  "block_content",
						Type:  form.FORM_FIELD_TYPE_TEXTAREA,
						Value: data.block.Content(),
						Help:  "No admin provider registered for block type: " + blockType,
					}),
				}
			}
		}
	}

	fieldsContent = append(fieldsContent,
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
	)

	return fieldsContent
}

func (controller blockUpdateController) fieldsSettings(data blockUpdateControllerData) []form.FieldInterface {
	fieldSiteID := &form.Field{
		Label: "Belongs to Site",
		Name:  "block_site_id",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formSiteID,
		Help:  "The site that this block belongs to",
		OptionsF: func() []form.FieldOption {
			options := []form.FieldOption{
				{
					Value: "- no site selected -",
					Key:   "",
				},
			}
			for _, site := range data.siteList {
				name := site.Name()
				status := site.Status()
				options = append(options, form.FieldOption{
					Value: name + " (" + status + ")",
					Key:   site.ID(),
				})
			}
			return options
		},
	}

	fieldStatus := form.NewField(form.FieldOptions{
		Label: "Status",
		Name:  "block_status",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formStatus,
		Help:  "The status of this block. Published blocks will be displayed on the site.",
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
	})

	// Get block type and create display value
	blockType := data.block.Type()
	if blockType == "" {
		blockType = cmsstore.BLOCK_TYPE_HTML
	}

	// Determine if type field should be editable (only for draft blocks)
	isTypeEditable := data.block.Status() == cmsstore.BLOCK_STATUS_DRAFT

	var fieldType form.FieldInterface
	if isTypeEditable {
		// Create editable type field for draft blocks
		fieldType = form.NewField(form.FieldOptions{
			Label: "Block Type",
			Name:  "block_type",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: blockType,
			Help:  "Block type determines how the block is rendered. Can only be changed while in draft status.",
			OptionsF: func() []form.FieldOption {
				options := []form.FieldOption{}

				// Add block types from global registry
				globalTypes := cmsstore.GetAllBlockTypes()
				for typeKey, blockType := range globalTypes {
					options = append(options, form.FieldOption{
						Value: blockType.TypeLabel(),
						Key:   typeKey,
					})
				}

				// If no global types, fall back to basic types
				if len(options) == 0 {
					options = []form.FieldOption{
						{Value: "HTML", Key: cmsstore.BLOCK_TYPE_HTML},
						{Value: "Menu", Key: cmsstore.BLOCK_TYPE_MENU},
						{Value: "Navbar", Key: cmsstore.BLOCK_TYPE_NAVBAR},
						{Value: "Breadcrumbs", Key: cmsstore.BLOCK_TYPE_BREADCRUMBS},
					}
				}

				return options
			},
		})
	} else {
		// Create readonly type field for published blocks
		fieldType = form.NewField(form.FieldOptions{
			Label:    "Block Type",
			Name:     "block_type",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    blockType, // Use actual block type value, not display label
			Readonly: true,
			Help:     "Block type cannot be changed after publication. This determines how the block is rendered.",
		})
	}

	fieldBlockName := form.NewField(form.FieldOptions{
		Label: "Block Name (Internal)",
		Name:  "block_name",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formName,
		Help:  "The name of the block as displayed in the admin panel. This is not vsible to the block vistors",
	})

	fieldMemo := form.NewField(form.FieldOptions{
		Label: "Admin Notes (Internal)",
		Name:  "block_memo",
		Type:  form.FORM_FIELD_TYPE_TEXTAREA,
		Value: data.formMemo,
		Help:  "Admin notes for this block. These notes will not be visible to the public.",
	})

	fieldBlockID := form.NewField(form.FieldOptions{
		Label:    "Block Reference / ID",
		Name:     "block_id",
		Type:     form.FORM_FIELD_TYPE_STRING,
		Value:    data.blockID,
		Readonly: true,
		Help:     "The reference number (ID) of the block. This is used to identify the block in the system and should not be changed.",
	})

	fieldView := form.NewField(form.FieldOptions{
		Label:    "View",
		Name:     "view",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.view,
		Readonly: true,
	})

	fieldsSettings := []form.FieldInterface{
		fieldStatus,
		fieldType,
		fieldBlockName,
		fieldSiteID,
		fieldMemo,
		fieldBlockID,
		fieldView,
	}

	return fieldsSettings
}

func (controller blockUpdateController) saveBlock(r *http.Request, data blockUpdateControllerData) (d blockUpdateControllerData, errorMessage string) {
	data.formContent = req.GetStringTrimmed(r, "block_content")
	data.formMemo = req.GetStringTrimmed(r, "block_memo")
	data.formName = req.GetStringTrimmed(r, "block_name")
	data.formSiteID = req.GetStringTrimmed(r, "block_site_id")
	data.formStatus = req.GetStringTrimmed(r, "block_status")
	data.formTitle = req.GetStringTrimmed(r, "block_title")
	data.formType = req.GetStringTrimmed(r, "block_type")

	if data.view == VIEW_SETTINGS {
		if data.formStatus == "" {
			data.formErrorMessage = "Status is required"
			return data, ""
		}

		if data.formName == "" {
			data.formErrorMessage = "Name is required"
			return data, ""
		}

		// Validate block type change constraints
		if data.formType != "" && data.formType != data.block.Type() {
			// Only allow type changes for draft blocks
			if data.block.Status() != cmsstore.BLOCK_STATUS_DRAFT {
				data.formErrorMessage = "Block type can only be changed while the block is in draft status"
				return data, ""
			}

			// Validate that the new type exists in the registry
			globalBlockType := cmsstore.GetBlockType(data.formType)
			if globalBlockType == nil {
				// Check if it's a basic fallback type
				validBasicTypes := map[string]bool{
					cmsstore.BLOCK_TYPE_HTML:        true,
					cmsstore.BLOCK_TYPE_MENU:        true,
					cmsstore.BLOCK_TYPE_NAVBAR:      true,
					cmsstore.BLOCK_TYPE_BREADCRUMBS: true,
				}
				if !validBasicTypes[data.formType] {
					data.formErrorMessage = "Invalid block type: " + data.formType
					return data, ""
				}
			}
		}
	}

	if data.view == VIEW_SETTINGS {
		data.block.SetMemo(data.formMemo)
		data.block.SetName(data.formName)
		data.block.SetSiteID(data.formSiteID)
		data.block.SetStatus(data.formStatus)

		// Apply block type change if validated
		if data.formType != "" && data.formType != data.block.Type() {
			// Clear content and metadata when changing type to prevent conflicts
			data.block.SetContent("")
			data.block.SetMetas(map[string]string{})
			data.block.SetType(data.formType)
		}
	}

	if data.view == VIEW_CONTENT {
		blockType := data.block.Type()
		if blockType == "" {
			blockType = cmsstore.BLOCK_TYPE_HTML
		}

		// First, check global BlockType registry
		globalBlockType := cmsstore.GetBlockType(blockType)
		if globalBlockType != nil {
			err := globalBlockType.SaveAdminFields(r, data.block)
			if err != nil {
				data.formErrorMessage = err.Error()
				return data, ""
			}
		} else {
			// Fall back to local admin provider registry
			registry := controller.ui.BlockAdminRegistry()
			provider := registry.GetProvider(blockType)

			if provider != nil {
				err := provider.SaveContentFields(r, data.block)
				if err != nil {
					data.formErrorMessage = err.Error()
					return data, ""
				}
			} else {
				// Fallback to basic content save if no provider
				data.block.SetContent(data.formContent)
			}
		}
	}

	err := controller.ui.Store().BlockUpdate(r.Context(), data.block)

	if err != nil {
		//config.LogStore.ErrorWithContext("At blockUpdateController > prepareDataAndValidate", err.Error())
		data.formErrorMessage = "System error. Saving block failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "block updated successfully"

	data.formRedirectURL = shared.URLR(data.request, shared.PathBlocksBlockUpdate, map[string]string{
		"block_id": data.blockID,
		"view":     data.view,
	})

	return data, ""
}

func (controller blockUpdateController) prepareDataAndValidate(r *http.Request) (data blockUpdateControllerData, errorMessage string) {
	data.request = r
	data.action = req.GetStringTrimmed(r, "action")
	data.blockID = req.GetStringTrimmed(r, "block_id")
	data.view = req.GetStringTrimmed(r, "view")

	if data.view == "" {
		data.view = VIEW_CONTENT
	}

	if data.blockID == "" {
		return data, "block id is required"
	}

	var err error
	data.block, err = controller.ui.Store().BlockFindByID(r.Context(), data.blockID)

	if err != nil {
		controller.ui.Logger().Error("At blockUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if data.block == nil {
		return data, "block not found"
	}

	data.siteList, err = controller.ui.Store().SiteList(r.Context(), cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		return data, "Site list failed to be retrieved" + err.Error()
	}

	data.formContent = data.block.Content()
	data.formName = data.block.Name()
	data.formMemo = data.block.Memo()
	data.formSiteID = data.block.SiteID()
	data.formStatus = data.block.Status()
	data.formType = data.block.Type()

	if r.Method != http.MethodPost {
		return data, ""
	}

	return controller.saveBlock(r, data)
}

type blockUpdateControllerData struct {
	request *http.Request
	action  string
	blockID string
	block   cmsstore.BlockInterface
	view    string

	siteList []cmsstore.SiteInterface

	formErrorMessage   string
	formRedirectURL    string
	formSuccessMessage string
	formContent        string
	formName           string
	formMemo           string
	formSiteID         string
	formStatus         string
	formTitle          string
	formType           string
}

// codeSnippet creates a code snippet card with copy-to-clipboard functionality
func (controller blockUpdateController) codeSnippet(title, code string) hb.TagInterface {
	snippetID := "snippet-" + strings.ReplaceAll(strings.ToLower(title), " ", "-")

	return hb.Div().
		Class("card mb-3").
		Child(hb.Div().
			Class("card-header d-flex justify-content-between align-items-center").
			Child(hb.Span().Text(title)).
			Child(hb.Button().
				Class("btn btn-sm btn-outline-primary").
				Child(hb.I().Class("bi bi-clipboard")).
				HTML(" Copy").
				Attr("onclick", `
					const code = document.getElementById('`+snippetID+`').textContent;
					navigator.clipboard.writeText(code).then(() => {
						const btn = event.target.closest('button');
						const originalHTML = btn.innerHTML;
						btn.innerHTML = '<i class="bi bi-check"></i> Copied!';
						btn.classList.remove('btn-outline-primary');
						btn.classList.add('btn-success');
						setTimeout(() => {
							btn.innerHTML = originalHTML;
							btn.classList.remove('btn-success');
							btn.classList.add('btn-outline-primary');
						}, 2000);
					});
				`))).
		Child(hb.Div().
			Class("card-body p-0").
			Child(hb.PRE().
				Class("mb-0").
				Style("background-color: #f8f9fa; padding: 1rem; border-radius: 0;").
				Child(hb.Code().
					ID(snippetID).
					Text(code))))
}
