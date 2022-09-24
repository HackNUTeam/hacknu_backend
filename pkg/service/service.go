package service

import (
	"hacknu/model"
	"hacknu/pkg/repository"
)

type User interface {
	CreateReading(location *model.LocationData) error
	GetHistory(user *model.GetLocationRequest) (*model.LocationData, error)
}
type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
	}
}
