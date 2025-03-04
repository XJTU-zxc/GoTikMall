package model

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	ProductId    uint32  `gorm:"type:int(11)"`
	OrderIdRefer string  `gorm:"type:varchar(100);index"`
	Quantity     int32   `gorm:"type:int(11)"`
	Cost         float32 `gorm:"type:float"`
}

func (OrderItem) TableName() string {
	return "order_item"
}
