package token

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshClaims struct {
	jwt.RegisteredClaims
	TokenID  int64    `json:"token_id"`
	DeviceID uint     `json:"device_id"`
	UserID   int64    `json:"user_id"`
	Roles    []string `json:"roles"`
}

// Creates new refresh without ID
func NewRefreshClaims(deviceId uint, userId int64, roles []string, expiresAt time.Time) *RefreshClaims {
	t := new(RefreshClaims)
	t.DeviceID = deviceId
	t.UserID = userId
	t.ExpiresAt = jwt.NewNumericDate(expiresAt)
	t.Roles = roles
	return t
}

// Parses token from token string
func RefreshFromString(token string, secret string) (*RefreshClaims, error) {
	t, err := ParseWithClaims(token, &RefreshClaims{}, secret)
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*RefreshClaims); ok {
		return claims, nil
	}
	return nil, errors.New("unknown claims type, cannot proceed")
}

// Encodes RefreshToken model into a JWT string
func (c *RefreshClaims) TokenString(secret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token_string, err := token.SignedString([]byte(secret))

	return token_string, err
}

// Inserts row that identifies token into the database (not token string)
func (t *RefreshClaims) InsertToDb(db *sql.DB) (int64, error) {
	const query = `INSERT INTO refresh_tokens (device_id, expires_at, user_id) 
    VALUES ($1,$2,$3) RETURNING id`
	row := db.QueryRow(query, t.DeviceID, t.ExpiresAt.Unix(), t.UserID)
	var lastIndertId int64
	if err := row.Scan(&lastIndertId); err != nil {
		return 0, err
	}
	return lastIndertId, nil
}

// Removes refresh token with specified id from the database
func DeleteRefresh(db *sql.DB, id int64) error {
	const query = "DELETE FROM refresh_tokens WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (t *RefreshClaims) Exists(db *sql.DB) (bool, error) {
	var exists bool
	const query = "" +
		"SELECT EXISTS(SELECT * FROM refresh_tokens " +
		"WHERE id=$1 AND user_id=$2 AND device_id=$3 AND expires_at=$4)"
	res, err := db.Query(query, t.TokenID, t.UserID, t.DeviceID, t.ExpiresAt.Unix())
	if err != nil {
		return false, err
	}
	if !res.Next() {
		return false, nil
	}
	res.Scan(&exists)
	return exists, err
}

func (t *RefreshClaims) FindID(db *sql.DB) (int64, bool, error) {
	const query = "" +
		"SELECT id FROM refresh_tokens " +
		"WHERE user_id=$1 AND device_id=$2"
	res, err := db.Query(query, t.UserID, t.DeviceID)
	if err != nil {
		return 0, false, err
	}
	if !res.Next() {
		return 0, false, nil
	}
	var id int64
	res.Scan(&id)
	return id, true, nil
}

// Updates expiredAt field of token in the database
func (t *RefreshClaims) Update(db *sql.DB, expiresAt uint64) (*RefreshClaims, error) {
	const q = "UPDATE refresh_tokens SET expires_at=$1 WHERE id=$2 AND user_id=$3 AND device_id=$4"
	_, err := db.Exec(q, expiresAt, t.TokenID, t.UserID, t.DeviceID)
	return t, err
}

func (t *RefreshClaims) InsertOrUpdate(db *sql.DB) (int64, error) {
	id, exists, err := t.FindID(db)
	if err != nil {
		return 0, errors.Join(errors.New("error trying to find refresh token ID"), err)
	}
	if !exists {
		if id, err = t.InsertToDb(db); err != nil {
			return 0, errors.Join(errors.New("error trying to insert refresh token to DB"), err)
		}
		return id, nil
	}
	old_id := t.TokenID
	t.TokenID = id
	if _, err = t.Update(db, uint64(t.ExpiresAt.Unix())); err != nil {
		return 0, errors.Join(errors.New("error trying to update refresh token in the DB"), err)
	}
	t.TokenID = old_id
	return id, nil
}
