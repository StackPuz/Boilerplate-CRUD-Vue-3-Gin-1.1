package models

import (
    "app/types"
)

type Brand struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    Name string `gorm:"column:name" binding:"required,max=50"`
}

func (Brand) TableName() string {
    return "Brand"
}

type BrandUpdate struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    Name string `gorm:"column:name" binding:"required,max=50"`
}

func (BrandUpdate) TableName() string {
    return "Brand"
}
