package models

import (
    "app/types"
)

type Product struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    Name string `gorm:"column:name" binding:"required,max=50"`
    Price types.Float64 `gorm:"column:price"`
    BrandId types.Int32 `gorm:"column:brand_id"`
    Image *string `gorm:"column:image" binding:"omitempty,max=50"`
    CreateUser *types.Int32 `gorm:"column:create_user"`
    CreateDate *types.Date `gorm:"column:create_date"`
}

func (Product) TableName() string {
    return "Product"
}

type ProductUpdate struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    Name string `gorm:"column:name" binding:"required,max=50"`
    Price types.Float64 `gorm:"column:price"`
    BrandId types.Int32 `gorm:"column:brand_id"`
    Image *string `gorm:"column:image" binding:"omitempty,max=50"`
}

func (ProductUpdate) TableName() string {
    return "Product"
}
