package admin

import (
	"net/http"
	"slices"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/gouniverse/router"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const VIEW_SETTINGS = "settings"
const VIEW_SEO = "seo"
const ACTION_REPEATER_ADD = "repeater_add"
const ACTION_REPEATER_DELETE = "repeater_delete"
const ACTION_REPEATER_MOVE_UP = "repeater_move_up"
const ACTION_REPEATER_MOVE_DOWN = "repeater_move_down"

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

	if data.action == ACTION_REPEATER_MOVE_UP {
		return controller.form(data).ToHTML()
	}

	if data.action == ACTION_REPEATER_MOVE_DOWN {
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
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), data.request)

	breadcrumbs := shared.AdminBreadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Site Manager",
			URL:  shared.URLR(data.request, shared.PathSitesSiteManager, nil),
		},
		{
			Name: "Edit Site",
			URL:  shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{"site_id": data.siteID}),
		},
	}, struct{ SiteList []cmsstore.SiteInterface }{
		SiteList: data.siteList,
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormpageUpdate").
		HxPost(shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{"site_id": data.siteID})).
		HxTarget("#FormpageUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URLR(data.request, shared.PathSitesSiteManager, nil))

	badgeStatus := hb.Div().
		Class("badge fs-6 ms-3").
		ClassIf(data.site.Status() == cmsstore.SITE_STATUS_ACTIVE, "bg-success").
		ClassIf(data.site.Status() == cmsstore.SITE_STATUS_INACTIVE, "bg-secondary").
		ClassIf(data.site.Status() == cmsstore.SITE_STATUS_DRAFT, "bg-warning").
		Text(data.site.Status())

	pageTitle := hb.Heading1().
		Text("CMS. Edit Site:").
		Text(" ").
		Text(data.site.Name()).
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
					HTMLIf(data.view == VIEW_SETTINGS, "Site Settings").
					HTMLIf(data.view == VIEW_SEO, "Site SEO").
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
				Href(shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{
					"site_id": data.siteID,
					"view":    VIEW_SETTINGS,
				})).
				HTML("Settings"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SEO, "active").
				Href(shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{
					"site_id": data.siteID,
					"view":    VIEW_SEO,
				})).
				HTML("SEO")))

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

func (controller siteUpdateController) form(data siteUpdateControllerData) hb.TagInterface {
	fieldsSettings := controller.fieldsSettings(data)
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

func (siteUpdateController) fieldsSettings(data siteUpdateControllerData) []form.FieldInterface {
	fieldDomainNames := form.NewRepeater(form.RepeaterOptions{
		Label: "Domain Names",
		Name:  "site_domain_names",
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
		RepeaterAddUrl: shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{
			"site_id": data.siteID,
			"view":    VIEW_SETTINGS,
			"action":  ACTION_REPEATER_ADD,
		}),
		RepeaterRemoveUrl: shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{
			"site_id": data.siteID,
			"view":    VIEW_SETTINGS,
			"action":  ACTION_REPEATER_DELETE,
		}),
		RepeaterMoveUpUrl: shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{
			"site_id": data.siteID,
			"view":    VIEW_SETTINGS,
			"action":  ACTION_REPEATER_MOVE_UP,
		}),
		RepeaterMoveDownUrl: shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{
			"site_id": data.siteID,
			"view":    VIEW_SETTINGS,
			"action":  ACTION_REPEATER_MOVE_DOWN,
		}),
	})

	fieldMemo := form.NewField(form.FieldOptions{
		Label: "Admin Notes (Internal)",
		Name:  "site_memo",
		Type:  form.FORM_FIELD_TYPE_TEXTAREA,
		Value: data.formMemo,
		Help:  "Admin notes for this site. These notes will not be visible to the public.",
	})

	fieldSiteName := form.NewField(form.FieldOptions{
		Label: "Site Name (Internal)",
		Name:  "site_name",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formName,
		Help:  "The name of the site as displayed in the admin panel. This is not vsible to the site vistors",
	})

	fieldStatus := form.NewField(form.FieldOptions{
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
				Key:   cmsstore.SITE_STATUS_DRAFT,
			},
			{
				Value: "Published",
				Key:   cmsstore.SITE_STATUS_ACTIVE,
			},
			{
				Value: "Unpublished",
				Key:   cmsstore.SITE_STATUS_INACTIVE,
			},
		},
	})

	fieldSiteID := form.NewField(form.FieldOptions{
		Label:    "Site Reference / ID",
		Name:     "site_id",
		Type:     form.FORM_FIELD_TYPE_STRING,
		Value:    data.siteID,
		Readonly: true,
		Help:     "The reference number (ID) of the site. This is used to identify the site in the system and should not be changed.",
	})

	// !!! required, so that the correct view is shown/saved
	fieldView := form.NewField(form.FieldOptions{
		Label:    "View",
		Name:     "view",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.view,
		Readonly: true,
	})

	fieldsSettings := []form.FieldInterface{
		fieldStatus,
		fieldSiteName,
		fieldDomainNames,
		fieldMemo,
		fieldSiteID,
		fieldView,
	}

	return fieldsSettings
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
	data.formMemo = req.GetStringTrimmed(r, "site_memo")
	data.formName = req.GetStringTrimmed(r, "site_name")
	data.formStatus = req.GetStringTrimmed(r, "site_status")
	data.formTitle = req.GetStringTrimmed(r, "site_title")
	data.formDomainNames = controller.requestMapToDomainNames(r)

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

	err := controller.ui.Store().SiteUpdate(data.request.Context(), data.site)

	if err != nil {
		//config.LogStore.ErrorWithContext("At siteUpdateController > prepareDataAndValidate", err.Error())
		data.formErrorMessage = "System error. Saving site failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "site saved successfully"

	data.formRedirectURL = shared.URLR(data.request, shared.PathSitesSiteUpdate, map[string]string{
		"site_id": data.siteID,
		"view":    data.view,
	})

	return data, ""
}

