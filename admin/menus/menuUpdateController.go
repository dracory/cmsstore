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
	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
)

const VIEW_SETTINGS = "settings"
const VIEW_MENU_ITEMS = "menu_items"
const ACTION_TREEEDITOR_HANDLE = "treeditor_handle"

// == CONTROLLER ==============================================================

type menuUpdateController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*menuUpdateController)(nil)

// == CONSTRUCTOR =============================================================

func NewMenuUpdateController(ui UiInterface) *menuUpdateController {
	return &menuUpdateController{
		ui: ui,
	}
}

func (controller *menuUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		return api.Error(errorMessage).ToString()
	}

	if data.action == ACTION_TREEEDITOR_HANDLE {
		return controller.treeEditorHandle(r, data).ToHTML()
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
			cdn.Htmx_2_0_0(),
		},
	}

	return controller.ui.Layout(w, r, "Edit Menu | CMS", html.ToHTML(), options)
}

func (controller *menuUpdateController) treeEditorHandle(r *http.Request, data menuUpdateControllerData) hb.TagInterface {
	pageList, err := controller.ui.Store().PageList(cmsstore.PageQuery().
		SetSiteID(data.menu.SiteID()).
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC))

	if err != nil {
		controller.ui.Logger().Error("At menuUpdateController > treeEditorHandle", "error", err.Error())
		return hb.Div().Text(`ERROR: ` + err.Error())
	}

	treeControl := &treeControl{
		treeJSON:         data.formMenuItemsJSON,
		targetTextareaID: "menu_items",
		renderURL: shared.URL(controller.ui.Endpoint(), shared.PathMenusMenuUpdate, map[string]string{
			"menu_id": data.menuID,
			"action":  ACTION_TREEEDITOR_HANDLE,
		}),
		pageList: pageList,
	}

	return treeControl.Render(r)
	// jsonString := data.formMenuItems
	// tree, err := NewTreeFromJSON(jsonString)

	// if err != nil {
	// 	return hb.Div().Text(`ERROR: ` + err.Error())
	// }

	// jsonString, err = tree.ToJSON()

	// if err != nil {
	// 	return hb.Div().Text(`ERROR: ` + err.Error())
	// }

	// return hb.Div().ID("TreeEditor").Text("TREE: " + jsonString)
}

