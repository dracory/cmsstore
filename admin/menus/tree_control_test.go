package admin

import (
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func initTreeControl(treeJSON string, renderURL string, targetTextareaID string) *treeControl {
	return &treeControl{
		id:               "test-tree",
		renderURL:        renderURL,
		treeJSON:         treeJSON,
		targetTextareaID: targetTextareaID,
		pageList:         []cmsstore.PageInterface{},
	}
}

func Test_TreeControl_Render_EmptyTree(t *testing.T) {
	control := initTreeControl("[]", "/test", "menu_items")

	req, _ := test.NewRequest("GET", "/test", test.NewRequestOptions{})

	body := control.Render(req)

	assert.Contains(t, body.ToHTML(), "No menu items")
	assert.Contains(t, body.ToHTML(), "New Menu Item")
}

func Test_TreeControl_Render_WithNodes(t *testing.T) {
	treeJSON := `[
		{"id":"1","name":"Home","page_id":"","parent_id":"","sequence":0,"target":"","url":"/"},
		{"id":"2","name":"About","page_id":"","parent_id":"","sequence":1,"target":"","url":"/about"}
	]`

	control := initTreeControl(treeJSON, "/test", "menu_items")

	req, _ := test.NewRequest("GET", "/test", test.NewRequestOptions{})

	body := control.Render(req)

	assert.Contains(t, body.ToHTML(), "Home")
	assert.Contains(t, body.ToHTML(), "About")
	assert.Contains(t, body.ToHTML(), "New Menu Item")
	assert.Contains(t, body.ToHTML(), "tree")
}

func Test_TreeControl_Render_NodeAdd(t *testing.T) {
	treeJSON := `[
		{"id":"1","name":"Home","page_id":"","parent_id":"","sequence":0,"target":"","url":"/"}
	]`

	control := initTreeControl(treeJSON, "/test", "menu_items")

	req, _ := test.NewRequest("POST", "/test", test.NewRequestOptions{
		GetValues: map[string][]string{
			"treectl_action":    {"node_add"},
			"treectl_parent_id": {"1"},
		},
	})

	body := control.Render(req)

	assert.Contains(t, body.ToHTML(), "tree")
	assert.NotContains(t, body.ToHTML(), "ERROR:")
}

func Test_TreeControl_Render_NodeDelete(t *testing.T) {
	treeJSON := `[
		{"id":"1","name":"Home","page_id":"","parent_id":"","sequence":0,"target":"","url":"/"},
		{"id":"2","name":"About","page_id":"","parent_id":"","sequence":1,"target":"","url":"/about"}
	]`

	control := initTreeControl(treeJSON, "/test", "menu_items")

	req, _ := test.NewRequest("POST", "/test", test.NewRequestOptions{
		GetValues: map[string][]string{
			"treectl_action":  {"node_delete"},
			"treectl_node_id": {"2"},
		},
	})

	body := control.Render(req)

	assert.Contains(t, body.ToHTML(), "tree")
	assert.NotContains(t, body.ToHTML(), "About") // Node should be deleted
}

func Test_TreeControl_Render_NodeUpdateModal(t *testing.T) {
	treeJSON := `[
		{"id":"1","name":"Home","page_id":"","parent_id":"","sequence":0,"target":"","url":"/"}
	]`

	control := initTreeControl(treeJSON, "/test", "menu_items")

	req, _ := test.NewRequest("POST", "/test", test.NewRequestOptions{
		GetValues: map[string][]string{
			"treectl_action":  {"node_update_modal"},
			"treectl_node_id": {"1"},
		},
	})

	body := control.Render(req)

	assert.Contains(t, body.ToHTML(), "New Menu")
	assert.Contains(t, body.ToHTML(), "Menu Item")
	assert.Contains(t, body.ToHTML(), "ModalNodeUpdate")
	assert.Contains(t, body.ToHTML(), "modal-backdrop")
}

func Test_TreeControl_Render_NodeUpdate(t *testing.T) {
	treeJSON := `[
		{"id":"1","name":"Home","page_id":"","parent_id":"","sequence":0,"target":"","url":"/"}
	]`

	control := initTreeControl(treeJSON, "/test", "menu_items")

	req, _ := test.NewRequest("POST", "/test", test.NewRequestOptions{
		GetValues: map[string][]string{
			"treectl_action":  {"node_update"},
			"treectl_node_id": {"1"},
			"treectl_name":    {"Updated Home"},
			"treectl_url":     {"/updated"},
		},
	})

	body := control.Render(req)

	assert.Contains(t, body.ToHTML(), "tree")
	assert.NotContains(t, body.ToHTML(), "ERROR:")
}

func Test_TreeControl_Render_NodeUpdate_NotFound(t *testing.T) {
	treeJSON := `[
		{"id":"1","name":"Home","page_id":"","parent_id":"","sequence":0,"target":"","url":"/"}
	]`

	control := initTreeControl(treeJSON, "/test", "menu_items")

	req, _ := test.NewRequest("POST", "/test", test.NewRequestOptions{
		GetValues: map[string][]string{
			"treectl_action":  {"node_update"},
			"treectl_node_id": {"999"},
		},
	})

	body := control.Render(req)

	assert.Contains(t, body.ToHTML(), "ERROR: Node not found")
}

func Test_TreeControl_Render_InvalidJSON(t *testing.T) {
	control := initTreeControl("invalid json", "/test", "menu_items")

	req, _ := test.NewRequest("GET", "/test", test.NewRequestOptions{})

	body := control.Render(req)

	assert.Contains(t, body.ToHTML(), "ERROR:")
}

func Test_TreeControl_Render_IDGeneration(t *testing.T) {
	control := &treeControl{
		renderURL:        "/test",
		treeJSON:         "[]",
		targetTextareaID: "menu_items",
		pageList:         []cmsstore.PageInterface{},
	}

	req, _ := test.NewRequest("GET", "/test", test.NewRequestOptions{})

	body := control.Render(req)

	// Should generate an ID when none is provided
	assert.Contains(t, body.ToHTML(), "treectl_")
}

func Test_TreeControl_Render_WithExistingID(t *testing.T) {
	control := &treeControl{
		id:               "existing-id",
		renderURL:        "/test",
		treeJSON:         "[]",
		targetTextareaID: "menu_items",
		pageList:         []cmsstore.PageInterface{},
	}

	req, _ := test.NewRequest("GET", "/test", test.NewRequestOptions{})

	body := control.Render(req)

	// Should use the existing ID
	assert.Contains(t, body.ToHTML(), "existing-id")
}

func Test_TreeControl_Render_JSONOutput(t *testing.T) {
	treeJSON := `[
		{"id":"1","name":"Home","page_id":"","parent_id":"","sequence":0,"target":"","url":"/"}
	]`

	control := initTreeControl(treeJSON, "/test", "menu_items")

	req, _ := test.NewRequest("GET", "/test", test.NewRequestOptions{})

	body := control.Render(req)

	// Should include JavaScript to update the textarea
	assert.Contains(t, body.ToHTML(), "JSON.stringify")
	assert.Contains(t, body.ToHTML(), "menu_items")
}

func Test_TreeControl_Render_JSONOutput2(t *testing.T) {
	treeJSON := `[
		{"id":"1","name":"Home","page_id":"","parent_id":"","sequence":0,"target":"","url":"/"}
	]`

	control := initTreeControl(treeJSON, "/test", "menu_items")

	req, _ := test.NewRequest("GET", "/test", test.NewRequestOptions{})

	body := control.Render(req)

	// Should include JavaScript to update the textarea
	assert.Contains(t, body.ToHTML(), "JSON.stringify")
	assert.Contains(t, body.ToHTML(), "menu_items")
}