func (controller siteUpdateController) prepareDataAndValidate(r *http.Request) (data siteUpdateControllerData, errorMessage string) {
	var err error
	data.request = r
	data.action = req.GetStringTrimmed(r, "action")
	data.siteID = req.GetStringTrimmed(r, "site_id")
	data.view = req.GetStringTrimmed(r, "view")

	if data.view == "" {
		data.view = VIEW_SETTINGS
	}

	if data.siteID == "" {
		return data, "site id is required"
	}

	data.site, err = controller.ui.Store().SiteFindByID(data.request.Context(), data.siteID)

	if err != nil {
		controller.ui.Logger().Error("At siteUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if data.site == nil {
		return data, "site not found"
	}

	data.siteList, err = controller.ui.Store().SiteList(data.request.Context(), cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		return data, "Site list failed to be retrieved" + err.Error()
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

	data.formDomainNames = controller.requestMapToDomainNames(r)

	if data.action == ACTION_REPEATER_ADD {
		data.formDomainNames = append(data.formDomainNames, "")
		return data, ""
	}

	if data.action == ACTION_REPEATER_DELETE {
		repeatableRemoveIndex := req.GetStringTrimmed(r, "repeatable_remove_index")

		if repeatableRemoveIndex == "" {
			return data, ""
		}

		data.formDomainNames = slices.Delete(data.formDomainNames, cast.ToInt(repeatableRemoveIndex), cast.ToInt(repeatableRemoveIndex)+1)

		return data, ""
	}

	if data.action == ACTION_REPEATER_MOVE_UP {
		repeatableMoveUpIndex := cast.ToInt(req.GetStringTrimmed(r, "repeatable_move_up_index"))

		if repeatableMoveUpIndex == 0 {
			return data, ""
		}

		current := data.formDomainNames[repeatableMoveUpIndex]
		upper := data.formDomainNames[repeatableMoveUpIndex-1]

		data.formDomainNames[repeatableMoveUpIndex] = upper
		data.formDomainNames[repeatableMoveUpIndex-1] = current

		return data, ""
	}

	if data.action == ACTION_REPEATER_MOVE_DOWN {
		repeatableMoveDownIndex := cast.ToInt(req.GetStringTrimmed(r, "repeatable_move_down_index"))

		if repeatableMoveDownIndex == len(data.formDomainNames)-1 {
			return data, ""
		}

		current := data.formDomainNames[repeatableMoveDownIndex]
		lower := data.formDomainNames[repeatableMoveDownIndex+1]

		data.formDomainNames[repeatableMoveDownIndex] = lower
		data.formDomainNames[repeatableMoveDownIndex+1] = current

		return data, ""
	}

	return controller.saveSite(r, data)
}

func (controller siteUpdateController) requestMapToDomainNames(r *http.Request) []string {
	formDomainNames := req.GetMaps(r, "site_domain_names", []map[string]string{})
	domainNames := []string{}

	for _, formDomainName := range formDomainNames {
		domainNames = append(domainNames, formDomainName["site_domain_name"])
	}

	return domainNames
}

type siteUpdateControllerData struct {
	request  *http.Request
	action   string
	siteID   string
	site     cmsstore.SiteInterface
	siteList []cmsstore.SiteInterface
	view     string

	formErrorMessage   string
	formRedirectURL    string
	formSuccessMessage string
	formHandler        string
	formName           string
	formDomainNames    []string
	formMemo           string
	formStatus         string
	formTitle          string
}
