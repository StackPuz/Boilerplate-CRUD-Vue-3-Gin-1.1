package models

import (
    "app/types"
)

type OrderHeader struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    CustomerId types.Int32 `gorm:"column:customer_id"`
    OrderDate types.Date `gorm:"column:order_date" binding:"required"`
}

func (OrderHeader) TableName() string {
    return "OrderHeader"
}

type OrderHeaderUpdate struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    CustomerId types.Int32 `gorm:"column:customer_id"`
    OrderDate types.Date `gorm:"column:order_date" binding:"required"`
}

func (OrderHeaderUpdate) TableName() string {
    return "OrderHeader"
}
