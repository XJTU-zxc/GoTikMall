package model

import (
	"context"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name string `json:"name"`
	Description string `json:"description"`
	Picture string `json:"picture"`
	Price float32 `json:"price"`

	Categories  []Category `json:"categories" gorm:"many2many:product_category"`
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (p Product) TableName() string {
	return "product"
}

func NewProductQuery(ctx context.Context, db *gorm.DB) *ProductQuery {
	return &ProductQuery{
		ctx: ctx,
		db:  db,
	}
}
func (p *ProductQuery) GetById(productId uint) (product *Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Where("ID = ?", productId).First(&product).Error
	return
}

func (p *ProductQuery) SearchProducts(q string) (products []Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Where("name LIKE ? or description LIKE ?", "%"+q+"%", "%"+q+"%").Find(&products).Error
	return
}