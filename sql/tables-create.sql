-- foreign_keys constraints are not on by default
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS tracker (
    tracker_id INTEGER NOT NULL UNIQUE,

    -- name to identify this tracker
    tracker_name TEXT NOT NULL UNIQUE,

    -- markdown formated notes
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
    -- examples: exersise, read status
    field_type TEXT CHECK(field_type in ('number', 'option')) NOT NULL DEFAULT 'number',

    -- name to identify this field
    field_name TEXT NOT NULL,

    -- markdown formated notes
    field_notes TEXT NOT NULL DEFAULT "",

    -- a tracker can not have duplicate field_name's
    UNIQUE(tracker_id, field_name),

    PRIMARY KEY(field_id),
    FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS number (
    number_id INTEGER NOT NULL UNIQUE,

    -- parent field
    field_id INTEGER NOT NULL,

    -- max/min value
    max_flag INTEGER NOT NULL DEFAULT false,
    max_value INTEGER NOT NULL DEFAULT 1000,
    min_flag INTEGER NOT NULL DEFAULT false,
    min_value INTEGER NOT NULL DEFAULT 1,

    -- 0 round to integer
    -- 2 round to 2 decimal places. example money
    -- -3 is thousands
    decimal_places INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY(number_id)
    FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
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

    PRIMARY KEY(option_id)
    FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
);
