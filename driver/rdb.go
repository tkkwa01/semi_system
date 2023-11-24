package driver

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"semi_systems/attendance/domain"
	"semi_systems/config"
	"time"
)

var db *gorm.DB

func NewDB() (*gorm.DB, error) {
	// 環境変数からMySQLの接続情報を取得
	cfg := config.Env.DB
	dsn := fmt.Sprintf("root:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		//cfg.User,
		cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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

func init() {
	connectWithRetry()
}

func connectWithRetry() {
	var err error

	operation := func() error {
		var con string

		if config.Env.DB.Socket != "" {
			con = fmt.Sprintf("unix(%s)", config.Env.DB.Socket)
		} else {
			con = fmt.Sprintf("tcp(%s:%d)", config.Env.DB.Host, config.Env.DB.Port)
		}

		dsn := fmt.Sprintf(
			"%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Env.DB.User,
			config.Env.DB.Password,
			con,
			config.Env.DB.Name,
		)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		return err
	}

	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = 2 * time.Minute // 例: 最大2分間再試行します

	if err := backoff.Retry(operation, backOff); err != nil {
		log.Fatalf("Failed to connect to database after retries: %v", err)
	}
}

func GetRDB() *gorm.DB {
	return db
}
