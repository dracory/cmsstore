package admin

import (
	"net/http"
	"slices"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/cms/types"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/router"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const VIEW_SETTINGS = "settings"
const VIEW_SEO = "seo"
const ACTION_REPEATER_ADD = "repeater_add"
const ACTION_REPEATER_DELETE = "repeater_delete"

// == CONTROLLER ==============================================================

type siteUpdateController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*siteUpdateController)(nil)

// == CONSTRUCTOR =============================================================

func NewSiteUpdateController(ui UiInterface) *siteUpdateController {
	return &siteUpdateController{
		ui: ui,
	}
}

func (controller *siteUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		return api.Error(errorMessage).ToString()
	}

	if data.action == ACTION_REPEATER_ADD {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_REPEATER_DELETE {
		return controller.form(data).ToHTML()
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
		Styles:    []string{},
		StyleURLs: []string{},
		Scripts:   []string{},
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
	}

	return controller.ui.Layout(w, r, "Edit Site | CMS", html.ToHTML(), options)
}

func (controller siteUpdateController) page(data siteUpdateControllerData) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs([]shared.Breadcrumb{
		{
			Name: "Home",
			URL:  controller.ui.URL(controller.ui.Endpoint(), "", nil),
		},
		{
			Name: "CMS",
			URL:  controller.ui.URL(controller.ui.Endpoint(), "", nil),
		},
		{
			Name: "Site Manager",
			URL:  controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathSiteManager(), nil),
		},
		{
			Name: "Edit Site",
			URL:  controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathSiteUpdate(), map[string]string{"site_id": data.siteID}),
		},
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormpageUpdate").
		HxPost(controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathSiteUpdate(), map[string]string{"site_id": data.siteID})).
		HxTarget("#FormpageUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathSiteManager(), nil))

	heading := hb.Heading1().
		Text("CMS. Edit Site:").
		Text(" ").
		Text(data.site.Name()).
		Child(buttonSave).
		Child(buttonCancel)

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTMLIf(data.view == VIEW_SETTINGS, "Web Site Settings").
					HTMLIf(data.view == VIEW_SEO, "Web Site SEO").
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
				ClassIf(data.view == VIEW_SETTINGS, "active").
				Href(controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathSiteUpdate(), map[string]string{
					"site_id": data.siteID,
					"view":    VIEW_SETTINGS,
				})).
				HTML("Settings"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SEO, "active").
				Href(controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathSiteUpdate(), map[string]string{
					"site_id": data.siteID,
					"view":    VIEW_SEO,
				})).
				HTML("SEO")))

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		// HTML(header).
		Child(heading).
		// HTML(breadcrumbs).
		// Child(pageTitle).
		Child(tabs).
		Child(card)
}

func (controller siteUpdateController) form(data siteUpdateControllerData) hb.TagInterface {
	fieldsSettings := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Status",
			Name:  "site_status",
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
					Key:   types.WEBPAGE_STATUS_DRAFT,
				},
				{
					Value: "Published",
					Key:   types.WEBPAGE_STATUS_ACTIVE,
				},
				{
					Value: "Unpublished",
					Key:   types.WEBPAGE_STATUS_INACTIVE,
				},
				{
					Value: "In Trash Bin",
					Key:   types.WEBPAGE_STATUS_DELETED,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Website Name (Internal)",
			Name:  "site_name",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formName,
			Help:  "The name of the site as displayed in the admin panel. This is not vsible to the site vistors",
		}),
		form.NewRepeater(form.RepeaterOptions{
			Label: "Domain Names",
			Fields: []form.FieldInterface{
				form.NewField(form.FieldOptions{
					Label: "Domain Name",
					Name:  "site_domain_name",
					Type:  form.FORM_FIELD_TYPE_STRING,
				}),
			},
			Values: lo.Map(data.formDomainNames, func(domainName string, _ int) map[string]string {
				return map[string]string{
					"site_domain_name": domainName,
				}
			}),
			RepeaterAddUrl: controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathSiteUpdate(), map[string]string{
				"site_id": data.siteID,
				"view":    VIEW_SETTINGS,
				"action":  ACTION_REPEATER_ADD,
			}),
			RepeaterRemoveUrl: controller.ui.URL(controller.ui.Endpoint(), controller.ui.PathSiteUpdate(), map[string]string{
				"site_id": data.siteID,
				"view":    VIEW_SETTINGS,
				"action":  ACTION_REPEATER_DELETE,
			}),
		}),
		form.NewField(form.FieldOptions{
			Label: "Admin Notes (Internal)",
			Name:  "site_memo",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formMemo,
			Help:  "Admin notes for this site. These notes will not be visible to the public.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "Website ID",
			Name:     "site_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.siteID,
			Readonly: true,
			Help:     "The reference number (ID) of the website. This is used to identify the website in the system and should not be changed.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.view,
			Readonly: true,
		}),
	}

	fieldsSEO := controller.fieldsSEO(data)

	formpageUpdate := form.NewForm(form.FormOptions{
		ID: "FormpageUpdate",
	})

	if data.view == VIEW_SETTINGS {
		formpageUpdate.SetFields(fieldsSettings)
	}

	if data.view == VIEW_SEO {
		formpageUpdate.SetFields(fieldsSEO)
	}

	if data.formErrorMessage != "" {
		formpageUpdate.AddField(&form.Field{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
		})
	}

	if data.formSuccessMessage != "" {
		formpageUpdate.AddField(&form.Field{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}).ToHTML(),
		})
	}

	return formpageUpdate.Build()
}

