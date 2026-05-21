package models

import "time"

type Player struct {
	PlayerID  string `json:"playerID"`
	IsCaptain bool   `json:"isCaptain"`
	Phone     string `json:"phone"`
	Name      string `json:"name"`
}

type MatchData struct {
	MatchID string
	TeamAID string
	TeamBID string
}

type CreateMatchRequest struct {
	TeamAName string `json:"teamAName"`
	TeamBName string `json:"teamBName"`

	//HostID    string `json:"host_id"` //take from middleware
	ScorerID1 *string `json:"scorerID1"`
	ScorerID2 *string `json:"scorerID2"`

	Overs int `json:"overs"`

	TeamAPlayers []Player `json:"teamAPlayers"`
	TeamBPlayers []Player `json:"teamBPlayers"`

	TossWinner   string `json:"tossWinner"`
	TossDecision string `json:"tossDecision"`
}

type MatchCard struct {
	MatchID string `db:"match_id" json:"matchID"`

	TeamAID   string `db:"team_a_id" json:"teamAID"`
	TeamAName string `db:"team_a_name" json:"teamAName"`

	TeamBID   string `db:"team_b_id" json:"teamBID"`
	TeamBName string `db:"team_b_name" json:"teamBName"`

	CurrentScore int `db:"current_score" json:"currentScore"`
	Wickets      int `db:"wickets" json:"wickets"`
	LegalBalls   int `db:"legal_balls" json:"legalBalls"`

	MatchStatus string `db:"match_status" json:"matchStatus"`

	OversPerSide int `db:"overs_per_side" json:"oversPerSide"`

	StartTime time.Time `db:"start_time" json:"startTime"`

	StrikerID    *string `db:"striker_id" json:"strikerID"`
	StrikerName  *string `db:"striker_name" json:"strikerName"`
	StrikerRuns  *int    `db:"striker_runs" json:"strikerRuns"`
	StrikerBalls *int    `db:"striker_balls" json:"strikerBalls"`

	BowlerID        *string `db:"bowler_id" json:"bowlerID "`
	BowlerName      *string `db:"bowler_name" json:"bowlerName"`
	BowlerWickets   *int    `db:"bowler_wickets" json:"bowlerWickets"`
	BowlerRunsGiven *int    `db:"bowler_runs_given" json:"bowlerRunsGiven"`
}
