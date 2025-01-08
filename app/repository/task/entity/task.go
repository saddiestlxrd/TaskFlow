package entity

import "time"

type task struct {
	id          uint64
	title       string
	description string
	status      bool
	created_at  time.Time
	updated_at  time.Time
	user_id     uint64
}
