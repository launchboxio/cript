package storage

import (
	"github.com/launchboxio/cript/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlAdapter(conf config.Mysql) (Storage, error) {
	db, err := gorm.Open(mysql.Open(conf.Url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DatabaseAdapter{
		DB: db,
	}, nil
}
