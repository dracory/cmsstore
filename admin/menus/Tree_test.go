package admin

import (
	"testing"

	"github.com/gouniverse/cmsstore"
)

func Test_Tree_FromJSON(t *testing.T) {
	tree, err := NewTreeFromJSON(`[
		{"id":"1","name":"","page_id":"","parent_id":"","sequence":0,"target":"","url":""},
		{"id":"2","name":"","page_id":"","parent_id":"","sequence":1,"target":"","url":""},
		{"id":"3","name":"","page_id":"","parent_id":"2","sequence":0,"target":"","url":""},
		{"id":"4","name":"","page_id":"","parent_id":"2","sequence":1,"target":"","url":""},
		{"id":"5","name":"","page_id":"","parent_id":"","sequence":2,"target":"","url":""}
	]`)

	if err != nil {
		t.Fatal(err)
	}

	if len(tree.list) != 5 {
		t.Fatal("tree list length is not 5")
	}
}

func Test_Tree_FromMenuItems(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("").SetSequenceInt(2),
	})

	if len(tree.list) != 5 {
		t.Fatal("tree list length is not 5")
	}
}

func Test_Tree_Children(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("").SetSequenceInt(2),
	})

	chldren := tree.Children("2")

	if len(chldren) != 2 {
		t.Fatal("Expected 2 children, got:", len(chldren))
	}

	if chldren[0].ID != "3" {
		t.Fatal("Expected child ID '3', got:", chldren[0].ID)
	}

	if chldren[1].ID != "4" {
		t.Fatal("Expected child ID '4', got:", chldren[1].ID)
	}
}

func Test_Tree_Duplicate(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("").SetSequenceInt(2),
	})

	if len(tree.list) != 5 {
		t.Fatal("Expected 5 nodes, got:", len(tree.list))
	}

	tree.Duplicate("2")

	if len(tree.list) != 8 {
		t.Fatal("Expected 8 nodes, got:", len(tree.list))
	}

	duplicatedNode := tree.list[5]

	if duplicatedNode.ParentID != "" {
		t.Fatal("Expected parent ID '', got:", tree.list[5].ParentID)
	}

	if duplicatedNode.Sequence != 4 {
		t.Fatal("Expected sequence 4, got:", tree.list[6].Sequence)
	}

	if len(tree.Children(duplicatedNode.ID)) != 2 {
		t.Fatal("Expected 2 children, got:", len(tree.Children(duplicatedNode.ID)))
	}
}

func Test_Tree_Find(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("").SetSequenceInt(2),
	})

	node := tree.Find("4")

	if node == nil {
		t.Fatal("Expected node, got nil")
	}

	if node.ID != "4" {
		t.Fatal("Expected node ID '4', got:", node.ID)
	}
}

func Test_Tree_FindNextSibling(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("2").SetSequenceInt(2),
		cmsstore.NewMenuItem().SetID("6").SetParentID("").SetSequenceInt(3),
	})

	node := tree.FindNextSibling("4")

	if node == nil {
		t.Fatal("Expected node, got nil")
	}

	if node.ID != "5" {
		t.Fatal("Expected node ID '5', got:", node.ID)
	}
}

func Test_Tree_FindPreviousSibling(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("2").SetSequenceInt(2),
		cmsstore.NewMenuItem().SetID("6").SetParentID("").SetSequenceInt(3),
	})

	node := tree.FindPreviousSibling("4")

	if node == nil {
		t.Fatal("Expected node, got nil")
	}

	if node.ID != "3" {
		t.Fatal("Expected node ID '3', got:", node.ID)
	}
}

// func Test_Tree_FlatNodeToNode(t *testing.T) {
// 	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
// 		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
// 		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
// 		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
// 		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
// 		cmsstore.NewMenuItem().SetID("5").SetParentID("").SetSequenceInt(2),
// 	})

// 	node := tree.flatNodeToNode(tree.list[1])

// 	if node == nil {
// 		t.Fatal("Expected node, got nil")
// 	}

// 	if node.ID() != "2" {
// 		t.Fatal("Expected node ID '2', got:", node.ID())
// 	}

// 	if len(node.Children()) != 2 {
// 		t.Fatal("Expected 2 child, got:", len(node.Children()))
// 	}
// }

func Test_Tree_MoveDown(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("2").SetSequenceInt(2),
		cmsstore.NewMenuItem().SetID("6").SetParentID("").SetSequenceInt(3),
	})

	tree.MoveDown("4")

	if len(tree.Children("2")) != 3 {
		t.Fatal("Expected 3 child, got:", len(tree.Children("2")))
	}

	if tree.Children("2")[2].ID != "4" {
		t.Fatal("Expected child ID '4', got:", tree.Children("2")[2].ID)
	}

	if tree.Children("2")[1].ID != "5" {
		t.Fatal("Expected child ID '5', got:", tree.Children("2")[1].ID)
	}

	if tree.Children("2")[0].ID != "3" {
		t.Fatal("Expected child ID '3', got:", tree.Children("2")[0].ID)
	}
}

