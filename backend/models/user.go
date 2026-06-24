package models

import "time"

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"size:100;not null"`
	Email         string    `json:"email" gorm:"size:200;uniqueIndex;not null"`
	Phone         string    `json:"phone" gorm:"size:20;not null"`
	PasswordHash  string    `json:"-" gorm:"not null"`
	Role          string    `json:"role" gorm:"size:20;default:user"`
	EmailVerified bool      `json:"emailVerified" gorm:"default:false"`
	CreatedAt     time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
