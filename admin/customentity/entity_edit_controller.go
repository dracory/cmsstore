package customentity

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
)

// EntityEditController handles entity editing
type EntityEditController struct {
	ui         UiInterface
	definition cmsstore.CustomEntityDefinition
}

// NewEntityEditController creates a new edit controller
func NewEntityEditController(ui UiInterface, definition cmsstore.CustomEntityDefinition) *EntityEditController {
	return &EntityEditController{
		ui:         ui,
		definition: definition,
	}
}

// Handler handles both GET (show form) and POST (process form)
func (c *EntityEditController) Handler(w http.ResponseWriter, r *http.Request) string {
	entityID := r.URL.Query().Get("id")
	if entityID == "" {
		return c.errorResponse("Entity ID is required")
	}

	if r.Method == http.MethodPost {
		return c.handleSubmit(w, r, entityID)
	}
	return c.showForm(r, entityID, nil, "")
}

func (c *EntityEditController) showForm(_ *http.Request, entityID string, formData map[string]string, errorMessage string) string {
	ctx := context.Background()
	customStore := c.ui.Store().CustomEntityStore()

	// Load entity
	entity, err := customStore.FindByID(ctx, entityID)
	if err != nil || entity == nil {
		return c.errorResponse("Entity not found")
	}

	// Load attribute values if not provided
	if formData == nil {
		formData = make(map[string]string)
		for _, attr := range c.definition.Attributes {
			attrValue, _ := customStore.Inner().AttributeFind(ctx, entityID, attr.Name)
			if attrValue != nil {
				formData[attr.Name] = attrValue.AttributeValue()
			}
		}
	}

	submitUrl := "/admin/custom-entity/" + c.definition.Type + "/edit?id=" + entityID

	// Build form fields
	fields := []form.FieldInterface{}
	for _, attr := range c.definition.Attributes {
		fieldType := c.getFieldType(attr.Type)

		field := form.NewField(form.FieldOptions{
			Label:    attr.Label,
			Name:     attr.Name,
			Type:     fieldType,
			Value:    formData[attr.Name],
			Required: attr.Required,
			Help:     attr.Help,
		})
		fields = append(fields, field)
	}

	formObj := form.NewForm(form.FormOptions{
		ID:     "FormEntityEdit",
		Fields: fields,
	})

	modalID := "ModalEntityEdit"
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

func (c *EntityEditController) modalHeader() hb.TagInterface {
	return hb.Div().Class("modal-header").
		Child(hb.Heading5().Class("modal-title").HTML("Edit " + c.definition.TypeLabel)).
		Child(hb.Button().
			Type("button").
			Class("btn-close").
			Data("bs-dismiss", "modal").
			Attr("aria-label", "Close"))
}

func (c *EntityEditController) modalBody(formObj *form.Form, submitUrl, errorMessage string) hb.TagInterface {
	body := hb.Div().Class("modal-body")

	if errorMessage != "" {
		body.Child(hb.Div().Class("alert alert-danger").HTML(errorMessage))
	}

	body.Child(hb.Form().
		ID("FormEntityEdit").
		Method(http.MethodPost).
		Action(submitUrl).
		Child(formObj.Build()))

	return body
}

func (c *EntityEditController) modalFooter(_ string) hb.TagInterface {
	return hb.Div().Class("modal-footer").
		Child(hb.Button().
			Type("button").
			Class("btn btn-secondary").
			Data("bs-dismiss", "modal").
			HTML("Close")).
		Child(hb.Button().
			Type("submit").
			Class("btn btn-primary").
			Attr("form", "FormEntityEdit").
			HTML("Update " + c.definition.TypeLabel))
}

func (c *EntityEditController) handleSubmit(_ http.ResponseWriter, r *http.Request, entityID string) string {
	ctx := context.Background()

	if err := r.ParseForm(); err != nil {
		return c.showForm(r, entityID, nil, "Error parsing form: "+err.Error())
	}

	customStore := c.ui.Store().CustomEntityStore()
	if customStore == nil {
		return c.showForm(r, entityID, nil, "Custom entity store is not available")
	}

	// Load entity
	entity, err := customStore.FindByID(ctx, entityID)
	if err != nil || entity == nil {
		return c.showForm(r, entityID, nil, "Entity not found")
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

	// Update entity
	err = customStore.Update(ctx, entity, attrs)
	if err != nil {
		return c.showForm(r, entityID, formData, "Error updating entity: "+err.Error())
	}

	// Success - close modal and refresh page
	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{
			Icon: "success",
			Text: c.definition.TypeLabel + " updated successfully",
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

func (c *EntityEditController) getFieldType(attrType string) string {
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

func (c *EntityEditController) errorResponse(message string) string {
	return hb.Swal(hb.SwalOptions{
		Icon: "error",
		Text: message,
	}).ToHTML()
}
