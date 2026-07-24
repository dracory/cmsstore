package frontend

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

// MediaHandler serves CMS media files via /cms/media/{mediaId}.{ext}
// It looks up the media record by ID and serves the content from whatever
// URL() the record holds (data URI, HTTP URL, file path — doesn't matter).
func (frontend *frontend) MediaHandler(w http.ResponseWriter, r *http.Request) string {
	if frontend.store == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "CMS store not configured"
	}

	mediaID := extractMediaID(r.URL.Path)
	if mediaID == "" {
		w.WriteHeader(http.StatusNotFound)
		return "Media not found"
	}

	media, err := frontend.store.MediaFindByID(context.Background(), mediaID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err.Error()
	}

	if media == nil {
		w.WriteHeader(http.StatusNotFound)
		return "Media not found"
	}

	if !media.IsActive() {
		w.WriteHeader(http.StatusNotFound)
		return "Media not found"
	}

	url := media.URL()

	// Data URI — decode base64 and serve directly
	if strings.HasPrefix(url, "data:") {
		return frontend.serveDataURI(w, r, url, media)
	}

	// HTTP(S) URL — redirect
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		http.Redirect(w, r, url, http.StatusFound)
		return ""
	}

	// File path — attempt to read and serve
	return frontend.serveFilePath(w, r, url, media)
}

func (frontend *frontend) serveDataURI(w http.ResponseWriter, r *http.Request, dataURI string, media interface {
	Type() string
	Extension() string
	Size() string
	ID() string
}) string {
	content, err := decodeDataURL(dataURI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "Failed to decode media content"
	}

	contentType := media.Type()
	if contentType == "" {
		contentType = mimeTypeFromExtension(media.Extension())
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", media.Size())
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header().Set("ETag", media.ID())

	if _, err := w.Write(content); err != nil {
		return "Failed to write media content: " + err.Error()
	}

	return ""
}

func (frontend *frontend) serveFilePath(w http.ResponseWriter, r *http.Request, filePath string, media interface {
	Type() string
	Extension() string
	ID() string
}) string {
	contentType := media.Type()
	if contentType == "" {
		contentType = mimeTypeFromExtension(media.Extension())
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header().Set("ETag", media.ID())

	http.ServeFile(w, r, filePath)
	return ""
}

// extractMediaID parses the media ID from URL paths in either format:
//   - /cms/media/<id>.<ext>           (e.g. /cms/media/1jq6fby3kzj.png)
//   - /cms/media/<id>/<handle>.<ext>  (e.g. /cms/media/1jq6fby3kzj/dulydo.png)
//
// The handle in the second form is cosmetic — lookup is always by ID.
func extractMediaID(urlPath string) string {
	path := strings.TrimPrefix(urlPath, "/cms/media/")
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	// If there's a slash, the ID is the first segment
	if idx := strings.Index(path, "/"); idx > 0 {
		return path[:idx]
	}

	// No slash — strip the extension to get the ID
	if idx := strings.LastIndex(path, "."); idx >= 0 {
		return path[:idx]
	}
	return path
}

// isMediaURL checks if the given path is a CMS media URL.
func isMediaURL(path string) bool {
	return strings.HasPrefix(path, "/cms/media/")
}

// decodeDataURL decodes a base64 data URL (e.g. "data:image/png;base64,...")
// and returns the raw binary content.
func decodeDataURL(dataURL string) ([]byte, error) {
	idx := strings.Index(dataURL, "base64,")
	if idx < 0 {
		return nil, errInvalidDataURL
	}
	b64data := dataURL[idx+7:]
	return base64.StdEncoding.DecodeString(b64data)
}

// mimeTypeFromExtension returns the MIME type for a given file extension.
func mimeTypeFromExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

var errInvalidDataURL = errors.New("invalid data URL: missing base64 prefix")
