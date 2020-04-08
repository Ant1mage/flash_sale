package datamodels

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	NickName     string `json:"nickName" form:"nick_name" gorm:"column:nick_name"`
	UserName     string `json:"userName" form:"user_name" gorm:"column:user_name"`
	HashPassword string `json:"-" form:"password" gorm:"column:password"`
}
