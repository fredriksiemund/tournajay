package postgres

import (
	"context"

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
