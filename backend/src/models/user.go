package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
	Photo    string `json:"photo"`
	Phone    string `json:"phone"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (user *User) Count(db *gorm.DB) int64 {
	var total int64

	db.Model(&User{}).Count(&total)

	return total
}

func (user *User) Take(db *gorm.DB, limit int, offset int) interface{} {
	var recipes []User

	db.Offset(offset).Limit(limit).Find(&recipes)

	return recipes
}
