package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: not matching record found")

type Tournament struct {
	Id   int
	Name string
	Date time.Time
}
