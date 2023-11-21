package driver

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error

	dsn := fmt.Sprintf(
		"root:password@tcp(db:3306)/attendance-management?charset=utf8mb4&parseTime=True&loc=Local",
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

}

func GetRDB() *gorm.DB {
	return db
}
