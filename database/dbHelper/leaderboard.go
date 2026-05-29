package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"
)

func GetBattingLeaderboard(page int, limit int) (leaderboard []models.BattingLeaderboard, err error) {
	offset := (page - 1) * limit
	query := `SELECT
					bsc.player_id,
					u.name AS player_name,
					COUNT(DISTINCT i.match_id) AS matches,
					COUNT(*) AS innings,
					COALESCE(SUM(bsc.runs), 0) AS runs,
					COALESCE(SUM(bsc.balls), 0) AS balls,
					COALESCE(SUM(bsc.fours), 0) AS fours,
					COALESCE(SUM(bsc.sixes), 0) AS sixes,
		
					ROUND(
						COALESCE(
							(SUM(bsc.runs)::numeric * 100) /
							NULLIF(SUM(bsc.balls), 0),
						0),
					2) AS strike_rate
				FROM batting_scorecards bsc
				INNER JOIN innings i
					ON i.id = bsc.innings_id
				INNER JOIN player_stats ps
					ON ps.id = bsc.player_id
				INNER JOIN users u
					ON u.id = ps.user_id
				GROUP BY
					bsc.player_id,
					u.name
				ORDER BY
					runs DESC,
					strike_rate DESC
		
				LIMIT $1 OFFSET $2`

	err = database.DB.Select(&leaderboard, query, limit, offset)
	return leaderboard, err
}
