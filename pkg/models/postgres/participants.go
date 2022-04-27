package postgres

import (
	"errors"
	"fmt"
	"strings"

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

func (m *ParticipantModel) AssignTeams(tournamentId int, participants []models.User, teamsIds []int) error {
	var placeholders []string
	var values []interface{}

	for i := 0; i < len(participants); i++ {
		placeholders = append(placeholders, fmt.Sprintf("($%d::int, $%d::varchar, $%d::int)", 3*i+1, 3*i+2, 3*i+3))
		values = append(values, tournamentId, participants[i].Id, teamsIds[i%len(teamsIds)])
	}

	stmt := fmt.Sprintf(
		`UPDATE participants SET team_id = tmp.team_id
		FROM (VALUES %s) AS tmp(tournament_id, user_id, team_id) 
		WHERE participants.tournament_id = tmp.tournament_id AND participants.user_id = tmp.user_id`,
		strings.Join(placeholders, ", "),
	)

	_, err := m.Db.Exec(ctx, stmt, values...)
	if err != nil {
		return err
	}

	return nil
}

func (m *ParticipantModel) Delete(tournamentId int, userId string) error {
	stmt := "DELETE FROM participants WHERE tournament_id = $1 AND user_id = $2"

	_, err := m.Db.Exec(ctx, stmt, tournamentId, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrNoRecord
		} else {
			return err
		}
	}

	return nil
}
