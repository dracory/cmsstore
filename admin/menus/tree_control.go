package admin

import (
	"net/http"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

type treeControl struct {
	id               string
	renderURL        string
	treeJSON         string
	targetTextareaID string
	pageList         []cmsstore.PageInterface
}

func (t *treeControl) Render(r *http.Request) hb.TagInterface {
	tree, err := NewTreeFromJSON(t.treeJSON)

	if err != nil {
		return hb.Div().Text(`ERROR: ` + err.Error())
	}

	// If no ID is provided, generate one, while making sure
	// we retain any existing ID
	if t.id == "" {
		treeID := utils.Req(r, "treectl_id", "")
		if treeID != "" {
			t.id = treeID
		} else {
			t.id = "treectl_" + uid.HumanUid()
		}
	}

	action := utils.Req(r, "treectl_action", "")

	if action == "node_add" {
		parentID := utils.Req(r, "treectl_parent_id", "")
		newNode := Node{ID: uid.HumanUid(), Name: "New"}
		tree.Add(parentID, newNode)
	}

	if action == "node_delete" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		tree.Remove(nodeID)
	}

	if action == "node_move_up" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		tree.MoveUp(nodeID)
	}

	if action == "node_move_down" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		tree.MoveDown(nodeID)
	}

	if action == "node_move_out_up" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		parent := tree.Parent(nodeID)
		if parent != nil {
			tree.MoveToPosition(nodeID, parent.ParentID, parent.Sequence)
		}
	}

	if action == "node_move_out_down" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		parent := tree.Parent(nodeID)
		if parent != nil {
			tree.MoveToPosition(nodeID, parent.ParentID, parent.Sequence+1)
		}
	}

	if action == "node_move_in_up" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		sibling := tree.FindPreviousSibling(nodeID)
		if sibling != nil {
			tree.MoveToParent(nodeID, sibling.ID)
		}
	}

	if action == "node_move_in_down" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		sibling := tree.FindNextSibling(nodeID)
		if sibling != nil {
			tree.MoveToParent(nodeID, sibling.ID)
		}
	}

	if action == "node_update_modal" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		node := tree.Find(nodeID)
		if node == nil {
			return hb.Div().Text(`ERROR: Node not found`)
		}
		return t.modalNodeUpdate(*node)
	}

	if action == "node_update" {
		nodeID := utils.Req(r, "treectl_node_id", "")
		pageID := utils.Req(r, "treectl_page_id", "")
		url := utils.Req(r, "treectl_url", "")
		target := utils.Req(r, "treectl_target", "")
		name := utils.Req(r, "treectl_name", "")
		node := tree.Find(nodeID)

		if node == nil {
			return hb.Div().Text(`ERROR: Node not found`)
		}

		node.PageID = pageID
		node.URL = url
		node.Target = target
		node.Name = name

		tree.Update(*node)
	}

	return t.renderControl(*tree)
}

func (t *treeControl) renderControl(tree Tree) hb.TagInterface {
	jsonString, err := tree.ToJSON()

	if err != nil {
		return hb.Div().Text(`ERROR: ` + err.Error())
	}

	buttonAddNode := hb.Button().
		ID("ButtonAddNode").
		Text("New Menu Item").
		HxPost(t.renderURL + "&treectl_action=node_add").
		HxTarget("#" + t.id).
		Class("btn btn-sm btn-primary")

	return hb.Div().
		ID(t.id).
		Class("card").
		Child(hb.Div().
			Class("card-header").
			Child(buttonAddNode)).
		Child(hb.Div().
			Class("card-body").
			Child(t.renderTree(tree)).
			Child(hb.NewScript(`document.querySelector('[name="` + t.targetTextareaID + `"]').value = JSON.stringify(` + jsonString + `);`)))
}

func (t *treeControl) renderTree(tree Tree) hb.TagInterface {
	roots := tree.Children("")
	hasRoots := len(roots) > 0

	treeView := hb.Div().Class("tree")

	if !hasRoots {
		return treeView.Text("No menu items. Please use the 'New Menu Item' button above to add a new menu item.")
	}

	for _, root := range roots {
		treeView.Child(t.renderNode(tree, root, 0))
	}

	return treeView
}

