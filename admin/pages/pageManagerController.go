package admin

import (
	"net/http"
	"strings"

	// "project/config"
	// "project/controllers/admin/cms/shared"
	// "project/internal/helpers"
	// "project/internal/layouts"
	// "project/internal/links"
	// "project/pkg/cmsstore"

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
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const ActionModalPageFilterShow = "modal_site_filter_show"

// == CONTROLLER ==============================================================

type pageManagerController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*pageManagerController)(nil)

// == CONSTRUCTOR =============================================================

func NewPageManagerController(ui UiInterface) *pageManagerController {
	return &pageManagerController{
		ui: ui,
	}
}

func (controller *pageManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return api.Error(errorMessage).ToString()
	}

	if data.action == ActionModalPageFilterShow {
		return controller.onModalRecordFilterShow(data).ToHTML()
	}

	options := struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Sweetalert2_11(),
		},
	}

	return controller.ui.Layout(w, r, "Page  Manager | CMS", controller.page(data).ToHTML(), options)
}

func (controller *pageManagerController) onModalRecordFilterShow(data pageManagerControllerData) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.Heading5().
		Text("Filters").
		Style(`margin:0px;padding:0px;`)

	buttonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonOk := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Apply").
		Class("btn btn-primary float-end").
		OnClick(`FormFilters.submit();` + modalCloseScript)

	filterForm := form.NewForm(form.FormOptions{
		ID:        "FormFilters",
		Method:    http.MethodGet,
		ActionURL: shared.URL(controller.ui.Endpoint(), shared.PathPagesPageManager, nil),
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Status",
				Name:  "status",
				Type:  form.FORM_FIELD_TYPE_SELECT,
				Help:  `The status of the site.`,
				Value: data.formStatus,
				Options: []form.FieldOption{
					{
						Value: "",
						Key:   "",
					},
					{
						Value: "Active",
						Key:   cmsstore.SITE_STATUS_ACTIVE,
					},
					{
						Value: "Inactive",
						Key:   cmsstore.SITE_STATUS_INACTIVE,
					},
					{
						Value: "Draft",
						Key:   cmsstore.SITE_STATUS_DRAFT,
					},
				},
			}),
			form.NewField(form.FieldOptions{
				Label: "Name",
				Name:  "name",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formName,
				Help:  `Filter by name.`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: `<div class="row">`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: `<div class="col-6">`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created From",
				Name:  "created_from",
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedFrom,
				Help:  `Filter by creation date.`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: `</div>`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: `<div class="col-6">`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created To",
				Name:  "created_to",
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedTo,
				Help:  `Filter by creation date.`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: `</div>`,
			}),
			form.NewField(form.FieldOptions{
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: `</div>`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Site ID",
				Name:  "site_id",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formSiteID,
				Help:  `Find site by reference number (ID).`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Path",
				Name:  "path",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: shared.PathPagesPageManager,
				Help:  `Path to this page.`,
			}),
		},
	}).Build()

	modal := bs.Modal().
		ID("ModalMessage").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						title,
						buttonModalClose,
					}),

					bs.ModalBody().
						Child(filterForm),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonOk),
				}),
			}),
		})

	backdrop := hb.Div().
		ID("ModalBackdrop").
		Class("modal-backdrop fade show").
		Style("display:block;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})

}

func (controller *pageManagerController) page(data pageManagerControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), controller.ui.Endpoint())

	breadcrumbs := controller.ui.AdminBreadcrumbs(controller.ui.Endpoint(), []shared.Breadcrumb{
		{
			Name: "Page Manager",
			URL:  shared.URL(controller.ui.Endpoint(), shared.PathPagesPageManager, nil),
		},
	})

	buttonPageNew := hb.Button().
		Class("btn btn-primary float-end").
		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("New Page").
		HxGet(shared.URL(controller.ui.Endpoint(), shared.PathPagesPageCreate, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	pageTitle := hb.Heading1().
		HTML("CMS. Page Manager").
		Child(buttonPageNew)

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader).
		Child(hb.HR()).
		Child(pageTitle).
		Child(controller.tableRecords(data))
}