func (controller menuUpdateController) page(data menuUpdateControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), controller.ui.Endpoint())

	breadcrumbs := controller.ui.AdminBreadcrumbs(controller.ui.Endpoint(), []shared.Breadcrumb{
		{
			Name: "Menu Manager",
			URL:  shared.URL(controller.ui.Endpoint(), shared.PathMenusMenuManager, nil),
		},
		{
			Name: "Edit Menu",
			URL:  shared.URL(controller.ui.Endpoint(), shared.PathMenusMenuUpdate, map[string]string{"menu_id": data.menuID}),
		},
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormMenuUpdate").
		HxPost(shared.URL(controller.ui.Endpoint(), shared.PathMenusMenuUpdate, map[string]string{"menu_id": data.menuID})).
		HxTarget("#FormMenuUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URL(controller.ui.Endpoint(), shared.PathMenusMenuManager, nil))

	badgeStatus := hb.Div().
		Class("badge fs-6 ms-3").
		ClassIf(data.menu.Status() == cmsstore.TEMPLATE_STATUS_ACTIVE, "bg-success").
		ClassIf(data.menu.Status() == cmsstore.TEMPLATE_STATUS_INACTIVE, "bg-secondary").
		ClassIf(data.menu.Status() == cmsstore.TEMPLATE_STATUS_DRAFT, "bg-warning").
		Text(data.menu.Status())

	pageTitle := hb.Heading1().
		Text("CMS. Edit Menu:").
		Text(" ").
		Text(data.menu.Name()).
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
					HTMLIf(data.view == VIEW_MENU_ITEMS, "Menu Items").
					HTMLIf(data.view == VIEW_SETTINGS, "Menu Settings").
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
				ClassIf(data.view == VIEW_MENU_ITEMS, "active").
				Href(shared.URL(controller.ui.Endpoint(), shared.PathMenusMenuUpdate, map[string]string{
					"menu_id": data.menuID,
					"view":    VIEW_MENU_ITEMS,
				})).
				HTML("Menu Items"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SETTINGS, "active").
				Href(shared.URL(controller.ui.Endpoint(), shared.PathMenusMenuUpdate, map[string]string{
					"menu_id": data.menuID,
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

func (controller menuUpdateController) form(data menuUpdateControllerData) hb.TagInterface {
	fieldsMenuItems := controller.fieldsMenuItems(data)
	fieldsSettings := controller.fieldsSettings(data)

	formpageUpdate := form.NewForm(form.FormOptions{
		ID: "FormMenuUpdate",
	})

	if data.view == VIEW_SETTINGS {
		formpageUpdate.SetFields(fieldsSettings)
	}

	if data.view == VIEW_MENU_ITEMS {
		formpageUpdate.SetFields(fieldsMenuItems)
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

func (c menuUpdateController) fieldsMenuItems(data menuUpdateControllerData) []form.FieldInterface {
	url := shared.URL(c.ui.Endpoint(), shared.PathMenusMenuUpdate, map[string]string{
		"action":  ACTION_TREEEDITOR_HANDLE,
		"menu_id": data.menuID,
	})
	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Type: form.FORM_FIELD_TYPE_RAW,
			Value: hb.Div().
				Text(`Loading menu items tree...`).
				HxPost(url).
				HxTrigger("load").
				HxInclude("#FormMenuUpdate").
				ToHTML(),
		}),
		form.NewField(form.FieldOptions{
			Label: "Menu Items",
			Name:  "menu_items",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formMenuItemsJSON,
		}),
		form.NewField(form.FieldOptions{
			Label:    "Menu ID",
			Name:     "menu_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.menuID,
			Readonly: true,
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_MENU_ITEMS,
			Readonly: true,
		}),
	}

	contentScript := hb.Script(`
function codeMirrorSelector() {
	return 'textarea[name="menu_content"]';
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

func (controller menuUpdateController) fieldsSettings(data menuUpdateControllerData) []form.FieldInterface {
	fieldsSettings := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Status",
			Name:  "menu_status",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formStatus,
			Help:  "The status of this webpage. Published pages will be displayed on the webmenu.",
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
			Label: "Menu Name (Internal)",
			Name:  "menu_name",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formName,
			Help:  "The name of the translation as displayed in the admin panel. This is not vsible to the public.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Belongs to Site",
			Name:  "menu_site_id",
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
			Name:  "menu_memo",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formMemo,
			Help:  "Admin notes for this menu. These notes will not be visible to the public.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "Menu Reference / ID",
			Name:     "menu_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.menuID,
			Readonly: true,
			Help:     "The reference number (ID) of the menu. This is used to identify the menu in the system and should not be changed.",
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

func (controller menuUpdateController) saveMenu(r *http.Request, data menuUpdateControllerData) (d menuUpdateControllerData, errorMessage string) {
	data.formMenuItemsJSON = utils.Req(r, "menu_items", "")
	data.formMemo = utils.Req(r, "menu_memo", "")
	data.formName = utils.Req(r, "menu_name", "")
	data.formStatus = utils.Req(r, "menu_status", "")
	data.formHandle = utils.Req(r, "menu_title", "")
	data.formSiteID = utils.Req(r, "menu_site_id", "")

	refreshPage := false

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
		data.menu.SetMemo(data.formMemo)

		if data.formName != data.menu.Name() {
			refreshPage = true // name has changed, must refersh the page
		}
		data.menu.SetName(data.formName)

		data.menu.SetSiteID(data.formSiteID)

		if data.formStatus != data.menu.Status() {
			refreshPage = true // status has changed, must refersh the page
		}

		data.menu.SetStatus(data.formStatus)
	}

	if data.view == VIEW_MENU_ITEMS {
		err := controller.saveMenuItems(data)

		if err != nil {
			data.formErrorMessage = err.Error()
			return data, ""
		}

	}

	err := controller.ui.Store().MenuUpdate(data.menu)

	if err != nil {
		//config.LogStore.ErrorWithContext("At menuUpdateController > prepareDataAndValidate", err.Error())
		data.formErrorMessage = "System error. Saving menu failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "Menu saved successfully"
	if refreshPage {
		data.formRedirectURL = shared.URL(controller.ui.Endpoint(), shared.PathMenusMenuUpdate, map[string]string{
			"menu_id": data.menuID,
			"view":    data.view,
		})
	}

	return data, ""
}

func (controller menuUpdateController) saveMenuItems(data menuUpdateControllerData) error {
	tree, err := NewTreeFromJSON(data.formMenuItemsJSON)

	if err != nil {
		return err
	}

	menuItemNodes := tree.List()

	idsToRemove := []string{}

	for _, existingMenuItem := range data.menuItemList {
		if !tree.Exists(existingMenuItem.ID()) {
			idsToRemove = append(idsToRemove, existingMenuItem.ID())
		}
	}

	for _, node := range menuItemNodes {
		menuItem, err := controller.ui.Store().MenuItemFindByID(node.ID)

		if err != nil {
			return err
		}

		if menuItem == nil {
			menuItem = cmsstore.NewMenuItem()
			menuItem.SetID(node.ID)
			menuItem.SetMenuID(data.menuID)
			menuItem.SetName(node.Name)
			menuItem.SetParentID(node.ParentID)
			menuItem.SetSequenceInt(node.Sequence)
			menuItem.SetPageID(node.PageID)
			menuItem.SetURL(node.URL)
			menuItem.SetTarget(node.Target)

			err = controller.ui.Store().MenuItemCreate(menuItem)

			if err != nil {
				return err
			}
		}

		if menuItem.Name() != node.Name {
			menuItem.SetName(node.Name)
		}

		if menuItem.ParentID() != node.ParentID {
			menuItem.SetParentID(node.ParentID)
		}

		if menuItem.SequenceInt() != node.Sequence {
			menuItem.SetSequenceInt(node.Sequence)
		}

		if menuItem.PageID() != node.PageID {
			menuItem.SetPageID(node.PageID)
		}

		if menuItem.URL() != node.URL {
			menuItem.SetURL(node.URL)
		}

		if menuItem.Target() != node.Target {
			menuItem.SetTarget(node.Target)
		}

		err = controller.ui.Store().MenuItemUpdate(menuItem)

		if err != nil {
			return err
		}
	}

	for _, id := range idsToRemove {
		menuItem, err := controller.ui.Store().MenuItemFindByID(id)

		if err != nil {
			return err
		}

		if menuItem == nil {
			continue // nothing to do, menu item not found
		}

		err = controller.ui.Store().MenuItemSoftDelete(menuItem)

		if err != nil {
			return err
		}
	}

	return nil
}

func (controller menuUpdateController) prepareDataAndValidate(r *http.Request) (data menuUpdateControllerData, errorMessage string) {
	data.action = utils.Req(r, "action", "")
	data.menuID = utils.Req(r, "menu_id", "")
	data.view = utils.Req(r, "view", "")

	if data.view == "" {
		data.view = VIEW_MENU_ITEMS
	}

	if data.menuID == "" {
		return data, "menu id is required"
	}

	var err error
	data.menu, err = controller.ui.Store().MenuFindByID(data.menuID)

	if err != nil {
		controller.ui.Logger().Error("At menuUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if data.menu == nil {
		return data, "menu not found"
	}

	data.siteList, err = controller.ui.Store().SiteList(cmsstore.SiteQuery())

	if err != nil {
		controller.ui.Logger().Error("At translationUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.menuItemList, err = controller.ui.Store().MenuItemList(cmsstore.MenuItemQuery().
		SetMenuID(data.menuID))

	if err != nil {
		controller.ui.Logger().Error("At menuUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	// 2. Populate form data
	menuItemsJson, err := controller.buildMenuItemsJson(data.menuItemList)

	if err != nil {
		controller.ui.Logger().Error("At menuUpdateController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	data.formMenuItemsJSON = menuItemsJson
	data.formName = data.menu.Name()
	data.formMemo = data.menu.Memo()
	data.formSiteID = data.menu.SiteID()
	data.formStatus = data.menu.Status()

	// 3. Show the webpage, if GET request
	if r.Method != http.MethodPost {
		return data, ""
	}

	if data.action == ACTION_TREEEDITOR_HANDLE {
		return data, "" // nothing to do, this is handled in the tree editor
	}

	// 4. Save the data
	return controller.saveMenu(r, data)
}

func (controller menuUpdateController) buildMenuItemsJson(menuItems []cmsstore.MenuItemInterface) (string, error) {
	tree := NewTreeFromMenuItems(menuItems)

	jsonString, err := tree.ToJSON()

	if err != nil {
		return "", err
	}

	return jsonString, nil
}

type menuUpdateControllerData struct {
	action       string
	menuID       string
	menu         cmsstore.MenuInterface
	menuItemList []cmsstore.MenuItemInterface
	siteList     []cmsstore.SiteInterface
	view         string

	formErrorMessage   string
	formRedirectURL    string
	formSuccessMessage string
	formMenuItemsJSON  string
	formHandle         string
	formName           string
	formMemo           string
	formSiteID         string
	formStatus         string
}
