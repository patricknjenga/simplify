package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(dsn string, name string, models ...any) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name)).Error
	if err != nil {
		return nil, err
	}
	err = db.Exec(fmt.Sprintf("USE %s", name)).Error
	if err != nil {
		return nil, err
	}
	return db, db.AutoMigrate(models...)
}
