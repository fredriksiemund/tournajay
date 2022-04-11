package postgres

import (
	"context"
	"errors"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/jackc/pgx/v4"
)

type UserModel struct {
	Db *pgx.Conn
}

var ctx = context.Background()

func (m *UserModel) Insert(id, name, email, picture string) error {
	stmt := `INSERT INTO users (id, name, email, picture)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO UPDATE SET name = $2, email = $3, picture = $4`

	_, err := m.Db.Exec(ctx, stmt, id, name, email, picture)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Get(id string) (*models.User, error) {
	stmt := "SELECT u.id, u.name, u.email, u.picture FROM users u WHERE u.id = $1"
	row := m.Db.QueryRow(context.Background(), stmt, id)

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
