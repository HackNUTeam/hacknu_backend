package repository

import (
	"database/sql"
	"hacknu/model"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type User interface {
	CreateReading(location *model.LocationData) error
	GetHistoryLocation(user *model.GetLocationRequest) ([]*model.LocationData, error)
	CreateUser(name string) (int64, error)
}

type Repository struct {
	User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: NewUserDB(db),
	}
}
