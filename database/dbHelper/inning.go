package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func StartInnings(tx *sqlx.Tx, matchID string, battingTeamID string, bowlingTeamID string, inningsOrder int, inningsType string) (InningsID string, err error) {
	query := `INSERT INTO innings(match_id, batting_team_id, bowling_team_id, innings_order, innings_type) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = tx.Get(&InningsID, query, matchID, battingTeamID, bowlingTeamID, inningsOrder, inningsType)
	return InningsID, err
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

func CreateBattingScorecards(tx *sqlx.Tx, inningID string, players []models.Player) error {

	query := `INSERT INTO batting_scorecards(innings_id,player_id)
				VALUES($1, $2)`

	for _, player := range players {
		_, err := tx.Exec(query, inningID, player.PlayerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateBowlingScorecards(tx *sqlx.Tx, inningID string, players []models.Player) error {

	query := `INSERT INTO bowling_scorecards(innings_id,player_id)
				VALUES($1, $2)`

	for _, player := range players {
		_, err := tx.Exec(query, inningID, player.PlayerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func CompleteInnings(tx *sqlx.Tx, inningsID string) error {
	query := `UPDATE innings
			SET
				is_completed = true,
				updated_at = NOW()
		WHERE id = $1`

	_, err := tx.Exec(query, inningsID)
	return err
}

func OverDetails(inningID string) (overDetails []models.OversDetails, err error) {
	query := `SELECT
				over_number,
				ball_in_over,
				is_free_hit,
				runs_batter,
				runs_extra,
				extra_type,
				is_wicket
			FROM balls
			WHERE innings_id = $1 AND archived_at IS NULL
			ORDER BY ball_sequence ASC`

	err = database.DB.Select(&overDetails, query, inningID)
	return overDetails, err
}
