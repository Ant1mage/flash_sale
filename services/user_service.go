package services

import (
	"errors"
	"flash-sale/dao"
	"flash-sale/datamodels"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Login(userName string, pwd string) (user *datamodels.User, isLogin bool)
	AddUser(user *datamodels.User) (userID int64, err error)
}

type UserService struct {
	UserRepository repositories.UserRepository
}

func (u *UserService) Login(userName string, pwd string) (user *datamodels.User, isLogin bool) {
	var err error
	user, err = u.UserRepository.SelectByName(userName)
	if err != nil {
		return
	}
	isLogin, _ = ValidatePassword(pwd, user.HashPassword)

	if !isLogin {
		return &datamodels.User{}, false
	}
	return user, true
}

func (u *UserService) AddUser(user *datamodels.User) (userID int64, err error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return userID, err
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func NewUserService(repository repositories.UserRepository) IUserService {
	return &UserService{
		UserRepository:repository,
	}
}

func ValidatePassword(password string, hashed string) (isValidated bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}