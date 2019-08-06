package repositories

import (
	"flash-sale/datamodels"
	"fmt"
	"github.com/jinzhu/gorm"
)

type ProductRepository interface {
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductDao struct {
	db *gorm.DB
}

func (p *ProductDao) Insert(product *datamodels.Product) (productID int64, err error) {
	if err = p.db.Create(&product).Error; err != nil {
		return
	}
	productID = product.ID
	return
}

func (p *ProductDao) Delete(productID int64) bool {
	if err := p.db.Where(&datamodels.Product{ID: productID}).Delete(datamodels.Order{}).Error; err != nil {
		return false
	}
	return true
}

func (p *ProductDao) Update(product *datamodels.Product) (err error) {
	if err = p.db.Model(&product).Where("ID = ?", product.ID).Updates(&product).Error; err != nil {
		return
	}
	return
}

func (p *ProductDao) SelectByKey(productID int64) (product *datamodels.Product, err error) {
	product = &datamodels.Product{}
	if err = p.db.Where("ID = ?", productID).Find(product).Error; err != nil {
		return
	}
	return
}

func (p *ProductDao) SelectAll() (productArray []*datamodels.Product, err error) {
	productArray = []*datamodels.Product{}
	if err = p.db.Find(&productArray).Error; err != nil {
		fmt.Println(err)
		return
	}
	return
}

func NewProductDao(db *gorm.DB) ProductRepository {
	return &ProductDao{
		db: db,
	}
}
