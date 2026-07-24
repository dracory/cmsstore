package page_update

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/cmsstore"
	"github.com/dracory/req"
)

func handleAjaxLoadMedia(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	files, err := store.MediaListByEntityID(r.Context(), pageID, "page")
	if err != nil {
		slog.Error("Failed to load page media", "error", err)
		return api.Error("Failed to load media").ToString()
	}

	fileList := []map[string]any{}
	for _, f := range files {
		fileList = append(fileList, map[string]any{
			"id":        f.ID(),
			"name":      f.Title(),
			"url":       f.URL(),
			"type":      f.Type(),
			"size":      f.Size(),
			"extension": f.Extension(),
			"sequence":  f.SequenceInt(),
		})
	}

	return api.SuccessWithData("Media loaded successfully", map[string]any{
		"files": fileList,
	}).ToString()
}

func handleAjaxUploadMedia(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		return api.Error("Failed to parse upload: " + err.Error()).ToString()
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		files = r.MultipartForm.File["upload_file"]
	}
	if len(files) == 0 {
		return api.Error("No files uploaded").ToString()
	}

	existingFiles, _ := store.MediaListByEntityID(r.Context(), pageID, "page")
	startSequence := len(existingFiles)

	uploaded := []map[string]any{}

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return api.Error("Failed to open file: " + err.Error()).ToString()
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return api.Error("Failed to read file: " + err.Error()).ToString()
		}

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		dataURI := "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(data)

		media := cmsstore.NewMedia().
			SetEntityID(pageID).
			SetEntityType("page").
			SetTitle(fileHeader.Filename).
			SetURL(dataURI).
			SetType(contentType).
			SetSize(strconv.FormatInt(fileHeader.Size, 10)).
			SetExtension(ext).
			SetSequenceInt(startSequence + i).
			SetStatus(cmsstore.MEDIA_STATUS_ACTIVE)

		if err := store.MediaCreate(r.Context(), media); err != nil {
			return api.Error("Failed to save file record: " + err.Error()).ToString()
		}

		uploaded = append(uploaded, map[string]any{
			"id":        media.ID(),
			"name":      media.Title(),
			"url":       media.URL(),
			"type":      media.Type(),
			"size":      media.Size(),
			"extension": media.Extension(),
			"sequence":  media.SequenceInt(),
		})
	}

	return api.SuccessWithData("Files uploaded successfully", map[string]any{
		"files": uploaded,
	}).ToString()
}

func handleAjaxSaveMedia(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	var reqData struct {
		PageID string `json:"page_id"`
		Files  []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Sequence int    `json:"sequence"`
		} `json:"files"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	for _, item := range reqData.Files {
		media, err := store.MediaFindByID(r.Context(), item.ID)
		if err != nil || media == nil {
			continue
		}
		if item.Name != "" {
			media.SetTitle(item.Name)
		}
		media.SetSequenceInt(item.Sequence)
		if err := store.MediaUpdate(r.Context(), media); err != nil {
			slog.Error("Failed to update media", "error", err, "media_id", item.ID)
		}
	}

	return api.Success("Media saved successfully").ToString()
}

func handleAjaxDeleteMedia(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	fileID := req.GetStringTrimmed(r, "file_id")
	if fileID == "" {
		return api.Error("file_id is required").ToString()
	}

	media, err := store.MediaFindByID(r.Context(), fileID)
	if err != nil {
		return api.Error("Failed to find media: " + err.Error()).ToString()
	}
	if media == nil {
		return api.Error("Media not found").ToString()
	}

	if err := store.MediaDeleteByID(r.Context(), fileID); err != nil {
		return api.Error("Failed to delete media: " + err.Error()).ToString()
	}

	return api.Success("Media deleted successfully").ToString()
}

func handleAjaxAddMedia(store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	pageID := reqGetString(r, "page_id")
	if pageID == "" {
		return api.Error("Page ID is required").ToString()
	}

	mediaURL := req.GetStringTrimmed(r, "media_url")
	if mediaURL == "" {
		return api.Error("media_url is required").ToString()
	}

	mediaFileName := req.GetStringTrimmed(r, "media_file_name")
	mediaType := req.GetStringTrimmed(r, "media_type")

	existingFiles, _ := store.MediaListByEntityID(r.Context(), pageID, "page")
	startSequence := len(existingFiles)

	media := cmsstore.NewMedia().
		SetEntityID(pageID).
		SetEntityType("page").
		SetTitle(mediaFileName).
		SetURL(mediaURL).
		SetType(mediaType).
		SetSequenceInt(startSequence).
		SetStatus(cmsstore.MEDIA_STATUS_ACTIVE)

	if err := store.MediaCreate(r.Context(), media); err != nil {
		return api.Error("Failed to save media: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Media added successfully", map[string]any{
		"file": map[string]any{
			"id":       media.ID(),
			"name":     media.Title(),
			"url":      media.URL(),
			"type":     media.Type(),
			"sequence": media.SequenceInt(),
		},
	}).ToString()
}
