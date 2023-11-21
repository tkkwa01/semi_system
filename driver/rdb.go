package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"semi_systems/attendance/domain"
)

func NewDB() (*gorm.DB, error) {
	// MySQLの接続情報
	dsn := "username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Infoレベルのログを出力
	})
	if err != nil {
		return nil, err
	}

	// テーブルの作成
	err = db.AutoMigrate(&domain.Attendance{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
