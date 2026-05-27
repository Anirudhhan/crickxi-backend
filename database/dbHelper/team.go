package dbHelper

import (
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func CreateTeam(tx *sqlx.Tx, name string, createdBy string) (teamID string, err error) {

	query := `INSERT INTO teams(name,created_by)
				VALUES($1, $2) RETURNING id`

	err = tx.Get(&teamID, query, name, createdBy)
	return teamID, err
}

func AddPlayersToTeam(tx *sqlx.Tx, teamID string, players []models.Player) error {

	query := `INSERT INTO team_players( team_id, player_id,	is_captain)
			VALUES($1, $2, $3)`

	for _, player := range players {
		_, err := tx.Exec(query, teamID, player.PlayerID, player.IsCaptain)
		if err != nil {
			return err
		}
	}
	return nil
}
