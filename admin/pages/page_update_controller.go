package admin

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/dracory/api"
	"github.com/dracory/blockeditor"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const VIEW_SETTINGS = "settings"
const VIEW_CONTENT = "content"
const VIEW_MIDDLEWARES = "middlewares"
const VIEW_SEO = "seo"
const ACTION_BLOCKEDITOR_HANDLE = "blockeditor_handle"
const ACTION_VERSION_HISTORY_SHOW = "action_version_history_show"
const ACTION_MIDDLEWARES_BEFORE_REPEATER_ADD = "action_middlewares_before_repeater_add"
const ACTION_MIDDLEWARES_BEFORE_REPEATER_DELETE = "action_middlewares_before_repeater_delete"
const ACTION_MIDDLEWARES_BEFORE_REPEATER_MOVE_UP = "action_middlewares_before_repeater_move_up"
const ACTION_MIDDLEWARES_BEFORE_REPEATER_MOVE_DOWN = "action_middlewares_before_repeater_move_down"
const ACTION_MIDDLEWARES_AFTER_REPEATER_ADD = "action_middlewares_after_repeater_add"
const ACTION_MIDDLEWARES_AFTER_REPEATER_DELETE = "action_middlewares_after_repeater_delete"
const ACTION_MIDDLEWARES_AFTER_REPEATER_MOVE_UP = "action_middlewares_after_repeater_move_up"
const ACTION_MIDDLEWARES_AFTER_REPEATER_MOVE_DOWN = "action_middlewares_after_repeater_move_down"
const ACTION_MIDDLEWARES_REPEATER_ADD = "action_middlewares_repeater_add"

// == CONTROLLER ==============================================================

type pageUpdateController struct {
	ui UiInterface
}

// == CONSTRUCTOR =============================================================

func NewPageUpdateController(ui UiInterface) *pageUpdateController {
	return &pageUpdateController{
		ui: ui,
	}
}

func (controller *pageUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		//return helpers.ToFlashError(w, r, errorMessage, shared.NewLinks().Pages(map[string]string{}), 10)
		return api.Error(errorMessage).ToString()
	}

	if data.action == ACTION_BLOCKEDITOR_HANDLE {
		return blockeditor.Handle(w, r, controller.ui.BlockEditorDefinitions())
	}

	if data.action == ACTION_VERSION_HISTORY_SHOW {
	}

	if data.action == ACTION_MIDDLEWARES_BEFORE_REPEATER_ADD {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_MIDDLEWARES_BEFORE_REPEATER_DELETE {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_MIDDLEWARES_BEFORE_REPEATER_MOVE_UP {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_MIDDLEWARES_BEFORE_REPEATER_MOVE_DOWN {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_MIDDLEWARES_AFTER_REPEATER_ADD {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_MIDDLEWARES_AFTER_REPEATER_DELETE {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_MIDDLEWARES_AFTER_REPEATER_MOVE_UP {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_MIDDLEWARES_AFTER_REPEATER_MOVE_DOWN {
		return controller.form(data).ToHTML()
	}

	if r.Method == http.MethodPost {
		return controller.form(data).ToHTML()
	}

	html := controller.page(data)

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

	options := struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		StyleURLs: []string{
			codemirrorCss,
			cdn.TrumbowygCss_2_27_3(),
		},
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Sweetalert2_11(),
			cdn.Jquery_3_7_1(),
			cdn.TrumbowygJs_2_27_3(),
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
		Styles: []string{
			`.CodeMirror {
				border: 1px solid #eee;
				height: auto;
			}
			`,
		},
		Scripts: []string{
			controller.script(),
		},
	}

	return controller.ui.Layout(w, r, "Edit page | CMS", html.ToHTML(), options)
}

func (controller pageUpdateController) script() string {
	js := ``
	return js
}

func (controller pageUpdateController) page(data pageUpdateControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), data.request)

	breadcrumbs := shared.AdminBreadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Page Manager",
			URL:  shared.URLR(data.request, shared.PathPagesPageManager, nil),
		},
		{
			Name: "Edit Page",
			URL:  shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{"page_id": data.pageID}),
		},
	}, struct{ SiteList []cmsstore.SiteInterface }{
		SiteList: data.siteList,
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormpageUpdate").
		HxPost(shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{"page_id": data.pageID})).
		HxTarget("#FormpageUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URLR(data.request, shared.PathPagesPageManager, nil))

	badgeStatus := hb.Div().
		Class("badge fs-6 ms-3").
		ClassIf(data.page.Status() == cmsstore.PAGE_STATUS_ACTIVE, "bg-success").
		ClassIf(data.page.Status() == cmsstore.PAGE_STATUS_INACTIVE, "bg-secondary").
		ClassIf(data.page.Status() == cmsstore.PAGE_STATUS_DRAFT, "bg-warning").
		Text(data.page.Status())

	buttonVersion := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-code-slash").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Version History").
		HxGet(shared.URLR(data.request, shared.PathPagesPageVersioning, map[string]string{
			"page_id": data.pageID,
			"action":  ACTION_VERSION_HISTORY_SHOW,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	pageTitle := hb.Heading1().
		Text("Edit Page:").
		Text(" ").
		Text(data.page.Name()).
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
					HTMLIf(data.view == VIEW_CONTENT, "Page Contents").
					HTMLIf(data.view == VIEW_SEO, "Page SEO").
					HTMLIf(data.view == VIEW_MIDDLEWARES, "Page Middlewares").
					HTMLIf(data.view == VIEW_SETTINGS, "Page Settings").
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
				Href(shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
					"page_id": data.pageID,
					"view":    VIEW_CONTENT,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SEO, "active").
				Href(shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
					"page_id": data.pageID,
					"view":    VIEW_SEO,
				})).
				HTML("SEO"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_MIDDLEWARES, "active").
				Href(shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
					"page_id": data.pageID,
					"view":    VIEW_MIDDLEWARES,
				})).
				HTML("Middlewares"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SETTINGS, "active").
				Href(shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
					"page_id": data.pageID,
					"view":    VIEW_SETTINGS,
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
		Child(card)
}

