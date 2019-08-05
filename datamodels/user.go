package datamodels

type User struct {
	ID int64 `json:"ID" form:"ID" gorm:"column:ID"`
	NickName string `json:"nickName" form:"nickName" gorm:"column:nickName"`
	UserName string `json:"userName" form:"userName" gorm:"column:userName"`
	HashPassword string `json:"-" form:"passWord" gorm:"column:passWord"`
}