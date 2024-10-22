package models

import (
    "app/types"
)

type UserRole struct {
    UserId types.Int32 `gorm:"column:user_id;primaryKey"`
    RoleId types.Int32 `gorm:"column:role_id;primaryKey"`
}

func (UserRole) TableName() string {
    return "UserRole"
}
