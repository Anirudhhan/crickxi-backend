package dbHelper

import (
	"crickxi-backend/database"
	"crickxi-backend/models"
)

func GetInningDetails(matchID string, inningOrder int) (inningDetails models.MatchScoreCard, err error) {
	query := `SELECT
				bt.id AS batting_team_id,
				bt.name AS batting_team_name,
				bwt.id AS bowling_team_id,
				bwt.name AS bowling_team_name
			FROM innings i
			JOIN teams bt
				ON bt.id = i.batting_team_id
			JOIN teams bwt
				ON bwt.id = i.bowling_team_id
			WHERE
				i.match_id = $1
				AND i.innings_order = $2`

	err = database.DB.Get(&inningDetails, query, matchID, inningOrder)
	return inningDetails, err
}

func GetBattingScorecardByMatchIDAndInning(matchID string, inningOrder int) (
	battingScoreCard []models.BattingScoreCard, err error) {
	query := `SELECT bsc.player_id, u.name, bsc.runs, bsc.balls, bsc.fours, bsc.sixes, bsc.is_out,
				bsc.dismissal_type, du.name AS dismissal_by_name
						from batting_scorecards bsc
						JOIN player_stats ps
							ON ps.id = bsc.player_id
						JOIN users u
							ON u.id = ps.user_id
						JOIN innings i
							ON i.id = bsc.innings_id
						LEFT JOIN player_stats dps
							ON dps.id = bsc.dismissal_by
						LEFT JOIN users du
							ON du.id = dps.user_id
				WHERE i.match_id = $1 AND i.innings_order = $2
				ORDER BY bsc.batting_order_position`

	err = database.DB.Select(&battingScoreCard, query, matchID, inningOrder)
	return battingScoreCard, err
}

func GetBowlingScorecardByMatchIDAndInning(matchID string, inningOrder int) (
	bowlingScoreCard []models.BowlingScoreCard, err error) {
	query := `SELECT bwsc.player_id, u.name, bwsc.legal_balls, bwsc.maidens, bwsc.runs_given, bwsc.no_balls, 
				   bwsc.wides, bwsc.wickets
			from bowling_scorecards bwsc
					 JOIN player_stats ps
						  ON ps.id = bwsc.player_id
					 JOIN users u
						  ON u.id = ps.user_id
					 JOIN innings i
						  ON i.id = bwsc.innings_id
			WHERE i.match_id = $1 AND i.innings_order = $2
			ORDER BY bwsc.legal_balls DESC, bwsc.runs_given ASC`

	err = database.DB.Select(&bowlingScoreCard, query, matchID, inningOrder)
	return bowlingScoreCard, err
}
