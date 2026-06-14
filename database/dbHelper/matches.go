package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"

	"github.com/jmoiron/sqlx"
)

func CreateMatch(tx *sqlx.Tx, req models.CreateMatchRequest, hostID string, tossWinnerTeamID string, teamAID string, teamBID string) (matchID string, err error) {

	query := `
		INSERT INTO matches(toss_winner_team_id, team_a_id, team_b_id, toss_decision, host_id, 
		                    scorer1_id, scorer2_id, current_innings_no, overs_per_side, match_status, start_time)
		VALUES($1, $2, $3, $4, $5, $6, $7, 1, $8, 'live', NOW())
-- 		TODO: take start time from user and match_status
		RETURNING id`

	err = tx.Get(&matchID, query, tossWinnerTeamID, teamAID, teamBID, req.TossDecision,
		hostID, req.ScorerID1, req.ScorerID2, req.Overs)

	return matchID, err
}

func StartLiveMatch(tx *sqlx.Tx, matchID string, inningsID string, req models.CreateMatchRequest) (err error) {
	query := `INSERT INTO live_match(match_id, current_innings_id, striker_id, non_striker_id, current_bowler_id)
				VALUES ($1, $2, $3,$4, $5)`

	_, err = tx.Exec(query, matchID, inningsID, req.StrikerID, req.NonStrikerID, req.CurrentBowlerID)
	return err
}

func GetMatches(search string, status string, hostID string, playerID string, page int, limit int) (matches []models.MatchCard, err error) {
	query := `SELECT
				m.id AS match_id,
				m.toss_winner_team_id AS toss_winner_team_id,
				m.winner_team_id AS winner_team_id, 
				m.host_id AS host_id,
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
				bwsc.wickets AS bowler_wickets,
			    COALESCE(pi.total_runs, 0)
				AS previous_innings_score
			
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
				ON bsc.player_id = sps.id AND lm.current_innings_id = bsc.innings_id
			LEFT JOIN player_stats bps
				ON lm.current_bowler_id = bps.id
			LEFT JOIN innings pi
				ON pi.match_id = m.id
				AND pi.innings_order = m.current_innings_no - 1
			LEFT JOIN users bu
				ON bps.user_id = bu.id
			LEFT JOIN bowling_scorecards bwsc
				ON bwsc.player_id = bps.id
				AND lm.current_innings_id = bwsc.innings_id
			
			WHERE
				m.archived_at IS NULL
			AND ($1 = '' OR
				ta.name ILIKE '%' || $1 || '%' OR
				tb.name ILIKE '%' || $1 || '%')
			AND ($2 = '' OR
				m.match_status = $2::match_status)
			AND ($5 = '' OR m.host_id = $5::uuid)
			AND ($6 = '' OR EXISTS (
				SELECT 1 FROM team_players tp 
				WHERE (tp.team_id = m.team_a_id OR tp.team_id = m.team_b_id) 
				AND tp.player_id = $6::uuid
			))
			ORDER BY m.created_at DESC
			LIMIT $3 OFFSET $4`

	offset := (page - 1) * limit

	err = database.DB.Select(&matches, query, search, status, limit, offset, hostID, playerID)
	return matches, err
}

func GetMatchByID(matchID string) (matchCard models.MatchCard, err error) {
	query := `
			SELECT
				m.id AS match_id,
				m.toss_winner_team_id AS toss_winner_team_id,
				m.winner_team_id AS winner_team_id,
				m.toss_decision AS toss_decision,
				m.host_id AS host_id,
				m.scorer1_id AS scorer1_id,
				m.scorer2_id AS scorer2_id,
				m.current_innings_no AS current_innings_no,
				m.match_status AS match_status,
				m.overs_per_side AS overs_per_side,
				m.start_time AS start_time,
				m.end_time AS end_time,
				m.updated_at AS match_updated_at,
				ta.id AS team_a_id,
				ta.name AS team_a_name,
				tb.id AS team_b_id,
				tb.name AS team_b_name,
				COALESCE(lm.current_score, 0)
					AS current_score,
				COALESCE(lm.wickets, 0)
					AS wickets,
				COALESCE(lm.legal_balls, 0)
					AS legal_balls,
			    COALESCE(pi.total_runs, 0)
				AS previous_innings_score,
				COALESCE(pi.legal_balls, 0)
					AS previous_innings_legal_balls,
				lm.current_bowler_id AS bowler_id,
				lm.current_innings_id AS current_innings_id,
				lm.striker_id AS striker_id,
				lm.non_striker_id AS non_striker_id,
				lm.is_free_hit AS is_free_hit,
				su.name AS striker_name,
				bsc.runs AS striker_runs,
				bsc.balls AS striker_balls,
				nsu.name AS non_striker_name,
				bnsc.runs AS non_striker_runs,
				bnsc.balls AS non_striker_balls,
				bu.name AS bowler_name,
				bwsc.runs_given AS bowler_runs_given,
				bwsc.legal_balls AS bowler_legal_balls,
				bwsc.wickets AS bowler_wickets
			
			FROM matches m
			
			JOIN teams ta
				ON ta.id = m.team_a_id
			
			JOIN teams tb
				ON tb.id = m.team_b_id
			
			LEFT JOIN live_match lm
				ON lm.match_id = m.id
			    
			LEFT JOIN innings pi
				ON pi.match_id = m.id
				AND pi.innings_order = m.current_innings_no - 1
			
			LEFT JOIN player_stats sps
				ON lm.striker_id = sps.id
			
			LEFT JOIN users su
				ON sps.user_id = su.id
			
			LEFT JOIN batting_scorecards bsc
				ON bsc.player_id = sps.id
				AND lm.current_innings_id = bsc.innings_id
			
			LEFT JOIN player_stats nsps
				ON lm.non_striker_id = nsps.id
			
			LEFT JOIN users nsu
				ON nsps.user_id = nsu.id
			
			LEFT JOIN batting_scorecards bnsc
				ON bnsc.player_id = nsps.id
				AND lm.current_innings_id = bnsc.innings_id
			
			LEFT JOIN player_stats bps
				ON lm.current_bowler_id = bps.id
			
			LEFT JOIN users bu
				ON bps.user_id = bu.id
			
			LEFT JOIN bowling_scorecards bwsc
				ON bwsc.player_id = bps.id
				AND lm.current_innings_id = bwsc.innings_id
			
			WHERE
				m.id = $1
				AND m.archived_at IS NULL
			
			ORDER BY m.created_at DESC`

	err = database.DB.Get(&matchCard, query, matchID)
	return matchCard, err
}

