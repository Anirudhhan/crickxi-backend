package models

import "time"

type PlayerStats struct {
	Id            string    `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Runs          int64     `db:"runs" json:"runs"`
	Catches       int64     `db:"catches" json:"catches"`
	RunOuts       int64     `db:"run_outs" json:"run_outs"`
	Wickets       int64     `db:"wickets" json:"wickets"`
	MatchesPlayed int64     `db:"matches_played" json:"matches_played"`
	BowlingStyle  *string   `db:"bowling_style" json:"bowling_style"`
	BattingStyle  *string   `db:"batting_style" json:"batting_style"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
