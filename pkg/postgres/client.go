package postgres

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"warehouse-system/config"
)

type Client struct {
	log *log.Logger
	db  *sql.DB
}

func (client *Client) Close() {
	if err := client.db.Close(); err != nil {
		client.log.Printf("Unable to close postgres client: %s\n", err)
	}
}

func NewClient(log *log.Logger, config *config.AppConfig) *Client {
	log.SetPrefix("[postgres client] ")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresPort,
		config.PostgresDB, config.PostgresSslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Unable to open connection: %s\n", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Printf("Unable to ping postgres: %s\n", err)
		return nil
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("Unable to create driver for migrations: %s\n", err)
		return nil
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		config.PostgresMigrationsPath, config.PostgresDB, driver)
	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Printf("Unable to make migrations: %s\n", err)
		return nil
	}

	return &Client{
		db:  db,
		log: log,
	}
}
