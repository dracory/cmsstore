package shared

import "net/http"

type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
