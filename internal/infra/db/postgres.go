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

// ConnectToPostgres establishes a connection to the PostgreSQL database using the provided configuration.
func ConnectToPostgres(cfg config.Config) (*gorm.DB, func() error, error) {
	postgresDB, err := GetPostgresInstance(cfg)
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := postgresDB.DB()
	if err != nil {
		return nil, nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, nil, err
	}

	return postgresDB, sqlDB.Close, nil
}

// GetPostgresInstance returns a singleton instance of the PostgreSQL database connection.
func GetPostgresInstance(cfg config.Config) (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dbInstance, err = gorm.Open(postgres.Open(cfg.PostgresConnectionString), &gorm.Config{})
	})

	return dbInstance, err
}

// RunPostgresMigrations runs the database migrations for the PostgreSQL database.
func RunPostgresMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&postgresmodel.User{},
		&postgresmodel.Post{},
		&postgresmodel.Comment{},
		&postgresmodel.UserFollow{},
	); err != nil {
		return err
	}

	return nil
}
