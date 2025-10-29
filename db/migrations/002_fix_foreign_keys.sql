-- Fix foreign key references after table renames
-- SQLite doesn't update foreign key constraints when tables are renamed,
-- so votes and board_squares still reference events_new instead of events.
-- This migration recreates those tables with correct foreign key references.

BEGIN TRANSACTION;

-- Recreate board_squares with correct FK reference to events
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
