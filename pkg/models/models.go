package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDuplicate = errors.New("models: a record with this primary key already exists")
var ErrInvalidTree = errors.New("tournaments: each node has to have both left and right defined or none at all")

type User struct {
	Id      string `bson:"_id"`
	Name    string `bson:"name"`
	Email   string `bson:"email"`
	Picture string `bson:"picture"`
}

type TournamentType struct {
	Id   int
	Name string
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

type Team struct {
	Id      int
	Name    string
	Members []User
}

type Game struct {
	Id              int
	TeamIds         []int
	PreviousGameIds []int
	Depth           int
}
