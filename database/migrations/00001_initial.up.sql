BEGIN;

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

CREATE TYPE innings_type AS ENUM (
    'normal',
    'super_over'
);

CREATE TYPE wicket_type AS ENUM (
    'bowled',
    'caught',
    'lbw',
    'run_out',
    'stumped',
    'hit_wicket',
    'retired_hurt',
    'retired_out'
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
    phone_no TEXT NOT NULL,
    password TEXT,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS active_user_phone
    ON users(TRIM(phone_no))
    WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS player_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),

    -- batting
    runs BIGINT DEFAULT 0,
    balls_faced BIGINT DEFAULT 0,
    innings_batted BIGINT DEFAULT 0,
    not_outs BIGINT DEFAULT 0,

    fours BIGINT DEFAULT 0,
    sixes BIGINT DEFAULT 0,

    highest_score INT DEFAULT 0,

    ducks BIGINT DEFAULT 0,
    golden_ducks BIGINT DEFAULT 0,

    fifties BIGINT DEFAULT 0,
    hundreds BIGINT DEFAULT 0,

    -- bowling
    wickets BIGINT DEFAULT 0,

    balls_bowled BIGINT DEFAULT 0,
    runs_conceded BIGINT DEFAULT 0,

    maiden_overs BIGINT DEFAULT 0,

    wides BIGINT DEFAULT 0,
    no_balls BIGINT DEFAULT 0,

    best_bowling_wickets INT DEFAULT 0,
    best_bowling_runs INT DEFAULT 0,

    innings_bowled BIGINT DEFAULT 0,

    -- fielding
    catches BIGINT DEFAULT 0,
    run_outs BIGINT DEFAULT 0,
    stumpings BIGINT DEFAULT 0,

    -- match stats
    matches_played BIGINT DEFAULT 0,
    matches_won BIGINT DEFAULT 0,
    matches_lost BIGINT DEFAULT 0,

    -- fantasy/game points
    total_points BIGINT DEFAULT 0,
    mvps BIGINT DEFAULT 0,

    -- styles
    bowling_style bowling_style,
    batting_style batting_style,

    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_players_user_id
    ON player_stats(user_id);

CREATE TABLE IF NOT EXISTS teams (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     name TEXT NOT NULL,
     created_by UUID REFERENCES users(id),
     created_at TIMESTAMPTZ DEFAULT NOW(),
     archived_at TIMESTAMPTZ
);


CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    toss_winner_team_id UUID,
    team_a_id UUID REFERENCES teams(id),
    team_b_id UUID REFERENCES teams(id),
    winner_team_id UUID,
    toss_decision toss_decision,
    host_id UUID REFERENCES users(id),
    scorer1_id UUID REFERENCES users(id),
    scorer2_id UUID REFERENCES users(id),
    current_inning_no INT NOT NULL,
    match_status match_status DEFAULT 'live',
    overs_per_side INT DEFAULT 10,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

ALTER TABLE matches
    ADD CONSTRAINT chk_different_teams
        CHECK (team_a_id != team_b_id);

CREATE INDEX IF NOT EXISTS idx_matches_status
    ON matches(match_status);

CREATE INDEX IF NOT EXISTS idx_matches_created_at
    ON matches(created_at DESC);

ALTER TABLE matches
    ADD CONSTRAINT fk_matches_toss_team
        FOREIGN KEY (toss_winner_team_id)
            REFERENCES teams(id);

ALTER TABLE matches
    ADD CONSTRAINT fk_matches_winner_team
        FOREIGN KEY (winner_team_id)
            REFERENCES teams(id);

CREATE TABLE IF NOT EXISTS innings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID REFERENCES matches(id) NOT NULL,
    batting_team_id UUID REFERENCES teams(id) NOT NULL,
    bowling_team_id UUID REFERENCES teams(id) NOT NULL,
    innings_order INT NOT NULL,
    innings_type innings_type DEFAULT 'normal',
    total_runs INT DEFAULT 0,
    wickets INT DEFAULT 0,
    legal_balls INT DEFAULT 0,
    extras INT DEFAULT 0,
    extras_wides INT DEFAULT 0,
    extras_no_balls INT DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_innings_match_id
    ON innings(match_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_innings_match_number
    ON innings(match_id, innings_order);

CREATE TABLE IF NOT EXISTS team_players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID REFERENCES teams(id) NOT NULL,
    player_id UUID REFERENCES player_stats(id) NOT NULL,
    batting_order INT,
    is_captain BOOLEAN DEFAULT FALSE,
    is_wicket_keeper BOOLEAN DEFAULT FALSE,
    is_common_player BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ,

    UNIQUE(team_id, player_id)
);

