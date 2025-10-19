-- Fresh database schema
BEGIN TRANSACTION;

CREATE TABLE games (
    game_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    is_active BOOLEAN DEFAULT 0,
    grid_size INTEGER DEFAULT 4
);

CREATE UNIQUE INDEX idx_active_game ON games(is_active) WHERE is_active = 1;

CREATE TABLE events (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    display_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'CLOSED')) DEFAULT 'OPEN',
    UNIQUE(game_id, display_id),
    FOREIGN KEY (game_id) REFERENCES games(game_id)
);

CREATE INDEX idx_events_game ON events(game_id);
CREATE INDEX idx_events_display ON events(game_id, display_id);

CREATE TABLE boards (
    board_id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    grid_size INTEGER DEFAULT 4,
    UNIQUE(game_id, user_id),
    FOREIGN KEY (game_id) REFERENCES games(game_id)
);

CREATE TABLE board_squares (
    board_id INTEGER,
    row INTEGER,
    column INTEGER,
    event_id INTEGER NOT NULL,
    PRIMARY KEY (board_id, row, column),
    FOREIGN KEY (board_id) REFERENCES boards(board_id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(event_id)
);

CREATE TABLE votes (
    event_id INTEGER,
    user_id INTEGER,
    voted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id),
    FOREIGN KEY (event_id) REFERENCES events(event_id)
);

CREATE INDEX idx_votes_event ON votes(event_id);

PRAGMA foreign_keys = ON;

COMMIT;