func (t *treeControl) renderNode(tree Tree, node Node, level int) hb.TagInterface {

	children := tree.Children(node.ID)
	hasChildren := len(children) > 0
	isRoot := node.ParentID == ""

	iconClass := "bi bi-chevron-right me-2"

	if hasChildren {
		iconClass = "bi bi-chevron-down me-2"
	}

	icon := hb.I().Class(iconClass)

	buttonDelete := hb.Button().
		Type("button").
		Class("btn btn-sm btn-danger float-end").
		Child(hb.I().Class("bi bi-trash")).
		Title("Delete").
		HxPost(t.renderURL + `&treectl_action=node_delete&treectl_id=` + t.id + `&treectl_node_id=` + node.ID).
		HxTarget("#" + t.id).
		HxConfirm("Are you sure?")

	buttonAddNode := hb.Button().
		Type("button").
		Class("btn btn-sm btn-primary float-end me-2").
		Child(hb.I().Class("bi bi-plus-circle")).
		Title("Add Child").
		HxPost(t.renderURL + `&treectl_action=node_add&treectl_id=` + t.id + `&treectl_parent_id=` + node.ID).
		HxTarget("#" + t.id)

	buttonMoveUp := hb.Button().
		Type("button").
		Class("btn btn-sm btn-primary float-end me-2").
		Child(hb.I().Class("bi bi-arrow-up")).
		Title("Move Up").
		HxPost(t.renderURL + `&treectl_action=node_move_up&treectl_id=` + t.id + `&treectl_node_id=` + node.ID).
		HxTarget("#" + t.id)

	buttonMoveDown := hb.Button().
		Type("button").
		Class("btn btn-sm btn-primary float-end me-2").
		Child(hb.I().Class("bi bi-arrow-down")).
		Title("Move Down").
		HxPost(t.renderURL + `&treectl_action=node_move_down&treectl_id=` + t.id + `&treectl_node_id=` + node.ID).
		HxTarget("#" + t.id)

	buttonMoveOutUp := hb.Button().
		Type("button").
		Class("btn btn-sm btn-primary float-end me-2").
		Child(hb.I().Class("bi bi-box-arrow-up-left")).
		Title("Move Out Up").
		HxPost(t.renderURL + `&treectl_action=node_move_out_up&treectl_id=` + t.id + `&treectl_node_id=` + node.ID).
		HxTarget("#" + t.id)

	buttonMoveOutDown := hb.Button().
		Type("button").
		Class("btn btn-sm btn-primary float-end me-2").
		Child(hb.I().Class("bi bi-box-arrow-down-left")).
		Title("Move Out Down").
		HxPost(t.renderURL + `&treectl_action=node_move_out_down&treectl_id=` + t.id + `&treectl_node_id=` + node.ID).
		HxTarget("#" + t.id)

	buttonMoveInUp := hb.Button().
		Type("button").
		Class("btn btn-sm btn-primary float-end me-2").
		Child(hb.I().Class("bi bi-box-arrow-in-up-right")).
		Title("Move In Up").
		HxPost(t.renderURL + `&treectl_action=node_move_in_up&treectl_id=` + t.id + `&treectl_node_id=` + node.ID).
		HxTarget("#" + t.id)

	buttonMoveInDown := hb.Button().
		Type("button").
		Class("btn btn-sm btn-primary float-end me-2").
		Child(hb.I().Class("bi bi-box-arrow-in-down-right")).
		Title("Move In Down").
		HxPost(t.renderURL + `&treectl_action=node_move_in_down&treectl_id=` + t.id + `&treectl_node_id=` + node.ID).
		HxTarget("#" + t.id)

	buttonEdit := hb.Button().
		Type("button").
		Class("btn btn-sm btn-success float-end me-2").
		Child(hb.I().Class("bi bi-pencil-square")).
		Title("Setings").
		HxPost(t.renderURL + `&treectl_action=node_update_modal&treectl_id=` + t.id + `&treectl_node_id=` + node.ID).
		HxTarget("#" + t.id).
		HxSwap(`beforeend`)

	padding := lo.Ternary(isRoot, 0, 20)
	backgroundOpacity := `0.0` + (utils.ToString(1 + level*2))

	nodeView := hb.Div().
		Class("tree-node").
		ClassIf(isRoot, "tree-node-root").
		Style(`border: 1px solid lavender; border-radius: 10px; margin: 5px 0px; padding: 5px`).
		Style("margin-left: " + utils.ToString(padding) + "px;").
		Style("background-color: rgba(0, 149, 182, " + backgroundOpacity + ");").
		Child(icon).
		Child(hb.Span().
			Class("tree-node-name").
			Style(`font-size: 20px;`).
			Text(node.Name)).
		Child(buttonDelete).
		Child(buttonMoveUp).
		Child(buttonMoveDown).
		Child(buttonMoveOutUp).
		Child(buttonMoveOutDown).
		Child(buttonMoveInUp).
		Child(buttonMoveInDown).
		Child(buttonAddNode).
		Child(buttonEdit)

	for _, child := range children {
		nodeView.Child(t.renderNode(tree, child, level+1))
	}

	return nodeView
}

