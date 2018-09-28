package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	TelegramID   int
	ChatID       int64
	FirstName    string
	LastName     string
	UserName     string
	LanguageCode string
	IsBot        bool
	Services     []Service `gorm:"many2many:subscriptions;"`
}

func (u *User) Display() {
	fmt.Println("telegram id   = ", u.TelegramID)
	fmt.Println("first name    = ", u.FirstName)
	fmt.Println("last name     = ", u.LastName)
	fmt.Println("user name     = ", u.UserName)
	fmt.Println("language code = ", u.LanguageCode)
	fmt.Println("is bot        = ", u.IsBot)
	fmt.Println()
}
