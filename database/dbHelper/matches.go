package dbHelper

import (
	"crickxi-backend/database"
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
	//--  todo: after implementing upcoming remove the start time

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

func GetMatches() (matches []models.MatchCard, err error) {
	query := `SELECT
				m.id AS match_id,
				ta.id AS team_a_id,
				ta.name AS team_a_name,
				tb.id AS team_b_id,
				tb.name AS team_b_name,
				COALESCE(lm.current_score, 0) AS current_score,
				COALESCE(lm.wickets, 0) AS wickets,
				COALESCE(lm.legal_balls, 0) AS legal_balls,
				m.match_status AS match_status,
				m.overs_per_side AS overs_per_side,
				m.start_time AS start_time,
				lm.striker_id AS striker_id,
				su.name AS striker_name,
				bsc.runs AS striker_runs,
				bsc.balls AS striker_balls,
				lm.current_bowler_id AS bowler_id,
				bu.name AS bowler_name,
				bwsc.runs_given AS bowler_runs_given,
				bwsc.wickets AS bowler_wickets
			
			FROM matches m
	
			JOIN teams ta
				ON ta.id = m.team_a_id		
			JOIN teams tb
				ON tb.id = m.team_b_id
			LEFT JOIN live_match lm
				ON lm.match_id = m.id
			LEFT JOIN player_stats sps
				ON lm.striker_id = sps.id
			LEFT JOIN users su
				ON sps.user_id = su.id
			LEFT JOIN batting_scorecards bsc
				ON bsc.player_id = sps.id AND lm.current_inning_id = bsc.innings_id
			LEFT JOIN player_stats bps
				ON lm.current_bowler_id = bps.id
			LEFT JOIN users bu
				ON bps.user_id = bu.id
			LEFT JOIN bowling_scorecards bwsc
				ON bwsc.player_id = bps.id
				AND lm.current_inning_id = bwsc.innings_id
			
			WHERE
				m.archived_at IS NULL
			ORDER BY m.created_at DESC`

	err = database.DB.Select(&matches, query)
	return matches, err
}
