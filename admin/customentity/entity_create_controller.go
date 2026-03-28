package customentity

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
)

// EntityCreateController handles entity creation
type EntityCreateController struct {
	ui         UiInterface
	definition cmsstore.CustomEntityDefinition
}

// NewEntityCreateController creates a new create controller
func NewEntityCreateController(ui UiInterface, definition cmsstore.CustomEntityDefinition) *EntityCreateController {
	return &EntityCreateController{
		ui:         ui,
		definition: definition,
	}
}

// Handler handles both GET (show form) and POST (process form)
func (c *EntityCreateController) Handler(w http.ResponseWriter, r *http.Request) string {
	if r.Method == http.MethodPost {
		return c.handleSubmit(w, r)
	}
	return c.showForm(r, nil, "")
}

func (c *EntityCreateController) showForm(r *http.Request, data map[string]string, errorMessage string) string {
	if data == nil {
		data = make(map[string]string)
	}

	submitUrl := "/admin/custom-entity/" + c.definition.Type + "/create"

	// Build form fields
	fields := []form.FieldInterface{}
	for _, attr := range c.definition.Attributes {
		fieldType := c.getFieldType(attr.Type)

		field := form.NewField(form.FieldOptions{
			Label:    attr.Label,
			Name:     attr.Name,
			Type:     fieldType,
			Value:    data[attr.Name],
			Required: attr.Required,
			Help:     attr.Help,
		})
		fields = append(fields, field)
	}

	formObj := form.NewForm(form.FormOptions{
		ID:     "FormEntityCreate",
		Fields: fields,
	})

	modalID := "ModalEntityCreate"
	modalBackdropClass := "ModalBackdrop"

	modal := hb.Div().
		ID(modalID).
		Class("modal fade show").
		Style("display: block;").
		Child(hb.Div().Class("modal-dialog modal-lg").
			Child(hb.Div().Class("modal-content").
				Child(c.modalHeader()).
				Child(c.modalBody(formObj, submitUrl, errorMessage)).
				Child(c.modalFooter(modalID))))

	backdrop := hb.Div().Class(modalBackdropClass).Style("display: block;")

	return hb.Wrap().Children([]hb.TagInterface{modal, backdrop}).ToHTML()
}

func (c *EntityCreateController) modalHeader() hb.TagInterface {
	return hb.Div().Class("modal-header").
		Child(hb.Heading5().Class("modal-title").HTML("New " + c.definition.TypeLabel)).
		Child(hb.Button().
			Type("button").
			Class("btn-close").
			Data("bs-dismiss", "modal").
			Attr("aria-label", "Close"))
}

func (c *EntityCreateController) modalBody(formObj *form.Form, submitUrl, errorMessage string) hb.TagInterface {
	body := hb.Div().Class("modal-body")

	if errorMessage != "" {
		body.Child(hb.Div().Class("alert alert-danger").HTML(errorMessage))
	}

	body.Child(hb.Form().
		ID("FormEntityCreate").
		Method(http.MethodPost).
		Action(submitUrl).
		Child(formObj.Build()))

	return body
}

func (c *EntityCreateController) modalFooter(modalID string) hb.TagInterface {
	return hb.Div().Class("modal-footer").
		Child(hb.Button().
			Type("button").
			Class("btn btn-secondary").
			Data("bs-dismiss", "modal").
			HTML("Close")).
		Child(hb.Button().
			Type("submit").
			Class("btn btn-primary").
			Attr("form", "FormEntityCreate").
			HTML("Create " + c.definition.TypeLabel))
}

func (c *EntityCreateController) handleSubmit(w http.ResponseWriter, r *http.Request) string {
	ctx := context.Background()

	if err := r.ParseForm(); err != nil {
		return c.showForm(r, nil, "Error parsing form: "+err.Error())
	}

	customStore := c.ui.Store().CustomEntityStore()
	if customStore == nil {
		return c.showForm(r, nil, "Custom entity store is not available")
	}

	// Build attributes map from form data
	attrs := make(map[string]interface{})
	formData := make(map[string]string)

	for _, attr := range c.definition.Attributes {
		value := r.FormValue(attr.Name)
		formData[attr.Name] = value

		if value != "" {
			// Convert value based on type
			switch attr.Type {
			case "int":
				if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
					attrs[attr.Name] = intVal
				} else {
					attrs[attr.Name] = value
				}
			case "float":
				if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
					attrs[attr.Name] = floatVal
				} else {
					attrs[attr.Name] = value
				}
			case "bool":
				attrs[attr.Name] = value == "true" || value == "1" || value == "on"
			default:
				attrs[attr.Name] = value
			}
		}
	}

	// Create entity
	_, err := customStore.Create(ctx, c.definition.Type, attrs, nil, nil)
	if err != nil {
		return c.showForm(r, formData, "Error creating entity: "+err.Error())
	}

	// Success - close modal and refresh page
	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{
			Icon: "success",
			Text: c.definition.TypeLabel + " created successfully",
		})).
		Child(hb.Script(`
			setTimeout(() => {
				document.querySelector('.modal').remove();
				document.querySelector('.ModalBackdrop').remove();
				window.location.reload();
			}, 1500);
		`)).
		ToHTML()
}

func (c *EntityCreateController) getFieldType(attrType string) string {
	switch attrType {
	case "int":
		return form.FORM_FIELD_TYPE_NUMBER
	case "float":
		return form.FORM_FIELD_TYPE_NUMBER
	case "bool":
		return form.FORM_FIELD_TYPE_SELECT
	case "json":
		return form.FORM_FIELD_TYPE_TEXTAREA
	default:
		return form.FORM_FIELD_TYPE_STRING
	}
}
