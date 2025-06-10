package rest

import (
	"net/http"
	"strings"

	"github.com/gouniverse/cmsstore"
)

// RestAPI represents the REST API for the CMS store
type RestAPI struct {
	store cmsstore.StoreInterface
}

// NewRestAPI creates a new REST API instance
func NewRestAPI(store cmsstore.StoreInterface) *RestAPI {
	return &RestAPI{
		store: store,
	}
}

// Handler returns an http.HandlerFunc that can be attached to any router
func (api *RestAPI) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set common headers
		w.Header().Set("Content-Type", "application/json")

		// Parse the path to determine the resource and action
		path := strings.TrimPrefix(r.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		if len(pathParts) < 2 {
			http.Error(w, `{"success":false,"error":"Invalid API path"}`, http.StatusBadRequest)
			return
		}

		// Check if this is an API request
		if pathParts[0] != "api" {
			http.Error(w, `{"success":false,"error":"Not an API request"}`, http.StatusBadRequest)
			return
		}

		// Handle different resources
		switch pathParts[1] {
		case "pages":
			api.handlePagesEndpoint(w, r, pathParts[2:])
		case "menus":
			api.handleMenusEndpoint(w, r, pathParts[2:])
		case "sites":
			api.handleSitesEndpoint(w, r, pathParts[2:])
		case "templates":
			api.handleTemplatesEndpoint(w, r, pathParts[2:])
		case "blocks":
			api.handleBlocksEndpoint(w, r, pathParts[2:])
		case "translations":
			api.handleTranslationsEndpoint(w, r, pathParts[2:])
		default:
			http.Error(w, `{"success":false,"error":"Unknown resource"}`, http.StatusNotFound)
		}
	}
}
