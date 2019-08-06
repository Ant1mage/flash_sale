package datamodels

import (
	"time"
)

type Order struct {
	ID          int64     `gorm:"column:ID"`
	OrderNum    string    `gorm:"column:orderNum"`
	UserID      int64     `gorm:"column:userID"`
	ProductID   int64     `gorm:"column:productID"`
	OrderStatus int64     `gorm:"column:orderStatus"`
	CreatedAt   time.Time `gorm:"column:createdAt"`
}

const (
	OrderWaited = iota
	OrderSuccess
	OrderFailed
)
