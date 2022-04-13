package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDuplicate = errors.New("models: a record with this primary key already exists")

type User struct {
	Id      string
	Name    string
	Email   string
	Picture string
}

type TournamentType struct {
	Id    int
	Title string
}

type Tournament struct {
	Id           int
	Title        string
	Description  string
	DateTime     time.Time
	Type         TournamentType
	Creator      User
	Participants []User
}
