package models

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
