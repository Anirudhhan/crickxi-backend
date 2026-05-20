package models

import "time"

type PlayerStats struct {
	Id     string `db:"id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`

	Name string `db:"name" json:"name"`

	// batting
	Runs          int64 `db:"runs" json:"runs"`
	BallsFaced    int64 `db:"balls_faced" json:"balls_faced"`
	InningsBatted int64 `db:"innings_batted" json:"innings_batted"`
	NotOuts       int64 `db:"not_outs" json:"not_outs"`

	Fours int64 `db:"fours" json:"fours"`
	Sixes int64 `db:"sixes" json:"sixes"`

	HighestScore int64 `db:"highest_score" json:"highest_score"`

	Ducks       int64 `db:"ducks" json:"ducks"`
	GoldenDucks int64 `db:"golden_ducks" json:"golden_ducks"`

	Fifties  int64 `db:"fifties" json:"fifties"`
	Hundreds int64 `db:"hundreds" json:"hundreds"`

	// bowling
	Wickets int64 `db:"wickets" json:"wickets"`

	BallsBowled  int64 `db:"balls_bowled" json:"balls_bowled"`
	RunsConceded int64 `db:"runs_conceded" json:"runs_conceded"`

	MaidenOvers int64 `db:"maiden_overs" json:"maiden_overs"`

	Wides   int64 `db:"wides" json:"wides"`
	NoBalls int64 `db:"no_balls" json:"no_balls"`

	BestBowlingWickets int64 `db:"best_bowling_wickets" json:"best_bowling_wickets"`
	BestBowlingRuns    int64 `db:"best_bowling_runs" json:"best_bowling_runs"`

	InningsBowled int64 `db:"innings_bowled" json:"innings_bowled"`

	// fielding
	Catches   int64 `db:"catches" json:"catches"`
	RunOuts   int64 `db:"run_outs" json:"run_outs"`
	Stumpings int64 `db:"stumpings" json:"stumpings"`

	// match stats
	MatchesPlayed int64 `db:"matches_played" json:"matches_played"`
	MatchesWon    int64 `db:"matches_won" json:"matches_won"`
	MatchesLost   int64 `db:"matches_lost" json:"matches_lost"`

	// fantasy/game points
	TotalPoints int64 `db:"total_points" json:"total_points"`
	MVPs        int64 `db:"mvps" json:"mvps"`

	// styles
	BowlingStyle *string `db:"bowling_style" json:"bowling_style"`
	BattingStyle *string `db:"batting_style" json:"batting_style"`

	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ArchivedAt *time.Time `db:"archived_at" json:"archived_at"`
}

type UpdateProfileRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	BattingStyle string `json:"battingStyle"`
	BowlingStyle string `json:"bowlingStyle"`
}

type SearchPlayer struct {
	UserID   string `db:"user_id" json:"userID"`
	PlayerID string `db:"player_id" json:"playerID"`
	Name     string `db:"name" json:"name"`
	Phone    string `db:"phone_no" json:"phone"`
}

type CreateGuestPlayerRequest struct {
	Name  string `json:"name" binding:"required,min=2"`
	Phone string `json:"phone" binding:"required"`
}
