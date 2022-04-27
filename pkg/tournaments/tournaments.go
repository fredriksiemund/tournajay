package tournaments

import (
	"fmt"
	"io"
)

type TeamNode struct {
	Left   *TeamNode
	Right  *TeamNode
	TeamId int
}

func (t *TeamNode) insert() []*TeamNode {
	t.Left = &TeamNode{TeamId: -1}
	t.Right = &TeamNode{TeamId: -1}
	return []*TeamNode{t.Left, t.Right}
}

func Print(w io.Writer, node *TeamNode, ns int, ch rune) {
	if node == nil {
		return
	}
	for i := 0; i < ns; i++ {
		fmt.Fprint(w, " ")
	}
	fmt.Fprintf(w, "%c:%v\n", ch, node.TeamId)
	Print(w, node.Left, ns+2, 'L')
	Print(w, node.Right, ns+2, 'R')
}

func NewSingleElimination(teamIds []int) *TeamNode {
	winnersBracket := &TeamNode{TeamId: -1}
	leafNodes := []*TeamNode{winnersBracket}

	// While length of leafNodes is less than teamIds
	for len(leafNodes) < len(teamIds) {
		// Pop first leafNode
		nextLeafNode := leafNodes[0]
		leafNodes = leafNodes[1:]
		// Create new leaf nodes
		newLeafNodes := nextLeafNode.insert()
		// Append the leaf nodes to the leafNode array
		leafNodes = append(leafNodes, newLeafNodes...)
	}

	// Populate leafNodes
	for i, leafNode := range leafNodes {
		leafNode.TeamId = teamIds[i]
	}

	return winnersBracket
}
