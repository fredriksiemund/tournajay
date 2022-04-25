package tournaments

import (
	"fmt"
	"io"
)

type TeamNode struct {
	left   *TeamNode
	right  *TeamNode
	teamId int
}

func (t *TeamNode) insert() []*TeamNode {
	t.left = &TeamNode{teamId: -1}
	t.right = &TeamNode{teamId: -1}
	return []*TeamNode{t.left, t.right}
}

func print(w io.Writer, node *TeamNode, ns int, ch rune) {
	if node == nil {
		return
	}
	for i := 0; i < ns; i++ {
		fmt.Fprint(w, " ")
	}
	fmt.Fprintf(w, "%c:%v\n", ch, node.teamId)
	print(w, node.left, ns+2, 'L')
	print(w, node.right, ns+2, 'R')
}

func createSingleEliminationTrounament(teamIds []int) *TeamNode {
	bracket := &TeamNode{}
	leafNodes := []*TeamNode{bracket}

	// While length of leafNodes is less than teamIds
	for len(leafNodes) != len(teamIds) {
		// Pop first leafNode
		nextLeafNode := leafNodes[0]
		// Create new leaf nodes
		newLeafNodes := nextLeafNode.insert()
		// Append the leaf nodes to the leafNode array
		leafNodes = leafNodes[1:]
		leafNodes = append(leafNodes, newLeafNodes...)
	}

	// Populate leafNodes
	for i, leafNode := range leafNodes {
		leafNode.teamId = teamIds[i]
	}

	return bracket
}
