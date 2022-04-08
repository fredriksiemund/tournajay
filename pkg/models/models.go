package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: not matching record found")

type TournamentType struct {
	Id    int
	Title string
}

type Tournament struct {
	Id          int
	Title       string
	Description string
	DateTime    time.Time
	Type        string
}

type User struct {
	Id      int
	Name    string
	Email   string
	Picture string
}
