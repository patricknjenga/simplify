package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(user string, pass string, host string, port string, database string, models ...any) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/sys?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database)).Error
	if err != nil {
		return nil, err
	}
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, database)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, db.AutoMigrate(models...)
}
