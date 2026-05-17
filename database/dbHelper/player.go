package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func GetPlayerProfileByID(playerStatsID string) (playerStats models.PlayerStats, err error) {
	query := `
	SELECT
		ps.id,
		ps.user_id,
		u.name,
		ps.runs,
		ps.balls_faced,
		ps.innings_batted,
		ps.not_outs,
		ps.fours,
		ps.sixes,
		ps.highest_score,
		ps.ducks,
		ps.golden_ducks,
		ps.fifties,
		ps.hundreds,
		ps.wickets,
		ps.balls_bowled,
		ps.runs_conceded,
		ps.maiden_overs,
		ps.wides,
		ps.no_balls,
		ps.best_bowling_wickets,
		ps.best_bowling_runs,
		ps.innings_bowled,
		ps.catches,
		ps.run_outs,
		ps.stumpings,
		ps.matches_played,
		ps.matches_won,
		ps.matches_lost,
		ps.total_points,
		ps.mvps,
		ps.bowling_style,
		ps.batting_style,
		ps.updated_at,
		ps.created_at,
		ps.archived_at
	FROM users u
	INNER JOIN player_stats ps
		ON ps.user_id = u.id
		AND ps.archived_at IS NULL
	WHERE
		ps.id = $1
		AND u.archived_at IS NULL
	`

	err = database.DB.Get(
		&playerStats,
		query,
		playerStatsID,
	)

	return playerStats, err
}

func UpdateUserProfile(tx *sqlx.Tx, userID string, name string, battingStyle string, bowlingStyle string) error {
	userQuery := `UPDATE users
				SET
					name = $1
					WHERE id = $2`

	_, err := tx.Exec(userQuery, name, userID)
	if err != nil {
		return err
	}

	playerQuery := `UPDATE player_stats
				SET
					batting_style = $1,
					bowling_style = $2,
					updated_at = NOW()
				WHERE user_id = $3`

	_, err = tx.Exec(playerQuery, battingStyle, bowlingStyle, userID)
	return err
}
