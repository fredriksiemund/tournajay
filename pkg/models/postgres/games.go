package postgres

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/fredriksiemund/tournament-planner/pkg/tournaments"
	"github.com/jackc/pgx/v4"
)

type GameModel struct {
	Db *pgx.Conn
}

type game struct {
	id          int
	contestants []int
	left        *game
	right       *game
}

type games struct {
	tournamentId int
	plan         map[int][]*game
}

func (g games) iterate(node *tournaments.TeamNode, depth int) (*game, error) {
	if node.Left == nil && node.Right == nil {
		return nil, nil
	} else if node.Left == nil || node.Right == nil {
		return nil, models.ErrInvalidTree
	}

	leftGame, err := g.iterate(node.Left, depth+1)
	if err != nil {
		return nil, err
	}

	rightGame, err := g.iterate(node.Right, depth+1)
	if err != nil {
		return nil, err
	}

	newGame := &game{left: leftGame, right: rightGame}
	if node.Left.TeamId != -1 {
		newGame.contestants = append(newGame.contestants, node.Left.TeamId)
	}
	if node.Right.TeamId != -1 {
		newGame.contestants = append(newGame.contestants, node.Right.TeamId)
	}

	g.plan[depth] = append(g.plan[depth], newGame)

	return newGame, nil
}

func (m *GameModel) InsertSingleEliminationGames(tournamentId int, finalNode *tournaments.TeamNode) error {
	g := &games{tournamentId: tournamentId, plan: make(map[int][]*game)}

	_, err := g.iterate(finalNode, 0)
	if err != nil {
		return err
	}

	// Sort keys
	depths := make([]int, 0, len(g.plan))
	for k := range g.plan {
		depths = append(depths, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(depths)))

	for _, depth := range depths {
		for _, game := range g.plan[depth] {
			// Insert game
			stmt := "INSERT INTO games (tournament_id, depth) VALUES ($1, $2) RETURNING id"
			row := m.Db.QueryRow(ctx, stmt, tournamentId, depth)

			var id int
			err := row.Scan(&id)
			if err != nil {
				return err
			}
			game.id = id

			// Insert contestants
			if len(game.contestants) > 0 {
				var placeholders []string
				var values []interface{}
				for i := 0; i < len(game.contestants); i++ {
					placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
					values = append(values, id, game.contestants[i])
				}
				stmt = fmt.Sprintf("INSERT INTO contestants (game_id, team_id) VALUES %s", strings.Join(placeholders, ", "))
				_, err := m.Db.Exec(ctx, stmt, values...)
				if err != nil {
					return err
				}
			}

			// Insert path
			if game.left != nil {
				stmt = "INSERT INTO game_paths (from_game_id, to_game_id, result_type_id) VALUES ($1, $2, $3)"
				_, err := m.Db.Exec(ctx, stmt, game.left.id, game.id, 1)
				if err != nil {
					return err
				}
			}
			if game.right != nil {
				stmt = "INSERT INTO game_paths (from_game_id, to_game_id, result_type_id) VALUES ($1, $2, $3)"
				_, err := m.Db.Exec(ctx, stmt, game.right.id, game.id, 1)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, v := range depths {
		fmt.Printf("%d ->", v)
		for _, va := range g.plan[v] {
			fmt.Printf(" %v", *va)
		}
		fmt.Println()
	}

	return nil
}
