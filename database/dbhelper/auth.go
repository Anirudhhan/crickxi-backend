package dbhelper

import (
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
