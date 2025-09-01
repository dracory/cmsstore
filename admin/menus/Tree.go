package admin

import (
	"encoding/json"
	"sort"

	"github.com/dracory/cmsstore"
	"github.com/dracory/uid"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Node struct {
	ID       string
	Name     string
	ParentID string
	Sequence int
	PageID   string
	URL      string
	Target   string
}

type Tree struct {
	list []Node
}

func NewTreeFromMenuItems(menuItems []cmsstore.MenuItemInterface) *Tree {
	nodes := traverseMenuItems(menuItems, "")
	return &Tree{
		list: nodes,
	}
}

func NewTreeFromJSON(jsonString string) (*Tree, error) {
	var maps []map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &maps)
	if err != nil {
		return nil, err
	}

	// nodesAny, err := utils.FromJSON(json, nil)

	// if err != nil {
	// 	return nil, err
	// }

	// if nodesAny == nil {
	// 	nodesAny = []Node{}
	// }

	// maps := maputils.AnyToArrayMapStringAny(nodesAny)

	nodes := []Node{}

	for _, m := range maps {
		nodes = append(nodes, nodeFromMap(m))
	}

	return &Tree{list: nodes}, err
}

func (tree *Tree) Add(parentID string, node Node) {
	children := tree.Children(parentID)
	node.Sequence = len(children)
	node.ParentID = parentID
	tree.list = append(tree.list, node)

	tree.RecalculateSequences(parentID)
}

// AddBlock adds a new ui.BlockInterface to the Tree
func (tree *Tree) AddMenuItem(parentID string, menuItem cmsstore.MenuItemInterface) {
	children := tree.Children(parentID)

	node := menuItemToNodeWithParentAndSequence(menuItem, parentID, len(children))

	tree.list = append(tree.list, node)

	tree.RecalculateSequences(parentID)
}

// Children returns the children of the Node with the given parentID
func (tree *Tree) Children(parentID string) []Node {
	childrenExt := make([]Node, 0)

	sequence := []int{}
	for _, node := range tree.list {
		if node.ParentID == parentID {
			sequence = append(sequence, node.Sequence)
			childrenExt = append(childrenExt, node)
		}
	}

	sortedSequence := sort.IntSlice(sequence)
	sortedSequence.Sort()

	sortedChildren := make([]Node, 0)

	for _, seq := range sortedSequence {
		for _, node := range childrenExt {
			if node.Sequence == seq {
				sortedChildren = append(sortedChildren, node)
			}
		}
	}

	return sortedChildren
}

// Clone creates a shallow clone of a Node (no children)
//
// This is used to create a clone of a Node, so that the original Node
// is not modified, but we can modify the clone safely
//
// Remember to update the ID, Sequence, and ParentID of the copy with new values
func (tree *Tree) Clone(node Node) Node {
	return Node{
		ID:       node.ID,
		Name:     node.Name,
		ParentID: node.ParentID,
		Sequence: node.Sequence,
	}
}

// Duplicate creates a deep clone of a Node (with children)
// and adds it to the tree, under the same parent
//
// Business Logic:
// - travserses the tree to find all blocks to be duplicated
// - makes a map with current IDs as keys, newly generated IDs as values
// - clones each block, and replaces the ID with the new ID
// - assignes the correct mapped IDs and ParentIDs
// - adds the cloned blocks to the tree directly (using list)
// - moves the duplicated block under the block being duplicated
func (tree *Tree) Duplicate(blockID string) {
	block := tree.Find(blockID)

	if block == nil {
		return
	}

	blocks := tree.Traverse(blockID)

	if len(blocks) == 0 {
		return
	}

	mapIDs := make(map[string]string)

	for _, block := range blocks {
		mapIDs[block.ID] = uid.HumanUid()
	}

	clonedBlocks := make([]Node, 0)
	for _, block := range blocks {
		newID := lo.ValueOr(mapIDs, block.ID, block.ID)
		newParentID := lo.ValueOr(mapIDs, block.ParentID, block.ParentID)
		clonedBlock := tree.Clone(block)
		clonedBlock.ID = newID
		clonedBlock.ParentID = newParentID
		clonedBlocks = append(clonedBlocks, clonedBlock)
	}

	tree.list = append(tree.list, clonedBlocks...)

	newID := mapIDs[blockID]

	tree.MoveToPosition(newID, block.ParentID, block.Sequence+1)
}

func (tree *Tree) Exists(nodeID string) bool {
	for _, node := range tree.list {
		if node.ID == nodeID {
			return true
		}
	}
	return false
}

func (tree *Tree) Find(nodeID string) *Node {
	if nodeID == "" {
		return nil
	}

	for _, node := range tree.list {
		if node.ID == nodeID {
			return &node
		}
	}
	return nil
}

