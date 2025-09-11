package admin

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
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

type templateUpdateController struct {
	ui UiInterface
}

// == CONSTRUCTOR =============================================================

func NewTemplateUpdateController(ui UiInterface) *templateUpdateController {
	return &templateUpdateController{
		ui: ui,
	}
}

func (controller *templateUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
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

	return controller.ui.Layout(w, r, "Edit Template | CMS", html.ToHTML(), options)
}

func (controller templateUpdateController) page(data templateUpdateControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), data.request)

	breadcrumbs := shared.AdminBreadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Template Manager",
			URL:  shared.URLR(data.request, shared.PathTemplatesTemplateManager, nil),
		},
		{
			Name: "Edit Template",
			URL:  shared.URLR(data.request, shared.PathTemplatesTemplateUpdate, map[string]string{"template_id": data.templateID}),
		},
	}, struct{ SiteList []cmsstore.SiteInterface }{
		SiteList: data.siteList,
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormTemplateUpdate").
		HxPost(shared.URLR(data.request, shared.PathTemplatesTemplateUpdate, map[string]string{"template_id": data.templateID})).
		HxTarget("#FormTemplateUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URLR(data.request, shared.PathTemplatesTemplateManager, nil))

	badgeStatus := hb.Div().
		Class("badge fs-6 ms-3").
		ClassIf(data.template.Status() == cmsstore.TEMPLATE_STATUS_ACTIVE, "bg-success").
		ClassIf(data.template.Status() == cmsstore.TEMPLATE_STATUS_INACTIVE, "bg-secondary").
		ClassIf(data.template.Status() == cmsstore.TEMPLATE_STATUS_DRAFT, "bg-warning").
		Text(data.template.Status())

	pageTitle := hb.Heading1().
		Text("CMS. Edit Template:").
		Text(" ").
		Text(data.template.Name()).
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
					HTMLIf(data.view == VIEW_CONTENT, "Template Content").
					HTMLIf(data.view == VIEW_SETTINGS, "Template Settings").
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
				Href(shared.URLR(data.request, shared.PathTemplatesTemplateUpdate, map[string]string{
					"template_id": data.templateID,
					"view":        VIEW_CONTENT,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SETTINGS, "active").
				Href(shared.URLR(data.request, shared.PathTemplatesTemplateUpdate, map[string]string{
					"template_id": data.templateID,
					"view":        VIEW_SETTINGS,
				})).
				HTML("Settings")))

	toolsInfo := hb.NewParagraph().Class("alert alert-info").
		HTML("Tools: ").
		Child(hb.NewHyperlink().HTML("Google Translate").Href("https://translate.google.com").Target("_blank")).
		HTML(", ").
		Child(hb.NewHyperlink().HTML("Bing Translate").Href("https://www.bing.com/translator").Target("_blank")).
		HTML(", ").
		Child(hb.NewHyperlink().HTML("Translateking").Href("https://translateking.com/").Target("_blank")).
		HTML(", ").
		Child(hb.NewHyperlink().HTML("Baidu Translate").Href("https://fanyi.baidu.com/").Target("_blank")).
		HTML(", ").
		Child(hb.NewHyperlink().HTML("Yandex Translate").Href("https://translate.yandex.com").Target("_blank")).
		HTML(", ").
		Child(hb.NewHyperlink().HTML("Yandex Translate").Href("https://www.reverso.net/text-translation").Target("_blank"))

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader).
		Child(hb.HR()).
		Child(pageTitle).
		Child(toolsInfo).
		Child(tabs).
		Child(card)
}

