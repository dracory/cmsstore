package customentity

import (
	"context"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
)

// EntityDeleteController handles entity deletion
type EntityDeleteController struct {
	ui         UiInterface
	definition cmsstore.CustomEntityDefinition
}

// NewEntityDeleteController creates a new delete controller
func NewEntityDeleteController(ui UiInterface, definition cmsstore.CustomEntityDefinition) *EntityDeleteController {
	return &EntityDeleteController{
		ui:         ui,
		definition: definition,
	}
}

// Handler handles both GET (show confirmation) and POST (process deletion)
func (c *EntityDeleteController) Handler(w http.ResponseWriter, r *http.Request) string {
	entityID := r.URL.Query().Get("id")
	if entityID == "" {
		return c.errorResponse("Entity ID is required")
	}

	if r.Method == http.MethodPost {
		return c.handleDelete(w, r, entityID)
	}
	return c.showConfirmation(r, entityID)
}

func (c *EntityDeleteController) showConfirmation(_ *http.Request, entityID string) string {
	ctx := context.Background()
	customStore := c.ui.Store().CustomEntityStore()

	// Load entity to verify it exists
	entity, err := customStore.FindByID(ctx, entityID)
	if err != nil || entity == nil {
		return c.errorResponse("Entity not found")
	}

	// Get entity display name (first attribute value)
	displayName := entityID
	if len(c.definition.Attributes) > 0 {
		firstAttr := c.definition.Attributes[0]
		attrValue, _ := customStore.Inner().AttributeFind(ctx, entityID, firstAttr.Name)
		if attrValue != nil && attrValue.GetValue() != "" {
			displayName = attrValue.GetValue()
		}
	}

	deleteUrl := "/admin/custom-entity/" + c.definition.Type + "/delete?id=" + entityID

	modalID := "ModalEntityDelete"
	modalBackdropClass := "ModalBackdrop"

	modal := hb.Div().
		ID(modalID).
		Class("modal fade show").
		Style("display: block;").
		Child(hb.Div().Class("modal-dialog").
			Child(hb.Div().Class("modal-content").
				Child(c.modalHeader()).
				Child(c.modalBody(displayName, deleteUrl)).
				Child(c.modalFooter(modalID, deleteUrl))))

	backdrop := hb.Div().Class(modalBackdropClass).Style("display: block;")

	return hb.Wrap().Children([]hb.TagInterface{modal, backdrop}).ToHTML()
}

func (c *EntityDeleteController) modalHeader() hb.TagInterface {
	return hb.Div().Class("modal-header bg-danger text-white").
		Child(hb.Heading5().Class("modal-title").HTML("Delete " + c.definition.TypeLabel)).
		Child(hb.Button().
			Type("button").
			Class("btn-close btn-close-white").
			Data("bs-dismiss", "modal").
			Attr("aria-label", "Close"))
}

func (c *EntityDeleteController) modalBody(displayName, deleteUrl string) hb.TagInterface {
	body := hb.Div().Class("modal-body")

	body.Child(hb.Div().Class("alert alert-warning").
		Child(hb.I().Class("bi bi-exclamation-triangle me-2")).
		Child(hb.Span().HTML("This action cannot be undone.")))

	body.Child(hb.Paragraph().HTML("Are you sure you want to delete this " + c.definition.TypeLabel + "?"))

	body.Child(hb.Div().Class("mb-3").
		Child(hb.Strong().HTML("Name: ")).
		Child(hb.Span().HTML(displayName)))

	body.Child(hb.Form().
		ID("FormEntityDelete").
		Method(http.MethodPost).
		Action(deleteUrl))

	return body
}

func (c *EntityDeleteController) modalFooter(_, _ string) hb.TagInterface {
	return hb.Div().Class("modal-footer").
		Child(hb.Button().
			Type("button").
			Class("btn btn-secondary").
			Data("bs-dismiss", "modal").
			HTML("Cancel")).
		Child(hb.Button().
			Type("submit").
			Class("btn btn-danger").
			Attr("form", "FormEntityDelete").
			Child(hb.I().Class("bi bi-trash me-2")).
			Child(hb.Span().HTML("Delete")))
}

func (c *EntityDeleteController) handleDelete(_ http.ResponseWriter, _ *http.Request, entityID string) string {
	ctx := context.Background()
	customStore := c.ui.Store().CustomEntityStore()

	if customStore == nil {
		return c.errorResponse("Custom entity store is not available")
	}

	// Delete entity
	err := customStore.Delete(ctx, entityID)
	if err != nil {
		return c.errorResponse("Error deleting entity: " + err.Error())
	}

	// Success - close modal and refresh page
	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{
			Icon: "success",
			Text: c.definition.TypeLabel + " deleted successfully",
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

func (c *EntityDeleteController) errorResponse(message string) string {
	return hb.Swal(hb.SwalOptions{
		Icon: "error",
		Text: message,
	}).ToHTML()
}
