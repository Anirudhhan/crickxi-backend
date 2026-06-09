package models

type BallEventReq struct {
	Runs           int     `json:"runs"`
	ExtraRuns      int     `json:"extraRuns"`
	ExtraType      *string `json:"extraType"`
	IsWicket       *bool   `json:"isWicket"`
	WicketType     *string `json:"wicketType"`
	WicketPlayerID *string `json:"wicketPlayerID"`
	FielderID      *string `json:"fielderID"`
	NextBatterID   *string `json:"nextBatterID"`
}
type Delivery struct {
	InningsID       string  `db:"innings_id"`
	BallSequence    int     `db:"ball_sequence"`
	OverNumber      int     `db:"over_number"`
	BallInOver      int     `db:"ball_in_over"`
	IsFreeHit       bool    `db:"is_free_hit"`
	IsLegalDelivery bool    `db:"is_legal_delivery"`
	LegalBalls      int     `db:"legal_balls"`
	StrikerID       string  `db:"striker_id"`
	NonStrikerID    *string `db:"non_striker_id"`
	BowlerID        string  `db:"bowler_id"`
	RunsBatter      int     `db:"runs_batter"`
	RunsExtra       int     `db:"runs_extra"`
	ExtraType       *string `db:"extra_type"`
	IsWicket        bool    `db:"is_wicket"`
	WicketType      *string `db:"wicket_type"`
	WicketPlayerID  *string `db:"wicket_player_id"`
	FielderID       *string `db:"fielder_id"`
	NextBatterID    *string `db:"next_batter_id"`
}