func (controller *pageManagerController) tableRecords(data pageManagerControllerData) hb.TagInterface {
	table := hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Name", cmsstore.COLUMN_NAME)).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Alias", cmsstore.COLUMN_ALIAS)).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Reference", cmsstore.COLUMN_ID)).
						Style(`cursor: pointer;`),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Status", cmsstore.COLUMN_STATUS)).
						Style("width: 200px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Created", cmsstore.COLUMN_CREATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Modified", cmsstore.COLUMN_UPDATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.recordList, func(page cmsstore.PageInterface, _ int) hb.TagInterface {

				pageName := page.Name()
				pageAlias := page.Alias()

				siteLink := hb.Hyperlink().
					Text(pageName).
					Href(shared.URL(controller.ui.Endpoint(), shared.PathPagesPageUpdate, map[string]string{
						"page_id": page.ID(),
					}))

				status := hb.Span().
					Style(`font-weight: bold;`).
					StyleIf(page.IsActive(), `color:green;`).
					StyleIf(page.IsSoftDeleted(), `color:silver;`).
					StyleIf(page.IsInactive(), `color:red;`).
					HTML(page.Status())

				buttonEdit := hb.Hyperlink().
					Class("btn btn-primary me-2").
					Child(hb.I().Class("bi bi-pencil-square")).
					Title("Edit").
					Href(shared.URL(controller.ui.Endpoint(), shared.PathPagesPageUpdate, map[string]string{
						"page_id": page.ID(),
					}))

				buttonDelete := hb.Hyperlink().
					Class("btn btn-danger").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete").
					HxGet(shared.URL(controller.ui.Endpoint(), shared.PathPagesPageDelete, map[string]string{
						"page_id": page.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				// buttonImpersonate := hb.Hyperlink().
				// 	Class("btn btn-warning me-2").
				// 	Child(hb.I().Class("bi bi-shuffle")).
				// 	Title("Impersonate").
				// 	Href(links.NewAdminLinks().UsersUserImpersonate(map[string]string{"site_id": site.ID()}))

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().
						Child(hb.Div().Child(siteLink)).
						Child(hb.Div().
							Style("font-size: 11px;").
							HTML("Alias: ").
							HTML(pageAlias)).
						Child(hb.Div().
							Style("font-size: 11px;").
							HTML("Ref: ").
							HTML(page.ID())),
					hb.TD().
						Child(status),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(page.CreatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(page.UpdatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(buttonEdit).
						// Child(buttonImpersonate).
						Child(buttonDelete),
				})
			})),
		})

	// cfmt.Successln("Table: ", table)

	return hb.Wrap().Children([]hb.TagInterface{
		controller.tableFilter(data),
		table,
		controller.tablePagination(data, int(data.recordCount), data.pageInt, data.perPage),
	})
}

