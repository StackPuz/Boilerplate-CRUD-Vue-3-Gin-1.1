package models

import (
    "app/types"
)

type Role struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    Name string `gorm:"column:name" binding:"required,max=50"`
}

func (Role) TableName() string {
    return "Role"
}
