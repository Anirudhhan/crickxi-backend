package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func RegisterUser(tx *sqlx.Tx, name string, phone string, hashedPassword string) (userID string, err error) {
	query := `INSERT INTO users(name, phone, password) 
	          VALUES ($1, $2, $3) 
	          RETURNING id`

	err = tx.Get(&userID, query, name, phone, hashedPassword)
	return userID, err
}

func RegisterPlayerStats(tx *sqlx.Tx, userID string) (playerID string, err error) {
	query := `INSERT INTO player_stats(user_id) 
	          VALUES ($1) 
	          RETURNING id`

	err = tx.Get(&playerID, query, userID)
	return playerID, err
}

func CreateUserSession(userID string) (sessionID string, err error) {
	query := `INSERT INTO user_sessions(user_id) VALUES($1) RETURNING id`

	err = database.DB.Get(&sessionID, query, userID)
	return sessionID, err
}

func GetLoginDetailsByPhone(phone string) (userDetails models.LoginUserDetails, err error) {
	query := `SELECT id, password
			FROM users
			WHERE phone = $1
			AND archived_at IS NULL`

	err = database.DB.Get(&userDetails, query, phone)
	return userDetails, err
}

func GetUserIDByActiveSession(sessionID string) (string, error) {
	query := `SELECT user_id 
		FROM user_sessions 
		WHERE id = $1 AND archived_at IS NULL`

	var userID string
	err := database.DB.Get(&userID, query, sessionID)
	return userID, err
}

func ArchiveUserSession(sessionID string) error {
	query := `UPDATE user_sessions SET archived_at = NOW() 
            WHERE id = $1 AND archived_at IS NULL`

	_, err := database.DB.Exec(query, sessionID)
	return err
}
