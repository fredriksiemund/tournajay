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
	Id       int
	Title    string
	DateTime time.Time
	Type     int
}