func (controller pageUpdateController) form(data pageUpdateControllerData) hb.TagInterface {
	fieldsSettings := controller.fieldsSettings(data)

	fieldsContent, errorMessage := controller.fieldsContent(data)

	if errorMessage != "" {
		hb.Div().Class("alert alert-danger").Text(errorMessage)
	}

	fieldsSEO := controller.fieldsSEO(data)

	formpageUpdate := form.NewForm(form.FormOptions{
		ID: "FormpageUpdate",
	})

	if data.view == VIEW_CONTENT {
		formpageUpdate.SetFields(fieldsContent)
	}

	if data.view == VIEW_MIDDLEWARES {
		formpageUpdate.SetFields(controller.fieldsMiddlewares(data))
	}

	if data.view == VIEW_SEO {
		formpageUpdate.SetFields(fieldsSEO)
	}

	if data.view == VIEW_SETTINGS {
		formpageUpdate.SetFields(fieldsSettings)
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

	// required := hb.Sup().HTML("required").Style("color:red;margin-left:10px;")

	// // Status
	// fomrGroupStatus := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Status").Child(required)).
	// 	Child(bs.FormSelect().
	// 		Name("page_status").
	// 		Child(bs.FormSelectOption("", "").
	// 			AttrIf(data.formStatus == "", "selected", "")).
	// 		Child(bs.FormSelectOption(blogstore.page_STATUS_DRAFT, "Draft").
	// 			AttrIf(data.formStatus == blogstore.page_STATUS_DRAFT, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.page_STATUS_PUBLISHED, "Published").
	// 			AttrIf(data.formStatus == blogstore.page_STATUS_PUBLISHED, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.page_STATUS_UNPUBLISHED, "Unpublished").
	// 			AttrIf(data.formStatus == blogstore.page_STATUS_UNPUBLISHED, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.page_STATUS_TRASH, "Trashed").
	// 			AttrIf(data.formStatus == blogstore.page_STATUS_TRASH, "selected", "selected")),
	// 	)

	// // Admin Notes
	// formGroupMemo := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Admin Notes")).
	// 	Child(bs.FormTextArea().
	// 		Name("page_memo").
	// 		HTML(data.formMemo).
	// 		Style("height:100px;"),
	// 	)

	// // page ID
	// formGrouppageId := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("page ID")).
	// 	Child(bs.FormInput().
	// 		Type(hb.TYPE_TEXT).
	// 		Name("page_id").
	// 		Value(data.pageID).
	// 		Attr("readonly", ""),
	// 	)

	// // Title
	// formGroupTitle := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Title").Child(required)).
	// 	Child(bs.FormInput().
	// 		Type("text").
	// 		Name("page_title").
	// 		Value(data.formTitle).
	// 		Style("width:100%;"),
	// 	)

	// // Summary
	// formGroupSummary := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Summary")).
	// 	Child(bs.FormTextArea().
	// 		Type("text").
	// 		Name("page_summary").
	// 		HTML(data.formSummary).
	// 		Style("width:100%;"),
	// 	)

	// // Published Date
	// formGroupPublishedAt := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Published Date")).
	// 	Child(bs.FormInput().
	// 		Type(hb.TYPE_TEXT).
	// 		Name("page_published_at").
	// 		Value(data.formPublishedAt).
	// 		Style("width:100%;"),
	// 	)

	// // Featured
	// formGroupFeatured := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Featured")).
	// 	Child(bs.FormSelect().
	// 		Name("page_featured").
	// 		Child(bs.FormSelectOption("", "").
	// 			AttrIf(data.formFeatured == "", "selected", "")).
	// 		Child(bs.FormSelectOption("yes", "Yes").
	// 			AttrIf(data.formFeatured == "yes", "selected", "selected")).
	// 		Child(bs.FormSelectOption("no", "No").
	// 			AttrIf(data.formFeatured == "no", "selected", "selected")),
	// 	)

	// form := hb.Form().
	// 	ID("FormpageUpdate").
	// 	Child(formGroupTitle).
	// 	Child(fomrGroupStatus).
	// 	Child(formGroupSummary).
	// 	Child(formGroupPublishedAt).
	// 	Child(formGroupFeatured).
	// 	Child(formGroupMemo).
	// 	Child(formGrouppageId)

	// if data.formErrorMessage != "" {
	// 	form.Child(hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}))
	// }

	// if data.formSuccessMessage != "" {
	// 	form.Child(hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}))
	// }

	// return form
}

func (pageUpdateController) fieldsSEO(data pageUpdateControllerData) []form.FieldInterface {
	fieldsSEO := []form.FieldInterface{
		&form.Field{
			Label: "Alias / Path / User Friendly URL",
			Name:  "page_alias",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formAlias,
			Help:  "The relative path on the website where this page will be visible to the vistors. Once set do not change it as search engines will look for this path.",
		},
		&form.Field{
			Label: "Meta Description",
			Name:  "page_meta_description",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formMetaDescription,
			Help:  "The description of this webpage as will be seen in search engines.",
		},
		&form.Field{
			Label: "Meta Keywords",
			Name:  "page_meta_keywords",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formMetaKeywords,
			Help:  "Specifies the keywords that will be used by the search engines to find this webpage. Separate keywords with commas.",
		},
		&form.Field{
			Label: "Meta Robots",
			Name:  "page_meta_robots",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formMetaRobots,
			Help:  "Specifies if this webpage should be indexed by the search engines. Index, Follow, means all. NoIndex, NoFollow means none.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "INDEX, FOLLOW",
					Key:   "INDEX, FOLLOW",
				},
				{
					Value: "NOINDEX, FOLLOW",
					Key:   "NOINDEX, FOLLOW",
				},
				{
					Value: "INDEX, NOFOLLOW",
					Key:   "INDEX, NOFOLLOW",
				},
				{
					Value: "NOINDEX, NOFOLLOW",
					Key:   "NOINDEX, NOFOLLOW",
				},
			},
		},
		&form.Field{
			Label: "Canonical URL",
			Name:  "page_canonical_url",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formCanonicalURL,
			Help:  "The canonical URL for this webpage. This is used by the search engines to display the preferred version of the web page in search results.",
		},
		&form.Field{
			Label:    "Webpage ID",
			Name:     "page_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.pageID,
			Readonly: true,
		},
		&form.Field{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_SEO,
			Readonly: true,
		},
	}
	return fieldsSEO
}

