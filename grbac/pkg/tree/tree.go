package tree

// Tree is a Radix search tree that supports wildcards
type Tree struct {
	root *Node
}

// NewTree is used to initialize a wildcard tree
func NewTree() *Tree {
	root := NewNode("ROOT", nil)
	return &Tree{
		root: root,
	}
}

// Query is used to query the current tree by args
func (tree *Tree) Query(args []string) ([]Data, error) {
	var data []Data

	parents := []*Node{tree.root}
	for i, arg := range args {
		eof := i == len(args)-1

		var nodes []*Node
		for _, parent := range parents {
			children, childData, err := parent.Find(arg)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, children...)
			if eof {
				data = append(data, childData...)
			}
		}
		parents = nodes
	}
	return data, nil
}

// Insert is used to insert a node into the current tree
func (tree *Tree) Insert(args []string, data Data) {
	parent := tree.root

	var nodeData Data
	for i, arg := range args {
		eof := i == len(args)-1
		if eof {
			nodeData = data
		}
		child := NewNode(arg, nodeData)
		parent.Insert(child)
		parent = child
	}
}
