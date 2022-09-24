package service

import (
	"hacknu/model"
	"hacknu/pkg/repository"
	"time"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) CreateReading(location *model.LocationData) error {
	location.Timestamp = time.Now().Unix()
	return u.repo.CreateReading(location)
}

func (u *UserService) GetHistory(user *model.GetLocationRequest) (*model.LocationData, error) {
	return u.repo.GetHistoryLocation(user)
}
