-- Fix foreign key references after table renames
-- SQLite doesn't update foreign key constraints when tables are renamed,
-- so tables still reference old table names (events_new, games_new, boards_new).
-- This migration recreates affected tables with correct foreign key references.

BEGIN TRANSACTION;

-- Recreate events with correct FK reference to games
CREATE TABLE events_fixed (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    display_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'CLOSED')) DEFAULT 'OPEN',
    UNIQUE(game_id, display_id),
    FOREIGN KEY (game_id) REFERENCES games(game_id)
);
INSERT INTO events_fixed SELECT * FROM events;
DROP TABLE events;
ALTER TABLE events_fixed RENAME TO events;
CREATE INDEX idx_events_game ON events(game_id);
CREATE INDEX idx_events_display ON events(game_id, display_id);

-- Recreate boards with correct FK reference to games
CREATE TABLE boards_fixed (
    board_id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    grid_size INTEGER DEFAULT 4,
    UNIQUE(game_id, user_id),
    FOREIGN KEY (game_id) REFERENCES games(game_id)
);
INSERT INTO boards_fixed SELECT * FROM boards;
DROP TABLE boards;
ALTER TABLE boards_fixed RENAME TO boards;

-- Recreate board_squares with correct FK references to boards and events
CREATE TABLE board_squares_fixed (
    board_id INTEGER,
    row INTEGER,
    column INTEGER,
    event_id INTEGER NOT NULL,
    PRIMARY KEY (board_id, row, column),
    FOREIGN KEY (board_id) REFERENCES boards(board_id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(event_id)
);
INSERT INTO board_squares_fixed SELECT * FROM board_squares;
DROP TABLE board_squares;
ALTER TABLE board_squares_fixed RENAME TO board_squares;

-- Recreate votes with correct FK reference to events
CREATE TABLE votes_fixed (
    event_id INTEGER,
    user_id INTEGER,
    voted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id),
    FOREIGN KEY (event_id) REFERENCES events(event_id)
);
INSERT INTO votes_fixed SELECT * FROM votes;
DROP TABLE votes;
ALTER TABLE votes_fixed RENAME TO votes;
CREATE INDEX idx_votes_event ON votes(event_id);

COMMIT;

