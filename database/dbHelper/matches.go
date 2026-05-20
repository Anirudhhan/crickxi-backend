package dbHelper

import (
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func CreateMatch(tx *sqlx.Tx, req models.CreateMatchRequest, hostID string) (matchData models.MatchData, err error) {
	// create match
	matchQuery := `INSERT INTO matches(host_id, scorer1_id, scorer2_id, overs_per_side, start_time)
					VALUES($1, $2, $3, $4, NOW()) RETURNING id`

	err = tx.Get(&matchData.MatchID, matchQuery, hostID, req.ScorerID1, req.ScorerID2, req.Overs)
	if err != nil {
		return matchData, err
	}

	// create team A
	teamAQuery := `INSERT INTO teams(match_id,name,created_by)
					VALUES($1, $2, $3) RETURNING id`

	err = tx.Get(&matchData.TeamAID, teamAQuery, matchData.MatchID, req.TeamAName, hostID)
	if err != nil {
		return matchData, err
	}

	// create team B
	teamBQuery := `INSERT INTO teams(match_id,name,created_by)
					VALUES($1, $2, $3) RETURNING id`

	err = tx.Get(&matchData.TeamBID, teamBQuery, matchData.MatchID, req.TeamBName, hostID)
	if err != nil {
		return matchData, err
	}

	teamPlayerQuery := `INSERT INTO team_players(team_id, player_id, is_captain)
						VALUES($1, $2, $3)`
	// adding team A players
	for _, player := range req.TeamAPlayers {

		_, err = tx.Exec(teamPlayerQuery, matchData.TeamAID, player.PlayerID, player.IsCaptain)
		if err != nil {
			return matchData, err
		}
	}

	// adding team B players
	for _, player := range req.TeamBPlayers {

		_, err = tx.Exec(teamPlayerQuery, matchData.TeamBID, player.PlayerID, player.IsCaptain)
		if err != nil {
			return matchData, err
		}
	}

	// decide toss winner team id
	var tossWinnerTeamID string

	if req.TossWinner == "A" {
		tossWinnerTeamID = matchData.TeamAID
	} else {
		tossWinnerTeamID = matchData.TeamBID
	}

	// update toss details
	updateMatchQuery := `UPDATE matches
					SET toss_winner_team_id = $1, toss_decision = $2
					WHERE id = $3`

	_, err = tx.Exec(updateMatchQuery, tossWinnerTeamID, req.TossDecision, matchData.MatchID)
	if err != nil {
		return matchData, err
	}

	return matchData, nil
}
