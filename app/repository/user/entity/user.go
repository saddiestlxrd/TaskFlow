package entity

import "time"

type User struct {
	id         uint64
	username   string
	email      string
	password   string
	created_at time.Time
	updated_at time.Time
}
