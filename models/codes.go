package models

import (
	"database/sql"
	"time"
)

type SentCode struct {
	UserId int64
	Code   string
	SentAt time.Time
}

func SendCode(db *sql.DB, phone string, code int) error {
	const query = `
    INSERT INTO sms_codes (user_id, code)
    SELECT u.id, $1
    FROM users u
    WHERE u.phone = $2
    ON CONFLICT (user_id) DO UPDATE
    SET code = $1, sent_at = CURRENT_TIMESTAMP
    RETURNING sent_at;
    `
	var sentAt time.Time
	row := db.QueryRow(query, code, phone)
	if err := row.Scan(&sentAt); err != nil {
		return err
	}
	return nil
}

func LastCodeSentAt(db *sql.DB, phone string) (*time.Time, error) {
	const query = `
    SELECT sent_at FROM sms_codes 
    WHERE 
    user_id = (SELECT id FROM users WHERE phone = $1)
    `
	var sentAt *time.Time
	row := db.QueryRow(query, phone)
	if err := row.Scan(sentAt); err != nil {
		return nil, err
	}
	return sentAt, nil
}

func LastCode(db *sql.DB, userId int64) (*SentCode, error) {
	const query = `
    SELECT user_id, code, sent_at FROM sms_codes
    WHERE user_id = $1
    `
	sentCode := new(SentCode)
	row := db.QueryRow(query, userId)
	err := row.Scan(&sentCode.UserId, &sentCode.Code, &sentCode.SentAt)
	if err != nil {
		return nil, err
	}
	return sentCode, nil
}

func DeleteCode(db *sql.DB, userId int64) error {
	const query = `
    DELETE FROM sms_codes 
    WHERE user_id = $1
    `
	_, err := db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}
