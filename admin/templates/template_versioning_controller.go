package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"sort"
	"strings"

	"github.com/dracory/bs"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONTROLLER ==============================================================

type templateVersioningController struct {
	ui UiInterface
}

type templateVersioningControllerData struct {
	request        *http.Request
	templateID     string
	versionings    []cmsstore.VersioningInterface
	versioningID   string
	versioning     cmsstore.VersioningInterface
	successMessage string
}

// == CONSTRUCTOR =============================================================

func NewTemplateVersioningController(ui UiInterface) *templateVersioningController {
	return &templateVersioningController{
		ui: ui,
	}
}

func (controller templateVersioningController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		return hb.Swal(hb.SwalOptions{
			Icon: "error",
			Text: errorMessage,
		}).ToHTML()
	}

	if data.successMessage != "" {
		return hb.Wrap().
			Child(hb.Swal(hb.SwalOptions{
				Icon: "success",
				Text: data.successMessage,
			})).
			Child(hb.Script("setTimeout(() => {window.location.href = window.location.href}, 2000)")).
			ToHTML()
	}

	return controller.
		modal(data).
		ToHTML()
}

func (controller *templateVersioningController) modal(data templateVersioningControllerData) hb.TagInterface {
	submitUrl := shared.URLR(data.request, shared.PathTemplatesTemplateVersioning, map[string]string{
		"template_id":   data.templateID,
		"versioning_id": data.versioningID,
	})

	modalID := "ModalTemplateVersioning"
	modalBackdropClass := "ModalBackdrop"

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("Template Revisions").Style(`margin:0px;`)
	if data.versioning != nil {
		name := carbon.Parse(data.versioning.CreatedAt(), carbon.UTC).Format("Y-m-d H:i")
		modalHeading = hb.Heading5().HTML("Template Revision: " + name).Style(`margin:0px;`)
	}

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalTemplateVersioning').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

	buttonSend := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Restore Selected Attributes").
		Class("btn btn-primary float-end").
		HxInclude("#" + modalID).
		HxPost(submitUrl).
		HxSelectOob("#ModalTemplateVersioning").
		HxTarget("body").
		HxSwap("beforeend")

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Close").
		Class("btn btn-secondary float-start").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	table := controller.tableRevisions(data)

	if data.versioning != nil {
		table = controller.tableRevision(data)
	}

	modal := bs.Modal().
		ID(modalID).
		Class("fade show modal-lg").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Child(hb.Script(jsCloseFn)).
		Child(bs.ModalDialog().
			Child(bs.ModalContent().
				Child(
					bs.ModalHeader().
						Child(modalHeading).
						Child(modalClose)).
				Child(
					bs.ModalBody().
						Child(table)).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(buttonCancel).
					ChildIf(data.versioning != nil, buttonSend)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})
}

func (controller *templateVersioningController) tableRevision(data templateVersioningControllerData) hb.TagInterface {
	versioning := data.versioning

	content := versioning.Content()

	if content == "" {
		return hb.Div().Class("alert alert-danger").HTML("Revision is empty. It has no content!")
	}

	dataAny := map[string]any{}
	if err := json.Unmarshal([]byte(content), &dataAny); err != nil {
		return hb.Div().Class("alert alert-danger").HTML(err.Error())
	}

	dataMap := cast.ToStringMapString(dataAny)
	keys := lo.Keys(dataMap)

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().Style("width:1px;text-align:center;").HTML("Apply"),
					hb.TH().Style("width:1px;text-align:center;").HTML("Attribute"),
					hb.TH().HTML("Value"),
				}),
			}),

			hb.Tbody().Children(lo.Map(keys, func(key string, _ int) hb.TagInterface {
				if !slices.Contains(controller.supportedAttributes(), key) {
					return nil
				}

				value := dataMap[key]

				checkbox := hb.Div().
					Class("form-check").
					Child(
						hb.Input().
							Type("checkbox").
							Class("form-check-input").
							Name("revision_attributes").
							Value(key),
					)

				valueContainer := hb.Input().
					Class("form-control w-100").
					Style(`background-color:#eee;`).
					Attr("readonly", "readonly").
					Value(value)
				
				if key == cmsstore.COLUMN_CONTENT {
					valueContainer = hb.TextArea().
						Class("form-control w-100").
						Style(`background-color:#eee;`).
						Attr("readonly", "readonly").
						Attr("rows", "5").
						Text(value)
				}

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Style("text-align:center;").Child(checkbox),
					hb.TD().Style("text-align:center;").Text(key),
					hb.TD().Child(valueContainer),
				})
			})),
		})
}

