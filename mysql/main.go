package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(dsn string, name string, models ...any) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name)).Error
	if err != nil {
		return err
	}
	err = db.Exec(fmt.Sprintf("USE %s", name)).Error
	if err != nil {
		return err
	}
	return db.AutoMigrate(models...)
}
