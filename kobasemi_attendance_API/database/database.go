package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"semi_systems/kobasemi_attendance_API/domain"
)

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("attendance.db"), &gorm.Config{
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
