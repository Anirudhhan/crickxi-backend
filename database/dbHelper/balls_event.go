package dbHelper

import (
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
