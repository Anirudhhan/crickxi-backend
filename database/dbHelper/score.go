package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func IsPlayerOut(inningsID string, playerID string) (isOut bool, err error) {
	query := `SELECT is_out FROM batting_scorecards
				WHERE innings_id = $1 AND player_id = $2`

	err = database.DB.Get(&isOut, query, inningsID, playerID)
	return isOut, err
}

func GetInningsDetails(matchID string, inningOrder int) (inningDetails models.MatchScoreCard, err error) {
	query := `SELECT
				bt.id AS batting_team_id,
				bt.name AS batting_team_name,
				i.total_runs AS total_runs,
				i.extras AS extras,
				bwt.id AS bowling_team_id,
				bwt.name AS bowling_team_name
			FROM innings i
			JOIN teams bt
				ON bt.id = i.batting_team_id
			JOIN teams bwt
				ON bwt.id = i.bowling_team_id
			WHERE
				i.match_id = $1
				AND i.innings_order = $2`

	err = database.DB.Get(&inningDetails, query, matchID, inningOrder)
	return inningDetails, err
}

func GetBattingScorecardByMatchIDAndInnings(tx *sqlx.Tx, matchID string, inningOrder int) (
	battingScoreCard []models.BattingScoreCard, err error) {
	query := `SELECT bsc.player_id, u.name, i.batting_team_id, bsc.runs, bsc.balls, bsc.fours, bsc.sixes, bsc.is_out,
				bsc.dismissal_type, du.name AS dismissal_by_name
						from batting_scorecards bsc
						JOIN player_stats ps
							ON ps.id = bsc.player_id
						JOIN users u
							ON u.id = ps.user_id
						JOIN innings i
							ON i.id = bsc.innings_id
						LEFT JOIN player_stats dps
							ON dps.id = bsc.dismissal_by
						LEFT JOIN users du
							ON du.id = dps.user_id
				WHERE i.match_id = $1 AND i.innings_order = $2
				ORDER BY bsc.batting_order_position`

	if tx != nil {
		err = tx.Select(&battingScoreCard, query, matchID, inningOrder)
	} else {
		err = database.DB.Select(&battingScoreCard, query, matchID, inningOrder)
	}
	return battingScoreCard, err
}

func GetBowlingScorecardByMatchIDAndInnings(tx *sqlx.Tx, matchID string, inningOrder int) (
	bowlingScoreCard []models.BowlingScoreCard, err error) {
	query := `SELECT bwsc.player_id, u.name, i.bowling_team_id, bwsc.legal_balls, bwsc.maidens, bwsc.runs_given, bwsc.no_balls, 
				   bwsc.wides, bwsc.wickets
			from bowling_scorecards bwsc
					 JOIN player_stats ps
						  ON ps.id = bwsc.player_id
					 JOIN users u
						  ON u.id = ps.user_id
					 JOIN innings i
						  ON i.id = bwsc.innings_id
			WHERE i.match_id = $1 AND i.innings_order = $2
			ORDER BY bwsc.legal_balls DESC, bwsc.runs_given ASC`

	if tx != nil {

		err = tx.Select(&bowlingScoreCard, query, matchID, inningOrder)
	} else {
		err = database.DB.Select(&bowlingScoreCard, query, matchID, inningOrder)
	}
	return bowlingScoreCard, err
}

func GetFieldingStatsByMatchID(tx *sqlx.Tx, matchID string) (stats []models.FieldingStats, err error) {
	query := `
			SELECT
				b.fielder_id,
		
				COUNT(*) FILTER (
					WHERE b.wicket_type = 'caught'
				) AS catches,
		
				COUNT(*) FILTER (
					WHERE b.wicket_type = 'run_out'
				) AS run_outs,
		
				COUNT(*) FILTER (
					WHERE b.wicket_type = 'stumped'
				) AS stumpings
		
			FROM balls b
			JOIN innings i
				ON i.id = b.innings_id
		
			WHERE
				i.match_id = $1
				AND b.fielder_id IS NOT NULL
		
			GROUP BY b.fielder_id`

	if tx != nil {
		err = tx.Select(&stats, query, matchID)
	} else {
		err = database.DB.Select(&stats, query, matchID)
	}
	return stats, err
}

func UpdateBatterScorecard(tx *sqlx.Tx, delivery models.Delivery, balls int, fours int, sixes int) error {
	battingQuery := `UPDATE batting_scorecards
					SET
						runs = runs + $1,
						balls = balls + $2,
						fours = fours + $3,
						sixes = sixes + $4,
						updated_at = NOW()
					WHERE innings_id = $5 AND player_id = $6`

	_, err := tx.Exec(battingQuery, delivery.RunsBatter, balls, fours, sixes, delivery.InningsID, delivery.StrikerID)
	return err
}

func UpdateDismissedBatter(tx *sqlx.Tx, delivery models.Delivery, dismissalBy *string, isOut bool) error {
	battingQuery := `UPDATE batting_scorecards
					SET
						is_out = $1,
						dismissal_type = $2,
						dismissal_by = $3,
						updated_at = NOW()
					WHERE innings_id = $4 AND player_id = $5`

	_, err := tx.Exec(battingQuery, isOut, delivery.WicketType, dismissalBy, delivery.InningsID, delivery.WicketPlayerID)
	return err
}

func ClearBattersDismissal(tx *sqlx.Tx, inningsID string, strikerID string, nonStrikerID *string) error {
	query := `UPDATE batting_scorecards 
              SET dismissal_type = NULL, dismissal_by = NULL, is_out = false, updated_at = NOW() 
              WHERE innings_id = $1 AND (player_id = $2 OR player_id = $3)`

	_, err := tx.Exec(query, inningsID, strikerID, nonStrikerID)
	return err
}

func UpdateBowlingScorecard(tx *sqlx.Tx, delivery models.Delivery, legalBall int, runs int, wides int, noBall int, wicket int) error {
	query := `UPDATE bowling_scorecards
				SET legal_balls = legal_balls + $1,
					maidens = maidens + $2,
					runs_given = runs_given + $3,
					wides = wides + $4,
					no_balls = no_balls + $5,
					wickets = wickets + $6
				WHERE innings_id = $7 AND player_id = $8`

	_, err := tx.Exec(query, legalBall, 0, runs, wides, noBall, wicket, delivery.InningsID, delivery.BowlerID)

	return err
}
