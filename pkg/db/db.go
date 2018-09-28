package db

import (
	"github.com/jinzhu/gorm"
	"github.com/mgutz/logxi/v1"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/binarycode/trewoga/pkg/model"
)

type Scope func(*gorm.DB) *gorm.DB

var db *gorm.DB

func Open(path string) {
	var err error

	db, err = gorm.Open("sqlite3", path)
	if err != nil {
		log.Fatal("Unable to open DB", "err", err)
	}

	db.AutoMigrate(&model.Service{})
	db.AutoMigrate(&model.User{})
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func GetService(q model.Service) (service model.Service, err error) {
	err = db.Where(q).First(&service).Error
	return
}

func ListServices() (services []model.Service, err error) {
	err = db.Find(&services).Error
	return
}

func ListScopedServices(scope Scope) (services []model.Service, err error) {
	err = db.Scopes(scope).Find(&services).Error
	return
}

func ListSubscribedServices(user model.User) (services []model.Service, err error) {
	err = db.Model(&user).Association("Services").Find(&services).Error
	return
}

func SaveService(service *model.Service) error {
	return db.Save(service).Error
}

func DestroyService(service *model.Service) error {
	return db.Delete(service).Error
}

func GetUser(q model.User) (user model.User, err error) {
	err = db.Where(q).First(&user).Error
	return
}

func ListUsers() (users []model.User, err error) {
	err = db.Find(&users).Error
	return
}

func ListSubscribedUsers(service model.Service) (users []model.User, err error) {
	err = db.Model(&service).Related(&users, "Users").Error
	return
}

func SaveUser(user *model.User) error {
	return db.Save(user).Error
}

func DestroyUser(user *model.User) error {
	return db.Delete(user).Error
}

func Subscribe(user model.User, service model.Service) error {
	return db.Model(&user).Association("Services").Append(&service).Error
}

func Unsubscribe(user model.User, service model.Service) error {
	return db.Model(&user).Association("Services").Delete(&service).Error
}
