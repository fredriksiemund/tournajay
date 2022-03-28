package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: not matching record found")

type Tournament struct {
	Id    int
	Title string
	Date  time.Time
	Type  int
}
