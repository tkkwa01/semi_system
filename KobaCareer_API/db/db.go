package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func NewDB() *gorm.DB {
	var err error

	url := fmt.Sprintf(
		"root:password@tcp(db:3306)/koba-career?charset=utf8mb4&parseTime=True&loc=Local",
	)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to MySQL")
	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
