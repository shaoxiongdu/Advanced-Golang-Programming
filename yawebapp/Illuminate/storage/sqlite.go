package storage

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLite(path string) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(path+"/app.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, err
}
