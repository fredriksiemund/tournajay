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
func (m *TournamentModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// This will return a specific tournament based on its id.
func (m *TournamentModel) Get(id int) (*models.Tournament, error) {
	return nil, nil
}

// This will return the 10 most recently created tournaments.
func (m *TournamentModel) Latest() ([]*models.Tournament, error) {
	stmt := `SELECT * FROM tournaments ORDER BY date ASC`

	rows, err := m.Db.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tournaments := []*models.Tournament{}

	for rows.Next() {
		t := &models.Tournament{}

		err = rows.Scan(&t.Id, &t.Title, &t.Date, &t.Type)
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
