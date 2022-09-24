package service

import "hacknu/pkg/repository"

type User interface {
}
type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
	}
}
