package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

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

func GetLastBall(matchID string) (ball models.Delivery, err error) {
	query := `SELECT 
                b.innings_id, b.ball_sequence, b.over_number, b.ball_in_over, b.is_free_hit,
                b.is_legal_delivery, b.striker_id, b.non_striker_id, b.bowler_id,
                b.runs_batter, b.runs_extra, b.extra_type, b.is_wicket, b.wicket_type,
                b.wicket_player_id, b.fielder_id
              FROM balls b
              JOIN innings i ON i.id = b.innings_id
              WHERE i.match_id = $1 AND b.archived_at IS NULL
              ORDER BY i.innings_order DESC, b.ball_sequence DESC
              LIMIT 1`

	err = database.DB.Get(&ball, query, matchID)
	return ball, err
}

func ArchiveBall(tx *sqlx.Tx, inningsID string, ballSequence int) error {
	query := `UPDATE balls SET archived_at = NOW() WHERE innings_id = $1 AND ball_sequence = $2`
	_, err := tx.Exec(query, inningsID, ballSequence)
	return err
}
