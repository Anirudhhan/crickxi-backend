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

func UpdateUserProfile(tx *sqlx.Tx, userID string, name string, battingStyle *string, bowlingStyle *string) error {
	userQuery := `UPDATE users
				SET
					name = $1, updated_at = NOW()
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

func SearchPlayers(search string) (players []models.SearchPlayer, err error) {
	query := `SELECT ps.id AS player_id, u.id AS user_id, u.name, u.phone_no FROM users u 
			JOIN player_stats ps
				ON ps.user_id = u.id
			WHERE
				u.archived_at IS NULL AND ps.archived_at IS NULL
				AND (LOWER(u.name) ILIKE '%' || LOWER($1) || '%' OR u.phone_no ILIKE '%' || $2 || '%')
			ORDER BY u.name LIMIT 20`

	err = database.DB.Select(&players, query, search, search)
	return players, err
}

func GetPlayersByTeamID(teamID string) (players []models.Player, err error) {

	query := `SELECT player_id FROM team_players WHERE team_id = $1`

	err = database.DB.Select(&players, query, teamID)
	return players, err
}

func UpdatePlayerStats(tx *sqlx.Tx, playerID string, stats models.UpdatePlayerStats) error {
	query := `UPDATE player_stats
				SET
					-- batting
					runs = runs + COALESCE($1, 0),
					balls_faced = balls_faced + COALESCE($2, 0),
					innings_batted = innings_batted + COALESCE($3, 0),
					not_outs = not_outs + COALESCE($4, 0),
					fours = fours + COALESCE($5, 0),
					sixes = sixes + COALESCE($6, 0),
				
					highest_score = CASE
						WHEN $7::int IS NULL THEN highest_score
						ELSE GREATEST(COALESCE(highest_score, 0), $7::int)
					END,
				
					ducks = ducks + COALESCE($8, 0),
					golden_ducks = golden_ducks + COALESCE($9, 0),
					fifties = fifties + COALESCE($10, 0),
					hundreds = hundreds + COALESCE($11, 0),
				
					-- bowling
					wickets = wickets + COALESCE($12, 0),
					balls_bowled = balls_bowled + COALESCE($13, 0),
					runs_conceded = runs_conceded + COALESCE($14, 0),
					maiden_overs = maiden_overs + COALESCE($15, 0),
					wides = wides + COALESCE($16, 0),
					no_balls = no_balls + COALESCE($17, 0),
				
					best_bowling_wickets = CASE
						WHEN $18::int IS NULL THEN best_bowling_wickets
						WHEN COALESCE(best_bowling_wickets, 0) < $18::int THEN $18::int
						ELSE best_bowling_wickets
					END,
					
					best_bowling_runs = CASE
						WHEN $18::int IS NULL OR $19::int IS NULL THEN best_bowling_runs
					
						WHEN COALESCE(best_bowling_wickets, 0) < $18::int THEN $19::int
					
						WHEN COALESCE(best_bowling_wickets, 0) = $18::int THEN
							LEAST(
								COALESCE(best_bowling_runs, $19::int),
								$19::int
							)
					
						ELSE best_bowling_runs
					END,
				
					innings_bowled = innings_bowled + COALESCE($20, 0),
				
					-- field
					catches = catches + COALESCE($21, 0),
					run_outs = run_outs + COALESCE($22, 0),
					stumpings = stumpings + COALESCE($23, 0),
				
					-- matches
					matches_played = matches_played + COALESCE($24, 0),
					matches_won = matches_won + COALESCE($25, 0),
					matches_lost = matches_lost + COALESCE($26, 0),
					updated_at = NOW()
				
				WHERE id = $27;`

	_, err := tx.Exec(
		query,
		stats.Runs,
		stats.BallsFaced,
		stats.InningsBatted,
		stats.NotOuts,
		stats.Fours,
		stats.Sixes,
		stats.HighestScore,
		stats.Ducks,
		stats.GoldenDucks,
		stats.Fifties,
		stats.Hundreds,
		stats.Wickets,
		stats.BallsBowled,
		stats.RunsConceded,
		stats.MaidenOvers,
		stats.Wides,
		stats.NoBalls,
		stats.BestWickets,
		stats.BestRuns,
		stats.InningsBowled,
		stats.Catches,
		stats.RunOuts,
		stats.Stumpings,
		stats.MatchesPlayed,
		stats.MatchesWon,
		stats.MatchesLost,
		playerID,
	)

	return err
}
