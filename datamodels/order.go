package datamodels

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	OrderNum    string    `gorm:"column:order_num"`
	UserID      int64     `gorm:"column:user_id"`
	ProductID   int64     `gorm:"column:product_id"`
	OrderStatus int64     `gorm:"column:order_status"`
}

const (
	OrderWaited = iota
	OrderSuccess
	OrderFailed
)
