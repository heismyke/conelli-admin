package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/conelli/admin-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDBStore struct {
	db *gorm.DB
}

func NewPostgresDBStore(db *gorm.DB) *PostgresDBStore {
	return &PostgresDBStore{
		db: db,
	}
}

func (s *PostgresDBStore) DB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return s.db
	}

	return s.db.WithContext(ctx)
}

func NewPool() (*gorm.DB, error) {
	dsn := config.Envs.DATABASE_URL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.Envs.HOST,
			config.Envs.DB_USER,
			config.Envs.PASSWORD,
			config.Envs.DB_NAME,
			config.Envs.DB_PORT,
			config.Envs.SSLMODE,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Print("connected successfully to database")

	return db, nil
}
