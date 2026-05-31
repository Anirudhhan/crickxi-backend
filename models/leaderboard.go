package models

type BattingLeaderboard struct {
	Rank       int    `json:"rank"`
	PlayerID   string `db:"player_id" json:"playerID"`
	PlayerName string `db:"player_name" json:"playerName"`

	Matches int `db:"matches" json:"matches"`
	Innings int `db:"innings" json:"innings"`

	Runs  int `db:"runs" json:"runs"`
	Balls int `db:"balls" json:"balls"`
	Fours int `db:"fours" json:"fours"`
	Sixes int `db:"sixes" json:"sixes"`

	StrikeRate float64 `db:"strike_rate" json:"strikeRate"`
}

type BowlingLeaderboard struct {
	Rank       int    `json:"rank"`
	PlayerID   string `db:"player_id" json:"playerID"`
	PlayerName string `db:"player_name" json:"playerName"`

	Matches int `db:"matches" json:"matches"`
	Innings int `db:"innings" json:"innings"`

	Wickets      int `db:"wickets" json:"wickets"`
	BallsBowled  int `db:"balls_bowled" json:"ballsBowled"`
	RunsConceded int `db:"runs_conceded" json:"runsConceded"`

	Average float64 `db:"average" json:"average"`
	Economy float64 `db:"economy" json:"economy"`
}

type FieldingLeaderboard struct {
	Rank       int    `json:"rank"`
	PlayerID   string `db:"player_id" json:"playerID"`
	PlayerName string `db:"player_name" json:"playerName"`

	Matches int `db:"matches" json:"matches"`

	Catches   int `db:"catches" json:"catches"`
	RunOuts   int `db:"run_outs" json:"runOuts"`
	Stumpings int `db:"stumpings" json:"stumpings"`

	Dismissals int `db:"dismissals" json:"dismissals"`
}
