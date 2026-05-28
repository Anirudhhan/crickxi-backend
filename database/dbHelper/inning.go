package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func StartInning(tx *sqlx.Tx, matchID string, battingTeamID string, bowlingTeamID string, inningsOrder int, inningsType string) (InningID string, err error) {
	query := `INSERT INTO innings(match_id, batting_team_id, bowling_team_id, innings_order, innings_type) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = tx.Get(&InningID, query, matchID, battingTeamID, bowlingTeamID, inningsOrder, inningsType)
	return InningID, err
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
			WHERE innings_id = $1
			ORDER BY ball_sequence ASC`

	err = database.DB.Select(&overDetails, query, inningID)
	return overDetails, err
}
