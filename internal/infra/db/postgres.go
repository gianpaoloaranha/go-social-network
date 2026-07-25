package db

import (
	"sync"

	postgresmodel "github.com/gianpaoloaranha/go-social-network/internal/adapters/out/db/postgres/model"
	"github.com/gianpaoloaranha/go-social-network/internal/infra/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func GetPostgresInstance(cfg config.Config) (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dbInstance, err = gorm.Open(postgres.Open(cfg.PostgresConnectionString), &gorm.Config{})
	})

	return dbInstance, err
}

func RunPostgresMigrations(cfg config.Config) error {
	db, err := GetPostgresInstance(cfg)
	if err != nil {
		return err
	}

	if err = db.AutoMigrate(
		&postgresmodel.User{},
		&postgresmodel.Post{},
		&postgresmodel.Comment{},
		&postgresmodel.UserFollow{},
	); err != nil {
		return err
	}

	return nil
}
