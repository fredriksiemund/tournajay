package postgres

import (
	"context"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/jackc/pgx/v4"
)

type TournamentTypeModel struct {
	Db *pgx.Conn
}

func (m *TournamentTypeModel) All() ([]*models.TournamentType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := `SELECT id, name FROM tournament_types`

	rows, err := m.Db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	types := []*models.TournamentType{}

	for rows.Next() {
		t := &models.TournamentType{}

		err = rows.Scan(&t.Id, &t.Name)
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
