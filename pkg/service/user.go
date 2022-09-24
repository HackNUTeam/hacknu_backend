package service

import (
	"hacknu/model"
	"hacknu/pkg/repository"
	"log"
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
	log.Print(location)
	id, err := u.repo.CreateUser(location.Identifier)
	if err != nil {
		log.Println(err)
		return err
	}
	location.UserID = id
	return u.repo.CreateReading(location)
}

func (u *UserService) GetHistory(user *model.GetLocationRequest) ([]*model.LocationData, error) {
	return u.repo.GetHistoryLocation(user.Name, user.Timestamp)
}