func (c pageUpdateController) fieldsContent(data pageUpdateControllerData) (fields []form.FieldInterface, errorMessage string) {
	editor := lo.IfF(data.page != nil, func() string { return data.page.Editor() }).Else("")

	fieldContent := form.Field{
		Label:   "Content",
		Name:    "page_content",
		Type:    form.FORM_FIELD_TYPE_TEXTAREA,
		Value:   data.formContent,
		Help:    "The content of this webpage. This will be displayed in the browser. If template is selected, the content will be displayed inside the template.",
		Options: []form.FieldOption{},
	}

	if editor == cmsstore.PAGE_EDITOR_BLOCKAREA {
		//fieldContent.Type = form.FORM_FIELD_TYPE_CODEMIRROR
		fieldContent.Options = []form.FieldOption{}
	}

	// For HTML Area editor, configure the Trumbowyg editor
	if editor == cmsstore.PAGE_EDITOR_HTMLAREA {
		htmlAreaFieldOptions := []form.FieldOption{
			{
				Key: "config",
				Value: `{
	btns: [
		['viewHTML'],
		['undo', 'redo'],
		['formatting'],
		['strong', 'em', 'del'],
		['superscript', 'subscript'],
		['link','justifyLeft','justifyRight','justifyCenter','justifyFull'],
		['unorderedList', 'orderedList'],
		['insertImage'],
		['removeformat'],
		['horizontalRule'],
		['fullscreen'],
	],
	autogrow: true,
	removeformatPasted: true,
	tagsToRemove: ['script', 'link', 'embed', 'iframe', 'input'],
	tagsToKeep: ['hr', 'img', 'i'],
	autogrowOnEnter: true,
	linkTargets: ['_blank'],
	}`,
			}}
		fieldContent.Type = form.FORM_FIELD_TYPE_HTMLAREA
		fieldContent.Options = htmlAreaFieldOptions
	}

	if editor == cmsstore.PAGE_EDITOR_BLOCKEDITOR {
		value := fieldContent.Value

		if value == "" {
			value = `[]`
		}

		editor, err := blockeditor.NewEditor(blockeditor.NewEditorOptions{
			// ID:    "blockeditor" + uid.HumanUid(),
			Name:  fieldContent.Name,
			Value: value,
			HandleEndpoint: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
				"page_id": data.pageID,
				"action":  ACTION_BLOCKEDITOR_HANDLE,
			}),
			BlockDefinitions: c.ui.BlockEditorDefinitions(),
		})

		if err != nil {
			return nil, "Error creating blockeditor: " + err.Error()
		}

		fieldContent.Type = form.FORM_FIELD_TYPE_BLOCKEDITOR
		fieldContent.CustomInput = editor
	}

	fieldsContent := []form.FieldInterface{
		&form.Field{
			Label: "Title",
			Name:  "page_title",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formTitle,
			Help:  "The title of this webpage. This will be displayed in the browser tab",
		},
		&fieldContent,
		&form.Field{
			Label:    "page ID",
			Name:     "page_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.pageID,
			Readonly: true,
		},
		&form.Field{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_CONTENT,
			Readonly: true,
		},
	}

	if editor == cmsstore.PAGE_EDITOR_MARKDOWN {
		contentScript := hb.Script(`
setTimeout(() => {
	const textArea = document.querySelector('textarea[name="page_content"]');
	textArea.style.height = '300px';
}, 2000)
			`).
			ToHTML()

		fieldsContent = append(fieldsContent, &form.Field{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		})
	}

	if editor == cmsstore.PAGE_EDITOR_CODEMIRROR {
		contentScript := hb.Script(`
function codeMirrorSelector() {
	return 'textarea[name="page_content"]';
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
	}

	return fieldsContent, ""
}

func (c pageUpdateController) fieldsMiddlewares(data pageUpdateControllerData) []form.FieldInterface {
	fieldMiddlewaresBefore := form.NewRepeater(form.RepeaterOptions{
		Label: "Middlewares Executed Before Page Content",
		Name:  "page_middlewares_before",
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Middleware",
				Name:  "page_middleware",
				Type:  form.FORM_FIELD_TYPE_STRING,
			}),
		},
		Values: lo.Map(data.formMiddlewaresBefore, func(middlewareName string, _ int) map[string]string {
			return map[string]string{
				"page_middleware": middlewareName,
			}
		}),
		RepeaterAddUrl: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
			"page_id": data.pageID,
			"view":    VIEW_MIDDLEWARES,
			"action":  ACTION_MIDDLEWARES_BEFORE_REPEATER_ADD,
		}),
		RepeaterRemoveUrl: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
			"page_id": data.pageID,
			"view":    VIEW_MIDDLEWARES,
			"action":  ACTION_MIDDLEWARES_BEFORE_REPEATER_DELETE,
		}),
		RepeaterMoveUpUrl: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
			"page_id": data.pageID,
			"view":    VIEW_MIDDLEWARES,
			"action":  ACTION_MIDDLEWARES_BEFORE_REPEATER_MOVE_UP,
		}),
		RepeaterMoveDownUrl: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
			"page_id": data.pageID,
			"view":    VIEW_MIDDLEWARES,
			"action":  ACTION_MIDDLEWARES_BEFORE_REPEATER_MOVE_DOWN,
		}),
	})

	formMiddlewaresAfter := form.NewRepeater(form.RepeaterOptions{
		Label: "Middlewares Executed After Page Content",
		Name:  "page_middlewares_after",
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Middleware",
				Name:  "page_middleware",
				Type:  form.FORM_FIELD_TYPE_STRING,
			}),
		},
		Values: lo.Map(data.formMiddlewaresAfter, func(middlewareName string, _ int) map[string]string {
			return map[string]string{
				"page_middleware": middlewareName,
			}
		}),
		RepeaterAddUrl: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
			"page_id": data.pageID,
			"view":    VIEW_MIDDLEWARES,
			"action":  ACTION_MIDDLEWARES_AFTER_REPEATER_ADD,
		}),
		RepeaterRemoveUrl: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
			"page_id": data.pageID,
			"view":    VIEW_MIDDLEWARES,
			"action":  ACTION_MIDDLEWARES_AFTER_REPEATER_DELETE,
		}),
		RepeaterMoveUpUrl: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
			"page_id": data.pageID,
			"view":    VIEW_MIDDLEWARES,
			"action":  ACTION_MIDDLEWARES_AFTER_REPEATER_MOVE_UP,
		}),
		RepeaterMoveDownUrl: shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
			"page_id": data.pageID,
			"view":    VIEW_MIDDLEWARES,
			"action":  ACTION_MIDDLEWARES_AFTER_REPEATER_MOVE_DOWN,
		}),
	})

	fieldView := &form.Field{
		Label:    "View",
		Name:     "view",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    VIEW_MIDDLEWARES,
		Readonly: true,
	}

	fieldSeparator := &form.Field{
		Label: "Separator",
		Name:  "separator",
		Type:  form.FORM_FIELD_TYPE_RAW,
		Value: hb.HR().ToHTML(),
	}

	fieldInfo := &form.Field{
		Label: "Info",
		Name:  "info",
		Type:  form.FORM_FIELD_TYPE_RAW,
		Value: hb.Div().
			Class(`alert alert-info`).
			Text(`You can add middlewares here that will be executed before or after the page content.`).
			Text(`Before middlewares can be used i.e. to authenticate the user.`).
			Text(`After middlewares can be used to change the page content after it has been generated.`).
			ToHTML(),
	}

	return []form.FieldInterface{
		fieldInfo,
		fieldMiddlewaresBefore,
		fieldSeparator,
		formMiddlewaresAfter,
		fieldView, // required
	}
}

func (c pageUpdateController) fieldsSettings(data pageUpdateControllerData) []form.FieldInterface {
	fieldEditor := &form.Field{
		Label: "Editor",
		Name:  "page_editor",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formEditor,
		Help:  "The content editor that will be used while editing this webpage content. Once set, this should not be changed, or the content may be lost. If left empty, the default editor (textarea) will be used. Note you will need to save and refresh to activate",
		OptionsF: func() []form.FieldOption {
			options := []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
			}

			options = append(options, form.FieldOption{
				Value: "CodeMirror (HTML Source Editor)",
				Key:   cmsstore.PAGE_EDITOR_CODEMIRROR,
			})

			if len(c.ui.BlockEditorDefinitions()) > 0 {
				options = append(options, form.FieldOption{
					Value: "BlockEditor (Visual Editor using Blocks)",
					Key:   cmsstore.PAGE_EDITOR_BLOCKEDITOR,
				})
			}

			options = append(options, form.FieldOption{
				Value: "Markdown (Simple Textarea)",
				Key:   cmsstore.PAGE_EDITOR_MARKDOWN,
			})

			options = append(options, form.FieldOption{
				Value: "HTML Area (WYSIWYG)",
				Key:   cmsstore.PAGE_EDITOR_HTMLAREA,
			})

			options = append(options, form.FieldOption{
				Value: "Text Area",
				Key:   cmsstore.PAGE_EDITOR_TEXTAREA,
			})

			return options
		},
	}

	fieldMemo := form.NewField(form.FieldOptions{
		Label: "Admin Notes (Internal)",
		Name:  "page_memo",
		Type:  form.FORM_FIELD_TYPE_TEXTAREA,
		Value: data.formMemo,
		Help:  "Admin notes for this page. These notes will not be visible to the public.",
	})

	fieldPageID := &form.Field{
		Label:    "Page Reference / ID",
		Name:     "page_id",
		Type:     form.FORM_FIELD_TYPE_STRING,
		Value:    data.pageID,
		Readonly: true,
		Help:     "The reference number (ID) of the page. This is used to identify the page in the system and should not be changed.",
	}

	fieldPageName := &form.Field{
		Label: "Page Name",
		Name:  "page_name",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formName,
		Help:  "The name of the page as displayed in the admin panel. This is not vsible to the page vistors",
	}

	fieldSiteID := &form.Field{
		Label: "Belongs to Site",
		Name:  "page_site_id",
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

	fieldStatus := &form.Field{
		Label: "Status",
		Name:  "page_status",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formStatus,
		Help:  "The status of this webpage. Published pages will be displayed on the website.",
		Options: []form.FieldOption{
			{
				Value: "- not selected -",
				Key:   "",
			},
			{
				Value: "Draft",
				Key:   cmsstore.PAGE_STATUS_DRAFT,
			},
			{
				Value: "Published",
				Key:   cmsstore.PAGE_STATUS_ACTIVE,
			},
			{
				Value: "Unpublished",
				Key:   cmsstore.PAGE_STATUS_INACTIVE,
			},
		},
	}

	fieldTemplateID := &form.Field{
		Label: "Template ID",
		Name:  "page_template_id",
		Type:  form.FORM_FIELD_TYPE_SELECT,
		Value: data.formTemplateID,
		Help:  "The template that this page content will be displayed in. This feature is useful if you want to implement consistent layouts. Leaving the template empty will display the page content as it is, standalone",
		OptionsF: func() []form.FieldOption {
			options := []form.FieldOption{
				{
					Value: "- not template selected, page content will be displayed as it is -",
					Key:   "",
				},
			}
			for _, template := range data.templateList {
				name := template.Name()
				options = append(options, form.FieldOption{
					Value: name,
					Key:   template.ID(),
				})
			}
			return options

		},
	}

	fieldView := &form.Field{
		Label:    "View",
		Name:     "view",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.view,
		Readonly: true,
	}

	fieldsSettings := []form.FieldInterface{
		fieldStatus,
		fieldTemplateID,
		fieldEditor,
		fieldPageName,
		fieldSiteID,
		fieldMemo,
		fieldPageID,
		fieldView, // required
	}

	return fieldsSettings
}

func (controller pageUpdateController) savePage(r *http.Request, data pageUpdateControllerData) (d pageUpdateControllerData, errorMessage string) {
	data.formAlias = req.GetStringTrimmed(r, "page_alias")
	data.formCanonicalURL = req.GetStringTrimmed(r, "page_canonical_url")
	data.formContent = req.GetStringTrimmed(r, "page_content")
	data.formEditor = req.GetStringTrimmed(r, "page_editor")
	data.formMemo = req.GetStringTrimmed(r, "page_memo")
	data.formMetaDescription = req.GetStringTrimmed(r, "page_meta_description")
	data.formMetaKeywords = req.GetStringTrimmed(r, "page_meta_keywords")
	data.formMetaRobots = req.GetStringTrimmed(r, "page_meta_robots")
	data.formName = req.GetStringTrimmed(r, "page_name")
	data.formSummary = req.GetStringTrimmed(r, "page_summary")
	data.formStatus = req.GetStringTrimmed(r, "page_status")
	data.formSiteID = req.GetStringTrimmed(r, "page_site_id")
	data.formTitle = req.GetStringTrimmed(r, "page_title")
	data.formTemplateID = req.GetStringTrimmed(r, "page_template_id")
	data.formMiddlewaresAfter = controller.requestMapToMiddlewaresAfter(r)
	data.formMiddlewaresBefore = controller.requestMapToMiddlewaresBefore(r)

	if data.view == VIEW_SETTINGS {
		if data.formStatus == "" {
			data.formErrorMessage = "Status is required"
			return data, ""
		}
	}

	if data.view == VIEW_CONTENT {
		if data.formTitle == "" {
			data.formErrorMessage = "Title is required"
			return data, ""
		}
	}

	if data.view == VIEW_MIDDLEWARES {
		data.page.SetMiddlewaresAfter(data.formMiddlewaresAfter)
		data.page.SetMiddlewaresBefore(data.formMiddlewaresBefore)
	}

	if data.view == VIEW_SETTINGS {
		// make sure the date is in the correct format
		// data.formPublishedAt = lo.Substring(strings.ReplaceAll(data.formPublishedAt, " ", "T")+":00", 0, 19)
		// publishedAt := lo.Ternary(data.formPublishedAt == "", sb.NULL_DATE, carbon.Parse(data.formPublishedAt).ToDateTimeString(carbon.UTC))
		data.page.SetEditor(data.formEditor)
		data.page.SetMemo(data.formMemo)
		data.page.SetName(data.formName)
		data.page.SetSiteID(data.formSiteID)
		data.page.SetStatus(data.formStatus)
		data.page.SetTemplateID(data.formTemplateID)
	}

	if data.view == VIEW_CONTENT {
		data.page.SetContent(data.formContent)
		data.page.SetTitle(data.formTitle)
	}

	if data.view == VIEW_SEO {
		data.page.SetAlias(data.formAlias)
		data.page.SetCanonicalUrl(data.formCanonicalURL)
		data.page.SetMetaDescription(data.formMetaDescription)
		data.page.SetMetaKeywords(data.formMetaKeywords)
		data.page.SetMetaRobots(data.formMetaRobots)
	}

	// versioning is handled by the store automatically

	err := controller.ui.Store().PageUpdate(data.request.Context(), data.page)

	if err != nil {
		controller.ui.Logger().Error("At pageUpdateController > prepareDataAndValidate", "error", err.Error())
		data.formErrorMessage = "System error. Saving page failed. " + err.Error()
		return data, ""
	}

	err = controller.movePageBlocks(r, data.page.ID(), data.page.SiteID())

	if err != nil {
		controller.ui.Logger().Error("At pageUpdateController > prepareDataAndValidate > movePageBlocks", "error", err.Error())
		data.formErrorMessage = "System error. Saving page failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "page saved successfully"

	data.formRedirectURL = shared.URLR(data.request, shared.PathPagesPageUpdate, map[string]string{
		"page_id": data.pageID,
		"view":    data.view,
	})

	return data, ""
}

func (controller pageUpdateController) createVersioning(ctx context.Context, page cmsstore.PageInterface) error {
	if !controller.ui.Store().VersioningEnabled() {
		return nil
	}

	if page == nil {
		return errors.New("page is nil")
	}

	lastVersioningList, err := controller.ui.Store().VersioningList(ctx, cmsstore.NewVersioningQuery().
		SetEntityType(cmsstore.VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder(sb.DESC).
		SetLimit(1))

	if err != nil {
		return err
	}

	content, err := page.MarshalToVersioning()

	if err != nil {
		return err
	}

	if controller.isLastVersioningSame(content, lastVersioningList) {
		return nil // nothing to do
	}

	entityID := page.ID()

	return controller.ui.Store().VersioningCreate(ctx, cmsstore.NewVersioning().
		SetEntityID(entityID).
		SetEntityType(cmsstore.VERSIONING_TYPE_PAGE).
		SetContent(content))
}

func (controller pageUpdateController) isLastVersioningSame(
	pageVersioningContent string,
	lastVersioningList []cmsstore.VersioningInterface,
) bool {
	lastVersioning := lo.IfF[cmsstore.VersioningInterface](len(lastVersioningList) > 0, func() cmsstore.VersioningInterface {
		return lastVersioningList[0]
	}).ElseF(func() cmsstore.VersioningInterface {
		return nil
	})

	if lastVersioning == nil {
		return false
	}

	lastVersioningContent := lastVersioning.Content()

	if lastVersioningContent == pageVersioningContent {
		cfmt.Infoln("No changes detected")
		return true
	}

	cfmt.Infoln("Changes detected")

	// cfmt.Infoln("last versioning content", lastVersioningContent)
	// cfmt.Warningln("new versioning content", pageVersioningContent)

	return false
}

// movePageBlocks moves all blocks from the current site to the new site
// if the page is moved to a different site
func (controller pageUpdateController) movePageBlocks(request *http.Request, pageID string, siteID string) error {
	blocks, err := controller.ui.Store().BlockList(request.Context(), cmsstore.BlockQuery().
		SetPageID(pageID))

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

// prepareDataAndValidate prepares the data and validates it
//
// Business Logic:
// - checks if the page exists
// - checks if the view is valid, and sets the default if not provided
// - retrieves the site list
// - retrieves the template list
// - if its a GET request, returns the data, (form data is from the database)
// - if its a POST request, saves the page and returns the data (form data is from the POST request)
//
// Parameters:
// - r *http.Request - the HTTP request object
//
// Returns:
// - data pageUpdateControllerData - the data for the current controller request
// - errorMessage string - the error message, or emty string if no error
func (controller pageUpdateController) prepareDataAndValidate(r *http.Request) (data pageUpdateControllerData, errorMessage string) {
	data.request = r
	data.action = req.GetStringTrimmed(r, "action")
	data.pageID = req.GetStringTrimmed(r, "page_id")
	data.view = req.GetStringTrimmedOr(r, "view", VIEW_CONTENT)

	if data.view == "" {
		data.view = VIEW_CONTENT
	}

	if data.pageID == "" {
		return data, "page ID is required"
	}

	var err error
	data.page, err = controller.ui.Store().PageFindByID(r.Context(), data.pageID)

	if err != nil {
		return data, err.Error()

	}

	if data.page == nil {
		return data, "page not found"
	}

	data.formAlias = data.page.Alias()
	data.formCanonicalURL = data.page.CanonicalUrl()
	data.formContent = data.page.Content()
	data.formEditor = data.page.Editor()
	data.formMetaDescription = data.page.MetaDescription()
	data.formMetaKeywords = data.page.MetaKeywords()
	data.formMetaRobots = data.page.MetaRobots()
	data.formName = data.page.Name()
	data.formMemo = data.page.Memo()
	data.formSiteID = data.page.SiteID()
	data.formStatus = data.page.Status()
	data.formTemplateID = data.page.TemplateID()
	data.formTitle = data.page.Title()
	data.formMiddlewaresAfter = data.page.MiddlewaresAfter()
	data.formMiddlewaresBefore = data.page.MiddlewaresBefore()

	data.siteList, err = controller.ui.Store().SiteList(data.request.Context(), cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		return data, "Site list failed to be retrieved" + err.Error()
	}

	templateList, err := controller.ui.Store().TemplateList(data.request.Context(), cmsstore.TemplateQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		return data, "Template list failed to be retrieved" + err.Error()
	}

	data.templateList = templateList

	if r.Method != http.MethodPost {
		return data, ""
	}

	data.formMiddlewaresAfter = controller.requestMapToMiddlewaresAfter(r)
	data.formMiddlewaresBefore = controller.requestMapToMiddlewaresBefore(r)

	if data.action == ACTION_MIDDLEWARES_AFTER_REPEATER_ADD {
		data.formMiddlewaresAfter = append(data.formMiddlewaresAfter, "")
		return data, ""
	}

	if data.action == ACTION_MIDDLEWARES_AFTER_REPEATER_DELETE {
		repeatableRemoveIndex := req.GetStringTrimmed(r, "repeatable_remove_index")

		if repeatableRemoveIndex == "" {
			return data, ""
		}

		data.formMiddlewaresAfter = slices.Delete(data.formMiddlewaresAfter, cast.ToInt(repeatableRemoveIndex), cast.ToInt(repeatableRemoveIndex)+1)

		return data, ""
	}

	if data.action == ACTION_MIDDLEWARES_AFTER_REPEATER_MOVE_UP {
		repeatableMoveUpIndex := cast.ToInt(req.GetStringTrimmed(r, "repeatable_move_up_index"))

		if repeatableMoveUpIndex == 0 {
			return data, ""
		}

		current := data.formMiddlewaresAfter[repeatableMoveUpIndex]
		upper := data.formMiddlewaresAfter[repeatableMoveUpIndex-1]

		data.formMiddlewaresAfter[repeatableMoveUpIndex] = upper
		data.formMiddlewaresAfter[repeatableMoveUpIndex-1] = current

		return data, ""
	}

	if data.action == ACTION_MIDDLEWARES_AFTER_REPEATER_MOVE_DOWN {
		repeatableMoveDownIndex := cast.ToInt(req.GetStringTrimmed(r, "repeatable_move_down_index"))

		if repeatableMoveDownIndex == len(data.formMiddlewaresAfter)-1 {
			return data, ""
		}

		current := data.formMiddlewaresAfter[repeatableMoveDownIndex]
		lower := data.formMiddlewaresAfter[repeatableMoveDownIndex+1]

		data.formMiddlewaresAfter[repeatableMoveDownIndex] = lower
		data.formMiddlewaresAfter[repeatableMoveDownIndex+1] = current

		return data, ""
	}

	if data.action == ACTION_MIDDLEWARES_BEFORE_REPEATER_ADD {
		data.formMiddlewaresBefore = append(data.formMiddlewaresBefore, "")
		return data, ""
	}

	if data.action == ACTION_MIDDLEWARES_BEFORE_REPEATER_DELETE {
		repeatableRemoveIndex := req.GetStringTrimmed(r, "repeatable_remove_index")

		if repeatableRemoveIndex == "" {
			return data, ""
		}

		data.formMiddlewaresBefore = slices.Delete(data.formMiddlewaresBefore, cast.ToInt(repeatableRemoveIndex), cast.ToInt(repeatableRemoveIndex)+1)

		return data, ""
	}

	if data.action == ACTION_MIDDLEWARES_BEFORE_REPEATER_MOVE_UP {
		repeatableMoveUpIndex := cast.ToInt(req.GetStringTrimmed(r, "repeatable_move_up_index"))

		if repeatableMoveUpIndex == 0 {
			return data, ""
		}

		current := data.formMiddlewaresBefore[repeatableMoveUpIndex]
		upper := data.formMiddlewaresBefore[repeatableMoveUpIndex-1]

		data.formMiddlewaresBefore[repeatableMoveUpIndex] = upper
		data.formMiddlewaresBefore[repeatableMoveUpIndex-1] = current

		return data, ""
	}

	if data.action == ACTION_MIDDLEWARES_BEFORE_REPEATER_MOVE_DOWN {
		repeatableMoveDownIndex := cast.ToInt(req.GetStringTrimmed(r, "repeatable_move_down_index"))

		if repeatableMoveDownIndex == len(data.formMiddlewaresBefore)-1 {
			return data, ""
		}

		current := data.formMiddlewaresBefore[repeatableMoveDownIndex]
		lower := data.formMiddlewaresBefore[repeatableMoveDownIndex+1]

		data.formMiddlewaresBefore[repeatableMoveDownIndex] = lower
		data.formMiddlewaresBefore[repeatableMoveDownIndex+1] = current

		return data, ""
	}

	return controller.savePage(r, data)
}

// func (pageUpdateController) pageMiddlewaresFromMeta(page cmsstore.PageInterface) []string {
// 	meta := page.Meta("middlewares")

// 	if meta == "" {
// 		return []string{}
// 	}

// 	m, err := utils.FromJSON(page.Meta("middlewares"), []string{})

// 	if err != nil {
// 		cfmt.Error(err)
// 		return []string{}
// 	}

// 	return lo.Map(m.([]interface{}), func(v interface{}, _ int) string {
// 		return v.(string)
// 	})
// }

func (controller pageUpdateController) requestMapToMiddlewaresBefore(r *http.Request) []string {
	formMiddlewares := req.GetMaps(r, "page_middlewares_before", []map[string]string{})
	middlewares := []string{}

	for _, formMiddleware := range formMiddlewares {
		middlewares = append(middlewares, formMiddleware["page_middleware"])
	}

	return middlewares
}

func (controller pageUpdateController) requestMapToMiddlewaresAfter(r *http.Request) []string {
	formMiddlewares := req.GetMaps(r, "page_middlewares_after", []map[string]string{})
	middlewares := []string{}

	for _, formMiddleware := range formMiddlewares {
		middlewares = append(middlewares, formMiddleware["page_middleware"])
	}

	return middlewares
}

type pageUpdateControllerData struct {
	request *http.Request
	action  string
	pageID  string
	page    cmsstore.PageInterface
	view    string

	siteList     []cmsstore.SiteInterface
	templateList []cmsstore.TemplateInterface

	formErrorMessage      string
	formRedirectURL       string
	formSuccessMessage    string
	formAlias             string
	formCanonicalURL      string
	formContent           string
	formName              string
	formEditor            string
	formMemo              string
	formMetaDescription   string
	formMetaKeywords      string
	formMetaRobots        string
	formMiddlewaresAfter  []string
	formMiddlewaresBefore []string
	formSiteID            string
	formStatus            string
	formTemplateID        string
	formSummary           string
	formTitle             string
}
