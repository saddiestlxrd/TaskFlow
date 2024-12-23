package models

import "time"

type User struct {
	id        int32     `json:"id"`
	username  string    `json:"username"`
	password  string    `json:"password"`
	email     string    `json:"email"`
	createdAt time.Time `json:"created_at"`
	updatedAt time.Time `json:"updated_at"`
}

//func (user *User) Id() int32 {
//	return user.id
//}

type Task struct {
	id          int32     `json:"id"`
	title       string    `json:"title"`
	description string    `json:"description"`
	status      bool      `json:"status"`
	createdAt   time.Time `json:"created_at"`
	updatedAt   time.Time `json:"updated_at"`
	userId      int32     `json:"user_id"`
}

//func (task *Task) Id() int32 {
//	return task.id
//}
