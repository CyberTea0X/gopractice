package models

import (
	"database/sql"
	"encoding/json"
	"io"
	"os"

	"github.com/lib/pq"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type User struct {
	Id    int64    `json:"id"`
	Phone string   `json:"phone"`
	Roles []string `json:"roles"`
}

// Ignores conflicts
func CreateOrUpdateUser(db *sql.DB, phone string, roles []string) (int64, error) {
	const query = `INSERT INTO users (phone, roles)
    VALUES ($1,$2)
    ON CONFLICT (phone) DO UPDATE
    SET phone = EXCLUDED.phone
    RETURNING id;
    `
	var id int64
	row := db.QueryRow(query, phone, pq.Array(roles))
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func GetUserByPhone(db *sql.DB, phone string) (*User, error) {
	user := new(User)
	user.Phone = phone
	const query = `SELECT id, roles FROM users WHERE phone=$1`
	rows := db.QueryRow(query, phone)
	err := rows.Scan(&user.Id, pq.Array(&user.Roles))
	if err != nil {
		return nil, err
	}
	return user, nil

}

func ParseUsersFromJson(file string) ([]*User, error) {
	var users []*User
	usersJson, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	usersRaw, err := io.ReadAll(usersJson)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(usersRaw, &users); err != nil {
		return nil, err
	}
	return users, nil
}
