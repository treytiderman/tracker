-- --------------------------------------------------------------------------------------
-- Drop
-- --------------------------------------------------------------------------------------

DROP TABLE IF EXISTS log;
DROP TABLE IF EXISTS entry;
DROP TABLE IF EXISTS option;
DROP TABLE IF EXISTS number;
DROP TABLE IF EXISTS field;
DROP TABLE IF EXISTS tracker;

-- --------------------------------------------------------------------------------------
-- Create
-- --------------------------------------------------------------------------------------

-- foreign_keys constraints are not on by default
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS tracker (
    tracker_id INTEGER NOT NULL UNIQUE,

    -- name to identify this tracker
    tracker_name TEXT NOT NULL UNIQUE,

    -- markdown formatted notes
    tracker_notes TEXT NOT NULL DEFAULT "",

    PRIMARY KEY (tracker_id)
);

CREATE TABLE IF NOT EXISTS field (
    field_id INTEGER NOT NULL UNIQUE,

    -- parent tracker
    tracker_id INTEGER NOT NULL,

    -- use "number" to track a signed whole number
    -- examples: weight, height...
    -- use "option" to a list of options
    -- examples: exercise, read status
    field_type TEXT CHECK(field_type in ('number', 'option')) NOT NULL DEFAULT 'number',

    -- name to identify this field
    field_name TEXT NOT NULL,

    -- markdown formatted notes
    field_notes TEXT NOT NULL DEFAULT "",

    -- a tracker can not have duplicate field_name's
    -- but multiple trackers can have the same field_name
    UNIQUE(tracker_id, field_name),

    PRIMARY KEY (field_id),
    FOREIGN KEY (tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS number (
    number_id INTEGER NOT NULL UNIQUE,

    -- parent field
    field_id INTEGER NOT NULL,

    -- 0 = round to integer
    -- 2 = round to 2 decimal places. example money
    -- -3 = round to thousands
    decimal_places INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (number_id),
    FOREIGN KEY (field_id) REFERENCES field (field_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS option (
    option_id INTEGER NOT NULL UNIQUE,

    -- parent field
    field_id INTEGER NOT NULL,

    -- key value pair
    option_value INTEGER NOT NULL DEFAULT 0,
    option_name TEXT NOT NULL DEFAULT "value",

    -- an option can not have duplicate option_name's
    UNIQUE(field_id, option_name),

    PRIMARY KEY (option_id),
    FOREIGN KEY (field_id) REFERENCES field (field_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS entry (
    entry_id INTEGER NOT NULL UNIQUE,

    -- parent tracker
    tracker_id INTEGER NOT NULL,

    -- when this was inserted
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- markdown formatted notes
    entry_notes TEXT NOT NULL DEFAULT "",

    PRIMARY KEY (entry_id),
    FOREIGN KEY (tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS log (
    log_id INTEGER NOT NULL UNIQUE,

    -- parent field
    entry_id INTEGER NOT NULL,

    -- field info
    field_id INTEGER NOT NULL,

    -- value
    log_value INTEGER NOT NULL,

    PRIMARY KEY (log_id),
    FOREIGN KEY (entry_id) REFERENCES entry (entry_id) ON DELETE CASCADE,
    FOREIGN KEY (field_id) REFERENCES field (field_id) ON DELETE CASCADE
);

-- --------------------------------------------------------------------------------------
-- Tracker 1
-- --------------------------------------------------------------------------------------

-- Create Tracker
INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("Walk Dog", "Take the dog out");

-- Add Entry
INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog found a turtle");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog learned to fly");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "Dog ran away");

INSERT INTO entry (tracker_id, entry_notes) VALUES (1, "No dog still walk");

-- --------------------------------------------------------------------------------------
-- Tracker 2
-- --------------------------------------------------------------------------------------

-- Create Tracker
INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("Money", "Spending habits");

-- Add Data Fields to Tracker
INSERT INTO field (tracker_id, field_type, field_name, field_notes)
VALUES (2, "number", "Transactions", "Transactions of different types");
INSERT INTO number (field_id, decimal_places)
VALUES (1, 2);

INSERT INTO field (tracker_id, field_type, field_name, field_notes)
VALUES (2, "option", "Card", "What credit/debit card was used?");
INSERT INTO option (field_id, option_value, option_name)
VALUES (2, 1, "Visa"),
       (2, 2, "Discover"),
       (2, 3, "American Express");

-- Log values
INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Gum");
INSERT INTO log (entry_id, field_id, log_value) VALUES (5, 1, 1099);
INSERT INTO log (entry_id, field_id, log_value) VALUES (5, 2, 2);

INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Turtle");
INSERT INTO log (entry_id, field_id, log_value) VALUES (6, 1, 42050);
INSERT INTO log (entry_id, field_id, log_value) VALUES (6, 2, 3);

INSERT INTO entry (tracker_id, entry_notes) VALUES (2, "Coffee");
INSERT INTO log (entry_id, field_id, log_value) VALUES (7, 1, 799);
INSERT INTO log (entry_id, field_id, log_value) VALUES (7, 2, 1);
