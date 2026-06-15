package models

import "time"

type Player struct {
	PlayerID  string `json:"playerID" db:"player_id"`
	IsCaptain bool   `json:"isCaptain" db:"is_captain"`
	Phone     string `json:"phone"`
	Name      string `json:"name"`
}
type PlayerStats struct {
	ID     string `db:"id" json:"id"`
	UserID string `db:"user_id" json:"user_id"`

	Name string `db:"name" json:"name"`

	// batting
	Runs          int `db:"runs" json:"runs"`
	BallsFaced    int `db:"balls_faced" json:"balls_faced"`
	InningsBatted int `db:"innings_batted" json:"innings_batted"`
	NotOuts       int `db:"not_outs" json:"not_outs"`

	Fours int `db:"fours" json:"fours"`
	Sixes int `db:"sixes" json:"sixes"`

	HighestScore int `db:"highest_score" json:"highest_score"`

	Ducks       int `db:"ducks" json:"ducks"`
	GoldenDucks int `db:"golden_ducks" json:"golden_ducks"`

	Fifties  int `db:"fifties" json:"fifties"`
	Hundreds int `db:"hundreds" json:"hundreds"`

	// bowling
	Wickets int `db:"wickets" json:"wickets"`

	BallsBowled  int `db:"balls_bowled" json:"balls_bowled"`
	RunsConceded int `db:"runs_conceded" json:"runs_conceded"`

	MaidenOvers int `db:"maiden_overs" json:"maiden_overs"`

	Wides   int `db:"wides" json:"wides"`
	NoBalls int `db:"no_balls" json:"no_balls"`

	BestBowlingWickets int `db:"best_bowling_wickets" json:"best_bowling_wickets"`
	BestBowlingRuns    int `db:"best_bowling_runs" json:"best_bowling_runs"`

	InningsBowled int `db:"innings_bowled" json:"innings_bowled"`

	// fielding
	Catches   int `db:"catches" json:"catches"`
	RunOuts   int `db:"run_outs" json:"run_outs"`
	Stumpings int `db:"stumpings" json:"stumpings"`

	// match stats
	MatchesPlayed int `db:"matches_played" json:"matches_played"`
	MatchesWon    int `db:"matches_won" json:"matches_won"`
	MatchesLost   int `db:"matches_lost" json:"matches_lost"`

	// fantasy/game points
	TotalPoints int `db:"total_points" json:"total_points"`
	MVPs        int `db:"mvps" json:"mvps"`

	// styles
	BowlingStyle *string `db:"bowling_style" json:"bowling_style"`
	BattingStyle *string `db:"batting_style" json:"batting_style"`

	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ArchivedAt *time.Time `db:"archived_at" json:"archived_at"`
}

type UpdateProfileRequest struct {
	Name         string  `json:"name" binding:"required,min=2"`
	BattingStyle *string `json:"battingStyle"`
	BowlingStyle *string `json:"bowlingStyle"`
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

type FieldingStats struct {
	PlayerID  string `db:"fielder_id"`
	Catches   int    `db:"catches"`
	RunOuts   int    `db:"run_outs"`
	Stumpings int    `db:"stumpings"`
}
type UpdatePlayerStats struct {
	// batting
	Runs          *int `json:"runs"`
	BallsFaced    *int `json:"ballsFaced"`
	InningsBatted *int `json:"inningsBatted"`
	NotOuts       *int `json:"notOuts"`
	Fours         *int `json:"fours"`
	Sixes         *int `json:"sixes"`

	HighestScore *int `json:"highestScore"`

	Ducks       *int `json:"ducks"`
	GoldenDucks *int `json:"goldenDucks"`
	Fifties     *int `json:"fifties"`
	Hundreds    *int `json:"hundreds"`

	// bowling
	Wickets       *int `json:"wickets"`
	BallsBowled   *int `json:"ballsBowled"`
	RunsConceded  *int `json:"runsConceded"`
	MaidenOvers   *int `json:"maidenOvers"`
	Wides         *int `json:"wides"`
	NoBalls       *int `json:"noBalls"`
	BestWickets   *int `json:"bestBowlingWickets"`
	BestRuns      *int `json:"bestBowlingRuns"`
	InningsBowled *int `json:"inningsBowled"`

	// fielding
	Catches   *int `json:"catches"`
	RunOuts   *int `json:"runOuts"`
	Stumpings *int `json:"stumpings"`

	// matches
	MatchesPlayed *int `json:"matchesPlayed"`
	MatchesWon    *int `json:"matchesWon"`
	MatchesLost   *int `json:"matchesLost"`
}
