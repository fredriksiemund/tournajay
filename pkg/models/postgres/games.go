package postgres

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fredriksiemund/tournament-planner/pkg/tournaments"
	"github.com/jackc/pgx/v4"
)

type GameModel struct {
	Db *pgx.Conn
}

func (m *GameModel) InsertSingleEliminationGames(tournamentId int, games map[int][]*tournaments.GameNode) error {
	// Sort keys
	depths := make([]int, 0, len(games))
	for k := range games {
		depths = append(depths, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(depths)))

	for _, depth := range depths {
		for _, game := range games[depth] {
			// Insert game
			stmt := "INSERT INTO games (tournament_id, depth) VALUES ($1, $2) RETURNING id"
			row := m.Db.QueryRow(ctx, stmt, tournamentId, depth)

			var id int
			err := row.Scan(&id)
			if err != nil {
				return err
			}
			game.Id = id

			// Insert contestants
			if len(game.Contestants) > 0 {
				var placeholders []string
				var values []interface{}
				for i := 0; i < len(game.Contestants); i++ {
					placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
					values = append(values, id, game.Contestants[i])
				}
				stmt = fmt.Sprintf("INSERT INTO contestants (game_id, team_id) VALUES %s", strings.Join(placeholders, ", "))
				_, err := m.Db.Exec(ctx, stmt, values...)
				if err != nil {
					return err
				}
			}

			// Insert path
			if game.Left != nil {
				stmt = "INSERT INTO game_paths (from_game_id, to_game_id, result_type_id) VALUES ($1, $2, $3)"
				_, err := m.Db.Exec(ctx, stmt, game.Left.Id, game.Id, 1)
				if err != nil {
					return err
				}
			}
			if game.Right != nil {
				stmt = "INSERT INTO game_paths (from_game_id, to_game_id, result_type_id) VALUES ($1, $2, $3)"
				_, err := m.Db.Exec(ctx, stmt, game.Right.Id, game.Id, 1)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *GameModel) Exists(tournamentId int) (bool, error) {
	stmt := "SELECT EXISTS (SELECT 1 from games WHERE tournament_id = $1)"

	row := m.Db.QueryRow(ctx, stmt, tournamentId)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