func (siteUpdateController) fieldsSEO(data siteUpdateControllerData) []form.FieldInterface {
	fieldsSEO := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label:    "Website ID",
			Name:     "site_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.siteID,
			Readonly: true,
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_SEO,
			Readonly: true,
		}),
	}
	return fieldsSEO
}

func (controller siteUpdateController) saveSite(r *http.Request, data siteUpdateControllerData) (d siteUpdateControllerData, errorMessage string) {
	data.formMemo = utils.Req(r, "site_memo", "")
	data.formName = utils.Req(r, "site_name", "")
	data.formStatus = utils.Req(r, "site_status", "")
	data.formTitle = utils.Req(r, "site_title", "")
	data.formDomainNames = utils.ReqArray(r, "site_domain_name", []string{})

	if data.view == VIEW_SETTINGS {
		if data.formStatus == "" {
			data.formErrorMessage = "Status is required"
			return data, ""
		}
	}

	if data.view == VIEW_SETTINGS {
		data.site.SetMemo(data.formMemo)
		data.site.SetName(data.formName)
		data.site.SetStatus(data.formStatus)
		_, err := data.site.SetDomainNames(data.formDomainNames)

		if err != nil {
			data.formErrorMessage = err.Error()
			return data, ""
		}
	}

	if data.view == VIEW_SEO {
		// nothing here yet
	}

	err := controller.ui.Store().SiteUpdate(data.site)

	if err != nil {
		//config.LogStore.ErrorWithContext("At siteUpdateController > prepareDataAndValidate", err.Error())
		data.formErrorMessage = "System error. Saving site failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "site saved successfully"

	return data, ""
}

func (controller siteUpdateController) prepareDataAndValidate(r *http.Request) (data siteUpdateControllerData, errorMessage string) {
	data.action = utils.Req(r, "action", "")
	data.siteID = utils.Req(r, "site_id", "")
	data.view = utils.Req(r, "view", "")

	if data.view == "" {
		data.view = VIEW_SETTINGS
	}

	if data.siteID == "" {
		return data, "site id is required"
	}

	var err error
	data.site, err = controller.ui.Store().SiteFindByID(data.siteID)

	if err != nil {
		controller.ui.Logger().Error("At siteUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if data.site == nil {
		return data, "site not found"
	}

	data.formName = data.site.Name()
	data.formMemo = data.site.Memo()
	data.formStatus = data.site.Status()
	data.formDomainNames, err = data.site.DomainNames()

	if err != nil {
		controller.ui.Logger().Error("At siteUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if r.Method != http.MethodPost {
		return data, ""
	}

	if data.action == ACTION_REPEATER_ADD {
		data.formDomainNames = append(data.formDomainNames, "")
		return data, ""
	}

	if data.action == ACTION_REPEATER_DELETE {
		repeatableRemoveIndex := utils.Req(r, "repeatable_remove_index", "")

		if repeatableRemoveIndex == "" {
			return data, ""
		}

		data.formDomainNames = slices.Delete(data.formDomainNames, cast.ToInt(repeatableRemoveIndex), cast.ToInt(repeatableRemoveIndex)+1)

		return data, ""
	}

	return controller.saveSite(r, data)
}

type siteUpdateControllerData struct {
	action string
	siteID string
	site   cmsstore.SiteInterface
	view   string

	formErrorMessage   string
	formSuccessMessage string
	formHandler        string
	formName           string
	formDomainNames    []string
	formMemo           string
	formStatus         string
	formTitle          string
}
