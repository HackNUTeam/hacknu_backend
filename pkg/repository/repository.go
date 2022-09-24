package repository

import "database/sql"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type User interface {
}

type Repository struct {
	User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: NewUserDB(db),
	}
}
