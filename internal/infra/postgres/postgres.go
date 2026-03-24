package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // driver import
	"github.com/pressly/goose/v3"
)

type Config struct {
	Host               string
	Port               int
	User               string
	Password           string
	DBName             string
	SSLMode            string
	ConnectTimeout     time.Duration
	MaxConnections     int
	MaxIdleConnections int
	MaxConnLifetime    time.Duration
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func Connect(ctx context.Context, c *Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", c.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(c.MaxConnections)
	db.SetMaxIdleConns(c.MaxIdleConnections)
	db.SetConnMaxLifetime(c.MaxConnLifetime)

	ctx, cancel := context.WithTimeout(ctx, c.ConnectTimeout)
	defer cancel()

	if pingErr := db.PingContext(ctx); pingErr != nil {
		return nil, fmt.Errorf("failed to ping database: %w", pingErr)
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
