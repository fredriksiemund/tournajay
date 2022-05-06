package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fredriksiemund/tournament-planner/pkg/models"
	"github.com/jackc/pgx/v4"
)

type TeamModel struct {
	Db *pgx.Conn
}

func (m *TeamModel) InsertOne(tournamentId int, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := "INSERT INTO teams (tournament_id, name) VALUES ($1, $2)"

	_, err := m.Db.Exec(ctx, stmt, tournamentId, name)
	if err != nil {
		return err
	}

	return nil
}

func (m *TeamModel) InsertMany(tournamentId int, nbrOfTeams int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var placeholders []string
	var values []interface{}

	for i := 0; i < nbrOfTeams; i++ {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", 2*i+1, 2*i+2))
		values = append(values, tournamentId, fmt.Sprintf("Team %d", i+1))
	}

	stmt := fmt.Sprintf("INSERT INTO teams (tournament_id, name) VALUES %s RETURNING id", strings.Join(placeholders, ", "))

	rows, err := m.Db.Query(ctx, stmt, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (m *TeamModel) All(tournamentId int) (map[int]*models.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt := `SELECT u.id, u.name, u.email, u.picture, t.id, t.name FROM participants p
	INNER JOIN users u ON p.user_id = u.id
	INNER JOIN teams t ON p.team_id = t.id
	WHERE p.tournament_id = $1
	ORDER BY t.id`

	rows, err := m.Db.Query(ctx, stmt, tournamentId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	teams := make(map[int]*models.Team)

	for rows.Next() {
		t := &models.Team{}
		u := &models.User{}

		err = rows.Scan(&u.Id, &u.Name, &u.Email, &u.Picture, &t.Id, &t.Name)
		if err != nil {
			return nil, err
		}

		team, exists := teams[t.Id]
		if exists {
			team.Members = append(team.Members, *u)
		} else {
			t.Members = append(t.Members, *u)
			teams[t.Id] = t
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