func Test_Tree_MoveToParent(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("2").SetSequenceInt(2),
		cmsstore.NewMenuItem().SetID("6").SetParentID("").SetSequenceInt(3),
	})

	tree.MoveToParent("6", "2")

	if len(tree.Children("2")) != 4 {
		t.Fatal("Expected 4 child, got:", len(tree.Children("2")))
	}

	if tree.Children("2")[0].ID != "3" {
		t.Fatal("Expected child ID '3', got:", tree.Children("2")[0].ID)
	}

	if tree.Children("2")[1].ID != "4" {
		t.Fatal("Expected child ID '4', got:", tree.Children("2")[1].ID)
	}

	if tree.Children("2")[2].ID != "5" {
		t.Fatal("Expected child ID '5', got:", tree.Children("2")[1].ID)
	}

	if tree.Children("2")[3].ID != "6" {
		t.Fatal("Expected child ID '6', got:", tree.Children("2")[3].ID)
	}
}

func Test_Tree_MoveToPosition(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("2").SetSequenceInt(2),
		cmsstore.NewMenuItem().SetID("6").SetParentID("").SetSequenceInt(3),
	})

	tree.MoveToPosition("6", "2", 2)

	if len(tree.Children("2")) != 4 {
		t.Fatal("Expected 4 child, got:", len(tree.Children("2")))
	}

	if tree.Children("2")[0].ID != "3" {
		t.Fatal("Expected child ID '3', got:", tree.Children("2")[0].ID)
	}

	if tree.Children("2")[1].ID != "4" {
		t.Fatal("Expected child ID '4', got:", tree.Children("2")[1].ID)
	}

	if tree.Children("2")[2].ID != "6" {
		t.Fatal("Expected child ID '6', got:", tree.Children("2")[2].ID)
	}

	if tree.Children("2")[3].ID != "5" {
		t.Fatal("Expected child ID '5', got:", tree.Children("2")[3].ID)
	}

	if len(tree.Children("")) != 2 {
		t.Fatal("Expected 2 child, got:", len(tree.Children("")))
	}
}

func Test_Tree_MoveUp(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("2").SetSequenceInt(2),
		cmsstore.NewMenuItem().SetID("6").SetParentID("").SetSequenceInt(3),
	})

	tree.MoveUp("4")

	if len(tree.Children("2")) != 3 {
		t.Fatal("Expected 3 child, got:", len(tree.Children("2")))
	}

	if tree.Children("2")[0].ID != "4" {
		t.Fatal("Expected child ID '4', got:", tree.Children("2")[0].ID)
	}

	if tree.Children("2")[1].ID != "3" {
		t.Fatal("Expected child ID '3', got:", tree.Children("2")[1].ID)
	}

	if tree.Children("2")[2].ID != "5" {
		t.Fatal("Expected child ID '5', got:", tree.Children("2")[1].ID)
	}
}

func Test_Tree_Parent(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("").SetSequenceInt(2),
	})

	node := tree.Parent("4")

	if node == nil {
		t.Fatal("Expected node, got nil")
	}

	if node.ID != "2" {
		t.Fatal("Expected node ID '2', got:", node.ID)
	}

	if len(tree.Children("2")) != 2 {
		t.Fatal("Expected 2 child, got:", len(tree.Children("2")))
	}
}

func Test_Tree_Remove(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("2").SetSequenceInt(2),
		cmsstore.NewMenuItem().SetID("6").SetParentID("").SetSequenceInt(3),
	})

	tree.Remove("4")

	if len(tree.Children("2")) != 2 {
		t.Fatal("Expected 2 child, got:", len(tree.Children("2")))
	}

	if tree.Children("2")[0].ID != "3" {
		t.Fatal("Expected child ID '3', got:", tree.Children("2")[0].ID)
	}

	if tree.Children("2")[0].Sequence != 0 {
		t.Fatal("Expected child sequence 0, got:", tree.Children("2")[0].Sequence)
	}

	if tree.Children("2")[1].ID != "5" {
		t.Fatal("Expected child ID '5', got:", tree.Children("2")[1].ID)
	}

	if tree.Children("2")[1].Sequence != 1 {
		t.Fatal("Expected child sequence 1, got:", tree.Children("2")[1].Sequence)
	}
}