func (tree *Tree) FindNextSibling(nodeID string) *Node {
	block := tree.Find(nodeID)

	if block == nil {
		return nil
	}

	children := tree.Children(block.ParentID)

	_, index, found := lo.FindIndexOf(children, func(bExt Node) bool {
		return bExt.ID == nodeID
	})

	if !found {
		return nil
	}

	if index == len(children)-1 {
		return nil
	}

	return &children[index+1]
}

func (tree *Tree) FindPreviousSibling(nodeID string) *Node {
	block := tree.Find(nodeID)

	if block == nil {
		return nil
	}

	children := tree.Children(block.ParentID)

	_, index, found := lo.FindIndexOf(children, func(bExt Node) bool {
		return bExt.ID == nodeID
	})

	if !found {
		return nil
	}

	if index == 0 {
		return nil
	}

	return &children[index-1]
}

func (tree *Tree) MoveDown(nodeID string) {
	block := tree.Find(nodeID)

	if block == nil {
		return
	}

	next := tree.FindNextSibling(block.ID)

	if next == nil {
		return
	}

	nextSequence := next.Sequence
	sequence := block.Sequence

	block.Sequence = nextSequence
	next.Sequence = sequence

	tree.Update(*block)
	tree.Update(*next)

	tree.RecalculateSequences(block.ParentID)
}

func (tree *Tree) MoveToParent(nodeID string, parentID string) {
	block := tree.Find(nodeID)

	if block == nil {
		return
	}

	if block.ParentID == parentID {
		return
	}

	children := tree.Children(block.ParentID)

	block.ParentID = parentID
	block.Sequence = len(children)

	tree.Update(*block)

	//tree.Remove(nodeID)
	//tree.Add(parentID, *block)

	tree.RecalculateSequences(parentID)
}

func (tree *Tree) MoveToPosition(nodeID string, parentID string, position int) {
	tree.MoveToParent(nodeID, parentID)

	block := tree.Find(nodeID)

	if block == nil {
		return
	}

	if block.Sequence == position {
		return
	}

	if position < 0 {
		return // position already at the top
	}

	if position > len(tree.Children(parentID)) {
		return // position already at the bottom
	}

	if block.Sequence < position {
		// move down
		for i := block.Sequence; i < position; i++ {
			tree.MoveDown(nodeID)
		}
	} else {
		// move up
		for i := block.Sequence; i > position; i-- {
			tree.MoveUp(nodeID)
		}
	}
}

func (tree *Tree) MoveUp(nodeID string) {
	block := tree.Find(nodeID)

	if block == nil {
		return
	}

	previous := tree.FindPreviousSibling(block.ID)

	if previous == nil {
		return
	}

	previousSequence := previous.Sequence
	sequence := block.Sequence

	block.Sequence = previousSequence
	previous.Sequence = sequence

	tree.Update(*block)
	tree.Update(*previous)

	tree.RecalculateSequences(block.ParentID)
}

func (tree *Tree) Parent(nodeID string) *Node {
	block := tree.Find(nodeID)

	if block == nil {
		return nil
	}

	return tree.Find(block.ParentID)
}

func (tree *Tree) RecalculateSequences(blockID string) {
	children := tree.Children(blockID)

	for i, child := range children {
		child.Sequence = i
		tree.Update(child)
	}
}

func (tree *Tree) List() []Node {
	return tree.list
}

// Remove removes the block with the given id
//
// Buisiness Logic:
// - checks if the block exists, if not, do nothing
// - removes the block from the list
// - recalculates the sequences of the parent's children
func (tree *Tree) Remove(nodeID string) {
	node := tree.Find(nodeID)

	if node == nil {
		return
	}

	parentID := node.ParentID
	for i, ext := range tree.list {
		if ext.ID == nodeID {
			tree.list = append(tree.list[:i], tree.list[i+1:]...)
		}
	}

	tree.RemoveOrphans()

	tree.RecalculateSequences(parentID)
}

// RemoveOrphans removes all orphaned blocks that have no parent
//
// Buisiness Logic:
// - finds and creates a new list without orphaned blocks
// - non orphaned blocks are the ones that have a parent or root blocks
// - updates the list with the new list
//
// Parameters:
// - none
//
// Returns:
// - none
func (tree *Tree) RemoveOrphans() {
	nonOrphans := make([]Node, 0)

	for _, block := range tree.list {
		if block.ParentID == "" {
			nonOrphans = append(nonOrphans, block)
			continue
		}

		parent := tree.Parent(block.ID)

		if parent != nil {
			nonOrphans = append(nonOrphans, block)
		}
	}

	tree.list = nonOrphans
}

func (tree *Tree) Traverse(blockID string) []Node {
	block := tree.Find(blockID)

	if block == nil {
		return []Node{}
	}

	travsersed := make([]Node, 0)
	travsersed = append(travsersed, *block)

	for _, child := range tree.Children(blockID) {
		travsersed = append(travsersed, tree.Traverse(child.ID)...)
	}

	return travsersed
}

