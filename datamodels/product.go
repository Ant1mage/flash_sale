package datamodels

type Product struct {
	ID           int64  `json:"id" gorm:"column:id"`
	ProductName  string `json:"ProductName" gorm:"column:product_name"`
	ProductNum   int64  `json:"ProductNum" gorm:"column:product_num"`
	ProductImage string `json:"ProductImage" gorm:"column:product_image"`
	ProductUrl   string `json:"ProductUrl" gorm:"column:product_url"`
}
