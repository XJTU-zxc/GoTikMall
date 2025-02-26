package model

import (
	"time"

	"gorm.io/gorm"
)

type ChargeReq struct {
	gorm.Model
	UserId     uint32 
	OrderId    string 
	TransactionId string 
	Amount     float32 
	PayTime   time.Time
}