func Test_Tree_RemoveOrphans(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("2").SetSequenceInt(2),
		cmsstore.NewMenuItem().SetID("6").SetParentID("").SetSequenceInt(3),
	})

	// append orphan nodes
	tree.list = append(tree.list, Node{ID: "77", ParentID: "43", Sequence: 1})
	tree.list = append(tree.list, Node{ID: "84", ParentID: "43", Sequence: 0})
	// prepend orphan nodes
	tree.list = append([]Node{{ID: "73", ParentID: "43", Sequence: 1}}, tree.list...)
	tree.list = append([]Node{{ID: "86", ParentID: "43", Sequence: 0}}, tree.list...)
	// insert orphan nodes at the middle
	tree.list = append(tree.list[:2], append([]Node{{ID: "99", ParentID: "43", Sequence: 1}}, tree.list[2:]...)...)

	orphanIDs := []string{"77", "84", "73", "86", "99"}

	for _, orphanID := range orphanIDs {
		if !tree.Exists(orphanID) {
			t.Fatal("Expected orphan node, got nil, ID:", orphanID)
		}
	}

	if len(tree.list) != 11 {
		t.Fatal("Expected 11 nodes, got:", len(tree.list))
	}

	tree.RemoveOrphans()

	for _, orphanID := range orphanIDs {
		if tree.Exists(orphanID) {
			t.Fatal("Expected orphan node to be removed, ID:", orphanID)
		}
	}

	if len(tree.list) != 6 {
		t.Fatal("Expected 6 nodes, got:", len(tree.list))
	}
}

func Test_Tree_RecalculateSequences(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("4").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("5").SetParentID("").SetSequenceInt(2),
	})

	tree.RecalculateSequences("2")

	if len(tree.Children("2")) != 2 {
		t.Fatal("Expected 2 children, got:", len(tree.Children("2")))
	}

	if tree.Children("2")[0].Sequence != 0 {
		t.Fatal("Expected sequence 0, got:", tree.Children("2")[0].Sequence)
	}

	if tree.Children("2")[1].Sequence != 1 {
		t.Fatal("Expected sequence 1, got:", tree.Children("2")[1].Sequence)
	}

	if tree.Children("2")[0].ID != "3" {
		t.Fatal("Expected child ID '3', got:", tree.Children("2")[0].ID)
	}

	if tree.Children("2")[1].ID != "4" {
		t.Fatal("Expected child ID '4', got:", tree.Children("2")[1].ID)
	}
}

func Test_Tree_Traverse(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetParentID("").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2").SetParentID("").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("2.1").SetParentID("2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2.1.1").SetParentID("2.1").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2.1.2").SetParentID("2.1").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("2.2").SetParentID("2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("2.2.1").SetParentID("2.2").SetSequenceInt(0),
		cmsstore.NewMenuItem().SetID("2.2.2").SetParentID("2.2").SetSequenceInt(1),
		cmsstore.NewMenuItem().SetID("3").SetParentID("").SetSequenceInt(2),
	})

	nodes := tree.Traverse("2")

	if len(nodes) != 7 {
		t.Fatal("Expected 7 nodes, got:", len(nodes))
	}

	if nodes[0].ID != "2" {
		t.Fatal("Expected node ID '2.1', got:", nodes[0].ID)
	}

	if nodes[1].ID != "2.1" {
		t.Fatal("Expected node ID '2.1', got:", nodes[1].ID)
	}

	if nodes[2].ID != "2.1.1" {
		t.Fatal("Expected node ID '2.1.1', got:", nodes[2].ID)
	}

	if nodes[3].ID != "2.1.2" {
		t.Fatal("Expected node ID '2.1.2', got:", nodes[3].ID)
	}

	if nodes[4].ID != "2.2" {
		t.Fatal("Expected node ID '2.2', got:", nodes[4].ID)
	}

	if nodes[5].ID != "2.2.1" {
		t.Fatal("Expected node ID '2.2.1', got:", nodes[5].ID)
	}

	if nodes[6].ID != "2.2.2" {
		t.Fatal("Expected node ID '2.2.2', got:", nodes[6].ID)
	}
}

func Test_Tree_Update(t *testing.T) {
	tree := NewTreeFromMenuItems([]cmsstore.MenuItemInterface{
		cmsstore.NewMenuItem().SetID("1").SetName("test").SetParentID("").SetSequenceInt(0),
	})

	node := tree.Find(`1`)

	if node == nil {
		t.Fatal("Expected node, got nil")
	}

	node.Name = "updated"

	tree.Update(*node)

	foundNode := tree.Find(`1`)

	if foundNode == nil {
		t.Fatal("Expected node, got nil")
	}

	if foundNode.Name != "updated" {
		t.Fatal("Expected type 'updated', got:", foundNode.Name)
	}
}
