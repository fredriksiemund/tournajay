package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/jackc/pgx/v4"
)

type UserModel struct {
	Db *pgx.Conn
}

func (m *UserModel) Insert(id, name, email, picture string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := `INSERT INTO users (id, name, email, picture)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO UPDATE SET name = $2, email = $3, picture = $4`

	_, err := m.Db.Exec(ctx, stmt, id, name, email, picture)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) One(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := "SELECT u.id, u.name, u.email, u.picture FROM users u WHERE u.id = $1"
	row := m.Db.QueryRow(ctx, stmt, id)

	u := &models.User{}
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Picture)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return u, nil
}
