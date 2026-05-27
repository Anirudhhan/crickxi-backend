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

func CreateBallEvent(tx *sqlx.Tx, delivery models.Delivery) error {
	query := `INSERT INTO balls(innings_id ,ball_sequence, over_number, ball_in_over, is_free_hit, 
                  is_legal_delivery, striker_id, non_striker_id, bowler_id, runs_batter,
                  runs_extra, extra_type, is_wicket, wicket_type, wicket_player_id, fielder_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err := tx.Exec(query, delivery.InningsID, delivery.BallSequence, delivery.OverNumber, delivery.BallInOver,
		delivery.IsFreeHit, delivery.IsLegalDelivery, delivery.StrikerID, delivery.NonStrikerID, delivery.BowlerID,
		delivery.RunsBatter, delivery.RunsExtra, delivery.ExtraType, delivery.IsWicket, delivery.WicketType, delivery.WicketPlayerID, delivery.FielderID)

	return err
}

func UpdateInnings(tx *sqlx.Tx, delivery models.Delivery, totalRuns int, wicket int, legalBall int, extraWide int, extraNoBall int) error {

	query := `UPDATE innings
			SET total_runs = total_runs + $1,
				wickets = wickets + $2,
				legal_balls = legal_balls + $3,
				extras = extras + $4,
				extras_wides = extras_wides + $5,
				extras_no_balls = extras_no_balls + $6,
				updated_at = NOW()
			WHERE id = $7`

	_, err := tx.Exec(query, totalRuns, wicket, legalBall, delivery.RunsExtra, extraWide, extraNoBall, delivery.InningsID)

	return err
}

func UpdateBatterStats(tx *sqlx.Tx, delivery models.Delivery, balls int, fours int, sixes int) error {
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

func UpdateDismissedBatter(tx *sqlx.Tx, delivery models.Delivery, dismissalBy *string) error {
	battingQuery := `UPDATE batting_scorecards
					SET
						is_out = true,
						dismissal_type = $1,
						dismissal_by = $2,
						updated_at = NOW()
					WHERE innings_id = $3 AND player_id = $4`

	_, err := tx.Exec(battingQuery, delivery.WicketType, dismissalBy, delivery.InningsID, delivery.WicketPlayerID)
	return err
}

func UpdateBowlingScoreCard(tx *sqlx.Tx, delivery models.Delivery, legalBall int, runs int, wides int, noBall int, wicket int) error {
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

func UpdateLiveMatch(tx *sqlx.Tx, delivery models.Delivery, matchID string, totalRuns int, wickets int, legalBalls int,
	strikerID string, nonStrikerID *string, nextFreeHit bool) error {

	query := `UPDATE live_match
				SET
					current_score = current_score + $1,
					wickets = wickets + $2,
					legal_balls = legal_balls + $3,
					current_ball_sequence = current_ball_sequence + 1,
					striker_id = $4,
					non_striker_id = $5,
					current_bowler_id = $6,
					is_free_hit = $7,
					updated_at = NOW()
				WHERE match_id = $8`

	_, err := tx.Exec(query, totalRuns, wickets, legalBalls, strikerID, nonStrikerID, delivery.BowlerID, nextFreeHit, matchID)
	return err
}
