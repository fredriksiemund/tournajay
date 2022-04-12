package postgres

import (
	"errors"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type ParticipantModel struct {
	Db *pgx.Conn
}

func (m *ParticipantModel) Insert(tournamentId int, userId string) error {
	stmt := "INSERT INTO participants (tournament_id, user_id) VALUES ($1, $2)"

	_, err := m.Db.Exec(ctx, stmt, tournamentId, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			return models.ErrDuplicate
		} else {
			return err
		}
	}

	return nil
}
