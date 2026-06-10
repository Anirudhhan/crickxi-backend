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

func OverDetails(matchID string, inningsOrder int) (overDetails []models.OversDetails, err error) {
	query := `
		SELECT
		    b.ball_sequence,
			b.over_number,
			b.ball_in_over,
			b.is_free_hit,
			b.runs_batter,
			b.runs_extra,
			b.extra_type,
			b.is_wicket,
			b.wicket_type,
			su.name AS striker_name,
			bu.name AS bowler_name,
			wu.name AS wicket_player_name,
			fu.name AS fielder_name

		FROM balls b
		JOIN innings i 
		    ON i.id = b.innings_id
		JOIN player_stats sps
			ON sps.id = b.striker_id
		JOIN users su
			ON su.id = sps.user_id
		JOIN player_stats bps
			ON bps.id = b.bowler_id
		JOIN users bu
			ON bu.id = bps.user_id
		LEFT JOIN player_stats wps
			ON wps.id = b.wicket_player_id
		LEFT JOIN users wu
			ON wu.id = wps.user_id
		LEFT JOIN player_stats fps
			ON fps.id = b.fielder_id
		LEFT JOIN users fu
			ON fu.id = fps.user_id
		WHERE
			i.match_id = $1
			AND i.innings_order = $2
			AND b.archived_at IS NULL
		ORDER BY b.ball_sequence ASC`

	err = database.DB.Select(&overDetails, query, matchID, inningsOrder)
	return overDetails, err
}
