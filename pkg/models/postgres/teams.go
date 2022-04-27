package postgres

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
)

type TeamModel struct {
	Db *pgx.Conn
}

func (m *TeamModel) Insert(nbrOfTeams int, tournamentId int) ([]int, error) {
	var placeholders []string
	var values []interface{}

	for i := 0; i < nbrOfTeams; i++ {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
		values = append(values, tournamentId, fmt.Sprintf("Team %d", i+1))
	}

	stmt := fmt.Sprintf("INSERT INTO teams (tournament_id, name) VALUES %s RETURNING id", strings.Join(placeholders, ", "))

	rows, err := m.Db.Query(ctx, stmt, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
