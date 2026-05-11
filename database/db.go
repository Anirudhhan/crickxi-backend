package database

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sqlx.DB

type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
)

func ConnectAndMigrate(host, port, dbName, user, password string, sslMode SSLMode) error {
	connectionStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbName, user, password, sslMode,
	)

	db, err := sqlx.Connect("postgres", connectionStr)
	if err != nil {
		return err
	}
	DB = db

	return migrateUp(DB)
}

func migrateUp(db *sqlx.DB) error {
	fmt.Println("Starting database migrations...")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

func Tx(fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}

		err = tx.Commit()
	}()

	err = fn(tx)
	return
}
