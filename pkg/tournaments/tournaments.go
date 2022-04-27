package tournaments

import "github.com/fredriksiemund/tournament-planner/pkg/models"

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

func generateTeamTree(teamIds []int) *TeamNode {
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

type GameNode struct {
	Id          int
	Contestants []int
	Left        *GameNode
	Right       *GameNode
}

type gameTree struct {
	plan map[int][]*GameNode
}

func (g gameTree) generateGameTree(node *TeamNode, depth int) (*GameNode, error) {
	if node.Left == nil && node.Right == nil {
		return nil, nil
	} else if node.Left == nil || node.Right == nil {
		return nil, models.ErrInvalidTree
	}

	leftGame, err := g.generateGameTree(node.Left, depth+1)
	if err != nil {
		return nil, err
	}

	rightGame, err := g.generateGameTree(node.Right, depth+1)
	if err != nil {
		return nil, err
	}

	newGame := &GameNode{Left: leftGame, Right: rightGame}
	if node.Left.TeamId != -1 {
		newGame.Contestants = append(newGame.Contestants, node.Left.TeamId)
	}
	if node.Right.TeamId != -1 {
		newGame.Contestants = append(newGame.Contestants, node.Right.TeamId)
	}

	g.plan[depth] = append(g.plan[depth], newGame)

	return newGame, nil
}

func NewSingleElimination(teamIds []int) (map[int][]*GameNode, error) {
	teamTree := generateTeamTree(teamIds)

	g := &gameTree{plan: make(map[int][]*GameNode)}
	_, err := g.generateGameTree(teamTree, 0)
	if err != nil {
		return nil, err
	}

	return g.plan, nil
}
