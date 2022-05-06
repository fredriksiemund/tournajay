package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/jackc/pgx/v4"
)

// Define a TournamentModel type which wraps a pgx.Conn connection pool.
type TournamentModel struct {
	Db *pgx.Conn
}

// This will insert a new tournament into the database.
func (m *TournamentModel) Insert(title, description, date, tournamentTypeId, creatorId string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := `INSERT INTO tournaments (title, description, date, tournament_type_id, creator_id)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	row := m.Db.QueryRow(ctx, stmt, title, description, date, tournamentTypeId, creatorId)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// This will return a specific tournament based on its id.
func (m *TournamentModel) One(id int) (*models.Tournament, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := `SELECT t.id, t.title, t.description, t.date, tt.id, tt.name, u.id, u.name, u.email, u.picture
	FROM tournaments t
	INNER JOIN tournament_types tt ON t.tournament_type_id = tt.id
	INNER JOIN users u ON t.creator_id = u.id
	WHERE t.id = $1`

	row := m.Db.QueryRow(ctx, stmt, id)

	t := &models.Tournament{}

	err := row.Scan(
		&t.Id,
		&t.Title,
		&t.Description,
		&t.Date,
		&t.Type.Id,
		&t.Type.Name,
		&t.Creator.Id,
		&t.Creator.Name,
		&t.Creator.Email,
		&t.Creator.Picture,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	stmt = `SELECT u.id, u.name, u.email, u.picture FROM participants p
	INNER JOIN users u ON p.user_id = u.id
	WHERE p.tournament_id = $1`

	rows, err := m.Db.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	participants := []models.User{}

	for rows.Next() {
		u := models.User{}

		err = rows.Scan(&u.Id, &u.Name, &u.Email, &u.Picture)
		if err != nil {
			return nil, err
		}

		participants = append(participants, u)
	}

	t.Participants = participants

	return t, nil
}

// This will return all tournaments.
func (m *TournamentModel) All() ([]*models.Tournament, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := `SELECT t.id, t.title, t.description, t.date, tt.id, tt.name, u.id, u.name, u.email, u.picture
	FROM tournaments t
	INNER JOIN tournament_types tt ON t.tournament_type_id = tt.id
	INNER JOIN users u ON t.creator_id = u.id
	ORDER BY t.date ASC`

	rows, err := m.Db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tournaments := []*models.Tournament{}

	for rows.Next() {
		t := &models.Tournament{}

		err = rows.Scan(
			&t.Id,
			&t.Title,
			&t.Description,
			&t.Date,
			&t.Type.Id,
			&t.Type.Name,
			&t.Creator.Id,
			&t.Creator.Name,
			&t.Creator.Email,
			&t.Creator.Picture,
		)
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

func (m *TournamentModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := "DELETE FROM tournaments WHERE id = $1"

	_, err := m.Db.Exec(ctx, stmt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNoRecord
		} else {
			return err
		}
	}

	return nil
}
