package model

import (
	"context"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
    Name string `json:"name"`
	Description string `json:"description"`

	Products []Product `json:"products" gorm:"many2many:product_category"`
}

type CategoryQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (c *Category) TableName() string {
	return "category"
}

func NewCategoryQuery(ctx context.Context, db *gorm.DB) *CategoryQuery {
	return &CategoryQuery{
		ctx: ctx,
		db:  db,
	}
}

func (c *CategoryQuery) GetProductsByCategoryName(categoryName string) (category []Category, err error) {
	err = c.db.WithContext(c.ctx).Model(&Category{}).Where("name = ?", categoryName).Preload("Products").Find(&category).Error
	return
}
