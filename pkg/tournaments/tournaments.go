package tournaments

import (
	"fmt"
	"sort"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
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

type GameTemplate struct {
	Id          int
	Contestants []string
}

type RoundTemplate struct {
	Title string
	Games []GameTemplate
}

func roundName(depth int, maxDepth int) string {
	if depth == 0 {
		return "Final"
	} else if depth == 1 {
		return "Semi finals"
	} else {
		return fmt.Sprintf("Round %d", maxDepth-depth+1)
	}
}

func SingleEliminationTemplate(games map[int][]*models.Game, teams map[int]*models.Team) []*RoundTemplate {
	depths := make([]int, 0, len(games))
	for k := range games {
		depths = append(depths, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(depths)))

	var rounds []*RoundTemplate
	for _, depth := range depths {
		round := &RoundTemplate{Title: roundName(depth, depths[0])}

		for _, g := range games[depth] {
			game := &GameTemplate{Id: g.Id}

			for i, id := range g.TeamIds {
				if id != -1 {
					game.Contestants = append(game.Contestants, teams[id].Name)
				} else {
					game.Contestants = append(game.Contestants, fmt.Sprintf("Winner of #%d", g.PreviousGameIds[i]))
				}
			}

			round.Games = append(round.Games, *game)
		}

		rounds = append(rounds, round)
	}

	return rounds
}
