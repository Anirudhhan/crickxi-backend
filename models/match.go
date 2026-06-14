package models

import "time"

type MatchData struct {
	MatchID          string
	CurrentInningsID string
	TeamAID          string
	TeamBID          string
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

	StrikerID       string  `json:"strikerID"`
	NonStrikerID    *string `json:"nonStrikerID"`
	CurrentBowlerID string  `json:"currentBowlerID"`
	WicketKeeperID  *string `json:"wicketKeeperID"`
	StartTime       *string `json:"startTime"`

	TossWinner   string `json:"tossWinner"`
	TossDecision string `json:"tossDecision"`
}

type MatchCard struct {
	MatchID string `db:"match_id" json:"matchID"`

	TossWinnerTeamID string  `db:"toss_winner_team_id" json:"tossWinnerTeamID"`
	WinnerTeamID     *string `db:"winner_team_id" json:"winnerTeamID"`

	TossDecision *string `db:"toss_decision" json:"tossDecision"`

	HostID    string  `db:"host_id" json:"hostID"`
	Scorer1ID *string `db:"scorer1_id" json:"scorer1ID"`
	Scorer2ID *string `db:"scorer2_id" json:"scorer2ID"`

	CurrentInningsNo *int `db:"current_innings_no" json:"currentInningsNo"`

	MatchStatus string `db:"match_status" json:"matchStatus"`

	OversPerSide int `db:"overs_per_side" json:"oversPerSide"`

	StartTime *time.Time `db:"start_time" json:"startTime"`
	EndTime   *time.Time `db:"end_time" json:"endTime"`

	MatchUpdatedAt *time.Time `db:"match_updated_at" json:"matchUpdatedAt"`

	IsFreeHit bool `db:"is_free_hit" json:"isFreeHit"`

	TeamAID   string `db:"team_a_id" json:"teamAID"`
	TeamAName string `db:"team_a_name" json:"teamAName"`

	TeamBID   string `db:"team_b_id" json:"teamBID"`
	TeamBName string `db:"team_b_name" json:"teamBName"`

	CurrentScore int `db:"current_score" json:"currentScore"`
	Wickets      int `db:"wickets" json:"wickets"`
	LegalBalls   int `db:"legal_balls" json:"legalBalls"`

	PreviousInningsScore      *int `db:"previous_innings_score" json:"previousInningsScore"`
	PreviousInningsLegalBalls *int `db:"previous_innings_legal_balls" json:"previousInningsLegalBalls"`

	CurrentInningsID *string `db:"current_innings_id" json:"currentInningsID"`

	StrikerID    *string `db:"striker_id" json:"strikerID"`
	StrikerName  *string `db:"striker_name" json:"strikerName"`
	StrikerRuns  *int    `db:"striker_runs" json:"strikerRuns"`
	StrikerBalls *int    `db:"striker_balls" json:"strikerBalls"`

	NonStrikerID    *string `db:"non_striker_id" json:"nonStrikerID"`
	NonStrikerName  *string `db:"non_striker_name" json:"nonStrikerName"`
	NonStrikerRuns  *int    `db:"non_striker_runs" json:"nonStrikerRuns"`
	NonStrikerBalls *int    `db:"non_striker_balls" json:"nonStrikerBalls"`

	BowlerID         *string `db:"bowler_id" json:"bowlerID"`
	BowlerName       *string `db:"bowler_name" json:"bowlerName"`
	BowlerRunsGiven  *int    `db:"bowler_runs_given" json:"bowlerRunsGiven"`
	BowlerLegalBalls *int    `db:"bowler_legal_balls" json:"bowlerLegalBalls"`
	BowlerWickets    *int    `db:"bowler_wickets" json:"bowlerWickets"`
}

type LiveMatchDetails struct {
	CurrentInningsID    string  `db:"current_innings_id"`
	StrikerID           string  `db:"striker_id"`
	NonStrikerID        *string `db:"non_striker_id"`
	CurrentBowlerID     string  `db:"current_bowler_id"`
	LegalBalls          int     `db:"legal_balls"`
	CurrentBallSequence int     `db:"current_ball_sequence"`
	CurrentScore        int     `db:"current_score"`
	Wickets             int     `db:"wickets"`
	IsFreeHit           bool    `db:"is_free_hit"`

	OversPerSide         int    `db:"overs_per_side"`
	CurrentInningsNo     int    `db:"current_innings_no"`
	BattingTeamID        string `db:"batting_team_id"`
	BowlingTeamID        string `db:"bowling_team_id"`
	BattingPlayerCount   int    `db:"batting_player_count"`
	PreviousInningsScore *int   `db:"previous_innings_score"`

	IsCompleted bool    `db:"is_completed"`
	EndTime     *string `db:"end_time"`
}

type StartNextInningsReq struct {
	StrikerID    string  `json:"strikerID"`
	NonStrikerID *string `json:"nonStrikerID"`
	BowlerID     string  `json:"bowlerID"`
}
