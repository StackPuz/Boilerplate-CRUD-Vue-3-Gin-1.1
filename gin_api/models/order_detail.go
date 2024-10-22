package models

import (
    "app/types"
)

type OrderDetail struct {
    OrderId types.Int32 `gorm:"column:order_id;primaryKey"`
    No types.Int32 `gorm:"column:no;primaryKey"`
    ProductId types.Int32 `gorm:"column:product_id"`
    Qty types.Int32 `gorm:"column:qty"`
}

func (OrderDetail) TableName() string {
    return "OrderDetail"
}

type OrderDetailUpdate struct {
    OrderId types.Int32 `gorm:"column:order_id;primaryKey"`
    No types.Int32 `gorm:"column:no;primaryKey"`
    ProductId types.Int32 `gorm:"column:product_id"`
    Qty types.Int32 `gorm:"column:qty"`
}

func (OrderDetailUpdate) TableName() string {
    return "OrderDetail"
}
