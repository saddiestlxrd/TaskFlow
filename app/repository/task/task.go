package task

import (
	"context"
	"errors"

	"TaskFlow/app/repository/task/entity"
	"github.com/jackc/pgx/v5"
)

type Task interface {
	Add()
	Delete()
	Update()
	Get(ctx context.Context) ([]entity.Task, error)
}

type task struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) Task {
	return &task{
		db: db,
	}
}

func (t *task) Add() {

}

func (t *task) Delete() {

}

func (t *task) Update() {

}

func (t *task) Get(ctx context.Context) ([]entity.Task, error) {
	query := "SELECT * FROM tasks"
	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Created_at, &task.Updated_at, &task.User_id); err != nil {
			return nil, errors.New("failed to scan task: " + err.Error())
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
