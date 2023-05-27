package models

import (
	"time"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID        int32     `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	Username  string    `gorm:"column:username;type:text;not null;unique" json:"username"`
	Email     string    `gorm:"column:email;type:text;not null;unique" json:"email"`
	Password  string    `gorm:"column:password;type:text;not null" json:"password"`
	Bio       *string   `gorm:"column:bio;type:text" json:"bio"`
	Image     *string   `gorm:"column:image;type:text" json:"image"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type UserRegister struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type UserRegisterRequest struct {
	User UserRegister `json:"user" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	User UserLogin `json:"user" binding:"required"`
}

type UserUpdate struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Username *string `json:"username" `
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

type UserUpdateRequest struct {
	User UserUpdate `json:"user" binding:"required"`
}

type UserProfile struct {
	Username  string  `json:"username"`
	Bio       *string `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

type UserProfileResponse struct {
	Profile UserProfile `json:"profile"`
}

type UserCommon struct {
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
	Token    string  `json:"token"`
}

type UserResponse struct {
	User *UserCommon `json:"user"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