func (controller *pageManagerController) sortableColumnLabel(data pageManagerControllerData, tableLabel string, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == sb.ASC, sb.DESC).Else(sb.ASC)

	if !isSelected {
		direction = sb.ASC
	}

	link := shared.URL(controller.ui.Endpoint(), shared.PathPagesPageManager, map[string]string{
		"page":      "0",
		"by":        columnName,
		"sort":      direction,
		"date_from": data.formCreatedFrom,
		"date_to":   data.formCreatedTo,
		"status":    data.formStatus,
		"site_id":   data.formSiteID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(controller.sortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func (controller *pageManagerController) sortingIndicator(columnName string, sortByColumnName string, sortOrder string) hb.TagInterface {
	isSelected := strings.EqualFold(sortByColumnName, columnName)

	direction := lo.If(isSelected && sortOrder == "asc", "up").
		ElseIf(isSelected && sortOrder == "desc", "down").
		Else("none")

	sortingIndicator := hb.Span().
		Class("sorting").
		HTMLIf(direction == "up", "&#8595;").
		HTMLIf(direction == "down", "&#8593;").
		HTMLIf(direction != "down" && direction != "up", "")

	return sortingIndicator
}

func (controller *pageManagerController) tableFilter(data pageManagerControllerData) hb.TagInterface {
	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(shared.URL(controller.ui.Endpoint(), shared.PathPagesPageManager, map[string]string{
			"action":       ActionModalPageFilterShow,
			"name":         data.formName,
			"status":       data.formStatus,
			"site_id":      data.formSiteID,
			"created_from": data.formCreatedFrom,
			"created_to":   data.formCreatedTo,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing pages").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formName != "" {
		description = append(description, hb.Span().Text("and name: "+data.formName).ToHTML())
	}

	if data.formSiteID != "" {
		description = append(description, hb.Span().Text("and ID: "+data.formSiteID).ToHTML())
	}

	if data.formCreatedFrom != "" && data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created between: "+data.formCreatedFrom+" and "+data.formCreatedTo).ToHTML())
	} else if data.formCreatedFrom != "" {
		description = append(description, hb.Span().Text("and created after: "+data.formCreatedFrom).ToHTML())
	} else if data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created before: "+data.formCreatedTo).ToHTML())
	}

	return hb.Div().
		Class("card bg-light mb-3").
		Style("").
		Children([]hb.TagInterface{
			hb.Div().Class("card-body").
				Child(buttonFilter).
				Child(hb.Span().
					HTML(strings.Join(description, " "))),
		})
}

func (controller *pageManagerController) tablePagination(data pageManagerControllerData, count int, page int, perPage int) hb.TagInterface {
	url := shared.URL(controller.ui.Endpoint(), shared.PathPagesPageManager, map[string]string{
		"status":       data.formStatus,
		"name":         data.formName,
		"created_from": data.formCreatedFrom,
		"created_to":   data.formCreatedTo,
		"by":           data.sortBy,
		"order":        data.sortOrder,
	})

	url = lo.Ternary(strings.Contains(url, "?"), url+"&page=", url+"?page=") // page must be last

	pagination := bs.Pagination(bs.PaginationOptions{
		NumberItems:       count,
		CurrentPageNumber: page,
		PagesToShow:       5,
		PerPage:           perPage,
		URL:               url,
	})

	return hb.Div().
		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
		HTML(pagination)
}

func (controller *pageManagerController) prepareData(r *http.Request) (data pageManagerControllerData, errorMessage string) {
	var err error
	initialPerPage := 20
	data.request = r
	data.action = utils.Req(r, "action", "")
	data.page = utils.Req(r, "page", "0")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(utils.Req(r, "per_page", cast.ToString(initialPerPage)))
	data.sortOrder = utils.Req(r, "sort", sb.DESC)
	data.sortBy = utils.Req(r, "by", cmsstore.COLUMN_CREATED_AT)
	data.formName = utils.Req(r, "name", "")
	data.formStatus = utils.Req(r, "status", "")
	data.formCreatedFrom = utils.Req(r, "created_from", "")
	data.formCreatedTo = utils.Req(r, "created_to", "")

	recordList, recordCount, err := controller.fetchRecordList(data)

	if err != nil {
		controller.ui.Logger().Error("At pageManagerController > prepareData", "error", err.Error())
		return data, "error retrieving web sites"
	}

	data.recordList = recordList
	data.recordCount = recordCount

	return data, ""
}

func (controller *pageManagerController) fetchRecordList(data pageManagerControllerData) (records []cmsstore.PageInterface, recordCount int64, err error) {
	pageIDs := []string{}

	if data.formPageID != "" {
		pageIDs = append(pageIDs, data.formPageID)
	}

	query := cmsstore.PageQuery().
		SetLimit(data.perPage).
		SetOffset(data.pageInt * data.perPage).
		SetOrderBy(data.sortBy).
		SetSortOrder(data.sortOrder)

	if data.formCreatedFrom != "" {
		query = query.SetCreatedAtGte(data.formCreatedFrom + " 00:00:00")
	}

	if data.formCreatedTo != "" {
		query = query.SetCreatedAtLte(data.formCreatedTo + " 23:59:59")
	}

	if data.formSiteID != "" {
		query = query.SetSiteID(data.formSiteID)
	}

	if len(pageIDs) > 0 {
		query = query.SetIDIn(pageIDs)
	}

	if data.formStatus != "" {
		query = query.SetStatus(data.formStatus)
	}

	if data.formAlias != "" {
		query = query.SetAliasLike(data.formAlias)
	}

	if data.formName != "" {
		query = query.SetNameLike(data.formName)
	}

	recordList, err := controller.ui.Store().PageList(query)

	if err != nil {
		return []cmsstore.PageInterface{}, 0, err
	}

	recordCount, err = controller.ui.Store().PageCount(query)

	if err != nil {
		return []cmsstore.PageInterface{}, 0, err
	}

	return recordList, recordCount, nil
}

type pageManagerControllerData struct {
	request         *http.Request
	action          string
	page            string
	pageInt         int
	perPage         int
	sortOrder       string
	sortBy          string
	formStatus      string
	formName        string
	formAlias       string
	formCreatedFrom string
	formCreatedTo   string
	formPageID      string
	formSiteID      string
	recordList      []cmsstore.PageInterface
	recordCount     int64
}
