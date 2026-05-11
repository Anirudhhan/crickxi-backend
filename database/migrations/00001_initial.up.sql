BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE batting_style AS ENUM (
    'right_hand',
    'left_hand'
);

CREATE TYPE bowling_style AS ENUM (
    'fast',
    'medium',
    'spin',
    'off_spin',
    'leg_spin'
);

CREATE TYPE toss_decision AS ENUM (
    'bat',
    'bowl'
);

CREATE TYPE match_status AS ENUM (
    'upcoming',
    'live',
    'completed',
    'abandoned'
);

CREATE TYPE innings_number AS ENUM (
    'first',
    'second'
);

CREATE TYPE wicket_type AS ENUM (
    'bowled',
    'caught',
    'lbw',
    'run_out',
    'stumped',
    'hit_wicket',
    'retired_hurt'
);

CREATE TYPE extra_type AS ENUM (
    'wide',
    'no_ball',
    'bye',
    'leg_bye'
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS active_user_phone
    ON users(TRIM(phone))
    WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS user_session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    name TEXT NOT NULL,
    phone TEXT,
    runs BIGINT DEFAULT 0,
    wickets BIGINT DEFAULT 0,
    matches_played BIGINT DEFAULT 0,
    bowling_style bowling_style,
    batting_style batting_style,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_players_user_id
    ON players(user_id);

CREATE INDEX IF NOT EXISTS idx_players_phone
    ON players(phone);

CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    toss_winner_match_team_id UUID,
    winner_match_team_id UUID,
    toss_decision toss_decision,
    host_id UUID REFERENCES players(id),
    scorer1_id UUID REFERENCES players(id),
    scorer2_id UUID REFERENCES players(id),
    match_status match_status DEFAULT 'upcoming',
    overs_per_side INT DEFAULT 10,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_matches_status
    ON matches(match_status);

CREATE INDEX IF NOT EXISTS idx_matches_created_at
    ON matches(created_at DESC);

CREATE TABLE IF NOT EXISTS match_teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID REFERENCES matches(id) NOT NULL,
    name TEXT NOT NULL,
    color_hex TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_match_teams_match_id
    ON match_teams(match_id);

ALTER TABLE matches
    ADD CONSTRAINT fk_matches_toss_team
        FOREIGN KEY (toss_winner_match_team_id)
            REFERENCES match_teams(id);

ALTER TABLE matches
    ADD CONSTRAINT fk_matches_winner_team
        FOREIGN KEY (winner_match_team_id)
            REFERENCES match_teams(id);

CREATE TABLE IF NOT EXISTS innings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID REFERENCES matches(id) NOT NULL,
    batting_team_id UUID REFERENCES match_teams(id) NOT NULL,
    bowling_team_id UUID REFERENCES match_teams(id) NOT NULL,
    innings_number innings_number NOT NULL,
    total_runs INT DEFAULT 0,
    wickets INT DEFAULT 0,
    legal_balls INT DEFAULT 0,
    extras INT DEFAULT 0,
    extras_wides INT DEFAULT 0,
    extras_no_balls INT DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
);

CREATE INDEX IF NOT EXISTS idx_innings_match_id
    ON innings(match_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_innings_match_number
    ON innings(match_id, innings_number);

CREATE TABLE IF NOT EXISTS match_team_players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_team_id UUID REFERENCES match_teams(id) NOT NULL,
    player_id UUID REFERENCES players(id) NOT NULL,
    batting_order INT,
    is_captain BOOLEAN DEFAULT FALSE,
    is_wicket_keeper BOOLEAN DEFAULT FALSE,
    is_common_player BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ,

    UNIQUE(match_team_id, player_id)
);

CREATE INDEX IF NOT EXISTS idx_match_team_players_player
    ON match_team_players(player_id);

CREATE TABLE IF NOT EXISTS balls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    innings_id UUID REFERENCES innings(id) NOT NULL,
    ball_sequence INT NOT NULL,
    over_number INT NOT NULL,
    ball_in_over INT NOT NULL,
    is_legal_delivery BOOLEAN DEFAULT TRUE,
    striker_id UUID REFERENCES players(id) NOT NULL,
    non_striker_id UUID REFERENCES players(id) NOT NULL,
    bowler_id UUID REFERENCES players(id) NOT NULL,
    runs_batter INT DEFAULT 0,
    runs_extra INT DEFAULT 0,
    extra_type extra_type,
    is_wicket BOOLEAN DEFAULT FALSE,
    wicket_type wicket_type,
    wicket_player_id UUID REFERENCES players(id),
    fielder_id UUID REFERENCES players(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ,

    UNIQUE(innings_id, ball_sequence)
);

CREATE INDEX IF NOT EXISTS idx_balls_over
    ON balls(innings_id, over_number);

CREATE INDEX IF NOT EXISTS idx_balls_striker
    ON balls(striker_id);

CREATE INDEX IF NOT EXISTS idx_balls_bowler
    ON balls(bowler_id);

CREATE TABLE IF NOT EXISTS batting_scorecards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    innings_id UUID REFERENCES innings(id) NOT NULL,
    player_id UUID REFERENCES players(id) NOT NULL,
    batting_order_position INT,
    runs INT DEFAULT 0,
    balls INT DEFAULT 0,
    fours INT DEFAULT 0,
    sixes INT DEFAULT 0,
    strike_rate DECIMAL(5,2) DEFAULT 0,
    dismissal_type wicket_type,
    dismissal_by UUID REFERENCES players(id),
    is_out BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,

    UNIQUE(innings_id, player_id),
);

CREATE TABLE IF NOT EXISTS bowling_scorecards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    innings_id UUID REFERENCES innings(id) NOT NULL,
    player_id UUID REFERENCES players(id) NOT NULL,
    legal_balls INT DEFAULT 0,
    maidens INT DEFAULT 0,
    runs_given INT DEFAULT 0,
    wides INT DEFAULT 0,
    no_balls INT DEFAULT 0,
    wickets INT DEFAULT 0,
    economy DECIMAL(5,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,

    UNIQUE(innings_id, player_id),
);

CREATE TABLE IF NOT EXISTS live_match (
    match_id UUID PRIMARY KEY REFERENCES matches(id),
    current_inning_id UUID REFERENCES innings(id),
    striker_id UUID REFERENCES players(id),
    non_striker_id UUID REFERENCES players(id),
    current_bowler_id UUID REFERENCES players(id),
    current_score INT DEFAULT 0,
    wickets INT DEFAULT 0,
    legal_balls INT DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
);

COMMIT;