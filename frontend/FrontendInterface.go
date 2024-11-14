package frontend

import "net/http"

type FrontendInterface interface {
	Handler(w http.ResponseWriter, r *http.Request)
	StringHandler(w http.ResponseWriter, r *http.Request) string
}
