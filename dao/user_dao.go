package repositories

import (
	"flash-sale/datamodels"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/core/errors"
)

type UserRepository interface {
	SelectByID(userID int64) (user *datamodels.User, err error)
	SelectByName(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userID int64, err error)
}

type UserDao struct {
	db *gorm.DB
}

func (u *UserDao) SelectByID(userID int64) (user *datamodels.User, err error) {
	user = &datamodels.User{}
	if err = u.db.Where("ID=?", userID).Find(&user).Error; err != nil {
		return
	}
	return
}

func (u *UserDao) SelectByName(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空")
	}
	user = &datamodels.User{}
	if err = u.db.Where("userName=?", userName).Find(&user).Error; err != nil {
		return
	}
	return
}

func (u UserDao) Insert(user *datamodels.User) (userID int64, err error) {
	if err = u.db.Create(&user).Error; err != nil {
		return
	}
	userID = user.ID
	return
}

func NewUserDao(db *gorm.DB) UserRepository {
	return &UserDao{
		db: db,
	}
}