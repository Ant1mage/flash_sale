package repositories

import (
	"database/sql"
	"flash-sale/common"
	"flash-sale/datamodels"
	"fmt"
	"github.com/jinzhu/gorm"
)

type OrderRepository interface {
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (orderMap map[int]map[string]string, err error)
}

type OrderDao struct {
	db *gorm.DB
}

func NewOrderDao(db *gorm.DB) OrderRepository {
	return &OrderDao{
		db: db,
	}
}

func (o *OrderDao) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.db.Create(&order).Error; err != nil {
		return
	}
	productID = order.ProductID
	return
}

func (o *OrderDao) Delete(productID int64) bool {
	if err := o.db.Where(&datamodels.Order{ProductID: productID}).Delete(datamodels.Order{}).Error; err != nil {
		return false
	}
	return true
}

func (o *OrderDao) Update(order *datamodels.Order) (err error) {
	if err = o.db.Model(&order).Where("orderNum = ?", order.OrderNum).Updates(&order).Error; err != nil {
		return
	}
	fmt.Println(order)
	return
}

func (o *OrderDao) SelectByKey(orderNum int64) (order *datamodels.Order, err error) {
	order = &datamodels.Order{}
	if err = o.db.Where("orderNum = ?", orderNum).Find(order).Error; err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (o *OrderDao) SelectAll() (orderArray []*datamodels.Order, err error) {
	orderArray = []*datamodels.Order{}
	if err = o.db.Find(&orderArray).Error; err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (o *OrderDao) SelectAllWithInfo() (orderMap map[int]map[string]string, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = o.db.Table("orders").Select("orders.ID,products.productName,orders.orderStatus").Joins("left join products on orders.productID=products.ID").Rows(); err != nil {
		return
	}
	orderMap = common.GetResultRows(rows)
	return
}
