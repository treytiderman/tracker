DROP TABLE IF EXISTS tracker_3;
DELETE FROM tracker WHERE tracker_name = "lift";

INSERT INTO tracker (tracker_name, tracker_notes)
VALUES ("lift", "lifting habits");

INSERT INTO field 
    (tracker_id, field_type, field_name)
VALUES 
    (3, "option", "exersise"),
    (3, "number", "weight"),
    (3, "number", "reps");

INSERT INTO option
    (field_id, option_value, option_name)
VALUES 
    (2, 1, "bench"),
    (2, 2, "sqaut"),
    (2, 3, "deadlift");

INSERT INTO number (field_id)
VALUES (3), (4);

SELECT tracker_id FROM tracker WHERE tracker_name = "lift";

CREATE TABLE IF NOT EXISTS tracker_3 (
    id INTEGER NOT NULL UNIQUE,

    -- parent tracker
    tracker_id INTEGER NOT NULL DEFAULT 3,

    -- timestamp when this was set
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- markdown formated notes
    notes TEXT NOT NULL DEFAULT "",

    -- custom fields from the field table
    exersise INTEGER NOT NULL DEFAULT 0,
    weight INTEGER NOT NULL DEFAULT 0,
    reps INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY(id),
    FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id)
);

INSERT INTO tracker_3
    (notes, exersise, weight, reps)
VALUES
    ("", 0, 135, 5),
    ("finally two plates!!", 1, 225, 5),
    ("", 0, 135, 8),
    ("broke back after first rep :(", 1, 225, 1);

SELECT * FROM tracker_3;