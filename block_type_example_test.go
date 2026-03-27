package cmsstore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
	"github.com/dracory/req"
)

// GalleryBlockType is a complete example of a custom block type.
//
// This single struct implements both frontend rendering and admin UI,
// ensuring they stay in sync and preventing registration errors.
type GalleryBlockType struct {
	store cmsstore.StoreInterface
}

// TypeKey returns the unique identifier stored in the database.
func (t *GalleryBlockType) TypeKey() string {
	return "gallery"
}

// TypeLabel returns the display name shown in the admin UI.
func (t *GalleryBlockType) TypeLabel() string {
	return "Gallery Block"
}

// Render implements frontend rendering logic.
func (t *GalleryBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	// Parse gallery data from block content
	var images []GalleryImage
	if err := json.Unmarshal([]byte(block.Content()), &images); err != nil {
		return "<!-- Invalid gallery data -->", nil
	}

	// Get configuration from metadata
	layout := block.Meta("layout")
	if layout == "" {
		layout = "grid"
	}
	columns := block.Meta("columns")
	if columns == "" {
		columns = "3"
	}

	// Generate HTML
	html := fmt.Sprintf(`<div class="gallery gallery-layout-%s" data-columns="%s">`, layout, columns)
	for _, img := range images {
		html += fmt.Sprintf(`
  <div class="gallery-item">
    <img src="%s" alt="%s" loading="lazy">
    <div class="caption">%s</div>
  </div>`, img.URL, img.Alt, img.Caption)
	}
	html += `</div>
<script>
  // Initialize gallery (use your preferred library)
  document.querySelectorAll('.gallery').forEach(gallery => {
    initGallery(gallery);
  });
</script>`

	return html, nil
}

// GetAdminFields implements admin UI form fields.
func (t *GalleryBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	return []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label:    "Gallery Images (JSON)",
			Name:     "gallery_images",
			Type:     form.FORM_FIELD_TYPE_TEXTAREA,
			Value:    block.Content(),
			Required: true,
			Help:     "Enter images as JSON array: [{\"url\":\"...\",\"alt\":\"...\",\"caption\":\"...\"}]",
		}),
		form.NewField(form.FieldOptions{
			Label: "Layout Style",
			Name:  "gallery_layout",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("layout"),
			Options: []form.FieldOption{
				{Value: "Grid", Key: "grid"},
				{Value: "Masonry", Key: "masonry"},
				{Value: "Carousel", Key: "carousel"},
				{Value: "Slideshow", Key: "slideshow"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Columns (for grid layout)",
			Name:  "gallery_columns",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta("columns"),
			Options: []form.FieldOption{
				{Value: "2 columns", Key: "2"},
				{Value: "3 columns", Key: "3"},
				{Value: "4 columns", Key: "4"},
				{Value: "6 columns", Key: "6"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Enable Lightbox",
			Name:  "gallery_lightbox",
			Type:  form.FORM_FIELD_TYPE_CHECKBOX,
			Value: block.Meta("lightbox"),
			Help:  "Open images in a lightbox when clicked",
		}),
	}
}

// SaveAdminFields implements form submission handling.
func (t *GalleryBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	// Get form values
	images := req.GetStringTrimmed(r, "gallery_images")
	layout := req.GetStringTrimmed(r, "gallery_layout")
	columns := req.GetStringTrimmed(r, "gallery_columns")
	lightbox := req.GetStringTrimmed(r, "gallery_lightbox")

	// Validate JSON
	if images == "" {
		return &ValidationError{Message: "Gallery images are required"}
	}

	var imageData []GalleryImage
	if err := json.Unmarshal([]byte(images), &imageData); err != nil {
		return &ValidationError{Message: "Invalid JSON format: " + err.Error()}
	}

	if len(imageData) == 0 {
		return &ValidationError{Message: "At least one image is required"}
	}

	// Save to block
	block.SetContent(images)
	block.SetMeta("layout", layout)
	block.SetMeta("columns", columns)
	block.SetMeta("lightbox", lightbox)

	return nil
}

// GalleryImage represents a single image in the gallery.
type GalleryImage struct {
	URL     string `json:"url"`
	Alt     string `json:"alt"`
	Caption string `json:"caption"`
}

// ValidationError represents a form validation error.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ExampleRegisterCustomBlockType shows how to register a custom block type.
//
// This single registration makes the block type available in both:
//   - Frontend rendering (when blocks are displayed on pages)
//   - Admin UI (block type dropdown, edit forms, save logic)
func ExampleRegisterCustomBlockType() {
	// Create your block type instance
	galleryType := &GalleryBlockType{
		store: nil, // Pass your store instance here
	}

	// Register it globally - that's it!
	cmsstore.RegisterCustomBlockType(galleryType)

	// Now the "gallery" block type is fully functional:
	// - Appears in admin UI block type dropdown as "Gallery Block"
	// - Shows custom form fields when editing
	// - Renders with custom HTML on the frontend
	// - Validates and saves form data correctly

	fmt.Println("Gallery block type registered successfully")
	// Output: Gallery block type registered successfully
}

// ExampleBlockType_interactiveWithVue shows a Vue.js interactive block.
type VueTreeBlockType struct{}

func (t *VueTreeBlockType) TypeKey() string {
	return "vue_tree"
}

func (t *VueTreeBlockType) TypeLabel() string {
	return "Interactive Tree (Vue.js)"
}

func (t *VueTreeBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	treeData := block.Content()
	blockID := block.ID()
	cssClass := block.Meta("css_class")

	return fmt.Sprintf(`
<div id="vue-tree-%s" class="%s"></div>

<script type="module">
import { createApp } from 'https://unpkg.com/vue@3/dist/vue.esm-browser.js'

createApp({
  data() {
    return {
      treeData: %s,
      expanded: {}
    }
  },
  methods: {
    toggleNode(nodeId) {
      this.expanded[nodeId] = !this.expanded[nodeId]
    }
  },
  template: `+"`"+`
    <div class="tree-container">
      <tree-node 
        v-for="node in treeData" 
        :key="node.id"
        :node="node"
        :expanded="expanded"
        @toggle="toggleNode"
      />
    </div>
  `+"`"+`
}).mount('#vue-tree-%s')
</script>
`, blockID, cssClass, treeData, blockID), nil
}

func (t *VueTreeBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	return []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label:    "Tree Data (JSON)",
			Name:     "tree_data",
			Type:     form.FORM_FIELD_TYPE_TEXTAREA,
			Value:    block.Content(),
			Required: true,
			Help:     "Tree structure in JSON format",
		}),
		form.NewField(form.FieldOptions{
			Label: "CSS Class",
			Name:  "tree_css_class",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta("css_class"),
			Help:  "Custom CSS class for styling",
		}),
	}
}

func (t *VueTreeBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	treeData := req.GetStringTrimmed(r, "tree_data")
	cssClass := req.GetStringTrimmed(r, "tree_css_class")

	// Validate JSON
	if treeData != "" {
		var test interface{}
		if err := json.Unmarshal([]byte(treeData), &test); err != nil {
			return &ValidationError{Message: "Invalid JSON: " + err.Error()}
		}
	}

	block.SetContent(treeData)
	block.SetMeta("css_class", cssClass)
	return nil
}
