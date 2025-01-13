package user

import "github.com/jackc/pgx/v5"

type User interface {
	Add()
}

type user struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) User {
	return &user{
		db: db,
	}
}

func (u *user) Add() {
}
