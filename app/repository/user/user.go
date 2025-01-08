package user

import "github.com/jackc/pgx/v5"

type User interface {
	Add()
}

type user struct {
	db         *pgx.Conn
	username   *user.username
	email      *user.email
	password   *user.password
	created_at *user.created_ats
	update_at  *user.updated_at
}

func New(db *pgx.Conn) User {
	return &user{
		db: db,
	}
}

func (u *user) Add() {
}
