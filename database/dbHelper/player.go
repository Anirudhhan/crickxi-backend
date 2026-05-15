package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"
)

func GetPlayerProfileByUserID(userID string) (playerStats models.PlayerStats, err error) {
	query := `
	SELECT
		u.name,
		ps.id,
		ps.runs,
		ps.catches,
		ps.run_outs,
		ps.wickets,
		ps.matches_played,
		ps.bowling_style,
		ps.batting_style,
		ps.updated_at,
		ps.created_at
	FROM users u
	INNER JOIN player_stats ps
		ON ps.user_id = u.id
		AND ps.archived_at IS NULL
	WHERE
		u.id = $1
		AND u.archived_at IS NULL
`

	err = database.DB.Get(
		&playerStats,
		query,
		userID,
	)

	return playerStats, err
}
