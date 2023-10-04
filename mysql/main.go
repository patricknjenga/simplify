package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	DB     *gorm.DB
	DSN    string
	Models []any
}

func New(dsn string, models ...any) *Mysql {
	return &Mysql{&gorm.DB{}, dsn, models}
}

func (m *Mysql) Open() error {
	var err error
	m.DB, err = gorm.Open(mysql.Open(m.DSN), &gorm.Config{})
	if err != nil {
		return err
	}
	err = m.DB.AutoMigrate(m.Models...)
	return err
}