func (controller templateUpdateController) form(data templateUpdateControllerData) hb.TagInterface {

	fieldsContent := controller.fieldsContent(data)
	fieldsSettings := controller.fieldsSettings(data)

	formpageUpdate := form.NewForm(form.FormOptions{
		ID: "FormTemplateUpdate",
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
				Timer:             5000,
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

func (templateUpdateController) fieldsContent(data templateUpdateControllerData) []form.FieldInterface {
	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Type: form.FORM_FIELD_TYPE_RAW,
			Value: hb.Div().
				Class(`alert alert-info`).
				Child(hb.Text("Available variables: [[PageContent]], [[PageCanonicalUrl]], [[PageMetaDescription]], [[PageMetaKeywords]], [[PageMetaRobots]], [[PageTitle]]")).
				ToHTML(),
		}),
		form.NewField(form.FieldOptions{
			Label: "Content (HTML)",
			Name:  "template_content",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formContent,
		}),
		form.NewField(form.FieldOptions{
			Label:    "Template ID",
			Name:     "template_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.templateID,
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
	return 'textarea[name="template_content"]';
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

func (controller templateUpdateController) fieldsSettings(data templateUpdateControllerData) []form.FieldInterface {
	fieldMemo := form.NewField(form.FieldOptions{
		Label: "Admin Notes (Internal)",
		Name:  "template_memo",
		Type:  form.FORM_FIELD_TYPE_TEXTAREA,
		Value: data.formMemo,
		Help:  "Admin notes for this template. These notes will not be visible to the public.",
	})

	fieldSiteID := &form.Field{
		Label: "Belongs to Site",
		Name:  "template_site_id",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formSiteID,
		Help:  "The site that this page belongs to",
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
		Name:  "template_status",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formStatus,
		Help:  "The status of this webpage. Published pages will be displayed on the webtemplate.",
		Options: []form.FieldOption{
			{
				Value: "- not selected -",
				Key:   "",
			},
			{
				Value: "Draft",
				Key:   cmsstore.TEMPLATE_STATUS_DRAFT,
			},
			{
				Value: "Published",
				Key:   cmsstore.TEMPLATE_STATUS_ACTIVE,
			},
			{
				Value: "Unpublished",
				Key:   cmsstore.TEMPLATE_STATUS_INACTIVE,
			},
		},
	})

	fieldTemplateID := form.NewField(form.FieldOptions{
		Label:    "Template ID",
		Name:     "template_id",
		Type:     form.FORM_FIELD_TYPE_STRING,
		Value:    data.templateID,
		Readonly: true,
		Help:     "The reference number (ID) of the template. This is used to identify the template in the system and should not be changed.",
	})

	fieldTemplateName := form.NewField(form.FieldOptions{
		Label: "Template Name (Internal)",
		Name:  "template_name",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formName,
		Help:  "The name of the template as displayed in the admin panel. This is not vsible to the template vistors",
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
		fieldTemplateName,
		fieldSiteID,
		fieldMemo,
		fieldTemplateID,
		fieldView,
	}

	return fieldsSettings
}

func (controller templateUpdateController) saveTemplate(data templateUpdateControllerData) (templateUpdateControllerData, string) {
	data.formContent = req.GetStringTrimmed(data.request, "template_content")
	data.formMemo = req.GetStringTrimmed(data.request, "template_memo")
	data.formName = req.GetStringTrimmed(data.request, "template_name")
	data.formSiteID = req.GetStringTrimmed(data.request, "template_site_id")
	data.formStatus = req.GetStringTrimmed(data.request, "template_status")
	data.formTitle = req.GetStringTrimmed(data.request, "template_title")

	if data.view == VIEW_SETTINGS {
		if data.formStatus == "" {
			data.formErrorMessage = "Status is required"
			return data, ""
		}
	}

	if data.view == VIEW_SETTINGS {
		data.template.SetMemo(data.formMemo)
		data.template.SetName(data.formName)
		data.template.SetSiteID(data.formSiteID)
		data.template.SetStatus(data.formStatus)
	}

	if data.view == VIEW_CONTENT {
		data.template.SetContent(data.formContent)
	}

	err := controller.ui.Store().TemplateUpdate(data.request.Context(), data.template)

	if err != nil {
		controller.ui.Logger().Error("At templateUpdateController > prepareDataAndValidate", "error", err.Error())
		data.formErrorMessage = "System error. Saving template failed. " + err.Error()
		return data, ""
	}

	err = controller.moveTemplateBlocks(data.request, data.template.ID(), data.formSiteID)

	if err != nil {
		controller.ui.Logger().Error("At templateUpdateController > prepareDataAndValidate", "error", err.Error())
		data.formErrorMessage = "System error. Saving template failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "template saved successfully"
	data.formRedirectURL = shared.URLR(data.request, shared.PathTemplatesTemplateUpdate, map[string]string{
		"template_id": data.template.ID(),
		"view":        data.view,
	})

	return data, ""
}

func (controller templateUpdateController) moveTemplateBlocks(request *http.Request, templateID string, siteID string) error {
	blocks, err := controller.ui.Store().BlockList(request.Context(), cmsstore.BlockQuery().
		SetPageID(templateID))

	if err != nil {
		return err
	}

	for _, block := range blocks {
		if block.SiteID() == siteID {
			continue // already in the right site
		}

		block.SetSiteID(siteID)

		err := controller.ui.Store().BlockUpdate(request.Context(), block)

		if err != nil {
			return err
		}
	}

	return nil
}

func (controller templateUpdateController) prepareDataAndValidate(r *http.Request) (data templateUpdateControllerData, errorMessage string) {
	data.request = r
	data.action = req.GetStringTrimmed(r, "action")
	data.templateID = req.GetStringTrimmed(r, "template_id")
	data.view = req.GetStringTrimmed(r, "view")

	if data.view == "" {
		data.view = VIEW_CONTENT
	}

	if data.templateID == "" {
		return data, "template id is required"
	}

	var err error
	data.template, err = controller.ui.Store().TemplateFindByID(r.Context(), data.templateID)

	if err != nil {
		controller.ui.Logger().Error("At templateUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if data.template == nil {
		return data, "template not found"
	}

	siteList, err := controller.ui.Store().SiteList(r.Context(), cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		return data, "Site list failed to be retrieved" + err.Error()
	}

	data.siteList = siteList

	data.formContent = data.template.Content()
	data.formName = data.template.Name()
	data.formMemo = data.template.Memo()
	data.formSiteID = data.template.SiteID()
	data.formStatus = data.template.Status()

	if r.Method != http.MethodPost {
		return data, ""
	}

	return controller.saveTemplate(data)
}

type templateUpdateControllerData struct {
	request    *http.Request
	action     string
	templateID string
	template   cmsstore.TemplateInterface
	view       string

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
}
