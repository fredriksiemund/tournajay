package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: not matching record found")

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
	Id          int
	Title       string
	Description string
	DateTime    time.Time
	Type        TournamentType
	Creator     User
}
