package entity

import "time"

type Task struct {
	ID          uint64    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Status      bool      `db:"status"`
	Created_at  time.Time `db:"created_at"`
	Updated_at  time.Time `db:"updated_at"`
	User_id     uint64    `db:"user_id"`
}
