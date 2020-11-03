package config

import (
	"fmt"
	"mvc_golang/app/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBInit() *gorm.DB {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("root:@/digitalent_bank?charset=utf8&parseTime=True&loc=Local")), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database" + err.Error())
	}
	db.AutoMigrate(new(model.Account), new(model.Transaction))
	return db
}
