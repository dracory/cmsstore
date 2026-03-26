package admin

import (
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

// HTMLAdminProvider provides admin UI for HTML block types.
type HTMLAdminProvider struct{}

// NewHTMLAdminProvider creates a new HTML block admin provider.
func NewHTMLAdminProvider() *HTMLAdminProvider {
	return &HTMLAdminProvider{}
}

// GetContentFields returns form fields for HTML block content editing.
func (p *HTMLAdminProvider) GetContentFields(block cmsstore.BlockInterface, r *http.Request) []form.FieldInterface {
	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Content (HTML)",
			Name:  "block_content",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: block.Content(),
		}),
	}

	// Add CodeMirror initialization script
	contentScript := hb.Script(`
function codeMirrorSelector() {
	return 'textarea[name="block_content"]';
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

// GetTypeLabel returns the display label for HTML blocks.
func (p *HTMLAdminProvider) GetTypeLabel() string {
	return "HTML Block"
}

// SaveContentFields processes form data and updates the HTML block.
func (p *HTMLAdminProvider) SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error {
	content := req.GetStringTrimmed(r, "block_content")
	block.SetContent(content)
	return nil
}
