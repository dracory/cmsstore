package frontend

import "net/http"

type FrontendInterface interface {
	// Handler renders the frontend
	Handler(w http.ResponseWriter, r *http.Request)

	// StringHandler return the frontend as a HTML string
	StringHandler(w http.ResponseWriter, r *http.Request) string

	// TemplateRenderHtmlByID builds the HTML of a template based on its ID
	TemplateRenderHtmlByID(r *http.Request, templateID string, options TemplateRenderHtmlByIDOptions) (string, error)
}

type TemplateRenderHtmlByIDOptions struct {
	PageContent         string
	PageCanonicalURL    string
	PageMetaDescription string
	PageMetaKeywords    string
	PageMetaRobots      string
	PageTitle           string
	Language            string
}
