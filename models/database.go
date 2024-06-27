package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

func SetupDatabase(c *DatabaseConfig) (*sql.DB, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
	fmt.Println(url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.Join(errors.New("failed to connect to the database"), err)
	}
	return db, nil
}

func MigrateDatabase(db *sql.DB) error {
	ddls := tableDDLS
	for _, ddl := range ddls {
		if _, err := db.Exec(ddl); err != nil {
			return err
		}
	}
	return nil
}

func CleanDatabase(db *sql.DB) error {
	tables := []string{"sms_codes", "users", "refresh_tokens"}
	_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", strings.Join(tables, ",")))
	if err != nil {
		return err
	}
	return nil
}