CREATE INDEX IF NOT EXISTS idx_team_players_player
    ON team_players(player_id);

CREATE TABLE IF NOT EXISTS balls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    innings_id UUID REFERENCES innings(id) NOT NULL,
    ball_sequence INT NOT NULL,
    over_number INT NOT NULL,
    ball_in_over INT NOT NULL,
    is_free_hit BOOLEAN DEFAULT FALSE,
    is_legal_delivery BOOLEAN DEFAULT TRUE,
    striker_id UUID REFERENCES player_stats(id) NOT NULL,
    non_striker_id UUID REFERENCES player_stats(id),
    bowler_id UUID REFERENCES player_stats(id) NOT NULL,
    runs_batter INT DEFAULT 0,
    runs_extra INT DEFAULT 0,
    extra_type extra_type,
    is_wicket BOOLEAN DEFAULT FALSE,
    wicket_type wicket_type,
    wicket_player_id UUID REFERENCES player_stats(id),
    fielder_id UUID REFERENCES player_stats(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_balls_innings_sequence_active
    ON balls(innings_id, ball_sequence)
    WHERE archived_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_balls_over
    ON balls(innings_id, over_number);

CREATE INDEX IF NOT EXISTS idx_balls_striker
    ON balls(striker_id);

CREATE INDEX IF NOT EXISTS idx_balls_bowler
    ON balls(bowler_id);

CREATE TABLE IF NOT EXISTS batting_scorecards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    innings_id UUID REFERENCES innings(id) NOT NULL,
    player_id UUID REFERENCES player_stats(id) NOT NULL,
    batting_order_position INT,
    runs INT DEFAULT 0,
    balls INT DEFAULT 0,
    fours INT DEFAULT 0,
    sixes INT DEFAULT 0,
    dismissal_type wicket_type,
    dismissal_by UUID REFERENCES player_stats(id),
    is_out BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ,

    UNIQUE(innings_id, player_id)
);

CREATE TABLE IF NOT EXISTS bowling_scorecards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    innings_id UUID REFERENCES innings(id) NOT NULL,
    player_id UUID REFERENCES player_stats(id) NOT NULL,
    legal_balls INT DEFAULT 0,
    maidens INT DEFAULT 0,
    runs_given INT DEFAULT 0,
    wides INT DEFAULT 0,
    no_balls INT DEFAULT 0,
    wickets INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    archived_at TIMESTAMPTZ,

    UNIQUE(innings_id, player_id)
);

CREATE TABLE IF NOT EXISTS live_match (
    match_id UUID PRIMARY KEY REFERENCES matches(id),
    current_inning_id UUID REFERENCES innings(id),
    striker_id UUID REFERENCES player_stats(id),
    non_striker_id UUID REFERENCES player_stats(id),
    current_bowler_id UUID REFERENCES player_stats(id),
    current_score INT DEFAULT 0,
    wickets INT DEFAULT 0,
    legal_balls INT DEFAULT 0,
    current_ball_sequence INT DEFAULT 0,
    is_free_hit BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS match_player_points (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   match_id UUID REFERENCES matches(id) NOT NULL,
   player_id UUID REFERENCES player_stats(id) NOT NULL,
   batting_points INT DEFAULT 0,
   bowling_points INT DEFAULT 0,
   fielding_points INT DEFAULT 0,
   result_points INT DEFAULT 0,
   total_points INT DEFAULT 0,
   created_at TIMESTAMPTZ DEFAULT NOW(),
   updated_at TIMESTAMPTZ DEFAULT NOW(),

   UNIQUE(match_id, player_id)
);

COMMIT;