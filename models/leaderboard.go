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
