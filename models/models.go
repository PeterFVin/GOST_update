package models

type Record struct {
	ID     int    `db:"id"`
	Number string `db:"number"`
	Name   string `db:"name"`
	State  string `db:"state"`
}
