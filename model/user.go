package model

import (
	"github.com/jinzhu/gorm"
)

// user model本体
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Gender   int    `json:"gender"` //male 1 female 2
}

func (User) TableName() string {
	return "sys.user"
}
