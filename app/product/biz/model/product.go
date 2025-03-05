package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       float32 `json:"price"`

	Categories []Category `json:"categories" gorm:"many2many:product_category"`
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

type CachedProductQuery struct {
	productQuery ProductQuery
	cacheClient  *redis.Client
	prefix       string
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

func NewCachedProductQuery(ctx context.Context, db *gorm.DB, cacheClient *redis.Client) *CachedProductQuery {
	return &CachedProductQuery{
		productQuery: *NewProductQuery(ctx, db),
		cacheClient:  cacheClient,
		prefix:       "TiktokMall",
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

func (c CachedProductQuery) GetById(productId uint) (product *Product, err error) {
	cacheKey := fmt.Sprintf("%s_%s_%d", c.prefix, "product_by_id", productId)
	cachedResult := c.cacheClient.Get(c.productQuery.ctx, cacheKey)

	err = func() error {
		if err = cachedResult.Err(); err != nil {
			return err
		}
		cachedResultByte, err := cachedResult.Bytes()
		if err != nil {
			return err
		}
		err = json.Unmarshal(cachedResultByte, &product)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		product, err = c.productQuery.GetById(productId)
		if err != nil {
			return &Product{}, err
		}
		encoded, err := json.Marshal(product)
		if err != nil {
			return product, nil
		}
		_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return
}

func (c CachedProductQuery) SearchProducts(q string) (products []Product, err error) {
	return c.productQuery.SearchProducts(q)
}
