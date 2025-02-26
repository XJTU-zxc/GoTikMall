package model

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserId    uint32 `gorm:"type:int(11);not null;index:idx_user_id"`
	ProductId uint32 `gorm:"type:int(11);not null"`
	Quantity  int32  `gorm:"type:int(11);not null"`
}

func (Cart) TableName() string {
	return "cart"
}

func AddItem(db *gorm.DB, ctx context.Context, cart *Cart) (err error) {
	var row Cart
	err = db.WithContext(ctx).Model(&Cart{}).Where(&Cart{UserId: cart.UserId, ProductId: cart.ProductId}).First(&row).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if row.ID != 0 {
		err = db.WithContext(ctx).Model(&Cart{}).Where(&Cart{UserId: cart.UserId, ProductId: cart.ProductId}).UpdateColumn("quantity", gorm.Expr("quantity + ?", cart.Quantity)).Error
	} else {
		err = db.WithContext(ctx).Model(&Cart{}).Create(cart).Error
	}
	return
}

func EmptyCart(db *gorm.DB, ctx context.Context, userId uint32) (err error) {
	if userId == 0 {
		return errors.New("user_id is required")
	}
	return db.WithContext(ctx).Delete(&Cart{}, "user_id = ?", userId).Error
}

func GetCartByUserId(db *gorm.DB, ctx context.Context, userId uint32) (cartList []*Cart, err error) {
	err = db.WithContext(ctx).Model(&Cart{}).Where(&Cart{UserId: userId}).Find(&cartList).Error
	return
}
