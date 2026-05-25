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

func UpdateInnings(tx *sqlx.Tx, delivery models.Delivery) error {

	query := `UPDATE innings
			SET total_runs = total_runs + $1,
				wickets = wickets + $2,
				legal_balls = legal_balls + $3,
				extras = extras + $4,
				extras_wides = extras_wides + $5,
				extras_no_balls = extras_no_balls + $6,
				is_completed = COALESCE($7, is_completed),
				updated_at = NOW()
			WHERE id = $8`

	wicket := 0
	if delivery.IsWicket &&
		delivery.WicketType != nil {

		switch *delivery.WicketType {

		case "bowled",
			"caught",
			"lbw",
			"stumped",
			"hit_wicket",
			"run_out":

			wicket = 1
		}
	}

	legalBall := 0
	if delivery.IsLegalDelivery {
		legalBall = 1
	}

	extraWide := 0
	extraNoBall := 0

	if delivery.ExtraType != nil {

		if *delivery.ExtraType == "wide" {
			extraWide = delivery.RunsExtra
		}

		if *delivery.ExtraType == "no_ball" {
			extraNoBall = delivery.RunsExtra
		}
	}

	totalRuns := delivery.RunsBatter + delivery.RunsExtra

	_, err := tx.Exec(query, totalRuns, wicket, legalBall, delivery.RunsExtra, extraWide, extraNoBall, false, delivery.InningsID)

	return err
}

func UpdateBattingScoreCard(tx *sqlx.Tx, delivery models.Delivery) error {
	query := `UPDATE batting_scorecards
				SET runs = runs + $1,
					balls = balls + $2,
					fours = fours + $3,
					sixes = sixes + $4,
					dismissal_type = $5,
					dismissal_by = $6,
					is_out = $7,
					updated_at = NOW()
				WHERE innings_id = $8 AND player_id = $9`

	balls := 0
	if delivery.IsLegalDelivery {
		balls += 1
	}

	fours := 0
	if delivery.RunsBatter == 4 {
		fours += 1
	}

	sixes := 0
	if delivery.RunsBatter == 6 {
		sixes += 1
	}

	var dismissalBy any

	if delivery.WicketType != nil {
		switch *delivery.WicketType {
		case "bowled", "caught", "lbw", "stumped", "hit_wicket":
			dismissalBy = delivery.BowlerID

		default:
			dismissalBy = nil
		}
	}

	dismissedPlayerID := delivery.StrikerID
	if delivery.WicketPlayerID != nil {
		dismissedPlayerID = *delivery.WicketPlayerID
	}

	_, err := tx.Exec(query, delivery.RunsBatter, balls, fours, sixes, delivery.WicketType, dismissalBy, delivery.IsWicket, delivery.InningsID, dismissedPlayerID)
	return err
}

func UpdateBowlingScoreCard(tx *sqlx.Tx, delivery models.Delivery) error {
	query := `UPDATE bowling_scorecards
				SET legal_balls = legal_balls + $1,
					maidens = maidens + $2,
					runs_given = runs_given + $3,
					wides = wides + $4,
					no_balls = no_balls + $5,
					wickets = wickets + $6
				WHERE innings_id = $7 AND player_id = $8`

	runs := delivery.RunsBatter + delivery.RunsExtra
	legalBall := 0
	if delivery.IsLegalDelivery {
		legalBall += 1
	}

	wides := 0
	noBall := 0
	if !delivery.IsLegalDelivery {
		if delivery.ExtraType != nil {
			switch *delivery.ExtraType {
			case "no_ball":
				noBall += 1
			default:
				wides += 1
			}
		}
	}

	wicket := 0
	if delivery.IsWicket && delivery.WicketType != nil {
		switch *delivery.WicketType {
		case "bowled", "caught", "lbw", "stumped", "hit_wicket":
			wicket = 1
		}
	}

	_, err := tx.Exec(query, legalBall, 0, runs, wides, noBall, wicket, delivery.InningsID, delivery.BowlerID)

	return err
}

func UpdateLiveMatch(tx *sqlx.Tx, delivery models.Delivery, matchID string) error {

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

	totalRuns := delivery.RunsBatter + delivery.RunsExtra

	wickets := 0
	if delivery.IsWicket &&
		delivery.WicketType != nil {
		switch *delivery.WicketType {
		case "bowled", "caught", "lbw", "stumped", "hit_wicket", "run_out":
			wickets = 1
		}
	}

	legalBalls := 0
	if delivery.IsLegalDelivery {
		legalBalls = 1
	}

	nextFreeHit := false
	if delivery.ExtraType != nil {
		if *delivery.ExtraType == "no_ball" {
			nextFreeHit = true
		}
	}

	strikerID := delivery.StrikerID
	nonStrikerID := delivery.NonStrikerID

	// strike rotation
	if delivery.RunsBatter%2 == 1 {
		strikerID, nonStrikerID = nonStrikerID, strikerID
	}

	// new batter assign
	if delivery.IsWicket && delivery.NextBatterID != nil &&
		delivery.WicketPlayerID != nil {
		if *delivery.WicketPlayerID == strikerID {
			strikerID = *delivery.NextBatterID
		}

		if *delivery.WicketPlayerID == nonStrikerID {
			nonStrikerID = *delivery.NextBatterID
		}
	}

	// over completed
	newLegalBalls := delivery.LegalBalls + legalBalls
	if legalBalls == 1 && newLegalBalls%6 == 0 {
		strikerID, nonStrikerID = nonStrikerID, strikerID
	}

	_, err := tx.Exec(query, totalRuns, wickets, legalBalls, strikerID, nonStrikerID, delivery.BowlerID, nextFreeHit, matchID)
	return err
}
