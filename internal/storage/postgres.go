package storage

import (
	"github.com/launchboxio/cript/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresAdapter(conf config.Postgres) (Storage, error) {
	db, err := gorm.Open(postgres.Open(conf.Url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DatabaseAdapter{
		DB: db,
	}, nil
}
