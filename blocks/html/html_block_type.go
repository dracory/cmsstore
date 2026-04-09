package html

import (
	"context"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
)

// HTMLBlockType provides both frontend rendering and admin UI for HTML blocks.
//
// This is a built-in block type that renders raw HTML content.
// It includes a CodeMirror editor in the admin UI for syntax highlighting.
type HTMLBlockType struct{}

// NewHTMLBlockType creates a new HTML block type.
func NewHTMLBlockType() *HTMLBlockType {
	return &HTMLBlockType{}
}

// TypeKey returns the unique identifier for HTML blocks.
func (t *HTMLBlockType) TypeKey() string {
	return cmsstore.BLOCK_TYPE_HTML
}

// TypeLabel returns the display name for HTML blocks.
func (t *HTMLBlockType) TypeLabel() string {
	return "HTML Block"
}

// Render renders an HTML block by returning its content as-is.
// Supports optional 'wrap' attribute to wrap content in an HTML element.
func (t *HTMLBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, opts ...cmsstore.RenderOption) (string, error) {
	if block == nil {
		return "<!-- Block is nil -->", nil
	}
	content := block.Content()
	if content == "" {
		return "<!-- Empty block content -->", nil
	}

	// Parse render options
	options := &cmsstore.RenderOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Apply wrapper if specified
	if wrapper := options.Attributes["wrap"]; wrapper != "" {
		return "<" + wrapper + ">" + content + "</" + wrapper + ">", nil
	}

	return content, nil
}

// GetAdminFields returns form fields for editing HTML block content.
func (t *HTMLBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Content (HTML)",
			Name:  "block_content",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: block.Content(),
		}),
	}

	// Add CodeMirror initialization script for syntax highlighting
	contentScript := hb.Script(`
function codeMirrorSelector() {
	return 'textarea[name="block_content"]';
}
function getCodeMirrorEditor() {
	return document.querySelector(codeMirrorSelector());
}
setTimeout(function () {
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

// GetCustomVariables returns nil as HTML blocks do not set any custom variables.
func (t *HTMLBlockType) GetCustomVariables() []cmsstore.BlockCustomVariable {
	return nil
}

// SaveAdminFields processes form submission and updates the HTML block.
func (t *HTMLBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	content := req.GetStringTrimmed(r, "block_content")
	block.SetContent(content)
	return nil
}
