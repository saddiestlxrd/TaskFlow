package repository

import (
	"TaskFlow/app/repository/user"
	"github.com/jackc/pgx/v5"
)

// repository -> user repo
//            -> task repo

type Repository interface {
	GetUser() user.User
}

type repository struct {
	userRepository user.User
}

func New(db *pgx.Conn) Repository {
	return &repository{
		userRepository: user.New(db),
	}
}

func (r *repository) GetUser() user.User {
	return r.userRepository
}