func GetLiveMatchDetails(matchID string) (liveMatchData models.LiveMatchDetails, err error) {
	query := `SELECT
				lm.current_innings_id,
				lm.striker_id,
				lm.non_striker_id,
				lm.current_bowler_id,
				lm.legal_balls,
				lm.current_ball_sequence,
				lm.current_score,
				lm.wickets,
				lm.is_free_hit,
				m.overs_per_side,
				m.current_innings_no,
				m.end_time,
				i.batting_team_id,
				i.bowling_team_id,
				i.is_completed,
				COUNT(bsc.player_id) AS batting_player_count,
				pi.total_runs AS previous_innings_score
			FROM live_match lm
					 LEFT JOIN matches m
							   ON m.id = lm.match_id
					 LEFT JOIN innings i
							   ON i.id = lm.current_innings_id
					 LEFT JOIN batting_scorecards bsc
							   ON bsc.innings_id = i.id
					 LEFT JOIN innings pi
							   ON pi.match_id = m.id AND pi.innings_order = m.current_innings_no - 1
			WHERE lm.match_id = $1
			GROUP BY
				lm.current_innings_id,
				lm.striker_id,
				lm.non_striker_id,
				lm.current_bowler_id,
				lm.legal_balls,
				lm.current_ball_sequence,
				lm.current_score,
				lm.wickets,
				lm.is_free_hit,
				m.overs_per_side,
				m.current_innings_no,
				m.end_time,
				i.batting_team_id,
				i.bowling_team_id,
				i.is_completed,
				pi.total_runs`

	err = database.DB.Get(&liveMatchData, query, matchID)
	return liveMatchData, err
}

func ResetLiveMatchForNextInnings(tx *sqlx.Tx, matchID string, inningsID string, req models.StartNextInningsReq) error {

	query := `UPDATE live_match
				SET
					current_innings_id = $1,
					current_score = 0,
					wickets = 0,
					legal_balls = 0,
					current_ball_sequence = 0,
					striker_id = $2,
					non_striker_id = $3,
					current_bowler_id = $4,
					is_free_hit = false,
					updated_at = NOW()
				WHERE match_id = $5`

	_, err := tx.Exec(query, inningsID, req.StrikerID, req.NonStrikerID, req.BowlerID, matchID)
	return err
}

func UpdateMatchInningsNo(tx *sqlx.Tx, matchID string, inningsNo int) error {
	query := `UPDATE matches
				SET
					current_innings_no = $1,	updated_at = NOW()
				WHERE id = $2`

	_, err := tx.Exec(query, inningsNo, matchID)
	return err
}

func ValidateBowlerID(matchID string, bowlerID string) (isValid bool, err error) {
	query := `SELECT EXISTS (
				SELECT 1
				FROM live_match lm
				INNER JOIN innings i
					ON i.id = lm.current_innings_id
				INNER JOIN team_players tp
					ON tp.team_id = i.bowling_team_id
				WHERE lm.match_id = $1
				  AND tp.player_id = $2
				  AND (
					lm.striker_id IS NULL
					OR lm.striker_id != $2
				  )
				  AND (
					lm.non_striker_id IS NULL
					OR lm.non_striker_id != $2
				  )
			)`

	err = database.DB.Get(&isValid, query, matchID, bowlerID)
	return isValid, err
}

func ChangeBowler(matchID string, bowlerID string) error {
	query := `UPDATE live_match
				SET
					current_bowler_id = $1,
					updated_at = NOW()
				WHERE match_id = $2`

	_, err := database.DB.Exec(query, bowlerID, matchID)
	return err
}

func UpdateLiveMatch(tx *sqlx.Tx, delivery models.Delivery, matchID string, totalRuns int, wickets int, legalBalls int,
	strikerID string, nonStrikerID *string, nextFreeHit bool) error {

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

	_, err := tx.Exec(query, totalRuns, wickets, legalBalls, strikerID, nonStrikerID, delivery.BowlerID, nextFreeHit, matchID)
	return err
}

func CompleteMatch(tx *sqlx.Tx, matchID string, winnerTeamID *string) error {
	query := `UPDATE matches
				SET
					match_status = 'completed',
					winner_team_id = $1,
					end_time = NOW(),
					updated_at = NOW()
				WHERE id = $2`

	_, err := tx.Exec(query, winnerTeamID, matchID)
	return err
}
