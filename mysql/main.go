package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(dsn string, name string, models ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf(dsn, "")), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.Exec(fmt.Sprintf("create database if not exists %s", name)).Error
	if err != nil {
		return nil, err
	}
	db, err = gorm.Open(mysql.Open(fmt.Sprintf(dsn, name)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, db.AutoMigrate(models...)
}
