package models

type BattingScoreCard struct {
	PlayerID        string  `db:"player_id" json:"playerID"`
	PlayerName      string  `db:"name" json:"playerName"`
	Runs            int     `db:"runs" json:"runs"`
	Balls           int     `db:"balls" json:"balls"`
	Fours           int     `db:"fours" json:"fours"`
	Sixes           int     `db:"sixes" json:"sixes"`
	IsOut           bool    `db:"is_out" json:"isOut"`
	DismissalType   *string `db:"dismissal_type" json:"dismissalType"`
	DismissalByName *string `db:"dismissal_by_name" json:"dismissalByName"`
}
type BowlingScoreCard struct {
	PlayerID   string `db:"player_id" json:"playerID"`
	PlayerName string `db:"name" json:"playerName"`
	LegalBalls int    `db:"legal_balls" json:"legalBalls"`
	Maidens    int    `db:"maidens" json:"maidens"`
	RunsGiven  int    `db:"runs_given" json:"runsGiven"`
	NoBalls    int    `db:"no_balls" json:"noBalls"`
	Wides      int    `db:"wides" json:"wides"`
	Wickets    int    `db:"wickets" json:"wickets"`
}

type MatchScoreCard struct {
	InningsOrder    int    `json:"inningsOrder"`
	BattingTeamID   string `json:"battingTeamID" db:"batting_team_id"`
	BattingTeamName string `json:"battingTeamName" db:"batting_team_name"`

	BowlingTeamID   string `json:"bowlingTeamID" db:"bowling_team_id"`
	BowlingTeamName string `json:"bowlingTeamName" db:"bowling_team_name"`

	BattingScoreCard []BattingScoreCard `json:"battingScoreCard"`
	BowlingScoreCard []BowlingScoreCard `json:"bowlingScoreCard"`
}
