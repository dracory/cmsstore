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
	"github.com/samber/lo"
)

const VIEW_SETTINGS = "settings"
const VIEW_CONTENT = "content"

// == CONTROLLER ==============================================================

type translationUpdateController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*translationUpdateController)(nil)

// == CONSTRUCTOR =============================================================

func NewTranslationUpdateController(ui UiInterface) *translationUpdateController {
	return &translationUpdateController{
		ui: ui,
	}
}

func (controller *translationUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
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
		StyleURLs: []string{},
		Scripts:   []string{},
		ScriptURLs: []string{
			cdn.Sweetalert2_11(),
			cdn.Htmx_2_0_0(),
		},
	}

	return controller.ui.Layout(w, r, "Edit Translation | CMS", html.ToHTML(), options)
}

func (controller translationUpdateController) page(data translationUpdateControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), data.request)

	breadcrumbs := shared.AdminBreadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Translation Manager",
			URL:  shared.URL(shared.Endpoint(data.request), shared.PathTranslationsTranslationManager, nil),
		},
		{
			Name: "Edit Translation",
			URL:  shared.URL(shared.Endpoint(data.request), shared.PathTranslationsTranslationUpdate, map[string]string{"translation_id": data.translationID}),
		},
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormTranslationUpdate").
		HxPost(shared.URL(shared.Endpoint(data.request), shared.PathTranslationsTranslationUpdate, map[string]string{"translation_id": data.translationID})).
		HxTarget("#FormTranslationUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URL(shared.Endpoint(data.request), shared.PathTranslationsTranslationManager, nil))

	badgeStatus := hb.Div().
		Class("badge fs-6 ms-3").
		ClassIf(data.translation.Status() == cmsstore.TEMPLATE_STATUS_ACTIVE, "bg-success").
		ClassIf(data.translation.Status() == cmsstore.TEMPLATE_STATUS_INACTIVE, "bg-secondary").
		ClassIf(data.translation.Status() == cmsstore.TEMPLATE_STATUS_DRAFT, "bg-warning").
		Text(data.translation.Status())

	pageTitle := hb.Heading1().
		Text("CMS. Edit Translation:").
		Text(" ").
		Text(data.translation.Name()).
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
					HTMLIf(data.view == VIEW_CONTENT, "Translation Content").
					HTMLIf(data.view == VIEW_SETTINGS, "Translation Settings").
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
				Href(shared.URL(shared.Endpoint(data.request), shared.PathTranslationsTranslationUpdate, map[string]string{
					"translation_id": data.translationID,
					"view":           VIEW_CONTENT,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SETTINGS, "active").
				Href(shared.URL(shared.Endpoint(data.request), shared.PathTranslationsTranslationUpdate, map[string]string{
					"translation_id": data.translationID,
					"view":           VIEW_SETTINGS,
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
			Text("To use this translation in your website use the following shortcode:").
			Child(hb.BR())).
		Child(hb.PRE().
			Child(hb.Code().
				Text(`<!-- START: Translation: ` + data.translation.Name() + ` -->`).
				Text("\n").
				Text(`[[TRANSLATION_` + data.translationID + `]]`).
				Text("\n").
				Text(`<!-- END: Translation: ` + data.translation.Name() + ` -->`))).
		Child(hb.PRE().
			Child(hb.Code().
				Text(`<!-- START: Translation: ` + data.translation.Name() + ` -->`).
				Text("\n").
				Text(`[[TRANSLATION_` + data.translation.Handle() + `]]`).
				Text("\n").
				Text(`<!-- END: Translation: ` + data.translation.Name() + ` -->`)))
}

func (controller translationUpdateController) form(data translationUpdateControllerData) hb.TagInterface {

	fieldsContent := controller.fieldsContent(data)
	fieldsSettings := controller.fieldsSettings(data)

	formpageUpdate := form.NewForm(form.FormOptions{
		ID: "FormTranslationUpdate",
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

func (c translationUpdateController) fieldsContent(data translationUpdateControllerData) []form.FieldInterface {
	contentFields := []form.FieldInterface{}

	for languageCode, languageName := range c.ui.Store().TranslationLanguages() {
		isDefault := languageCode == c.ui.Store().TranslationLanguageDefault()
		contentField := form.NewField(form.FieldOptions{
			Label: languageName + " Translation " + lo.Ternary(isDefault, " (Default)", ""),
			Name:  "translation_content[" + languageCode + "]",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formContent[languageCode],
		})

		contentFields = append(contentFields, contentField)
	}

	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Key / Handle",
			Name:  "translation_handle",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formHandle,
			Help:  "Must be unique per site, lowercase, no spaces, no hyphens, and no dots allowed. This is used to identify this translation in human friendly fashion.",
		}),
	}

	fieldsContent = append(fieldsContent, contentFields...)

	fieldsContent = append(fieldsContent, []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label:    "Translation ID",
			Name:     "translation_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.translationID,
			Readonly: true,
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_CONTENT,
			Readonly: true,
		}),
	}...)

	contentScript := hb.Script(`
function codeMirrorSelector() {
	return 'textarea[name="translation_content"]';
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

func (controller translationUpdateController) fieldsSettings(data translationUpdateControllerData) []form.FieldInterface {
	fieldsSettings := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Status",
			Name:  "translation_status",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formStatus,
			Help:  "The status of this translation. Published translations will be visible on the website.",
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
		}),
		form.NewField(form.FieldOptions{
			Label: "Translation Name (Internal)",
			Name:  "translation_name",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formName,
			Help:  "The name of the translation as displayed in the admin panel. This is not vsible to the public.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Belongs to Site",
			Name:  "translation_site_id",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formSiteID,
			Help:  "The site to which this translation belongs to.",
			OptionsF: func() []form.FieldOption {
				options := []form.FieldOption{
					{
						Value: "- not site selected -",
						Key:   "",
					},
				}
				for _, site := range data.siteList {
					name := site.Name()
					status := site.Status()
					options = append(options, form.FieldOption{
						Value: name + ` (` + status + `)`,
						Key:   site.ID(),
					})
				}
				return options

			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Admin Notes (Internal)",
			Name:  "translation_memo",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formMemo,
			Help:  "Admin notes for this translation. These notes will not be visible to the public.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "Translation Reference (ID)",
			Name:     "translation_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.translationID,
			Readonly: true,
			Help:     "The reference number (ID) of the translation. This is used to identify the translation in the system and should not be changed.",
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

func (controller translationUpdateController) saveTranslation(r *http.Request, data translationUpdateControllerData) (d translationUpdateControllerData, errorMessage string) {
	data.formContent = utils.ReqMap(r, "translation_content")
	data.formMemo = utils.Req(r, "translation_memo", "")
	data.formName = utils.Req(r, "translation_name", "")
	data.formStatus = utils.Req(r, "translation_status", "")
	data.formHandle = utils.Req(r, "translation_handle", "")
	data.formSiteID = utils.Req(r, "translation_site_id", "")

	refreshPage := false

	if data.view == VIEW_CONTENT {
		if data.formHandle == "" {
			data.formErrorMessage = "Key is required"
			return data, ""
		}
	}

	if data.view == VIEW_SETTINGS {
		if data.formStatus == "" {
			data.formErrorMessage = "Status is required"
			return data, ""
		}

		if data.formSiteID == "" {
			data.formErrorMessage = "Site is required"
			return data, ""
		}
	}

	if data.view == VIEW_SETTINGS {
		data.translation.SetMemo(data.formMemo)

		if data.formName != data.translation.Name() {
			refreshPage = true // name has changed, must refersh the page
		}

		data.translation.SetName(data.formName)

		// if data.formStatus != data.translation.Status() {
		// 	refreshPage = true // status has changed, must refersh the page
		// }

		data.translation.SetSiteID(data.formSiteID)

		if data.formStatus != data.translation.Status() {
			refreshPage = true // status has changed, must refersh the page
		}

		data.translation.SetStatus(data.formStatus)
	}

	if data.view == VIEW_CONTENT {
		err := data.translation.SetContent(data.formContent)

		if err != nil {
			data.formErrorMessage = "System error. Saving translation failed. " + err.Error()
			return data, ""
		}

		if data.formHandle != data.translation.Handle() {
			refreshPage = true // handle has changed, must refersh the page
		}

		data.translation.SetHandle(data.formHandle)
	}

	err := controller.ui.Store().TranslationUpdate(data.translation)

	if err != nil {
		//config.LogStore.ErrorWithContext("At translationUpdateController > prepareDataAndValidate", err.Error())
		data.formErrorMessage = "System error. Saving translation failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "translation saved successfully"
	if refreshPage {
		data.formRedirectURL = shared.URL(shared.Endpoint(data.request), shared.PathTranslationsTranslationUpdate, map[string]string{
			"translation_id": data.translationID,
			"view":           data.view,
		})
	}

	return data, ""
}

func (controller translationUpdateController) prepareDataAndValidate(r *http.Request) (data translationUpdateControllerData, errorMessage string) {
	data.request = r
	data.action = utils.Req(r, "action", "")
	data.translationID = utils.Req(r, "translation_id", "")
	data.view = utils.Req(r, "view", "")

	if data.view == "" {
		data.view = VIEW_CONTENT
	}

	if data.translationID == "" {
		return data, "translation id is required"
	}

	// 1. Fetch required data

	var err error
	data.translation, err = controller.ui.Store().TranslationFindByID(data.translationID)

	if err != nil {
		controller.ui.Logger().Error("At translationUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if data.translation == nil {
		return data, "translation not found"
	}

	data.siteList, err = controller.ui.Store().SiteList(cmsstore.SiteQuery())

	if err != nil {
		controller.ui.Logger().Error("At translationUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	// 2. Populate form data

	data.formContent, err = data.translation.Content()

	if err != nil {
		controller.ui.Logger().Error("At translationUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.formHandle = data.translation.Handle()
	data.formName = data.translation.Name()
	data.formMemo = data.translation.Memo()
	data.formSiteID = data.translation.SiteID()
	data.formStatus = data.translation.Status()

	// 3. Show the webpage, if GET request
	if r.Method != http.MethodPost {
		return data, ""
	}

	// 4. Save the data
	return controller.saveTranslation(r, data)
}

type translationUpdateControllerData struct {
	request       *http.Request
	action        string
	translationID string
	translation   cmsstore.TranslationInterface
	siteList      []cmsstore.SiteInterface
	view          string

	formErrorMessage   string
	formRedirectURL    string
	formSuccessMessage string
	formContent        map[string]string
	formHandle         string
	formName           string
	formMemo           string
	formSiteID         string
	formStatus         string
}
