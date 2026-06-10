package models

type OversDetails struct {
	BallSequence     int     `db:"ball_sequence"`
	OverNumber       int     `db:"over_number"`
	BallInOver       int     `db:"ball_in_over"`
	IsFreeHit        bool    `db:"is_free_hit"`
	RunsBatter       int     `db:"runs_batter"`
	RunsExtra        int     `db:"runs_extra"`
	ExtraType        *string `db:"extra_type"`
	IsWicket         bool    `db:"is_wicket"`
	WicketType       *string `db:"wicket_type"`
	StrikerName      string  `db:"striker_name"`
	BowlerName       string  `db:"bowler_name"`
	WicketPlayerName *string `db:"wicket_player_name"`
	FielderName      *string `db:"fielder_name"`
}
type OverBall struct {
	Ball    string `json:"ball"`
	Display string `json:"display"`
}

type OverResponse struct {
	OverNumber int        `json:"overNumber"`
	Balls      []OverBall `json:"balls"`
}
