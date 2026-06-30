package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func New(connectionURL string) (*Client, error) {
	db, err := gorm.Open(postgres.Open(connectionURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to access sql database: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *Client) DB() *gorm.DB {
	return c.db
}

func (c *Client) SQLDB() (*sql.DB, error) {
	return c.db.DB()
}

func (c *Client) Ping() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
