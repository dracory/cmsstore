package frontend

import "net/http"

type FrontendInterface interface {
	Handler(w http.ResponseWriter, r *http.Request)
	StringHandler(w http.ResponseWriter, r *http.Request) string
	TemplateRenderHtmlByID(r *http.Request, templateID string, options struct {
		PageContent         string
		PageCanonicalURL    string
		PageMetaDescription string
		PageMetaKeywords    string
		PageMetaRobots      string
		PageTitle           string
		Language            string
	}) (string, error)
}
