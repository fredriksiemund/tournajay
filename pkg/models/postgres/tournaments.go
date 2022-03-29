package postgres

import (
	"context"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/jackc/pgx/v4"
)

// Define a TournamentModel type which wraps a pgx.Conn connection pool.
type TournamentModel struct {
	Db *pgx.Conn
}

// This will insert a new tournament into the database.
func (m *TournamentModel) Insert(title, datetime, tournament_type string) (int, error) {
	stmt := `INSERT INTO tournaments (title, datetime, tournament_type)
	VALUES ($1, $2, $3) RETURNING id`

	row := m.Db.QueryRow(context.Background(), stmt, title, datetime, tournament_type)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// This will return a specific tournament based on its id.
func (m *TournamentModel) Get(id int) (*models.Tournament, error) {
	return nil, nil
}

// This will return all tournaments.
func (m *TournamentModel) All() ([]*models.Tournament, error) {
	stmt := `SELECT * FROM tournaments ORDER BY datetime ASC`

	rows, err := m.Db.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tournaments := []*models.Tournament{}

	for rows.Next() {
		t := &models.Tournament{}

		err = rows.Scan(&t.Id, &t.Title, &t.DateTime, &t.Type)
		if err != nil {
			return nil, err
		}

		tournaments = append(tournaments, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tournaments, nil
}
