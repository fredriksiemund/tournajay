package postgres

import (
	"fmt"
	"strings"

	"github.com/fredriksiemund/tournament-planner/pkg/tournaments"
	"github.com/jackc/pgx/v4"
)

type GameModel struct {
	Db *pgx.Conn
}

type games struct {
	tournamentId int
}

func (g games) iterate(node *tournaments.TeamNode, depth int, db *pgx.Conn) int {
	var previousGameIds []int
	if node.Left == nil && node.Right == nil {
		return -1
	} else if node.Left != nil {
		prevGame := g.iterate(node.Left, depth+1, db)
		if prevGame != -1 {
			previousGameIds = append(previousGameIds, prevGame)
		}
	} else if node.Right != nil {
		prevGame := g.iterate(node.Right, depth+1, db)
		if prevGame != -1 {
			previousGameIds = append(previousGameIds, prevGame)
		}
	}

	// Create a game with the team on the left and right side
	stmt := "INSERT INTO games (tournament_id, depth) VALUES ($1, $2) RETURNING id"
	row := db.QueryRow(ctx, stmt, g.tournamentId, depth)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1
	}

	// Insert contestants
	var contestants []int
	if node.Left.TeamId != -1 {
		contestants = append(contestants, node.Left.TeamId)
	} else if node.Right.TeamId != -1 {
		contestants = append(contestants, node.Right.TeamId)
	}
	if len(contestants) > 0 {
		var placeholders []string
		var values []interface{}
		for i := 0; i < len(contestants); i++ {
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
			values = append(values, id, contestants[i])
		}
		stmt = fmt.Sprintf("INSERT INTO contestants (game_id, team_id) VALUES %s", strings.Join(placeholders, ", "))
		_, err := db.Exec(ctx, stmt, values...)
		if err != nil {
			return -1
		}
	}

	// Insert path

	return id
}

func (m *GameModel) InsertSingleEliminationGames(finalNode *tournaments.TeamNode) ([]int, error) {
	// var placeholders []string
	// var values []interface{}

	// for i := 0; i < nbrOfTeams; i++ {
	// 	placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
	// 	values = append(values, tournamentId, fmt.Sprintf("Team %d", i+1))
	// }

	// stmt := fmt.Sprintf("INSERT INTO teams (tournament_id, name) VALUES %s RETURNING id", strings.Join(placeholders, ", "))

	// rows, err := m.Db.Query(ctx, stmt, values...)
	// if err != nil {
	// 	return nil, err
	// }

	// defer rows.Close()

	// var ids []int
	// for rows.Next() {
	// 	var id int
	// 	err = rows.Scan(&id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	ids = append(ids, id)
	// }

	// if err = rows.Err(); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
