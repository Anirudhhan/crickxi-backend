package dbHelper

import (
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func ResetInningsStats(tx *sqlx.Tx, inningsID string) error {
	query := `UPDATE innings 
              SET total_runs = 0, wickets = 0, legal_balls = 0, extras = 0, 
                  extras_wides = 0, extras_no_balls = 0, is_completed = false, updated_at = NOW() 
              WHERE id = $1`
	_, err := tx.Exec(query, inningsID)
	return err
}

func ResetScorecards(tx *sqlx.Tx, inningsID string) error {
	battingQuery := `UPDATE batting_scorecards 
                     SET runs = 0, balls = 0, fours = 0, sixes = 0, is_out = false, 
                         dismissal_type = NULL, dismissal_by = NULL, updated_at = NOW() 
                     WHERE innings_id = $1`
	_, err := tx.Exec(battingQuery, inningsID)
	if err != nil {
		return err
	}

	bowlingQuery := `UPDATE bowling_scorecards 
                     SET legal_balls = 0, maidens = 0, runs_given = 0, wides = 0, no_balls = 0, wickets = 0, updated_at = NOW() 
                     WHERE innings_id = $1`
	_, err = tx.Exec(bowlingQuery, inningsID)
	return err
}

func ResetLiveMatchStats(tx *sqlx.Tx, matchID string, strikerID string, nonStrikerID *string, bowlerID string) error {
	query := `UPDATE live_match 
              SET current_score = 0, wickets = 0, legal_balls = 0, current_ball_sequence = 0, 
                  is_free_hit = false, striker_id = $1, non_striker_id = $2, current_bowler_id = $3, updated_at = NOW() 
              WHERE match_id = $4`
	_, err := tx.Exec(query, strikerID, nonStrikerID, bowlerID, matchID)
	return err
}

func GetAllActiveBalls(tx *sqlx.Tx, inningsID string) (balls []models.Delivery, err error) {
	query := `SELECT innings_id, ball_sequence, over_number, ball_in_over, is_free_hit, 
                     is_legal_delivery, striker_id, non_striker_id, bowler_id, runs_batter, 
                     runs_extra, extra_type, is_wicket, wicket_type, wicket_player_id, fielder_id, next_batter_id
              FROM balls 
              WHERE innings_id = $1 AND archived_at IS NULL 
              ORDER BY ball_sequence ASC`
	err = tx.Select(&balls, query, inningsID)
	return balls, err
}

func GetInitialInningsState(tx *sqlx.Tx, inningsID string) (strikerID string, nonStrikerID *string, bowlerID string, err error) {
	var result struct {
		StrikerID    string  `db:"striker_id"`
		NonStrikerID *string `db:"non_striker_id"`
		BowlerID     string  `db:"bowler_id"`
	}

	query := `SELECT striker_id, non_striker_id, bowler_id 
              FROM balls 
              WHERE innings_id = $1 AND archived_at IS NULL 
              ORDER BY ball_sequence ASC LIMIT 1`
	err = tx.Get(&result, query, inningsID)
	return result.StrikerID, result.NonStrikerID, result.BowlerID, err
}