func (controller *templateVersioningController) tableRevisions(data templateVersioningControllerData) hb.TagInterface {
	return hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().HTML("Version"),
					hb.TH().HTML("Created"),
					hb.TH().HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.versionings, func(versioning cmsstore.VersioningInterface, _ int) hb.TagInterface {
				name := carbon.Parse(versioning.CreatedAt(), carbon.UTC).Format("Y-m-d H:i")
				ago := carbon.Parse(versioning.CreatedAt(), carbon.UTC).DiffForHumans()

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().
						Text(name),
					hb.TD().
						Text(ago),
					hb.TD().Children([]hb.TagInterface{
						hb.Button().
							Class("btn btn-sm btn-primary").
							Child(hb.I().Class("bi bi-eye me-2")).
							Text("Preview").
							HxGet(shared.URLR(data.request, shared.PathTemplatesTemplateVersioning, map[string]string{
								"template_id":   data.templateID,
								"versioning_id": versioning.ID(),
							})).
							HxTarget("#" + "ModalTemplateVersioning").
							HxSwap("outerHTML"),
					}),
				})
			})),
		})

}

func (controller *templateVersioningController) prepareDataAndValidate(r *http.Request) (data templateVersioningControllerData, errorMessage string) {
	var err error
	data.request = r
	data.templateID = strings.TrimSpace(req.GetStringTrimmed(r, "template_id"))
	data.versioningID = strings.TrimSpace(req.GetStringTrimmed(r, "versioning_id"))

	if data.templateID == "" {
		return data, "template id is required"
	}

	template, err := controller.ui.Store().TemplateFindByID(data.request.Context(), data.templateID)

	if err != nil {
		controller.ui.Logger().Error("At templateVersioningController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if template == nil {
		return data, "Template not found"
	}

	data.versionings, err = controller.ui.Store().VersioningList(data.request.Context(), cmsstore.NewVersioningQuery().
		SetEntityType(cmsstore.VERSIONING_TYPE_TEMPLATE).
		SetEntityID(data.templateID).
		SetOrderBy(cmsstore.COLUMN_CREATED_AT).
		SetSortOrder(sb.DESC))

	if err != nil {
		controller.ui.Logger().Error("At templateVersioningController > prepareDataAndValidate", "error", err.Error())
		return data, err.Error()
	}

	if data.versioningID != "" {
		data.versioning, err = controller.ui.Store().VersioningFindByID(data.request.Context(), data.versioningID)

		if err != nil {
			controller.ui.Logger().Error("At templateVersioningController > prepareDataAndValidate", "error", err.Error())
			return data, err.Error()
		}
	}

	if r.Method != http.MethodPost {
		return data, ""
	}

	attrs := req.GetArray(r, "revision_attributes", []string{})

	if len(attrs) < 1 {
		return data, "No revision attributes were selected. Aborted"
	}

	controller.restoreRevisionAttributes(data.request.Context(), template, data.versioning, attrs)

	data.successMessage = "revision attributes restored successfully."

	return data, ""
}

func (controller *templateVersioningController) restoreRevisionAttributes(ctx context.Context, template cmsstore.TemplateInterface, versioning cmsstore.VersioningInterface, attrs []string) error {
	if template == nil {
		return errors.New("template is nil")
	}

	content := versioning.Content()

	if content == "" {
		return errors.New("revision is empty. it has no content!")
	}

	dataAny := map[string]any{}
	err := json.Unmarshal([]byte(content), &dataAny)

	if err != nil {
		return err
	}

	dataMap := cast.ToStringMapString(dataAny)

	for _, attr := range attrs {
		if !slices.Contains(controller.supportedAttributes(), attr) {
			continue
		}

		value := dataMap[attr]

		if attr == cmsstore.COLUMN_CONTENT {
			template.SetContent(value)
		}

		if attr == cmsstore.COLUMN_HANDLE {
			template.SetHandle(value)
		}

		if attr == cmsstore.COLUMN_MEMO {
			template.SetMemo(value)
		}

		if attr == cmsstore.COLUMN_NAME {
			template.SetName(value)
		}

		if attr == cmsstore.COLUMN_STATUS {
			template.SetStatus(value)
		}
		
		if attr == cmsstore.COLUMN_SITE_ID {
			template.SetSiteID(value)
		}
	}

	err = controller.ui.Store().TemplateUpdate(ctx, template)

	if err != nil {
		return err
	}

	return nil
}

func (controller *templateVersioningController) supportedAttributes() []string {
	return []string{
		cmsstore.COLUMN_CONTENT,
		cmsstore.COLUMN_HANDLE,
		cmsstore.COLUMN_MEMO,
		cmsstore.COLUMN_NAME,
		cmsstore.COLUMN_STATUS,
		cmsstore.COLUMN_SITE_ID,
	}
}
