package task

import "github.com/jackc/pgx/v5"

type Task interface {
	Add()
}

type task struct {
	db          *pgx.Conn
	id          *task.id
	title       *task.title
	description *task.description
	status      *task.status
	created_at  *task.created_at
	updated_at  *task.updated_at
	user_id     *task.user_id
}

func New(db *pgx.Conn) Task {
	return &task{
		db: db,
	}
}

func (t *task) Add() {

}