func (tree *Tree) Update(node Node) {
	for i, ext := range tree.list {
		if ext.ID == node.ID {
			tree.list[i] = node
		}
	}
}

func (tree *Tree) ToJSON() (jsonString string, err error) {
	nodes := tree.list

	maps := make([]map[string]interface{}, 0)

	for _, node := range nodes {
		maps = append(maps, nodeToMap(node))
	}

	jsonBytes, err := json.Marshal(maps)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func nodeToMap(node Node) map[string]any {
	return map[string]any{
		"id":        node.ID,
		"name":      node.Name,
		"parent_id": node.ParentID,
		"sequence":  node.Sequence,
		"page_id":   node.PageID,
		"url":       node.URL,
		"target":    node.Target,
	}
}

func nodeFromMap(nodeMap map[string]interface{}) Node {
	return Node{
		ID:       nodeMap["id"].(string),
		Name:     nodeMap["name"].(string),
		ParentID: nodeMap["parent_id"].(string),
		Sequence: cast.ToInt(nodeMap["sequence"]),
		PageID:   nodeMap["page_id"].(string),
		URL:      nodeMap["url"].(string),
		Target:   nodeMap["target"].(string),
	}
}

// func (tree *Tree) ToJSON() (json string, err error) {
// 	parentBlocks := tree.Children("")

// 	blocks := make([]map[string]interface{}, 0)

// 	for _, node := range parentBlocks {
// 		blocks = append(blocks, tree.nodeToMap(node))
// 	}

// 	return utils.ToJSON(blocks)
// }

// func (tree *Tree) nodeToMap(node Node) map[string]any {
// 	children := tree.Children(node.ID)

// 	childrenMaps := []map[string]any{}
// 	for _, child := range children {
// 		childrenMaps = append(childrenMaps, tree.nodeToMap(child))
// 	}

// 	block := ui.NewFromMap(map[string]interface{}{
// 		"id":         node.ID,
// 		"type":       node.Type,
// 		"parameters": node.Parameters,
// 		"children":   children,
// 	})

// 	return block
// }

func traverseMenuItems(menuItems []cmsstore.MenuItemInterface, parentID string) []Node {
	children := childrenMenuItems(menuItems, parentID)

	list := []Node{}

	for _, child := range children {
		node := menuItemToNodeWithParentAndSequence(child, parentID, child.SequenceInt())
		list = append(list, node)
		list = append(list, traverseMenuItems(menuItems, child.ID())...)
	}

	// for index, menuItem := range menuItems {
	// 	children := childrenMenuItems(menuItems, menuItem.ID())
	// 	node := menuItemToNodeWithParentAndSequence(menuItem, parentID, index)

	// 	list = append(list, node)
	// 	list = append(list, traverseMenuItems(children, menuItem.ID())...)
	// }

	return list
}

func childrenMenuItems(menuItems []cmsstore.MenuItemInterface, parentID string) []cmsstore.MenuItemInterface {
	unorderedChildren := lo.Filter(menuItems, func(item cmsstore.MenuItemInterface, _ int) bool {
		return item.ParentID() == parentID
	})

	unorderedSequences := lo.Map(unorderedChildren, func(item cmsstore.MenuItemInterface, _ int) int {
		return item.SequenceInt()
	})

	uniqUnorderedSequences := lo.Uniq(unorderedSequences) // fix any sequence issues

	sortedSequences := sort.IntSlice(uniqUnorderedSequences)
	sortedSequences.Sort()

	sortedChildren := make([]cmsstore.MenuItemInterface, 0)

	for _, sequence := range sortedSequences {
		for _, child := range unorderedChildren {
			if child.SequenceInt() == sequence {
				child.SetSequenceInt(len(sortedChildren)) // fix any sequence issues
				sortedChildren = append(sortedChildren, child)
			}
		}
	}

	return sortedChildren
}

func menuItemToNodeWithParentAndSequence(menuItem cmsstore.MenuItemInterface, parentID string, sequence int) Node {
	return Node{
		ID:       menuItem.ID(),
		Name:     menuItem.Name(),
		ParentID: parentID,
		Sequence: sequence,
		PageID:   menuItem.PageID(),
		URL:      menuItem.URL(),
		Target:   menuItem.Target(),
	}
}

func menuItemToNode(menuItem cmsstore.MenuItemInterface) Node {
	return Node{
		ID:       menuItem.ID(),
		Name:     menuItem.Name(),
		ParentID: menuItem.ParentID(),
		Sequence: menuItem.SequenceInt(),
		PageID:   menuItem.PageID(),
		URL:      menuItem.URL(),
		Target:   menuItem.Target(),
	}
}
