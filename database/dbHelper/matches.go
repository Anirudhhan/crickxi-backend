package dbHelper

import (
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func CreateMatch(tx *sqlx.Tx, req models.CreateMatchRequest, hostID string) (matchData models.MatchData, err error) {

	// create team A
	teamAQuery := `INSERT INTO teams(name, created_by)
					VALUES($1, $2) RETURNING id`

	err = tx.Get(&matchData.TeamAID, teamAQuery, req.TeamAName, hostID)
	if err != nil {
		return matchData, err
	}

	// create team B
	teamBQuery := `INSERT INTO teams(name, created_by)
					VALUES($1, $2) RETURNING id`

	err = tx.Get(&matchData.TeamBID, teamBQuery, req.TeamBName, hostID)
	if err != nil {
		return matchData, err
	}

	// decide toss winner
	var tossWinnerTeamID string

	if req.TossWinner == "A" {
		tossWinnerTeamID = matchData.TeamAID
	} else {
		tossWinnerTeamID = matchData.TeamBID
	}

	// create match
	matchQuery := `INSERT INTO matches(
			toss_winner_team_id, team_a_id, team_b_id, toss_decision, host_id,
			scorer1_id,	scorer2_id, current_inning_no, overs_per_side, match_status, start_time)
			VALUES($1, $2, $3, $4, $5, $6, $7, 1, $8, 'upcoming', NOW()) RETURNING id`

	err = tx.Get(&matchData.MatchID, matchQuery, tossWinnerTeamID, matchData.TeamAID, matchData.TeamBID, req.TossDecision, hostID, req.ScorerID1, req.ScorerID2, req.Overs)
	if err != nil {
		return matchData, err
	}

	// insert team players
	teamPlayerQuery := `INSERT INTO team_players(team_id, player_id, is_captain)
						VALUES($1, $2, $3)`

	// team A players
	for _, player := range req.TeamAPlayers {

		_, err = tx.Exec(teamPlayerQuery, matchData.TeamAID, player.PlayerID, player.IsCaptain)
		if err != nil {
			return matchData, err
		}
	}

	// team B players
	for _, player := range req.TeamBPlayers {

		_, err = tx.Exec(teamPlayerQuery, matchData.TeamBID, player.PlayerID, player.IsCaptain)
		if err != nil {
			return matchData, err
		}
	}

	return matchData, nil
}

//{
//"matchID": "uuid",
//
//"matchStatus": "live",
//
//"teamA": {
//"name": "Thunder Bolts",
//"shortName": "TB",
//"score": "142/2",
//"overs": "28.4"
//},
//
//"teamB": {
//"teamID": "uuid",
//"name": "Knight Riders",
//"shortName": "KR"
//},
//
//"currentRunRate": 4.95,
//
//"target": 288,
//
//"striker": {
//"playerID": "uuid",
//"name": "Babar Azam",
//"runs": 68,
//"balls": 84
//},
//
//"bowler": {
//"playerID": "uuid",
//"name": "Ravindra Jadeja",
//"wickets": 1,
//"runsGiven": 28
//}
//}
