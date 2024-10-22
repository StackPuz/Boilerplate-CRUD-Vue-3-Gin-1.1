package models

import (
    "app/types"
)

type UserAccount struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    Name string `gorm:"column:name" binding:"required,max=50"`
    Email string `gorm:"column:email" binding:"required,max=50"`
    Password string `gorm:"column:password"`
    PasswordResetToken *string `gorm:"column:password_reset_token" binding:"omitempty,max=100"`
    Active types.Bit `gorm:"column:active"`
}

func (UserAccount) TableName() string {
    return "UserAccount"
}

type UserAccountUpdate struct {
    Id types.Int32 `gorm:"column:id;primaryKey;autoIncrement"`
    Name string `gorm:"column:name" binding:"required,max=50"`
    Email string `gorm:"column:email" binding:"required,max=50"`
    Password string `gorm:"column:password"`
    Active types.Bit `gorm:"column:active"`
}

func (UserAccountUpdate) TableName() string {
    return "UserAccount"
}