func (t *treeControl) modalNodeUpdate(node Node) hb.TagInterface {
	submitUrl := t.renderURL + `&treectl_action=node_update&treectl_id=` + t.id + `&treectl_node_id=` + node.ID

	name := node.Name
	pageID := node.PageID
	url := node.URL
	target := node.Target

	form := form.NewForm(form.FormOptions{
		ID: "FormMenuUpdate",
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label:    "Menu Item name",
				Name:     "treectl_name",
				Type:     form.FORM_FIELD_TYPE_STRING,
				Value:    name,
				Required: true,
			}),
			form.NewField(form.FieldOptions{
				Label:    "Page",
				Name:     "treectl_page_id",
				Type:     form.FORM_FIELD_TYPE_SELECT,
				Value:    pageID,
				Required: true,
				Help:     "Select a page to link to, if you want to link to a page",
				Options: append([]form.FieldOption{
					{
						Value: "Select site",
						Key:   "",
					},
				},
					lo.Map(t.pageList, func(page cmsstore.PageInterface, index int) form.FieldOption {
						return form.FieldOption{
							Value: page.Name(),
							Key:   page.ID(),
						}
					})...),
			}),
			form.NewField(form.FieldOptions{
				Label:    "Menu Item URL",
				Name:     "treectl_url",
				Type:     form.FORM_FIELD_TYPE_STRING,
				Value:    url,
				Required: true,
				Help:     "The URL to link to (if page is not selected)",
			}),
			form.NewField(form.FieldOptions{
				Label:    "Target",
				Name:     "treectl_target",
				Type:     form.FORM_FIELD_TYPE_SELECT,
				Value:    target,
				Required: true,
				Options: []form.FieldOption{
					{
						Value: "_self",
						Key:   "_self",
					},
					{
						Value: "_blank",
						Key:   "_blank",
					},
					{
						Value: "_parent",
						Key:   "_parent",
					},
				},
			}),
		},
	})

	modalID := "ModalNodeUpdate"
	modalBackdropClass := "ModalBackdrop"

	modalCloseScript := `closeModal` + modalID + `();`

	modalHeading := hb.Heading5().HTML("New Menu").Style(`margin:0px;`)

	modalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalNodeUpdate').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

	buttonSend := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Update Menu Item").
		Class("btn btn-primary float-end").
		HxInclude("#" + modalID).
		HxPost(submitUrl).
		HxTarget(`#` + t.id)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Close").
		Class("btn btn-secondary float-start").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Child(hb.Script(jsCloseFn)).
		Child(bs.ModalDialog().
			Child(bs.ModalContent().
				Child(
					bs.ModalHeader().
						Child(modalHeading).
						Child(modalClose)).
				Child(
					bs.ModalBody().
						Child(form.Build())).
				Child(bs.ModalFooter().
					Style(`display:flex;justify-content:space-between;`).
					Child(buttonCancel).
					Child(buttonSend)),
			))

	backdrop := hb.Div().Class(modalBackdropClass).
		Class("modal-backdrop fade show").
		Style("display:block;z-index:1000;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})
}
