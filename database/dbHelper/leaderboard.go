package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"
)

func GetBattingLeaderboard(page int, limit int) (leaderboard []models.BattingLeaderboard, err error) {
	offset := (page - 1) * limit

	query := `SELECT
				ps.id AS player_id,
				u.name AS player_name,
	
				ps.matches_played AS matches,
				ps.innings_batted AS innings,
	
				ps.runs,
				ps.balls_faced AS balls,
				ps.fours,
				ps.sixes,
	
				ROUND(
					COALESCE(
						(ps.runs::numeric * 100) /
						NULLIF(ps.balls_faced, 0),
						0
					),
				2) AS strike_rate
	
			FROM player_stats ps
			JOIN users u
				ON u.id = ps.user_id
	
			ORDER BY
				ps.runs DESC,
				strike_rate DESC
	
			LIMIT $1 OFFSET $2`

	err = database.DB.Select(&leaderboard, query, limit, offset)

	for i := range leaderboard {
		leaderboard[i].Rank = offset + i + 1
	}

	return leaderboard, err
}

func GetBowlingLeaderboard(page int, limit int) (leaderboard []models.BowlingLeaderboard, err error) {
	offset := (page - 1) * limit
	query := `SELECT
				ps.id AS player_id,
				u.name AS player_name,
	
				ps.matches_played AS matches,
				ps.innings_bowled AS innings,
	
				ps.wickets,
				ps.balls_bowled,
				ps.runs_conceded,
	
				ROUND(
					COALESCE(
						ps.runs_conceded::numeric /
						NULLIF(ps.wickets, 0),
						0
					),
				2) AS average,
	
				ROUND(
					COALESCE(
						(ps.runs_conceded::numeric * 6) /
						NULLIF(ps.balls_bowled, 0),
						0
					),
				2) AS economy
	
			FROM player_stats ps
			JOIN users u
				ON u.id = ps.user_id
	
			ORDER BY
				ps.wickets DESC,
				average ASC
	
			LIMIT $1 OFFSET $2`

	err = database.DB.Select(&leaderboard, query, limit, offset)

	for i := range leaderboard {
		leaderboard[i].Rank = offset + i + 1
	}

	return leaderboard, err
}

func GetFieldingLeaderboard(page int, limit int) (leaderboard []models.FieldingLeaderboard, err error) {
	offset := (page - 1) * limit
	query := `SELECT
				ps.id AS player_id,
				u.name AS player_name,
	
				ps.matches_played AS matches,
	
				ps.catches,
				ps.run_outs,
				ps.stumpings,
	
				(ps.catches + ps.run_outs + ps.stumpings)
					AS dismissals
	
			FROM player_stats ps
			JOIN users u
				ON u.id = ps.user_id
	
			ORDER BY
				dismissals DESC,
				ps.catches DESC
	
			LIMIT $1 OFFSET $2`

	err = database.DB.Select(&leaderboard, query, limit, offset)

	for i := range leaderboard {
		leaderboard[i].Rank = offset + i + 1
	}

	return leaderboard, err
}
