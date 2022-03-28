package postgres

import (
	"context"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/jackc/pgx/v4"
)

type TournamentTypeModel struct {
	Db *pgx.Conn
}

func (m *TournamentTypeModel) All() ([]*models.TournamentType, error) {
	stmt := `SELECT * FROM tournament_types`

	rows, err := m.Db.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	types := []*models.TournamentType{}

	for rows.Next() {
		t := &models.TournamentType{}

		err = rows.Scan(&t.Id, &t.Title)
		if err != nil {
			return nil, err
		}

		types = append(types, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return types, nil
}
