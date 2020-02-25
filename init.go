package skl

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //default
	"github.com/parnurzeal/gorequest"
)

var (
	DB  *gorm.DB
	app = gorequest.New()
)

func InitDB(dial string) error {
	var err error
	DB, err = gorm.Open("mysql", dial)
	if err != nil {
		return err
	}
	return initDB()
}

func initDB() error {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Group{})

	DB.Model(&User{}).Related(&Group{}, "Groups")
	DB.Model(&Group{}).Related(&User{}, "Users")
	return nil
}
