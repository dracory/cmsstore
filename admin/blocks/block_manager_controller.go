package admin

import (
	"net/http"
	"strings"

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

const ActionModalPageFilterShow = "modal_block_filter_show"

// == CONTROLLER ==============================================================

type blockManagerController struct {
	ui UiInterface
}

var _ router.HTMLControllerInterface = (*blockManagerController)(nil)

// == CONSTRUCTOR =============================================================

func NewBlockManagerController(ui UiInterface) *blockManagerController {
	return &blockManagerController{
		ui: ui,
	}
}

func (controller *blockManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
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
	return controller.ui.Layout(w, r, "Block Manager | CMS", controller.page(data).ToHTML(), options)
}

func (controller *blockManagerController) onModalRecordFilterShow(data blockManagerControllerData) *hb.Tag {
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

	fieldSiteID := form.NewField(form.FieldOptions{
		Label: "Site ID",
		Name:  "filter_site_id",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formSiteID,
		Help:  `Find site by reference number (ID).`,
	})

	filterForm := form.NewForm(form.FormOptions{
		ID:        "FormFilters",
		Method:    http.MethodGet,
		ActionURL: shared.URLR(data.request, shared.PathBlocksBlockManager, nil),
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Status",
				Name:  "status",
				Type:  form.FORM_FIELD_TYPE_SELECT,
				Help:  `The status of the block.`,
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
				Value: `</div><div class="col-6">`,
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
				Value: `</div></div>`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Block ID",
				Name:  "block_id",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formBlockID,
				Help:  `Find block by reference number (ID).`,
			}),
			fieldSiteID,
			// !!! Needed or it loses the path from the get submission
			form.NewField(form.FieldOptions{
				Label: "Path",
				Name:  "path",
				Type:  form.FORM_FIELD_TYPE_HIDDEN,
				Value: shared.PathBlocksBlockManager,
				Help:  `The path of the page.`,
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

func (controller *blockManagerController) page(data blockManagerControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.ui.Store(), controller.ui.Logger(), data.request)

	breadcrumbs := shared.AdminBreadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Block Manager",
			URL:  shared.URLR(data.request, shared.PathBlocksBlockManager, nil),
		},
	}, struct{ SiteList []cmsstore.SiteInterface }{
		SiteList: data.siteList,
	})

	buttonPageNew := hb.Button().
		Class("btn btn-primary float-end").
		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("New Block").
		HxGet(shared.URLR(data.request, shared.PathBlocksBlockCreate, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	pageTitle := hb.Heading1().
		HTML("CMS. Block Manager").
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

func (controller *blockManagerController) tableRecords(data blockManagerControllerData) hb.TagInterface {
	table := hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Name", cmsstore.COLUMN_NAME)).
						Text(", ").
						// Child(controller.sortableColumnLabel(data, "Site", cmsstore.COLUMN_SITE_ID)).
						// Text(", ").
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
			hb.Tbody().Children(lo.Map(data.recordList, func(block cmsstore.BlockInterface, _ int) hb.TagInterface {
				site, siteFound := lo.Find(data.siteList, func(site cmsstore.SiteInterface) bool {
					return site.ID() == block.SiteID()
				})

				siteName := lo.IfF(siteFound, func() string { return site.Name() }).Else("none")

				blockName := block.Name()

				blockLink := hb.Hyperlink().
					Text(blockName).
					Href(shared.URLR(data.request, shared.PathBlocksBlockUpdate, map[string]string{
						"block_id": block.ID(),
					}))

				status := hb.Span().
					Style(`font-weight: bold;`).
					StyleIf(block.IsActive(), `color:green;`).
					StyleIf(block.IsSoftDeleted(), `color:silver;`).
					StyleIf(block.IsInactive(), `color:red;`).
					HTML(block.Status())

				buttonEdit := hb.Hyperlink().
					Class("btn btn-primary me-2").
					Child(hb.I().Class("bi bi-pencil-square")).
					Title("Edit").
					Href(shared.URLR(data.request, shared.PathBlocksBlockUpdate, map[string]string{
						"block_id": block.ID(),
					}))

				buttonDelete := hb.Hyperlink().
					Class("btn btn-danger").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete").
					HxGet(shared.URLR(data.request, shared.PathBlocksBlockDelete, map[string]string{
						"block_id": block.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().
						Child(hb.Div().Child(blockLink)).
						Child(hb.Div().
							Style("font-size: 11px;").
							HTML("Site: ").
							HTML(siteName)).
						Child(hb.Div().
							Style("font-size: 11px;").
							HTML("Ref: ").
							HTML(block.ID())),
					hb.TD().
						Child(status),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(block.CreatedAtCarbon().Format("d M Y"))),
					hb.TD().
						Child(hb.Div().
							Style("font-size: 13px;white-space: nowrap;").
							HTML(block.UpdatedAtCarbon().Format("d M Y"))),
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

func (controller *blockManagerController) sortableColumnLabel(data blockManagerControllerData, tableLabel string, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == sb.ASC, sb.DESC).Else(sb.ASC)

	if !isSelected {
		direction = sb.ASC
	}

	link := shared.URLR(data.request, shared.PathBlocksBlockManager, map[string]string{
		"page":      "0",
		"by":        columnName,
		"sort":      direction,
		"date_from": data.formCreatedFrom,
		"date_to":   data.formCreatedTo,
		"status":    data.formStatus,
		"block_id":  data.formBlockID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(controller.sortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func (controller *blockManagerController) sortingIndicator(columnName string, sortByColumnName string, sortOrder string) hb.TagInterface {
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

func (controller *blockManagerController) tableFilter(data blockManagerControllerData) hb.TagInterface {
	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info text-white me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(shared.URLR(data.request, shared.PathBlocksBlockManager, map[string]string{
			"action":       ActionModalPageFilterShow,
			"name":         data.formName,
			"status":       data.formStatus,
			"block_id":     data.formBlockID,
			"created_from": data.formCreatedFrom,
			"created_to":   data.formCreatedTo,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing blocks").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formName != "" {
		description = append(description, hb.Span().Text("and name: "+data.formName).ToHTML())
	}

	if data.formBlockID != "" {
		description = append(description, hb.Span().Text("and ID: "+data.formBlockID).ToHTML())
	}

	if data.formSiteID != "" {
		description = append(description, shared.FilterDescriptionSite(data.request.Context(), controller.ui.Store(), data.formSiteID).ToHTML())
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

func (controller *blockManagerController) tablePagination(data blockManagerControllerData, count int, page int, perPage int) hb.TagInterface {
	url := shared.URLR(data.request, shared.PathBlocksBlockManager, map[string]string{
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

func (controller *blockManagerController) prepareData(r *http.Request) (data blockManagerControllerData, errorMessage string) {
	var err error
	initialPerPage := 20
	data.request = r
	data.action = utils.Req(r, "action", "")
	data.page = utils.Req(r, "page", "0")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(utils.Req(r, "per_page", cast.ToString(initialPerPage)))
	data.sortOrder = utils.Req(r, "sort", sb.DESC)
	data.sortBy = utils.Req(r, "by", cmsstore.COLUMN_CREATED_AT)

	data.formBlockID = utils.Req(r, "filter_block_id", "")
	data.formCreatedFrom = utils.Req(r, "filter_created_from", "")
	data.formCreatedTo = utils.Req(r, "filter_created_to", "")
	data.formName = utils.Req(r, "filter_name", "")
	data.formSiteID = utils.Req(r, "filter_site_id", "")
	data.formStatus = utils.Req(r, "filter_status", "")

	recordList, recordCount, err := controller.fetchRecordList(data)

	if err != nil {
		controller.ui.Logger().Error("At blockManagerController > prepareData", "error", err.Error())
		return data, "error retrieving web blocks"
	}

	data.recordList = recordList
	data.recordCount = recordCount

	data.siteList, err = controller.ui.Store().SiteList(r.Context(), cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		controller.ui.Logger().Error("At translationManagerController > prepareData", "error", err.Error())
		return data, "error retrieving sites"
	}

	return data, ""
}

func (controller *blockManagerController) fetchRecordList(data blockManagerControllerData) (records []cmsstore.BlockInterface, recordCount int64, err error) {
	blockIDs := []string{}

	if data.formBlockID != "" {
		blockIDs = append(blockIDs, data.formBlockID)
	}

	query := cmsstore.BlockQuery()

	if len(blockIDs) > 0 {
		query.SetIDIn(blockIDs)
	}

	if data.formName != "" {
		query.SetNameLike(data.formName)
	}

	if data.formStatus != "" {
		query.SetStatus(data.formStatus)
	}

	if data.formSiteID != "" {
		query.SetSiteID(data.formSiteID)
	}

	// if data.formCreatedFrom != "" {
	// 	query.SetCreatedAtGte(data.formCreatedFrom + " 00:00:00")
	// }

	// if data.formCreatedTo != "" {
	// 	query.CreatedAtLte = data.formCreatedTo + " 23:59:59"
	// }

	// query, err = query.SetLimit(data.perPage)

	// if err != nil {
	// 	return []cmsstore.BlockInterface{}, 0, err
	// }

	// query, err = query.SetOffset(data.pageInt * data.perPage)

	// if err != nil {
	// 	return []cmsstore.BlockInterface{}, 0, err
	// }

	// query, err = query.SetOrderBy(data.sortBy)

	// if err != nil {
	// 	return []cmsstore.BlockInterface{}, 0, err
	// }

	// query, err = query.SetSortOrder(data.sortOrder)

	// if err != nil {
	// 	return []cmsstore.BlockInterface{}, 0, err
	// }

	recordList, err := controller.ui.Store().BlockList(data.request.Context(), query)

	if err != nil {
		return []cmsstore.BlockInterface{}, 0, err
	}

	recordCount, err = controller.ui.Store().BlockCount(data.request.Context(), query)

	if err != nil {
		return []cmsstore.BlockInterface{}, 0, err
	}

	return recordList, recordCount, nil
}

type blockManagerControllerData struct {
	request *http.Request
	action  string

	page      string
	pageInt   int
	perPage   int
	sortOrder string
	sortBy    string

	siteList []cmsstore.SiteInterface

	formBlockID     string
	formCreatedFrom string
	formCreatedTo   string
	formName        string
	formSiteID      string
	formStatus      string

	recordList  []cmsstore.BlockInterface
	recordCount int64
}
