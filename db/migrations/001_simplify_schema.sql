-- Migration: Simplify schema, add display_id, fix foreign keys
-- This migration:
-- 1. Adds display_id to events for user-facing numbering
-- 2. Simplifies events primary key from (event_id, game_id) to event_id
-- 3. Adds proper foreign key constraints
-- 4. Enforces single active game via unique index
-- 5. Adds CHECK constraint for event status

BEGIN TRANSACTION;

-- Create new tables with improved schema
CREATE TABLE games_new (
    game_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    is_active BOOLEAN DEFAULT 0,
    grid_size INTEGER DEFAULT 4
);

CREATE UNIQUE INDEX idx_active_game ON games_new(is_active) WHERE is_active = 1;

CREATE TABLE events_new (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    display_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'CLOSED')) DEFAULT 'OPEN',
    UNIQUE(game_id, display_id),
    FOREIGN KEY (game_id) REFERENCES games_new(game_id)
);

CREATE INDEX idx_events_game ON events_new(game_id);
CREATE INDEX idx_events_display ON events_new(game_id, display_id);

CREATE TABLE boards_new (
    board_id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    grid_size INTEGER DEFAULT 4,
    UNIQUE(game_id, user_id),
    FOREIGN KEY (game_id) REFERENCES games_new(game_id)
);

CREATE TABLE board_squares_new (
    board_id INTEGER,
    row INTEGER,
    column INTEGER,
    event_id INTEGER NOT NULL,
    PRIMARY KEY (board_id, row, column),
    FOREIGN KEY (board_id) REFERENCES boards_new(board_id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events_new(event_id)
);

CREATE TABLE votes_new (
    event_id INTEGER,
    user_id INTEGER,
    voted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id),
    FOREIGN KEY (event_id) REFERENCES events_new(event_id)
);

CREATE INDEX idx_votes_event ON votes_new(event_id);

-- Copy data from old tables
-- Games: unchanged structure
INSERT INTO games_new SELECT * FROM games;

-- Events: need to add display_id (assuming events are already ordered by event_id per game)
-- This assumes event_id currently represents the display order within each game
INSERT INTO events_new (event_id, game_id, display_id, description, status)
SELECT event_id, game_id, event_id as display_id, description, status
FROM events;

-- Boards: unchanged structure
INSERT INTO boards_new SELECT * FROM boards;

-- Board squares: unchanged structure
INSERT INTO board_squares_new SELECT * FROM board_squares;

-- Votes: remove game_id from primary key (it's redundant now)
INSERT INTO votes_new (event_id, user_id, voted_at)
SELECT event_id, user_id, voted_at FROM votes;

-- Drop old tables
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS board_squares;
DROP TABLE IF EXISTS boards;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS games;

-- Rename new tables to original names
ALTER TABLE games_new RENAME TO games;
ALTER TABLE events_new RENAME TO events;
ALTER TABLE boards_new RENAME TO boards;
ALTER TABLE board_squares_new RENAME TO board_squares;
ALTER TABLE votes_new RENAME TO votes;

-- Enable foreign keys (must be set per connection)
PRAGMA foreign_keys = ON;

COMMIT;

-- Verification queries (run these after migration to check integrity)
-- Check all board_squares reference valid events:
--   SELECT COUNT(*) FROM board_squares bs LEFT JOIN events e ON bs.event_id = e.event_id WHERE e.event_id IS NULL;
--   Expected: 0
--
-- Check all votes reference valid events:
--   SELECT COUNT(*) FROM votes v LEFT JOIN events e ON v.event_id = e.event_id WHERE e.event_id IS NULL;
--   Expected: 0
--
-- Check only one active game:
--   SELECT COUNT(*) FROM games WHERE is_active = 1;
--   Expected: 0 or 1
--
-- Check all events have valid display_ids:
--   SELECT game_id, COUNT(*), MIN(display_id), MAX(display_id) FROM events GROUP BY game_id;
--   Display IDs should be sequential per game

